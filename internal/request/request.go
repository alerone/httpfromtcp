package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	reqBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	reqParts := strings.Split(string(reqBytes), "\r\n")
	requestLine, err := parseRequestLine(reqParts[0])

	if err != nil {
		return nil, err
	}

	return &Request{RequestLine: *requestLine}, nil
}

func parseRequestLine(data []byte) (*RequestLine, error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return nil, fmt.Errorf("could not find CRLF in request-line")
	}

	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, err
	}

	return requestLine, nil
}

func requestLineFromString(line string) (*RequestLine,  error) {
	parts := strings.Fields(line)
	if len( parts ) != 3 {
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
	if httpPart != "1.1" {
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

