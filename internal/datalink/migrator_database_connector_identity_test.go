package datalink

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestEnsureSQLiteDatabaseConnectorIdentityRevisionMigratesAndBackfills(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	ctx := t.Context()

	_, err = db.ExecContext(ctx, `
		CREATE TABLE database_connectors (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL
		);
		INSERT INTO database_connectors (id, name) VALUES ('legacy', 'Legacy');
	`)
	require.NoError(t, err)

	require.NoError(t, ensureSQLiteDatabaseConnectorIdentityRevision(db))
	var firstRevision string
	require.NoError(t, db.QueryRowContext(ctx, `SELECT identity_revision FROM database_connectors WHERE id = 'legacy'`).Scan(&firstRevision))
	require.NotEmpty(t, firstRevision)

	require.NoError(t, ensureSQLiteDatabaseConnectorIdentityRevision(db))
	var secondRevision string
	require.NoError(t, db.QueryRowContext(ctx, `SELECT identity_revision FROM database_connectors WHERE id = 'legacy'`).Scan(&secondRevision))
	require.Equal(t, firstRevision, secondRevision)
}

func TestMigrator_MigrateBackfillsEmptyDatabaseConnectorIdentityRevision(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	ctx := t.Context()
	migrator := NewMigrator()

	require.NoError(t, migrator.Migrate(db))
	_, err = db.ExecContext(ctx, `
		INSERT INTO database_connectors
		(id, name, kind, connection_config, status, enabled, created_at, updated_at)
		VALUES ('empty-revision', 'Empty Revision', 'sqlite', '{}', 'ready', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`)
	require.NoError(t, err)
	require.NoError(t, migrator.Migrate(db))

	var revision string
	require.NoError(t, db.QueryRowContext(ctx, `SELECT identity_revision FROM database_connectors WHERE id = 'empty-revision'`).Scan(&revision))
	require.NotEmpty(t, revision)
}
