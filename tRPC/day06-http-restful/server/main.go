// Package main 是 day06-http-restful 的服务端入口。
//
// 核心演示：同一份 UserService 实现，通过 trpc_go.yaml 配置两个 service：
//   - port 8001: protocol: trpc（tRPC 私有协议）
//   - port 8080: protocol: http（HTTP/JSON，RESTful 风格）
//
// HTTP 端口的路由由框架根据 proto method 自动生成：
//
//	POST /trpc.study.user.UserService/CreateUser
//	POST /trpc.study.user.UserService/GetUser
package main

import (
	pb "day06-http-restful/stub/trpc/study/user"

	trpc "git.code.oa.com/trpc-go/trpc-go"
	_ "git.code.oa.com/trpc-go/trpc-go/http" // 注册 HTTP 传输层
	"git.code.oa.com/trpc-go/trpc-go/log"
)

func main() {
	s := trpc.NewServer()
	// 同一个 impl 注册到 server，框架自动分发到两个 service（tRPC + HTTP）
	pb.RegisterUserServiceService(s, newUserImpl())
	log.Infof("UserService starting: tRPC=:8001 HTTP=:8080")
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
