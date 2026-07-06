package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	// don't know hwat state to save here
	listener net.Listener
	port     int

	serverOpen *atomic.Bool
}

func (s *Server) Close() error {
	fmt.Println("Close called")
	s.serverOpen.Store(false)
	err := s.listener.Close()
	if err != nil {
		return err
	}
	return nil

}

func (s *Server) listen() {
	for {

		fmt.Println("Listen called")
		if !s.serverOpen.Load() {
			log.Fatalf("Server is closed, Not listening to requests")

		}
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("Some error in connection accept", conn)
		}

		go s.handle(conn)

	}

}

func (s *Server) handle(conn net.Conn) {
	fmt.Println("Handle called")
	/*
		HTTP/1.1 200 OK
		Content-Type: text/plain
		Content-Length: 13

		Hello World!
	*/
	// 	writeString := `
	// HTTP/1.1 200 OK
	// Content-Type: text/plain
	// Content-Length: 13
	//
	// Hello World!
	// 	`
	n, err := conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!"))
	if err != nil {
		fmt.Println("Error on connection write", err)
	}

	fmt.Println("Bytes written", n)
	conn.Close()
}

func Serve(port int) (*Server, error) {
	fmt.Println("Serve called")
	// create a new server with the port
	portStr := fmt.Sprintf(":%d", port)
	fmt.Println(portStr)

	tcpListener, err := net.Listen("tcp", portStr)
	if err != nil {
		return nil, fmt.Errorf("Failed to start tcp, %v", err)
	}
	var serverRunning atomic.Bool
	serverRunning.Store(true)
	server := Server{
		port:       port,
		listener:   tcpListener,
		serverOpen: &serverRunning,
	}
	go server.listen()
	return &server, nil
}
