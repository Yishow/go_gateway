package adapters

import (
	"context"
	"errors"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/protocol/wirebudget"
	"io"
	"net"
	"time"
)

// AcquisitionReadTimeout exposes the existing MC timeout without connection data.
func (c *MC3EConnector) AcquisitionReadTimeout() time.Duration {
	if c.config.Timeout <= 0 {
		return 5 * time.Second
	}
	return time.Duration(c.config.Timeout) * time.Second
}

func (c *MC3EConnector) readContext(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(ctx, c.AcquisitionReadTimeout())
	return wirebudget.With(ctx, mcRetryAttempts+1, mcRetryBaseDelay), cancel
}

const (
	mcRetryAttempts  = 3
	mcRetryBaseDelay = 300 * time.Millisecond
)

// withConnectionRetry 以乾淨重連重試 op，直到成功或耗盡重試次數。
// FX5U 系列的 MC 服務在 port 連線額度被佔用時會「接受 TCP 連線後直接關閉」
// （讀取得到 EOF）；這類瞬時佔用通常在短暫退避後釋放，重連重試即可恢復。
func (c *MC3EConnector) withConnectionRetry(ctx context.Context, op func() error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	err := op()
	for err != nil && mcConnectionFailure(err) {
		budget := wirebudget.From(ctx)
		if budget == nil || budget.Remaining() == 0 {
			return err
		}
		if waitErr := wirebudget.Wait(ctx); waitErr != nil {
			return waitErr
		}
		if reconnectErr := c.Reconnect(ctx); reconnectErr != nil {
			return reconnectErr
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		err = op()
	}
	return err
}

// retryExecuteRead 以乾淨重連重試 executeRead（首次失敗即進入重連重試）
func (c *MC3EConnector) retryExecuteRead(ctx context.Context, req connector.ReadRequest) (connector.ReadResult, error) {
	var result connector.ReadResult
	var err error
	if werr := c.withConnectionRetry(ctx, func() error {
		result, err = c.executeReadContext(ctx, req)
		return err
	}); werr != nil {
		return result, werr
	}
	return result, nil
}

// mcConnectionFailure never guesses retry eligibility from address/error text.
func mcConnectionFailure(err error) bool {
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, net.ErrClosed) {
		return true
	}
	var networkError *net.OpError
	return errors.As(err, &networkError)
}

// IsConnectedContext bounds connection-manager liveness checks.
func (c *MC3EConnector) IsConnectedContext(ctx context.Context) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	if !c.connected || c.client == nil {
		return false, nil
	}
	return c.client.IsConnectedContext(ctx)
}

// CloseContext preserves the same transport gate during canceled retirement.
func (c *MC3EConnector) CloseContext(ctx context.Context) error {
	if c.client == nil {
		return ctx.Err()
	}
	if err := c.client.CloseContext(ctx); err != nil {
		return err
	}
	c.connected = false
	return nil
}
