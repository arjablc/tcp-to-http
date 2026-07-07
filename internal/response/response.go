package response

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"

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

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	const chunkSize int64 = 1024
	regNurseBytes := []byte("\r\n")
	chunkSizeHex := strconv.FormatInt(chunkSize, 16)
	// lets assume like 16 bytes chunk
	chunkSizeLine := fmt.Appendf(nil, "%s\r\n", chunkSizeHex)
	var writtenIdx int
	for i := 0; i < len(p); i += int(chunkSize) {
		n, err := w.Buf.Write(chunkSizeLine)
		if err != nil {
			return 0, fmt.Errorf("Error Writing chunk: %v", err)
		}
		writtenIdx += n
		end := min(i+int(chunkSize), len(p))
		data := append(p[i:end], regNurseBytes...)
		n, err = w.Buf.Write(data)
		if err != nil {
			return 0, fmt.Errorf("Error Writing chunk: %v", err)
		}
		writtenIdx += n
	}
	return writtenIdx, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	endBytes := fmt.Appendf(nil, "%s\r\n", strconv.FormatInt(0, 16))
	i, err := w.Buf.Write(endBytes)
	if err != nil {
		return 0, fmt.Errorf("Error Writing chunk Done : %v", err)
	}
	j, err := w.Buf.Write([]byte("\r\n"))
	if err != nil {
		return 0, fmt.Errorf("Error Writing chunk Done : %v", err)
	}
	return i + j, nil
}
