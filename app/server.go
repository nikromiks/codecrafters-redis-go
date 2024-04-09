package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"

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

	var mu sync.RWMutex
	kvs := make(map[string]string)

	for {
		// Block until we receive an incoming connection
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Handle client connection
		go handleClient(&conn, &mu, &kvs)
	}
}

func handleClient(conn *net.Conn, mu *sync.RWMutex, kvs *map[string]string) {
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
		}

		// ECHO
		if v.Type() == resp.Array && strings.Contains(val, "[echo") {
			wr.WriteString(v.Array()[1].String())
		}

		// SET
		if v.Type() == resp.Array && strings.Contains(val, "[set") {
			key := v.Array()[1].String()

			mu.Lock()
			(*kvs)[key] = v.Array()[2].String()
			mu.Unlock()

			wr.WriteSimpleString("OK")
		}

		// GET
		if v.Type() == resp.Array && strings.Contains(val, "[get") {
			key := v.Array()[1].String()

			mu.Lock()
			val := (*kvs)[key]
			mu.Unlock()

			wr.WriteString(val)
		}

		(*conn).Write(wbuf.Bytes())
	}
}
