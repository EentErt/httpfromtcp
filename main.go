package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	currentLine := ""

	chars := make([]byte, 8)
	for {
		charCount, err := file.Read(chars)
		if err != nil && err != io.EOF {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}

		parts := strings.Split(string(chars[:charCount]), "\n")
		if len(parts) == 1 && err != io.EOF {
			currentLine += string(chars[:charCount])
			continue
		} else if len(parts) > 1 {
			currentLine += parts[0]
			fmt.Printf("read: %s\n", currentLine)
			for i := 1; i < len(parts); i++ {
				if i < len(parts)-1 {
					fmt.Printf("read: %s\n", parts[i])
				} else {
					currentLine = fmt.Sprint(parts[i])
				}
			}
		} else {
			currentLine += string(chars[:charCount])
			fmt.Printf("read: %s\n", currentLine)
			break
		}
	}
}
