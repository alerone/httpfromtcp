package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

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

var routes = map[string]server.Handler{
	"/myproblem":   myProblemRoute,
	"/yourproblem": yourProblemRoute,
	"/httpbin":     httpbinRoute,
}

func routeServing(w *response.Writer, r *request.Request) {
	for route, handler := range routes {
		if strings.HasPrefix(r.RequestLine.RequestTarget, route) {
			handler(w, r)
			return
		}
	}
	w.WriteStatusLine(response.OkStatus)
	bdy := `<html>
			  <head>
				<title>200 OK</title>
			  </head>
			  <body>
				<h1>Success!</h1>
				<p>Your request was an absolute banger.</p>
			  </body>
	        </html>
	`
	hdrs := response.GetDefaultHeaders(len(bdy))
	hdrs.Set("Content-Type", "text/html")
	w.WriteHeaders(hdrs)
	w.WriteBody([]byte(bdy))
}

func httpbinRoute(w *response.Writer, r *request.Request) {
	path := strings.TrimPrefix(r.RequestLine.RequestTarget, "/httpbin")
	url := "https://httpbin.org" + path
	res, err := http.Get(url)
	if err != nil {
		w.WriteStatusLine(400)
		bdy := `<html>
			  <head>
				<title>400 Bad request</title>
			  </head>
			  <body>
				<h1>Not found on httpbin.org</h1>
			  </body>
	        </html>
		`
		hdrs := response.GetDefaultHeaders(len(bdy))
		hdrs.Set("Content-Type", "text/html")
		w.WriteHeaders(hdrs)
		w.WriteBody([]byte(bdy))
		return
	}
	buf := make([]byte, 1024)
	w.WriteStatusLine(response.OkStatus)
	hdrs := response.GetDefaultHeaders(0)
	hdrs.Remove("content-length")
	hdrs.Set("Transfer-Encoding", "chunked")
	hdrs.Set("Host", "httpbin.org")
	w.WriteHeaders(hdrs)
	for {
		n, err := res.Body.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				w.WriteChunkedBodyDone()
				return
			}
			fmt.Println(err.Error())
			return
		}

		fmt.Println("Read: ", n)

		w.WriteChunkedBody(buf)
	}

}

func yourProblemRoute(w *response.Writer, r *request.Request) {
	w.WriteStatusLine(400)
	bdy := `<html>
	<head>
	<title>400 Bad Request</title>
	</head>
	<body>
	<h1>Bad Request</h1>
	<p>Your request honestly kinda sucked.</p>
	</body>
	</html>
	`
	hdrs := response.GetDefaultHeaders(len(bdy))
	hdrs.Set("Content-Type", "text/html")
	w.WriteHeaders(hdrs)
	w.WriteBody([]byte(bdy))
}
func myProblemRoute(w *response.Writer, r *request.Request) {
	w.WriteStatusLine(500)

	bdy := `<html>
	<head>
	<title>500 Internal Server Error</title>
	</head>
	<body>
	<h1>Internal Server Error</h1>
	<p>Okay, you know what? This one is on me.</p>
	</body>
	</html>
	`
	hdrs := response.GetDefaultHeaders(len(bdy))
	hdrs.Set("Content-Type", "text/html")
	w.WriteHeaders(hdrs)
	w.WriteBody([]byte(bdy))
}
