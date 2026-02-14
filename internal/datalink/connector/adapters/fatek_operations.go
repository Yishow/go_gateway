package adapters

import (
	"context"
	"fmt"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/fatek"
)

// TestConnection 測試連線
func (c *FatekConnector) TestConnection(ctx context.Context) error {
	if err := c.ensureConnection(); err != nil {
		return err
	}

	defer c.afterOperation()

	// 改用讀取一個暫存器來測試連線（例如讀取 D0）
	// 原因：某些 FATEK PLC 型號不支援 Loopback Test (Cmd 4E)
	// 讀取指令 (Cmd 46) 是所有 FATEK PLC 都支援的標準指令
	// 注意：如果 D0 不存在，可以改用 R0 或其他暫存器
	_, err := c.client.ReadRegisters("D", 0, 1)
	if err != nil {
		// 如果讀取 D0 失敗，嘗試讀取 R0（系統暫存器，通常總是存在）
		_, err2 := c.client.ReadRegisters("R", 0, 1)
		if err2 != nil {
			return fmt.Errorf("連線測試失敗 (嘗試讀取 D0 和 R0 都失敗): D0錯誤=%w, R0錯誤=%w", err, err2)
		}
		// R0 讀取成功，連線正常
		return nil
	}

	return nil
}

// Read 讀取資料
func (c *FatekConnector) Read(ctx context.Context, req connector.ReadRequest) (connector.ReadResult, error) {
	if err := c.ensureConnection(); err != nil {
		return connector.ReadResult{
			Quality:   schema.QualityBad,
			Timestamp: time.Now(),
			Error:     err.Error(),
		}, err
	}

	defer c.afterOperation()

	result := connector.ReadResult{
		Timestamp: time.Now(),
		Quality:   schema.QualityGood,
	}

	// 解析地址 (格式: "D0100", "R0", "M100", "X0", "Y0", etc.)
	symbol, address, err := parseFatekAddress(req.Address)
	if err != nil {
		result.Quality = schema.QualityBad
		result.Error = err.Error()
		return result, err
	}

	// 取得元件類型
	comp, err := fatek.GetComponentType(symbol)
	if err != nil {
		result.Quality = schema.QualityBad
		result.Error = err.Error()
		return result, err
	}

	count := 1
	if req.Count > 0 {
		count = req.Count
	}

	if comp.IsDiscrete {
		// 離散元件 (位元)
		values, err := c.client.ReadStatus(symbol, address, count)
		if err != nil {
			result.Quality = schema.QualityBad
			result.Error = err.Error()
			return result, err
		}
		if len(values) > 0 {
			if count == 1 {
				result.Value = values[0]
			} else {
				result.Value = values
			}
		}
	} else {
		// 暫存器元件
		values, err := c.client.ReadRegisters(symbol, address, count)
		if err != nil {
			result.Quality = schema.QualityBad
			result.Error = err.Error()
			return result, err
		}
		result.RawBytes = intSliceToBytes(values)
		result.Value = convertFatekValue(values, req.DataType, comp.Width)
	}

	return result, nil
}

// Write 寫入資料
func (c *FatekConnector) Write(ctx context.Context, req connector.WriteRequest) error {
	if err := c.ensureConnection(); err != nil {
		return err
	}

	defer c.afterOperation()

	// 解析地址
	symbol, address, err := parseFatekAddress(req.Address)
	if err != nil {
		return err
	}

	// 取得元件類型
	comp, err := fatek.GetComponentType(symbol)
	if err != nil {
		return err
	}

	if comp.IsDiscrete {
		// 離散元件寫入
		switch v := req.Value.(type) {
		case bool:
			return c.client.WriteStatus(symbol, address, []bool{v})
		case []bool:
			return c.client.WriteStatus(symbol, address, v)
		default:
			return fmt.Errorf("離散元件需要布林值")
		}
	} else {
		// 暫存器元件寫入
		values, err := toIntSlice(req.Value)
		if err != nil {
			return err
		}
		return c.client.WriteRegisters(symbol, address, values)
	}
}
