package main

import (
	"log"
	//"net/http"
	_ "net/http/pprof"
	"time"
)

/**
	timer 在执行timer.Stop之后, 并不会关闭timer.C通道, 以防止其它go程从中读取数据
 */
func testTimer_bak()  {
	i := 0;
	timer := time.NewTimer(time.Second * 10)
	for {
		log.Println(i, "main func")
		time.Sleep(time.Second)
		go func(index int) {
			log.Println(index, "go func...")
			select {
			case <-timer.C:
				log.Println(index, "timer.C 关闭了吗?", timer.Stop())
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
func testTimer()  {
	log.Println("begin")
	timer := time.NewTimer(time.Second * 10)
	ticker := time.NewTicker(time.Second * 2)

	go func() {
		for {
			log.Println("debug")
			_, data := <- ticker.C
			log.Println("子go程", data)
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
			ticker.Stop()
			timer.Stop()
			_, ok := <- timer.C
			log.Println("timer触发后, C 的状态: ", ok)
		}

	}()

	for {
		_, data := <- ticker.C
		log.Println("我是读取ticker的父go程", data)
	}

}
/**
 */
func main() {

	go testTimer_bak()

	go testTimer()

	for {
		select {}
	}

}