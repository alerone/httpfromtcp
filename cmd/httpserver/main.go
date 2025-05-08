package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alerone/httpfromtcp/internal/headers"
	"github.com/alerone/httpfromtcp/internal/request"
	"github.com/alerone/httpfromtcp/internal/response"
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
func routeServing(w *response.Writer, r *request.Request) {
	switch r.RequestLine.RequestTarget {
	case "/yourproblem":
		{
			w.WriteStatusLine(400)
			hdrs := headers.NewHeaders()
			hdrs.Set("Content-Type", "text/html")

			bdy := `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`
			w.WriteHeaders(hdrs)
			w.WriteBody([]byte(bdy))
		}
	case "/myproblem":
		{
			w.WriteStatusLine(500)
			hdrs := headers.NewHeaders()
			hdrs.Set("Content-Type", "text/html")

			bdy := `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`
			w.WriteHeaders(hdrs)
			w.WriteBody([]byte(bdy))
		}
	default:
		w.WriteStatusLine(500)
		hdrs := headers.NewHeaders()
		hdrs.Set("Content-Type", "text/html")
		bdy := `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`

		w.WriteHeaders(hdrs)
		w.WriteBody([]byte(bdy))

	}

}
