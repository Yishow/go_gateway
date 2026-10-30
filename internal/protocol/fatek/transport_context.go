package fatek

import (
	"bufio"
	"context"
	"go-gateway/internal/protocol/wirebudget"
	"io"
	"net"
	"strconv"
	"time"
)

func boundedContext(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	return context.WithTimeout(ctx, timeout)
}

// ConnectContext replaces the captured connection under the same cancelable gate.
func (t *TCPTransport) ConnectContext(ctx context.Context) error {
	if err := t.mu.LockContext(ctx); err != nil {
		return err
	}
	defer t.mu.Unlock()
	ctx, cancel := boundedContext(ctx, t.Timeout)
	defer cancel()
	if err := t.internalClose(); err != nil {
		return err
	}
	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", net.JoinHostPort(t.Host, strconv.Itoa(t.Port)))
	if err != nil {
		return err
	}
	if err = ctx.Err(); err != nil {
		_ = conn.Close()
		return err
	}
	t.conn = conn
	t.reader = bufio.NewReader(conn)
	return nil
}
func (t *TCPTransport) CloseContext(ctx context.Context) error {
	if err := t.mu.LockContext(ctx); err != nil {
		return err
	}
	defer t.mu.Unlock()
	return t.internalClose()
}
func (t *TCPTransport) IsConnectedContext(ctx context.Context) (bool, error) {
	if err := t.mu.LockContext(ctx); err != nil {
		return false, err
	}
	defer t.mu.Unlock()
	return t.conn != nil, nil
}
func (t *TCPTransport) SendReceiveContext(ctx context.Context, data []byte) (response []byte, err error) {
	if err = t.mu.LockContext(ctx); err != nil {
		return nil, err
	}
	defer t.mu.Unlock()
	ctx, cancel := boundedContext(ctx, t.Timeout)
	defer cancel()
	if t.conn == nil {
		return nil, ErrConnectionClosed
	}
	conn := t.conn
	deadline, _ := ctx.Deadline()
	if err = conn.SetDeadline(deadline); err != nil {
		return nil, err
	}
	done := make(chan struct{})
	stop := context.AfterFunc(ctx, func() { defer close(done); _ = conn.SetDeadline(time.Now()) })
	defer func() {
		if !stop() {
			<-done
		}
		if ctx.Err() != nil {
			err = ctx.Err()
		}
		if err != nil {
			_ = t.internalClose()
		}
	}()
	if err = wirebudget.Take(ctx); err != nil {
		return nil, err
	}
	n, err := conn.Write(data)
	if err != nil {
		return nil, err
	}
	if n != len(data) {
		return nil, io.ErrShortWrite
	}
	return t.reader.ReadBytes(ETX)
}
