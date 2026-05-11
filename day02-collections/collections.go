package main

import "fmt"

// Day 2 目标：
// 1. 掌握数组、切片、映射的基本用法。
// 2. 理解 slice 的三元组结构：指向底层数组的指针、长度 len、容量 cap。
// 3. 理解 map 是引用类型，以及并发读写 map 的风险。
func main() {
	demoSliceGrowth()
	demoSliceSharing()
	demoMapPitfalls()
}

func demoSliceGrowth() {
	fmt.Println("== slice 扩容观察 ==")

	// make([]int, 0, 2) 创建长度为 0、容量为 2 的切片。
	// 长度 len 表示当前可访问元素数；容量 cap 表示从起点到底层数组末尾还能容纳多少元素。
	nums := make([]int, 0, 10)

	for i := 1; i <= 8; i++ {
		nums = append(nums, i)
		fmt.Printf("append %d -> len=%d cap=%d data=%v\n", i, len(nums), cap(nums), nums)
	}

	// 扩容原理：
	// - 当 append 后长度不超过容量，Go 直接写入原底层数组。
	// - 当容量不足，运行时会分配更大的底层数组，把旧元素拷贝过去，再返回新的 slice。
	// - 扩容倍数不是永远 2 倍；Go 运行时会根据元素大小和当前容量做权衡。
	// 工程建议：
	// - 已知规模时优先 make([]T, 0, n)，减少扩容和拷贝。
}

func demoSliceSharing() {
	fmt.Println("\n== slice 共享底层数组踩坑 ==")

	origin := []string{"A", "B", "C", "D"}
	part := origin[:2] // part 与 origin 共享同一个底层数组。

	part[0] = "X"
	fmt.Println("origin after part[0] = X:", origin)

	// 这里 part 的 cap 仍然覆盖到 origin 后面的空间。
	// append(part, "Y") 可能直接覆盖 origin[2]，因为底层数组还有容量。
	part = append(part, "Y")
	fmt.Println("origin after append(part, Y):", origin)
	fmt.Println("part:", part)

	// 如果希望创建独立副本，使用 copy 或 append 到 nil 切片。
	clone := append([]string(nil), origin...)
	clone[0] = "CLONE"
	fmt.Println("origin unchanged by clone:", origin)
	fmt.Println("clone:", clone)
	clonend := make([]string, len(origin))
	copy(clonend, origin)
	fmt.Println("origin unchanged by clonend:", origin)
	fmt.Println("clonend:", clonend)
}

func demoMapPitfalls() {
	fmt.Println("\n== map 使用与踩坑 ==")

	// map 的零值是 nil。nil map 可以读，不能写。
	var nilMap map[string]int
	fmt.Println("read nilMap:", nilMap["missing"])

	// nilMap["x"] = 1 会 panic，因此写入前必须 make。
	score := make(map[string]int)
	score["alice"] = 95
	score["bob"] = 82

	// 读取 map 时，第二个返回值 ok 用来区分“键不存在”和“值刚好是零值”。
	value, ok := score["carol"]
	fmt.Printf("carol value=%d exists=%t\n", value, ok)

	// 删除不存在的 key 是安全的。
	delete(score, "nobody")

	for name, s := range score {
		// map 遍历顺序是随机的，Go 故意不保证顺序。
		// 这能避免开发者无意依赖哈希表内部实现。
		fmt.Printf("%s => %d\n", name, s)
	}
	city := make(map[string]string)
	city["城市"] = "北京"
	city["城市2"] = "上海"
	cityVal, ok2 := city["城市1"]
	fmt.Printf("value=%s exists=%t\n", cityVal, ok2)
	// 并发提醒：
	// 普通 map 不支持并发读写。多个 goroutine 同时写 map 会导致 fatal error。
	// 工程上可选 sync.Mutex 保护 map，或使用 sync.Map 处理特定高并发读多写少场景。
}
