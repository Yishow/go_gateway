package dbtarget

import (
	"os"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestManagedDestinationSQLiteRejectsNonMainBeforeInspect(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.db")
	f := newManagedDestinationFixture(t, schema.DatabaseConnectorKindSQLite, ConnectionConfig{"dsn": path})
	_, err := f.inspector.InspectTable(t.Context(), f.connector.ID, "other", "readings")
	require.ErrorIs(t, err, ErrManagedDestinationInvalid)
	require.NoFileExists(t, path)
}

func TestManagedInspectionKeepsRevisionBoundSnapshotDuringConnectorUpdate(t *testing.T) {
	original := newInspectionSQLiteTarget(t, `CREATE TABLE readings (original_value TEXT)`)
	replacement := newInspectionSQLiteTarget(t, `CREATE TABLE readings (replacement_value TEXT)`)
	f := newManagedDestinationFixture(t, schema.DatabaseConnectorKindSQLite, ConnectionConfig{"dsn": original})
	base := f.service.repo
	repo := &schemaExecutionInterleaveRepository{ConnectorRepository: base}
	f.service.repo = repo
	repo.afterRead = func() {
		changed := *f.connector
		changed.IdentityRevision = "replacement-revision"
		config, err := serializeConnectionConfig(ConnectionConfig{"dsn": replacement})
		require.NoError(t, err)
		changed.ConnectionConfig = config
		require.NoError(t, base.UpdateWithExpectedIdentityRevision(t.Context(), &changed, f.connector.IdentityRevision))
	}
	result, err := f.inspector.InspectTableAtRevision(t.Context(), f.connector.ID, f.connector.IdentityRevision, "main", "readings")
	require.NoError(t, err)
	require.Equal(t, TableInspectionExists, result.Status)
	require.Len(t, result.Columns, 1)
	require.Equal(t, "original_value", result.Columns[0].Name)
	_, err = os.Stat(original)
	require.NoError(t, err)
}
