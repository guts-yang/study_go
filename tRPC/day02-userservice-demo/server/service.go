package main

import (
	"context"
	"strings"
	"sync"
	"time"

	pb "day02-userservice-demo/stub/trpc/study/user"

	"git.code.oa.com/trpc-go/trpc-go/errs"
	"git.code.oa.com/trpc-go/trpc-go/log"
)

// userImpl 内存版用户存储，对应 day07 的 userStore。
type userImpl struct {
	mu    sync.RWMutex
	seq   uint64
	store map[uint64]*pb.User
}

func newUserImpl() *userImpl {
	return &userImpl{store: make(map[uint64]*pb.User)}
}

// CreateUser 创建新用户。
func (u *userImpl) CreateUser(ctx context.Context, req *pb.CreateUserReq) (*pb.CreateUserRsp, error) {
	if strings.TrimSpace(req.GetName()) == "" {
		return nil, errs.New(errs.RetClientValidateFail, "name is required")
	}
	u.mu.Lock()
	u.seq++
	user := &pb.User{
		Id:        u.seq,
		Name:      req.GetName(),
		CreatedAt: time.Now().Unix(),
	}
	u.store[u.seq] = user
	u.mu.Unlock()

	log.Infof("CreateUser: id=%d name=%s", user.Id, user.Name)
	return &pb.CreateUserRsp{User: user}, nil
}

// GetUser 按 ID 查询用户。
func (u *userImpl) GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.GetUserRsp, error) {
	u.mu.RLock()
	user, ok := u.store[req.GetId()]
	u.mu.RUnlock()

	if !ok {
		return nil, errs.New(errs.RetServerNoFunc, "user not found")
	}
	log.Infof("GetUser: id=%d name=%s", user.Id, user.Name)
	return &pb.GetUserRsp{User: user}, nil
}
