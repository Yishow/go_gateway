package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (s *WriteGroupService) reviewRecordingPlanMigration(
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
		return nil, normalizeWriteGroupServiceError("begin recording plan migration review", err)
	}
	defer func() {
		if rollbackErr := setupTx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			result = nil
			err = errors.Join(err, fmt.Errorf("rollback recording plan migration review: %w", rollbackErr))
		}
	}()
	record, err := setupTx.Get(ctx)
	if errors.Is(err, ErrNotFound) {
		return nil, writeGroupNotFound("review recording plan migration workspace")
	}
	if err != nil {
		return nil, normalizeWriteGroupServiceError("read recording plan migration workspace", err)
	}
	if record.ID != request.WorkspaceID {
		return nil, writeGroupNotFound("review recording plan migration workspace")
	}
	if record.DatabaseSetupRevision != request.ExpectedWorkspaceRevision {
		return nil, writeGroupRevisionConflict("recording plan migration workspace revision is stale")
	}
	connector, err := readRecordingPlanMigrationConnector(ctx, setupTx.SQLTx(), s.repo, record.DatabaseConnectorID)
	if err != nil {
		return nil, err
	}
	if connector.identityRevision != request.ExpectedConnectorRevision {
		return nil, writeGroupRevisionConflict("recording plan migration connector revision is stale")
	}
	preview, err := buildRecordingPlanMigrationPreview(ctx, setupTx.SQLTx(), s.repo, record, connector, ids)
	if err != nil {
		return nil, normalizeWriteGroupServiceError("recompute recording plan migration review", err)
	}
	if preview.ReviewDigest != request.ReviewDigest {
		return nil, writeGroupRevisionConflict("recording plan migration review digest is stale")
	}
	for _, item := range preview.Items {
		if item.Status == WriteGroupMigrationPreviewStatusBlocked || item.CandidateGroup == nil {
			return nil, fmt.Errorf("recording plan %s: %w", item.SourceID, errors.Join(ErrWriteGroupPlanMigrationBlocked, ErrWriteGroupValidation))
		}
	}
	return nil, fmt.Errorf("recording plan migration review is unsupported: %w", errors.Join(ErrWriteGroupPlanMigrationBlocked, ErrWriteGroupValidation))
}
