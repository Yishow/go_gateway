package modbus

import (
	"context"
	"fmt"
	"go-gateway/internal/protocol/wirebudget"
)

// WithContext creates a call-local view, without copying the client's mutex or
// its transport's transaction counter. Wire serialization stays at transport.
func (c *ModbusClient) WithContext(ctx context.Context) *ModbusClient {
	transport := c.transport
	if bound, ok := transport.(*contextTransport); ok {
		transport = bound.Transport
	}
	return NewClient(&contextTransport{Transport: transport, ctx: ctx}, c.unitID)
}
func (c *ModbusClient) ConnectContext(ctx context.Context) error { return c.WithContext(ctx).Connect() }
func (c *ModbusClient) CloseContext(ctx context.Context) error   { return c.WithContext(ctx).Close() }
func (c *ModbusClient) IsConnectedContext(ctx context.Context) (bool, error) {
	return c.WithContext(ctx).transport.(*contextTransport).connected()
}

type contextTransport struct {
	Transport
	ctx context.Context
}

func (t *contextTransport) Connect() error {
	if err := t.ctx.Err(); err != nil {
		return err
	}
	if transport, ok := t.Transport.(interface{ ConnectContext(context.Context) error }); ok {
		return transport.ConnectContext(t.ctx)
	}
	return t.Transport.Connect()
}
func (t *contextTransport) Close() error {
	if err := t.ctx.Err(); err != nil {
		return err
	}
	if transport, ok := t.Transport.(interface{ CloseContext(context.Context) error }); ok {
		return transport.CloseContext(t.ctx)
	}
	return t.Transport.Close()
}
func (t *contextTransport) connected() (bool, error) {
	if err := t.ctx.Err(); err != nil {
		return false, err
	}
	if transport, ok := t.Transport.(interface {
		IsConnectedContext(context.Context) (bool, error)
	}); ok {
		return transport.IsConnectedContext(t.ctx)
	}
	if transport, ok := t.Transport.(interface{ IsConnected() bool }); ok {
		return transport.IsConnected(), t.ctx.Err()
	}
	return true, nil
}
func (t *contextTransport) IsConnected() bool { live, _ := t.connected(); return live }
func (t *contextTransport) SendReceive(data []byte) ([]byte, error) {
	if err := t.ctx.Err(); err != nil {
		return nil, err
	}
	if transport, ok := t.Transport.(interface {
		SendReceiveContext(context.Context, []byte) ([]byte, error)
	}); ok {
		return transport.SendReceiveContext(t.ctx, data)
	}
	if err := wirebudget.Take(t.ctx); err != nil {
		return nil, err
	}
	return t.Transport.SendReceive(data)
}

// Unwrapping is essential: RTU uses CRC framing, never an MBAP header.
func (t *contextTransport) GetOriginalTransport() any {
	if unwrapper, ok := t.Transport.(transportUnwrapper); ok {
		return unwrapper.GetOriginalTransport()
	}
	return t.Transport
}
func (t *contextTransport) GetNextTransactionID() uint16 { id, _ := t.nextTransactionID(); return id }
func (t *contextTransport) nextTransactionID() (uint16, error) {
	if err := t.ctx.Err(); err != nil {
		return 0, err
	}
	if transport, ok := t.Transport.(interface {
		GetNextTransactionIDContext(context.Context) (uint16, error)
	}); ok {
		return transport.GetNextTransactionIDContext(t.ctx)
	}
	if transport, ok := t.Transport.(interface{ GetNextTransactionID() uint16 }); ok {
		return transport.GetNextTransactionID(), t.ctx.Err()
	}
	return 0, fmt.Errorf("不支援的傳輸類型")
}
