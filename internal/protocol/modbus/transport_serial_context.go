package modbus

import (
	"context"
	"go.bug.st/serial"
	"io"
	"sync"
)

type contextSerialReader struct {
	ctx    context.Context
	reader io.Reader
}

func (r contextSerialReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	n, err := r.reader.Read(p)
	if n == 0 && err == nil {
		if r.ctx.Err() != nil {
			return 0, r.ctx.Err()
		}
		return 0, io.ErrNoProgress
	}
	return n, err
}

// Only the captured handle is closed. Join the cancellation callback before
// releasing the operation gate so an old callback cannot affect a replacement.
func serialContext(ctx context.Context, port serial.Port) func(bool) error {
	var once sync.Once
	var closeErr error
	closePort := func() { once.Do(func() { closeErr = port.Close() }) }
	done := make(chan struct{})
	stop := context.AfterFunc(ctx, func() { defer close(done); closePort() })
	return func(failed bool) error {
		if !stop() {
			<-done
		}
		if failed || ctx.Err() != nil {
			closePort()
		}
		return closeErr
	}
}
