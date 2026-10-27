package recordingplan

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// UnconfirmedSQLServerStorageDigest is the seed digest for a newly saved SQL
// Server storage scope. The same scope hash is shared with the managed table
// identity; the purpose string is the only SQL Server-specific part.
func UnconfirmedSQLServerStorageDigest(scope StorageScope) string {
	return mysqlStorageIdentityHash(scope, "sqlserver-storage-unconfirmed-v1")
}

// recordSQLServerStorageConfirmation freezes the layout and its canonical
// digest after a verified SQL Server schema confirmation. Unlike the generic
// SQLite/PostgreSQL path, SQL Server seeds its scope with an unconfirmed
// digest, so the confirmed layout digest must advance with the layout.
func recordSQLServerStorageConfirmation(ctx context.Context, q DBTX, scope StorageScope, c StorageConfirmation, now time.Time) (err error) {
	if db, ok := q.(*sql.DB); ok {
		var tx *sql.Tx
		tx, err = db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin sqlserver storage confirmation: %w", err)
		}
		defer func() {
			if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
				err = errors.Join(err, fmt.Errorf("rollback sqlserver storage confirmation: %w", rollbackErr))
			}
		}()
		if err := recordSQLServerStorageConfirmation(ctx, tx, scope, c, now); err != nil {
			return err
		}
		return tx.Commit()
	}
	layout, err := ValidateGroupSchemaLayout(c.Layout)
	if err != nil {
		return err
	}
	encoded, err := json.Marshal(layout)
	if err != nil {
		return fmt.Errorf("encode confirmed sqlserver layout: %w", err)
	}
	sum := sha256.Sum256(encoded)
	layoutDigest := hex.EncodeToString(sum[:])
	current, err := ReadStorageLineage(ctx, q, scope)
	if err != nil {
		return err
	}
	if current.ConfirmedOperationID == c.OperationID {
		return nil
	}
	timestamp := now.UTC().Format(storageLineageTimeLayout)
	if current.Confirmed() && c.CreatedTable {
		if _, err := q.ExecContext(ctx, numberedPlaceholders(`UPDATE write_group_storage_scopes SET generation = generation + 1,
			history_may_be_lost = 1, confirmed_layout = $7, layout_digest = $8, confirmed_connector_revision = $9,
			confirmed_operation_id = $10, confirmed_digest = $11, repair_state = '', repair_reason = ''
			WHERE `+storageScopeWhere+` AND generation = $12`),
			append(scopeArgs(scope), string(encoded), layoutDigest, c.ConnectorRevision, c.OperationID, c.VerifiedDigest, current.Generation)...); err != nil {
			return fmt.Errorf("record recreated sqlserver storage: %w", err)
		}
		if err := recordStorageEvent(ctx, q, scope, current.Generation+1, c.GroupID, StorageEventRecreatedByConfirm, "", "table-created-again", c.OperationID, timestamp); err != nil {
			return err
		}
		if layout.SQLServer != nil && layout.SQLServer.NativeMonthly != nil {
			oldLayout := GroupSchemaLayout{}
			if current.Layout != nil {
				oldLayout = *current.Layout
			}
			return resetSQLServerNativePartitionSidecar(ctx, q, scope, oldLayout, layout, current.ConnectorRevision, c.ConnectorRevision, current.Generation, current.Generation+1, c.OperationID, now)
		}
		return nil
	}
	if _, err := q.ExecContext(ctx, numberedPlaceholders(`UPDATE write_group_storage_scopes SET confirmed_layout = $7,
		layout_digest = $8, confirmed_connector_revision = $9, confirmed_operation_id = $10, confirmed_digest = $11,
		repair_state = '', repair_reason = '' WHERE `+storageScopeWhere),
		append(scopeArgs(scope), string(encoded), layoutDigest, c.ConnectorRevision, c.OperationID, c.VerifiedDigest)...); err != nil {
		return fmt.Errorf("record sqlserver storage confirmation: %w", err)
	}
	if current.Confirmed() && layout.SQLServer != nil && layout.SQLServer.NativeMonthly != nil {
		oldLayout := GroupSchemaLayout{}
		if current.Layout != nil {
			oldLayout = *current.Layout
		}
		if err := rebindSQLServerNativePartitionSidecar(ctx, q, scope, oldLayout, layout,
			current.ConnectorRevision, c.ConnectorRevision, current.Generation, c.OperationID, now); err != nil {
			return err
		}
	}
	return recordStorageEvent(ctx, q, scope, current.Generation, c.GroupID, StorageEventConfirmed, "", "", c.OperationID, timestamp)
}

