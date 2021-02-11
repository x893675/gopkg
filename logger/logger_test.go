package logger

import (
	"testing"

	"go.uber.org/zap"
)

func BenchmarkInfoWithConsoleEncode(b *testing.B) {
	ApplyLogger()
	defer FlushLogs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for pb.Next() {
				Info("test",
					zap.Int64("int64-1", int64(1)),
					zap.Int64("int64-2", int64(2)),
					zap.Float64("float64", 1.0),
					zap.String("string1", "\n"),
					zap.String("string2", "💩"),
					zap.String("string3", "🤔"),
					zap.String("string4", "🙊"),
					zap.Bool("bool", true),
					zap.Any("request", struct {
						Method  string `json:"method"`
						Timeout int    `json:"timeout"`
						secret  string `json:"secret"`
					}{
						Method:  "GET",
						Timeout: 10,
						secret:  "pony",
					}))
			}
		}
	})
	b.StopTimer()
}

func BenchmarkInfoWithJsonEncode(b *testing.B) {
	_logging.encodeType = JSONEncode
	ApplyLogger()
	defer FlushLogs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for pb.Next() {
				Info("test",
					zap.Int64("int64-1", int64(1)),
					zap.Int64("int64-2", int64(2)),
					zap.Float64("float64", 1.0),
					zap.String("string1", "\n"),
					zap.String("string2", "💩"),
					zap.String("string3", "🤔"),
					zap.String("string4", "🙊"),
					zap.Bool("bool", true),
					zap.Any("request", struct {
						Method  string `json:"method"`
						Timeout int    `json:"timeout"`
						secret  string `json:"secret"`
					}{
						Method:  "GET",
						Timeout: 10,
						secret:  "pony",
					}))
			}
		}
	})
	b.StopTimer()
}
