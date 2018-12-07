package iotmqtt

import (
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	proto "github.com/huin/mqtt"
	"github.com/jeffallen/mqtt"
)

//Client mqtt
func Client(serverAddr, dataJSON []byte, mqttwg *sync.WaitGroup, i int) {
	log.Print("starting client ", i)
	conn, err := net.Dial("tcp", string(serverAddr))
	if err != nil {
		log.Fatal("dial: ", err)
	}
	cc := mqtt.NewClientConn(conn)
	cc.Dump = false

	if err := cc.Connect("", ""); err != nil {
		log.Fatalf("connect: %v\n", err)
		os.Exit(1)
	}

	half := int32(*pace / 2)

	for {
		cc.Publish(&proto.Publish{
			Header:    proto.Header{},
			TopicName: topic,
			Payload:   payload,
		})
		sltime := rand.Int31n(half) - (half / 2) + int32(*pace)
		time.Sleep(time.Duration(sltime) * time.Second)
	}
}
