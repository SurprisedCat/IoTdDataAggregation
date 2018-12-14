package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	//"../auth"
	"../config"
	"../iotcoap"
	"../iothttp"
	"../iotmqtt"
	"../rawsocket"
	"../utils"
	proto "github.com/huin/mqtt"
)

var wgAuth sync.WaitGroup

func main() {

	serverAddr := config.ServerAddr
	protocolType := config.ProtocolType
	httpPort := []byte("8080")
	aggreLength := config.AggreLength //the total bytes for send
	var timeout int64 = 6000000000    //nanosecond
	coapPort := []byte("5683")
	mqttPort := []byte("1883")
	ifaceName := config.AggreIfaceName //raw socket device
	backPort := config.BackPort

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
		for {
			lengthPacketForSend := len(iothttp.ContentChan) // the number of packets in the channel
			recPost := make([]map[string][]byte, lengthPacketForSend)
			if lengthPacketForSend*config.TotalPayloadJSONLength > aggreLength {
				for i := 0; i < lengthPacketForSend; i++ {
					recPost[i] = <-iothttp.ContentChan
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
			if time.Now().UnixNano()-timeStart > timeout && lengthPacketForSend > 0 { //纳秒为单位
				fmt.Printf("TimeOut：%d\n", lengthPacketForSend)
				for i := 0; i < lengthPacketForSend; i++ {
					recPost[i] = <-iothttp.ContentChan
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
		for {
			lengthPacketForSend := len(iotcoap.ContentChan) // the number of packets in the channel
			recPost := make([]map[string][]byte, lengthPacketForSend)
			// the totoal bytes in the channel is larger than threhold. 1200 is the maximum payload of COAP in this condition
			if lengthPacketForSend*config.TotalPayloadJSONLength > utils.Min(aggreLength, 1100) { // the totoal bytes in the channel is larger than threhold
				for i := 0; i < lengthPacketForSend; i++ {
					recPost[i] = <-iotcoap.ContentChan
					fmt.Printf("COAP aggregator %d \n", i)
				}
				fmt.Printf("%v\n", recPost)
				dataJSON, err := json.Marshal(recPost)
				if err != nil {
					utils.CheckErr(err, "COAP POST error")
				}
				//fmt.Printf("%v\n", dataJSON)
				coapwg.Add(1)
				go iotcoap.AggregatorSend(serverAddr, coapPort, dataJSON, &coapwg)
				timeStart = time.Now().UnixNano()
			}
			if time.Now().UnixNano()-timeStart > timeout && lengthPacketForSend > 0 { //纳秒为单位
				fmt.Printf("COAP TimeOut：%d\n", lengthPacketForSend)
				for i := 0; i < lengthPacketForSend; i++ {
					recPost[i] = <-iotcoap.ContentChan
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
		for {

			var payload proto.Payload
			lengthPacketForSend := len(iotmqtt.ContentChan) // the number of packets in the channel
			recPost := make([]map[string][]byte, lengthPacketForSend)
			if lengthPacketForSend*config.TotalPayloadJSONLength > aggreLength { // the totoal bytes in the channel is larger than threhold
				for i := 0; i < lengthPacketForSend; i++ {
					recPost[i] = <-iotmqtt.ContentChan
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
			if time.Now().UnixNano()-timeStart > timeout && lengthPacketForSend > 0 { //纳秒为单位
				fmt.Printf("MQTT TimeOut：%d\n", lengthPacketForSend)
				for i := 0; i < lengthPacketForSend; i++ {
					recPost[i] = <-iotmqtt.ContentChan
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
		for {
			lengthPacketForSend := len(rawsocket.ContentChan) // the number of packets in the channel
			recPost := make([]map[string][]byte, lengthPacketForSend)
			// the totoal bytes in the channel is larger than threhold. 1200 is the maximum payload of COAP in this condition
			if bytes.Compare(backPort, coapPort) == 0 {
				aggreLength = utils.Min(aggreLength, 1100)
			}
			if lengthPacketForSend*config.TotalPayloadJSONLength > aggreLength { // the totoal bytes in the channel is larger than threhold
				for i := 0; i < lengthPacketForSend; i++ {
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
			if time.Now().UnixNano()-timeStart > timeout && lengthPacketForSend > 0 { //纳秒为单位
				fmt.Printf("RawSocket TimeOut：%d\n", lengthPacketForSend)
				for i := 0; i < lengthPacketForSend; i++ {
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
