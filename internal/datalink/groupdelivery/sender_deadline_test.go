package groupdelivery

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const deadlineSettlementBudget = 100 * time.Millisecond

// The resolver signals when the remote attempt has started, then remains
// blocked until the test has observed a complete settlement budget elapse.
type deadlineBarrier struct {
	target  Target
	err     error
	entered chan context.Context
	release chan struct{}
}

func (b *deadlineBarrier) Resolve(ctx context.Context, _ OutboxItem) (Target, error) {
	b.entered <- ctx
	select {
	case <-b.release:
		return b.target, b.err
	case <-ctx.Done():
		return Target{}, ctx.Err()
	}
}

func newDeadlineFixture(t *testing.T, kind string, receipt, rejectRow bool) (fixture *deliveryFixture, effectKey string, rowCount func() int) {
	t.Helper()
	capability := DedupeNone
	if receipt {
		capability = DedupeReceipt
	}
	if kind == "postgres" {
		f, namespace, admin := newPostgresDeliveryFixture(t, receipt)
		if rejectRow {
			_, err := admin.ExecContext(t.Context(), `ALTER TABLE "`+namespace+`".readings ADD CHECK (counter < 0)`)
			require.NoError(t, err)
		}
		effect := f.enqueueSchema(t, namespace, capability, 5)
		return f, effect, func() int { return pgRows(t, admin, namespace) }
	}
	ddl := plainTargetTable
	if rejectRow {
		ddl = strings.Replace(ddl, "counter INTEGER NOT NULL", "counter INTEGER NOT NULL CHECK(counter < 0)", 1)
	}
	f := newDeliveryFixture(t, ddl, receipt)
	effect := f.enqueue(t, capability, "", t0, 5)
	return f, effect, func() int { return f.targetRows(t) }
}

func deliverAfterRemoteBudget(ctx context.Context, t *testing.T, f *deliveryFixture, effect string) (DeliveryResult, error) {
	t.Helper()
	gate := &deadlineBarrier{target: f.resolver.target, err: f.resolver.err,
		entered: make(chan context.Context, 1), release: make(chan struct{})}
	release := sync.OnceFunc(func() { close(gate.release) })
	sender := NewSender(f.store, gate, SenderConfig{Owner: "deadline-test/current", LeaseTTL: time.Millisecond,
		DeliveryTimeout: 5 * time.Second, SettleTimeout: deadlineSettlementBudget, Backoff: f.sender.backoff})
	type outcome struct {
		result DeliveryResult
		err    error
	}
	finished := make(chan outcome, 1)
	var wg sync.WaitGroup
	wg.Go(func() {
		result, err := sender.Deliver(ctx, effect)
		finished <- outcome{result, err}
	})
	defer wg.Wait()
	// Release before joining the goroutine, including assertion-failure paths.
	defer release()
	var remote context.Context
	select {
	case remote = <-gate.entered:
	case result := <-finished:
		t.Fatalf("delivery ended before reaching remote barrier: %+v", result)
	case <-time.After(10 * time.Second):
		t.Fatal("delivery did not reach remote barrier")
	}
	budget := time.NewTimer(2 * deadlineSettlementBudget)
	defer budget.Stop()
	select {
	case <-budget.C:
	case <-remote.Done():
		t.Fatal("attempt ended before the remote barrier was released")
	}
	require.NoError(t, remote.Err(), "remote attempt remains legal after a complete settlement budget")
	// An undersized requested lease must still protect this live remote attempt.
	var lease string
	require.NoError(t, f.local.QueryRowContext(ctx, `SELECT claim_expires_at FROM wg_delivery_outbox WHERE effect_key = ?`, effect).Scan(&lease))
	until, err := time.Parse(time.RFC3339Nano, lease)
	require.NoError(t, err)
	require.Greater(t, time.Until(until), 5*time.Second, "lease includes attempt, settlement and existing margin")
	release()
	select {
	case result := <-finished:
		return result.result, result.err
	case <-time.After(10 * time.Second):
		t.Fatal("delivery did not settle after remote release")
		return DeliveryResult{}, nil
	}
}

func requireDeadlineReceipt(t *testing.T, f *deliveryFixture, effect string) {
	t.Helper()
	var key, digest string
	require.NoError(t, f.local.QueryRowContext(t.Context(), `SELECT effect_key, payload_digest FROM wg_delivery_receipts`).Scan(&key, &digest))
	require.Equal(t, effect, key)
	require.Equal(t, f.state(t, effect).PayloadDigest, digest)
}

func TestSenderDeadlineSlowRemoteResultsAreDurablySettled(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		for _, scenario := range []struct {
			name, state, code string
			receipt, reject   bool
			resolveErr        error
			loseCommit        bool
			rows              int
		}{
			{name: "success", state: StateCommitted, receipt: true, rows: 1},
			{name: "transient", state: StateRetrying, code: "target-unavailable", resolveErr: errors.New("offline")},
			{name: "target-blocked", state: StateBlocked, code: "target-blocked", resolveErr: ErrTargetBlocked},
			{name: "row-rejected", state: StateQuarantined, code: "destination-rejected-row", reject: true},
			{name: "ambiguous-with-receipt", state: StateRetrying, code: "commit-ambiguous", receipt: true, loseCommit: true, rows: 1},
			{name: "ambiguous-without-dedupe", state: StateUnknown, code: "commit-ambiguous", loseCommit: true, rows: 1},
		} {
			t.Run(kind+"/"+scenario.name, func(t *testing.T) {
				f, effect, rows := newDeadlineFixture(t, kind, scenario.receipt, scenario.reject)
				f.resolver.err = scenario.resolveErr
				f.lose.Store(scenario.loseCommit)
				result, err := deliverAfterRemoteBudget(t.Context(), t, f, effect)
				require.NoError(t, err)
				require.Equal(t, scenario.state, result.State)
				item := f.state(t, effect)
				require.Equal(t, effect, item.EffectKey)
				require.Equal(t, scenario.state, item.State)
				require.Equal(t, scenario.code, item.LastErrorCode)
				require.Equal(t, scenario.rows, rows())
				if scenario.state == StateCommitted {
					requireDeadlineReceipt(t, f, effect)
				} else {
					require.Zero(t, countRows(t, f.local, "wg_delivery_receipts"))
				}
			})
		}
	}
}
