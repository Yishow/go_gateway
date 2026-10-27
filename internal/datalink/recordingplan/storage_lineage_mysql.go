package recordingplan

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// UnconfirmedMySQLStorageDigest identifies a saved draft scope. It is never a
// verified layout digest: confirmation atomically replaces it with the hash of
// the backend-probed, catalog-verified family-bound layout.
func UnconfirmedMySQLStorageDigest(scope StorageScope) string {
	return mysqlStorageIdentityHash(scope, "mysql-storage-unconfirmed-v1")
}

func mysqlStorageIdentityHash(scope StorageScope, purpose string) string {
	input := strings.Join([]string{purpose, scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, scope.StorageMode, scope.LayoutVersion}, "\x00")
	sum := sha256.Sum256([]byte(input))
	return hex.EncodeToString(sum[:])
}

func recordMySQLStorageConfirmation(ctx context.Context, q DBTX, scope StorageScope, c StorageConfirmation, now time.Time) (err error) {
	if db, ok := q.(*sql.DB); ok {
		var tx *sql.Tx
		tx, err = db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin mysql storage confirmation: %w", err)
		}
		defer func() {
			if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
				err = errors.Join(err, fmt.Errorf("rollback mysql storage confirmation: %w", rollbackErr))
			}
		}()
		if err := recordMySQLStorageConfirmation(ctx, tx, scope, c, now); err != nil {
			return err
		}
		return tx.Commit()
	}
	if _, ok := q.(*sql.Tx); !ok {
		return fmt.Errorf("mysql storage confirmation requires a local transaction")
	}
	if scope.WorkspaceID == "" || scope.ConnectorID == "" || scope.Database == "" || scope.Database != scope.TableSchema ||
		scope.StorageMode != "plain" && scope.StorageMode != "native_monthly" || scope.LayoutVersion != "managed_long_v1" ||
		c.ConnectorRevision == "" || c.OperationID == "" || c.GroupID == "" || len(c.VerifiedDigest) != 64 {
		return ErrInvalidGroupSchemaLayout
	}
	layout, err := ValidateGroupSchemaLayout(c.Layout)
	if err != nil {
		return err
	}
	if layout.MySQL == nil || (scope.StorageMode == "native_monthly") != (layout.MySQL.NativeMonthly != nil) {
		return ErrInvalidGroupSchemaLayout
	}
	table := "gw_long_" + mysqlStorageIdentityHash(scope, "table")[:16]
	owner := "_gw_owner_" + mysqlStorageIdentityHash(scope, "owner")[:24]
	fixed, err := ManagedMySQLLongGroupLayout(table, owner, *layout.MySQL)
	if err != nil || !groupSchemaLayoutsEqual(fixed, layout) {
		return ErrInvalidGroupSchemaLayout
	}
	current, err := ReadStorageLineage(ctx, q, scope)
	if err != nil {
		return err
	}
	if current.TableName != table || current.OwnerColumn != owner {
		return ErrIncompatibleExistingTable
	}
	encoded, err := json.Marshal(layout)
	if err != nil {
		return fmt.Errorf("encode mysql storage layout: %w", err)
	}
	sum := sha256.Sum256(encoded)
	digest := hex.EncodeToString(sum[:])
	if current.Confirmed() {
		previous, err := json.Marshal(current.Layout)
		if err != nil {
			return fmt.Errorf("encode previous mysql storage layout: %w", err)
		}
		previousHash := sha256.Sum256(previous)
		if current.Layout.MySQL == nil || current.LayoutDigest != hex.EncodeToString(previousHash[:]) {
			return ErrIncompatibleExistingTable
		}
		if current.ConfirmedOperationID == c.OperationID {
			if digest != current.LayoutDigest || c.ConnectorRevision != current.ConnectorRevision || c.VerifiedDigest != current.ConfirmedDigest {
				return ErrSchemaOperationMismatch
			}
			return nil
		}
	} else if current.LayoutDigest != UnconfirmedMySQLStorageDigest(scope) {
		return ErrIncompatibleExistingTable
	}
	args := append(scopeArgs(scope), digest, current.Generation, current.LayoutDigest, current.ConfirmedOperationID)
	result, err := q.ExecContext(ctx, numberedPlaceholders(`UPDATE write_group_storage_scopes SET layout_digest = $7
		WHERE `+storageScopeWhere+` AND generation = $8 AND layout_digest = $9 AND confirmed_operation_id = $10`), args...)
	if err != nil {
		return fmt.Errorf("bind mysql storage layout: %w", err)
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("confirm mysql storage layout update: %w", err)
	}
	if changed != 1 {
		return ErrSchemaOperationNotOwned
	}
	c.Layout = layout
	return recordStorageConfirmation(ctx, q, scope, c, now)
}

// rebindMySQLNativePartitionSidecar retains the authorized current partition
// projection while an explicit verified confirmation changes only the immutable
// base layout digest and connector revision. It runs in the confirmation
// transaction; a missing sidecar is a no-op.
func rebindMySQLNativePartitionSidecar(
	ctx context.Context,
	q DBTX,
	scope StorageScope,
	oldLayout, newLayout GroupSchemaLayout,
	oldConnectorRevision, newConnectorRevision string,
	generation int,
	operationID string,
	now time.Time,
) error {
	if scope.StorageMode != mysqlPartitionStorageMode || oldLayout.MySQL == nil || newLayout.MySQL == nil ||
		newLayout.MySQL.NativeMonthly == nil {
		return nil
	}
	oldDigest := digestGroupSchemaLayout(oldLayout)
	newDigest := digestGroupSchemaLayout(newLayout)
	if oldDigest == "" || newDigest == "" || (oldDigest == newDigest && oldConnectorRevision == newConnectorRevision) {
		return nil
	}
	state, err := ReadMySQLPartitionMaintenanceState(ctx, q, scope)
	if err != nil {
		return err
	}
	if state == nil {
		return nil
	}
	if state.Generation != generation || state.BaseLayoutDigest != oldDigest || state.ConnectorRevision != oldConnectorRevision {
		return fmt.Errorf("%w: mysql partition sidecar does not match the old confirmed binding", ErrMySQLPartitionMaintenanceStateMismatch)
	}
	timestamp := now.UTC().Format(time.RFC3339Nano)
	result, err := q.ExecContext(ctx, numberedPlaceholders(`UPDATE `+mysqlPartitionMaintenanceStateTable+`
		SET base_layout_digest = $7, connector_revision = $8, revision = revision + 1,
			last_operation_id = $9, status = $10, reason = '', updated_at = $11
		WHERE `+storageScopeWhere+` AND base_layout_digest = $12 AND connector_revision = $13
			AND generation = $14 AND revision = $15`),
		append(scopeArgs(scope), newDigest, newConnectorRevision, operationID, mysqlPartitionMaintenanceReady,
			timestamp, oldDigest, oldConnectorRevision, state.Generation, state.Revision)...)
	if err != nil {
		return fmt.Errorf("rebind mysql partition sidecar: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read mysql partition sidecar rebind: %w", err)
	}
	if affected != 1 {
		return ErrMySQLPartitionMaintenanceStateMismatch
	}
	return nil
}
