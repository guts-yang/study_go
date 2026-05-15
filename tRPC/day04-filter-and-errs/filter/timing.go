package filter

import (
	"context"
	"time"

	"git.code.oa.com/trpc-go/trpc-go/filter"
	"git.code.oa.com/trpc-go/trpc-go/log"
)

// Timing 客户端耗时 filter。
func Timing(ctx context.Context, req, rsp interface{}, next filter.ClientHandleFunc) error {
	start := time.Now()
	err := next(ctx, req, rsp)
	log.InfoContextf(ctx, "rpc cost=%s err=%v", time.Since(start), err)
	return err
}

func init() {
	// Register(name, serverFilter, clientFilter)：第二参数传 nil 表示只注册客户端 filter
	filter.Register("timing", nil, Timing)
}
