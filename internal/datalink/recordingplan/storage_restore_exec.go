package recordingplan

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go-gateway/internal/datalink/common"
)

// claimStorageRepair claims the next attempt for the scope's current
// generation. The attempt number is the count of earlier operations, so every
// attempt has its own identity and no finished attempt is ever run twice.
func (s *Service) claimStorageRepair(ctx context.Context, repo storageRepairRepository, req StorageRestoreRequest, lineage *StorageLineage, layout GroupSchemaLayout) (*SchemaOperation, ClaimOutcome, error) {
	prefix := fmt.Sprintf("repair:%s:%d:", storageScopeFingerprint(req.Scope), lineage.Generation)
	attempt, err := repo.countStorageRepairAttempts(ctx, req.Scope.WorkspaceID, prefix)
	if err != nil {
		return nil, "", err
	}
	token := fmt.Sprintf("%s%d", prefix, attempt)
	sum := sha256.Sum256([]byte(token))
	detail, err := json.Marshal(storageRepairDetail{
		Kind: StorageRepairAction, Generation: lineage.Generation, ToGeneration: lineage.Generation + 1,
		GroupID: req.GroupID, ConnectorID: req.Scope.ConnectorID, Table: layout.TableName,
		Scope: req.Scope, Dialect: req.Dialect, OwnerColumn: layout.OwnerColumn,
		SQLServer: req.SQLServer,
	})
	if err != nil {
		return nil, "", fmt.Errorf("encode storage repair detail: %w", err)
	}
	owner, err := common.NewUUID()
	if err != nil {
		return nil, "", fmt.Errorf("create storage repair owner: %w", err)
	}
	now := time.Now().UTC()
	scopeKey := (&SchemaPreviewToken{
		WorkspaceID: req.Scope.WorkspaceID, ConnectorID: req.Scope.ConnectorID, Database: req.Scope.Database,
		Schema: req.Scope.TableSchema, GroupLayout: &layout,
	}).ScopeKey()
	return s.repo.ClaimSchemaOperation(ctx, &SchemaOperation{
		OperationID: "repair-" + hex.EncodeToString(sum[:16]), Token: token, WorkspaceID: req.Scope.WorkspaceID,
		ScopeKey: scopeKey, Owner: "claim-" + owner, Action: StorageRepairAction, Status: SchemaOperationRunning,
		Detail: string(detail), CreatedAt: now, UpdatedAt: now,
	})
}

func storageScopeFingerprint(scope StorageScope) string {
	sum := sha256.Sum256([]byte(scope.WorkspaceID + "\x1f" + scope.ConnectorID + "\x1f" + scope.Database + "\x1f" +
		scope.TableSchema + "\x1f" + scope.StorageMode + "\x1f" + scope.LayoutVersion))
	return hex.EncodeToString(sum[:8])
}

