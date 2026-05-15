package main

import (
	"context"
	"strings"
	"sync"
	"time"

	pb "day04-filter-and-errs/stub/trpc/study/user"

	"git.code.oa.com/trpc-go/trpc-go/errs"
	"git.code.oa.com/trpc-go/trpc-go/log"
)

type userImpl struct {
	mu    sync.RWMutex
	seq   uint64
	store map[uint64]*pb.User
}

func newUserImpl() *userImpl {
	return &userImpl{store: make(map[uint64]*pb.User)}
}

func (u *userImpl) CreateUser(ctx context.Context, req *pb.CreateUserReq) (*pb.CreateUserRsp, error) {
	if strings.TrimSpace(req.GetName()) == "" {
		// errs.RetClientValidateFail = 51（客户端参数校验失败区间）
		return nil, errs.New(errs.RetClientValidateFail, "name is required")
	}
	u.mu.Lock()
	u.seq++
	user := &pb.User{Id: u.seq, Name: req.GetName(), CreatedAt: time.Now().Unix()}
	u.store[u.seq] = user
	u.mu.Unlock()
	log.InfoContextf(ctx, "CreateUser ok: id=%d", user.Id)
	return &pb.CreateUserRsp{User: user}, nil
}

func (u *userImpl) GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.GetUserRsp, error) {
	u.mu.RLock()
	user, ok := u.store[req.GetId()]
	u.mu.RUnlock()
	if !ok {
		// errs.RetServerNoFunc = 12（服务端逻辑错误区间）
		return nil, errs.New(errs.RetServerNoFunc, "user not found")
	}
	return &pb.GetUserRsp{User: user}, nil
}
