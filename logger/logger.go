package logger

import (
	"fmt"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var _logging = defaultLoggingT()

type loggingT struct {
	// Boolean flags. Not handled atomically because the flag.Value interface
	// does not let us avoid the =true, and that shorthand is necessary for
	// compatibility. TODO: does this matter enough to fix? Seems unlikely.
	toStderr     bool // The -logtostderr flag.
	alsoToStderr bool // The -alsologtostderr flag.

	// If non-empty, specifies the path of the file to write logs. mutually exclusive
	// with the log_dir option.
	logFile string

	// When logFile is specified, this limiter makes sure the logFile won't exceeds a certain size. When exceeds, the
	// logFile will be cleaned up. If this value is 0, no size limitation will be applied to logFile.
	logFileMaxSizeMB int

	flushInterval time.Duration

	// log level, debug,info,error
	Level zapcore.Level

	// encode type
	encodeType EncodeType

	// mu protects the remaining elements of this structure and is
	// used to synchronize logging.
	mu sync.Mutex

	// NB: this looks very similar to zap.SugaredLogger, but
	// deals with our desire to have multiple verbosity levels.
	l *zap.Logger
}

func (l *loggingT) ApplyZapLogger() {
	var multiWriteSyncer []zapcore.WriteSyncer

	if !l.toStderr {
		lumberJackLogger := &lumberjack.Logger{
			Filename:   l.logFile,
			MaxSize:    l.logFileMaxSizeMB,
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   false,
			LocalTime:  true,
		}
		multiWriteSyncer = append(multiWriteSyncer, zapcore.Lock(zapcore.AddSync(lumberJackLogger)))
	} else {
		multiWriteSyncer = append(multiWriteSyncer, os.Stderr)
	}

	if !l.toStderr && l.alsoToStderr {
		multiWriteSyncer = append(multiWriteSyncer, os.Stderr)
	}

	core := zapcore.NewCore(newDefaultProductionLogEncoder(true, l.encodeType),
		zapcore.NewMultiWriteSyncer(multiWriteSyncer...),
		l.Level)
	zl := zap.New(core)
	if l.Level == zapcore.DebugLevel {
		// caller skip set 1
		// 使得DEBUG模式下caller的值为调用当前package的代码路径
		zl = zl.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))
	}
	l.l = zl
}

// lockAndFlushAll is like flushAll but locks l.mu first.
func (l *loggingT) lockAndFlushAll() {
	l.mu.Lock()
	l.flushAll()
	l.mu.Unlock()
}

func (l *loggingT) flushAll() {
	_ = l.l.Sync()
}

func ApplyLogger() {
	_logging.mu.Lock()
	defer _logging.mu.Unlock()
	_logging.ApplyZapLogger()
}

func defaultLoggingT() *loggingT {
	l := &loggingT{
		toStderr:         true,
		alsoToStderr:     false,
		logFile:          "",
		logFileMaxSizeMB: 100,
		Level:            zapcore.InfoLevel,
		flushInterval:    5 * time.Second,
		encodeType:       ConsoleEncode,
	}
	//l.ApplyZapLogger()
	return l
}

func newDefaultProductionLogEncoder(colorize bool, encodeType EncodeType) zapcore.Encoder {
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.EncodeTime = func(ts time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(ts.Format("2006-01-02T15:04:05Z07:00"))
	}
	switch encodeType {
	case ConsoleEncode:
		encCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapcore.NewConsoleEncoder(encCfg)
	case JSONEncode:
		fallthrough
	default:
		return zapcore.NewJSONEncoder(encCfg)
	}
}

func Info(msg string, fields ...zap.Field) {
	_logging.l.Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	_logging.l.Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	_logging.l.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	_logging.l.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	_logging.l.Fatal(msg, fields...)
}

func Infof(format string, args ...interface{}) {
	_logging.l.Info(fmt.Sprintf(format, args...))
}

func Debugf(format string, args ...interface{}) {
	_logging.l.Debug(fmt.Sprintf(format, args...))
}

func Warnf(format string, args ...interface{}) {
	_logging.l.Warn(fmt.Sprintf(format, args...))
}

func Errorf(format string, args ...interface{}) {
	_logging.l.Error(fmt.Sprintf(format, args...))
}

func Fatalf(format string, args ...interface{}) {
	_logging.l.Fatal(fmt.Sprintf(format, args...))
}

func FlushLogs() {
	_logging.lockAndFlushAll()
}
