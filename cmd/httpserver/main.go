package main

import(
	"tcphttp/internal/server"
	"log"
	"syscall"
	"os"
	"os/signal"
	"tcphttp/internal/request"
	"io"
	"tcphttp/internal/response"
	"fmt"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w io.Writer, r *request.Request) *server.HandlerError {
/*
If the request target (path) is /yourproblem return a 400 and the message "Your problem is not my problem\n"
If the request target (path) is /myproblem return a 500 and the message "Woopsie, my bad\n"
Otherwise, it should just write the string "All good, frfr\n" to the response body.
	*/
	fmt.Println("Inside handler func")
	if r.Target() == "/yourproblem" {
		err := server.HandlerError {
			StatusCode: response.StatusBadRequest,
			Message: "Your problem is not my problem\n",
		}
		return &err
	} else if r.Target() == "/myproblem"{
		err := server.HandlerError {
			StatusCode: response.StatusInternalServerError,
			Message: "Woopsie, my bad\n",
		}
		return &err
	} else {
		w.Write([]byte("All good, frfr\n"))
	}

	return nil
}
