package hsllogic

import (
	"math"
	"testing"
)

// =============================================================================
// ByteTransform 測試 (目標覆蓋率 90%+)
// =============================================================================

func TestNewByteTransform(t *testing.T) {
	bt := NewByteTransform(DataFormatABCD)
	if bt.Format != DataFormatABCD {
		t.Errorf("expected ABCD, got %s", bt.Format)
	}

	bt2 := NewByteTransform(DataFormatCDAB)
	if bt2.Format != DataFormatCDAB {
		t.Errorf("expected CDAB, got %s", bt2.Format)
	}
}

func TestDefaultTransform(t *testing.T) {
	if DefaultTransform.Format != DataFormatABCD {
		t.Errorf("default transform should be ABCD")
	}
}

// =============================================================================
// Uint16 轉換測試
// =============================================================================

func TestTransformUint16_ABCD(t *testing.T) {
	bt := NewByteTransform(DataFormatABCD)
	buffer := []byte{0x12, 0x34}
	result := bt.TransformUint16(buffer, 0)
	expected := uint16(0x1234)
	if result != expected {
		t.Errorf("ABCD Uint16: expected 0x%04X, got 0x%04X", expected, result)
	}
}

func TestTransformUint16_DCBA(t *testing.T) {
	bt := NewByteTransform(DataFormatDCBA)
	buffer := []byte{0x34, 0x12}
	result := bt.TransformUint16(buffer, 0)
	expected := uint16(0x1234)
	if result != expected {
		t.Errorf("DCBA Uint16: expected 0x%04X, got 0x%04X", expected, result)
	}
}

func TestTransformUint16_BufferTooShort(t *testing.T) {
	bt := NewByteTransform(DataFormatABCD)
	buffer := []byte{0x12}
	result := bt.TransformUint16(buffer, 0)
	if result != 0 {
		t.Errorf("buffer too short should return 0")
	}
}

func TestTransformInt16(t *testing.T) {
	bt := NewByteTransform(DataFormatABCD)
	buffer := []byte{0xFF, 0xFE} // -2 in big-endian
	result := bt.TransformInt16(buffer, 0)
	expected := int16(-2)
	if result != expected {
		t.Errorf("Int16: expected %d, got %d", expected, result)
	}
}

// =============================================================================
// Uint32 轉換測試
// =============================================================================

func TestTransformUint32_ABCD(t *testing.T) {
	bt := NewByteTransform(DataFormatABCD)
	buffer := []byte{0x12, 0x34, 0x56, 0x78}
	result := bt.TransformUint32(buffer, 0)
	expected := uint32(0x12345678)
	if result != expected {
		t.Errorf("ABCD Uint32: expected 0x%08X, got 0x%08X", expected, result)
	}
}

func TestTransformUint32_DCBA(t *testing.T) {
	bt := NewByteTransform(DataFormatDCBA)
	buffer := []byte{0x78, 0x56, 0x34, 0x12}
	result := bt.TransformUint32(buffer, 0)
	expected := uint32(0x12345678)
	if result != expected {
		t.Errorf("DCBA Uint32: expected 0x%08X, got 0x%08X", expected, result)
	}
}

func TestTransformUint32_BADC(t *testing.T) {
	bt := NewByteTransform(DataFormatBADC)
	buffer := []byte{0x34, 0x12, 0x78, 0x56}
	result := bt.TransformUint32(buffer, 0)
	expected := uint32(0x12345678)
	if result != expected {
		t.Errorf("BADC Uint32: expected 0x%08X, got 0x%08X", expected, result)
	}
}

func TestTransformUint32_CDAB(t *testing.T) {
	bt := NewByteTransform(DataFormatCDAB)
	buffer := []byte{0x56, 0x78, 0x12, 0x34}
	result := bt.TransformUint32(buffer, 0)
	expected := uint32(0x12345678)
	if result != expected {
		t.Errorf("CDAB Uint32: expected 0x%08X, got 0x%08X", expected, result)
	}
}

func TestTransformUint32_BufferTooShort(t *testing.T) {
	bt := NewByteTransform(DataFormatABCD)
	buffer := []byte{0x12, 0x34}
	result := bt.TransformUint32(buffer, 0)
	if result != 0 {
		t.Errorf("buffer too short should return 0")
	}
}

