package main

import (
	"errors"
	"fmt"
	"github.com/gookit/goutil/dump"
	"log"
	"time"
)

// 模拟一个执行缓慢的函数
func slowFunc(s int, c chan string) {
	fmt.Printf("sleep begin : %v \n", time.Now().Format("2006-01-02 15:04:05"))
	msg := fmt.Sprintf("我是一个执行了%d秒的函数, 我刚刚执行完毕", s)
	time.Sleep(time.Duration(s) * time.Second)
	fmt.Printf("sleep end : %v \n", time.Now().Format("2006-01-02 15:04:05"))
	c <- msg
}

/**
channel 是goroutine之间, 相互通信的桥梁, 可以在各个goroutine之间收发消息
通道的特性:
	1. 同一时刻, 只能有一个goroutine对通道进行读写操作,
	2. 通道内的消息, 遵循先入先出的队列特征(保证收发消息的顺序)
	3. 通道类型的空值是nil , 使用make函数创建一个通道

使用无缓冲通道接收数据: <== 注意是无缓冲通道
	1. 通道内数据的收发, 在两个不同的goroutine之间进行(由于通道内的数据在没有接收方处理时, 发送方将一直保持阻塞状态, 故:无法在同一个goroutine内同时收发)
	2. 接收方将持续阻塞, 直到发送方发送数据(接收方接收时, 如果通道内没有发送方发送的数据, 接收方也会持续阻塞)
	3. 通道内一次只能接收一个数据元素

使用有缓冲通道收发数据: <-- 注意是有缓冲通道
	1. 给有缓冲通道发送数据, 发送长度满了之后, 将持续阻塞, 直到有其他goroutine从该通道内读取出数据来
	2. 从有缓冲通道接收数据, 如果读取的长度达到设定的长度之后, 仍继续读取, 将会持续阻塞, 直到有其他goroutine写入数据至该通道内
	3. 有缓冲的通道, 关闭后, 继续读取的话, 会返回通道内类型的零值
*/

func baseChannel() {
	/**
	通道使用make语法来创建, 如下:
	关键字chan后边定义的类型, 代表该通道将用来存储何种类型的数据
	make函数默认创建出来的通道是无缓冲通道, 无缓冲 channel 在同步发送和接收操作。 即使使用并发，通信也是同步的。
	*/
	chan1 := make(chan string)

	// chan2 := make(chan bool)
	// chan3 := make(chan []*int)
	// chan4 := make(chan interface{}) // 定义一个可以接收任何数据类型的通道
	// fmt.Printf("%#v\n %#v \n %#v\n", chan4)

	go slowFunc(3, chan1)

	msg, ok := <-chan1 // 非阻塞接收数据
	fmt.Printf("msg: %v, ok: %v \n", msg, ok)
	fmt.Printf("我是主线程, 我刚刚执行完毕\n")
	panic(errors.New("测试抛出一个panic"))
}

func sendMsgWhenChannelClosed() {

	defer func() {
		if err := recover(); err != nil {
			log.Println("触发一个panic: ", err)
		}
	}()

	ch := make(chan string, 5)

	ticker := time.NewTicker(time.Second)

	i := 0

	for {
		select {
		case <-ticker.C:
			if i == 2 { // 当第3秒的时候, 关闭通道
				log.Println("close ch channel")
				close(ch)
			}
			i += 1
			msg := time.Now().String()
			ch <- msg
			log.Println("send msg: ", msg)
		}
	}

}

func recMsgOnClosedChannel() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("get an panic: ", err)
		}
	}()

	type Msg struct {
		Msg    string
		Status bool
		b      byte
		u      rune
		list   []int
	}
	ch := make(chan Msg, 4)

	timer := time.NewTimer(time.Second * 3)
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-timer.C:
			go func() {
				for { // 如果不退出这个协程, 将会一直从一个已关闭的chan读取出它的零值
					time.Sleep(time.Second)
					val, ok := <-ch
					log.Printf("ch长度: %#v, val: %v, ok: %v \n", len(ch), val, ok)
				}
			}()
			log.Println("时间到, 关闭ch")
			ticker.Stop()
			timer.Stop()
			close(ch)
			return
		case <-ticker.C:
			msg := Msg{
				Msg:    time.Now().Format("2006-01-02 15:04:05"),
				Status: true,
				b:      'm',
				u:      0100,
				list:   []int{1, 2, 3},
			}
			log.Println("向ch发送消息: ", msg)
			ch <- msg
		}

		log.Println("ping...")
	}
	log.Println("当前函数执行的是recMsgOnClosedChannel")
}

func transformDataByChannel() {

	mystr := "hello world"
	ch := make(chan string)
	fmt.Printf("go启动传递之前, str的地址: %p \n", &mystr)
	go func(ch chan string, str string) {
		fmt.Printf("传递之前, str的地址: %p \n", &str)
		ch <- str
	}(ch, mystr)

	mystr = <-ch

	fmt.Printf("传递之后, str的地址: %+v \n", &mystr)

}

func testHowToUseReturnInForSelect() (int, bool) {
	ticker := time.NewTicker(time.Second)

	i := 0
	for {
		select {
		case <-ticker.C:
			i += 1
			fmt.Println(i)
			if i == 3 {
				//return 999, true
				break
			}
			fmt.Println("第" + fmt.Sprintf("%v", i) + "次ping...")
		}
	}

	return 0, false
	//fmt.Println("func : testHowToUseReturnInForSelect")
}
func main() {

	//go baseChannel() // 基本使用

	//go sendMsgWhenChannelClosed() // 当有缓冲通道关闭后, 再向其发送消息, 将会收到一个panic

	/**
	有缓冲通道读取消息
		1. 当通道中没有消息, 读取时, 会一直阻塞直到有消息进来
		2. 当通道中有消息, 且通道没关闭, 此时可写入消息, 读取将会读出消息, 直到读完后阻塞
		3. 当通道中有消息, 此时通道关闭, 此时不可写入消息, 读取将会读出内部消息, 继续读取将会读取chan中类型的零值,
		4. 当通道中没有消息, 且通道没关闭, 此时会一直阻塞读取, 直到有消息写入
	*/
	go recMsgOnClosedChannel()

	go transformDataByChannel() // 这里怎么验证数据在通道内传递, 是拷贝传递的呢? 比如把a发送给chan, 另一个goroutine在读取a时, 就是拷贝??

	/**
	for 循环中的select 中的case里, 使用return, 相当于break掉for循环后直接return了
	当在这个case中使用brake时, 只是跳过了这一次的for循环, 相当于continue, 它还是会继续for循环
	*/
	//res, ok := testHowToUseReturnInForSelect()
	//fmt.Println(res, ok)

	//另一种阻塞代码的方式
	//ch := make(chan bool)
	//<-ch

	// 测试一个有缓冲通道, 先阻塞读, 再关闭通道, 看发生什么(仍然是读取类型的零值)
	//go readingACacheChannelOnclose()

	for {
		time.Sleep(time.Second) // 这里用来确保main不挂
	}

}

func readingACacheChannelOnclose() {
	ch := make(chan int, 3)
	go func(ch chan int) {
		ticker := time.NewTicker(time.Second * 2)
		for {
			select {
			case <-ticker.C:
				fmt.Println("阻塞读取中...")
			case val, ok := <-ch:
				time.Sleep(time.Second)
				dump.P(val, ok)
			}
		}
	}(ch)

	time.Sleep(time.Second * 5)

	close(ch)
}
