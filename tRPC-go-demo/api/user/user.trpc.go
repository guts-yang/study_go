// Package user 中本文件是「手写桩代码」，等价于 protoc-gen-trpc 生成的 *.trpc.go。
//
// 内容包含五大块：
//   1. UserService          —— 服务端业务方需要实现的接口。
//   2. UserService_*_Handler —— 三个 Method.Func 闭包：负责反序列化 + filter 链 + 调用业务方法。
//   3. UserServer_ServiceDesc —— 注册元信息（ServiceName / HandlerType / Methods）。
//   4. RegisterUserService   —— 业务侧 main 函数中调用，把实现注册到 trpc.NewServer()。
//   5. UserClientProxy       —— 客户端代理接口及其默认实现（NewUserClientProxy）。
//
// 学习要点：把 protoc 自动生成的 .pb.go 摆开来人手写一遍，能看清 tRPC-Go 框架的契约边界——
// 框架要的就这些东西，剩下的都是它的运行时（transport / codec / filter / router）。
package user

import (
	"context"
	"fmt"

	"trpc.group/trpc-go/trpc-go/client"
	"trpc.group/trpc-go/trpc-go/codec"
	"trpc.group/trpc-go/trpc-go/server"
)

// ============================================================================
// 1. 服务端业务接口
// ============================================================================

// UserService 是业务方需要实现的接口。
// 服务端 main 中会 new 一个实现该接口的对象，传给 RegisterUserService。
type UserService interface {
	GetUser(ctx context.Context, req *GetUserReq) (*GetUserRsp, error)
	CreateUser(ctx context.Context, req *CreateUserReq) (*CreateUserRsp, error)
	ListUser(ctx context.Context, req *ListUserReq) (*ListUserRsp, error)
}

// ============================================================================
// 2. 三个 Method.Func 闭包（Server-side handler）
// ============================================================================
//
// Method.Func 的标准签名（来自 server.Method 结构体）：
//   func(svr interface{}, ctx context.Context, f server.FilterFunc) (rspBody interface{}, err error)
//
// 其中 server.FilterFunc 的类型为：
//   type FilterFunc func(reqBody interface{}) (filter.ServerChain, error)
//
// 在 handler 内部：
//   1) 创建空的 req 指针。
//   2) 调用 f(req)：框架会用注册的 codec 把网络字节填进 req，并返回 server filter chain。
//   3) 定义最终业务调用闭包 handleFunc：把 req 类型断言到具体类型 → 调业务接口。
//   4) 走 filter chain：filters.Filter(ctx, req, handleFunc)。
//
// 这是手写桩代码的"灵魂"——理解了这一段，整个 tRPC-Go 服务端处理链路就一清二楚了。

// UserService_GetUser_Handler 处理 GetUser RPC。
func UserService_GetUser_Handler(svr interface{}, ctx context.Context, f server.FilterFunc) (interface{}, error) {
	req := &GetUserReq{}
	filters, err := f(req) // ← 反序列化 req + 拿到 filter chain
	if err != nil {
		return nil, err
	}
	handleFunc := func(ctx context.Context, reqBody interface{}) (interface{}, error) {
		return svr.(UserService).GetUser(ctx, reqBody.(*GetUserReq))
	}
	return filters.Filter(ctx, req, handleFunc)
}

// UserService_CreateUser_Handler 处理 CreateUser RPC。
func UserService_CreateUser_Handler(svr interface{}, ctx context.Context, f server.FilterFunc) (interface{}, error) {
	req := &CreateUserReq{}
	filters, err := f(req)
	if err != nil {
		return nil, err
	}
	handleFunc := func(ctx context.Context, reqBody interface{}) (interface{}, error) {
		return svr.(UserService).CreateUser(ctx, reqBody.(*CreateUserReq))
	}
	return filters.Filter(ctx, req, handleFunc)
}

// UserService_ListUser_Handler 处理 ListUser RPC。
func UserService_ListUser_Handler(svr interface{}, ctx context.Context, f server.FilterFunc) (interface{}, error) {
	req := &ListUserReq{}
	filters, err := f(req)
	if err != nil {
		return nil, err
	}
	handleFunc := func(ctx context.Context, reqBody interface{}) (interface{}, error) {
		return svr.(UserService).ListUser(ctx, reqBody.(*ListUserReq))
	}
	return filters.Filter(ctx, req, handleFunc)
}

// ============================================================================
// 3. ServiceDesc：注册元信息
// ============================================================================
//
// ServiceName 必须与 trpc_go.yaml 里的 server.service[].name 完全一致，
// 否则服务不会被正确启动。
//
// Method.Name 是「带斜杠的全限定 RPC 名」，格式为 /<ServiceName>/<MethodName>，
// 客户端 codec.Message.WithClientRPCName(...) 也要传一致的字符串。
var UserServer_ServiceDesc = server.ServiceDesc{
	ServiceName: "trpc.demo.user.User",
	HandlerType: ((*UserService)(nil)),
	Methods: []server.Method{
		{
			Name: "/trpc.demo.user.User/GetUser",
			Func: UserService_GetUser_Handler,
		},
		{
			Name: "/trpc.demo.user.User/CreateUser",
			Func: UserService_CreateUser_Handler,
		},
		{
			Name: "/trpc.demo.user.User/ListUser",
			Func: UserService_ListUser_Handler,
		},
	},
}

