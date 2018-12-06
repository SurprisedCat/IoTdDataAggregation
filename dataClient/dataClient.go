package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"../auth"
	"../simssl"
	"../utils"
)

func main() {
	/**********************Authentication****************/
	contents, err := ioutil.ReadFile("key.txt") //read the key.txt
	authTime := 0
	for err != nil || len(contents) == 0 {
		if !auth.ClientAuth([]byte("127.0.0.1")) { //connect server for autentication
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
		auth.ClientAuth([]byte("127.0.0.1"))
	}
	/*********************Check Expiration************/

	/*********************data generation************/
	clientID := utils.GetClientID()
	origData := append([]byte("I am "), clientID...)
	encryptedData, err := simssl.AesEncrypt(origData, encryptKey)
	if err != nil {
		log.Fatal("simssl.AesEncrypt:", err)
	}
	dataForSend := map[string][]byte{"ID": clientID, "data": encryptedData} //data is in the form of map[string][]byte
	dataJSON, err := json.Marshal(dataForSend)
	/*********************data generation************/

	/********************send with http*************/

	/********************send with http*************/

	/********************send with mqtt*************/

	/********************send with mqtt*************/

	/********************send with coap*************/

	/********************send with coap*************/

	/********************send with socketraw*************/

	/********************send with socketraw*************/

	/*****************data decoding****************/
	fmt.Println(string(dataJSON))
	dec := map[string][]byte{}                   //data is in the form of map[string][]byte
	err = json.Unmarshal([]byte(dataJSON), &dec) //json decode
	if err != nil {
		log.Fatal("JSON error:", err)
	}
	auth.GetKeyClient()
	decrypted, _ := simssl.AesDecrypt(dec["data"], encryptKey) //aes decoding
	fmt.Println(string(dec["ID"]))
	fmt.Println(string(decrypted))

}

/*****************data decoding****************/
