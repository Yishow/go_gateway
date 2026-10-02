package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
)

const writeGroupMigrationSourceKindLegacySingleMapping = "legacy-single-mapping"

type legacyWritePair struct {
	tagID       string
	connectorID string
}

func (s *WriteGroupService) checkLegacyTargetWriteWithRecord(
	ctx context.Context,
	tx *sql.Tx,
	record *Record,
	mappingID, tagID, connectorID string,
) error {
	if record == nil {
		return nil
	}
	groups, err := s.readLegacyOwnershipGroups(ctx, tx, record.ID)
	if err != nil {
		return err
	}
	requestedPair := legacyWritePair{tagID: strings.TrimSpace(tagID), connectorID: strings.TrimSpace(connectorID)}
	for _, group := range groups {
		pairs, err := s.validateLegacyWriteGroup(ctx, tx, group)
		if err != nil {
			return err
		}
		if mappingID != "" && slices.ContainsFunc(group.Migration.SourceIDs, func(sourceID string) bool {
			return strings.TrimSpace(sourceID) == mappingID
		}) {
			return legacyWriteConflict()
		}
		if requestedPair.tagID != "" && requestedPair.connectorID != "" {
			if strings.TrimSpace(group.Destination.ConnectorID) == requestedPair.connectorID {
				for _, member := range group.Members {
					if strings.TrimSpace(member.TagID) == requestedPair.tagID {
						return legacyWriteConflict()
					}
				}
			}
			for _, pair := range pairs {
				if pair == requestedPair {
					return legacyWriteConflict()
				}
			}
		}
	}
	return nil
}

func (s *WriteGroupService) readLegacyOwnershipGroups(ctx context.Context, tx *sql.Tx, workspaceID string) ([]*WriteGroup, error) {
	if tx == nil || strings.TrimSpace(workspaceID) == "" {
		return nil, legacyWriteMalformed("canonical ownership scope is incomplete")
	}
	groups, err := s.repo.ListInTx(ctx, tx, workspaceID)
	if err != nil {
		return nil, legacyWriteOperationalReadError("canonical ownership could not be read", err)
	}
	return groups, nil
}

