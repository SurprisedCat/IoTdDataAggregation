package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"unsafe"

	"./simssl"
)

//错误处理函数
func checkErr(err error, extra string) bool {
	if err != nil {
		formatStr := " Err : %s\n"
		if extra != "" {
			formatStr = extra + formatStr
		}

		fmt.Fprintf(os.Stderr, formatStr, err.Error())
		return true
	}

	return false
}

//连接处理函数
func svrConnHandler(conn net.Conn) {
	fmt.Println("Client connect success :", conn.RemoteAddr().String())
	conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
	request := make([]byte, 128)
	defer conn.Close()
	for {
		readLen, err := conn.Read(request)
		if checkErr(err, "Read") {
			break
		}

		//socket被关闭了
		if readLen == 0 {
			fmt.Println("Client connection close!")
			break
		} else {
			//输出接收到的信息
			fmt.Println(string(request[:readLen]))

			time.Sleep(time.Second)
			//发送
			conn.Write([]byte("World !"))
		}

		request = make([]byte, 128)
	}
}

type test struct {
	test1 byte
	test2 uint64
}

func main() {
	simssl.GenerateClientHello([]byte("TEST"))
	test0 := test{12, 12312}
	fmt.Println(unsafe.Sizeof(test0))
	/*
		//解析地址
		tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:6666")
		if checkErr(err, "ResolveTCPAddr") {
			return
		}

		//设置监听地址
		listener, err := net.ListenTCP("tcp", tcpAddr)
		if checkErr(err, "ListenTCP") {
			return
		}

		for {
			//监听
			fmt.Println("Start wait for client.")
			conn, err := listener.Accept()
			if checkErr(err, "Accept") {
				continue
			}

			//消息处理函数
			go svrConnHandler(conn)
		}
	*/
}
