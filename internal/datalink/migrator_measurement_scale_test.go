package datalink

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestExactRecordingScaleMigration(t *testing.T) {
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "measurement-scale.db"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	ctx := t.Context()
	migrator := NewMigrator()
	require.NoError(t, migrator.Migrate(db))
	exists, err := sqliteColumnExists(db, "measurement_definitions", "exact_recording_scale")
	require.NoError(t, err)
	require.True(t, exists, "production startup must install the exact scale column")

	// Model an existing configuration database from before the additive column.
	_, err = db.ExecContext(ctx, `ALTER TABLE measurement_definitions DROP COLUMN exact_recording_scale`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO measurement_definitions
		(id, workspace_id, device_id, point_id, equipment_id, definition_revision,
		source_binding_revision, series_epoch, name, quantity, semantic_kind)
		VALUES ('legacy', 'workspace', 'device', 'point', 'equipment', 'rev-1',
		'rev-1', 'epoch-1', 'Energy', 'energy', 'counter')`)
	require.NoError(t, err)
	require.NoError(t, migrator.Migrate(db))
	var scale sql.NullString
	require.NoError(t, db.QueryRowContext(ctx,
		`SELECT exact_recording_scale FROM measurement_definitions WHERE id='legacy'`).Scan(&scale))
	require.False(t, scale.Valid, "migration cannot infer scale for existing definitions")

	_, err = db.ExecContext(ctx,
		`UPDATE measurement_definitions SET exact_recording_scale='0.10' WHERE id='legacy'`)
	require.NoError(t, err)
	require.NoError(t, migrator.Migrate(db), "startup migration must be idempotent")
	require.NoError(t, db.QueryRowContext(ctx,
		`SELECT exact_recording_scale FROM measurement_definitions WHERE id='legacy'`).Scan(&scale))
	require.True(t, scale.Valid)
	require.Equal(t, "0.10", scale.String, "repeated startup preserves saved decimal text")
}
