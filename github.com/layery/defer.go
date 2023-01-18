package main

import (
	"log"
	"time"
)

func slow() {
	defer trace("模拟一个耗时操作")()
	time.Sleep(time.Second * 2)

	log.Println("执行完成")
}

func trace(msg string) func() {
	start := time.Now()
	log.Println("begin enter ", msg)
	return func() {
		log.Printf("out %s (%s)", msg, time.Since(start))
	}

}

func main() {
	// defer对具名返回值函数的影响

	slow()

}
