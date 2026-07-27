package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type QueueMaster[T any] struct {
	dataCh     chan T
	ctx        context.Context    // 内聚 ctx，让组件自身掌控生命周期
	cancel     context.CancelFunc // 内聚 cancel，避免锁内回调
	producerWg sync.WaitGroup     // 追踪活跃生产者
	closed     bool
	mu         sync.Mutex // 保护 closed 标记和 AddProducer
}

// NewQueueMaster 创建消息队列，将根 context 传入
func NewQueueMaster[T any](parentCtx context.Context, bufSize int) *QueueMaster[T] {
	ctx, cancel := context.WithCancel(parentCtx)
	return &QueueMaster[T]{
		dataCh: make(chan T, bufSize),
		ctx:    ctx,
		cancel: cancel,
	}
}

// Send 完美版：单层 select 彻底免疫死锁
func (q *QueueMaster[T]) Send(val T) bool {
	select {
	case <-q.ctx.Done(): // 收到停止信号
		return false
	case q.dataCh <- val: // 没满直接写；满了在此阻塞，若此时 ctx 取消也能瞬间逃生
		return true
	}
}

func (q *QueueMaster[T]) Receive() (T, bool) {
	val, ok := <-q.dataCh
	return val, ok
}

// AddProducer 核心改进：通过返回值告知外部是否允许添加，彻底杜绝 WaitGroup 滥用 Panic
func (q *QueueMaster[T]) AddProducer() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return false // 队列正准备下线，拒绝新增生产者
	}
	q.producerWg.Add(1)
	return true
}

func (q *QueueMaster[T]) ProducerDone() {
	q.producerWg.Done()
}

// channelDown 清爽版：无需外部传参，锁内只做纯粹的状态变更
func (q *QueueMaster[T]) channelDown() {
	q.mu.Lock()
	if q.closed {
		q.mu.Unlock()
		return
	}
	q.closed = true
	// 锁内只改变信号，不执行外部未知闭包，出锁后安全
	q.mu.Unlock()
	q.cancel()

	// 1. 阻塞等待：所有生产者安全退出
	q.producerWg.Wait()

	// 2. 安全关闭通道
	close(q.dataCh)
}

// ============ 测试演示 ============
func main() {
	// 创建队列，传入根 Context
	queue := NewQueueMaster[string](context.Background(), 4)
	var consumerWg sync.WaitGroup

	// 启动 3 个生产者
	producerNum := 3
	for i := 0; i < producerNum; i++ {
		// 必须先判断是否允许添加生产者
		if !queue.AddProducer() {
			continue
		}
		go func(pid int) {
			defer queue.ProducerDone()
			for j := 111; ; j++ {
				msg := fmt.Sprintf("生产者: %d, 消息: %d", pid, j)
				// 发送时无需再重复传入 ctx
				ok := queue.Send(msg)
				if !ok {
					log.Printf("生产者: %d, 收到停止信号，不再写入\n", pid)
					return
				}
				log.Println(msg)
				time.Sleep(2 * time.Millisecond)
			}
		}(i)
	}

	// 启动 5 个慢消费者
	consumerNum := 5
	for i := 0; i < consumerNum; i++ {
		consumerWg.Add(1)
		go func(cid int) {
			defer consumerWg.Done()
			for {
				val, ok := queue.Receive()
				if !ok {
					log.Printf("消费者: %d, channel已关闭，没有更多数据，退出\n", cid)
					return
				}
				log.Printf("消费者: %d, 收到消息: %s\n", cid, val)
				time.Sleep(2 * time.Second)
			}
		}(i)
	}

	time.Sleep(3 * time.Second)
	fmt.Println("==== 调用 channelDown 开始优雅关闭 ====")
	queue.channelDown() // 干净利落的关闭

	consumerWg.Wait()
	fmt.Println("所有消费者处理完毕，队列完全下线")
}
