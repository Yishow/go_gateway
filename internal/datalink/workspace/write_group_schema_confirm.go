package workspace

import (
	"context"
	"database/sql"
	"fmt"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
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
	return s.ConfirmManagedSchemaWithStorage(ctx, id, mutation, operationID, digest, nil)
}

// StorageEvidence is what a verified confirmation of managed long storage
// proves about the shared table, recorded with the group's own confirmation.
type StorageEvidence struct {
	// Layout is the layout the operation verified, including the extension
	// schema of a hypertable.
	Layout recordingplan.GroupSchemaLayout
	// CreatedTable is true when the confirmed batch created the shared table.
	CreatedTable bool
}

// ConfirmManagedSchemaWithStorage records the verified result and, for a
// managed long group, freezes the storage lineage in the same local
// transaction: the verified layout, the connector revision and the operation
// that proved them. A table created again after an earlier confirmation starts
// the next generation and marks earlier history as possibly lost. The group's
// own confirmation is separate from every other group's: this never confirms
// anything on another group's behalf.
func (s *WriteGroupService) ConfirmManagedSchemaWithStorage(ctx context.Context, id string, mutation WriteGroupMutation, operationID, digest string, storage *StorageEvidence) (*WriteGroupSaveResult, error) {
	if operationID == "" || len(digest) != 64 {
		return nil, fmt.Errorf("managed schema proof is incomplete: %w", ErrWriteGroupValidation)
	}
	current, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if isMixedWriteGroup(current.Group) && managedLongGroup(current.Group) && storage == nil {
		return nil, ErrWriteGroupValidation
	}
	if current.Group.WorkspaceID != mutation.WorkspaceID {
		return nil, writeGroupNotFound("confirm group schema")
	}
	if err := s.requireUnreservedWriteGroup(ctx, mutation.WorkspaceID, id); err != nil {
		return nil, err
	}
	if current.Group.Destination.SchemaRevision == operationID && current.Group.Destination.SchemaDigest == digest {
		if isMixedWriteGroup(current.Group) {
			return s.replayMixedSchemaConfirmation(ctx, id, mutation, operationID, digest, storage)
		}
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
		if err := s.repo.requireUnreservedWriteGroupInTx(ctx, tx, record.ID, id); err != nil {
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
		if isMixedWriteGroup(group) {
			if err := s.verifyMixedSchemaProofInTx(ctx, tx, record, group, operationID, digest, storage); err != nil {
				return err
			}
		}
		group.Destination.SchemaRevision = operationID
		group.Destination.SchemaDigest = digest
		group.UpdatedAt = s.clockNow()
		if err := s.repo.updateGroup(ctx, tx, group); err != nil {
			return err
		}
		if storage != nil && managedLongGroup(group) {
			err := recordingplan.RecordStorageConfirmation(ctx, tx, StorageScopeOf(group), recordingplan.StorageConfirmation{
				Layout: storage.Layout, ConnectorRevision: group.Destination.ConnectorRevision, OperationID: operationID,
				VerifiedDigest: digest, CreatedTable: storage.CreatedTable, GroupID: group.ID,
			}, s.clockNow())
			if err != nil {
				return fmt.Errorf("record storage lineage: %w", err)
			}
		}
		saved = cloneWriteGroup(group)
		return nil
	})
	if err != nil {
		return nil, normalizeWriteGroupServiceError("confirm group schema", err)
	}
	return &WriteGroupSaveResult{WorkspaceRevision: updated.DatabaseSetupRevision, Group: saved}, nil
}
