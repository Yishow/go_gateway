package modbus

import (
	"context"
	"errors"
	"go-gateway/internal/protocol/wirebudget"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var errCancelDeadlineFault = errors.New("cancel deadline fault")
var errCancelCloseFault = errors.New("cancel close fault")

type cancellationFaultConn struct {
	net.Conn
	entered      chan struct{}
	closed       chan struct{}
	once         sync.Once
	deadlines    atomic.Int32
	writes       atomic.Int32
	closes       atomic.Int32
	closeFailure bool
}

func (c *cancellationFaultConn) SetDeadline(time.Time) error {
	if c.deadlines.Add(1) > 1 {
		return errCancelDeadlineFault
	}
	return nil
}
func (c *cancellationFaultConn) Write(b []byte) (int, error) {
	c.writes.Add(1)
	close(c.entered)
	return len(b), nil
}
func (c *cancellationFaultConn) Read([]byte) (int, error) {
	select {
	case <-c.closed:
	case <-time.After(time.Second):
	}
	return 0, net.ErrClosed
}
func (c *cancellationFaultConn) Close() error {
	c.closes.Add(1)
	c.once.Do(func() { close(c.closed) })
	if c.closeFailure {
		return errCancelCloseFault
	}
	return nil
}
func TestModbusCancelDeadlineFailureRetainsErrorsAndSocketOwnership(t *testing.T) {
	for _, closeFailure := range []bool{false, true} {
		t.Run(map[bool]string{false: "close-success", true: "close-failure"}[closeFailure], func(t *testing.T) {
			old := &cancellationFaultConn{entered: make(chan struct{}), closed: make(chan struct{}), closeFailure: closeFailure}
			transport := NewTCPTransport("127.0.0.1", 1)
			transport.conn = old

			base, cancel := context.WithCancel(t.Context())
			defer cancel()
			ctx := wirebudget.With(base, 2, 0)
			done := make(chan error, 1)
			go func() { _, err := transport.SendReceiveContext(ctx, []byte{1}); done <- err }()
			<-old.entered
			cancel()
			var err error
			select {
			case err = <-done:
			case <-time.After(2 * time.Second):
				t.Fatal("canceled operation did not join")
			}
			require.ErrorIs(t, err, context.Canceled)
			require.ErrorIs(t, err, errCancelDeadlineFault)
			require.ErrorIs(t, err, net.ErrClosed)
			if closeFailure {
				require.ErrorIs(t, err, errCancelCloseFault)
			}
			require.EqualValues(t, 1, old.writes.Load())
			require.Equal(t, 1, wirebudget.From(ctx).Used())
			replacement := &cancellationFaultConn{entered: make(chan struct{}), closed: make(chan struct{})}
			transport.conn = replacement
			oldCloses := old.closes.Load()
			cancel()
			require.Equal(t, oldCloses, old.closes.Load())
			require.Zero(t, replacement.closes.Load())
			require.Zero(t, replacement.deadlines.Load())
			require.Zero(t, replacement.writes.Load())
		})
	}
}
