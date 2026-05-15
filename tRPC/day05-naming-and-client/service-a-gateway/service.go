package main

import (
	"context"
	"time"

	pb "day05-naming-and-client/stub/trpc/study/user"

	"git.code.oa.com/trpc-go/trpc-go/client"
	"git.code.oa.com/trpc-go/trpc-go/errs"
	"git.code.oa.com/trpc-go/trpc-go/log"
)

// gatewayImpl 持有下游 User 服务的 proxy。
//
// 关键设计：
//  1. proxy 在构造时一次性创建，handler 里复用 —— 利用 selector 内置的连接池
//  2. WithTarget 用 example:// 走自定义 selector
//  3. WithTimeout 是该 callee 的上限，与 context.WithTimeout 取最小
type gatewayImpl struct {
	userCli pb.UserServiceClientProxy
}

func newGatewayImpl() *gatewayImpl {
	return &gatewayImpl{
		userCli: pb.NewUserServiceClientProxy(
			client.WithTarget("example://trpc.study.user.UserService"),
			client.WithTimeout(500*time.Millisecond),
		),
	}
}

func (g *gatewayImpl) GreetUser(ctx context.Context, req *pb.GreetUserReq) (*pb.GreetUserRsp, error) {
	rsp, err := g.userCli.GetUser(ctx, &pb.GetUserReq{Id: req.GetId()})
	if err != nil {
		// 把下游错误包装成网关层错误：错误码 502（业务区间），保留 chain
		log.WarnContextf(ctx, "call User.GetUser failed: %v", err)
		return nil, errs.Wrap(err, 502, "downstream user service failed")
	}
	return &pb.GreetUserRsp{Greeting: "Hello, " + rsp.GetUser().GetName()}, nil
}
