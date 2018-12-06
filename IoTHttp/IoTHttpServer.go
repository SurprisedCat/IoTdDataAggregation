package linuxdashboard

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

//RouterRegister Register router information
func RouterRegister() *gin.Engine {
	fmt.Println("Linux Dashboard")
	router := gin.Default()

	router.LoadHTMLFiles("linuxdashboard/godashboard/index.html")

	router.GET("/", IndexAPI)
	router.GET("/proc/stat", ProcStat)

	return router
}

//IndexAPI 显示主界面
func IndexAPI(c *gin.Context) {

	c.HTML(http.StatusOK, "index.html", gin.H{})
}

/*
* Data are encoded into json format
 */

//ProcStat data of cpu utilization
func ProcStat(c *gin.Context) {
	res := CmdExec("cat /proc/stat | head -n 1 | awk '{$1=\"\";print}'")
	resArray := strings.Split(res[0], " ")
	var cpu []int64
	var totalcpu, idlecpu int64
	for _, v := range resArray {
		temp, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			cpu = append(cpu, temp)
			totalcpu = totalcpu + temp
		}
	}
	idlecpu = cpu[3]
	c.JSON(http.StatusOK, gin.H{
		"totalcpu": totalcpu,
		"idlecpu":  idlecpu,
	})
}
