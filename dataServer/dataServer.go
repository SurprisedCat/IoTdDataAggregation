package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"../common"
	"../simssl"
)

func main() {

	//解析地址
	tcpAddr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:7676")
	if common.CheckErr(err, "ResolveTCPAddr") {
		return
	}

	//设置监听地址
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if common.CheckErr(err, "ListenTCP") {
		return
	}

	fmt.Println("Start wait for client.")
	for {
		//监听
		conn, err := listener.Accept()
		if common.CheckErr(err, "Accept") {
			continue
		}

		//消息处理函数
		go svrConnHandler(conn)
	}

}

//连接处理函数
func svrConnHandler(conn net.Conn) {
	if conn == nil {
		fmt.Println("Client connect failed")
		return
	}
	fmt.Println("Client connect success :", conn.RemoteAddr().String())
	conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
	//buf := make([]byte, 2048)
	defer conn.Close()
	//输出接收到的信息

	dec := gob.NewDecoder(conn)
	recPacket := &simssl.SimSsl{}
	dec.Decode(recPacket)

	/****************发送 serverReply****************/
	//get serverID=hostname+uid
	hostname, err := os.Hostname()
	if common.CheckErr(err, "cannot get Hostname ") {
		return
	}
	serverID := []byte(hostname + strconv.FormatInt(int64(os.Getuid()), 10))

	var serverReply simssl.SimSsl
	switch recPacket.ContentType {
	case 0x01:
		//generate the server hello packet
		serverReply, err = simssl.GenerateServerHello(recPacket.ClientID, serverID, recPacket.RandomInit, recPacket.EncryptKey, recPacket.ExpirationTime)
		if common.CheckErr(err, "simssl.GenerateServerHello") {
			return
		}
		fmt.Printf("Received : %+v", recPacket)

	case 0x03:
		//generate the server erase packet
		serverReply, err = simssl.GenerateServerErase(recPacket.ClientID, serverID)
		if common.CheckErr(err, "simssl.GenerateServerErase") {
			return
		}
		fmt.Printf("Received : %+v", recPacket)

	default:
		serverReply = simssl.SimSsl{}
	}
	//发送
	enc := gob.NewEncoder(conn)
	err = enc.Encode(&serverReply)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	/****************发送 serverReply****************/
	return
}
