package modbusshare

import (
	"math"
	"strconv"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestModbusWordSplitPreservesRawBits(t *testing.T) {
	for _, tc := range []struct {
		raw   uint32
		words []uint16
	}{
		{0, []uint16{0, 0}}, {0xffffffff, []uint16{0xffff, 0xffff}},
		{0x80000000, []uint16{0x8000, 0}}, {0x7fffffff, []uint16{0x7fff, 0xffff}},
		{0x7fc12345, []uint16{0x7fc1, 0x2345}}, {0x7fa12345, []uint16{0x7fa1, 0x2345}},
		{0xffc12345, []uint16{0xffc1, 0x2345}}, {0x12345678, []uint16{0x1234, 0x5678}},
	} {
		require.Equal(t, tc.words, uint32Words(tc.raw))
	}
	for _, tc := range []struct {
		raw   uint64
		words []uint16
	}{
		{0, []uint16{0, 0, 0, 0}}, {0xffffffffffffffff, []uint16{0xffff, 0xffff, 0xffff, 0xffff}},
		{0x8000000000000000, []uint16{0x8000, 0, 0, 0}}, {0x7fffffffffffffff, []uint16{0x7fff, 0xffff, 0xffff, 0xffff}},
		{0x7ff8123456789abc, []uint16{0x7ff8, 0x1234, 0x5678, 0x9abc}}, {0x7ff0123456789abc, []uint16{0x7ff0, 0x1234, 0x5678, 0x9abc}},
		{0xfff8123456789abc, []uint16{0xfff8, 0x1234, 0x5678, 0x9abc}}, {0x123456789abcdef0, []uint16{0x1234, 0x5678, 0x9abc, 0xdef0}},
	} {
		require.Equal(t, tc.words, uint64Words(tc.raw))
	}
}

func TestModbusEncodingKeepsSignedAndFloatBoundaries(t *testing.T) {
	for _, tc := range []struct {
		kind  schema.DataType
		value any
		words []uint16
	}{
		{schema.DataTypeInt32, int32(math.MinInt32), []uint16{0x8000, 0}},
		{schema.DataTypeInt32, int32(-1), []uint16{0xffff, 0xffff}},
		{schema.DataTypeUint32, uint32(math.MaxUint32), []uint16{0xffff, 0xffff}},
		{schema.DataTypeInt64, int64(math.MinInt64), []uint16{0x8000, 0, 0, 0}},
		{schema.DataTypeInt64, int64(-1), []uint16{0xffff, 0xffff, 0xffff, 0xffff}},
		{schema.DataTypeUint64, uint64(math.MaxUint64), []uint16{0xffff, 0xffff, 0xffff, 0xffff}},
		{schema.DataTypeFloat32, math.Copysign(0, -1), []uint16{0x8000, 0}},
		{schema.DataTypeFloat64, math.Copysign(0, -1), []uint16{0x8000, 0, 0, 0}},
		{schema.DataTypeFloat64, math.Float64frombits(0x7ff8123456789abc), []uint16{0x7ff8, 0x1234, 0x5678, 0x9abc}},
	} {
		got, err := encodeToWords(tc.kind, tc.value)
		require.NoError(t, err)
		require.Equal(t, tc.words, got)
	}
	maxUint := ^uint(0)
	got, err := toInt64(maxUint)
	if strconv.IntSize == 64 {
		require.EqualError(t, err, "uint value 18446744073709551615 overflows int64")
		require.Zero(t, got)
	} else {
		require.NoError(t, err)
		require.Equal(t, int64(4294967295), got)
	}
}
