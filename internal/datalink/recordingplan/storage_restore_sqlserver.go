package recordingplan

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const sqlserverStorageModePlain = "plain"

// SQLServerStorageRestoreEvidence contains the original frozen confirmed physical
// inspections for a SQL Server managed storage scope.
type SQLServerStorageRestoreEvidence struct {
	Group   TargetTableInspection `json:"group"`
	Receipt TargetTableInspection `json:"receipt"`
}

func sqlServerStorageLineageRefusal(lineage *StorageLineage, req StorageRestoreRequest) *restoreRefusal {
	return sqlServerStorageLineageRefusalMode(lineage, req, true)
}

func sqlServerStorageLineageRefusalForVerification(lineage *StorageLineage, req StorageRestoreRequest) *restoreRefusal {
	return sqlServerStorageLineageRefusalMode(lineage, req, true)
}

func sqlServerStorageLineageRefusalMode(lineage *StorageLineage, req StorageRestoreRequest, allowNative bool) *restoreRefusal {
	refusal := &restoreRefusal{StorageCodeOwnerMismatch, StorageReasonOwnershipMismatch}
	if lineage == nil || !lineage.Confirmed() {
		return &restoreRefusal{StorageCodeOwnerMismatch, StorageReasonConfirmationMissing}
	}
	if code, reason := lineageRefusal(lineage, req); code != "" {
		return &restoreRefusal{code, reason}
	}
	scope, layout := req.Scope, lineage.Layout
	if req.Dialect != dialectSQLServer || layout == nil || layout.SQLServer == nil ||
		scope.WorkspaceID == "" || scope.ConnectorID == "" ||
		scope.Database == "" || scope.TableSchema == "" ||
		(scope.StorageMode != sqlserverStorageModePlain && (!allowNative || scope.StorageMode != sqlServerPartitionStorageMode || layout.SQLServer.NativeMonthly == nil)) ||
		(scope.StorageMode == sqlserverStorageModePlain && layout.SQLServer.NativeMonthly != nil) || scope.LayoutVersion != mysqlPartitionLayoutVersion ||
		lineage.Generation < 1 {
		return refusal
	}
	opts := layout.SQLServer
	version16 := strings.HasPrefix(opts.ServerVersion, "16.") || opts.ServerVersion == "16"
	version17 := strings.HasPrefix(opts.ServerVersion, "17.") || opts.ServerVersion == "17"
	if opts.Database != scope.Database || opts.Schema != scope.TableSchema ||
		(!version16 && !version17) ||
		(version16 && opts.CompatibilityLevel != 160) ||
		(version17 && opts.CompatibilityLevel != 170) ||
		opts.CapabilityRevision == "" ||
		opts.ServerEdition == "" {
		return refusal
	}
	table := "gw_long_" + mysqlStorageIdentityHash(scope, "table")[:16]
	owner := "_gw_owner_" + mysqlStorageIdentityHash(scope, "owner")[:24]
	fixed, err := ManagedSQLServerLongGroupLayout(table, owner, *opts)
	if err != nil || !groupSchemaLayoutsEqual(fixed, *layout) || req.TableName != table || req.OwnerColumn != owner {
		return refusal
	}
	encoded, err := json.Marshal(layout)
	if err != nil {
		return refusal
	}
	sum := sha256.Sum256(encoded)
	if lineage.LayoutDigest != hex.EncodeToString(sum[:]) || len(lineage.ConfirmedDigest) != sha256.Size*2 {
		return refusal
	}
	if _, err := hex.DecodeString(lineage.ConfirmedDigest); err != nil {
		return refusal
	}
	if req.SQLServer == nil {
		return refusal
	}
	origDigest, err := VerifyManagedSQLServerCatalog(*layout, scope.Database, req.SQLServer.Group, req.SQLServer.Receipt)
	if err != nil || origDigest != lineage.ConfirmedDigest {
		return refusal
	}
	return nil
}

