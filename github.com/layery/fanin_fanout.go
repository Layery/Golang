package main

import (
	"context"
	"fmt"
	"sync"
)

// 数据源
func generator(ctx context.Context, nums []int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case <-ctx.Done():
				return
			case out <- n:
			}
		}
	}()
	return out
}

// 单任务处理函数
func worker(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case <-ctx.Done():
				return
			case out <- n * n:
			}
		}
	}()
	return out
}

// Fan-in：合并多个channel输出
func fanIn(ctx context.Context, channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		defer wg.Done()
		for val := range c {
			select {
			case <-ctx.Done():
				return
			case out <- val:
			}
		}
	}

	wg.Add(len(channels))
	for _, ch := range channels {
		go output(ch)
	}

	// 所有worker完成后关闭输出通道
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	datalist := make([]int, 100)
	for i := 0; i < 100; i++ {
		datalist[i] = i
	}
	source := generator(ctx, datalist)

	// Fan-out：启动3个并发worker
	w1 := worker(ctx, source)
	w2 := worker(ctx, source)
	w3 := worker(ctx, source)
	w4 := worker(ctx, source)
	w5 := worker(ctx, source)

	// Fan-in：汇总结果
	result := fanIn(ctx, w1, w2, w3, w4, w5)

	// 消费汇总后的数据
	for v := range result {
		fmt.Printf("get result: %d\n", v)
	}

	fmt.Println("all stream finished")
}
