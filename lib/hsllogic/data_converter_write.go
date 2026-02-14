package hsllogic

import (
	"encoding/binary"
	"math"
)

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
