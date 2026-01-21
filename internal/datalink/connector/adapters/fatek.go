package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/fatek"
)

// =============================================================================
// FATEK FBs 連接器
// =============================================================================

// FatekConnector FATEK FBs 協議連接器適配器
type FatekConnector struct {
	client         *fatek.FatekClient
	transport      fatek.Transport
	config         schema.ConnectionConfigFatekFBs
	connected      bool
	persistentMode bool // 長連接模式標誌
}

// NewFatekConnector 建立新的 FATEK 連接器
func NewFatekConnector() connector.Protocol {
	return &FatekConnector{}
}

// Connect 建立連線
func (c *FatekConnector) Connect(ctx context.Context, configJSON string) error {
	// 解析配置
	if err := json.Unmarshal([]byte(configJSON), &c.config); err != nil {
		return fmt.Errorf("解析 FATEK 配置失敗: %w", err)
	}

	// 設定預設值
	if c.config.Timeout == 0 {
		c.config.Timeout = 5
	}

	// 根據模式建立傳輸層
	switch strings.ToLower(c.config.Mode) {
	case "tcp":
		if c.config.Port == 0 {
			c.config.Port = 500
		}
		transport := fatek.NewTCPTransport(c.config.Host, c.config.Port)
		transport.Timeout = time.Duration(c.config.Timeout) * time.Second
		c.transport = transport
	case "serial":
		if c.config.BaudRate == 0 {
			c.config.BaudRate = 9600
		}
		c.transport = fatek.NewSerialTransport(
			c.config.SerialPort,
			c.config.BaudRate,
			7,      // DataBits (FATEK 預設 7)
			1,      // StopBits
			"even", // Parity (FATEK 預設 Even)
			time.Duration(c.config.Timeout)*time.Second,
		)
	default:
		return fmt.Errorf("不支援的 FATEK 連線模式: %s", c.config.Mode)
	}

	// 建立客戶端
	c.client = fatek.NewClient(c.transport, int(c.config.StationNo))

	// 預設使用長連接模式
	if !c.persistentMode {
		c.persistentMode = true
	}

	// 連線（僅在長連接模式下立即連線）
	if c.persistentMode {
		if err := c.client.Connect(); err != nil {
			return fmt.Errorf("FATEK 連線失敗: %w", err)
		}
		c.connected = true
	}

	return nil
}

// Close 關閉連線
func (c *FatekConnector) Close() error {
	if c.client != nil {
		c.connected = false
		return c.client.Close()
	}
	return nil
}

// IsConnected 檢查連線狀態
func (c *FatekConnector) IsConnected() bool {
	return c.connected
}

// ProtocolType 取得協議類型
func (c *FatekConnector) ProtocolType() schema.ProtocolType {
	return schema.ProtocolFatekFBs
}

