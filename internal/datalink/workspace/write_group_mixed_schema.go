package workspace

import (
	"go-gateway/internal/datalink/schema"
	"reflect"
)

// mixedSchemaShapeSupported validates storage intent, not runtime capability.
// It deliberately does not extend WideProfileShapeSupported or the runtime gate.
func mixedSchemaShapeSupported(group *WriteGroup, types map[string]schema.DataType) bool {
	if !isMixedWriteGroup(group) || group.Destination.StorageStrategy != WriteGroupStorageStrategyManaged || (!group.isWideProfile() && !group.isLongLayout()) {
		return false
	}
	if group.StorageMode != "" && group.StorageMode != WriteGroupStorageModePlain {
		return false
	}
	if ValidateMemberCalculationIntent(group) != nil || validateMixedMemberCadences(group) != nil || len(group.Members) == 0 {
		return false
	}
	if group.isWideProfile() && (group.RowPolicy.RecordKeyColumn == "" || group.RowPolicy.ProvenanceColumn == "" || group.RowPolicy.DueAtColumn == "" || group.RowPolicy.BucketStartColumn != "" || group.RowPolicy.IncompletePolicy != "partial") {
		return false
	}
	if _, err := mixedRecordProfileKey(group); err != nil {
		return false
	}
	for i, m := range group.Members {
		typ, ok := types[m.TagID]
		if !ok || m.SourceValueType == "" || m.ResultValueType == "" || m.EntityKey != "" {
			return false
		}
		expected := m
		if applyMixedRecordingMetadata(&expected, i, typ) != nil || expected.SourceValueType != m.SourceValueType || expected.ResultValueType != m.ResultValueType {
			return false
		}
	}
	return true
}

// sameMixedSchemaIntent adds calculation-policy and full member identity to the
// existing proof invalidation rule, while leaving all legacy comparisons alone.
func sameMixedSchemaIntent(a, b *WriteGroup) bool {
	if !isMixedWriteGroup(a) && !isMixedWriteGroup(b) {
		return true
	}
	if !isMixedWriteGroup(a) || !isMixedWriteGroup(b) || a.WorkspaceID != b.WorkspaceID || a.ID != b.ID || a.RecordID != b.RecordID || a.LayoutVersion != b.LayoutVersion || a.StorageMode != b.StorageMode || a.RowPolicy.DueAtColumn != b.RowPolicy.DueAtColumn {
		return false
	}
	if len(a.Members) != len(b.Members) {
		return false
	}
	for _, m := range a.Members {
		found := false
		for _, n := range b.Members {
			if m.DeviceID == n.DeviceID && m.PointID == n.PointID && m.TagID == n.TagID {
				found = m.Required == n.Required && reflect.DeepEqual(m.MeasurementID, n.MeasurementID) && m.SourceRevision == n.SourceRevision && m.MappingRevision == n.MappingRevision && reflect.DeepEqual(m.MaxAgeSeconds, n.MaxAgeSeconds)
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// mixedReplayComparable resolves only server-generated omissions for an already
// scoped canonical create replay. Explicit stale source revisions are rejected
// by sameMixedMemberRevisions before this helper is used. Schema comparisons
// for saved updates remain strict and are not weakened for this create path.
func mixedReplayComparable(candidate, existing *WriteGroup) *WriteGroup {
	comparable := cloneWriteGroup(candidate)
	comparable.ID = existing.ID
	for i := range comparable.Members {
		m := &comparable.Members[i]
		for _, n := range existing.Members {
			if m.DeviceID == n.DeviceID && m.PointID == n.PointID && m.TagID == n.TagID {
				if m.SourceRevision == "" {
					m.SourceRevision = n.SourceRevision
				}
				if m.MappingRevision == "" {
					m.MappingRevision = n.MappingRevision
				}
				break
			}
		}
	}
	return comparable
}
