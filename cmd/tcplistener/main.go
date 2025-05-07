package main

import (
	"fmt"
	"log"
	"net"

	"github.com/alerone/httpfromtcp/internal/request"
)

func main() {
	port := ":42069"
	typeConn := "tcp"

	listener, err := net.Listen(typeConn, port)
	if err != nil {
		log.Fatalf("couldn't listen on port %s: %s", port, err.Error())
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("couldnt accept connection: %s", err.Error())
		}
		fmt.Println("Connection accepted")
		rq, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", rq.RequestLine.Method)
		fmt.Printf("- Target: %s\n", rq.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", rq.RequestLine.HttpVersion)

		fmt.Println("Headers:")
		for key, val := range rq.Headers {
			fmt.Printf("- %s: %s\n", key, val)
		}

		fmt.Println("Body:")
		fmt.Println(string(rq.Body))
		fmt.Println("Connection closed")
	}
}
