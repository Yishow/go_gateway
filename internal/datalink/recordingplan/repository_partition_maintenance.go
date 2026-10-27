package recordingplan

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const mysqlPartitionMaintenanceStateTable = "write_group_mysql_partition_maintenance"

type mysqlPartitionMaintenanceRepository interface {
	storageRepairRepository
	ensureMySQLPartitionMaintenanceState(context.Context, *MySQLPartitionMaintenanceState) error
	countMySQLPartitionMaintenanceAttempts(context.Context, string, string) (int, error)
	finishMySQLPartitionMaintenance(context.Context, string, string, SchemaOperationResult, func(context.Context, DBTX) error) (*SchemaOperation, error)
	finishStaleMySQLPartitionMaintenance(context.Context, string, SchemaOperationResult, time.Time, func(context.Context, DBTX) error) (*SchemaOperation, error)
}

var _ mysqlPartitionMaintenanceRepository = (*SQLRepository)(nil)

// ReadMySQLPartitionMaintenanceState reads the authorized current partition
// projection. A missing row is distinct from an unreadable or malformed row.
func ReadMySQLPartitionMaintenanceState(ctx context.Context, q DBTX, scope StorageScope) (*MySQLPartitionMaintenanceState, error) {
	if q == nil || !validateMySQLPartitionMaintenanceScope(scope) {
		return nil, ErrMySQLPartitionMaintenanceRequest
	}
	query := adaptPlaceholders(`SELECT base_layout_digest, connector_revision, generation, current_evidence,
		revision, last_operation_id, status, reason, updated_at FROM ` + mysqlPartitionMaintenanceStateTable + ` WHERE ` + storageScopeWhere)
	state := &MySQLPartitionMaintenanceState{Scope: scope}
	var encoded, updatedAt string
	err := q.QueryRowContext(ctx, query, scopeArgs(scope)...).Scan(
		&state.BaseLayoutDigest, &state.ConnectorRevision, &state.Generation, &encoded, &state.Revision,
		&state.LastOperationID, &state.Status, &state.Reason, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read mysql partition maintenance state: %w", err)
	}
	var current MySQLPartitionEvidence
	if err := json.Unmarshal([]byte(encoded), &current); err != nil {
		return nil, fmt.Errorf("decode mysql partition maintenance evidence: %w", err)
	}
	if err := validateCurrentMySQLPartitionEvidence(&current); err != nil {
		return nil, fmt.Errorf("validate mysql partition maintenance evidence: %w", err)
	}
	state.Current = &current
	if state.UpdatedAt, err = time.Parse(time.RFC3339Nano, updatedAt); err != nil {
		return nil, fmt.Errorf("decode mysql partition maintenance time: %w", err)
	}
	if err := validateMySQLPartitionMaintenanceState(state); err != nil {
		return nil, err
	}
	return state, nil
}

func (r *SQLRepository) ensureMySQLPartitionMaintenanceState(ctx context.Context, state *MySQLPartitionMaintenanceState) error {
	if err := validateMySQLPartitionMaintenanceState(state); err != nil {
		return err
	}
	encoded, err := json.Marshal(state.Current)
	if err != nil {
		return fmt.Errorf("encode mysql partition maintenance evidence: %w", err)
	}
	query := adaptPlaceholders(`INSERT INTO ` + mysqlPartitionMaintenanceStateTable + ` (
		workspace_id, connector_id, database_name, table_schema, storage_mode, layout_version,
		base_layout_digest, connector_revision, generation, current_evidence, revision,
		last_operation_id, status, reason, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	ON CONFLICT (workspace_id, connector_id, database_name, table_schema, storage_mode, layout_version) DO NOTHING`)
	_, err = r.db.ExecContext(ctx, query, append(scopeArgs(state.Scope), state.BaseLayoutDigest, state.ConnectorRevision,
		state.Generation, string(encoded), state.Revision, state.LastOperationID, state.Status, state.Reason,
		state.UpdatedAt.UTC().Format(time.RFC3339Nano))...)
	if err != nil {
		return fmt.Errorf("bootstrap mysql partition maintenance state: %w", err)
	}
	return nil
}

func updateMySQLPartitionMaintenanceState(ctx context.Context, q DBTX, expectedRevision int64, next *MySQLPartitionMaintenanceState) error {
	if q == nil || expectedRevision < 1 || validateMySQLPartitionMaintenanceState(next) != nil {
		return ErrMySQLPartitionMaintenanceStateMismatch
	}
	encoded, err := json.Marshal(next.Current)
	if err != nil {
		return fmt.Errorf("encode mysql partition maintenance evidence: %w", err)
	}
	result, err := q.ExecContext(ctx, numberedPlaceholders(`UPDATE `+mysqlPartitionMaintenanceStateTable+` SET
		base_layout_digest = $7, connector_revision = $8, generation = $9, current_evidence = $10,
		revision = $11, last_operation_id = $12, status = $13, reason = $14, updated_at = $15
		WHERE `+storageScopeWhere+` AND base_layout_digest = $16 AND connector_revision = $17
		AND generation = $18 AND revision = $19`), append(scopeArgs(next.Scope), next.BaseLayoutDigest, next.ConnectorRevision,
		next.Generation, string(encoded), next.Revision, next.LastOperationID, next.Status, next.Reason,
		next.UpdatedAt.UTC().Format(time.RFC3339Nano), next.BaseLayoutDigest, next.ConnectorRevision,
		next.Generation, expectedRevision)...)
	if err != nil {
		return fmt.Errorf("update mysql partition maintenance state: %w", err)
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read mysql partition maintenance state update: %w", err)
	}
	if changed != 1 {
		return ErrMySQLPartitionMaintenanceStateMismatch
	}
	return nil
}

func (r *SQLRepository) countMySQLPartitionMaintenanceAttempts(ctx context.Context, workspaceID, tokenPrefix string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, adaptPlaceholders(`SELECT COUNT(*) FROM managed_schema_operations
		WHERE workspace_id = $1 AND action = $2 AND substr(token, 1, length($3)) = $3`),
		workspaceID, MySQLPartitionMaintenanceAction, tokenPrefix, tokenPrefix).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count mysql partition maintenance attempts: %w", err)
	}
	return count, nil
}

func (r *SQLRepository) finishMySQLPartitionMaintenance(ctx context.Context, operationID, owner string, result SchemaOperationResult, apply func(context.Context, DBTX) error) (*SchemaOperation, error) {
	return r.finishStorageRepair(ctx, operationID, owner, result, apply)
}

func (r *SQLRepository) finishStaleMySQLPartitionMaintenance(ctx context.Context, operationID string, result SchemaOperationResult, staleBefore time.Time, apply func(context.Context, DBTX) error) (*SchemaOperation, error) {
	return r.finishStaleStorageRepair(ctx, operationID, result, staleBefore, apply)
}
