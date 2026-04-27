package main

import "fmt"

// Day 1 目标：
// 1. 理解 Go 程序的最小组成：package、import、func main。
// 2. 理解 Go Modules 是现代 Go 项目的标准工程组织方式。
//
// 工程说明：
// - 当前项目根目录的 go.mod 声明了 module study_go。
// - 每个 dayXX 文件夹都是一个独立的 main package，因此可以在对应目录执行 go run .。
// - Go 的设计哲学之一是“简单、显式、工具友好”：目录即包，go.mod 即依赖边界。
func main() {
	fmt.Println("Day 1: Go 基础语法与 Go Modules")

	// variables.go 中定义的函数。Go 同一目录下、同一 package 的文件会被一起编译。
	demoVariables()
	demoControlFlow()
}
