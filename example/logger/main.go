package main

import (
	"flag"
	"strings"

	"github.com/x893675/gopkg/logger"
	"go.uber.org/zap"
)

func main() {
	filter := sampleFilter{}
	logger.InitFlags(nil)
	flag.Parse()
	logger.SetFilter(filter)
	logger.ApplyLogger()
	defer logger.FlushLogs()
	logger.Debug("hello world", zap.String("request_id", "1234567"))
	logger.Debugf("hello world, request_id=%s", "1234567")

	logger.Info("hello world filter me", zap.String("request_id", "1234567"))
	logger.Infof("hello world %s", "filter me")
	logger.Infof("hello world, filter me, request_id=%s", "1234567")

	logger.Warn("hello world", zap.String("request_id", "1234567"))
	logger.Warnf("hello world, request_id=%s", "1234567")

	logger.Error("hello world", zap.String("request_id", "1234567"))
	logger.Errorf("hello world, request_id=%s", "1234567")

	logger.Fatal("hello world", zap.String("request_id", "1234567"))
}

var _ logger.LogFilter = (*sampleFilter)(nil)

type sampleFilter struct{}

func (s sampleFilter) Filter(args []interface{}) []interface{} {
	for i, arg := range args {
		v, ok := arg.(string)
		if ok && strings.Contains(v, "filter me") {
			args[i] = "[FILTERED]"
		}
	}
	return args
}

func (s sampleFilter) FilterF(format string, args []interface{}) (string, []interface{}) {
	return strings.Replace(format, "filter me", "[FILTERED]", 1), s.Filter(args)
}
