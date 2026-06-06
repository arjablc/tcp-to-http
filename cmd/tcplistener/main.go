package main

import (
	"fmt"
	"log"
	"net"

	"github.com/arjablc/tcp-to-http/internal/readers"
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
		fmt.Printf("Connection created\n")
		linesChan := internals.GetLinesChannel(connection)
		for line := range linesChan {
			fmt.Println(line)
		}

	}
}
