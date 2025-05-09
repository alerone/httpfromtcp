package server

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/alerone/httpfromtcp/internal/request"
	"github.com/alerone/httpfromtcp/internal/response"
)

type Server struct {
	closed   atomic.Bool
	listener net.Listener
	handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("starting http server error: %s", err.Error())
	}

	server := &Server{
		handler:  handler,
		listener: listener,
	}

	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	err := s.listener.Close()
	if err != nil {
		return fmt.Errorf("closing server error: %s", err.Error())
	}

	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	rq, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Println(err.Error())
		errorMsg := fmt.Sprintf("HTTP/1.1 400 Bad Request\r\nContent-Length: %d\r\n\r\n%s", len(err.Error()), err.Error())
		conn.Write([]byte(errorMsg))
		return
	}

	buf := new(bytes.Buffer)
	writer := response.NewWriter(buf)
	s.handler(&writer, rq)
	buf.WriteTo(conn)
}
