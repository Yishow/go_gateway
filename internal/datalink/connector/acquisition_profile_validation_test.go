package connector

import (
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionConfigRejectsAmbiguousDeclaration(t *testing.T) {
	for _, raw := range []string{
		`null`, `{}`, `{"version":1}`, `{"profile":"safe-v1"}`, `{"version":1,"profile":null}`,
		`{"version":1,"profile":"safe-v1","extra":0}`, `{"Version":1,"profile":"safe-v1"}`,
		`{"version":1,"version":1,"profile":"safe-v1"}`, `{"version":1.0,"profile":"safe-v1"}`,
		`{"version":1,"profile":"safe-v1","limits":null}`, `{"version":1,"profile":"safe-v1","limits":{"native_units":null}}`,
		`{"version":1,"profile":"safe-v1","limits":{"native_units":0}}`, `{"version":1,"profile":"safe-v1","limits":{"native_units":1.5}}`,
		`{"version":1,"profile":"safe-v1","limits":{"request_bytes":-1}}`, `{"version":1,"profile":"safe-v1","limits":{"response_bytes":9223372036854775808}}`,
		`{"version":1,"profile":"safe-v1","limits":{"random_limit":1}}`,
		`{"version":1,"profile":"safe-v1","coverage":null}`, `{"version":1,"profile":"safe-v1","coverage":{"revision":" "}}`,
		`{"version":1,"profile":"safe-v1","coverage":{"revision":"v","readable":null}}`,
		`{"version":1,"profile":"safe-v1","coverage":{"revision":"v","readable":[{"area":"D","end":5}]}}`,
		`{"version":1,"profile":"safe-v1","coverage":{"revision":"v","readable":[{"area":"D","start":0,"end":5,"count":5}]}}`,
	} {
		t.Run(raw, func(t *testing.T) {
			_, err := ResolveAcquisitionConfig(schema.ProtocolMC3E, `{"acquisition":`+raw+`}`)
			require.ErrorIs(t, err, ErrAcquisitionConfig)
		})
	}
}
func TestAcquisitionConfigCoverageUsesPhysicalProtocolBounds(t *testing.T) {
	for _, tc := range []struct {
		protocol   schema.ProtocolType
		area       string
		start, end int
		valid      bool
	}{
		{schema.ProtocolModbusTCP, "HR", 65535, 65536, true}, {schema.ProtocolModbusTCP, "HR", 0, 65537, false},
		{schema.ProtocolModbusTCP, "COIL", 0, 10, true}, {schema.ProtocolModbusTCP, "C", 0, 10, false},
		{schema.ProtocolModbusTCP, "IR", 0, 10, true}, {schema.ProtocolModbusTCP, "hr", 0, 10, false},
		{schema.ProtocolMC3E, "D", 0, 16777216, true}, {schema.ProtocolMC3E, "D", 0, 16777217, false},
		{schema.ProtocolMC3E, "W", 0, 10, true}, {schema.ProtocolMC3E, "UNKNOWN", 0, 10, false},
		{schema.ProtocolFatekFBs, "DR", 1, 65536, true}, {schema.ProtocolFatekFBs, "DR", 0, 65537, false},
		{schema.ProtocolFatekFBs, "WY", 1, 9999, true}, {schema.ProtocolFatekFBs, "WY", 0, 10001, false},
		{schema.ProtocolFatekFBs, "DZZ", 0, 5, false}, {schema.ProtocolFatekFBs, "D", -1, 5, false},
		{schema.ProtocolFatekFBs, "D", 5, 5, false},
	} {
		raw, err := json.Marshal(map[string]any{"acquisition": map[string]any{"version": 1, "profile": "safe-v1", "coverage": map[string]any{"revision": "bounds", "readable": []any{map[string]any{"area": tc.area, "start": tc.start, "end": tc.end}}}}})
		require.NoError(t, err)
		_, err = ResolveAcquisitionConfig(tc.protocol, string(raw))
		require.Equal(t, tc.valid, err == nil, string(raw))
	}
}
func TestAcquisitionConfigRejectsForeignFamilyAndRoutedCPU(t *testing.T) {
	for _, raw := range []string{
		`{"acquisition":{"version":1,"profile":"safe-v1"},"acquisition":{"version":1,"profile":"mitsubishi-fx5-cpu-3e-binary-v1"}}`,
		`{"acquisition":{"version":1,"profile":"fatek-fbs-ascii-v1"}}`,
		`{"network_no":1,"acquisition":{"version":1,"profile":"mitsubishi-fx5-cpu-3e-binary-v1"}}`,
		`{"pc_no":1,"acquisition":{"version":1,"profile":"mitsubishi-fx5-cpu-3e-binary-v1"}}`,
		`{"io_no":1,"acquisition":{"version":1,"profile":"mitsubishi-fx5-cpu-3e-binary-v1"}}`,
		`{"station_no":1,"acquisition":{"version":1,"profile":"mitsubishi-fx5-cpu-3e-binary-v1"}}`,
		`{"timeout":9223372037}`, `{"timeout":1.1}`, `{"station_no":256}`,
	} {
		_, err := ResolveAcquisitionConfig(schema.ProtocolMC3E, raw)
		require.ErrorIs(t, err, ErrAcquisitionConfig, raw)
	}
	_, err := ResolveAcquisitionConfig(schema.ProtocolMQTT, `{"acquisition":{"version":1,"profile":"safe-v1"}}`)
	require.ErrorIs(t, err, ErrAcquisitionConfig)
}
func TestAcquisitionProtocolSchemaSharesDeclarationProfiles(t *testing.T) {
	for _, protocol := range []schema.ProtocolType{schema.ProtocolMC3E, schema.ProtocolFatekFBs, schema.ProtocolModbusTCP} {
		info, ok := GetProtocolInfo(protocol)
		require.True(t, ok)
		var config struct {
			Properties map[string]json.RawMessage `json:"properties"`
		}
		require.NoError(t, json.Unmarshal(info.ConfigSchema, &config))
		require.JSONEq(t, string(AcquisitionConfigSchema(protocol)), string(config.Properties["acquisition"]))
		var declaration struct {
			Properties struct {
				Profile struct {
					Enum []string `json:"enum"`
				} `json:"profile"`
			} `json:"properties"`
		}
		require.NoError(t, json.Unmarshal(config.Properties["acquisition"], &declaration))
		for _, profile := range declaration.Properties.Profile.Enum {
			raw, _ := json.Marshal(map[string]any{"acquisition": map[string]any{"version": 1, "profile": profile}})
			_, err := ResolveAcquisitionConfig(protocol, string(raw))
			require.NoError(t, err)
		}
	}
	info, _ := GetProtocolInfo(schema.ProtocolMQTT)
	require.NotContains(t, string(info.ConfigSchema), "acquisition")
}
