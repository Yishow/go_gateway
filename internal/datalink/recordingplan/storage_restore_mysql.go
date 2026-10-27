package recordingplan

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
)

func mysqlStorageLineageRefusal(lineage *StorageLineage, req StorageRestoreRequest) *restoreRefusal {
	refusal := &restoreRefusal{StorageCodeOwnerMismatch, StorageReasonOwnershipMismatch}
	if lineage == nil || !lineage.Confirmed() {
		return &restoreRefusal{StorageCodeOwnerMismatch, StorageReasonConfirmationMissing}
	}
	if code, reason := lineageRefusal(lineage, req); code != "" {
		return &restoreRefusal{code, reason}
	}
	scope, layout := req.Scope, lineage.Layout
	if req.Dialect != dialectMySQL || scope.WorkspaceID == "" || scope.ConnectorID == "" ||
		scope.Database == "" || scope.Database != scope.TableSchema ||
		(scope.StorageMode != "plain" && scope.StorageMode != mysqlPartitionStorageMode) ||
		scope.LayoutVersion != "managed_long_v1" || lineage.Generation < 1 || !mysqlSchemaApplyBaselineSupported(layout) {
		return refusal
	}
	if (scope.StorageMode == mysqlPartitionStorageMode) != (layout.MySQL.NativeMonthly != nil) {
		return refusal
	}
	table := "gw_long_" + mysqlStorageIdentityHash(scope, "table")[:16]
	owner := "_gw_owner_" + mysqlStorageIdentityHash(scope, "owner")[:24]
	fixed, err := ManagedMySQLLongGroupLayout(table, owner, *layout.MySQL)
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
	return nil
}

func inspectStorageForRestore(ctx context.Context, req StorageRestoreRequest, layout GroupSchemaLayout, target StorageRestoreTarget, allowPartial bool) (groupInspection, *restoreRefusal, error) {
	if req.Dialect != dialectMySQL && layout.MySQL == nil {
		if req.Dialect == dialectSQLServer || layout.SQLServer != nil {
			return inspectSQLServerStorageForRestore(ctx, req, layout, target, allowPartial)
		}
		return inspectForRestore(ctx, req.Dialect, layout, target)
	}
	if req.Dialect != dialectMySQL || layout.MySQL == nil || target.Inspect == nil {
		return groupInspection{}, nil, ErrSchemaAdapterUnverified
	}
	found, err := inspectMySQLGroupTables(ctx, layout, req.Scope.Database, target.Inspect)
	switch {
	case errors.Is(err, ErrTargetPermissionDenied):
		return groupInspection{}, &restoreRefusal{StorageCodeRepairBlocked, StorageReasonPermissionDenied}, nil
	case errors.Is(err, ErrPreviewTokenStale):
		return groupInspection{}, &restoreRefusal{StorageCodeRepairBlocked, StorageReasonVerificationMismatch}, nil
	case errors.Is(err, ErrIncompatibleExistingTable):
		return groupInspection{}, &restoreRefusal{StorageCodeOwnerMismatch, StorageReasonTableOccupied}, nil
	case err != nil:
		return groupInspection{}, nil, err
	case found.pendingReceipt:
		return groupInspection{}, &restoreRefusal{StorageCodeRepairBlocked, StorageReasonReceiptsMissing}, nil
	case !allowPartial && found.pendingGroup && found.tables[0].Action != SchemaTableActionCreate:
		// Only a table proved missing starts an automatic repair. A missing
		// index on an existing table requires a fresh explicit confirmation.
		return groupInspection{}, &restoreRefusal{StorageCodeOwnerMismatch, StorageReasonTableOccupied}, nil
	}
	return found, nil, nil
}

