// 下游 User 服务，逻辑同 day02。
package main

import (
	pb "day05-naming-and-client/stub/trpc/study/user"

	trpc "git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/log"
)

func main() {
	s := trpc.NewServer()
	pb.RegisterUserServiceService(s, newUserImpl())
	log.Info("user service starting on :8001")
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
