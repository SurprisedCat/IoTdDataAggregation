package common

import (
	"fmt"
	"os"
)

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
