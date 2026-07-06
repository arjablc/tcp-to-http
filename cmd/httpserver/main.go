package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/arjablc/tcp-to-http/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port)
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
