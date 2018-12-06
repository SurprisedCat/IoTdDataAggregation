package auth

import (
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"../database"
	"../simssl"
	"../utils"
)

/*
ClientAuth client authenticates the key and expiration time
*/
func ClientAuth(serveraddr []byte) bool {
	clientID := utils.GetClientID()
	/*send client hello*/
	clientHello, err := simssl.GenerateClientHello(clientID)
	if utils.CheckErr(err, "simssl.GenerateClientHello") {
		return false
	}

	//Open a socket for transmission
	conn, err := net.Dial("tcp", string(append(serveraddr, []byte(":7676")...)))
	if err != nil {
		log.Printf("Connection error: %v", err)
		return false
	}
	clientConnHandler(conn, &clientHello)
	fmt.Println("dataClient is authenticated")
	return true
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
	fmt.Printf("Received 1: %+v\n", recPacket)

	if recPacket.ContentType == 0x02 {
		//return false means something is wrong
		if !simssl.CheckKey(message, recPacket) {

			clientID := utils.GetClientID()
			clientFailed, err := simssl.GenerateClientFailed(clientID, []byte("Known"))
			if utils.CheckErr(err, "simssl.GenerateClientFailed") {
				return
			}
			/*发送失败包 0x03*/
			err = enc.Encode(clientFailed)
			if err != nil {
				log.Fatal("encode error:", err)
			}
			//fmt.Printf("Send 2: %+v\n", clientFailed)

			/*******************接收0x04***********************/
			recPacket = &simssl.SimSsl{}
			dec.Decode(recPacket)
			if recPacket.ContentType == 0x04 {
				//fmt.Printf("Received 2:%+v\n", recPacket)
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
		}
	} else {
		panic("Unexpected packet type")
	}
	/*******************接收***********************/
}

/*
GetKeyClient get the encrypt key
*/
func GetKeyClient() ([]byte, int64) {
	/* 从文件中读取测试*/
	contents, err := ioutil.ReadFile("key.txt")
	if err != nil {
		panic("Key Read failed")
	}
	timeEx, _ := binary.Varint(contents[16:24])
	return contents[:16], timeEx
}

/**-------------server-------------------**/

/*
SvrListen 开启监听认证信息
*/
func SvrListen(wg *sync.WaitGroup) {
	//解析地址
	tcpAddr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:7676")
	if utils.CheckErr(err, "ResolveTCPAddr") {
		wg.Done()
		return
	}

	//设置监听地址
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if utils.CheckErr(err, "ListenTCP") {
		wg.Done()
		return
	}

	fmt.Println("Start wait for client.")
	for {
		//监听
		conn, err := listener.Accept()
		if utils.CheckErr(err, "Accept") {
			continue
		}
		//消息处理函数
		wg.Add(1)
		go svrConnHandler(conn, wg)
	}
	//wg.Done()
}

//连接处理函数
/*
svrConnHandler 服务端处理函数
*/
func svrConnHandler(conn net.Conn, wg *sync.WaitGroup) {
	if conn == nil {
		fmt.Println("Client connect failed")
		wg.Done()
		return
	}
	fmt.Println("Client connect success :", conn.RemoteAddr().String())
	conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
	//buf := make([]byte, 2048)

	//加了这个函数防止意外退出服务
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Runtime error caught :", r)
		}
	}()
	defer wg.Done()
	defer conn.Close()

	//输出接收到的信息

	dec := gob.NewDecoder(conn)
	recPacket := &simssl.SimSsl{}
	dec.Decode(recPacket)

	//get serverID=hostname+uid
	hostname, err := os.Hostname()
	if utils.CheckErr(err, "cannot get Hostname ") {
		wg.Done()
		return
	}
	serverID := []byte(hostname + strconv.FormatInt(int64(os.Getuid()), 10))

	/****************发送 serverReply 0x02****************/
	var serverReply simssl.SimSsl
	if recPacket.ContentType == 0x01 {
		//generate the server hello packet
		serverReply, err = simssl.GenerateServerHello(recPacket.ClientID, serverID, recPacket.RandomInit, recPacket.EncryptKey, recPacket.ExpirationTime)
		if utils.CheckErr(err, "simssl.GenerateServerHello") {
			wg.Done()
			return
		}
		fmt.Printf("Received 1: %+v", recPacket)
	} else {
		serverReply = simssl.SimSsl{}
	}
	//发送0x02 或者 空的包
	enc := gob.NewEncoder(conn)
	err = enc.Encode(&serverReply)
	if err != nil {
		wg.Done()
		log.Fatal("encode error:", err)
	}
	//	Write the client info into redis
	if !database.DataServerWriteAuth(recPacket.ClientID[:], recPacket.EncryptKey[:], recPacket.ExpirationTime-time.Now().Unix()) {
		wg.Done()
		fmt.Println("DataBase write error")
	}

	/****************发送 serverReply 0x02****************/

	/*错误处理 0x03*/
	recPacket = &simssl.SimSsl{}
	dec.Decode(recPacket)
	fmt.Printf("Received 2: %+v", recPacket)

	/****************发送 serverReply 0x04****************/
	if recPacket.ContentType == 0x03 {

		//Erase the client info from redis
		//	Write the client info into redis
		if !database.DataServerEraseAuth(recPacket.ClientID[:]) {
			wg.Done()
			fmt.Println("DataBase delete error")
		}
		//generate the server erase packet
		serverReply, err = simssl.GenerateServerErase(recPacket.ClientID, serverID)
		if utils.CheckErr(err, "simssl.GenerateServerErase") {
			wg.Done()
			return
		}
		//send 0x04
		err = enc.Encode(&serverReply)
		if err != nil {
			wg.Done()
			log.Fatal("encode error:", err)
		}
	}
	/****************发送 serverReply 0x04****************/

	return
}

/*
GetValidationKeyServer the check the existation and expiration of encryption key
*/
func GetValidationKeyServer(clientID []byte) ([]byte, bool) {
	encryptedKey, vali := database.DataServerGetKey(clientID)
	if vali == false || encryptedKey == nil {
		return nil, false
	}
	return encryptedKey, true
}
