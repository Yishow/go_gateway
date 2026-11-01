package optimizer

import (
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionPointSpanUsesCanonicalPhysicalWidth(t *testing.T) {
	for _, tc := range []struct {
		protocol, address string
		dataType          schema.DataType
		want              NativeRange
	}{
		{"modbus", "40001", schema.DataTypeFloat64, NativeRange{"HR", 0, 4}},
		{"mc", "W0F", schema.DataTypeUint32, NativeRange{"W", 15, 17}},
		{"fatek", "WY0", schema.DataTypeUint16, NativeRange{"WY", 0, 16}},
		{"fatek", "DR0", schema.DataTypeFloat64, NativeRange{"DR", 0, 4}},
	} {
		t.Run(tc.protocol+tc.address, func(t *testing.T) {
			got, err := AcquisitionPointSpan(tc.protocol, EffectivePoint{Address: tc.address, DataType: tc.dataType})
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
	for _, address := range []string{"D999999999999999999999999", "D16777215"} {
		_, err := AcquisitionPointSpan("mc", EffectivePoint{Address: address, DataType: schema.DataTypeUint32})
		require.Error(t, err)
	}
}
