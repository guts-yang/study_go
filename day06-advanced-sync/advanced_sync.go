package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Day 6 目标：
// 1. 使用 sync.Mutex 保护共享状态。
// 2. 使用 sync.WaitGroup 等待一组 goroutine 完成。
// 3. 使用 context 传递取消信号和截止时间。
// 4. 编写可测试的小函数，理解 go test 的基本方式。

type SafeCounter struct {
	mu    sync.Mutex
	value int
}

func (c *SafeCounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func CountConcurrently(workers int, timesPerWorker int) int {
	var wg sync.WaitGroup
	counter := &SafeCounter{}

	for i := 0; i < workers; i++ {

		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < timesPerWorker; j++ {
				counter.Inc()
			}
		}()
	}

	// WaitGroup 内部维护计数器。Add 增加任务数，Done 减少任务数，Wait 阻塞直到计数器归零。
	wg.Wait()
	return counter.Value()
}

func FetchWithContext(ctx context.Context) error {
	select {
	case <-time.After(10 * time.Millisecond):
		fmt.Println("fetch success")
		return nil
	case <-ctx.Done():
		// ctx.Done() 返回一个 channel。
		// 当超时、取消或父 context 被取消时，该 channel 会被关闭，从而唤醒这里的接收。
		return ctx.Err()
	}
}

func main() {
	total := CountConcurrently(100, 1000)
	fmt.Println("counter =", total)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	if err := FetchWithContext(ctx); err != nil {
		fmt.Println("fetch canceled:", err)
	}

	// 工程建议：
	// - goroutine 泄漏常见原因是没有退出信号或 channel 永远无人接收。
	// - 后端服务中，HTTP 请求、RPC 调用、数据库访问都应尽量传递 context。
	// - 共享内存需要同步；更 Go 的方式是“不要通过共享内存通信，而要通过通信共享内存”。
}
