package response

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/arjablc/tcp-to-http/internal/headers"
)

type StatusCode int
type WriterStatus int

var INVALID_WRITE_ORDER error = errors.New("Response Write Called with wrong order")

const (
	StatusOk         StatusCode = 200
	StatusBadRequest StatusCode = 400
	StatusInternal   StatusCode = 500
)

const (
	WriterStatusInit WriterStatus = iota
	WriterStatusLine
	WriterStatusHeaders
	WriterStatusBody
)

type Writer struct {
	Buf bytes.Buffer

	writerStatus WriterStatus
}

func NewWriter() *Writer {
	return &Writer{
		Buf:          *bytes.NewBuffer([]byte{}),
		writerStatus: WriterStatusInit,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerStatus != WriterStatusInit {
		return INVALID_WRITE_ORDER

	}
	statusLine := fmt.Sprintf("HTTP/1.1 %d\r\n", statusCode)
	_, err := w.Buf.Write([]byte(statusLine))
	w.writerStatus = WriterStatusLine
	if err != nil {
		return fmt.Errorf("Write Error: %v", err)
	}
	return nil
}
func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerStatus != WriterStatusLine {
		return INVALID_WRITE_ORDER
	}
	for key, value := range headers {
		headerString := fmt.Sprintf("%s: %s\r\n", key, value)
		_, err := w.Buf.Write([]byte(headerString))
		if err != nil {
			return fmt.Errorf("Header Write Error: %v", err)
		}
	}
	w.Buf.Write([]byte("\r\n"))
	w.writerStatus = WriterStatusHeaders
	return nil

}
func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerStatus != WriterStatusHeaders {
		return 0, INVALID_WRITE_ORDER
	}
	n, err := w.Buf.Write(p)
	w.writerStatus = WriterStatusBody
	return n, err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")
	headers.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	return headers
}
