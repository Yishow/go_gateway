package datalink

import (
	"database/sql"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/schema/migrations"

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
	// Build a genuine pre-029/pre-047 database. Dropping a key table after the
	// latest migration would be lineage corruption, not a historical schema.
	for _, name := range []string{"024_write_groups_sqlite.up.sql", "028_write_group_basic_keys_sqlite.up.sql"} {
		raw, err := migrations.FS.ReadFile(name)
		require.NoError(t, err)
		_, err = db.ExecContext(ctx, string(raw))
		require.NoError(t, err)
	}

	groupColumns := []string{"layout_version", "profile_kind", "profile_interval", "profile_unit", "storage_mode"}
	memberColumns := []string{"circuit_id", "circuit_name"}
	require.Zero(t, sqliteColumnCount(t, db, "write_groups", groupColumns...))
	require.Zero(t, sqliteColumnCount(t, db, "write_group_members", memberColumns...))
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

// TestRecordingReductionMigrationKeepsLegacyRowsNullable pins the 038 columns:
// the reduction and the member role/type columns exist after a migrate, are
// created again for a member table the 024 repair path recreated, stay NULL
// for legacy rows, and re-running the migrator is idempotent.
func TestRecordingReductionMigrationKeepsLegacyRowsNullable(t *testing.T) {
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "recording-reduction.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	ctx := t.Context()
	migrator := NewMigrator()
	require.NoError(t, migrator.Migrate(db))

	groupColumns := []string{"sample_reduction"}
	memberColumns := []string{"value_role", "source_value_type", "result_value_type"}
	require.Equal(t, len(groupColumns), sqliteColumnCount(t, db, "write_groups", groupColumns...))
	require.Equal(t, len(memberColumns), sqliteColumnCount(t, db, "write_group_members", memberColumns...))

	_, err = db.ExecContext(ctx, `INSERT INTO write_groups
		(id,workspace_id,revision,name,destination_connector_id,destination_connector_revision,destination_table_name,row_policy)
		VALUES ('legacy-group','workspace','revision','Legacy','connector','connector-revision','values','{"interval_seconds":60}')`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO write_group_members
		(group_id,member_index,device_id,point_id,tag_id,target_column)
		VALUES ('legacy-group',0,'device','point','tag','value')`)
	require.NoError(t, err)

	require.NoError(t, migrator.Migrate(db), "re-running must be idempotent")
	var sampleReduction sql.NullString
	require.NoError(t, db.QueryRowContext(ctx,
		`SELECT sample_reduction FROM write_groups WHERE id = 'legacy-group'`).Scan(&sampleReduction))
	require.False(t, sampleReduction.Valid, "a legacy group's reduction stays NULL")
	var valueRole, sourceType, resultType sql.NullString
	require.NoError(t, db.QueryRowContext(ctx,
		`SELECT value_role, source_value_type, result_value_type FROM write_group_members WHERE group_id = 'legacy-group'`).
		Scan(&valueRole, &sourceType, &resultType))
	require.False(t, valueRole.Valid || sourceType.Valid || resultType.Valid, "legacy member columns stay NULL")

	// A member table the 024 repair path recreated must be completed with the
	// role and type columns again.
	_, err = db.ExecContext(ctx, `DROP TABLE write_group_members`)
	require.NoError(t, err)
	require.NoError(t, migrator.Migrate(db))
	require.Equal(t, len(memberColumns), sqliteColumnCount(t, db, "write_group_members", memberColumns...))
}
