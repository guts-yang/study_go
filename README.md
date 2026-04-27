# Go 后端 7 天核心工程实践

这个项目是一套面向后端开发面试和工程入门的 Go 语言学习路径。它按 7 天拆分，每天一个独立目录，每个目录都是可运行的小项目，帮助你从语法、集合、面向对象、接口、并发、测试一路推进到 `net/http` Web 服务。

## 目录结构

```text
study_go/
├── go.mod
├── day01-basics/
│   ├── README.md
│   ├── main.go
│   └── variables.go
├── day02-collections/
│   ├── README.md
│   └── collections.go
├── day03-oop/
│   ├── README.md
│   └── oop.go
├── day04-interface-error/
│   ├── README.md
│   └── interface_err.go
├── day05-concurrency/
│   ├── README.md
│   └── concurrency.go
├── day06-advanced-sync/
│   ├── README.md
│   ├── advanced_sync.go
│   └── user_test.go
└── day07-web-server/
    ├── README.md
    └── server.go
```

## 环境准备

建议使用 Go 1.22 或更新版本，因为第 7 天的 `net/http` 路由示例使用了 Go 1.22 标准库的路径参数能力。

检查版本：

```powershell
go version
```

验证整个项目：

```powershell
go test ./...
```

## 推荐学习方式

每天按这个节奏学习：

1. 先阅读当天目录的 `README.md`，了解知识点和面试重点。
2. 再阅读 `.go` 文件中的中文注释，把注释当成讲义。
3. 执行 `go run ./dayXX-xxx` 观察输出。
4. 修改代码做当天 README 里的练习。
5. 用自己的话复述“为什么 Go 要这样设计”。

不要只看代码。Go 的难点不是语法复杂，而是工程边界、并发模型、错误处理和标准库风格。

## 7 天路线

| 天数 | 目录 | 主题 | 核心目标 |
| --- | --- | --- | --- |
| Day 1 | `day01-basics` | 基础语法与 Go Modules | 能运行 Go 项目，理解变量、控制流、包和 module |
| Day 2 | `day02-collections` | Slice 与 Map | 理解 slice 底层数组、扩容、共享，以及 map 常见坑 |
| Day 3 | `day03-oop` | Struct、方法与指针 | 掌握 Go 的组合式面向对象 |
| Day 4 | `day04-interface-error` | Interface 与错误处理 | 理解隐式接口、多态、error as value |
| Day 5 | `day05-concurrency` | Goroutine 与 Channel | 理解轻量并发和 channel 阻塞语义 |
| Day 6 | `day06-advanced-sync` | sync、Context、测试 | 掌握 Mutex、WaitGroup、Context 和 go test |
| Day 7 | `day07-web-server` | net/http Web 服务 | 写出一个并发安全的极简后端服务 |

## 常用命令

运行某一天：

```powershell
go run ./day01-basics
```

运行所有测试：

```powershell
go test ./...
```

格式化代码：

```powershell
gofmt -w .
```

启动第 7 天服务：

```powershell
go run ./day07-web-server
```

请求服务：

```powershell
curl http://localhost:8080/health
curl http://localhost:8080/users
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d "{\"name\":\"Charlie\"}"
curl http://localhost:8080/users/1
```

## 面试复盘清单

学完后你应该能回答：

- `go.mod` 的作用是什么？module path 是什么？
- Go 的零值机制解决了什么工程问题？
- slice 的 `len` 和 `cap` 有什么区别？
- `append` 什么时候会影响原 slice？
- map 为什么不能并发读写？
- 值接收者和指针接收者如何选择？
- Go 为什么偏好组合而不是继承？
- interface 的隐式实现有什么好处？
- `error` 为什么是普通返回值而不是异常？
- goroutine 和线程有什么区别？
- 无缓冲 channel 为什么会阻塞？
- `sync.Mutex` 和 channel 分别适合什么场景？
- `context.Context` 在后端服务里解决什么问题？
- `net/http` handler 为什么要注意并发安全？

## 学习建议

这一周的目标不是背语法，而是形成后端工程直觉：数据结构怎么传递，并发状态怎么保护，失败路径怎么表达，请求生命周期怎么取消。Go 的美感在于少量概念反复组合，写得越多，你越能感受到它对大型工程的克制和锋利。
