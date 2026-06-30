package headers

import (
	"bytes"
	"errors"
	"fmt"
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
	val, ok := h[key]
	if !ok {
		h[key] = value
		return
	}
	val = fmt.Sprintf("%s, %s", val, value)
	h[key] = val

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

	parts := bytes.SplitN(data[:crlfIdx], headerSep, 2)
	if len(parts) < 2 {
		return 0, false, errors.New(": not found")
	}
	headerKey := strings.ToLower(string(parts[0]))
	headerValue := bytes.TrimSpace(parts[1])

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

	h.Set(headerKey, string(headerValue))
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
