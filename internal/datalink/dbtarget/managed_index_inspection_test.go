package dbtarget

import (
	"os"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

const (
	indexTestTable = "gw_long_0123456789abcdef"
	indexTestOwner = "_gw_owner_0123456789abcdef01234567"
)

func longDDL(t *testing.T, dialect string) []string {
	t.Helper()
	layout, err := recordingplan.ManagedLongGroupLayout(dialect, indexTestTable, indexTestOwner)
	require.NoError(t, err)
	statements, err := recordingplan.GenerateGroupSchemaDDL(dialect, layout)
	require.NoError(t, err)
	return statements
}

func TestManagedIndexInspectionSQLite(t *testing.T) {
	statements := longDDL(t, "sqlite")
	target := newInspectionSQLiteTarget(t, append(statements,
		`CREATE UNIQUE INDEX idx_extra_unique ON `+indexTestTable+` (group_id, record_id)`,
		`CREATE INDEX idx_extra_partial ON `+indexTestTable+` (device_id) WHERE circuit_id IS NOT NULL`,
	)...)
	fixture := newManagedDestinationFixture(t, schema.DatabaseConnectorKindSQLite, ConnectionConfig{"dsn": target})

	indexes, err := fixture.inspector.InspectTableIndexesAtRevision(t.Context(), fixture.connector.ID, "managed-revision", "main", indexTestTable)

	require.NoError(t, err)
	require.Equal(t, []IndexInfo{
		{Name: "idx_extra_unique", Columns: []string{"group_id", "record_id"}, Unique: true},
		{Name: "idx_" + indexTestTable + "_src_time", Columns: []string{"workspace_id", "tag_id", "bucket_start"}},
	}, indexes, "the primary-key index and the partial index are not reported; key order is preserved")

	t.Run("a missing table has no indexes and a missing file is not created", func(t *testing.T) {
		none, err := fixture.inspector.InspectTableIndexesAtRevision(t.Context(), fixture.connector.ID, "managed-revision", "main", "gw_missing")
		require.NoError(t, err)
		require.Empty(t, none)

		missing := filepath.Join(t.TempDir(), "absent.db")
		absent := newManagedDestinationFixture(t, schema.DatabaseConnectorKindSQLite, ConnectionConfig{"dsn": missing})
		none, err = absent.inspector.InspectTableIndexesAtRevision(t.Context(), absent.connector.ID, "managed-revision", "main", indexTestTable)
		require.NoError(t, err)
		require.Empty(t, none)
		_, statErr := os.Stat(missing)
		require.ErrorIs(t, statErr, os.ErrNotExist)
	})

	t.Run("a stale connector revision is refused", func(t *testing.T) {
		_, err := fixture.inspector.InspectTableIndexesAtRevision(t.Context(), fixture.connector.ID, "another-revision", "main", indexTestTable)
		require.Error(t, err)
	})

	t.Run("the gateway's own database is never inspected", func(t *testing.T) {
		internal := newManagedDestinationFixture(t, schema.DatabaseConnectorKindSQLite, ConnectionConfig{"dsn": ""})
		internal.connector.ConnectionConfig = `{"dsn":"` + filepath.ToSlash(internal.internalPath) + `"}`
		require.NoError(t, NewSQLConnectorRepository(internal.internalDB).Update(t.Context(), internal.connector))
		_, err := internal.inspector.InspectTableIndexesAtRevision(t.Context(), internal.connector.ID, "managed-revision", "main", "write_groups")
		require.ErrorIs(t, err, ErrManagedDestinationInternal)
	})
}

func TestManagedIndexInspectionPostgres(t *testing.T) {
	dsn, fields := postgresKeywordDSN(t)
	admin := postgresAdminDB(t, dsn)
	schemaName := postgresScopedSchema(t, admin)
	tx, err := admin.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	_, err = tx.ExecContext(t.Context(), `SET LOCAL search_path TO `+pgIdentifier(t, schemaName))
	require.NoError(t, err)
	for _, statement := range longDDL(t, "postgres") {
		_, err = tx.ExecContext(t.Context(), statement)
		require.NoError(t, err, statement)
	}
	_, err = tx.ExecContext(t.Context(), `CREATE UNIQUE INDEX idx_extra_unique ON `+indexTestTable+` (group_id, record_id)`)
	require.NoError(t, err)
	_, err = tx.ExecContext(t.Context(), `CREATE INDEX idx_extra_partial ON `+indexTestTable+` (device_id) WHERE circuit_id IS NOT NULL`)
	require.NoError(t, err)
	_, err = tx.ExecContext(t.Context(), `CREATE INDEX idx_extra_expression ON `+indexTestTable+` (lower(device_id))`)
	require.NoError(t, err)
	require.NoError(t, tx.Commit())

	fixture := newManagedDestinationFixture(t, schema.DatabaseConnectorKindPostgres,
		postgresConnectorConfig(fields, schemaName, fields["user"], fields["password"]))
	indexes, err := fixture.inspector.InspectTableIndexesAtRevision(t.Context(), fixture.connector.ID, "managed-revision", schemaName, indexTestTable)
	require.NoError(t, err)
	require.Equal(t, []IndexInfo{
		{Name: "idx_extra_unique", Columns: []string{"group_id", "record_id"}, Unique: true},
		{Name: "idx_" + indexTestTable + "_src_time", Columns: []string{"workspace_id", "tag_id", "bucket_start"}},
	}, indexes, "partial and expression indexes cannot satisfy the layout and are not reported")

	none, err := fixture.inspector.InspectTableIndexesAtRevision(t.Context(), fixture.connector.ID, "managed-revision", schemaName, "gw_missing")
	require.NoError(t, err)
	require.Empty(t, none)
}

// TestManagedIndexInspectionPostgresIgnoresUnusableIndexes proves that an index
// with the right name and columns but without usable catalog state does not
// count as the confirmed index. pg_index flags are flipped on the disposable
// fixture to reproduce the states a failed concurrent build, an in-progress
// build and a dropping index leave behind.
func TestManagedIndexInspectionPostgresIgnoresUnusableIndexes(t *testing.T) {
	dsn, fields := postgresKeywordDSN(t)
	admin := postgresAdminDB(t, dsn)
	schemaName := postgresScopedSchema(t, admin)
	tx, err := admin.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	_, err = tx.ExecContext(t.Context(), `SET LOCAL search_path TO `+pgIdentifier(t, schemaName))
	require.NoError(t, err)
	for _, statement := range longDDL(t, "postgres") {
		_, err = tx.ExecContext(t.Context(), statement)
		require.NoError(t, err, statement)
	}
	require.NoError(t, tx.Commit())
	fixture := newManagedDestinationFixture(t, schema.DatabaseConnectorKindPostgres,
		postgresConnectorConfig(fields, schemaName, fields["user"], fields["password"]))
	indexName := "idx_" + indexTestTable + "_src_time"
	want := []IndexInfo{{Name: indexName, Columns: []string{"workspace_id", "tag_id", "bucket_start"}}}
	inspect := func() []IndexInfo {
		indexes, err := fixture.inspector.InspectTableIndexesAtRevision(t.Context(), fixture.connector.ID, "managed-revision", schemaName, indexTestTable)
		require.NoError(t, err)
		return indexes
	}
	setFlag := func(flag string, value bool) {
		_, err := admin.ExecContext(t.Context(), `UPDATE pg_catalog.pg_index SET `+flag+` = $1
			WHERE indexrelid = (quote_ident($2) || '.' || quote_ident($3))::regclass`, value, schemaName, indexName)
		require.NoError(t, err)
	}

	require.Equal(t, want, inspect(), "a healthy index is reported")
	for _, flag := range []string{"indisvalid", "indisready", "indislive"} {
		t.Run(flag, func(t *testing.T) {
			setFlag(flag, false)
			require.Empty(t, inspect(), "an index with %s=false is not usable and must not be reported", flag)
			setFlag(flag, true)
			require.Equal(t, want, inspect(), "restoring %s makes the index usable again", flag)
		})
	}
}
