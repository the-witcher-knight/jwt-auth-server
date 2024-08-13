package tracing

import (
	"go.uber.org/zap"
)

type AttributeType uint

const (
	AttributeTypeString AttributeType = iota
	AttributeTypeInt
)

type Attribute struct {
	Key   string
	Value interface{}
	Type  AttributeType
}

func String(k, v string) Attribute {
	return Attribute{k, v, AttributeTypeString}
}

func Int(k string, v int) Attribute {
	return Attribute{k, v, AttributeTypeInt}
}

func toZapField(attr Attribute) zap.Field {
	switch attr.Type {
	case AttributeTypeString:
		return zap.String(attr.Key, attr.Value.(string))
	case AttributeTypeInt:
		return zap.Int(attr.Key, attr.Value.(int))
	default:
		return zap.Any(attr.Key, attr.Value)
	}
}
