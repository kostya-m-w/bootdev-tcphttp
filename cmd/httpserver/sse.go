package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"tcphttp/internal/request"
	"tcphttp/internal/response"
	"tcphttp/internal/server"
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

func sseStream(w *response.Writer, r *request.Request, pipes server.SsePipeStorage) {
	h := response.GetDefaultHeaders(0)
	connectionId, ok := r.QueryParam("id")
	if !ok {
		w.WriteStatusLine(response.StatusBadRequest)
		return
	}
	pipe, ok := pipes[connectionId]
	if !ok {
		pr, pw := io.Pipe()
		pipe = server.SsePipe{Reader: pr, Writer: pw}
		pipes[connectionId] = pipe
	}
	reader := bufio.NewReader(pipe.Reader)

	h.Remove("Content-Length")
	h.HardSet("Content-type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	h.HardSet("Connection", "keep-alive")
	w.WriteStatusLine(response.StatusOk)
	w.WriteHeaders(h)

	for {
		data, err := reader.ReadBytes('\n')

		if err != nil && errors.Is(err, io.EOF) {
			return
		}
		if err != nil {

			fmt.Printf("Error readint sse pipe(%v): %v", connectionId, err)
			return
		}
		data = bytes.TrimSuffix(data, []byte("\n"))
		w.WriteSse(data)
	}
}

func recieveForSse(w *response.Writer, r *request.Request, pipes server.SsePipeStorage) {
	h := response.GetDefaultHeaders(0)
	id, id_ok := r.QueryParam("id")
	message, message_ok := r.QueryParam("message")

	w.WriteHeaders(h)
	if !id_ok && !message_ok {
		w.WriteStatusLine(response.StatusBadRequest)
	}
	pipe, ok := pipes[id]
	if !ok {
		w.WriteStatusLine(response.StatusBadRequest)
	}

	writer := bufio.NewWriter(pipe.Writer)
	_, err := writer.Write([]byte(fmt.Sprintf("%v\n", message)))
	if err != nil {
		w.WriteStatusLine(response.StatusInternalServerError)
	}
	if err != writer.Flush() {
		w.WriteStatusLine(response.StatusInternalServerError)
	}

	w.WriteStatusLine(response.StatusOk)
}
