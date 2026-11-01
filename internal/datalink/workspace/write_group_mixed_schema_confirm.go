package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"go-gateway/internal/datalink/recordingplan"
	"reflect"
)

// The existing operation ledger remains the only proof authority. Mixed
// confirmation cannot relabel an old successful operation using a new CAS.
func (s *WriteGroupService) verifyMixedSchemaProofInTx(ctx context.Context, tx *sql.Tx, record *Record, group *WriteGroup, operationID, digest string, storage *StorageEvidence) error {
	types, err := loadWriteGroupTagTypes(ctx, tx, s.repo.query, group)
	if err != nil {
		return err
	}
	var dialect string
	if err := tx.QueryRowContext(ctx, s.repo.query(`SELECT kind FROM database_connectors WHERE id=$1`), group.Destination.ConnectorID).Scan(&dialect); err != nil {
		return err
	}
	layout, err := managedSchemaLayout(group, types, dialect)
	if err != nil {
		return err
	}
	proofs, issues := s.counterBindingEvidenceInTx(ctx, tx, group, types, nil)
	if len(issues) > 0 {
		return errors.Join(ErrWriteGroupValidation, ErrCounterBindingProofConflict)
	}
	source, err := managedSourceDigest(group, types, proofs)
	if err != nil {
		return err
	}
	// Ordinary confirmations retain their original digest contract. A different
	// digest is accepted only as the authenticated original successor reservation.
	recorded, err := recordingplan.SucceededSchemaTokenInTx(ctx, tx, record.ID, operationID, digest)
	if err != nil {
		return err
	}
	if recorded.SourceDigest != source {
		source, err = s.mixedSuccessorSchemaSourceInTx(ctx, tx, group, source, recorded)
		if err != nil {
			return err
		}
	}
	scope := recordingplan.SchemaPreviewScope{WorkspaceID: record.ID, WorkspaceRevision: record.DatabaseSetupRevision, PlanID: group.ID, PlanRevision: group.Revision, ConnectorID: group.Destination.ConnectorID, ConnectorRevision: group.Destination.ConnectorRevision, Dialect: dialect, Database: group.Destination.Database, Schema: group.Destination.TableSchema, GroupLayout: layout, SourceDigest: source, SchemaRevision: group.Destination.SchemaRevision, SchemaDigest: group.Destination.SchemaDigest}
	if managedLongGroup(group) && group.WritePolicy.RecreateMissing != nil {
		v := *group.WritePolicy.RecreateMissing
		scope.RecreateMissing = &v
	}
	token, err := recordingplan.VerifiedSchemaProofInTx(ctx, tx, scope, operationID, digest)
	if err != nil {
		return err
	}
	if storage == nil {
		return nil
	}
	if !reflect.DeepEqual(layout, &storage.Layout) {
		return recordingplan.ErrPreviewTokenStale
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

// Exact replay is a verified read, not a new workspace update or new schema
// consent. All scope facts except the original operation's own advancement bind.
func (s *WriteGroupService) replayMixedSchemaConfirmation(ctx context.Context, id string, mutation WriteGroupMutation, operationID, digest string, storage *StorageEvidence) (result *WriteGroupSaveResult, err error) {
	store, ok := s.workspaceSvc.repo.(DatabaseSetupStore)
	if !ok {
		return nil, ErrWriteGroupReadinessUnavailable
	}
	setup, err := store.BeginDatabaseSetup(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := setup.Rollback(); e != nil && !errors.Is(e, sql.ErrTxDone) {
			result = nil
			err = errors.Join(err, e)
		}
	}()
	record, err := setup.Get(ctx)
	if err != nil {
		return nil, err
	}
	group, err := s.repo.GetInTx(ctx, setup.SQLTx(), mutation.WorkspaceID, id)
	if err != nil {
		return nil, err
	}
	if record.ID != mutation.WorkspaceID || group.Revision != mutation.ExpectedGroupRevision || group.Destination.ConnectorRevision != mutation.ExpectedConnectorRevision || !isMixedWriteGroup(group) {
		return nil, recordingplan.ErrPreviewTokenStale
	}
	if err := s.repo.requireUnreservedWriteGroupInTx(ctx, setup.SQLTx(), record.ID, id); err != nil {
		return nil, err
	}
	if err := s.repo.validate(ctx, setup.SQLTx(), group); err != nil {
		return nil, err
	}
	if group.Destination.SchemaRevision != operationID || group.Destination.SchemaDigest != digest {
		return nil, recordingplan.ErrPreviewTokenStale
	}
	if err := s.verifyMixedSchemaProofInTx(ctx, setup.SQLTx(), record, group, operationID, digest, storage); err != nil {
		return nil, err
	}
	return &WriteGroupSaveResult{WorkspaceRevision: record.DatabaseSetupRevision, Group: cloneWriteGroup(group)}, nil
}

// Ledger details select a request ID only; ReadSuccessorReserveInTx and the
// original decoders authenticate it before it can supply any proof authority.
func (s *WriteGroupService) mixedSuccessorSchemaSourceInTx(ctx context.Context, tx *sql.Tx, group *WriteGroup, fresh string, token *recordingplan.SchemaPreviewToken) (string, error) {
	meta, found, err := s.repo.loadWriteGroupLineage(ctx, tx, group.WorkspaceID, group.ID)
	if err != nil {
		return "", err
	}
	if !found || meta.CurrentGroupID != group.ID || group.AppliedRevision == "" {
		return "", recordingplan.ErrPreviewTokenStale
	}
	rows, err := tx.QueryContext(ctx, s.repo.query(`SELECT detail FROM managed_schema_operations WHERE workspace_id=$1 AND action=$2`), group.WorkspaceID, recordingplan.SuccessorReserveAction)
	if err != nil {
		return "", err
	}
	var requestID string
	for rows.Next() {
		var raw []byte
		if err := rows.Scan(&raw); err != nil {
			rows.Close()
			return "", err
		}
		var selector struct {
			Scope struct {
				GroupID   string `json:"successor_group_id"`
				RequestID string `json:"request_id"`
			} `json:"scope"`
		}
		if json.Unmarshal(raw, &selector) != nil {
			rows.Close()
			return "", ErrWriteGroupLineageCorrupt
		}
		if selector.Scope.GroupID != group.ID {
			continue
		}
		if requestID != "" || selector.Scope.RequestID == "" {
			rows.Close()
			return "", ErrWriteGroupLineageCorrupt
		}
		requestID = selector.Scope.RequestID
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return "", err
	}
	if err := rows.Close(); err != nil {
		return "", err
	}
	if requestID == "" {
		return "", recordingplan.ErrPreviewTokenStale
	}
	root, err := recordingplan.ReadSuccessorReserveInTx(ctx, tx, group.WorkspaceID, requestID)
	if err != nil {
		return "", err
	}
	if root == nil || root.Scope.SuccessorGroupID != group.ID || root.Scope.LineageID != meta.LineageID {
		return "", ErrWriteGroupLineageCorrupt
	}
	_, reserved, err := decodeSuccessorReserveRecord(root)
	if err != nil {
		return "", err
	}
	binding, err := decodeSuccessorReserveSchemaBinding(root.SchemaSourceBinding, reserved)
	if err != nil {
		return "", err
	}
	if reserved.Group.Revision != group.Revision || reserved.Group.RecordID != group.RecordID || binding.ManagedSourceDigest != fresh {
		return "", recordingplan.ErrPreviewTokenStale
	}
	if err := successorStatusTokenScope(token, reserved, binding); err != nil {
		return "", err
	}
	return binding.SourceDigest, nil
}
