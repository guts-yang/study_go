package main

import (
	"context"
	"fmt"
	"time"

	pb "day03-config-and-admin/stub/trpc/study/user"

	_ "git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/client"
)

func main() {
	proxy := pb.NewUserServiceClientProxy(
		client.WithTarget("ip://127.0.0.1:8001"),
		client.WithTimeout(time.Second),
	)
	ctx := context.Background()

	crsp, err := proxy.CreateUser(ctx, &pb.CreateUserReq{Name: "Bob"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("CreateUser → id=%d name=%s\n", crsp.GetUser().GetId(), crsp.GetUser().GetName())

	grsp, err := proxy.GetUser(ctx, &pb.GetUserReq{Id: crsp.GetUser().GetId()})
	if err != nil {
		panic(err)
	}
	fmt.Printf("GetUser    → id=%d name=%s\n", grsp.GetUser().GetId(), grsp.GetUser().GetName())
}
