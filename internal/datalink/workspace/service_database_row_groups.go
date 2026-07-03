package workspace

import (
	"context"
	"fmt"
	"slices"
	"strings"
)

// DatabaseRowGroup describes one Step 4 shared-row planning group in a workspace.
type DatabaseRowGroup struct {
	ID               string   `json:"id"`
	ConnectorID      string   `json:"connector_id,omitempty"`
	TableSchema      string   `json:"table_schema,omitempty"`
	TableName        string   `json:"table_name"`
	MemberPointIDs   []string `json:"member_point_ids"`
	GroupKeyColumns  []string `json:"group_key_columns,omitempty"`
	UniqueKeyColumns []string `json:"unique_key_columns,omitempty"`
}

// DatabaseTargetRef stores workspace-only metadata for a persisted target row.
type DatabaseTargetRef struct {
	PointID    string `json:"point_id"`
	RowGroupID string `json:"row_group_id,omitempty"`
}

// ValidateDatabaseRowGroups checks whether row-group payload matches the expected database scope.
func ValidateDatabaseRowGroups(connectorID string, tableSchema string, tableName string, groups []DatabaseRowGroup) error {
	_, err := normalizeDatabaseRowGroups(connectorID, tableSchema, tableName, groups)
	return err
}

// SaveDatabaseRowGroups persists the Step 4 row-group plan without changing target rows.
func (s *Service) SaveDatabaseRowGroups(ctx context.Context, connectorID string, tableSchema string, tableName string, groups []DatabaseRowGroup) (*Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, err := s.getOrCreate(ctx)
	if err != nil {
		return nil, err
	}

	normalized, err := normalizeDatabaseRowGroups(connectorID, tableSchema, tableName, groups)
	if err != nil {
		return nil, err
	}
	record.DatabaseRowGroups = normalized
	record.DatabaseTargetRefs = filterDatabaseTargetRefs(record.DatabaseTargetRefs, normalized)
	record.UpdatedAt = s.now()
	if err := s.repo.Save(ctx, record); err != nil {
		return nil, fmt.Errorf("save workspace database row groups: %w", err)
	}
	return cloneRecord(record), nil
}

// SaveDatabaseTargetReference attaches workspace metadata to a point target.
func (s *Service) SaveDatabaseTargetReference(ctx context.Context, pointID string, rowGroupID string) (*Record, error) {
	pointID = strings.TrimSpace(pointID)
	rowGroupID = strings.TrimSpace(rowGroupID)
	if pointID == "" {
		return nil, workspaceValidationError("database target point id 不能為空")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	record, err := s.getOrCreate(ctx)
	if err != nil {
		return nil, err
	}
	if rowGroupID != "" {
		if !databaseRowGroupExists(record.DatabaseRowGroups, rowGroupID) {
			return nil, workspaceValidationError("database target row group does not exist")
		}
		if !databaseRowGroupContainsPoint(record.DatabaseRowGroups, rowGroupID, pointID) {
			return nil, workspaceValidationError("database target row group does not contain point")
		}
	}

	refs := make([]DatabaseTargetRef, 0, len(record.DatabaseTargetRefs)+1)
	for _, ref := range record.DatabaseTargetRefs {
		if ref.PointID != pointID {
			refs = append(refs, ref)
		}
	}
	if rowGroupID != "" {
		refs = append(refs, DatabaseTargetRef{PointID: pointID, RowGroupID: rowGroupID})
	}
	record.DatabaseTargetRefs = refs
	record.UpdatedAt = s.now()
	if err := s.repo.Save(ctx, record); err != nil {
		return nil, fmt.Errorf("save workspace database target reference: %w", err)
	}
	return cloneRecord(record), nil
}

func normalizeDatabaseRowGroups(connectorID string, tableSchema string, tableName string, groups []DatabaseRowGroup) ([]DatabaseRowGroup, error) {
	connectorID = strings.TrimSpace(connectorID)
	tableSchema = strings.TrimSpace(tableSchema)
	tableName = strings.TrimSpace(tableName)
	seen := map[string]struct{}{}
	out := make([]DatabaseRowGroup, 0, len(groups))
	for _, group := range groups {
		group.ID = strings.TrimSpace(group.ID)
		if group.ID == "" {
			return nil, workspaceValidationError("database row group id 不能為空")
		}
		if _, exists := seen[group.ID]; exists {
			return nil, workspaceValidationError("database row group id 不能重複")
		}
		seen[group.ID] = struct{}{}
		group.ConnectorID = defaultScopedValue(group.ConnectorID, connectorID)
		group.TableSchema = defaultScopedValue(group.TableSchema, tableSchema)
		group.TableName = defaultScopedValue(group.TableName, tableName)
		if group.ConnectorID != connectorID || group.TableSchema != tableSchema || group.TableName != tableName {
			return nil, workspaceValidationError("database row group scope must match workspace database config")
		}
		group.MemberPointIDs = normalizeStringSet(group.MemberPointIDs)
		group.GroupKeyColumns = normalizeStringSet(group.GroupKeyColumns)
		group.UniqueKeyColumns = normalizeStringSet(group.UniqueKeyColumns)
		out = append(out, group)
	}
	return out, nil
}

func filterDatabaseTargetRefs(refs []DatabaseTargetRef, groups []DatabaseRowGroup) []DatabaseTargetRef {
	out := make([]DatabaseTargetRef, 0, len(refs))
	for _, ref := range refs {
		pointID := strings.TrimSpace(ref.PointID)
		rowGroupID := strings.TrimSpace(ref.RowGroupID)
		if pointID == "" {
			continue
		}
		if rowGroupID == "" {
			out = append(out, DatabaseTargetRef{PointID: pointID})
			continue
		}
		if databaseRowGroupContainsPoint(groups, rowGroupID, pointID) {
			out = append(out, DatabaseTargetRef{PointID: pointID, RowGroupID: rowGroupID})
		}
	}
	return out
}

func databaseRowGroupExists(groups []DatabaseRowGroup, id string) bool {
	return slices.ContainsFunc(groups, func(group DatabaseRowGroup) bool {
		return group.ID == id
	})
}

func databaseRowGroupContainsPoint(groups []DatabaseRowGroup, id string, pointID string) bool {
	id = strings.TrimSpace(id)
	pointID = strings.TrimSpace(pointID)
	if id == "" || pointID == "" {
		return false
	}
	return slices.ContainsFunc(groups, func(group DatabaseRowGroup) bool {
		return group.ID == id && slices.Contains(group.MemberPointIDs, pointID)
	})
}

func defaultScopedValue(value string, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func normalizeStringSet(values []string) []string {
	out := make([]string, 0, len(values))
	seen := map[string]struct{}{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}

func cloneDatabaseRowGroups(groups []DatabaseRowGroup) []DatabaseRowGroup {
	if groups == nil {
		return nil
	}
	out := make([]DatabaseRowGroup, len(groups))
	for index, group := range groups {
		out[index] = group
		out[index].MemberPointIDs = append([]string{}, group.MemberPointIDs...)
		out[index].GroupKeyColumns = append([]string{}, group.GroupKeyColumns...)
		out[index].UniqueKeyColumns = append([]string{}, group.UniqueKeyColumns...)
	}
	return out
}

func cloneDatabaseTargetRefs(refs []DatabaseTargetRef) []DatabaseTargetRef {
	if refs == nil {
		return nil
	}
	return append([]DatabaseTargetRef{}, refs...)
}
