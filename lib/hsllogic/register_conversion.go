package hsllogic

import "math"

// =============================================================================
// 暫存器陣列轉換
// =============================================================================

// RegistersToValue 將暫存器陣列轉換為指定類型
func (c *DataConverter) RegistersToValue(registers []uint16, dataType DataType) interface{} {
	if len(registers) == 0 {
		return nil
	}

	switch dataType {
	case DataTypeBool:
		return registers[0] != 0
	case DataTypeInt16:
		return int16(registers[0]) // #nosec G115 -- reinterpret the uint16 register bits as signed int16.
	case DataTypeUint16:
		return registers[0]
	case DataTypeInt32:
		return c.transform.RegistersToInt32(registers)
	case DataTypeUint32:
		return c.transform.RegistersToUint32(registers)
	case DataTypeFloat32:
		return c.transform.RegistersToFloat32(registers)
	case DataTypeInt64:
		if len(registers) >= 4 {
			return registersToInt64(registers, c.transform.Format)
		}
	case DataTypeUint64:
		if len(registers) >= 4 {
			return uint64(registersToInt64(registers, c.transform.Format)) // #nosec G115 -- reinterpret the assembled 64-bit register bits as unsigned.
		}
	case DataTypeFloat64:
		if len(registers) >= 4 {
			bits := uint64(registersToInt64(registers, c.transform.Format)) // #nosec G115 -- preserve the assembled 64-bit register bits for IEEE-754 decoding.
			return math.Float64frombits(bits)
		}
	}
	return registers[0]
}

// registersToInt64 將 4 個暫存器轉換為 int64
func registersToInt64(registers []uint16, format DataFormat) int64 {
	if len(registers) < 4 {
		return 0
	}

	switch format {
	case DataFormatCDAB, DataFormatDCBA:
		return int64(uint64(registers[3])<<48 | uint64(registers[2])<<32 |
			uint64(registers[1])<<16 | uint64(registers[0])) // #nosec G115 -- reinterpret assembled 64-bit register bits as signed.
	default: // ABCD, BADC
		return int64(uint64(registers[0])<<48 | uint64(registers[1])<<32 |
			uint64(registers[2])<<16 | uint64(registers[3])) // #nosec G115 -- reinterpret assembled 64-bit register bits as signed.
	}
}

// =============================================================================
// 類型轉換輔助函數
// =============================================================================

// RegistersToValues 將暫存器陣列轉換為指定類型的值切片
// 用於批量讀取時，將連續暫存器解析為多個值
// 例如: 讀取 10 個 Int32，傳入 20 個 uint16 暫存器，返回 10 個 int32 值
func (c *DataConverter) RegistersToValues(registers []uint16, dataType DataType, count int) []interface{} {
	if len(registers) == 0 || count <= 0 {
		return nil
	}

	regPerValue := RegisterCountForDataType(dataType)
	results := make([]interface{}, 0, count)

	for i := 0; i < count && i*regPerValue < len(registers); i++ {
		startIdx := i * regPerValue
		endIdx := startIdx + regPerValue
		if endIdx > len(registers) {
			break
		}
		subRegs := registers[startIdx:endIdx]
		val := c.RegistersToValue(subRegs, dataType)
		results = append(results, val)
	}

	return results
}
