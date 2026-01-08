package main

import(
	"fmt"
	"os"
	"errors"
	"io"
	"strings"
)

func main(){
	file, err := os.Open("./messages.txt")
	reader := getLinesChannel(file)
	if err != nil {
		fmt.Printf("Error reading file: %v", err)
	}

	for line := range(reader) {
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
			//fmt.Printf("read: \"%q\"\n", currentLine)
		}

	}()
	return out
}
