package datalink

import (
	"database/sql"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestDefaultEmbeddedSQLiteDSNEnablesForeignKeysOnReusedConnection(t *testing.T) {
	databasePath := filepath.ToSlash(filepath.Join(t.TempDir(), "foreign-keys.db"))
	dsn := strings.Replace(DefaultEmbeddedSQLiteDSN, "datalink.db", databasePath, 1)
	db, err := sql.Open("sqlite", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	ApplySQLitePoolDefaults(db)
	ctx := t.Context()

	require.NoError(t, db.PingContext(ctx))
	assertSQLiteForeignKeysEnabled(t, db)
	require.NoError(t, db.PingContext(ctx))
	assertSQLiteForeignKeysEnabled(t, db)

	_, err = db.ExecContext(ctx, `
		CREATE TABLE fk_parent (id INTEGER PRIMARY KEY);
		CREATE TABLE fk_child (
			id INTEGER PRIMARY KEY,
			parent_id INTEGER NOT NULL REFERENCES fk_parent(id) ON DELETE CASCADE
		);
		INSERT INTO fk_parent (id) VALUES (1);
		INSERT INTO fk_child (id, parent_id) VALUES (1, 1);
		DELETE FROM fk_parent WHERE id = 1;
	`)
	require.NoError(t, err)

	var childCount int
	require.NoError(t, db.QueryRowContext(ctx, "SELECT COUNT(*) FROM fk_child").Scan(&childCount))
	require.Zero(t, childCount)
}

func assertSQLiteForeignKeysEnabled(t *testing.T, db *sql.DB) {
	t.Helper()
	var foreignKeys int
	require.NoError(t, db.QueryRowContext(t.Context(), "PRAGMA foreign_keys").Scan(&foreignKeys))
	require.Equal(t, 1, foreignKeys)
}
