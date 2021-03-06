package iotmqtt

import (
	"flag"
	"log"
	"net"

	"github.com/jeffallen/mqtt"
)

//StartMqttServer start mqtt server
func StartMqttServer(mqttPort []byte) {
	var addr = flag.String("addr", "0.0.0.0:"+string(mqttPort), "listen address of broker")

	flag.Parse()

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Print("listen: ", err)
		return
	}
	svr := mqtt.NewServer(l)
	svr.Start()
	<-svr.Done
}
