package main

import (
	"crypto/sha256"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"tcphttp/internal/headers"
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

	} else if target == "/sse-client" {
		sseClient(w, r)
	} else if target == "/sse-stream" {
		sseStream(w, r)
	} else if target == "/video" {
		h := response.GetDefaultHeaders(0)
		h.HardSet("Content-Type", "video/mp4")
		w.WriteStatusLine(response.StatusOk)
		video, err := os.ReadFile("./assets/vim.mp4")
		if err != nil {
			fmt.Printf("Error loading video: %v", err)
			return
		}
		h.HardSet("content-length", strconv.Itoa(len(video)))
		w.WriteHeaders(h)
		fmt.Printf("File size: %v", len(video))
		w.WriteBody(video)
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
		h.Set("Trailer", "X-Content-SHA256")
		h.Set("Trailer", "X-Content-Length")
		w.WriteHeaders(h)
		bytesRead := 0
		fullBody := []byte{}
		for {
			chunk := make([]byte, 32)
			n, err := resp.Body.Read(chunk)
			bytesRead += n
			fmt.Printf("HTTPbin read: n: %v(%X), %q\n", n, n, chunk[:n])
			if err != nil {
				break
			}

			n, _ = w.WriteChunkedBody(chunk[:n])
			fullBody = append(fullBody, chunk[:n]...)
		}

		fmt.Println("End of stream")
		fmt.Printf("bytesRead: %v\n", bytesRead)
		_, _ = w.WriteChunkedBodyDone()
		trailers := headers.NewHeaders()
		out := sha256.Sum256(fullBody)
		trailers.Set("X-Content-SHA256", sha256ToStr(out[:]))
		trailers.Set("X-Content-Length", fmt.Sprintf("%d", len(fullBody)))
		w.WriteTrailers(trailers)

		w.WrapResponse()
	} else {
		defHandler(w, func() (response.StatusCode, string) {
			statusCode = response.StatusOk
			body = bodyHtml("200 OK", "Success!", "Your request was an absolute banger.")
			return statusCode, body
		})
	}

}

func sha256ToStr(hash []byte) string {
	encoded := ""
	for _, b := range hash {
		encoded += fmt.Sprintf("%02x", b)
	}
	return encoded

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
