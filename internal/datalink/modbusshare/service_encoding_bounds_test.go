package modbusshare

import (
	"go-gateway/internal/datalink/schema"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeWordsRejectsValuesOutsideDeclaredWidth(t *testing.T) {
	for _, tc := range []struct {
		kind  schema.DataType
		value any
	}{
		{schema.DataTypeInt16, 32768}, {schema.DataTypeInt16, -32769},
		{schema.DataTypeUint16, 65536}, {schema.DataTypeUint16, -1},
		{schema.DataTypeInt32, int64(math.MaxInt32) + 1}, {schema.DataTypeInt32, int64(math.MinInt32) - 1},
		{schema.DataTypeUint32, uint64(math.MaxUint32) + 1},
		{schema.DataTypeInt64, float64(math.MaxInt64)}, {schema.DataTypeInt64, math.NaN()},
		{schema.DataTypeUint64, float64(math.MaxUint64)}, {schema.DataTypeUint64, math.NaN()},
	} {
		words, err := encodeToWords(tc.kind, tc.value)
		require.Error(t, err, "type=%s value=%v", tc.kind, tc.value)
		require.Nil(t, words)
	}
}

func TestEncodeWordsPreservesSignAndIEEEBits(t *testing.T) {
	for _, tc := range []struct {
		kind  schema.DataType
		value any
		words []uint16
	}{
		{schema.DataTypeInt16, -1, []uint16{65535}},
		{schema.DataTypeInt32, int32(math.MinInt32), []uint16{32768, 0}},
		{schema.DataTypeInt64, int64(math.MinInt64), []uint16{32768, 0, 0, 0}},
		{schema.DataTypeUint64, uint64(math.MaxUint64), []uint16{65535, 65535, 65535, 65535}},
		{schema.DataTypeFloat32, float32(-1), []uint16{0xbf80, 0}},
		{schema.DataTypeFloat64, float64(-1), []uint16{0xbff0, 0, 0, 0}},
	} {
		words, err := encodeToWords(tc.kind, tc.value)
		require.NoError(t, err)
		require.Equal(t, tc.words, words)
	}
}
