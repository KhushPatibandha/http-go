package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {

	directory := flag.String("directory", "", "abs file path");
	flag.Parse();

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	
	for {
		connection, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		
		go handleConn(connection, *directory);
	}
}

func handleConn(connection net.Conn, directory string) {

	defer connection.Close();
	req, err := http.ReadRequest(bufio.NewReader(connection));
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		return
	}

	fmt.Println("Request method: ", req.Method);
	fmt.Println("Request url: ", req.URL.Path);
	fmt.Println("Request header: ", req.Header.Values("User-Agent"));
	fmt.Println(directory);

	if req.URL.Path == "/" {
		connection.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"));
	} else if strings.Contains(req.URL.Path, "/echo") {
		content := req.URL.Path[6:];
		contentLen := len(content);

		strToReturn := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(contentLen) + "\r\n\r\n" + content;

		connection.Write([]byte(strToReturn));
	} else if strings.Contains(req.URL.Path, "/user-agent") {
		headerContent := req.Header.Values("User-Agent")[0];
		headerContentLen := len(headerContent);

		strToReturn := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(headerContentLen) + "\r\n\r\n" + headerContent;

		connection.Write([]byte(strToReturn));
	} else if strings.Contains(req.URL.Path, "/files") {
		fileName := req.URL.Path[7:];

		fileContent, err := os.ReadFile("/" + directory + "/" + fileName);
		if err != nil {
			connection.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"));
			return;
		}
		fileContentLen := len(fileContent);

		strToReturn := "HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: " + strconv.Itoa(fileContentLen) + "\r\n\r\n" + string(fileContent);

		connection.Write([]byte(strToReturn));
	} else {
		connection.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"));
	}
	
}
