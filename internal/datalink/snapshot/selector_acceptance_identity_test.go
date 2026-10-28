package snapshot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSelectorAcceptedOriginsKeepOriginalIDDistinct(t *testing.T) {
	for _, reverse := range []bool{false, true} {
		s := newFixtureSelector()
		transferred := goodSample("s-1", "temperature", 3, 7)
		transferred.AcceptanceOrigin = AcceptanceOrigin{"workspace", "A", "revision-A"}
		native := goodSample("s-1", "temperature", 4, 9)
		inputs := []Sample{transferred, native}
		if reverse {
			inputs[0], inputs[1] = inputs[1], inputs[0]
		}
		for _, sample := range inputs {
			require.Contains(t, []OfferOutcome{OfferSelected, OfferNotSelected}, s.Offer(sample, at(9)).Outcome)
		}
		require.Equal(t, []Sample{transferred, native}, s.Accepted("temperature"))
		require.Len(t, s.counterReadings("temperature"), 2, "cumulative inputs cannot collapse original IDs across origins")
		for _, sample := range inputs {
			require.Equal(t, OfferDuplicate, s.Offer(sample, at(9)).Outcome)
		}
		changed := native
		changed.ObservedAt = at(5)
		require.Equal(t, OfferConflict, s.Offer(changed, at(9)).Outcome, "native same-ID conflict remains unchanged")
		changed = transferred
		changed.ObservedAt = at(5)
		require.Equal(t, OfferConflict, s.Offer(changed, at(9)).Outcome, "same original acceptance cannot select a new observation")
	}
}

func TestSelectorMultipleTransferredOriginsDoNotCollapse(t *testing.T) {
	s := newFixtureSelector()
	first := goodSample("s-1", "temperature", 3, 7)
	first.AcceptanceOrigin = AcceptanceOrigin{"workspace", "A", "revision-A"}
	second := goodSample("s-1", "temperature", 4, 9)
	second.AcceptanceOrigin = AcceptanceOrigin{"workspace", "B", "revision-B"}
	require.Equal(t, OfferSelected, s.Offer(first, at(9)).Outcome)
	require.Equal(t, OfferSelected, s.Offer(second, at(9)).Outcome)
	require.Len(t, s.Accepted("temperature"), 2)
	require.Len(t, s.counterReadings("temperature"), 2)
}

func TestSampleAcceptanceOriginCannotBeSelectedThroughJSON(t *testing.T) {
	sample := goodSample("s-1", "temperature", 3, 7)
	sample.AcceptanceOrigin = AcceptanceOrigin{"workspace", "A", "revision-A"}
	raw, err := json.Marshal(sample)
	require.NoError(t, err)
	require.NotContains(t, string(raw), "AcceptanceOrigin")
	var decoded Sample
	require.NoError(t, json.Unmarshal([]byte(`{"SampleID":"s-1","AcceptanceOrigin":{"WorkspaceID":"forged","GroupID":"A","GroupRevision":"rev"}}`), &decoded))
	require.Equal(t, AcceptanceOrigin{}, decoded.AcceptanceOrigin)
}
