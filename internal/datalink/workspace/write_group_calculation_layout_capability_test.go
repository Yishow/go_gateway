package workspace

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMixedExactScaleCapabilityRestrictsWideWithoutDenyingLongDelta(t *testing.T) {
	raw, err := json.Marshal(WriteGroupSaveResult{})
	require.NoError(t, err)
	var wire struct {
		Capabilities struct {
			Restrictions []map[string]any `json:"unsupported_combinations"`
		} `json:"calculation_capabilities"`
	}
	require.NoError(t, json.Unmarshal(raw, &wire))
	found := false
	for _, restriction := range wire.Capabilities.Restrictions {
		if restriction["method"] == "delta" && restriction["phase"] == "draft" {
			found = true
			require.Equal(t, "managed_wide_profile_v1", restriction["layout"], "long delta's normal scoped save is proven; wide exact scale still refuses")
		} else {
			require.NotContains(t, restriction, "layout", "omission retains existing all-layout restriction")
		}
	}
	require.True(t, found, "wide must not be accidentally advertised by removing the entire delta restriction")
}
