package golog

import (
	"go.uber.org/zap"
)

func NewStringField(key string, value string) zap.Field {
	return zap.String(key, value)
}

func NewIntField(key string, value int) zap.Field {
	return zap.Int(key, value)
}

func NewInt64Field(key string, value int64) zap.Field {
	return zap.Int64(key, value)
}

func NewObjectField(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

func NewArrayStringField(key string, value []string) zap.Field {
	return zap.Strings(key, value)
}

func NewBooleanField(key string, value bool) zap.Field {
	return zap.Bool(key, value)
}
