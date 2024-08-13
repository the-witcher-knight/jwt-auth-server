package tracing

import (
	"context"
)

type contextKey string

const (
	contextKeyTracer contextKey = "tracer"
)

func SetInContext(ctx context.Context, tracer *Tracer) context.Context {
	return context.WithValue(ctx, contextKeyTracer, tracer)
}

func FromContext(ctx context.Context) *Tracer {
	tracer := ctx.Value(contextKeyTracer)
	if tracer == nil {
		return New(nil, Noop())
	}

	return tracer.(*Tracer)
}
