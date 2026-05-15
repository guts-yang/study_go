// Package integration_test 演示集成测试：真正启动服务端，用真实 RPC 调用验证接口。
//
// 运行方式（需先启动服务端或在测试里内嵌启动）：
//
//	go test ./integration/... -v
package integration_test

import (
	"context"
	"testing"
	"time"

	pb "day07-test-and-best-practice/stub/trpc/study/user"

	_ "git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUserService_Integration 通过真实 RPC 调用验证 CreateUser 和 GetUser。
// 前置条件：需在另一个窗口先 go run ./server/ 启动服务端。
func TestUserService_Integration(t *testing.T) {
	proxy := pb.NewUserServiceClientProxy(
		client.WithTarget("ip://127.0.0.1:8001"),
		client.WithTimeout(2*time.Second),
	)
	ctx := context.Background()

	// 1. 创建用户
	crsp, err := proxy.CreateUser(ctx, &pb.CreateUserReq{Name: "IntegrationUser"})
	require.NoError(t, err)
	assert.NotZero(t, crsp.GetUser().GetId())
	assert.Equal(t, "IntegrationUser", crsp.GetUser().GetName())

	// 2. 查询用户
	grsp, err := proxy.GetUser(ctx, &pb.GetUserReq{Id: crsp.GetUser().GetId()})
	require.NoError(t, err)
	assert.Equal(t, "IntegrationUser", grsp.GetUser().GetName())

	// 3. 查询不存在的用户，期望报错
	_, err = proxy.GetUser(ctx, &pb.GetUserReq{Id: 99999})
	assert.Error(t, err)
}
