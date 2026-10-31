package grouppipeline

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/workspace"
)

func TestMixedFrozenReplayAndFinalizedStorageKeepDistinctProofs(t *testing.T) {
	for _, long := range []bool{false, true} {
		f := newFixture(t)
		g, l := mixedStorageProofFixture(t, long, dbtarget.SQLDialectSQLite)
		id := "measurement-mixed"
		m := &g.Members[0]
		m.Calculation = workspace.WriteGroupMemberCalculationDelta
		m.ValueRole = workspace.WriteGroupValueRoleCounterRaw
		m.MeasurementID = &id
		proof, err := counterproof.NewCounterBindingProof(counterproof.CounterBindingFacts{
			Member:        counterproof.CounterBindingMemberIdentity{WorkspaceID: g.WorkspaceID, GroupID: g.ID, AppliedRevision: g.Revision, DeviceID: m.DeviceID, PointID: m.PointID, TagID: m.TagID},
			MeasurementID: id, DefinitionRevision: "definition-1", SourceBindingRevision: "binding-1", SeriesEpoch: "epoch-1", Unit: "kWh", SemanticKind: measurement.SemanticKindCounter,
			SourceRevision: m.SourceRevision, MappingRevision: m.MappingRevision, SourceValueType: m.SourceValueType, ResultValueType: measurement.ExactFloat64, ScaleMode: counterproof.CounterScaleUnscaled, ExactScale: "1",
		})
		require.NoError(t, err)
		snap := &workspace.WriteGroupAppliedSnapshot{WorkspaceID: g.WorkspaceID, GroupID: g.ID, GroupRevision: g.Revision, AppliedRevision: g.AppliedRevision, EffectiveAt: at(0), Group: g, RuntimeLayout: l, CounterBindings: counterproof.CounterBindingProofs{proof.Member.Key(): proof}}
		raw, err := json.Marshal(snap)
		require.NoError(t, err)
		key := groupKey(snap)
		require.NoError(t, f.store.SaveRuntimeVersion(t.Context(), key, raw))
		loaded, err := loadFrozenSnapshot(t.Context(), f.store, key)
		require.NoError(t, err)
		require.NoError(t, workspace.ValidateMixedRecordingSnapshot(loaded.Group, loaded.RuntimeLayout.TagTypes, loaded.CounterBindings))
		projected, err := counterproof.SnapshotStorageProjection(raw)
		require.NoError(t, err)
		storage, err := decodeFrozenStorageSnapshot(projected, key)
		require.NoError(t, err)
		require.Empty(t, storage.CounterBindings)
		require.Error(t, workspace.ValidateMixedRecordingSnapshot(storage.Group, storage.RuntimeLayout.TagTypes, storage.CounterBindings), "storage projection cannot authorize replay")
		_, err = f.store.DB().Exec(`UPDATE wg_runtime_versions SET payload=json_remove(payload,'$.counter_bindings') WHERE group_id=?`, g.ID)
		require.NoError(t, err)
		_, err = loadFrozenSnapshot(t.Context(), f.store, key)
		require.ErrorIs(t, err, counterproof.ErrCounterBindingProofMissing)
		finalized, err := loadFinalizedStorageSnapshot(t.Context(), f.store, key)
		require.NoError(t, err)
		require.Equal(t, g.Destination, finalized.Group.Destination)
		_, err = f.store.DB().Exec(`UPDATE wg_runtime_versions SET payload=json_set(payload,'$.group.workspace_id','foreign') WHERE group_id=?`, g.ID)
		require.NoError(t, err)
		_, err = loadFinalizedStorageSnapshot(t.Context(), f.store, key)
		require.ErrorIs(t, err, errRuntimeSnapshotInvalid)
	}
}
