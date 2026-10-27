package recordingplan

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSQLServerNativeStorageConfirmationUpdatesLayoutDigest(t *testing.T) {
	db := lineageDB(t)
	migration, err := os.ReadFile(filepath.Join("..", "schema", "migrations", "042_sqlserver_partition_maintenance_sqlite.up.sql"))
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), string(migration))
	require.NoError(t, err)
	scope := StorageScope{
		WorkspaceID: "ws", ConnectorID: "conn", Database: "recordings", TableSchema: "dbo",
		StorageMode: "native_monthly", LayoutVersion: "managed_long_v1",
	}
	options := sqlServerSchemaOptionsForTest()
	native, err := NewSQLServerNativeMonthlyPartition(
		time.Date(2026, time.October, 15, 0, 0, 0, 0, time.UTC),
		"pf_sqlserver_long", "ps_sqlserver_long",
	)
	require.NoError(t, err)
	options.NativeMonthly = native
	layout, err := ManagedSQLServerLongGroupLayout("sqlserver_long", "_gw_owner_0123456789abcdef01234567", options)
	require.NoError(t, err)

	_, err = db.ExecContext(t.Context(), `INSERT INTO write_group_storage_scopes
		(workspace_id, connector_id, database_name, table_schema, storage_mode, layout_version, table_name, owner_column, layout_digest, generation)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 1)`,
		scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, scope.StorageMode, scope.LayoutVersion,
		layout.TableName, layout.OwnerColumn, UnconfirmedSQLServerStorageDigest(scope))
	require.NoError(t, err)

	now := time.Date(2026, time.October, 15, 1, 2, 3, 0, time.UTC)
	require.NoError(t, RecordStorageConfirmation(t.Context(), db, scope, StorageConfirmation{
		Layout: layout, ConnectorRevision: "rev-1", OperationID: "op-confirm",
		VerifiedDigest: strings.Repeat("a", 64), CreatedTable: true, GroupID: "group-1",
	}, now))
	lineage, err := ReadStorageLineage(t.Context(), db, scope)
	require.NoError(t, err)
	require.True(t, lineage.Confirmed())
	require.NotNil(t, lineage.Layout.SQLServer.NativeMonthly)
	encoded, err := json.Marshal(lineage.Layout)
	require.NoError(t, err)
	sum := sha256.Sum256(encoded)
	require.Equal(t, hex.EncodeToString(sum[:]), lineage.LayoutDigest)
	require.NotEqual(t, UnconfirmedSQLServerStorageDigest(scope), lineage.LayoutDigest)
	require.Equal(t, 1, lineage.Generation)

	require.NoError(t, RecordStorageConfirmation(t.Context(), db, scope, StorageConfirmation{
		Layout: layout, ConnectorRevision: "rev-1", OperationID: "op-confirm-2",
		VerifiedDigest: strings.Repeat("b", 64), CreatedTable: true, GroupID: "group-1",
	}, now.Add(time.Minute)))
	rebuilt, err := ReadStorageLineage(t.Context(), db, scope)
	require.NoError(t, err)
	require.Equal(t, 2, rebuilt.Generation)
	require.True(t, rebuilt.HistoryMayBeLost)
}
