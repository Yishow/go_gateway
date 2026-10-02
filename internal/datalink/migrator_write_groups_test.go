package datalink

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestWriteGroupsMigration_IdempotentAndPreservesRows(t *testing.T) {
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "write-groups.db"))
	require.NoError(t, err)
	defer db.Close()
	ctx := t.Context()
	migrator := NewMigrator()

	require.NoError(t, migrator.Migrate(db))
	_, err = db.ExecContext(ctx, `
		INSERT INTO write_groups (
			id, workspace_id, revision, name, destination_connector_id,
			destination_connector_revision, destination_table_name
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`, "group-1", "workspace-1", "revision-1", "group", "connector-1", "connector-revision-1", "raw_values")
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `
		INSERT INTO write_group_members (group_id, member_index, device_id, point_id, tag_id, target_column)
		VALUES (?, ?, ?, ?, ?, ?)
	`, "group-1", 0, "device-1", "point-1", "tag-1", "temperature")
	require.NoError(t, err)

	require.NoError(t, migrator.Migrate(db))
	var groupCount, memberCount int
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_groups`).Scan(&groupCount))
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_group_members`).Scan(&memberCount))
	require.Equal(t, 1, groupCount)
	require.Equal(t, 1, memberCount)
}

func TestWriteGroupsMigrationRepairsMissingMemberTableWithoutDroppingGroups(t *testing.T) {
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "write-groups-repair.db"))
	require.NoError(t, err)
	defer db.Close()
	ctx := t.Context()
	migrator := NewMigrator()

	require.NoError(t, migrator.Migrate(db))
	for _, id := range []string{"group-1", "group-2"} {
		_, err = db.ExecContext(ctx, `
			INSERT INTO write_groups (
				id, workspace_id, revision, name, destination_connector_id,
				destination_connector_revision, destination_table_name
			) VALUES (?, ?, ?, ?, ?, ?, ?)
		`, id, "workspace-1", "revision-"+id, id, "connector-1", "connector-revision-1", "raw_values")
		require.NoError(t, err)
	}
	_, err = db.ExecContext(ctx, `DROP TABLE write_group_versions`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `DROP TABLE write_group_members`)
	require.NoError(t, err)

	// The completed 024 migration is now partially absent. Re-running the
	// real migrator must recreate the missing table and retain groups already
	// stored in the surviving table.
	require.NoError(t, migrator.Migrate(db))
	rows, err := db.QueryContext(ctx, `SELECT id FROM write_groups ORDER BY id`)
	require.NoError(t, err)
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		require.NoError(t, rows.Scan(&id))
		ids = append(ids, id)
	}
	require.NoError(t, rows.Err())
	require.Equal(t, []string{"group-1", "group-2"}, ids)

	_, err = db.ExecContext(ctx, `
		INSERT INTO write_group_members (group_id, member_index, device_id, point_id, tag_id, target_column)
		VALUES (?, ?, ?, ?, ?, ?)
	`, "group-1", 0, "device-1", "point-1", "tag-1", "temperature")
	require.NoError(t, err)
	var versionTable string
	require.NoError(t, db.QueryRowContext(ctx,
		`SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'write_group_versions'`,
	).Scan(&versionTable))
	require.Equal(t, "write_group_versions", versionTable)
}

func TestWriteGroupsMigrationAddsIdentityMapToExistingGroupSchema(t *testing.T) {
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "write-group-map-upgrade.db"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	ctx := t.Context()
	migrator := NewMigrator()
	require.NoError(t, migrator.Migrate(db))
	_, err = db.ExecContext(ctx, `INSERT INTO write_groups
		(id,workspace_id,revision,name,destination_connector_id,destination_connector_revision,destination_table_name)
		VALUES ('retained-group','workspace','revision','Retained','connector','connector-revision','values')`)
	require.NoError(t, err)
	// Model a configuration database created before the reviewed migration map
	// was added: all other 024 tables already exist and contain their data.
	_, err = db.ExecContext(ctx, `DROP TABLE write_group_migration_maps`)
	require.NoError(t, err)
	require.NoError(t, migrator.Migrate(db))
	_, err = db.ExecContext(ctx, `INSERT INTO write_group_migration_maps
		(workspace_id,source_kind,source_id,source_revision,group_id,before_intent,review_digest,adapter_version)
		VALUES ('workspace','legacy-single-mapping','legacy-source','source-revision','retained-group','{}','review-digest','single-mapping-v1')`)
	require.NoError(t, err)
	require.NoError(t, migrator.Migrate(db))
	var groupName, mappedID string
	require.NoError(t, db.QueryRowContext(ctx, `SELECT name FROM write_groups WHERE id='retained-group'`).Scan(&groupName))
	require.Equal(t, "Retained", groupName)
	require.NoError(t, db.QueryRowContext(ctx, `SELECT group_id FROM write_group_migration_maps WHERE source_id='legacy-source'`).Scan(&mappedID))
	require.Equal(t, "retained-group", mappedID)
}

func TestWriteGroupsMigrationAddsEntityKeyToExistingMemberTable(t *testing.T) {
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "write-group-member-upgrade.db"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	ctx := t.Context()
	migrator := NewMigrator()
	require.NoError(t, migrator.Migrate(db))
	_, err = db.ExecContext(ctx, `INSERT INTO write_groups
		(id,workspace_id,revision,name,destination_connector_id,destination_connector_revision,destination_table_name)
		VALUES ('retained-group','workspace','revision','Retained','connector','connector-revision','values')`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `DROP TABLE write_group_members`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `CREATE TABLE write_group_members (
		group_id TEXT NOT NULL, member_index INTEGER NOT NULL, device_id TEXT NOT NULL,
		point_id TEXT NOT NULL, tag_id TEXT NOT NULL, source_revision TEXT NOT NULL DEFAULT '',
		mapping_revision TEXT NOT NULL DEFAULT '', measurement_id TEXT, target_column TEXT NOT NULL,
		required INTEGER NOT NULL DEFAULT 1, max_age_seconds INTEGER,
		PRIMARY KEY (group_id, member_index), FOREIGN KEY (group_id) REFERENCES write_groups(id) ON DELETE RESTRICT
	)`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO write_group_members
		(group_id,member_index,device_id,point_id,tag_id,source_revision,mapping_revision,target_column)
		VALUES ('retained-group',0,'device','point','tag','source','mapping','value')`)
	require.NoError(t, err)

	require.NoError(t, migrator.Migrate(db))
	var entityKey, sourceRevision string
	require.NoError(t, db.QueryRowContext(ctx, `SELECT entity_key, source_revision FROM write_group_members WHERE group_id='retained-group'`).Scan(&entityKey, &sourceRevision))
	require.Empty(t, entityKey)
	require.Equal(t, "source", sourceRevision)
	var entityColumnExists int
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM pragma_table_info('write_group_members') WHERE name='entity_key'`).Scan(&entityColumnExists))
	require.Equal(t, 1, entityColumnExists)
}
