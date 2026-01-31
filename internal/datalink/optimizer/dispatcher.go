package optimizer

import (
	"encoding/binary"
	"fmt"
	"math"

	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 類型定義
// =============================================================================

// ByteOrder 字節序
type ByteOrder int

const (
	// BigEndian 大端序
	BigEndian ByteOrder = iota
	// LittleEndian 小端序
	LittleEndian
)

// DispatcherConfig 分發器配置
type DispatcherConfig struct {
	ByteOrder ByteOrder
}

// DefaultDispatcherConfig 預設分發器配置
func DefaultDispatcherConfig() DispatcherConfig {
	return DispatcherConfig{
		ByteOrder: BigEndian,
	}
}

// =============================================================================
// Dispatcher 資料分發器
// =============================================================================

// Dispatcher 負責將聚合回應的原始數據分發回各個點位
type Dispatcher struct {
	config DispatcherConfig
}

// NewDispatcher 建立新的分發器
func NewDispatcher(config DispatcherConfig) *Dispatcher {
	return &Dispatcher{config: config}
}

// Dispatch 分發 Block 的原始數據到各個點位
// 返回 map[pointID]value
func (d *Dispatcher) Dispatch(block Block, rawBytes []byte) (map[string]interface{}, error) {
	// 驗證資料長度
	expectedLen := block.Length * 2 // 每個暫存器 2 bytes
	if len(rawBytes) < expectedLen {
		return nil, fmt.Errorf("資料長度不足: 預期 %d bytes，實際 %d bytes", expectedLen, len(rawBytes))
	}

	values := make(map[string]interface{})

	for _, point := range block.Points {
		offset, exists := block.PointOffsets[point.ID]
		if !exists {
			continue
		}

		byteOffset := offset * 2 // 轉換為字節偏移
		value, err := d.extractValue(rawBytes, byteOffset, point.DataType)
		if err != nil {
			return nil, fmt.Errorf("提取點位 %s 的值失敗: %w", point.ID, err)
		}

		values[point.ID] = value
	}

	return values, nil
}

// extractValue 從原始數據中提取指定類型的值
func (d *Dispatcher) extractValue(data []byte, offset int, dataType schema.DataType) (interface{}, error) {
	var order binary.ByteOrder
	if d.config.ByteOrder == BigEndian {
		order = binary.BigEndian
	} else {
		order = binary.LittleEndian
	}

	switch dataType {
	case schema.DataTypeBool:
		if offset >= len(data) {
			return nil, fmt.Errorf("偏移量超出範圍")
		}
		// 非零為 true
		return data[offset] != 0 || (offset+1 < len(data) && data[offset+1] != 0), nil

	case schema.DataTypeInt16:
		if offset+2 > len(data) {
			return nil, fmt.Errorf("資料不足以讀取 int16")
		}
		return int16(order.Uint16(data[offset:])), nil

	case schema.DataTypeUint16:
		if offset+2 > len(data) {
			return nil, fmt.Errorf("資料不足以讀取 uint16")
		}
		return order.Uint16(data[offset:]), nil

	case schema.DataTypeInt32:
		if offset+4 > len(data) {
			return nil, fmt.Errorf("資料不足以讀取 int32")
		}
		return int32(order.Uint32(data[offset:])), nil

	case schema.DataTypeUint32:
		if offset+4 > len(data) {
			return nil, fmt.Errorf("資料不足以讀取 uint32")
		}
		return order.Uint32(data[offset:]), nil

	case schema.DataTypeFloat32:
		if offset+4 > len(data) {
			return nil, fmt.Errorf("資料不足以讀取 float32")
		}
		bits := order.Uint32(data[offset:])
		return math.Float32frombits(bits), nil

	case schema.DataTypeInt64:
		if offset+8 > len(data) {
			return nil, fmt.Errorf("資料不足以讀取 int64")
		}
		return int64(order.Uint64(data[offset:])), nil

	case schema.DataTypeUint64:
		if offset+8 > len(data) {
			return nil, fmt.Errorf("資料不足以讀取 uint64")
		}
		return order.Uint64(data[offset:]), nil

	case schema.DataTypeFloat64:
		if offset+8 > len(data) {
			return nil, fmt.Errorf("資料不足以讀取 float64")
		}
		bits := order.Uint64(data[offset:])
		return math.Float64frombits(bits), nil

	case schema.DataTypeString:
		// 預設讀取 16 bytes (8 暫存器)
		endOffset := offset + 16
		if endOffset > len(data) {
			endOffset = len(data)
		}
		// 去除尾部的 null 字元
		str := string(data[offset:endOffset])
		for i := 0; i < len(str); i++ {
			if str[i] == 0 {
				str = str[:i]
				break
			}
		}
		return str, nil

	default:
		return nil, fmt.Errorf("不支援的資料類型: %s", dataType)
	}
}

// DispatchMultiple 批量分發多個 Block 的數據
func (d *Dispatcher) DispatchMultiple(blocks []Block, rawDataMap map[int][]byte) (map[string]interface{}, error) {
	allValues := make(map[string]interface{})

	for i, block := range blocks {
		rawBytes, exists := rawDataMap[i]
		if !exists {
			continue
		}

		values, err := d.Dispatch(block, rawBytes)
		if err != nil {
			return nil, fmt.Errorf("Block %d 分發失敗: %w", i, err)
		}

		for k, v := range values {
			allValues[k] = v
		}
	}

	return allValues, nil
}
