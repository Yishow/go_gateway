package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
)

// ReviewRowGroupMigration persists one or more reviewed workspace row-group
// conversions. It shares the single-mapping review envelope and keeps all
// group, map and compatibility projection writes in one local transaction.
func (s *WriteGroupService) ReviewRowGroupMigration(
	ctx context.Context,
	request WriteGroupMigrationReviewRequest,
) (result *WriteGroupMigrationReviewResult, err error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	request, ids, err := prepareWriteGroupMigrationReviewRequest(request)
	if err != nil {
		return nil, err
	}
	store, ok := s.workspaceSvc.repo.(DatabaseSetupStore)
	if !ok {
		return nil, ErrWriteGroupServiceUnavailable
	}

	s.workspaceSvc.mu.Lock()
	defer s.workspaceSvc.mu.Unlock()
	setupTx, err := store.BeginDatabaseSetup(ctx)
	if err != nil {
		return nil, normalizeWriteGroupServiceError("begin row-group migration review", err)
	}
	committed := false
	defer func() {
		if committed {
			return
		}
		if rollbackErr := setupTx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			result = nil
			err = errors.Join(err, fmt.Errorf("rollback row-group migration review: %w", rollbackErr))
		}
	}()

	record, err := setupTx.Get(ctx)
	if errors.Is(err, ErrNotFound) {
		return nil, writeGroupNotFound("review row-group migration workspace")
	}
	if err != nil {
		return nil, normalizeWriteGroupServiceError("read row-group migration workspace", err)
	}
	if record.ID != request.WorkspaceID {
		return nil, writeGroupNotFound("review row-group migration workspace")
	}
	if record.DatabaseSetupRevision != request.ExpectedWorkspaceRevision {
		return nil, writeGroupRevisionConflict("row-group migration workspace revision is stale")
	}
	tx := setupTx.SQLTx()
	connector, err := readMigrationPreviewConnector(ctx, tx, s.repo, record.DatabaseConnectorID)
	if err != nil {
		return nil, err
	}
	if connector.identityRevision != request.ExpectedConnectorRevision {
		return nil, writeGroupRevisionConflict("row-group migration connector revision is stale")
	}
	preview, err := buildRowGroupMigrationPreview(ctx, tx, s.repo, record, connector, ids)
	if err != nil {
		return nil, normalizeWriteGroupServiceError("recompute row-group migration review", err)
	}
	if preview.ReviewDigest != request.ReviewDigest {
		return nil, writeGroupRevisionConflict("row-group migration review digest is stale")
	}

	type reviewedItem struct {
		item         WriteGroupMigrationPreviewItem
		existing     *WriteGroup
		migrationMap *writeGroupMigrationMap
	}
	reviewed := make([]reviewedItem, 0, len(preview.Items))
	for _, item := range preview.Items {
		if item.Status == WriteGroupMigrationPreviewStatusBlocked {
			return nil, fmt.Errorf("row-group migration source %s is blocked: %w", item.SourceID, ErrWriteGroupValidation)
		}
		if item.Status != WriteGroupMigrationPreviewStatusNeedsReview || item.CandidateGroup == nil || item.BeforeRowGroupIntent == nil || item.SourceRevision == "" {
			return nil, fmt.Errorf("row-group migration source %s is not reviewable: %w", item.SourceID, ErrWriteGroupValidation)
		}
		candidate := item.CandidateGroup
		migrationMap, found, err := s.repo.getMigrationMapInTx(ctx, tx, record.ID, candidate.Migration.SourceKind, item.SourceID)
		if err != nil {
			return nil, err
		}
		if !found {
			reviewed = append(reviewed, reviewedItem{item: item})
			continue
		}
		if migrationMap.SourceRevision != item.SourceRevision {
			return nil, writeGroupRevisionConflict("row-group migration source revision changed after prior review")
		}
		if err := validateRowGroupMigrationMap(migrationMap, record.ID, item, preview.AdapterVersion); err != nil {
			return nil, fmt.Errorf("row-group migration identity map is incomplete: %w", ErrWriteGroupValidation)
		}
		group, err := s.repo.GetInTx(ctx, tx, record.ID, migrationMap.GroupID)
		if err != nil {
			return nil, fmt.Errorf("load row-group migration identity map group: %w", ErrWriteGroupValidation)
		}
		if err := validateReviewedRowGroupMigration(group, migrationMap, record.ID, item); err != nil {
			return nil, err
		}
		if group.Destination.ConnectorID != connector.id || group.Destination.ConnectorRevision != request.ExpectedConnectorRevision {
			return nil, writeGroupRevisionConflict("migrated row-group canonical connector revision changed")
		}
		reviewed = append(reviewed, reviewedItem{item: item, existing: group, migrationMap: migrationMap})
	}

	groups := make([]*WriteGroup, 0, len(reviewed))
	created := false
	for _, entry := range reviewed {
		if entry.existing != nil {
			groups = append(groups, cloneWriteGroup(entry.existing))
			continue
		}
		candidate := cloneWriteGroup(entry.item.CandidateGroup)
		candidate.ID = ""
		candidate.Revision = ""
		candidate.AppliedRevision = ""
		candidate.Status = WriteGroupStatusDraft
		candidate.WorkspaceID = record.ID
		candidate.Destination.ConnectorRevision = request.ExpectedConnectorRevision
		candidate.Migration.ReviewResult = writeGroupMigrationReviewResultConfirmed
		candidate.Migration.AdapterVersion = preview.AdapterVersion
		candidate.Migration.SourceRevision = entry.item.SourceRevision
		candidate.Migration.LegacyRowGroupID = entry.item.SourceID
		candidate.Migration.SourceIDs = slices.Clone(entry.item.CandidateGroup.Migration.SourceIDs)
		candidate.Migration.TargetMappingPoints = maps.Clone(entry.item.CandidateGroup.Migration.TargetMappingPoints)
		if err := s.repo.CreateInTx(ctx, tx, candidate); err != nil {
			return nil, normalizeWriteGroupServiceError("save migrated row-group write group", err)
		}
		if err := projectRowGroupWriteGroup(record, candidate, entry.item.SourceID); err != nil {
			return nil, err
		}
		beforeIntent, err := json.Marshal(entry.item.BeforeRowGroupIntent)
		if err != nil {
			return nil, fmt.Errorf("encode row-group migration before intent: %w", err)
		}
		now := s.now()
		if err := s.repo.insertMigrationMapInTx(ctx, tx, writeGroupMigrationMap{
			WorkspaceID:    record.ID,
			SourceKind:     candidate.Migration.SourceKind,
			SourceID:       entry.item.SourceID,
			SourceRevision: entry.item.SourceRevision,
			GroupID:        candidate.ID,
			BeforeIntent:   string(beforeIntent),
			ReviewDigest:   request.ReviewDigest,
			AdapterVersion: preview.AdapterVersion,
			CreatedAt:      now,
			UpdatedAt:      now,
		}); err != nil {
			return nil, err
		}
		groups = append(groups, cloneWriteGroup(candidate))
		created = true
	}

	if created {
		record.DatabaseSetupRevision = s.workspaceSvc.newID()
		record.UpdatedAt = s.workspaceSvc.now()
		if err := setupTx.Save(ctx, record); err != nil {
			return nil, fmt.Errorf("save row-group migration review: %w", err)
		}
		if err := setupTx.Commit(); err != nil {
			return nil, fmt.Errorf("commit row-group migration review: %w", err)
		}
		committed = true
	}
	return &WriteGroupMigrationReviewResult{
		WorkspaceID:       record.ID,
		WorkspaceRevision: record.DatabaseSetupRevision,
		ConnectorRevision: connector.identityRevision,
		Groups:            groups,
	}, nil
}

