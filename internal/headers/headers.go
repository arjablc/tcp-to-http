package headers

import (
	"bytes"
	"errors"
	"slices"
	"strings"
)

type Headers map[string]string

var crlf string = "\r\n"
var headerSep []byte = []byte(":")
var validSpecialCharacters []rune = []rune{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

func NewHeaders() Headers {
	return make(Headers)
}

const INVALID_HEADERS string = "Invalid headers"

// Better API
func (h Headers) Set(key, value string) {
	h[key] = value
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	crlfIdx := bytes.Index(data, []byte(crlf))
	if crlfIdx < 0 {
		return 0, false, nil
	}
	if crlfIdx == 0 {
		// even with the crlf on the 0,
		// we still need to return the crlf byte count
		return 2, true, nil

	}

	headerSepIdx := bytes.Index(data, headerSep)
	if headerSepIdx < 0 {
		return 0, false, errors.New(": not found")
	}
	headerKey := strings.ToLower(string(data[:headerSepIdx]))
	headerValue := string(data[headerSepIdx+2 : crlfIdx])

	// THis only checks if there are whitespaces
	// around header
	// NOTE: bute there must not be spaces within
	// the key as well
	if strings.TrimSpace(headerKey) != headerKey {
		return 0, false, errors.New("when key is invalid")
	}

	if strings.ContainsAny(headerKey, " \t\n") {
		return 0, false, errors.New("when key is invalid")
	}
	for _, char := range headerKey {
		if !isValid(char) {
			return 0, false, errors.New("Invalid char in Key")
		}
	}

	h.Set(headerKey, strings.TrimSpace(headerValue))
	return crlfIdx + 2, false, nil
}

func isValid(c rune) bool {
	if c >= 'a' && c <= 'z' {
		return true
	}
	if c >= 'A' && c <= 'Z' {
		return true
	}
	if c > '0' && c < '9' {
		return true
	}
	return slices.Contains(validSpecialCharacters, c)
}
