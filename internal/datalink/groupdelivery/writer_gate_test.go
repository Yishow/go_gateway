package groupdelivery

import (
	"context"
	"database/sql"
	"runtime"
	"testing"
	"time"
	"weak"

	"github.com/stretchr/testify/require"
)

func TestStoreWriterGateSharesOnlyDatabaseIdentity(t *testing.T) {
	a, b := new(sql.DB), new(sql.DB)
	require.Same(t, writerGateFor(a), writerGateFor(a))
	require.NotSame(t, writerGateFor(a), writerGateFor(b))
	held, release, err := acquireStoreWriter(t.Context(), a)
	require.NoError(t, err)
	defer release()
	for _, db := range []*sql.DB{a, b} {
		_, _, err := acquireStoreWriter(held, db)
		require.ErrorIs(t, err, ErrNestedStoreTransaction)
	}
}

func TestStoreWriterGateCancellationReleasesPermit(t *testing.T) {
	db := new(sql.DB)
	_, release, err := acquireStoreWriter(t.Context(), db)
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan error, 1)
	go func() {
		_, undo, err := acquireStoreWriter(ctx, db)
		if undo != nil {
			undo()
		}
		done <- err
	}()
	cancel()
	require.ErrorIs(t, <-done, context.Canceled)
	release()
	release()
	for i := 0; i < 100; i++ {
		_, undo, err := acquireStoreWriter(ctx, db)
		require.ErrorIs(t, err, context.Canceled)
		require.Nil(t, undo)
	}
	_, undo, err := acquireStoreWriter(t.Context(), db)
	require.NoError(t, err)
	undo()
}

func unreachableWriterDatabase() weak.Pointer[sql.DB] {
	db := new(sql.DB)
	writerGateFor(db)
	return weak.Make(db)
}

func TestStoreWriterGateDoesNotRetainDatabase(t *testing.T) {
	key := unreachableWriterDatabase()
	require.Eventually(t, func() bool {
		runtime.GC()
		if key.Value() != nil {
			return false
		}
		_, present := storeWriterGates.Load(key)
		return !present
	}, 5*time.Second, time.Millisecond)
}

func TestStoreWriterGateAcrossStoresAndClone(t *testing.T) {
	store, db, _ := newTestStore(t)
	stores := []*Store{store, NewStore(db), store.WithClock(time.Now)}
	entered := make(chan struct{})
	leave := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		done <- store.inTx(t.Context(), "held", func(ctx context.Context, tx *sql.Tx) error {
			close(entered)
			<-leave
			return nil
		})
	}()
	<-entered
	for _, other := range stores {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		err := other.inTx(ctx, "canceled", func(context.Context, *sql.Tx) error { t.Error("canceled callback ran"); return nil })
		require.ErrorIs(t, err, context.Canceled)
	}
	close(leave)
	require.NoError(t, <-done)
	for _, other := range stores {
		require.NoError(t, other.inTx(t.Context(), "next", func(ctx context.Context, tx *sql.Tx) error {
			err := store.inTx(ctx, "nested", func(context.Context, *sql.Tx) error { t.Error("nested callback ran"); return nil })
			require.ErrorIs(t, err, ErrNestedStoreTransaction)
			return nil
		}))
	}
}

// Cancel exactly after the select grants a free permit; this deterministically
// exercises the second cancellation check instead of relying on scheduler luck.
type cancelOnSecondErr struct {
	context.Context
	cancel context.CancelFunc
	calls  int
}

func (c *cancelOnSecondErr) Err() error {
	c.calls++
	if c.calls == 2 {
		c.cancel()
	}
	return c.Context.Err()
}
func TestStoreWriterGateCancelAfterGrantDoesNotLeak(t *testing.T) {
	db := new(sql.DB)
	base, cancel := context.WithCancel(t.Context())
	defer cancel()
	ctx := &cancelOnSecondErr{Context: base, cancel: cancel}
	_, undo, err := acquireStoreWriter(ctx, db)
	require.ErrorIs(t, err, context.Canceled)
	require.Nil(t, undo)
	require.Empty(t, writerGateFor(db).token)
}
