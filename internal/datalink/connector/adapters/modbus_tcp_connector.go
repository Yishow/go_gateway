package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/modbus"
)

// =============================================================================
// Modbus TCP 連接器
// =============================================================================

// ModbusTCPConnector Modbus TCP 協議連接器適配器
type ModbusTCPConnector struct {
	client         *modbus.ModbusClient
	transport      modbus.Transport
	config         schema.ConnectionConfigModbusTCP
	stateMu        sync.RWMutex
	connected      bool
	persistentMode bool // 長連接模式標誌
}

// NewModbusTCPConnector 建立新的 Modbus TCP 連接器
func NewModbusTCPConnector() connector.Protocol {
	return &ModbusTCPConnector{}
}

// Connect 建立連線
func (c *ModbusTCPConnector) Connect(ctx context.Context, configJSON string) error {
	// 解析配置
	if err := json.Unmarshal([]byte(configJSON), &c.config); err != nil {
		return fmt.Errorf("解析 Modbus TCP 配置失敗: %w", err)
	}

	// 設定預設值
	if c.config.Port == 0 {
		c.config.Port = 502
	}
	if c.config.SlaveID == 0 {
		c.config.SlaveID = 1
	}
	if c.config.Timeout == 0 {
		c.config.Timeout = 5
	}

	// 建立傳輸層
	c.transport = modbus.NewTCPTransport(
		c.config.Host,
		c.config.Port,
	)
	// 設定 Timeout (如果 Transport 支援)
	if tcp, ok := c.transport.(*modbus.TCPTransport); ok {
		tcp.Timeout = time.Duration(c.config.Timeout) * time.Second
	}

	// 建立客戶端
	c.client = modbus.NewClient(c.transport, c.config.SlaveID)

	// 預設使用長連接模式
	if !c.IsPersistentMode() {
		c.SetPersistentConnection(true)
	}

	// 連線（僅在長連接模式下立即連線）
	if c.IsPersistentMode() {
		if err := c.client.Connect(); err != nil {
			return fmt.Errorf("modbus TCP 連線失敗: %w", err)
		}
		c.setConnected(true)
	}

	return nil
}

// Close 關閉連線
func (c *ModbusTCPConnector) Close() error {
	if c.client != nil {
		c.setConnected(false)
		return c.client.Close()
	}
	return nil
}

// IsConnected 檢查連線狀態（同時確認標誌與底層 socket 存活）
func (c *ModbusTCPConnector) IsConnected() bool {
	if !c.isConnected() {
		return false
	}
	if c.client != nil && !c.client.IsConnected() {
		return false
	}
	return true
}

// ProtocolType 取得協議類型
func (c *ModbusTCPConnector) ProtocolType() schema.ProtocolType {
	return schema.ProtocolModbusTCP
}

// TestConnection 測試連線
func (c *ModbusTCPConnector) TestConnection(ctx context.Context) error {
	// 確保連線
	if err := c.ensureConnection(); err != nil {
		return err
	}

	// 短連接模式：測試完成後自動斷線
	defer c.afterOperation()

	// 嘗試讀取一個保持暫存器來測試連線
	_, err := c.client.ReadHoldingRegisters(0, 1)
	return err
}

// Read 讀取資料
func (c *ModbusTCPConnector) Read(ctx context.Context, req connector.ReadRequest) (connector.ReadResult, error) {
	// 確保連線（短連接模式下會自動重連，長連接模式下若底層斷開會自癒重連）
	if err := c.ensureConnection(); err != nil {
		return connector.ReadResult{
			Quality:   schema.QualityBad,
			Timestamp: time.Now(),
			Error:     err.Error(),
		}, err
	}

	// 短連接模式：操作完成後自動斷線
	defer c.afterOperation()

	result, err := c.executeRead(req)
	if err != nil && c.IsPersistentMode() && isConnectionError(err) {
		// 連線中斷時嘗試單次自癒重連與重試
		c.setConnected(false)
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
func (c *ModbusTCPConnector) SetPersistentConnection(enabled bool) {
	c.stateMu.Lock()
	c.persistentMode = enabled
	c.stateMu.Unlock()
}

// IsPersistentMode 檢查當前是否為長連接模式
func (c *ModbusTCPConnector) IsPersistentMode() bool {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.persistentMode
}

// Disconnect 顯式斷線
func (c *ModbusTCPConnector) Disconnect() error {
	return c.Close()
}

// Reconnect 重新連線，用於斷線後的自動恢復
func (c *ModbusTCPConnector) Reconnect(ctx context.Context) error {
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

// ensureConnection 確保連線已建立
func (c *ModbusTCPConnector) ensureConnection() error {
	if c.IsPersistentMode() {
		if !c.isConnected() {
			if c.client == nil {
				return fmt.Errorf("未連線")
			}
			if err := c.client.Connect(); err != nil {
				return fmt.Errorf("連線失敗: %w", err)
			}
			c.setConnected(true)
			return nil
		}
		// 若已標記連線但底層連線已斷開，自動自癒重連
		if c.client != nil && !c.client.IsConnected() {
			c.client.Close()
			if err := c.client.Connect(); err != nil {
				c.setConnected(false)
				return fmt.Errorf("連線失敗: %w", err)
			}
			c.setConnected(true)
		}
		return nil
	}

	// 短連接模式：每次操作前重新連線
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
func (c *ModbusTCPConnector) afterOperation() {
	if !c.IsPersistentMode() && c.isConnected() {
		// 短連接模式：操作完成後斷線
		c.client.Close()
		c.setConnected(false)
	}
}

func (c *ModbusTCPConnector) setConnected(connected bool) {
	c.stateMu.Lock()
	c.connected = connected
	c.stateMu.Unlock()
}

func (c *ModbusTCPConnector) isConnected() bool {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.connected
}

func isConnectionError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, modbus.ErrConnectionClosed) ||
		errors.Is(err, io.EOF) ||
		errors.Is(err, io.ErrUnexpectedEOF) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "連線已關閉") ||
		strings.Contains(msg, "closed network connection") ||
		strings.Contains(msg, "connection reset") ||
		strings.Contains(msg, "broken pipe") ||
		strings.Contains(msg, "eof")
}