// runStorageRepair executes the claimed attempt: look again under the claim,
// run the statements, verify the catalog, and end the operation together with
// its lineage change.
func (s *Service) runStorageRepair(ctx context.Context, repo storageRepairRepository, req StorageRestoreRequest, lineage *StorageLineage, layout GroupSchemaLayout, op *SchemaOperation, statements []string, target StorageRestoreTarget) (StorageRestoreResult, error) {
	now := time.Now()
	// Look again under the claim: somebody may have put the table back, or the
	// destination may have changed, since the first look. Nothing is run unless
	// this second look still shows exactly a missing group table.
	again, refusal, err := inspectStorageForRestore(ctx, req, layout, target, true)
	if err != nil {
		finished, finishErr := repo.finishStorageRepair(ctx, op.OperationID, op.Owner, SchemaOperationResult{
			Status: SchemaOperationFailed, Reason: SchemaReasonUnverifiable, NextAction: storageRepairNextActionConfirm,
		}, nil)
		if finishErr != nil {
			return StorageRestoreResult{}, fmt.Errorf("%w: %w", ErrSchemaOperationUnacknowledged, finishErr)
		}
		_ = finished
		return StorageRestoreResult{}, err
	}
	if refusal != nil || again.pendingReceipt {
		code, reason := StorageCodeRepairBlocked, StorageReasonReceiptsMissing
		if refusal != nil {
			code, reason = refusal.code, refusal.reason
		}
		finished, finishErr := repo.finishStorageRepair(ctx, op.OperationID, op.Owner, SchemaOperationResult{
			Status: SchemaOperationFailed, Reason: reason, NextAction: storageRepairNextActionConfirm,
		}, func(ctx context.Context, q DBTX) error {
			return RecordStorageRefusal(ctx, q, req.Scope, lineage.Generation, req.GroupID, code, reason, now)
		})
		if finishErr != nil {
			return StorageRestoreResult{}, fmt.Errorf("%w: %w", ErrSchemaOperationUnacknowledged, finishErr)
		}
		return StorageRestoreResult{Outcome: StorageBlocked, Code: code, Reason: reason, Operation: finished}, nil
	}
	if !again.pendingGroup {
		if layout.SQLServer != nil {
			currentDigest := schemaTablesDigest(again.tables)
			if currentDigest != lineage.ConfirmedDigest {
				authorizedDigest, authErr := s.authorizedSQLServerDigest(ctx, req.Scope, lineage)
				if authErr != nil || currentDigest != authorizedDigest {
					finished, finishErr := repo.finishStorageRepair(ctx, op.OperationID, op.Owner, SchemaOperationResult{
						Status: SchemaOperationFailed, Reason: StorageReasonVerificationMismatch, NextAction: storageRepairNextActionOccupied,
					}, func(ctx context.Context, q DBTX) error {
						return BlockStorageRepair(ctx, q, req.Scope, lineage.Generation, StorageReasonVerificationMismatch, op.OperationID, req.GroupID, now)
					})
					if finishErr != nil {
						return StorageRestoreResult{}, fmt.Errorf("%w: %w", ErrSchemaOperationUnacknowledged, finishErr)
					}
					return StorageRestoreResult{Outcome: StorageBlocked, Code: StorageCodeRepairBlocked, Reason: StorageReasonVerificationMismatch, Operation: finished}, nil
				}
			}
		}
		finished, finishErr := repo.finishStorageRepair(ctx, op.OperationID, op.Owner, SchemaOperationResult{
			Status: SchemaOperationSucceeded, Reason: SchemaReasonRecovered, VerifiedDigest: schemaTablesDigest(again.tables),
		}, nil)
		if finishErr != nil {
			return StorageRestoreResult{}, fmt.Errorf("%w: %w", ErrSchemaOperationUnacknowledged, finishErr)
		}
		return StorageRestoreResult{Outcome: StorageNotMissing, Operation: finished}, nil
	}

	var execution SchemaExecution
	var execErr error
	switch {
	case layout.MySQL != nil:
		execution, execErr = s.executeMySQLStorageRepair(ctx, req, layout, op, statements, target)
	case layout.SQLServer != nil:
		execution, execErr = s.executeSQLServerStorageRepair(ctx, req, layout, op, statements, target)
	default:
		execution, execErr = runSchemaStatements(ctx, target.Execute, statements)
	}
	result, apply, code := s.judgeStorageRepair(ctx, req, lineage, layout, op, execution, execErr, target, now)
	if layout.MySQL != nil || layout.SQLServer != nil {
		result.Detail = op.Detail
	}
	finished, err := repo.finishStorageRepair(ctx, op.OperationID, op.Owner, result, apply)
	if err != nil {
		return StorageRestoreResult{Operation: withoutOwner(*op)}, fmt.Errorf("%w: %w", ErrSchemaOperationUnacknowledged, err)
	}
	if result.Status == SchemaOperationSucceeded {
		return StorageRestoreResult{Outcome: StorageRestored, Operation: finished}, nil
	}
	return StorageRestoreResult{Outcome: StorageBlocked, Code: code, Reason: result.Reason, Operation: finished}, nil
}

