package iotcoap

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"../auth"
	"../simssl"
	"github.com/dustin/go-coap"
)

//ContentChan content type is json decoded as map[string][]byte
var ContentChan = make(chan map[string][]byte, 50)

/*
* Data are encoded into json format
 */

//ProcSingle process the single data packet
func ProcSingle(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
	//log.Printf("Got message in ProcSingle: path=%q: %#v from %v", m.Path(), m, a)
	log.Printf("Got message in ProcSingle: path=%q: from %v", m.Path(), a)
	errRes := &coap.Message{
		Type:      coap.Reset,
		Code:      coap.BadRequest,
		MessageID: m.MessageID,
		Token:     m.Token,
	}
	//json_decode
	recPost := map[string][]byte{}
	err := json.Unmarshal(m.Payload, &recPost)
	if err != nil {
		fmt.Println("Json decode error")
		return errRes
	}
	clientID := recPost["ID"]
	clientData := recPost["data"]
	eKey, vali := auth.GetValidationKeyServer([]byte(clientID))
	//fmt.Println(clientData)
	if vali == false {
		fmt.Println("AES key decode error")
		return errRes
	}
	originData, err := simssl.AesDecrypt([]byte(clientData), eKey)
	fmt.Println(string(originData))
	if err != nil {
		fmt.Println("AES decode error")
		return errRes
	}

	if m.IsConfirmable() {
		repJSON, _ := json.Marshal(map[string][]byte{"status": []byte("OK")})
		res := &coap.Message{
			Type:      coap.Acknowledgement,
			Code:      coap.Changed,
			MessageID: m.MessageID,
			Token:     m.Token,
			Payload:   repJSON,
		}
		res.SetOption(coap.ContentFormat, coap.AppJSON)

		//log.Printf("Transmitting from ProcSingle %#v", res)
		return res
	}
	return nil
}

//ProcCluster process the cluster upload data
func ProcCluster(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
	fmt.Println("Receive a clustered packet")
	errRes := &coap.Message{
		Type:      coap.Reset,
		Code:      coap.BadRequest,
		MessageID: m.MessageID,
		Token:     m.Token,
	}
	//json_decode
	recPost := []map[string][]byte{}
	err := json.Unmarshal(m.Payload, &recPost)
	if err != nil {
		fmt.Println(m.Payload)
		fmt.Println("Cluster Json decode error")
		return errRes
	}
	clientID := recPost[0]["ID"]
	eKey, vali := auth.GetValidationKeyServer([]byte(clientID))
	//fmt.Println(clientData)
	if vali == false {
		fmt.Println("AES key decode error")
		return errRes
	}
	for i := 0; i < len(recPost); i++ {
		clientData := recPost[i]["data"]
		originData, err := simssl.AesDecrypt([]byte(clientData), eKey)
		fmt.Println(string(recPost[i]["ID"]), string(originData))
		if err != nil {
			return errRes
		}
	}

	if m.IsConfirmable() {
		repJSON, _ := json.Marshal(map[string][]byte{"status": []byte("OK")})
		res := &coap.Message{
			Type:      coap.Acknowledgement,
			Code:      coap.Changed,
			MessageID: m.MessageID,
			Token:     m.Token,
			Payload:   repJSON,
		}
		res.SetOption(coap.ContentFormat, coap.AppJSON)

		//log.Printf("Transmitting from ProcSingle %#v", res)
		return res
	}
	return nil
}

//ProcAggre aggregator process
func ProcAggre(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
	//log.Printf("Got message in ProcSingle: path=%q: %#v from %v", m.Path(), m, a)
	log.Printf("Got message in ProcSingle: path=%q: from %v", m.Path(), a)
	errRes := &coap.Message{
		Type:      coap.Reset,
		Code:      coap.BadRequest,
		MessageID: m.MessageID,
		Token:     m.Token,
	}
	//json_decode
	recPost := map[string][]byte{}
	err := json.Unmarshal(m.Payload, &recPost)
	if err != nil {
		fmt.Println("Json decode error")
		return errRes
	}
	fmt.Println(recPost)
	//put it into channel
	ContentChan <- recPost
	if m.IsConfirmable() {
		repJSON, _ := json.Marshal(map[string][]byte{"status": []byte("OK")})
		res := &coap.Message{
			Type:      coap.Acknowledgement,
			Code:      coap.Changed,
			MessageID: m.MessageID,
			Token:     m.Token,
			Payload:   repJSON,
		}
		res.SetOption(coap.ContentFormat, coap.AppJSON)

		//log.Printf("Transmitting from ProcSingle %#v", res)
		return res
	}
	return nil
}

//StartCoapServer will listen on port 5683 aggregator 15683
func StartCoapServer(coapPort string) {
	mux := coap.NewServeMux()
	mux.Handle("/v1/upload/single", coap.FuncHandler(ProcSingle))
	mux.Handle("/v1/upload/cluster", coap.FuncHandler(ProcCluster))
	mux.Handle("/v1/upload/aggre", coap.FuncHandler(ProcAggre))

	log.Fatal(coap.ListenAndServe("udp", "0.0.0.0:"+coapPort, mux))
}
