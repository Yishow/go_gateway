package datalink

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func recordProfileSchema(t *testing.T, db *sql.DB) string {
	t.Helper()
	var schema string
	require.NoError(t, db.QueryRow(`SELECT sql FROM sqlite_master WHERE name='write_group_profile_keys'`).Scan(&schema))
	return schema
}

func TestRecordProfileKeysUnknownSchemaRefused(t *testing.T) {
	for _, alteration := range []string{
		`ALTER TABLE write_group_profile_keys ADD COLUMN unexpected TEXT`,
		`ALTER TABLE write_group_profile_keys RENAME TO old_keys; CREATE TABLE write_group_profile_keys AS SELECT * FROM old_keys`,
	} {
		t.Run(alteration, func(t *testing.T) {
			db := recordProfileDB(t)
			_, err := db.Exec(alteration)
			require.NoError(t, err)
			before := recordProfileSchema(t, db)
			require.ErrorContains(t, runSQLiteConfigurationExpansion(db, migrateMixedRecordProfileKeysTx), "unknown write_group_profile_keys schema")
			require.Equal(t, before, recordProfileSchema(t, db))
		})
	}
}

func addRecordProfileObjects(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.Exec(`CREATE INDEX record_key_index ON write_group_profile_keys(created_at);
 CREATE TABLE record_key_audit(value TEXT);
 CREATE TRIGGER record_key_trigger AFTER UPDATE OF created_at ON write_group_profile_keys BEGIN INSERT INTO record_key_audit VALUES(new.created_at); END;
 CREATE VIEW record_key_view AS SELECT group_id,created_at FROM write_group_profile_keys;
 CREATE TABLE record_key_child(group_id TEXT REFERENCES write_group_profile_keys(group_id));
 INSERT INTO record_key_child VALUES('G1');`)
	require.NoError(t, err)
}

func TestRecordProfileKeysSchemaObjectsAndConstraints(t *testing.T) {
	db := recordProfileDB(t)
	addRecordProfileObjects(t, db)
	require.NoError(t, runSQLiteConfigurationExpansion(db, migrateMixedRecordProfileKeysTx))
	_, err := db.Exec(`UPDATE write_group_profile_keys SET created_at='2026-10-30 01:02:03' WHERE group_id='G1'`)
	require.NoError(t, err)
	var value string
	require.NoError(t, db.QueryRow(`SELECT value FROM record_key_audit`).Scan(&value))
	require.Equal(t, "2026-10-30 01:02:03", value)
	require.NoError(t, db.QueryRow(`SELECT created_at FROM record_key_view WHERE group_id='G1'`).Scan(&value))
	require.Contains(t, value, "2026-10-30")
	var count int
	require.NoError(t, db.QueryRow(`SELECT count(*) FROM sqlite_master WHERE name='record_key_index'`).Scan(&count))
	require.Equal(t, 1, count)
	for _, query := range []string{
		`DELETE FROM write_group_profile_keys WHERE group_id='G1'`,
		`DELETE FROM write_groups WHERE id='profile-only'`,
		`INSERT INTO write_group_profile_keys(workspace_id,circuit_id,profile_kind,canonical_role,group_id) VALUES('workspace','circuit-P','cumulative','circuit','unkeyed')`,
		`INSERT INTO write_group_profile_keys(workspace_id,circuit_id,profile_kind,canonical_role,group_id) VALUES('workspace','other','cumulative','circuit','profile-only')`,
		`INSERT INTO write_group_profile_keys(workspace_id,circuit_id,profile_kind,canonical_role,group_id) VALUES('workspace','other','cumulative','circuit','missing')`,
	} {
		_, err = db.Exec(query)
		require.Error(t, err, query)
	}
}

