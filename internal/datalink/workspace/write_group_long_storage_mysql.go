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

// ensureMySQLStorageScopeInTx freezes an unconfirmed versioned seed for a new
// MySQL long group, without running any external probe inside the local
// transaction. An existing scope row with a confirmed layout proof is preserved
// and verified against the canonical MySQL long layout contract.
func (r *SQLWriteGroupRepository) ensureMySQLStorageScopeInTx(ctx context.Context, tx *sql.Tx, group *WriteGroup) error {
	scope := longStorageScopeOf(group)
	seed := recordingplan.UnconfirmedMySQLStorageDigest(StorageScopeOf(group))

	if _, err := tx.ExecContext(ctx, r.query(`
		INSERT INTO write_group_storage_scopes (
			workspace_id, connector_id, database_name, table_schema, storage_mode, layout_version,
			table_name, owner_column, layout_digest, generation
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 1)
		ON CONFLICT (workspace_id, connector_id, database_name, table_schema, storage_mode, layout_version) DO NOTHING
	`), scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, string(scope.Mode), string(scope.Layout),
		scope.tableName(), scope.ownerColumn(), seed); err != nil {
		return fmt.Errorf("save mysql storage scope: %w", err)
	}

	var tableName, ownerColumn, storedDigest, confirmedLayout string
	if err := tx.QueryRowContext(ctx, r.query(`
		SELECT table_name, owner_column, layout_digest, confirmed_layout FROM write_group_storage_scopes
		WHERE workspace_id = $1 AND connector_id = $2 AND database_name = $3 AND table_schema = $4
			AND storage_mode = $5 AND layout_version = $6
	`), scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, string(scope.Mode), string(scope.Layout),
	).Scan(&tableName, &ownerColumn, &storedDigest, &confirmedLayout); err != nil {
		return fmt.Errorf("confirm mysql storage scope: %w", err)
	}

	if tableName != scope.tableName() || ownerColumn != scope.ownerColumn() {
		return fmt.Errorf("%w: the saved storage scope identity does not match this layout", ErrWriteGroupValidation)
	}

	confirmedLayout = strings.TrimSpace(confirmedLayout)
	if confirmedLayout == "" {
		if storedDigest != seed {
			return fmt.Errorf("%w: the unconfirmed mysql storage scope does not match this seed", ErrWriteGroupValidation)
		}
		return nil
	}

	var confirmed recordingplan.GroupSchemaLayout
	if err := json.Unmarshal([]byte(confirmedLayout), &confirmed); err != nil {
		return fmt.Errorf("%w: invalid confirmed layout in storage scope: %w", ErrWriteGroupValidation, err)
	}
	if confirmed.MySQL == nil || (group.StorageMode == WriteGroupStorageModeNativeMonthly) != (confirmed.MySQL.NativeMonthly != nil) {
		return fmt.Errorf("%w: confirmed layout missing mysql options", ErrWriteGroupValidation)
	}
	normalized, err := recordingplan.ValidateGroupSchemaLayout(confirmed)
	if err != nil {
		return fmt.Errorf("%w: invalid confirmed layout in storage scope: %w", ErrWriteGroupValidation, err)
	}
	canonical, err := recordingplan.ManagedMySQLLongGroupLayout(scope.tableName(), scope.ownerColumn(), *normalized.MySQL)
	if err != nil {
		return fmt.Errorf("%w: build canonical mysql long layout: %w", ErrWriteGroupValidation, err)
	}
	canonicalBytes, err := json.Marshal(canonical)
	if err != nil {
		return fmt.Errorf("encode canonical storage layout: %w", err)
	}
	normalizedBytes, err := json.Marshal(normalized)
	if err != nil {
		return fmt.Errorf("encode normalized storage layout: %w", err)
	}
	if !bytes.Equal(canonicalBytes, normalizedBytes) {
		return fmt.Errorf("%w: confirmed layout does not match canonical mysql long layout", ErrWriteGroupValidation)
	}
	sum := sha256.Sum256(canonicalBytes)
	expectedDigest := hex.EncodeToString(sum[:])
	if storedDigest != expectedDigest || confirmed.TableName != scope.tableName() || confirmed.OwnerColumn != scope.ownerColumn() {
		return fmt.Errorf("%w: the confirmed mysql storage scope does not match this layout", ErrWriteGroupValidation)
	}
	return nil
}