func inspectSQLServerStorageForRestore(ctx context.Context, req StorageRestoreRequest, layout GroupSchemaLayout, target StorageRestoreTarget, allowPartial bool) (groupInspection, *restoreRefusal, error) {
	if target.Inspect == nil {
		return groupInspection{}, nil, ErrSchemaAdapterUnverified
	}
	if layout.SQLServer == nil {
		return groupInspection{}, nil, ErrInvalidGroupSchemaLayout
	}
	database := req.Scope.Database
	if layout.SQLServer != nil && layout.SQLServer.Database != "" {
		database = layout.SQLServer.Database
	}
	var found groupInspection
	var err error
	if layout.SQLServer.NativeMonthly != nil {
		found, err = inspectSQLServerNativeRestoreGroupTables(ctx, layout, database, target.Inspect)
	} else {
		found, err = inspectSQLServerGroupTables(ctx, layout, database, target.Inspect)
	}
	switch {
	case errors.Is(err, ErrTargetPermissionDenied):
		return groupInspection{}, &restoreRefusal{StorageCodeRepairBlocked, StorageReasonPermissionDenied}, nil
	case errors.Is(err, ErrPreviewTokenStale):
		return groupInspection{}, &restoreRefusal{StorageCodeRepairBlocked, StorageReasonVerificationMismatch}, nil
	case errors.Is(err, ErrIncompatibleExistingTable), errors.Is(err, ErrSQLServerPartitionUnavailable):
		return groupInspection{}, &restoreRefusal{StorageCodeOwnerMismatch, StorageReasonTableOccupied}, nil
	case err != nil:
		return groupInspection{}, nil, err
	case found.pendingReceipt:
		return groupInspection{}, &restoreRefusal{StorageCodeRepairBlocked, StorageReasonReceiptsMissing}, nil
	case !allowPartial && found.pendingGroup && found.tables[0].Action != SchemaTableActionCreate:
		return groupInspection{}, &restoreRefusal{StorageCodeOwnerMismatch, StorageReasonTableOccupied}, nil
	}
	return found, nil, nil
}

func validateSQLServerReconstructionOperation(op *SchemaOperation, scope StorageScope, lineage *StorageLineage) error {
	if op == nil || op.Action != StorageRepairAction || op.Status != SchemaOperationSucceeded || op.VerifiedDigest == "" {
		return fmt.Errorf("%w: invalid reconstruction operation", ErrSchemaOperationMismatch)
	}
	var detail storageRepairDetail
	if err := json.Unmarshal([]byte(op.Detail), &detail); err != nil {
		return fmt.Errorf("%w: invalid reconstruction detail: %w", ErrSchemaOperationMismatch, err)
	}
	if detail.Scope != scope || detail.ConnectorID != scope.ConnectorID || detail.Dialect != dialectSQLServer ||
		detail.Generation != lineage.Generation-1 || detail.ToGeneration != lineage.Generation ||
		detail.SQLServer == nil || detail.SQLServer.Group.SQLServer == nil || detail.SQLServer.Receipt.SQLServer == nil {
		return fmt.Errorf("%w: reconstruction detail scope or generation mismatch", ErrSchemaOperationMismatch)
	}
	if lineage.Layout != nil {
		if detail.Table != lineage.Layout.TableName || detail.OwnerColumn != lineage.Layout.OwnerColumn {
			return fmt.Errorf("%w: reconstruction detail table mismatch", ErrSchemaOperationMismatch)
		}
		origDigest, err := VerifyManagedSQLServerCatalog(*lineage.Layout, scope.Database, detail.SQLServer.Group, detail.SQLServer.Receipt)
		if err != nil || origDigest != lineage.ConfirmedDigest {
			return fmt.Errorf("%w: reconstruction detail original catalog does not match confirmed digest", ErrSchemaOperationMismatch)
		}
	}
	return nil
}

func (s *Service) authorizedSQLServerDigest(ctx context.Context, scope StorageScope, lineage *StorageLineage) (string, error) {
	if lineage == nil || !lineage.Confirmed() {
		return "", ErrStorageLineageNotFound
	}
	if lineage.Generation <= 1 || lineage.LastReconstructionOperationID == "" {
		return lineage.ConfirmedDigest, nil
	}
	lastOp, err := s.repo.GetSchemaOperation(ctx, scope.WorkspaceID, lineage.LastReconstructionOperationID)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrSchemaOperationMismatch, err)
	}
	if err := validateSQLServerReconstructionOperation(lastOp, scope, lineage); err != nil {
		return "", err
	}
	return lastOp.VerifiedDigest, nil
}

