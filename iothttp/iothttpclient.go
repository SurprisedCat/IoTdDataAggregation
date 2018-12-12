package iothttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"../utils"
)

/*
ClientSend send a specific packet to ip:port
*/
func ClientSend(serverAddr, httpPort, dataJSON []byte, httpwg *sync.WaitGroup) {
	defer httpwg.Done()
	var err error
	var resp *http.Response
	req := bytes.NewBuffer(dataJSON)

	if bytes.Compare(httpPort, []byte("18080")) == 0 {
		resp, err = http.Post("http://"+string(serverAddr)+":"+string(httpPort)+"/v1/upload/aggre", "application/json;charset=utf-8", req)

	} else {
		resp, err = http.Post("http://"+string(serverAddr)+":"+string(httpPort)+"/v1/upload/single", "application/json;charset=utf-8", req)

	}
	if err != nil {
		utils.CheckErr(err, "HTTP POST error")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.CheckErr(err, "HTTP POST error")
	}
	resp.Body.Close()
	httpResp := map[string][]byte{}
	json.Unmarshal(body, &httpResp)
	if bytes.Compare(httpResp["status"], []byte("OK")) == 0 {
		fmt.Printf("http response : %s\n", httpResp)
	} else {
		fmt.Printf("http response : %s\n", httpResp)
	}
}

//AggregatorSend send packets to server
func AggregatorSend(serverAddr, httpPort, dataJSON []byte, httpwg *sync.WaitGroup) {
	defer httpwg.Done()

	req := bytes.NewBuffer(dataJSON)
	resp, err := http.Post("http://"+string(serverAddr)+":"+string(httpPort)+"/v1/upload/cluster", "application/json;charset=utf-8", req)
	if err != nil {
		utils.CheckErr(err, "HTTP POST error")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.CheckErr(err, "HTTP POST error")
	}
	resp.Body.Close()
	httpResp := map[string][]byte{}
	json.Unmarshal(body, &httpResp)
	if bytes.Compare(httpResp["status"], []byte("OK")) == 0 {
		fmt.Printf("http response : %s\n", httpResp)
	} else {
		fmt.Printf("http response : %s\n", httpResp)
	}
}
