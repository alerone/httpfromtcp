package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/alerone/httpfromtcp/internal/request"
)

type Server struct {
	State    *atomic.Bool
	listener net.Listener
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("starting http server error: %s", err.Error())
	}

	var state atomic.Bool
	state.Store(true)

	server := &Server{
		State:    &state,
		listener: listener,
	}

	go func() {
		for {
			server.listen()
		}
	}()

	return server, nil
}

func (s *Server) Close() error {
	s.State.Store(false)
	err := s.listener.Close()
	if err != nil {
		return fmt.Errorf("closing server error: %s", err.Error())
	}

	return nil
}

func (s *Server) listen() {
	if s.State.Load() == false {
		return
	}
	conn, err := s.listener.Accept()
	if err != nil {
		// fmt.Println(err.Error())
		return
	}
	go s.handle(conn)
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	_, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	res := `HTTP/1.1 200 OK
	Content-Type: text/plain

	Hello World!
	`
	conn.Write([]byte(res))
}