func TestTransformInt32(t *testing.T) {
	bt := NewByteTransform(DataFormatABCD)
	buffer := []byte{0xFF, 0xFF, 0xFF, 0xFE} // -2 in big-endian
	result := bt.TransformInt32(buffer, 0)
	expected := int32(-2)
	if result != expected {
		t.Errorf("Int32: expected %d, got %d", expected, result)
	}
}

// =============================================================================
// Float32 轉換測試
// =============================================================================

func TestTransformFloat32_ABCD(t *testing.T) {
	bt := NewByteTransform(DataFormatABCD)
	// 3.14 in IEEE 754 big-endian: 0x4048F5C3
	buffer := []byte{0x40, 0x48, 0xF5, 0xC3}
	result := bt.TransformFloat32(buffer, 0)
	expected := float32(3.14)
	if math.Abs(float64(result-expected)) > 0.001 {
		t.Errorf("Float32 ABCD: expected %f, got %f", expected, result)
	}
}

func TestTransformFloat32_CDAB(t *testing.T) {
	bt := NewByteTransform(DataFormatCDAB)
	// 3.14 in CDAB format
	buffer := []byte{0xF5, 0xC3, 0x40, 0x48}
	result := bt.TransformFloat32(buffer, 0)
	expected := float32(3.14)
	if math.Abs(float64(result-expected)) > 0.001 {
		t.Errorf("Float32 CDAB: expected %f, got %f", expected, result)
	}
}

// =============================================================================
// Uint64 轉換測試
// =============================================================================

func TestTransformUint64_ABCD(t *testing.T) {
	bt := NewByteTransform(DataFormatABCD)
	buffer := []byte{0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE, 0xF0}
	result := bt.TransformUint64(buffer, 0)
	expected := uint64(0x123456789ABCDEF0)
	if result != expected {
		t.Errorf("ABCD Uint64: expected 0x%016X, got 0x%016X", expected, result)
	}
}

func TestTransformUint64_DCBA(t *testing.T) {
	bt := NewByteTransform(DataFormatDCBA)
	buffer := []byte{0xF0, 0xDE, 0xBC, 0x9A, 0x78, 0x56, 0x34, 0x12}
	result := bt.TransformUint64(buffer, 0)
	expected := uint64(0x123456789ABCDEF0)
	if result != expected {
		t.Errorf("DCBA Uint64: expected 0x%016X, got 0x%016X", expected, result)
	}
}

func TestTransformUint64_BufferTooShort(t *testing.T) {
	bt := NewByteTransform(DataFormatABCD)
	buffer := []byte{0x12, 0x34, 0x56, 0x78}
	result := bt.TransformUint64(buffer, 0)
	if result != 0 {
		t.Errorf("buffer too short should return 0")
	}
}

func TestTransformInt64(t *testing.T) {
	bt := NewByteTransform(DataFormatABCD)
	buffer := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE}
	result := bt.TransformInt64(buffer, 0)
	expected := int64(-2)
	if result != expected {
		t.Errorf("Int64: expected %d, got %d", expected, result)
	}
}

func TestTransformFloat64(t *testing.T) {
	bt := NewByteTransform(DataFormatABCD)
	// 3.14159265358979 in IEEE 754 big-endian
	buffer := []byte{0x40, 0x09, 0x21, 0xFB, 0x54, 0x44, 0x2D, 0x18}
	result := bt.TransformFloat64(buffer, 0)
	expected := 3.14159265358979
	if math.Abs(result-expected) > 0.0001 {
		t.Errorf("Float64: expected %f, got %f", expected, result)
	}
}

// =============================================================================
// ToBytes 轉換測試
// =============================================================================

func TestUint16ToBytes_ABCD(t *testing.T) {
	bt := NewByteTransform(DataFormatABCD)
	result := bt.Uint16ToBytes(0x1234)
	expected := []byte{0x12, 0x34}
	if result[0] != expected[0] || result[1] != expected[1] {
		t.Errorf("Uint16ToBytes ABCD: expected %v, got %v", expected, result)
	}
}

