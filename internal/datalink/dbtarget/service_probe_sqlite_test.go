package dbtarget

import (
	"database/sql"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestSQLiteConnectorCreateProbeDoesNotCreateMissingFile(t *testing.T) {
	target := filepath.Join(t.TempDir(), "missing.db")
	db := openMigratedTestDB(t)
	repo := NewSQLConnectorRepository(db)
	service := NewConnectorService(repo, NewSQLTargetMappingRepository(db))

	connector, err := service.Create(t.Context(), CreateConnectorRequest{
		Name: "missing sqlite",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: ConnectionConfig{
			"dsn": target,
		},
	})

	require.NoError(t, err)
	// A planned file that does not exist yet is usable: only the explicit schema
	// confirmation creates it, so it must not read as a lasting "unreachable".
	require.Equal(t, schema.DatabaseConnectorStatusReady, connector.Status)
	require.Empty(t, connector.LastCheckError)
	require.NoFileExists(t, target)
	saved, err := repo.GetByID(t.Context(), connector.ID)
	require.NoError(t, err)
	require.Equal(t, connector.Status, saved.Status)
}

func TestSQLiteConnectorUpdateProbeDoesNotCreateMissingFile(t *testing.T) {
	existing := filepath.Join(t.TempDir(), "existing.db")
	createProbeSQLiteFile(t, existing)
	missing := filepath.Join(t.TempDir(), "missing.db")
	db := openMigratedTestDB(t)
	repo := NewSQLConnectorRepository(db)
	service := NewConnectorService(repo, NewSQLTargetMappingRepository(db))

	connector, err := service.Create(t.Context(), CreateConnectorRequest{
		Name: "sqlite update",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: ConnectionConfig{
			"dsn": existing,
		},
	})
	require.NoError(t, err)

	expectedRevision := connector.IdentityRevision
	updated, err := service.Update(t.Context(), connector.ID, UpdateConnectorRequest{
		ExpectedIdentityRevision: &expectedRevision,
		ConnectionConfig: &ConnectionConfig{
			"dsn": missing,
		},
	})

	require.NoError(t, err)
	require.Equal(t, schema.DatabaseConnectorStatusReady, updated.Status)
	require.NoFileExists(t, missing)
	require.FileExists(t, existing)
}

func TestSQLiteProbeExistingFileUsesReadOnlyDSN(t *testing.T) {
	target := filepath.Join(t.TempDir(), "existing.db")
	createProbeSQLiteFile(t, target)
	before, err := os.Stat(target)
	require.NoError(t, err)

	originalOpen := openExternalDBManagerFunc
	var probedDSN string
	openExternalDBManagerFunc = func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
		probedDSN = stringConfigValue(config, "dsn")
		return originalOpen(kind, config)
	}
	t.Cleanup(func() { openExternalDBManagerFunc = originalOpen })

	status, _, message := probeConnector(t.Context(), schema.DatabaseConnectorKindSQLite, ConnectionConfig{"dsn": target})

	require.Equal(t, schema.DatabaseConnectorStatusReady, status)
	require.Empty(t, message)
	expectedDSN := (&url.URL{Scheme: "file", Path: target, RawQuery: "mode=ro"}).String()
	require.Equal(t, expectedDSN, probedDSN)
	after, err := os.Stat(target)
	require.NoError(t, err)
	require.Equal(t, before.Size(), after.Size())
	require.Equal(t, before.ModTime(), after.ModTime())
}

func TestSQLiteProbeExistingFileWithEscapedPlainPath(t *testing.T) {
	target := filepath.Join(t.TempDir(), "existing#name.db")
	createProbeSQLiteFile(t, target)

	originalOpen := openExternalDBManagerFunc
	var probedDSN string
	openExternalDBManagerFunc = func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
		probedDSN = stringConfigValue(config, "dsn")
		return originalOpen(kind, config)
	}
	t.Cleanup(func() { openExternalDBManagerFunc = originalOpen })

	status, _, message := probeConnector(t.Context(), schema.DatabaseConnectorKindSQLite, ConnectionConfig{"dsn": target})

	require.Equal(t, schema.DatabaseConnectorStatusReady, status)
	require.Empty(t, message)
	expected := (&url.URL{Scheme: "file", Path: target, RawQuery: "mode=ro"}).String()
	require.Equal(t, expected, probedDSN)
}

func TestSQLiteProbeRejectsInvalidModesWithoutCreatingFile(t *testing.T) {
	for _, testCase := range []struct {
		name   string
		suffix string
	}{
		{name: "unknown mode", suffix: "?mode=unsupported"},
		{name: "malformed query", suffix: "?%zz"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			target := filepath.Join(t.TempDir(), "invalid.db") + testCase.suffix
			status, _, message := probeConnector(t.Context(), schema.DatabaseConnectorKindSQLite, ConnectionConfig{"dsn": target})

			require.NotEqual(t, schema.DatabaseConnectorStatusReady, status)
			require.NotEmpty(t, message)
			require.NoFileExists(t, filepath.Join(filepath.Dir(target), "invalid.db"))
		})
	}
}

func TestSQLiteManagedDDLCreatesMissingFileAndTable(t *testing.T) {
	target := filepath.Join(t.TempDir(), "confirmed.db")
	fixture := newManagedDestinationFixture(t, schema.DatabaseConnectorKindSQLite, ConnectionConfig{"dsn": target})
	require.NoFileExists(t, target)

	execution, err := fixture.inspector.ExecuteSchemaStatements(
		t.Context(), fixture.connector.ID, fixture.connector.IdentityRevision, "main",
		[]string{`CREATE TABLE confirmed_record (id TEXT PRIMARY KEY)`},
	)

	require.NoError(t, err)
	require.Equal(t, 1, execution.Committed)
	require.FileExists(t, target)
	db, err := sql.Open("sqlite", target)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	require.True(t, sqliteTableExists(t, db, "confirmed_record"))
}

func createProbeSQLiteFile(t *testing.T, target string) {
	t.Helper()
	db, err := sql.Open("sqlite", target)
	require.NoError(t, err)
	require.NoError(t, db.PingContext(t.Context()))
	require.NoError(t, db.Close())
}

func TestSQLiteConnectorProbeMissingParentStaysUnreachableWithoutCreatingDirectories(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "absent-dir", "missing.db")
	db := openMigratedTestDB(t)
	service := NewConnectorService(NewSQLConnectorRepository(db), NewSQLTargetMappingRepository(db))

	connector, err := service.Create(t.Context(), CreateConnectorRequest{
		Name:             "missing parent",
		Kind:             schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: ConnectionConfig{"dsn": target},
	})

	require.NoError(t, err)
	require.Equal(t, schema.DatabaseConnectorStatusUnreachable, connector.Status)
	require.NotEmpty(t, connector.LastCheckError)
	require.NoDirExists(t, filepath.Dir(target))
	require.NoFileExists(t, target)
}
