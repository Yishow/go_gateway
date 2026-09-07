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
		c.config.Timeout = 5 // 秒；與 mcprotocol.NewTCPTransport 的預設逾時對齊（慢速 PLC 需餘裕）
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
	// io_no 為 0 時重映射為 3E 標準自站預設 0x03FF（JSON 未填 io_no 時預設 0）
	ioNo := c.config.IONo
	if ioNo == 0 {
		ioNo = 0x03FF
	}
	pcNo := c.config.PCNo
	if pcNo == 0 {
		// JSON 未填 pc_no 時預設 0，會被寫死覆蓋掉原本的 0xFF，需還原為 0xFF (自站)
		pcNo = 0xFF
	}
	frame := mcprotocol.NewRequestFrame(c.config.NetworkNo, pcNo, c.config.StationNo)
	frame.IONo = ioNo
	c.client.SetFrame(frame)
	// 逾時已於上方正規化為非零值（預設 5 秒），直接套用
	c.client.SetTimeout(time.Duration(c.config.Timeout) * time.Second)

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

// IsConnected 檢查連線狀態（同時確認標誌與底層 socket 存活）
func (c *MC3EConnector) IsConnected() bool {
	if !c.connected {
		return false
	}
	if c.client != nil && !c.client.IsConnected() {
		return false
	}
	return true
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

	result, err := c.executeRead(req)
	if err != nil && c.persistentMode && isConnectionError(err) {
		c.connected = false
		if reconnectErr := c.ensureConnection(); reconnectErr == nil {
			if retryResult, retryErr := c.executeRead(req); retryErr == nil {
				return retryResult, nil
			}
		}
	}

	if err != nil {
		result.Quality = schema.QualityBad
		result.Error = err.Error()
		return result, err
	}

	return result, nil
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

// ensureConnection 確保連線已建立
func (c *MC3EConnector) ensureConnection() error {
	if c.persistentMode {
		if !c.connected {
			if c.client == nil {
				return fmt.Errorf("未連線")
			}
			if err := c.client.Connect(); err != nil {
				return fmt.Errorf("連線失敗: %w", err)
			}
			c.connected = true
			return nil
		}
		// 若已標記連線但底層連線已斷開，自動自癒重連
		if c.client != nil && !c.client.IsConnected() {
			c.client.Close()
			if err := c.client.Connect(); err != nil {
				c.connected = false
				return fmt.Errorf("連線失敗: %w", err)
			}
			c.connected = true
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
