package main

import (
	"log"
	"sync"
	"time"
)

func chanTranslateUse[T any](ch chan T, data T, goModify func(T), mainModify func(*T), name string) {
	ch <- data

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		received := <-ch
		time.Sleep(time.Millisecond * 20)
		goModify(received)
	}()

	time.Sleep(time.Millisecond * 10)
	mainModify(&data)

	wg.Wait()
	log.Printf("[%s] 主 goroutine 的 data 最终值: %v\n\n", name, data)
}

func testValueType() {
	log.Println("=== 值类型 int — 预期：无 data race ===")
	ch := make(chan int, 1)
	chanTranslateUse(ch, 100,
		func(x int) {
			x = x - 20
			log.Printf("  goroutine 修改 x=%d", x)
		},
		func(p *int) {
			*p = *p - 10
			log.Printf("  主 goroutine 修改 *p=%d", *p)
		},
		"int(值类型)",
	)
}

func testSliceNoExpand() {
	log.Println("=== 切片(未扩容) — 预期：有 data race ===")
	ch := make(chan []int, 1)
	data := []int{1, 2, 3, 4, 5}
	chanTranslateUse(ch, data,
		func(x []int) {
			x[3] = 300
			log.Printf("  goroutine 修改 x[3]=%d, x=%v", x[3], x)
		},
		func(p *[]int) {
			(*p)[3] = 999
			log.Printf("  主 goroutine 修改 (*p)[3]=%d, data=%v", (*p)[3], *p)
		},
		"切片(未扩容)",
	)
}

func testSliceWithExpand() {
	log.Println("=== 切片(扩容) — 预期：无 data race ===")
	ch := make(chan []int, 1)
	data := make([]int, 5, 5)
	for i := range data {
		data[i] = i + 1
	}
	chanTranslateUse(ch, data,
		func(x []int) {
			x = append(x, 6, 7)
			x[3] = 300
			log.Printf("  goroutine 修改后 x=%v (len=%d, cap=%d)", x, len(x), cap(x))
		},
		func(p *[]int) {
			(*p)[3] = 999
			log.Printf("  主 goroutine 修改 (*p)[3]=%d, data=%v (len=%d, cap=%d)", (*p)[3], *p, len(*p), cap(*p))
		},
		"切片(扩容)",
	)
}

func testMapType() {
	log.Println("=== map — 预期：有 data race ===")
	ch := make(chan map[string]int, 1)
	data := map[string]int{"a": 1, "b": 2}
	chanTranslateUse(ch, data,
		func(x map[string]int) {
			x["a"] = 300
			log.Printf("  goroutine 修改 x[a]=%d, x=%v", x["a"], x)
		},
		func(p *map[string]int) {
			(*p)["a"] = 999
			log.Printf("  主 goroutine 修改 (*p)[a]=%d, data=%v", (*p)["a"], *p)
		},
		"map",
	)
}

func testPointerType() {
	log.Println("=== 指针 *int — 预期：有 data race ===")
	ch := make(chan *int, 1)
	val := 100
	chanTranslateUse(ch, &val,
		func(x *int) {
			*x = 200
			log.Printf("  goroutine 修改 *x=%d", *x)
		},
		func(p **int) {
			**p = 300
			log.Printf("  主 goroutine 修改 **p=%d", **p)
		},
		"*int(指针)",
	)
}

func testFuncNonClosure() {
	log.Println("=== func(非闭包) — 预期：无 data race ===")
	ch := make(chan interface{}, 1)

	add := func(a, b int) int { return a + b }
	ch <- add

	go func() {
		f := (<-ch).(func(int, int) int)
		r := f(1, 2)
		log.Printf("  goroutine 调用 add(1,2) = %d", r)
	}()

	r := add(3, 4)
	log.Printf("  主 goroutine 调用 add(3,4) = %d", r)
	time.Sleep(time.Millisecond * 50)
	log.Println("  → 无 data race（函数代码不可变，只读）\n")
}

func testFuncClosure() {
	log.Println("=== func(闭包捕获变量) — 预期：有 data race ===")
	ch := make(chan interface{}, 1)

	counter := 0
	incr := func() int { counter++; return counter }
	ch <- incr

	go func() {
		f := (<-ch).(func() int)
		_ = f()
		log.Println("  goroutine 调用 incr() 完成")
	}()

	_ = incr()
	log.Println("  主 goroutine 调用 incr() 完成")
	time.Sleep(time.Millisecond * 50)
	log.Println("  → 有 data race（闭包捕获的 counter 被并发修改）\n")
}

func main() {
	testValueType()
	testSliceNoExpand()
	testSliceWithExpand()
	testMapType()
	testPointerType()
	testFuncNonClosure()
	testFuncClosure()
}
