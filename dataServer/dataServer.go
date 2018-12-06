package main

import (
	"linuxdashboard"
	"sync"

	"../auth"
)

var wgAuth sync.WaitGroup

func main() {
	wgAuth.Add(1)
	go auth.SvrListen(&wgAuth)
	wgAuth.Wait()

	//路由部分
	router := linuxdashboard.RouterRegister()
	//静态资源
	router.Static("/static", "./linuxdashboard/godashboard")
	//运行的端口
	router.Run(":8080")

}