// TestConnection 測試連線
func (c *FatekConnector) TestConnection(ctx context.Context) error {
	if err := c.ensureConnection(); err != nil {
		return err
	}

	defer c.afterOperation()

	// 使用 Loopback Test (Cmd 4E) 測試連線
	ok, err := c.client.LoopbackTest("TEST")
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("Loopback 測試失敗")
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

// =============================================================================
// PersistentConnection 介面實作（長連接支援）
// =============================================================================

// SetPersistentConnection 設定是否使用長連接模式
func (c *FatekConnector) SetPersistentConnection(enabled bool) {
	c.persistentMode = enabled
}

// IsPersistentMode 檢查當前是否為長連接模式
func (c *FatekConnector) IsPersistentMode() bool {
	return c.persistentMode
}

// Disconnect 顯式斷線
func (c *FatekConnector) Disconnect() error {
	return c.Close()
}

// Reconnect 重新連線
func (c *FatekConnector) Reconnect(ctx context.Context) error {
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
func (c *FatekConnector) ensureConnection() error {
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
func (c *FatekConnector) afterOperation() {
	if !c.persistentMode && c.connected {
		c.client.Close()
		c.connected = false
	}
}

// =============================================================================
// 輔助函數
// =============================================================================

// parseFatekAddress 解析 FATEK 地址字串
// 格式: "D0100", "R0", "M100", "X0", "Y0", "D100", etc.
func parseFatekAddress(addressStr string) (symbol string, address int, err error) {
	addressStr = strings.TrimSpace(strings.ToUpper(addressStr))
	if len(addressStr) < 2 {
		return "", 0, fmt.Errorf("無效的 FATEK 地址: %s", addressStr)
	}

	// 支援的元件符號
	symbols := []string{
		"DR", "DD", "DWM", "DWS", "DWX", "DWY", // 32 位元
		"D", "R", "WM", "WS", "WX", "WY",         // 16 位元
		"M", "S", "T", "C", "X", "Y",             // 位元
	}

	for _, s := range symbols {
		if strings.HasPrefix(addressStr, s) {
			addrPart := addressStr[len(s):]
			addr, err := strconv.Atoi(addrPart)
			if err != nil {
				return "", 0, fmt.Errorf("無效的 FATEK 地址數字: %s", addrPart)
			}
			return s, addr, nil
		}
	}

	return "", 0, fmt.Errorf("無法識別的 FATEK 元件符號: %s", addressStr)
}

// convertFatekValue 將 FATEK 暫存器值轉換為指定型別
func convertFatekValue(values []int, dataType schema.DataType, width int) interface{} {
	if len(values) == 0 {
		return nil
	}

	switch dataType {
	case schema.DataTypeBool:
		return values[0] != 0
	case schema.DataTypeInt16:
		return int16(values[0])
	case schema.DataTypeUint16:
		return uint16(values[0])
	case schema.DataTypeInt32:
		if len(values) >= 1 && width == 32 {
			return int32(values[0])
		} else if len(values) >= 2 {
			return int32(values[0])<<16 | int32(values[1])
		}
	case schema.DataTypeUint32:
		if len(values) >= 1 && width == 32 {
			return uint32(values[0])
		} else if len(values) >= 2 {
			return uint32(values[0])<<16 | uint32(values[1])
		}
	case schema.DataTypeFloat32:
		if len(values) >= 1 && width == 32 {
			// FATEK 32 位元暫存器可能已經包含完整的 float
			// 這裡需要根據實際 PLC 配置處理
			return float32(values[0])
		}
	}

	// 預設返回第一個值
	return values[0]
}

// intSliceToBytes 將 int 切片轉換為位元組切片
func intSliceToBytes(values []int) []byte {
	bytes := make([]byte, len(values)*4)
	for i, v := range values {
		bytes[i*4] = byte(v >> 24)
		bytes[i*4+1] = byte(v >> 16)
		bytes[i*4+2] = byte(v >> 8)
		bytes[i*4+3] = byte(v)
	}
	return bytes
}

// toIntSlice 將任意值轉換為 int 切片
func toIntSlice(v interface{}) ([]int, error) {
	switch val := v.(type) {
	case int:
		return []int{val}, nil
	case []int:
		return val, nil
	case int16:
		return []int{int(val)}, nil
	case int32:
		return []int{int(val)}, nil
	case int64:
		return []int{int(val)}, nil
	case uint16:
		return []int{int(val)}, nil
	case uint32:
		return []int{int(val)}, nil
	case []interface{}:
		result := make([]int, len(val))
		for i, n := range val {
			switch num := n.(type) {
			case int:
				result[i] = num
			case float64:
				result[i] = int(num)
			default:
				return nil, fmt.Errorf("無法轉換元素為 int: %T", n)
			}
		}
		return result, nil
	default:
		return nil, fmt.Errorf("無法轉換為 int 切片: %T", v)
	}
}

// =============================================================================
// 註冊連接器
// =============================================================================

func init() {
	connector.Register(schema.ProtocolFatekFBs, NewFatekConnector)
}
