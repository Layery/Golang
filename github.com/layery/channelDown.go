package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type Queue struct {
	dataCh     chan string
	producerWg *sync.WaitGroup
	consumerWg *sync.WaitGroup
	mu         sync.Mutex
	ctx        context.Context
	cancel     func()
	isClosed   bool
}

func NewQueue(size int) *Queue {
	ctx, cancel := context.WithCancel(context.Background())
	return &Queue{
		dataCh:     make(chan string, size),
		ctx:        ctx,
		cancel:     cancel,
		producerWg: &sync.WaitGroup{},
		consumerWg: &sync.WaitGroup{},
	}
}

func (q *Queue) Send(data string) bool {
	select {
	case q.dataCh <- data:
		return true
	case <-q.ctx.Done(): // 收到停止信号
		return false
	}
}

func (q *Queue) Receive() (string, bool) {
	val, ok := <-q.dataCh
	return val, ok
}

func (q *Queue) Down() {
	q.mu.Lock()
	if q.isClosed {
		q.mu.Unlock()
		return
	}
	q.isClosed = true
	q.mu.Unlock()

	q.cancel()          // 发送关闭信号, 不再生产
	q.producerWg.Wait() // 等待全部都生产完
	close(q.dataCh)
}

func main() {
	queue := NewQueue(3)

	producerNumber := 8
	for i := 0; i < producerNumber; i++ {
		if queue.isClosed {
			log.Printf("队列已关闭, 不再创建生产者\n")
			return
		}
		queue.producerWg.Add(1)
		go func(i int) {
			defer queue.producerWg.Done()
			defer func() {
				if err := recover(); err != nil { // todo 生产者挂了, 尝试重新生产
					// todo
				}
			}()
			for j := 100; ; j++ {
				msg := fmt.Sprintf("producer: %d, msg: %d", i, j)
				res := queue.Send(msg)
				if !res {
					log.Printf("生产者: %d, 收到停止信号\n", i)
					return
				}
				log.Printf("生产者: %d 写入消息%s\n", i, msg)
			}
		}(i)
	}

	consumerNumber := 2
	for i := 0; i < consumerNumber; i++ {
		queue.consumerWg.Add(1)
		go func(i int) {
			defer queue.consumerWg.Done()
			for {
				val, ok := queue.Receive()
				if ok {
					log.Printf("消费者: %d, 接收消息: %v\n", i, val)
				} else {
					log.Printf("队列已关闭")
					return
				}
				time.Sleep(1 * time.Second)
			}
		}(i)
	}

	time.Sleep(1 * time.Second)

	queue.Down()

	queue.consumerWg.Wait()

	log.Printf("end")
}
