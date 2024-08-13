package tracing

import (
	"errors"
	"fmt"
	"syscall"

	"go.uber.org/zap"
)

func (tr *Tracer) Info(format string, args ...interface{}) {
	tr.zapLog.Info(fmt.Sprintf(format, args...))
}

func (tr *Tracer) Error(err error, format string, args ...interface{}) {
	logFields := []zap.Field{
		zap.Error(err),
	}

	if stackTraced, ok := err.(StackTracer); ok {
		logFields = append(logFields,
			zap.String("exception.stack", fmt.Sprintf("%v", stackTraced.StackTrace())),
		)
	}

	tr.zapLog.Error(fmt.Sprintf(format, args...), logFields...)
}

func (tr *Tracer) WithAttributes(attrs ...Attribute) *Tracer {
	cloned := *tr
	if len(attrs) > 0 {
		var zapFields []zap.Field
		for _, attr := range attrs {
			cloned.logTags = append(cloned.logTags, attr)
			zapFields = append(zapFields, toZapField(attr))
		}

		cloned.zapLog = cloned.zapLog.With(zapFields...)
	}
	
	return &cloned
}

func (tr *Tracer) Flush() error {
	if err := tr.zapLog.Sync(); err != nil {
		// Ignore this stderr https://github.com/uber-go/zap/issues/328
		if !errors.Is(err, syscall.ENOTTY) && !errors.Is(err, syscall.EINVAL) {
			return err
		}
	}

	return nil
}
