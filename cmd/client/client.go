package client

import (
	"log"
	"net"
)

type Client net.Conn

const (
	port = ":8000"
)

func main() {
	conn, err := net.Dial("tcp", "localhost"+port)
	if err != nil {
		log.Fatal("ERROR - Error connecting to server:", err)
	}

	cli := Client(conn)
}
