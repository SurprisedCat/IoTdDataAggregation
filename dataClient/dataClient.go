package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"../common"
	"../simssl"
)

func main() {

	//Get the hostname of the machine
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	//clientID = hostname + uid
	clientID := hostname + strconv.FormatInt(int64(os.Getuid()), 10)
	/*send client hello*/
	clientHello, err := simssl.GenerateClientHello([]byte(clientID))
	if common.CheckErr(err, "simssl.GenerateClientHello") {
		return
	}

	//Open a socket for transmission
	conn, err := net.Dial("tcp", "127.0.0.1:7676")
	if err != nil {
		log.Fatal("Connection error", err)
	}
	clientConnHandler(conn, &clientHello)
	fmt.Println("dataClient")
}

func clientConnHandler(c net.Conn, message *simssl.SimSsl) {
	defer c.Close()
	/*****************发送*********************/
	// Initialize the encoder and decoder.  Normally enc and dec would be
	// bound to network connections and the encoder and decoder would
	// run in different processes.
	enc := gob.NewEncoder(c) // Will write to network.
	dec := gob.NewDecoder(c) //WIll read from network
	// Encode (send) the value.
	err := enc.Encode(message)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	/*****************发送*********************/

	/*******************接受***********************/
	recPacket := &simssl.SimSsl{}
	dec.Decode(recPacket)
	fmt.Printf("Received : %+v", recPacket)

	if recPacket.ContentType == 0x02 {
		//return false means something is wrong
		if !simssl.CheckKey(message, recPacket) {
			//Get the hostname of the machine
			hostname, err := os.Hostname()
			if err != nil {
				panic(err)
			}
			//clientID = hostname + uid
			clientID := hostname + strconv.FormatInt(int64(os.Getuid()), 10)
			clientFailed, err := simssl.GenerateClientFailed([]byte(clientID), []byte("Known"))
			if common.CheckErr(err, "simssl.GenerateClientFailed") {
				return
			}
			/*发送失败包*/
			err = enc.Encode(clientFailed)
			if err != nil {
				log.Fatal("encode error:", err)
			}
			fmt.Printf("Received : %+v", clientFailed)
			recPacket = &simssl.SimSsl{}
			dec.Decode(recPacket)

			if recPacket.ContentType == 0x04 {
				fmt.Printf("%+v", recPacket)
				//在redis中清除秘钥
			}
		} else {
			fmt.Println("sd") //在redis中写入client和秘钥
		}
	} else {
		panic("Unexpected packet type")
	}
	/*******************接收***********************/

}