func resetSQLServerNativePartitionSidecar(
	ctx context.Context,
	q DBTX,
	scope StorageScope,
	oldLayout, newLayout GroupSchemaLayout,
	oldConnectorRevision, newConnectorRevision string,
	fromGeneration, toGeneration int,
	operationID string,
	now time.Time,
) error {
	if scope.StorageMode != sqlServerPartitionStorageMode || newLayout.SQLServer == nil || newLayout.SQLServer.NativeMonthly == nil {
		return nil
	}
	state, err := ReadSQLServerPartitionMaintenanceState(ctx, q, scope)
	if err != nil {
		return err
	}
	if state == nil {
		return nil
	}
	oldDigest := digestGroupSchemaLayout(oldLayout)
	newDigest := digestGroupSchemaLayout(newLayout)
	if oldDigest == "" || newDigest == "" || state.Generation != fromGeneration ||
		state.BaseLayoutDigest != oldDigest || state.ConnectorRevision != oldConnectorRevision ||
		toGeneration != fromGeneration+1 {
		return fmt.Errorf("%w: sqlserver partition sidecar does not match the confirmed generation", ErrSQLServerPartitionMaintenanceStateMismatch)
	}
	encoded, err := json.Marshal(CloneSQLServerPartitionEvidence(newLayout.SQLServer.NativeMonthly))
	if err != nil {
		return fmt.Errorf("encode sqlserver partition sidecar reset: %w", err)
	}
	timestamp := now.UTC().Format(time.RFC3339Nano)
	result, err := q.ExecContext(ctx, numberedPlaceholders(`UPDATE `+sqlServerPartitionMaintenanceStateTable+`
		SET base_layout_digest = $7, connector_revision = $8, current_evidence = $9, generation = $10,
			revision = $11, last_operation_id = $12, status = $13, reason = $14, updated_at = $15
		WHERE `+storageScopeWhere+` AND base_layout_digest = $16 AND connector_revision = $17
			AND generation = $18 AND revision = $19`),
		append(scopeArgs(scope), newDigest, newConnectorRevision, string(encoded), toGeneration,
			state.Revision+1, operationID, sqlServerPartitionMaintenanceReady, "", timestamp,
			oldDigest, oldConnectorRevision, state.Generation, state.Revision)...)
	if err != nil {
		return fmt.Errorf("reset sqlserver partition sidecar: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read sqlserver partition sidecar reset: %w", err)
	}
	if affected != 1 {
		return ErrSQLServerPartitionMaintenanceStateMismatch
	}
	return nil
}

func rebindSQLServerNativePartitionSidecar(
	ctx context.Context,
	q DBTX,
	scope StorageScope,
	oldLayout, newLayout GroupSchemaLayout,
	oldConnectorRevision, newConnectorRevision string,
	generation int,
	operationID string,
	now time.Time,
) error {
	if scope.StorageMode != sqlServerPartitionStorageMode || oldLayout.SQLServer == nil || newLayout.SQLServer == nil ||
		newLayout.SQLServer.NativeMonthly == nil {
		return nil
	}
	oldDigest := digestGroupSchemaLayout(oldLayout)
	newDigest := digestGroupSchemaLayout(newLayout)
	if oldDigest == "" || newDigest == "" || (oldDigest == newDigest && oldConnectorRevision == newConnectorRevision) {
		return nil
	}
	state, err := ReadSQLServerPartitionMaintenanceState(ctx, q, scope)
	if err != nil {
		return err
	}
	if state == nil {
		return nil
	}
	if state.Generation != generation || state.BaseLayoutDigest != oldDigest || state.ConnectorRevision != oldConnectorRevision {
		return fmt.Errorf("%w: sqlserver partition sidecar does not match the old confirmed binding", ErrSQLServerPartitionMaintenanceStateMismatch)
	}
	base := newLayout.SQLServer.NativeMonthly
	if !sameSQLServerPartitionEvidence(state.Current, base) && !sameSQLServerPartitionPrefix(base, state.Current) {
		return fmt.Errorf("%w: sqlserver partition sidecar current evidence is incompatible with new layout", ErrSQLServerPartitionMaintenanceStateMismatch)
	}
	timestamp := now.UTC().Format(time.RFC3339Nano)
	result, err := q.ExecContext(ctx, numberedPlaceholders(`UPDATE `+sqlServerPartitionMaintenanceStateTable+`
		SET base_layout_digest = $7, connector_revision = $8, revision = revision + 1,
			last_operation_id = $9, status = $10, reason = '', updated_at = $11
		WHERE `+storageScopeWhere+` AND base_layout_digest = $12 AND connector_revision = $13
			AND generation = $14 AND revision = $15`),
		append(scopeArgs(scope), newDigest, newConnectorRevision, operationID, sqlServerPartitionMaintenanceReady,
			timestamp, oldDigest, oldConnectorRevision, state.Generation, state.Revision)...)
	if err != nil {
		return fmt.Errorf("rebind sqlserver partition sidecar: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read sqlserver partition sidecar rebind: %w", err)
	}
	if affected != 1 {
		return ErrSQLServerPartitionMaintenanceStateMismatch
	}
	return nil
}
