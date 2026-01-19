package main

import(
	"tcphttp/internal/server"
	"log"
	"syscall"
	"os"
	"os/signal"
	"tcphttp/internal/request"
	//"io"
	"tcphttp/internal/response"
	"fmt"
	"strconv"
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

func handler(w *response.Writer, r *request.Request) {
/*
If the request target (path) is /yourproblem return a 400 and the message "Your problem is not my problem\n"
If the request target (path) is /myproblem return a 500 and the message "Woopsie, my bad\n"
Otherwise, it should just write the string "All good, frfr\n" to the response body.
	*/

	var statusCode response.StatusCode 
	var body string
	headers := response.GetDefaultHeaders(0)

	headers.HardSet("Content-Type", "text/html")
	fmt.Println("Inside handler func")
	if r.Target() == "/yourproblem" {
		statusCode = response.StatusBadRequest
		body = bodyHtml("400 Bad Request", "Bad Request", "Your request honestly kinda sucked.")
	} else if r.Target() == "/myproblem"{
		statusCode = response.StatusInternalServerError
		body = bodyHtml("500 Internal Server Error", "Internal Server Error", "Okay, you know what? This one is on me.")
	} else {
		statusCode = response.StatusOk
		body = bodyHtml("200 OK", "Success!", "Your request was an absolute banger.")
	}

	headers.HardSet("content-length", strconv.Itoa(len(body)))
	w.WriteStatusLine(statusCode)
	w.WriteHeaders(headers)
	w.WriteBody([]byte(body))
}

func bodyHtml(title, h1, message string) string{
	return fmt.Sprintf(`
<html>
  <head>
    <title>%v</title>
  </head>
  <body>
    <h1>%v</h1>
    <p>%v</p>
  </body>
</html>
`, title, h1, message)
}
