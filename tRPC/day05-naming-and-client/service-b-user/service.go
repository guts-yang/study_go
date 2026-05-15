// 业务实现完全沿用 day02。
package main

import (
	"context"
	"strings"
	"sync"
	"time"

	pb "day05-naming-and-client/stub/trpc/study/user"

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
		return nil, errs.New(400, "name is required")
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	u.seq++
	user := &pb.User{Id: u.seq, Name: req.GetName(), CreatedAt: time.Now().Unix()}
	u.store[user.Id] = user
	log.InfoContextf(ctx, "[user] CreateUser id=%d", user.Id)
	return &pb.CreateUserRsp{User: user}, nil
}

func (u *userImpl) GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.GetUserRsp, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	user, ok := u.store[req.GetId()]
	if !ok {
		return nil, errs.New(404, "user not found")
	}
	return &pb.GetUserRsp{User: user}, nil
}
