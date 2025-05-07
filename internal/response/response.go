package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/alerone/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	okStatus                  StatusCode = 200
	badRequestStatus          StatusCode = 400
	internalServerErrorStatus StatusCode = 500
)

const(
	defaultCntLen = "Content-Length"
	defaultConn = "Connection"
	defaultCntType = "Content-Type"
)

var codeReasons = map[StatusCode]string{
	okStatus:                  "OK",
	badRequestStatus:          "Bad Request",
	internalServerErrorStatus: "Internal Server Error",
}


func WriteStatusLine(w io.Writer, statusCode StatusCode) {
	reason, ok := codeReasons[statusCode]
	if !ok {
		response := fmt.Appendf(nil, "HTTP/1.1 %d\n", statusCode)
		w.Write(response)
		return
	}

	response := fmt.Appendf(nil, "HTTP/1.1 %d %s\n", statusCode, reason)
	w.Write(response)
	return
}

func WriteHeaders (w io.Writer, headers headers.Headers) error {
	val, ok := headers.Get(defaultCntLen)
	if !ok {
		return fmt.Errorf("error while getting default header: %s\n", defaultCntLen)
	}
	_, err := w.Write(fmt.Appendf(nil, "%s: %s\n", defaultCntLen, val))
	if err != nil {
		return fmt.Errorf("writing header error: %s\n", err.Error())
	}
	

	val, ok = headers.Get(defaultConn)
	if !ok {
		return fmt.Errorf("error while getting default header: %s\n", defaultConn)
	}
	_, err = w.Write(fmt.Appendf(nil, "%s: %s\n", defaultConn, val))

	if err != nil {
		return fmt.Errorf("writing header error: %s\n", err.Error())
	}

	val, ok = headers.Get(defaultCntType)
	if !ok {
		return fmt.Errorf("error while getting default header: %s\n", defaultCntType)
	}
	_, err = w.Write(fmt.Appendf(nil, "%s: %s\n", defaultCntType, val))

	if err != nil {
		return fmt.Errorf("writing header error: %s\n", err.Error())
	}

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	defaults := headers.NewHeaders()

	defaults["content-length"] = strconv.Itoa(contentLen)
	defaults["connection"] = "close"
	defaults["content-type"] = "text/plain"

	return defaults
}

