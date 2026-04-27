# Day 6: 并发进阶，sync 包、Context 控制与单元测试

本章对应文件：

- `advanced_sync.go`
- `user_test.go`

## 学习目标

- 使用 `sync.Mutex` 保护共享数据。
- 使用 `sync.WaitGroup` 等待多个 goroutine 完成。
- 使用 `context.Context` 控制取消和超时。
- 使用 `go test` 编写并运行单元测试。

## 运行方式

运行示例：

```powershell
go run ./day06-advanced-sync
```

运行测试：

```powershell
go test ./day06-advanced-sync
```

查看详细测试输出：

```powershell
go test -v ./day06-advanced-sync
```

## 知识点讲解

当多个 goroutine 访问同一份共享状态时，需要同步机制。`sync.Mutex` 是互斥锁，同一时间只允许一个 goroutine 进入临界区。

```go
c.mu.Lock()
defer c.mu.Unlock()
c.value++
```

`sync.WaitGroup` 用来等待一组 goroutine 完成。常见模式是：

```go
wg.Add(1)
go func() {
    defer wg.Done()
}()
wg.Wait()
```

`context.Context` 是后端开发里非常重要的请求生命周期控制工具。它可以传递取消信号、超时时间和少量请求级元数据。HTTP 请求、RPC 调用、数据库查询通常都应该接收 context。

测试文件以 `_test.go` 结尾，测试函数形如 `func TestXxx(t *testing.T)`。Go 的测试工具内置在标准工具链中，不需要额外框架。

## 重点代码

- `SafeCounter`：使用 Mutex 保护计数器。
- `CountConcurrently()`：使用 WaitGroup 并发计数。
- `FetchWithContext()`：使用 Context 响应超时取消。
- `TestCountConcurrently()`：验证并发计数正确性。
- `TestFetchWithContextTimeout()`：验证超时错误。

## 动手练习

1. 临时去掉 `SafeCounter` 中的锁，然后运行 `go test -race ./day06-advanced-sync` 观察数据竞争。
2. 把 worker 数量改成 100，观察结果是否稳定。
3. 给 `FetchWithContext` 增加一个成功场景测试。

## 复盘问题

- 什么是临界区？
- Mutex 和 channel 都能做同步，如何选择？
- WaitGroup 的 `Add` 为什么通常要在 goroutine 外调用？
- Context 为什么不应该存放大量业务数据？
- 单元测试应该测试行为还是实现细节？
