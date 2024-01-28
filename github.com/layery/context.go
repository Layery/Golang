package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"
)

func contextWithCancel(ctxBackground context.Context) {
	ctx, cancel := context.WithCancel(ctxBackground)

	defer func() {
		if err := recover(); err != nil {
			log.Println("system err: ", err)
		}
	}()

	i := 0
	for i <= 10 {
		go func(ctx context.Context, index int) {
			for {
				time.Sleep(time.Second)
				select {
				case <-ctx.Done():
					fmt.Println()
					log.Println(strconv.Itoa(index) + "收到停止信号")
					return // todo 这里如果不return, 当前这个子协程将不会自动退出, 它会一直能接收到停止信号
				default:
					log.Println("我是第" + strconv.Itoa(index) + "个协程")
				}
			}
		}(ctx, i)
		i += 1
	}

	time.Sleep(time.Second * 2)
	cancel()

	select {}
}

func contextWithValue(ctxBackground context.Context) {
	data := map[string]interface{}{
		"name": "llf",
		"age":  33,
	}
	ctx := context.WithValue(ctxBackground, "key", data)

	go func(ctx context.Context) {
		log.Println("我是子协程, 我获取到了来自父ctx的参数", ctx.Value("key"))
	}(ctx)
}

func contextWithTimeoutAutoCancel(ctxBackground context.Context) {
	log.Println("程序执行2秒后将超时自动退出")
	ctx, cancel := context.WithTimeout(ctxBackground, time.Second*2)
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
		time.Sleep(time.Second)

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
	go contextWithCancel(ctxBackground)

	/**
	withValue的应用场景, 适用于上下文之间传递参数,
	但不要传递关键参数, 一般就是签名, logID 之类的参数,
	传递的数据, 键 和 值 都是interface{}类型, 类型断言时, 记得保证程序的健壮性
	*/
	go contextWithValue(ctxBackground)

	/**
	withTimeout的应用场景, 适用于超时之后自动取消当前正在执行的某些操作
	*/
	go contextWithTimeoutAutoCancel(ctxBackground)

	/**
	withTimeout的应用场景, 适用于未过超时时间, 手动取消执行
	*/
	go contextWithTimeoutHandCancel(ctxBackground)

	select {}
}
