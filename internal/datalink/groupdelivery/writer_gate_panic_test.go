package groupdelivery

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStoreTransactionPanicRollsBackBeforeReturningConnection(t *testing.T) {
	store, db, _ := newTestStore(t)
	db.SetMaxOpenConns(1)
	_, err := db.Exec(`CREATE TABLE panic_probe(value INTEGER); INSERT INTO panic_probe VALUES(1)`)
	require.NoError(t, err)
	require.PanicsWithValue(t, "callback panic", func() {
		_ = store.inTx(t.Context(), "panic", func(ctx context.Context, tx *sql.Tx) error {
			_, err := tx.Exec(`UPDATE panic_probe SET value=2`)
			require.NoError(t, err)
			panic("callback panic")
		})
	})
	require.Zero(t, db.Stats().InUse, "panic must terminate transaction before releasing writer permit")
	var value int
	require.NoError(t, store.inTx(t.Context(), "after-panic", func(ctx context.Context, tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, `SELECT value FROM panic_probe`).Scan(&value)
	}))
	require.Equal(t, 1, value)
}
