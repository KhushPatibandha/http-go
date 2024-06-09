package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	
	connection, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	handleConn(connection);
}

func handleConn(connection net.Conn) {

	defer connection.Close();
	req, err := http.ReadRequest(bufio.NewReader(connection));
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		return
	}

	fmt.Println("Request method: ", req.Method);
	fmt.Println("Request url: ", req.URL.Path);

	if req.URL.Path == "/" {
		connection.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"));
	}
	
	connection.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"));
}
