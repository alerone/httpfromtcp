# HTTP/1.1 made in golang

This project is an implementation of the HyperText Transfer Protocol version 1.1 in golang. This is a project made step by step from a tutorial on boot.dev.

---

## Features

This is a project only for learning purposes.

- Parsing of HTTP requests.
- Generation of HTTP responses.
- Basic HTTP server creation.
- Support for 3 different status codes (200, 400, 500)
- Basic Connections (not keep-alive)
- Transfer chunked encoding

---

## Project Structure

```
 .
├──  assets
│   └──  vim.mp4
├──  cmd
│   ├──  httpserver
│   │   └──  main.go
│   ├──  tcplistener
│   │   └──  main.go
│   └──  udpsender
│       └──  main.go
├──  go.mod
├──  go.sum
├──  internal
│   ├──  headers
│   │   ├──  headers.go
│   │   └──  headers_test.go
│   ├──  request
│   │   ├──  request.go
│   │   └──  request_test.go
│   ├──  response
│   │   ├──  errors.go
│   │   └──  response.go
│   └──  server
│       ├──  handler.go
│       └──  server.go
├──  messages.txt
└──  README.md
```

---

## Requirements

- [Go](https://golang.org/dl/)  installed (version 1.18+)

---

## Start

To start the http example server just write on a terminal:

```bash
go run ./cmd/httpserver/
```

The server will listen on port `:42069`.

---

## Ejemplo de Uso


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

---
