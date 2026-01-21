// Package adapters 提供現有協議驅動到統一 connector.Protocol 介面的適配器。
package adapters

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
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

// =============================================================================
// Modbus RTU 連接器
// =============================================================================

// ModbusRTUConnector Modbus RTU 協議連接器適配器
type ModbusRTUConnector struct {
	client         *modbus.ModbusClient
	transport      modbus.Transport
	config         schema.ConnectionConfigModbusRTU
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
	if !c.persistentMode {
		c.persistentMode = true
	}

	// 連線（僅在長連接模式下立即連線）
	if c.persistentMode {
		if err := c.client.Connect(); err != nil {
			return fmt.Errorf("Modbus RTU 連線失敗: %w", err)
		}
		c.connected = true
	}

	return nil
}

// Close 關閉連線
func (c *ModbusRTUConnector) Close() error {
	if c.client != nil {
		c.connected = false
		return c.client.Close()
	}
	return nil
}

// IsConnected 檢查連線狀態
func (c *ModbusRTUConnector) IsConnected() bool {
	return c.connected
}

// ProtocolType 取得協議類型
func (c *ModbusRTUConnector) ProtocolType() schema.ProtocolType {
	return schema.ProtocolModbusRTU
}

// TestConnection 測試連線
func (c *ModbusRTUConnector) TestConnection(ctx context.Context) error {
	if !c.connected {
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
	tcpConn := &ModbusTCPConnector{client: c.client, connected: c.connected}
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

	tcpConn := &ModbusTCPConnector{client: c.client, connected: c.connected}
	return tcpConn.Write(ctx, req)
}

// =============================================================================
// ModbusRTU PersistentConnection 介面實作
// =============================================================================

// SetPersistentConnection 設定是否使用長連接模式
func (c *ModbusRTUConnector) SetPersistentConnection(enabled bool) {
	c.persistentMode = enabled
}

// IsPersistentMode 檢查當前是否為長連接模式
func (c *ModbusRTUConnector) IsPersistentMode() bool {
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
		c.connected = false
		return fmt.Errorf("重新連線失敗: %w", err)
	}

	c.connected = true
	return nil
}

