package main

import (
	"errors"
	"github.com/x893675/gopkg/logger"
)

func main() {
	logger.NewLoggerWithOptions(&logger.Options{
		Level:  "info",
		Output: "stderr",
		Encode: "raw",
	})
	l := logger.LogR()
	l.Info("hello world", "key", 1, "key2", "key2")
	l.Error(errors.New("this is error"), "this is message", "key", 1, "key2", "key2")
}
