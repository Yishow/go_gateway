package grouppipeline

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func healthyCounterRevision(t *testing.T, f *longPipeline, groupID, revision string) *workspace.WriteGroupAppliedSnapshot {
	t.Helper()
	raw, err := f.store.RuntimeVersion(t.Context(), groupKey(f.snapshot))
	require.NoError(t, err)
	var snap workspace.WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal(raw, &snap))
	snap.GroupID, snap.GroupRevision, snap.AppliedRevision = groupID, revision, revision
	snap.Group.ID, snap.Group.Revision, snap.Group.AppliedRevision = groupID, revision, revision
	snap.Group.Members = snap.Group.Members[:1]
	member := &snap.Group.Members[0]
	member.DeviceID = "healthy-device"
	member.PointID = "healthy-point"
	member.TagID = "healthy-tag"
	member.SourceRevision = "healthy-source"
	member.MappingRevision = "healthy-map"
	member.MeasurementID = ptr("healthy-measurement")
	member.ExactRecordingSource = nil
	member.SourceValueType, member.ResultValueType = "", ""
	snap.RuntimeLayout.TagTypes = map[string]schema.DataType{"healthy-tag": schema.DataTypeUint64}
	freezeFixtureCounters(t, &snap, map[string]runtime.CounterSource{"healthy-tag": {MeasurementID: "healthy-measurement", DefinitionRevision: "def-1", SourceBindingRevision: "bind-1", SeriesEpoch: "epoch-1", Unit: "kWh", Kind: aggregation.CounterKindMonotonic}})
	require.NoError(t, f.pipe.saveRuntimeSnapshot(t.Context(), &snap))
	return &snap
}

func TestCounterReplayQuarantineDoesNotBlockHealthyScopeOrClearOnCurrentChanges(t *testing.T) {
	for _, sameGroup := range []bool{false, true} {
		t.Run(map[bool]string{false: "other group", true: "new revision"}[sameGroup], func(t *testing.T) {
			f := newLongPipeline(t)
			original, err := f.store.RuntimeVersion(t.Context(), groupKey(f.snapshot))
			require.NoError(t, err)
			groupID := "healthy-group"
			if sameGroup {
				groupID = f.group.ID
			}
			healthy := healthyCounterRevision(t, f, groupID, "healthy-revision")
			require.NoError(t, f.reading(0, "old-pending", at(1), 100))
			damageCounterProof(t, f, "counter-replay-proof-missing")
			// An unrelated healthy scope can start, and even the same group's newly
			// applied revision must not clear an older revision's missing evidence.
			f.groups.groups = append(f.groups.groups, liveGroup(healthy.Group, workspace.WriteGroupStatusReady))
			f.groups.snapshots[groupID] = append(f.groups.snapshots[groupID], healthy)
			reader := &unavailableExactCounterReader{}
			f.pipe = New(Dependencies{Groups: f.groups, Destinations: f.dest, Store: f.store, Counters: reader}, Config{Now: func() time.Time { return f.clock }})
			require.NoError(t, f.pipe.Reconcile(t.Context()))
			require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(healthy.Group, "healthy-1", at(2), uint64(7))))
			var healthyInputs int
			require.NoError(t, f.local.QueryRow(`SELECT COUNT(*) FROM wg_delivery_samples WHERE group_id=? AND group_revision=?`, groupID, "healthy-revision").Scan(&healthyInputs))
			require.Equal(t, 1, healthyInputs)
			for range 2 {
				entries, err := f.pipe.CounterReplayQuarantines(t.Context(), f.group.ID)
				require.NoError(t, err)
				require.Len(t, entries, 1)
				require.Equal(t, "rev-1", entries[0].GroupRevision)
				require.Equal(t, 1, entries[0].PendingSamples)
				// Editable provider changes and a draft cannot repair original evidence.
				f.counters.sources[longKwhFirst] = runtime.CounterSource{Unit: "MWh", SeriesEpoch: "epoch-2"}
				require.NoError(t, f.pipe.Reconcile(t.Context()))
				require.Zero(t, reader.calls)
			}
			// Fixture-only restoration of the exact trusted bytes is not a repair API.
			_, err = f.local.Exec(`UPDATE wg_runtime_versions SET payload=? WHERE group_id=? AND group_revision='rev-1'`, string(original), f.group.ID)
			require.NoError(t, err)
			require.NoError(t, f.pipe.Reconcile(t.Context()))
			entries, err := f.pipe.CounterReplayQuarantines(t.Context(), f.group.ID)
			require.NoError(t, err)
			require.Empty(t, entries)
			require.Zero(t, reader.calls)
		})
	}
}

func TestCounterReplayQuarantineDoesNotTreatIncomingClaimsAsStoredDamage(t *testing.T) {
	f := newLongPipeline(t)
	require.NoError(t, f.reading(0, "accepted", at(1), 100))
	bad := envelope(f.group, "wrong-epoch", at(2), uint64(105))
	bad.SeriesEpoch = "epoch-foreign-secret"
	err := f.pipe.AcceptSample(t.Context(), bad)
	require.ErrorIs(t, err, runtime.ErrCounterSampleClaimConflict)
	require.NotContains(t, err.Error(), "foreign-secret")
	entries, err := f.pipe.CounterReplayQuarantines(t.Context(), f.group.ID)
	require.NoError(t, err)
	require.Empty(t, entries)
	require.Equal(t, StateActive, f.statusOf().State)
	bad.SeriesEpoch = ""
	require.NoError(t, f.pipe.AcceptSample(t.Context(), bad))
}