func TestRecordProfileKeysFaultRollback(t *testing.T) {
	for _, stage := range []string{"copy", "replace", "indexes", "foreign_keys", "commit"} {
		for _, fk := range []int{0, 1} {
			t.Run(fmt.Sprintf("%s/fk%d", stage, fk), func(t *testing.T) {
				db := recordProfileDB(t)
				addRecordProfileObjects(t, db)
				_, err := db.Exec(fmt.Sprintf("PRAGMA foreign_keys=%d", fk))
				require.NoError(t, err)
				columns := lineageSnapshotColumns(t, db, "write_group_profile_keys", "write_group_lineages", "write_group_generations", "record_key_child")
				before := lineageMigrationSnapshot(t, db, columns)
				schema := recordProfileSchema(t, db)
				fault := errors.New("record migration injected failure")
				err = runSQLiteConfigurationExpansion(db, func(ctx context.Context, tx *sql.Tx) error {
					return rebuildRecordProfileKeys(ctx, tx, func(at string) error {
						if at == stage {
							return fault
						}
						return nil
					})
				})
				require.ErrorIs(t, err, fault)
				require.Equal(t, schema, recordProfileSchema(t, db))
				require.Equal(t, before, lineageMigrationSnapshot(t, db, columns))
				var actual int
				require.NoError(t, db.QueryRow("PRAGMA foreign_keys").Scan(&actual))
				require.Equal(t, fk, actual)
				require.NoError(t, db.QueryRow(`SELECT count(*) FROM sqlite_master WHERE name='write_group_profile_keys_record_expand'`).Scan(&actual))
				require.Zero(t, actual)
				require.NoError(t, runSQLiteConfigurationExpansion(db, migrateMixedRecordProfileKeysTx))
			})
		}
	}
}

func TestRecordProfileKeysRestartAndLineageCorruption(t *testing.T) {
	path := filepath.Join(t.TempDir(), "record.db")
	db := openLineageMigrationDB(t, path)
	createPre047LineageSchema(t, db)
	seedLegacyLineageRows(t, db)
	require.NoError(t, ensureSQLiteWriteGroupLineages(db))
	require.NoError(t, runSQLiteConfigurationExpansion(db, migrateMixedRecordProfileKeysTx))
	_, err := db.Exec(`INSERT INTO write_group_profile_keys(workspace_id,selection_scope,profile_kind,canonical_role,group_id,lineage_id) VALUES('workspace','12345678-1234-1234-1234-123456789abc','mixed','record','unkeyed','record-lineage');
 INSERT INTO write_group_lineages(lineage_id,workspace_id,current_group_id) VALUES('record-lineage','workspace','unkeyed');
 INSERT INTO write_group_generations(group_id,lineage_id,generation) VALUES('unkeyed','record-lineage',1);`)
	require.NoError(t, err)
	require.NoError(t, ensureSQLiteWriteGroupLineages(db))
	columns := lineageSnapshotColumns(t, db, "write_group_profile_keys", "write_group_lineages", "write_group_generations")
	before := lineageMigrationSnapshot(t, db, columns)
	require.NoError(t, db.Close())
	db = openLineageMigrationDB(t, path)
	require.NoError(t, ensureSQLiteWriteGroupLineages(db))
	require.NoError(t, runSQLiteConfigurationExpansion(db, migrateMixedRecordProfileKeysTx))
	require.Equal(t, before, lineageMigrationSnapshot(t, db, columns))
	_, err = db.Exec(`DELETE FROM write_group_generations WHERE group_id='unkeyed'`)
	require.NoError(t, err)
	corrupt := lineageMigrationSnapshot(t, db, columns)
	require.ErrorContains(t, ensureSQLiteWriteGroupLineages(db), "corrupt")
	require.Equal(t, corrupt, lineageMigrationSnapshot(t, db, columns), "initialized corruption is never repaired")
}

func TestRecordProfileKeyValidatorRejectsMalformedTupleWithoutChecks(t *testing.T) {
	db := recordProfileDB(t)
	require.NoError(t, runSQLiteConfigurationExpansion(db, migrateMixedRecordProfileKeysTx))
	_, err := db.Exec(`PRAGMA ignore_check_constraints=ON`)
	require.NoError(t, err)
	for _, scope := range []string{"bad", "00000000-0000-0000-0000-000000000000", strings.ToUpper("12345678-1234-1234-1234-123456789abc"), "123456781234-1234-1234-1234-123456789abc", "12345678-1234-1234-1234-123456789abc\x00"} {
		_, err = db.Exec(`UPDATE write_group_profile_keys SET profile_kind='mixed',canonical_role='record',circuit_id='',selection_scope=? WHERE group_id='G1'`, scope)
		require.NoError(t, err)
		require.ErrorContains(t, validateSQLiteLineageKeyScopes(context.Background(), db), "profile key shape")
	}
}
