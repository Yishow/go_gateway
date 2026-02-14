package adapters

import (
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFatekAddress(t *testing.T) {
	symbol, address, err := parseFatekAddress(" d100 ")
	require.NoError(t, err)
	assert.Equal(t, "D", symbol)
	assert.Equal(t, 100, address)

	symbol, address, err = parseFatekAddress("dwx12")
	require.NoError(t, err)
	assert.Equal(t, "DWX", symbol)
	assert.Equal(t, 12, address)
}

func TestConvertFatekValue(t *testing.T) {
	value := convertFatekValue([]int{1, 2}, schema.DataTypeUint32, 16)
	assert.Equal(t, uint32(1)<<16|uint32(2), value)

	value = convertFatekValue([]int{123}, schema.DataTypeBool, 16)
	assert.Equal(t, true, value)
}

func TestToIntSlice(t *testing.T) {
	values, err := toIntSlice([]interface{}{1, float64(2)})
	require.NoError(t, err)
	assert.Equal(t, []int{1, 2}, values)

	_, err = toIntSlice("bad")
	require.Error(t, err)
}
