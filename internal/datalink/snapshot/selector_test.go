package snapshot

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

var (
	fixtureStart    = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	fixtureInterval = 10 * time.Second
)

func at(seconds int) time.Time { return fixtureStart.Add(time.Duration(seconds) * time.Second) }

func goodSample(id, member string, seconds int, value int64) Sample {
	return Sample{
		SampleID:   id,
		MemberKey:  member,
		ObservedAt: at(seconds),
		Quality:    schema.QualityGood,
		Value:      measurement.NewInt64(value),
	}
}

func newFixtureSelector(members ...string) *Selector {
	if len(members) == 0 {
		members = []string{"temperature"}
	}
	return NewSelector(SelectorConfig{
		Start:         fixtureStart,
		Interval:      fixtureInterval,
		MaxFutureSkew: 2 * time.Second,
		Members:       members,
	})
}

func TestUTCWindowBoundaryAndReorderingBucketStart(t *testing.T) {
	require.Equal(t, fixtureStart, BucketStart(at(0), fixtureInterval))
	require.Equal(t, fixtureStart, BucketStart(at(9).Add(999*time.Millisecond), fixtureInterval))
	require.Equal(t, at(10), BucketStart(at(10), fixtureInterval), "t=10 starts the next bucket")
	require.Equal(t, at(10), BucketStart(at(19), fixtureInterval))
	// Non-UTC input resolves to the same UTC bucket.
	zone := time.FixedZone("UTC+8", 8*3600)
	require.Equal(t, fixtureStart, BucketStart(at(3).In(zone), fixtureInterval))
	require.Equal(t, time.UTC, BucketStart(at(3).In(zone), fixtureInterval).Location())
	// Instants before the epoch still floor, never truncate toward zero.
	before := time.Date(1969, 12, 31, 23, 59, 55, 0, time.UTC)
	require.Equal(t, time.Date(1969, 12, 31, 23, 59, 50, 0, time.UTC), BucketStart(before, fixtureInterval))
}

func TestUTCWindowBoundaryAndReorderingOutOfOrderArrival(t *testing.T) {
	now := at(12)
	for name, order := range map[string][]Sample{
		"eight then three": {goodSample("s8", "temperature", 8, 80), goodSample("s3", "temperature", 3, 30)},
		"three then eight": {goodSample("s3", "temperature", 3, 30), goodSample("s8", "temperature", 8, 80)},
	} {
		selector := newFixtureSelector()
		for _, sample := range order {
			require.Contains(t, []OfferOutcome{OfferSelected, OfferNotSelected}, selector.Offer(sample, now).Outcome, name)
		}
		got, ok := selector.Selected("temperature")
		require.True(t, ok, name)
		require.Equal(t, "s8", got.SampleID, name)
		require.True(t, measurement.NewInt64(80).Equal(got.Value), name)
	}
}

func TestUTCWindowBoundaryAndReorderingOfferOutcomes(t *testing.T) {
	selector := newFixtureSelector()
	require.Equal(t, OfferSelected, selector.Offer(goodSample("s8", "temperature", 8, 80), at(12)).Outcome)
	require.Equal(t, OfferNotSelected, selector.Offer(goodSample("s3", "temperature", 3, 30), at(12)).Outcome)
	require.Equal(t, OfferSelected, selector.Offer(goodSample("s9", "temperature", 9, 90), at(12)).Outcome)
}

func TestUTCWindowBoundaryAndReorderingBoundarySampleBelongsToNextBucket(t *testing.T) {
	selector := newFixtureSelector()
	result := selector.Offer(goodSample("s10", "temperature", 10, 100), at(12))
	require.Equal(t, OfferRejected, result.Outcome)
	require.Equal(t, ReasonOutsideBucket, result.Reason)
	_, ok := selector.Selected("temperature")
	require.False(t, ok)

	next := NewSelector(SelectorConfig{Start: at(10), Interval: fixtureInterval, MaxFutureSkew: 2 * time.Second, Members: []string{"temperature"}})
	require.Equal(t, OfferSelected, next.Offer(goodSample("s10", "temperature", 10, 100), at(12)).Outcome)

	before := selector.Offer(goodSample("sm1", "temperature", -1, 1), at(12))
	require.Equal(t, OfferRejected, before.Outcome)
	require.Equal(t, ReasonOutsideBucket, before.Reason)
}