func (s *WriteGroupService) validateLegacyWriteGroup(
	ctx context.Context,
	tx *sql.Tx,
	group *WriteGroup,
) ([]legacyWritePair, error) {
	if group == nil || strings.TrimSpace(group.ID) == "" {
		return nil, legacyWriteMalformed("canonical group identity is incomplete")
	}
	migration := group.Migration
	if strings.TrimSpace(migration.SourceKind) == "" {
		if hasUnexpectedMigrationProvenance(migration) {
			return nil, legacyWriteMalformed("canonical migration provenance is incomplete")
		}
		return nil, nil
	}
	switch strings.TrimSpace(migration.SourceKind) {
	case writeGroupMigrationSourceKindLegacySingleMapping:
		if len(migration.SourceIDs) != 1 || strings.TrimSpace(migration.LegacyRowGroupID) != "" || len(migration.TargetMappingPoints) != 0 {
			return nil, legacyWriteMalformed("single-mapping ownership provenance is incomplete")
		}
		sourceID := strings.TrimSpace(migration.SourceIDs[0])
		migrationMap, found, err := s.repo.getMigrationMapInTx(ctx, tx, group.WorkspaceID, migration.SourceKind, sourceID)
		if err != nil {
			return nil, legacyWriteOperationalReadError("single-mapping ownership map could not be read", err)
		}
		if !found || migrationMap.GroupID != group.ID || migrationMap.SourceKind != migration.SourceKind || migrationMap.SourceID != sourceID {
			return nil, legacyWriteMalformed("single-mapping ownership map is incomplete")
		}
		var intent WriteGroupMigrationIntent
		if err := json.Unmarshal([]byte(migrationMap.BeforeIntent), &intent); err != nil {
			return nil, legacyWriteMalformed("single-mapping before intent is malformed")
		}
		if strings.TrimSpace(intent.SourceID) != sourceID {
			return nil, legacyWriteMalformed("single-mapping before intent identity is incomplete")
		}
		pair, ok := legacyWritePairFromIntent(intent)
		if !ok {
			return nil, legacyWriteMalformed("single-mapping before intent target is incomplete")
		}
		return []legacyWritePair{pair}, nil

	case writeGroupMigrationSourceKindLegacyRowGroup:
		legacyRowGroupID := strings.TrimSpace(migration.LegacyRowGroupID)
		if legacyRowGroupID == "" || len(migration.SourceIDs) == 0 || len(migration.TargetMappingPoints) != len(migration.SourceIDs) {
			return nil, legacyWriteMalformed("row-group ownership provenance is incomplete")
		}
		sourceIDs := make(map[string]struct{}, len(migration.SourceIDs))
		for _, rawSourceID := range migration.SourceIDs {
			sourceID := strings.TrimSpace(rawSourceID)
			if sourceID == "" {
				return nil, legacyWriteMalformed("row-group ownership source identity is incomplete")
			}
			if _, exists := sourceIDs[sourceID]; exists {
				return nil, legacyWriteMalformed("row-group ownership source identity is duplicated")
			}
			sourceIDs[sourceID] = struct{}{}
		}
		for rawMappingID, rawPointID := range migration.TargetMappingPoints {
			mappingID := strings.TrimSpace(rawMappingID)
			if mappingID == "" || strings.TrimSpace(rawPointID) == "" {
				return nil, legacyWriteMalformed("row-group target mapping provenance is incomplete")
			}
			if _, exists := sourceIDs[mappingID]; !exists {
				return nil, legacyWriteMalformed("row-group target mapping provenance has an unknown source")
			}
		}
		migrationMap, found, err := s.repo.getMigrationMapInTx(ctx, tx, group.WorkspaceID, migration.SourceKind, legacyRowGroupID)
		if err != nil {
			return nil, legacyWriteOperationalReadError("row-group ownership map could not be read", err)
		}
		if !found || migrationMap.GroupID != group.ID || migrationMap.SourceKind != migration.SourceKind || migrationMap.SourceID != legacyRowGroupID {
			return nil, legacyWriteMalformed("row-group ownership map is incomplete")
		}
		var intent WriteGroupRowGroupMigrationIntent
		if err := json.Unmarshal([]byte(migrationMap.BeforeIntent), &intent); err != nil {
			return nil, legacyWriteMalformed("row-group before intent is malformed")
		}
		if strings.TrimSpace(intent.SourceID) != legacyRowGroupID || len(intent.Members) == 0 {
			return nil, legacyWriteMalformed("row-group before intent identity is incomplete")
		}
		pairs := make([]legacyWritePair, 0, len(intent.Members))
		beforeIntentSourceIDs := make(map[string]struct{}, len(intent.Members))
		for _, member := range intent.Members {
			sourceID := strings.TrimSpace(member.SourceID)
			if sourceID == "" {
				return nil, legacyWriteMalformed("row-group before intent source identity is incomplete")
			}
			if _, exists := sourceIDs[sourceID]; !exists {
				return nil, legacyWriteMalformed("row-group before intent source identity is incomplete")
			}
			if _, ok := migration.TargetMappingPoints[sourceID]; !ok {
				return nil, legacyWriteMalformed("row-group target mapping provenance is incomplete")
			}
			if _, duplicate := beforeIntentSourceIDs[sourceID]; duplicate {
				return nil, legacyWriteMalformed("row-group before intent source identity is duplicated")
			}
			beforeIntentSourceIDs[sourceID] = struct{}{}
			pair, ok := legacyWritePairFromIntent(member)
			if !ok {
				return nil, legacyWriteMalformed("row-group before intent target is incomplete")
			}
			pairs = append(pairs, pair)
		}
		if len(pairs) != len(migration.SourceIDs) {
			return nil, legacyWriteMalformed("row-group before intent does not cover all mappings")
		}
		if len(beforeIntentSourceIDs) != len(sourceIDs) {
			return nil, legacyWriteMalformed("row-group before intent does not cover all mappings")
		}
		for sourceID := range sourceIDs {
			if _, found := beforeIntentSourceIDs[sourceID]; !found {
				return nil, legacyWriteMalformed("row-group before intent does not cover all mappings")
			}
		}
		return pairs, nil

	default:
		return nil, legacyWriteMalformed("canonical migration source kind is unknown")
	}
}

func legacyWritePairFromIntent(intent WriteGroupMigrationIntent) (legacyWritePair, bool) {
	tagID := strings.TrimSpace(intent.TagID)
	connectorID := strings.TrimSpace(intent.ConnectorID)
	return legacyWritePair{tagID: tagID, connectorID: connectorID}, tagID != "" && connectorID != ""
}

func hasUnexpectedMigrationProvenance(migration WriteGroupMigration) bool {
	return len(migration.SourceIDs) > 0 || strings.TrimSpace(migration.LegacyRowGroupID) != "" ||
		strings.TrimSpace(migration.SourceRevision) != "" || strings.TrimSpace(migration.AdapterVersion) != "" ||
		strings.TrimSpace(migration.ReviewResult) != "" || len(migration.TargetMappingPoints) > 0
}

