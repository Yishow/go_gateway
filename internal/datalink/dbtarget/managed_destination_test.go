package dbtarget

import (
	"database/sql"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type managedDestinationFixture struct {
	inspector    *ManagedTableInspector
	service      *ConnectorService
	connector    *schema.DatabaseConnector
	internalDB   *sql.DB
	internalPath string
}

func newManagedDestinationFixture(t *testing.T, kind schema.DatabaseConnectorKind, config ConnectionConfig) managedDestinationFixture {
	t.Helper()
	internalDB := openMigratedTestDB(t)
	datalinkbase.ApplySQLitePoolDefaults(internalDB)
	var sequence int
	var databaseName, internalPath string
	require.NoError(t, internalDB.QueryRowContext(t.Context(), `PRAGMA database_list`).Scan(&sequence, &databaseName, &internalPath))

	repo := NewSQLConnectorRepository(internalDB)
	serialized, err := serializeConnectionConfig(config)
	require.NoError(t, err)
	connector := &schema.DatabaseConnector{
		ID:               "managed-destination",
		Name:             "managed destination",
		Kind:             kind,
		ConnectionConfig: serialized,
		IdentityRevision: "managed-revision",
		Enabled:          true,
	}
	require.NoError(t, repo.Create(t.Context(), connector))
	service := NewConnectorService(repo)
	return managedDestinationFixture{
		inspector: NewManagedTableInspector(service, internalDB),
		service:   service, connector: connector, internalDB: internalDB, internalPath: internalPath,
	}
}

func TestManagedDestinationSQLiteMissingFileReturnsMissingWithoutCreate(t *testing.T) {
	target := filepath.Join(t.TempDir(), "missing.db")
	fixture := newManagedDestinationFixture(t, schema.DatabaseConnectorKindSQLite, ConnectionConfig{"dsn": target})

	result, err := fixture.inspector.InspectTable(t.Context(), fixture.connector.ID, "", "gw_record_samples")

	require.NoError(t, err)
	require.Equal(t, TableInspectionMissing, result.Status)
	require.Equal(t, "main", result.Schema)
	require.Equal(t, "gw_record_samples", result.Table)
	require.Empty(t, result.Columns)
	_, statErr := os.Stat(target)
	require.ErrorIs(t, statErr, os.ErrNotExist)
}

func TestManagedDestinationSQLiteReadsExistingFileReadOnly(t *testing.T) {
	target := newInspectionSQLiteTarget(t, `CREATE TABLE gw_record_samples (recorded_at TEXT NOT NULL, value REAL)`)
	fixture := newManagedDestinationFixture(t, schema.DatabaseConnectorKindSQLite, ConnectionConfig{"dsn": target})

	result, err := fixture.inspector.InspectTable(t.Context(), fixture.connector.ID, "main", "gw_record_samples")

	require.NoError(t, err)
	require.Equal(t, TableInspectionExists, result.Status)
	require.Equal(t, []string{"recorded_at", "value"}, []string{result.Columns[0].Name, result.Columns[1].Name})
}

func TestManagedDestinationSQLiteFileOpaqueURIReadsDecodedPath(t *testing.T) {
	target := newInspectionSQLiteTarget(t, `CREATE TABLE gw_record_samples (value REAL)`)
	fixture := newManagedDestinationFixture(t, schema.DatabaseConnectorKindSQLite, ConnectionConfig{
		"dsn": (&url.URL{Scheme: "file", Opaque: url.PathEscape(target)}).String(),
	})

	result, err := fixture.inspector.InspectTable(t.Context(), fixture.connector.ID, "main", "gw_record_samples")

	require.NoError(t, err)
	require.Equal(t, TableInspectionExists, result.Status)
}

func TestManagedDestinationSQLiteAllowsExternalSymlinkFile(t *testing.T) {
	target := newInspectionSQLiteTarget(t, `CREATE TABLE gw_record_samples (value REAL)`)
	alias := filepath.Join(t.TempDir(), "target-alias.db")
	require.NoError(t, os.Symlink(target, alias))
	fixture := newManagedDestinationFixture(t, schema.DatabaseConnectorKindSQLite, ConnectionConfig{"dsn": alias})

	result, err := fixture.inspector.InspectTable(t.Context(), fixture.connector.ID, "main", "gw_record_samples")

	require.NoError(t, err)
	require.Equal(t, TableInspectionExists, result.Status)
}

func TestManagedDestinationSQLiteRejectsInternalFileAliases(t *testing.T) {
	tests := []struct {
		name  string
		alias func(t *testing.T, internalPath string) string
	}{
		{name: "same path", alias: func(_ *testing.T, internalPath string) string { return internalPath }},
		{name: "relative path", alias: func(t *testing.T, internalPath string) string {
			cwd, err := os.Getwd()
			require.NoError(t, err)
			require.NoError(t, os.Chdir(filepath.Dir(internalPath)))
			t.Cleanup(func() { require.NoError(t, os.Chdir(cwd)) })
			return filepath.Base(internalPath)
		}},
		{name: "file opaque uri", alias: func(_ *testing.T, internalPath string) string {
			return (&url.URL{Scheme: "file", Opaque: url.PathEscape(internalPath)}).String()
		}},
		{name: "symlink", alias: func(t *testing.T, internalPath string) string {
			alias := filepath.Join(t.TempDir(), "internal-alias.db")
			require.NoError(t, os.Symlink(internalPath, alias))
			return alias
		}},
		{name: "hardlink", alias: func(t *testing.T, internalPath string) string {
			alias := filepath.Join(t.TempDir(), "internal-hardlink.db")
			require.NoError(t, os.Link(internalPath, alias))
			return alias
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture := newManagedDestinationFixture(t, schema.DatabaseConnectorKindSQLite, nil)
			alias := test.alias(t, fixture.internalPath)
			fixture.connector.ConnectionConfig, _ = serializeConnectionConfig(ConnectionConfig{"dsn": alias})
			require.NoError(t, fixture.service.repo.Update(t.Context(), fixture.connector))

			_, err := fixture.inspector.ValidateDestination(t.Context(), fixture.connector.ID, fixture.connector.IdentityRevision)

			require.ErrorIs(t, err, ErrManagedDestinationInternal)
			require.NotContains(t, err.Error(), fixture.internalPath)
		})
	}
}

func TestManagedDestinationSQLiteRejectsAttachedInternalFile(t *testing.T) {
	fixture := newManagedDestinationFixture(t, schema.DatabaseConnectorKindSQLite, nil)
	attachedPath := filepath.Join(t.TempDir(), "journal.db")
	_, err := fixture.internalDB.ExecContext(t.Context(), `ATTACH DATABASE ? AS journal`, attachedPath)
	require.NoError(t, err)
	fixture.connector.ConnectionConfig, err = serializeConnectionConfig(ConnectionConfig{"dsn": attachedPath})
	require.NoError(t, err)
	require.NoError(t, fixture.service.repo.Update(t.Context(), fixture.connector))

	_, err = fixture.inspector.ValidateDestination(t.Context(), fixture.connector.ID, fixture.connector.IdentityRevision)

	require.ErrorIs(t, err, ErrManagedDestinationInternal)
	require.NotContains(t, err.Error(), attachedPath)
}

func TestManagedDestinationSQLiteRejectsInvalidDestinations(t *testing.T) {
	tests := []struct {
		name   string
		config func(t *testing.T) ConnectionConfig
	}{
		{name: "memory", config: func(*testing.T) ConnectionConfig { return ConnectionConfig{"dsn": ":memory:"} }},
		{name: "host uri", config: func(t *testing.T) ConnectionConfig {
			return ConnectionConfig{"dsn": "file://remote" + filepath.ToSlash(filepath.Join(t.TempDir(), "target.db"))}
		}},
		{name: "unsafe query", config: func(t *testing.T) ConnectionConfig {
			return ConnectionConfig{"dsn": filepath.Join(t.TempDir(), "target.db") + "?_pragma=journal_mode(WAL)"}
		}},
		{name: "directory", config: func(t *testing.T) ConnectionConfig { return ConnectionConfig{"dsn": t.TempDir()} }},
		{name: "symlink directory", config: func(t *testing.T) ConnectionConfig {
			alias := filepath.Join(t.TempDir(), "directory-alias")
			require.NoError(t, os.Symlink(t.TempDir(), alias))
			return ConnectionConfig{"dsn": alias}
		}},
		{name: "missing parent", config: func(t *testing.T) ConnectionConfig {
			return ConnectionConfig{"dsn": filepath.Join(t.TempDir(), "missing", "target.db")}
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture := newManagedDestinationFixture(t, schema.DatabaseConnectorKindSQLite, test.config(t))
			_, err := fixture.inspector.ValidateDestination(t.Context(), fixture.connector.ID, fixture.connector.IdentityRevision)

			require.ErrorIs(t, err, ErrManagedDestinationInvalid)
			require.NotContains(t, err.Error(), "password")
		})
	}
}

func TestManagedDestinationSQLiteFailsClosedWithoutInternalDB(t *testing.T) {
	target := filepath.Join(t.TempDir(), "missing.db")
	fixture := newManagedDestinationFixture(t, schema.DatabaseConnectorKindSQLite, ConnectionConfig{"dsn": target})
	inspector := NewManagedTableInspector(fixture.service, nil)

	_, err := inspector.ValidateDestination(t.Context(), fixture.connector.ID, fixture.connector.IdentityRevision)

	require.ErrorIs(t, err, ErrManagedDestinationInvalid)
	_, statErr := os.Stat(target)
	require.ErrorIs(t, statErr, os.ErrNotExist)
}

func TestManagedDestinationSQLiteMissingFileRequiresWritableParent(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root can bypass directory permission checks")
	}
	parent := t.TempDir()
	require.NoError(t, os.Chmod(parent, 0o500))
	t.Cleanup(func() { _ = os.Chmod(parent, 0o700) })
	fixture := newManagedDestinationFixture(t, schema.DatabaseConnectorKindSQLite, ConnectionConfig{
		"dsn": filepath.Join(parent, "missing.db"),
	})

	_, err := fixture.inspector.ValidateDestination(t.Context(), fixture.connector.ID, fixture.connector.IdentityRevision)

	require.ErrorIs(t, err, ErrManagedDestinationInvalid)
}

