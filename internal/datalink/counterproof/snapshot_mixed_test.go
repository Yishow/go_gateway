package counterproof

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func mixedProofSnapshot(t *testing.T, method string, withProof bool) ([]byte, CounterBindingProof) {
	t.Helper()
	proof, err := NewCounterBindingProof(CounterBindingFacts{
		Member: CounterBindingMemberIdentity{WorkspaceID: "w", GroupID: "g", AppliedRevision: "r", DeviceID: "d", PointID: "p", TagID: "t"}, MeasurementID: "m", DefinitionRevision: "def", SourceBindingRevision: "binding", SeriesEpoch: "epoch", SemanticKind: measurement.SemanticKindCounter, Unit: "kWh", SourceRevision: "src", MappingRevision: "map", ScaleMode: CounterScaleUnscaled, ExactScale: "1", SourceValueType: schema.DataTypeUint64, ResultValueType: measurement.ExactUint64,
	})
	require.NoError(t, err)
	bindings := CounterBindingProofs{}
	if withProof {
		bindings[proof.Member.Key()] = proof
	}
	data, err := json.Marshal(map[string]any{
		"workspace_id": "w", "group_id": "g", "applied_revision": "r", "counter_bindings": bindings,
		"group":          map[string]any{"workspace_id": "w", "id": "g", "revision": "r", "applied_revision": "r", "layout_version": "managed_long_v1", "destination": map[string]any{"storage_strategy": "managed"}, "recording_profile": map[string]any{"kind": "mixed", "result_shape": "member-calculation-v1"}, "members": []any{map[string]any{"device_id": "d", "point_id": "p", "tag_id": "t", "measurement_id": "m", "source_revision": "src", "mapping_revision": "map", "source_value_type": "uint64", "result_value_type": "decimal", "calculation": method, "cadence": map[string]any{"interval": 15, "unit": "second", "alignment": "utc_epoch"}}}},
		"runtime_layout": map[string]any{"tag_types": map[string]any{"t": "uint64"}},
	})
	require.NoError(t, err)
	return data, proof
}

func TestMixedSnapshotCounterProofUsesDeltaOnlyAndInputType(t *testing.T) {
	data, proof := mixedProofSnapshot(t, "delta", true)
	got, err := ResolveSnapshot(data, SnapshotScope{"w", "g", "r"})
	require.NoError(t, err)
	require.Equal(t, proof, got[proof.Member.Key()], "decimal delta result must not be confused with uint64 counter input")
	for _, method := range []string{"latest", "mean"} {
		data, _ = mixedProofSnapshot(t, method, false)
		got, err = ResolveSnapshot(data, SnapshotScope{"w", "g", "r"})
		require.NoError(t, err)
		require.Empty(t, got)
		data, _ = mixedProofSnapshot(t, method, true)
		_, err = ResolveSnapshot(data, SnapshotScope{"w", "g", "r"})
		require.ErrorIs(t, err, ErrCounterBindingProofConflict)
	}
	data, _ = mixedProofSnapshot(t, "delta", false)
	_, err = ResolveSnapshot(data, SnapshotScope{"w", "g", "r"})
	require.ErrorIs(t, err, ErrCounterBindingProofMissing)
}