func validateRowGroupMigrationMap(
	migrationMap *writeGroupMigrationMap,
	workspaceID string,
	item WriteGroupMigrationPreviewItem,
	adapterVersion string,
) error {
	if migrationMap == nil ||
		strings.TrimSpace(migrationMap.WorkspaceID) != workspaceID ||
		strings.TrimSpace(migrationMap.SourceKind) != writeGroupMigrationSourceKindLegacyRowGroup ||
		strings.TrimSpace(migrationMap.SourceID) != strings.TrimSpace(item.SourceID) ||
		strings.TrimSpace(migrationMap.SourceRevision) == "" ||
		strings.TrimSpace(migrationMap.GroupID) == "" ||
		strings.TrimSpace(migrationMap.BeforeIntent) == "" ||
		strings.TrimSpace(migrationMap.ReviewDigest) == "" ||
		strings.TrimSpace(migrationMap.AdapterVersion) != strings.TrimSpace(adapterVersion) {
		return ErrWriteGroupValidation
	}
	var beforeIntent WriteGroupRowGroupMigrationIntent
	if err := json.Unmarshal([]byte(migrationMap.BeforeIntent), &beforeIntent); err != nil {
		return ErrWriteGroupValidation
	}
	if strings.TrimSpace(beforeIntent.SourceID) != strings.TrimSpace(item.SourceID) ||
		strings.TrimSpace(beforeIntent.ConnectorID) == "" ||
		strings.TrimSpace(beforeIntent.ConnectorRevision) == "" ||
		beforeIntent.RowGroup.ID != item.SourceID {
		return ErrWriteGroupValidation
	}
	return nil
}

func validateReviewedRowGroupMigration(
	group *WriteGroup,
	migrationMap *writeGroupMigrationMap,
	workspaceID string,
	item WriteGroupMigrationPreviewItem,
) error {
	if group == nil || migrationMap == nil ||
		strings.TrimSpace(group.ID) != strings.TrimSpace(migrationMap.GroupID) ||
		strings.TrimSpace(group.WorkspaceID) != strings.TrimSpace(workspaceID) ||
		strings.TrimSpace(group.Revision) == "" ||
		strings.TrimSpace(group.Migration.SourceKind) != writeGroupMigrationSourceKindLegacyRowGroup ||
		strings.TrimSpace(group.Migration.LegacyRowGroupID) != strings.TrimSpace(item.SourceID) ||
		strings.TrimSpace(group.Migration.SourceRevision) != strings.TrimSpace(migrationMap.SourceRevision) ||
		strings.TrimSpace(group.Migration.AdapterVersion) != strings.TrimSpace(migrationMap.AdapterVersion) ||
		group.Migration.ReviewResult != writeGroupMigrationReviewResultConfirmed ||
		!slices.Equal(group.Migration.SourceIDs, item.CandidateGroup.Migration.SourceIDs) ||
		!maps.Equal(group.Migration.TargetMappingPoints, item.CandidateGroup.Migration.TargetMappingPoints) {
		return fmt.Errorf("canonical row-group migration identity is incomplete: %w", ErrWriteGroupValidation)
	}
	if group.Status == WriteGroupStatusDeleted {
		return nil
	}
	return nil
}
