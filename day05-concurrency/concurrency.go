package main

import (
	"fmt"
	"time"
)

// Day 5 目标：
// 1. 理解 goroutine 是 Go 的轻量级并发执行单元。
// 2. 理解 channel 的阻塞语义：发送和接收如何同步。
// 3. 掌握 close、range channel、select 的基本用法。

func main() {
	demoGoroutineAndChannel()
	demoBufferedChannel()
	demoSelectTimeout()
	demodoubleroutine()
}

func demoGoroutineAndChannel() {
	fmt.Println("== unbuffered channel ==")

	ch := make(chan string)

	go func() {
		// goroutine 由 Go runtime 调度，不等同于操作系统线程。
		// Go runtime 会把大量 goroutine 多路复用到较少的 OS thread 上。
		time.Sleep(200 * time.Millisecond)
		ch <- "job done"
		// 无缓冲 channel 的发送会阻塞，直到另一个 goroutine 准备接收。
	}()

	msg := <-ch
	fmt.Println("receive:", msg)
}

func demoBufferedChannel() {
	fmt.Println("\n== buffered channel + close ==")

	jobs := make(chan int, 3)

	for i := 1; i <= 3; i++ {
		// 有缓冲 channel 在缓冲区未满时发送不阻塞。
		jobs <- i
	}
	close(jobs)

	// close 表示不会再发送新值。range 会持续接收，直到 channel 被关闭且缓冲区被读空。
	for job := range jobs {
		fmt.Println("handle job:", job)
	}

	// 注意：
	// - 向已关闭 channel 发送会 panic。
	// - 从已关闭 channel 接收会立刻返回零值和 ok=false。
	value, ok := <-jobs
	fmt.Printf("after close: value=%d ok=%t\n", value, ok)
}

func demoSelectTimeout() {
	fmt.Println("\n== select timeout ==")

	result := make(chan string)

	go func() {
		time.Sleep(300 * time.Millisecond)
		result <- "slow result"
	}()

	select {
	case v := <-result:
		fmt.Println("got:", v)
	case <-time.After(400 * time.Millisecond):
		// select 会等待多个 channel 操作中最先就绪的一个。
		// time.After 返回一个只读 channel，到时间后发送当前时间，常用于超时控制。
		fmt.Println("timeout")
	}
}

func demodoubleroutine() {
	ch := make(chan int)

	go func() {
		ch <- 1
	}()
	go func() {
		ch <- 2
	}()
	msg := <-ch
	fmt.Println("receive:", msg)
}
