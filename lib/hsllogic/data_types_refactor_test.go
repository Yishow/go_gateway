package hsllogic

import (
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
