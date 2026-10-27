package datalink

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestExactRecordingSourceMigration(t *testing.T) {
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "exact-source.db"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	ctx := t.Context()
	migrator := NewMigrator()
	require.NoError(t, migrator.Migrate(db))
	exists, err := sqliteColumnExists(db, "write_group_members", "exact_recording_source")
	require.NoError(t, err)
	require.True(t, exists, "startup must install the frozen source column")
	_, err = db.ExecContext(ctx, `INSERT INTO write_groups
		(id,workspace_id,revision,name,destination_connector_id,destination_connector_revision,destination_table_name)
		VALUES ('legacy','workspace','revision','Legacy','connector','connector-revision','values')`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO write_group_members
		(group_id,member_index,device_id,point_id,tag_id,target_column)
		VALUES ('legacy',0,'device','point','tag','value')`)
	require.NoError(t, err)
	var proof sql.NullString
	require.NoError(t, db.QueryRowContext(ctx,
		`SELECT exact_recording_source FROM write_group_members WHERE group_id='legacy'`).Scan(&proof))
	require.False(t, proof.Valid, "legacy members cannot acquire inferred scale evidence")
	const saved = `{"scale":"0.10","definition_revision":"def-1"}`
	_, err = db.ExecContext(ctx, `UPDATE write_group_members SET exact_recording_source=?`, saved)
	require.NoError(t, err)
	require.NoError(t, migrator.Migrate(db))
	require.NoError(t, db.QueryRowContext(ctx,
		`SELECT exact_recording_source FROM write_group_members WHERE group_id='legacy'`).Scan(&proof))
	require.Equal(t, saved, proof.String, "migration preserves existing proof bytes")
	_, err = db.ExecContext(ctx, `DROP TABLE write_group_members`)
	require.NoError(t, err)
	require.NoError(t, migrator.Migrate(db), "the existing member-table repair path must reinstall the additive column")
	exists, err = sqliteColumnExists(db, "write_group_members", "exact_recording_source")
	require.NoError(t, err)
	require.True(t, exists)
}
