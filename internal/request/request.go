package request

import(
	"io"
	"strings"
	"bytes"
	"tcphttp/internal/headers"
	"fmt"
	"strconv"
)

var SEPARATOR = []byte("\r\n")

type parserState int
const (
	StateInit parserState = 0
	StateParsingHeaders parserState = 1
	StateParsingBody parserState = 2
	StateDone parserState = 3
)

type RequestLine struct {
	Method string
	RequestTarget string
	HttpVersion string
}

type Request struct{
	RequestLine RequestLine
	Headers headers.Headers
	Body []byte
	state parserState
}

func newRequest() *Request {
	return &Request{
		state: StateInit,
		Headers: headers.NewHeaders(),
		Body: []byte{},
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
		fmt.Printf("Read loop iteration start bytesReadCount: %v, readBuf: %q", bytesReadCount, readBuf)
		n, err := reader.Read(readBuf[bytesReadCount:])
		eof := false
		fmt.Printf("\nRead: %q, n: %v, eof: %v\n", readBuf, n, eof)
		if err != nil && err == io.EOF {
			eof = true
			err = nil
		}
		if err != nil {
			return nil, err 
		}
		bytesReadCount += n

		if bytesReadCount >= len(readBuf) / 2 {
		//	fmt.Println("extending buffer")
			newSize := bytesReadCount*2
			newBuf := make([]byte, newSize, newSize)
			copy(newBuf, readBuf)
			readBuf = newBuf
		}

		fmt.Printf("Parse: %q, bytesReadCount: %v\n", readBuf[:bytesReadCount], bytesReadCount)
		bytesParsedCount, err = req.parse(readBuf[:bytesReadCount])
		//fmt.Printf("Parse: bytesParsedCount: %v\n", bytesParsedCount)


		if err != nil {
			return req, err 
		}

		if bytesParsedCount > 0 {
			copy(readBuf, readBuf[bytesParsedCount:])
			bytesReadCount -= bytesParsedCount
		}
		if eof {
			if req.state == StateParsingBody && !req.BodyDone() {
				return nil, fmt.Errorf("Partial body")
			}else{
				break
			}
		}
		fmt.Println("----")
	}
	fmt.Println("<<<<")
	return req, nil
}


func (r *Request) parse(data []byte) (int, error){
	//fmt.Printf("received to parse: %q\n", data) 
	switch r.state {
	case StateInit:
		reqLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		if n > 0 {
			r.RequestLine = *reqLine
			r.state = StateParsingHeaders
			//fmt.Printf("Finisehd parsing request line: %q\n", data)
		}
		return n, nil
	case StateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		fmt.Printf("Parse header: %q, done: %v, err: %v\n", data, done, err)
		if done {
			r.state = StateParsingBody
			_ , ok := r.Headers.Get("content-length")
			if !ok {
				r.state = StateDone
			}
		}
		return n, err
	case StateParsingBody:
		fmt.Printf("body piece: %q, accumulated body: %q\n", data, r.Body)
		contentLengthVal, ok := r.Headers.Get("content-length")

		if  !ok && len(data) == 0 {
			return 0, nil
		}else {
			contentLengthVal = "0"
		}

		contentLength, err := strconv.Atoi(contentLengthVal)
		if err != nil {
			return 0, fmt.Errorf("Error converting content-length to int")
		}

		r.Body = append(r.Body, data...)

		//fmt.Printf("body piece: %q, accumulated body: %q\n", data, r.Body)
		currentBodyLen := len(r.Body)
		if  contentLength > 0 && currentBodyLen > contentLength{
			return len(data), fmt.Errorf("body to long, content-length: %v, body length: %v", contentLength, currentBodyLen)
		} else if currentBodyLen == contentLength {
			r.state = StateDone
		}
		return len(data), nil
	case StateDone:
		return 0, nil
	default:
		return 0, nil
	}
}
func (r *Request) isDone() bool{
	return r.state == StateDone
}

func (r *Request) BodyDone() bool {
	declaredLength := r.Headers.ContentLength()
	actualLength := len(r.Body)
	//fmt.Printf("declared length: %v, actual length: %v\n", declaredLength, actualLength)

	return actualLength >= declaredLength
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

	return &reqLine, bytesReadCount + 2, nil
}

func (r *Request) Target() string {
	return r.RequestLine.RequestTarget
}
