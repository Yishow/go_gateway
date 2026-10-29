package datalink

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteGroupLineagesMigrationPartialArtifactsFailClosed(t *testing.T) {
	artifacts := []string{
		`ALTER TABLE write_group_basic_keys ADD COLUMN lineage_id TEXT`,
		`ALTER TABLE write_group_profile_keys ADD COLUMN lineage_id TEXT`,
		`CREATE TABLE write_group_lineages(marker TEXT)`,
		`CREATE TABLE write_group_generations(marker TEXT)`,
	}
	for mask := 1; mask < 15; mask++ {
		t.Run(fmt.Sprintf("artifacts-%04b", mask), func(t *testing.T) {
			db := openLineageMigrationDB(t, filepath.Join(t.TempDir(), "partial.db"))
			createPre047LineageSchema(t, db)
			for i, ddl := range artifacts {
				if mask&(1<<i) != 0 {
					_, err := db.Exec(ddl)
					require.NoError(t, err)
				}
			}
			before := lineageSchemaSnapshot(t, db)
			require.ErrorContains(t, ensureSQLiteWriteGroupsTables(db), "lineage migration is incomplete")
			require.Equal(t, before, lineageSchemaSnapshot(t, db), "partial initialization must not repair even an empty database")
		})
	}
	for _, table := range []string{"write_group_basic_keys", "write_group_profile_keys"} {
		t.Run("missing-"+table, func(t *testing.T) {
			db := openLineageMigrationDB(t, filepath.Join(t.TempDir(), "missing.db"))
			require.NoError(t, NewMigrator().Migrate(db))
			_, err := db.Exec(`DROP TABLE ` + table)
			require.NoError(t, err)
			before := lineageSchemaSnapshot(t, db)
			require.ErrorContains(t, NewMigrator().Migrate(db), "lineage migration is incomplete")
			require.Equal(t, before, lineageSchemaSnapshot(t, db), "preflight must precede old key-table recreation")
		})
	}
}

func TestWriteGroupLineagesMigrationRejectsInvalidLegacyAnchors(t *testing.T) {
	for _, corrupt := range []string{
		`UPDATE write_group_profile_keys SET workspace_id='foreign' WHERE group_id='G1'`,
		`UPDATE write_group_basic_keys SET group_id='missing' WHERE group_id='G1'`,
		`UPDATE write_group_basic_keys SET device_id='   ' WHERE group_id='G1'`,
		`UPDATE write_group_basic_keys SET canonical_role='unknown' WHERE group_id='G1'`,
	} {
		t.Run(corrupt, func(t *testing.T) {
			db := openLineageMigrationDB(t, filepath.Join(t.TempDir(), "legacy.db"))
			createPre047LineageSchema(t, db)
			seedLegacyLineageRows(t, db)
			_, err := db.Exec(corrupt)
			require.NoError(t, err)
			columns := lineageSnapshotColumns(t, db, lineagePreservedTables...)
			before := lineageMigrationSnapshot(t, db, columns)
			require.ErrorContains(t, ensureSQLiteWriteGroupLineages(db), "lineage metadata is corrupt")
			require.Equal(t, before, lineageMigrationSnapshot(t, db, columns))
			initialized, err := sqliteWriteGroupLineageState(t.Context(), db)
			require.NoError(t, err)
			require.False(t, initialized, "bad anchors must leave all expansion artifacts absent")
		})
	}
}

func TestWriteGroupLineagesMigrationMissingRegistriesPreservesKeyMarkers(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing-registries.db")
	db := upgradedLineageFixture(t, path)
	columns := lineageSnapshotColumns(t, db, lineagePreservedTables...)
	before := lineageMigrationSnapshot(t, db, columns)
	_, err := db.Exec(`DROP TABLE write_group_generations; DROP TABLE write_group_lineages`)
	require.NoError(t, err)
	schema := lineageSchemaSnapshot(t, db)
	require.ErrorContains(t, NewMigrator().Migrate(db), "lineage migration is incomplete")
	require.Equal(t, before, lineageMigrationSnapshot(t, db, columns))
	require.Equal(t, schema, lineageSchemaSnapshot(t, db))
	require.NoError(t, db.Close())
	db = openLineageMigrationDB(t, path)
	require.ErrorContains(t, NewMigrator().Migrate(db), "lineage migration is incomplete")
	require.Equal(t, before, lineageMigrationSnapshot(t, db, columns), "nonnull durable key references must never be rebound")
	require.Equal(t, schema, lineageSchemaSnapshot(t, db))
}