// ============================================================================
// 4. RegisterUserService：业务侧入口
// ============================================================================

// RegisterUserService 把业务实现注册到 server.Service。
// 第一个参数 s 通常是 trpc.NewServer() 返回的 *server.Server（它实现了 server.Service）。
//
// 内部只做两件事：
//   - 把 UserServer_ServiceDesc 的指针 + 业务实现传给 s.Register。
//   - 注册失败直接 panic（启动失败应该立刻让进程退出）。
func RegisterUserService(s server.Service, svr UserService) {
	if err := s.Register(&UserServer_ServiceDesc, svr); err != nil {
		panic(fmt.Sprintf("User service register error: %v", err))
	}
}

// ============================================================================
// 5. UserClientProxy：客户端代理
// ============================================================================
//
// 客户端代理的作用：把"调远程 RPC"这件事封装得像"调本地函数"。
// 内部三步走：
//   ① 拿到一个 codec.Message（codec.WithCloneMessage 创建一个干净的 msg 挂到 ctx 上）。
//   ② 把 RPC 元信息塞到 msg：CalleeServiceName、ClientRPCName、SerializationType...
//      这些信息会被 transport+codec 编码进 tRPC 协议的 PB 包头。
//   ③ 调用 client.Client.Invoke(ctx, req, rsp, opts...)，等返回。

// UserClientProxy 是客户端调用接口（与 UserService 镜像，但多了 opts 可变参数）。
type UserClientProxy interface {
	GetUser(ctx context.Context, req *GetUserReq, opts ...client.Option) (*GetUserRsp, error)
	CreateUser(ctx context.Context, req *CreateUserReq, opts ...client.Option) (*CreateUserRsp, error)
	ListUser(ctx context.Context, req *ListUserReq, opts ...client.Option) (*ListUserRsp, error)
}

// userClientProxyImpl 是 UserClientProxy 的默认实现。
type userClientProxyImpl struct {
	client client.Client   // 底层 client，调用其 Invoke 完成发送 / 收包 / 反序列化
	opts   []client.Option // 创建时的全局 opts，每次调用都会附加
}

// NewUserClientProxy 创建客户端代理。
// 是 var 而不是 func，方便测试时被覆盖以便 mock。
var NewUserClientProxy = func(opts ...client.Option) UserClientProxy {
	return &userClientProxyImpl{
		client: client.DefaultClient, // 复用全局单例
		opts:   opts,
	}
}

// invoke 是三个 RPC 方法共用的内部封装。
// rpcName 形如 "/trpc.demo.user.User/GetUser"。
func (c *userClientProxyImpl) invoke(
	ctx context.Context, rpcName, methodName string,
	req, rsp interface{}, opts ...client.Option,
) error {
	// 创建一份新的 codec.Message 挂到 ctx 上（每次 RPC 一份，避免并发污染）。
	ctx, msg := codec.WithCloneMessage(ctx)
	defer codec.PutBackMessage(msg)

	// —— RPC 路由信息：写在 msg 上（不是 client.Option！）——
	msg.WithClientRPCName(rpcName)
	msg.WithCalleeServiceName(UserServer_ServiceDesc.ServiceName)
	msg.WithCalleeApp("demo")
	msg.WithCalleeServer("user")
	msg.WithCalleeService("User")
	msg.WithCalleeMethod(methodName)
	// —— 序列化方式：JSON（值=2，框架默认已注册 JSON serializer）——
	msg.WithSerializationType(codec.SerializationTypeJSON)

	// 合并构造时的 opts 与本次调用的 opts。
	callopts := make([]client.Option, 0, len(c.opts)+len(opts))
	callopts = append(callopts, c.opts...)
	callopts = append(callopts, opts...)

	// 真正发送：底层 client 会做 codec.Encode → transport.Send → 收 → codec.Decode → 填 rsp。
	return c.client.Invoke(ctx, req, rsp, callopts...)
}

// GetUser 远程调用 /trpc.demo.user.User/GetUser。
func (c *userClientProxyImpl) GetUser(ctx context.Context, req *GetUserReq, opts ...client.Option) (*GetUserRsp, error) {
	rsp := &GetUserRsp{}
	if err := c.invoke(ctx, "/trpc.demo.user.User/GetUser", "GetUser", req, rsp, opts...); err != nil {
		return nil, err
	}
	return rsp, nil
}

// CreateUser 远程调用 /trpc.demo.user.User/CreateUser。
func (c *userClientProxyImpl) CreateUser(ctx context.Context, req *CreateUserReq, opts ...client.Option) (*CreateUserRsp, error) {
	rsp := &CreateUserRsp{}
	if err := c.invoke(ctx, "/trpc.demo.user.User/CreateUser", "CreateUser", req, rsp, opts...); err != nil {
		return nil, err
	}
	return rsp, nil
}

// ListUser 远程调用 /trpc.demo.user.User/ListUser。
func (c *userClientProxyImpl) ListUser(ctx context.Context, req *ListUserReq, opts ...client.Option) (*ListUserRsp, error) {
	rsp := &ListUserRsp{}
	if err := c.invoke(ctx, "/trpc.demo.user.User/ListUser", "ListUser", req, rsp, opts...); err != nil {
		return nil, err
	}
	return rsp, nil
}
