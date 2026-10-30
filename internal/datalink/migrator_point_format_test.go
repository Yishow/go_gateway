package datalink

import (
	"database/sql"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema/migrations"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestMigratorPointFormatLegacyAndRepeat(t *testing.T) {
	path := filepath.Join(t.TempDir(), "legacy-config.sqlite")
	db, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	defer db.Close()
	db.SetMaxOpenConns(1)
	for _, file := range []string{"001_initial_schema_sqlite.sql", "003_align_points_unique_function_sqlite.up.sql"} {
		raw, err := migrations.FS.ReadFile(file)
		require.NoError(t, err)
		_, err = db.ExecContext(t.Context(), string(raw))
		require.NoError(t, err)
	}
	_, err = db.ExecContext(t.Context(), `ALTER TABLE points ADD COLUMN data_format TEXT;
 INSERT INTO devices(id,name,protocol,status,connection_config) VALUES('device','legacy device','mc_3e','active','{"data_format":" cdab "}');
 INSERT INTO points(id,device_id,name,address,data_type,data_format,mode) VALUES('old','device','old point','D0','uint32','ABCD','read');
 INSERT INTO points(id,device_id,name,address,data_type,data_format,mode) VALUES('untouched','device','untouched point','D2','uint32','invalid-old-format','read');`)
	require.NoError(t, err)
	require.NoError(t, NewMigrator().Migrate(db))
	var columns int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM pragma_table_info('points') WHERE name IN ('format_mode','format_override','format_version','revision_id')`).Scan(&columns))
	require.Equal(t, 4, columns)
	var mode, override sql.NullString
	var version int
	var revision, format, connection string
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT format_mode,format_override,format_version,revision_id,data_format FROM points WHERE id='old'`).Scan(&mode, &override, &version, &revision, &format))
	require.False(t, mode.Valid)
	require.False(t, override.Valid)
	require.Zero(t, version)
	require.Empty(t, revision)
	require.Equal(t, "ABCD", format)
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT connection_config FROM devices WHERE id='device'`).Scan(&connection))
	require.Equal(t, `{"data_format":" cdab "}`, connection)
	old, err := point.NewSQLRepository(db).GetByID(t.Context(), "old")
	require.NoError(t, err)
	require.Empty(t, old.FormatMode)
	require.Empty(t, old.FormatOverride)
	require.Zero(t, old.FormatVersion)
	_, err = db.ExecContext(t.Context(), `UPDATE points SET format_mode='override',format_override='ABCD',format_version=1,revision_id='explicit-revision' WHERE id='old'`)
	require.NoError(t, err)
	require.NoError(t, db.Close())
	db, err = sql.Open("sqlite", path)
	require.NoError(t, err)
	defer db.Close()
	db.SetMaxOpenConns(1)
	require.NoError(t, NewMigrator().Migrate(db))
	untouched, err := point.NewSQLRepository(db).GetByID(t.Context(), "untouched")
	require.NoError(t, err)
	require.Empty(t, untouched.FormatMode)
	require.Empty(t, untouched.FormatOverride)
	require.Zero(t, untouched.FormatVersion)
	require.Empty(t, untouched.RevisionID)
	require.Equal(t, "invalid-old-format", untouched.DataFormat)
	restored, err := point.NewSQLRepository(db).GetByID(t.Context(), "old")
	require.NoError(t, err)
	require.Equal(t, "override", string(restored.FormatMode))
	require.Equal(t, "ABCD", restored.FormatOverride)
	require.Equal(t, 1, restored.FormatVersion)
	require.Equal(t, "explicit-revision", restored.RevisionID)
	require.Equal(t, "ABCD", restored.DataFormat)
}

func TestMigratorPointFormatRejectsPartialSchema(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()
	raw, err := migrations.FS.ReadFile("001_initial_schema_sqlite.sql")
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), string(raw))
	require.NoError(t, err)
	// Align old uniqueness first so a legacy table rebuild cannot remove this marker.
	raw, err = migrations.FS.ReadFile("003_align_points_unique_function_sqlite.up.sql")
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), string(raw))
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), "ALTER TABLE points ADD COLUMN format_mode TEXT")
	require.NoError(t, err)
	err = NewMigrator().Migrate(db)
	require.ErrorContains(t, err, "point format migration is incomplete")
	var columns int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM pragma_table_info('points') WHERE name IN ('format_mode','format_override','format_version','revision_id')`).Scan(&columns))
	require.Equal(t, 1, columns, "partial metadata must not be silently repaired or inferred")
}
