package main

import (
	"fmt"
	"time"
)

func selectForWithBreak() {
	i := 0
	for {
		select {
		case <-time.After(time.Second * time.Duration(2)):
			i++
			if i == 5 {
				fmt.Println("[selectForWithBreak]当i == 5时,执行break,仍然无法跳出for循环")
				break
			} else {
				fmt.Println("[selectForWithBreak]inside the select, curr i : ", i)
			}
		}
	}
}

func selectForWithReturn() {

	i := 0
	for {
		select {
		case <-time.After(time.Second):
			if i == 5 {
				fmt.Println("[selectForWithReturn]当i == 5时,执行return, 跳出for循环")
				return
			} else {
				fmt.Println("[selectForWithReturn]inside the select, curr i : ", i)
			}
			i += 1
		}
	}
}

func main() {

	// for循环中的select, 使用break, 不能退出for循环
	go selectForWithBreak()

	// for 循环中的select, 使用return, 可以退出for循环
	go selectForWithReturn()

	select {}
}
