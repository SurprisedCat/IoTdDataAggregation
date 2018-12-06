package main

import (
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"../auth"
	"../iothttp"
	"../simssl"
	"../utils"
)

func main() {
	//common parameters
	serverAddr := []byte("127.0.0.1")
	origData := []byte("DATA")
	protocolType := "http"
	httpPort := []byte("8080")

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
		for i := 0; i < 100; i++ {
			httpwg.Add(1)
			go iothttp.ClientSend(serverAddr, httpPort, dataJSON, &httpwg)
		}
		httpwg.Wait()
	}

	/********************send with http*************/

	/********************send with mqtt*************/

	/********************send with mqtt*************/

	/********************send with coap*************/

	/********************send with coap*************/

	/********************send with socketraw*************/

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
