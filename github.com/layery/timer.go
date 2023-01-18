package main

import (
	"fmt"
	"log"

	//"net/http"
	_ "net/http/pprof"
	"time"
)

/**
timer 在执行timer.Stop之后, 并不会关闭timer.C通道, 以防止其它go程从中读取数据
*/
func testTimer_bak() {
	i := 0
	timer := time.NewTimer(time.Second * 10)
	for {
		log.Println(i, "main func")
		time.Sleep(time.Second * 5)
		go func(index int) {
			log.Println(index, "go func...")
			select {
			case <-timer.C:
				log.Printf("%#v, %#v \n", index, "timer.C 关闭了吗?", timer.Stop())
			}
			log.Println(index, "go func end ...")
			return
		}(i)
		i += 1
	}
}

/**
timer 在执行timer.Stop之后, 并不会关闭timer.C通道, 以防止其它go程从中读取数据
*/
func testTimer() {
	log.Println("begin")
	timer := time.NewTimer(time.Second * 10)
	ticker := time.NewTicker(time.Second * 1)

	go func() {
		for {
			select {
			case <-ticker.C:
				_, ok := <-timer.C
				log.Println("子go程检测timer.C状态", ok)
			}
		}
	}()

	go func() {

		select {
		/**
		     调用timer.Stop将会暂停触发计时器, 并且返回true, 除非timer到期, 或者已经被stop了
		     将会返回false
			 延时结束, 当timer.C接收到消息时, 说明timer已经过期了,
			 此时在这个case调用timer.Stop返回的就是false, 但是timer.C的通道不会被关闭, 这将导致
			 timer.C阻塞读取, go程一直无法销毁
		*/
		case <-timer.C:
			log.Println("timer is end !")
			timer.Stop()
			fmt.Println("hello")

			_, ok := <-timer.C
			log.Println("timer.stop 触发后, timer.C 的状态: ", ok)
		}

	}()

	for {
		<-ticker.C
	}

}

// 验证timer.stop后, 通道会不会关闭
func checkTimerChan() {
	log.Println("begin")

	var ch = make(chan interface{}, 2)
	var timer *time.Timer = nil
	ticker := time.NewTicker(time.Second)

	go func() {
		select {
		case <-ticker.C:
			log.Printf("当前timer.c状态: %#v", timer.C)
		}
	}()

	go func() {
		timer = time.NewTimer(time.Second * 20) // 延时10秒结束
	}()

	ch <- 1
	go func() {
		select { // select 收到一条消息后, 将不再阻塞, 整个select语句就此结束
		case <-ch:
			log.Println("ch 收到信号, 执行timer.stop()")
			// timer.Stop()
		case <-timer.C:
			log.Println("时间到")
		}
	}()

	go func() {
		time.Sleep(time.Second * 15)

		log.Println("15秒后, 检查timer.C的状态, %#v", timer.C)

	}()
	for {
		fmt.Println("the end")
		time.Sleep(time.Second)
		_, ok := <-timer.C
		log.Printf("ping...%#v", ok)

	}
}

/**
 */
func main() {

	// testTimer_bak()

	// go testTimer()

	go checkTimerChan()
	for {
		select {}
	}

}
