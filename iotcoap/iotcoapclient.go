package iotcoap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/dustin/go-coap"
)

/*
ClientSend send a specific packet to ip:port
*/
func ClientSend(serverAddr, coapPort, dataJSON []byte, coapwg *sync.WaitGroup) {

	rand.Seed(time.Now().Unix())
	defer coapwg.Done()

	req := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.POST,
		MessageID: uint16(rand.Uint32()),
		Payload:   dataJSON,
	}

	var path string
	if bytes.Compare(coapPort, []byte("15683")) == 0 {
		path = "/v1/upload/aggre"
	} else {
		path = "/v1/upload/single"
	}

	req.SetOption(coap.ContentFormat, coap.AppJSON)
	req.SetOption(coap.Accept, coap.AppJSON)

	req.SetOption(coap.MaxAge, 60)
	req.SetPathString(path)

	c, err := coap.Dial("udp", string(serverAddr)+":"+string(coapPort))
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}

	rv, err := c.Send(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	if rv != nil {
		coapResp := map[string][]byte{}
		json.Unmarshal(rv.Payload, &coapResp)
		if bytes.Compare(coapResp["status"], []byte("OK")) == 0 {
			fmt.Printf("coap response : %s\n", coapResp)
		} else {
			fmt.Printf("coap response : %s\n", coapResp)
		}
	}
}

//AggregatorSend send packets to server
func AggregatorSend(serverAddr, coapPort, dataJSON []byte, coapwg *sync.WaitGroup) {
	rand.Seed(time.Now().Unix())
	defer coapwg.Done()

	req := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.POST,
		MessageID: uint16(rand.Uint32()),
		Payload:   dataJSON,
	}

	path := "/v1/upload/cluster"

	req.SetOption(coap.ContentFormat, coap.AppJSON)
	req.SetOption(coap.Accept, coap.AppJSON)

	req.SetOption(coap.MaxAge, 60)
	req.SetPathString(path)

	c, err := coap.Dial("udp", string(serverAddr)+":"+string(coapPort))
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}

	rv, err := c.Send(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	if rv != nil {
		coapResp := map[string][]byte{}
		json.Unmarshal(rv.Payload, &coapResp)
		if bytes.Compare(coapResp["status"], []byte("OK")) == 0 {
			fmt.Printf("coap response : %s\n", coapResp)
		} else {
			fmt.Printf("coap response : %s\n", coapResp)
		}
	}
}
