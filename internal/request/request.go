package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/alerone/httpfromtcp/internal/headers"
)

const (
	rqStateInitialized requestState = iota
	rqStateParsingHeaders
	rqStateParsingBody
	rqStateDone
	bufferSize int = 8
	crlf           = "\r\n"
)

type requestState int

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	state       requestState
}

func (r *Request) parse(data []byte) (n int, err error) {
	totalBytesParsed := 0

	for r.state != rqStateDone {
		n, err = r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}

		if n == 0 {
			return totalBytesParsed, nil
		}

		totalBytesParsed += n
	}
	return n, nil
}

func (r *Request) parseSingle(data []byte) (n int, err error) {
	switch r.state {
	case rqStateInitialized:
		{
			n, rl, err := parseRequestLine(data)
			if err != nil {
				return 0, err
			}
			if n == 0 {
				return 0, nil
			}

			r.RequestLine = *rl
			r.state = rqStateParsingHeaders
			return n, nil
		}
	case rqStateParsingHeaders:
		{
			n, done, err := r.Headers.Parse(data)
			if err != nil {
				return 0, err
			}
			if done {
				r.state = rqStateParsingBody
			}
			return n, nil
		}
	case rqStateParsingBody:
		{
			cl, ok := r.Headers.Get("Content-Length")
			if !ok {
				r.state = rqStateDone
				return 0, nil
			}
			r.Body = append(r.Body, data...)
			clNum, err := strconv.Atoi(cl)
			if err != nil {
				return 0, fmt.Errorf("invalid content length not an integer: %s", cl)
			}
			if len(r.Body) > clNum {
				return 0, fmt.Errorf("body length greater than content length")
			} else if len(r.Body) == clNum {
				r.state = rqStateDone
			}
			return len(data), nil
		}
	case rqStateDone:
		return -1, fmt.Errorf("error: trying to parse data in done state")
	default:
		return -1, fmt.Errorf("error: unknown state")
	}
}

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := &Request{
		state:   rqStateInitialized,
		Headers: headers.NewHeaders(),
		Body: []byte(""),
	}
	readToIndex := 0
	buf := make([]byte, bufferSize, bufferSize)
	for request.state != rqStateDone {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, 2*len(buf))
			copy(newBuf, buf)
			buf = newBuf
		}
		readCount, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				request.state = rqStateDone
				break
			}
			return nil, fmt.Errorf("error while reading request: %s", err.Error())
		}

		readToIndex += readCount

		pn, err := request.parse(buf[:readToIndex])
		if err != nil {
			return nil, fmt.Errorf("error while parsing request: %s", err.Error())
		}

		copy(buf, buf[pn:])
		readToIndex -= pn
	}
	return request, nil
}

func parseRequestLine(data []byte) (n int, res *RequestLine, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, nil, nil
	}

	requestLineText := string(data[:idx])
	res, err = requestLineFromString(requestLineText)
	if err != nil {
		return -1, nil, err
	}

	return idx + 2, res, nil
}

func requestLineFromString(line string) (*RequestLine, error) {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", line)
	}

	method := parts[0]
	if checkMethodIsUpper(method) != true {
		return nil, fmt.Errorf("http bad request: method not all upper chars -> method: %s", method)
	}

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", line)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-protocol %s", httpPart)
	}

	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: parts[1],
		HttpVersion:   version,
	}, nil
}

func checkMethodIsUpper(method string) bool {
	for _, r := range method {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
