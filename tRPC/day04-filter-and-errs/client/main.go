// Package main 是 day04-filter-and-errs 的客户端。
package main

import (
	"context"
	"fmt"
	"time"

	pb "day04-filter-and-errs/stub/trpc/study/user"

	_ "git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/client"
	"git.code.oa.com/trpc-go/trpc-go/codec"
	"git.code.oa.com/trpc-go/trpc-go/errs"
)

func main() {
	baseOpts := []client.Option{
		client.WithTarget("ip://127.0.0.1:8001"),
		client.WithTimeout(time.Second),
	}

	// --- 场景 1：携带有效 token ---
	ctx1, msg1 := codec.WithNewMessage(context.Background())
	msg1.WithClientMetaData(codec.MetaData{"token": []byte("valid-token")})

	proxy := pb.NewUserServiceClientProxy(baseOpts...)
	crsp, err := proxy.CreateUser(ctx1, &pb.CreateUserReq{Name: "Charlie"})
	if err != nil {
		fmt.Printf("[有效token] CreateUser error: %v\n", err)
	} else {
		fmt.Printf("[有效token] CreateUser ok: id=%d name=%s\n", crsp.GetUser().GetId(), crsp.GetUser().GetName())
	}

	// --- 场景 2：不携带 token ---
	_, err = proxy.CreateUser(context.Background(), &pb.CreateUserReq{Name: "Dave"})
	if err != nil {
		fmt.Printf("[无token]   CreateUser error: code=%d msg=%s\n", errs.Code(err), errs.Msg(err))
	}
}
