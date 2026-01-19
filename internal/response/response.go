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

type ResponseWriteState string
const (
	WriteStatusLineState ResponseWriteState = "Write Status Line"
	WriteHeadersState ResponseWriteState = "Write Headers"
	WriteBodyState ResponseWriteState = "Write Body"
	WriteDoneState ResponseWriteState = "Write Done"
)

type Writer struct {
	conn io.Writer
	WriteState ResponseWriteState
	
}
func NewWriter(w io.Writer) Writer{
	return Writer{
		conn: w,
		WriteState: WriteStatusLineState,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	statusLine := []byte(fmt.Sprintf("HTTP/1.1 %v ", statusCode))

	switch statusCode {
	case StatusOk:
		statusLine = append(statusLine, []byte("OK")...)
	case StatusBadRequest:
		statusLine = append(statusLine, []byte("Bad Request")...)
	case StatusInternalServerError:
		statusLine = append(statusLine, []byte("Internal Server Error")...)
	}
	w.conn.Write(statusLine)
	w.conn.Write([]byte("\r\n"))
	w.WriteState = WriteHeadersState
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers{
	headers := headers.NewHeaders()
	headers.Set("Content-Length", strconv.Itoa(contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")

	return headers
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.WriteState != WriteHeadersState {
		return fmt.Errorf("Wrong writer state - expected: %v, actual: %v", WriteHeadersState, w.WriteState)
	}
	for header, value := range(headers) {
		headerLine := fmt.Sprintf("%v: %v\r\n", header, value)
		fmt.Printf("Writing header: %q\n", headerLine)
		w.conn.Write([]byte(headerLine))
	}
	w.conn.Write([]byte("\r\n"))
	w.WriteState = WriteBodyState
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) { 
	if w.WriteState != WriteBodyState {
		return 0, fmt.Errorf("Wrong writer state - expected: %v, actual: %v", WriteBodyState, w.WriteState)
	}
	w.conn.Write(p)

	w.WriteState = WriteDoneState
	return len(p), nil
}
