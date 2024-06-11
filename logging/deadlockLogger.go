package logging

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/xboshy/go-deadlock"
)

type deadlockLogger struct {
	Logger
	*bytes.Buffer
	bufferSync     chan struct{}
	panic          func()
	reportDeadlock sync.Once
}

// Panic is defined here just so we can emulate the usage of the deadlockLogger
func (dlogger *deadlockLogger) Panic() {
	dlogger.Logger.Panic("potential deadlock detected")
}

// Write implements the io.Writer interface, ensuring that the write is syncronized.
func (dlogger *deadlockLogger) Write(p []byte) (n int, err error) {
	dlogger.bufferSync <- struct{}{}
	n, err = dlogger.Buffer.Write(p)
	<-dlogger.bufferSync
	return
}

// captureCallstack captures the callstack and return a byte array of the output.
func captureCallstack() []byte {
	// Capture all goroutine stacks
	var buf []byte
	bufferSize := 256 * 1024
	for {
		buf = make([]byte, bufferSize)
		if writtenBytes := runtime.Stack(buf, true); writtenBytes < bufferSize {
			buf = buf[:writtenBytes]
			break
		}
		bufferSize *= 2
	}
	return buf
}

// onPotentialDeadlock is the handler to be used by the deadlock library.
func (dlogger *deadlockLogger) onPotentialDeadlock() {
	// The deadlock reporting is done only once; this would prevent recursive deadlock issues.
	// in practive, once we report the deadlock, we panic and abort anyway, so it won't be an issue.
	dlogger.reportDeadlock.Do(func() {
		// Capture all goroutine stacks
		buf := captureCallstack()

		dlogger.bufferSync <- struct{}{}
		loggedString := dlogger.String()
		<-dlogger.bufferSync

		fmt.Fprintln(os.Stderr, string(buf))

		// logging the logged string to the logger has to happen in a separate go-routine, since the
		// logger itself ( for instance, the CyclicLogWriter ) is using a mutex of it's own.
		go func() {
			dlogger.Error(loggedString)
			dlogger.panic()
		}()
	})
}

func SetupDeadlockLogger(logger Logger) *deadlockLogger {
	dlogger := &deadlockLogger{
		Logger:     logger,
		Buffer:     bytes.NewBuffer(make([]byte, 0)),
		bufferSync: make(chan struct{}, 1),
	}

	dlogger.panic = dlogger.Panic
	deadlock.Opts.LogBuf = dlogger
	deadlock.Opts.OnPotentialDeadlock = dlogger.onPotentialDeadlock
	return dlogger
}
