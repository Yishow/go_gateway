package audit

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	datalinkbase "go-gateway/internal/datalink"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestSQLRepositoryPersistsWorkspaceAuditHistoryAcrossRestart(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "audit.db")
	db := openAuditTestDB(t, dbPath)

	occurredAt := time.Date(2026, 3, 16, 10, 30, 0, 0, time.UTC)
	svc := NewService(NewSQLRepository(db))
	require.NoError(t, svc.Record(context.Background(), RecordEvent{
		WorkspaceID: "ws-1",
		EventType:   EventTypeWorkspaceActivation,
		Result:      ResultPartialSuccess,
		Scope:       "devices:dev-A,dev-B",
		ReferenceID: "activation-1",
		Details: map[string]any{
			"success": []string{"dev-A"},
			"failed":  []string{"dev-B"},
		},
		OccurredAt: occurredAt,
	}))
	require.NoError(t, db.Close())

	restartedDB := openAuditTestDB(t, dbPath)
	restartedSvc := NewService(NewSQLRepository(restartedDB))
	entries, err := restartedSvc.List(context.Background(), ListFilter{WorkspaceID: "ws-1", Limit: 10})
	require.NoError(t, err)
	require.Len(t, entries, 1)

	entry := entries[0]
	require.Equal(t, "ws-1", entry.WorkspaceID)
	require.Equal(t, EventTypeWorkspaceActivation, entry.EventType)
	require.Equal(t, ResultPartialSuccess, entry.Result)
	require.Equal(t, "devices:dev-A,dev-B", entry.Scope)
	require.Equal(t, "activation-1", entry.ReferenceID)
	require.Equal(t, occurredAt, entry.OccurredAt)
	require.Contains(t, entry.Details, "dev-B")
}

func openAuditTestDB(t *testing.T, dbPath string) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	require.NoError(t, datalinkbase.NewMigrator().Migrate(db))
	return db
}
