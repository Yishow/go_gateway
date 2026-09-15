// Package hsllogic 提供仿照 HslCommunication 的工業級數據轉換與地址解析功能。
// 此套件為 go_gateway 的核心共通邏輯層，所有協議適配器都應整合此套件。
package hsllogic

import (
	"encoding/binary"
	"math"
)

// =============================================================================
// DataFormat 字節序格式定義
// =============================================================================

// DataFormat 定義字節序格式，仿照 HSL 的 IByteTransform
type DataFormat string

const (
	// DataFormatABCD Big-Endian (預設，Modbus 標準)
	DataFormatABCD DataFormat = "ABCD"
	// DataFormatBADC Big-Endian with Byte Swap
	DataFormatBADC DataFormat = "BADC"
	// DataFormatCDAB Little-Endian with Word Swap (許多 PLC 如 Siemens 使用)
	DataFormatCDAB DataFormat = "CDAB"
	// DataFormatDCBA Little-Endian (Intel 格式)
	DataFormatDCBA DataFormat = "DCBA"
)

// =============================================================================
// ByteTransform 字節轉換器
// =============================================================================

// ByteTransform 提供字節序轉換的核心功能
type ByteTransform struct {
	Format DataFormat
}

// NewByteTransform 建立新的字節轉換器
func NewByteTransform(format DataFormat) *ByteTransform {
	return &ByteTransform{Format: format}
}

// DefaultTransform 使用預設格式 (ABCD) 的轉換器
var DefaultTransform = NewByteTransform(DataFormatABCD)

// =============================================================================
// 從 []byte 轉換為各類型
// =============================================================================

// TransformUint16 將 2 字節轉換為 uint16
func (t *ByteTransform) TransformUint16(buffer []byte, index int) uint16 {
	if len(buffer) < index+2 {
		return 0
	}
	b := buffer[index : index+2]

	switch t.Format {
	case DataFormatBADC, DataFormatDCBA:
		// 交換字節
		return uint16(b[1])<<8 | uint16(b[0])
	default: // ABCD, CDAB
		return uint16(b[0])<<8 | uint16(b[1])
	}
}

// TransformInt16 將 2 字節轉換為 int16
func (t *ByteTransform) TransformInt16(buffer []byte, index int) int16 {
	return int16(t.TransformUint16(buffer, index)) // #nosec G115 -- fixed-width bit reinterpretation preserves signed 16-bit wire data.
}

// TransformUint32 將 4 字節轉換為 uint32
func (t *ByteTransform) TransformUint32(buffer []byte, index int) uint32 {
	if len(buffer) < index+4 {
		return 0
	}
	b := buffer[index : index+4]

	switch t.Format {
	case DataFormatABCD:
		return binary.BigEndian.Uint32(b)
	case DataFormatDCBA:
		return binary.LittleEndian.Uint32(b)
	case DataFormatBADC:
		// AB CD -> BA DC (每個 word 內部交換)
		return uint32(b[1])<<24 | uint32(b[0])<<16 | uint32(b[3])<<8 | uint32(b[2])
	case DataFormatCDAB:
		// AB CD -> CD AB (word 順序交換)
		return uint32(b[2])<<24 | uint32(b[3])<<16 | uint32(b[0])<<8 | uint32(b[1])
	default:
		return binary.BigEndian.Uint32(b)
	}
}

// TransformInt32 將 4 字節轉換為 int32
func (t *ByteTransform) TransformInt32(buffer []byte, index int) int32 {
	return int32(t.TransformUint32(buffer, index)) // #nosec G115 -- fixed-width bit reinterpretation preserves signed 32-bit wire data.
}

// TransformFloat32 將 4 字節轉換為 float32
func (t *ByteTransform) TransformFloat32(buffer []byte, index int) float32 {
	bits := t.TransformUint32(buffer, index)
	return math.Float32frombits(bits)
}

// TransformUint64 將 8 字節轉換為 uint64
func (t *ByteTransform) TransformUint64(buffer []byte, index int) uint64 {
	if len(buffer) < index+8 {
		return 0
	}
	b := buffer[index : index+8]

	switch t.Format {
	case DataFormatABCD:
		return binary.BigEndian.Uint64(b)
	case DataFormatDCBA:
		return binary.LittleEndian.Uint64(b)
	case DataFormatBADC:
		return uint64(b[1])<<56 | uint64(b[0])<<48 | uint64(b[3])<<40 | uint64(b[2])<<32 |
			uint64(b[5])<<24 | uint64(b[4])<<16 | uint64(b[7])<<8 | uint64(b[6])
	case DataFormatCDAB:
		return uint64(b[6])<<56 | uint64(b[7])<<48 | uint64(b[4])<<40 | uint64(b[5])<<32 |
			uint64(b[2])<<24 | uint64(b[3])<<16 | uint64(b[0])<<8 | uint64(b[1])
	default:
		return binary.BigEndian.Uint64(b)
	}
}

