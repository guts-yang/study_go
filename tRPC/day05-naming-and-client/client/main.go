// Package main 是 day05 的演示客户端。
//
// 流程：
//  1. 直接调 User 服务，先 CreateUser 一个 Alice；
//  2. 再调 Gateway 的 GreetUser，让 Gateway 内部去问 User 拿名字，拼出问候语。
package main

import (
	"context"
	"fmt"
	"time"

	pb "day05-naming-and-client/stub/trpc/study/user"

	_ "git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/client"
)

func main() {
	userCli := pb.NewUserServiceClientProxy(
		client.WithTarget("ip://127.0.0.1:8001"),
		client.WithTimeout(time.Second),
	)
	gatewayCli := pb.NewGatewayServiceClientProxy(
		client.WithTarget("ip://127.0.0.1:8002"),
		client.WithTimeout(time.Second),
	)

	ctx := context.Background()

	crsp, err := userCli.CreateUser(ctx, &pb.CreateUserReq{Name: "Alice"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("CreateUser via User    → id=%d\n", crsp.GetUser().GetId())

	grsp, err := gatewayCli.GreetUser(ctx, &pb.GreetUserReq{Id: crsp.GetUser().GetId()})
	if err != nil {
		panic(err)
	}
	fmt.Printf("GreetUser via Gateway  → %q\n", grsp.GetGreeting())
}
