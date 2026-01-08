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
	defer file.Close()
	if err != nil {
		fmt.Printf("error reading file: %v", err)
	}


	var currentLine string
	for {
		buf := make([]byte, 8)
		_, err := file.Read(buf)
		chunk := strings.Trim(string(buf), "\x00")
		parts := strings.Split(chunk, "\n")
		
		for i, part := range(parts) {
		//	fmt.Printf("part: %q\n", strings.Trim(part, "\x00"))
			currentLine += part
			if i < len(parts) - 1 {
				// add part to current line
				// 
				//currentLine += part
				fmt.Printf("read: %s\n", currentLine)
				currentLine = ""
				//currentLine += part
			}
		}
	//	fmt.Println("----")

		if errors.Is(err, io.EOF) {
			break
		}
	}
	if len(currentLine) > 0 {
		fmt.Printf("read: \"%q\"\n", currentLine)
	}
}
