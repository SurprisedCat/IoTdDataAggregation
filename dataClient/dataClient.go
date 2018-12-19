package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	proto "github.com/huin/mqtt"

	"../auth"
	"../config"
	"../iotcoap"
	"../iothttp"
	"../iotmqtt"
	"../rawsocket"
	"../simssl"
	"../utils"
)

func main() {
	//common parameters
	selfID := config.SelfID
	serverAddr := config.ServerAddr
	aggregatorAddr := config.AggregatorAddr
	origData := config.PayloadData
	protocolType := config.ProtocolType
	totalReq := config.TotalReq
	cluster := config.Cluster
	timeGap := config.TimeGap //微秒为单位

	var httpPort, coapPort, mqttPort []byte
	if cluster == true {
		httpPort = []byte("18080")
		coapPort = []byte("15683")
		mqttPort = []byte("11883")
	} else {
		httpPort = []byte("8080")
		coapPort = []byte("5683")
		mqttPort = []byte("1883")
	}
	ifaceName := config.ClientIfaceName
	dstMac := config.DstMac

	if config.Debug == true {
		del := os.Remove("./key.txt")
		if del != nil {
			fmt.Println(del)
			fmt.Println("Start to authenticate and create key.txt")
		}
	}

	/**********************Authentication****************/
	contents, err := ioutil.ReadFile("key.txt") //read the key.txt
	authTime := 0
	for err != nil || len(contents) == 0 {
		if !auth.ClientAuth(serverAddr, selfID) { //connect server for autentication
			authTime++
			log.Printf("Authentication fails for %d", authTime)
			time.Sleep(time.Duration(3*authTime) * time.Second)
		}
		if authTime > 5 { // for 5 times, th client will stop
			log.Fatal("Authenticaion fails")
		}
		contents, err = ioutil.ReadFile("key.txt")
	}
	/**********************Authentication****************/

	/*********************Check Expiration************/
	encryptKey := contents[:16]
	expirationTime, numbers := binary.Varint(contents[16:24])
	if expirationTime < time.Now().Unix() || numbers <= 0 {
		auth.ClientAuth(serverAddr, selfID)
	}
	/*********************Check Expiration************/

	/*********************data generation************/
	clientID := utils.GetClientID(selfID)
	fmt.Println(clientID)

	encryptedData, err := simssl.AesEncrypt(origData, encryptKey)
	if err != nil {
		log.Fatal("simssl.AesEncrypt:", err)
	}
	dataForSend := map[string][]byte{"ID": clientID, "data": encryptedData} //data is in the form of map[string][]byte
	dataJSON, err := json.Marshal(dataForSend)
	/*********************data generation************/

	/********************send with http*************/
	if protocolType == "http" {
		var httpwg sync.WaitGroup
		before := time.Now().UnixNano()

		for i := 0; i < totalReq; i++ {
			if cluster == true {
				httpwg.Add(1)
				go iothttp.ClientSend(aggregatorAddr, httpPort, dataJSON, &httpwg)
			} else {
				httpwg.Add(1)
				go iothttp.ClientSend(serverAddr, httpPort, dataJSON, &httpwg)
			}
			//这里可以控制发包频率
			time.Sleep(time.Duration(timeGap) * time.Microsecond)
		}
		httpwg.Wait()
		fmt.Println(time.Now().UnixNano() - before)
	}
	/********************send with http*************/

	/********************send with coap*************/
	if protocolType == "coap" {
		var coapwg sync.WaitGroup
		before := time.Now().UnixNano()

		for i := 0; i < totalReq; i++ {
			if cluster == true {
				coapwg.Add(1)
				go iotcoap.ClientSend(aggregatorAddr, coapPort, dataJSON, &coapwg)
			} else {
				coapwg.Add(1)
				go iotcoap.ClientSend(serverAddr, coapPort, dataJSON, &coapwg)
			}
			//这里可以控制发包频率
			time.Sleep(time.Duration(timeGap) * time.Microsecond)
		}
		coapwg.Wait()
		fmt.Println(time.Now().UnixNano() - before)
	}
	/********************send with coap*************/

	/********************send with mqtt*************/
	if protocolType == "mqtt" {
		var conns = flag.Int("conns", 10, "how many conns (0 means infinite)")
		//var host = flag.String("host", string(serverAddr)+":1883", "hostname of broker")
		//var user = flag.String("user", "", "username")
		//var pass = flag.String("pass", "", "password")
		//var dump = flag.Bool("dump", false, "dump messages?")
		var wait = flag.Int("wait", 10, "ms to wait between client connects")
		var pace = flag.Int("pace", 1000000, "sleep time")

		var payload proto.Payload
		var topic string

		flag.Parse()

		if flag.NArg() != 2 {
			topic = config.MqttTopic
			payload = proto.BytesPayload(dataJSON)
		} else {
			topic = flag.Arg(0)
			payload = proto.BytesPayload([]byte(flag.Arg(1)))
		}

		var mqttwg sync.WaitGroup
		i := 1
		for ; i != *conns; i++ {
			if cluster == true {
				mqttwg.Add(1)
				go iotmqtt.ClientPublisher(i, aggregatorAddr, mqttPort, topic, &payload, *pace, &mqttwg)
			} else {
				mqttwg.Add(1)
				go iotmqtt.ClientPublisher(i, serverAddr, mqttPort, topic, &payload, *pace, &mqttwg)
			}

			time.Sleep(time.Duration(*wait) * time.Millisecond)

		}
		mqttwg.Wait()
		// sleep forever
		//<-make(chan struct{})
	}
	/********************send with mqtt*************/

	/********************send with socketraw*************/
	if protocolType == "7676" {

		var llwg sync.WaitGroup
		for i := 0; i < totalReq; i++ {
			llwg.Add(1)
			go rawsocket.SendLinkLayer(ifaceName, dstMac, 0x7676, dataJSON, &llwg)
			time.Sleep(time.Duration(timeGap) * time.Microsecond)
		}
		llwg.Wait()
	}

	/********************send with socketraw*************/

	/*****************data decoding****************/
	// fmt.Println(string(dataJSON))
	// dec := map[string][]byte{}                   //data is in the form of map[string][]byte
	// err = json.Unmarshal([]byte(dataJSON), &dec) //json decode
	// if err != nil {
	// 	log.Fatal("JSON error:", err)
	// }
	// auth.GetKeyClient()
	// decrypted, _ := simssl.AesDecrypt(dec["data"], encryptKey) //aes decoding
	// fmt.Println(string(dec["ID"]))
	// fmt.Println(string(decrypted))

}

/*****************data decoding****************/