// judgeStorageRepair decides the honest outcome from the catalog after the
// statements ran, never from the adapter's own report. A result nobody can
// vouch for blocks further automatic restores of the scope until an explicit
// confirmation checks the table again.
func (s *Service) judgeStorageRepair(ctx context.Context, req StorageRestoreRequest, lineage *StorageLineage, layout GroupSchemaLayout, op *SchemaOperation, execution SchemaExecution, execErr error, target StorageRestoreTarget, now time.Time) (result SchemaOperationResult, apply func(context.Context, DBTX) error, code string) {
	if layout.SQLServer != nil {
		return s.judgeSQLServerStorageRepair(ctx, req, lineage, layout, op, execution, execErr, target, now)
	}
	result = SchemaOperationResult{ExecutedStatements: execution.Committed}
	block := func(reason, next string) (SchemaOperationResult, func(context.Context, DBTX) error, string) {
		result.Status, result.Reason, result.NextAction = SchemaOperationPartial, reason, next
		return result, func(ctx context.Context, q DBTX) error {
			return BlockStorageRepair(ctx, q, req.Scope, lineage.Generation, reason, op.OperationID, req.GroupID, now)
		}, StorageCodeRepairBlocked
	}
	failed := func(reason, next string) (SchemaOperationResult, func(context.Context, DBTX) error, string) {
		result.Status, result.Reason, result.NextAction = SchemaOperationFailed, reason, next
		return result, func(ctx context.Context, q DBTX) error {
			return RecordStorageRefusal(ctx, q, req.Scope, lineage.Generation, req.GroupID, StorageCodeRepairBlocked, reason, now)
		}, StorageCodeRepairBlocked
	}

	inspected, refusal, err := inspectStorageForRestore(ctx, req, layout, target, true)
	switch {
	case err != nil:
		// The catalog cannot be read: what the statements did is unknown.
		return block(StorageReasonUnknownResult, storageRepairNextActionConfirm)
	case refusal != nil:
		// A table is there now, but not the confirmed one.
		return block(StorageReasonVerificationMismatch, storageRepairNextActionOccupied)
	case inspected.pendingGroup:
		switch {
		case layout.MySQL != nil && execution.Committed == 0 && errors.Is(execErr, ErrTargetPermissionDenied):
			// The real catalog still shows no table and permission was denied;
			// this is a failed attempt, without claiming any DDL rollback.
			return failed(StorageReasonPermissionDenied, storageRepairNextActionPermission)
		case execErr != nil && execution.RolledBack && errors.Is(execErr, ErrTargetPermissionDenied):
			return failed(StorageReasonPermissionDenied, storageRepairNextActionPermission)
		case execErr != nil && execution.RolledBack:
			return failed(StorageReasonStatementFailed, storageRepairNextActionStatement)
		case execErr != nil:
			// The commit's own result was lost and the catalog shows nothing yet.
			return block(StorageReasonUnknownResult, storageRepairNextActionConfirm)
		}
		return block(StorageReasonVerificationMismatch, storageRepairNextActionConfirm)
	case inspected.pendingReceipt:
		return block(StorageReasonReceiptsMissing, storageRepairNextActionReceipts)
	}
	digest := schemaTablesDigest(inspected.tables)
	if digest != lineage.ConfirmedDigest {
		return block(StorageReasonVerificationMismatch, storageRepairNextActionConfirm)
	}
	result.Status, result.VerifiedDigest = SchemaOperationSucceeded, digest
	if execErr != nil {
		// The table is there and matches although the adapter reported an error
		// (a lost commit acknowledgement): the catalog is the evidence.
		result.Reason = SchemaReasonRecovered
	}
	return result, func(ctx context.Context, q DBTX) error {
		_, err := RecordStorageReconstruction(ctx, q, req.Scope, lineage.Generation, op.OperationID, req.GroupID, now)
		return err
	}, ""
}

// resolveStaleStorageRepair ends a repair whose owner never came back, from
// the catalog alone. It never runs a statement.
func (s *Service) resolveStaleStorageRepair(ctx context.Context, repo storageRepairRepository, req StorageRestoreRequest, stale *SchemaOperation, target StorageRestoreTarget) error {
	var detail storageRepairDetail
	if err := json.Unmarshal([]byte(stale.Detail), &detail); err != nil || detail.Kind != StorageRepairAction {
		return fmt.Errorf("storage repair %s has no usable detail", stale.OperationID)
	}
	lineage, err := ReadStorageLineage(ctx, repo.storageDB(), req.Scope)
	if err != nil {
		return err
	}
	if !lineage.Confirmed() {
		return fmt.Errorf("storage repair %s has no confirmed lineage", stale.OperationID)
	}
	layout := *lineage.Layout
	inspected, refusal, inspectErr := inspectStorageForRestore(ctx, req, layout, target, true)
	if inspectErr != nil {
		return inspectErr
	}
	cutoff := time.Now().UTC().Add(-SchemaOperationLease)
	now := time.Now()
	var result SchemaOperationResult
	if layout.MySQL != nil {
		if refusal := mysqlStorageLineageRefusal(lineage, req); refusal != nil {
			return fmt.Errorf("%w: mysql repair confirmation does not match its scope", ErrSchemaOperationMismatch)
		}
		if detail.Scope != req.Scope || detail.ConnectorID != req.Scope.ConnectorID || detail.Dialect != req.Dialect ||
			detail.Table != req.TableName || detail.OwnerColumn != req.OwnerColumn || detail.Generation != lineage.Generation || detail.ToGeneration != lineage.Generation+1 {
			return ErrSchemaOperationMismatch
		}
		if refusal == nil && !inspected.pendingGroup && schemaTablesDigest(inspected.tables) != lineage.ConfirmedDigest {
			refusal = &restoreRefusal{StorageCodeRepairBlocked, StorageReasonVerificationMismatch}
		}
	} else if layout.SQLServer != nil {
		return s.resolveStaleSQLServerStorageRepair(ctx, repo, req, stale, target, lineage, layout)
	}
	var apply func(context.Context, DBTX) error
	switch {
	case refusal == nil && !inspected.pendingGroup && !inspected.pendingReceipt:
		result = SchemaOperationResult{Status: SchemaOperationSucceeded, Reason: SchemaReasonRecovered, VerifiedDigest: schemaTablesDigest(inspected.tables)}
		apply = func(ctx context.Context, q DBTX) error {
			_, err := RecordStorageReconstruction(ctx, q, req.Scope, detail.Generation, stale.OperationID, detail.GroupID, now)
			return err
		}
	case refusal == nil && inspected.pendingGroup && (layout.MySQL == nil || inspected.tables[0].Action == SchemaTableActionCreate):
		// The table is still missing and nothing proves a statement is in
		// flight: the attempt is closed as abandoned, so a new one may follow.
		result = SchemaOperationResult{Status: SchemaOperationFailed, Reason: StorageReasonAttemptAbandoned, NextAction: storageRepairNextActionConfirm}
	default:
		result = SchemaOperationResult{Status: SchemaOperationPartial, Reason: StorageReasonVerificationMismatch, NextAction: storageRepairNextActionOccupied}
		apply = func(ctx context.Context, q DBTX) error {
			return BlockStorageRepair(ctx, q, req.Scope, detail.Generation, StorageReasonVerificationMismatch, stale.OperationID, detail.GroupID, now)
		}
	}
	if layout.MySQL != nil {
		result.Detail = stale.Detail
	}
	if _, err := repo.finishStaleStorageRepair(ctx, stale.OperationID, result, cutoff, apply); err != nil {
		return fmt.Errorf("resolve stale storage repair %s: %w", stale.OperationID, err)
	}
	return nil
}

