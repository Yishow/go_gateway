package wirebudget

import (
	"context"
	"github.com/stretchr/testify/require"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestBudgetSharedConcurrentSendCap(t *testing.T) {
	ctx := With(t.Context(), 4, 0)
	same := With(ctx, 100, 0)
	require.Same(t, From(ctx), From(same))
	var wg sync.WaitGroup
	var sent atomic.Int32
	for range 40 {
		wg.Go(func() {
			if Take(ctx) == nil {
				sent.Add(1)
			}
		})
	}
	wg.Wait()
	require.EqualValues(t, 4, sent.Load())
	require.Equal(t, 0, From(ctx).Remaining())
	require.ErrorIs(t, Take(ctx), ErrExhausted)
}
func TestBudgetCanceledAndDeadlineBackoff(t *testing.T) {
	ctx, cancel := context.WithCancel(With(t.Context(), 4, time.Hour))
	cancel()
	require.ErrorIs(t, Take(ctx), context.Canceled)
	require.Zero(t, From(ctx).Used())
	ctx, cancel = context.WithTimeout(With(t.Context(), 4, time.Hour), 10*time.Millisecond)
	defer cancel()
	require.NoError(t, Take(ctx))
	require.ErrorIs(t, Wait(ctx), context.DeadlineExceeded)
	require.Equal(t, 1, From(ctx).Used())
}
