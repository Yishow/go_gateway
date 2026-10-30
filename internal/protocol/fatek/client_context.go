package fatek

import (
	"context"
	"go-gateway/internal/protocol/wirebudget"
)

type contextTransport struct {
	Transport
	ctx context.Context
}

func (t contextTransport) SendReceive(data []byte) ([]byte, error) {
	if sender, ok := t.Transport.(interface {
		SendReceiveContext(context.Context, []byte) ([]byte, error)
	}); ok {
		return sender.SendReceiveContext(t.ctx, data)
	}
	if err := wirebudget.Take(t.ctx); err != nil {
		return nil, err
	}
	return t.Transport.SendReceive(data)
}

// WithContext builds a call-local view without copying a client mutex.
func (c *FatekClient) WithContext(ctx context.Context) *FatekClient {
	transport := c.transport
	for {
		bound, ok := transport.(contextTransport)
		if !ok {
			break
		}
		transport = bound.Transport
	}
	return NewClient(contextTransport{transport, ctx}, c.station)
}
func (c *FatekClient) ConnectContext(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if t, ok := c.transport.(interface{ ConnectContext(context.Context) error }); ok {
		return t.ConnectContext(ctx)
	}
	return c.Connect()
}
func (c *FatekClient) CloseContext(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if t, ok := c.transport.(interface{ CloseContext(context.Context) error }); ok {
		return t.CloseContext(ctx)
	}
	return c.Close()
}
func (c *FatekClient) IsConnectedContext(ctx context.Context) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	if t, ok := c.transport.(interface {
		IsConnectedContext(context.Context) (bool, error)
	}); ok {
		return t.IsConnectedContext(ctx)
	}
	return true, nil
}
