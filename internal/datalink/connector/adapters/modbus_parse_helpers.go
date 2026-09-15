package adapters

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"strings"

	"go-gateway/internal/datalink/schema"
	"go-gateway/lib/hsllogic"
)

// =============================================================================
// 輔助函數
// =============================================================================

// parseModbusAddress 解析 Modbus 地址字串
// 支援格式:
//   - "40001" (傳統格式，4xxxx = 保持暫存器)
//   - "HR100" (別名格式)
//   - "100" (純數字，需搭配 function 參數)
func parseModbusAddress(addressStr, function string) (address uint16, functionCode string, parseErr error) {
	addressStr = strings.TrimSpace(addressStr)
	function = strings.TrimSpace(strings.ToLower(function))

	// 別名格式處理
	upperAddr := strings.ToUpper(addressStr)
	switch {
	case strings.HasPrefix(upperAddr, "HR"):
		// 保持暫存器
		addr, err := strconv.ParseUint(addressStr[2:], 10, 16)
		return uint16(addr), "03", err
	case strings.HasPrefix(upperAddr, "IR"):
		// 輸入暫存器
		addr, err := strconv.ParseUint(addressStr[2:], 10, 16)
		return uint16(addr), "04", err
	case strings.HasPrefix(upperAddr, "C") || strings.HasPrefix(upperAddr, "CO"):
		// 線圈
		prefix := "C"
		if strings.HasPrefix(upperAddr, "CO") {
			prefix = "CO"
		}
		addr, err := strconv.ParseUint(addressStr[len(prefix):], 10, 16)
		return uint16(addr), "01", err
	case strings.HasPrefix(upperAddr, "DI"):
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
	switch {
	case addr >= 40001 && addr <= 49999:
		return uint16(addr - 40001), "03", nil
	case addr >= 30001 && addr <= 39999:
		return uint16(addr - 30001), "04", nil
	case addr >= 10001 && addr <= 19999:
		return uint16(addr - 10001), "02", nil
	case addr >= 1 && addr <= 9999:
		return uint16(addr - 1), "01", nil
	}

	// 純數字，使用指定的 function
	if function == "" {
		function = "03" // 預設為保持暫存器
	}
	if addr > math.MaxUint16 {
		return 0, "", fmt.Errorf("modbus 地址超出 16 位元範圍: %s", addressStr)
	}
	return uint16(addr), function, nil
}

// parseModbusDataFormat 將字串規格化為 hsllogic 字節序；空字串或未知值沿用歷史 Modbus 行為（ABCD）。
func parseModbusDataFormat(format string) hsllogic.DataFormat {
	switch strings.ToUpper(strings.TrimSpace(format)) {
	case string(hsllogic.DataFormatABCD):
		return hsllogic.DataFormatABCD
	case string(hsllogic.DataFormatBADC):
		return hsllogic.DataFormatBADC
	case string(hsllogic.DataFormatCDAB):
		return hsllogic.DataFormatCDAB
	case string(hsllogic.DataFormatDCBA):
		return hsllogic.DataFormatDCBA
	default:
		return hsllogic.DataFormatABCD
	}
}

func effectiveModbusDataFormat(requestFormat, connectionFormat string) string {
	if strings.TrimSpace(requestFormat) != "" {
		return requestFormat
	}
	return connectionFormat
}

// convertModbusValue 將 Modbus 暫存器值轉換為指定型別（依 dataFormat 解多暫存器數值）。
func convertModbusValue(registers []uint16, dataType schema.DataType, dataFormat string) interface{} {
	if len(registers) == 0 {
		return nil
	}
	converter := hsllogic.NewDataConverter(parseModbusDataFormat(dataFormat))
	return converter.RegistersToValue(registers, hsllogic.DataType(dataType))
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
		if val >= 0 && val <= math.MaxUint16 {
			return uint16(val), nil
		}
	case int16:
		if val >= 0 {
			return uint16(val), nil
		}
	case int32:
		if val >= 0 && val <= math.MaxUint16 {
			return uint16(val), nil
		}
	case int64:
		if val >= 0 && val <= math.MaxUint16 {
			return uint16(val), nil
		}
	case uint:
		if val <= math.MaxUint16 {
			return uint16(val), nil
		}
	case uint16:
		return val, nil
	case uint32:
		if val <= math.MaxUint16 {
			return uint16(val), nil
		}
	case uint64:
		if val <= math.MaxUint16 {
			return uint16(val), nil
		}
	case float32:
		if val >= 0 && val <= math.MaxUint16 && math.Trunc(float64(val)) == float64(val) {
			return uint16(val), nil
		}
	case float64:
		if val >= 0 && val <= math.MaxUint16 && math.Trunc(float64(val)) == float64(val) {
			return uint16(val), nil
		}
	}
	return 0, fmt.Errorf("uint16 值必須是 0 到 65535 之間的整數")
}

// toUint16Slice 將任意切片轉換為 uint16 切片
func toUint16Slice(v interface{}) ([]uint16, error) {
	switch val := v.(type) {
	case []uint16:
		return val, nil
	case []int:
		result := make([]uint16, len(val))
		for i, n := range val {
			u, err := toUint16(n)
			if err != nil {
				return nil, err
			}
			result[i] = u
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
