package golog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewStringField(key string, value string) zap.Field {
	return zap.Field{
		Key:       key,
		Type:      zapcore.StringType,
		Interface: value,
	}
}

func NewIntField(key string, value int) zap.Field {
	return zap.Field{
		Key:       key,
		Type:      zapcore.Int32Type,
		Interface: int64(value),
	}
}

func NewInt64Field(key string, value int64) zap.Field {
	return zap.Field{
		Key:       key,
		Type:      zapcore.Int64Type,
		Interface: value,
	}
}

func NewObjectField(key string, value interface{}) zap.Field {
	return zap.Field{
		Key:       key,
		Type:      zapcore.ObjectMarshalerType,
		Interface: value,
	}
}

func NewArrayField(key string, value []interface{}) zap.Field {
	return zap.Field{
		Key:       key,
		Type:      zapcore.ArrayMarshalerType,
		Interface: value,
	}
}

func NewBooleanField(key string, value bool) zap.Field {
	return zap.Field{
		Key:       key,
		Type:      zapcore.BoolType,
		Interface: value,
	}
}
