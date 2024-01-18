package golog

import (
	"context"

	"go.uber.org/zap"
)

type LoggerInterface interface {
	Debug(ctx context.Context, message string, fields ...zap.Field)
	Info(ctx context.Context, message string, fields ...zap.Field)
	Warn(ctx context.Context, message string, fields ...zap.Field)
	Error(ctx context.Context, message string, err error, fields ...zap.Field)
	Fatal(ctx context.Context, message string, err error, fields ...zap.Field)
	Panic(ctx context.Context, message string, err error, fields ...zap.Field)
	TDR(ctx context.Context, tdr LogModel)
}
