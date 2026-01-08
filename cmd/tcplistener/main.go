package main

import(
	"fmt"
	//"os"
	"errors"
	"io"
	"strings"
	"net"
)

func main(){
	ln, err := net.Listen("tcp", ":42069")

	if err != nil {
		fmt.Printf("Error start listening: %v", err)
	}

	fmt.Printf("Listening to : %v\n", 42069)
	conn, err := ln.Accept()
	if err != nil {
		fmt.Printf("Error accepting connection: %v", err)
	}

	fmt.Printf("Connection Established\n")

	readChan := getLinesChannel(conn)
	for line := range(readChan) {
		fmt.Printf("read: %s\n", line)
	}
}
func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string)

	go func(){
		defer f.Close()
		defer close(out)
		var currentLine string

		for {
			buf := make([]byte, 8)
			n, err := f.Read(buf)
			//chunk := strings.Trim(string(buf[:n]), "\x00")
			chunk := string(buf[:n])
			parts := strings.Split(chunk, "\n")
			
			for i, part := range(parts) {
				currentLine += part
				if i < len(parts) - 1 {
					//fmt.Printf("read: %s\n", currentLine)
					out <- currentLine
					currentLine = ""
				}
			}

			if errors.Is(err, io.EOF) {
				break
			}
		}
		if len(currentLine) > 0 {
			out <- currentLine
		}

	}()
	return out
}
