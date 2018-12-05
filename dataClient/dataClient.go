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

type dec struct {
	page []byte `json:"page"`
}

func main() {
	contents, err := ioutil.ReadFile("key.txt")
	if err != nil || len(contents) == 0 {
		auth.ClientAuth([]byte("127.0.0.1"))
	}
	encryptKey := contents[:16]
	expirationTime, numbers := binary.Varint(contents[16:24])
	if expirationTime < time.Now().Unix() || numbers <= 0 {
		auth.ClientAuth([]byte("127.0.0.1"))
	}
	origData := []byte("I am a phd")
	encryptedData, err := simssl.AesEncrypt(origData, encryptKey)
	if err != nil {
		log.Fatal("simssl.AesEncrypt:", err)
	}

	clientID := utils.GetClientID()
	dataForSend := "{\"" + clientID + "\":\"" + string(encryptedData) + "\"}"
	dataForSend = `{"page": "1"}`
	fmt.Println(dataForSend)
	decTest := dec{}
	err = json.Unmarshal([]byte(dataForSend), &decTest)
	if err != nil {
		log.Fatal("JSON error:", err)
	}
	fmt.Println(decTest)

}
