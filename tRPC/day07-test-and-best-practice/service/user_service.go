// Package service 演示"可测试"的业务实现：
// 把下游 client 抽象成接口字段，单测里用 gomock 替换。
package service

import (
	"context"

	pb "day07-test-and-best-practice/stub/trpc/study/user"

	"git.code.oa.com/trpc-go/trpc-go/errs"
)

// AggregatorService 聚合下游 UserService 提供问候服务。
type AggregatorService struct {
	Downstream pb.UserServiceClientProxy
}

// Greet 拼接 "Hello, <name>"；下游报错时按错误码分别处理。
func (a *AggregatorService) Greet(ctx context.Context, id uint64) (string, error) {
	if a.Downstream == nil {
		return "", errs.New(500, "downstream not configured")
	}
	rsp, err := a.Downstream.GetUser(ctx, &pb.GetUserReq{Id: id})
	if err != nil {
		if errs.Code(err) == 404 {
			return "", errs.New(404, "user not found")
		}
		return "", errs.Wrap(err, 502, "downstream failed")
	}
	return "Hello, " + rsp.GetUser().GetName(), nil
}
