package groupdelivery

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStoreBusyCommitCleansOwnedConnectionBeforeReuse(t *testing.T) {
	store, db, path := newTestStore(t)
	db.SetMaxOpenConns(1)
	_, err := db.Exec(`PRAGMA journal_mode=DELETE; PRAGMA busy_timeout=1; CREATE TABLE tx_probe(value INTEGER); INSERT INTO tx_probe VALUES(1)`)
	require.NoError(t, err)
	reader, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	defer reader.Close()
	held, err := reader.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer held.Rollback()
	var value int
	require.NoError(t, held.QueryRow(`SELECT value FROM tx_probe`).Scan(&value))
	require.Equal(t, 1, value)
	err = store.inTx(t.Context(), "busy-test", func(tx *sql.Tx) error { _, err := tx.Exec(`UPDATE tx_probe SET value=2`); return err })
	require.Error(t, err, "read lock must make DELETE-journal COMMIT busy")
	require.NoError(t, held.Rollback())
	require.NoError(t, store.inTx(t.Context(), "after-busy", func(tx *sql.Tx) error { return tx.QueryRow(`SELECT value FROM tx_probe`).Scan(&value) }), "failed commit must not return an active native transaction to the pool")
	require.Equal(t, 1, value, "unacknowledged busy transaction was rolled back, not guessed committed")
}