func TestUTCWindowBoundaryAndReorderingSameTimeUsesGreatestSampleID(t *testing.T) {
	samples := []Sample{
		goodSample("sample-a", "temperature", 5, 1),
		goodSample("sample-c", "temperature", 5, 3),
		goodSample("sample-b", "temperature", 5, 2),
	}
	rng := rand.New(rand.NewSource(7))
	for run := 0; run < 20; run++ {
		shuffled := append([]Sample(nil), samples...)
		rng.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
		selector := newFixtureSelector()
		for _, sample := range shuffled {
			selector.Offer(sample, at(8))
		}
		got, ok := selector.Selected("temperature")
		require.True(t, ok)
		require.Equal(t, "sample-c", got.SampleID, fmt.Sprintf("run %d order %v", run, shuffled))
	}
}

func TestUTCWindowBoundaryAndReorderingLatestObservationWinsEvenWhenBad(t *testing.T) {
	selector := newFixtureSelector()
	selector.Offer(goodSample("good", "temperature", 4, 40), at(12))
	bad := Sample{SampleID: "bad", MemberKey: "temperature", ObservedAt: at(9), Quality: schema.QualityBad, QualityReason: "read-failed"}
	require.Equal(t, OfferSelected, selector.Offer(bad, at(12)).Outcome)
	got, _ := selector.Selected("temperature")
	require.Equal(t, schema.QualityBad, got.Quality, "a later failed read must not be hidden by an older good value")
}

func TestUTCWindowBoundaryAndReorderingDuplicateSamples(t *testing.T) {
	selector := newFixtureSelector()
	original := goodSample("s8", "temperature", 8, 80)
	require.Equal(t, OfferSelected, selector.Offer(original, at(12)).Outcome)

	require.Equal(t, OfferDuplicate, selector.Offer(original, at(12)).Outcome, "same id and payload is a no-op")
	got, _ := selector.Selected("temperature")
	require.Equal(t, original, got)

	changed := original
	changed.Value = measurement.NewInt64(81)
	conflict := selector.Offer(changed, at(12))
	require.Equal(t, OfferConflict, conflict.Outcome)
	got, _ = selector.Selected("temperature")
	require.True(t, measurement.NewInt64(80).Equal(got.Value), "a conflicting resend must not change the selection")

	// A duplicate of a sample that lost selection is still a duplicate.
	older := goodSample("s3", "temperature", 3, 30)
	require.Equal(t, OfferNotSelected, selector.Offer(older, at(12)).Outcome)
	require.Equal(t, OfferDuplicate, selector.Offer(older, at(12)).Outcome)
	olderChanged := older
	olderChanged.ObservedAt = at(4)
	require.Equal(t, OfferConflict, selector.Offer(olderChanged, at(12)).Outcome)
}

func TestUTCWindowBoundaryAndReorderingFutureSkew(t *testing.T) {
	selector := newFixtureSelector()
	now := at(5)
	within := goodSample("within", "temperature", 7, 70) // exactly now + 2s skew
	require.Equal(t, OfferSelected, selector.Offer(within, now).Outcome)

	beyond := goodSample("beyond", "temperature", 8, 80) // now + 3s
	result := selector.Offer(beyond, now)
	require.Equal(t, OfferRejected, result.Outcome)
	require.Equal(t, ReasonFutureSkew, result.Reason)
	got, _ := selector.Selected("temperature")
	require.Equal(t, "within", got.SampleID)

	strict := NewSelector(SelectorConfig{Start: fixtureStart, Interval: fixtureInterval, Members: []string{"temperature"}})
	require.Equal(t, OfferRejected, strict.Offer(goodSample("f", "temperature", 6, 1), now).Outcome, "unset skew tolerates no future time")
	require.Equal(t, OfferSelected, strict.Offer(goodSample("p", "temperature", 5, 1), now).Outcome)
}

func TestUTCWindowBoundaryAndReorderingRejectsUnknownMemberAndMissingIdentity(t *testing.T) {
	selector := newFixtureSelector("temperature")
	result := selector.Offer(goodSample("x", "pressure", 3, 1), at(5))
	require.Equal(t, OfferRejected, result.Outcome)
	require.Equal(t, ReasonUnknownMember, result.Reason)

	result = selector.Offer(goodSample("", "temperature", 3, 1), at(5))
	require.Equal(t, OfferRejected, result.Outcome)
	require.Equal(t, ReasonMissingIdentity, result.Reason)
}
