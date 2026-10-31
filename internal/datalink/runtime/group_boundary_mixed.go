package runtime

import (
	"context"
	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"
	"time"
)

func mixedBoundaryGroup(group *workspace.WriteGroup) bool {
	return group != nil && group.RecordingProfile != nil && group.RecordingProfile.Kind == workspace.WriteGroupProfileMixed
}

// newMixedGroupBoundary assembles the existing member reducers and encoders.
// Durable dispatch requires typed mixed progress, never a legacy common cursor.
func newMixedGroupBoundary(ctx context.Context, cfg GroupBoundaryConfig) (*GroupBoundary, error) {
	blocked := func(code string, err error) (*GroupBoundary, error) {
		return nil, &GroupBoundaryError{Code: code, Cause: err}
	}
	if cfg.Ledger != nil {
		capable, ok := cfg.Ledger.(interface{ PersistsMixedProgress() bool })
		if !ok || !capable.PersistsMixedProgress() {
			return blocked("mixed-durable-progress-unsupported", nil)
		}
	}
	group := cfg.Group
	if err := workspace.ValidateMixedRecordingSnapshot(group, cfg.TagTypes, cfg.CounterBindings); err != nil {
		return blocked(boundaryMemberMetadata, err)
	}
	if !cfg.FirstBucket.IsZero() || !cfg.Until.IsZero() {
		return blocked(boundarySnapshotInvalid, nil)
	}
	parts := boundaryParts{members: map[string]boundaryMember{}, counters: map[string]aggregation.CounterBinding{}, policy: snapshot.IncompletePartial, mode: snapshot.ModeMixedWide}
	wide := []dbtarget.GroupRowMember{}
	long := []dbtarget.LongRowMember{}
	for _, member := range group.Members {
		if member.DeviceID == "" || member.PointID == "" || member.TagID == "" || member.SourceRevision == "" || member.MappingRevision == "" || member.Cadence == nil {
			return blocked(boundaryMemberIdentity, nil)
		}
		interval, err := member.Cadence.Duration()
		if err != nil {
			return blocked(boundarySnapshotInvalid, err)
		}
		kind, ok := measurement.ExactTypeForTag(cfg.TagTypes[member.TagID])
		if !ok {
			return blocked(boundaryTagTypeUnsupport, nil)
		}
		key := memberKey(member.DeviceID, member.PointID, member.TagID)
		entry := boundaryMember{key: key, kind: kind, interval: interval, sourceRevision: member.SourceRevision, mappingRevision: member.MappingRevision}
		var binding *aggregation.CounterBinding
		sourceType := member.SourceValueType
		if member.Calculation == workspace.WriteGroupMemberCalculationDelta {
			bound, code := counterBindingFor(group, member, key, kind, cfg.Counters)
			if code != "" {
				return blocked(code, nil)
			}
			if cfg.Ledger != nil && !ledgerPersistsBaselines(cfg.Ledger) {
				return blocked(boundaryCounterDurable, nil)
			}
			binding = &bound
			entry.counter = binding
			parts.counters[key] = bound
			if member.ExactRecordingSource != nil {
				entry.exactScale = member.ExactRecordingSource.Scale
			}
		}
		parts.members[key] = entry
		var age time.Duration
		if member.MaxAgeSeconds != nil && *member.MaxAgeSeconds != 0 {
			age, err = snapshot.CheckedCadenceDuration(int64(*member.MaxAgeSeconds), "second")
			if err != nil {
				return blocked(boundarySnapshotInvalid, err)
			}
		}
		parts.snapshotMembers = append(parts.snapshotMembers, snapshot.Member{Key: key, Required: member.Required, MaxAge: age, SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision, SourceValueType: sourceType, Calculation: snapshot.Calculation(member.Calculation), Interval: interval})
		wide = append(wide, dbtarget.GroupRowMember{MemberKey: key, Column: member.TargetColumn, Type: member.ResultValueType, Required: member.Required, SourceValueType: member.SourceValueType, ValueRole: string(member.ValueRole), Calculation: snapshot.Calculation(member.Calculation), Interval: interval, Counter: binding})
		lm := dbtarget.LongRowMember{MemberKey: key, DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID, CircuitID: member.CircuitID, Type: member.ResultValueType, Required: member.Required, SourceValueType: member.SourceValueType, ValueRole: string(member.ValueRole), Calculation: snapshot.Calculation(member.Calculation), Interval: interval, Counter: binding}
		if member.MeasurementID != nil {
			lm.MeasurementID = *member.MeasurementID
		}
		long = append(long, lm)
	}
	var issues []dbtarget.LayoutIssue
	if group.EffectiveLayoutVersion() == workspace.WriteGroupLayoutManagedLongV1 {
		parts.mode = snapshot.ModeMixedLong
		encoder, found := dbtarget.NewLongRowLayout(dbtarget.LongRowSpec{ResultShape: string(group.RecordingProfile.ResultShape), Dialect: cfg.Dialect, Columns: cfg.Columns, WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: group.Revision, ProfileKind: "mixed", Members: long})
		parts.encoder = encoder
		issues = found
		// Long's fixed nullable value slots and provenance are checked by its layout.
		parts.storage = snapshot.StorageCapabilities{NullableValues: true, MemberQuality: true}
	} else {
		spec := boundaryRowSpec(cfg, wide, false)
		spec.ResultShape = string(group.RecordingProfile.ResultShape)
		spec.WorkspaceID = group.WorkspaceID
		spec.GroupRevision = group.Revision
		spec.DueAtColumn = group.RowPolicy.DueAtColumn
		spec.Partial = true
		encoder, found := dbtarget.NewGroupRowLayout(spec)
		parts.encoder = encoder
		issues = found
		if len(issues) == 0 {
			parts.storage = encoder.Capabilities()
		}
	}
	if len(issues) > 0 {
		return nil, &GroupBoundaryError{Code: boundaryLayoutBlocked, Issues: issues}
	}
	return assembleBoundary(ctx, cfg, parts)
}
