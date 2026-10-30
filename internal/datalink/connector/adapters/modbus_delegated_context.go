package adapters

import (
	"context"
	"go-gateway/internal/datalink/schema"
	"time"
)

func (c *ModbusRTUConnector) AcquisitionReadTimeout() time.Duration {
	return modbusReadTimeout(c.config.Timeout)
}
func (c *ModbusRTUConnector) readProxy() *ModbusTCPConnector {
	return &ModbusTCPConnector{client: c.client, transport: c.transport, connected: c.isConnected(), persistentMode: true, config: schema.ConnectionConfigModbusTCP{DataFormat: c.config.DataFormat, Timeout: c.config.Timeout}}
}
func (c *ModbusRTUConnector) IsConnectedContext(ctx context.Context) (bool, error) {
	return c.readProxy().IsConnectedContext(ctx)
}
func (c *ModbusRTUConnector) CloseContext(ctx context.Context) error {
	if c.client == nil {
		return ctx.Err()
	}
	if err := c.client.CloseContext(ctx); err != nil {
		return err
	}
	c.setConnected(false)
	return nil
}
func (c *ModbusRTUConnector) afterOperationContext(ctx context.Context) {
	if !c.IsPersistentMode() {
		_ = c.CloseContext(ctx)
	}
}

func (c *ModbusUDPConnector) AcquisitionReadTimeout() time.Duration {
	return modbusReadTimeout(c.config.Timeout)
}
func (c *ModbusUDPConnector) readProxy() *ModbusTCPConnector {
	return &ModbusTCPConnector{client: c.client, transport: c.transport, connected: c.isConnected(), persistentMode: true, config: schema.ConnectionConfigModbusTCP{DataFormat: c.config.DataFormat, Timeout: c.config.Timeout}}
}
func (c *ModbusUDPConnector) IsConnectedContext(ctx context.Context) (bool, error) {
	return c.readProxy().IsConnectedContext(ctx)
}
func (c *ModbusUDPConnector) CloseContext(ctx context.Context) error {
	if c.client == nil {
		return ctx.Err()
	}
	if err := c.client.CloseContext(ctx); err != nil {
		return err
	}
	c.setConnected(false)
	return nil
}
func (c *ModbusUDPConnector) afterOperationContext(ctx context.Context) {
	if !c.IsPersistentMode() {
		_ = c.CloseContext(ctx)
	}
}
