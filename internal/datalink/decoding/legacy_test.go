package decoding

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
)

func TestResolveLegacyPolicies(t *testing.T) {
	cases := []struct {
		name  string
		input LegacyInput
		want  Resolution
	}{
		{"modbus_unknown_point", LegacyInput{Protocol: schema.ProtocolModbusTCP, PointFormat: "invalid", ConnectionFormat: "CDAB"}, Resolution{"legacy", LegacyModbus, "ABCD", SourcePoint, false}},
		{"modbus_unknown_connection", LegacyInput{Protocol: schema.ProtocolModbusRTU, ConnectionFormat: "invalid"}, Resolution{"legacy", LegacyModbus, "ABCD", SourceConnection, false}},
		{"modbus_lowercase_connection", LegacyInput{Protocol: schema.ProtocolModbusUDP, ConnectionFormat: "cdab"}, Resolution{"legacy", LegacyModbus, "CDAB", SourceConnection, false}},
		{"modbus_blank_connection", LegacyInput{Protocol: schema.ProtocolModbusTCP, ConnectionFormat: " 	"}, Resolution{"legacy", LegacyModbus, "ABCD", SourceNative, false}},
		{"mc_unknown_connection", LegacyInput{Protocol: schema.ProtocolMC3E, ConnectionFormat: "invalid"}, Resolution{"legacy", LegacyMC3E, "ABCD", SourceConnection, false}},
		{"mc_lowercase_connection", LegacyInput{Protocol: schema.ProtocolMC3E, ConnectionFormat: "cdab"}, Resolution{"legacy", LegacyMC3E, "ABCD", SourceConnection, false}},
		{"mc_padded_connection", LegacyInput{Protocol: schema.ProtocolMC3E, ConnectionFormat: " CDAB "}, Resolution{"legacy", LegacyMC3E, "ABCD", SourceConnection, false}},
		{"mc_blank_connection", LegacyInput{Protocol: schema.ProtocolMC3E, ConnectionFormat: " 	"}, Resolution{"legacy", LegacyMC3E, "ABCD", SourceConnection, false}},
		{"modbus_native", LegacyInput{Protocol: schema.ProtocolModbusTCP}, Resolution{"legacy", LegacyModbus, "ABCD", SourceNative, false}},
		{"modbus_connection", LegacyInput{Protocol: schema.ProtocolModbusRTU, ConnectionFormat: "CDAB"}, Resolution{"legacy", LegacyModbus, "CDAB", SourceConnection, false}},
		{"modbus_point", LegacyInput{Protocol: schema.ProtocolModbusUDP, PointFormat: "BADC", ConnectionFormat: "CDAB"}, Resolution{"legacy", LegacyModbus, "BADC", SourcePoint, false}},
		{"modbus_blank_point", LegacyInput{Protocol: schema.ProtocolModbusTCP, PointFormat: " \t", ConnectionFormat: "DCBA"}, Resolution{"legacy", LegacyModbus, "DCBA", SourceConnection, false}},
		{"modbus_normalization", LegacyInput{Protocol: schema.ProtocolModbusTCP, PointFormat: " badc "}, Resolution{"legacy", LegacyModbus, "BADC", SourcePoint, false}},
		{"equal_format_not_inherit", LegacyInput{Protocol: schema.ProtocolModbusTCP, PointFormat: "ABCD", ConnectionFormat: "ABCD"}, Resolution{"legacy", LegacyModbus, "ABCD", SourcePoint, false}},
		{"mc_native_ignores_point", LegacyInput{Protocol: schema.ProtocolMC3E, PointFormat: "ABCD"}, Resolution{"legacy", LegacyMC3E, "CDAB", SourceNative, false}},
		{"mc_connection_ignores_point", LegacyInput{Protocol: schema.ProtocolMC3E, PointFormat: "ABCD", ConnectionFormat: "CDAB"}, Resolution{"legacy", LegacyMC3E, "CDAB", SourceConnection, false}},
		{"mc_ignored_invalid_point", LegacyInput{Protocol: schema.ProtocolMC3E, PointFormat: "invalid", ConnectionFormat: "BADC"}, Resolution{"legacy", LegacyMC3E, "BADC", SourceConnection, false}},
		{"fatek_native_not_fictional_format", LegacyInput{Protocol: schema.ProtocolFatekFBs, PointFormat: "CDAB", ConnectionFormat: "DCBA"}, Resolution{"legacy", LegacyFatek, "", SourceNative, true}},
		{"fatek_ignored_invalid_formats", LegacyInput{Protocol: schema.ProtocolFatekFBs, PointFormat: "invalid", ConnectionFormat: "invalid"}, Resolution{"legacy", LegacyFatek, "", SourceNative, true}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			before := tc.input
			got, err := ResolveLegacy(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
			require.Equal(t, before, tc.input)
			// Serialization/reconstruction is a unit-level restart simulation, not a DB migration.
			raw, err := json.Marshal(tc.input)
			require.NoError(t, err)
			var restored LegacyInput
			require.NoError(t, json.Unmarshal(raw, &restored))
			again, err := ResolveLegacy(restored)
			require.NoError(t, err)
			require.Equal(t, got, again)
			restored.Mode = "legacy"
			explicit, err := ResolveLegacy(restored)
			require.NoError(t, err)
			require.Equal(t, got, explicit)
		})
	}
}

func TestResolveLegacyPolicyNotFrozenDefault(t *testing.T) {
	input := LegacyInput{Protocol: schema.ProtocolMC3E, PointFormat: "ABCD", ConnectionFormat: "CDAB"}
	before, err := ResolveLegacy(input)
	require.NoError(t, err)
	input.ConnectionFormat = "ABCD"
	after, err := ResolveLegacy(input)
	require.NoError(t, err)
	require.Equal(t, before.Policy, after.Policy)
	require.Equal(t, "legacy", after.Mode)
	require.Equal(t, "CDAB", before.EffectiveFormat)
	require.Equal(t, "ABCD", after.EffectiveFormat)
}

func TestResolveLegacyFailsSafely(t *testing.T) {
	for _, tc := range []struct {
		input LegacyInput
		want  error
	}{
		{LegacyInput{Protocol: schema.ProtocolModbusTCP, Mode: "override"}, ErrUnsupportedMode},
		{LegacyInput{Protocol: schema.ProtocolModbusTCP, Mode: "inherit"}, ErrUnsupportedMode},
		{LegacyInput{Protocol: schema.ProtocolModbusTCP, Mode: "garbage"}, ErrUnsupportedMode},
		{LegacyInput{Protocol: schema.ProtocolMQTT}, ErrUnsupportedProtocol},
	} {
		got, err := ResolveLegacy(tc.input)
		require.ErrorIs(t, err, tc.want)
		require.Zero(t, got, "a failed resolution must not publish a plausible format")
		require.NotContains(t, err.Error(), "secret")
	}
}
