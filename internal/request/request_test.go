package request

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := min(cr.pos + cr.numBytesPerRead, len(cr.data))
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n 
	if n > cr.numBytesPerRead {
		n = cr.numBytesPerRead
		cr.pos -= n - cr.numBytesPerRead
	}
	return n, nil
}


func TestGoodGetRequestLine(t *testing.T) {
	testRequest := "GET / HTTP/1.1\r\nnHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"
	r, err := RequestFromReader(strings.NewReader(testRequest))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestGoodGetLineWithPath(t *testing.T) {
	testRequest := "GET /coffee HTTP/1.1\r\nnHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"
	r, err := RequestFromReader(strings.NewReader(testRequest))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func GoodPOSTRequestWithPath(t *testing.T) {
	testRequest := "POST /coffee HTTP/1.1\r\nnHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"
	r, err := RequestFromReader(strings.NewReader(testRequest))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "POST", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func InvalidNumberOfPartsInRequestLine(t *testing.T) {
	testRequest := "/coffee HTTP/1.1\r\nnHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"
	_, err := RequestFromReader(strings.NewReader(testRequest))
	assert.Error(t, err)
}

func InvalidMethodRequestLine(t *testing.T) {
	testRequest := "get /coffee HTTP/1.1\r\nnHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"
	_, err := RequestFromReader(strings.NewReader(testRequest))
	assert.Error(t, err)

	testRequest = "/coffee HTTP/1.1 POST\r\nnHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"
	_, err = RequestFromReader(strings.NewReader(testRequest))
	assert.Error(t, err)
}

func InvalidVersionInRequestLine(t *testing.T) {
	testRequest := "POST /coffee HTTP/1.4\r\nnHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"
	_, err := RequestFromReader(strings.NewReader(testRequest))
	assert.Error(t, err)

	testRequest = "POST /coffee HTTP/2.3\r\nnHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"
	_, err = RequestFromReader(strings.NewReader(testRequest))
	assert.Error(t, err)
}
