package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	n, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("cant resolve udp addr: %s\n", err.Error())
	}

	con, err := net.DialUDP("udp", nil, n)
	if err != nil {
		log.Fatalf("cant dial udp addr: %s\n", err.Error())
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString(byte('\n'))
		if err != nil {
			log.Fatalf("cant read from reader: %s", err.Error())
		}
		_, err = con.Write([]byte(input))
		if err != nil {
			log.Fatalf("cant write to addr: %s", err.Error())
		}
	}
}
