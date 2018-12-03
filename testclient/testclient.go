package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"

	"../simssl"
)

func main() {
	clientHello, _ := simssl.GenerateClientHello([]byte("TEST"))
	fmt.Println("start client")
	conn, err := net.Dial("tcp", "localhost:7676")
	if err != nil {
		log.Fatal("Connection error", err)
	}
	encoder := gob.NewEncoder(conn)
	encoder.Encode(&clientHello)
	conn.Close()
	fmt.Println("done")
}