func legacyWriteConflict() error {
	return fmt.Errorf("canonical write-group ownership conflict: open_write_groups: %w", ErrWriteGroupLegacyWriteConflict)
}

func legacyWriteMalformed(reason string) error {
	return fmt.Errorf("%s: %w", reason, errors.Join(ErrWriteGroupLegacyWriteConflict, ErrWriteGroupValidation))
}

func legacyWriteOperationalReadError(reason string, cause error) error {
	return fmt.Errorf("%s: %w", reason, errors.Join(ErrWriteGroupServiceUnavailable, cause))
}

func (s *WriteGroupService) checkLegacyRowGroupReplacementWithRecord(
	ctx context.Context,
	tx *sql.Tx,
	current *Record,
	connectorID string,
	groups []DatabaseRowGroup,
) error {
	if current == nil || strings.TrimSpace(current.ID) == "" {
		return legacyWriteMalformed("legacy row-group workspace identity is incomplete")
	}
	owned, err := s.readOwnedRowGroupIDs(ctx, tx, current.ID)
	if err != nil {
		return err
	}
	for rowGroupID := range owned {
		currentGroup, found := findDatabaseRowGroup(current.DatabaseRowGroups, rowGroupID)
		if !found {
			return legacyWriteMalformed("canonical row-group projection is incomplete")
		}
		requestedGroup, found := findDatabaseRowGroup(groups, rowGroupID)
		if !found || !databaseRowGroupsEquivalent(currentGroup, requestedGroup, connectorID) {
			return fmt.Errorf("canonical row-group %s is owned: %w", rowGroupID, ErrWriteGroupLegacyWriteConflict)
		}
	}
	return nil
}

func (s *WriteGroupService) readOwnedRowGroupIDs(ctx context.Context, tx *sql.Tx, workspaceID string) (map[string]struct{}, error) {
	groups, err := s.readLegacyOwnershipGroups(ctx, tx, workspaceID)
	if err != nil {
		return nil, err
	}
	owned := make(map[string]struct{}, len(groups))
	for _, group := range groups {
		if strings.TrimSpace(group.Migration.SourceKind) == writeGroupMigrationSourceKindLegacyRowGroup {
			if _, err := s.validateLegacyWriteGroup(ctx, tx, group); err != nil {
				return nil, err
			}
			owned[strings.TrimSpace(group.Migration.LegacyRowGroupID)] = struct{}{}
			continue
		}
		if _, err := s.validateLegacyWriteGroup(ctx, tx, group); err != nil {
			return nil, err
		}
		owned[strings.TrimSpace(group.ID)] = struct{}{}
	}
	return owned, nil
}

func databaseRowGroupsEquivalent(left, right DatabaseRowGroup, connectorID string) bool {
	left = normalizedDatabaseRowGroupForComparison(left, connectorID, left.TableSchema, left.TableName)
	right = normalizedDatabaseRowGroupForComparison(right, connectorID, left.TableSchema, left.TableName)
	return left.ID == right.ID && left.ConnectorID == right.ConnectorID && left.TableSchema == right.TableSchema &&
		left.TableName == right.TableName && equalStringSets(left.MemberPointIDs, right.MemberPointIDs) &&
		equalStringSets(left.GroupKeyColumns, right.GroupKeyColumns) && equalStringSets(left.UniqueKeyColumns, right.UniqueKeyColumns)
}

func normalizedDatabaseRowGroupForComparison(group DatabaseRowGroup, connectorID, tableSchema, tableName string) DatabaseRowGroup {
	group.ID = strings.TrimSpace(group.ID)
	group.ConnectorID = strings.TrimSpace(group.ConnectorID)
	if group.ConnectorID == "" {
		group.ConnectorID = strings.TrimSpace(connectorID)
	}
	group.TableSchema = strings.TrimSpace(group.TableSchema)
	if group.TableSchema == "" {
		group.TableSchema = strings.TrimSpace(tableSchema)
	}
	group.TableName = strings.TrimSpace(group.TableName)
	if group.TableName == "" {
		group.TableName = strings.TrimSpace(tableName)
	}
	group.MemberPointIDs = normalizeStringSetForComparison(group.MemberPointIDs)
	group.GroupKeyColumns = normalizeStringSetForComparison(group.GroupKeyColumns)
	group.UniqueKeyColumns = normalizeStringSetForComparison(group.UniqueKeyColumns)
	return group
}

func normalizeStringSetForComparison(values []string) []string {
	values = normalizeStringSet(values)
	slices.Sort(values)
	return values
}

func equalStringSets(left, right []string) bool {
	return slices.Equal(normalizeStringSetForComparison(left), normalizeStringSetForComparison(right))
}
