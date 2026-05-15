// 上游 Gateway 服务：在 handler 中调下游 User 服务。
package main

import (
	pb "day05-naming-and-client/stub/trpc/study/user"

	// 触发自定义 selector 注册（example://）
	_ "day05-naming-and-client/selector"

	trpc "git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/log"
)

func main() {
	s := trpc.NewServer()
	pb.RegisterGatewayServiceService(s, newGatewayImpl())
	log.Info("gateway service starting on :8002")
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
