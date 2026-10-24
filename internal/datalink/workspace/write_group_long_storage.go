package workspace

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
)

// longStorageScope is the owner of one shared managed_long_v1 table: the
// workspace plus the confirmed destination, storage mode and layout. A group,
// a label or a point count never takes part, so adding sources only adds rows.
type longStorageScope struct {
	WorkspaceID, ConnectorID, Database, TableSchema string
	Mode                                            WriteGroupStorageMode
	Layout                                          WriteGroupLayoutVersion
}

func longStorageScopeOf(group *WriteGroup) longStorageScope {
	return longStorageScope{
		WorkspaceID: group.WorkspaceID, ConnectorID: group.Destination.ConnectorID,
		Database: group.Destination.Database, TableSchema: group.Destination.TableSchema,
		Mode: group.StorageMode, Layout: group.EffectiveLayoutVersion(),
	}
}

func (s longStorageScope) digest(purpose string) string {
	sum := sha256.Sum256([]byte(strings.Join([]string{
		purpose, s.WorkspaceID, s.ConnectorID, s.Database, s.TableSchema, string(s.Mode), string(s.Layout),
	}, "\x00")))
	return hex.EncodeToString(sum[:])
}

// tableName is derived from the scope, never from a group or display text.
func (s longStorageScope) tableName() string { return "gw_long_" + s.digest("table")[:16] }

// ownerColumn is the ownership marker of the table. Another workspace,
// destination or layout derives another name, so its table cannot validate.
func (s longStorageScope) ownerColumn() string { return "_gw_owner_" + s.digest("owner")[:24] }

func (s longStorageScope) layout(dialect string) (*recordingplan.GroupSchemaLayout, error) {
	if s.Mode == WriteGroupStorageModeTimescale {
		if dialect != writeGroupConnectorKindPostgres {
			// Only PostgreSQL has hypertables. The selection is refused, never
			// rewritten to plain storage.
			return nil, fmt.Errorf("%w: %w", ErrWriteGroupValidation, &recordingplan.TimescaleUnavailableError{Reasons: []string{dbtarget.TimescaleReasonDialect}})
		}
		layout, err := recordingplan.ManagedLongTimescaleLayout(s.tableName(), s.ownerColumn())
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrWriteGroupValidation, err)
		}
		return &layout, nil
	}
	layout, err := recordingplan.ManagedLongGroupLayout(dialect, s.tableName(), s.ownerColumn())
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrWriteGroupValidation, err)
	}
	return &layout, nil
}

// ensureStorageScopeInTx freezes the scope's table identity, owner marker,
// layout digest and generation the first time a long group saves into it. It
// runs in the caller's transaction, only on create and update, never on reads.
func (r *SQLWriteGroupRepository) ensureStorageScopeInTx(ctx context.Context, tx *sql.Tx, group *WriteGroup) error {
	if !group.isLongLayout() {
		return nil
	}
	var kind string
	if err := tx.QueryRowContext(ctx, r.query(`SELECT kind FROM database_connectors WHERE id = $1`),
		group.Destination.ConnectorID).Scan(&kind); err != nil {
		return fmt.Errorf("read storage scope connector: %w", err)
	}
	scope := longStorageScopeOf(group)
	layout, err := scope.layout(strings.ToLower(strings.TrimSpace(kind)))
	if err != nil {
		return err
	}
	encoded, err := json.Marshal(layout)
	if err != nil {
		return fmt.Errorf("encode storage layout: %w", err)
	}
	sum := sha256.Sum256(encoded)
	digest := hex.EncodeToString(sum[:])
	if _, err := tx.ExecContext(ctx, r.query(`
		INSERT INTO write_group_storage_scopes (
			workspace_id, connector_id, database_name, table_schema, storage_mode, layout_version,
			table_name, owner_column, layout_digest, generation
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 1)
		ON CONFLICT (workspace_id, connector_id, database_name, table_schema, storage_mode, layout_version) DO NOTHING
	`), scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, string(scope.Mode), string(scope.Layout),
		scope.tableName(), scope.ownerColumn(), digest); err != nil {
		return fmt.Errorf("save storage scope: %w", err)
	}
	var tableName, ownerColumn, storedDigest string
	if err := tx.QueryRowContext(ctx, r.query(`
		SELECT table_name, owner_column, layout_digest FROM write_group_storage_scopes
		WHERE workspace_id = $1 AND connector_id = $2 AND database_name = $3 AND table_schema = $4
			AND storage_mode = $5 AND layout_version = $6
	`), scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, string(scope.Mode), string(scope.Layout),
	).Scan(&tableName, &ownerColumn, &storedDigest); err != nil {
		return fmt.Errorf("confirm storage scope: %w", err)
	}
	if tableName != scope.tableName() || ownerColumn != scope.ownerColumn() || storedDigest != digest {
		return fmt.Errorf("%w: the saved storage scope does not match this layout", ErrWriteGroupValidation)
	}
	return nil
}

// StorageScopeOf returns the key of the local storage scope row of a long
// group. It is meaningful for managed long groups only.
func StorageScopeOf(group *WriteGroup) recordingplan.StorageScope {
	scope := longStorageScopeOf(group)
	return recordingplan.StorageScope{
		WorkspaceID: scope.WorkspaceID, ConnectorID: scope.ConnectorID, Database: scope.Database, TableSchema: scope.TableSchema,
		StorageMode: string(scope.Mode), LayoutVersion: string(scope.Layout),
	}
}

// LongStorageLayout returns the fixed layout the shared managed_long_v1 table
// of the group's storage scope must have in the given dialect ("sqlite" or
// "postgres"). A hypertable layout carries no extension schema yet: that is
// probed evidence frozen at confirmation.
func LongStorageLayout(group *WriteGroup, dialect string) (*recordingplan.GroupSchemaLayout, error) {
	if !managedLongGroup(group) {
		return nil, fmt.Errorf("%w: not a managed long group", ErrWriteGroupValidation)
	}
	return longStorageScopeOf(group).layout(dialect)
}

// LongStorageIdentity returns the table name and owner-marker column the shared
// managed_long_v1 table of a long group must have. Both are derived from the
// storage scope (workspace, confirmed destination, storage mode, layout), so a
// runtime that recovers the group can prove it is writing the table it was
// confirmed for. ok is false for a group that is not a long group.
func LongStorageIdentity(group *WriteGroup) (table, ownerColumn string, ok bool) {
	if group == nil || !group.isLongLayout() {
		return "", "", false
	}
	scope := longStorageScopeOf(group)
	return scope.tableName(), scope.ownerColumn(), true
}
