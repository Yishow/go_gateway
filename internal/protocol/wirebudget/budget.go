// Package wirebudget carries one read request's actual-send allowance through
// scheduler, adapter and transport layers. Reconnects never mint more attempts.
package wirebudget

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

var ErrExhausted = errors.New("wire attempt budget exhausted")

type contextKey struct{}

// Budget is shared only within one logical request. Limits are immutable.
type Budget struct {
	limit int64
	used  atomic.Int64
	delay time.Duration
}

// With installs a budget only when the caller does not already own one.
func With(ctx context.Context, attempts int, delay time.Duration) context.Context {
	if From(ctx) != nil {
		return ctx
	}
	if attempts < 1 {
		attempts = 1
	}
	if delay < 0 {
		delay = 0
	}
	return context.WithValue(ctx, contextKey{}, &Budget{limit: int64(attempts), delay: delay})
}
func From(ctx context.Context) *Budget { b, _ := ctx.Value(contextKey{}).(*Budget); return b }
func (b *Budget) Used() int {
	if b == nil {
		return 0
	}
	return int(b.used.Load())
}
func (b *Budget) Remaining() int {
	if b == nil {
		return 0
	}
	return int(b.limit - b.used.Load())
}

// Take must run immediately before an actual send, not before a reconnect.
func Take(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	b := From(ctx)
	if b == nil {
		return nil
	}
	for {
		n := b.used.Load()
		if n >= b.limit {
			return ErrExhausted
		}
		if b.used.CompareAndSwap(n, n+1) {
			return nil
		}
	}
}

// Wait uses the configured delay, exponentially bounded at thirty seconds.
// The original context deadline includes this wait and is never extended.
func Wait(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	b := From(ctx)
	if b == nil {
		return nil
	}
	if b.Remaining() == 0 {
		return ErrExhausted
	}
	delay := b.delay
	for i := 1; i < b.Used() && delay < 30*time.Second; i++ {
		if delay > 15*time.Second {
			delay = 30 * time.Second
			break
		}
		delay *= 2
	}
	if delay > 30*time.Second {
		delay = 30 * time.Second
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return ctx.Err()
	}
}
