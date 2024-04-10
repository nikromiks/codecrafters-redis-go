package controller

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/db"
	"github.com/tidwall/resp"
)

func Handle(v *resp.Value, wr *resp.Writer, d *db.DB) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered after error: ", r)
		}
	}()
	command := ""

	if v.Type() == resp.Array && len(v.Array()) != 0 {
		command = strings.ToLower(v.Array()[0].String())
	}

	switch {
	case command == "ping":
		ping(wr)
	case command == "echo":
		echo(v, wr)
	case command == "set":
		set(v, wr, d)
	case command == "get":
		get(v, wr, d)
	case command == "info":
		info(v, wr, d)
	}
}

func ping(wr *resp.Writer) {
	wr.WriteSimpleString("PONG")
}

func echo(v *resp.Value, wr *resp.Writer) {
	if len(v.Array()) < 1 {
		log.Panic("Incorrect echo")
		return
	}
	wr.WriteString(v.Array()[1].String())
}

func get(v *resp.Value, wr *resp.Writer, d *db.DB) {
	key := v.Array()[1].String()

	val, ok := d.Get(key)
	if !ok {
		wr.WriteNull()
		return
	}
	if val.Exp != -1 && val.Exp < time.Now().UnixMilli() {
		wr.WriteNull()
		return
	}

	wr.WriteString(val.Val)
}

func set(v *resp.Value, wr *resp.Writer, d *db.DB) {
	key := v.Array()[1].String()

	newVal := db.Value{
		Val: v.Array()[2].String(),
		Exp: -1,
	}

	if len(v.Array()) >= 5 && v.Array()[3].String() == "px" {
		expVal, err := strconv.ParseInt(v.Array()[4].String(), 10, 64)
		if err != nil {
			log.Panicf("Error input value %s", v.Array()[4].String())
		}

		newVal.Exp = time.Now().UnixMilli() + expVal
	}
	d.Set(key, newVal)

	wr.WriteSimpleString("OK")
}

func info(v *resp.Value, wr *resp.Writer, d *db.DB) {
	wr.WriteString("role:master")
	wr.WriteString("connected_slaves:0")
	wr.WriteString("master_replid:123")
	wr.WriteString("master_repl_offset:0")
}