func (s *Service) executeSQLServerStorageRepair(ctx context.Context, req StorageRestoreRequest, layout GroupSchemaLayout, op *SchemaOperation, statements []string, target StorageRestoreTarget) (SchemaExecution, error) {
	if target.ExecuteSQLServerRepair == nil {
		return SchemaExecution{}, ErrSchemaAdapterUnverified
	}
	for _, stmt := range statements {
		if strings.Contains(stmt, managedEffectReceiptTable) {
			return SchemaExecution{}, fmt.Errorf("%w: a storage repair cannot recreate receipts", ErrInvalidGroupSchemaLayout)
		}
	}
	progress, ok := s.repo.(operationExecutionRepository)
	if !ok {
		return SchemaExecution{}, ErrSchemaOperationUnacknowledged
	}
	var detail storageRepairDetail
	if err := json.Unmarshal([]byte(op.Detail), &detail); err != nil || detail.Kind != StorageRepairAction {
		return SchemaExecution{}, ErrSchemaOperationMismatch
	}
	repairStatements, err := s.sqlServerNativeRepairStatements(ctx, req, layout, op, statements)
	if err != nil {
		return SchemaExecution{RolledBack: true}, err
	}
	detail.SQLServerStarted = true
	encoded, err := json.Marshal(detail)
	if err != nil {
		return SchemaExecution{}, fmt.Errorf("encode sqlserver storage repair detail: %w", err)
	}
	if err := progress.SaveSchemaOperationProgress(ctx, op.OperationID, op.Owner, string(encoded)); err != nil {
		return SchemaExecution{}, fmt.Errorf("%w: %w", ErrSchemaOperationUnacknowledged, err)
	}
	op.Detail = string(encoded)

	return target.ExecuteSQLServerRepair(ctx, layout.TableName, op.OperationID, repairStatements)
}

