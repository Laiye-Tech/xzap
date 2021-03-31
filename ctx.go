package xzap

import (
	"context"
)

type (
	msgIDKey struct{}
	debugKey struct{}
)

func InjectCtx(ctx context.Context, opts ...CtxOption) context.Context {
	if ctx.Value(debugKey{}) == nil {
		if debugInHead(ctx) {
			opts = append(opts, WithDebug())
		}
	}
	for _, opt := range opts {
		ctx = opt.Apply(ctx)
	}
	return ctx
}
