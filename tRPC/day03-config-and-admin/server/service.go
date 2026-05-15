package main

import (
	"context"
	"strings"
	"sync"
	"time"

	pb "day03-config-and-admin/stub/trpc/study/user"

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
		return nil, errs.New(errs.RetClientValidateFail, "name is required")
	}
	u.mu.Lock()
	u.seq++
	user := &pb.User{Id: u.seq, Name: req.GetName(), CreatedAt: time.Now().Unix()}
	u.store[u.seq] = user
	u.mu.Unlock()

	// 使用框架 log，而非 fmt.Printf；框架会自动附加 caller/traceID 等字段
	log.InfoContextf(ctx, "CreateUser: id=%d name=%s", user.Id, user.Name)
	return &pb.CreateUserRsp{User: user}, nil
}

func (u *userImpl) GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.GetUserRsp, error) {
	u.mu.RLock()
	user, ok := u.store[req.GetId()]
	u.mu.RUnlock()
	if !ok {
		return nil, errs.New(errs.RetServerNoFunc, "user not found")
	}
	log.InfoContextf(ctx, "GetUser: id=%d name=%s", user.Id, user.Name)
	return &pb.GetUserRsp{User: user}, nil
}