func TestManagedDestinationPostgresUsesSavedConnectorAndSchema(t *testing.T) {
	stubInspectionDatabase(t, nil, postgresCatalogSetup(t.TempDir())...)
	fixture := newManagedDestinationFixture(t, schema.DatabaseConnectorKindPostgres, ConnectionConfig{
		"host": "db.internal", "user": "reader", "database": "metrics", "schema": "public",
	})

	resolved, err := fixture.inspector.ValidateDestination(t.Context(), fixture.connector.ID, fixture.connector.IdentityRevision)
	require.NoError(t, err)
	require.Equal(t, fixture.connector.ID, resolved.ID)

	result, err := fixture.inspector.InspectTable(t.Context(), fixture.connector.ID, "public", "sensor_values")
	require.NoError(t, err)
	require.Equal(t, TableInspectionMissing, result.Status)
	require.Equal(t, "public", result.Schema)
}

func TestManagedDestinationValidateDestinationPreservesSavedRevisionScope(t *testing.T) {
	fixture := newManagedDestinationFixture(t, schema.DatabaseConnectorKindPostgres, ConnectionConfig{
		"host": "db.internal", "user": "reader", "database": "metrics",
	})

	_, err := fixture.inspector.ValidateDestination(t.Context(), fixture.connector.ID, "stale")
	require.ErrorIs(t, err, ErrConnectorRevisionConflict)

	fixture.connector.Enabled = false
	require.NoError(t, fixture.service.repo.Update(t.Context(), fixture.connector))
	_, err = fixture.inspector.ValidateDestination(t.Context(), fixture.connector.ID, fixture.connector.IdentityRevision)
	require.ErrorIs(t, err, ErrConnectorDisabled)
}

