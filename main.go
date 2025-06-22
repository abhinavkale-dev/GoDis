package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	fmt.Println("Listening on port :6969")

	l, err := net.Listen("tcp", ":6969")
	if err != nil {
		fmt.Println(err)
		return
	}

	aof, err := NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

	aof.Read(func(value Value) {
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		fmt.Printf("Replaying commands of User: %s", command)
		for _, v := range args {
			fmt.Printf(" %q", v.bulk)
		}
		fmt.Println()

		handler, ok := handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			return
		}

		handler(args)
	})

	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		if value.typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewWriter(conn)

		handler, ok := handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		if command == "SET" || command == "HSET" {
			aof.Write(value)
		}

		result := handler(args)
		writer.Write(result)
	}
}
