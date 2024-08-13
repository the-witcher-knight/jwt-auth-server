package tracing

import (
	pkgerrors "github.com/pkg/errors"
)

type StackTracer interface {
	StackTrace() pkgerrors.StackTrace
}
