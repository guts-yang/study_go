# Day 5: 并发基石，Goroutine 与 Channel 初探

本章对应文件：

- `concurrency.go`

## 学习目标

- 理解 goroutine 是 Go 的轻量级并发单元。
- 理解无缓冲 channel 和有缓冲 channel。
- 掌握 `close`、`range channel`、`select`、超时控制。
- 建立并发程序的阻塞直觉。

## 运行方式

```powershell
go run ./day05-concurrency
```

## 知识点讲解

goroutine 是由 Go runtime 调度的轻量级执行单元。创建方式很简单：

```go
go func() {
    // 并发执行
}()
```

它不是操作系统线程。Go runtime 会把大量 goroutine 调度到较少的 OS thread 上执行，这也是 Go 适合高并发服务的关键原因之一。

channel 是 goroutine 之间通信的管道。无缓冲 channel 的发送和接收必须同时准备好，否则会阻塞：

```go
ch := make(chan string)
ch <- "hello"
msg := <-ch
```

有缓冲 channel 可以暂存一定数量的元素。缓冲区没满时，发送不阻塞；缓冲区为空时，接收阻塞。

`select` 可以同时等待多个 channel 操作，常用于超时、取消和多路复用。

## 重点代码

- `demoGoroutineAndChannel()`：无缓冲 channel 的同步语义。
- `demoBufferedChannel()`：缓冲 channel、关闭和遍历。
- `demoSelectTimeout()`：使用 `select` 和 `time.After` 控制超时。

## 动手练习

1. 把超时时间从 `100ms` 改为 `500ms`，观察输出变化。
2. 写两个 goroutine 同时向同一个 channel 发送数据。
3. 尝试向已关闭 channel 发送数据，观察 panic。

## 复盘问题

- goroutine 和线程有什么区别？
- 无缓冲 channel 为什么能实现同步？
- close channel 的真正含义是什么？
- `select` 在后端服务中常用于哪些场景？
