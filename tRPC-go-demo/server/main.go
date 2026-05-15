// Package main 是 tRPC-Go demo 服务端入口。
//
// 启动流程（仅 3 行核心代码）：
//   1. trpc.NewServer()                         读取 trpc_go.yaml，构造 *server.Server
//   2. user.RegisterUserService(s, impl)        把业务实现注册到 ServiceDesc
//   3. s.Serve()                                启动监听 + accept 循环（阻塞）
//
// 框架内部会做：
//   - 读 yaml 中 server.service 列表，按 protocol=trpc 创建对应的 transport / codec
//   - 监听 ip:port，accept 新连接
//   - 收到帧后：framer 分帧 → codec.Decode（PB 包头 + JSON 包体）
//   - 路由到 ServiceDesc.Methods 中匹配的 Method.Func（手写桩闭包）
//   - 闭包内反序列化 req → 走 filter 链 → 调用业务方法 → 返回 rsp
//   - codec.Encode 打回 → transport.Write
package main

import (
	trpc "trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"

	"trpc-go-demo/api/user"
)

func main() {
	// trpc.NewServer 会自动加载当前目录下的 trpc_go.yaml（默认路径）。
	// 如果要指定其他配置文件，可用 trpc.NewServer(server.WithConfigPath("xxx.yaml"))。
	s := trpc.NewServer()

	// 把业务实现注册到框架。s 是 *server.Server，实现了 server.Service 接口。
	user.RegisterUserService(s, NewUserServiceImpl())

	log.Info("tRPC-Go demo server starting...")
	log.Info("ServiceName: trpc.demo.user.User")
	log.Info("Listen     : 0.0.0.0:8000  (protocol=trpc, serialization=json)")
	log.Info("RPC paths  : /trpc.demo.user.User/{GetUser|CreateUser|ListUser}")

	// Serve 会阻塞直到收到信号或致命错误。
	if err := s.Serve(); err != nil {
		log.Fatalf("server.Serve error: %v", err)
	}
}
