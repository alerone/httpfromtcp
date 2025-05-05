package request

import (
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

func parseRequestLine(line string) (*RequestLine, error) {
	parts := strings.Fields(line)
	if len( parts ) != 3 {
		return nil, fmt.Errorf("http bad request: request-line with length not equal to 3 parts")
	}

	parseRes := RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   parts[2][len(parts[2])-3:],
	}
	if checkMethodIsUpper(parseRes.Method) != true {
		return nil, fmt.Errorf("http bad request: method not all upper chars -> method: %s", parseRes.Method)
	}

	if parts[2] != "HTTP/1.1" {
		return nil, fmt.Errorf("the only version of http is: HTTP/1.1")
	}
	
	return &parseRes, nil
}

func checkMethodIsUpper(method string) bool {
	for _, r := range method {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

