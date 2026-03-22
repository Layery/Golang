package main

import (
	"fmt"
	"sort"
	"sync"
)

func quickSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	// 选择基准值（pivot），这里选择数组的第一个元素
	pivot := arr[0]
	left := []int{}  // 存储小于基准值的元素
	right := []int{} // 存储大于基准值的元素

	// 遍历数组，根据元素与基准值的大小关系，分配到 left 或 right
	for _, value := range arr[1:] {
		if value <= pivot {
			left = append(left, value)
		} else {
			right = append(right, value)
		}
	}

	// 递归地对 left 和 right 进行排序，并拼接结果
	return append(append(quickSort(left), pivot), quickSort(right)...)
}

func popSort(arr []int) []int {
	cnt := len(arr)
	for i := 0; i < cnt-1; i++ {
		for j := 0; j < cnt-i-1; j++ {
			if arr[j] > arr[j+1] {
				temp := arr[j]
				arr[j] = arr[j+1]
				arr[j+1] = temp
			}
		}
	}
	return arr
}

func deferFunc() (res string) {
	defer func() {
		res = "test"
	}()
	return "llf"
}

func bufferProducer() {
	// 1. 单生产者按需生产1-30的整数,3个消费者消费这个队列,不保证有序, 输出这些整数

	// 2. 使用无缓冲的chan来实现一个生产者, 消费者

	ch := make(chan int, 50)
	wg := &sync.WaitGroup{}

	type CntStruct struct {
		cnt int
		mu  *sync.Mutex
	}
	cnt := CntStruct{cnt: 0, mu: &sync.Mutex{}}

	go func(ch chan int) {
		defer close(ch)
		for i := 0; i < 20; i++ {
			ch <- i
		}
	}(ch)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(goIndex int) {
			defer wg.Done()
			for msg := range ch {
				cnt.mu.Lock()
				cnt.cnt += 1
				cnt.mu.Unlock()
				fmt.Printf("goIndex : %v, msg: %v \n", goIndex, msg)
			}
		}(i)
	}
	wg.Wait()

	fmt.Println("共", cnt.cnt)
}

func noBufferProducer() {
	ch := make(chan int)
	//wg := sync.WaitGroup{}

	for i := 0; i < 20; i++ {
		ch <- i
		fmt.Println("写入消息", i)
	}
	go func(ch chan int) {
		for msg := range ch {
			fmt.Println(msg)
		}
	}(ch)
}

func filterSlice(param1, param2 []int) (res1 []int, res2 []int) {
	if len(param1) == 0 && len(param2) == 0 {
		return []int{}, []int{}
	}
	for _, val := range param1 {
		if val%2 == 0 {
			continue
		}
		res1 = append(res1, val*2)
	}

	for _, val := range param2 {
		if val%2 == 0 {
			continue
		}
		res2 = append(res2, val*2)
	}

	if len(res1) == 0 {
		defer func() {
			res1 = []int{}
		}()
	}
	if len(res2) == 0 {
		defer func() {
			res2 = []int{}
		}()
	}

	sort.Ints(res1)
	sort.Ints(res2)
	return res1, res2
}

func sliceTest(wg *sync.WaitGroup) {
	defer wg.Done()
	// 全部是偶数
	input1 := []int{2, 6, 0, 4, 8, 10, 1}
	// 有基数, 有偶数
	input2 := []int{1, 7, 8, 9, 10, 2, 3, 4, 5, 6}
	res1, res2 := filterSlice(input2, input1)
	fmt.Println(res1, res2)
}

//
//

func main() {
	// 从网上抓取简历
	// job:2015-01-01 ~ 2017-01-01
	// job:2018-01-01 ~ 2021-01-01
	// job:2016-01-01 ~ 2017-01-01

	// 每个不同的job,

}
