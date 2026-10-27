package workspace

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/recordingplan"
)

// ensureSQLServerStorageScopeInTx freezes an unconfirmed seed for a new SQL
// Server long group. A confirmed scope is preserved and checked against the
// canonical SQL Server layout; no destination probe runs inside local storage.
func (r *SQLWriteGroupRepository) ensureSQLServerStorageScopeInTx(ctx context.Context, tx *sql.Tx, group *WriteGroup) error {
	scope := longStorageScopeOf(group)
	seed := recordingplan.UnconfirmedSQLServerStorageDigest(StorageScopeOf(group))
	if _, err := tx.ExecContext(ctx, r.query(`
		INSERT INTO write_group_storage_scopes (
			workspace_id, connector_id, database_name, table_schema, storage_mode, layout_version,
			table_name, owner_column, layout_digest, generation
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 1)
		ON CONFLICT (workspace_id, connector_id, database_name, table_schema, storage_mode, layout_version) DO NOTHING
	`), scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, string(scope.Mode), string(scope.Layout),
		scope.tableName(), scope.ownerColumn(), seed); err != nil {
		return fmt.Errorf("save sqlserver storage scope: %w", err)
	}
	var tableName, ownerColumn, storedDigest, confirmedLayout string
	if err := tx.QueryRowContext(ctx, r.query(`
		SELECT table_name, owner_column, layout_digest, confirmed_layout FROM write_group_storage_scopes
		WHERE workspace_id = $1 AND connector_id = $2 AND database_name = $3 AND table_schema = $4
			AND storage_mode = $5 AND layout_version = $6
	`), scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, string(scope.Mode), string(scope.Layout),
	).Scan(&tableName, &ownerColumn, &storedDigest, &confirmedLayout); err != nil {
		return fmt.Errorf("confirm sqlserver storage scope: %w", err)
	}
	if tableName != scope.tableName() || ownerColumn != scope.ownerColumn() {
		return fmt.Errorf("%w: the saved sqlserver storage scope identity does not match this layout", ErrWriteGroupValidation)
	}
	if strings.TrimSpace(confirmedLayout) == "" {
		if storedDigest != seed {
			return fmt.Errorf("%w: the unconfirmed sqlserver storage scope does not match this seed", ErrWriteGroupValidation)
		}
		return nil
	}
	var confirmed recordingplan.GroupSchemaLayout
	if err := json.Unmarshal([]byte(confirmedLayout), &confirmed); err != nil {
		return fmt.Errorf("%w: invalid confirmed sqlserver layout in storage scope: %w", ErrWriteGroupValidation, err)
	}
	if confirmed.SQLServer == nil || (group.StorageMode == WriteGroupStorageModeNativeMonthly) != (confirmed.SQLServer.NativeMonthly != nil) {
		return fmt.Errorf("%w: confirmed layout missing sqlserver options", ErrWriteGroupValidation)
	}
	normalized, err := recordingplan.ValidateGroupSchemaLayout(confirmed)
	if err != nil {
		return fmt.Errorf("%w: invalid confirmed sqlserver layout in storage scope: %w", ErrWriteGroupValidation, err)
	}
	canonical, err := recordingplan.ManagedSQLServerLongGroupLayout(scope.tableName(), scope.ownerColumn(), *normalized.SQLServer)
	if err != nil {
		return fmt.Errorf("%w: build canonical sqlserver long layout: %w", ErrWriteGroupValidation, err)
	}
	canonicalBytes, err := json.Marshal(canonical)
	if err != nil {
		return fmt.Errorf("encode canonical sqlserver storage layout: %w", err)
	}
	normalizedBytes, err := json.Marshal(normalized)
	if err != nil {
		return fmt.Errorf("encode normalized sqlserver storage layout: %w", err)
	}
	if !bytes.Equal(canonicalBytes, normalizedBytes) {
		return fmt.Errorf("%w: confirmed layout does not match canonical sqlserver long layout", ErrWriteGroupValidation)
	}
	sum := sha256.Sum256(canonicalBytes)
	if storedDigest != hex.EncodeToString(sum[:]) || confirmed.TableName != scope.tableName() || confirmed.OwnerColumn != scope.ownerColumn() {
		return fmt.Errorf("%w: the confirmed sqlserver storage scope does not match this layout", ErrWriteGroupValidation)
	}
	return nil
}
