package wirebudget

import (
	"context"
	"time"
)

// WaitCooperator temporarily yields execution capacity while retaining request ownership.
// It must be used only by a single scan goroutine, never actual-send admission.
type WaitCooperator interface {
	Suspend()
	Resume(context.Context) error
}
type waitCooperatorKey struct{}

// WithWaitCooperator attaches optional execution-capacity cooperation to retry waits.
func WithWaitCooperator(ctx context.Context, cooperator WaitCooperator) context.Context {
	return context.WithValue(ctx, waitCooperatorKey{}, cooperator)
}

// WaitDuration preserves the caller's original deadline and retry delay. Optional
// cooperation releases only active capacity; the caller still owns its device gate.
func WaitDuration(ctx context.Context, delay time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if delay <= 0 {
		return ctx.Err()
	}
	cooperation, ok := ctx.Value(waitCooperatorKey{}).(WaitCooperator)
	if ok && cooperation != nil {
		cooperation.Suspend()
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if ok && cooperation != nil {
		return cooperation.Resume(ctx)
	}
	return nil
}

// ResumeAfterWait re-enters active execution before a subsequent request when a
// previous request expired during yielded backoff. It may block, so callers must
// invoke it outside short admission/observer callbacks and scheduler/manager
// global locks. The original device operation gate remains owned.
func ResumeAfterWait(ctx context.Context) error {
	cooperation, ok := ctx.Value(waitCooperatorKey{}).(WaitCooperator)
	if ok && cooperation != nil {
		return cooperation.Resume(ctx)
	}
	return ctx.Err()
}
