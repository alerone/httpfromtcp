package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alerone/httpfromtcp/internal/request"
	"github.com/alerone/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, routeServing)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
	defer server.Close()

	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
func routeServing(w io.Writer, r *request.Request) *server.HandlerError {
	switch r.RequestLine.RequestTarget {
	case "/yourproblem":
		{
			err := &server.HandlerError{
				StatusCode: 400,
				Message:    "Your problem is not my problem\n",
			}
			return err
		}
	case "/myproblem":
		{
			err := &server.HandlerError{
				StatusCode: 500,
				Message:    "Woopsie, my bad\n",
			}
			return err
		}
	default:
		w.Write([]byte("All good, frfr\n"))
	}

	return nil
}
