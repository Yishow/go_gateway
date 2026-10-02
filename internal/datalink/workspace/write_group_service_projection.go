package workspace

import (
	"fmt"
	"slices"
	"strings"
)

// projectWriteGroup mirrors canonical group identity and source points into
// the existing workspace compatibility metadata. It deliberately preserves
// unrelated legacy groups and the workspace's connector ownership field.
func projectWriteGroup(record *Record, group *WriteGroup) error {
	if record == nil || group == nil || strings.TrimSpace(group.ID) == "" {
		return fmt.Errorf("project write group: %w", ErrWriteGroupValidation)
	}
	if strings.TrimSpace(group.Migration.SourceKind) == writeGroupMigrationSourceKindLegacyRowGroup && strings.TrimSpace(group.Migration.LegacyRowGroupID) != "" {
		return projectRowGroupWriteGroup(record, group, group.Migration.LegacyRowGroupID)
	}
	memberPointIDs := make([]string, 0, len(group.Members))
	for _, member := range group.Members {
		if pointID := strings.TrimSpace(member.PointID); pointID != "" && !slices.Contains(memberPointIDs, pointID) {
			memberPointIDs = append(memberPointIDs, pointID)
		}
	}
	rowGroup := DatabaseRowGroup{
		ID:             group.ID,
		ConnectorID:    group.Destination.ConnectorID,
		TableSchema:    group.Destination.TableSchema,
		TableName:      group.Destination.TableName,
		MemberPointIDs: memberPointIDs,
	}
	replaced := false
	for index, existing := range record.DatabaseRowGroups {
		if existing.ID != group.ID {
			continue
		}
		rowGroup.GroupKeyColumns = append([]string(nil), existing.GroupKeyColumns...)
		rowGroup.UniqueKeyColumns = append([]string(nil), existing.UniqueKeyColumns...)
		record.DatabaseRowGroups[index] = rowGroup
		replaced = true
		break
	}
	if !replaced {
		record.DatabaseRowGroups = append(record.DatabaseRowGroups, rowGroup)
	}

	refs := make([]DatabaseTargetRef, 0, len(record.DatabaseTargetRefs)+len(memberPointIDs))
	for _, ref := range record.DatabaseTargetRefs {
		pointID := strings.TrimSpace(ref.PointID)
		rowGroupID := strings.TrimSpace(ref.RowGroupID)
		if pointID == "" {
			continue
		}
		if rowGroupID == group.ID {
			continue
		}
		refs = append(refs, DatabaseTargetRef{PointID: pointID, RowGroupID: rowGroupID})
	}
	for _, pointID := range memberPointIDs {
		if slices.ContainsFunc(refs, func(ref DatabaseTargetRef) bool {
			return ref.PointID == pointID
		}) {
			continue
		}
		refs = append(refs, DatabaseTargetRef{PointID: pointID, RowGroupID: group.ID})
	}
	record.DatabaseTargetRefs = refs
	return nil
}
