package hsllogic

import (
	"math"
	"testing"
)

const (
	plcFloat32Bits  uint32 = 0x3F800000
	plcUint32Bits   uint32 = 0x01020304
	plcSigned32Bits uint32 = 0xFEDCBA98
	plcFloat64Bits  uint64 = 0x3FF0000000000000
	plcUint64Bits   uint64 = 0x0102030405060708
	plcSigned64Bits uint64 = 0xFEDCBA9876543210
)

type plcByteOrderVector struct {
	format   DataFormat
	float32  []uint16
	uint32   []uint16
	signed32 []uint16
	float64  []uint16
	uint64   []uint16
	signed64 []uint16
}

func plcByteOrderVectors() []plcByteOrderVector {
	return []plcByteOrderVector{
		{
			format:   DataFormatABCD,
			float32:  []uint16{0x3F80, 0x0000},
			uint32:   []uint16{0x0102, 0x0304},
			signed32: []uint16{0xFEDC, 0xBA98},
			float64:  []uint16{0x3FF0, 0x0000, 0x0000, 0x0000},
			uint64:   []uint16{0x0102, 0x0304, 0x0506, 0x0708},
			signed64: []uint16{0xFEDC, 0xBA98, 0x7654, 0x3210},
		},
		{
			format:   DataFormatBADC,
			float32:  []uint16{0x803F, 0x0000},
			uint32:   []uint16{0x0201, 0x0403},
			signed32: []uint16{0xDCFE, 0x98BA},
			float64:  []uint16{0xF03F, 0x0000, 0x0000, 0x0000},
			uint64:   []uint16{0x0201, 0x0403, 0x0605, 0x0807},
			signed64: []uint16{0xDCFE, 0x98BA, 0x5476, 0x1032},
		},
		{
			format:   DataFormatCDAB,
			float32:  []uint16{0x0000, 0x3F80},
			uint32:   []uint16{0x0304, 0x0102},
			signed32: []uint16{0xBA98, 0xFEDC},
			float64:  []uint16{0x0000, 0x0000, 0x0000, 0x3FF0},
			uint64:   []uint16{0x0708, 0x0506, 0x0304, 0x0102},
			signed64: []uint16{0x3210, 0x7654, 0xBA98, 0xFEDC},
		},
		{
			format:   DataFormatDCBA,
			float32:  []uint16{0x0000, 0x803F},
			uint32:   []uint16{0x0403, 0x0201},
			signed32: []uint16{0x98BA, 0xDCFE},
			float64:  []uint16{0x0000, 0x0000, 0x0000, 0xF03F},
			uint64:   []uint16{0x0807, 0x0605, 0x0403, 0x0201},
			signed64: []uint16{0x1032, 0x5476, 0x98BA, 0xDCFE},
		},
	}
}

func TestPLCByteOrderRegisterDecoding(t *testing.T) {
	for _, vector := range plcByteOrderVectors() {
		t.Run(string(vector.format), func(t *testing.T) {
			transform := NewByteTransform(vector.format)
			converter := NewDataConverter(vector.format)

			t.Run("uint32", func(t *testing.T) {
				if got := transform.RegistersToUint32(vector.uint32); got != plcUint32Bits {
					t.Errorf("got 0x%08X, want 0x%08X", got, plcUint32Bits)
				}
			})
			t.Run("int32", func(t *testing.T) {
				got := transform.RegistersToInt32(vector.signed32)
				if got >= 0 || uint32(got) != plcSigned32Bits {
					t.Errorf("got %d (0x%08X), want negative 0x%08X", got, uint32(got), plcSigned32Bits)
				}
			})
			t.Run("float32", func(t *testing.T) {
				got := transform.RegistersToFloat32(vector.float32)
				if got != math.Float32frombits(plcFloat32Bits) {
					t.Errorf("got %v, want %v", got, math.Float32frombits(plcFloat32Bits))
				}
			})
			t.Run("uint64", func(t *testing.T) {
				got, ok := converter.RegistersToValue(vector.uint64, DataTypeUint64).(uint64)
				if !ok || got != plcUint64Bits {
					t.Errorf("got %#v, want 0x%016X", got, plcUint64Bits)
				}
			})
			t.Run("int64", func(t *testing.T) {
				got, ok := converter.RegistersToValue(vector.signed64, DataTypeInt64).(int64)
				if !ok || got >= 0 || uint64(got) != plcSigned64Bits {
					t.Errorf("got %#v (0x%016X), want negative 0x%016X", got, uint64(got), plcSigned64Bits)
				}
			})
			t.Run("float64", func(t *testing.T) {
				got, ok := converter.RegistersToValue(vector.float64, DataTypeFloat64).(float64)
				if !ok || got != math.Float64frombits(plcFloat64Bits) {
					t.Errorf("got %#v, want %v", got, math.Float64frombits(plcFloat64Bits))
				}
			})
		})
	}
}

func TestPLCByteOrderSingleWordRegisterDecodingUnchanged(t *testing.T) {
	for _, format := range []DataFormat{
		DataFormatABCD,
		DataFormatBADC,
		DataFormatCDAB,
		DataFormatDCBA,
	} {
		t.Run(string(format), func(t *testing.T) {
			converter := NewDataConverter(format)
			int16Value, ok := converter.RegistersToValue([]uint16{0xFFFE}, DataTypeInt16).(int16)
			if !ok || int16Value != -2 {
				t.Fatalf("int16: got %#v, want -2", int16Value)
			}
			uint16Value, ok := converter.RegistersToValue([]uint16{0xFFFE}, DataTypeUint16).(uint16)
			if !ok || uint16Value != 0xFFFE {
				t.Fatalf("uint16: got %#v, want 0xFFFE", uint16Value)
			}
		})
	}
}
