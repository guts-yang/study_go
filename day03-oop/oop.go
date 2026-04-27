package main

import "fmt"

// Day 3 目标：
// 1. 理解 Go 没有 class，但用 struct + method 组合出面向对象能力。
// 2. 区分值接收者和指针接收者。
// 3. 理解组合优于继承的 Go 风格。

type User struct {
	ID   int
	Name string
	Age  int
}

// Rename 使用指针接收者。
// 原理：方法接收者本质上也是参数。指针接收者传入的是地址，可以修改原对象。
func (u *User) Rename(name string) {
	u.Name = name
}

// IsAdult 使用值接收者。
// 原理：值接收者会拷贝一份 User，适合只读、小对象、不可变语义的方法。
func (u User) IsAdult() bool {
	return u.Age >= 18
}

type AuditInfo struct {
	CreatedBy string
}

type Admin struct {
	User      // 匿名字段，形成组合。Admin 会“提升”User 的字段和方法。
	AuditInfo // 继续组合审计信息。
	Level     int
}

func main() {
	user := User{ID: 1, Name: "Alice", Age: 20}
	fmt.Printf("before rename: %+v\n", user)

	user.Rename("Alice Zhang")
	fmt.Printf("after rename: %+v adult=%t\n", user, user.IsAdult())

	admin := Admin{
		User:      User{ID: 2, Name: "Root", Age: 30},
		AuditInfo: AuditInfo{CreatedBy: "system"},
		Level:     10,
	}

	// 由于 User 是匿名字段，可以直接访问 admin.Name，也可以显式写 admin.User.Name。
	admin.Rename("Super Admin")
	fmt.Printf("admin name=%s level=%d createdBy=%s\n", admin.Name, admin.Level, admin.CreatedBy)

	// Go 的设计哲学：
	// - 不提供传统继承，避免复杂的层级关系。
	// - 鼓励通过小结构体组合能力，让类型关系更扁平、更容易测试和重构。
}
