package datalink

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemberCadenceMigrationLegacyNullAndSemanticStorage(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()
	db.SetMaxOpenConns(1)
	require.NoError(t, NewMigrator().Migrate(db))
	require.NoError(t, NewMigrator().Migrate(db))
	var count int
	require.NoError(t, db.QueryRow(`SELECT count(*) FROM pragma_table_info('write_group_members') WHERE name='cadence' AND type='TEXT' AND "notnull"=0`).Scan(&count))
	require.Equal(t, 1, count)
	// No row backfill or fabricated interval is performed by the migration.
	var defaultValue sql.NullString
	require.NoError(t, db.QueryRow(`SELECT dflt_value FROM pragma_table_info('write_group_members') WHERE name='cadence'`).Scan(&defaultValue))
	require.False(t, defaultValue.Valid)
}

func TestMemberCadenceMigrationRollsBackWithUnknownKeySchema(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()
	db.SetMaxOpenConns(1)
	require.NoError(t, NewMigrator().Migrate(db))
	_, err = db.Exec(`ALTER TABLE write_group_members DROP COLUMN cadence; ALTER TABLE write_group_profile_keys ADD COLUMN unknown_extra TEXT; PRAGMA foreign_keys=ON`)
	require.NoError(t, err)
	var beforeTemp int
	require.NoError(t, db.QueryRow(`SELECT count(*) FROM pragma_database_list WHERE name='temp'`).Scan(&beforeTemp))
	require.ErrorContains(t, ensureSQLiteMemberCadence(db), "unknown write_group_profile_keys schema")
	var count, foreignKeys int
	require.NoError(t, db.QueryRow(`SELECT count(*) FROM pragma_table_info('write_group_members') WHERE name='cadence'`).Scan(&count))
	require.Zero(t, count)
	require.NoError(t, db.QueryRow(`PRAGMA foreign_keys`).Scan(&foreignKeys))
	require.Equal(t, 1, foreignKeys)
	var databases int
	require.NoError(t, db.QueryRow(`SELECT count(*) FROM pragma_database_list WHERE name='temp'`).Scan(&databases))
	require.Equal(t, beforeTemp, databases, "pre-existing temp schema must not be removed")
}

func TestMemberCadenceMigrationRejectsUnknownColumnShape(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()
	db.SetMaxOpenConns(1)
	require.NoError(t, NewMigrator().Migrate(db))
	_, err = db.Exec(`ALTER TABLE write_group_members DROP COLUMN cadence; ALTER TABLE write_group_members ADD COLUMN cadence INTEGER`)
	require.NoError(t, err)
	require.ErrorContains(t, ensureSQLiteMemberCadence(db), "unknown member cadence column")
}

func TestConfigurationExpansionBusyCommitKeepsConnectionReusable(t *testing.T) {
	path := t.TempDir() + "/busy-migration.db"
	db, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	defer db.Close()
	db.SetMaxOpenConns(1)
	require.NoError(t, NewMigrator().Migrate(db))
	_, err = db.Exec(`CREATE TABLE probe(value INTEGER); INSERT INTO probe VALUES(1); PRAGMA foreign_keys=ON`)
	require.NoError(t, err)
	readers, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	defer readers.Close()
	var reader *sql.Tx
	err = runSQLiteConfigurationExpansion(db, func(ctx context.Context, tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `UPDATE probe SET value=2`); err != nil {
			return err
		}
		var err error
		reader, err = readers.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		var value int
		return reader.QueryRowContext(ctx, `SELECT value FROM probe`).Scan(&value)
	})
	require.ErrorContains(t, err, "SQLITE_BUSY")
	require.NotNil(t, reader)
	require.NoError(t, reader.Rollback())
	require.NoError(t, runSQLiteConfigurationExpansion(db, func(ctx context.Context, tx *sql.Tx) error {
		var value int
		if err := tx.QueryRowContext(ctx, `SELECT value FROM probe`).Scan(&value); err != nil {
			return err
		}
		require.Equal(t, 1, value)
		return nil
	}))
	var foreignKeys int
	require.NoError(t, db.QueryRow(`PRAGMA foreign_keys`).Scan(&foreignKeys))
	require.Equal(t, 1, foreignKeys)
}
