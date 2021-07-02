package logger

import "fmt"

func Out(data interface{}) {
	fmt.Printf("当前程序返回的是: %#v \n", data)
}