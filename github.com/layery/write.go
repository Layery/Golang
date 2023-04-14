package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

func test1() {
	file, err := os.Create("./write_output.txt")
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}
	defer file.Close()
	// 用于协调所有goroutine的通道
	order := make(chan int)
	// 启动10个goroutine并发写入文件
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int, ch chan int) {
			defer wg.Done()
			// 等待前面的goroutine完成写入
			currNum := <-order
			// 写入当前goroutine的数据
			currTime := time.Now().UnixNano()
			_, err := fmt.Fprintln(file, fmt.Sprintf("%v: Hello from goroutine %d", currTime, currNum))
			if err != nil {
				fmt.Println("Failed to write to file:", err)
				return
			}
			currNum += 1
			order <- currNum

		}(i, order)

	}
	// 向通道发送第一个信号，开始写入
	order <- 0
	// 等待所有goroutine完成写入
	wg.Wait()

	b, err := ioutil.ReadFile("./write_output.txt")
	if err != nil {
		log.Print(err)
	}
	fmt.Println(string(b))
}

func test2(ch chan int) {

	fmt.Println("start")
	ch <- 222
	fmt.Println("能写入吗? ")
}
func main() {

	var ch = make(chan int)
	// go test1()

	go test2(ch)

	go func(ch chan int) {
		rs := <-ch

		fmt.Println("res : ", rs)

	}(ch)

	select {}

}
