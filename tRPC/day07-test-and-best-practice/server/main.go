// Package main 演示生产级 server：
//   - 全局 recovery filter
//   - 通过 admin /metrics 暴露指标
package main

import (
	"day07-test-and-best-practice/service"
	pb "day07-test-and-best-practice/stub/trpc/study/user"

	trpc "git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/log"
)

func main() {
	s := trpc.NewServer()
	pb.RegisterUserServiceService(s, &userImpl{
		agg: &service.AggregatorService{},
	})
	log.Info("server starting on :8001, admin on :11014")
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
