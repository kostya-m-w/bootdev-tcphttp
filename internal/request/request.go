package request

import(
	"io"
	"fmt"
	"strings"
	"bytes"
)

var SEPARATOR = []byte("\r\n")

type parserState int
const (
	StateInit parserState = 0
	StateDone parserState = 1
)

type RequestLine struct {
	Method string
	RequestTarget string
	HttpVersion string
}

type Request struct{
	RequestLine RequestLine
	state parserState
}

func newRequest() *Request {
	return &Request{
		state: StateInit,
	}
}

const buffSize = 8
func RequestFromReader(reader io.Reader) (*Request, error){
	bytesReadCount := 0
	bytesParsedCount := 0
	readBuf := make([]byte, buffSize, buffSize)

	fmt.Println(">>>>")
	req := newRequest()
	for !req.isDone() {
		n, err := reader.Read(readBuf[bytesReadCount:])
		fmt.Printf("Read: %q, n: %v\n", readBuf, n)
		if err != nil && err == io.EOF {
			break
		}
		if err != nil {
			return nil, err 
		}
		bytesReadCount += n

		if bytesReadCount >= len(readBuf) / 2 {

			fmt.Println("extending buffer")
			newSize := bytesReadCount*2
			newBuf := make([]byte, newSize, newSize)
			copy(newBuf, readBuf)
			readBuf = newBuf
		}

		fmt.Printf("Parse: %q, bytesReadCount: %v\n", readBuf[:bytesReadCount], bytesReadCount)
		bytesParsedCount, err = req.parse(readBuf[:bytesReadCount])


		if err != nil {
			return req, err 
		}

		if bytesParsedCount > 0 {
			fmt.Errorf("Removing parsed part")
			copy(readBuf, readBuf[bytesParsedCount:])
			bytesReadCount -= bytesParsedCount
		}
		fmt.Println("----")
	}
	fmt.Println("<<<<")
	return req, nil
}


func (r *Request) parse(data []byte) (int, error){
	fmt.Printf("received to parse: %q\n", data) 
	switch r.state {
	case StateInit:
		reqLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		if n > 0 {
			r.RequestLine = *reqLine
			r.state = StateDone
		}
		return n, nil
	case StateDone:
		return 0, nil
	default:
		return 0, nil
	}
}
func (r *Request) isDone() bool{
	return r.state == StateDone
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	var reqLine RequestLine
	
	sepIndex := bytes.Index(data, SEPARATOR)

	if sepIndex == -1 {
		return &reqLine, 0, nil
	}


	rawLine := data[:sepIndex]
	//leftOver := data[sepindex + len(SEPARATOR):]
	bytesReadCount := len(rawLine)

	parts := strings.Split(string(rawLine), " ")

	if len(parts) != 3 {
		return &reqLine, bytesReadCount, fmt.Errorf("Error parsing request line")
	}

	split_version := strings.Split(parts[2], "/")

	reqLine.Method = parts[0]
	reqLine.RequestTarget = parts[1]
	reqLine.HttpVersion = split_version[1]

	return &reqLine, bytesReadCount, nil
}
