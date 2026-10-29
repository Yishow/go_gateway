package handlers

import (
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/workspace"
)

// projectRuntimeSourceDelivery consumes only the source proof from the same delivery read.
// A producer clock crossing a boundary may leave these facts unavailable while current stages remain valid.
func projectRuntimeSourceDelivery(raw *grouppipeline.SourceEvidenceView, frozen *workspace.WriteGroupAppliedSnapshot, current *runtimeview.CurrentDelivery) error {
	if raw == nil || raw.State == grouppipeline.SourceEvidenceUnavailable {
		return nil
	}
	group := frozen.Group
	layout := group.EffectiveLayoutVersion()
	if group.RecordingProfile == nil || (layout != workspace.WriteGroupLayoutManagedLongV1 && layout != workspace.WriteGroupLayoutManagedWideProfileV1) {
		return nil
	}
	if raw.State != grouppipeline.SourceEvidenceAvailable || raw.UnavailableReason != "" || !runtimeDeliveryIDs(raw.GroupRevision, raw.ConnectorID, raw.ConnectorRevision) {
		return runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	if raw.GroupRevision != frozen.AppliedRevision || raw.ConnectorID != group.Destination.ConnectorID || raw.ConnectorRevision != group.Destination.ConnectorRevision {
		return runtimeEvidenceError(runtimeview.ReasonStale)
	}
	members := make(map[[3]string]workspace.WriteGroupMember, len(group.Members))
	for _, m := range group.Members {
		members[[3]string{m.DeviceID, m.PointID, m.TagID}] = m
	}
	sources := make([]runtimeview.SourceEvidence, 0, len(raw.Sources))
	seen := make(map[[3]string]bool, len(raw.Sources))
	for _, s := range raw.Sources {
		if !runtimeDeliveryIDs(s.DeviceID, s.PointID, s.TagID) || (s.CircuitID != "" && !runtimeDeliveryIDs(s.CircuitID)) {
			return runtimeEvidenceError(runtimeview.ReasonUnavailable)
		}
		key := [3]string{s.DeviceID, s.PointID, s.TagID}
		member, owned := members[key]
		if !owned || (s.CircuitID != "" && s.CircuitID != member.CircuitID) {
			return runtimeEvidenceError(runtimeview.ReasonForeign)
		}
		if seen[key] {
			return runtimeEvidenceError(runtimeview.ReasonUnavailable)
		}
		seen[key] = true
		latest, err := runtimeSourceDeliveryResult(s.Latest, frozen, member, current.LastCommit != nil)
		if err != nil {
			return err
		}
		committed, err := runtimeSourceDeliveryResult(s.LastSQLCommitted, frozen, member, current.LastCommit != nil)
		if err != nil {
			return err
		}
		if s.LastSQLCommitted != nil && (committed == nil || committed.Receipt == nil) {
			return runtimeEvidenceError(runtimeview.ReasonUnavailable)
		}
		sources = append(sources, runtimeview.SourceEvidence{DeviceID: s.DeviceID, PointID: s.PointID, TagID: s.TagID, Latest: latest, LastCommitted: committed})
	}
	current.Sources, current.SourceState = sources, "available"
	return nil
}

func runtimeSourceDeliveryResult(raw *groupdelivery.SourceRecord, frozen *workspace.WriteGroupAppliedSnapshot, member workspace.WriteGroupMember, groupCommit bool) (*runtimeview.SourceResult, error) {
	if raw == nil {
		return nil, nil
	}
	if !raw.Available || raw.VerifiedSourceRecord == nil {
		return &runtimeview.SourceResult{State: "unavailable"}, nil
	}
	r := raw.VerifiedSourceRecord
	if raw.Cause != "" || !runtimeDeliveryIDs(r.GroupRevision, r.ConnectorRevision, r.RecordID, r.EffectKey, r.PayloadDigest) {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	if r.GroupRevision != frozen.AppliedRevision || r.ConnectorRevision != frozen.Group.Destination.ConnectorRevision {
		return nil, runtimeEvidenceError(runtimeview.ReasonStale)
	}
	profile := frozen.Group.RecordingProfile
	if r.ProfileKind != string(profile.Kind) || (profile.Kind != workspace.WriteGroupProfileInstantaneous && profile.Kind != workspace.WriteGroupProfileCumulative) {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	state := "unknown"
	switch r.Status {
	case groupdelivery.SourceStatusOK:
		state = "available"
	case groupdelivery.SourceStatusPartial:
		state = "partial"
	case groupdelivery.SourceStatusMissing:
		state = "missing"
	case groupdelivery.SourceStatusBaseline, groupdelivery.SourceStatusBad, groupdelivery.SourceStatusStale, groupdelivery.SourceStatusInvalid, groupdelivery.SourceStatusUncertain, groupdelivery.SourceStatusGap:
	default:
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	switch raw.Delivery {
	case groupdelivery.SourceDeliveryQueued, groupdelivery.SourceDeliveryRetrying, groupdelivery.SourceDeliveryBlocked, groupdelivery.SourceDeliveryQuarantined, groupdelivery.SourceDeliveryUnknown, groupdelivery.SourceDeliveryCommitted, groupdelivery.SourceDeliverySkipped:
	default:
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	start, end := runtimeSourceDeliveryTime(r.WindowStart), runtimeSourceDeliveryTime(r.WindowEnd)
	if start == nil || end == nil || !r.WindowStart.Before(r.WindowEnd) {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	missing := r.Status == groupdelivery.SourceStatusMissing || r.Status == groupdelivery.SourceStatusBad || r.Status == groupdelivery.SourceStatusStale || r.Status == groupdelivery.SourceStatusInvalid
	reading := r.Status == groupdelivery.SourceStatusOK || r.Status == groupdelivery.SourceStatusPartial || r.Status == groupdelivery.SourceStatusBaseline
	if (missing && r.RawPresent) || (reading && !r.RawPresent) {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	result := &runtimeview.SourceResult{State: state, WindowStart: start, WindowEnd: end}
	if profile.Kind == workspace.WriteGroupProfileCumulative {
		if r.DeltaKnown == nil || r.Aggregation != nil || (*r.DeltaKnown && (missing || r.Status == groupdelivery.SourceStatusBaseline || r.Status == groupdelivery.SourceStatusUncertain)) {
			return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
		}
		known := *r.DeltaKnown
		result.DeltaKnown = &known
		if !known {
			result.State = "unknown"
		}
	} else if r.DeltaKnown != nil {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	if err := runtimeSourceDeliveryCoverage(r, profile, member, result); err != nil {
		return nil, err
	}
	committed := raw.Delivery == groupdelivery.SourceDeliveryCommitted
	if committed != (r.CommittedAt != nil) {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	if committed {
		if !groupCommit || runtimeSourceDeliveryTime(*r.CommittedAt) == nil {
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

func runtimeSourceDeliveryCoverage(r *groupdelivery.VerifiedSourceRecord, profile *workspace.WriteGroupRecordingProfile, member workspace.WriteGroupMember, result *runtimeview.SourceResult) error {
	expected := profile.Kind == workspace.WriteGroupProfileInstantaneous && profile.EffectiveSampleReduction() == workspace.WriteGroupSampleReductionMean && string(member.ValueRole) == "gauge"
	a := r.Aggregation
	if (a != nil) != expected {
		return runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	if a == nil {
		return nil
	}
	exact := a.Basis == string(measurement.MeanBasisExact) && a.ResultType == string(measurement.ExactDecimal)
	floating := a.Basis == string(measurement.MeanBasisFloat64) && a.ResultType == string(measurement.ExactFloat64)
	if a.Method != measurement.MeanAggregationMethod || a.Version != measurement.MeanAggregationVersion || (!exact && !floating) || a.SourceType == "" || a.SourceType != string(member.SourceValueType) || a.ResultType != string(member.ResultValueType) || a.ValidCount < 0 || a.RejectedCount < 0 {
		return runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	reading := r.Status == groupdelivery.SourceStatusOK || r.Status == groupdelivery.SourceStatusPartial
	if (a.ValidCount == 0) != (a.Coverage == "missing") || (a.ValidCount > 0) != reading {
		return runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	switch a.Coverage {
	case "partial", "unknown":
		coverage := a.Coverage
		result.Coverage, result.State = &coverage, coverage
	case "missing":
		result.State = "missing" // The existing DTO deliberately has no missing coverage enum.
	default:
		return runtimeEvidenceError(runtimeview.ReasonUnavailable) // Complete has no producer proof.
	}
	return nil
}

func runtimeSourceDeliveryTime(at time.Time) *string {
	if at.IsZero() {
		return nil
	}
	value := at.UTC().Format(time.RFC3339Nano)
	if _, err := time.Parse(time.RFC3339Nano, value); err != nil {
		return nil
	}
	return &value
}
