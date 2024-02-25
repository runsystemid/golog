package golog

import (
	"context"
	"sync"

	"go.uber.org/zap"
)

var (
	once      sync.Once
	singleton LoggerInterface
)

// Loader to construct logger. So logger will
func Load(config Config) LoggerInterface {
	once.Do(func() {
		singleton = NewLogger(config)
	})

	return singleton
}

// Debug logs a message at DebugLevel.
func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	singleton.Debug(ctx, msg, fields...)
}

// Info logs a message at InfoLevel.
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	singleton.Info(ctx, msg, fields...)
}

// Warn logs a message at WarnLevel.
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	singleton.Warn(ctx, msg, fields...)
}

// Error logs a message at ErrorLevel.
func Error(ctx context.Context, msg string, err error, fields ...zap.Field) {
	singleton.Error(ctx, msg, err, fields...)
}

// Fatal logs a message at FatalLevel.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func Fatal(ctx context.Context, msg string, err error, fields ...zap.Field) {
	singleton.Fatal(ctx, msg, err, fields...)
}

// Panic logs a message at PanicLevel.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func Panic(ctx context.Context, msg string, err error, fields ...zap.Field) {
	singleton.Panic(ctx, msg, err, fields...)
}

// TDR (Transaction Detail Request) consist of request and response log
func TDR(ctx context.Context, model LogModel) {
	singleton.TDR(ctx, model)
}
