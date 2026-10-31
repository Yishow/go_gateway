package modbusshare

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFormatRecoveryWitnessRejectsSameRevisionDirty(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	db.SetMaxOpenConns(1)
	_, err = db.Exec(`CREATE TABLE system_settings(key TEXT PRIMARY KEY,value TEXT,updated_at DATETIME); INSERT INTO system_settings VALUES('modbus_share_workspace_revision:ws','{"revision":"rev","dirty":false}',NULL)`)
	require.NoError(t, err)
	var witness string
	require.NoError(t, db.QueryRow(`SELECT value FROM system_settings WHERE key='modbus_share_workspace_revision:ws'`).Scan(&witness))
	ctx := WithFormatRecoveryWitness(t.Context(), "ws", witness)
	store := NewSQLWorkspaceRevisionStore(db)
	require.NoError(t, store.MarkDirty(t.Context(), "ws"))
	require.Error(t, store.CommitDesiredMappings(ctx, "ws", "rev", "rev", nil), "captured clean revision must not clear a later SAME revision unknown")
	_, dirty, err := store.GetRevision(t.Context(), "ws")
	require.NoError(t, err)
	require.True(t, dirty)
}
