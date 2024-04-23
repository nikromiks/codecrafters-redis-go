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

	listener := startListener(config.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn, db, config)
	}
}

func startListener(port int) net.Listener {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatal("Failed to bind to port ", port)
	}
	return listener
}

func handleConnection(conn net.Conn, db *db.DB, config *config.Config) {
	defer conn.Close()

	for {
		reader := resp.NewReader(conn)
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
