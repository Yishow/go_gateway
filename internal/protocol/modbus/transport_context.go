package modbus

import (
	"context"
	"errors"
	"net"
	"time"
)

func transportContext(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	return context.WithTimeout(ctx, timeout)
}

// The callback captures this socket, and is joined before its operation gate opens.
func socketContext(ctx context.Context, conn net.Conn) (func() error, error) {
	deadline, _ := ctx.Deadline()
	if err := conn.SetDeadline(deadline); err != nil {
		return nil, err
	}
	done := make(chan struct{})
	var cancelErr error
	stop := context.AfterFunc(ctx, func() {
		defer close(done)
		if deadlineErr := conn.SetDeadline(time.Now()); deadlineErr != nil {
			cancelErr = errors.Join(ctx.Err(), deadlineErr, conn.Close())
		}
	})
	return func() error {
		if !stop() {
			<-done
		}
		return cancelErr
	}, nil
}

// Socket deadlines may fire just before the context timer callback runs.
func transportContextError(ctx context.Context, err error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if deadline, ok := ctx.Deadline(); ok && !time.Now().Before(deadline) {
		return context.DeadlineExceeded
	}
	return err
}
