package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	requestData, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.New("Failed reading data")
	}
	reqString := string(requestData)
	lines := strings.Split(reqString, "\r\n")
	reqLine, err := parseReqLine(lines[0])
	if err != nil {
		return nil, err
	}
	var req Request
	req.RequestLine = *reqLine

	return &req, nil
}

func parseReqLine(line string) (*RequestLine, error) {
	var reqLine RequestLine
	splits := strings.Split(line, " ")
	if len(splits) != 3 {
		return nil, errors.New("Invalid line")
	}
	method, requestTarget, httpVersionString := splits[0], splits[1], splits[2]
	if method != strings.ToUpper(method) {
		return nil, errors.New("Invalid method verb")
	}

	reqLine.Method = method
	reqLine.RequestTarget = requestTarget
	httpVersionSplit := strings.Split(httpVersionString, "/")
	if len(httpVersionSplit) != 2 {
		return nil, errors.New("Invalid http version string")
	}

	reqLine.HttpVersion = httpVersionSplit[1]

	return &reqLine, nil

}
