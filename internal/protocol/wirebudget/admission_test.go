package wirebudget

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdmissionEveryActualSendPreservesBudget(t *testing.T) {
	eligible := true
	denied := errors.New("installation no longer current")
	checks, sends := 0, 0
	ctx := With(t.Context(), 3, 0)
	ctx = WithAdmission(ctx, func() error {
		checks++
		if !eligible {
			return denied
		}
		return nil
	})
	ctx = WithSendObserver(ctx, func() { sends++ })
	require.NoError(t, Take(ctx))
	eligible = false
	require.ErrorIs(t, Take(ctx), denied)
	require.Equal(t, 2, checks)
	require.Equal(t, 1, sends)
	require.Equal(t, 1, From(ctx).Used())
	eligible = true
	require.NoError(t, Take(ctx))
	require.Equal(t, 2, From(ctx).Used())
}

func TestAdmissionComposesBeforeObserverAndHonorsCancellation(t *testing.T) {
	var order []string
	ctx := WithAdmission(t.Context(), func() error { order = append(order, "parent"); return nil })
	ctx = WithAdmission(ctx, func() error { order = append(order, "child"); return nil })
	ctx = WithSendObserver(ctx, func() { order = append(order, "first-observer") })
	ctx = WithSendObserver(ctx, func() { order = append(order, "second-observer") })
	require.NoError(t, Take(ctx), "admission also protects callers with no budget")
	require.Equal(t, []string{"parent", "child", "first-observer", "second-observer"}, order)
	canceled, cancel := context.WithCancel(ctx)
	cancel()
	require.ErrorIs(t, Take(canceled), context.Canceled)
	require.Len(t, order, 4, "cancellation must reject before any callback")
}

func TestAdmissionFailureAndExhaustionNeverObserveSend(t *testing.T) {
	denied := errors.New("no lease")
	called := 0
	ctx := WithAdmission(With(t.Context(), 1, 0), func() error { return denied })
	ctx = WithSendObserver(ctx, func() { called++ })
	require.ErrorIs(t, Take(ctx), denied)
	require.Zero(t, From(ctx).Used())
	require.Zero(t, called)
	ctx = WithSendObserver(With(t.Context(), 1, 0), func() { called++ })
	require.NoError(t, Take(ctx))
	require.ErrorIs(t, Take(ctx), ErrExhausted)
	require.Equal(t, 1, called)
}

func TestSendObserverPanicDoesNotChangeAcquisition(t *testing.T) {
	ctx := WithSendObserver(With(t.Context(), 1, 0), func() { panic("optional diagnostic") })
	require.NotPanics(t, func() { require.NoError(t, Take(ctx)) })
	require.Equal(t, 1, From(ctx).Used())
}
