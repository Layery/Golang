package main

import (
	"fmt"
)

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
		make函数默认创建出来的通道是无缓冲通道, 无缓冲 channel 在同步发送和接收操作。 即使使用并发，通信也是同步的。

		如果想创建有缓冲的channel, 可以在make函数的第2个参数, 设置为通道的缓冲大小, 如下所示
	    当 channel 已满时，任何发送操作都将等待，直到有空间保存数据。
	               相反，如果 channel 是空的且存在读取操作，程序则会被阻止，直到有数据要读取。
	 */
	size := 2
	ch := make(chan string, size)
	send(ch, "one")
	send(ch, "two")  // 至此, 主协程已经将ch容量占满,
	go send(ch, "three") // 这两个子协程虽然在往里send数据, 只是没有send进去
	go send(ch, "four")
	//time.Sleep(10 * time.Second) // 主协程睡10秒
	fmt.Println("主协程数据send完毕 \n\n")

	// 主协程开始读取ch里的数据, 同时, 两个go子协程里, 因为for读取之后有了空余空间
	// 所以它们两个可以往ch里塞数据了
	// 所以这个for循环里, 可以将4条数据都读取出来
	for i := 0; i < 4 ; i ++ {
		fmt.Println("msg => " + <-ch)
	}
}

func send(ch chan string, message string) {
	fmt.Println("start:  " + message)
	ch <- message
	fmt.Println("end:  " + message)
}