package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
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
	fmt.Println("Request header: ", req.Header.Values("Content-Type"));
	fmt.Println("Request header: ", req.Header.Values("Accept-Encoding"));
	fmt.Println(directory);

	if req.URL.Path == "/" {
		connection.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"));
	} else if strings.Contains(req.URL.Path, "/echo") {
		headerContentEncoding := req.Header.Values("Accept-Encoding");
		content := req.URL.Path[6:];
		contentLen := len(content);

		var strToReturn string;

		if len(headerContentEncoding) > 0 && strings.Contains(headerContentEncoding[0], "gzip") {
			encodedContent := encodeGzip(content);
			encodedContentLen := len(encodedContent);
			strToReturn = "HTTP/1.1 200 OK\r\nContent-Encoding: gzip\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(encodedContentLen) + "\r\n\r\n" + string(encodedContent);
		} else {
			strToReturn = "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(contentLen) + "\r\n\r\n" + content;
		}

		connection.Write([]byte(strToReturn));
	} else if strings.Contains(req.URL.Path, "/user-agent") {
		headerUserAgentContent := req.Header.Values("User-Agent")[0];
		headerContentEncoding := req.Header.Values("Accept-Encoding");
		headerUserAgentContentLen := len(headerUserAgentContent);

		var strToReturn string;
		if len(headerContentEncoding) > 0 && strings.Contains(headerContentEncoding[0], "gzip") {
			encodedContent := encodeGzip(headerUserAgentContent);
			encodedContentLen := len(encodedContent);
			strToReturn = "HTTP/1.1 200 OK\r\nContent-Encoding: gzip\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(encodedContentLen) + "\r\n\r\n" + string(encodedContent);
		} else {
			strToReturn = "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(headerUserAgentContentLen) + "\r\n\r\n" + headerUserAgentContent;
		}

		connection.Write([]byte(strToReturn));
	} else if req.Method == "GET" && strings.Contains(req.URL.Path, "/files") {
		fileName := req.URL.Path[7:];

		fileContent, err := os.ReadFile("/" + directory + "/" + fileName);
		if err != nil {
			connection.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"));
			return;
		}
		fileContentLen := len(fileContent);

		strToReturn := "HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: " + strconv.Itoa(fileContentLen) + "\r\n\r\n" + string(fileContent);

		connection.Write([]byte(strToReturn));
	} else if req.Method == "POST" && strings.Contains(req.URL.Path, "/files") {
	
		fileName := req.URL.Path[7:];
		fileContent, _ := io.ReadAll(req.Body);

		_ = os.WriteFile("/" + directory + "/" + fileName, fileContent, 0644);

		connection.Write([]byte("HTTP/1.1 201 Created\r\n\r\n"));
	} else {
		connection.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"));
	}
	
}

func encodeGzip(content string) []byte {
	encodedContent := new(bytes.Buffer)
	gz := gzip.NewWriter(encodedContent)
	_, err := gz.Write([]byte(content))
	if err != nil {
		fmt.Println("Error encoding content: ", err.Error())
		return nil;
	}
	err = gz.Close()
	if err != nil {
		fmt.Println("Error closing gzip writer: ", err.Error())
		return nil;
	}
	return encodedContent.Bytes();
}
