package server

import(
	"net"
	"sync/atomic"
	"fmt"
	"tcphttp/internal/response"
	"tcphttp/internal/request"
	"bytes"
	"io"
)


type Server struct {
	listener net.Listener
	open *atomic.Bool
	handler Handler
}
type HandlerError struct {
	StatusCode response.StatusCode
	Message string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func Serve(port int, handler Handler) (*Server, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%v", port))

	fmt.Println("Listent to tcp port")
	if err != nil {
		return nil, fmt.Errorf("Error start listening: %v", err)
	}

	open := atomic.Bool{}
	open.Store(true)
	server := Server{
		listener: ln,
		open: &open,
		handler: handler,
	}
	
	server.listen()

	return &server, nil
}

func (s *Server) Close() error {
	s.open.Store(false)
	s.listener.Close()
	return nil
}

func (s *Server) listen() {
	fmt.Println("Start Accepting http connections")
	for {
		if s.open.Load() {
			conn, err := s.listener.Accept()

			if err == nil {
				fmt.Println("Connection Accepted")
				go s.handle(conn)
			}else{
				fmt.Printf("Error accepting connection: %v", err)
			}
		}else{
			break
		}
	}
}

func (s *Server) handle(conn net.Conn) {
	fmt.Println("Start Handling")
	request, err := request.RequestFromReader(conn)
	fmt.Println("RequestFromReader done")
	if err != nil {
		fmt.Printf("Request Header failde to parse parsed: %q", err)
		errStr := fmt.Sprintf("Error parsing request: %v\n", err)
		response.WriteStatusLine(conn, 400)
		headers := response.GetDefaultHeaders(len(errStr))
		response.WriteHeaders(conn, headers)

		conn.Write([]byte("\r\n"))
		conn.Write([]byte(errStr))
		conn.Close()
		return
	}

	fmt.Println("Request Header parsed")

	buffBack := make([]byte, 0)
	buff := bytes.NewBuffer(buffBack)

	fmt.Println("going to call handler func")
	handleErr := s.handler(buff, request)

	if handleErr != nil {
		errStr := handleErr.Message
		response.WriteStatusLine(conn, handleErr.StatusCode)
		headers := response.GetDefaultHeaders(len(errStr))
		response.WriteHeaders(conn, headers)

		conn.Write([]byte("\r\n"))
		conn.Write([]byte(errStr))
		conn.Close()
		return
	}



	handleResp := buff.Bytes()
	fmt.Printf("Now going to write to connection: %q", handleResp)
	response.WriteStatusLine(conn, 200)
	fmt.Println("Status line written")
	headers := response.GetDefaultHeaders(len(handleResp))
	response.WriteHeaders(conn, headers)
	fmt.Println("headers written")
	conn.Write([]byte("\r\n"))
	conn.Write(handleResp)
	fmt.Println("Now going to close connection")
	conn.Close()
}
