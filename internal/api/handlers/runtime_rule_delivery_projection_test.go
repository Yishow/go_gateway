package handlers

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/workspace"
)

func runtimeDeliveryProjectionFixture() (*workspace.WriteGroup, *workspace.WriteGroupAppliedSnapshot, runtimeview.Scope, *grouppipeline.DeliveryView, time.Time) {
	scope := runtimeview.Scope{WorkspaceID: "workspace-W", WorkspaceRevision: "workspace-R", DeviceID: "device-A"}
	at := time.Date(2026, 10, 1, 12, 0, 15, 0, time.UTC)
	saved := &workspace.WriteGroup{ID: "group-G", WorkspaceID: scope.WorkspaceID, Revision: "C", AppliedRevision: "B", Status: workspace.WriteGroupStatusReady,
		Destination:   workspace.WriteGroupDestination{ConnectorID: "C2", ConnectorRevision: "8", TableName: "new_table"},
		LayoutVersion: workspace.WriteGroupLayoutManagedWideV1, StorageMode: workspace.WriteGroupStorageModePlain,
		RowPolicy: workspace.WriteGroupRowPolicy{IntervalSeconds: 15, IncompletePolicy: "skip_row"}, WritePolicy: workspace.WriteGroupWritePolicy{Mode: "append"},
		Members: []workspace.WriteGroupMember{{DeviceID: "device-A", PointID: "point-A", TagID: "tag-A", TargetColumn: "reading"}}}
	effective := *saved
	effective.Revision, effective.AppliedRevision = "A", "A"
	effective.Destination = workspace.WriteGroupDestination{ConnectorID: "C1", ConnectorRevision: "7", TableName: "old_table"}
	effective.Members = append([]workspace.WriteGroupMember(nil), saved.Members...)
	snapshot := &workspace.WriteGroupAppliedSnapshot{WorkspaceID: scope.WorkspaceID, GroupID: saved.ID, GroupRevision: "A", AppliedRevision: "A", EffectiveAt: at.Add(-time.Minute), Group: &effective, RuntimeLayout: &workspace.WriteGroupRuntimeLayout{Dialect: dbtarget.SQLDialectSQLite}}
	raw := &grouppipeline.DeliveryView{GroupID: saved.ID, RevisionStages: &grouppipeline.RevisionStagesView{GroupRevision: "A", Stages: grouppipeline.StagesView{Collecting: 2, Queued: 3}}, CounterReplayQuarantines: []grouppipeline.CounterReplayQuarantine{}, Backlog: []grouppipeline.BacklogView{}}
	return saved, snapshot, scope, raw, at
}

func runtimeDeliveryMarshal(t *testing.T, saved *workspace.WriteGroup, snapshot *workspace.WriteGroupAppliedSnapshot, scope runtimeview.Scope, delivery *runtimeview.DeliveryEvidence, at time.Time) []byte {
	t.Helper()
	observed := at.Format(time.RFC3339Nano)
	group, err := projectRuntimeRecordingGroup(saved, snapshot, scope, observed)
	require.NoError(t, err)
	group.Delivery = delivery
	raw, err := runtimeview.MarshalRuleEvidence(runtimeview.RuleEvidenceView{ContractVersion: 1, WorkspaceID: scope.WorkspaceID, WorkspaceRevision: scope.WorkspaceRevision, DeviceID: scope.DeviceID, ObservedAt: &observed, Freshness: "fresh", Groups: []runtimeview.GroupView{group}}, scope)
	require.NoError(t, err)
	return raw
}

func TestRuntimeRuleDeliveryCollectingQueuedNeverSQLCommitted(t *testing.T) {
	saved, snapshot, scope, raw, at := runtimeDeliveryProjectionFixture()
	// Aggregate current+history counters are not current-revision evidence.
	raw.Stages.SQLCommitted = 91
	got, err := projectRuntimeRuleDelivery(raw, saved, snapshot, scope, nil)
	require.NoError(t, err)
	require.EqualValues(t, 2, got.Current.Stages.Collecting)
	require.EqualValues(t, 3, got.Current.Stages.Queued)
	require.Zero(t, got.Current.Stages.SQLCommitted)
	require.Nil(t, got.Current.LastCommit)
	require.Nil(t, got.Current.Intake)
	require.Equal(t, "group", got.Current.TotalsScope)
	require.Equal(t, "unavailable", got.Current.SourceState)
	wire := string(runtimeDeliveryMarshal(t, saved, snapshot, scope, got, at))
	for _, field := range []string{"next_due_at", "next_retry_at", "attempt_duration_ms", "oldest_pending_seconds", "quota", "storage", "test_write"} {
		require.Contains(t, wire, `"`+field+`":null`)
	}
	require.Contains(t, wire, `"sources":[]`)
	require.Contains(t, wire, `"counter_replay_quarantines":[]`)
}

