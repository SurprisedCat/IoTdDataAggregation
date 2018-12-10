#! /bin/bash
go get github.com/gomodule/redigo/redis
go get github.com/gin-gonic/gin
go get github.com/dustin/go-coap
go get github.com/jeffallen/mqtt
go get github.com/huin/mqtt
go get github.com/golang/net/bpf
go get github.com/golang/sys/unix 
GOPATH=`go env GOPATH`
mkdir $GOPATH/src/golang.org/x -p
cp $GOPATH/src/github.com/golang/net $GOPATH/src/golang.org/x/ -rf
cp $GOPATH/src/github.com/golang/sys $GOPATH/src/golang.org/x/ -rf
go get github.com/mdlayher/raw
go get github.com/mdlayher/ethernet