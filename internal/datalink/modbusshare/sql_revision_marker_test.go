package modbusshare

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWorkspaceRevisionWritersRevokeFormatDirtyOwnership(t *testing.T) {
	for _, operation := range []string{"dirty", "revision", "desired"} {
		t.Run(operation, func(t *testing.T) {
			db, err := sql.Open("sqlite", ":memory:")
			require.NoError(t, err)
			t.Cleanup(func() { _ = db.Close() })
			db.SetMaxOpenConns(1)
			_, err = db.Exec(`CREATE TABLE system_settings(key TEXT PRIMARY KEY,value TEXT,updated_at DATETIME); INSERT INTO system_settings VALUES('modbus_share_workspace_revision:ws','{"revision":"rev-current","dirty":true,"point_format_pending_v1":{"version":1}}',NULL)`)
			require.NoError(t, err)
			store := NewSQLWorkspaceRevisionStore(db)
			switch operation {
			case "dirty":
				err = store.MarkDirty(t.Context(), "ws")
			case "revision":
				err = store.UpdateRevision(t.Context(), "ws", "rev-current", "rev-next")
			case "desired":
				err = store.CommitDesiredMappings(t.Context(), "ws", "rev-current", "rev-next", nil)
			}
			require.NoError(t, err)
			var raw string
			require.NoError(t, db.QueryRow(`SELECT value FROM system_settings WHERE key='modbus_share_workspace_revision:ws'`).Scan(&raw))
			require.NotContains(t, raw, "point_format_pending_v1")
			revision, dirty, err := store.GetRevision(t.Context(), "ws")
			require.NoError(t, err)
			if operation == "dirty" {
				require.Equal(t, "rev-current", revision)
				require.True(t, dirty)
			} else {
				require.Equal(t, "rev-next", revision)
			}
		})
	}
}
