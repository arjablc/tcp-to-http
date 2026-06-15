// NOTE: This implementation is on older ReadAll approach
// which would block the execution untill all data is read
// when used with chuncked data from reader
package request

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

func archivedRequestFromReader(reader io.Reader) (*Request, error) {
	rawBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.New("Failed reading data")
	}
	reqLine, _, err := archivedParseReqLine(rawBytes)
	if err != nil {
		return nil, err
	}
	var req Request
	req.RequestLine = *reqLine

	return &req, nil
}

func archivedParseReqLine(rawBytes []byte) (*RequestLine, int, error) {
	crlfIdx := bytes.Index(rawBytes, []byte(crlf))
	if crlfIdx < 0 {
		return nil, 0, nil
	}
	reqLineString := string(rawBytes[:crlfIdx])
	reqLine, err := reqLineFromString(reqLineString)
	if err != nil {
		return nil, 0, err
	}
	return reqLine, len(rawBytes), nil
}

func archiveParseReqLine(rawBytes []byte) (*RequestLine, error) {
	// you need to check if there even is a crlf or not
	crlfIdx := bytes.Index(rawBytes, []byte(crlf))
	if crlfIdx < 0 {
		return nil, errors.New("Invalid format")
	}
	reqLineString := string(rawBytes[:crlfIdx])
	reqLine, err := reqLineFromString(reqLineString)
	if err != nil {
		return nil, err
	}
	return reqLine, nil
}

func archiveReqLineFromString(in string) (*RequestLine, error) {
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
