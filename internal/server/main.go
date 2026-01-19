package server

import(
	"net"
	"sync/atomic"
	"fmt"
	"tcphttp/internal/response"
//	"tcphttp/internal/headers"
	"tcphttp/internal/request"
//	"bytes"
//	"io"
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

type Handler func(w *response.Writer, req *request.Request)

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
	defer conn.Close()
	fmt.Println("Start Handling")
	request, err := request.RequestFromReader(conn)

	if err != nil {
		fmt.Printf("Request Header failde to parse parsed: %q", err)
		return 
	}

	fmt.Println("RequestFromReader done")
	responseWriter := response.NewWriter(conn)

	fmt.Println("going to call handler func")
	s.handler(&responseWriter, request)
}
