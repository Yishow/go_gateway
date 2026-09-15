package hsllogic

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDataConverterRegistersToValues(t *testing.T) {
	converter := NewDataConverter(DataFormatABCD)
	registers := []uint16{0, 100, 0, 200}

	values := converter.RegistersToValues(registers, DataTypeInt32, 2)
	require.Len(t, values, 2)
	assert.Equal(t, int32(100), values[0])
	assert.Equal(t, int32(200), values[1])
}

func TestToFloat64AndToInt64(t *testing.T) {
	f, ok := ToFloat64(true)
	require.True(t, ok)
	assert.Equal(t, 1.0, f)

	i, ok := ToInt64(float64(12.8))
	require.True(t, ok)
	assert.Equal(t, int64(12), i)

	_, ok = ToInt64("bad")
	assert.False(t, ok)
}

func TestToInt64RejectsOverflowAndNonFiniteValues(t *testing.T) {
	for _, value := range []interface{}{
		uint64(math.MaxInt64) + 1,
		float64(math.MaxInt64),
		float64(math.MaxInt64) * 2,
		float32(1 << 63),
		math.NaN(),
		math.Inf(1),
		math.Inf(-1),
	} {
		got, ok := ToInt64(value)
		assert.Falsef(t, ok, "ToInt64(%v) returned %d", value, got)
		assert.Zero(t, got)
	}
}

func TestToInt64AcceptsSignedRangeBoundaries(t *testing.T) {
	for _, tt := range []struct {
		value interface{}
		want  int64
	}{
		{value: uint64(math.MaxInt64), want: math.MaxInt64},
		{value: float64(math.MinInt64), want: math.MinInt64},
	} {
		got, ok := ToInt64(tt.value)
		assert.Truef(t, ok, "ToInt64(%v) rejected an in-range value", tt.value)
		assert.Equal(t, tt.want, got)
	}
}
