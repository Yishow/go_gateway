package groupdelivery

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// The destination attempt, including a storage restore and the insert that
// follows it, is bounded by DeliveryTimeout. Local settlement is a separate
// step with a full SettleTimeout of its own that starts only once all external
// work is over: a restore that takes longer than SettleTimeout must not leave
// the settlement an expired context.
const (
	recoverySettleBudget = 100 * time.Millisecond
	// recoverySlowRestore outlasts two settlement budgets and stays far inside
	// the delivery timeout, so a budget started before the restore has certainly
	// run out by the time the restore is over.
	recoverySlowRestore = 2*recoverySettleBudget + recoverySettleBudget/2
)

// slowRecoverySender returns a sender whose storage recovery answers decision
// after waiting recoverySlowRestore, and recreates the dropped table first when
// rebuild is set, the way a real restore does.
func slowRecoverySender(t *testing.T, f *deliveryFixture, decision RestoreDecision, restoreErr error, rebuild bool, during func()) *fakeRecovery {
	t.Helper()
	recovery := &fakeRecovery{decision: decision, restoreErr: restoreErr, restore: func(target Target) {
		if during != nil {
			during()
		}
		time.Sleep(recoverySlowRestore)
		if rebuild {
			_, err := target.DB.ExecContext(context.Background(), plainTargetTable)
			require.NoError(t, err)
		}
	}}
	f.sender = NewSender(f.store, f.resolver, SenderConfig{
		MaxRetries: 3, Backoff: f.sender.backoff, Recovery: recovery,
		DeliveryTimeout: 5 * time.Second, SettleTimeout: recoverySettleBudget,
	})
	return recovery
}

func TestSenderRestoreSlowerThanSettlementBudgetStillSettles(t *testing.T) {
	t.Run("the restored table takes the row and its receipt is recorded", func(t *testing.T) {
		f := newDeliveryFixture(t, plainTargetTable, true)
		effect := f.enqueue(t, DedupeReceipt, "", t0, 7)
		_, err := f.target.ExecContext(t.Context(), `DROP TABLE readings`)
		require.NoError(t, err)
		recovery := slowRecoverySender(t, f, RestoreDecision{Action: RestoreRetryNow}, nil, true, nil)

		result, err := f.sender.Deliver(t.Context(), effect)
		require.NoError(t, err, "the destination committed; settlement must not fail on a spent budget")
		require.Equal(t, StateCommitted, result.State)
		require.Equal(t, StateCommitted, f.state(t, effect).State, "not left in sending")
		require.EqualValues(t, 1, recovery.restores.Load())
		require.Equal(t, 1, f.targetRows(t))
		require.Equal(t, 1, countRows(t, f.target, "gw_effect_receipts"), "the destination's own receipt")
		require.Equal(t, 1, countRows(t, f.local, "wg_delivery_receipts"), "the local receipt")
	})

	t.Run("a second refusal after the slow restore is classified and recorded", func(t *testing.T) {
		f := newDeliveryFixture(t, plainTargetTable, true)
		effect := f.enqueue(t, DedupeReceipt, "", t0, 7)
		_, err := f.target.ExecContext(t.Context(), `DROP TABLE readings`)
		require.NoError(t, err)
		slowRecoverySender(t, f, RestoreDecision{Action: RestoreRetryNow}, nil, false, nil)

		result, err := f.sender.Deliver(t.Context(), effect)
		require.NoError(t, err)
		require.Equal(t, [2]string{StateBlocked, codeTargetUnusable}, [2]string{result.State, result.Code})
		require.Equal(t, [2]string{StateBlocked, codeTargetUnusable}, [2]string{f.state(t, effect).State, f.state(t, effect).LastErrorCode})
	})

	for name, tc := range map[string]struct {
		decision   RestoreDecision
		restoreErr error
		state      string
		code       string
	}{
		"waiting for another restore":    {RestoreDecision{Action: RestoreWait}, nil, StateRetrying, codeStorageRestoring},
		"a refusal":                      {RestoreDecision{Action: RestoreBlock, Code: "storage_owner_mismatch"}, nil, StateBlocked, "storage_owner_mismatch"},
		"a recovery that cannot decide":  {RestoreDecision{}, errors.New("catalog unavailable"), StateRetrying, codeStorageRecoveryUnavailable},
		"nothing to say about the table": {RestoreDecision{Action: RestoreNotApplicable}, nil, StateBlocked, codeTargetUnusable},
	} {
		t.Run(name+" is recorded with a full budget", func(t *testing.T) {
			f := newDeliveryFixture(t, plainTargetTable, true)
			effect := f.enqueue(t, DedupeReceipt, "", t0, 7)
			_, err := f.target.ExecContext(t.Context(), `DROP TABLE readings`)
			require.NoError(t, err)
			slowRecoverySender(t, f, tc.decision, tc.restoreErr, false, nil)

			result, err := f.sender.Deliver(t.Context(), effect)
			require.NoError(t, err)
			require.Equal(t, [2]string{tc.state, tc.code}, [2]string{result.State, result.Code})
			require.Equal(t, [2]string{tc.state, tc.code}, [2]string{f.state(t, effect).State, f.state(t, effect).LastErrorCode})
		})
	}

	t.Run("a caller canceled while the restore ran can still settle", func(t *testing.T) {
		f := newDeliveryFixture(t, plainTargetTable, true)
		effect := f.enqueue(t, DedupeReceipt, "", t0, 7)
		_, err := f.target.ExecContext(t.Context(), `DROP TABLE readings`)
		require.NoError(t, err)
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		slowRecoverySender(t, f, RestoreDecision{Action: RestoreWait}, nil, false, cancel)

		result, err := f.sender.Deliver(ctx, effect)
		require.NoError(t, err, "shutdown does not stop the local record of what happened")
		require.ErrorIs(t, ctx.Err(), context.Canceled)
		require.Equal(t, [2]string{StateRetrying, codeStorageRestoring}, [2]string{result.State, result.Code})
		require.Equal(t, StateRetrying, f.state(t, effect).State)
	})
}
