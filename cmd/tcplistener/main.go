package main

import (
	"fmt"
	"log"
	"net"

	"github.com/arjablc/tcp-to-http/internal/request"
)

func main() {
	tcpListener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("Failed to establish listener: %v", err)

	}
	defer tcpListener.Close()

	for {
		connection, err := tcpListener.Accept()
		if err != nil {
			fmt.Printf("Failed to wait for connection: %v\n", connection)
		}
		fmt.Printf("Connected\n")

		req, err := request.RequestFromReader(connection)
		if err != nil {
			fmt.Printf("Error, %v", err)

		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for key, value := range req.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}
		fmt.Println("Body:")
		fmt.Println(string(req.Body))

	}
}