func TestWriteGroupLineagesMigrationFutureGenerationsAndReservationSurvive(t *testing.T) {
	path := filepath.Join(t.TempDir(), "future.db")
	db := upgradedLineageFixture(t, path)
	seedFutureLineageRows(t, db)
	columns := lineageSnapshotColumns(t, db, append(lineagePreservedTables, "write_group_lineages", "write_group_generations")...)
	before := lineageMigrationSnapshot(t, db, columns)
	require.NoError(t, NewMigrator().Migrate(db))
	require.Equal(t, before, lineageMigrationSnapshot(t, db, columns))
	require.NoError(t, db.Close())
	db = openLineageMigrationDB(t, path)
	require.NoError(t, NewMigrator().Migrate(db))
	require.Equal(t, before, lineageMigrationSnapshot(t, db, columns))
	var anchor, current, reserved string
	var revision int
	require.NoError(t, db.QueryRow(`SELECT k.group_id,l.current_group_id,l.reserved_group_id,l.revision
		FROM write_group_basic_keys k JOIN write_group_lineages l USING(lineage_id) WHERE k.group_id='G1'`).Scan(&anchor, &current, &reserved, &revision))
	require.Equal(t, "G1", anchor)
	require.Equal(t, "G2", current)
	require.Equal(t, "G3", reserved)
	require.Equal(t, 7, revision)
}

func TestWriteGroupLineagesMigrationCorruptMetadataNeverRepaired(t *testing.T) {
	cases := map[string]string{
		"null-reference":             `UPDATE write_group_basic_keys SET lineage_id=NULL WHERE group_id='G1'`,
		"empty-reference":            `UPDATE write_group_basic_keys SET lineage_id='' WHERE group_id='G1'`,
		"missing-lineage":            `DELETE FROM write_group_lineages WHERE current_group_id='G2'`,
		"missing-genesis":            `DELETE FROM write_group_generations WHERE group_id='G1'`,
		"missing-current-generation": `DELETE FROM write_group_generations WHERE group_id='G2'`,
		"missing-current-group":      `DELETE FROM write_groups WHERE id='G2'`,
		"foreign-generation-scope":   `UPDATE write_groups SET workspace_id='foreign' WHERE id='G3'`,
		"foreign-key-scope":          `UPDATE write_group_basic_keys SET workspace_id='foreign' WHERE group_id='G1'`,
		"foreign-lineage-scope":      `UPDATE write_group_lineages SET workspace_id='foreign' WHERE current_group_id='G2'`,
		"contradictory-alias":        `UPDATE write_group_profile_keys SET lineage_id=(SELECT lineage_id FROM write_group_generations WHERE group_id='profile-only') WHERE group_id='G1'`,
		"non-genesis-anchor":         `UPDATE write_group_basic_keys SET group_id='G2' WHERE group_id='G1'`,
		"foreign-current":            `UPDATE write_group_lineages SET current_group_id='G1' WHERE current_group_id='basic-only'`,
		"foreign-reservation":        `UPDATE write_group_lineages SET reserved_group_id='basic-only' WHERE current_group_id='G2'`,
		"current-reservation-equal":  `UPDATE write_group_lineages SET reserved_group_id='G2' WHERE current_group_id='G2'`,
		"partial-reservation":        `UPDATE write_group_lineages SET reservation_operation_id=NULL WHERE current_group_id='G2'`,
		"blank-reservation-intent":   `UPDATE write_group_lineages SET reservation_intent_digest=' ' WHERE current_group_id='G2'`,
		"orphan-operation":           `UPDATE write_group_lineages SET reservation_operation_id='O2' WHERE current_group_id='basic-only'`,
		"foreign-predecessor":        `UPDATE write_group_generations SET predecessor_group_id='basic-only' WHERE group_id='G2'`,
		"skipped-predecessor":        `UPDATE write_group_generations SET predecessor_group_id='G1' WHERE group_id='G3'`,
		"genesis-predecessor":        `UPDATE write_group_generations SET predecessor_group_id='G2' WHERE group_id='G1'`,
		"missing-predecessor":        `UPDATE write_group_generations SET predecessor_group_id=NULL WHERE group_id='G2'`,
		"fractional-generation":      `UPDATE write_group_generations SET generation=2.5 WHERE group_id='G2'`,
		"fractional-revision":        `UPDATE write_group_lineages SET revision=1.5 WHERE current_group_id='G2'`,
		"zero-revision":              `UPDATE write_group_lineages SET revision=0 WHERE current_group_id='G2'`,
		"orphan-lineage":             `DELETE FROM write_group_basic_keys WHERE group_id='basic-only'`,
	}
	for name, corrupt := range cases {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "corrupt.db")
			db := upgradedLineageFixture(t, path)
			seedFutureLineageRows(t, db)
			_, err := db.Exec(`PRAGMA foreign_keys=OFF; PRAGMA ignore_check_constraints=ON; ` + corrupt)
			require.NoError(t, err)
			columns := lineageSnapshotColumns(t, db, append(lineagePreservedTables, "write_group_lineages", "write_group_generations")...)
			before := lineageMigrationSnapshot(t, db, columns)
			require.ErrorContains(t, NewMigrator().Migrate(db), "lineage metadata is corrupt")
			require.Equal(t, before, lineageMigrationSnapshot(t, db, columns))
			require.NoError(t, db.Close())
			db = openLineageMigrationDB(t, path)
			require.ErrorContains(t, NewMigrator().Migrate(db), "lineage metadata is corrupt")
			require.Equal(t, before, lineageMigrationSnapshot(t, db, columns), "reopen must never reseed or repair")
		})
	}
}

