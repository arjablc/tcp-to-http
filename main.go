package main

import (
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
		var currentLine string
		// suppose we only have a way to read 8 byte
		for {
			buffer := make([]byte, 8)
			_, err := f.Read(buffer)

			if err == io.EOF {
				strChan <- currentLine
				f.Close()
				close(strChan)
				os.Exit(0)
			}
			if err != nil {
				log.Fatalf("Failed reading file bytes: %v", err)

			}

			bufStr := string(buffer)
			parts := strings.Split(bufStr, "\n")

			if len(parts) > 1 {
				// that would mean we've hit new line
				currentLine += parts[0]
				strChan <- currentLine
				currentLine = parts[1]
				continue
			}
			// if no new line we'll concatinate the whole buffer
			currentLine += bufStr
		}
		// this is done to declare and call the func
		// at the same place

	}()
	return strChan
}
