package main

import (
	"fmt"
	"time"
)

// 模拟一个执行缓慢的函数
func slowFunc(s int, c chan string)  {
	fmt.Printf("sleep begin : %v \n", time.Now())
	msg := fmt.Sprintf("我是一个执行了%d秒的函数, 我刚刚执行完毕", s)
	time.Sleep(time.Duration(s) * time.Second)
	fmt.Printf("sleep end : %v \n", time.Now())
	c <- msg
}

/**
	channel 是goroutine之间, 相互通信的桥梁, 可以在各个goroutine之间收发消息
	通道的特性:
		1. 同一时刻, 只能有一个goroutine对通道进行读写操作,
		2. 通道内的消息, 遵循先入先出的队列特征(保证收发消息的顺序)
		3. 通道类型的空值是nil , 使用make函数创建一个通道
	使用通道接收数据:
		1. 通道内数据的收发, 在两个不同的goroutine之间进行(由于通道内的数据在没有接收方处理时, 发送方将一直保持阻塞状态, 故:无法在同一个goroutine内同时收发)
		2. 接收方将持续阻塞, 直到发送方发送数据(接收方接收时, 如果通道内没有发送方发送的数据, 接收方也会持续阻塞)
		3. 通道内一次只能接收一个数据元素
 */
func main() {

	/**
		通道使用make语法来创建, 如下:
		关键字chan后边定义的类型, 代表该通道将用来存储何种类型的数据
	 */
	chan1 := make(chan string)

	// chan2 := make(chan bool)
	// chan3 := make(chan []*int)
	// chan4 := make(chan interface{}) // 定义一个可以接收任何数据类型的通道
	// fmt.Printf("%#v\n %#v \n %#v\n", chan4)

	go slowFunc(8, chan1)

	msg, ok := <- chan1 // 非阻塞接收数据

	//msg := <- chan1 // 阻塞接收数据

	fmt.Printf("ok ====> %#v \n", ok)

	fmt.Printf("msg ====> %#v \n", msg)

	fmt.Printf("我是主线程, 我刚刚执行完毕\n")
	//time.Sleep(5 * time.Second)
}
