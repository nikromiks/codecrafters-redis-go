package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/tidwall/resp"
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

	for {
		// Read data
		buf := make([]byte, 1024)
		n, err := (*conn).Read(buf)
		if err != nil {
			return
		}

		log.Println("Received data", string(buf[:n]))

		rd := resp.NewReader(bytes.NewBuffer(buf))
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Read %s\n", v.String())

		var wbuf bytes.Buffer
		wr := resp.NewWriter(&wbuf)

		val := v.String()

		// PING
		if val == "[ping]" {
			wr.WriteSimpleString("PONG")
			(*conn).Write(wbuf.Bytes())
			continue
		}

		// ECHO
		if v.Type() == resp.Array && strings.Contains(val, "[echo") {
			wr.WriteString(v.Array()[1].String())
			(*conn).Write(wbuf.Bytes())
			continue
		}
	}
}
