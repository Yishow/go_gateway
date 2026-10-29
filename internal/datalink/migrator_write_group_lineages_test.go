package datalink

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"go-gateway/internal/datalink/schema/migrations"

	"github.com/stretchr/testify/require"
)

func TestWriteGroupLineagesMigrationFreshUpgradeAndReopen(t *testing.T) {
	for _, legacy := range []bool{false, true} {
		name := "fresh"
		if legacy {
			name = "pre047"
		}
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "lineage.db")
			db := openLineageMigrationDB(t, path)
			var columns map[string][]string
			var before map[string][]string
			var keySchema []string
			if legacy {
				createPre047LineageSchema(t, db)
				seedLegacyLineageRows(t, db)
				columns = lineageSnapshotColumns(t, db, lineagePreservedTables...)
				before = lineageMigrationSnapshot(t, db, columns)
				keySchema = lineageOriginalKeySchema(t, db)
			}
			_, err := db.Exec(`PRAGMA foreign_keys=ON`)
			require.NoError(t, err)
			require.NoError(t, NewMigrator().Migrate(db))
			var count int
			require.NoError(t, db.QueryRow(`SELECT count(*) FROM write_group_lineages`).Scan(&count))
			want := 0
			if legacy {
				want = 3
				require.Equal(t, before, lineageMigrationSnapshot(t, db, columns))
				require.Equal(t, keySchema, lineageOriginalKeySchema(t, db), "old column metadata, primary/unique keys and indexes must be unchanged")
				assertInitialLineageRows(t, db)
			}
			require.Equal(t, want, count)
			all := lineageSnapshotColumns(t, db, append(lineagePreservedTables, "write_group_lineages", "write_group_generations")...)
			snapshot := lineageMigrationSnapshot(t, db, all)
			require.NoError(t, NewMigrator().Migrate(db))
			require.Equal(t, snapshot, lineageMigrationSnapshot(t, db, all))
			require.NoError(t, db.Close())
			db = openLineageMigrationDB(t, path)
			require.NoError(t, NewMigrator().Migrate(db))
			require.Equal(t, snapshot, lineageMigrationSnapshot(t, db, all))
		})
	}
}

func openLineageMigrationDB(t *testing.T, path string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	return db
}

// This is a genuine pre-047 database, never latest-Migrate with deleted markers.
func createPre047LineageSchema(t *testing.T, db *sql.DB) {
	t.Helper()
	for _, name := range []string{
		"023_schema_operations_sqlite.up.sql", "026_operation_test_write_sqlite.up.sql",
		"024_write_groups_sqlite.up.sql", "025_write_group_delivery_sqlite.up.sql",
		"028_write_group_basic_keys_sqlite.up.sql", "029_write_group_circuit_profiles_sqlite.up.sql",
		"030_write_group_profile_keys_sqlite.up.sql", "031_write_group_storage_mode_sqlite.up.sql",
		"038_write_group_recording_reduction_sqlite.up.sql", "040_write_group_exact_recording_source_sqlite.up.sql",
		"045_write_group_sample_handoffs_sqlite.up.sql", "046_write_group_handoff_convergence_sqlite.up.sql",
	} {
		raw, err := migrations.FS.ReadFile(name)
		require.NoError(t, err)
		_, err = db.Exec(string(raw))
		require.NoError(t, err, name)
	}
}

func seedLegacyLineageRows(t *testing.T, db *sql.DB) {
	t.Helper()
	for _, group := range []string{"G1", "basic-only", "profile-only", "unkeyed"} {
		insertLineageMigrationGroup(t, db, group)
	}
	_, err := db.Exec(`UPDATE write_groups SET status='deleted' WHERE id='G1';
		INSERT INTO write_group_basic_keys(workspace_id,device_id,canonical_role,group_id,created_at)
		VALUES ('workspace','device-G1','basic','G1','2001-02-03 04:05:06'),
		('workspace','device-basic','basic','basic-only','2002-03-04 05:06:07');
		INSERT INTO write_group_profile_keys(workspace_id,circuit_id,selection_scope,profile_kind,canonical_role,group_id,created_at)
		VALUES ('workspace','','device-G1','instantaneous','unassigned','G1','2003-04-05 06:07:08'),
		('workspace','circuit-P','','cumulative','circuit','profile-only','2004-05-06 07:08:09');
		INSERT INTO write_group_members(group_id,member_index,device_id,point_id,tag_id,target_column)
		VALUES ('G1',0,'device-G1','point','tag','value');
		INSERT INTO write_group_versions(group_id,group_revision,effective_at,payload)
		VALUES ('G1','R1','2000-01-01',' { "frozen" : [1,2], "revision":"R1" } ');
		INSERT INTO wg_delivery_samples(workspace_id,group_id,group_revision,sample_id,member_key,observed_at,bucket_start,payload,payload_digest,created_at)
		VALUES ('workspace','G1','R1','sample','member','25','20',' { "value" : 25.0 } ','sample-digest','26');
		INSERT INTO wg_delivery_outbox(effect_key,record_id,workspace_id,group_id,group_revision,bucket_start,partition_key,destination_scope,connector_id,connector_revision,table_name,payload,payload_digest,next_retry_at,created_at,updated_at)
		VALUES ('effect-G1','record-G1','workspace','G1','R1','20','partition','destination','connector','C1','values',' { "value" : 25.0 } ','outbox-digest','30','26','27');
		INSERT INTO wg_delivery_receipts(effect_key,payload_digest,committed_at)
		VALUES ('effect-G1','receipt-digest','28');
		INSERT INTO managed_schema_operations(operation_id,token,workspace_id,scope_key,owner,action,status,created_at,updated_at,detail)
		VALUES ('operation-G1','token','workspace','scope','G1','create','completed','10','11',' { "unchanged" : true } ');`)
	require.NoError(t, err)
}

