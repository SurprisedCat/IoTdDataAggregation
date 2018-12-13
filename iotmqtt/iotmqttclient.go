package iotmqtt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"../auth"
	"../simssl"
	proto "github.com/huin/mqtt"
	"github.com/jeffallen/mqtt"
)

//ContentChan for subscribers
var ContentChan = make(chan map[string][]byte, 50)

//ClientPublisher mqtt publisher
func ClientPublisher(i int, serverAddr, mqttPort []byte, topic string, payload *proto.Payload, pace int, mqttwg *sync.WaitGroup) {
	defer mqttwg.Done()
	log.Print("starting client ", i)
	conn, err := net.Dial("tcp", string(serverAddr)+":"+string(mqttPort))
	if err != nil {
		log.Fatal("dial: ", err)
	}
	cc := mqtt.NewClientConn(conn)
	cc.Dump = false

	if err := cc.Connect("", ""); err != nil {
		log.Fatalf("connect: %v\n", err)
		os.Exit(1)
	}

	//half := int32(pace / 2)

	for {
		cc.Publish(&proto.Publish{
			Header:    proto.Header{},
			TopicName: topic,
			Payload:   *payload,
		})
		//	sltime := rand.Int31n(half) - (half / 2) + int32(pace)
		time.Sleep(time.Duration(pace) * time.Microsecond)
	}
}

//ServerSubscriberSingle the only one subscriber
func ServerSubscriberSingle(serverAddr, mqttPort []byte, topic string) {
	conn, err := net.Dial("tcp", string(serverAddr)+":"+string(mqttPort))
	if err != nil {
		log.Fatal("dial: ", err)
	}
	cc := mqtt.NewClientConn(conn)
	cc.Dump = false

	if err := cc.Connect("", ""); err != nil {
		log.Fatalf("connect: %v\n", err)
		os.Exit(1)
	}
	tq := make([]proto.TopicQos, 1)
	//for i := 0; i < flag.NArg(); i++ {
	tq[0].Topic = topic
	tq[0].Qos = proto.QosAtMostOnce
	//}
	cc.Subscribe(tq)
	for m := range cc.Incoming {

		//json_decode
		recPublish := map[string][]byte{}
		bufPayload := bytes.NewBuffer(make([]byte, 0))
		m.Payload.WritePayload(bufPayload)
		err := json.Unmarshal(bufPayload.Bytes(), &recPublish)
		if err != nil {
			fmt.Println("Json decode error")
			return
		}
		clientID := recPublish["ID"]
		clientData := recPublish["data"]
		eKey, vali := auth.GetValidationKeyServer([]byte(clientID))
		//fmt.Println(clientData)
		if vali == false {
			fmt.Println("AES key decode error")
			return
		}
		originData, err := simssl.AesDecrypt([]byte(clientData), eKey)
		fmt.Println(string(originData))
		if err != nil {
			fmt.Println("AES decode error")
			return
		}
		fmt.Print(m.TopicName, "\t")
		fmt.Printf("ID:%s\tdata:%s", clientID, originData)
		fmt.Println("\tr: ", m.Header.Retain)
	}
}

//ServerSubscriberCluster the only one subscriber
func ServerSubscriberCluster(serverAddr, mqttPort []byte, topic string) {
	conn, err := net.Dial("tcp", string(serverAddr)+":"+string(mqttPort))
	if err != nil {
		log.Fatal("dial: ", err)
	}
	cc := mqtt.NewClientConn(conn)
	cc.Dump = false

	if err := cc.Connect("", ""); err != nil {
		log.Fatalf("connect: %v\n", err)
		os.Exit(1)
	}
	tq := make([]proto.TopicQos, 1)
	//for i := 0; i < flag.NArg(); i++ {
	tq[0].Topic = topic
	tq[0].Qos = proto.QosAtMostOnce
	//}
	cc.Subscribe(tq)
	for m := range cc.Incoming {

		//json_decode
		recPublish := []map[string][]byte{}
		bufPayload := bytes.NewBuffer(make([]byte, 0))
		m.Payload.WritePayload(bufPayload)
		err := json.Unmarshal(bufPayload.Bytes(), &recPublish)
		if err != nil {
			fmt.Println("Json decode error")
			return
		}
		clientID := recPublish[0]["ID"]
		clientData := recPublish[0]["data"]
		eKey, vali := auth.GetValidationKeyServer([]byte(clientID))
		//fmt.Println(clientData)
		if vali == false {
			fmt.Println("AES key decode error")
			return
		}
		originData, err := simssl.AesDecrypt([]byte(clientData), eKey)
		fmt.Println(string(originData))
		if err != nil {
			fmt.Println("AES decode error")
			return
		}
		fmt.Print(m.TopicName, "\t")
		fmt.Printf("ID:%s\tdata:%s", clientID, originData)
		fmt.Println("\tr: ", m.Header.Retain)
	}
}

//AggregatorSubscriber the only one subscriber in the aggregator
func AggregatorSubscriber(serverAddr, mqttPort []byte, topic string) {
	conn, err := net.Dial("tcp", string(serverAddr)+":"+string(mqttPort))
	if err != nil {
		log.Fatal("dial: ", err)
	}
	cc := mqtt.NewClientConn(conn)
	cc.Dump = false

	if err := cc.Connect("", ""); err != nil {
		log.Fatalf("connect: %v\n", err)
		os.Exit(1)
	}
	tq := make([]proto.TopicQos, 1)
	//for i := 0; i < flag.NArg(); i++ {
	tq[0].Topic = topic
	tq[0].Qos = proto.QosAtMostOnce
	//}
	cc.Subscribe(tq)
	for m := range cc.Incoming {

		//json_decode
		recPublish := map[string][]byte{}
		bufPayload := bytes.NewBuffer(make([]byte, 0))
		m.Payload.WritePayload(bufPayload)
		err := json.Unmarshal(bufPayload.Bytes(), &recPublish)
		if err != nil {
			fmt.Println("Json decode error")
			return
		}
		clientID := recPublish["ID"]
		clientData := recPublish["data"]
		ContentChan <- recPublish
		fmt.Print(m.TopicName, "\t")
		fmt.Printf("ID:%s\tdata:%s", clientID, clientData)
		fmt.Println("\tr: ", m.Header.Retain)
	}
}
