package main

import (
	"fmt"
	"runtime"
	"time"
)

func testGoroutine()  {
	i := 0
	for {
		if i == 100 {
			break
		}
		fmt.Printf("当前程序执行结果: %#v\n", i)
		i += 1
	}
}

func say(hi string)  {
	for i := 0; i < 5; i++ {
		runtime.Gosched() // 这个函数表示让cpu把时间片让给别人, 下次某个时候继续恢复执行该goroutine
		fmt.Printf("当前程序执行结果: %#v\n", hi)
	}
}

func showFunc(ss int) {
	time.Sleep(time.Duration(ss) * time.Second)
	fmt.Printf("程序睡眠了%d秒 \n\n", ss)
}


/**
	Golang 多线程奉行的原则: 不要通过共享来通信, 要通过通信来共享!
 */
func main()  {
	/**
	 	此处只是实现了并发,
		在Go 1.5将标识并发系统线程个数的runtime.GOMAXPROCS的初始值由1改为了运行环境的CPU核数

		Golang 在runtime.GOMAXPROCS的数量与任务数量相等时, 可以做到并行执行, 一般情况下, 都是并发执行

		goroutine 属于是抢占式任务处理, 谁优先抢到了, 谁先处理
	 */
	runtime.GOMAXPROCS(4)

	//go say("Layery")
	//say("hello")

	// 获取当前运行环境的cpu核心数
	fmt.Printf("当前操作系统的核心数: %d \n\n", runtime.NumCPU())

	/**
	匿名自执行goroutine
	 */
	go func(ss int) {
		showFunc(ss)
	}(3)
	//go showFunc(5)
	fmt.Println("睡眠函数执行完之后, 我才开始执行")
	// 如果没有这一句, 由于goroutine的立即返回的特性, 程序将继续执行后边的println函数,并结束程序
	// 加了一个睡8秒函数, 程序将不会结束, 直到等待goroutine执行完毕
	time.Sleep(10 * time.Second)
}
