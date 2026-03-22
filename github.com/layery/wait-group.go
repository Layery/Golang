package main

import (
	"sync"
	//"time"
)

const N = 3

var wg = &sync.WaitGroup{}

func wait_group_misstake() {
	for i := 0; i < N; i++ {
		go func(i int) {
			wg.Add(1)
			println(i)
			defer wg.Done()
		}(i)
	}
	wg.Wait()
}

func wait_group_fixbug() {
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			println(i)
		}(i)
	}
	wg.Wait()
}

func main() {

	// waitgroup一个常见的错误,由于go执行的太快了, 可能导致go协程还没开始执行, main就执行完毕了
	// 多次执行可能会发生没有任何输出的情况
	// 解决这个问题, 可以参考
	//wait_group_misstake()

	// 如何解决?, 把wg.Add()方法, 提到go func()的外边来,确保每注册一个待执行的协程, 就会+1,
	// 避免wg.Add(1)和循环产生的协程数量不对等
	wait_group_fixbug()
}
