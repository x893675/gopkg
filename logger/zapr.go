package logger

import (
	"github.com/go-logr/logr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var _ logr.Logger = (*zapLogger)(nil)

func LogR() logr.Logger {
	return logger
}

func (z *zapLogger) Enabled() bool {
	return true
}

func (z *zapLogger) Info(msg string, keysAndValues ...interface{}) {
	entry := zapcore.Entry{
		Time:    time.Now(),
		Message: msg,
	}
	checkedEntry := z.l.Core().Check(entry, nil)
	checkedEntry.Write(z.handleFields(keysAndValues)...)
	//z.l.Info(msg, z.handleFields(keysAndValues)...)
}

func (z *zapLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	entry := zapcore.Entry{
		Level:   zapcore.ErrorLevel,
		Time:    time.Now(),
		Message: msg,
	}
	checkedEntry := z.l.Core().Check(entry, nil)
	checkedEntry.Write(z.handleFields(keysAndValues, handleError(err))...)
}

func (z *zapLogger) V(level int) logr.Logger {
	return &zapLogger{
		lvl: z.lvl + level,
		l:   z.l,
	}
}

func (z *zapLogger) WithValues(keysAndValues ...interface{}) logr.Logger {
	z.l = z.l.With(z.handleFields(keysAndValues)...)
	return z
}

func (z *zapLogger) WithName(name string) logr.Logger {
	z.l = z.l.Named(name)
	return z
}

// handleFields converts a bunch of arbitrary key-value pairs into Zap fields.  It takes
// additional pre-converted Zap fields, for use with automatically attached fields, like
// `error`.
func (z *zapLogger) handleFields(args []interface{}, additional ...zap.Field) []zap.Field {
	// a slightly modified version of zap.SugaredLogger.sweetenFields
	if len(args) == 0 {
		// fast-return if we have no suggared fields.
		return append(additional, zap.Int("v", z.lvl))
	}

	// unlike Zap, we can be pretty sure users aren't passing structured
	// fields (since logr has no concept of that), so guess that we need a
	// little less space.
	fields := make([]zap.Field, 0, len(args)/2+len(additional)+1)
	fields = append(fields, zap.Int("v", z.lvl))
	for i := 0; i < len(args)-1; i += 2 {
		// check just in case for strongly-typed Zap fields, which is illegal (since
		// it breaks implementation agnosticism), so we can give a better error message.
		if _, ok := args[i].(zap.Field); ok {
			z.dPanic("strongly-typed Zap Field passed to logr")
			break
		}

		// process a key-value pair,
		// ensuring that the key is a string
		key, val := args[i], args[i+1]
		keyStr, isString := key.(string)
		if !isString {
			// if the key isn't a string, stop logging
			z.dPanic("non-string key argument passed to logging, ignoring all later arguments")
			break
		}

		fields = append(fields, zap.Any(keyStr, val))
	}

	return append(fields, additional...)
}

func handleError(err error) zap.Field {
	return zap.NamedError("err", err)
}

func (z *zapLogger) dPanic(msg string) {
	entry := zapcore.Entry{
		Level:   zapcore.DPanicLevel,
		Time:    time.Now(),
		Message: msg,
	}
	checkedEntry := z.l.Core().Check(entry, nil)
	checkedEntry.Write(zap.Int("v", z.lvl))
}
