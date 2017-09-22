package glog

import (
	"io"
	"os"
	"syscall"
)

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

type dump struct{}

func (this *dump) Write(p []byte) (n int, err error) {
	logging.printf(dumpLog, " %v", string(p))
	return 0, nil
}

func Dump() io.Writer {
	RedirectStderr()
	return new(dump)
}

func RedirectStderr() {
	errFile, _ := os.OpenFile(logDir+"/err."+program, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0644)
	syscall.Dup2(int(errFile.Fd()), 2)
}
