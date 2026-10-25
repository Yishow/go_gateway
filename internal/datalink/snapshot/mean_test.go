package snapshot

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func meanConfig(mutate func(*Config)) Config {
	config := fixtureConfig()
	config.Mode = ModeInstantaneousProfile
	config.Members = []Member{{Key: "alpha", SampleReduction: SampleReductionMean}}
	if mutate != nil {
		mutate(&config)
	}
	return config
}

func newMeanAssembler(t *testing.T, mutate func(*Config)) *Assembler {
	t.Helper()
	assembler, err := NewAssembler(meanConfig(mutate))
	require.NoError(t, err)
	return assembler
}

// newLatestAssembler keeps the same one-member profile shape on the existing
// latest rule, so both reductions can be compared on identical input.
func newLatestAssembler(t *testing.T) *Assembler {
	t.Helper()
	assembler, err := NewAssembler(meanConfig(func(c *Config) { c.Members[0].SampleReduction = SampleReductionLatest }))
	require.NoError(t, err)
	return assembler
}

func acceptedSample(id string, seconds int, value measurement.ExactValue) Sample {
	return Sample{SampleID: id, MemberKey: "alpha", ObservedAt: at(seconds), Quality: schema.QualityGood, Value: value}
}

func intSample(id string, seconds int, value int64) Sample {
	return acceptedSample(id, seconds, measurement.NewInt64(value))
}

func accept(t *testing.T, a *Assembler, samples ...Sample) {
	t.Helper()
	for _, sample := range samples {
		require.Contains(t, []OfferOutcome{OfferSelected, OfferNotSelected}, a.Offer(sample, at(10)).Outcome, sample.SampleID)
	}
}

func closeOneWindow(t *testing.T, a *Assembler, now int) MemberResult {
	t.Helper()
	outcomes := a.Tick(at(now))
	require.Len(t, outcomes, 1, "one member closes one record per window")
	require.Len(t, outcomes[0].Members, 1)
	return outcomes[0].Members[0]
}

// exactMean fails the test and returns nil when the member carries no reduction,
// so a RED run reports the missing result instead of dereferencing nil.
func exactMean(t *testing.T, member MemberResult) *measurement.MeanResult {
	t.Helper()
	require.NotNil(t, member.Reduction, "a mean member carries an independent reduction result")
	return member.Reduction
}

func digitsOfExact(t *testing.T, value measurement.ExactValue) string {
	t.Helper()
	digits, ok := value.DecimalDigits()
	require.True(t, ok)
	return digits
}

func mustFloatValue(t *testing.T, value float64) measurement.ExactValue {
	t.Helper()
	exact, err := measurement.NewFloat64(value)
	require.NoError(t, err)
	return exact
}