func (s *Service) judgeSQLServerStorageRepair(ctx context.Context, req StorageRestoreRequest, lineage *StorageLineage, layout GroupSchemaLayout, op *SchemaOperation, execution SchemaExecution, execErr error, target StorageRestoreTarget, now time.Time) (result SchemaOperationResult, apply func(context.Context, DBTX) error, code string) {
	result = SchemaOperationResult{ExecutedStatements: execution.Committed, Detail: op.Detail}
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

	if target.Inspect == nil {
		return block(StorageReasonUnknownResult, storageRepairNextActionConfirm)
	}
	groupState, groupErr := target.Inspect(ctx, layout.TableName)
	receiptState, receiptErr := target.Inspect(ctx, managedEffectReceiptTable)
	if groupErr != nil || receiptErr != nil {
		return block(StorageReasonUnknownResult, storageRepairNextActionConfirm)
	}

	repo, ok := s.repo.(storageRepairRepository)
	if !ok {
		return block(StorageReasonUnknownResult, storageRepairNextActionConfirm)
	}
	normalizedGroup, err := normalizeSQLServerNativeRepairGroup(ctx, repo.storageDB(), req, lineage, layout, groupState)
	if err != nil {
		return block(StorageReasonVerificationMismatch, storageRepairNextActionOccupied)
	}
	database := req.Scope.Database
	if layout.SQLServer != nil && layout.SQLServer.Database != "" {
		database = layout.SQLServer.Database
	}
	inspected, err := inspectSQLServerGroupSnapshot(layout, database, normalizedGroup, receiptState)
	switch {
	case errors.Is(err, ErrTargetPermissionDenied):
		return block(StorageReasonPermissionDenied, storageRepairNextActionPermission)
	case errors.Is(err, ErrPreviewTokenStale):
		return block(StorageReasonVerificationMismatch, storageRepairNextActionOccupied)
	case errors.Is(err, ErrIncompatibleExistingTable):
		return block(StorageReasonTableOccupied, storageRepairNextActionOccupied)
	case err != nil:
		return block(StorageReasonUnknownResult, storageRepairNextActionConfirm)
	case inspected.pendingReceipt:
		return block(StorageReasonReceiptsMissing, storageRepairNextActionReceipts)
	case inspected.pendingGroup:
		switch {
		case execErr != nil && execution.RolledBack && errors.Is(execErr, ErrTargetPermissionDenied):
			return failed(StorageReasonPermissionDenied, storageRepairNextActionPermission)
		case execErr != nil && execution.RolledBack:
			return failed(StorageReasonStatementFailed, storageRepairNextActionStatement)
		case execErr != nil:
			return block(StorageReasonUnknownResult, storageRepairNextActionConfirm)
		default:
			return block(StorageReasonVerificationMismatch, storageRepairNextActionConfirm)
		}
	}

	if execution.RolledBack {
		reason, next := StorageReasonStatementFailed, storageRepairNextActionStatement
		if execErr != nil && errors.Is(execErr, ErrTargetPermissionDenied) {
			reason, next = StorageReasonPermissionDenied, storageRepairNextActionPermission
		}
		return failed(reason, next)
	}

	currReceiptDigest, err := sqlserverCatalogDigest(receiptState)
	if err != nil {
		return block(StorageReasonUnknownResult, storageRepairNextActionConfirm)
	}
	origReceiptDigest, err := sqlserverCatalogDigest(req.SQLServer.Receipt)
	if err != nil || currReceiptDigest != origReceiptDigest {
		return block(StorageReasonVerificationMismatch, storageRepairNextActionOccupied)
	}

	if groupState.SQLServer == nil || groupState.SQLServer.ObjectID <= 0 {
		return block(StorageReasonVerificationMismatch, storageRepairNextActionConfirm)
	}
	if groupState.SQLServer.RepairOperationID != op.OperationID {
		return block(StorageReasonVerificationMismatch, storageRepairNextActionOccupied)
	}

	digest := schemaTablesDigest(inspected.tables)
	result.Status, result.VerifiedDigest = SchemaOperationSucceeded, digest
	if execErr != nil {
		result.Reason = SchemaReasonRecovered
	}
	return result, func(ctx context.Context, q DBTX) error {
		changed, err := RecordStorageReconstruction(ctx, q, req.Scope, lineage.Generation, op.OperationID, req.GroupID, now)
		if err != nil {
			return err
		}
		if !changed {
			return fmt.Errorf("%w: storage reconstruction generation CAS miss for %s", ErrSchemaOperationMismatch, op.OperationID)
		}
		return nil
	}, ""
}

