package groupdelivery

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStoreBusyCommitCleansOwnedConnectionBeforeReuse(t *testing.T) {
	store, db, path := newTestStore(t)
	db.SetMaxOpenConns(1)
	_, err := db.ExecContext(t.Context(), `PRAGMA journal_mode=DELETE; PRAGMA busy_timeout=1; CREATE TABLE tx_probe(value INTEGER); INSERT INTO tx_probe VALUES(1)`)
	require.NoError(t, err)
	reader, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	defer reader.Close()
	held, err := reader.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer held.Rollback()
	var value int
	require.NoError(t, held.QueryRowContext(t.Context(), `SELECT value FROM tx_probe`).Scan(&value))
	require.Equal(t, 1, value)
	err = store.inTx(t.Context(), "busy-test", func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE tx_probe SET value=2`)
		return err
	})
	require.Error(t, err, "read lock must make DELETE-journal COMMIT busy")
	require.NoError(t, held.Rollback())
	require.NoError(t, store.inTx(t.Context(), "after-busy", func(ctx context.Context, tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, `SELECT value FROM tx_probe`).Scan(&value)
	}), "failed commit must not return an active native transaction to the pool")
	require.Equal(t, 1, value, "unacknowledged busy transaction was rolled back, not guessed committed")
}
