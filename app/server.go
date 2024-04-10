package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/controller"
	"github.com/codecrafters-io/redis-starter-go/app/db"

	"github.com/tidwall/resp"
)

func main() {
	port := "6379"
	if len(os.Args) > 2 && os.Args[1] == "--port" {
		port = os.Args[2]
	}

	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Fatal("Failed to bind to port ", port)
	}
	defer l.Close()

	log.Println("Server is listening on port ", port)

	db := db.New()

	for {
		// Block until we receive an incoming connection
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error:", err)
			continue
		}

		go handleClient(&conn, db)
	}
}

func handleClient(conn *net.Conn, db *db.DB) {
	// Ensure we close the connection after we're done
	defer (*conn).Close()

	for {
		rd := resp.NewReader((*conn))
		v, _, err := rd.ReadValue()
		if err == io.EOF || err != nil {
			log.Println(err)
			break
		}

		log.Printf("Read %s\n", v.String())

		var buf bytes.Buffer
		wr := resp.NewWriter(&buf)

		controller.Handle(&v, wr, db)

		if buf.Len() != 0 {
			(*conn).Write(buf.Bytes())
		}
	}
}
