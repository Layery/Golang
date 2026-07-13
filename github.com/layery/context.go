package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func contextWithCancel(ctxBackground context.Context) {
	ctx, cancel := context.WithCancel(ctxBackground)

	wg := sync.WaitGroup{}

	defer func() {
		if err := recover(); err != nil {
			log.Println("system err: ", err)
		}
	}()

	ticker := time.NewTicker(time.Second)
	i := 0
	for i < 3 {
		wg.Add(1)
		go func(ctx context.Context, index int) {
			defer func() {
				if err := recover(); err != nil {
					log.Println("sub goroutine panic: ", err)
				}
				wg.Done()
			}()

			for {
				select {
				case <-ctx.Done():
					log.Println("停止: " + strconv.Itoa(index) + "收到停止信号")
					return // todo 这里如果不return, 当前这个子协程将不会自动退出, 它会一直能接收到停止信号
				case <-ticker.C:
					log.Println("我是第" + strconv.Itoa(index) + "个协程")
				}
			}
		}(ctx, i)
		i += 1
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\n捕获到Ctrl+C终止信号，准备停止所有协程...")
		// 触发context取消，通知所有子协程
		cancel()
	}()
	wg.Wait()
	fmt.Println("优雅退出")
}

func contextWithValue(ctxBackground context.Context) {
	type customKey struct{}

	key := customKey{}
	wg := sync.WaitGroup{}

	data := map[string]interface{}{
		"name": "llf",
		"age":  33,
	}
	ctxBackground = context.WithValue(ctxBackground, key, data)
	wg.Go(func() { // 这个wg.Go方法, 是1.25版本之后才有的, 入参是个函数且需要确保不会panic, 并且仍然需要显式执行wg.Wait()
		log.Println("我是子协程, 我获取到了来自父ctx的参数", ctxBackground.Value(key))
	})
	wg.Wait()
	fmt.Println("执行完毕")
}

func contextWithTimeoutAutoCancel(ctxBackground context.Context) {
	tick := 5
	log.Println("程序执行", tick, "秒后将超时自动退出")
	ctx, cancel := context.WithTimeout(ctxBackground, time.Duration(tick)*time.Second)
	defer cancel()
	for range time.Tick(time.Second) {
		select {
		case <-ctx.Done():
			log.Println("超时时间到, 收到done信号: ", ctx.Err())
			return // todo 不写return, 当前goroutine不会退出
		default:
			log.Println()
		}
	}
}

func contextWithTimeoutHandCancel(ctxBackground context.Context) {
	ctx, cancel := context.WithTimeout(ctxBackground, time.Second*6)

	for i := 0; i <= 10; i++ {
		time.Sleep(time.Second * 3)

		select {
		case <-ctx.Done():
			log.Println("由于手动cancel, 提前收到done信号: ", ctx.Err())
			return // todo 这里写不写return有啥区别??

		default:
			log.Println("curr times : ", i)
			cancel()
		}
	}

}

func main() {

	ctxBackground := context.Background()

	ctxTodo := context.TODO()

	fmt.Println(ctxBackground, ctxTodo, "两者互为别名, 区别不大")

	/**
	withCancel的应用场景, 适用于在多个goroutine同时工作的时候, 由他们的父协程来控制取消
	*/
	contextWithCancel(ctxBackground)

	/**
	withValue的应用场景, 适用于上下文之间传递参数,
	但不要传递关键参数, 一般就是签名, logID 之类的参数,
	传递的数据, 键 和 值 都是interface{}类型, 类型断言时, 记得保证程序的健壮性
	*/
	contextWithValue(ctxBackground)

	/**
	withTimeout的应用场景, 适用于超时之后自动取消当前正在执行的某些操作
	*/
	contextWithTimeoutAutoCancel(ctxBackground)

	/**
	withTimeout的应用场景, 适用于未过超时时间, 手动取消执行
	*/
	contextWithTimeoutHandCancel(ctxBackground)

	fmt.Println("main exit")
}