// ResolveStaleStorageRepair ends an abandoned repair of the scope, if there is
// one, from the catalog alone: the table proves success, a still missing table
// closes the attempt as abandoned, anything else blocks. It starts no new
// attempt and runs no statement. The bool says whether one was resolved.
func (s *Service) ResolveStaleStorageRepair(ctx context.Context, req StorageRestoreRequest, target StorageRestoreTarget) (bool, error) {
	repo, ok := s.repo.(storageRepairRepository)
	if !ok {
		return false, ErrStorageRepairUnsupported
	}
	lineage, err := ReadStorageLineage(ctx, repo.storageDB(), req.Scope)
	if err != nil {
		if errors.Is(err, ErrStorageLineageNotFound) {
			return false, nil
		}
		return false, err
	}
	if !lineage.Confirmed() {
		return false, nil
	}
	scopeKey := (&SchemaPreviewToken{
		WorkspaceID: req.Scope.WorkspaceID, ConnectorID: req.Scope.ConnectorID, Database: req.Scope.Database,
		Schema: req.Scope.TableSchema, GroupLayout: lineage.Layout,
	}).ScopeKey()
	active, err := s.repo.FindActiveSchemaOperation(ctx, req.Scope.WorkspaceID, scopeKey)
	if errors.Is(err, ErrSchemaOperationNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if active.Action != StorageRepairAction || !schemaOperationLeaseExpired(active) {
		return false, nil
	}
	return true, s.resolveStaleStorageRepair(ctx, repo, req, active, target)
}

// StaleStorageRepairs lists repairs whose owner has not touched them within
// the operation lease, so a restarted gateway can resolve them.
func (s *Service) StaleStorageRepairs(ctx context.Context) ([]SchemaOperation, error) {
	repo, ok := s.repo.(storageRepairRepository)
	if !ok {
		return nil, ErrStorageRepairUnsupported
	}
	operations, err := repo.listStaleStorageRepairs(ctx, time.Now().UTC().Add(-SchemaOperationLease))
	if err != nil {
		return nil, err
	}
	for i := range operations {
		operations[i].Owner = ""
	}
	return operations, nil
}

// StorageRepairRequestOf rebuilds the restore request an abandoned repair was
// made for, from the service-owned detail of its operation. The connector
// revision is what the scope was last confirmed on.
func StorageRepairRequestOf(op SchemaOperation, connectorRevision string) (StorageRestoreRequest, bool) {
	var detail storageRepairDetail
	if err := json.Unmarshal([]byte(op.Detail), &detail); err != nil || detail.Kind != StorageRepairAction {
		return StorageRestoreRequest{}, false
	}
	return StorageRestoreRequest{
		Scope: detail.Scope, Dialect: detail.Dialect, ConnectorRevision: connectorRevision,
		TableName: detail.Table, OwnerColumn: detail.OwnerColumn, GroupID: detail.GroupID,
		SQLServer: detail.SQLServer,
	}, true
}
