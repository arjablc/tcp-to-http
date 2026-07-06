package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/arjablc/tcp-to-http/internal/headers"
)

type RequestState int

const contentLength = "Content-Length"

const (
	RequestStateInit RequestState = iota
	RequestStateParsingHeaders
	RequestStateBodyParsing
	RequestStateDone
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte

	ReqState       RequestState
	bodyLengthRead int
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.ReqState {
	case RequestStateInit:
		reqLine, n, err := parseReqLine(data)
		if err != nil {
			return 0, err
		}
		if reqLine != nil {
			r.RequestLine = *reqLine
			r.ReqState = RequestStateParsingHeaders
		}
		return n, nil
	case RequestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.ReqState = RequestStateBodyParsing
		}
		return n, nil

	case RequestStateBodyParsing:
		contentLength := r.Headers.Get(contentLength)
		if len(contentLength) == 0 {
			r.ReqState = RequestStateDone
			return len(data), nil
		}
		contentLengthI, err := strconv.Atoi(contentLength)
		if err != nil {
			return 0, err
		}
		r.Body = append(r.Body, data...)
		r.bodyLengthRead += len(data)

		//NOTE: Handling only the more than because
		// when less than the contentlength we might be
		// waiting on more data to come from the stream

		// FIX: wrong because the state of the bodyread is not saved
		// this assumes body parsing happens on sinlgle parse call

		// if len(r.Body) > contentLengthI {
		// 	return 0, errors.New("Content Lenght Mismatch")
		// }

		if r.bodyLengthRead > contentLengthI {
			return 0, errors.New("Content Length Mismatch")
		}

		if len(r.Body) == contentLengthI {
			r.ReqState = RequestStateDone
		}
		return len(data), nil

	case RequestStateDone:
		return 0, fmt.Errorf("Trying to read data when data stream is done")
	default:
		return 0, fmt.Errorf("Invalid state of Parser")
	}

}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesRead := 0
	for r.ReqState != RequestStateDone {
		n, err := r.parseSingle(data[totalBytesRead:])
		if err != nil {
			return n, err
		}
		if n == 0 {
			return totalBytesRead, nil
		}

		totalBytesRead += n
	}

	return totalBytesRead, nil
}

var crlf string = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	bytesRead := 0
	request := Request{
		Headers:  headers.NewHeaders(),
		ReqState: RequestStateInit,
	}
	readBuffer := make([]byte, 8)

	for request.ReqState != RequestStateDone {
		if bytesRead >= len(readBuffer) {
			newBuf := make([]byte, len(readBuffer)*2)
			copy(newBuf, readBuffer)
			readBuffer = newBuf
		}
		readIdx, err := reader.Read(readBuffer[bytesRead:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if request.ReqState == RequestStateBodyParsing {
					return nil, fmt.Errorf("Stream EOF but Parser not done")
				}

				return &request, nil
			}
			return nil, err
		}
		// update the number of bytes read or the index of the last read byte
		bytesRead += readIdx
		parsedIdx, err := request.parse(readBuffer[:bytesRead])
		if err != nil {
			return nil, err
		}
		// Done to remove the parsed stuff and then prepare a new buffer
		// with unparsed data(might be other parts of a request)
		// if not done then nothing will be moved out of the buffer
		copy(readBuffer, readBuffer[parsedIdx:])
		bytesRead -= parsedIdx
	}
	return &request, nil
}

func parseReqLine(rawBytes []byte) (*RequestLine, int, error) {
	crlfIdx := bytes.Index(rawBytes, []byte(crlf))
	if crlfIdx < 0 {
		return nil, 0, nil
	}
	reqLineString := string(rawBytes[:crlfIdx])
	reqLine, err := reqLineFromString(reqLineString)
	if err != nil {
		return nil, 0, err
	}
	return reqLine, crlfIdx + 2, nil
}

func reqLineFromString(in string) (*RequestLine, error) {
	splits := strings.Split(in, " ")
	if len(splits) != 3 {
		return nil, errors.New("Invalid line")
	}
	method, requestTarget, httpVersionString := splits[0], splits[1], splits[2]
	for _, mChar := range method {

		if mChar < 'A' || mChar > 'Z' {
			return nil, errors.New("Invalid method verb")
		}
	}

	httpVersionSplit := strings.Split(httpVersionString, "/")
	if len(httpVersionSplit) != 2 {
		return nil, errors.New("Invalid http version string")
	}
	if httpVersionSplit[1] != "1.1" {
		return nil, errors.New("Invalid HTTP version")
	}

	var reqLine RequestLine
	reqLine.Method = method
	reqLine.RequestTarget = requestTarget
	reqLine.HttpVersion = httpVersionSplit[1]

	return &reqLine, nil

}
