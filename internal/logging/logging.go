package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type Logger struct {
	slog *slog.Logger
}

func NewLogger() (*Logger, error) {
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return &Logger{
		slog: l,
	}, nil
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...Field) {
	attrs := make([]interface{}, len(fields))
	for i, f := range fields {
		attrs[i] = toSlogField(f)
	}

	l.slog.Info(msg, attrs...)
}

func (l *Logger) Error(ctx context.Context, err error, msg string, fields ...Field) {
	fields = append(fields,
		AttributeString("error.message", err.Error()),
		AttributeString("error.type", fmt.Sprintf("%T", err)),
	)

	attrs := make([]interface{}, len(fields))
	for i, f := range fields {
		attrs[i] = toSlogField(f)
	}

	l.slog.Error(msg, attrs...)
}

func (l *Logger) With(fields ...Field) *Logger {
	attrs := make([]interface{}, len(fields))
	for i, f := range fields {
		attrs[i] = toSlogField(f)
	}

	return &Logger{
		slog: l.slog.With(attrs...),
	}
}

func NewNoop() *Logger {
	return &Logger{
		slog: &slog.Logger{},
	}
}

type Field struct {
	Key   string
	Value interface{}
	Type  FieldType
}

type FieldType uint8

const (
	FieldTypeString FieldType = iota
	FieldTypeInt
)

func AttributeString(key, value string) Field {
	return Field{
		Key:   key,
		Value: value,
		Type:  FieldTypeString,
	}
}

func toSlogField(f Field) slog.Attr {
	switch f.Type {
	case FieldTypeString:
		return slog.String(f.Key, f.Value.(string))
	case FieldTypeInt:
		return slog.Int(f.Key, f.Value.(int))
	default:
		return slog.Any(f.Key, f.Value)
	}
}
