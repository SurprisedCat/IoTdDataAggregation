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
	req := bytes.NewBuffer(dataJSON)
	resp, err := http.Post("http://"+string(serverAddr)+":"+string(httpPort)+"/v1/upload/single", "application/json;charset=utf-8", req)
	if err != nil {
		utils.CheckErr(err, "HTTP POST error")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.CheckErr(err, "HTTP POST error")
	}
	resp.Body.Close()
	httpRespTest := map[string][]byte{}
	json.Unmarshal(body, &httpRespTest)
	fmt.Printf("http response : %s\n", httpRespTest)
	httpwg.Done()

}
