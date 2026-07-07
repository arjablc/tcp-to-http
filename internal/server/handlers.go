package server

import (
	"github.com/arjablc/tcp-to-http/internal/request"
	"github.com/arjablc/tcp-to-http/internal/response"
)

type HandlerError struct {
	StatusCode int
	Message    string
}

func NewHandlerError(statusCode int, message string) *HandlerError {
	return &HandlerError{
		StatusCode: statusCode,
		Message:    message,
	}
}

type Handler func(w *response.Writer, req *request.Request)
