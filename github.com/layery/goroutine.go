package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"
)


var wg sync.WaitGroup

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


func SayGreetings(greeting string, times, label int) {
	for i := 0; i <= times; i++ {
		//fmt.Printf("%v\n", greeting)
		//d := time.Second * time.Duration(rand.Intn(5)) / 2
		d := time.Second * 2
		log.Println(fmt.Sprintf("label: %d, i: %d, seelp: %d, value: ====>%v", label, i,  d,  greeting))
		time.Sleep(d) // 随机睡眠0到2.5秒
	}
	wg.Done()
}

func main() {

	myTestFunc6()
	//myTestFunc5()

	//myTestFunc4()

	//myTestFunc3()
	//myTestFunc2()
	//myTestFunc1()



}

func getApi(api string, ch chan string) {
	_, err := http.Get(api)
	if err != nil {
		ch<- fmt.Sprintf("curr api %v is error %v \n", api, err)
		return
	}
	ch <- fmt.Sprintf("curr api is success %v \n", api)
}

func myTestFunc6() {
	start := time.Now()

	// 定义一个channel
	ch := make(chan string)
	apis := []string{
		"https://management.azure.com",
		"https://dev.azure.com",
		"https://api.github.com",
		"https://outlook.office.com/",
		"https://api.somewhereintheinternet.com/",
		"https://graph.microsoft.com",
	}
	for _, api := range apis {
		go getApi(api, ch)
	}
	//msg <-ch
	//
	//fmt.Printf("channel msg is: %v \n", msg)

	fmt.Print(<-ch)
	fmt.Print(<-ch)
	fmt.Print(<-ch)
	fmt.Print(<-ch)
	fmt.Print(<-ch)
	fmt.Print(<-ch)
	fmt.Print(<-ch)
	fmt.Print(<-ch)
	fmt.Print(<-ch)

	limit := time.Since(start)
	fmt.Printf("Done! it took %v seconds \n", limit.Seconds())
}



func myTestFunc5() {
	start := time.Now()

	apis := []string{
		"https://management.azure.com",
		"https://dev.azure.com",
		"https://api.github.com",
		"https://outlook.office.com/",
		"https://api.somewhereintheinternet.com/",
		"https://graph.microsoft.com",
	}
	for _, api := range apis {
		_, err := http.Get(api)

		if err != nil {
			fmt.Printf("error: %s is down \n", api)
			continue
		}
		fmt.Printf("success: %s is up and running \n", api)
	}
	limit := time.Since(start)
	fmt.Printf("Done! it took %v seconds \n", limit.Seconds())
}




func myTestFunc4() {
	var slice []func()

	sli := []int{1, 2, 3, 4, 5}
	for _, v := range sli {
		fmt.Println(&v)
		temp := v
		slice = append(slice, func(){
			fmt.Println(temp * temp) // 直接打印结果
		})
	}

	for _, val  := range slice {
		val()
	}
	// 输出 25 25 25 25 25
}



func myTestFunc3()  {
	func() {
		for i := 0; i < 3; i++ {
			defer fmt.Println("a:", i)
		}
	}()
	fmt.Println()
	func() {
		for i := 0; i < 3; i++ {
			defer func() {
				fmt.Println("b:", i)
			}()
		}
	}()
}


func myTestFunc2 () {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	//c := make(chan string)

	// WaitGroup类型有3个方法, Add, Done 和Wait ,
	// Add 用来注册需要完成的线程数(任务数)
	// Done 某个线程(任务)执行完后通知主线程
	// Wait 将会一直阻塞, 知道注册过的线程都执行完毕之后, 才会继续执行后续的代码

	fmt.Println(time.Now())
	wg.Add(2)
	go SayGreetings("hi", 10, 1)
	go SayGreetings("weidingyi", 10, 2)

	wg.Wait()

	fmt.Println(time.Now())
	fmt.Println("两个线程都执行完毕了\n")
}


/**
	Golang 多线程奉行的原则: 不要通过共享来通信, 要通过通信来共享!
 */
func myTestFunc1()  {
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
