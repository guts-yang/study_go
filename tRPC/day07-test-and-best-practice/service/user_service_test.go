package service_test

import (
	"context"
	"testing"

	"day07-test-and-best-practice/service"
	pb "day07-test-and-best-practice/stub/trpc/study/user"

	"git.code.oa.com/trpc-go/trpc-go/errs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAggregatorService_Greet_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProxy := pb.NewMockUserServiceClientProxy(ctrl)
	mockProxy.EXPECT().
		GetUser(gomock.Any(), &pb.GetUserReq{Id: 1}).
		Return(&pb.GetUserRsp{User: &pb.User{Id: 1, Name: "Alice"}}, nil)

	svc := &service.AggregatorService{Downstream: mockProxy}
	greeting, err := svc.Greet(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, "Hello, Alice", greeting)
}

func TestAggregatorService_Greet_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 用框架 errs.New 构造带错误码的错误，让 errs.Code 能正确识别
	notFoundErr := errs.New(404, "user not found")

	mockProxy := pb.NewMockUserServiceClientProxy(ctrl)
	mockProxy.EXPECT().
		GetUser(gomock.Any(), &pb.GetUserReq{Id: 999}).
		Return(nil, notFoundErr)

	svc := &service.AggregatorService{Downstream: mockProxy}
	_, err := svc.Greet(context.Background(), 999)
	require.Error(t, err)
	assert.Equal(t, 404, errs.Code(err))
}
