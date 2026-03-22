package utils

import (
	"fmt"
)

func ConsoleLog(params ...interface{}) {
	for _, v := range params {
		fmt.Printf("%#v \n", v)
	}
}
