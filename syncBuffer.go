package glog

import (
	"bufio"
	"io"
	"os"
	"strings"
	"time"
)

// flushSyncWriter is the interface satisfied by logging destinations.
type flushSyncWriter interface {
	Flush() error
	Sync() error
	io.Writer
}

// syncBuffer joins a bufio.Writer to its underlying file, providing access to the
// file's Sync method and providing a wrapper for the Write method that provides log
// file rotation. There are conflicting methods, so the file cannot be embedded.
// l.mu is held for all its methods.
type syncBuffer struct {
	logger *loggingT
	*bufio.Writer
	file     *os.File
	sev      severity
	nbytes   uint64 // The number of bytes written to this file
	lastHour int
}

func (sb *syncBuffer) Sync() error {
	return sb.file.Sync()
}

func (sb *syncBuffer) Write(p []byte) (n int, err error) {
	now := time.Now()
	curHour := int(now.Unix()) / 3600
	if sb.lastHour != curHour {
		if err := sb.rotateFile(now); err != nil {
			sb.logger.exit(err)
		}
		sb.lastHour = curHour
	}
	n, err = sb.Writer.Write(p)
	sb.nbytes += uint64(n)
	if err != nil {
		sb.logger.exit(err)
	}
	return
}

// rotateFile closes the syncBuffer's file and starts a new one.
func (sb *syncBuffer) rotateFile(now time.Time) error {
	if sb.file != nil {
		sb.Flush()
		sb.file.Close()
	}
	var err error
	sb.file, _, err = create(strings.ToLower(severityName[sb.sev]), now)
	sb.nbytes = 0
	if err != nil {
		return err
	}

	sb.Writer = bufio.NewWriterSize(sb.file, bufferSize)
	return err
}