func TestCounterReplayQuarantineClassifiesLivePersistedProofFailure(t *testing.T) {
	f := newLongPipeline(t)
	require.NoError(t, f.reading(0, "accepted", at(1), 100))
	damageCounterProof(t, f, "counter-replay-proof-conflict")
	before := counterDurableState(t, f.store)
	err := f.reading(0, "refused", at(2), 105)
	require.ErrorIs(t, err, counterproof.ErrCounterBindingProofConflict)
	entries, err := f.pipe.CounterReplayQuarantines(t.Context(), f.group.ID)
	require.NoError(t, err)
	require.Len(t, entries, 1)
	f.closeUntil(at(1200))
	require.Equal(t, before, counterDurableState(t, f.store))
	require.False(t, f.pipe.WantsSample("device-1", "point-1", longKwhFirst))
}

func TestCounterReplayQuarantineReadFailureIsNotAnEmptySuccess(t *testing.T) {
	f := newLongPipeline(t)
	require.NoError(t, f.reading(0, "pending", at(1), 100))
	damageCounterProof(t, f, "counter-replay-proof-missing")
	f.reconcileStatus()
	injected := errors.New("private counting storage failure")
	reject := func(query string) error {
		if strings.Contains(query, "COUNT(*) FROM wg_delivery_samples WHERE workspace_id") {
			return injected
		}
		return nil
	}
	f.localFault.Reject.Store(&reject)
	view, err := f.pipe.Delivery(t.Context(), f.group.ID)
	require.Nil(t, view)
	require.ErrorIs(t, err, injected)
	f.localFault.Reject.Store(nil)
	entries, err := f.pipe.CounterReplayQuarantines(t.Context(), f.group.ID)
	require.NoError(t, err)
	require.Len(t, entries, 1)
}

func TestCounterReplayRecoveryInfrastructureFailureRemainsAnError(t *testing.T) {
	f := newLongPipeline(t)
	require.NoError(t, f.reading(0, "pending", at(1), 100))
	f.groups.groups = nil
	f.pipe.boundaries = map[string]*managed{}
	injected := errors.New("private runtime database unavailable")
	reject := func(query string) error {
		if strings.Contains(query, "SELECT payload FROM wg_runtime_versions") {
			return injected
		}
		return nil
	}
	f.localFault.Reject.Store(&reject)
	err := f.pipe.Reconcile(t.Context())
	require.ErrorIs(t, err, injected, "an unclassified storage failure is not a proof quarantine or successful recovery")
	f.localFault.Reject.Store(nil)
	entries, err := f.pipe.CounterReplayQuarantines(t.Context(), f.group.ID)
	require.NoError(t, err)
	require.Empty(t, entries)
}

func TestCounterReplayQuarantineUsesFailedOriginalRevisionNotLatestAppliedRevision(t *testing.T) {
	f := newLongPipeline(t)
	require.NoError(t, f.reading(0, "pending", at(1), 100))
	damageCounterProof(t, f, "counter-replay-proof-missing")
	// A newly applied revision can become effective later; ResolveAppliedAt
	// still selects the original revision for the present intake boundary.
	f.groups.groups[0].AppliedRevision = "future-revision"
	f.reconcileStatus()
	entries, err := f.pipe.CounterReplayQuarantines(t.Context(), f.group.ID)
	require.NoError(t, err)
	require.Equal(t, []CounterReplayQuarantine{{WorkspaceID: "workspace-A", GroupID: "group-C", GroupRevision: "rev-1", PendingSamples: 1, Reason: "counter-replay-proof-missing", Action: proofRepairAction}}, entries)
}

func TestCounterReplayQuarantineClearingIsLimitedToTheValidatedFullScope(t *testing.T) {
	f := newLongPipeline(t)
	healthy := healthyCounterRevision(t, f, f.group.ID, "rev-2")
	restored, reason := f.pipe.build(t.Context(), healthy)
	require.Empty(t, reason)
	key := groupKey(f.snapshot)
	require.True(t, f.pipe.setCounterQuarantine(key, "counter-replay-proof-missing"))
	foreign := key
	foreign.WorkspaceID = "workspace-B"
	require.True(t, f.pipe.setCounterQuarantine(foreign, "counter-replay-proof-conflict"))
	require.True(t, f.pipe.setCounterQuarantine(groupKey(healthy), "counter-replay-proof-unsupported"))
	f.pipe.installManaged(restored)
	entries, err := f.pipe.CounterReplayQuarantines(t.Context(), f.group.ID)
	require.NoError(t, err)
	require.Equal(t, []CounterReplayQuarantine{
		{WorkspaceID: "workspace-A", GroupID: "group-C", GroupRevision: "rev-1", Reason: "counter-replay-proof-missing", Action: proofRepairAction},
		{WorkspaceID: "workspace-B", GroupID: "group-C", GroupRevision: "rev-1", Reason: "counter-replay-proof-conflict", Action: proofRepairAction},
	}, entries, "validating workspace-A/rev-2 cannot clear any other original scope")
	require.False(t, f.pipe.setCounterQuarantine(key, "counter-sample-claim-conflict"))
	foreign.WorkspaceID = ""
	require.False(t, f.pipe.setCounterQuarantine(foreign, "counter-replay-proof-missing"))
}
