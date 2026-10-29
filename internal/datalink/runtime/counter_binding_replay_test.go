package runtime

import (
	"fmt"
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

// These acceptance tests prime real durable state before refusing a new claim;
// recovery/provider substitution is tested separately when that path is wired.
func TestCounterBindingClaimsRefuseWithoutAdvancingAcceptedState(t *testing.T) {
	for _, cumulative := range []bool{true, false} {
		for _, scaled := range []bool{false, true} {
			t.Run(fmt.Sprintf("cumulative=%t/scaled=%t", cumulative, scaled), func(t *testing.T) {
				fixture := counterClaimsFixture(t, cumulative, scaled)
				f, store, key, boundary, interval, proof := fixture.f, fixture.store, fixture.key, fixture.boundary, fixture.interval, fixture.proof
				initial, err := store.Restore(t.Context(), key)
				require.NoError(t, err)
				require.Equal(t, at(interval), initial.NextClose)
				if cumulative {
					require.Len(t, initial.Baselines, 1)
				} else {
					require.Empty(t, initial.Baselines)
				}
				for _, field := range []string{"measurement_id", "series_epoch", "definition_revision", "source_binding_revision"} {
					for _, othersMatch := range []bool{false, true} {
						t.Run(fmt.Sprintf("%s/other_claims=%t", field, othersMatch), func(t *testing.T) {
							before, err := store.Restore(t.Context(), key)
							require.NoError(t, err)
							buckets, baseline := boundary.assembler.OpenBuckets(), boundary.assembler.CounterCheckpoints()
							count := countJournal(t, store.DB())
							id := fmt.Sprintf("%s-%t", field, othersMatch)
							envelope := f.envelope(0, id, at(interval+time.Second), int64(120004))
							if othersMatch {
								matchingCounterClaims(&envelope)
							}
							changeCounterClaim(&envelope, field, "foreign-secret-claim")
							err = boundary.AcceptSample(t.Context(), envelope)
							assertCounterClaimRefusal(t, err, field)
							require.Equal(t, count, countJournal(t, store.DB()))
							after, err := store.Restore(t.Context(), key)
							require.NoError(t, err)
							require.Equal(t, before, after, "no sample, baseline, or checkpoint change")
							require.Equal(t, buckets, boundary.assembler.OpenBuckets())
							require.Equal(t, baseline, boundary.assembler.CounterCheckpoints())
							require.Equal(t, before.NextClose, boundary.NextClose())
							changeCounterClaim(&envelope, field, map[string]string{
								"measurement_id": "meas-int", "series_epoch": "epoch-1",
								"definition_revision": "def-1", "source_binding_revision": "bind-1",
							}[field])
							require.NoError(t, boundary.AcceptSample(t.Context(), envelope), "a rejected ID must still accept its corrected sample")
							require.Equal(t, count+1, countJournal(t, store.DB()))
							changeCounterClaim(&envelope, field, "foreign-secret-claim")
							assertCounterClaimRefusal(t, boundary.AcceptSample(t.Context(), envelope), field)
							require.Equal(t, count+1, countJournal(t, store.DB()), "duplicate IDs cannot bypass the gate")
						})
					}
				}
				// Empty claims inherit this proof reference, not today's definition.
				empty := f.envelope(0, "empty-claims", at(interval+2*time.Second), int64(120004))
				require.NoError(t, boundary.AcceptSample(t.Context(), empty))
				restored, err := store.Restore(t.Context(), key)
				require.NoError(t, err)
				last := restored.Samples[len(restored.Samples)-1]
				require.Equal(t, proof.Digest, last.CounterProof.Digest)
				require.Equal(t, proof.Member, last.CounterProof.Member)
				if scaled {
					digits, ok := last.Value.DecimalDigits()
					require.True(t, ok)
					require.Equal(t, "12000.4", digits)
				} else {
					require.Equal(t, measurement.ExactInt64, last.Value.Type())
					require.Equal(t, int64(120004), last.Value.Value())
				}
				require.Equal(t, initial.Baselines, restored.Baselines)
				require.Equal(t, initial.NextClose, restored.NextClose)
			})
		}
	}
}

func assertCounterClaimRefusal(t *testing.T, err error, field string) {
	t.Helper()
	require.ErrorIs(t, err, ErrCounterSampleClaimConflict)
	var scoped *CounterSampleClaimError
	require.ErrorAs(t, err, &scoped)
	require.Equal(t, field, scoped.Claim)
	require.Equal(t, `["workspace-A","group-L","rev-1","device-A","point-int","tag-int"]`, scoped.Member.Key())
	var refused *GroupSampleError
	require.ErrorAs(t, err, &refused)
	require.Equal(t, snapshot.OfferRejected, refused.Outcome)
	require.Equal(t, "counter-sample-claim-conflict", refused.Reason)
	require.NotContains(t, err.Error(), "foreign-secret-claim")
}