func TestWriteGroupLineagesMigrationConstraintsProtectReferences(t *testing.T) {
	db := upgradedLineageFixture(t, filepath.Join(t.TempDir(), "constraints.db"))
	seedFutureLineageRows(t, db)
	_, err := db.Exec(`PRAGMA foreign_keys=ON`)
	require.NoError(t, err)
	for _, statement := range []string{
		`UPDATE write_group_basic_keys SET lineage_id='missing' WHERE group_id='G1'`,
		`UPDATE write_group_profile_keys SET lineage_id='missing' WHERE group_id='G1'`,
		`UPDATE write_group_basic_keys SET lineage_id='' WHERE group_id='G1'`,
		`UPDATE write_group_lineages SET revision=0 WHERE current_group_id='G2'`,
		`UPDATE write_group_lineages SET revision=1.5 WHERE current_group_id='G2'`,
		`UPDATE write_group_generations SET generation=2.5 WHERE group_id='G2'`,
		`UPDATE write_group_lineages SET reservation_operation_id=NULL WHERE current_group_id='G2'`,
		`UPDATE write_group_lineages SET reserved_group_id='G2' WHERE current_group_id='G2'`,
		`UPDATE write_group_generations SET predecessor_group_id=NULL WHERE group_id='G2'`,
		`UPDATE write_group_generations SET predecessor_group_id='G1' WHERE group_id='G1'`,
		`UPDATE write_group_generations SET generation=2 WHERE group_id='G3'`,
		`DELETE FROM write_group_lineages WHERE current_group_id='G2'`,
		`DELETE FROM write_group_generations WHERE group_id='G1'`,
		`DELETE FROM write_groups WHERE id='G1'`,
	} {
		_, err := db.Exec(statement)
		require.Error(t, err, statement)
	}
	for _, table := range []string{"write_group_basic_keys", "write_group_profile_keys"} {
		var nullable, references int
		require.NoError(t, db.QueryRow(`SELECT "notnull" FROM pragma_table_info(?) WHERE name='lineage_id'`, table).Scan(&nullable))
		require.Zero(t, nullable, "expand reference is nullable only during migration")
		require.NoError(t, db.QueryRow(`SELECT count(*) FROM pragma_foreign_key_list(?)
			WHERE "from" IN ('group_id','lineage_id') AND on_delete='RESTRICT'`, table).Scan(&references))
		require.Equal(t, 2, references, "original group FK and new lineage FK must both be retained")
	}
	require.NoError(t, NewMigrator().Migrate(db))
}

func upgradedLineageFixture(t *testing.T, path string) *sql.DB {
	t.Helper()
	db := openLineageMigrationDB(t, path)
	createPre047LineageSchema(t, db)
	seedLegacyLineageRows(t, db)
	require.NoError(t, NewMigrator().Migrate(db))
	return db
}

func seedFutureLineageRows(t *testing.T, db *sql.DB) {
	t.Helper()
	insertLineageMigrationGroup(t, db, "G2")
	insertLineageMigrationGroup(t, db, "G3")
	_, err := db.Exec(`INSERT INTO write_group_generations(group_id,lineage_id,generation,predecessor_group_id)
		SELECT 'G2',lineage_id,2,'G1' FROM write_group_generations WHERE group_id='G1';
		INSERT INTO write_group_generations(group_id,lineage_id,generation,predecessor_group_id)
		SELECT 'G3',lineage_id,3,'G2' FROM write_group_generations WHERE group_id='G1';
		UPDATE write_group_lineages SET current_group_id='G2',reserved_group_id='G3',reservation_operation_id='O3',reservation_intent_digest=' frozen intent bytes ',revision=7 WHERE current_group_id='G1';`)
	require.NoError(t, err)
}

func lineageSchemaSnapshot(t *testing.T, db *sql.DB) []string {
	t.Helper()
	rows, err := db.Query(`SELECT type||':'||name||':'||COALESCE(sql,'') FROM sqlite_master ORDER BY type,name`)
	require.NoError(t, err)
	defer rows.Close()
	var result []string
	for rows.Next() {
		var value string
		require.NoError(t, rows.Scan(&value))
		result = append(result, value)
	}
	require.NoError(t, rows.Err())
	return result
}
