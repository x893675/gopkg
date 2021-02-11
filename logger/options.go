package logger

import "flag"

// InitFlags is for explicitly initializing the flags.
func InitFlags(flagset *flag.FlagSet) {
	if flagset == nil {
		flagset = flag.CommandLine
	}

	flagset.StringVar(&_logging.logFile, "log-file", _logging.logFile, "If non-empty, use this log file")
	flagset.IntVar(&_logging.logFileMaxSizeMB, "log-file-max-size", _logging.logFileMaxSizeMB,
		"Defines the maximum size a log file can grow to. Unit is megabytes. "+
			"If the value is 0, the maximum file size is unlimited.")
	flagset.BoolVar(&_logging.toStderr, "logtostderr", _logging.toStderr, "log to standard error instead of files")
	flagset.BoolVar(&_logging.alsoToStderr, "alsologtostderr", _logging.alsoToStderr, "log to standard error as well as files")
	flagset.Var(&_logging.Level, "level", "the number of the log level verbosity")
	flagset.Var(&_logging.encodeType, "encode-type", "the number of the log encode type")
	flagset.DurationVar(&_logging.flushInterval, "log-flush-frequency", _logging.flushInterval, "Maximum number of seconds between log flushes")
}
