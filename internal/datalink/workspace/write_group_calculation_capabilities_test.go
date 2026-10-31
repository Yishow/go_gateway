package workspace

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/snapshot"
)

func TestWriteGroupReadbackAdvertisesPreparationNotActivation(t *testing.T) {
	for _, result := range []any{WriteGroupSaveResult{}, WriteGroupListResult{}} {
		raw, err := json.Marshal(result)
		require.NoError(t, err)
		var data map[string]any
		require.NoError(t, json.Unmarshal(raw, &data))
		cap, ok := data["calculation_capabilities"].(map[string]any)
		require.True(t, ok, "canonical responses must tell consumers whether mixed is actually available")
		require.Equal(t, true, cap["draft_supported"])
		require.Equal(t, false, cap["activation_supported"])
		require.Equal(t, false, cap["editor_supported"])
		require.Equal(t, "production_acceptance_pending", cap["reason"])
		cadence := cap["cadence"].(map[string]any)
		require.EqualValues(t, snapshot.MaxCadenceSeconds, cadence["max_seconds"])
		require.EqualValues(t, snapshot.MaxCadenceMinutes, cadence["max_minutes"])
		require.Equal(t, "utc_epoch", cadence["alignment"])
	}
}

func TestWriteGroupReadbackIgnoresCallerCapabilityOverride(t *testing.T) {
	forged := &WriteGroupCalculationCapabilities{Version: 99, ActivationSupported: true, EditorSupported: true, Reason: "forged"}
	for _, pair := range [][2]any{
		{WriteGroupSaveResult{}, WriteGroupSaveResult{Capabilities: forged}},
		{WriteGroupListResult{}, WriteGroupListResult{Capabilities: forged}},
	} {
		expected, err := json.Marshal(pair[0])
		require.NoError(t, err)
		actual, err := json.Marshal(pair[1])
		require.NoError(t, err)
		require.Equal(t, string(expected), string(actual), "caller field cannot change wire bytes or backend capability authority")
		require.NotContains(t, string(actual), "forged")
	}
	require.True(t, forged.ActivationSupported, "marshalling does not mutate caller state")
}
