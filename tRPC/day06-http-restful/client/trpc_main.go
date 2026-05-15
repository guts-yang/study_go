// Package main 用 tRPC 协议直接调（端口 8003）。
//
// 命令：go run .\client\ -mode trpc
// 注意：client 目录下有两个 main 文件，靠 build tag 控制只编一个。
//go:build trpc

package main

import (
	"context"
	"fmt"
	"time"

	pb "day06-http-restful/stub/trpc/study/user"

	_ "trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/client"
)

func main() {
	proxy := pb.NewUserServiceClientProxy(
		client.WithTarget("ip://127.0.0.1:8003"),
		client.WithTimeout(time.Second),
	)
	rsp, err := proxy.CreateUser(context.Background(), &pb.CreateUserReq{Name: "Alice"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("[trpc] CreateUser id=%d\n", rsp.GetUser().GetId())
}