// TestWindowArithmeticMean pins the window reducer: equal weight for every
// unique eligible observation inside the real window, no interpolation, no
// filled value, and the newest-reading result untouched beside it.
func TestWindowArithmeticMean(t *testing.T) {
	t.Run("10, 20 and 30 average to 20 while the newest reading stays 30", func(t *testing.T) {
		a := newMeanAssembler(t, nil)
		accept(t, a, intSample("a10", 1, 10), intSample("a20", 2, 20), intSample("a30", 3, 30))
		require.Equal(t, OfferDuplicate, a.Offer(intSample("a10", 1, 10), at(4)).Outcome, "a retransmission is not counted twice")
		bad := intSample("bad", 4, 1000)
		bad.Quality, bad.QualityReason = schema.QualityBad, "read-timeout"
		accept(t, a, bad)

		member := closeOneWindow(t, a, 10)
		reduction := exactMean(t, member)
		if reduction != nil {
			require.Equal(t, "20", digitsOfExact(t, reduction.Value))
			require.Equal(t, int64(3), reduction.Count)
			require.Equal(t, "60", digitsOfExact(t, reduction.Sum))
			require.Equal(t, measurement.MeanBasisExact, reduction.Basis)
			require.Equal(t, MemberOK, member.Status)
			require.Equal(t, "a30", member.Sample.SampleID, "the newest eligible reading is still reported")
			require.False(t, reduction.Value.Equal(member.Sample.Value), "the mean is not passed off as a source reading")
		}
	})

	t.Run("arrival order never changes the mean", func(t *testing.T) {
		samples := []Sample{intSample("a10", 1, 10), intSample("a20", 2, 20), intSample("a30", 3, 30)}
		forward, backward := newMeanAssembler(t, nil), newMeanAssembler(t, nil)
		accept(t, forward, samples...)
		accept(t, backward, samples[2], samples[0], samples[1])
		first, second := closeOneWindow(t, forward, 10), closeOneWindow(t, backward, 10)
		if exactMean(t, first) != nil {
			require.Equal(t, first.Reduction, second.Reduction)
			require.Equal(t, "20", digitsOfExact(t, first.Reduction.Value))
		}
	})

	t.Run("a conflicting identity is refused and never counted", func(t *testing.T) {
		a := newMeanAssembler(t, nil)
		accept(t, a, intSample("a10", 1, 10), intSample("a20", 2, 20))
		require.Equal(t, OfferConflict, a.Offer(intSample("a10", 1, 999), at(3)).Outcome)

		member := closeOneWindow(t, a, 10)
		if reduction := exactMean(t, member); reduction != nil {
			require.Equal(t, "15", digitsOfExact(t, reduction.Value))
			require.Equal(t, int64(2), reduction.Count)
		}
	})

	t.Run("an observation at the window end belongs to the next window", func(t *testing.T) {
		a := newMeanAssembler(t, nil)
		accept(t, a, intSample("a1", 1, 10), intSample("a9", 9, 20), intSample("a10", 10, 30))
		first := closeOneWindow(t, a, 10)
		if reduction := exactMean(t, first); reduction != nil {
			require.Equal(t, "15", digitsOfExact(t, reduction.Value))
			require.Equal(t, int64(2), reduction.Count)
		}
		second := closeOneWindow(t, a, 20)
		if reduction := exactMean(t, second); reduction != nil {
			require.Equal(t, "30", digitsOfExact(t, reduction.Value))
			require.Equal(t, int64(1), reduction.Count)
		}
	})

	t.Run("stale and bad readings do not enter the denominator", func(t *testing.T) {
		a := newMeanAssembler(t, func(c *Config) { c.Members[0].MaxAge = 2 * time.Second })
		accept(t, a, intSample("old", 1, 10), intSample("a8", 8, 20), intSample("a9", 9, 30))
		member := closeOneWindow(t, a, 10)
		if reduction := exactMean(t, member); reduction != nil {
			require.Equal(t, "25", digitsOfExact(t, reduction.Value), "the 9s-old reading is stale at the window end")
			require.Equal(t, int64(2), reduction.Count)
			require.Equal(t, "a9", member.Sample.SampleID)
		}
	})

	t.Run("a cutover start filters to the real window and keeps the aligned identity", func(t *testing.T) {
		cut := newMeanAssembler(t, func(c *Config) { c.StartsAt = at(7) })
		accept(t, cut, intSample("before", 1, 100), intSample("a8", 8, 200), intSample("a9", 9, 300))
		outcomes := cut.Tick(at(10))
		require.Len(t, outcomes, 1)
		member := outcomes[0].Members[0]
		if reduction := exactMean(t, member); reduction != nil {
			require.Equal(t, "250", digitsOfExact(t, reduction.Value), "only observations inside [7s, 10s) contribute")
			require.Equal(t, int64(2), reduction.Count)
		}
		require.Equal(t, at(0), outcomes[0].BucketStart, "record identity keeps the aligned bucket")
		require.Equal(t, at(7), outcomes[0].WindowStart)
		require.Equal(t, at(10), outcomes[0].WindowEnd)
		require.Equal(t, "start", outcomes[0].Truncated)
	})

	t.Run("a cutover end filters to the real window", func(t *testing.T) {
		cut := newMeanAssembler(t, func(c *Config) { c.EndsAt, c.Until = at(8), at(10) })
		accept(t, cut, intSample("a3", 3, 10), intSample("a7", 7, 20), intSample("a9", 9, 30))
		outcomes := cut.Tick(at(10))
		require.Len(t, outcomes, 1)
		member := outcomes[0].Members[0]
		if reduction := exactMean(t, member); reduction != nil {
			require.Equal(t, "15", digitsOfExact(t, reduction.Value), "only observations inside [0s, 8s) contribute")
			require.Equal(t, int64(2), reduction.Count)
		}
		require.Equal(t, at(8), outcomes[0].WindowEnd)
		require.Equal(t, "end", outcomes[0].Truncated)
	})

	t.Run("a mean window with no eligible observation in the real window is not usable", func(t *testing.T) {
		cut := newMeanAssembler(t, func(c *Config) { c.StartsAt = at(7) })
		// The only good reading sits before the real start, so the real window has
		// no observation at all: the raw reading must not become the window value.
		accept(t, cut, intSample("before", 1, 42))
		outcomes := cut.Tick(at(10))
		require.Len(t, outcomes, 1)
		member := outcomes[0].Members[0]
		require.Nil(t, member.Reduction)
		require.Nil(t, member.Sample, "there is no observation inside the real window")
		require.Equal(t, MemberMissing, member.Status)
		require.Equal(t, ReasonNoSamples, member.Reason)
		require.False(t, member.Usable(), "a window without a mean is never written from the raw reading")
	})

	t.Run("a mean window reports only observations inside the real window", func(t *testing.T) {
		cut := newMeanAssembler(t, func(c *Config) { c.EndsAt, c.Until = at(8), at(10) })
		// The good reading at 9s is outside the real [0s, 8s) window; it must not
		// become the reported sample even though the nominal bucket holds it.
		accept(t, cut, intSample("a3", 3, 10), intSample("a7", 7, 20))
		require.Equal(t, OfferSelected, cut.Offer(intSample("a9", 9, 30), at(10)).Outcome)
		outcomes := cut.Tick(at(10))
		require.Len(t, outcomes, 1)
		member := outcomes[0].Members[0]
		if reduction := exactMean(t, member); reduction != nil {
			require.Equal(t, "15", digitsOfExact(t, reduction.Value))
			require.Equal(t, int64(2), reduction.Count)
		}
		require.Equal(t, MemberOK, member.Status)
		require.Equal(t, "a7", member.Sample.SampleID, "the reported reading comes from the real window")
	})

	t.Run("an unusable real window keeps its actual evidence", func(t *testing.T) {
		cut := newMeanAssembler(t, func(c *Config) { c.StartsAt = at(7) })
		bad := intSample("bad8", 8, 1000)
		bad.Quality, bad.QualityReason = schema.QualityBad, "read-timeout"
		accept(t, cut, bad)
		member := closeOneWindow(t, cut, 10)
		require.Nil(t, member.Reduction)
		require.Equal(t, MemberBad, member.Status)
		require.Equal(t, "bad8", member.Sample.SampleID)
		require.False(t, member.Usable())
	})

	t.Run("an ineligible reading outside the real window cannot make the mean unusable", func(t *testing.T) {
		cut := newMeanAssembler(t, func(c *Config) { c.EndsAt, c.Until = at(8), at(10) })
		accept(t, cut, intSample("a3", 3, 10), intSample("a7", 7, 20))
		bad := intSample("bad9", 9, 30)
		bad.Quality, bad.QualityReason = schema.QualityBad, "read-timeout"
		accept(t, cut, bad)
		member := closeOneWindow(t, cut, 10)
		if reduction := exactMean(t, member); reduction != nil {
			require.Equal(t, "15", digitsOfExact(t, reduction.Value))
		}
		require.Equal(t, MemberOK, member.Status)
		require.Equal(t, "a7", member.Sample.SampleID)
	})

	t.Run("a float window follows observed_at then sample_id, not arrival or value", func(t *testing.T) {
		build := func(order []Sample) *measurement.MeanResult {
			a := newMeanAssembler(t, nil)
			accept(t, a, order...)
			return exactMean(t, closeOneWindow(t, a, 10))
		}
		// The t=2 tie orders "a" before "b" by sample_id, whatever the value, so
		// the float sum is 1e16 + 1 - 1e16 + 1 = 1 and the mean is 0.25. Ordering
		// by value or by arrival would produce different sums.
		samples := []Sample{
			acceptedSample("m", 1, mustFloatValue(t, 1e16)),
			acceptedSample("a", 2, mustFloatValue(t, 1)),
			acceptedSample("b", 2, mustFloatValue(t, -1e16)),
			acceptedSample("n", 3, mustFloatValue(t, 1)),
		}
		reversed := []Sample{samples[3], samples[2], samples[1], samples[0]}
		forward, backward := build(samples), build(reversed)
		if forward != nil {
			require.Equal(t, forward, backward, "arrival order never changes the float mean")
			require.Equal(t, measurement.MeanBasisFloat64, forward.Basis)
			require.Equal(t, measurement.ExactFloat64, forward.Value.Type())
			require.Equal(t, 1.0, forward.Sum.Value())
			require.Equal(t, 0.25, forward.Value.Value())
			require.Equal(t, int64(4), forward.Count)
		}
	})

	t.Run("the 15 second contract window averages 10, 20 and 30 to 20", func(t *testing.T) {
		a := newMeanAssembler(t, func(c *Config) { c.Interval = 15 * time.Second })
		accept(t, a, intSample("a3", 3, 10), intSample("a6", 6, 20), intSample("a9", 9, 30))
		require.Equal(t, OfferSelected, a.Offer(intSample("a15", 15, 40), at(15)).Outcome, "the sample at the window end opens the next window")
		first := closeOneWindow(t, a, 15)
		if reduction := exactMean(t, first); reduction != nil {
			require.Equal(t, "20", digitsOfExact(t, reduction.Value))
			require.Equal(t, int64(3), reduction.Count)
			require.Equal(t, "60", digitsOfExact(t, reduction.Sum))
		}
		second := closeOneWindow(t, a, 30)
		if reduction := exactMean(t, second); reduction != nil {
			require.Equal(t, "40", digitsOfExact(t, reduction.Value), "the sample exactly at 15s belongs to the next window")
			require.Equal(t, int64(1), reduction.Count)
		}
	})

	t.Run("a non-numeric mean is invalid instead of falling back to the latest reading", func(t *testing.T) {
		a := newMeanAssembler(t, nil)
		accept(t, a, acceptedSample("state", 3, measurement.NewText("on")))
		member := closeOneWindow(t, a, 10)
		require.Nil(t, member.Reduction, "no reduction value is invented")
		require.Equal(t, MemberInvalid, member.Status)
		require.Equal(t, "non-numeric", member.Reason)
		require.Equal(t, "state", member.Sample.SampleID, "the reading stays visible as evidence")
	})

	t.Run("latest is unchanged and mean never replaces it", func(t *testing.T) {
		samples := []Sample{intSample("a10", 1, 10), intSample("a20", 2, 20), intSample("a30", 3, 30)}
		mean, latest := newMeanAssembler(t, nil), newLatestAssembler(t)
		accept(t, mean, samples...)
		accept(t, latest, samples...)
		meanMember := closeOneWindow(t, mean, 10)
		latestMember := closeOneWindow(t, latest, 10)
		if reduction := exactMean(t, meanMember); reduction != nil {
			require.Equal(t, "20", digitsOfExact(t, reduction.Value))
		}
		require.Nil(t, latestMember.Reduction, "latest mode carries no reduction result")
		require.Equal(t, "a30", latestMember.Sample.SampleID)
		require.Equal(t, "30", digitsOfExact(t, latestMember.Sample.Value))
	})
}

func TestMeanWindowConfigValidation(t *testing.T) {
	// Mean belongs to the instantaneous profile; cumulative and wide keep latest.
	_, err := NewAssembler(meanConfig(func(c *Config) { c.Mode = ModeWide }))
	require.ErrorIs(t, err, ErrInvalidConfig)
	require.ErrorContains(t, err, "mean")

	_, err = NewAssembler(meanConfig(func(c *Config) { c.Mode = ModeCumulativeProfile }))
	require.ErrorIs(t, err, ErrInvalidConfig)
	require.ErrorContains(t, err, "mean")

	_, err = NewAssembler(meanConfig(func(c *Config) { c.Members[0].SampleReduction = "average" }))
	require.ErrorIs(t, err, ErrInvalidConfig)

	_, err = NewAssembler(meanConfig(func(c *Config) { c.Members[0].SourceValueType = schema.DataType("decimal") }))
	require.ErrorIs(t, err, ErrInvalidConfig, "an unknown frozen source enum is refused")
}