func TestUint16ToBytes_DCBA(t *testing.T) {
	bt := NewByteTransform(DataFormatDCBA)
	result := bt.Uint16ToBytes(0x1234)
	expected := []byte{0x34, 0x12}
	if result[0] != expected[0] || result[1] != expected[1] {
		t.Errorf("Uint16ToBytes DCBA: expected %v, got %v", expected, result)
	}
}

func TestInt16ToBytes(t *testing.T) {
	bt := NewByteTransform(DataFormatABCD)
	result := bt.Int16ToBytes(-2)
	expected := []byte{0xFF, 0xFE}
	if result[0] != expected[0] || result[1] != expected[1] {
		t.Errorf("Int16ToBytes: expected %v, got %v", expected, result)
	}
}

func TestUint32ToBytes_AllFormats(t *testing.T) {
	tests := []struct {
		format   DataFormat
		value    uint32
		expected []byte
	}{
		{DataFormatABCD, 0x12345678, []byte{0x12, 0x34, 0x56, 0x78}},
		{DataFormatDCBA, 0x12345678, []byte{0x78, 0x56, 0x34, 0x12}},
		{DataFormatBADC, 0x12345678, []byte{0x34, 0x12, 0x78, 0x56}},
		{DataFormatCDAB, 0x12345678, []byte{0x56, 0x78, 0x12, 0x34}},
	}

	for _, tt := range tests {
		bt := NewByteTransform(tt.format)
		result := bt.Uint32ToBytes(tt.value)
		for i, b := range tt.expected {
			if result[i] != b {
				t.Errorf("%s Uint32ToBytes[%d]: expected 0x%02X, got 0x%02X", tt.format, i, b, result[i])
			}
		}
	}
}

func TestFloat32ToBytes(t *testing.T) {
	bt := NewByteTransform(DataFormatABCD)
	result := bt.Float32ToBytes(3.14)
	// 驗證可以轉回來
	back := bt.TransformFloat32(result, 0)
	if math.Abs(float64(back-3.14)) > 0.001 {
		t.Errorf("Float32ToBytes roundtrip failed: got %f", back)
	}
}

// =============================================================================
// Registers 轉換測試
// =============================================================================

func TestRegistersToInt32_ABCD(t *testing.T) {
	bt := NewByteTransform(DataFormatABCD)
	registers := []uint16{0x1234, 0x5678}
	result := bt.RegistersToInt32(registers)
	expected := int32(0x12345678)
	if result != expected {
		t.Errorf("RegistersToInt32 ABCD: expected 0x%08X, got 0x%08X", expected, result)
	}
}

func TestRegistersToInt32_CDAB(t *testing.T) {
	bt := NewByteTransform(DataFormatCDAB)
	registers := []uint16{0x5678, 0x1234}
	result := bt.RegistersToInt32(registers)
	expected := int32(0x12345678)
	if result != expected {
		t.Errorf("RegistersToInt32 CDAB: expected 0x%08X, got 0x%08X", expected, result)
	}
}

func TestRegistersToInt32_TooShort(t *testing.T) {
	bt := NewByteTransform(DataFormatABCD)
	registers := []uint16{0x1234}
	result := bt.RegistersToInt32(registers)
	if result != 0 {
		t.Errorf("registers too short should return 0")
	}
}

func TestRegistersToUint32(t *testing.T) {
	bt := NewByteTransform(DataFormatABCD)
	registers := []uint16{0x1234, 0x5678}
	result := bt.RegistersToUint32(registers)
	expected := uint32(0x12345678)
	if result != expected {
		t.Errorf("RegistersToUint32: expected 0x%08X, got 0x%08X", expected, result)
	}
}

func TestRegistersToFloat32(t *testing.T) {
	bt := NewByteTransform(DataFormatABCD)
	// 3.14 = 0x4048F5C3
	registers := []uint16{0x4048, 0xF5C3}
	result := bt.RegistersToFloat32(registers)
	expected := float32(3.14)
	if math.Abs(float64(result-expected)) > 0.001 {
		t.Errorf("RegistersToFloat32: expected %f, got %f", expected, result)
	}
}
