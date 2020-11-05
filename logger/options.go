package logger

import "github.com/spf13/pflag"

type Options struct {
	Level  string `json:"level" yaml:"level"`
	Output string `json:"output" yaml:"output"`
	Encode string `json:"encode" yaml:"encode"`
}

func NewLogOptions() *Options {
	return &Options{
		Level:  "info",
		Output: "stderr",
		Encode: "json",
	}
}

func (s *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.Level, "log-level", s.Level, "log level (debug,info,warn,error)")
	fs.StringVar(&s.Output, "log-output", s.Output, "log output (stdout or stderr)")
	fs.StringVar(&s.Encode, "log-encode", s.Encode, "log encoder (raw or json)")
}

func (s *Options) Validate() []error {
	return nil
}
