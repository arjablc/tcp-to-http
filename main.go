package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("./message.txt")
	if err != nil {
		log.Fatalf("Failed reading file error: %v", err)
	}
	lineChan := getLinesChannel(file)
	for {
		line, ok := <-lineChan
		if !ok {
			fmt.Println("Channel cloased")
		}
		fmt.Printf("read: %s\n", line)
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	strChan := make(chan string)
	go func() {
		defer f.Close()
		defer close(strChan)

		var currentLine string
		// suppose we only have a way to read 8 byte
		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)

			if err != nil {
				if currentLine != "" {
					strChan <- currentLine
				}
				if errors.Is(err, io.EOF) {
					break

				}
				log.Fatalf("Failed reading file bytes: %v", err)
				return
			}

			bufStr := string(buffer[:n])
			parts := strings.Split(bufStr, "\n")
			for i := 0; i < len(parts)-1; i++ {
				strChan <- (currentLine + parts[i])
				currentLine = ""
			}

			// if no new line we'll concatinate the whole buffer
			currentLine += parts[len(parts)-1]
		}
		// this is done to declare and call the func
		// at the same place

	}()
	return strChan
}
