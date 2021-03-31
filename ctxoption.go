package xzap

import (
	"context"

	"google.golang.org/grpc/metadata"
)

type CtxOption interface {
	Apply(ctx context.Context) context.Context
}

type ctxFuncOption func(context.Context) context.Context

func (f ctxFuncOption) Apply(c context.Context) context.Context {
	return f(c)
}

func WithMsgID(msgID uint64) CtxOption {
	return ctxFuncOption(func(c context.Context) context.Context {
		return context.WithValue(c, msgIDKey{}, msgID)
	})
}

func msgID(ctx context.Context) uint64 {
	v := ctx.Value(msgIDKey{})
	if v != nil {
		return v.(uint64)
	}
	return 0
}

func WithDebug() CtxOption {
	return ctxFuncOption(func(c context.Context) context.Context {
		return context.WithValue(c, debugKey{}, debugKey{})
	})
}

func debug(ctx context.Context) bool {
	v := ctx.Value(debugKey{})
	if v != nil {
		return true
	}
	return debugInHead(ctx)
}

func debugInHead(ctx context.Context) bool {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		_, ok := md["debug"]
		return ok
	}
	return false
}
