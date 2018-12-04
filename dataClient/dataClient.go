package main

import (
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"

	"../common"
	"../database"
	"../simssl"
)

func main() {

	database.ConnTest()
	clientAuth()
}

func clientAuth() {
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

	/*******************接收0x02***********************/
	recPacket := &simssl.SimSsl{}
	dec.Decode(recPacket)
	//fmt.Printf("Received 1: %+v\n", recPacket)

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
			/*发送失败包 0x03*/
			err = enc.Encode(clientFailed)
			if err != nil {
				log.Fatal("encode error:", err)
			}
			fmt.Printf("Send 2: %+v\n", clientFailed)

			/*******************接收0x04***********************/
			recPacket = &simssl.SimSsl{}
			dec.Decode(recPacket)
			if recPacket.ContentType == 0x04 {
				fmt.Printf("Received 2:%+v\n", recPacket)
				//在文件中中清除秘钥
				if ioutil.WriteFile("key.txt", nil, 0644) != nil {
					panic("Key cannot be erased from files")
				}

			}
		} else {
			//秘钥写入文件
			buf := make([]byte, 8, 24)
			binary.PutVarint(buf, recPacket.ExpirationTime)
			//concat the encrypt key and ExpirationTime
			buf = append(message.EncryptKey[:], buf...)
			if ioutil.WriteFile("key.txt", buf, 0644) != nil {
				panic("Key cannot be written into files")
			}
			/* 从文件中读取测试*/
			contents, err := ioutil.ReadFile("key.txt")
			if err != nil {
				panic("Key Read failed")
			}
			fmt.Printf("%d\n", contents[:15])
			timeEx, _ := binary.Varint(contents[16:23])
			fmt.Printf("%d，%d\n", timeEx, recPacket.ExpirationTime)
		}
	} else {
		panic("Unexpected packet type")
	}
	/*******************接收***********************/

}
