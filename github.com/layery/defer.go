package main

import (
	"fmt"
	"github.com/gookit/goutil/dump"
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

func mianshiti() {
	// 原题中, str这个变量非常具有迷惑性, 让人无法一眼分辨出他的长度是多少
	// 忘记是啥了, 但是绝对是有[]byte这几个字符

	str := []byte("hello")
	if len(str) == 2 {
		defer fmt.Println(1)
	} else {
		defer fmt.Println(2)
	}
	defer fmt.Println(3)
}

func get_sum(a, b int) (sum int) {
	defer func() {
		sum = 45
	}()
	sum = a + b
	return sum
}

func main() {
	/**
	defer对具名返回值函数的影响
	因为具名返回值函数, 在函数声明阶段就已经声明了一个返回值变量
	而defer执行在return之前, 所以可以修改这个返回值变量
	*/
	res := get_sum(3, 4)
	dump.P(res)

	slow()

}
