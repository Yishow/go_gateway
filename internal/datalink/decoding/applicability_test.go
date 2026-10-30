package decoding

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func TestResolvePointRejectsInapplicableOverride(t *testing.T) {
	for _, tc := range []struct {
		protocol schema.ProtocolType
		address  string
		dtype    schema.DataType
	}{
		{schema.ProtocolModbusTCP, "HR0", schema.DataTypeUint16},
		{schema.ProtocolMC3E, "D0", schema.DataTypeBool},
		{schema.ProtocolFatekFBs, "D0", schema.DataTypeString},
		{schema.ProtocolModbusRTU, "COIL0", schema.DataTypeUint32},
		{schema.ProtocolMC3E, "M0", schema.DataTypeUint32},
		{schema.ProtocolMC3E, "D0.1", schema.DataTypeUint32},
		{schema.ProtocolFatekFBs, "M0", schema.DataTypeUint32},
		{schema.ProtocolModbusTCP, "", schema.DataTypeUint32},
		{schema.ProtocolFatekFBs, "secret://bad", schema.DataTypeUint32},
	} {
		t.Run(string(tc.protocol)+tc.address+string(tc.dtype), func(t *testing.T) {
			p := schema.Point{Address: tc.address, DataType: tc.dtype, FormatMode: schema.PointFormatOverride, FormatOverride: "CDAB", FormatVersion: 1}
			got, err := ResolvePoint(p, FormatContext{Protocol: tc.protocol})
			require.Error(t, err)
			require.Zero(t, got)
			if err != nil {
				require.NotContains(t, err.Error(), "secret")
			}
		})
	}
}

func TestFormatCapabilityProtocolTypeWidthMatrix(t *testing.T) {
	for _, tc := range []struct {
		protocol schema.ProtocolType
		address  string
		width    int
	}{
		{schema.ProtocolModbusTCP, "HR0", 16}, {schema.ProtocolModbusRTU, "IR0", 16}, {schema.ProtocolModbusUDP, "40001", 16},
		{schema.ProtocolMC3E, "D0", 16}, {schema.ProtocolFatekFBs, "D0", 16}, {schema.ProtocolFatekFBs, "DR0", 32},
	} {
		for _, dtype := range []schema.DataType{schema.DataTypeInt32, schema.DataTypeUint32, schema.DataTypeFloat32, schema.DataTypeInt64, schema.DataTypeUint64, schema.DataTypeFloat64} {
			for _, format := range []string{"ABCD", "BADC", "CDAB", "DCBA"} {
				t.Run(string(tc.protocol)+tc.address+string(dtype)+format, func(t *testing.T) {
					cap, err := FormatCapability(tc.protocol, tc.address, dtype)
					require.NoError(t, err)
					require.Equal(t, tc.width, cap.NativeBits)
					require.Equal(t, scalarBits(dtype)/tc.width, cap.Components)
					require.Equal(t, scalarBits(dtype)/16, cap.Words)
					got, err := ResolvePoint(schema.Point{Address: tc.address, DataType: dtype, FormatMode: schema.PointFormatOverride, FormatOverride: format}, FormatContext{Protocol: tc.protocol})
					require.NoError(t, err)
					require.True(t, got.Applicable)
					require.Equal(t, format, got.EffectiveFormat)
				})
			}
		}
	}
}

func TestFormatCapabilityNativeSpanAndUnsupported(t *testing.T) {
	for _, tc := range []struct {
		protocol schema.ProtocolType
		address  string
	}{
		{schema.ProtocolModbusTCP, "HR65535"}, {schema.ProtocolMC3E, "D16777215"},
		{schema.ProtocolFatekFBs, "DR99999"}, {schema.ProtocolMC3E, "UNKNOWN0"},
		{schema.ProtocolModbusTCP, "HR0.0"}, {schema.ProtocolModbusTCP, "DI0"},
		{schema.ProtocolMQTT, "HR0"},
	} {
		got, err := FormatCapability(tc.protocol, tc.address, schema.DataTypeUint64)
		require.Error(t, err, tc.address)
		require.Zero(t, got)
	}
	for _, dtype := range []schema.DataType{schema.DataTypeBool, schema.DataTypeInt16, schema.DataTypeUint16, schema.DataTypeString, "unknown"} {
		p := schema.Point{Address: "HR0", DataType: dtype, DataFormat: "BADC", RevisionID: "original"}
		mode, format := schema.PointFormatOverride, "CDAB"
		got, err := WithIntent(p, IntentEdit{Mode: &mode, ConfiguredFormat: &format})
		require.ErrorIs(t, err, ErrInapplicableFormat)
		require.Equal(t, p, got)
	}
}

func TestInvalidInheritedScalarFailsButLegacyMetadataRemains(t *testing.T) {
	p := schema.Point{Address: "M0", DataType: schema.DataTypeUint64, DataFormat: "BADC"}
	before := p
	got, err := ResolvePoint(p, FormatContext{Protocol: schema.ProtocolMC3E, ConnectionFormat: "CDAB"})
	require.NoError(t, err)
	require.False(t, got.Applicable)
	require.Equal(t, "CDAB", got.EffectiveFormat)
	require.Equal(t, before, p)
	p.FormatMode = schema.PointFormatInherit
	got, err = ResolvePoint(p, FormatContext{Protocol: schema.ProtocolMC3E})
	require.ErrorIs(t, err, ErrInapplicableFormat)
	require.Zero(t, got)
}
