package tracing

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Tracer struct {
	logTags []Attribute

	zapLog *zap.Logger
}

func New(opts ...Option) *Tracer {
	l := &Tracer{
		zapLog: newZap(nil),
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func newZap(logTags []Attribute) *zap.Logger {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), os.Stdout, zap.InfoLevel)
	zapLog := zap.New(core, zap.AddCaller())

	if len(logTags) > 0 {
		var logFields []zap.Field
		for _, attr := range logTags {
			logFields = append(logFields, toZapField(attr))
		}

		zapLog = zapLog.With(logFields...)
	}

	return zapLog
}
