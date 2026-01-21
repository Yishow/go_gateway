package hsllogic

import (
	"encoding/binary"
	"math"
)

// =============================================================================
// DataType 數據類型定義
// =============================================================================

// DataType 數據類型常數
type DataType string

const (
	// DataTypeBool 布林值 (1 位元)
	DataTypeBool DataType = "bool"
	// DataTypeInt8 有符號 8 位元整數
	DataTypeInt8 DataType = "int8"
	// DataTypeUint8 無符號 8 位元整數
	DataTypeUint8 DataType = "uint8"
	// DataTypeInt16 有符號 16 位元整數
	DataTypeInt16 DataType = "int16"
	// DataTypeUint16 無符號 16 位元整數
	DataTypeUint16 DataType = "uint16"
	// DataTypeInt32 有符號 32 位元整數
	DataTypeInt32 DataType = "int32"
	// DataTypeUint32 無符號 32 位元整數
	DataTypeUint32 DataType = "uint32"
	// DataTypeInt64 有符號 64 位元整數
	DataTypeInt64 DataType = "int64"
	// DataTypeUint64 無符號 64 位元整數
	DataTypeUint64 DataType = "uint64"
	// DataTypeFloat32 32 位元浮點數
	DataTypeFloat32 DataType = "float32"
	// DataTypeFloat64 64 位元浮點數
	DataTypeFloat64 DataType = "float64"
	// DataTypeString ASCII 字串
	DataTypeString DataType = "string"
	// DataTypeBytes 原始位元組
	DataTypeBytes DataType = "bytes"
)

// =============================================================================
// DataTypeInfo 數據類型資訊
// =============================================================================

// DataTypeInfo 數據類型的元資料
type DataTypeInfo struct {
	// Type 數據類型
	Type DataType
	// ByteSize 位元組大小
	ByteSize int
	// RegisterCount 暫存器數量 (16 位元暫存器)
	RegisterCount int
	// IsSigned 是否為有符號類型
	IsSigned bool
	// IsFloating 是否為浮點類型
	IsFloating bool
}

// dataTypeInfoMap 數據類型資訊映射
var dataTypeInfoMap = map[DataType]DataTypeInfo{
	DataTypeBool:    {DataTypeBool, 1, 1, false, false},
	DataTypeInt8:    {DataTypeInt8, 1, 1, true, false},
	DataTypeUint8:   {DataTypeUint8, 1, 1, false, false},
	DataTypeInt16:   {DataTypeInt16, 2, 1, true, false},
	DataTypeUint16:  {DataTypeUint16, 2, 1, false, false},
	DataTypeInt32:   {DataTypeInt32, 4, 2, true, false},
	DataTypeUint32:  {DataTypeUint32, 4, 2, false, false},
	DataTypeInt64:   {DataTypeInt64, 8, 4, true, false},
	DataTypeUint64:  {DataTypeUint64, 8, 4, false, false},
	DataTypeFloat32: {DataTypeFloat32, 4, 2, false, true},
	DataTypeFloat64: {DataTypeFloat64, 8, 4, false, true},
	DataTypeString:  {DataTypeString, 0, 0, false, false}, // 可變長度
	DataTypeBytes:   {DataTypeBytes, 0, 0, false, false},  // 可變長度
}

// GetDataTypeInfo 取得數據類型資訊
func GetDataTypeInfo(dt DataType) DataTypeInfo {
	if info, ok := dataTypeInfoMap[dt]; ok {
		return info
	}
	return DataTypeInfo{Type: dt}
}

// ByteSizeForDataType 取得數據類型的位元組大小
func ByteSizeForDataType(dt DataType) int {
	return GetDataTypeInfo(dt).ByteSize
}

// RegisterCountForDataType 取得數據類型需要的暫存器數量
func RegisterCountForDataType(dt DataType) int {
	info := GetDataTypeInfo(dt)
	if info.RegisterCount == 0 {
		return 1 // 預設至少 1 個暫存器
	}
	return info.RegisterCount
}

// =============================================================================
// DataConverter 數據轉換器
// =============================================================================

// DataConverter 提供完整的數據類型轉換功能
type DataConverter struct {
	transform *ByteTransform
}

// NewDataConverter 建立新的數據轉換器
func NewDataConverter(format DataFormat) *DataConverter {
	return &DataConverter{
		transform: NewByteTransform(format),
	}
}

