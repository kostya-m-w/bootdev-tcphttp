package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"tcphttp/internal/request"
	"tcphttp/internal/response"
	"tcphttp/internal/server"

	//"io"
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

type injectedHandler func() (response.StatusCode, string)

func defHandler(w *response.Writer, injected injectedHandler) {
	headers := response.GetDefaultHeaders(0)
	statusCode, body := injected()
	headers.HardSet("content-length", strconv.Itoa(len(body)))
	w.WriteStatusLine(statusCode)
	w.WriteHeaders(headers)
	w.WriteBody([]byte(body))
}

func handler(w *response.Writer, r *request.Request) {
	var statusCode response.StatusCode
	var body string
	h := response.GetDefaultHeaders(0)

	h.HardSet("Content-Type", "text/html")
	fmt.Println("Inside handler func")
	target := r.Target()
	if target == "/yourproblem" {
		defHandler(w, func() (response.StatusCode, string) {
			statusCode = response.StatusBadRequest
			body = bodyHtml("400 Bad Request", "Bad Request", "Your request honestly kinda sucked.")
			return statusCode, body
		})
	} else if target == "/myproblem" {
		defHandler(w, func() (response.StatusCode, string) {
			statusCode = response.StatusInternalServerError
			body = bodyHtml("500 Internal Server Error", "Internal Server Error", "Okay, you know what? This one is on me.")
			return statusCode, body
		})
	} else if strings.HasPrefix(target, "/httpbin") {
		subTarget := strings.TrimPrefix(target, "/httpbin")
		proxyTarget := fmt.Sprintf("https://httpbin.org%v", subTarget)
		resp, err := http.Get(proxyTarget)
		if err != nil {
			w.WriteStatusLine(response.StatusInternalServerError)
			w.WriteHeaders(h)

			defHandler(w, func() (response.StatusCode, string) {
				return response.StatusInternalServerError, bodyHtml("500 Internal Server Error", "Internal Server Error", fmt.Sprintf("%v", err))
			})
			return
		}
		fmt.Println("HTTPbin response headers")
		for k, v := range resp.Header {
			fmt.Printf("%v: %q\n", k, v)
		}
		h.Remove("content-length")
		h.Set("Transfer-Encoding", "chunked")

		w.WriteStatusLine(response.StatusOk)
		w.WriteHeaders(h)
		buff := make([]byte, 1024)
		for {
			n, err := resp.Body.Read(buff)
			fmt.Printf("HTTPbin read: n: %v(%X), %q\n", n, n, buff[:n])
			if err != nil {

				fmt.Println("Error reading httpbin")
				_, _ = w.WriteChunkedBodyDone()
				break
			}

			_, _ = w.WriteChunkedBody(buff[:n])
			if n < 1024 {
				fmt.Println("End of stream")
				_, _ = w.WriteChunkedBodyDone()
				break
			}
		}

	} else {
		defHandler(w, func() (response.StatusCode, string) {
			statusCode = response.StatusOk
			body = bodyHtml("200 OK", "Success!", "Your request was an absolute banger.")
			return statusCode, body
		})
	}

}

func bodyHtml(title, h1, message string) string {
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
