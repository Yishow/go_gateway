package handlers

import (
	"math"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/workspace"
)

func projectRuntimeRuleDelivery(raw *grouppipeline.DeliveryView, saved *workspace.WriteGroup, snapshot *workspace.WriteGroupAppliedSnapshot, scope runtimeview.Scope, intake *runtimeview.IntakeEvidence) (*runtimeview.DeliveryEvidence, error) {
	if raw == nil {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	if raw.GroupID != saved.ID {
		return nil, runtimeEvidenceError(runtimeview.ReasonForeign)
	}
	revision := ""
	if snapshot != nil {
		revision = snapshot.AppliedRevision
	}
	result := &runtimeview.DeliveryEvidence{OldestPendingScope: "group-all-revisions", Backlog: []runtimeview.Backlog{}}
	age := raw.OldestPendingSeconds
	if math.IsNaN(age) || math.IsInf(age, 0) || age < 0 || age > 9007199254740991 {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	// The producer has no presence bit: zero cannot distinguish a new row from no pending row.
	if age > 0 {
		result.OldestPendingSeconds = &age
	}
	if raw.CounterReplayQuarantines != nil {
		quarantines := make([]runtimeview.CounterReplayQuarantine, 0, len(raw.CounterReplayQuarantines))
		seen := map[string]bool{}
		for _, q := range raw.CounterReplayQuarantines {
			if q.WorkspaceID != scope.WorkspaceID || q.GroupID != saved.ID {
				return nil, runtimeEvidenceError(runtimeview.ReasonForeign)
			}
			reason := runtimeview.ParseReason(q.Reason)
			if !runtimeDeliveryIDs(q.GroupRevision) || !runtimeDeliveryCount(q.PendingSamples) || seen[q.GroupRevision] || q.Action != grouppipeline.CounterReplayRepairAction || (reason != runtimeview.ReasonCounterMissing && reason != runtimeview.ReasonCounterConflict && reason != runtimeview.ReasonCounterUnsupported) {
				return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
			}
			seen[q.GroupRevision] = true
			quarantines = append(quarantines, runtimeview.CounterReplayQuarantine{WorkspaceID: scope.WorkspaceID, GroupID: saved.ID, GroupRevision: q.GroupRevision, PendingSamples: int64(q.PendingSamples), Reason: reason, Action: q.Action})
		}
		result.CounterReplayQuarantines = &quarantines
	}
	for _, b := range raw.Backlog {
		if !runtimeDeliveryIDs(b.GroupRevision, b.ConnectorID, b.ConnectorRevision, b.TableName) || (b.TableSchema != "" && !runtimeDeliveryIDs(b.TableSchema)) || !runtimeDeliveryCount(b.Pending) {
			return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
		}
		if revision != "" && b.GroupRevision == revision {
			destination := snapshot.Group.Destination
			if b.ConnectorID != destination.ConnectorID || b.ConnectorRevision != destination.ConnectorRevision || b.TableSchema != destination.TableSchema || b.TableName != destination.TableName {
				return nil, runtimeEvidenceError(runtimeview.ReasonStale)
			}
			// Current output reasons are not intake reasons or historical backlog.
			continue
		}
		reasons := make([]runtimeview.Reason, 0, len(b.ErrorCodes))
		for _, code := range b.ErrorCodes {
			reasons = append(reasons, runtimeview.ParseReason(code))
		}
		result.Backlog = append(result.Backlog, runtimeview.Backlog{Binding: runtimeview.Binding{WorkspaceID: scope.WorkspaceID, WorkspaceRevision: scope.WorkspaceRevision, EntityRevision: b.GroupRevision}, Destination: runtimeview.Destination{ConnectorID: b.ConnectorID, ConnectorRevision: b.ConnectorRevision, Schema: runtimeEvidenceString(b.TableSchema), Table: b.TableName, Kind: "unavailable"}, Pending: int64(b.Pending), Reasons: reasons})
	}
	if revision == "" {
		if raw.RevisionStages != nil || raw.LastSQLCommittedEffect != nil {
			return nil, runtimeEvidenceError(runtimeview.ReasonStale)
		}
		return result, nil
	}
	if raw.RevisionStages == nil {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	if raw.RevisionStages.GroupRevision != revision {
		return nil, runtimeEvidenceError(runtimeview.ReasonStale)
	}
	stages := raw.RevisionStages.Stages
	for _, count := range []int{stages.Collecting, stages.Queued, stages.Retrying, stages.Blocked, stages.Quarantined, stages.Unknown, stages.SQLCommitted, stages.Skipped} {
		if !runtimeDeliveryCount(count) {
			return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
		}
	}
	receipt, err := runtimeDeliveryReceipt(raw.LastSQLCommittedEffect, stages.SQLCommitted, snapshot)
	if err != nil {
		return nil, err
	}
	quota, err := runtimeDeliveryQuota(raw.Quota)
	if err != nil {
		return nil, err
	}
	devices := map[string]bool{}
	for _, member := range snapshot.Group.Members {
		devices[member.DeviceID] = true
	}
	result.Current = &runtimeview.CurrentDelivery{Binding: runtimeview.Binding{WorkspaceID: scope.WorkspaceID, WorkspaceRevision: scope.WorkspaceRevision, EntityRevision: revision}, TotalsScope: "group", Shared: len(devices) > 1, Intake: intake,
		Stages:     runtimeview.StageCounts{Collecting: int64(stages.Collecting), Queued: int64(stages.Queued), Retrying: int64(stages.Retrying), Blocked: int64(stages.Blocked), Quarantined: int64(stages.Quarantined), Unknown: int64(stages.Unknown), SQLCommitted: int64(stages.SQLCommitted), Skipped: int64(stages.Skipped)},
		LastCommit: receipt, Sources: []runtimeview.SourceEvidence{}, SourceState: "unavailable", Quota: quota}
	// Storage remains unavailable: editable-group storage is not frozen lineage proof.
	if err := projectRuntimeSourceDelivery(raw.SourceEvidence, snapshot, result.Current); err != nil {
		return nil, err
	}
	return result, nil
}

func runtimeDeliveryQuota(raw groupdelivery.QuotaStatusView) (*runtimeview.QuotaEvidence, error) {
	if !raw.Configured {
		if (raw.State != "" && raw.State != "unconfigured") || raw.Scope != "" || raw.UsedBytes != 0 || raw.MaxBytes != 0 || raw.IntakeRefused {
			return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
		}
		return nil, nil
	}
	if (raw.Scope != groupdelivery.QuotaScopeGroup && raw.Scope != groupdelivery.QuotaScopeGlobal) ||
		(raw.State != groupdelivery.QuotaOK && raw.State != groupdelivery.QuotaWarning && raw.State != groupdelivery.QuotaHard) ||
		raw.UsedBytes < 0 || raw.UsedBytes > 9007199254740991 || raw.MaxBytes <= 0 || raw.MaxBytes > 9007199254740991 ||
		raw.IntakeRefused != (raw.State == groupdelivery.QuotaHard) {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	// Measured group/global usage spans all revisions, not the current revision or device.
	return &runtimeview.QuotaEvidence{Scope: string(raw.Scope), State: raw.State, UsedBytes: raw.UsedBytes, MaxBytes: raw.MaxBytes, IntakeRefused: raw.IntakeRefused}, nil
}

func runtimeDeliveryReceipt(effect *groupdelivery.CommittedEffect, committed int, snapshot *workspace.WriteGroupAppliedSnapshot) (*runtimeview.Receipt, error) {
	if (committed > 0) != (effect != nil) {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	if effect == nil {
		return nil, nil
	}
	if effect.GroupRevision != snapshot.AppliedRevision || effect.ConnectorRevision != snapshot.Group.Destination.ConnectorRevision {
		return nil, runtimeEvidenceError(runtimeview.ReasonStale)
	}
	if !runtimeDeliveryIDs(effect.RecordID, effect.EffectKey, effect.PayloadDigest) || effect.CommittedAt.IsZero() {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	return &runtimeview.Receipt{GroupRevision: effect.GroupRevision, ConnectorRevision: effect.ConnectorRevision, RecordID: effect.RecordID, EffectKey: effect.EffectKey, PayloadDigest: effect.PayloadDigest, CommittedAt: *runtimeEvidenceTime(effect.CommittedAt), Provenance: "runtime"}, nil
}
func runtimeDeliveryIDs(values ...string) bool {
	for _, value := range values {
		if value == "" || runtimeSafeID(value) != value {
			return false
		}
	}
	return true
}
func runtimeDeliveryCount(value int) bool { return value >= 0 && int64(value) <= 9007199254740991 }
