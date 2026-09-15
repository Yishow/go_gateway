package hsllogic

// =============================================================================
// DataConverter 數據轉換器
// =============================================================================

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
	return int8(buffer[index]) // #nosec G115 -- preserve the byte's two's-complement bit pattern as int8.
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
