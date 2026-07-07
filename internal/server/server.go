package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/arjablc/tcp-to-http/internal/request"
	"github.com/arjablc/tcp-to-http/internal/response"
)

type Server struct {
	// don't know hwat state to save here
	listener net.Listener
	port     int
	Handler  Handler

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
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	resWriter := response.NewWriter()
	request, err := request.RequestFromReader(conn)
	if err != nil {
		resWriter.WriteStatusLine(response.StatusBadRequest)
		body := fmt.Appendf(nil, "Error parsing request: %v", err)
		resWriter.WriteHeaders(response.GetDefaultHeaders(len(body)))
		resWriter.WriteBody(body)
		return
	}

	s.Handler(resWriter, request)
	_, err = conn.Write(resWriter.Buf.Bytes())
	if err != nil {
		fmt.Println("Error writing Body")
	}
}

func Serve(port int, handler Handler) (*Server, error) {
	// create a new server with the port
	portStr := fmt.Sprintf(":%d", port)
	tcpListener, err := net.Listen("tcp", portStr)
	if err != nil {
		return nil, fmt.Errorf("Failed to start tcp, %v", err)
	}
	var serverRunning atomic.Bool
	server := Server{
		port:     port,
		listener: tcpListener,
		closed:   &serverRunning,
		Handler:  handler,
	}
	go server.listen()
	return &server, nil
}

// func writeHandlerError(w io.Writer, handlerErr *HandlerError) {
// 	messageBytes := []byte(handlerErr.Message)
// 	headers := response.GetDefaultHeaders(len(messageBytes))
// 	err := response.WriteStatusLine(w, response.StatusCode(handlerErr.StatusCode))
// 	if err != nil {
// 		fmt.Println("Error writing status line")
// 	}
// 	err = response.WriteHeaders(w, headers)
// 	if err != nil {
// 		fmt.Println("Error writing headers")
// 	}
// 	w.Write(messageBytes)
//
// }
