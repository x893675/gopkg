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
	flagset.Var(&_logging.encodeType, "encode-type", "the number of the log encode type, console or json")
	flagset.IntVar(&_logging.maxBackups, "max-backups", _logging.maxBackups, ""+
		"MaxBackups is the maximum number of old log files to retain."+
		"The default is to retain all old log files (though MaxAge may still cause them to get deleted.)")
	flagset.IntVar(&_logging.maxAge, "max-age", _logging.maxAge, ""+
		"MaxAge is the maximum number of days to retain old log files based on the timestamp encoded in their filename. "+
		"Note that a day is defined as 24 hours and may not exactly correspond to calendar days due to daylight savings, "+
		"leap seconds, etc. The default is not to remove old log files based on age.")
	flagset.BoolVar(&_logging.compress, "compress", _logging.compress, ""+
		"Compress determines if the rotated log files should be compressed using gzip. The default is not to perform compression.")
	flagset.BoolVar(&_logging.useLocalTimeBack, "use-localtime", _logging.useLocalTimeBack, ""+
		"LocalTime determines if the time used for formatting the timestamps in backup files is the computer's local time. "+
		"false mean to use UTC time.")
}
