// Package main 是 day02-userservice-demo 的服务端入口。
//
// 对标 day07 HTTP 版：用 tRPC 实现同样的 CreateUser / GetUser 语义。
// 区别在于：路由、序列化、并发都由框架处理，业务代码只需实现接口方法。
package main

import (
	pb "day02-userservice-demo/stub/trpc/study/user"

	trpc "git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/log"
)

func main() {
	s := trpc.NewServer()
	pb.RegisterUserServiceService(s, newUserImpl())
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
