package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/arjablc/tcp-to-http/internal/headers"
	"github.com/arjablc/tcp-to-http/internal/request"
	"github.com/arjablc/tcp-to-http/internal/response"
	"github.com/arjablc/tcp-to-http/internal/server"
)

const port = 42069

func requestHandler(w *response.Writer, req *request.Request) {
	log.Println("Inside handler")
	log.Println("Request Targer", req.RequestLine.RequestTarget)
	if after, ok := strings.CutPrefix(req.RequestLine.RequestTarget, "/httpbin"); ok {
		log.Println("Http bin branch hit")
		forwardTarget := after
		url := fmt.Sprintf("https://httpbin.org/%s", forwardTarget)
		res, err := http.Get(url)
		if err != nil {
			writeResp(w, response.StatusInternal, "Failed to forward request")
		}
		buff := make([]byte, 1024)

		headers := headers.NewHeaders()
		headers.Set("Connection", "close")
		headers.Set("Transfer-Encoding", "chunked")
		err = w.WriteStatusLine(200)
		if err != nil {
			fmt.Println(err)
		}
		err = w.WriteHeaders(headers)
		if err != nil {
			fmt.Println(err)
		}
		readIdx := 0
		for {
			n, err := res.Body.Read(buff)
			fmt.Println("Incoming Chunk Size Reads:", n)
			fmt.Println("-----------Data:-----------")
			fmt.Print(string(buff))
			fmt.Println("\n----------------------")
			if err != nil {
				if errors.Is(err, io.EOF) {
					_, err = w.WriteChunkedBodyDone()
					if err != nil {
						fmt.Println(err)
					}
					return
				}
				fmt.Println(err)
				return
			}
			n, err = w.WriteChunkedBody(buff)
			if err != nil {
				fmt.Println(err)
			}
			readIdx += n
		}

	}

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		writeResp(w, response.StatusBadRequest, badReqHtml)
	case "/myproblem":
		writeResp(w, response.StatusInternal, internalServerHtml)
	default:
		writeResp(w, response.StatusOk, okHtml)
	}

}

func writeResp(w *response.Writer, statusCode response.StatusCode, body string) {

	bodyBytes := []byte(body)
	headers := headers.NewHeaders()
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/html")
	headers.Set("Content-Length", fmt.Sprintf("%d", len(bodyBytes)))

	err := w.WriteStatusLine(statusCode)
	if err != nil {
		fmt.Println(err)
	}
	err = w.WriteHeaders(headers)
	if err != nil {
		fmt.Println(err)
	}
	_, err = w.WriteBody(bodyBytes)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	server, err := server.Serve(port, requestHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)

	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	// waiting or signal from os, for sigint and sigterm
	// when the user closes the server using ctrl+c or
	// when os closes the server
	<-sigChan
	log.Println("Server gracefully stopped")

	// NOTE:
	// why even use the sigChan

	// this will run the Serve method immidiately
	// serve will use goroutines to handle requests
	// if we don't do the sigchan then this main function will return
	// causing the server to stop
}
