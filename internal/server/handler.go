package server

import (
	"github.com/alerone/httpfromtcp/internal/request"
	"github.com/alerone/httpfromtcp/internal/response"
)

type Handler func(w *response.Writer, req *request.Request) 

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}



