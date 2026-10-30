package decoding

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	"go-gateway/lib/hsllogic"
	"testing"
)

func TestSingleWordPoliciesKeepRegisterAndDirectByteSemanticsSeparate(t *testing.T) {
	for _, protocol := range []schema.ProtocolType{schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP, schema.ProtocolMC3E, schema.ProtocolFatekFBs} {
		for _, mode := range []schema.PointFormatMode{schema.PointFormatLegacy, schema.PointFormatInherit} {
			for _, format := range []string{"BADC", "DCBA"} {
				for _, dtype := range []schema.DataType{schema.DataTypeUint16, schema.DataTypeInt16} {
					p := schema.Point{Address: "D0", DataType: dtype, FormatMode: mode, DataFormat: format}
					before := p
					got, err := ResolvePoint(p, FormatContext{Protocol: protocol, RuleOwned: true, RuleFormat: format, ConnectionFormat: format})
					require.NoError(t, err)
					require.False(t, got.Applicable)
					require.Equal(t, before, p)
					if protocol == schema.ProtocolFatekFBs {
						require.True(t, got.NativeCompatibility)
						require.Empty(t, got.EffectiveFormat)
						continue
					}
					require.Equal(t, format, got.EffectiveFormat)
					converter := hsllogic.NewDataConverter(hsllogic.DataFormat(got.EffectiveFormat))
					if dtype == schema.DataTypeUint16 {
						require.Equal(t, uint16(0x0102), converter.RegistersToValue([]uint16{0x0102}, hsllogic.DataTypeUint16))
					} else {
						require.Equal(t, int16(-32767), converter.RegistersToValue([]uint16{0x8001}, hsllogic.DataTypeInt16))
					}
					require.Equal(t, uint16(0x0201), hsllogic.NewByteTransform(hsllogic.DataFormat(format)).TransformUint16([]byte{1, 2}, 0))
				}
			}
		}
	}
}