// DefaultConverter 使用預設格式的轉換器
var DefaultConverter = NewDataConverter(DataFormatABCD)

// =============================================================================
// 從 []byte 讀取指定類型
// =============================================================================

// ReadBool 從位元組讀取布林值
func (c *DataConverter) ReadBool(buffer []byte, index int) bool {
	if len(buffer) <= index {
		return false
	}
	return buffer[index] != 0
}

// ReadInt8 從位元組讀取 int8
func (c *DataConverter) ReadInt8(buffer []byte, index int) int8 {
	if len(buffer) <= index {
		return 0
	}
	return int8(buffer[index])
}

// ReadUint8 從位元組讀取 uint8
func (c *DataConverter) ReadUint8(buffer []byte, index int) uint8 {
	if len(buffer) <= index {
		return 0
	}
	return buffer[index]
}

// ReadInt16 從位元組讀取 int16
func (c *DataConverter) ReadInt16(buffer []byte, index int) int16 {
	return c.transform.TransformInt16(buffer, index)
}

// ReadUint16 從位元組讀取 uint16
func (c *DataConverter) ReadUint16(buffer []byte, index int) uint16 {
	return c.transform.TransformUint16(buffer, index)
}

// ReadInt32 從位元組讀取 int32
func (c *DataConverter) ReadInt32(buffer []byte, index int) int32 {
	return c.transform.TransformInt32(buffer, index)
}

// ReadUint32 從位元組讀取 uint32
func (c *DataConverter) ReadUint32(buffer []byte, index int) uint32 {
	return c.transform.TransformUint32(buffer, index)
}

// ReadInt64 從位元組讀取 int64
func (c *DataConverter) ReadInt64(buffer []byte, index int) int64 {
	return c.transform.TransformInt64(buffer, index)
}

// ReadUint64 從位元組讀取 uint64
func (c *DataConverter) ReadUint64(buffer []byte, index int) uint64 {
	return c.transform.TransformUint64(buffer, index)
}

// ReadFloat32 從位元組讀取 float32
func (c *DataConverter) ReadFloat32(buffer []byte, index int) float32 {
	return c.transform.TransformFloat32(buffer, index)
}

// ReadFloat64 從位元組讀取 float64
func (c *DataConverter) ReadFloat64(buffer []byte, index int) float64 {
	return c.transform.TransformFloat64(buffer, index)
}

// ReadString 從位元組讀取 ASCII 字串
func (c *DataConverter) ReadString(buffer []byte, index, length int) string {
	if len(buffer) < index+length {
		length = len(buffer) - index
	}
	if length <= 0 {
		return ""
	}
	// 移除尾部的空字元
	data := buffer[index : index+length]
	for i := len(data) - 1; i >= 0; i-- {
		if data[i] != 0 {
			return string(data[:i+1])
		}
	}
	return ""
}

// =============================================================================
// 讀取指定數據類型 (動態)
// =============================================================================

// ReadValue 根據數據類型讀取值
func (c *DataConverter) ReadValue(buffer []byte, index int, dataType DataType) interface{} {
	switch dataType {
	case DataTypeBool:
		return c.ReadBool(buffer, index)
	case DataTypeInt8:
		return c.ReadInt8(buffer, index)
	case DataTypeUint8:
		return c.ReadUint8(buffer, index)
	case DataTypeInt16:
		return c.ReadInt16(buffer, index)
	case DataTypeUint16:
		return c.ReadUint16(buffer, index)
	case DataTypeInt32:
		return c.ReadInt32(buffer, index)
	case DataTypeUint32:
		return c.ReadUint32(buffer, index)
	case DataTypeInt64:
		return c.ReadInt64(buffer, index)
	case DataTypeUint64:
		return c.ReadUint64(buffer, index)
	case DataTypeFloat32:
		return c.ReadFloat32(buffer, index)
	case DataTypeFloat64:
		return c.ReadFloat64(buffer, index)
	default:
		return nil
	}
}

// =============================================================================
// 寫入指定類型到 []byte
// =============================================================================

// WriteBool 寫入布林值
func (c *DataConverter) WriteBool(value bool) []byte {
	if value {
		return []byte{1}
	}
	return []byte{0}
}

// WriteInt8 寫入 int8
func (c *DataConverter) WriteInt8(value int8) []byte {
	return []byte{byte(value)}
}

