package api

import (
	"database/sql"
	"path/filepath"
	"testing"

	"go-gateway/internal/testutil/lostcommit"

	"github.com/stretchr/testify/require"
)

// Exercise the same fault configuration locally: the first actual commit
// succeeds but returns ErrResponseLost, retries succeed, and their commits
// cannot overwrite the captured pre-retry payload. This is not SQL evidence.
func TestSQLServerCircuitAckFaultCapturesFirstCommitOnly(t *testing.T) {
	local, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "local.db"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, local.Close()) })
	_, err = local.ExecContext(t.Context(), `CREATE TABLE wg_delivery_outbox(group_id TEXT PRIMARY KEY,payload TEXT);
		INSERT INTO wg_delivery_outbox VALUES('group','before')`)
	require.NoError(t, err)
	faults := &lostcommit.Faults{}
	target := lostcommit.OpenFaults(filepath.Join(t.TempDir(), "target.db"), faults)
	t.Cleanup(func() { require.NoError(t, target.Close()) })
	_, err = target.ExecContext(t.Context(), `CREATE TABLE effects(id INTEGER PRIMARY KEY)`)
	require.NoError(t, err)
	capture := installSQLServerCircuitAckFault(t, faults, local, "group")
	tx, err := target.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	_, err = tx.ExecContext(t.Context(), `INSERT INTO effects VALUES(1)`)
	require.NoError(t, err)
	require.ErrorIs(t, tx.Commit(), lostcommit.ErrResponseLost)
	var count int
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM effects`).Scan(&count))
	require.Equal(t, 1, count, "the first transaction really committed")
	_, err = local.ExecContext(t.Context(), `UPDATE wg_delivery_outbox SET payload='after'`)
	require.NoError(t, err)
	tx, err = target.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	_, err = tx.ExecContext(t.Context(), `INSERT INTO effects VALUES(2)`)
	require.NoError(t, err)
	require.NoError(t, tx.Commit(), "a later real commit must not lose another ACK")
	capture.mu.Lock()
	defer capture.mu.Unlock()
	require.NoError(t, capture.err)
	require.Equal(t, "before", capture.payload, "retry must not overwrite the earlier captured bytes")
}
