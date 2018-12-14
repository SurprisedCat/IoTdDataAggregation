package config

import (
	"encoding/json"
	"net"

	"../utils"
)

//README 说明，
var README = "使用rawsocket发送/接收的时候需要root权限，或者sudo；聚合的指标从包书数量，改成了字节数"

//Debug true at present
var Debug = true

//ServerAddr 云端的地址
var ServerAddr = []byte("10.103.238.174")

//AggregatorAddr 聚合节点的IP地址
var AggregatorAddr = []byte("10.112.17.170")

//AggreIfaceName 聚合节点的MAC地址
var AggreIfaceName = "wlp58s0" //raw socket device for aggregator
//BackPort 聚合节点传输给服务器所使用的协议
var BackPort = []byte("8080") //8080 http 5683 coap

//ProtocolType 测试的协议类型 http coap mqtt 7676(代表rawsocket)
var ProtocolType = "coap"

//聚合的包数目
// var AggreLength int = 50

//AggreLength 聚合的字节长度，字节超过20汇聚节点发送
var AggreLength = 1200

//ClientIfaceName 客户节点的MAC地址
var ClientIfaceName = "wlp4s0" //raw socket device for aggregator

//DstMac 发送的聚合节点的MAC地址
var DstMac = net.HardwareAddr{0xf8, 0x63, 0x3f, 0x42, 0x04, 0x00}

//TotalReq 总请求量
var TotalReq = 150

//TimeGap 时间间隔
var TimeGap = 100000 //微秒为单位

//Cluster 是否启用聚合方式 true/false
var Cluster = true

/**********payload***************/

//SelfID 产生ID的字符串，会影响payload长度,可以不更改
var SelfID = "cx"

//PayloadData 具体的payload信息
var PayloadData = []byte("I am OAI cx")

//TotalPayloadJSONLength 计算负载的整体长度，不用改动
var TotalPayloadJSONLength int

//初始化整体payload信息，不要改这里
func init() {
	//添加额外的解密载荷
	if ProtocolType == "http" {
		PayloadData = append(PayloadData, []byte("00000000000000000000000000000000")...)
	} else if ProtocolType == "coap" {
		PayloadData = append(PayloadData, []byte("0000000000000000")...)
	} else if ProtocolType == "mqtt" {
		PayloadData = append(PayloadData, []byte("00000000000000000000000000000000")...)
	} else if ProtocolType == "7676" {

	} else {

	}
	TotalPayloadJSON, _ := json.Marshal(map[string][]byte{"ID": utils.GetClientID(SelfID), "data": PayloadData})
	TotalPayloadJSONLength = len(TotalPayloadJSON)

}
