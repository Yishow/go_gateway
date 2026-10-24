package groupdelivery

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

// fakeRecovery answers the way the production recovery would, and lets a test
// change the destination the way a restore does.
type fakeRecovery struct {
	admit      string
	admitErr   error
	decision   RestoreDecision
	restoreErr error
	restore    func(Target)
	blocked    string
	admits     atomic.Int32
	restores   atomic.Int32
	blockedAsk atomic.Int32
	lastCause  atomic.Value
}

func (r *fakeRecovery) Admit(context.Context, OutboxItem) (string, error) {
	r.admits.Add(1)
	return r.admit, r.admitErr
}

func (r *fakeRecovery) Restore(_ context.Context, _ OutboxItem, target Target, cause error) (RestoreDecision, error) {
	r.restores.Add(1)
	r.lastCause.Store(cause.Error())
	if r.restore != nil {
		r.restore(target)
	}
	return r.decision, r.restoreErr
}

func (r *fakeRecovery) BlockedCode(context.Context, OutboxItem) string {
	r.blockedAsk.Add(1)
	return r.blocked
}

func newRecoveringSender(t *testing.T, f *deliveryFixture, recovery StorageRecovery) {
	t.Helper()
	config := f.sender.config
	config.Recovery = recovery
	f.sender = NewSender(f.store, f.resolver, config)
}

