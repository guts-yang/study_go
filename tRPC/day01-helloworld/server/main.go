// Package main 是 day01-helloworld 的服务端入口。
//
// ⚠️ 编译前置条件：
//
//	cd .\tRPC\day01-helloworld
//	trpc create -p .\proto\helloworld.proto -o . --rpconly
//	go mod tidy
//
// 上面命令会在 ./stub/trpc/helloworld/ 目录下生成 *.pb.go / *.trpc.go / *_mock.go。
// 之前直接 `go run` 会因为找不到 stub 包而编译失败，这是预期的——
// tRPC 的工作流就是 "先生成桩代码，再写业务代码"。
package main

import (
	"context"

	pb "day01-helloworld/stub/trpc/helloworld"

	trpc "git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/log"
)

// greeterImpl 实现 pb.GreeterService 接口。
//
// 由 trpc-cmdline 生成的 helloworld.trpc.go 中会有类似这样的接口定义：
//
//	type GreeterService interface {
//	    SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error)
//	}
//
// 我们只需提供这个接口的实现即可，框架会负责协议解析、路由、并发分发。
type greeterImpl struct{}

// SayHello 处理客户端的问候请求。
// ctx 携带了超时、调用链 metadata、对端信息等；req 已经是反序列化好的强类型对象。
func (g *greeterImpl) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Infof("收到 SayHello 请求: msg=%s", req.GetMsg())
	return &pb.HelloReply{Msg: "Hello, " + req.GetMsg()}, nil
}

func main() {
	// trpc.NewServer() 默认读取启动目录下的 trpc_go.yaml；
	// 也可以用 -conf 指定路径： go run .\server\ -conf .\trpc_go.yaml
	s := trpc.NewServer()

	// 把实现注册到 server。注册方法名 RegisterGreeterService 由桩代码生成。
	pb.RegisterGreeterService(s, &greeterImpl{})

	// Serve 阻塞运行，直到收到 SIGINT/SIGTERM；框架内部已处理优雅退出。
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
