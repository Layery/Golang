package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"net/http"
	"runtime"
	"time"
)

// 定义接收响应的chan, 记录响应时间, api
type Result struct {
	RunningTime interface{}
	Desc        string
}

func testGoroutine() {
	i := 0
	for {
		if i == 100 {
			break
		}
		fmt.Printf("当前程序执行结果: %#v\n", i)
		i += 1
	}
}

func say(hi string) {
	for i := 0; i < 5; i++ {
		runtime.Gosched() // 这个函数表示让cpu把时间片让给别人, 下次某个时候继续恢复执行该goroutine
		fmt.Printf("当前程序执行结果: %#v\n", hi)
	}
}

func showFunc(ss int) {
	time.Sleep(time.Duration(ss) * time.Second)
	fmt.Printf("程序睡眠了%d秒 \n\n", ss)
}

func SayGreetings(greeting string, times, label int, wg *sync.WaitGroup) {
	for i := 0; i <= times; i++ {
		//fmt.Printf("%v\n", greeting)
		//d := time.Second * time.Duration(rand.Intn(5)) / 2
		d := time.Second * 2
		log.Println(fmt.Sprintf("label: %d, i: %d, seelp: %d, value: ====>%v", label, i, d, greeting))
		time.Sleep(d) // 随机睡眠0到2.5秒
	}
	wg.Done()
}


func main() {
	log.Println("begin")
	myTestFunc7()

	//myTestFunc6()
	//myTestFunc5()

	//myTestFunc4()

	//myTestFunc3()
	//myTestFunc2()
	//myTestFunc1()

	select {}

}

func getApi(api string, ch chan Result, wg *sync.WaitGroup, timeout time.Duration) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, api, nil)
	if err != nil {
		panic(err)
	}
	_, err = http.DefaultClient.Do(req)

	res := Result{
		RunningTime: time.Since(start).Seconds(),
	}
	select {
	case <-ctx.Done(): // 单个api超时
			res.Desc = "error: api is timeout, api: " + api
	default:
		if err != nil {
			res.Desc = "error: api响应异常, api: " + api
		} else {
			res.Desc = "success: api响应正常, api: " + api
		}
	}
	ch <- res
	wg.Done()
	return
}

/**
1. 并发获取资源
2. 有api超时, 需要知道是哪个超时
3. 等待所有的api获取完毕后, 获得接口响应时间
*/
func myTestFunc7() {
	apis := []string{
		"https://management.azure.com",
		"https://dev.azure.com",
		"http://proceeding.cc/api/index/index",
		"https://api.github.com",
		"https://outlook.office.com/",
		"https://api.somewhereintheinternet.com/",
		"https://graph.microsoft.com",
	}

	// 定义一个chan, 接收每个api的响应
	ch := make(chan Result, len(apis))

	wg := sync.WaitGroup{}

	for _, api := range apis {
		wg.Add(1)
		go getApi(api, ch, &wg, time.Second * 1) // 单个api 1秒后超时
	}
	wg.Wait()

	for resp := range ch{
		log.Println(resp)
	}
}


//func myTestFunc6() {
//	start := time.Now()
//
//	// 定义一个channel
//	ch := make(chan string, 6)
//	apis := []string{
//		"https://management.azure.com",
//		"https://dev.azure.com",
//		"https://api.github.com",
//		"https://outlook.office.com/",
//		"https://api.somewhereintheinternet.com/",
//		"https://graph.microsoft.com",
//	}
//	for _, api := range apis {
//		go getApi(api, ch)
//	}
//	fmt.Println(<-ch)
//	fmt.Println(<-ch)
//	fmt.Println(<-ch)
//	fmt.Println(<-ch)
//	fmt.Println(<-ch)
//	fmt.Println(<-ch)
//	close(ch)
//
//	limit := time.Since(start)
//	fmt.Printf("Done! it took %v seconds \n", limit.Seconds())
//}

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

/**
本质上还是for range 在同一个地址空间遍历数据的问题
*/
func myTestFunc4() {
	var slice []func()

	sli := []int{1, 2, 3, 4, 5}
	for _, v := range sli {
		fmt.Println(&v)
		temp := v
		slice = append(slice, func() {
			fmt.Println(temp * temp) // 直接打印结果
		})
	}

	for _, val := range slice {
		val()
	}
	// 输出 25 25 25 25 25
}

func myTestFunc3() {
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

func myTestFunc2() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	//c := make(chan string)

	// WaitGroup类型有3个方法, Add, Done 和Wait ,
	// Add 用来注册需要完成的线程数(任务数)
	// Done 某个线程(任务)执行完后通知主线程
	// Wait 将会一直阻塞, 直到注册过的线程都执行完毕之后, 才会继续执行后续的代码
	wg := sync.WaitGroup{}
	fmt.Println(time.Now())
	wg.Add(2)
	go SayGreetings("hi", 10, 1, &wg)
	go SayGreetings("weidingyi", 10, 2, &wg)

	wg.Wait()

	fmt.Println(time.Now())
	fmt.Println("两个线程都执行完毕了\n")
}

/**
Golang 多线程奉行的原则: 不要通过共享来通信, 要通过通信来共享!
*/
func myTestFunc1() {
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
