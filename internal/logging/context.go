package logging

import (
	"context"
)

const (
	loggerCtxKey = "logger"
)

func MayBeOf(ctx context.Context) *Logger {
	logger, ok := ctx.Value(loggerCtxKey).(*Logger)
	if !ok {
		logger = NewNoop()
	}

	return logger
}

func SetInCtx(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}
