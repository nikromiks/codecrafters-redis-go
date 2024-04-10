package main

import (
	"bytes"
	"io"
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/controller"
	"github.com/codecrafters-io/redis-starter-go/app/db"

	"github.com/tidwall/resp"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		log.Fatal("Failed to bind to port 6379")
	}
	defer l.Close()

	log.Println("Server is listening on port 6379")

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
