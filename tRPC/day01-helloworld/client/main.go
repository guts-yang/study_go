// Package main 是 day01-helloworld 的客户端入口。
//
// 客户端不需要 trpc_go.yaml —— 所有配置都可以通过 client.With* Options 直接传入。
// 实际项目中也常常把 callee 配置写到 yaml 的 client 段里，等 day05 再展开。
package main

import (
	"context"
	"fmt"
	"time"

	pb "day01-helloworld/stub/trpc/helloworld"

	// 导入 trpc-go 触发其 init() 完成全局插件注册。
	_ "git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/client"
)

func main() {
	// NewGreeterClientProxy 来自桩代码，返回的 proxy 是并发安全的。
	// 在真实项目中通常做成 package 级单例，避免每次调用都新建。
	proxy := pb.NewGreeterClientProxy(
		// ip:// 表示直连，无需任何名字服务；适合本地调试。
		client.WithTarget("ip://127.0.0.1:8000"),
		// 单次 RPC 整体超时：包括连接、序列化、网络、服务端处理、反序列化。
		client.WithTimeout(time.Second),
	)

	rsp, err := proxy.SayHello(context.Background(), &pb.HelloRequest{Msg: "world"})
	if err != nil {
		// tRPC 的 err 通常是 *errs.Error，里面带 code + msg；day04 会展开错误码体系。
		panic(err)
	}
	fmt.Println(rsp.GetMsg()) // 期望：Hello, world
}
