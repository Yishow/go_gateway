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

// newWideProfileGroupBoundary activates an applied managed_wide_profile_v1 group:
// one row per fixed UTC window with one column per member. Every member keeps its
// own role: a numeric gauge averages the window, a state or a confirmed
// counter_raw keeps its newest eligible reading, and each is stored as its frozen
// result type, never as the Tag type. The rows are always partial, so an optional
// member without a usable value is NULL, and the destination must prove it can
// hold that. A cumulative profile keeps every member's raw reading in its column
// and closes each against its own confirmed counter baseline, through the same
// engine and the same durable closure as every other layout. Every unknown,
// partial or contradictory fact is refused instead of falling back to another
// writer.
func newWideProfileGroupBoundary(ctx context.Context, cfg GroupBoundaryConfig) (*GroupBoundary, error) {
	blocked := func(code string, issues ...dbtarget.LayoutIssue) error {
		return &GroupBoundaryError{Code: code, Issues: issues}
	}
	group := cfg.Group
	profile := group.RecordingProfile
	if profile == nil || (profile.Kind != workspace.WriteGroupProfileInstantaneous && profile.Kind != workspace.WriteGroupProfileCumulative) {
		return nil, blocked(boundaryProfileUnsupported)
	}
	if profile.IntervalSeconds() != group.RowPolicy.IntervalSeconds {
		return nil, blocked(boundaryProfileInterval)
	}
	reduction := profile.EffectiveSampleReduction()
	mode := snapshot.ModeWideProfile
	cumulative := profile.Kind == workspace.WriteGroupProfileCumulative
	if cumulative {
		mode = snapshot.ModeWideCumulative
		if reduction != workspace.WriteGroupSampleReductionLatest {
			// A cumulative window counts its own interval; it never averages.
			return nil, blocked(boundaryMemberMetadata)
		}
		if cfg.Ledger != nil && !ledgerPersistsBaselines(cfg.Ledger) {
			// A ledger that cannot carry the baseline with the closure would restart
			// every counter silently. Refuse rather than look durable.
			return nil, blocked(boundaryCounterDurable)
		}
	}
	counters := make(map[string]aggregation.CounterBinding)
	members := make(map[string]boundaryMember, len(group.Members))
	layoutMembers := make([]dbtarget.GroupRowMember, 0, len(group.Members))
	snapshotMembers := make([]snapshot.Member, 0, len(group.Members))
	for _, member := range group.Members {
		if member.DeviceID == "" || member.PointID == "" || member.TagID == "" || member.TargetColumn == "" {
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
		stored, memberReduction, code := wideMemberShape(cumulative, reduction, member, tagType, kind)
		if code != "" {
			return nil, blocked(code)
		}
		key := memberKey(member.DeviceID, member.PointID, member.TagID)
		entry := boundaryMember{key: key, kind: kind}
		var binding *aggregation.CounterBinding
		if cumulative {
			bound, code := counterBindingFor(group, member, key, kind, cfg.Counters)
			if code != "" {
				return nil, blocked(code)
			}
			counters[key], binding, entry.counter = bound, &bound, &bound
		} else if member.ValueRole == workspace.WriteGroupValueRoleCounterRaw {
			// The latest raw value reuses the confirmed binding for its identity check
			// only: no baseline, delta or window is computed here.
			binding, code := counterBindingFor(group, member, key, kind, cfg.Counters)
			if code != "" {
				return nil, blocked(code)
			}
			entry.counter = &binding
		}
		members[key] = entry
		layoutMembers = append(layoutMembers, dbtarget.GroupRowMember{
			MemberKey: key, Column: member.TargetColumn, Type: stored, Required: member.Required,
			Reduction: string(memberReduction), ValueRole: string(member.ValueRole), SourceValueType: member.SourceValueType,
			Counter: binding,
		})
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
	// Readiness and the frozen proof judge the whole recording shape with this same
	// rule, so a direct activation can never serve a shape they refuse.
	if !workspace.WideProfileShapeSupported(group, cfg.TagTypes) {
		return nil, blocked(boundaryProfileUnsupported)
	}
	spec := boundaryRowSpec(cfg, layoutMembers, false)
	spec.Partial = true
	layout, issues := dbtarget.NewGroupRowLayout(spec)
	if len(issues) > 0 {
		return nil, blocked(boundaryLayoutBlocked, issues...)
	}
	parts := boundaryParts{
		members: members, snapshotMembers: snapshotMembers, encoder: layout,
		policy: snapshot.IncompletePartial, mode: mode, storage: layout.Capabilities(),
	}
	if cumulative {
		parts.counters = counters
	}
	return assembleBoundary(ctx, cfg, parts)
}

// wideMemberShape resolves one member's frozen facts to the stored result type
// and the effective reduction. A cumulative member carries no role and no
// reduction: its column keeps the raw reading of its own numeric Tag type, and
// the source and result types must both say so. Every other member is judged by
// the same rule as a long member.
func wideMemberShape(
	cumulative bool, reduction workspace.WriteGroupSampleReduction, member workspace.WriteGroupMember,
	tagType schema.DataType, source measurement.ExactType,
) (measurement.ExactType, snapshot.SampleReduction, string) {
	if !cumulative {
		return longMemberShape(reduction, member, tagType, source)
	}
	if member.ValueRole != "" || member.SourceValueType == "" || member.SourceValueType != tagType || member.ResultValueType != source {
		return "", snapshot.SampleReductionLatest, boundaryMemberMetadata
	}
	return source, snapshot.SampleReductionLatest, ""
}
