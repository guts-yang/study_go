# Day 2: 切片 Slice 与映射 Map 的底层踩坑指南

本章对应文件：

- `collections.go`

## 学习目标

- 区分数组和切片。
- 理解 slice 的底层结构：指针、长度、容量。
- 理解 `append` 触发扩容时会发生什么。
- 掌握 map 的常见坑：nil map、随机遍历顺序、并发读写风险。

## 运行方式

```powershell
go run ./day02-collections
```

## 知识点讲解

数组是定长值类型，例如 `[3]int`。切片是更常用的动态视图，可以理解成一个描述底层数组片段的结构：

```text
slice = pointer + len + cap
```

- `pointer` 指向底层数组的某个位置。
- `len` 表示当前可访问元素数量。
- `cap` 表示从起点到底层数组末尾的容量。

当 `append` 后长度没有超过容量时，新元素会写入原底层数组。当容量不够时，Go runtime 会分配更大的底层数组，拷贝旧元素，再返回新的 slice。

这会带来一个经典坑：多个 slice 可能共享同一个底层数组，修改一个 slice 可能影响另一个 slice。

map 是哈希表，常用于 key-value 查询。它的零值是 `nil`，可以读取但不能写入。写入前需要 `make(map[K]V)`。

普通 map 不支持并发读写。如果多个 goroutine 同时写 map，程序可能直接崩溃。真实后端项目中，常见处理方式是加 `sync.Mutex`，或者在特定场景下使用 `sync.Map`。

## 重点代码

- `demoSliceGrowth()`：观察 `len` 和 `cap` 的变化。
- `demoSliceSharing()`：演示共享底层数组导致的数据覆盖。
- `demoMapPitfalls()`：演示 nil map、`ok` 判断、遍历顺序。

## 动手练习

1. 把 `make([]int, 0, 2)` 改成 `make([]int, 0, 10)`，观察扩容次数。
2. 在 `demoSliceSharing()` 中使用 `copy` 创建独立切片。
3. 写一个 `map[string][]string`，模拟“城市 -> 用户列表”。

## 复盘问题

- slice 为什么不是简单的动态数组？
- `len` 和 `cap` 的区别是什么？
- 为什么 `append` 后一定要接收返回值？
- 如何区分 map 中 key 不存在和值为零值？
- map 并发读写为什么危险？
