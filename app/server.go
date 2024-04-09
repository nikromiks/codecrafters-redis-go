package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// Listen for incoming connections
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
		return
	}

	// Ensure we teardown the server when the program exits
	defer l.Close()

	fmt.Println("Server is listening on port 6379")

	for {
		// Block until we receive an incoming connection
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Handle client connection
		go handleClient(&conn)
	}
}

func handleClient(conn *net.Conn) {
	// Ensure we close the connection after we're done
	defer (*conn).Close()

	// Read data
	buf := make([]byte, 1024)
	for {
		n, err := (*conn).Read(buf)
		if err != nil {
			return
		}

		log.Println("Received data", string(buf[:n]))

		if bytes.Equal(buf[:n], []byte("*1\r\n$4\r\nping\r\n")) {
			(*conn).Write([]byte("+PONG\r\n"))
		}
	}
}
