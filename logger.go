package golog

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Log struct {
	logger    *zap.Logger
	loggerTDR *zap.Logger
}

func NewLogger(conf Config) LoggerInterface {
	rotator := &lumberjack.Logger{
		Filename:   conf.FileLocation + "/system.log",
		MaxSize:    conf.FileMaxSize, // megabytes
		MaxBackups: conf.FileMaxBackup,
		MaxAge:     conf.FileMaxAge, // days
	}

	rotatorTDR := &lumberjack.Logger{
		Filename:   conf.FileLocation + "/tdr.log",
		MaxSize:    conf.FileMaxSize, // megabytes
		MaxBackups: conf.FileMaxBackup,
		MaxAge:     conf.FileMaxAge, // days
	}

	encoderConfig := zap.NewDevelopmentEncoderConfig()

	if conf.Env == "production" {
		encoderConfig = zap.NewProductionEncoderConfig()
	}

	encoderConfig.TimeKey = "timestamp"
	encoderConfig.LevelKey = "logLevel"
	encoderConfig.MessageKey = "message"
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewCore(
		jsonEncoder,
		zapcore.AddSync(rotator),
		zap.NewAtomicLevelAt(zap.InfoLevel),
	)

	coreTDR := zapcore.NewCore(
		jsonEncoder,
		zapcore.AddSync(rotatorTDR),
		zap.NewAtomicLevelAt(zap.InfoLevel),
	)

	if conf.Stdout {
		core = zapcore.NewTee(
			core,
			zapcore.NewCore(
				consoleEncoder,
				zapcore.AddSync(os.Stdout),
				zap.NewAtomicLevelAt(zap.InfoLevel),
			),
		)

		coreTDR = zapcore.NewTee(
			coreTDR,
			zapcore.NewCore(
				consoleEncoder,
				zapcore.AddSync(os.Stdout),
				zap.NewAtomicLevelAt(zap.InfoLevel),
			),
		)
	}

	logger := zap.New(core, zap.AddStacktrace(zap.ErrorLevel), zap.AddCallerSkip(2)).With(
		zap.String("app", conf.App),
		zap.String("appVer", conf.AppVer),
		zap.String("env", conf.Env),
	)

	loggerTDR := zap.New(coreTDR, zap.AddStacktrace(zap.ErrorLevel), zap.AddCallerSkip(2)).With(
		zap.String("app", conf.App),
		zap.String("appVer", conf.AppVer),
		zap.String("env", conf.Env),
	)

	return &Log{
		logger:    logger,
		loggerTDR: loggerTDR,
	}
}

func (l *Log) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	ctxField := populateFieldFromContext(ctx)
	fields = append(fields, ctxField...)
	l.logger.Debug(msg, fields...)
	defer l.logger.Sync()
}

func (l *Log) Info(ctx context.Context, msg string, fields ...zap.Field) {
	ctxField := populateFieldFromContext(ctx)
	fields = append(fields, ctxField...)
	l.logger.Info(msg, fields...)
	defer l.logger.Sync()
}

func (l *Log) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	ctxField := populateFieldFromContext(ctx)
	fields = append(fields, ctxField...)
	l.logger.Warn(msg, fields...)
	defer l.logger.Sync()
}

func (l *Log) Error(ctx context.Context, msg string, err error, fields ...zap.Field) {
	ctxField := populateFieldFromContext(ctx)
	fields = append(fields, ctxField...)
	fields = append(fields, zap.Any("error", toJSON(err)))
	l.logger.Error(msg, fields...)
	defer l.logger.Sync()
}

func (l *Log) Fatal(ctx context.Context, msg string, err error, fields ...zap.Field) {
	ctxField := populateFieldFromContext(ctx)
	fields = append(fields, ctxField...)
	fields = append(fields, zap.Any("error", toJSON(err)))
	l.logger.Fatal(msg, fields...)
	defer l.logger.Sync()
}

func (l *Log) Panic(ctx context.Context, msg string, err error, fields ...zap.Field) {
	ctxField := populateFieldFromContext(ctx)
	fields = append(fields, ctxField...)
	fields = append(fields, zap.Any("error", toJSON(err)))
	l.logger.Panic(msg, fields...)
	defer l.logger.Sync()
}

func (l *Log) TDR(ctx context.Context, log LogModel) {
	fields := populateFieldFromContext(ctx)

	fields = append(fields, zap.String("correlationID", log.CorrelationID))
	fields = append(fields, zap.Any("header", removeAuth(log.Header)))
	fields = append(fields, zap.Any("request", toJSON(sanitizeBody(log.Request))))
	fields = append(fields, zap.String("statusCode", log.StatusCode))
	fields = append(fields, zap.Uint64("httpStatus", log.HttpStatus))
	fields = append(fields, zap.Any("response", toJSON(log.Response)))
	fields = append(fields, zap.Int64("rt", log.ResponseTime.Milliseconds()))
	fields = append(fields, zap.Any("error", toJSON(log.Error)))
	fields = append(fields, zap.Any("otherData", toJSON(log.OtherData)))

	l.loggerTDR.Info(":", fields...)

	defer l.loggerTDR.Sync()
}

func toJSON(object interface{}) interface{} {
	if object == nil {
		return nil
	}

	if w, ok := object.(string); ok {
		var jsonobj map[string]interface{}
		if err := json.Unmarshal([]byte(w), &jsonobj); err != nil {
			return w
		}
		return jsonobj
	}
	return object
}

func removeAuth(header interface{}) interface{} {
	if mapHeader, ok := header.(map[string]string); ok {
		delete(mapHeader, "Authorization")
		return mapHeader
	}

	return header
}

func sanitizeBody(reqBody interface{}) interface{} {
	if reqByte, ok := reqBody.([]byte); ok {
		bodyMap := make(map[string]interface{}, 0)
		if err := json.Unmarshal(reqByte, &bodyMap); err != nil {
			return string(reqByte)
		}

		for key, value := range bodyMap {
			if strings.Contains(key, "password") {
				bodyMap[key] = strings.Repeat("*", len(value.(string)))
			} else {
				bodyMap[key] = value
			}
		}

		return bodyMap
	}

	return reqBody
}

func populateFieldFromContext(ctx context.Context) []zap.Field {
	fieldFromCtx := make([]zap.Field, 0)

	if v, ok := ctx.Value("traceId").(string); ok {
		fieldFromCtx = append(fieldFromCtx, zap.String("traceId", v))
	}

	if v, ok := ctx.Value("srcIP").(string); ok {
		fieldFromCtx = append(fieldFromCtx, zap.String("srcIP", v))
	}

	if v, ok := ctx.Value("port").(string); ok {
		fieldFromCtx = append(fieldFromCtx, zap.String("port", v))
	}

	if v, ok := ctx.Value("path").(string); ok {
		fieldFromCtx = append(fieldFromCtx, zap.String("path", v))
	}

	return fieldFromCtx
}
