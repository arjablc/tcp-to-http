package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {

	udpAddr, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		fmt.Printf("Failed to resolve udp address: %v", err)
	}
	connection, err := net.DialUDP("udp", nil, udpAddr)

	if err != nil {
		fmt.Printf("Error establishing UDP connection: %v", err)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error taking input: %v", err)
		}

		connection.Write([]byte(input))

	}

}
