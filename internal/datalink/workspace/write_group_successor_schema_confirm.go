package workspace

import (
	"context"
	"database/sql"
	"reflect"

	"go-gateway/internal/datalink/recordingplan"
)

// ConfirmSuccessorSchema saves verified schema evidence for one reservation.
// It never applies the group, transfers the head, or enables collection.
func (s *WriteGroupService) ConfirmSuccessorSchema(ctx context.Context, id string, req WriteGroupSuccessorStageRequest, operationID, digest string, storage *StorageEvidence) (*WriteGroupSaveResult, error) {
	if operationID == "" || operationID == req.Envelope.Request.OperationID || len(digest) != 64 || storage == nil {
		return nil, ErrWriteGroupValidation
	}
	if err := s.validate(); err != nil {
		return nil, err
	}
	current, err := s.Get(ctx, req.Current.GroupID)
	if err != nil {
		return nil, err
	}
	if current.Group.Destination.SchemaRevision == operationID && current.Group.Destination.SchemaDigest == digest {
		// Only exact already-persisted proof permits read-only replay after the
		// workspace revision advanced. All other current scope facts still bind.
		replay := req
		replay.Current.WorkspaceRevision = current.WorkspaceRevision
		snapshot, err := s.readSuccessorSchemaSnapshotChecked(ctx, id, replay, func(tx *sql.Tx, snapshot *successorSchemaSnapshot) error {
			return s.verifySuccessorSchemaProof(ctx, tx, snapshot, operationID, digest, storage)
		})
		if err != nil {
			return nil, err
		}
		return &WriteGroupSaveResult{WorkspaceRevision: snapshot.reserved.WorkspaceRevision, Group: snapshot.reserved.Group}, nil
	}
	// Fresh external capability work is outside the configuration transaction.
	// The exact replay above remains a read of the already persisted proof.
	before, err := s.readSuccessorSchemaSnapshot(ctx, id, req)
	if err != nil {
		return nil, err
	}
	var freshLayout *recordingplan.GroupSchemaLayout
	if before.dialect == writeGroupConnectorKindMySQL || before.dialect == writeGroupConnectorKindSQLServer {
		scope, e := s.SuccessorSchemaScope(ctx, id, req)
		if e != nil {
			return nil, e
		}
		freshLayout = scope.GroupLayout
	}
	var saved *WriteGroup
	updated, err := s.workspaceSvc.UpdateDatabaseSetup(ctx, req.Current.WorkspaceRevision, func(ctx context.Context, tx *sql.Tx, record *Record) error {
		snapshot, err := s.successorSchemaSnapshotInTx(ctx, tx, record, id, req)
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(before, snapshot) {
			return recordingplan.ErrPreviewTokenStale
		}
		if freshLayout != nil && !reflect.DeepEqual(freshLayout, &storage.Layout) {
			return recordingplan.ErrPreviewTokenStale
		}
		if err := s.verifySuccessorSchemaProof(ctx, tx, snapshot, operationID, digest, storage); err != nil {
			return err
		}
		group := snapshot.reserved.Group
		group.Destination.SchemaRevision, group.Destination.SchemaDigest = operationID, digest
		group.UpdatedAt = s.clockNow()
		if err := s.repo.updateGroup(ctx, tx, group); err != nil {
			return err
		}
		if managedLongGroup(group) {
			if err := recordingplan.RecordStorageConfirmation(ctx, tx, StorageScopeOf(group), recordingplan.StorageConfirmation{Layout: storage.Layout, ConnectorRevision: group.Destination.ConnectorRevision, OperationID: operationID, VerifiedDigest: digest, CreatedTable: storage.CreatedTable, GroupID: group.ID}, s.clockNow()); err != nil {
				return err
			}
		}
		saved = cloneWriteGroup(group)
		return nil
	})
	if err != nil {
		return nil, normalizeWriteGroupServiceError("confirm successor schema", err)
	}
	return &WriteGroupSaveResult{WorkspaceRevision: updated.DatabaseSetupRevision, Group: saved}, nil
}

// The existing durable ledger and canonical token verifier, not caller input,
// own complete schema proof. No second token decoder or scope authority lives here.
func (s *WriteGroupService) verifySuccessorSchemaProof(ctx context.Context, tx *sql.Tx, snapshot *successorSchemaSnapshot, operationID, digest string, storage *StorageEvidence) error {
	group := snapshot.reserved.Group
	source, err := successorSchemaSourceDigest(snapshot)
	if err != nil {
		return err
	}
	layout := &storage.Layout
	if snapshot.dialect == writeGroupConnectorKindMySQL || snapshot.dialect == writeGroupConnectorKindSQLServer {
		recorded, e := recordingplan.SucceededSchemaTokenInTx(ctx, tx, group.WorkspaceID, operationID, digest)
		if e != nil {
			return e
		}
		if snapshot.dialect == writeGroupConnectorKindSQLServer {
			layout, e = successorSQLServerProofLayoutInTx(ctx, tx, group, recorded.GroupLayout)
		} else {
			layout, e = successorMySQLProofLayout(group, snapshot.sourceFamily, recorded.GroupLayout)
		}
		if e != nil {
			return e
		}
		if !reflect.DeepEqual(layout, &storage.Layout) {
			return recordingplan.ErrPreviewTokenStale
		}
	}
	scope := recordingplan.SchemaPreviewScope{WorkspaceID: group.WorkspaceID, WorkspaceRevision: snapshot.reserved.WorkspaceRevision, PlanID: group.ID, PlanRevision: group.Revision, ConnectorID: group.Destination.ConnectorID, ConnectorRevision: group.Destination.ConnectorRevision, Dialect: snapshot.dialect, Database: group.Destination.Database, Schema: group.Destination.TableSchema, GroupLayout: layout, SourceDigest: source, SchemaRevision: group.Destination.SchemaRevision, SchemaDigest: group.Destination.SchemaDigest}
	if managedLongGroup(group) && group.WritePolicy.RecreateMissing != nil {
		v := *group.WritePolicy.RecreateMissing
		scope.RecreateMissing = &v
	}
	token, err := recordingplan.VerifiedSchemaProofInTx(ctx, tx, scope, operationID, digest)
	if err != nil {
		return err
	}
	if token.GroupLayout == nil || token.GroupLayout.TableName != group.Destination.TableName {
		return ErrWriteGroupValidation
	}
	created := false
	for _, table := range token.Tables {
		if table.Name == group.Destination.TableName && table.Action == recordingplan.SchemaTableActionCreate {
			created = true
		}
	}
	if created != storage.CreatedTable {
		return ErrWriteGroupValidation
	}
	return nil
}
