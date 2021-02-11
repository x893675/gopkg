package logger

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogFilter is a collection of functions that can filter all logging calls,
// e.g. for sanitization of arguments and prevent accidental leaking of secrets.
type LogFilter interface {
	Filter(args []interface{}) []interface{}
	FilterF(format string, args []interface{}) (string, []interface{})
	FilterS(msg string, keysAndValues []interface{}) (string, []interface{})
}

var logging = defaultLoggingT()

func defaultLoggingT() *loggingT {
	l := &loggingT{
		toStderr:         true,
		alsoToStderr:     false,
		logDir:           "",
		logFile:          "",
		logFileMaxSizeMB: 100,
		skipHeaders:      false,
		addDirHeader:     false,
		Level:            zapcore.InfoLevel,
	}
	l.initZapLogger()
	return l
}

func InitLogger() {
	logging.mu.Lock()
	defer logging.mu.Unlock()
	logging.initZapLogger()
}

func (l *loggingT) initZapLogger() {
	var multiWriteSyncer []zapcore.WriteSyncer

	if !l.toStderr {
		lumberJackLogger := &lumberjack.Logger{
			Filename:   "./test.log",
			MaxSize:    l.logFileMaxSizeMB,
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   false,
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

type loggingT struct {
	// Boolean flags. Not handled atomically because the flag.Value interface
	// does not let us avoid the =true, and that shorthand is necessary for
	// compatibility. TODO: does this matter enough to fix? Seems unlikely.
	toStderr     bool // The -logtostderr flag.
	alsoToStderr bool // The -alsologtostderr flag.
	// If non-empty, overrides the choice of directory in which to write logs.
	// See createLogDirs for the full list of possible destinations.
	logDir string

	// If non-empty, specifies the path of the file to write logs. mutually exclusive
	// with the log_dir option.
	logFile string

	// When logFile is specified, this limiter makes sure the logFile won't exceeds a certain size. When exceeds, the
	// logFile will be cleaned up. If this value is 0, no size limitation will be applied to logFile.
	logFileMaxSizeMB int

	// If true, do not add the prefix headers, useful when used with SetOutput
	skipHeaders bool

	// If true, do not add the headers to log files
	skipLogHeaders bool

	// If true, add the file directory to the header
	addDirHeader bool

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

// InitFlags is for explicitly initializing the flags.
func InitFlags(flagset *flag.FlagSet) {
	if flagset == nil {
		flagset = flag.CommandLine
	}

	flagset.StringVar(&logging.logDir, "log_dir", logging.logDir, "If non-empty, write log files in this directory")
	flagset.StringVar(&logging.logFile, "log_file", logging.logFile, "If non-empty, use this log file")
	flagset.IntVar(&logging.logFileMaxSizeMB, "log_file_max_size", logging.logFileMaxSizeMB,
		"Defines the maximum size a log file can grow to. Unit is megabytes. "+
			"If the value is 0, the maximum file size is unlimited.")
	flagset.BoolVar(&logging.toStderr, "logtostderr", logging.toStderr, "log to standard error instead of files")
	flagset.BoolVar(&logging.alsoToStderr, "alsologtostderr", logging.alsoToStderr, "log to standard error as well as files")
	flagset.Var(&logging.Level, "level", "the number of the log level verbosity")
	flagset.Var(&logging.encodeType, "encode_type", "the number of the log encode type")
	flagset.BoolVar(&logging.addDirHeader, "add_dir_header", logging.addDirHeader, "If true, adds the file directory to the header of the log messages")
	flagset.BoolVar(&logging.skipHeaders, "skip_headers", logging.skipHeaders, "If true, avoid header prefixes in the log messages")
	flagset.BoolVar(&logging.skipLogHeaders, "skip_log_headers", logging.skipLogHeaders, "If true, avoid headers when opening log files")
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
	logging.l.Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	logging.l.Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logging.l.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logging.l.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logging.l.Fatal(msg, fields...)
}

func Infof(format string, args ...interface{}) {
	logging.l.Info(fmt.Sprintf(format, args...))
}

func Debugf(format string, args ...interface{}) {
	logging.l.Debug(fmt.Sprintf(format, args...))
}

func Warnf(format string, args ...interface{}) {
	logging.l.Warn(fmt.Sprintf(format, args...))
}

func Errorf(format string, args ...interface{}) {
	logging.l.Error(fmt.Sprintf(format, args...))
}

func Fatalf(format string, args ...interface{}) {
	logging.l.Fatal(fmt.Sprintf(format, args...))
}
