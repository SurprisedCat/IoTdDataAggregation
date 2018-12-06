package utils

import (
	"fmt"
	"os"
	"strconv"
)

/*
ClinetData parameters for main function
*/
type ClinetData struct {
	id         string
	dataOrigin string
}

//ClientParameter = Parameter{"127.0.0.1","I am a phd"}

/*
CheckErr 错误处理函数
*/
func CheckErr(err error, extra string) bool {
	if err != nil {
		formatStr := " Err : %s\n"
		if extra != "" {
			formatStr = extra + formatStr
		}

		fmt.Fprintf(os.Stderr, formatStr, err.Error())
		return true
	}

	return false
}

/*
GetClientID get the hostname and uid of the user
*/
func GetClientID() []byte {
	//Get the hostname of the machine
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	//clientID = hostname + uid
	clientID := hostname + strconv.FormatInt(int64(os.Getuid()), 10)
	return []byte(clientID)
}
