package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/modbus"
)

// =============================================================================
// Modbus UDP 連接器
// =============================================================================

// ModbusUDPConnector Modbus UDP 協議連接器適配器
type ModbusUDPConnector struct {
	client         *modbus.ModbusClient
	transport      modbus.Transport
	config         schema.ConnectionConfigModbusUDP
	connected      bool
	persistentMode bool // 長連接模式標誌
}

// NewModbusUDPConnector 建立新的 Modbus UDP 連接器
func NewModbusUDPConnector() connector.Protocol {
	return &ModbusUDPConnector{}
}

// Connect 建立連線
func (c *ModbusUDPConnector) Connect(ctx context.Context, configJSON string) error {
	if err := json.Unmarshal([]byte(configJSON), &c.config); err != nil {
		return fmt.Errorf("解析 Modbus UDP 配置失敗: %w", err)
	}

	if c.config.Port == 0 {
		c.config.Port = 502
	}
	if c.config.SlaveID == 0 {
		c.config.SlaveID = 1
	}
	if c.config.Timeout == 0 {
		c.config.Timeout = 5
	}

	c.transport = modbus.NewUDPTransport(
		c.config.Host,
		c.config.Port,
	)
	// 設定 Timeout
	if udp, ok := c.transport.(*modbus.UDPTransport); ok {
		udp.Timeout = time.Duration(c.config.Timeout) * time.Second
	}

	c.client = modbus.NewClient(c.transport, c.config.SlaveID)

	// 預設使用長連接模式
	if !c.persistentMode {
		c.persistentMode = true
	}

	// 連線（僅在長連接模式下立即連線）
	if c.persistentMode {
		if err := c.client.Connect(); err != nil {
			return fmt.Errorf("Modbus UDP 連線失敗: %w", err)
		}
		c.connected = true
	}

	return nil
}

// Close 關閉連線
func (c *ModbusUDPConnector) Close() error {
	if c.client != nil {
		c.connected = false
		return c.client.Close()
	}
	return nil
}

// IsConnected 檢查連線狀態
func (c *ModbusUDPConnector) IsConnected() bool {
	return c.connected
}

// ProtocolType 取得協議類型
func (c *ModbusUDPConnector) ProtocolType() schema.ProtocolType {
	return schema.ProtocolModbusUDP
}

// TestConnection 測試連線
func (c *ModbusUDPConnector) TestConnection(ctx context.Context) error {
	if !c.connected {
		return fmt.Errorf("未連線")
	}
	_, err := c.client.ReadHoldingRegisters(0, 1)
	return err
}

// Read 讀取資料
func (c *ModbusUDPConnector) Read(ctx context.Context, req connector.ReadRequest) (connector.ReadResult, error) {
	// 確保連線
	if err := c.ensureConnection(); err != nil {
		return connector.ReadResult{Quality: schema.QualityBad, Timestamp: time.Now(), Error: err.Error()}, err
	}

	// 短連接模式：操作完成後自動斷線
	defer c.afterOperation()

	tcpConn := &ModbusTCPConnector{client: c.client, connected: c.connected}
	return tcpConn.Read(ctx, req)
}

// Write 寫入資料
func (c *ModbusUDPConnector) Write(ctx context.Context, req connector.WriteRequest) error {
	// 確保連線
	if err := c.ensureConnection(); err != nil {
		return err
	}

	// 短連接模式：操作完成後自動斷線
	defer c.afterOperation()

	tcpConn := &ModbusTCPConnector{client: c.client, connected: c.connected}
	return tcpConn.Write(ctx, req)
}

// =============================================================================
// ModbusUDP PersistentConnection 介面實作
// =============================================================================

// SetPersistentConnection 設定是否使用長連接模式
func (c *ModbusUDPConnector) SetPersistentConnection(enabled bool) {
	c.persistentMode = enabled
}

// IsPersistentMode 檢查當前是否為長連接模式
func (c *ModbusUDPConnector) IsPersistentMode() bool {
	return c.persistentMode
}

// Disconnect 顯式斷線
func (c *ModbusUDPConnector) Disconnect() error {
	return c.Close()
}

// Reconnect 重新連線
func (c *ModbusUDPConnector) Reconnect(ctx context.Context) error {
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
func (c *ModbusUDPConnector) ensureConnection() error {
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
func (c *ModbusUDPConnector) afterOperation() {
	if !c.persistentMode && c.connected {
		c.client.Close()
		c.connected = false
	}
}
