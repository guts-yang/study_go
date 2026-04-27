# Day 4: Interface 与多态，以及错误处理

本章对应文件：

- `interface_err.go`

## 学习目标

- 理解 interface 是行为契约。
- 理解 Go 的接口是隐式实现。
- 使用接口实现多态。
- 掌握 `error`、自定义错误、错误包装和 `errors.As`。

## 运行方式

```powershell
go run ./day04-interface-error
```

## 知识点讲解

接口描述行为，而不是描述继承关系：

```go
type Notifier interface {
    Notify(message string) error
}
```

任何类型只要实现了 `Notify(string) error`，就自动满足 `Notifier`。不需要写 `implements`。这让业务代码可以依赖抽象行为，而不是依赖具体实现。

Go 的错误处理哲学是 `error as value`。错误是普通返回值，不是隐藏的异常控制流。调用方通过 `if err != nil` 明确处理失败路径。

`fmt.Errorf("xxx: %w", err)` 可以包装错误并保留错误链。`errors.As` 可以从错误链中提取某类错误，适合处理自定义错误类型。

接口还有一个高频坑：接口值底层可以理解为 `type + value`。只有二者都为 nil 时，接口才等于 nil。因此把一个 nil 指针放进接口里，接口本身可能并不等于 nil。

## 重点代码

- `Notifier`：接口定义。
- `EmailNotifier`、`SMSNotifier`：不同实现。
- `ServiceError`：自定义错误类型。
- `sendWelcome()`：依赖接口并返回错误。

## 动手练习

1. 新增一个 `WebhookNotifier`，实现 `Notify` 方法。
2. 让 `sendWelcome` 接收不同 notifier，观察多态效果。
3. 新增一个错误类型 `ValidationError`，并用 `errors.As` 识别它。

## 复盘问题

- Go 的接口为什么是隐式实现？
- 接口适合作为参数还是返回值？
- 为什么 Go 不使用异常作为常规错误处理方式？
- `%w` 和 `%v` 包装错误有什么区别？
