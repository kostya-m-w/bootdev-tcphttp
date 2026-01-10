package headers

import(
	"bytes"
	"fmt"
)
var SEPARATOR = []byte("\r\n")
var NAME_VALUE_SEPARATOR = []byte(":")
var WS = " "

var HEADER_PARSE_ERROR = fmt.Errorf("Error parsing header")
type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	headersSepIndex := bytes.Index(data, SEPARATOR)
	if headersSepIndex == -1 {
		return 0, false, nil
	}else if headersSepIndex == 0{
		return 0, true, nil
	}

	headerLine := data[:headersSepIndex]
	
	sepIndex := bytes.Index(headerLine, NAME_VALUE_SEPARATOR)

	name := headerLine[:sepIndex]
	val := headerLine[sepIndex+1:]

	if bytes.HasSuffix(name, []byte(WS))  {
		return 0, false, HEADER_PARSE_ERROR
	}

	h[string(bytes.Trim(name, WS))] = string(bytes.Trim(val, WS))
	
	return headersSepIndex + 2, false, nil
}