func TestManagedDestinationInternalDBRowsAreClosed(t *testing.T) {
	target := newInspectionSQLiteTarget(t, `CREATE TABLE gw_record_samples (value REAL)`)
	fixture := newManagedDestinationFixture(t, schema.DatabaseConnectorKindSQLite, ConnectionConfig{"dsn": target})
	fixture.internalDB.SetMaxOpenConns(1)

	_, err := fixture.inspector.ValidateDestination(t.Context(), fixture.connector.ID, fixture.connector.IdentityRevision)
	require.NoError(t, err)

	var version int
	require.NoError(t, fixture.internalDB.QueryRowContext(t.Context(), `PRAGMA user_version`).Scan(&version))
}

func TestManagedDestinationErrorsDoNotExposePrivatePaths(t *testing.T) {
	target := filepath.Join(t.TempDir(), "missing.db")
	fixture := newManagedDestinationFixture(t, schema.DatabaseConnectorKindSQLite, ConnectionConfig{"dsn": target})
	fixture.connector.ConnectionConfig = "{\"dsn\":\"" + target + "\",\"password\":\"secret\""
	require.NoError(t, fixture.service.repo.Update(t.Context(), fixture.connector))

	_, err := fixture.inspector.ValidateDestination(t.Context(), fixture.connector.ID, fixture.connector.IdentityRevision)

	require.Error(t, err)
	require.True(t, errors.Is(err, ErrManagedDestinationInvalid) || errors.Is(err, ErrConnectorRevisionConflict))
	require.NotContains(t, err.Error(), target)
	require.NotContains(t, err.Error(), "secret")
	require.NotContains(t, err.Error(), "password")
}
