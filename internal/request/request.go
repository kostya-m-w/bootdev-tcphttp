package request

import(
	"io"
	"fmt"
	"strings"
)

type Request struct{
	RequestLine RequestLine
	state int
}

type RequestLine struct {
	Method string
	RequestTarget string
	HttpVersion string
	state int
}

func RequestFromReader(reader io.Reader) (*Request, error){
	bytesReadCount := 0
	bytesParsedCount := 0

	req := Request{ state: 0 }
	readBuf := make([]byte, 8)
	accBuf := make([]byte, 0) 
	for {
		n, err := reader.Read(readBuf)
		fmt.Printf("Read: %q, n: %v\n", readBuf[:n], n)
		bytesReadCount += n
		if err != nil && err == io.EOF {
			break
		}
		if err != nil {
			return &req, err 
		}

		accBuf = append(accBuf, readBuf[:n]...)
		n, err = req.parse(accBuf)
		bytesParsedCount += n

		if err != nil {
			return &req, err 
		}

		if req.isDone() {
			break
		}
	}
	return &req, nil
}


func (r *Request) parse(data []byte) (int, error){
	reqLine, n, err := parseRequestLine(data)

	if err != nil {
		return 0, err
	}

	if n > 0 {
		r.RequestLine = reqLine
		r.state = 1
	}

	return n, nil
}
func (r *Request) isDone() bool{
	return r.state == 1
}

func parseRequestLine(data []byte) (RequestLine, int, error) {
	var reqLine RequestLine
	bytesReadCount := 0
	lineEnded := false
	rawLine := ""
	for i, b := range(data) {
		bytesReadCount += 1
		rawLine += string(b)
		if i > 0 && data[i-1] == '\r' && data[i] == '\n' {
			lineEnded = true
			break
		}
	}

	if !lineEnded {
		return reqLine, 0, nil
	}

	parsed := strings.Trim(rawLine, "\r\n")
	parts := strings.Split(parsed, " ")

	if len(parts) != 3 {
		return reqLine, bytesReadCount, fmt.Errorf("Error parsing request line")
	}

	split_version := strings.Split(parts[2], "/")

	reqLine.Method = parts[0]
	reqLine.RequestTarget = parts[1]
	reqLine.HttpVersion = split_version[1]
	return reqLine, bytesReadCount, nil
}
