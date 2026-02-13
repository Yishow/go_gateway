package adapters

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"strings"

	"go-gateway/internal/datalink/schema"
)

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
