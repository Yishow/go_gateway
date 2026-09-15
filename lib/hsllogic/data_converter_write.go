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
	return []byte{byte(value)} // #nosec G115 -- preserve the int8 two's-complement byte representation.
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
		binary.BigEndian.PutUint64(b, uint64(value)) // #nosec G115 -- write the int64 bit pattern as an unsigned fixed-width word.
	case DataFormatDCBA:
		binary.LittleEndian.PutUint64(b, uint64(value)) // #nosec G115 -- write the int64 bit pattern as an unsigned fixed-width word.
	case DataFormatBADC:
		bits := uint64(value)   // #nosec G115 -- preserve the int64 two's-complement bit pattern.
		b[1] = byte(bits >> 56) // #nosec G115 -- extract a byte from a fixed-width uint64 wire value.
		b[0] = byte(bits >> 48) // #nosec G115 -- extract a byte from a fixed-width uint64 wire value.
		b[3] = byte(bits >> 40) // #nosec G115 -- extract a byte from a fixed-width uint64 wire value.
		b[2] = byte(bits >> 32) // #nosec G115 -- extract a byte from a fixed-width uint64 wire value.
		b[5] = byte(bits >> 24) // #nosec G115 -- extract a byte from a fixed-width uint64 wire value.
		b[4] = byte(bits >> 16) // #nosec G115 -- extract a byte from a fixed-width uint64 wire value.
		b[7] = byte(bits >> 8)  // #nosec G115 -- extract a byte from a fixed-width uint64 wire value.
		b[6] = byte(bits)       // #nosec G115 -- extract a byte from a fixed-width uint64 wire value.
	case DataFormatCDAB:
		bits := uint64(value)   // #nosec G115 -- preserve the int64 two's-complement bit pattern.
		b[6] = byte(bits >> 56) // #nosec G115 -- extract a byte from a fixed-width uint64 wire value.
		b[7] = byte(bits >> 48) // #nosec G115 -- extract a byte from a fixed-width uint64 wire value.
		b[4] = byte(bits >> 40) // #nosec G115 -- extract a byte from a fixed-width uint64 wire value.
		b[5] = byte(bits >> 32) // #nosec G115 -- extract a byte from a fixed-width uint64 wire value.
		b[2] = byte(bits >> 24) // #nosec G115 -- extract a byte from a fixed-width uint64 wire value.
		b[3] = byte(bits >> 16) // #nosec G115 -- extract a byte from a fixed-width uint64 wire value.
		b[0] = byte(bits >> 8)  // #nosec G115 -- extract a byte from a fixed-width uint64 wire value.
		b[1] = byte(bits)       // #nosec G115 -- extract a byte from a fixed-width uint64 wire value.
	default:
		binary.BigEndian.PutUint64(b, uint64(value)) // #nosec G115 -- write the int64 bit pattern as an unsigned fixed-width word.
	}
	return b
}

// WriteUint64 寫入 uint64
func (c *DataConverter) WriteUint64(value uint64) []byte {
	return c.WriteInt64(int64(value)) // #nosec G115 -- preserve the uint64 bit pattern through the shared fixed-width writer.
}

// WriteFloat64 寫入 float64
func (c *DataConverter) WriteFloat64(value float64) []byte {
	return c.WriteInt64(int64(math.Float64bits(value))) // #nosec G115 -- reinterpret IEEE-754 bits as the shared int64 wire word.
}

// WriteString 寫入 ASCII 字串 (補齊到指定長度)
func (c *DataConverter) WriteString(value string, length int) []byte {
	b := make([]byte, length)
	copy(b, value)
	return b
}