func (s *Service) resolveStaleSQLServerStorageRepair(ctx context.Context, repo storageRepairRepository, req StorageRestoreRequest, stale *SchemaOperation, target StorageRestoreTarget, lineage *StorageLineage, layout GroupSchemaLayout) error {
	if refusal := sqlServerStorageLineageRefusal(lineage, req); refusal != nil {
		return fmt.Errorf("%w: sqlserver repair confirmation does not match its scope", ErrSchemaOperationMismatch)
	}
	var detail storageRepairDetail
	if err := json.Unmarshal([]byte(stale.Detail), &detail); err != nil || detail.Kind != StorageRepairAction {
		return fmt.Errorf("storage repair %s has no usable detail", stale.OperationID)
	}
	if detail.Scope != req.Scope || detail.ConnectorID != req.Scope.ConnectorID || detail.Dialect != req.Dialect ||
		detail.Table != req.TableName || detail.OwnerColumn != req.OwnerColumn || detail.Generation != lineage.Generation ||
		detail.ToGeneration != lineage.Generation+1 {
		return ErrSchemaOperationMismatch
	}
	if detail.SQLServer == nil || detail.SQLServer.Group.SQLServer == nil || detail.SQLServer.Receipt.SQLServer == nil {
		return ErrSchemaOperationMismatch
	}
	detailOrigDigest, err := VerifyManagedSQLServerCatalog(layout, req.Scope.Database, detail.SQLServer.Group, detail.SQLServer.Receipt)
	if err != nil || detailOrigDigest != lineage.ConfirmedDigest {
		return ErrSchemaOperationMismatch
	}
	detailReceiptDigest, err := sqlserverCatalogDigest(detail.SQLServer.Receipt)
	if err != nil {
		return ErrSchemaOperationMismatch
	}
	reqReceiptDigest, err := sqlserverCatalogDigest(req.SQLServer.Receipt)
	if err != nil || detailReceiptDigest != reqReceiptDigest {
		return ErrSchemaOperationMismatch
	}

	if target.Inspect == nil {
		return ErrSchemaAdapterUnverified
	}
	groupState, groupErr := target.Inspect(ctx, layout.TableName)
	receiptState, receiptErr := target.Inspect(ctx, managedEffectReceiptTable)
	if groupErr != nil || receiptErr != nil {
		return fmt.Errorf("inspect sqlserver storage for stale repair: %w", errors.Join(groupErr, receiptErr))
	}
	normalizedGroup, normalizeErr := normalizeSQLServerNativeRepairGroup(ctx, repo.storageDB(), req, lineage, layout, groupState)
	if normalizeErr != nil {
		return normalizeErr
	}
	database := req.Scope.Database
	if layout.SQLServer != nil && layout.SQLServer.Database != "" {
		database = layout.SQLServer.Database
	}
	inspected, inspectErr := inspectSQLServerGroupSnapshot(layout, database, normalizedGroup, receiptState)

	cutoff := time.Now().UTC().Add(-SchemaOperationLease)
	now := time.Now()
	var result SchemaOperationResult
	var apply func(context.Context, DBTX) error

	var refusalReason string
	switch {
	case errors.Is(inspectErr, ErrTargetPermissionDenied):
		refusalReason = StorageReasonPermissionDenied
	case errors.Is(inspectErr, ErrPreviewTokenStale) || errors.Is(inspectErr, ErrIncompatibleExistingTable):
		refusalReason = StorageReasonVerificationMismatch
	case inspectErr != nil:
		return inspectErr
	case inspected.pendingReceipt:
		refusalReason = StorageReasonReceiptsMissing
	}

	switch {
	case refusalReason != "":
		result = SchemaOperationResult{Status: SchemaOperationPartial, Reason: refusalReason, NextAction: storageRepairNextActionOccupied, Detail: stale.Detail}
		apply = func(ctx context.Context, q DBTX) error {
			return BlockStorageRepair(ctx, q, req.Scope, detail.Generation, refusalReason, stale.OperationID, detail.GroupID, now)
		}
	case inspected.pendingGroup:
		result = SchemaOperationResult{Status: SchemaOperationFailed, Reason: StorageReasonAttemptAbandoned, NextAction: storageRepairNextActionConfirm, Detail: stale.Detail}
	default:
		currentDigest := schemaTablesDigest(inspected.tables)
		isCurrentAuthorized := false
		if currentDigest == lineage.ConfirmedDigest {
			isCurrentAuthorized = true
		} else if lineage.Generation > 1 {
			if authDigest, err := s.authorizedSQLServerDigest(ctx, req.Scope, lineage); err == nil && currentDigest == authDigest {
				isCurrentAuthorized = true
			}
		}

		if isCurrentAuthorized {
			result = SchemaOperationResult{
				Status:     SchemaOperationFailed,
				Reason:     StorageReasonAttemptAbandoned,
				NextAction: storageRepairNextActionConfirm,
				Detail:     stale.Detail,
			}
			if _, err := repo.finishStaleStorageRepair(ctx, stale.OperationID, result, cutoff, nil); err != nil {
				return fmt.Errorf("resolve stale storage repair %s: %w", stale.OperationID, err)
			}
			return nil
		}

		currReceiptDigest, err := sqlserverCatalogDigest(receiptState)
		if err != nil || currReceiptDigest != reqReceiptDigest || !detail.SQLServerStarted ||
			groupState.SQLServer == nil || groupState.SQLServer.RepairOperationID != stale.OperationID {
			result = SchemaOperationResult{Status: SchemaOperationPartial, Reason: StorageReasonVerificationMismatch, NextAction: storageRepairNextActionOccupied, Detail: stale.Detail}
			apply = func(ctx context.Context, q DBTX) error {
				return BlockStorageRepair(ctx, q, req.Scope, detail.Generation, StorageReasonVerificationMismatch, stale.OperationID, detail.GroupID, now)
			}
		} else {
			result = SchemaOperationResult{Status: SchemaOperationSucceeded, Reason: SchemaReasonRecovered, VerifiedDigest: currentDigest, Detail: stale.Detail}
			apply = func(ctx context.Context, q DBTX) error {
				changed, err := RecordStorageReconstruction(ctx, q, req.Scope, detail.Generation, stale.OperationID, detail.GroupID, now)
				if err != nil {
					return err
				}
				if !changed {
					return fmt.Errorf("%w: storage reconstruction generation CAS miss for stale %s", ErrSchemaOperationMismatch, stale.OperationID)
				}
				return nil
			}
		}
	}

	if _, err := repo.finishStaleStorageRepair(ctx, stale.OperationID, result, cutoff, apply); err != nil {
		return fmt.Errorf("resolve stale storage repair %s: %w", stale.OperationID, err)
	}
	return nil
}

