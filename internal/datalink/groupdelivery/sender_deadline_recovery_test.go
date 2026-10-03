package groupdelivery

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSenderDeadlineSlowCommitResponseGetsAFreshSettlementBudget(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			f, effect, rows := newDeadlineFixture(t, kind, true, false)
			entered, resume := make(chan struct{}), make(chan struct{})
			release := sync.OnceFunc(func() { close(resume) })
			afterCommit := func() {
				close(entered)
				<-resume // the real destination has committed; its reply is held
			}
			f.faults.AfterCommit.Store(&afterCommit)
			sender := NewSender(f.store, f.resolver, SenderConfig{
				DeliveryTimeout: 5 * time.Second, SettleTimeout: deadlineSettlementBudget})
			type outcome struct {
				result DeliveryResult
				err    error
			}
			finished := make(chan outcome, 1)
			var wg sync.WaitGroup
			wg.Go(func() {
				result, err := sender.Deliver(t.Context(), effect)
				finished <- outcome{result, err}
			})
			defer wg.Wait()
			defer release()
			select {
			case <-entered:
			case result := <-finished:
				t.Fatalf("delivery did not reach real commit: %+v", result)
			case <-time.After(10 * time.Second):
				t.Fatal("destination did not commit")
			}
			require.Equal(t, 1, rows(), "independent SELECT observes the committed row before its reply")
			require.Equal(t, StateSending, f.state(t, effect).State)
			require.Zero(t, countRows(t, f.local, "wg_delivery_receipts"))
			budget := time.NewTimer(2 * deadlineSettlementBudget)
			defer budget.Stop()
			select {
			case <-budget.C:
			case result := <-finished:
				t.Fatalf("delivery returned while commit response was held: %+v", result)
			}
			release()
			select {
			case result := <-finished:
				require.NoError(t, result.err)
				require.Equal(t, StateCommitted, result.result.State)
			case <-time.After(10 * time.Second):
				t.Fatal("local settlement did not finish")
			}
			requireDeadlineReceipt(t, f, effect)
			require.Equal(t, 1, rows())
		})
	}
}

func TestSenderDeadlineCancelledCallerAfterSlowCommitStillRecordsTheEffect(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			f, effect, rows := newDeadlineFixture(t, kind, true, false)
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			afterCommit := func() { cancel() }
			f.faults.AfterCommit.Store(&afterCommit)
			result, err := deliverAfterRemoteBudget(ctx, t, f, effect)
			require.ErrorIs(t, ctx.Err(), context.Canceled, "cancellation happened after a real target commit")
			require.NoError(t, err)
			require.Equal(t, StateCommitted, result.State)
			require.Equal(t, StateCommitted, f.state(t, effect).State)
			require.Equal(t, 1, rows())
			requireDeadlineReceipt(t, f, effect)
		})
	}
}

func TestSenderDeadlineLocalReceiptFailureRetainsEvidenceAndRecoversSafely(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		for _, receipt := range []bool{false, true} {
			name := "without-dedupe"
			if receipt {
				name = "with-receipt"
			}
			t.Run(kind+"/"+name, func(t *testing.T) {
				f, effect, rows := newDeadlineFixture(t, kind, receipt, false)
				_, err := f.local.ExecContext(t.Context(), `CREATE TRIGGER no_local_receipt BEFORE INSERT ON wg_delivery_receipts
					BEGIN SELECT RAISE(ABORT, 'injected local failure'); END`)
				require.NoError(t, err)
				result, err := deliverAfterRemoteBudget(t.Context(), t, f, effect)
				require.ErrorIs(t, err, ErrLocalReceiptFailed)
				require.Equal(t, StateSending, result.State)
				require.Equal(t, StateSending, f.state(t, effect).State)
				require.Equal(t, effect, f.state(t, effect).EffectKey)
				require.Equal(t, 1, rows(), "remote SQL committed even though local receipt storage failed")
				require.Zero(t, countRows(t, f.local, "wg_delivery_receipts"))
				_, err = f.sender.Deliver(t.Context(), effect)
				require.ErrorIs(t, err, ErrNotDeliverable, "sending is never blindly sent a second time")

				resolved, err := f.store.ResolveInterrupted(t.Context(), effect)
				require.NoError(t, err)
				_, err = f.local.ExecContext(t.Context(), `DROP TRIGGER no_local_receipt`)
				require.NoError(t, err)
				if receipt {
					require.Equal(t, StateRetrying, resolved)
					result, err = f.sender.Deliver(t.Context(), effect)
					require.NoError(t, err)
					require.Equal(t, StateCommitted, result.State)
					requireDeadlineReceipt(t, f, effect)
				} else {
					require.Equal(t, StateUnknown, resolved)
					_, err = f.sender.Deliver(t.Context(), effect)
					require.ErrorIs(t, err, ErrNotDeliverable)
					require.Zero(t, countRows(t, f.local, "wg_delivery_receipts"))
				}
				require.Equal(t, 1, rows(), "recovery keeps exactly one committed SQL effect")
			})
		}
	}
}
