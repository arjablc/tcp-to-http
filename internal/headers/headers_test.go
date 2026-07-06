package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid 2 headers with existing headers
	headers = map[string]string{"host": "localhost:42069"}
	data = []byte("User-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, "curl/7.81.0", headers["user-agent"])
	assert.Equal(t, 25, n)
	assert.False(t, done)

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("\r\n a bunch of other stuff")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Empty(t, headers)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Invalid spacing header (leading spaces)
	headers = NewHeaders()
	data = []byte("       Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid spacing header (internal spaces)
	headers = NewHeaders()
	data = []byte("Hos\nt: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	//Test: Invalid characters in the header key

	headers = NewHeaders()
	data = []byte("Ho@st: localhost:42069\r\n Content-Type:               application/json\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.False(t, done)
	assert.Equal(t, 0, n)

	//Test: Multiple headers with same name
	headers = map[string]string{"host": "localhost:42069"}
	data = []byte("Host: loli:1234\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069, loli:1234", headers["host"])
	assert.False(t, done)
}

func TestHeaderKeysWithDigits(t *testing.T) {
	// RFC 7230 defines a header field-name as a token, and a token
	// is composed of tchars which INCLUDE DIGIT. So header keys
	// containing the digits '0' and '9' must be accepted.
	// This currently FAILS against isValid in headers.go because
	// the digit check is `c > '0' && c < '9'` — an off-by-one that
	// wrongly rejects both '0' and '9'.

	// Test: header key containing digit '0'
	headers := NewHeaders()
	data := []byte("X-Custom0: foo\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "foo", headers["x-custom0"])
	assert.False(t, done)
	assert.Equal(t, len("X-Custom0: foo\r\n"), n)

	// Test: header key containing digit '9'
	headers = NewHeaders()
	data = []byte("X-Custom9: bar\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "bar", headers["x-custom9"])
	assert.False(t, done)
	assert.Equal(t, len("X-Custom9: bar\r\n"), n)

	// Test: header key containing all boundary digits '0'..'9'
	headers = NewHeaders()
	data = []byte("X-0123456789: baz\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "baz", headers["x-0123456789"])
	assert.False(t, done)
	assert.Equal(t, len("X-0123456789: baz\r\n"), n)

	// Test: header key that is entirely digits (still a valid token)
	headers = NewHeaders()
	data = []byte("0: qux\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "qux", headers["0"])
	assert.False(t, done)
	assert.Equal(t, len("0: qux\r\n"), n)
}

func TestIsValidDigitBoundaries(t *testing.T) {
	// Locks in that every decimal digit is a valid header-name character.
	// Currently fails for '0' and '9' because of the off-by-one in
	// headers.go isValid: `c > '0' && c < '9'`.
	for r := '0'; r <= '9'; r++ {
		assert.True(t, isValid(r), "isValid should accept digit %q", string(r))
	}
}
