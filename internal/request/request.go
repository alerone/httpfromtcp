package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

const (
	rqStateInitialized requestState = iota
	rqStateDone
	bufferSize int = 8
)

type requestState int

type Request struct {
	RequestLine RequestLine
	state       requestState
}

func (r *Request) parse(data []byte) (n int, err error) {
	if r.state == rqStateInitialized {
		n, rl, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}

		r.RequestLine = *rl
		r.state = rqStateDone
	} else if r.state == rqStateDone {
		return -1, fmt.Errorf("error: trying to parse data in done state")
	} else {
		return -1, fmt.Errorf("error: unknown state")
	}
	return n, nil
}

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := &Request{
		state: rqStateInitialized,
	}
	readToIndex := 0
	buf := make([]byte, bufferSize, bufferSize)
	for request.state != rqStateDone{
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
	idx := bytes.Index(data, []byte("\r\n"))
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
