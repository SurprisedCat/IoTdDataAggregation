package iothttp

import (
	"fmt"
	"net/http"

	"../auth"
	"../simssl"
	"github.com/gin-gonic/gin"
)

//RouterRegister Register router information
func RouterRegister() *gin.Engine {
	fmt.Println("IoT data aggregation HTTP RESTful")
	router := gin.Default()

	router.LoadHTMLFiles("../iothttp/index.html")

	router.GET("/", IndexAPI)
	router.POST("/v1/upload/single", ProcSingle)
	router.POST("/v1/upload/cluster", ProcCluster)

	return router
}

//IndexAPI 显示主界面
func IndexAPI(c *gin.Context) {

	c.HTML(http.StatusOK, "index.html", gin.H{})
}

/*
* Data are encoded into json format
 */

//ProcSingle process the single data packet
func ProcSingle(c *gin.Context) {

	recPost := map[string][]byte{}
	err := c.BindJSON(&recPost)

	clientID := recPost["ID"]
	clientData := recPost["data"]
	eKey, vali := auth.GetValidationKeyServer([]byte(clientID))
	if vali == false {
		c.JSON(http.StatusOK, gin.H{
			"status": []byte("error"),
		})
		return
	}
	originData, err := simssl.AesDecrypt([]byte(clientData), eKey)
	fmt.Println(originData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": []byte("error"),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": []byte("OK"),
	})
}

//ProcCluster process the cluster upload data
func ProcCluster(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"totalcpu": 1,
		"idlecpu":  1,
	})
}
