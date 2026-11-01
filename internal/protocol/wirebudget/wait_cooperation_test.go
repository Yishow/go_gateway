package wirebudget

import (
	"context"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/require"
)

type testWaitCooperator struct{ suspended, resumed atomic.Int32 }

func (c *testWaitCooperator) Suspend()                         { c.suspended.Add(1) }
func (c *testWaitCooperator) Resume(ctx context.Context) error { c.resumed.Add(1); return ctx.Err() }

func TestWaitCooperationRetainsBudgetAndDeadline(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		cooperation := &testWaitCooperator{}
		ctx, cancel := context.WithTimeout(t.Context(), time.Second)
		defer cancel()
		ctx = WithWaitCooperator(With(ctx, 3, 100*time.Millisecond), cooperation)
		require.NoError(t, Take(ctx))
		deadline, _ := ctx.Deadline()
		done := make(chan error, 1)
		go func() { done <- Wait(ctx) }()
		synctest.Wait()
		require.EqualValues(t, 1, cooperation.suspended.Load())
		require.Zero(t, cooperation.resumed.Load())
		require.Equal(t, 1, From(ctx).Used())
		time.Sleep(100 * time.Millisecond)
		require.NoError(t, <-done)
		require.EqualValues(t, 1, cooperation.resumed.Load())
		after, _ := ctx.Deadline()
		require.Equal(t, deadline, after)
		require.Equal(t, 1, From(ctx).Used())
	})
}

func TestWaitCooperationCancelDoesNotReacquire(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		cooperation := &testWaitCooperator{}
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		ctx = WithWaitCooperator(With(ctx, 3, time.Second), cooperation)
		require.NoError(t, Take(ctx))
		done := make(chan error, 1)
		go func() { done <- Wait(ctx) }()
		synctest.Wait()
		cancel()
		require.ErrorIs(t, <-done, context.Canceled)
		require.EqualValues(t, 1, cooperation.suspended.Load())
		require.Zero(t, cooperation.resumed.Load())
		require.Equal(t, 1, From(ctx).Used())
	})
}