// VerifySQLServerStorageGeneration validates the storage generation and catalog facts of a SQL Server scope.
// It reads local lineage and operation records only, performing zero DDL.
func (s *Service) VerifySQLServerStorageGeneration(ctx context.Context, req StorageRestoreRequest, group, receipt TargetTableInspection) error {
	repo, ok := s.repo.(storageRepairRepository)
	if !ok {
		return ErrStorageRepairUnsupported
	}
	lineage, err := ReadStorageLineage(ctx, repo.storageDB(), req.Scope)
	if err != nil {
		return err
	}
	if !lineage.Confirmed() {
		return ErrTargetInspectionUnconfirmed
	}
	if lineage.RepairState != "" {
		return fmt.Errorf("%w: storage repair is blocked: %s", ErrTargetInspectionUnconfirmed, lineage.RepairReason)
	}
	if lineage.ConnectorRevision != req.ConnectorRevision {
		return fmt.Errorf("%w: connector revision changed", ErrTargetInspectionUnconfirmed)
	}
	if lineage.Generation < 1 {
		return ErrTargetInspectionUnconfirmed
	}
	if req.SQLServer == nil {
		return ErrTargetInspectionUnconfirmed
	}
	if refusal := sqlServerStorageLineageRefusalForVerification(lineage, req); refusal != nil {
		return ErrTargetInspectionUnconfirmed
	}

	if lineage.Layout.SQLServer.NativeMonthly != nil {
		group, err = normalizeSQLServerNativeRepairGroup(ctx, repo.storageDB(), req, lineage, *lineage.Layout, group)
		if err != nil {
			return err
		}
	}
	currentDigest, err := VerifyManagedSQLServerCatalog(*lineage.Layout, req.Scope.Database, group, receipt)
	if err != nil {
		return err
	}

	currReceiptDigest, err := sqlserverCatalogDigest(receipt)
	if err != nil {
		return err
	}
	origReceiptDigest, err := sqlserverCatalogDigest(req.SQLServer.Receipt)
	if err != nil || currReceiptDigest != origReceiptDigest {
		return fmt.Errorf("%w: receipt table does not match confirmed proof", ErrIncompatibleExistingTable)
	}

	if currentDigest == lineage.ConfirmedDigest {
		return nil
	}

	if lineage.Generation <= 1 || lineage.LastReconstructionOperationID == "" {
		return ErrTargetInspectionUnconfirmed
	}
	lastOp, err := s.repo.GetSchemaOperation(ctx, req.Scope.WorkspaceID, lineage.LastReconstructionOperationID)
	if err != nil || validateSQLServerReconstructionOperation(lastOp, req.Scope, lineage) != nil || lastOp.VerifiedDigest != currentDigest {
		return ErrTargetInspectionUnconfirmed
	}
	return nil
}
