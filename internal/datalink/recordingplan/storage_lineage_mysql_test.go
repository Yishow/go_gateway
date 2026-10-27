package recordingplan

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMySQLStorageConfirmationBindsActualLayoutAfterDraft(t *testing.T) {
	db := lineageDB(t)
	scope := StorageScope{WorkspaceID: "ws", ConnectorID: "conn", Database: "recordings", TableSchema: "recordings", StorageMode: "plain", LayoutVersion: "managed_long_v1"}
	digest := func(purpose string) string {
		return fmt.Sprintf("%x", sha256.Sum256([]byte(strings.Join([]string{purpose, scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, scope.StorageMode, scope.LayoutVersion}, "\x00"))))
	}
	layout, err := ManagedMySQLLongGroupLayout("gw_long_"+digest("table")[:16], "_gw_owner_"+digest("owner")[:24], MySQLGroupSchemaOptions{
		ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "cap-1", Engine: "InnoDB", Collation: "utf8mb4_0900_bin",
	})
	require.NoError(t, err)
	seed := digest("mysql-storage-unconfirmed-v1")
	_, err = db.ExecContext(t.Context(), `INSERT INTO write_group_storage_scopes (workspace_id, connector_id, database_name, table_schema, storage_mode, layout_version,
		table_name, owner_column, layout_digest, generation) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 1)`,
		scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, scope.StorageMode, scope.LayoutVersion, layout.TableName, layout.OwnerColumn, seed)
	require.NoError(t, err)
	confirmation := StorageConfirmation{Layout: layout, ConnectorRevision: "rev-1", OperationID: "op-1", VerifiedDigest: strings.Repeat("d", 64), CreatedTable: true, GroupID: "group-1"}
	_, err = db.ExecContext(t.Context(), `CREATE TRIGGER reject_mysql_confirmation BEFORE INSERT ON write_group_storage_events BEGIN SELECT RAISE(ABORT, 'event unavailable'); END`)
	require.NoError(t, err)
	require.Error(t, RecordStorageConfirmation(t.Context(), db, scope, confirmation, time.Now()))
	unconfirmed, err := ReadStorageLineage(t.Context(), db, scope)
	require.NoError(t, err)
	require.False(t, unconfirmed.Confirmed(), "failed proof persistence rolls back the digest as well")
	require.Equal(t, seed, unconfirmed.LayoutDigest)
	_, err = db.ExecContext(t.Context(), `DROP TRIGGER reject_mysql_confirmation`)
	require.NoError(t, err)
	require.NoError(t, RecordStorageConfirmation(t.Context(), db, scope, confirmation, time.Now()))
	stored, err := ReadStorageLineage(t.Context(), db, scope)
	require.NoError(t, err)
	encoded, err := json.Marshal(layout)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("%x", sha256.Sum256(encoded)), stored.LayoutDigest)
	require.NotEqual(t, seed, stored.LayoutDigest, "the unverified draft seed cannot stand in for the verified family-bound layout")
	require.True(t, stored.Confirmed())
	require.Equal(t, 1, stored.Generation)
	require.Equal(t, layout, *stored.Layout)
	wrong := confirmation
	wrong.Layout.TableName = "gw_long_different"
	require.Error(t, RecordStorageConfirmation(t.Context(), db, scope, wrong, time.Now()), "same operation cannot accept altered evidence")
	after, err := ReadStorageLineage(t.Context(), db, scope)
	require.NoError(t, err)
	require.Equal(t, stored.LayoutDigest, after.LayoutDigest)
	require.Equal(t, *stored.Layout, *after.Layout)
}
