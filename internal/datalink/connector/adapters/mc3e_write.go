package adapters

import (
	"context"
	"fmt"

	"go-gateway/internal/datalink/connector"
	"go-gateway/lib/hsllogic"
)

// Write 寫入資料
func (c *MC3EConnector) Write(ctx context.Context, req connector.WriteRequest) error {
	if err := c.ensureConnectionContext(ctx); err != nil {
		return err
	}

	defer c.afterOperation()

	err := c.executeWrite(ctx, req)
	if err != nil && c.IsPersistentMode() && isConnectionError(err) {
		c.connected = false
		if reconnectErr := c.ensureConnectionContext(ctx); reconnectErr == nil {
			if retryErr := c.executeWrite(ctx, req); retryErr == nil {
				return nil
			}
		}
	}

	return err
}

func (c *MC3EConnector) executeWrite(ctx context.Context, req connector.WriteRequest) error {
	// 使用 hsllogic 解析地址
	parsedAddr, err := hsllogic.ParseAddress(hsllogic.ProtocolMitsubishi, req.Address)
	if err != nil {
		return err
	}

	if parsedAddr.IsBitDevice {
		// 位元設備寫入
		switch v := req.Value.(type) {
		case bool:
			return c.client.WithContext(ctx).BatchWriteBit(parsedAddr.DeviceType, parsedAddr.Offset, []bool{v})
		case []bool:
			return c.client.WithContext(ctx).BatchWriteBit(parsedAddr.DeviceType, parsedAddr.Offset, v)
		default:
			return fmt.Errorf("位元設備需要布林值")
		}
	}

	// 字組設備寫入
	values, err := toIntSlice(req.Value)
	if err != nil {
		return err
	}
	return c.client.WithContext(ctx).BatchWriteWord(parsedAddr.DeviceType, parsedAddr.Offset, values)
}
