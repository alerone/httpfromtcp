package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
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
		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Println(line)
		}
		fmt.Println("Connection closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {
		defer f.Close()
		defer close(ch)
		currentLine := ""
		for {
			buffer := make([]byte, 8, 8)
			n, err := f.Read(buffer)
			if err != nil {
				if currentLine != "" {
					ch <- currentLine
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				return
			}
			str := string(buffer[:n])
			parts := strings.Split(str, "\n")
			for i := range len(parts) - 1 {
				ch <- fmt.Sprintf("%s%s", currentLine, parts[i])
				currentLine = ""
			}
			currentLine += parts[len(parts)-1]
		}
	}()

	return ch
}

func readFromFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("could not open %s: %s\n", filepath, err.Error())
	}
	defer file.Close()
	for {
		b := make([]byte, 8, 8)
		n, err := file.Read(b)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("error: %s\n", err.Error())
		}
		str := string(b[:n])
		fmt.Printf("read: %s\n", str)
	}

	return nil
}

func readFromFileFullLines(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("could not open %s: %s\n", filepath, err.Error())
	}
	defer f.Close()

	currentLine := ""
	for {
		buffer := make([]byte, 8, 8)
		n, err := f.Read(buffer)
		if err != nil {
			if currentLine != "" {
				fmt.Printf("read: %s\n", currentLine)
				currentLine = ""
			}
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("error: %s\n", err.Error())
		}
		str := string(buffer[:n])
		parts := strings.Split(str, "\n")
		for i := range len(parts) - 1 {
			fmt.Printf("read: %s%s\n", currentLine, parts[i])
			currentLine = ""
		}
		currentLine += parts[len(parts)-1]
	}

	return nil
}
