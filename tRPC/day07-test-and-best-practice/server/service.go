package main

import (
	"context"
	"strings"
	"sync"
	"time"

	"day07-test-and-best-practice/service"
	pb "day07-test-and-best-practice/stub/trpc/study/user"

	"git.code.oa.com/trpc-go/trpc-go/errs"
	"git.code.oa.com/trpc-go/trpc-go/log"
)

type userImpl struct {
	agg  *service.AggregatorService
	mu   sync.RWMutex
	seq  uint64
	mem  map[uint64]*pb.User
}

func (u *userImpl) init() {
	if u.mem == nil {
		u.mem = make(map[uint64]*pb.User)
	}
}

func (u *userImpl) CreateUser(ctx context.Context, req *pb.CreateUserReq) (*pb.CreateUserRsp, error) {
	if strings.TrimSpace(req.GetName()) == "" {
		return nil, errs.New(errs.RetClientValidateFail, "name is required")
	}
	u.mu.Lock()
	u.init()
	u.seq++
	user := &pb.User{Id: u.seq, Name: req.GetName(), CreatedAt: time.Now().Unix()}
	u.mem[u.seq] = user
	u.mu.Unlock()
	log.InfoContextf(ctx, "CreateUser: id=%d", user.Id)
	return &pb.CreateUserRsp{User: user}, nil
}

func (u *userImpl) GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.GetUserRsp, error) {
	u.mu.RLock()
	user, ok := u.mem[req.GetId()]
	u.mu.RUnlock()
	if !ok {
		return nil, errs.New(404, "user not found")
	}
	log.InfoContextf(ctx, "GetUser: id=%d", user.Id)
	return &pb.GetUserRsp{User: user}, nil
}
