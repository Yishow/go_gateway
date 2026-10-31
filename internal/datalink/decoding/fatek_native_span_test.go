package decoding

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
)

func TestFatekFormatCapabilityUsesPhysicalComponentStep(t *testing.T) {
	for _, tc := range []struct {
		address    string
		dtype      schema.DataType
		components int
		valid      bool
	}{
		{"WY8", schema.DataTypeUint32, 2, true}, {"DWM8", schema.DataTypeUint64, 2, true},
		{"DR1", schema.DataTypeUint64, 2, true}, {"DR65532", schema.DataTypeUint64, 2, true},
		{"DR65533", schema.DataTypeUint64, 2, false}, {"WY9", schema.DataTypeUint32, 2, false},
	} {
		t.Run(tc.address, func(t *testing.T) {
			cap, err := FormatCapability(schema.ProtocolFatekFBs, tc.address, tc.dtype)
			if !tc.valid {
				require.ErrorIs(t, err, ErrInapplicableFormat)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.components, cap.Components)
			for _, mode := range []schema.PointFormatMode{schema.PointFormatInherit, schema.PointFormatOverride} {
				point := schema.Point{Address: tc.address, DataType: tc.dtype, FormatMode: mode, FormatVersion: 1}
				if mode == schema.PointFormatOverride {
					point.FormatOverride = "CDAB"
				}
				resolution, err := ResolvePoint(point, FormatContext{Protocol: schema.ProtocolFatekFBs})
				require.NoError(t, err)
				require.True(t, resolution.Applicable)
			}
		})
	}
}