// WriteUint8 寫入 uint8
func (c *DataConverter) WriteUint8(value uint8) []byte {
	return []byte{value}
}

// WriteInt16 寫入 int16
func (c *DataConverter) WriteInt16(value int16) []byte {
	return c.transform.Int16ToBytes(value)
}

// WriteUint16 寫入 uint16
func (c *DataConverter) WriteUint16(value uint16) []byte {
	return c.transform.Uint16ToBytes(value)
}

// WriteInt32 寫入 int32
func (c *DataConverter) WriteInt32(value int32) []byte {
	return c.transform.Int32ToBytes(value)
}

// WriteUint32 寫入 uint32
func (c *DataConverter) WriteUint32(value uint32) []byte {
	return c.transform.Uint32ToBytes(value)
}

// WriteFloat32 寫入 float32
func (c *DataConverter) WriteFloat32(value float32) []byte {
	return c.transform.Float32ToBytes(value)
}

// WriteInt64 寫入 int64
func (c *DataConverter) WriteInt64(value int64) []byte {
	b := make([]byte, 8)
	switch c.transform.Format {
	case DataFormatABCD:
		binary.BigEndian.PutUint64(b, uint64(value))
	case DataFormatDCBA:
		binary.LittleEndian.PutUint64(b, uint64(value))
	case DataFormatBADC:
		bits := uint64(value)
		b[1] = byte(bits >> 56)
		b[0] = byte(bits >> 48)
		b[3] = byte(bits >> 40)
		b[2] = byte(bits >> 32)
		b[5] = byte(bits >> 24)
		b[4] = byte(bits >> 16)
		b[7] = byte(bits >> 8)
		b[6] = byte(bits)
	case DataFormatCDAB:
		bits := uint64(value)
		b[6] = byte(bits >> 56)
		b[7] = byte(bits >> 48)
		b[4] = byte(bits >> 40)
		b[5] = byte(bits >> 32)
		b[2] = byte(bits >> 24)
		b[3] = byte(bits >> 16)
		b[0] = byte(bits >> 8)
		b[1] = byte(bits)
	default:
		binary.BigEndian.PutUint64(b, uint64(value))
	}
	return b
}

// WriteUint64 寫入 uint64
func (c *DataConverter) WriteUint64(value uint64) []byte {
	return c.WriteInt64(int64(value))
}

// WriteFloat64 寫入 float64
func (c *DataConverter) WriteFloat64(value float64) []byte {
	return c.WriteInt64(int64(math.Float64bits(value)))
}

// WriteString 寫入 ASCII 字串 (補齊到指定長度)
func (c *DataConverter) WriteString(value string, length int) []byte {
	b := make([]byte, length)
	copy(b, []byte(value))
	return b
}

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
		return int16(registers[0])
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
			return uint64(registersToInt64(registers, c.transform.Format))
		}
	case DataTypeFloat64:
		if len(registers) >= 4 {
			bits := uint64(registersToInt64(registers, c.transform.Format))
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
			uint64(registers[1])<<16 | uint64(registers[0]))
	default: // ABCD, BADC
		return int64(uint64(registers[0])<<48 | uint64(registers[1])<<32 |
			uint64(registers[2])<<16 | uint64(registers[3]))
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

// ToFloat64 將任意數值轉換為 float64
func ToFloat64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case bool:
		if v {
			return 1, true
		}
		return 0, true
	case int8:
		return float64(v), true
	case uint8:
		return float64(v), true
	case int16:
		return float64(v), true
	case uint16:
		return float64(v), true
	case int32:
		return float64(v), true
	case uint32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint64:
		return float64(v), true
	case int:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
	}
}

// ToInt64 將任意數值轉換為 int64
func ToInt64(value interface{}) (int64, bool) {
	switch v := value.(type) {
	case bool:
		if v {
			return 1, true
		}
		return 0, true
	case int8:
		return int64(v), true
	case uint8:
		return int64(v), true
	case int16:
		return int64(v), true
	case uint16:
		return int64(v), true
	case int32:
		return int64(v), true
	case uint32:
		return int64(v), true
	case int64:
		return v, true
	case uint64:
		return int64(v), true
	case int:
		return int64(v), true
	case float32:
		return int64(v), true
	case float64:
		return int64(v), true
	default:
		return 0, false
	}
}
