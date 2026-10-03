package groupdelivery

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type delayedResolver struct {
	target Target
	err    error
	delay  time.Duration
}

func (r delayedResolver) Resolve(ctx context.Context, _ OutboxItem) (Target, error) {
	timer := time.NewTimer(r.delay)
	defer timer.Stop()
	select {
	case <-timer.C:
		return r.target, r.err
	case <-ctx.Done():
		return Target{}, ctx.Err()
	}
}

func immediateSettlementRetry(int) (string, time.Time) {
	return StateRetrying, time.Now().UTC().Add(-time.Second)
}

func TestSenderSlowRemoteSuccessStartsSettlementTimeoutAfterAttempt(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, true)
	effect := f.enqueueRaw(t, DedupeReceipt, t0, 5)
	sender := NewSender(f.store, delayedResolver{target: Target{DB: f.target, Kind: f.resolver.target.Kind}, delay: 20 * time.Millisecond}, SenderConfig{
		Owner:           "node-1/a",
		DeliveryTimeout: 200 * time.Millisecond,
		SettleTimeout:   5 * time.Millisecond,
		Backoff:         immediateSettlementRetry,
	})

	result, err := sender.Deliver(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, StateCommitted, result.State)
	require.Equal(t, 1, countRows(t, f.local, "wg_delivery_receipts"))
}

func TestSenderSlowRemoteFailureStartsSettlementTimeoutAfterAttempt(t *testing.T) {
	f := newDeliveryFixture(t, plainTargetTable, false)
	effect := f.enqueueRaw(t, DedupeNone, t0, 5)
	sender := NewSender(f.store, delayedResolver{err: errors.New("destination unavailable"), delay: 20 * time.Millisecond}, SenderConfig{
		Owner:           "node-1/a",
		DeliveryTimeout: 200 * time.Millisecond,
		SettleTimeout:   5 * time.Millisecond,
		Backoff:         immediateSettlementRetry,
	})

	result, err := sender.Deliver(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, StateRetrying, result.State)
	require.Equal(t, codeTargetUnavailable, f.state(t, effect).LastErrorCode)
}
