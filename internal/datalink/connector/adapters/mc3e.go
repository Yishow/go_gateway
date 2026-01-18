package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/mcprotocol"
)

// =============================================================================
// Mitsubishi MC 3E 連接器
// =============================================================================

// MC3EConnector Mitsubishi MC Protocol 3E Frame 連接器適配器
type MC3EConnector struct {
	client    *mcprotocol.MCClient
	config    schema.ConnectionConfigMC3E
	connected bool
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

	// 建立客戶端
	c.client = mcprotocol.NewClient(c.config.Host, c.config.Port)

	// 連線
	if err := c.client.Connect(); err != nil {
		return fmt.Errorf("MC 3E 連線失敗: %w", err)
	}

	c.connected = true
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
	if !c.connected {
		return fmt.Errorf("未連線")
	}

	// 嘗試讀取 D0 來測試連線
	_, err := c.client.BatchReadWord("D", 0, 1)
	return err
}

// Read 讀取資料
func (c *MC3EConnector) Read(ctx context.Context, req connector.ReadRequest) (connector.ReadResult, error) {
	if !c.connected {
		return connector.ReadResult{
			Quality:   schema.QualityBad,
			Timestamp: time.Now(),
			Error:     "未連線",
		}, fmt.Errorf("未連線")
	}

	result := connector.ReadResult{
		Timestamp: time.Now(),
		Quality:   schema.QualityGood,
	}

	// 解析地址 (格式: "D100", "M0", "X0", "Y0", "R100", "W100", etc.)
	device, address, err := parseMC3EAddress(req.Address)
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
	if isBitDevice(device) {
		// 位元設備
		values, err := c.client.BatchReadBit(device, address, count)
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
		// 計算需要讀取的字組數量
		wordCount := schema.RegisterCountForDataType(req.DataType)
		if req.Count > 1 {
			wordCount = count
		}

		values, err := c.client.BatchReadWord(device, address, wordCount)
		if err != nil {
			result.Quality = schema.QualityBad
			result.Error = err.Error()
			return result, err
		}
		result.RawBytes = intSliceToBytes(values)
		result.Value = convertMC3EValue(values, req.DataType)
	}

	return result, nil
}

// Write 寫入資料
func (c *MC3EConnector) Write(ctx context.Context, req connector.WriteRequest) error {
	if !c.connected {
		return fmt.Errorf("未連線")
	}

	// 解析地址
	device, address, err := parseMC3EAddress(req.Address)
	if err != nil {
		return err
	}

	if isBitDevice(device) {
		// 位元設備寫入
		switch v := req.Value.(type) {
		case bool:
			return c.client.BatchWriteBit(device, address, []bool{v})
		case []bool:
			return c.client.BatchWriteBit(device, address, v)
		default:
			return fmt.Errorf("位元設備需要布林值")
		}
	} else {
		// 字組設備寫入
		values, err := toIntSlice(req.Value)
		if err != nil {
			return err
		}
		return c.client.BatchWriteWord(device, address, values)
	}
}

// =============================================================================
// 輔助函數
// =============================================================================

// parseMC3EAddress 解析 MC Protocol 地址字串
// 格式: "D100", "M0", "X0", "Y0", "R100", "W100", etc.
func parseMC3EAddress(addressStr string) (device string, address int, err error) {
	addressStr = strings.TrimSpace(strings.ToUpper(addressStr))
	if len(addressStr) < 2 {
		return "", 0, fmt.Errorf("無效的 MC 地址: %s", addressStr)
	}

	// 支援的設備碼 (按長度排序以優先匹配較長的名稱)
	devices := []string{
		"ZR", "SD", "SW", "SB", "SM", // 特殊設備
		"D", "W", "R", "B", "F", // 字組設備
		"M", "L", "S", "X", "Y", // 位元設備
		"T", "C", "ST", "CC", "TC", // 計時器/計數器
	}

	for _, d := range devices {
		if strings.HasPrefix(addressStr, d) {
			addrPart := addressStr[len(d):]
			// 處理十六進位地址 (X, Y 設備通常使用八進位)
			var addr int
			var parseErr error
			if d == "X" || d == "Y" {
				// 八進位解析
				addr64, err := strconv.ParseInt(addrPart, 8, 32)
				if err != nil {
					// 嘗試十進位
					addr64, err = strconv.ParseInt(addrPart, 10, 32)
					if err != nil {
						return "", 0, fmt.Errorf("無效的 MC 地址數字: %s", addrPart)
					}
				}
				addr = int(addr64)
			} else {
				addr, parseErr = strconv.Atoi(addrPart)
				if parseErr != nil {
					// 嘗試十六進位
					addr64, err := strconv.ParseInt(addrPart, 16, 32)
					if err != nil {
						return "", 0, fmt.Errorf("無效的 MC 地址數字: %s", addrPart)
					}
					addr = int(addr64)
				}
			}
			return d, addr, nil
		}
	}

	return "", 0, fmt.Errorf("無法識別的 MC 設備碼: %s", addressStr)
}

// isBitDevice 判斷是否為位元設備
func isBitDevice(device string) bool {
	bitDevices := map[string]bool{
		"M": true, "L": true, "S": true,
		"X": true, "Y": true, "B": true,
		"F": true, "SB": true, "SM": true,
	}
	return bitDevices[device]
}

// convertMC3EValue 將 MC Protocol 字組值轉換為指定型別
func convertMC3EValue(values []int, dataType schema.DataType) interface{} {
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
		if len(values) >= 2 {
			// 三菱為 Little Endian (低位在前)
			return int32(uint32(values[1])<<16 | uint32(values[0])&0xFFFF)
		}
		return int32(values[0])
	case schema.DataTypeUint32:
		if len(values) >= 2 {
			return uint32(values[1])<<16 | uint32(values[0])&0xFFFF
		}
		return uint32(values[0])
	case schema.DataTypeFloat32:
		if len(values) >= 2 {
			bits := uint32(values[1])<<16 | uint32(values[0])&0xFFFF
			return math.Float32frombits(bits)
		}
		return float32(values[0])
	case schema.DataTypeInt64:
		if len(values) >= 4 {
			val := uint64(values[3])<<48 | uint64(values[2])<<32 |
				uint64(values[1])<<16 | uint64(values[0])&0xFFFF
			return int64(val)
		}
	case schema.DataTypeUint64:
		if len(values) >= 4 {
			return uint64(values[3])<<48 | uint64(values[2])<<32 |
				uint64(values[1])<<16 | uint64(values[0])&0xFFFF
		}
	case schema.DataTypeFloat64:
		if len(values) >= 4 {
			bits := uint64(values[3])<<48 | uint64(values[2])<<32 |
				uint64(values[1])<<16 | uint64(values[0])&0xFFFF
			return math.Float64frombits(bits)
		}
	}

	// 預設返回第一個值
	return values[0]
}

// =============================================================================
// 註冊連接器
// =============================================================================

func init() {
	connector.Register(schema.ProtocolMC3E, NewMC3EConnector)
}
