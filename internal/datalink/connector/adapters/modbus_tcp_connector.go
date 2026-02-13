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
// Modbus TCP 連接器
// =============================================================================

// ModbusTCPConnector Modbus TCP 協議連接器適配器
type ModbusTCPConnector struct {
	client         *modbus.ModbusClient
	transport      modbus.Transport
	config         schema.ConnectionConfigModbusTCP
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
	if !c.persistentMode {
		c.persistentMode = true
	}

	// 連線（僅在長連接模式下立即連線）
	if c.persistentMode {
		if err := c.client.Connect(); err != nil {
			return fmt.Errorf("Modbus TCP 連線失敗: %w", err)
		}
		c.connected = true
	}

	return nil
}

// Close 關閉連線
func (c *ModbusTCPConnector) Close() error {
	if c.client != nil {
		c.connected = false
		return c.client.Close()
	}
	return nil
}

// IsConnected 檢查連線狀態
func (c *ModbusTCPConnector) IsConnected() bool {
	return c.connected
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
	// 確保連線（短連接模式下會自動重連）
	if err := c.ensureConnection(); err != nil {
		return connector.ReadResult{
			Quality:   schema.QualityBad,
			Timestamp: time.Now(),
			Error:     err.Error(),
		}, err
	}

	// 短連接模式：操作完成後自動斷線
	defer c.afterOperation()

	result := connector.ReadResult{
		Timestamp: time.Now(),
		Quality:   schema.QualityGood,
	}

	// 解析地址 (格式: "40001" 或 "HR0" 或 "0")
	address, function, err := parseModbusAddress(req.Address, req.Function)
	if err != nil {
		result.Quality = schema.QualityBad
		result.Error = err.Error()
		return result, err
	}

	// 計算需要讀取的暫存器數量
	count := schema.RegisterCountForDataType(req.DataType)
	if req.Count > 0 {
		count = req.Count
	}

	// 根據功能碼讀取
	switch function {
	case "01", "coil", "FC01":
		// 讀取線圈
		values, err := c.client.ReadCoils(address, uint16(count))
		if err != nil {
			result.Quality = schema.QualityBad
			result.Error = err.Error()
			return result, err
		}
		if len(values) > 0 {
			result.Value = values[0]
		}

	case "02", "discrete", "FC02":
		// 讀取離散輸入
		values, err := c.client.ReadDiscreteInputs(address, uint16(count))
		if err != nil {
			result.Quality = schema.QualityBad
			result.Error = err.Error()
			return result, err
		}
		if len(values) > 0 {
			result.Value = values[0]
		}

	case "03", "holding", "FC03", "":
		// 讀取保持暫存器 (預設)
		values, err := c.client.ReadHoldingRegisters(address, uint16(count))
		if err != nil {
			result.Quality = schema.QualityBad
			result.Error = err.Error()
			return result, err
		}
		result.RawBytes = uint16SliceToBytes(values)
		result.Value = convertModbusValue(values, req.DataType)

	case "04", "input", "FC04":
		// 讀取輸入暫存器
		values, err := c.client.ReadInputRegisters(address, uint16(count))
		if err != nil {
			result.Quality = schema.QualityBad
			result.Error = err.Error()
			return result, err
		}
		result.RawBytes = uint16SliceToBytes(values)
		result.Value = convertModbusValue(values, req.DataType)

	default:
		result.Quality = schema.QualityBad
		result.Error = fmt.Sprintf("不支援的功能碼: %s", function)
		return result, fmt.Errorf("不支援的功能碼: %s", function)
	}

	return result, nil
}

// Write 寫入資料
func (c *ModbusTCPConnector) Write(ctx context.Context, req connector.WriteRequest) error {
	// 確保連線（短連接模式下會自動重連）
	if err := c.ensureConnection(); err != nil {
		return err
	}

	// 短連接模式：操作完成後自動斷線
	defer c.afterOperation()

	// 解析地址
	address, function, err := parseModbusAddress(req.Address, req.Function)
	if err != nil {
		return err
	}

	switch function {
	case "coil", "05", "FC05", "01", "FC01":
		// 寫入單個線圈
		value, ok := req.Value.(bool)
		if !ok {
			return fmt.Errorf("線圈寫入需要布林值")
		}
		return c.client.WriteSingleCoil(address, value)

	case "holding", "06", "FC06", "", "03", "FC03":
		// 寫入單個保持暫存器
		value, err := toUint16(req.Value)
		if err != nil {
			return err
		}
		return c.client.WriteSingleRegister(address, value)

	case "15", "FC15":
		// 寫入多個線圈
		values, ok := req.Value.([]bool)
		if !ok {
			return fmt.Errorf("多線圈寫入需要布林切片")
		}
		return c.client.WriteMultipleCoils(address, values)

	case "16", "FC16":
		// 寫入多個暫存器
		values, err := toUint16Slice(req.Value)
		if err != nil {
			return err
		}
		return c.client.WriteMultipleRegisters(address, values)

	default:
		return fmt.Errorf("不支援的寫入功能碼: %s", function)
	}
}

// =============================================================================
// PersistentConnection 介面實作（長連接支援）
// =============================================================================

/**
 * SetPersistentConnection 設定是否使用長連接模式
 * @param enabled true 啟用長連接，false 使用短連接
 *
 * 長連接模式：連線保持開啟，多次操作重用同一連線
 * 短連接模式：每次操作後自動斷線，下次操作重新連線
 */
func (c *ModbusTCPConnector) SetPersistentConnection(enabled bool) {
	c.persistentMode = enabled
}

/**
 * IsPersistentMode 檢查當前是否為長連接模式
 * @returns bool 是否為長連接模式
 */
func (c *ModbusTCPConnector) IsPersistentMode() bool {
	return c.persistentMode
}

/**
 * Disconnect 顯式斷線
 * @returns error 斷線錯誤
 *
 * 在長連接模式下，此方法可用於手動關閉連線
 * 在短連接模式下，此方法等同於 Close()
 */
func (c *ModbusTCPConnector) Disconnect() error {
	return c.Close()
}

/**
 * Reconnect 重新連線
 * @param ctx 上下文
 * @returns error 重新連線錯誤
 *
 * 用於斷線後的自動恢復
 */
func (c *ModbusTCPConnector) Reconnect(ctx context.Context) error {
	// 檢查 client 是否已初始化
	if c.client == nil {
		return fmt.Errorf("客戶端未初始化，請先呼叫 Connect")
	}

	// 先關閉現有連線
	c.client.Close()

	// 重新連線
	if err := c.client.Connect(); err != nil {
		c.connected = false
		return fmt.Errorf("重新連線失敗: %w", err)
	}

	c.connected = true
	return nil
}

/**
 * ensureConnection 確保連線已建立（用於短連接模式）
 * @returns error 連線錯誤
 */
func (c *ModbusTCPConnector) ensureConnection() error {
	if c.persistentMode {
		// 長連接模式：檢查連線是否有效
		if !c.connected {
			return fmt.Errorf("未連線")
		}
		return nil
	}

	// 短連接模式：每次操作前重新連線
	if c.client == nil {
		return fmt.Errorf("客戶端未初始化")
	}

	// 檢查是否需要重新連線
	if !c.connected {
		if err := c.client.Connect(); err != nil {
			return fmt.Errorf("連線失敗: %w", err)
		}
		c.connected = true
	}

	return nil
}

/**
 * afterOperation 操作後處理（用於短連接模式）
 */
func (c *ModbusTCPConnector) afterOperation() {
	if !c.persistentMode && c.connected {
		// 短連接模式：操作完成後斷線
		c.client.Close()
		c.connected = false
	}
}
