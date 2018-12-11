package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"sync"
	"time"

	proto "github.com/huin/mqtt"

	"../auth"
	"../iotcoap"
	"../iothttp"
	"../iotmqtt"
	"../rawsocket"
	"../simssl"
	"../utils"
)

func main() {
	//common parameters
	serverAddr := []byte("127.0.0.1")
	aggragatorAddr := []byte("127.0.0.1")
	origData := []byte("DATA")
	protocolType := "http"
	totalReq := 15000
	cluster := false
	httpPort := []byte("8080")
	coapPort := []byte("5683")
	mqttPort := []byte("1883")
	ifaceName := "wlp4s0"
	dstMac := net.HardwareAddr{0x34, 0xe6, 0xad, 0x09, 0xc6, 0x3f}

	/**********************Authentication****************/
	contents, err := ioutil.ReadFile("key.txt") //read the key.txt
	authTime := 0
	for err != nil || len(contents) == 0 {
		if !auth.ClientAuth(serverAddr) { //connect server for autentication
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
		auth.ClientAuth(serverAddr)
	}
	/*********************Check Expiration************/

	/*********************data generation************/
	clientID := utils.GetClientID()
	origData = append([]byte("I am "), clientID...)
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
		if cluster == true {
			for i := 0; i < totalReq; i++ {
				httpwg.Add(1)
				go iothttp.ClientSend(aggragatorAddr, httpPort, dataJSON, &httpwg)
				//这里可以控制发包频率
				time.Sleep(time.Duration(time.Microsecond * 300))
			}
		} else {
			for i := 0; i < totalReq; i++ {
				httpwg.Add(1)
				go iothttp.ClientSend(serverAddr, httpPort, dataJSON, &httpwg)
				//这里可以控制发包频率
				time.Sleep(time.Duration(time.Microsecond * 300))
			}
		}
		httpwg.Wait()
		fmt.Println(time.Now().UnixNano() - before)
	}
	/********************send with http*************/

	/********************send with coap*************/
	if protocolType == "coap" {
		var coapwg sync.WaitGroup
		for i := 0; i < 10; i++ {
			coapwg.Add(1)
			time.Sleep(time.Duration(100 * time.Microsecond))
			go iotcoap.ClientSend(serverAddr, coapPort, dataJSON, &coapwg)
		}
		coapwg.Wait()
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
		var pace = flag.Int("pace", 60, "send a message on average once every pace seconds")

		var payload proto.Payload
		var topic string

		flag.Parse()

		if flag.NArg() != 2 {
			topic = "clientSingle"
			payload = proto.BytesPayload(dataJSON)
		} else {
			topic = flag.Arg(0)
			payload = proto.BytesPayload([]byte(flag.Arg(1)))
		}

		var mqttwg sync.WaitGroup
		i := 1
		for ; i != *conns; i++ {
			mqttwg.Add(1)
			go iotmqtt.ClientPublisher(i, serverAddr, mqttPort, topic, &payload, pace, &mqttwg)

			time.Sleep(time.Duration(*wait) * time.Millisecond)
		}
		mqttwg.Add(1)
		go iotmqtt.ClientSubscriber(serverAddr, mqttPort, topic, &mqttwg)
		mqttwg.Wait()
		// sleep forever
		//<-make(chan struct{})
	}
	/********************send with mqtt*************/

	/********************send with socketraw*************/
	if protocolType == "7676" {

		var llwg sync.WaitGroup
		for i := 0; i < 10; i++ {
			llwg.Add(1)
			go rawsocket.SendLinkLayer(ifaceName, dstMac, 0x7676, dataJSON, &llwg)
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
