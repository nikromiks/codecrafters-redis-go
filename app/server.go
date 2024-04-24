package main

import (
	"bytes"
	"fmt"
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/controller"
	"github.com/codecrafters-io/redis-starter-go/app/db"
	"github.com/tidwall/resp"
)

func main() {
	config := config.New()
	db := db.New()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		log.Fatal("Failed to bind to port ", config.Port)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		log.Println("New connection: ", conn.RemoteAddr())

		go handleConnection(conn, db, config)
	}
}

func handleConnection(conn net.Conn, db *db.DB, config *config.Config) {
	defer conn.Close()

	reader := resp.NewReader(conn)
	for {
		value, _, err := reader.ReadValue()
		if err != nil {
			break
		}
		var buf bytes.Buffer
		writer := resp.NewWriter(&buf)

		controller.Handle(&value, writer, db, config)

		if buf.Len() != 0 {
			conn.Write(buf.Bytes())
		}
	}
}
