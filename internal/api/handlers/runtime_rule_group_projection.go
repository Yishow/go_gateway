package handlers

import (
	"time"

	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/workspace"
)

func projectRuntimeRecordingGroup(saved *workspace.WriteGroup, snapshot *workspace.WriteGroupAppliedSnapshot, scope runtimeview.Scope, observed string) (runtimeview.GroupView, error) {
	result := runtimeview.GroupView{GroupID: saved.ID, SavedRevision: runtimeEvidenceString(saved.Revision),
		AppliedRevision: runtimeEvidenceString(saved.AppliedRevision), AsOf: &observed, State: "saved-only", Proof: "unavailable"}
	var err error
	result.Saved, err = runtimeRecordingFacts(saved, nil, scope)
	if err != nil {
		return result, err
	}
	if saved.AppliedRevision == "" {
		return result, nil
	}
	result.State = "unavailable"
	if snapshot == nil {
		return result, nil
	}
	result.Effective, err = runtimeRecordingFacts(snapshot.Group, snapshot, scope)
	if err != nil {
		return result, err
	}
	result.State, result.Proof = "effective", "applied-snapshot-and-layout"
	result.EffectiveRevision = runtimeEvidenceString(snapshot.AppliedRevision)
	if saved.AppliedRevision == snapshot.AppliedRevision {
		result.AppliedEffectiveAt = runtimeEvidenceTime(snapshot.EffectiveAt)
	}
	return result, nil
}

func runtimeRecordingFacts(group *workspace.WriteGroup, snapshot *workspace.WriteGroupAppliedSnapshot, scope runtimeview.Scope) (*runtimeview.GroupFacts, error) {
	if group == nil || group.WorkspaceID != scope.WorkspaceID {
		return nil, runtimeEvidenceError(runtimeview.ReasonForeign)
	}
	for _, value := range []string{group.ID, group.Revision, group.Destination.ConnectorID, group.Destination.ConnectorRevision, group.Destination.TableName} {
		if value == "" || runtimeSafeID(value) != value {
			return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
		}
	}
	switch group.Status {
	case workspace.WriteGroupStatusDraft, workspace.WriteGroupStatusReady, workspace.WriteGroupStatusRunning, workspace.WriteGroupStatusDisabled, workspace.WriteGroupStatusDeleted:
	default:
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	result := &runtimeview.GroupFacts{
		Binding:     runtimeview.Binding{WorkspaceID: scope.WorkspaceID, WorkspaceRevision: scope.WorkspaceRevision, EntityRevision: group.Revision},
		Enabled:     group.Status != workspace.WriteGroupStatusDisabled && group.Status != workspace.WriteGroupStatusDeleted,
		Destination: runtimeview.Destination{ConnectorID: group.Destination.ConnectorID, ConnectorRevision: group.Destination.ConnectorRevision, Table: group.Destination.TableName, Kind: "unavailable"},
		Members:     []runtimeview.Member{}, Layout: runtimeRecordingEnum(string(group.EffectiveLayoutVersion()), "managed_wide_v1", "managed_long_v1", "managed_wide_profile_v1"),
		StorageMode: runtimeRecordingEnum(string(group.StorageMode), "plain", "timescale", "native_monthly"),
		Reduction:   "unavailable", WriteMode: runtimeRecordingEnum(group.WritePolicy.Mode, "append", "latest"),
		IncompletePolicy: runtimeRecordingEnum(group.RowPolicy.IncompletePolicy, "skip_row", "partial"),
	}
	if schema := group.Destination.TableSchema; schema != "" {
		if runtimeSafeID(schema) != schema {
			return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
		}
		result.Destination.Schema = &schema
	}
	if group.RowPolicy.IntervalSeconds > 0 {
		seconds := int64(group.RowPolicy.IntervalSeconds)
		result.IntervalSeconds = &seconds
	}
	if group.RowPolicy.AllowedLatenessSeconds >= 0 {
		seconds := int64(group.RowPolicy.AllowedLatenessSeconds)
		result.AllowedLatenessSeconds = &seconds
	}
	if profile := group.RecordingProfile; profile != nil && profile.Kind == workspace.WriteGroupProfileMixed {
		if err := runtimeMixedGroupFacts(group, result); err != nil {
			return nil, err
		}
	} else if profile != nil {
		result.RecordingProfile = &runtimeview.RecordingProfile{Kind: string(profile.Kind), Interval: int64(profile.Interval), Unit: string(profile.Unit)}
		if profile.Kind == workspace.WriteGroupProfileCumulative {
			result.Reduction = "cumulative_delta"
		} else {
			result.Reduction = runtimeRecordingEnum(string(profile.EffectiveSampleReduction()), "latest", "mean")
		}
	} else if group.EffectiveLayoutVersion() == workspace.WriteGroupLayoutManagedWideV1 {
		result.Reduction = "latest" // The existing legacy-wide authority defines latest.
	}
	if snapshot != nil {
		result.EffectiveAt, result.EndsAt = runtimeEvidenceTime(snapshot.EffectiveAt), runtimeEvidenceTime(snapshot.EndsAt)
		result.Destination.Kind = runtimeRecordingEnum(string(snapshot.RuntimeLayout.Dialect), "sqlite", "postgres", "mysql", "sqlserver")
	}
	for _, member := range group.Members {
		for _, value := range []string{member.DeviceID, member.PointID, member.TagID, member.TargetColumn} {
			if value == "" || runtimeSafeID(value) != value {
				return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
			}
		}
		if member.CircuitID != "" && runtimeSafeID(member.CircuitID) != member.CircuitID {
			return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
		}
		valueType := string(member.SourceValueType)
		if valueType == "" && snapshot != nil {
			valueType = string(snapshot.RuntimeLayout.TagTypes[member.TagID])
		}
		projected := runtimeview.Member{DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID,
			Column: member.TargetColumn, CircuitID: runtimeEvidenceString(member.CircuitID)}
		if runtimeSafeEnum(valueType, "bool", "int16", "uint16", "int32", "uint32", "int64", "uint64", "float32", "float64", "string") != "" {
			projected.SourceValueType = &valueType
		}
		if result.MixedRecording != nil {
			if err := runtimeMixedMemberFacts(member, &projected); err != nil {
				return nil, err
			}
		}
		result.Members = append(result.Members, projected)
	}
	return result, nil
}

func runtimeRecordingEnum(value string, allowed ...string) string {
	if result := runtimeSafeEnum(value, allowed...); result != "" {
		return result
	}
	return "unavailable"
}

func runtimeEvidenceString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func runtimeEvidenceTime(value time.Time) *string {
	if value.IsZero() {
		return nil
	}
	return runtimeEvidenceString(value.UTC().Format(time.RFC3339Nano))
}
