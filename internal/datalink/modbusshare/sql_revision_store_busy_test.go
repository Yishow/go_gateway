package modbusshare

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestShareRevisionCommitBusyLeavesOriginalRevisionAndCleanPool(t *testing.T) {
	for _, operation := range []string{"desired", "revision", "dirty"} {
		t.Run(operation, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "share.db")
			db, err := sql.Open("sqlite", path)
			require.NoError(t, err)
			defer db.Close()
			db.SetMaxOpenConns(1)
			_, err = db.Exec(`PRAGMA journal_mode=DELETE;PRAGMA busy_timeout=0;CREATE TABLE system_settings(key TEXT PRIMARY KEY,value TEXT,updated_at DATETIME)`)
			require.NoError(t, err)
			reader, err := sql.Open("sqlite", path)
			require.NoError(t, err)
			defer reader.Close()
			held, err := reader.Begin()
			require.NoError(t, err)
			defer held.Rollback()
			var count int
			require.NoError(t, held.QueryRow(`SELECT count(*) FROM system_settings`).Scan(&count))
			store := NewSQLWorkspaceRevisionStore(db)
			if operation == "desired" {
				err = store.CommitDesiredMappings(t.Context(), "ws", "rev-1", "rev-2", nil)
			} else if operation == "dirty" {
				err = store.MarkDirty(t.Context(), "ws")
			} else {
				err = store.UpdateRevision(t.Context(), "ws", "rev-1", "rev-2")
			}
			require.Error(t, err)
			rev, _, err := store.GetRevision(t.Context(), "ws")
			require.NoError(t, err)
			require.Equal(t, "rev-1", rev)
			require.NoError(t, held.Rollback())
			next, err := db.Begin()
			require.NoError(t, err)
			require.NoError(t, next.Rollback())
		})
	}
}
