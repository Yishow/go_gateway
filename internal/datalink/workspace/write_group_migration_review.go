package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const writeGroupMigrationReviewResultConfirmed = "confirmed_snapshot_conversion"

// WriteGroupMigrationReviewRequest binds a reviewed single-mapping snapshot
// conversion to the exact workspace, connector, source revisions, and preview
// digest the operator inspected.
type WriteGroupMigrationReviewRequest struct {
	WorkspaceID               string   `json:"workspace_id"`
	ExpectedWorkspaceRevision string   `json:"expected_workspace_revision"`
	ExpectedConnectorRevision string   `json:"expected_connector_revision"`
	ReviewDigest              string   `json:"review_digest"`
	SourceIDs                 []string `json:"source_ids"`
	ConfirmSnapshotConversion bool     `json:"confirm_snapshot_conversion"`
}

// WriteGroupMigrationReviewResult returns the canonical draft groups created
// or reused by a reviewed single-mapping snapshot conversion.
type WriteGroupMigrationReviewResult struct {
	WorkspaceID       string        `json:"workspace_id"`
	WorkspaceRevision string        `json:"workspace_revision"`
	ConnectorRevision string        `json:"connector_revision"`
	Groups            []*WriteGroup `json:"groups"`
}

// ReviewSingleMappingMigration persists a reviewed single-mapping snapshot
// conversion. The implementation is intentionally added in the review
// domain file so API callers have a stable seam while preview remains read
// only.
func (s *WriteGroupService) ReviewSingleMappingMigration(
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
		return nil, normalizeWriteGroupServiceError("begin migration review", err)
	}
	committed := false
	defer func() {
		if committed {
			return
		}
		if rollbackErr := setupTx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			result = nil
			err = errors.Join(err, fmt.Errorf("rollback migration review: %w", rollbackErr))
		}
	}()

	record, err := setupTx.Get(ctx)
	if errors.Is(err, ErrNotFound) {
		return nil, writeGroupNotFound("review migration workspace")
	}
	if err != nil {
		return nil, normalizeWriteGroupServiceError("read migration review workspace", err)
	}
	if record.ID != request.WorkspaceID {
		return nil, writeGroupNotFound("review migration workspace")
	}
	if record.DatabaseSetupRevision != request.ExpectedWorkspaceRevision {
		return nil, writeGroupRevisionConflict("migration review workspace revision is stale")
	}

	tx := setupTx.SQLTx()
	connector, err := readMigrationPreviewConnector(ctx, tx, s.repo, record.DatabaseConnectorID)
	if err != nil {
		return nil, err
	}
	if connector.identityRevision != request.ExpectedConnectorRevision {
		return nil, writeGroupRevisionConflict("migration review connector revision is stale")
	}
	preview, err := buildSingleMappingMigrationPreview(ctx, tx, s.repo, record, connector, ids)
	if err != nil {
		return nil, normalizeWriteGroupServiceError("recompute migration review", err)
	}
	if preview.ReviewDigest != request.ReviewDigest {
		return nil, writeGroupRevisionConflict("migration review digest is stale")
	}

	type reviewedItem struct {
		item         WriteGroupMigrationPreviewItem
		existing     *WriteGroup
		migrationMap *writeGroupMigrationMap
	}
	reviewed := make([]reviewedItem, 0, len(preview.Items))
	for _, item := range preview.Items {
		if item.Status == WriteGroupMigrationPreviewStatusBlocked {
			return nil, fmt.Errorf("migration review source %s is blocked: %w", item.SourceID, ErrWriteGroupValidation)
		}
		if item.Status != WriteGroupMigrationPreviewStatusNeedsReview || item.CandidateGroup == nil || item.BeforeIntent == nil || item.SourceRevision == "" {
			return nil, fmt.Errorf("migration review source %s is not reviewable: %w", item.SourceID, ErrWriteGroupValidation)
		}
		migrationMap, found, err := s.repo.getMigrationMapInTx(ctx, tx, record.ID, item.CandidateGroup.Migration.SourceKind, item.SourceID)
		if err != nil {
			return nil, err
		}
		if !found {
			reviewed = append(reviewed, reviewedItem{item: item})
			continue
		}
		if migrationMap.SourceRevision != item.SourceRevision {
			return nil, writeGroupRevisionConflict("migration source revision changed after prior review")
		}
		if err := validateWriteGroupMigrationMap(migrationMap, record.ID, item, preview.AdapterVersion); err != nil {
			return nil, fmt.Errorf("migration identity map is incomplete: %w", ErrWriteGroupValidation)
		}
		group, err := s.repo.GetInTx(ctx, tx, record.ID, migrationMap.GroupID)
		if err != nil {
			return nil, fmt.Errorf("load migration identity map group: %w", ErrWriteGroupValidation)
		}
		if err := validateReviewedMigrationGroup(group, migrationMap, record.ID, item); err != nil {
			return nil, err
		}
		if group.Destination.ConnectorID != connector.id || group.Destination.ConnectorRevision != request.ExpectedConnectorRevision {
			return nil, writeGroupRevisionConflict("migrated canonical connector revision changed")
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
		candidate.Migration.SourceIDs = []string{entry.item.SourceID}
		if err := s.repo.CreateInTx(ctx, tx, candidate); err != nil {
			return nil, normalizeWriteGroupServiceError("save migrated write group", err)
		}
		if err := projectWriteGroup(record, candidate); err != nil {
			return nil, err
		}
		beforeIntent, err := json.Marshal(entry.item.BeforeIntent)
		if err != nil {
			return nil, fmt.Errorf("encode migration before intent: %w", err)
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
			return nil, fmt.Errorf("save workspace migration review: %w", err)
		}
		if err := setupTx.Commit(); err != nil {
			return nil, fmt.Errorf("commit migration review: %w", err)
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

func prepareWriteGroupMigrationReviewRequest(
	request WriteGroupMigrationReviewRequest,
) (WriteGroupMigrationReviewRequest, []string, error) {
	request.WorkspaceID = strings.TrimSpace(request.WorkspaceID)
	request.ExpectedWorkspaceRevision = strings.TrimSpace(request.ExpectedWorkspaceRevision)
	request.ExpectedConnectorRevision = strings.TrimSpace(request.ExpectedConnectorRevision)
	request.ReviewDigest = strings.TrimSpace(request.ReviewDigest)
	if request.WorkspaceID == "" || request.ExpectedConnectorRevision == "" || request.ReviewDigest == "" {
		return WriteGroupMigrationReviewRequest{}, nil, fmt.Errorf("migration review revision envelope: %w", ErrWriteGroupValidation)
	}
	if !request.ConfirmSnapshotConversion {
		return WriteGroupMigrationReviewRequest{}, nil, fmt.Errorf("migration review requires explicit snapshot confirmation: %w", ErrWriteGroupValidation)
	}
	ids, err := normalizeMigrationPreviewIDs(request.SourceIDs)
	if err != nil {
		return WriteGroupMigrationReviewRequest{}, nil, err
	}
	request.SourceIDs = append([]string(nil), ids...)
	return request, ids, nil
}

func validateWriteGroupMigrationMap(
	migrationMap *writeGroupMigrationMap,
	workspaceID string,
	item WriteGroupMigrationPreviewItem,
	adapterVersion string,
) error {
	if migrationMap == nil ||
		strings.TrimSpace(migrationMap.WorkspaceID) != workspaceID ||
		strings.TrimSpace(migrationMap.SourceKind) != strings.TrimSpace(item.CandidateGroup.Migration.SourceKind) ||
		strings.TrimSpace(migrationMap.SourceID) != strings.TrimSpace(item.SourceID) ||
		strings.TrimSpace(migrationMap.SourceRevision) == "" ||
		strings.TrimSpace(migrationMap.GroupID) == "" ||
		strings.TrimSpace(migrationMap.BeforeIntent) == "" ||
		strings.TrimSpace(migrationMap.ReviewDigest) == "" ||
		strings.TrimSpace(migrationMap.AdapterVersion) != strings.TrimSpace(adapterVersion) {
		return ErrWriteGroupValidation
	}
	var beforeIntent WriteGroupMigrationIntent
	if err := json.Unmarshal([]byte(migrationMap.BeforeIntent), &beforeIntent); err != nil {
		return ErrWriteGroupValidation
	}
	if strings.TrimSpace(beforeIntent.SourceID) != strings.TrimSpace(item.SourceID) ||
		strings.TrimSpace(beforeIntent.ConnectorID) == "" ||
		strings.TrimSpace(beforeIntent.ConnectorRevision) == "" {
		return ErrWriteGroupValidation
	}
	return nil
}

func validateReviewedMigrationGroup(
	group *WriteGroup,
	migrationMap *writeGroupMigrationMap,
	workspaceID string,
	item WriteGroupMigrationPreviewItem,
) error {
	if group == nil || migrationMap == nil ||
		strings.TrimSpace(group.ID) != strings.TrimSpace(migrationMap.GroupID) ||
		strings.TrimSpace(group.WorkspaceID) != strings.TrimSpace(workspaceID) ||
		strings.TrimSpace(group.Revision) == "" ||
		strings.TrimSpace(group.Migration.SourceKind) != strings.TrimSpace(migrationMap.SourceKind) ||
		len(group.Migration.SourceIDs) != 1 ||
		strings.TrimSpace(group.Migration.SourceIDs[0]) != strings.TrimSpace(item.SourceID) ||
		strings.TrimSpace(group.Migration.SourceRevision) != strings.TrimSpace(migrationMap.SourceRevision) ||
		strings.TrimSpace(group.Migration.AdapterVersion) != strings.TrimSpace(migrationMap.AdapterVersion) ||
		group.Migration.ReviewResult != writeGroupMigrationReviewResultConfirmed {
		return fmt.Errorf("canonical migration identity is incomplete: %w", ErrWriteGroupValidation)
	}
	return nil
}
