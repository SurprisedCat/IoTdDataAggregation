package config

import (
	"net"
)

//ServerAddr 云端的地址
var ServerAddr []byte = []byte("127.0.0.1")

//AggregatorAddr 聚合节点的IP地址
var AggregatorAddr []byte = []byte("10.112.43.97")

//聚合节点的MAC地址
var AggreIfaceName string = "wlp61s0" //raw socket device for aggregator
var BackPort []byte = []byte("5683")  //8080 http 5683 coap

//ProtocolType 测试的协议类型
var ProtocolType string = "coap"

//聚合的包数目
var AggreLength int = 20

//客户节点的MAC地址
var ClientIfaceName string = "wlp61s0" //raw socket device for aggregator
//聚合节点的MAC地址
var DstMac net.HardwareAddr = net.HardwareAddr{0xa0, 0x88, 0x69, 0x16, 0xda, 0xb4}

//总请求量
var TotalReq int = 150
var TimeGap int = 100000 //微秒为单位

//是否启用聚合方式
var Cluster bool = true

//payload
var PayloadData []byte = []byte("I am OAI cx")
