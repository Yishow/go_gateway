package adapters

import (
	"context"
	"fmt"

	"go-gateway/internal/datalink/connector"
)

// Write 寫入資料
func (c *ModbusTCPConnector) Write(ctx context.Context, req connector.WriteRequest) error {
	// 確保連線（長連接模式下斷線會自動重連，短連接模式下會建立連線）
	if err := c.ensureConnectionContext(ctx); err != nil {
		return err
	}

	// 短連接模式：操作完成後自動斷線
	defer c.afterOperation()

	err := c.executeWrite(ctx, req)
	if err != nil && c.IsPersistentMode() && isConnectionError(err) {
		// 連線中斷時嘗試單次自癒重連與重試
		c.setConnected(false)
		if reconnectErr := c.ensureConnectionContext(ctx); reconnectErr == nil {
			if retryErr := c.executeWrite(ctx, req); retryErr == nil {
				return nil
			}
		}
	}

	return err
}

func (c *ModbusTCPConnector) executeWrite(ctx context.Context, req connector.WriteRequest) error {
	// 解析地址
	address, function, err := parseModbusAddress(req.Address, req.Function)
	if err != nil {
		return err
	}

	switch function {
	case modbusFunctionCoil, "05", "FC05", "01", "FC01":
		// 寫入單個線圈
		value, ok := req.Value.(bool)
		if !ok {
			return fmt.Errorf("線圈寫入需要布林值")
		}
		return c.client.WithContext(ctx).WriteSingleCoil(address, value)

	case modbusFunctionHolding, "06", "FC06", "", "03", "FC03":
		// 寫入單個保持暫存器
		value, err := toUint16(req.Value)
		if err != nil {
			return err
		}
		return c.client.WithContext(ctx).WriteSingleRegister(address, value)

	case "15", "FC15":
		// 寫入多個線圈
		values, ok := req.Value.([]bool)
		if !ok {
			return fmt.Errorf("多線圈寫入需要布林切片")
		}
		return c.client.WithContext(ctx).WriteMultipleCoils(address, values)

	case "16", "FC16":
		// 寫入多個暫存器
		values, err := toUint16Slice(req.Value)
		if err != nil {
			return err
		}
		return c.client.WithContext(ctx).WriteMultipleRegisters(address, values)

	default:
		return fmt.Errorf("不支援的寫入功能碼: %s", function)
	}
}
