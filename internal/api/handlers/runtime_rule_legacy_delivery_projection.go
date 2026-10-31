package handlers

import (
	"encoding/json"
	"slices"
	"strings"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/schema"
)

// projectLegacyRuntimeDelivery keeps the native delivery wire shape, not its raw
// diagnostic strings. It neither reinterprets revisions nor authorizes recovery.
func projectLegacyRuntimeDelivery(raw *grouppipeline.DeliveryView, expectedGroupID string) (*grouppipeline.DeliveryView, error) {
	unavailable := runtimeEvidenceError(runtimeview.ReasonUnavailable)
	if raw == nil || raw.GroupID != expectedGroupID || !runtimeDeliveryIDs(raw.GroupID) {
		return nil, unavailable
	}
	// Detach every nested pointer/slice before normalization. This copies only the
	// existing value-free DTO; no store, configuration or frozen payload is read.
	encoded, err := json.Marshal(raw)
	if err != nil {
		return nil, unavailable
	}
	var result grouppipeline.DeliveryView
	if json.Unmarshal(encoded, &result) != nil {
		return nil, unavailable
	}
	if !legacyDeliveryChoice(result.Intake.State, "", "active", "retiring", "blocked", "not_running") {
		return nil, unavailable
	}
	result.Intake.Reason = legacyDeliveryReason(result.Intake.Reason)
	if r := result.RevisionStages; r != nil && !runtimeDeliveryIDs(r.GroupRevision) {
		return nil, unavailable
	}
	if r := result.LastSQLCommittedEffect; r != nil && !legacyDeliveryReceiptSafe(r) {
		return nil, unavailable
	}
	for i := range result.Backlog {
		b := &result.Backlog[i]
		if !runtimeDeliveryIDs(b.GroupRevision, b.ConnectorID, b.ConnectorRevision, b.TableName) || !legacyDeliveryOptionalID(b.TableSchema) {
			return nil, unavailable
		}
		for j := range b.ErrorCodes {
			b.ErrorCodes[j] = legacyDeliveryReason(b.ErrorCodes[j])
		}
	}
	for i := range result.Attention {
		a := &result.Attention[i]
		if !runtimeDeliveryIDs(a.EffectKey, a.PayloadDigest, a.GroupRevision, a.ConnectorID, a.ConnectorRevision, a.TableName) || !legacyDeliveryOptionalID(a.TableSchema) || !legacyDeliveryChoice(a.State, "blocked", "quarantined", "unknown") {
			return nil, unavailable
		}
		a.ErrorCode = legacyDeliveryReason(a.ErrorCode)
	}
	for i := range result.CounterReplayQuarantines {
		q := &result.CounterReplayQuarantines[i]
		if !runtimeDeliveryIDs(q.WorkspaceID, q.GroupID, q.GroupRevision) || q.GroupID != expectedGroupID || !legacyDeliveryChoice(q.Reason, string(runtimeview.ReasonCounterMissing), string(runtimeview.ReasonCounterConflict), string(runtimeview.ReasonCounterUnsupported)) || q.Action != grouppipeline.CounterReplayRepairAction {
			return nil, unavailable
		}
	}
	for i := range result.RecentBucketIssues {
		b := &result.RecentBucketIssues[i]
		if !runtimeDeliveryIDs(b.GroupRevision) || !legacyDeliveryChoice(b.Kind, "skipped", "no_data", "truncated") {
			return nil, unavailable
		}
		for j, cause := range b.Causes {
			if !legacyDeliveryChoice(cause, "missing", "bad", "stale", "invalid", "no_data", "unavailable", "boundary") {
				b.Causes[j] = "unavailable"
			}
		}
	}
	q := &result.Quota
	if !legacyDeliveryChoice(q.State, "", "unconfigured", "ok", "warning", "hard_limit") || !legacyDeliveryChoice(string(q.Scope), "", "global", "group") {
		return nil, unavailable
	}
	if !legacyDeliveryChoice(q.LossRiskNotice, "", "new readings for this scope are not recorded until space is available; already accepted data is kept", "capacity is running low; intake will be refused at the hard limit") {
		q.LossRiskNotice = "unavailable"
	}
	if s := result.Storage; s != nil {
		if !legacyDeliveryChoice(s.RepairState, "", "blocked") {
			return nil, unavailable
		}
		s.RepairReason = legacyDeliveryReason(s.RepairReason)
		for i := range s.Events {
			e := &s.Events[i]
			if !legacyDeliveryOptionalID(e.GroupID) || !legacyDeliveryOptionalID(e.OperationID) || !legacyDeliveryChoice(e.Kind, "confirmed", "recreated_by_confirmation", "reconstructed", "repair_blocked") {
				return nil, unavailable
			}
			e.Code, e.Reason = legacyDeliveryReason(e.Code), legacyDeliveryReason(e.Reason)
		}
	}
	if !legacyDeliverySourceSafe(result.SourceEvidence) {
		return nil, unavailable
	}
	return &result, nil
}
func legacyDeliveryChoice(value string, allowed ...string) bool {
	return slices.Contains(allowed, value)
}
func legacyDeliveryOptionalID(value string) bool { return value == "" || runtimeDeliveryIDs(value) }
func legacyDeliveryReceiptSafe(r *groupdelivery.CommittedEffect) bool {
	return runtimeDeliveryIDs(r.GroupRevision, r.ConnectorRevision, r.RecordID, r.EffectKey, r.PayloadDigest)
}