func TestRuntimeRuleDeliveryDisabledHistoricalDestinationNotRelabeled(t *testing.T) {
	saved, snapshot, scope, raw, at := runtimeDeliveryProjectionFixture()
	saved.Status = workspace.WriteGroupStatusDisabled
	raw.Backlog = []grouppipeline.BacklogView{
		{GroupRevision: "A", ConnectorID: "C1", ConnectorRevision: "7", TableName: "old_table", Pending: 9, ErrorCodes: []string{"recording-ownership-conflict"}},
		{GroupRevision: "A0", ConnectorID: "C1", ConnectorRevision: "7", TableName: "old_table", Pending: 3, ErrorCodes: []string{"recording-ownership-unavailable", "recording-ownership-conflict", "recording-handoff-proof-conflict", "recording-cutover-finalized-overlap", "password=private-secret"}},
	}
	got, err := projectRuntimeRuleDelivery(raw, saved, snapshot, scope, nil)
	require.NoError(t, err)
	require.Len(t, got.Backlog, 1)
	b := got.Backlog[0]
	require.Equal(t, "A0", b.Binding.EntityRevision)
	require.Equal(t, runtimeview.Destination{ConnectorID: "C1", ConnectorRevision: "7", Table: "old_table", Kind: "unavailable"}, b.Destination)
	require.EqualValues(t, 3, b.Pending)
	require.Equal(t, []runtimeview.Reason{runtimeview.ReasonOwnershipUnavailable, runtimeview.ReasonOwnershipConflict, runtimeview.ReasonHandoffConflict, runtimeview.ReasonFinalizedOverlap, runtimeview.ReasonUnavailable}, b.Reasons)
	require.NotContains(t, string(runtimeDeliveryMarshal(t, saved, snapshot, scope, got, at)), "private-secret")
}

func TestRuntimeRuleDeliverySharedFrozenMemberTotalsNotPerTag(t *testing.T) {
	saved, snapshot, scope, raw, at := runtimeDeliveryProjectionFixture()
	snapshot.Group.Members = append(snapshot.Group.Members, workspace.WriteGroupMember{DeviceID: "device-B", PointID: "point-B", TagID: "tag-B", TargetColumn: "other"})
	got, err := projectRuntimeRuleDelivery(raw, saved, snapshot, scope, nil)
	require.NoError(t, err)
	require.True(t, got.Current.Shared)
	require.EqualValues(t, 3, got.Current.Stages.Queued)
	require.Empty(t, got.Current.Sources)
	runtimeDeliveryMarshal(t, saved, snapshot, scope, got, at)
}

func TestRuntimeRuleDeliveryReplayInputsSeparateFromFinalizedQuarantine(t *testing.T) {
	for _, reason := range []string{"counter-replay-proof-missing", "counter-replay-proof-conflict", "counter-replay-proof-unsupported"} {
		t.Run(reason, func(t *testing.T) {
			saved, snapshot, scope, raw, at := runtimeDeliveryProjectionFixture()
			raw.RevisionStages.Stages.Quarantined = 2
			raw.CounterReplayQuarantines = []grouppipeline.CounterReplayQuarantine{{WorkspaceID: scope.WorkspaceID, GroupID: saved.ID, GroupRevision: "A0", PendingSamples: 7, Reason: reason, Action: grouppipeline.CounterReplayRepairAction}}
			got, err := projectRuntimeRuleDelivery(raw, saved, snapshot, scope, nil)
			require.NoError(t, err)
			require.EqualValues(t, 2, got.Current.Stages.Quarantined)
			require.EqualValues(t, 7, (*got.CounterReplayQuarantines)[0].PendingSamples)
			require.Equal(t, runtimeview.ActionInspectImmutableEvidence, runtimeview.ActionForReason((*got.CounterReplayQuarantines)[0].Reason))
			runtimeDeliveryMarshal(t, saved, snapshot, scope, got, at)
		})
	}
}

