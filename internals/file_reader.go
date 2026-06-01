package internals

import (
	"fmt"
	"log"
	"os"
)

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