// TestSenderRestoresAConfirmedTableOnlyWhenRecoveryAllows: a destination that
// refuses an insert because its table is gone is handed to storage recovery,
// which alone decides; every other failure never reaches it.
func TestSenderRestoresAConfirmedTableOnlyWhenRecoveryAllows(t *testing.T) {
	t.Run("a restored table gets the insert tried again, once, and the row commits", func(t *testing.T) {
		f := newDeliveryFixture(t, plainTargetTable, true)
		effect := f.enqueue(t, DedupeReceipt, "", t0, 7)
		_, err := f.target.ExecContext(t.Context(), `DROP TABLE readings`)
		require.NoError(t, err)
		recovery := &fakeRecovery{decision: RestoreDecision{Action: RestoreRetryNow}, restore: func(target Target) {
			_, err := target.DB.ExecContext(t.Context(), plainTargetTable)
			require.NoError(t, err)
		}}
		newRecoveringSender(t, f, recovery)

		result, err := f.sender.Deliver(t.Context(), effect)
		require.NoError(t, err)
		require.Equal(t, StateCommitted, result.State)
		require.EqualValues(t, 1, recovery.restores.Load())
		require.Equal(t, 1, f.targetRows(t))
		require.Equal(t, 1, countRows(t, f.target, "gw_effect_receipts"))
		require.Equal(t, 1, countRows(t, f.local, "wg_delivery_receipts"))
	})

	t.Run("a table that is still gone after a restore is not restored in a loop", func(t *testing.T) {
		f := newDeliveryFixture(t, plainTargetTable, true)
		effect := f.enqueue(t, DedupeReceipt, "", t0, 7)
		_, err := f.target.ExecContext(t.Context(), `DROP TABLE readings`)
		require.NoError(t, err)
		recovery := &fakeRecovery{decision: RestoreDecision{Action: RestoreRetryNow}}
		newRecoveringSender(t, f, recovery)

		result, err := f.sender.Deliver(t.Context(), effect)
		require.NoError(t, err)
		require.Equal(t, StateBlocked, result.State, "the second refusal is the ordinary one")
		require.Equal(t, codeTargetUnusable, result.Code)
		require.EqualValues(t, 1, recovery.restores.Load(), "recovery is asked once per attempt")
	})

	for name, tc := range map[string]struct {
		decision RestoreDecision
		state    string
		code     string
	}{
		"a refusal blocks with its own stable code": {RestoreDecision{Action: RestoreBlock, Code: "storage_owner_mismatch"}, StateBlocked, "storage_owner_mismatch"},
		"another execution restoring makes it wait": {RestoreDecision{Action: RestoreWait}, StateRetrying, "storage-restore-in-progress"},
		"nothing to say keeps the ordinary block":   {RestoreDecision{Action: RestoreNotApplicable}, StateBlocked, codeTargetUnusable},
	} {
		t.Run(name, func(t *testing.T) {
			f := newDeliveryFixture(t, plainTargetTable, true)
			effect := f.enqueue(t, DedupeReceipt, "", t0, 7)
			_, err := f.target.ExecContext(t.Context(), `DROP TABLE readings`)
			require.NoError(t, err)
			recovery := &fakeRecovery{decision: tc.decision}
			newRecoveringSender(t, f, recovery)

			result, err := f.sender.Deliver(t.Context(), effect)
			require.NoError(t, err)
			require.Equal(t, [2]string{tc.state, tc.code}, [2]string{result.State, result.Code})
			require.Equal(t, [2]string{tc.state, tc.code}, [2]string{f.state(t, effect).State, f.state(t, effect).LastErrorCode})
			require.Equal(t, 0, countRows(t, f.local, "wg_delivery_receipts"), "nothing was recorded as delivered")
		})
	}

	t.Run("a recovery that cannot decide leaves the row to retry", func(t *testing.T) {
		f := newDeliveryFixture(t, plainTargetTable, true)
		effect := f.enqueue(t, DedupeReceipt, "", t0, 7)
		_, err := f.target.ExecContext(t.Context(), `DROP TABLE readings`)
		require.NoError(t, err)
		newRecoveringSender(t, f, &fakeRecovery{restoreErr: errors.New("catalog unavailable")})
		result, err := f.sender.Deliver(t.Context(), effect)
		require.NoError(t, err)
		require.Equal(t, [2]string{StateRetrying, "storage-recovery-unavailable"}, [2]string{result.State, result.Code})
	})

	t.Run("failures that are not about the table's shape never reach recovery", func(t *testing.T) {
		f := newDeliveryFixture(t, uniqueTargetTable, true)
		effect := f.enqueue(t, DedupeNone, "record_id", t0, 7)
		// A foreign row with the same key makes the insert fail on a constraint.
		_, err := f.target.ExecContext(t.Context(), `INSERT INTO readings (record_id, counter) VALUES ('anything', 1)`)
		require.NoError(t, err)
		_, err = f.target.ExecContext(t.Context(), `CREATE TRIGGER refuse BEFORE INSERT ON readings BEGIN SELECT RAISE(ABORT, 'refused'); END`)
		require.NoError(t, err)
		recovery := &fakeRecovery{decision: RestoreDecision{Action: RestoreRetryNow}}
		newRecoveringSender(t, f, recovery)

		result, err := f.sender.Deliver(t.Context(), effect)
		require.NoError(t, err)
		require.Equal(t, StateQuarantined, result.State, "a rule of the destination that rejects the row is the row's fault")
		require.Zero(t, recovery.restores.Load(), "storage recovery was never asked")
	})

	t.Run("admission can block before the destination is touched", func(t *testing.T) {
		f := newDeliveryFixture(t, plainTargetTable, true)
		effect := f.enqueue(t, DedupeReceipt, "", t0, 7)
		recovery := &fakeRecovery{admit: "schema_repair_blocked"}
		newRecoveringSender(t, f, recovery)

		result, err := f.sender.Deliver(t.Context(), effect)
		require.NoError(t, err)
		require.Equal(t, [2]string{StateBlocked, "schema_repair_blocked"}, [2]string{result.State, result.Code})
		require.Zero(t, f.resolver.calls.Load(), "the destination was not opened")
		require.Equal(t, 0, f.targetRows(t))

		failing := &fakeRecovery{admitErr: errors.New("local database busy")}
		newRecoveringSender(t, f, failing)
		other := f.enqueue(t, DedupeReceipt, "", t0.Add(10*1e9), 8)
		_, err = f.sender.Deliver(t.Context(), other)
		require.Error(t, err, "an admission that cannot be answered does not guess")
		require.Equal(t, StatePending, f.state(t, other).State)
	})

	t.Run("a destination that cannot be opened for its identity is named by recovery", func(t *testing.T) {
		f := newDeliveryFixture(t, plainTargetTable, true)
		effect := f.enqueue(t, DedupeReceipt, "", t0, 7)
		f.resolver.err = ErrTargetBlocked
		recovery := &fakeRecovery{blocked: "schema_repair_blocked"}
		newRecoveringSender(t, f, recovery)
		result, err := f.sender.Deliver(t.Context(), effect)
		require.NoError(t, err)
		require.Equal(t, [2]string{StateBlocked, "schema_repair_blocked"}, [2]string{result.State, result.Code})

		plain := newDeliveryFixture(t, plainTargetTable, true)
		other := plain.enqueue(t, DedupeReceipt, "", t0, 7)
		plain.resolver.err = ErrTargetBlocked
		result, err = plain.sender.Deliver(t.Context(), other)
		require.NoError(t, err)
		require.Equal(t, codeTargetBlocked, result.Code, "without recovery the code is the delivery layer's own")
	})
}
