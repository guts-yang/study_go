// Package main 是 day03-config-and-admin 的服务端入口。
//
// 本 day 重点：
//  1. trpc_go.yaml 中启用 admin 端口（11014）；
//  2. 配置多个 log writer（console + file）；
//  3. 在业务代码里用框架 log 而不是 fmt.Printf。
//
// 启动后可访问 admin：
//
//	curl http://127.0.0.1:11014/is_healthy/
//	curl http://127.0.0.1:11014/cmds
//	curl http://127.0.0.1:11014/version
package main

import (
	pb "day03-config-and-admin/stub/trpc/study/user"

	trpc "git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/log"
)

func main() {
	s := trpc.NewServer()
	pb.RegisterUserServiceService(s, newUserImpl())
	log.Infof("UserService starting, admin port=11014")
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
