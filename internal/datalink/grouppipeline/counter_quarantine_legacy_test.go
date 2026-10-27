package grouppipeline

import (
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func legacyQuarantineFixture(t *testing.T) *longPipeline {
	t.Helper()
	f := newLongPipelineOn(t, nil, true)
	key := groupKey(f.snapshot)
	raw, err := f.store.RuntimeVersion(t.Context(), key)
	require.NoError(t, err)
	var snap workspace.WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal(raw, &snap))
	snap.CounterBindings = nil
	legacy, err := json.Marshal(snap)
	require.NoError(t, err)
	_, err = f.local.Exec(`DELETE FROM wg_runtime_versions`)
	require.NoError(t, err)
	value, err := measurement.NewDecimal("100.0")
	require.NoError(t, err)
	sample := snapshot.Sample{SampleID: "legacy-original", MemberKey: `["device-1","point-1","tag-kwh-1"]`, ObservedAt: at(1), Quality: schema.QualityGood, Value: value, SourceRevision: "src-1", MappingRevision: "map-1"}
	require.NoError(t, f.store.AppendSample(t.Context(), key, at(0), sample))
	require.NoError(t, f.store.SaveRuntimeVersion(t.Context(), key, legacy))
	return f
}

func TestCounterReplayQuarantinePreservesOriginalLegacyJournalAndAssociationAudit(t *testing.T) {
	f := legacyQuarantineFixture(t)
	f.reconcileStatus()
	require.Equal(t, 1, f.count(f.local, `SELECT COUNT(*) FROM wg_counter_proof_associations`))
	// A conflicting persisted audit is evidence, not a reason to replace or delete it.
	_, err := f.local.Exec(`UPDATE wg_counter_proof_associations SET proof_digest='conflicting-original-evidence'`)
	require.NoError(t, err)
	before := counterDurableState(t, f.store)
	reader := &unavailableExactCounterReader{}
	for range 2 {
		f.pipe = New(Dependencies{Groups: f.groups, Destinations: f.dest, Store: f.store, Counters: reader}, Config{Now: func() time.Time { return f.clock }})
		require.NoError(t, f.pipe.Reconcile(t.Context()))
		entries, err := f.pipe.CounterReplayQuarantines(t.Context(), f.group.ID)
		require.NoError(t, err)
		require.Equal(t, []CounterReplayQuarantine{{WorkspaceID: "workspace-A", GroupID: "group-C", GroupRevision: "rev-1", PendingSamples: 1, Reason: "counter-replay-proof-conflict", Action: proofRepairAction}}, entries)
		f.closeUntil(at(1200))
		require.Equal(t, before, counterDurableState(t, f.store))
		require.Zero(t, reader.calls)
	}
}