func legacyDeliverySourceSafe(s *grouppipeline.SourceEvidenceView) bool {
	if s == nil {
		return true
	}
	if !legacyDeliveryChoice(s.State, grouppipeline.SourceEvidenceAvailable, grouppipeline.SourceEvidenceUnavailable) || !runtimeDeliveryIDs(s.GroupRevision) || !legacyDeliveryOptionalID(s.ConnectorID) || !legacyDeliveryOptionalID(s.ConnectorRevision) {
		return false
	}
	if s.State == grouppipeline.SourceEvidenceAvailable && (!runtimeDeliveryIDs(s.ConnectorID, s.ConnectorRevision) || s.UnavailableReason != "") {
		return false
	}
	if !legacyDeliveryChoice(s.UnavailableReason, "", "not_long_layout", grouppipeline.SourceEvidenceNoSnapshot, grouppipeline.SourceEvidenceRevisionMismatch) {
		return false
	}
	for i := range s.Sources {
		source := &s.Sources[i]
		if !runtimeDeliveryIDs(source.DeviceID, source.PointID, source.TagID) || !legacyDeliveryOptionalID(source.CircuitID) || !legacyDeliverySourceRecordSafe(source.Latest) || !legacyDeliverySourceRecordSafe(source.LastSQLCommitted) {
			return false
		}
	}
	return true
}
func legacyDeliverySourceRecordSafe(r *groupdelivery.SourceRecord) bool {
	if r == nil {
		return true
	}
	if !legacyDeliveryChoice(string(r.Delivery), "", "queued", "retrying", "blocked", "quarantined", "unknown", "sql_committed", "operator_skipped", "window_skipped") || !legacyDeliveryChoice(string(r.Cause), "", "payload_invalid", "digest_mismatch", "identity_mismatch", "status_unknown", "window_invalid", "calculation_projection_unsupported") {
		return false
	}
	v := r.VerifiedSourceRecord
	if v == nil {
		return !r.Available
	}
	if v.Calculation != nil || v.ProfileKind == "mixed" {
		return legacyDeliveryMixedRecordSafe(r)
	}
	if !r.Available || !runtimeDeliveryIDs(v.GroupRevision, v.ConnectorRevision, v.RecordID, v.EffectKey, v.PayloadDigest) || !legacyDeliveryChoice(v.ProfileKind, "instantaneous", "cumulative") || !legacyDeliveryChoice(string(v.Status), "ok", "partial", "baseline", "missing", "bad", "stale", "invalid", "reset_or_uncertain", "gap") {
		return false
	}
	v.Reason = legacyDeliverySourceReason(v.Reason)
	if a := v.Aggregation; a != nil {
		if a.Method != measurement.MeanAggregationMethod || a.Version != measurement.MeanAggregationVersion || !legacyDeliveryChoice(a.Basis, string(measurement.MeanBasisExact), string(measurement.MeanBasisFloat64)) || !schema.DataType(a.SourceType).IsValid() || !legacyDeliveryChoice(a.ResultType, "decimal", "float64") || !legacyDeliveryChoice(a.Coverage, "partial", "unknown", "missing") {
			return false
		}
		a.Reason = legacyDeliverySourceReason(a.Reason)
	}
	return true
}

// These code vocabularies mirror the existing producers: sender/reconcile,
// snapshot builder, recordingplan storage restore and groupdelivery source proof.
// Unknown text is never promoted to a success/empty reason. No prefix wildcard.
var legacyDeliveryReasons = strings.Fields(`
 unavailable none snapshot-unavailable snapshot-config-invalid destination-unavailable destination-missing
 destination-revision-changed destination-unsupported destination-table-missing destination-table-unavailable
 tag-unavailable interval-invalid counter-semantics-unavailable
 group-missing sink-missing group-not-applied member-identity-incomplete tag-type-unsupported
 layout-blocked destination-scope-invalid restore-failed dedupe-unsupported
 insert-failed storage-restore-in-progress storage-recovery-unavailable commit-ambiguous
 destination-identity-conflict target-blocked destination-rejected-row target-unusable target-unavailable
 payload-digest-mismatch payload-undecodable schema_repair_blocked storage_owner_mismatch
 confirmation-missing ownership-mismatch repair-needs-confirmation connector-revision-changed table-occupied
 recovery-disabled recovery-policy-unconfirmed receipts-missing timescale-unavailable timescale-extension-moved
 permission-denied statement-failed verification-mismatch unknown-ddl-result repair_attempt_abandoned
 table-created-again table-rebuilt
`)

func legacyDeliveryReason(value string) string {
	if value == "" || slices.Contains(legacyDeliveryReasons, value) {
		return value
	}
	if reason := runtimeview.ParseReason(value); reason != runtimeview.ReasonUnavailable {
		return string(reason)
	}
	return "unavailable"
}

var legacyDeliverySourceReasons = strings.Fields(`
 reason_unavailable not_due no-samples type-mismatch counter-unavailable missing-value max-age-exceeded mapping-failed
 quality-bad quality-missing quality-stale quality-invalid quality-uncertain coverage-unverified bad invalid stale non-numeric
 first-reading binding-changed no-eligible-reading decrease-unproven contradictory-same-time non-monotonic-time
 mixed-numeric-basis rollover-evidence-insufficient rollover-range-invalid rollover-multiple-wraps-possible
 rollover-rate-bound-violated non-numeric-reading delta-not-representable
`)

func legacyDeliverySourceReason(value string) string {
	if value == "" || slices.Contains(legacyDeliverySourceReasons, value) {
		return value
	}
	return groupdelivery.SourceReasonUnavailable
}
