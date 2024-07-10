package golog

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var SENSITIVE_HEADER = []string{
	"Authorization",
	"Signature",
	"Apikey",
	"Content-Disposition",
}

var SENSITIVE_ATTR = map[string]bool{
	"password":                 true,
	"license":                  true,
	"license_code":             true,
	"token":                    true,
	"access_token":             true,
	"refresh_token":            true,
	"bank_ac_no":               true, // bank account number
	"id_number":                true, // id number, equivalent to KTP in Indonesia
	"mobile":                   true, // mobile number
	"npwp":                     true, // tax number, equivalent to NPWP in Indonesia
	"phone":                    true, // phone number
	"card_no":                  true, // social security number
	"basic_salary":             true, // basic salary payslip
	"brutto":                   true, // brutto payslip
	"employment_deduction":     true, // employment deduction payslip
	"functional_allowance":     true, // functional allowance payslip
	"health_allowance":         true, // health allowance payslip
	"health_deduction":         true, // health deduction payslip
	"incentive_income_tax_21":  true, // incentive income tax payslip
	"income_tax_21":            true, // income tax payslip
	"jht_allowance":            true, // jht allowance payslip
	"jkk_allowance":            true, // jkk allowance payslilp
	"jkm_allowance":            true, // jkm allowance payslip
	"jkn_allowance":            true, // jkn allowance payslip
	"loan":                     true, // loan payslip
	"other_deduction":          true, // other deduction payslip
	"position_allowance":       true, // position allowance payslip
	"skill_allowance":          true, // skill allowance payslip
	"special_region_allowance": true, // special region allowance payslip
	"take_home_pay":            true, // take home pay payslip
	"total_allowance":          true, // total allowance payslip
	"total_deduction":          true, // total deduction payslip
	"total_wages":              true, // total wages payslip
}

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
	encoderConfig.StacktraceKey = "stacktrace"
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

	appVer := conf.AppVer

	content, err := os.ReadFile("version.txt")
	if err == nil {
		c := string(content)
		appVer = strings.TrimSuffix(c, "\n")
	}

	logger := zap.New(core, zap.AddStacktrace(zap.ErrorLevel), zap.AddCallerSkip(2)).With(
		zap.String("app", conf.App),
		zap.String("appVer", appVer),
		zap.String("env", conf.Env),
	)

	loggerTDR := zap.New(coreTDR, zap.AddCallerSkip(2)).With(
		zap.String("app", conf.App),
		zap.String("appVer", appVer),
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
	fields = append(fields, zap.Any("request", toJSON(maskField(log.Request))))
	fields = append(fields, zap.String("statusCode", log.StatusCode))
	fields = append(fields, zap.Uint64("httpStatus", log.HttpStatus))
	fields = append(fields, zap.Any("response", toJSON(maskField(log.Response))))
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
	// Fasthttp
	if mapHeader, ok := header.(fasthttp.RequestHeader); ok {
		for _, val := range SENSITIVE_HEADER {
			mapHeader.Del(val)
		}
		return string(mapHeader.Header())
	}

	// Http
	if mapHeader, ok := header.(http.Header); ok {
		for _, val := range SENSITIVE_HEADER {
			mapHeader.Del(val)
		}
	}

	return header
}

func maskField(body interface{}) interface{} {
	if bodyByte, ok := body.([]byte); ok {
		bodyMap := make(map[string]interface{}, 0)
		if err := json.Unmarshal(bodyByte, &bodyMap); err != nil {
			return string(bodyByte)
		}

		for key, value := range bodyMap {
			switch value.(type) {
			case map[string]interface{}:
				valueByte, _ := json.Marshal(value)
				bodyMap[key] = maskField(valueByte)
			default:
				if isSensitiveField(key) {
					bodyMap[key] = strings.Repeat("*", 5)
				} else {
					bodyMap[key] = value
				}
			}
		}

		return bodyMap
	}

	return body
}

func isSensitiveField(key string) bool {
	if _, ok := SENSITIVE_ATTR[strings.ToLower(key)]; ok {
		return true
	}
	return false
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
