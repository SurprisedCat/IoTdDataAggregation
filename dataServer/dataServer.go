package main

import (
	"sync"

	"../auth"
	"../iotcoap"
	"../iothttp"
	"../iotmqtt"
)

var wgAuth sync.WaitGroup

func main() {
	wgAuth.Add(1)
	go auth.SvrListen(&wgAuth)

	//http
	//路由部分
	router := iothttp.RouterRegister()
	//静态资源
	//router.Static("/static", "./linuxdashboard/godashboard")
	//运行的端口
	go router.Run(":8080")

	//coap
	go iotcoap.StartCoapServer() //port 5683

	//mqtt
	go iotmqtt.StartMqttServer() //port 1883

	//raw socket

	wgAuth.Wait()

}
