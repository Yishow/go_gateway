package workspace

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMixedRowPolicyOmitsLegacyIntervalOnly(t *testing.T) {
	legacy := &WriteGroup{RowPolicy: WriteGroupRowPolicy{IntervalSeconds: 0, AllowedLatenessSeconds: 3}}
	legacyWire, err := json.Marshal(legacy)
	require.NoError(t, err)
	require.Contains(t, string(legacyWire), `"interval_seconds":0`)
	mixed := cloneWriteGroup(legacy)
	mixed.RecordingProfile = &WriteGroupRecordingProfile{Kind: WriteGroupProfileMixed, ResultShape: WriteGroupResultShapeMemberCalculationV1}
	wire, err := json.Marshal(mixed)
	require.NoError(t, err)
	require.NotContains(t, string(wire), "interval_seconds")
	require.Contains(t, string(wire), `"allowed_lateness_seconds":3`)
	stored, err := marshalWriteGroupRowPolicy(mixed)
	require.NoError(t, err)
	var document map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(wire, &document))
	require.JSONEq(t, string(document["row_policy"]), string(stored))
	var decoded WriteGroup
	require.NoError(t, json.Unmarshal(wire, &decoded))
	require.Zero(t, decoded.RowPolicy.IntervalSeconds)
}

func TestMixedJSONRejectsExplicitGroupWindow(t *testing.T) {
	for _, interval := range []string{"0", "60", "null"} {
		raw := `{"recording_profile":{"kind":"mixed","result_shape":"member-calculation-v1"},"row_policy":{"interval_seconds":` + interval + `}}`
		var group WriteGroup
		require.Error(t, json.Unmarshal([]byte(raw), &group), "explicit legacy group window must not become a mixed default")
	}
}
