package grouppipeline

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"go-gateway/internal/datalink/counterproof"

	"github.com/stretchr/testify/require"
)

func TestCounterDuplicateRefusesDamagedDurableProof(t *testing.T) {
	for _, reason := range []error{counterproof.ErrCounterBindingProofMissing, counterproof.ErrCounterBindingProofConflict, counterproof.ErrCounterBindingProofUnsupported} {
		t.Run(reason.Error(), func(t *testing.T) {
			f := newLongPipeline(t)
			require.NoError(t, f.reading(0, "accepted", at(1), 100))
			damageCounterProof(t, f, reason.Error())
			before := counterDurableState(t, f.store)

			err := f.reading(0, "accepted", at(1), 100)
			require.ErrorIs(t, err, reason)
			entries, err := f.pipe.CounterReplayQuarantines(t.Context(), f.group.ID)
			require.NoError(t, err)
			require.Len(t, entries, 1)
			require.Equal(t, 1, entries[0].PendingSamples)
			require.Equal(t, before, counterDurableState(t, f.store))
		})
	}
}

func TestCounterClosureRetainsPendingInputWhenProofIsDamaged(t *testing.T) {
	for _, reason := range []error{counterproof.ErrCounterBindingProofMissing, counterproof.ErrCounterBindingProofConflict, counterproof.ErrCounterBindingProofUnsupported} {
		t.Run(reason.Error(), func(t *testing.T) {
			f := newLongPipeline(t)
			require.NoError(t, f.reading(0, "accepted", at(1), 100))
			damageCounterProof(t, f, reason.Error())
			before := counterDurableState(t, f.store)

			f.closeUntil(at(600))
			count, err := f.store.PendingSampleCount(t.Context(), groupKey(f.snapshot))
			require.NoError(t, err)
			require.Equal(t, 1, count, "closure must not consume an input whose durable proof cannot be verified")
			entries, err := f.pipe.CounterReplayQuarantines(t.Context(), f.group.ID)
			require.NoError(t, err)
			require.Len(t, entries, 1)
			require.Equal(t, reason.Error(), entries[0].Reason)
			require.Equal(t, before, counterDurableState(t, f.store))
		})
	}
}

func TestCounterDuplicateDoesNotRecreateMissingJournal(t *testing.T) {
	f := newLongPipeline(t)
	require.NoError(t, f.reading(0, "accepted", at(1), 100))
	_, err := f.local.ExecContext(t.Context(), `DELETE FROM wg_delivery_samples WHERE sample_id=?`, "accepted")
	require.NoError(t, err)
	before := counterDurableState(t, f.store)

	err = f.reading(0, "accepted", at(1), 100)
	require.ErrorIs(t, err, counterproof.ErrCounterBindingProofMissing)
	entries, err := f.pipe.CounterReplayQuarantines(t.Context(), f.group.ID)
	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Equal(t, before, counterDurableState(t, f.store))
}

func TestCounterClosureProofReadFailureRetainsPendingInput(t *testing.T) {
	f := newLongPipeline(t)
	require.NoError(t, f.reading(0, "accepted", at(1), 100))
	before := counterDurableState(t, f.store)
	injected := errors.New("proof storage unavailable")
	reject := func(query string) error {
		if strings.Contains(query, "SELECT payload FROM wg_runtime_versions") {
			return injected
		}
		return nil
	}
	f.localFault.Reject.Store(&reject)
	f.closeUntil(at(600))
	f.localFault.Reject.Store(nil)

	require.Len(t, f.errs, 1)
	require.ErrorIs(t, f.errs[0], injected)
	entries, err := f.pipe.CounterReplayQuarantines(t.Context(), f.group.ID)
	require.NoError(t, err)
	require.Empty(t, entries, "an infrastructure failure is not corrupt proof")
	require.Equal(t, before, counterDurableState(t, f.store))
}

func TestCounterClosureRetainsInputWhenJournalProofIsDamaged(t *testing.T) {
	f := newLongPipeline(t)
	require.NoError(t, f.reading(0, "accepted", at(1), 100))
	var payload string
	require.NoError(t, f.local.QueryRowContext(t.Context(), `SELECT payload FROM wg_delivery_samples WHERE sample_id=?`, "accepted").Scan(&payload))
	var sample map[string]any
	require.NoError(t, json.Unmarshal([]byte(payload), &sample))
	sample["counter_proof"].(map[string]any)["digest"] = strings.Repeat("0", 64)
	damaged, err := json.Marshal(sample)
	require.NoError(t, err)
	_, err = f.local.ExecContext(t.Context(), `UPDATE wg_delivery_samples SET payload=? WHERE sample_id=?`, string(damaged), "accepted")
	require.NoError(t, err)
	before := counterDurableState(t, f.store)

	f.closeUntil(at(600))
	count, err := f.store.PendingSampleCount(t.Context(), groupKey(f.snapshot))
	require.NoError(t, err)
	require.Equal(t, 1, count, "closure must retain an input whose journal proof was damaged")
	entries, err := f.pipe.CounterReplayQuarantines(t.Context(), f.group.ID)
	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Equal(t, counterproof.ErrCounterBindingProofConflict.Error(), entries[0].Reason)
	require.Equal(t, before, counterDurableState(t, f.store))
}

func TestCounterClosureRetainsLegacyInputWhenAssociationIsMissing(t *testing.T) {
	f := legacyQuarantineFixture(t)
	f.reconcileStatus()
	result, err := f.local.ExecContext(t.Context(), `DELETE FROM wg_counter_proof_associations WHERE sample_id=?`, "legacy-original")
	require.NoError(t, err)
	deleted, err := result.RowsAffected()
	require.NoError(t, err)
	require.EqualValues(t, 1, deleted)
	before := counterDurableState(t, f.store)

	f.closeUntil(at(600))
	count, err := f.store.PendingSampleCount(t.Context(), groupKey(f.snapshot))
	require.NoError(t, err)
	require.Equal(t, 1, count, "closure must not recreate or bypass the original association")
	entries, err := f.pipe.CounterReplayQuarantines(t.Context(), f.group.ID)
	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Equal(t, counterproof.ErrCounterBindingProofMissing.Error(), entries[0].Reason)
	require.Equal(t, before, counterDurableState(t, f.store))
}

func TestCounterLegacyDuplicateRefusesMissingAssociation(t *testing.T) {
	f := legacyQuarantineFixture(t)
	f.reconcileStatus()
	result, err := f.local.ExecContext(t.Context(), `DELETE FROM wg_counter_proof_associations WHERE sample_id=?`, "legacy-original")
	require.NoError(t, err)
	deleted, err := result.RowsAffected()
	require.NoError(t, err)
	require.EqualValues(t, 1, deleted)
	before := counterDurableState(t, f.store)

	err = f.reading(0, "legacy-original", at(1), 1000)
	require.ErrorIs(t, err, counterproof.ErrCounterBindingProofMissing)
	entries, err := f.pipe.CounterReplayQuarantines(t.Context(), f.group.ID)
	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Equal(t, before, counterDurableState(t, f.store))
}
