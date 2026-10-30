package adapters

import (
	"fmt"
	"math"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/fatek"
)

// =============================================================================
// 輔助函數
// =============================================================================

// parseFatekAddress 解析 FATEK 地址字串
// 格式: "D0100", "R0", "M100", "X0", "Y0", "D100", etc.
func parseFatekAddress(addressStr string) (symbol string, address int, err error) {
	return fatek.ParseComponentAddress(addressStr)
}

// convertFatekValue 將 FATEK 暫存器值轉換為指定型別
func convertFatekValue(values []int, dataType schema.DataType, width int) (any, error) {
	if len(values) == 0 {
		return nil, fmt.Errorf("FATEK conversion requires at least one component")
	}

	switch dataType {
	case schema.DataTypeBool:
		return values[0] != 0, nil
	case schema.DataTypeInt16:
		return int16(values[0]), nil // #nosec G115 -- Decode the fixed-width FATEK register bit pattern, including its sign bit.
	case schema.DataTypeUint16:
		return uint16(values[0]), nil // #nosec G115 -- Decode the fixed-width FATEK register bit pattern, including its sign bit.
	case schema.DataTypeInt32:
		bits, err := assembleFatekBits(values, dataType, width, 32)
		if err != nil {
			return nil, err
		}
		return int32(uint32(bits)), nil // #nosec G115 -- Reinterpret the assembled 32-bit FATEK bits as signed.
	case schema.DataTypeUint32:
		bits, err := assembleFatekBits(values, dataType, width, 32)
		if err != nil {
			return nil, err
		}
		return uint32(bits), nil // #nosec G115 -- assembleFatekBits with targetBits 32 produces at most 32 bits.
	case schema.DataTypeFloat32:
		bits, err := assembleFatekBits(values, dataType, width, 32)
		if err != nil {
			return nil, err
		}
		return math.Float32frombits(uint32(bits)), nil // #nosec G115 -- assembleFatekBits with targetBits 32 produces at most 32 bits.
	case schema.DataTypeInt64:
		bits, err := assembleFatekBits(values, dataType, width, 64)
		if err != nil {
			return nil, err
		}
		return int64(bits), nil // #nosec G115 -- Reinterpret the assembled 64-bit FATEK bits as signed.
	case schema.DataTypeUint64:
		bits, err := assembleFatekBits(values, dataType, width, 64)
		if err != nil {
			return nil, err
		}
		return bits, nil
	case schema.DataTypeFloat64:
		bits, err := assembleFatekBits(values, dataType, width, 64)
		if err != nil {
			return nil, err
		}
		return math.Float64frombits(bits), nil
	}

	return nil, fmt.Errorf("unsupported FATEK data type: %s", dataType)
}

func assembleFatekBits(values []int, dataType schema.DataType, width, targetBits int) (uint64, error) {
	if width != 16 && width != 32 {
		return 0, fmt.Errorf("unsupported FATEK component width %d for %s", width, dataType)
	}
	required := targetBits / width
	if len(values) < required {
		return 0, fmt.Errorf("insufficient FATEK components for %s: need %d, got %d", dataType, required, len(values))
	}

	mask := (uint64(1) << uint(width)) - 1
	var bits uint64
	for i := range required {
		component := uint64(values[i]) & mask // #nosec G115 -- Preserve the fixed-width FATEK component bits.
		bits = bits<<uint(width) | component
	}
	return bits, nil
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
func toIntSlice(v any) ([]int, error) {
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
	case []any:
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
