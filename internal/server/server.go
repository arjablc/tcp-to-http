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

	closed *atomic.Bool
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		err := s.listener.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error Accepting Connection: %v", err)
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 13\r\n\r\n" +
		"Hello World!\n"

	_, err := conn.Write([]byte(response))

	if err != nil {
		fmt.Println("Error Writing to Connection", err)
	}
}

func Serve(port int) (*Server, error) {
	// create a new server with the port
	portStr := fmt.Sprintf(":%d", port)
	fmt.Println(portStr)

	tcpListener, err := net.Listen("tcp", portStr)
	if err != nil {
		return nil, fmt.Errorf("Failed to start tcp, %v", err)
	}
	var serverRunning atomic.Bool
	server := Server{
		port:     port,
		listener: tcpListener,
		closed:   &serverRunning,
	}
	go server.listen()
	return &server, nil
}