func TestRuntimeRuleDeliveryEmptyRevisionNoCurrentZeros(t *testing.T) {
	saved, _, scope, raw, at := runtimeDeliveryProjectionFixture()
	raw.RevisionStages = nil
	raw.CounterReplayQuarantines = nil
	raw.Backlog = []grouppipeline.BacklogView{{GroupRevision: "A", ConnectorID: "C1", ConnectorRevision: "7", TableName: "old_table", Pending: 3}}
	got, err := projectRuntimeRuleDelivery(raw, saved, nil, scope, nil)
	require.NoError(t, err)
	require.Nil(t, got.Current)
	require.Nil(t, got.CounterReplayQuarantines)
	require.Len(t, got.Backlog, 1)
	runtimeDeliveryMarshal(t, saved, nil, scope, got, at)
}

func TestRuntimeRuleDeliveryMissingOrMismatchedCommitProofFailsClosed(t *testing.T) {
	cases := map[string]struct {
		mutate func(*grouppipeline.DeliveryView)
		reason runtimeview.Reason
	}{
		"missing":                  {func(v *grouppipeline.DeliveryView) { v.LastSQLCommittedEffect = nil }, runtimeview.ReasonUnavailable},
		"zero-with-receipt":        {func(v *grouppipeline.DeliveryView) { v.RevisionStages.Stages.SQLCommitted = 0 }, runtimeview.ReasonUnavailable},
		"wrong-group-revision":     {func(v *grouppipeline.DeliveryView) { v.LastSQLCommittedEffect.GroupRevision = "B" }, runtimeview.ReasonStale},
		"wrong-connector-revision": {func(v *grouppipeline.DeliveryView) { v.LastSQLCommittedEffect.ConnectorRevision = "8" }, runtimeview.ReasonStale},
		"unsafe-id":                {func(v *grouppipeline.DeliveryView) { v.LastSQLCommittedEffect.RecordID = "postgres://private-secret" }, runtimeview.ReasonUnavailable},
		"missing-time":             {func(v *grouppipeline.DeliveryView) { v.LastSQLCommittedEffect.CommittedAt = time.Time{} }, runtimeview.ReasonUnavailable},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			saved, snapshot, scope, raw, at := runtimeDeliveryProjectionFixture()
			raw.RevisionStages.Stages.SQLCommitted = 1
			raw.LastSQLCommittedEffect = &groupdelivery.CommittedEffect{GroupRevision: "A", ConnectorRevision: "7", RecordID: "record-1", EffectKey: "effect-1", PayloadDigest: "digest-1", CommittedAt: at}
			got, err := projectRuntimeRuleDelivery(raw, saved, snapshot, scope, nil)
			require.NoError(t, err)
			require.Equal(t, "runtime", got.Current.LastCommit.Provenance)
			runtimeDeliveryMarshal(t, saved, snapshot, scope, got, at)
			tc.mutate(raw)
			_, err = projectRuntimeRuleDelivery(raw, saved, snapshot, scope, nil)
			require.True(t, runtimeview.IsContractError(err, tc.reason), "%v", err)
		})
	}
}

