# Day 3: 面向对象，结构体 Struct、方法与指针

本章对应文件：

- `oop.go`

## 学习目标

- 使用 `struct` 建模业务对象。
- 使用方法给类型绑定行为。
- 区分值接收者和指针接收者。
- 理解 Go 的组合优于继承。

## 运行方式

```powershell
go run ./day03-oop
```

## 知识点讲解

Go 没有传统意义上的 `class`，也没有继承关键字。它通过 `struct + method + interface` 构建面向对象能力。

结构体负责表达数据：

```go
type User struct {
    ID   int
    Name string
    Age  int
}
```

方法负责表达行为：

```go
func (u *User) Rename(name string) {
    u.Name = name
}
```

方法接收者本质上也是函数参数。值接收者会拷贝对象，适合只读、小对象、不可变语义。指针接收者传递地址，适合修改对象、避免大对象拷贝、保持方法集合一致。

Go 鼓励组合。`Admin` 通过匿名字段组合 `User` 和 `AuditInfo`，获得它们的字段和方法。这种方式比继承层级更扁平，更容易测试和重构。

## 重点代码

- `User`：基础业务实体。
- `Rename()`：指针接收者，修改原对象。
- `IsAdult()`：值接收者，只读判断。
- `Admin`：通过匿名字段组合能力。

## 动手练习

1. 给 `User` 增加 `Email` 字段。
2. 增加 `ChangeAge(age int)` 方法，并思考是否应该用指针接收者。
3. 新增一个 `Manager` 类型，组合 `User`，并添加 `Department` 字段。

## 复盘问题

- Go 没有 class，为什么仍然能写面向对象代码？
- 值接收者和指针接收者的选择标准是什么？
- 组合相比继承有什么工程优势？
