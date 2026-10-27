package runtime

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestCounterBindingClaimsRequireCompleteFrozenIdentity(t *testing.T) {
	for name, remove := range map[string]func(*counterClaimIdentity){
		"measurement": func(c *counterClaimIdentity) { c.measurementID = "" },
		"definition":  func(c *counterClaimIdentity) { c.definitionRevision = "" },
		"binding":     func(c *counterClaimIdentity) { c.sourceBindingRevision = "" },
		"epoch":       func(c *counterClaimIdentity) { c.seriesEpoch = "" },
	} {
		t.Run(name, func(t *testing.T) {
			f, store, key, boundary, interval, proof := counterClaimsFixture(t, true, false)
			memberKey := proof.Member.SampleMemberKey()
			member := boundary.members[memberKey]
			remove(&member.counterClaims)
			boundary.members[memberKey] = member
			before, err := store.Restore(t.Context(), key)
			require.NoError(t, err)
			err = boundary.AcceptSample(t.Context(), f.envelope(0, "empty-without-proof", at(interval+time.Second), int64(120004)))
			require.ErrorIs(t, err, counterproof.ErrCounterBindingProofMissing, "empty claims cannot invent missing applied identity")
			after, err := store.Restore(t.Context(), key)
			require.NoError(t, err)
			require.Equal(t, before, after)
			require.Zero(t, boundary.assembler.OpenBuckets())
		})
	}
}

func TestCounterBindingClaimsCopyAppliedIdentityNotCallerConfig(t *testing.T) {
	f, store, key, boundary, interval, original := counterClaimsFixture(t, true, false)
	facts := original.CounterBindingFacts
	facts.DefinitionRevision, facts.SeriesEpoch = "def-2", "epoch-2"
	changed, err := workspace.NewCounterBindingProof(facts)
	require.NoError(t, err)
	f.config.CounterBindings[original.Member.Key()] = changed
	source := f.config.Counters["tag-int"]
	source.DefinitionRevision, source.SeriesEpoch = "def-2", "epoch-2"
	f.config.Counters["tag-int"] = source
	envelope := f.envelope(0, "still-applied", at(interval+time.Second), int64(120004))
	require.NoError(t, boundary.AcceptSample(t.Context(), envelope))
	restored, err := store.Restore(t.Context(), key)
	require.NoError(t, err)
	require.Len(t, restored.Samples, 1)
	require.Equal(t, original.Digest, restored.Samples[0].CounterProof.Digest)
	envelope.SeriesEpoch = "epoch-2"
	assertCounterClaimRefusal(t, boundary.AcceptSample(t.Context(), envelope), "series_epoch")
	require.Equal(t, 2, countJournal(t, store.DB()), "only baseline plus original-proof sample were accepted")
}
