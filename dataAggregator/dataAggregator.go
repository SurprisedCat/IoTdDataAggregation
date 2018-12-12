package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	//"../auth"
	"../iotcoap"
	"../iothttp"
	"../iotmqtt"
	"../rawsocket"
	"../utils"
	proto "github.com/huin/mqtt"
)

var wgAuth sync.WaitGroup

func main() {

	serverAddr := []byte("127.0.0.1")
	protocolType := "mqtt"
	httpPort := []byte("8080")
	aggreLength := 10 //maximum 50
	var timeout int64 = 9000000000
	coapPort := []byte("5683")
	mqttPort := []byte("1883")
	ifaceName := "wlp61s0" //raw socket device
	backPort := coapPort

	// http
	// 路由部分
	router := iothttp.RouterRegister()
	//静态资源
	//router.Static("/static", "./linuxdashboard/godashboard")
	//运行的端口
	go router.Run(":18080")
	if protocolType == "http" {
		var httpwg sync.WaitGroup
		timeStart := time.Now().UnixNano()
		recPost := make([]map[string][]byte, aggreLength)
		for {
			if len(iothttp.ContentChan) > aggreLength {
				for i := 0; i < aggreLength; i++ {
					recPost[i] = <-iothttp.ContentChan
					recPost[i]["data"] = append(recPost[i]["data"], []byte("00000000000000000000000000000000")...)
				}
				dataJSON, err := json.Marshal(recPost)
				if err != nil {
					utils.CheckErr(err, "HTTP POST error")
				}
				fmt.Printf("%v\n", dataJSON)
				httpwg.Add(1)
				go iothttp.AggregatorSend(serverAddr, httpPort, dataJSON, &httpwg)
				timeStart = time.Now().UnixNano()
			}
			if time.Now().UnixNano()-timeStart > timeout { //纳秒为单位
				fmt.Printf("TimeOut：%d\n", len(iothttp.ContentChan))
				for i := 0; i < len(iothttp.ContentChan); i++ {
					recPost[i] = <-iothttp.ContentChan
					recPost[i]["data"] = append(recPost[i]["data"], []byte("00000000000000000000000000000000")...)
				}
				dataJSON, err := json.Marshal(recPost)
				if err != nil {
					utils.CheckErr(err, "HTTP POST error")
				}
				fmt.Printf("%v\n", dataJSON)
				httpwg.Add(1)
				go iothttp.AggregatorSend(serverAddr, httpPort, dataJSON, &httpwg)
				timeStart = time.Now().UnixNano()
			}
		}
	}
	//coap
	go iotcoap.StartCoapServer("15683") //port 15683
	if protocolType == "coap" {
		var coapwg sync.WaitGroup
		timeStart := time.Now().UnixNano()
		recPost := make([]map[string][]byte, aggreLength)
		for {
			if len(iotcoap.ContentChan) > aggreLength {
				for i := 0; i < aggreLength; i++ {
					recPost[i] = <-iotcoap.ContentChan
					recPost[i]["data"] = append(recPost[i]["data"], []byte("00000000000000000000000000000000")...)
				}
				fmt.Printf("%v\n", recPost)
				dataJSON, err := json.Marshal(recPost)
				if err != nil {
					utils.CheckErr(err, "COAP POST error")
				}
				fmt.Printf("%v\n", dataJSON)
				coapwg.Add(1)
				go iotcoap.AggregatorSend(serverAddr, coapPort, dataJSON, &coapwg)
				timeStart = time.Now().UnixNano()
			}
			if time.Now().UnixNano()-timeStart > timeout { //纳秒为单位
				fmt.Printf("COAP TimeOut：%d\n", len(iotcoap.ContentChan))
				for i := 0; i < len(iotcoap.ContentChan); i++ {
					recPost[i] = <-iotcoap.ContentChan
					recPost[i]["data"] = append(recPost[i]["data"], []byte("00000000000000000000000000000000")...)
				}
				dataJSON, err := json.Marshal(recPost)
				if err != nil {
					utils.CheckErr(err, "COAP POST error")
				}
				fmt.Printf("%v\n", dataJSON)
				coapwg.Add(1)
				go iotcoap.AggregatorSend(serverAddr, coapPort, dataJSON, &coapwg)
				timeStart = time.Now().UnixNano()
			}
		}
	}
	//mqtt
	go iotmqtt.StartMqttServer([]byte("11883")) //port 11883
	time.Sleep(time.Duration(2) * time.Second)
	go iotmqtt.AggregatorSubscriber([]byte("127.0.0.1"), []byte("11883"), string(utils.GetClientID("cx")))
	if protocolType == "mqtt" {
		var mqttwg sync.WaitGroup
		timeStart := time.Now().UnixNano()
		recPost := make([]map[string][]byte, aggreLength)
		for {

			var payload proto.Payload

			if len(iotmqtt.ContentChan) > aggreLength {
				for i := 0; i < aggreLength; i++ {
					recPost[i] = <-iotmqtt.ContentChan
					recPost[i]["data"] = append(recPost[i]["data"], []byte("00000000000000000000000000000000")...)
				}
				fmt.Printf("%v\n", recPost)
				dataJSON, err := json.Marshal(recPost)
				payload = proto.BytesPayload(dataJSON)
				if err != nil {
					utils.CheckErr(err, "MQTT aggregator error")
				}
				fmt.Printf("%v\n", dataJSON)
				mqttwg.Add(1)
				go iotmqtt.ClientPublisher(1, serverAddr, mqttPort, string(utils.GetClientID("cx")), &payload, 0, &mqttwg)
				timeStart = time.Now().UnixNano()
			}
			if time.Now().UnixNano()-timeStart > timeout { //纳秒为单位
				fmt.Printf("MQTT TimeOut：%d\n", len(iotcoap.ContentChan))
				for i := 0; i < len(iotmqtt.ContentChan); i++ {
					recPost[i] = <-iotmqtt.ContentChan
					recPost[i]["data"] = append(recPost[i]["data"], []byte("00000000000000000000000000000000")...)
				}
				dataJSON, err := json.Marshal(recPost)
				payload = proto.BytesPayload(dataJSON)
				if err != nil {
					utils.CheckErr(err, "COAP POST error")
				}
				fmt.Printf("%v\n", dataJSON)
				mqttwg.Add(1)
				go iotmqtt.ClientPublisher(1, serverAddr, mqttPort, string(utils.GetClientID("cx")), &payload, 0, &mqttwg)
				timeStart = time.Now().UnixNano()
			}
		}
	}

	go rawsocket.RecLinkLayer(ifaceName, 0x7676)
	//raw socket
	if protocolType == "7676" {
		var rswg sync.WaitGroup
		timeStart := time.Now().UnixNano()
		recPost := make([]map[string][]byte, aggreLength)
		for {
			if len(rawsocket.ContentChan) > aggreLength {
				for i := 0; i < aggreLength; i++ {
					recPost[i] = <-rawsocket.ContentChan
				}
				dataJSON, err := json.Marshal(recPost)
				if err != nil {
					utils.CheckErr(err, "rawsocket json decode error")
				}
				fmt.Printf("%v\n", recPost)
				if bytes.Compare(backPort, httpPort) == 0 {
					rswg.Add(1)
					go iothttp.AggregatorSend(serverAddr, backPort, dataJSON, &rswg)
				} else if bytes.Compare(backPort, coapPort) == 0 {
					rswg.Add(1)
					go iotcoap.AggregatorSend(serverAddr, backPort, dataJSON, &rswg)
				}

				timeStart = time.Now().UnixNano()
			}
			if time.Now().UnixNano()-timeStart > timeout { //纳秒为单位
				fmt.Printf("RawSocket TimeOut：%d\n", len(rawsocket.ContentChan))
				for i := 0; i < len(rawsocket.ContentChan); i++ {
					recPost[i] = <-rawsocket.ContentChan
				}
				dataJSON, err := json.Marshal(recPost)
				if err != nil {
					utils.CheckErr(err, "COAP POST error")
				}
				fmt.Printf("%v\n", dataJSON)
				if bytes.Compare(backPort, httpPort) == 0 {
					rswg.Add(1)
					go iothttp.AggregatorSend(serverAddr, backPort, dataJSON, &rswg)
				} else if bytes.Compare(backPort, coapPort) == 0 {
					rswg.Add(1)
					go iotcoap.AggregatorSend(serverAddr, backPort, dataJSON, &rswg)
				}
				timeStart = time.Now().UnixNano()
			}
		}
	}

	//wgAuth.Wait()
}
