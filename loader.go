package logger

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
func Debug(ctx context.Context, msg string, otherData ...interface{}) {
	if otherData == nil {
		singleton.Debug(ctx, msg)
		return
	}

	singleton.Debug(ctx, msg, zap.Any("otherData", otherData))
}

// Info logs a message at InfoLevel.
func Info(ctx context.Context, msg string, otherData ...interface{}) {
	if otherData == nil {
		singleton.Info(ctx, msg)
		return
	}

	singleton.Info(ctx, msg, zap.Any("otherData", otherData))
}

// Warn logs a message at WarnLevel.
func Warn(ctx context.Context, msg string, otherData ...interface{}) {
	if otherData == nil {
		singleton.Warn(ctx, msg)
		return
	}

	singleton.Warn(ctx, msg, zap.Any("otherData", otherData))
}

// Error logs a message at ErrorLevel.
func Error(ctx context.Context, msg string, err error, otherData ...interface{}) {
	if otherData == nil {
		singleton.Error(ctx, msg, err)
		return
	}

	singleton.Error(ctx, msg, err, zap.Any("otherData", otherData))
}

// Fatal logs a message at FatalLevel.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func Fatal(ctx context.Context, msg string, err error, otherData ...interface{}) {
	if otherData == nil {
		singleton.Fatal(ctx, msg, err)
		return
	}

	singleton.Fatal(ctx, msg, err, zap.Any("otherData", otherData))
}

// Panic logs a message at PanicLevel.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func Panic(ctx context.Context, msg string, err error, otherData ...interface{}) {
	if otherData == nil {
		singleton.Panic(ctx, msg, err)
		return
	}

	singleton.Panic(ctx, msg, err, zap.Any("otherData", otherData))
}

// TDR (Transaction Detail Request) consist of request and response log
func TDR(ctx context.Context, model LogModel) {
	singleton.TDR(ctx, model)
}
