# Day 1: Go 基础语法与 Go Modules 工程初始化

本章对应文件：

- `main.go`
- `variables.go`

## 学习目标

- 理解 Go 程序入口：`package main` 与 `func main()`。
- 理解同一目录、同一 package 下的多个 `.go` 文件会一起编译。
- 掌握变量、常量、零值、类型推导、`if`、`for`、`switch`。
- 理解 `go.mod` 是 Go Modules 项目的依赖和模块边界。

## 运行方式

在项目根目录执行：

```powershell
go run ./day01-basics
```

也可以进入目录后执行：

```powershell
cd day01-basics
go run .
```

## 知识点讲解

Go 项目通常由一个 `go.mod` 管理。`go.mod` 中的 `module study_go` 表示当前项目的模块名，`go 1.22` 表示该项目使用的 Go 语言版本语义。

Go 的入口文件不要求必须叫 `main.go`，真正重要的是：

```go
package main

func main() {
}
```

一个目录通常对应一个 package。`main.go` 可以调用 `variables.go` 里的函数，是因为它们都声明了 `package main`。

Go 的变量有零值。比如 `int` 的零值是 `0`，`bool` 的零值是 `false`，`string` 的零值是空字符串。这让很多结构体可以直接声明后使用，减少未初始化状态带来的混乱。

Go 只有 `for` 一种循环语句。它可以写成传统三段式，也可以写成类似 `while` 的形式。这是 Go 追求语法简洁的体现。

## 重点代码

- `demoVariables()`：变量、常量、零值、类型推导。
- `demoControlFlow()`：条件判断、循环、分支。

## 动手练习

1. 新增一个 `demoArray()` 函数，声明 `[3]string` 并打印。
    - 注意声明形式
2. 把 `score := 88` 改成不同分数，观察分支输出。
    - 已完成，分支输出符合逻辑
3. 尝试声明一个变量但不使用，观察 Go 编译器报错。
    - .\variables.go:60:6: declared and not used: numbers
    
## 复盘问题

- Go 为什么不允许有未使用变量？
- `var x int` 和 `x := 1` 分别适合什么场景？
- `go.mod` 对一个后端项目意味着什么？
