package logger

import (
	"context"
	"fmt"
	ctxpkg "github.com/x893675/gopkg/ctx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

const (
	TraceIDKey   = "trace_id"
	UserIDKey    = "user_id"
	RequestIDKey = "request_id"
	StackKey     = "stack"
)

var (
	logger *zap.Logger
)

func init() {
	logger, _ = zap.NewProduction(zap.AddCallerSkip(1))
}

func Logger() *zap.Logger {
	return logger
}

func NewLogger(opt *Options) {
	var level zapcore.Level
	switch opt.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "info":
		fallthrough
	default:
		level = zapcore.InfoLevel
	}
	var syncer zapcore.WriteSyncer
	switch opt.Output {
	case "stdout":
		syncer = zapcore.NewMultiWriteSyncer(os.Stdout)
	case "stderr":
		fallthrough
	default:
		syncer = zapcore.NewMultiWriteSyncer(os.Stdout)
	}
	core := zapcore.NewCore(
		newDefaultProductionLogEncoder(true, opt.Encode),
		syncer,
		level,
	)
	logger = zap.New(core)
	if level == zapcore.DebugLevel {
		// caller skip set 1
		// 使得DEBUG模式下caller的值为调用当前package的代码路径
		logger = logger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))
	}
}

func newDefaultProductionLogEncoder(colorize bool, encodeType string) zapcore.Encoder {
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.EncodeTime = func(ts time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(ts.Format("2006/01/02 15:04:05.000"))
	}
	switch encodeType {
	case "raw":
		encCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapcore.NewConsoleEncoder(encCfg)
	case "json":
		fallthrough
	default:
		return zapcore.NewJSONEncoder(encCfg)
	}
}

func Enabled(l zapcore.Level) bool {
	return logger.Core().Enabled(l)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	f := ExtraContext(ctx)
	logger.Info(msg, append(fields, f...)...)
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	f := ExtraContext(ctx)
	logger.Debug(msg, append(fields, f...)...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	f := ExtraContext(ctx)
	logger.Warn(msg, append(fields, f...)...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	f := ExtraContext(ctx)
	logger.Error(msg, append(fields, f...)...)
}

func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	f := ExtraContext(ctx)
	logger.Fatal(msg, append(fields, f...)...)
}

func ExtraContext(ctx context.Context) []zap.Field {
	if ctx == nil {
		return nil
	}
	var fields []zap.Field
	if v := ctxpkg.FromTraceIDContext(ctx); v != "" {
		fields = append(fields, zap.String(TraceIDKey, v))
	}
	if v := ctxpkg.FromUserIDContext(ctx); v != "" {
		fields = append(fields, zap.String(UserIDKey, v))
	}
	if v := ctxpkg.FromRequestIDContext(ctx); v != "" {
		fields = append(fields, zap.String(RequestIDKey, v))
	}
	if v := ctxpkg.FromStackContext(ctx); v != nil {
		fields = append(fields, zap.String(StackKey, fmt.Sprintf("%+v", v)))
	}
	return fields
}

func Infof(format string, args ...interface{}) {
	logger.Info(fmt.Sprintf(format, args...))
}