func TestRuntimeRuleDeliveryInvalidCountsScopeAndQuarantineFailClosed(t *testing.T) {
	cases := map[string]struct {
		mutate func(*grouppipeline.DeliveryView)
		reason runtimeview.Reason
	}{
		"foreign-group":      {func(v *grouppipeline.DeliveryView) { v.GroupID = "foreign" }, runtimeview.ReasonForeign},
		"wrong-revision":     {func(v *grouppipeline.DeliveryView) { v.RevisionStages.GroupRevision = "B" }, runtimeview.ReasonStale},
		"missing-stages":     {func(v *grouppipeline.DeliveryView) { v.RevisionStages = nil }, runtimeview.ReasonUnavailable},
		"negative-count":     {func(v *grouppipeline.DeliveryView) { v.RevisionStages.Stages.Queued = -1 }, runtimeview.ReasonUnavailable},
		"unsafe-count":       {func(v *grouppipeline.DeliveryView) { v.RevisionStages.Stages.SQLCommitted = 9007199254740992 }, runtimeview.ReasonUnavailable},
		"foreign-quarantine": {func(v *grouppipeline.DeliveryView) { v.CounterReplayQuarantines[0].WorkspaceID = "foreign" }, runtimeview.ReasonForeign},
		"bad-quarantine-reason": {func(v *grouppipeline.DeliveryView) {
			v.CounterReplayQuarantines[0].Reason = "recording-ownership-conflict"
		}, runtimeview.ReasonUnavailable},
		"bad-quarantine-action":      {func(v *grouppipeline.DeliveryView) { v.CounterReplayQuarantines[0].Action = "retry-now" }, runtimeview.ReasonUnavailable},
		"negative-samples":           {func(v *grouppipeline.DeliveryView) { v.CounterReplayQuarantines[0].PendingSamples = -1 }, runtimeview.ReasonUnavailable},
		"unsafe-backlog-destination": {func(v *grouppipeline.DeliveryView) { v.Backlog[0].TableName = "/private/secret.sqlite" }, runtimeview.ReasonUnavailable},
		"negative-backlog":           {func(v *grouppipeline.DeliveryView) { v.Backlog[0].Pending = -1 }, runtimeview.ReasonUnavailable},
		"duplicate-quarantine": {func(v *grouppipeline.DeliveryView) {
			v.CounterReplayQuarantines = append(v.CounterReplayQuarantines, v.CounterReplayQuarantines[0])
		}, runtimeview.ReasonUnavailable},
		"unsafe-quarantine-revision": {func(v *grouppipeline.DeliveryView) { v.CounterReplayQuarantines[0].GroupRevision = "/private/revision" }, runtimeview.ReasonUnavailable},
		"foreign-quarantine-group":   {func(v *grouppipeline.DeliveryView) { v.CounterReplayQuarantines[0].GroupID = "other-group" }, runtimeview.ReasonForeign},
		"unsafe-backlog-schema":      {func(v *grouppipeline.DeliveryView) { v.Backlog[0].TableSchema = "password=secret" }, runtimeview.ReasonUnavailable},
		"current-destination-conflict": {func(v *grouppipeline.DeliveryView) {
			v.Backlog[0].GroupRevision = "A"
			v.Backlog[0].ConnectorRevision = "8"
		}, runtimeview.ReasonStale},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			saved, snapshot, scope, raw, _ := runtimeDeliveryProjectionFixture()
			raw.CounterReplayQuarantines = []grouppipeline.CounterReplayQuarantine{{WorkspaceID: scope.WorkspaceID, GroupID: saved.ID, GroupRevision: "A0", PendingSamples: 7, Reason: "counter-replay-proof-missing", Action: grouppipeline.CounterReplayRepairAction}}
			raw.Backlog = []grouppipeline.BacklogView{{GroupRevision: "A0", ConnectorID: "C1", ConnectorRevision: "7", TableName: "old_table", Pending: 3}}
			before, _ := json.Marshal(raw)
			tc.mutate(raw)
			mutated, _ := json.Marshal(raw)
			require.NotEqual(t, before, mutated)
			_, err := projectRuntimeRuleDelivery(raw, saved, snapshot, scope, nil)
			require.True(t, runtimeview.IsContractError(err, tc.reason), "%v", err)
			after, _ := json.Marshal(raw)
			require.JSONEq(t, string(mutated), string(after), "read projection must not mutate borrowed authority values")
		})
	}
}

func TestRuntimeRuleDeliveryMissingProofRejectsUnexpectedCurrentEvidence(t *testing.T) {
	saved, _, scope, raw, _ := runtimeDeliveryProjectionFixture()
	_, err := projectRuntimeRuleDelivery(raw, saved, nil, scope, nil)
	require.True(t, runtimeview.IsContractError(err, runtimeview.ReasonStale))
	raw.RevisionStages = nil
	raw.LastSQLCommittedEffect = &groupdelivery.CommittedEffect{GroupRevision: "A"}
	_, err = projectRuntimeRuleDelivery(raw, saved, nil, scope, nil)
	require.True(t, runtimeview.IsContractError(err, runtimeview.ReasonStale))
}