// TransformInt64 將 8 字節轉換為 int64
func (t *ByteTransform) TransformInt64(buffer []byte, index int) int64 {
	return int64(t.TransformUint64(buffer, index)) // #nosec G115 -- fixed-width bit reinterpretation preserves signed 64-bit wire data.
}

// TransformFloat64 將 8 字節轉換為 float64
func (t *ByteTransform) TransformFloat64(buffer []byte, index int) float64 {
	bits := t.TransformUint64(buffer, index)
	return math.Float64frombits(bits)
}

// =============================================================================
// 從各類型轉換為 []byte
// =============================================================================

// Uint16ToBytes 將 uint16 轉換為 2 字節
func (t *ByteTransform) Uint16ToBytes(value uint16) []byte {
	b := make([]byte, 2)
	switch t.Format {
	case DataFormatBADC, DataFormatDCBA:
		b[0] = byte(value)      // #nosec G115 -- extract the low byte of a fixed-width uint16 wire value.
		b[1] = byte(value >> 8) // #nosec G115 -- extract the high byte of a fixed-width uint16 wire value.
	default:
		b[0] = byte(value >> 8) // #nosec G115 -- extract the high byte of a fixed-width uint16 wire value.
		b[1] = byte(value)      // #nosec G115 -- extract the low byte of a fixed-width uint16 wire value.
	}
	return b
}

// Int16ToBytes 將 int16 轉換為 2 字節
func (t *ByteTransform) Int16ToBytes(value int16) []byte {
	return t.Uint16ToBytes(uint16(value)) // #nosec G115 -- preserve the two's-complement 16-bit representation.
}

// Uint32ToBytes 將 uint32 轉換為 4 字節
func (t *ByteTransform) Uint32ToBytes(value uint32) []byte {
	b := make([]byte, 4)
	switch t.Format {
	case DataFormatABCD:
		binary.BigEndian.PutUint32(b, value)
	case DataFormatDCBA:
		binary.LittleEndian.PutUint32(b, value)
	case DataFormatBADC:
		b[1] = byte(value >> 24) // #nosec G115 -- extract a byte from a fixed-width uint32 wire value.
		b[0] = byte(value >> 16) // #nosec G115 -- extract a byte from a fixed-width uint32 wire value.
		b[3] = byte(value >> 8)  // #nosec G115 -- extract a byte from a fixed-width uint32 wire value.
		b[2] = byte(value)       // #nosec G115 -- extract a byte from a fixed-width uint32 wire value.
	case DataFormatCDAB:
		b[2] = byte(value >> 24) // #nosec G115 -- extract a byte from a fixed-width uint32 wire value.
		b[3] = byte(value >> 16) // #nosec G115 -- extract a byte from a fixed-width uint32 wire value.
		b[0] = byte(value >> 8)  // #nosec G115 -- extract a byte from a fixed-width uint32 wire value.
		b[1] = byte(value)       // #nosec G115 -- extract a byte from a fixed-width uint32 wire value.
	default:
		binary.BigEndian.PutUint32(b, value)
	}
	return b
}

// Int32ToBytes 將 int32 轉換為 4 字節
func (t *ByteTransform) Int32ToBytes(value int32) []byte {
	return t.Uint32ToBytes(uint32(value)) // #nosec G115 -- preserve the two's-complement 32-bit representation.
}

// Float32ToBytes 將 float32 轉換為 4 字節
func (t *ByteTransform) Float32ToBytes(value float32) []byte {
	bits := math.Float32bits(value)
	return t.Uint32ToBytes(bits)
}

// =============================================================================
// Registers (uint16 陣列) 轉換
// =============================================================================

// RegistersToInt32 將 2 個暫存器轉換為 int32
func (t *ByteTransform) RegistersToInt32(registers []uint16) int32 {
	if len(registers) < 2 {
		return 0
	}
	switch t.Format {
	case DataFormatCDAB, DataFormatDCBA:
		return int32(uint32(registers[1])<<16 | uint32(registers[0])) // #nosec G115 -- reinterpret assembled 32-bit register bits as signed.
	default:
		return int32(uint32(registers[0])<<16 | uint32(registers[1])) // #nosec G115 -- reinterpret assembled 32-bit register bits as signed.
	}
}

// RegistersToUint32 將 2 個暫存器轉換為 uint32
func (t *ByteTransform) RegistersToUint32(registers []uint16) uint32 {
	if len(registers) < 2 {
		return 0
	}
	switch t.Format {
	case DataFormatCDAB, DataFormatDCBA:
		return uint32(registers[1])<<16 | uint32(registers[0])
	default:
		return uint32(registers[0])<<16 | uint32(registers[1])
	}
}

// RegistersToFloat32 將 2 個暫存器轉換為 float32
func (t *ByteTransform) RegistersToFloat32(registers []uint16) float32 {
	bits := t.RegistersToUint32(registers)
	return math.Float32frombits(bits)
}
