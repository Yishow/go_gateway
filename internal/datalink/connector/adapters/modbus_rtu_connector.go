package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/modbus"
)

// =============================================================================
// Modbus RTU 連接器
// =============================================================================

// ModbusRTUConnector Modbus RTU 協議連接器適配器
type ModbusRTUConnector struct {
	client         *modbus.ModbusClient
	transport      modbus.Transport
	config         schema.ConnectionConfigModbusRTU
	stateMu        sync.RWMutex
	connected      bool
	persistentMode bool // 長連接模式標誌
}

// NewModbusRTUConnector 建立新的 Modbus RTU 連接器
func NewModbusRTUConnector() connector.Protocol {
	return &ModbusRTUConnector{}
}

// Connect 建立連線
func (c *ModbusRTUConnector) Connect(ctx context.Context, configJSON string) error {
	if err := json.Unmarshal([]byte(configJSON), &c.config); err != nil {
		return fmt.Errorf("解析 Modbus RTU 配置失敗: %w", err)
	}

	// 設定預設值
	if c.config.BaudRate == 0 {
		c.config.BaudRate = 9600
	}
	if c.config.DataBits == 0 {
		c.config.DataBits = 8
	}
	if c.config.StopBits == 0 {
		c.config.StopBits = 1
	}
	if c.config.Parity == "" {
		c.config.Parity = "none"
	}
	if c.config.SlaveID == 0 {
		c.config.SlaveID = 1
	}
	if c.config.Timeout == 0 {
		c.config.Timeout = 5
	}

	// 建立 RTU 傳輸層
	c.transport = modbus.NewRTUTransport(
		c.config.SerialPort,
		c.config.BaudRate,
		c.config.DataBits,
		c.config.StopBits,
		c.config.Parity,
		time.Duration(c.config.Timeout)*time.Second,
	)

	c.client = modbus.NewClient(c.transport, c.config.SlaveID)

	// 預設使用長連接模式
	if !c.IsPersistentMode() {
		c.SetPersistentConnection(true)
	}

	// 連線（僅在長連接模式下立即連線）
	if c.IsPersistentMode() {
		if err := c.client.Connect(); err != nil {
			return fmt.Errorf("Modbus RTU 連線失敗: %w", err)
		}
		c.setConnected(true)
	}

	return nil
}

// Close 關閉連線
func (c *ModbusRTUConnector) Close() error {
	if c.client != nil {
		c.setConnected(false)
		return c.client.Close()
	}
	return nil
}

// IsConnected 檢查連線狀態
func (c *ModbusRTUConnector) IsConnected() bool {
	return c.isConnected()
}

// ProtocolType 取得協議類型
func (c *ModbusRTUConnector) ProtocolType() schema.ProtocolType {
	return schema.ProtocolModbusRTU
}

// TestConnection 測試連線
func (c *ModbusRTUConnector) TestConnection(ctx context.Context) error {
	if !c.isConnected() {
		return fmt.Errorf("未連線")
	}
	_, err := c.client.ReadHoldingRegisters(0, 1)
	return err
}

// Read 讀取資料 (委託給 TCP 版本的邏輯，因為 Modbus Client 已經抽象化)
func (c *ModbusRTUConnector) Read(ctx context.Context, req connector.ReadRequest) (connector.ReadResult, error) {
	// 確保連線
	if err := c.ensureConnection(); err != nil {
		return connector.ReadResult{Quality: schema.QualityBad, Timestamp: time.Now(), Error: err.Error()}, err
	}

	// 短連接模式：操作完成後自動斷線
	defer c.afterOperation()

	// 複用 TCP 版本的讀取邏輯
	tcpConn := &ModbusTCPConnector{client: c.client, connected: true, persistentMode: true}
	return tcpConn.Read(ctx, req)
}

// Write 寫入資料
func (c *ModbusRTUConnector) Write(ctx context.Context, req connector.WriteRequest) error {
	// 確保連線
	if err := c.ensureConnection(); err != nil {
		return err
	}

	// 短連接模式：操作完成後自動斷線
	defer c.afterOperation()

	tcpConn := &ModbusTCPConnector{client: c.client, connected: true, persistentMode: true}
	return tcpConn.Write(ctx, req)
}

// =============================================================================
// ModbusRTU PersistentConnection 介面實作
// =============================================================================

// SetPersistentConnection 設定是否使用長連接模式
func (c *ModbusRTUConnector) SetPersistentConnection(enabled bool) {
	c.stateMu.Lock()
	c.persistentMode = enabled
	c.stateMu.Unlock()
}

// IsPersistentMode 檢查當前是否為長連接模式
func (c *ModbusRTUConnector) IsPersistentMode() bool {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.persistentMode
}

// Disconnect 顯式斷線
func (c *ModbusRTUConnector) Disconnect() error {
	return c.Close()
}

// Reconnect 重新連線
func (c *ModbusRTUConnector) Reconnect(ctx context.Context) error {
	if c.client == nil {
		return fmt.Errorf("客戶端未初始化，請先呼叫 Connect")
	}

	c.client.Close()

	if err := c.client.Connect(); err != nil {
		c.setConnected(false)
		return fmt.Errorf("重新連線失敗: %w", err)
	}

	c.setConnected(true)
	return nil
}

// ensureConnection 確保連線已建立（用於短連接模式）
func (c *ModbusRTUConnector) ensureConnection() error {
	if c.IsPersistentMode() {
		if !c.isConnected() {
			return fmt.Errorf("未連線")
		}
		return nil
	}

	if c.client == nil {
		return fmt.Errorf("客戶端未初始化")
	}

	if !c.isConnected() {
		if err := c.client.Connect(); err != nil {
			return fmt.Errorf("連線失敗: %w", err)
		}
		c.setConnected(true)
	}

	return nil
}

// afterOperation 操作後處理（用於短連接模式）
func (c *ModbusRTUConnector) afterOperation() {
	if !c.IsPersistentMode() && c.isConnected() {
		c.client.Close()
		c.setConnected(false)
	}
}

func (c *ModbusRTUConnector) setConnected(connected bool) {
	c.stateMu.Lock()
	c.connected = connected
	c.stateMu.Unlock()
}

func (c *ModbusRTUConnector) isConnected() bool {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.connected
}
