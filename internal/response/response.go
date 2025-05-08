package response

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/alerone/httpfromtcp/internal/headers"
)

type StatusCode int
type writerState int

const (
	OkStatus                  StatusCode = 200
	BadRequestStatus          StatusCode = 400
	InternalServerErrorStatus StatusCode = 500
)

const (
	defaultCntLen  = "Content-Length"
	defaultConn    = "Connection"
	defaultCntType = "Content-Type"
)

const (
	initState writerState = iota
	writingStatus
	writingHdrs
	writingBody
)

var codeReasons = map[StatusCode]string{
	OkStatus:                  "OK",
	BadRequestStatus:          "Bad Request",
	InternalServerErrorStatus: "Internal Server Error",
}

type Writer struct {
	statusCode StatusCode
	Headers    headers.Headers
	body       []byte
	out        io.Writer
	state      writerState
}

func NewWriter(out io.Writer) Writer {
	return Writer{
		state: initState,
		out:   out,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != initState {
		return &InvalidOrderResponseWriter{
			expectedState: initState,
			actual:        w.state,
		}
	}
	w.statusCode = statusCode
	w.state = writingStatus
	writeStatusLine(w.out, statusCode)
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != writingStatus {
		return &InvalidOrderResponseWriter{
			expectedState: writingStatus,
			actual:        w.state,
		}
	}
	w.state = writingHdrs
	w.Headers = headers
	for key, val := range headers {
		w.out.Write(fmt.Appendf(nil, "%s: %s\r\n", key, val))
	}

	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != writingHdrs {
		return 0, &InvalidOrderResponseWriter{
			expectedState: writingStatus,
			actual:        w.state,
		}
	}
	bodyCount := len(p)
	if _, ok := w.Headers.Get("Content-Length"); !ok {
		w.Headers.Set("Content-Length", strconv.Itoa(bodyCount))
		w.out.Write(fmt.Appendf(nil, "Content-Length: %d\r\n", bodyCount))
	}
	w.out.Write([]byte("\r\n"))
	w.state = writingBody
	w.body = p
	w.out.Write(p)
	return len(p)+2, nil
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.state != writingBody {
		if w.state != writingHdrs {
			return 0, &InvalidOrderResponseWriter{
				expectedState: writingHdrs,
				actual:        w.state,
			}
		}
		w.out.Write([]byte("\r\n"))
	}
	w.state = writingBody

	encoding := fmt.Appendf(nil, "%X\r\n%s\r\n", len(p), string(p))
	w.out.Write(encoding)
	return len(encoding), nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.state != writingBody {
		return 0, &InvalidOrderResponseWriter{
			expectedState: writingBody,
			actual:        w.state,
		}
	}
	encoding := fmt.Appendf(nil, "%X\r\n\r\n", 0)
	w.out.Write(encoding)

	return len(encoding), nil
}

func writeStatusLine(w io.Writer, statusCode StatusCode) {
	reason, ok := codeReasons[statusCode]
	if !ok {
		response := fmt.Appendf(nil, "HTTP/1.1 %d\r\n", statusCode)
		w.Write(response)
		return
	}

	response := fmt.Appendf(nil, "HTTP/1.1 %d %s\r\n", statusCode, reason)
	w.Write(response)
	return
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	defaults := headers.NewHeaders()

	defaults.Set("Content-Length", strconv.Itoa(contentLen))
	defaults.Set("Connection", "close")
	defaults.Set("Content-Type", "text/plain")

	return defaults
}

func (ws *writerState) String() string {
	out := new(bytes.Buffer)

	switch *ws {
	case initState:
		out.WriteString("not writing")
	case writingStatus:
		out.WriteString("writing status line")
	case writingHdrs:
		out.WriteString("writing headers")
	case writingBody:
		out.WriteString("writing body")
	default:
		out.WriteString("error order unknown")
	}

	return out.String()
}
