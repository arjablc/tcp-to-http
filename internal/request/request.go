package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/arjablc/tcp-to-http/internal/headers"
)

type RequestState int

const (
	RequestStateInit RequestState = iota
	RequestStateDone
	RequestStateParsingHeaders
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	ReqState    RequestState
}

func (r *Request) parseSingle(data []byte) (int, error) {
	log.Println("Inside Parse Single")
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
		fmt.Println("ParseSingle: headers parsing")
		fmt.Println("ParseSingle: data len:", len(data))
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.ReqState = RequestStateDone
		}
		return n, nil
	case RequestStateDone:
		return 0, fmt.Errorf("Trying to read data when data stream is done")
	default:
		return 0, fmt.Errorf("Invalid state of Parser")
	}

}

func (r *Request) parse(data []byte) (int, error) {
	fmt.Println("Inside Parse Method")
	totalBytesRead := 0
	for r.ReqState != RequestStateDone {
		fmt.Printf("Current request State: %v\n", r.ReqState)
		fmt.Printf("Current totalBytesRead: %v\n", totalBytesRead)
		n, err := r.parseSingle(data[totalBytesRead:])
		fmt.Printf("Returning bytes from parse single: %d", n)

		if n == 0 {
			return 0, nil
		}
		if err != nil {
			return n, err
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
		fmt.Println("================BUFFER==================")
		fmt.Print(string(readBuffer))
		fmt.Println()
		fmt.Println("================BUFFER==================")
		if err != nil {
			if errors.Is(err, io.EOF) {
				request.ReqState = RequestStateDone
				break
			}
			return nil, err
		}
		// update the number of bytes read or the index of the last read byte
		bytesRead += readIdx
		fmt.Println("Calling parse")
		parsedIdx, err := request.parse(readBuffer[:bytesRead])
		fmt.Println()
		fmt.Println(" parse Call ended")
		if err != nil {
			return nil, err
		}
		// Done to remove the parsed stuff and then prepare a new buffer
		// with unparsed data(might be other parts of a request)
		// if not done then nothing will be moved out of the buffer
		copy(readBuffer, readBuffer[parsedIdx:])
		bytesRead -= parsedIdx
	}
	fmt.Println("Returning request")
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
