package glog

func Inf(format string, args ...interface{}) {
	logging.printf(infoLog, format, args...)
}

func Dbg(format string, args ...interface{}) {
	logging.printf(debugLog, format, args...)
}

func Wrn(format string, args ...interface{}) {
	logging.printf(warningLog, format, args...)
}

func Err(format string, args ...interface{}) {
	logging.printf(errorLog, format, args...)
}

func Ftl(format string, args ...interface{}) {
	logging.printf(fatalLog, format, args...)
}
