#! /bin/bash
echo "TEST sudo go [parameters] "
TESTRES=`sudo go env | grep GOPATH`
if test -z "$TESTRES" 
then
    echo "Need sudo to execute the raw socket programs"
    echo "Please add go path to \"Defaults secure_path\" of /etc/sudoers"
    exit
fi
echo "go get github.com/gomodule/redigo/redis"
go get github.com/gomodule/redigo/redis
echo "go get github.com/gin-gonic/gin"
go get github.com/gin-gonic/gin
echo "go get github.com/dustin/go-coap" 
go get github.com/dustin/go-coap
echo "go get github.com/jeffallen/mqtt"
go get github.com/jeffallen/mqtt
echo "go get github.com/huin/mqtt"
go get github.com/huin/mqtt
echo "go get github.com/golang/net/bpf"
go get github.com/golang/net/bpf
echo "Create golang/net/bpf later. Pass this warning."
echo "go get github.com/golang/sys/unix" 
go get github.com/golang/sys/unix 
echo "Create golang/sys/unix later. Pass this warning."
echo "Create the golang.org/x/"
GOPATH=`go env GOPATH`
mkdir $GOPATH/src/golang.org/x -p
cp $GOPATH/src/github.com/golang/net $GOPATH/src/golang.org/x/ -rf
cp $GOPATH/src/github.com/golang/sys $GOPATH/src/golang.org/x/ -rf
echo "OK. golang/net/bpf and golang/sys/unix have been created."
echo "go get github.com/mdlayher/raw"
go get github.com/mdlayher/raw
echo "go get github.com/mdlayher/ethernet"
go get github.com/mdlayher/ethernet
echo "Finish!!"
