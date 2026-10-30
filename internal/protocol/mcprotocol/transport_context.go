package mcprotocol

import (
	"context"
	"fmt"
	"go-gateway/internal/protocol/wirebudget"
	"io"
	"net"
	"strconv"
	"time"
)

// SendReceiveContext bounds the complete call, including reconnects. Cancellation
// only touches the captured socket; its callback ends before the lock is released.
func (t *TCPTransport) SendReceiveContext(ctx context.Context, req []byte) ([]byte, error) {
	if err := t.mu.LockContext(ctx); err != nil {
		return nil, err
	}
	defer t.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	timeout := t.Timeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if budget := wirebudget.From(ctx); budget != nil && budget.Remaining() == 0 {
		return nil, wirebudget.ErrExhausted
	}
	if t.conn == nil {
		if err := t.internalConnectContext(ctx); err != nil {
			return nil, err
		}
	}
	response, err := t.sendReceiveContextLocked(ctx, req)
	if err == nil {
		return response, nil
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if b := wirebudget.From(ctx); b != nil {
		// Retry ownership remains at the adapter/scheduler, which share this budget.
		return nil, err
	}
	// Preserve the legacy direct-client reconnect behavior, but never extend its
	// deadline. Budgeted acquisition never performs this hidden resend.
	if t.conn == nil {
		if e := t.internalConnectContext(ctx); e == nil {
			return t.sendReceiveContextLocked(ctx, req)
		}
	}
	return nil, err
}

func (t *TCPTransport) sendReceiveContextLocked(ctx context.Context, req []byte) ([]byte, error) {
	conn := t.conn
	if conn == nil {
		return nil, fmt.Errorf("connection is closed")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	deadline, _ := ctx.Deadline()
	if err := conn.SetDeadline(deadline); err != nil {
		t.internalClose()
		return nil, err
	}
	callbackDone := make(chan struct{})
	stop := context.AfterFunc(ctx, func() { defer close(callbackDone); _ = conn.SetDeadline(time.Now()) })
	defer func() {
		if !stop() {
			<-callbackDone
		}
	}()
	if err := wirebudget.Take(ctx); err != nil {
		return nil, err
	}
	n, err := conn.Write(req)
	if err == nil && n != len(req) {
		err = io.ErrShortWrite
	}
	if err != nil {
		t.internalClose()
		return nil, err
	}
	var header [9]byte
	if _, err = io.ReadFull(conn, header[:]); err != nil {
		t.internalClose()
		return nil, err
	}
	size, err := ParseResponseHeader(header[:])
	if err != nil {
		t.internalClose()
		return nil, err
	}
	if size < 2 || size > t.MaxBodySize {
		t.internalClose()
		return nil, fmt.Errorf("invalid response body size")
	}
	body := make([]byte, size)
	if _, err = io.ReadFull(conn, body); err != nil {
		t.internalClose()
		return nil, err
	}
	if err = ctx.Err(); err != nil {
		t.internalClose()
		return nil, err
	}
	return body, nil
}

func (t *TCPTransport) internalConnectContext(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	t.internalClose()
	timeout := t.Timeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	dialCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	conn, err := (&net.Dialer{}).DialContext(dialCtx, "tcp", net.JoinHostPort(t.Host, strconv.Itoa(t.Port)))
	if err != nil {
		return err
	}
	if err = ctx.Err(); err != nil {
		_ = conn.Close()
		return err
	}
	t.conn = conn
	return nil
}
