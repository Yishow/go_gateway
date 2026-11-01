package handlers

import (
	"encoding/hex"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"
	"time"
)

const runtimeStatusNotDue = "not_due"
const runtimeStatusMissing = "missing"
const runtimeStatusPartial = "partial"
const runtimeStatusUnknown = "unknown"
const runtimeMethodDelta = "delta"

func runtimeMixedSourceResult(raw *groupdelivery.SourceRecord, frozen *workspace.WriteGroupAppliedSnapshot, member workspace.WriteGroupMember, groupCommit bool) (*runtimeview.SourceResult, error) {
	r := raw.VerifiedSourceRecord
	c := r.Calculation
	if c == nil || r.ProfileKind != runtimeRecordingMixed || member.Cadence == nil || frozen.Group.RecordingProfile == nil || frozen.Group.RecordingProfile.ResultShape != workspace.WriteGroupResultShapeMemberCalculationV1 {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	duration, err := member.Cadence.Duration()
	if err != nil || c.Method != string(member.Calculation) || c.IntervalSeconds != int64(duration/time.Second) || c.Alignment != member.Cadence.Alignment || c.ResultShape != string(frozen.Group.RecordingProfile.ResultShape) || c.SourceType != string(member.SourceValueType) || c.ResultType != string(member.ResultValueType) {
		return nil, runtimeEvidenceError(runtimeview.ReasonStale)
	}
	if r.GroupRevision != frozen.AppliedRevision || r.ConnectorRevision != frozen.Group.Destination.ConnectorRevision {
		return nil, runtimeEvidenceError(runtimeview.ReasonStale)
	}
	if !runtimeDeliveryIDs(r.GroupRevision, r.ConnectorRevision, r.RecordID, c.OutcomeDigest) || raw.Cause != "" || !legacyDeliveryMixedRecordSafe(raw) {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	wide := frozen.Group.EffectiveLayoutVersion() == workspace.WriteGroupLayoutManagedWideProfileV1
	expected, err := snapshot.MixedDispatchRecordID(frozen.WorkspaceID, frozen.Group.ID, frozen.AppliedRevision, c.DueAt)
	if !wide {
		if c.AlignedStart == nil || c.Status == runtimeStatusNotDue {
			return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
		}
		expected, err = snapshot.ProfileRecordID(frozen.WorkspaceID, frozen.Group.ID, frozen.AppliedRevision, groupdelivery.SourceEntityKey(member.DeviceID, member.PointID, member.TagID), runtimeRecordingMixed, *c.AlignedStart)
	}
	if err != nil {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	if r.RecordID != expected {
		return nil, runtimeEvidenceError(runtimeview.ReasonForeign)
	}
	projected := &runtimeview.CalculationResult{ResultShape: c.ResultShape, Method: c.Method, IntervalSeconds: c.IntervalSeconds, Alignment: c.Alignment, DueAt: c.DueAt.UTC().Format(time.RFC3339Nano), NextDue: mixedEvidenceTime(c.NextDue), OutcomeDigest: c.OutcomeDigest, OutcomeKind: c.OutcomeKind, Status: c.Status, Reason: runtimeEvidenceString(c.Reason), ResultPresent: c.ResultPresent, RowEmitted: c.RowEmitted, RawPresent: c.RawPresent, SourceType: c.SourceType, ResultType: c.ResultType, AlignedStart: mixedEvidenceTime(c.AlignedStart), AlignedEnd: mixedEvidenceTime(c.AlignedEnd), ActualStart: mixedEvidenceTime(c.ActualStart), ActualEnd: mixedEvidenceTime(c.ActualEnd), Rounded: c.Rounded}
	result := &runtimeview.SourceResult{State: runtimeStatusUnknown, Calculation: projected, WindowStart: projected.ActualStart, WindowEnd: projected.ActualEnd, DeltaKnown: r.DeltaKnown}
	switch {
	case c.Status == runtimeStatusNotDue:
		result.State = runtimeStatusNotDue
	case c.ResultPresent:
		result.State = runtimeEvidenceAvailable
	case c.Status == runtimeStatusMissing && r.Status == groupdelivery.SourceStatusMissing:
		result.State = runtimeStatusMissing
	}
	if r.Aggregation != nil {
		coverage := r.Aggregation.Coverage
		switch coverage {
		case runtimeStatusPartial:
			result.State = runtimeStatusPartial
			result.Coverage = &coverage
		case runtimeStatusUnknown:
			result.Coverage = &coverage
		case runtimeStatusMissing:
		default:
			return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
		}
	}
	if c.Method == runtimeMethodDelta && r.DeltaKnown != nil && !*r.DeltaKnown && result.State != runtimeStatusNotDue {
		result.State = runtimeStatusUnknown
	}
	committed := raw.Delivery == groupdelivery.SourceDeliveryCommitted
	if committed != (r.CommittedAt != nil) {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	if committed {
		if !groupCommit {
			return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
		}
		receipt, err := runtimeDeliveryReceipt(&groupdelivery.CommittedEffect{GroupRevision: r.GroupRevision, ConnectorRevision: r.ConnectorRevision, RecordID: r.RecordID, EffectKey: r.EffectKey, PayloadDigest: r.PayloadDigest, CommittedAt: *r.CommittedAt}, 1, frozen)
		if err != nil {
			return nil, err
		}
		result.Receipt = receipt
	}
	return result, nil
}
func mixedEvidenceTime(value *time.Time) *string {
	if value == nil {
		return nil
	}
	return runtimeSourceDeliveryTime(*value)
}

func legacyDeliveryMixedRecordSafe(raw *groupdelivery.SourceRecord) bool {
	r := raw.VerifiedSourceRecord
	if r == nil || r.Calculation == nil || r.ProfileKind != runtimeRecordingMixed || !raw.Available {
		return false
	}
	c := r.Calculation
	if c.Rounded != nil && c.Method != "mean" {
		return false
	}
	for _, at := range []*time.Time{&c.DueAt, c.NextDue, c.AlignedStart, c.AlignedEnd, c.ActualStart, c.ActualEnd} {
		if at == nil {
			continue
		}
		if _, err := snapshot.CheckedWindowAdd(*at, 0, "source evidence"); err != nil {
			return false
		}
	}
	if c.NextDue != nil && !c.NextDue.After(c.DueAt) {
		return false
	}
	hash, err := hex.DecodeString(c.OutcomeDigest)
	if err != nil || len(hash) != 32 || !legacyDeliveryChoice(c.Status, "ok", runtimeStatusMissing, "bad", runtimeStatusStale, "invalid", runtimeStatusNotDue) || !schema.DataType(c.SourceType).IsValid() || !legacyDeliveryChoice(c.ResultType, "bool", "text", "int64", "uint64", "decimal", "float64") {
		return false
	}
	if c.ResultPresent && c.Status != "ok" {
		return false
	}
	if !legacyDeliveryChoice(string(raw.Delivery), "queued", "retrying", "blocked", "quarantined", runtimeStatusUnknown, "sql_committed", "operator_skipped", "window_skipped") {
		return false
	}
	if c.ResultShape != "member-calculation-v1" || !legacyDeliveryChoice(c.Method, runtimeMethodLatest, "mean", runtimeMethodDelta) || c.Alignment != "utc_epoch" || c.DueAt.IsZero() || !runtimeDeliveryIDs(r.RecordID, r.GroupRevision, r.ConnectorRevision, c.OutcomeDigest) || c.RawPresent != r.RawPresent {
		return false
	}
	if _, err := snapshot.CheckedCadenceDuration(c.IntervalSeconds, "second"); err != nil {
		return false
	}
	if !legacyDeliveryChoice(c.OutcomeKind, "row", "skipped", "no_data") || c.RowEmitted != (c.OutcomeKind == "row") {
		return false
	}
	if c.RowEmitted {
		if !runtimeDeliveryIDs(r.EffectKey, r.PayloadDigest) || raw.Delivery == groupdelivery.SourceDeliveryWindowSkipped {
			return false
		}
	} else if r.EffectKey != "" || r.PayloadDigest != "" || r.CommittedAt != nil || raw.Delivery != groupdelivery.SourceDeliveryWindowSkipped {
		return false
	}
	if c.Status == runtimeStatusNotDue {
		if c.AlignedStart != nil || c.AlignedEnd != nil || c.ResultPresent || c.RawPresent || c.ActualStart != nil || c.ActualEnd != nil {
			return false
		}
	} else if c.ActualStart == nil || c.ActualEnd == nil || c.AlignedStart == nil || c.AlignedEnd == nil || !c.ActualStart.Before(*c.ActualEnd) || c.ActualStart.Before(*c.AlignedStart) || c.ActualEnd.After(*c.AlignedEnd) {
		return false
	}
	if c.Method == runtimeMethodDelta {
		if r.DeltaKnown == nil || *r.DeltaKnown != c.ResultPresent {
			return false
		}
	} else if r.DeltaKnown != nil {
		return false
	}
	c.Reason = legacyDeliverySourceReason(c.Reason)
	r.Reason = legacyDeliverySourceReason(r.Reason)
	return true
}
