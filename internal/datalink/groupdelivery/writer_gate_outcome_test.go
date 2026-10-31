package groupdelivery

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"go-gateway/internal/datalink/sqltx"

	"github.com/stretchr/testify/require"
)

func TestStoreWriterGateUnknownCommitDoesNotRetry(t *testing.T) {
	store, db, _ := newTestStore(t)
	db.SetMaxOpenConns(1)
	_, err := db.ExecContext(t.Context(), `CREATE TABLE gate_outcome(value INTEGER); INSERT INTO gate_outcome VALUES(0)`)
	require.NoError(t, err)
	calls := 0
	err = store.inTx(t.Context(), "lost-terminal-ack", func(ctx context.Context, tx *sql.Tx) error {
		calls++
		_, err := tx.ExecContext(ctx, `UPDATE gate_outcome SET value=value+1`)
		if err != nil {
			return err
		}
		// Inject an already completed native commit, so the owner's terminal
		// acknowledgment cannot establish its outcome. Do not retry the callback.
		return tx.Commit()
	})
	require.ErrorIs(t, err, sqltx.ErrCommitOutcomeUnknown)
	require.Equal(t, 1, calls)
	require.Zero(t, db.Stats().InUse)
	require.NoError(t, store.inTx(t.Context(), "after-unknown", func(ctx context.Context, tx *sql.Tx) error {
		var value int
		if err := tx.QueryRowContext(ctx, `SELECT value FROM gate_outcome`).Scan(&value); err != nil {
			return err
		}
		require.Equal(t, 1, value)
		// Even a different DB gate cannot be nested from a marked callback.
		foreign := NewStore(new(sql.DB))
		err := foreign.inTx(ctx, "foreign", func(context.Context, *sql.Tx) error { t.Error("nested callback ran"); return nil })
		require.ErrorIs(t, err, ErrNestedStoreTransaction)
		return nil
	}))
}

func TestStoreWriterGateConcurrentStoresAndClones(t *testing.T) {
	store, db, _ := newTestStore(t)
	db.SetMaxOpenConns(8)
	_, err := db.ExecContext(t.Context(), `CREATE TABLE gate_count(value INTEGER); INSERT INTO gate_count VALUES(0)`)
	require.NoError(t, err)
	stores := []*Store{store, NewStore(db), store.WithClock(time.Now)}
	results := make(chan error, 12)
	for i := 0; i < 12; i++ {
		s := stores[i%len(stores)]
		go func() {
			for j := 0; j < 8; j++ {
				err := s.inTx(t.Context(), "concurrent-clone", func(ctx context.Context, tx *sql.Tx) error {
					_, err := tx.ExecContext(ctx, `UPDATE gate_count SET value=value+1`)
					return err
				})
				if err != nil {
					results <- err
					return
				}
			}
			results <- nil
		}()
	}
	for i := 0; i < 12; i++ {
		require.NoError(t, <-results)
	}
	var value int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT value FROM gate_count`).Scan(&value))
	require.Equal(t, 96, value)
	require.Zero(t, db.Stats().InUse)
}
