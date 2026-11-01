package adapters

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go-gateway/internal/protocol/wirebudget"
)

func (c *FatekConnector) AcquisitionReadTimeout() time.Duration {
	if c.config.Timeout <= 0 {
		return 5 * time.Second
	}
	return time.Duration(c.config.Timeout) * time.Second
}
func (c *FatekConnector) readContext(ctx context.Context, attempts int) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(ctx, c.AcquisitionReadTimeout())
	return wirebudget.With(ctx, attempts, 0), cancel
}
func (c *FatekConnector) ensureConnectionContext(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if c.client == nil {
		return fmt.Errorf("客戶端未初始化")
	}
	if c.persistentMode && !c.connected {
		return fmt.Errorf("未連線")
	}
	live, err := c.client.IsConnectedContext(ctx)
	if err != nil {
		return err
	}
	if !c.connected || !live {
		if err := c.acquisitionGeneration.connect(ctx, c.client.ConnectContext); err != nil {
			return err
		}
		c.connected = true
	}
	return nil
}
func (c *FatekConnector) IsConnectedContext(ctx context.Context) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	if !c.connected || c.client == nil {
		return false, nil
	}
	return c.client.IsConnectedContext(ctx)
}
func (c *FatekConnector) CloseContext(ctx context.Context) error {
	if c.client == nil {
		return ctx.Err()
	}
	if err := c.client.CloseContext(ctx); err != nil {
		return err
	}
	c.connected = false
	return nil
}

func (c *FatekConnector) afterOperationContext(ctx context.Context) {
	if !c.persistentMode && c.connected {
		if err := c.CloseContext(ctx); err != nil {
			slog.ErrorContext(ctx, adapterCleanupDiagnosticMessage, adapterCleanupProtocolLogKey, "fatek_fbs")
		}
	}
}
