package snapshot

import (
	"strings"
	"testing"

	"go-gateway/internal/datalink/counterproof"

	"github.com/stretchr/testify/require"
)

func TestSelectorCounterProofIdentityCannotChangeOrDowngrade(t *testing.T) {
	sample := goodSample("counter-1", "temperature", 2, 7)
	sample.CounterProof = counterproof.CounterBindingReference{Version: 1, Digest: strings.Repeat("a", 64)}
	selector := newFixtureSelector()
	require.Equal(t, OfferSelected, selector.Offer(sample, at(3)).Outcome)
	changed := sample
	changed.CounterProof.Digest = strings.Repeat("b", 64)
	require.Equal(t, OfferConflict, selector.Offer(changed, at(3)).Outcome)
	changed.CounterProof = counterproof.CounterBindingReference{}
	require.Equal(t, OfferConflict, selector.Offer(changed, at(3)).Outcome)
	require.Equal(t, OfferDuplicate, selector.Offer(sample, at(3)).Outcome)
}

func TestSelectorLegacyCounterResendDoesNotRelabelHistory(t *testing.T) {
	legacy := goodSample("counter-1", "temperature", 2, 7)
	selector := newFixtureSelector()
	require.Equal(t, OfferSelected, selector.Offer(legacy, at(3)).Outcome)
	resent := legacy
	resent.CounterProof = counterproof.CounterBindingReference{Version: 1, Digest: strings.Repeat("a", 64)}
	require.Equal(t, OfferDuplicate, selector.Offer(resent, at(3)).Outcome)
	selected, ok := selector.Selected("temperature")
	require.True(t, ok)
	require.True(t, selected.CounterProof.IsZero())
}
