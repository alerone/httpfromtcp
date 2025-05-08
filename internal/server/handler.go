package server

import (
	"io"
	"log"

	"github.com/alerone/httpfromtcp/internal/request"
	"github.com/alerone/httpfromtcp/internal/response"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func (h *HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, h.StatusCode)
	hdrs := response.GetDefaultHeaders(len(h.Message))
	err := response.WriteHeaders(w, hdrs)
	w.Write([]byte(h.Message))
	if err != nil {
		log.Printf("error writing response: %s",err)
	}
}


