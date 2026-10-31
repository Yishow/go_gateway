package adapters

import (
	"context"
	"errors"
	"fmt"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/modbus"
	"go-gateway/internal/protocol/wirebudget"
	"io"
	"net"
	"time"
)

func modbusReadTimeout(seconds int) time.Duration {
	if seconds <= 0 {
		return 5 * time.Second
	}
	return time.Duration(seconds) * time.Second
}
func (c *ModbusTCPConnector) AcquisitionReadTimeout() time.Duration {
	return modbusReadTimeout(c.config.Timeout)
}
func (c *ModbusTCPConnector) readContext(ctx context.Context, attempts int) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(ctx, c.AcquisitionReadTimeout())
	return wirebudget.With(ctx, attempts, 0), cancel
}
func (c *ModbusTCPConnector) IsConnectedContext(ctx context.Context) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	if !c.isConnected() || c.client == nil {
		return false, nil
	}
	return c.client.IsConnectedContext(ctx)
}
func (c *ModbusTCPConnector) CloseContext(ctx context.Context) error {
	if c.client == nil {
		return ctx.Err()
	}
	if err := c.client.CloseContext(ctx); err != nil {
		return err
	}
	c.setConnected(false)
	return nil
}
func (c *ModbusTCPConnector) ensureConnectionContext(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if c.client == nil {
		return fmt.Errorf("未連線")
	}
	connected, err := c.IsConnectedContext(ctx)
	if err != nil {
		return err
	}
	if connected {
		return nil
	}
	if err = c.client.ConnectContext(ctx); err != nil {
		return err
	}
	c.setConnected(true)
	return nil
}
func (c *ModbusTCPConnector) afterOperationContext(ctx context.Context) {
	if !c.IsPersistentMode() {
		_ = c.CloseContext(ctx)
	}
}

// Read shares the caller's actual-send budget with every reconnect and retry.
func (c *ModbusTCPConnector) Read(ctx context.Context, req connector.ReadRequest) (connector.ReadResult, error) {
	ctx, cancel := c.readContext(ctx, 2)
	defer cancel()
	if err := c.ensureConnectionContext(ctx); err != nil {
		return connector.ReadResult{Quality: schema.QualityBad, Timestamp: time.Now(), Error: err.Error()}, err
	}
	defer c.afterOperationContext(ctx)
	var result connector.ReadResult
	err := c.withReadConnectionRetry(ctx, func() error {
		var readErr error
		result, readErr = c.executeReadContext(ctx, req)
		return readErr
	})
	if err != nil {
		result.Quality = schema.QualityBad
		result.Error = err.Error()
	}
	return result, err
}

// withReadConnectionRetry is shared by scalar and planned reads. Only the same
// typed connection failures retry, and every send consumes the caller's budget.
func (c *ModbusTCPConnector) withReadConnectionRetry(ctx context.Context, read func() error) error {
	err := read()
	for err != nil && c.IsPersistentMode() && modbusConnectionFailure(err) {
		budget := wirebudget.From(ctx)
		if budget.Remaining() == 0 {
			break
		}
		if waitErr := wirebudget.Wait(ctx); waitErr != nil {
			err = waitErr
			break
		}
		if reconnectErr := c.Reconnect(ctx); reconnectErr != nil {
			err = reconnectErr
			break
		}
		err = read()
	}
	return err
}

func modbusConnectionFailure(err error) bool {
	if errors.Is(err, modbus.ErrConnectionClosed) || errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, net.ErrClosed) {
		return true
	}
	var op *net.OpError
	return errors.As(err, &op)
}
