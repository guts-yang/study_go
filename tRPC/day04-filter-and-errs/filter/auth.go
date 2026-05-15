package filter

import (
	"context"

	trpc "git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/errs"
	"git.code.oa.com/trpc-go/trpc-go/filter"
	"git.code.oa.com/trpc-go/trpc-go/log"
)

const validToken = "valid-token"

func authFilter(ctx context.Context, req interface{}, next filter.ServerHandleFunc) (interface{}, error) {
	msg := trpc.Message(ctx)
	token := msg.ServerMetaData()["token"]
	log.DebugContextf(ctx, "authFilter: token=%s", token)
	if string(token) != validToken {
		return nil, errs.New(errs.RetServerAuthFail, "invalid token")
	}
	return next(ctx, req)
}

func init() {
	filter.Register("auth", authFilter, nil)
}
