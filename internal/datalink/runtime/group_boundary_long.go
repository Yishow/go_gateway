package runtime

import (
	"context"
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"
)

// Activation codes specific to the long layout.
const (
	boundaryProfileUnsupported = "profile-unsupported"
	boundaryProfileInterval    = "profile-interval-mismatch"
	boundaryMemberMetadata     = "member-metadata-unsupported"
)

// newLongGroupBoundary activates an applied managed_long_v1 group. Each source
// is an independent record per window, written into the shared fixed table;
// nothing here depends on a per-source column. An instantaneous group records
// the newest eligible reading. A cumulative group additionally needs confirmed
// counter semantics for every member and compares readings through exact
// baselines; any other profile, or a missing profile, is refused so a group can
// never be written under the wrong semantics.
func newLongGroupBoundary(ctx context.Context, cfg GroupBoundaryConfig) (*GroupBoundary, error) {
	blocked := func(code string, issues ...dbtarget.LayoutIssue) error {
		return &GroupBoundaryError{Code: code, Issues: issues}
	}
	group := cfg.Group
	profile := group.RecordingProfile
	if profile == nil {
		return nil, blocked(boundaryProfileUnsupported)
	}
	var mode snapshot.Mode
	switch profile.Kind {
	case workspace.WriteGroupProfileInstantaneous:
		mode = snapshot.ModeInstantaneousProfile
	case workspace.WriteGroupProfileCumulative:
		mode = snapshot.ModeCumulativeProfile
	default:
		return nil, blocked(boundaryProfileUnsupported)
	}
	if profile.IntervalSeconds() != group.RowPolicy.IntervalSeconds {
		return nil, blocked(boundaryProfileInterval)
	}
	reduction := profile.EffectiveSampleReduction()
	if reduction == workspace.WriteGroupSampleReductionMean && profile.Kind != workspace.WriteGroupProfileInstantaneous {
		// A cumulative window counts its own interval; it never averages.
		return nil, blocked(boundaryMemberMetadata)
	}
	cumulative := mode == snapshot.ModeCumulativeProfile
	if cumulative && cfg.Ledger != nil && !ledgerPersistsBaselines(cfg.Ledger) {
		// A ledger that cannot carry the baseline with the closure would restart
		// every counter silently. Refuse rather than look durable.
		return nil, blocked(boundaryCounterDurable)
	}
	members := make(map[string]boundaryMember, len(group.Members))
	layoutMembers := make([]dbtarget.LongRowMember, 0, len(group.Members))
	snapshotMembers := make([]snapshot.Member, 0, len(group.Members))
	counters := make(map[string]aggregation.CounterBinding)
	for _, member := range group.Members {
		if member.DeviceID == "" || member.PointID == "" || member.TagID == "" {
			return nil, blocked(boundaryMemberIdentity)
		}
		tagType, ok := cfg.TagTypes[member.TagID]
		if !ok {
			return nil, blocked(boundaryTagTypeUnsupport)
		}
		kind, ok := measurement.ExactTypeForTag(tagType)
		if !ok {
			return nil, blocked(boundaryTagTypeUnsupport)
		}
		stored, memberReduction, code := longMemberShape(reduction, member, tagType, kind)
		if code != "" {
			return nil, blocked(code)
		}
		key := memberKey(member.DeviceID, member.PointID, member.TagID)
		layoutMember := dbtarget.LongRowMember{
			MemberKey: key, DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID,
			CircuitID: member.CircuitID, Type: stored,
			Reduction: string(memberReduction), ValueRole: string(member.ValueRole), SourceValueType: member.SourceValueType,
		}
		if member.MeasurementID != nil {
			layoutMember.MeasurementID = *member.MeasurementID
		}
		entry := boundaryMember{key: key, kind: kind}
		if !cumulative && member.ValueRole == workspace.WriteGroupValueRoleCounterRaw {
			// The latest raw counter value reuses the confirmed binding for its
			// identity check only: no baseline, delta or window is computed here.
			binding, code := counterBindingFor(group, member, key, kind, cfg.Counters)
			if code != "" {
				return nil, blocked(code)
			}
			entry.counter = &binding
		}
		if cumulative {
			binding, code := counterBindingFor(group, member, key, kind, cfg.Counters)
			if code != "" {
				return nil, blocked(code)
			}
			counters[key] = binding
			entry.counter, layoutMember.Counter, layoutMember.MeasurementID = &binding, &binding, binding.MeasurementID
		}
		members[key] = entry
		layoutMembers = append(layoutMembers, layoutMember)
		var maxAge time.Duration
		if member.MaxAgeSeconds != nil {
			maxAge = time.Duration(*member.MaxAgeSeconds) * time.Second
		}
		snapshotMembers = append(snapshotMembers, snapshot.Member{
			Key: key, Required: member.Required, MaxAge: maxAge,
			SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision,
			SampleReduction: memberReduction, SourceValueType: member.SourceValueType,
		})
	}
	// Readiness and the frozen proof judge the whole recording shape with this
	// same rule. Judging it here too means a reduction this runtime does not know
	// can never fall back to the legacy latest rule through a direct activation.
	if !workspace.LongRecordingShapeSupported(group, cfg.TagTypes) {
		return nil, blocked(boundaryProfileUnsupported)
	}
	layout, issues := dbtarget.NewLongRowLayout(dbtarget.LongRowSpec{
		Dialect: cfg.Dialect, Columns: cfg.Columns,
		WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: group.Revision,
		ProfileKind: string(profile.Kind), ProfileInterval: profile.Interval, ProfileUnit: string(profile.Unit),
		Members: layoutMembers,
	})
	if len(issues) > 0 {
		return nil, blocked(boundaryLayoutBlocked, issues...)
	}
	parts := boundaryParts{members: members, snapshotMembers: snapshotMembers, encoder: layout, mode: mode}
	if cumulative {
		parts.counters = counters
	}
	return assembleBoundary(ctx, cfg, parts)
}

