package optimizer

import (
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionPointShapeUsesPlannerAddressAndTypeAuthority(t *testing.T) {
	for _, tc := range []struct {
		protocol, address, function string
		dtype                       schema.DataType
		bits, width                 int
	}{
		{"mc", "D0", "", schema.DataTypeFloat64, 16, 4},
		{"mc", "D0", "", schema.DataTypeBool, 16, 1},
		{"mc", "M0", "", schema.DataTypeBool, 1, 1},
		{acquisitionProtocolFatek, "DR0", "", schema.DataTypeFloat64, 32, 2},
		{acquisitionProtocolFatek, "D0", "", schema.DataTypeFloat32, 16, 2},
		{acquisitionProtocolModbus, "0", "04", schema.DataTypeUint32, 16, 2},
		{acquisitionProtocolModbus, "0", "01", schema.DataTypeBool, 1, 1},
	} {
		point := EffectivePoint{Address: tc.address, Function: tc.function, DataType: tc.dtype}
		bits, width, err := AcquisitionPointShape(tc.protocol, point)
		require.NoError(t, err, "%+v", tc)
		require.Equal(t, tc.bits, bits)
		require.Equal(t, tc.width, width)
	}
	for _, point := range []EffectivePoint{
		{Address: "UNKNOWN0", DataType: schema.DataTypeUint16},
		{Address: "D0", DataType: schema.DataTypeString},
		{Address: "M0", DataType: schema.DataTypeUint16},
		{Address: "D0.3", DataType: schema.DataTypeUint16},
		{Address: "D0.3", DataType: schema.DataTypeBool},
		{Address: "D16777216", DataType: schema.DataTypeUint16},
		{Address: "D16777215", DataType: schema.DataTypeUint32},
	} {
		_, _, err := AcquisitionPointShape("mc", point)
		require.Error(t, err)
	}
}

func TestAcquisitionCoverageValidationUsesCanonicalPhysicalRanges(t *testing.T) {
	for _, field := range []string{"readable", "forbidden", "excluded"} {
		t.Run(field, func(t *testing.T) {
			coverage := ReadableCoverage{Revision: "v"}
			ranges := []NativeRange{{Area: "HR", Start: 65535, End: 65536}}
			switch field {
			case "readable":
				coverage.Readable = ranges
			case "forbidden":
				coverage.Forbidden = ranges
			case "excluded":
				coverage.Excluded = ranges
			}
			require.True(t, ValidateAcquisitionCoverage(acquisitionProtocolModbus, coverage))
			ranges[0].End++
			require.False(t, ValidateAcquisitionCoverage(acquisitionProtocolModbus, coverage))
		})
	}
}
