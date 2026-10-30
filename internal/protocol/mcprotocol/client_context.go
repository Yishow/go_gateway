package mcprotocol

import (
	"context"
	"errors"
	"go-gateway/internal/protocol/wirebudget"
)

type contextTransport interface {
	SendReceiveContext(context.Context, []byte) ([]byte, error)
}
type contextConnector interface{ ConnectContext(context.Context) error }
type boundTransport struct {
	Transport
	ctx context.Context
}

func (t boundTransport) SendReceive(req []byte) ([]byte, error) {
	if err := t.ctx.Err(); err != nil {
		return nil, err
	}
	if transport, ok := t.Transport.(contextTransport); ok {
		return transport.SendReceiveContext(t.ctx, req)
	}
	// Legacy injected transports expose one synchronous call. Production TCP uses
	// the context interface above; this does not claim arbitrary legacy I/O cancels.
	if err := wirebudget.Take(t.ctx); err != nil {
		return nil, err
	}
	body, err := t.Transport.SendReceive(req)
	if t.ctx.Err() != nil {
		return nil, t.ctx.Err()
	}
	return body, err
}

// WithContext returns an isolated call view; shared client/frame are unchanged.
func (c *MCClient) WithContext(ctx context.Context) *MCClient {
	copy := *c
	transport := c.transport
	if bound, ok := transport.(boundTransport); ok {
		transport = bound.Transport
	}
	copy.transport = boundTransport{Transport: transport, ctx: ctx}
	return &copy
}

// ConnectContext forwards cancellation to transports that support it.
func (c *MCClient) ConnectContext(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	transport := c.transport
	if bound, ok := transport.(boundTransport); ok {
		transport = bound.Transport
	}
	if transport, ok := transport.(contextConnector); ok {
		return transport.ConnectContext(ctx)
	}
	err := transport.Connect()
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return err
}

// IsCommandRejected recognizes only the 3E command/subcommand rejection code.
// C059 also covers malformed command specification: callers must have validated
// their command before treating this as eligibility for a safe fallback plan.
// Source: Mitsubishi SH-080891, Error codes for MC protocol communication.
// https://www.mitsubishielectric.com/dl/fa/document/manual/plc/sh080891eng/sh080891engs.pdf
func IsCommandRejected(err error) bool { var e *MCError; return errors.As(err, &e) && e.Code == 0xC059 }

// IsConnectedContext avoids an uninterruptible status lock before a bounded read.
func (c *MCClient) IsConnectedContext(ctx context.Context) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	transport := c.transport
	if bound, ok := transport.(boundTransport); ok {
		transport = bound.Transport
	}
	if checker, ok := transport.(interface {
		IsConnectedContext(context.Context) (bool, error)
	}); ok {
		return checker.IsConnectedContext(ctx)
	}
	if checker, ok := transport.(interface{ IsConnected() bool }); ok {
		return checker.IsConnected(), ctx.Err()
	}
	return true, nil
}

// ReconnectContext lets context-aware transports replace their socket under the
// same cancelable gate, rather than first entering an unconditional Close lock.
func (c *MCClient) ReconnectContext(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	transport := c.transport
	if bound, ok := transport.(boundTransport); ok {
		transport = bound.Transport
	}
	if transport, ok := transport.(contextConnector); ok {
		return transport.ConnectContext(ctx)
	}
	if err := transport.Close(); err != nil {
		return err
	}
	return c.ConnectContext(ctx)
}

// CloseContext keeps protocol cleanup from waiting past its caller's deadline.
func (c *MCClient) CloseContext(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	transport := c.transport
	if bound, ok := transport.(boundTransport); ok {
		transport = bound.Transport
	}
	if closer, ok := transport.(interface{ CloseContext(context.Context) error }); ok {
		return closer.CloseContext(ctx)
	}
	return transport.Close()
}
