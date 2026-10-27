package runtime

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestCounterBindingReplayRejectsReferenceFromDifferentFrozenMeaning(t *testing.T) {
	f, store, key, payload := counterAcceptanceFixture(t)
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), key, payload))
	original := f.build(t)
	require.NoError(t, original.AcceptSample(t.Context(), f.envelope(0, "original-1", at(time.Second), int64(100))))
	before, err := store.Restore(t.Context(), key)
	require.NoError(t, err)
	// The replacement is internally valid and agrees with its source adapter,
	// but cannot reinterpret the already accepted original proof reference.
	for key, proof := range f.config.CounterBindings {
		facts := proof.CounterBindingFacts
		facts.DefinitionRevision = "def-2"
		facts.SeriesEpoch = "epoch-2"
		facts.Unit = "MWh"
		f.config.CounterBindings[key], err = workspace.NewCounterBindingProof(facts)
		require.NoError(t, err)
	}
	f.config.Counters = nil
	boundary, err := NewGroupBoundary(f.config)
	require.Nil(t, boundary)
	require.ErrorIs(t, err, counterproof.ErrCounterBindingProofConflict)
	after, err := store.Restore(t.Context(), key)
	require.NoError(t, err)
	require.Equal(t, before, after)
}

func TestCounterBindingDurableActivationNeverUsesUnprovenCounterConfig(t *testing.T) {
	f, store, key, payload := counterAcceptanceFixture(t)
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), key, payload))
	f.config.CounterBindings = nil
	boundary, err := NewGroupBoundary(f.config)
	require.Nil(t, boundary)
	require.ErrorIs(t, err, counterproof.ErrCounterBindingProofMissing)
	require.Zero(t, countJournal(t, store.DB()))
	// The pre-existing no-ledger utility is explicitly non-durable. It is not
	// how the sole production pipeline caller starts a group.
	f.config.Ledger = nil
	boundary, err = NewGroupBoundary(f.config)
	require.NoError(t, err)
	require.NotNil(t, boundary)
}
