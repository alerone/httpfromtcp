<div align="center">

# HTTP/1.1 made in golang

![GitHub top language](https://img.shields.io/github/languages/top/alerone/httpfromtcp?color=%2377CDFF)
![GitHub last commit](https://img.shields.io/github/last-commit/alerone/httpfromtcp?color=%23bc0bbf)
![GitHub Created At](https://img.shields.io/github/created-at/alerone/httpfromtcp?color=%230dba69)
![GitHub repo size](https://img.shields.io/github/repo-size/alerone/httpfromtcp?color=%23390385)

<br>

<img src="https://github.com/user-attachments/assets/1bc11981-5738-4a70-aa24-2a9a6eea563a" alt="golang y http" width="250" height="250"/>

</div>

This project is an implementation of the HyperText Transfer Protocol version 1.1 in golang, made following the course from ![boot.dev](https://www.boot.dev/courses/learn-http-protocol-golang)


## Features

This is a project only for learning purposes.

- Parsing of HTTP requests.
- Generation of HTTP responses.
- Basic HTTP server creation.
- Support for 3 different status codes (200, 400, 500)
- Basic Connections (not keep-alive)
- Transfer chunked encoding


## Project Structure

```
.
├── assets
│   └── vim.mp4
├── cmd
│   ├── httpserver
│   │   └── main.go
│   ├──  tcplistener
│   │   └── main.go
│   └──  udpsender
│       └──  main.go
├── go.mod
├── go.sum
├── internal
│   ├── headers

│   │   ├── headers.go
│   │   └── headers_test.go
│   ├── request
│   │   ├── request.go
│   │   └── request_test.go
│   ├── response
│   │   ├── errors.go
│   │   └── response.go
│   └── server
│       ├── handler.go
│       └── server.go
├── messages.txt
└── README.md
```

## Requirements

- [Go](https://golang.org/dl/)  installed (version 1.18+)

## Start

To start the http example server just write on a terminal:

```bash
go run ./cmd/httpserver/
```

The server will listen on port `:42069`.

## Usage

You can test the server with `curl`
```bash
curl http://localhost:42069/
```
There are 5 routes in this example server:
- /myproblem
- /yourproblem
- /httpbin/...
- /video/

Any other request to the server will respond with a 200 OK and a message body. There is an implementation of a reverse proxy on 
`/httpbin` where you can redirect the request to `httpbin.org`.

## Create your own server

To create your own server you can initialize the Server class listening on any port like this:

```go
server, err := server.Serve(port, handleFunction)
if err != nil {
    log.Fatalf("Error starting server: %s", err)
}
defer server.Close()
```

To handle the requests from the server u must pass a `Handler` function to the Serve func. a `Handler` function has this structure:

```go
type Handler func(w *response.Writer, req *request.Request) 
```
The response.Writer lets the user manage the response Status Line (status code), the Headers, the Body, an optional
Chunked Body and optional Trailers for this optional Chunked Body.

