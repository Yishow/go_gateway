package modbus

import (
	"context"
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
func socketContext(ctx context.Context, conn net.Conn) (func(), error) {
	deadline, _ := ctx.Deadline()
	if err := conn.SetDeadline(deadline); err != nil {
		return nil, err
	}
	done := make(chan struct{})
	stop := context.AfterFunc(ctx, func() { defer close(done); _ = conn.SetDeadline(time.Now()) })
	return func() {
		if !stop() {
			<-done
		}
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