func insertLineageMigrationGroup(t *testing.T, db *sql.DB, group string) {
	t.Helper()
	_, err := db.Exec(`INSERT INTO write_groups(id,workspace_id,revision,name,destination_connector_id,destination_connector_revision,destination_table_name,row_policy,write_policy,migration,created_at,updated_at)
		VALUES (?,'workspace','R1',?,'connector','C1','values',' { "interval_seconds" : 60 } ',' { "enabled" : true } ',' { "history" : "retain" } ','1999-01-01','2000-01-01')`, group, group)
	require.NoError(t, err)
}

func assertInitialLineageRows(t *testing.T, db *sql.DB) {
	t.Helper()
	for _, group := range []string{"G1", "basic-only", "profile-only"} {
		var id, workspace, current string
		var revision, generation int
		var reserved, operation, intent, predecessor sql.NullString
		require.NoError(t, db.QueryRow(`SELECT l.lineage_id,l.workspace_id,l.current_group_id,l.revision,l.reserved_group_id,l.reservation_operation_id,l.reservation_intent_digest,g.generation,g.predecessor_group_id
			FROM write_group_generations g JOIN write_group_lineages l USING(lineage_id) WHERE g.group_id=?`, group).
			Scan(&id, &workspace, &current, &revision, &reserved, &operation, &intent, &generation, &predecessor))
		require.Regexp(t, `^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$`, id)
		require.Equal(t, "workspace", workspace)
		require.Equal(t, group, current)
		require.Equal(t, 1, revision)
		require.Equal(t, 1, generation)
		require.False(t, reserved.Valid || operation.Valid || intent.Valid || predecessor.Valid)
	}
	var basicID, profileID string
	require.NoError(t, db.QueryRow(`SELECT lineage_id FROM write_group_basic_keys WHERE group_id='G1'`).Scan(&basicID))
	require.NoError(t, db.QueryRow(`SELECT lineage_id FROM write_group_profile_keys WHERE group_id='G1'`).Scan(&profileID))
	require.Equal(t, basicID, profileID, "dual aliases share one lineage")
	var count int
	require.NoError(t, db.QueryRow(`SELECT count(*) FROM write_group_generations WHERE group_id='unkeyed'`).Scan(&count))
	require.Zero(t, count, "unkeyed legacy groups receive no invented lineage")
}

var lineagePreservedTables = []string{
	"write_groups", "write_group_members", "write_group_versions", "write_group_basic_keys", "write_group_profile_keys",
	"wg_delivery_samples", "wg_delivery_outbox", "wg_delivery_receipts", "managed_schema_operations",
}

func lineageSnapshotColumns(t *testing.T, db *sql.DB, tables ...string) map[string][]string {
	t.Helper()
	result := make(map[string][]string)
	for _, table := range tables {
		rows, err := db.Query(`SELECT name FROM pragma_table_info(?) ORDER BY cid`, table)
		require.NoError(t, err)
		for rows.Next() {
			var name string
			require.NoError(t, rows.Scan(&name))
			result[table] = append(result[table], `"`+name+`"`)
		}
		require.NoError(t, rows.Err())
		require.NoError(t, rows.Close())
	}
	return result
}

func lineageMigrationSnapshot(t *testing.T, db *sql.DB, columns map[string][]string) map[string][]string {
	t.Helper()
	result := make(map[string][]string)
	for table, names := range columns {
		rows, err := db.Query(`SELECT ` + strings.Join(names, ",") + ` FROM "` + table + `"`)
		require.NoError(t, err)
		result[table] = nil
		for rows.Next() {
			values, pointers := make([]any, len(names)), make([]any, len(names))
			for i := range values {
				pointers[i] = &values[i]
			}
			require.NoError(t, rows.Scan(pointers...))
			raw, err := json.Marshal(values)
			require.NoError(t, err)
			result[table] = append(result[table], string(raw))
		}
		require.NoError(t, rows.Err())
		require.NoError(t, rows.Close())
		sort.Strings(result[table])
	}
	return result
}

func lineageOriginalKeySchema(t *testing.T, db *sql.DB) []string {
	t.Helper()
	var result []string
	for _, table := range []string{"write_group_basic_keys", "write_group_profile_keys"} {
		rows, err := db.Query(`SELECT json_array(cid,name,type,"notnull",dflt_value,pk) FROM pragma_table_info(?) WHERE name<>'lineage_id'
			UNION ALL SELECT json_array(name,sql) FROM sqlite_master WHERE type='index' AND tbl_name=? ORDER BY 1`, table, table)
		require.NoError(t, err)
		for rows.Next() {
			var row string
			require.NoError(t, rows.Scan(&row))
			result = append(result, table+":"+row)
		}
		require.NoError(t, rows.Err())
		require.NoError(t, rows.Close())
	}
	return result
}
