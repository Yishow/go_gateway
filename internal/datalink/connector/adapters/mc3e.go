package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/mcprotocol"
	"go-gateway/lib/hsllogic"
)

// =============================================================================
// Mitsubishi MC 3E 連接器
// =============================================================================

// MC3EConnector Mitsubishi MC Protocol 3E Frame 連接器適配器
type MC3EConnector struct {
	client         *mcprotocol.MCClient
	config         schema.ConnectionConfigMC3E
	connected      bool
	dataConverter  *hsllogic.DataConverter // hsllogic 數據轉換器
	persistentMode bool                    // 長連接模式標誌
}

// NewMC3EConnector 建立新的 MC 3E 連接器
func NewMC3EConnector() connector.Protocol {
	return &MC3EConnector{}
}

// Connect 建立連線
func (c *MC3EConnector) Connect(ctx context.Context, configJSON string) error {
	// 解析配置
	if err := json.Unmarshal([]byte(configJSON), &c.config); err != nil {
		return fmt.Errorf("解析 MC 3E 配置失敗: %w", err)
	}

	// 設定預設值
	if c.config.Port == 0 {
		c.config.Port = 5000
	}
	if c.config.Timeout == 0 {
		c.config.Timeout = 5
	}

	// 初始化 hsllogic 數據轉換器 (三菱使用 CDAB 格式：低位 word 在前)
	dataFormat := hsllogic.DataFormatCDAB
	if c.config.DataFormat != "" {
		dataFormat = hsllogic.DataFormat(c.config.DataFormat)
	}
	c.dataConverter = hsllogic.NewDataConverter(dataFormat)

	// 建立客戶端
	c.client = mcprotocol.NewClient(c.config.Host, c.config.Port)
	// 套用 station/pc/network/io_no 設定（原本被 hardcode 0 忽略的 bug）
	ioNo := c.config.IONo
	if ioNo == 0 {
		ioNo = 0x03FF
	}
	pcNo := c.config.PCNo
	if pcNo == 0 {
		// JSON 未填 pc_no 時預設 0，會被寫死覆蓋掉原本的 0xFF，需還原為 0xFF (自站)
		pcNo = 0xFF
	}
	c.client.SetFrame(c.config.NetworkNo, pcNo, c.config.StationNo, ioNo, 0)
	// extend timeout to 5s for slow PLC (raw is instant but framework needs margin)
	if c.config.Timeout > 0 {
		c.client.SetTimeout(time.Duration(c.config.Timeout) * time.Second)
	} else {
		c.client.SetTimeout(5 * time.Second)
	}

	// 預設使用長連接模式
	if !c.persistentMode {
		c.persistentMode = true
	}

	// 連線（僅在長連接模式下立即連線）
	if c.persistentMode {
		if err := c.client.Connect(); err != nil {
			return fmt.Errorf("MC 3E 連線失敗: %w", err)
		}
		c.connected = true
	}

	return nil
}

// Close 關閉連線
func (c *MC3EConnector) Close() error {
	if c.client != nil {
		c.connected = false
		return c.client.Close()
	}
	return nil
}

// IsConnected 檢查連線狀態
func (c *MC3EConnector) IsConnected() bool {
	return c.connected
}

// ProtocolType 取得協議類型
func (c *MC3EConnector) ProtocolType() schema.ProtocolType {
	return schema.ProtocolMC3E
}

// TestConnection 測試連線
func (c *MC3EConnector) TestConnection(ctx context.Context) error {
	if err := c.ensureConnection(); err != nil {
		return err
	}

	defer c.afterOperation()

	// 嘗試讀取 D0 來測試連線
	_, err := c.client.BatchReadWord("D", 0, 1)
	return err
}

