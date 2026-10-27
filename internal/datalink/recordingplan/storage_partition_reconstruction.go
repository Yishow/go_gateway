package recordingplan

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// resetMySQLNativePartitionSidecar removes a stale expanded projection after a
// verified native repair or explicit re-confirmation. It only advances the
// sidecar inside the caller's transaction; a missing sidecar is a no-op.
func resetMySQLNativePartitionSidecar(
	ctx context.Context,
	q DBTX,
	scope StorageScope,
	oldLayout, newLayout GroupSchemaLayout,
	oldConnectorRevision, newConnectorRevision string,
	fromGeneration, toGeneration int,
	operationID string,
	now time.Time,
) error {
	if scope.StorageMode != mysqlPartitionStorageMode || newLayout.MySQL == nil || newLayout.MySQL.NativeMonthly == nil {
		return nil
	}
	state, err := ReadMySQLPartitionMaintenanceState(ctx, q, scope)
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
		return fmt.Errorf("%w: mysql partition sidecar does not match the confirmed generation", ErrMySQLPartitionMaintenanceStateMismatch)
	}
	encoded, err := json.Marshal(cloneMySQLPartitionEvidence(newLayout.MySQL.NativeMonthly))
	if err != nil {
		return fmt.Errorf("encode mysql partition sidecar reset: %w", err)
	}
	timestamp := now.UTC().Format(time.RFC3339Nano)
	result, err := q.ExecContext(ctx, numberedPlaceholders(`UPDATE `+mysqlPartitionMaintenanceStateTable+`
		SET base_layout_digest = $7, connector_revision = $8, current_evidence = $9, generation = $10,
			revision = $11, last_operation_id = $12, status = $13, reason = $14, updated_at = $15
		WHERE `+storageScopeWhere+` AND base_layout_digest = $16 AND connector_revision = $17
			AND generation = $18 AND revision = $19`),
		append(scopeArgs(scope), newDigest, newConnectorRevision, string(encoded), toGeneration,
			state.Revision+1, operationID, mysqlPartitionMaintenanceReady, "", timestamp,
			oldDigest, oldConnectorRevision, state.Generation, state.Revision)...)
	if err != nil {
		return fmt.Errorf("reset mysql partition sidecar: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read mysql partition sidecar reset: %w", err)
	}
	if affected != 1 {
		return ErrMySQLPartitionMaintenanceStateMismatch
	}
	return nil
}
