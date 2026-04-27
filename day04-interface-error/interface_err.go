package main

import (
	"errors"
	"fmt"
)

// Day 4 目标：
// 1. 理解 interface 是行为契约：只关心“能做什么”，不关心“是什么”。
// 2. 理解 Go 的隐式实现：类型拥有接口要求的方法，就自动实现接口。
// 3. 掌握 error as value 的错误处理风格。

type Notifier interface {
	Notify(message string) error
}

type EmailNotifier struct {
	Address string
}

func (e EmailNotifier) Notify(message string) error {
	if e.Address == "" {
		return errors.New("email address is empty")
	}
	fmt.Printf("[email] to=%s msg=%s\n", e.Address, message)
	return nil
}

type SMSNotifier struct {
	Phone string
}

func (s SMSNotifier) Notify(message string) error {
	if s.Phone == "" {
		return errors.New("phone is empty")
	}
	fmt.Printf("[sms] to=%s msg=%s\n", s.Phone, message)
	return nil
}

// ServiceError 是自定义错误类型。
// 只要实现 Error() string 方法，就满足内置 error 接口。
type ServiceError struct {
	Code int
	Msg  string
}

func (e ServiceError) Error() string {
	return fmt.Sprintf("service error code=%d msg=%s", e.Code, e.Msg)
}

func sendWelcome(n Notifier, username string) error {
	if username == "" {
		return ServiceError{Code: 400, Msg: "username is empty"}
	}

	// Go 不使用异常作为常规控制流，而是显式返回 error。
	// 好处：调用方必须面对失败路径，工程上更可控。
	if err := n.Notify("Welcome, " + username); err != nil {
		// %w 用于错误包装，保留原始错误链，便于 errors.Is / errors.As 判断。
		return fmt.Errorf("send welcome failed: %w", err)
	}
	return nil
}

func main() {
	notifiers := []Notifier{
		EmailNotifier{Address: "alice@example.com"},
		SMSNotifier{Phone: "13800000000"},
		EmailNotifier{}, // 故意触发错误。
	}

	for _, n := range notifiers {
		if err := sendWelcome(n, "Alice"); err != nil {
			fmt.Println("error:", err)
		}
	}

	err := sendWelcome(EmailNotifier{Address: "bob@example.com"}, "")
	if err != nil {
		var serviceErr ServiceError
		if errors.As(err, &serviceErr) {
			fmt.Println("matched ServiceError:", serviceErr.Code, serviceErr.Msg)
		} else {
			fmt.Println("normal error:", err)
		}
	}

	// 面试提醒：
	// interface 的底层可粗略理解为 type + value。
	// 一个接口变量只有 type 和 value 都为 nil 时才等于 nil。
	// 这也是“返回 *MyError 给 error 接口后不等于 nil”的常见坑点来源。
}
