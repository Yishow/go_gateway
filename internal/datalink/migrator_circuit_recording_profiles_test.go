package datalink

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func sqliteColumnCount(t *testing.T, db *sql.DB, table string, columns ...string) int {
	t.Helper()
	count := 0
	for _, column := range columns {
		exists, err := sqliteColumnExists(db, table, column)
		require.NoError(t, err)
		if exists {
			count++
		}
	}
	return count
}

func TestCircuitRecordingProfilesMigrationKeepsLegacyRowsNullable(t *testing.T) {
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "circuit-profiles.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	ctx := t.Context()
	migrator := NewMigrator()
	require.NoError(t, migrator.Migrate(db))

	groupColumns := []string{"layout_version", "profile_kind", "profile_interval", "profile_unit", "storage_mode"}
	memberColumns := []string{"circuit_id", "circuit_name"}
	require.Equal(t, len(groupColumns), sqliteColumnCount(t, db, "write_groups", groupColumns...))
	require.Equal(t, len(memberColumns), sqliteColumnCount(t, db, "write_group_members", memberColumns...))

	// Model a database created before 029/030: remove only what they added.
	_, err = db.ExecContext(ctx, `DROP TABLE write_group_profile_keys`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `DROP TABLE write_group_storage_scopes`)
	require.NoError(t, err)
	for _, column := range groupColumns {
		_, err = db.ExecContext(ctx, `ALTER TABLE write_groups DROP COLUMN `+column)
		require.NoError(t, err)
	}
	for _, column := range memberColumns {
		_, err = db.ExecContext(ctx, `ALTER TABLE write_group_members DROP COLUMN `+column)
		require.NoError(t, err)
	}
	_, err = db.ExecContext(ctx, `INSERT INTO write_groups
		(id,workspace_id,revision,name,destination_connector_id,destination_connector_revision,destination_table_name,row_policy)
		VALUES ('legacy-group','workspace','revision','Legacy','connector','connector-revision','values','{"interval_seconds":60}')`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO write_group_members
		(group_id,member_index,device_id,point_id,tag_id,target_column)
		VALUES ('legacy-group',0,'device','point','tag','value')`)
	require.NoError(t, err)

	require.NoError(t, migrator.Migrate(db))
	require.NoError(t, migrator.Migrate(db), "re-running must be idempotent")

	require.Equal(t, len(groupColumns), sqliteColumnCount(t, db, "write_groups", groupColumns...))
	require.Equal(t, len(memberColumns), sqliteColumnCount(t, db, "write_group_members", memberColumns...))
	var layout, kind, unit, mode sql.NullString
	var interval sql.NullInt64
	var rowPolicy string
	require.NoError(t, db.QueryRowContext(ctx, `
		SELECT layout_version, profile_kind, profile_interval, profile_unit, storage_mode, row_policy
		FROM write_groups WHERE id = 'legacy-group'`).Scan(&layout, &kind, &interval, &unit, &mode, &rowPolicy))
	require.False(t, layout.Valid || kind.Valid || interval.Valid || unit.Valid || mode.Valid, "legacy profile columns must stay NULL")
	var scopes int
	require.NoError(t, db.QueryRowContext(ctx, `SELECT count(*) FROM write_group_storage_scopes`).Scan(&scopes))
	require.Zero(t, scopes)
	require.JSONEq(t, `{"interval_seconds":60}`, rowPolicy)
	var circuitID, circuitName sql.NullString
	require.NoError(t, db.QueryRowContext(ctx, `
		SELECT circuit_id, circuit_name FROM write_group_members WHERE group_id = 'legacy-group'`).Scan(&circuitID, &circuitName))
	require.False(t, circuitID.Valid || circuitName.Valid)
}

func TestCircuitRecordingProfilesMigrationCompletesRecreatedMemberTable(t *testing.T) {
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "circuit-profiles-repair.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	ctx := t.Context()
	migrator := NewMigrator()
	require.NoError(t, migrator.Migrate(db))

	// The 024 repair path recreates a missing member table without 029 columns.
	_, err = db.ExecContext(ctx, `DROP TABLE write_group_members`)
	require.NoError(t, err)
	require.NoError(t, migrator.Migrate(db))
	require.Equal(t, 2, sqliteColumnCount(t, db, "write_group_members", "circuit_id", "circuit_name"))
}

func TestCircuitRecordingProfileKeysRejectAmbiguousScope(t *testing.T) {
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "circuit-profile-keys.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	ctx := t.Context()
	require.NoError(t, NewMigrator().Migrate(db))
	_, err = db.ExecContext(ctx, `INSERT INTO write_groups
		(id,workspace_id,revision,name,destination_connector_id,destination_connector_revision,destination_table_name)
		VALUES ('group-1','workspace','revision','Group','connector','connector-revision','values')`)
	require.NoError(t, err)

	insert := func(circuitID, scope, kind, role string) error {
		_, err := db.ExecContext(ctx, `INSERT INTO write_group_profile_keys
			(workspace_id,circuit_id,selection_scope,profile_kind,canonical_role,group_id)
			VALUES ('workspace',?,?,?,?,'group-1')`, circuitID, scope, kind, role)
		return err
	}
	require.Error(t, insert("circuit", "device", "instantaneous", "circuit"), "circuit role must not carry a device scope")
	require.Error(t, insert("", "", "instantaneous", "circuit"), "circuit role needs a circuit id")
	require.Error(t, insert("circuit", "", "instantaneous", "unassigned"), "unassigned role must not carry a circuit id")
	require.Error(t, insert("", "device", "cumulative", "unassigned"), "unassigned scope records only instantaneous")
	require.NoError(t, insert("", "device", "instantaneous", "unassigned"))
}
