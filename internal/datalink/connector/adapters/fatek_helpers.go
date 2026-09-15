package adapters

import (
	"fmt"
	"strconv"
	"strings"

	"go-gateway/internal/datalink/schema"
)

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
		"D", "R", "WM", "WS", "WX", "WY", // 16 位元
		"M", "S", "T", "C", "X", "Y", // 位元
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
		return int16(values[0]) // #nosec G115 -- Decode the fixed-width FATEK register bit pattern, including its sign bit.
	case schema.DataTypeUint16:
		return uint16(values[0]) // #nosec G115 -- Decode the fixed-width FATEK register bit pattern, including its sign bit.
	case schema.DataTypeInt32:
		if len(values) >= 1 && width == 32 {
			return int32(values[0]) // #nosec G115 -- Decode the fixed-width FATEK register bit pattern, including its sign bit.
		} else if len(values) >= 2 {
			return int32(values[0])<<16 | int32(values[1]) // #nosec G115 -- Decode the fixed-width FATEK register bit pattern, including its sign bit.
		}
	case schema.DataTypeUint32:
		if len(values) >= 1 && width == 32 {
			return uint32(values[0]) // #nosec G115 -- Decode the fixed-width FATEK register bit pattern, including its sign bit.
		} else if len(values) >= 2 {
			return uint32(values[0])<<16 | uint32(values[1]) // #nosec G115 -- Decode the fixed-width FATEK register bit pattern, including its sign bit.
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
		bytes[i*4] = byte((v >> 24) & 0xff)
		bytes[i*4+1] = byte((v >> 16) & 0xff)
		bytes[i*4+2] = byte((v >> 8) & 0xff)
		bytes[i*4+3] = byte(v & 0xff)
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
