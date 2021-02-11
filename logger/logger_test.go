package logger

// import (
// 	"go.uber.org/zap"
// 	"testing"
// )

// func initLoggerWithRaw() {
// 	NewLoggerWithOptions(&Options{
// 		Level:  "info",
// 		Output: "stdout",
// 		Encode: "raw",
// 	})
// }

// func BenchmarkInfoWithJson(b *testing.B) {
// 	b.ResetTimer()
// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			Info("test",
// 				zap.Int64("int64-1", int64(1)),
// 				zap.Int64("int64-2", int64(2)),
// 				zap.Float64("float64", 1.0),
// 				zap.String("string1", "\n"),
// 				zap.String("string2", "💩"),
// 				zap.String("string3", "🤔"),
// 				zap.String("string4", "🙊"),
// 				zap.Bool("bool", true),
// 				zap.Any("request", struct {
// 					Method  string `json:"method"`
// 					Timeout int    `json:"timeout"`
// 					secret  string `json:"secret"`
// 				}{
// 					Method:  "GET",
// 					Timeout: 10,
// 					secret:  "pony",
// 				}))
// 		}
// 	})
// }

// func BenchmarkInfoWithRaw(b *testing.B) {
// 	initLoggerWithRaw()
// 	b.ResetTimer()
// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			Info("test",
// 				zap.Int64("int64-1", int64(1)),
// 				zap.Int64("int64-2", int64(2)),
// 				zap.Float64("float64", 1.0),
// 				zap.String("string1", "\n"),
// 				zap.String("string2", "💩"),
// 				zap.String("string3", "🤔"),
// 				zap.String("string4", "🙊"),
// 				zap.Bool("bool", true),
// 				zap.Any("request", struct {
// 					Method  string `json:"method"`
// 					Timeout int    `json:"timeout"`
// 					secret  string `json:"secret"`
// 				}{
// 					Method:  "GET",
// 					Timeout: 10,
// 					secret:  "pony",
// 				}))
// 		}
// 	})
// }