// longMemberShape resolves one member's frozen recording metadata to the stored
// result type and the effective snapshot reduction. A new shape proves its
// source/result/role facts; anything partial or contradictory is refused rather
// than degraded to latest or to a raw sample.
func longMemberShape(
	reduction workspace.WriteGroupSampleReduction, member workspace.WriteGroupMember,
	tagType schema.DataType, source measurement.ExactType,
) (measurement.ExactType, snapshot.SampleReduction, string) {
	newShape := reduction == workspace.WriteGroupSampleReductionMean ||
		member.ValueRole != "" || member.SourceValueType != "" || member.ResultValueType != ""
	if !newShape {
		return source, snapshot.SampleReductionLatest, ""
	}
	switch member.ValueRole {
	case workspace.WriteGroupValueRoleGauge, workspace.WriteGroupValueRoleCounterRaw, workspace.WriteGroupValueRoleState:
	default:
		return "", snapshot.SampleReductionLatest, boundaryMemberMetadata
	}
	if member.SourceValueType == "" || member.SourceValueType != tagType {
		return "", snapshot.SampleReductionLatest, boundaryMemberMetadata
	}
	stored, ok := workspace.LongMemberStoredType(reduction, member.ValueRole, tagType)
	if !ok || member.ResultValueType == "" || member.ResultValueType != stored {
		return "", snapshot.SampleReductionLatest, boundaryMemberMetadata
	}
	if member.ValueRole == workspace.WriteGroupValueRoleCounterRaw &&
		(member.MeasurementID == nil || *member.MeasurementID == "") {
		return "", snapshot.SampleReductionLatest, boundaryMemberMetadata
	}
	if reduction == workspace.WriteGroupSampleReductionMean && member.ValueRole == workspace.WriteGroupValueRoleGauge {
		return stored, snapshot.SampleReductionMean, ""
	}
	return stored, snapshot.SampleReductionLatest, ""
}
