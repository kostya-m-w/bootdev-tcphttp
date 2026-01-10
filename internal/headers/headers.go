package headers

import(
	"bytes"
	"fmt"
)
var SEPARATOR = []byte("\r\n")
var HEADERS_END = []byte("\r\n\r\n")
var NAME_VALUE_SEPARATOR = []byte(":")
//var WS = []byte(" ")
var WS = " "

var HEADER_PARSE_ERROR = fmt.Errorf("Error parsing header")
type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	headersEndIndex := bytes.Index(data, HEADERS_END)
	if headersEndIndex == -1 {
		return 0, false, nil
	}

	headerLines := bytes.Split(data[:headersEndIndex], SEPARATOR)
	
	for _, headerLine := range(headerLines) {
		sepIndex := bytes.Index(headerLine, []byte(":"))

		name := headerLine[:sepIndex]
		val := headerLine[sepIndex+1:]
		if bytes.HasSuffix(name, []byte(WS))  {
			return 0, false, HEADER_PARSE_ERROR
		}

		h[string(bytes.Trim(name, WS))] = string(bytes.Trim(val, WS))
	}
	
	return headersEndIndex + 2, true, nil
}
