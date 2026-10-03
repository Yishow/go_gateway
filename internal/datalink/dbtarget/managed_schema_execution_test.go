package dbtarget

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type schemaExecutionInterleaveRepository struct {
	ConnectorRepository
	afterRead func()
}

func (r *schemaExecutionInterleaveRepository) GetByID(ctx context.Context, id string) (*schema.DatabaseConnector, error) {
	connector, err := r.ConnectorRepository.GetByID(ctx, id)
	if err == nil && r.afterRead != nil {
		afterRead := r.afterRead
		r.afterRead = nil
		afterRead()
	}
	return connector, err
}

func TestManagedSchemaExecutionKeepsConfirmedTargetAfterConnectorInterleave(t *testing.T) {
	originalPath := filepath.Join(t.TempDir(), "confirmed.db")
	replacementPath := filepath.Join(t.TempDir(), "replacement.db")
	f := newManagedDestinationFixture(t, schema.DatabaseConnectorKindSQLite, ConnectionConfig{"dsn": originalPath})
	base := f.service.repo
	repo := &schemaExecutionInterleaveRepository{ConnectorRepository: base}
	f.service.repo = repo
	repo.afterRead = func() {
		changed := *f.connector
		changed.IdentityRevision = "replacement-revision"
		config, err := serializeConnectionConfig(ConnectionConfig{"dsn": replacementPath})
		require.NoError(t, err)
		changed.ConnectionConfig = config
		require.NoError(t, base.UpdateWithExpectedIdentityRevision(t.Context(), &changed, f.connector.IdentityRevision))
	}

	result, err := f.inspector.ExecuteSchemaStatements(t.Context(), f.connector.ID, f.connector.IdentityRevision, "main", []string{`CREATE TABLE confirmed_record (id TEXT PRIMARY KEY)`})
	require.NoError(t, err)
	require.Equal(t, 1, result.Committed)
	_, err = os.Stat(replacementPath)
	require.ErrorIs(t, err, os.ErrNotExist, "DDL must not use settings saved after the confirmed snapshot was resolved")
	db, err := sql.Open("sqlite", originalPath)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	require.True(t, sqliteTableExists(t, db, "confirmed_record"))
}

func TestManagedSchemaExecutionRejectsStaleOrDisabledBeforeOpen(t *testing.T) {
	for _, disabled := range []bool{false, true} {
		t.Run(map[bool]string{false: "stale", true: "disabled"}[disabled], func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "missing.db")
			f := newManagedDestinationFixture(t, schema.DatabaseConnectorKindSQLite, ConnectionConfig{"dsn": path})
			changed := *f.connector
			if disabled {
				changed.Enabled = false
			} else {
				changed.IdentityRevision = "new-revision"
			}
			require.NoError(t, f.service.repo.UpdateWithExpectedIdentityRevision(t.Context(), &changed, f.connector.IdentityRevision))
			_, err := f.inspector.ExecuteSchemaStatements(t.Context(), f.connector.ID, f.connector.IdentityRevision, "main", []string{`CREATE TABLE forbidden_record (id TEXT)`})
			if disabled {
				require.ErrorIs(t, err, ErrConnectorDisabled)
			} else {
				require.ErrorIs(t, err, ErrConnectorRevisionConflict)
			}
			_, err = os.Stat(path)
			require.ErrorIs(t, err, os.ErrNotExist)
		})
	}
}
