package response

import(
	"tcphttp/internal/headers"
	"io"
	"fmt"
	"strconv"
)


type StatusCode int
const (
	StatusOk StatusCode = 200
	StatusBadRequest StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error{
	statusLine := []byte(fmt.Sprintf("HTTP/1.1 %v ", statusCode))

	switch statusCode {
	case StatusOk:
		statusLine = append(statusLine, []byte("OK")...)
	case StatusBadRequest:
		statusLine = append(statusLine, []byte("Bad Request")...)
	case StatusInternalServerError:
		statusLine = append(statusLine, []byte("Internal Server Error")...)
	}
	w.Write(statusLine)
	w.Write([]byte("\r\n"))
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers{
	headers := headers.NewHeaders()
	headers.Set("Content-Length", strconv.Itoa(contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")

	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for header, value := range(headers) {
		headerLine := fmt.Sprintf("%v: %v\r\n", header, value)
		fmt.Printf("Writing header: %q\n", headerLine)
		w.Write([]byte(headerLine))
	}
	return nil
}

