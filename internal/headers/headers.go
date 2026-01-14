package headers

import(
	"bytes"
	"fmt"
)
var SEPARATOR = []byte("\r\n")
var NAME_VALUE_SEPARATOR = []byte(":")
var WS = " "
var UPPER_LETTERS = []byte("ABCDEFGHIJKLMNOPQRSTUVXYZ")
var LOWER_LETTERS = []byte("abcdefghijklmnopqrstuvxyz")
var NUMBERS = []byte("0123456789")
var ALLOWED_CHARS = []byte("!#$%&'*+-.^_`|~|")

var HEADER_PARSE_ERROR = fmt.Errorf("Error parsing header")
var INVALID_CHARACTER = fmt.Errorf("Error parsing header: not allower characters")
var HEADER_NAME_TO_SHORT = fmt.Errorf("Error parsing header: Header name to short")
type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	headersSepIndex := bytes.Index(data, SEPARATOR)
	if headersSepIndex == -1 {
		return 0, false, nil
	}else if headersSepIndex == 0{
		return 2, true, nil
	}

	headerLine := data[:headersSepIndex]
	
	sepIndex := bytes.Index(headerLine, NAME_VALUE_SEPARATOR)

	name := headerLine[:sepIndex]
	val := headerLine[sepIndex+1:]

	if bytes.HasSuffix(name, []byte(WS))  {
		return 0, false, HEADER_PARSE_ERROR
	}

	name = bytes.Trim(name, WS)

	if hasInvalidCharacters(name) {
		return 0, false, INVALID_CHARACTER
	}
	
	name = bytes.ToLower(name)
	
	if len(name) < 1 {
		return 0, false, HEADER_NAME_TO_SHORT
	}

	h[string(name)] = string(bytes.Trim(val, WS))
	
	return headersSepIndex + 2, false, nil
}

func hasInvalidCharacters(val []byte) bool{
	for i := 0; i < len(val); i++ {

		if bytes.IndexByte(UPPER_LETTERS, val[i]) < 0 && bytes.IndexByte(LOWER_LETTERS, val[i]) < 0 && bytes.IndexByte(NUMBERS, val[i]) < 0 && bytes.IndexByte(ALLOWED_CHARS, val[i]) < 0 {
			fmt.Printf("Invalid character: %q\n", val[i])
			return true
		}
	}
	return false
}