func (s *Service) executeMySQLStorageRepair(ctx context.Context, req StorageRestoreRequest, layout GroupSchemaLayout, op *SchemaOperation, statements []string, target StorageRestoreTarget) (SchemaExecution, error) {
	token := &SchemaPreviewToken{Dialect: dialectMySQL, Database: req.Scope.Database, Schema: req.Scope.TableSchema, GroupLayout: &layout, Statements: statements}
	steps, err := mysqlSchemaApplySteps(token)
	if err != nil {
		return SchemaExecution{}, err
	}
	for _, step := range steps {
		if step.Kind == mysqlStepKindReceiptTable {
			return SchemaExecution{}, fmt.Errorf("%w: a storage repair cannot recreate receipts", ErrInvalidGroupSchemaLayout)
		}
	}
	var detail storageRepairDetail
	if err := json.Unmarshal([]byte(op.Detail), &detail); err != nil || detail.Kind != StorageRepairAction {
		return SchemaExecution{}, ErrSchemaOperationMismatch
	}
	token.Digest = op.Token
	verified := 0
	for i := range steps {
		step := &steps[i]
		before, inspectErr := inspectMySQLApplyStep(ctx, token, layout,
			mysqlApplyStepRef{Kind: step.Kind, Name: step.Name, SQLHash: step.SQLHash}, target.Inspect)
		if inspectErr != nil {
			return SchemaExecution{Committed: verified}, inspectErr
		}
		if before.State == mysqlStepStateExact {
			step.Status, step.Fingerprint = mysqlSchemaStepVerified, before.Fingerprint
			if err := s.saveMySQLStorageRepairSteps(ctx, op, &detail, steps); err != nil {
				return SchemaExecution{Committed: verified}, err
			}
			verified++
			continue
		}
		if before.State != mysqlStepStateMissing {
			return SchemaExecution{Committed: verified}, ErrTargetInspectionUnconfirmed
		}
		step.Status = mysqlSchemaStepStarted
		if err := s.saveMySQLStorageRepairSteps(ctx, op, &detail, steps); err != nil {
			return SchemaExecution{Committed: verified}, err
		}
		if err := s.requireMySQLSchemaOperationOwner(ctx, op); err != nil {
			return SchemaExecution{Committed: verified}, err
		}
		// There is one DDL call, outside the row transaction. Its response alone
		// cannot prove application or rollback; the catalog below decides.
		_, execErr := runSchemaStatements(ctx, target.Execute, []string{statements[i]})
		after, inspectErr := inspectMySQLApplyStep(ctx, token, layout,
			mysqlApplyStepRef{Kind: step.Kind, Name: step.Name, SQLHash: step.SQLHash}, target.Inspect)
		if inspectErr != nil {
			return SchemaExecution{Committed: verified}, inspectErr
		}
		if after.State != mysqlStepStateExact {
			step.Status, step.FailureCode = mysqlSchemaStepFailed, mysqlSchemaDDLFailureCode(execErr)
			if err := s.saveMySQLStorageRepairSteps(ctx, op, &detail, steps); err != nil {
				return SchemaExecution{Committed: verified}, err
			}
			if execErr == nil {
				execErr = ErrTargetInspectionUnconfirmed
			}
			return SchemaExecution{Committed: verified}, execErr
		}
		step.Status, step.Fingerprint = mysqlSchemaStepVerified, after.Fingerprint
		if err := s.saveMySQLStorageRepairSteps(ctx, op, &detail, steps); err != nil {
			return SchemaExecution{Committed: verified}, err
		}
		verified++
	}
	return SchemaExecution{Committed: verified}, nil
}

func (s *Service) saveMySQLStorageRepairSteps(ctx context.Context, op *SchemaOperation, detail *storageRepairDetail, steps []mysqlSchemaApplyStep) error {
	progress, ok := s.repo.(operationExecutionRepository)
	if !ok {
		return ErrSchemaOperationUnacknowledged
	}
	detail.MySQLSteps = steps
	encoded, err := json.Marshal(detail)
	if err != nil {
		return fmt.Errorf("encode mysql storage repair steps: %w", err)
	}
	if err := progress.SaveSchemaOperationProgress(ctx, op.OperationID, op.Owner, string(encoded)); err != nil {
		return fmt.Errorf("%w: %w", ErrSchemaOperationUnacknowledged, err)
	}
	op.Detail = string(encoded)
	return nil
}
