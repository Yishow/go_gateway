package snapshot

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func boundsAssembler(t *testing.T, mode Mode, first, start, end time.Time) *Assembler {
	t.Helper()
	c := fixtureConfig()
	c.Mode, c.FirstBucket, c.StartsAt, c.EndsAt = mode, first, start, end
	c.Members = []Member{{Key: "temperature", Required: mode == ModeWide}}
	c.Storage = StorageCapabilities{NullableValues: true, MemberQuality: true}
	if !end.IsZero() {
		c.Until = BucketStart(end, c.Interval)
		if !c.Until.Equal(end) {
			c.Until = c.Until.Add(c.Interval)
		}
	}
	a, err := NewAssembler(c)
	require.NoError(t, err)
	return a
}

func TestActualBoundsRebuildRestoresEarlierLatestCandidate(t *testing.T) {
	for _, mode := range []Mode{ModeWide, ModeInstantaneousProfile, ModeWideProfile} {
		t.Run(string(mode), func(t *testing.T) {
			a := boundsAssembler(t, mode, at(20), time.Time{}, time.Time{})
			before := goodSample("before", "temperature", 22, 6)
			before.ObservedAt = at(22).Add(499 * time.Millisecond)
			exact := goodSample("exact", "temperature", 22, 8)
			exact.ObservedAt = at(22).Add(500 * time.Millisecond)
			after := goodSample("after", "temperature", 23, 7)
			for _, sample := range []Sample{after, exact, before, goodSample("earlier", "temperature", 21, 5)} {
				require.Contains(t, []OfferOutcome{OfferSelected, OfferNotSelected}, a.Offer(sample, at(23)).Outcome)
			}
			a.SetEndsAt(exact.ObservedAt)
			plan := a.Plan(at(30))
			require.Len(t, plan.Outcomes, 1)
			result := memberResult(t, plan.Outcomes[0], "temperature")
			require.Equal(t, before, *result.Sample, "actual end must rebuild from all retained candidates, not keep or delete only chosen latest")
			require.Equal(t, exact.ObservedAt, plan.Outcomes[0].WindowEnd)
			require.Equal(t, at(20), a.NextClose(), "rebuild cannot consume or checkpoint")
			require.Equal(t, OfferSelected, a.Offer(after, at(23)).Outcome, "pure assembler keeps its established buffer contract, not a durable ACK")
			require.Equal(t, before, *memberResult(t, a.Plan(at(30)).Outcomes[0], "temperature").Sample, "closure projection cannot reintroduce a successor observation")
			a.SetEndsAt(at(22))
			a.SetEndsAt(at(24)) // An end never extends again.
			require.Equal(t, "earlier", memberResult(t, a.Plan(at(30)).Outcomes[0], "temperature").Sample.SampleID)
		})
	}
}

func TestActualBoundsNewOwnerKeepsExactBoundaryAndLater(t *testing.T) {
	boundary := at(22).Add(500 * time.Millisecond)
	a := boundsAssembler(t, ModeInstantaneousProfile, at(20), boundary, at(30))
	before := goodSample("before", "temperature", 22, 6)
	before.ObservedAt = boundary.Add(-time.Nanosecond)
	exact := before
	exact.SampleID, exact.ObservedAt = "exact", boundary
	require.Equal(t, OfferSelected, a.Replay(before).Outcome)
	require.Nil(t, memberResult(t, a.Plan(at(40)).Outcomes[0], "temperature").Sample, "before-start input cannot contribute")
	require.Equal(t, OfferDuplicate, a.Check(before, before.ObservedAt).Outcome, "Plan projection must not mutate the standalone accepted buffer")
	require.Equal(t, OfferSelected, a.Replay(exact).Outcome)
	require.Equal(t, OfferSelected, a.Replay(goodSample("after", "temperature", 23, 7)).Outcome)
	require.Equal(t, "after", memberResult(t, a.Plan(at(40)).Outcomes[0], "temperature").Sample.SampleID)
}

func TestActualBoundsPreservesSameIDOriginalOriginsAndNoBoundsBehavior(t *testing.T) {
	a := boundsAssembler(t, ModeWide, at(0), time.Time{}, time.Time{})
	first := goodSample("same", "temperature", 8, 10)
	first.AcceptanceOrigin = AcceptanceOrigin{"workspace", "original", "rev"}
	second := goodSample("same", "temperature", 9, 20)
	future := goodSample("future", "temperature", 11, 100)
	for _, sample := range []Sample{future, second, first} {
		require.Contains(t, []OfferOutcome{OfferSelected, OfferNotSelected}, a.Offer(sample, at(12)).Outcome)
	}
	require.Equal(t, 2, a.OpenBuckets())
	a.SetEndsAt(at(10))
	require.Equal(t, 1, a.OpenBuckets(), "successor-only buckets are no longer old-owner eligible memory")
	selector := a.buckets[at(0)].selectors[""]
	require.Equal(t, []Sample{first, second}, selector.Accepted("temperature"))
	require.Equal(t, second, *memberResult(t, a.Plan(at(20)).Outcomes[0], "temperature").Sample)
	require.Equal(t, OfferDuplicate, a.Offer(second, at(12)).Outcome)
	require.Zero(t, a.Late())
}

func TestActualBoundsLatestFreshnessUsesActualEnd(t *testing.T) {
	for _, mode := range []Mode{ModeWide, ModeInstantaneousProfile, ModeWideProfile} {
		t.Run(string(mode), func(t *testing.T) {
			for _, tc := range []struct {
				name     string
				observed time.Time
				status   MemberStatus
			}{
				{"fresh", at(22).Add(499 * time.Millisecond), MemberOK},
				{"exact max age", at(20).Add(500 * time.Millisecond), MemberOK},
				{"truly stale", at(20), MemberStale},
			} {
				t.Run(tc.name, func(t *testing.T) {
					c := fixtureConfig()
					c.Mode, c.FirstBucket = mode, at(20)
					c.Members = []Member{{Key: "temperature", Required: mode == ModeWide, MaxAge: 2 * time.Second}}
					c.Storage = StorageCapabilities{NullableValues: true, MemberQuality: true}
					a, err := NewAssembler(c)
					require.NoError(t, err)
					sample := goodSample("candidate", "temperature", 22, 6)
					sample.ObservedAt = tc.observed
					require.Equal(t, OfferSelected, a.Offer(sample, at(23)).Outcome)
					require.Equal(t, OfferSelected, a.Offer(goodSample("cross-boundary", "temperature", 23, 7), at(23)).Outcome)
					a.SetEndsAt(at(22).Add(500 * time.Millisecond))
					plan := a.Plan(at(30))
					require.Len(t, plan.Outcomes, 1)
					result := memberResult(t, plan.Outcomes[0], "temperature")
					require.Equal(t, tc.status, result.Status, "freshness is measured at actual22.5, not aligned30")
					require.Equal(t, sample, *result.Sample)
				})
			}
		})
	}
}
