package xzap

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
)

type spanContext interface {
	IsSampled() bool
	TraceID() jaeger.TraceID
}

// getTraceSpanContext 获取 ctx 中 span context
func getTraceSpanContext(ctx context.Context) spanContext {
	if !opentracing.IsGlobalTracerRegistered() {
		return nil
	}
	if span := opentracing.SpanFromContext(ctx); span != nil {
		if jaegerSpan, ok := span.(*jaeger.Span); ok {
			return jaegerSpan.SpanContext()
		}
	}
	return nil
}