// ensureConnection 確保連線已建立（用於短連接模式）
func (c *ModbusRTUConnector) ensureConnection() error {
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
func (c *ModbusRTUConnector) afterOperation() {
	if !c.persistentMode && c.connected {
		c.client.Close()
		c.connected = false
	}
}

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

// =============================================================================
// 輔助函數
// =============================================================================

// parseModbusAddress 解析 Modbus 地址字串
// 支援格式:
//   - "40001" (傳統格式，4xxxx = 保持暫存器)
//   - "HR100" (別名格式)
//   - "100" (純數字，需搭配 function 參數)
func parseModbusAddress(addressStr, function string) (uint16, string, error) {
	addressStr = strings.TrimSpace(addressStr)
	function = strings.TrimSpace(strings.ToLower(function))

	// 別名格式處理
	upperAddr := strings.ToUpper(addressStr)
	if strings.HasPrefix(upperAddr, "HR") {
		// 保持暫存器
		addr, err := strconv.ParseUint(addressStr[2:], 10, 16)
		return uint16(addr), "03", err
	} else if strings.HasPrefix(upperAddr, "IR") {
		// 輸入暫存器
		addr, err := strconv.ParseUint(addressStr[2:], 10, 16)
		return uint16(addr), "04", err
	} else if strings.HasPrefix(upperAddr, "C") || strings.HasPrefix(upperAddr, "CO") {
		// 線圈
		prefix := "C"
		if strings.HasPrefix(upperAddr, "CO") {
			prefix = "CO"
		}
		addr, err := strconv.ParseUint(addressStr[len(prefix):], 10, 16)
		return uint16(addr), "01", err
	} else if strings.HasPrefix(upperAddr, "DI") {
		// 離散輸入
		addr, err := strconv.ParseUint(addressStr[2:], 10, 16)
		return uint16(addr), "02", err
	}

	// 數字格式處理
	addr, err := strconv.ParseUint(addressStr, 10, 32)
	if err != nil {
		return 0, "", fmt.Errorf("無效的地址格式: %s", addressStr)
	}

	// 傳統 Modbus 地址格式
	if addr >= 40001 && addr <= 49999 {
		return uint16(addr - 40001), "03", nil
	} else if addr >= 30001 && addr <= 39999 {
		return uint16(addr - 30001), "04", nil
	} else if addr >= 10001 && addr <= 19999 {
		return uint16(addr - 10001), "02", nil
	} else if addr >= 1 && addr <= 9999 {
		return uint16(addr - 1), "01", nil
	}

	// 純數字，使用指定的 function
	if function == "" {
		function = "03" // 預設為保持暫存器
	}
	return uint16(addr), function, nil
}

// convertModbusValue 將 Modbus 暫存器值轉換為指定型別
func convertModbusValue(registers []uint16, dataType schema.DataType) interface{} {
	if len(registers) == 0 {
		return nil
	}

	switch dataType {
	case schema.DataTypeBool:
		return registers[0] != 0
	case schema.DataTypeInt16:
		return int16(registers[0])
	case schema.DataTypeUint16:
		return registers[0]
	case schema.DataTypeInt32:
		if len(registers) >= 2 {
			return int32(uint32(registers[0])<<16 | uint32(registers[1]))
		}
	case schema.DataTypeUint32:
		if len(registers) >= 2 {
			return uint32(registers[0])<<16 | uint32(registers[1])
		}
	case schema.DataTypeFloat32:
		if len(registers) >= 2 {
			bits := uint32(registers[0])<<16 | uint32(registers[1])
			return math.Float32frombits(bits)
		}
	case schema.DataTypeInt64:
		if len(registers) >= 4 {
			val := uint64(registers[0])<<48 | uint64(registers[1])<<32 |
				uint64(registers[2])<<16 | uint64(registers[3])
			return int64(val)
		}
	case schema.DataTypeUint64:
		if len(registers) >= 4 {
			return uint64(registers[0])<<48 | uint64(registers[1])<<32 |
				uint64(registers[2])<<16 | uint64(registers[3])
		}
	case schema.DataTypeFloat64:
		if len(registers) >= 4 {
			bits := uint64(registers[0])<<48 | uint64(registers[1])<<32 |
				uint64(registers[2])<<16 | uint64(registers[3])
			return math.Float64frombits(bits)
		}
	}

	// 預設返回第一個暫存器值
	return registers[0]
}

// uint16SliceToBytes 將 uint16 切片轉換為位元組切片
func uint16SliceToBytes(values []uint16) []byte {
	bytes := make([]byte, len(values)*2)
	for i, v := range values {
		binary.BigEndian.PutUint16(bytes[i*2:], v)
	}
	return bytes
}

// toUint16 將任意值轉換為 uint16
func toUint16(v interface{}) (uint16, error) {
	switch val := v.(type) {
	case int:
		return uint16(val), nil
	case int16:
		return uint16(val), nil
	case int32:
		return uint16(val), nil
	case int64:
		return uint16(val), nil
	case uint:
		return uint16(val), nil
	case uint16:
		return val, nil
	case uint32:
		return uint16(val), nil
	case uint64:
		return uint16(val), nil
	case float32:
		return uint16(val), nil
	case float64:
		return uint16(val), nil
	default:
		return 0, fmt.Errorf("無法轉換為 uint16: %T", v)
	}
}

// toUint16Slice 將任意切片轉換為 uint16 切片
func toUint16Slice(v interface{}) ([]uint16, error) {
	switch val := v.(type) {
	case []uint16:
		return val, nil
	case []int:
		result := make([]uint16, len(val))
		for i, n := range val {
			result[i] = uint16(n)
		}
		return result, nil
	case []interface{}:
		result := make([]uint16, len(val))
		for i, n := range val {
			u, err := toUint16(n)
			if err != nil {
				return nil, err
			}
			result[i] = u
		}
		return result, nil
	default:
		return nil, fmt.Errorf("無法轉換為 uint16 切片: %T", v)
	}
}

// =============================================================================
// 註冊連接器
// =============================================================================

func init() {
	connector.Register(schema.ProtocolModbusTCP, NewModbusTCPConnector)
	connector.Register(schema.ProtocolModbusRTU, NewModbusRTUConnector)
	connector.Register(schema.ProtocolModbusUDP, NewModbusUDPConnector)
}
