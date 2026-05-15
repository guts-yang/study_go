package main

import (
	"context"
	"reflect"
	"testing"

	_ "git.code.oa.com/trpc-go/trpc-go/http"

	"github.com/golang/mock/gomock"

	pb "day05-naming-and-client/stub/trpc/study/user"
)

//go:generate go mod tidy
//go:generate mockgen -destination=stub/day05-naming-and-client/stub/trpc/study/user/user_mock.go -package=user -self_package=day05-naming-and-client/stub/trpc/study/user --source=stub/day05-naming-and-client/stub/trpc/study/user/user.trpc.go

func Test_userServiceImpl_CreateUser(t *testing.T) {
	// 开始写mock逻辑
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userServiceService := pb.NewMockUserServiceService(ctrl)
	var inorderClient []*gomock.Call
	// 预期行为
	m := userServiceService.EXPECT().CreateUser(gomock.Any(), gomock.Any()).AnyTimes()
	m.DoAndReturn(func(ctx context.Context, req *pb.CreateUserReq) (*pb.CreateUserRsp, error) {
		// 直接返回预设响应；如需调用真实实现，请在同一 package 下声明 impl 并注入
		return &pb.CreateUserRsp{User: &pb.User{Name: req.GetName()}}, nil
	})
	gomock.InOrder(inorderClient...)

	// 开始写单元测试逻辑
	type args struct {
		ctx context.Context
		req *pb.CreateUserReq
		rsp *pb.CreateUserRsp
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rsp *pb.CreateUserRsp
			var err error
			if rsp, err = userServiceService.CreateUser(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("userServiceImpl.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(rsp, tt.args.rsp) {
				t.Errorf("userServiceImpl.CreateUser() rsp got = %v, want %v", rsp, tt.args.rsp)
			}
		})
	}
}

func Test_userServiceImpl_GetUser(t *testing.T) {
	// 开始写mock逻辑
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userServiceService := pb.NewMockUserServiceService(ctrl)
	var inorderClient []*gomock.Call
	// 预期行为
	m := userServiceService.EXPECT().GetUser(gomock.Any(), gomock.Any()).AnyTimes()
	m.DoAndReturn(func(ctx context.Context, req *pb.GetUserReq) (*pb.GetUserRsp, error) {
		// 直接返回预设响应；如需调用真实实现，请在同一 package 下声明 impl 并注入
		return &pb.GetUserRsp{User: &pb.User{Id: req.GetId()}}, nil
	})
	gomock.InOrder(inorderClient...)

	// 开始写单元测试逻辑
	type args struct {
		ctx context.Context
		req *pb.GetUserReq
		rsp *pb.GetUserRsp
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rsp *pb.GetUserRsp
			var err error
			if rsp, err = userServiceService.GetUser(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("userServiceImpl.GetUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(rsp, tt.args.rsp) {
				t.Errorf("userServiceImpl.GetUser() rsp got = %v, want %v", rsp, tt.args.rsp)
			}
		})
	}
}
