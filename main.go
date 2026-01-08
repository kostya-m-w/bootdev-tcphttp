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
		var currentLine string

		for {
			buf := make([]byte, 8)
			_, err := f.Read(buf)
			chunk := strings.Trim(string(buf), "\x00")
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

		f.Close()
		close(out)
	}()
	return out
}
