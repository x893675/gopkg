package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var (
	logger *zapLogger
)

func init() {
	l, _ := zap.NewProduction(zap.AddCallerSkip(1))
	logger = &zapLogger{l: l}
}

// zapLogger is a logr.Logger that uses Zap to record log.
type zapLogger struct {
	// NB: this looks very similar to zap.SugaredLogger, but
	// deals with our desire to have multiple verbosity levels.
	l   *zap.Logger
	lvl int
}

func NewLoggerWithOptions(opt *Options) {
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
		syncer = zapcore.NewMultiWriteSyncer(os.Stderr)
	}
	core := zapcore.NewCore(
		newDefaultProductionLogEncoder(true, opt.Encode),
		syncer,
		level,
	)
	l := zap.New(core)
	if level == zapcore.DebugLevel {
		// caller skip set 1
		// 使得DEBUG模式下caller的值为调用当前package的代码路径
		l = l.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))
	}
	logger = &zapLogger{l: l}
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
	return logger.l.Core().Enabled(l)
}

func Info(msg string, fields ...zap.Field) {
	logger.l.Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	logger.l.Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.l.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.l.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.l.Fatal(msg, fields...)
}

func Infof(format string, args ...interface{}) {
	logger.l.Info(fmt.Sprintf(format, args...))
}

func Debugf(format string, args ...interface{}) {
	logger.l.Debug(fmt.Sprintf(format, args...))
}

func Warnf(format string, args ...interface{}) {
	logger.l.Warn(fmt.Sprintf(format, args...))
}

func Errorf(format string, args ...interface{}) {
	logger.l.Error(fmt.Sprintf(format, args...))
}

func Fatalf(format string, args ...interface{}) {
	logger.l.Fatal(fmt.Sprintf(format, args...))
}
