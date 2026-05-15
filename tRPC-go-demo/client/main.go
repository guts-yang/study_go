// Package main 是 tRPC-Go demo 客户端入口。
//
// 它演示了"无配置启动客户端"的最简范式：不需要 trpc_go.yaml，
// 仅靠几个 client.Option 就能直连下游：
//
//   client.WithTarget("ip://127.0.0.1:8000")  // 直连模式：ip://host:port
//   client.WithProtocol("trpc")               // 业务协议
//   client.WithTimeout(2 * time.Second)       // 整体超时
//
// 桩代码内部会自动设置 SerializationType=JSON（值=2），与服务端 trpc_go.yaml 一致。
//
// 调用顺序：
//   1. CreateUser("Alice")  → 期望返回 ID=1
//   2. CreateUser("Bob")    → 期望返回 ID=2
//   3. ListUser()           → 期望返回 2 条
//   4. GetUser(1)           → 期望返回 Alice
//   5. GetUser(999)         → 期望返回错误码 404（演示 errs 错误码透传）
package main

import (
	"context"
	"fmt"
	"time"

	// 必须 blank-import trpc-go 根包，触发 init() 向 codec 注册表注入
	// "trpc"、"http" 等内置协议的 framerBuilder / codec / serializer。
	// 如果不 import，client.Invoke 会报 "client: codec empty"。
	_ "trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/client"
	"trpc.group/trpc-go/trpc-go/errs"
	"trpc.group/trpc-go/trpc-go/log"

	"trpc-go-demo/api/user"
)

func main() {
	// 构造客户端代理：所有调用共享这些 options。
	proxy := user.NewUserClientProxy(
		client.WithTarget("ip://127.0.0.1:8000"),
		client.WithProtocol("trpc"),
		client.WithTimeout(2*time.Second),
	)

	ctx := context.Background()

	fmt.Println("==========================================================")
	fmt.Println("  tRPC-Go demo client → trpc.demo.user.User")
	fmt.Println("==========================================================")

	// 1) CreateUser Alice
	mustCreate(ctx, proxy, "Alice")

	// 2) CreateUser Bob
	mustCreate(ctx, proxy, "Bob")

	// 3) ListUser
	listUsers(ctx, proxy)

	// 4) GetUser(1)
	getUser(ctx, proxy, 1)

	// 5) GetUser(999) —— 演示错误码
	getUser(ctx, proxy, 999)

	fmt.Println("==========================================================")
	fmt.Println("  done.")
}

// mustCreate 创建用户，失败直接 panic（demo 简化错误处理）。
func mustCreate(ctx context.Context, proxy user.UserClientProxy, name string) {
	start := time.Now()
	rsp, err := proxy.CreateUser(ctx, &user.CreateUserReq{Name: name})
	elapsed := time.Since(start)
	if err != nil {
		log.Fatalf("[CreateUser %q] err: %v", name, err)
	}
	fmt.Printf("[CreateUser] req={name:%q} → rsp.user={id:%d, name:%q}  cost=%s\n",
		name, rsp.User.ID, rsp.User.Name, elapsed)
}

// listUsers 列举所有用户。
func listUsers(ctx context.Context, proxy user.UserClientProxy) {
	start := time.Now()
	rsp, err := proxy.ListUser(ctx, &user.ListUserReq{})
	elapsed := time.Since(start)
	if err != nil {
		log.Fatalf("[ListUser] err: %v", err)
	}
	fmt.Printf("[ListUser]   total=%d  cost=%s\n", rsp.Total, elapsed)
	for _, u := range rsp.Users {
		fmt.Printf("             - id=%d name=%q\n", u.ID, u.Name)
	}
}

// getUser 查询单个用户，区分成功 / 错误码两种情况打印。
func getUser(ctx context.Context, proxy user.UserClientProxy, id int64) {
	start := time.Now()
	rsp, err := proxy.GetUser(ctx, &user.GetUserReq{ID: id})
	elapsed := time.Since(start)
	if err != nil {
		// errs.Code(err) 取出业务错误码；errs.Msg(err) 取出错误信息。
		fmt.Printf("[GetUser]    req={id:%d} → ERROR code=%d msg=%q  cost=%s\n",
			id, errs.Code(err), errs.Msg(err), elapsed)
		return
	}
	fmt.Printf("[GetUser]    req={id:%d} → rsp.user={id:%d, name:%q}  cost=%s\n",
		id, rsp.User.ID, rsp.User.Name, elapsed)
}
