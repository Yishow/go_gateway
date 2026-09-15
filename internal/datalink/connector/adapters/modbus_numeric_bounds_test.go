package adapters

import (
	"go-gateway/internal/datalink/schema"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestModbusNumericInputBounds(t *testing.T) {
	for _, value := range []any{int(-1), int16(-1), int32(-1), int64(-1), int(65536), uint(65536), uint32(65536), uint64(math.MaxUint64), float32(1.5), 1.5, math.NaN(), math.Inf(1)} {
		_, err := toUint16(value)
		require.Error(t, err, "value=%v", value)
	}
	for _, value := range []any{uint64(65535), 65535.0, int32(65535), uint16(65535)} {
		actual, err := toUint16(value)
		require.NoError(t, err)
		require.Equal(t, uint16(65535), actual)
	}
	for _, value := range []any{[]int{1, -1}, []interface{}{1, 65536}} {
		values, err := toUint16Slice(value)
		require.Error(t, err)
		require.Nil(t, values)
	}
	for _, address := range []string{"65536", "4294967295", "HR65536", "-1"} {
		_, _, err := parseModbusAddress(address, "03")
		require.Error(t, err, "address=%s", address)
	}
	address, function, err := parseModbusAddress("65535", "03")
	require.NoError(t, err)
	require.Equal(t, uint16(65535), address)
	require.Equal(t, "03", function)
}

func TestFatekSignedWireBitsArePreserved(t *testing.T) {
	require.Equal(t, int16(-1), convertFatekValue([]int{65535}, schema.DataTypeInt16, 16))
	require.Equal(t, int32(-1), convertFatekValue([]int{65535, 65535}, schema.DataTypeInt32, 16))
	require.Equal(t, uint16(65535), convertFatekValue([]int{-1}, schema.DataTypeUint16, 16))
	require.Equal(t, []byte{255, 255, 255, 255, 128, 0, 0, 0}, intSliceToBytes([]int{-1, math.MinInt32}))
}
