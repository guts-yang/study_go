// Package main 是 tRPC-Go demo 服务端业务实现：UserServiceImpl。
//
// 内部使用 map[int64]*User + sync.RWMutex 做线程安全的内存存储。
// 这里的代码是纯业务，与 tRPC 框架无关——这正是好框架应有的样子：
// 业务代码完全感知不到协议、序列化、网络的存在。
package main

import (
	"context"
	"sync"
	"sync/atomic"

	"trpc.group/trpc-go/trpc-go/errs"
	"trpc.group/trpc-go/trpc-go/log"

	"trpc-go-demo/api/user"
)

// UserServiceImpl 实现了 user.UserService 接口。
type UserServiceImpl struct {
	mu     sync.RWMutex          // 保护 users 的并发读写
	users  map[int64]*user.User  // 内存存储
	nextID int64                 // 用 atomic 操作保证自增 ID 的并发安全
}

// NewUserServiceImpl 构造一个 UserServiceImpl 实例。
func NewUserServiceImpl() *UserServiceImpl {
	return &UserServiceImpl{
		users: make(map[int64]*user.User),
	}
}

// GetUser 根据 ID 查询用户。找不到则返回业务错误码 404。
func (s *UserServiceImpl) GetUser(ctx context.Context, req *user.GetUserReq) (*user.GetUserRsp, error) {
	log.Infof("[GetUser] req.ID=%d", req.ID)

	s.mu.RLock()
	defer s.mu.RUnlock()

	u, ok := s.users[req.ID]
	if !ok {
		// errs.New(code, msg) 创建带业务错误码的 tRPC 错误，
		// 客户端收到后可以通过 errs.Code(err) / errs.Msg(err) 取出来。
		return nil, errs.New(404, "user not found")
	}
	return &user.GetUserRsp{User: u}, nil
}

// CreateUser 创建用户，返回带新 ID 的 user 对象。
func (s *UserServiceImpl) CreateUser(ctx context.Context, req *user.CreateUserReq) (*user.CreateUserRsp, error) {
	log.Infof("[CreateUser] req.Name=%q", req.Name)

	if req.Name == "" {
		return nil, errs.New(400, "name is required")
	}

	id := atomic.AddInt64(&s.nextID, 1) // 原子自增，避免锁内自增
	u := &user.User{ID: id, Name: req.Name}

	s.mu.Lock()
	s.users[id] = u
	s.mu.Unlock()

	return &user.CreateUserRsp{User: u}, nil
}

// ListUser 返回当前所有用户的快照。
func (s *UserServiceImpl) ListUser(ctx context.Context, req *user.ListUserReq) (*user.ListUserRsp, error) {
	log.Infof("[ListUser] called")

	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*user.User, 0, len(s.users))
	for _, u := range s.users {
		users = append(users, u) // 注意：返回的是指针快照，业务方不应再修改
	}
	return &user.ListUserRsp{Users: users, Total: len(users)}, nil
}
