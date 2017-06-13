// modified from glog, see https://github.com/golang/glog for more details
package glog

import (
	"fmt"
	"os"
	"time"
)

// bufferSize sizes the buffer associated with each log file. It's large
// so that log records can accumulate without the logging thread blocking
// on disk I/O. The flushDaemon will block instead.
const bufferSize = 256 * 1024

const flushInterval = 1 * time.Second

const digits = "0123456789"

// Stats tracks the number of lines of output and number of bytes
// per severity level. Values must be read with atomic.LoadInt64.
var Stats struct {
	Info, Warning, Error OutputStats
}

var severityStats = [numSeverity]*OutputStats{
	infoLog:    &Stats.Info,
	warningLog: &Stats.Warning,
	errorLog:   &Stats.Error,
}

var logging loggingT

var timeNow = time.Now // Stubbed out for testing.

func init() {
	// Default stderrThreshold is ERROR.
	logging.stderrThreshold = errorLog

	logging.setVState(0, nil, false)
	go logging.flushDaemon()
}

// Flush flushes all pending log I/O.
func Flush() {
	logging.lockAndFlushAll()
}

// timeoutFlush calls Flush and returns when it completes or after timeout
// elapses, whichever happens first.  This is needed because the hooks invoked
// by Flush may deadlock when glog.Fatal is called from a hook that holds
// a lock.
func timeoutFlush(timeout time.Duration) {
	done := make(chan bool, 1)
	go func() {
		Flush() // calls logging.lockAndFlushAll()
		done <- true
	}()
	select {
	case <-done:
	case <-time.After(timeout):
		fmt.Fprintln(os.Stderr, "glog: Flush took longer than", timeout)
	}
}
