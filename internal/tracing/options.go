package tracing

import (
	"go.uber.org/zap"
)

type Option func(*Tracer)

func Noop() Option {
	return func(tracer *Tracer) {
		tracer.zapLog = zap.NewNop()
	}
}
