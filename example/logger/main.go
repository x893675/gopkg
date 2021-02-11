package main

import (
	"flag"

	"github.com/x893675/gopkg/logger"
	"go.uber.org/zap"
)

func main() {
	// logger.NewLoggerWithOptions(&logger.Options{
	// 	Level:  "debug",
	// 	Output: "stderr",
	// 	Encode: "json",
	// })
	logger.InitFlags(nil)
	flag.Parse()
	logger.InitLogger()
	logger.Debug("hello world", zap.String("request_id", "1234567"))
	logger.Debugf("hello world, request_id=%s", "1234567")

	logger.Info("hello world", zap.String("request_id", "1234567"))
	logger.Infof("hello world, request_id=%s", "1234567")

	logger.Warn("hello world", zap.String("request_id", "1234567"))
	logger.Warnf("hello world, request_id=%s", "1234567")

	logger.Error("hello world", zap.String("request_id", "1234567"))
	logger.Errorf("hello world, request_id=%s", "1234567")

	logger.Fatal("hello world", zap.String("request_id", "1234567"))
}
