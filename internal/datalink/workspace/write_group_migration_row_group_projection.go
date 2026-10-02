package workspace

import (
	"fmt"
	"slices"
	"strings"
)

// projectRowGroupWriteGroup replaces the original workspace row-group
// projection in place. The canonical UUID remains separate from the legacy
// row-group ID so old readers and target refs retain their stable identity.
func projectRowGroupWriteGroup(record *Record, group *WriteGroup, legacyRowGroupID string) error {
	if record == nil || group == nil || strings.TrimSpace(legacyRowGroupID) == "" {
		return fmt.Errorf("project row-group write group: %w", ErrWriteGroupValidation)
	}
	legacyRowGroupID = strings.TrimSpace(legacyRowGroupID)
	memberPointIDs := make([]string, 0, len(group.Members))
	for _, member := range group.Members {
		pointID := strings.TrimSpace(member.PointID)
		if pointID != "" && !slices.Contains(memberPointIDs, pointID) {
			memberPointIDs = append(memberPointIDs, pointID)
		}
	}
	rowGroup := DatabaseRowGroup{
		ID:               legacyRowGroupID,
		ConnectorID:      group.Destination.ConnectorID,
		TableSchema:      group.Destination.TableSchema,
		TableName:        group.Destination.TableName,
		MemberPointIDs:   memberPointIDs,
		GroupKeyColumns:  slices.Clone(group.RowPolicy.GroupKeyColumns),
		UniqueKeyColumns: slices.Clone(group.RowPolicy.UniqueKeyColumns),
	}
	replaced := false
	for index, existing := range record.DatabaseRowGroups {
		if existing.ID != legacyRowGroupID {
			continue
		}
		record.DatabaseRowGroups[index] = rowGroup
		replaced = true
		break
	}
	if !replaced {
		return fmt.Errorf("project row-group write group: legacy row group %s: %w", legacyRowGroupID, ErrWriteGroupNotFound)
	}

	memberSet := make(map[string]struct{}, len(memberPointIDs))
	for _, pointID := range memberPointIDs {
		memberSet[pointID] = struct{}{}
	}
	refs := make([]DatabaseTargetRef, 0, len(record.DatabaseTargetRefs)+len(memberPointIDs))
	for _, ref := range record.DatabaseTargetRefs {
		pointID := strings.TrimSpace(ref.PointID)
		if pointID == "" {
			continue
		}
		if ref.RowGroupID == legacyRowGroupID {
			continue
		}
		if _, member := memberSet[pointID]; member {
			continue
		}
		refs = append(refs, DatabaseTargetRef{PointID: pointID, RowGroupID: strings.TrimSpace(ref.RowGroupID)})
	}
	for _, pointID := range memberPointIDs {
		refs = append(refs, DatabaseTargetRef{PointID: pointID, RowGroupID: legacyRowGroupID})
	}
	record.DatabaseTargetRefs = refs
	return nil
}
