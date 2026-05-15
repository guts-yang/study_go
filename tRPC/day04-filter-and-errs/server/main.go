// Package main 是 day04-filter-and-errs 的服务端入口。
//
// 核心演示：
//  1. 注册鉴权 filter（通过 blank import 触发 filter/auth.go 的 init）；
//  2. trpc_go.yaml 中配置 filter 链：recovery → auth；
//  3. 业务代码用 errs.New(code, msg) 返回带错误码的错误。
package main

import (
	_ "day04-filter-and-errs/filter" // 触发 auth filter 的 init 注册

	pb "day04-filter-and-errs/stub/trpc/study/user"

	trpc "git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/log"
)

func main() {
	s := trpc.NewServer()
	pb.RegisterUserServiceService(s, newUserImpl())
	log.Infof("UserService with auth filter starting on :8001")
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
