package main

import (
	"fmt"
	"os"
	"strconv"
	"tcphttp/internal/request"
	"tcphttp/internal/response"
)

func sseClient(w *response.Writer, r *request.Request) {
	h := response.GetDefaultHeaders(0)
	h.HardSet("Content-Type", "text/html")
	html, err := os.ReadFile("./assets/sse-client.html")
	if err != nil {
		w.WriteStatusLine(response.StatusInternalServerError)
		w.WriteHeaders(h)
		return
	}

	h.HardSet("content-length", strconv.Itoa(len(html)))
	w.WriteStatusLine(response.StatusOk)
	w.WriteHeaders(h)
	w.WriteBody(html)
}

func sseStream(w *response.Writer, r *request.Request) {
	h := response.GetDefaultHeaders(0)
	h.Remove("Content-Length")
	h.HardSet("Content-type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	h.HardSet("Connection", "keep-alive")
	w.WriteStatusLine(response.StatusOk)
	w.WriteHeaders(h)

	i := 0
	for {
		st := fmt.Sprintf("data: message %v", i)
		w.WriteSse([]byte(st))
		i++
	}
}
