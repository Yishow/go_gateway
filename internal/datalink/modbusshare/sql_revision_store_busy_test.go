package modbusshare

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShareRevisionCommitBusyLeavesOriginalRevisionAndCleanPool(t *testing.T) {
	for _, operation := range []string{"desired", "revision", "dirty"} {
		t.Run(operation, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "share.db")
			db, err := sql.Open("sqlite", path)
			require.NoError(t, err)
			defer db.Close()
			db.SetMaxOpenConns(1)
			_, err = db.ExecContext(t.Context(), `PRAGMA journal_mode=DELETE;PRAGMA busy_timeout=0;CREATE TABLE system_settings(key TEXT PRIMARY KEY,value TEXT,updated_at DATETIME)`)
			require.NoError(t, err)
			reader, err := sql.Open("sqlite", path)
			require.NoError(t, err)
			defer reader.Close()
			held, err := reader.BeginTx(t.Context(), nil)
			require.NoError(t, err)
			defer held.Rollback()
			var count int
			require.NoError(t, held.QueryRowContext(t.Context(), `SELECT count(*) FROM system_settings`).Scan(&count))
			store := NewSQLWorkspaceRevisionStore(db)
			switch operation {
			case "desired":
				err = store.CommitDesiredMappings(t.Context(), "ws", "rev-1", "rev-2", nil)
			case "dirty":
				err = store.MarkDirty(t.Context(), "ws")
			default:
				err = store.UpdateRevision(t.Context(), "ws", "rev-1", "rev-2")
			}
			require.Error(t, err)
			rev, _, err := store.GetRevision(t.Context(), "ws")
			require.NoError(t, err)
			require.Equal(t, "rev-1", rev)
			require.NoError(t, held.Rollback())
			next, err := db.BeginTx(t.Context(), nil)
			require.NoError(t, err)
			require.NoError(t, next.Rollback())
		})
	}
}
