package main

import "fmt"

// demoVariables 演示变量、常量、类型推导和零值。
func demoVariables() {
	// Go 是静态类型语言，但支持类型推导。
	// var name string = "Go" 写全类型更显式；:= 更常用于函数内部的短变量声明。
	var language string = "Go"
	version := 1.22

	// 未显式初始化的变量会获得“零值”。
	// 这是 Go 简化工程代码的重要设计：避免未初始化内存导致的不可预测行为。
	var count int       // int 零值为 0
	var enabled bool    // bool 零值为 false
	var nickname string // string 零值为 ""

	const company = "Backend Team"

	fmt.Printf("language=%s version=%.2f company=%s\n", language, version, company)
	fmt.Printf("zero values: count=%d enabled=%t nickname=%q\n", count, enabled, nickname)

	// Go 不允许声明后不使用变量。
	// 这体现了 Go 对工程洁净度的坚持：无用代码应该尽早暴露，而不是沉积在项目里。
}

// demoControlFlow 演示 if、for、switch。
func demoControlFlow() {
	score := 88

	// if 的条件不需要括号，但代码块的大括号必须存在。
	if score >= 90 {
		fmt.Println("grade=A")
	} else if score >= 80 {
		fmt.Println("grade=B")
	} else {
		fmt.Println("grade=C")
	}

	// Go 只有 for，没有 while。通过不同写法覆盖传统 for/while/无限循环。
	sum := 0
	for i := 1; i <= 5; i++ {
		sum += i
	}
	fmt.Println("sum 1..5 =", sum)

	// switch 默认自带 break，不会像 C/Java 那样自动 fallthrough。
	// 这减少了隐式控制流带来的 bug。
	switch day := "Monday"; day {
	case "Monday":
		fmt.Println("start of week")
	case "Friday":
		fmt.Println("almost weekend")
	default:
		fmt.Println("normal day")
	}
}
