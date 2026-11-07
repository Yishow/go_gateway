package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
)

func (s *WriteGroupService) permanentDelete(ctx context.Context, id string, mutation WriteGroupMutation) (*WriteGroupSaveResult, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("delete write group: %w", ErrWriteGroupValidation)
	}
	mutation, err := prepareWriteGroupMutation(mutation, true, false)
	if err != nil {
		return nil, err
	}
	var result *WriteGroupSaveResult
	updated, err := s.workspaceSvc.UpdateDatabaseSetup(ctx, mutation.ExpectedWorkspaceRevision, func(ctx context.Context, tx *sql.Tx, record *Record) error {
		if record.ID != mutation.WorkspaceID {
			return writeGroupNotFound("delete write group")
		}
		group, err := s.repo.GetInTx(ctx, tx, record.ID, id)
		if err != nil {
			return err
		}
		if group.Revision != mutation.ExpectedGroupRevision {
			return writeGroupRevisionConflict("write group revision is stale")
		}
		if group.Destination.ConnectorRevision != mutation.ExpectedConnectorRevision {
			return writeGroupRevisionConflict("write group connector revision is stale")
		}
		scope, err := s.loadPermanentDeleteScope(ctx, tx, group)
		if err != nil {
			return err
		}
		if err := validateLiveWriteGroupConnector(ctx, tx, s.repo, group.Destination.ConnectorID, mutation.ExpectedConnectorRevision); err != nil {
			reason := writeGroupReasonConnectorUnavailable
			if errors.Is(err, ErrWriteGroupConnectorRevisionConflict) {
				reason = writeGroupReasonConnectorIdentityStale
			}
			return writeGroupRemoveActionError(group, reason, err)
		}
		result, err = s.permanentDeletePreparedInTx(ctx, tx, record, group, scope)
		return err
	})
	if err != nil {
		return nil, normalizeWriteGroupServiceError("delete write group", err)
	}
	result.WorkspaceRevision = updated.DatabaseSetupRevision
	s.recordPermanentDeleteAudit(ctx, result)
	return result, nil
}

// permanentDeleteInTx shares the caller's transaction with workspace changes.
func (s *WriteGroupService) permanentDeleteInTx(ctx context.Context, tx *sql.Tx, record *Record, group *WriteGroup) (*WriteGroupSaveResult, error) {
	scope, err := s.loadPermanentDeleteScope(ctx, tx, group)
	if err != nil {
		return nil, err
	}
	return s.permanentDeletePreparedInTx(ctx, tx, record, group, scope)
}

func (s *WriteGroupService) permanentDeletePreparedInTx(ctx context.Context, tx *sql.Tx, record *Record, group *WriteGroup, scope *permanentDeleteScope) (*WriteGroupSaveResult, error) {
	if err := s.checkPermanentDeleteOperationsForScope(ctx, tx, group, scope.ids, scope.lineage); err != nil {
		return nil, err
	}
	if err := s.checkPermanentDeleteSharedProof(ctx, tx, group, scope.ids); err != nil {
		return nil, err
	}
	if err := s.checkPermanentDeleteIntakeAndBacklog(ctx, tx, scope.groups); err != nil {
		return nil, err
	}
	for _, generation := range scope.groups {
		if generation.AppliedRevision != "" {
			if s.backlogGuard == nil {
				return nil, writeGroupRemoveActionError(generation, writeGroupReasonOwnershipUnavailable, ErrWriteGroupLifecycleBlocked)
			}
			if err := s.backlogGuard.CheckWriteGroupBacklog(ctx, tx, generation); err != nil {
				return nil, writeGroupRemoveActionError(generation, writeGroupReasonOwnershipUnavailable, errors.Join(ErrWriteGroupLifecycleBlocked, err))
			}
		}
	}
	lineageID := ""
	if scope.lineage != nil {
		lineageID = scope.lineage.LineageID
	}
	if err := s.purgeWriteGroupLineageInTx(ctx, tx, scope.groups, lineageID); err != nil {
		return nil, err
	}
	projectionIDs := make(map[string]bool, len(scope.ids))
	for _, generation := range scope.groups {
		projectionIDs[generation.ID] = true
		if generation.Migration.SourceKind == writeGroupMigrationSourceKindLegacyRowGroup {
			projectionIDs[generation.Migration.LegacyRowGroupID] = true
		}
	}
	record.DatabaseRowGroups = slices.DeleteFunc(record.DatabaseRowGroups, func(row DatabaseRowGroup) bool { return projectionIDs[row.ID] })
	record.DatabaseTargetRefs = slices.DeleteFunc(record.DatabaseTargetRefs, func(ref DatabaseTargetRef) bool { return projectionIDs[ref.RowGroupID] })
	return &WriteGroupSaveResult{Group: cloneWriteGroup(group), Deleted: true, GroupID: group.ID, LineageID: lineageID, DeletedGenerations: len(scope.groups)}, nil
}
