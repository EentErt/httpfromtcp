package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Printf("Error listening on port 42069: %v\n", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		fmt.Println("connection accepted")

		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("Error reading request: %v", err)
			conn.Close()
		}

		req.PrintRequest()
	}
}

/*
func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {
		currentLine := ""

		chars := make([]byte, 8)
		for {
			charCount, err := f.Read(chars)
			if err != nil && err != io.EOF {
				fmt.Printf("Error reading file: %v\n", err)
				close(ch)
			}
			defer f.Close()

			parts := strings.Split(string(chars[:charCount]), "\n")
			if len(parts) == 1 && err != io.EOF {
				currentLine += string(chars[:charCount])
				continue
			} else if len(parts) > 1 {
				currentLine += parts[0]
				ch <- currentLine
				currentLine = ""
				for i := 1; i < len(parts); i++ {
					if i < len(parts)-1 {
						ch <- parts[i]
					} else {
						currentLine = parts[i]
					}
				}
				continue
			} else {
				currentLine += string(chars[:charCount])
				ch <- currentLine
				close(ch)
				break
			}
		}
	}()
	return ch
}
*/
