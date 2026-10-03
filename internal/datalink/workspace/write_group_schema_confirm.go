package workspace

import (
	"context"
	"database/sql"
	"fmt"

	"go-gateway/internal/datalink/dbtarget"
)

// ManagedTableInspector protects managed schema preparation from opening the
// gateway's own configuration or journal file as an external destination.
func (s *WriteGroupService) ManagedTableInspector(targets *dbtarget.ConnectorService) *dbtarget.ManagedTableInspector {
	if s == nil || s.repo == nil {
		return nil
	}
	return dbtarget.NewManagedTableInspector(targets, s.repo.db)
}

// ConfirmManagedSchema records only the verified result of the existing
// schema ledger. It does not Apply the group or change an accepted snapshot.
func (s *WriteGroupService) ConfirmManagedSchema(ctx context.Context, id string, mutation WriteGroupMutation, operationID, digest string) (*WriteGroupSaveResult, error) {
	if operationID == "" || len(digest) != 64 {
		return nil, fmt.Errorf("managed schema proof is incomplete: %w", ErrWriteGroupValidation)
	}
	current, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if current.Group.WorkspaceID != mutation.WorkspaceID {
		return nil, writeGroupNotFound("confirm group schema")
	}
	if current.Group.Destination.SchemaRevision == operationID && current.Group.Destination.SchemaDigest == digest {
		return current, nil
	}
	var saved *WriteGroup
	updated, err := s.workspaceSvc.UpdateDatabaseSetup(ctx, mutation.ExpectedWorkspaceRevision, func(ctx context.Context, tx *sql.Tx, record *Record) error {
		if record.ID != mutation.WorkspaceID {
			return writeGroupNotFound("confirm group schema")
		}
		group, err := s.repo.GetInTx(ctx, tx, record.ID, id)
		if err != nil {
			return err
		}
		if group.Revision != mutation.ExpectedGroupRevision || group.Destination.ConnectorRevision != mutation.ExpectedConnectorRevision {
			return writeGroupRevisionConflict("managed schema scope changed")
		}
		if group.Status == WriteGroupStatusDeleted || group.Destination.StorageStrategy != WriteGroupStorageStrategyManaged {
			return ErrWriteGroupLifecycleBlocked
		}
		if err := s.repo.validate(ctx, tx, group); err != nil {
			return err
		}
		group.Destination.SchemaRevision = operationID
		group.Destination.SchemaDigest = digest
		group.UpdatedAt = s.clockNow()
		if err := s.repo.updateGroup(ctx, tx, group); err != nil {
			return err
		}
		saved = cloneWriteGroup(group)
		return nil
	})
	if err != nil {
		return nil, normalizeWriteGroupServiceError("confirm group schema", err)
	}
	return &WriteGroupSaveResult{WorkspaceRevision: updated.DatabaseSetupRevision, Group: saved}, nil
}
