package main

import (
	"sync"

	"../auth"
)

var wgAuth sync.WaitGroup

func main() {
	wgAuth.Add(1)
	go auth.SvrListen(&wgAuth)
	wgAuth.Wait()

}
