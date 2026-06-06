package internals

import (
	"fmt"
	"log"
	"os"
)

// why even create a file reader
// Because file and network behave similarly when it comes to data
func fileReader() {
	file, err := os.Open("./message.txt")
	if err != nil {
		log.Fatalf("Failed reading file error: %v", err)
	}
	lineChan := GetLinesChannel(file)

	for line := range lineChan {
		fmt.Printf("read: %s\n", line)
	}
}

// IO reader and writer are interfaces
// that work with both network and files
