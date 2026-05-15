// Package main 是 day02-userservice-demo 的客户端演示。
package main

import (
	"context"
	"fmt"
	"time"

	pb "day02-userservice-demo/stub/trpc/study/user"

	_ "git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/client"
)

func main() {
	proxy := pb.NewUserServiceClientProxy(
		client.WithTarget("ip://127.0.0.1:8001"),
		client.WithTimeout(time.Second),
	)

	ctx := context.Background()

	// 1. 创建用户
	crsp, err := proxy.CreateUser(ctx, &pb.CreateUserReq{Name: "Alice"})
	if err != nil {
		panic(fmt.Sprintf("CreateUser failed: %v", err))
	}
	fmt.Printf("CreateUser → id=%d name=%s\n", crsp.GetUser().GetId(), crsp.GetUser().GetName())

	// 2. 查询用户
	grsp, err := proxy.GetUser(ctx, &pb.GetUserReq{Id: crsp.GetUser().GetId()})
	if err != nil {
		panic(fmt.Sprintf("GetUser failed: %v", err))
	}
	fmt.Printf("GetUser    → id=%d name=%s\n", grsp.GetUser().GetId(), grsp.GetUser().GetName())

	// 3. 查询不存在的用户
	_, err = proxy.GetUser(ctx, &pb.GetUserReq{Id: 9999})
	if err != nil {
		fmt.Printf("GetUser 9999 → expected error: %v\n", err)
	}
}