// Read 讀取資料
func (c *MC3EConnector) Read(ctx context.Context, req connector.ReadRequest) (connector.ReadResult, error) {
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

	// 使用 hsllogic 解析地址 (格式: "D100", "M0", "X0", "Y0", "R100", "W100", etc.)
	parsedAddr, err := hsllogic.ParseAddress(hsllogic.ProtocolMitsubishi, req.Address)
	if err != nil {
		result.Quality = schema.QualityBad
		result.Error = err.Error()
		return result, err
	}

	count := 1
	if req.Count > 0 {
		count = req.Count
	}

	// 根據設備類型選擇讀取方式
	if parsedAddr.IsBitDevice {
		// 位元設備
		values, err := c.client.BatchReadBit(parsedAddr.DeviceType, parsedAddr.Offset, count)
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
		// 字組設備
		// 使用 hsllogic 計算需要讀取的暫存器數量 = count * 每個值需要的暫存器數
		regPerValue := hsllogic.RegisterCountForDataType(hsllogic.DataType(req.DataType))
		wordCount := count * regPerValue

		values, err := c.client.BatchReadWord(parsedAddr.DeviceType, parsedAddr.Offset, wordCount)
		if err != nil {
			result.Quality = schema.QualityBad
			result.Error = err.Error()
			return result, err
		}
		result.RawBytes = intSliceToBytes(values)

		// 使用 hsllogic 數據轉換器進行類型轉換
		registers := make([]uint16, len(values))
		for i, v := range values {
			registers[i] = uint16(v)
		}

		if count == 1 {
			result.Value = c.dataConverter.RegistersToValue(registers, hsllogic.DataType(req.DataType))
		} else {
			result.Value = c.dataConverter.RegistersToValues(registers, hsllogic.DataType(req.DataType), count)
		}
	}

	return result, nil
}

// Write 寫入資料
func (c *MC3EConnector) Write(ctx context.Context, req connector.WriteRequest) error {
	if err := c.ensureConnection(); err != nil {
		return err
	}

	defer c.afterOperation()

	// 使用 hsllogic 解析地址
	parsedAddr, err := hsllogic.ParseAddress(hsllogic.ProtocolMitsubishi, req.Address)
	if err != nil {
		return err
	}

	if parsedAddr.IsBitDevice {
		// 位元設備寫入
		switch v := req.Value.(type) {
		case bool:
			return c.client.BatchWriteBit(parsedAddr.DeviceType, parsedAddr.Offset, []bool{v})
		case []bool:
			return c.client.BatchWriteBit(parsedAddr.DeviceType, parsedAddr.Offset, v)
		default:
			return fmt.Errorf("位元設備需要布林值")
		}
	} else {
		// 字組設備寫入
		values, err := toIntSlice(req.Value)
		if err != nil {
			return err
		}
		return c.client.BatchWriteWord(parsedAddr.DeviceType, parsedAddr.Offset, values)
	}
}

// =============================================================================
// PersistentConnection 介面實作（長連接支援）
// =============================================================================

// SetPersistentConnection 設定是否使用長連接模式
func (c *MC3EConnector) SetPersistentConnection(enabled bool) {
	c.persistentMode = enabled
}

// IsPersistentMode 檢查當前是否為長連接模式
func (c *MC3EConnector) IsPersistentMode() bool {
	return c.persistentMode
}

// Disconnect 顯式斷線
func (c *MC3EConnector) Disconnect() error {
	return c.Close()
}

// Reconnect 重新連線
func (c *MC3EConnector) Reconnect(ctx context.Context) error {
	if c.client == nil {
		return fmt.Errorf("客戶端未初始化，請先呼叫 Connect")
	}

	c.client.Close()

	if err := c.client.Connect(); err != nil {
		c.connected = false
		return fmt.Errorf("重新連線失敗: %w", err)
	}

	c.connected = true
	return nil
}

// ensureConnection 確保連線已建立（用於短連接模式）
func (c *MC3EConnector) ensureConnection() error {
	if c.persistentMode {
		if !c.connected {
			return fmt.Errorf("未連線")
		}
		return nil
	}

	if c.client == nil {
		return fmt.Errorf("客戶端未初始化")
	}

	if !c.connected {
		if err := c.client.Connect(); err != nil {
			return fmt.Errorf("連線失敗: %w", err)
		}
		c.connected = true
	}

	return nil
}

// afterOperation 操作後處理（用於短連接模式）
func (c *MC3EConnector) afterOperation() {
	if !c.persistentMode && c.connected {
		c.client.Close()
		c.connected = false
	}
}

// 注意: intSliceToBytes 和 toIntSlice 已定義於 fatek.go

// =============================================================================
// 註冊連接器
// =============================================================================

func init() {
	connector.Register(schema.ProtocolMC3E, NewMC3EConnector)
}

