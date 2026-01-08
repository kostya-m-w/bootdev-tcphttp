package request

import(
	"io"
	"fmt"
	"strings"
)

type Request struct{
	RequestLine RequestLine
}

type RequestLine struct {
	Method string
	RequestTarget string
	HttpVersion string
}

func RequestFromReader(reader io.Reader) (*Request, error){
	var req Request
	rawMessage, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	
	message := strings.Split(string(rawMessage), "\r\n")
	reqLine, err := parseRequestLine(message[0])
	if err  != nil {
		return nil, err
	}
	req.RequestLine = reqLine
	return &req, nil
}

func parseRequestLine(rawLine string) (RequestLine, error) {
	var reqLine RequestLine
	parsed := strings.Split(rawLine, " ")

	if len(parsed) != 3 {
		return reqLine, fmt.Errorf("Error parsing request line")
	}

	split_version := strings.Split(parsed[2], "/")

	reqLine.Method = parsed[0]
	reqLine.RequestTarget = parsed[1]
	reqLine.HttpVersion = split_version[1]
	return reqLine, nil
}
