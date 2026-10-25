package snapshot

import (
	"math"
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func interval15(t *testing.T, mutate func(*Config)) *Assembler {
	t.Helper()
	return newMeanAssembler(t, func(c *Config) {
		c.Interval = 15 * time.Second
		if mutate != nil {
			mutate(c)
		}
	})
}

func evidenceOf(t *testing.T, member MemberResult) *MeanEvidence {
	t.Helper()
	require.NotNil(t, member.Evidence, "a mean member carries window evidence")
	return member.Evidence
}

func badSample(id string, seconds int, reason string) Sample {
	sample := intSample(id, seconds, 1000)
	sample.Quality, sample.QualityReason = schema.QualityBad, reason
	return sample
}

// TestMeanWindowEvidence pins the missing-data contract: the evidence reports
// only what was accepted, actual counts and observed bounds, and a coverage that
// is never complete without a verifiable whole-window proof.
func TestMeanWindowEvidence(t *testing.T) {
	t.Run("10, 20 and 30 with a retransmission and a bad reading", func(t *testing.T) {
		a := interval15(t, nil)
		accept(t, a, intSample("a3", 3, 10), intSample("a6", 6, 20), intSample("a9", 9, 30))
		require.Equal(t, OfferDuplicate, a.Offer(intSample("a3", 3, 10), at(12)).Outcome)
		accept(t, a, badSample("bad", 12, "connect postgres://user:secret@db/readings failed"))

		member := closeOneWindow(t, a, 15)
		evidence := evidenceOf(t, member)
		if evidence != nil {
			require.Equal(t, measurement.MeanAggregationVersion, evidence.Version)
			require.Equal(t, measurement.MeanAggregationMethod, evidence.Method)
			require.Equal(t, measurement.MeanBasisExact, evidence.Basis)
			require.Equal(t, measurement.ExactDecimal, evidence.ResultValueType)
			require.Equal(t, int64(3), evidence.ValidCount)
			require.Equal(t, int64(1), evidence.RejectedCount, "a retransmission is never a rejection")
			require.Equal(t, at(3), evidence.FirstObservedAt)
			require.Equal(t, at(12), evidence.LastObservedAt)
			require.Equal(t, CoveragePartial, evidence.Coverage)
			require.Equal(t, "bad", evidence.Reason)
			require.Equal(t, []string{"bad"}, evidence.RejectionReasons)
			require.NotContains(t, evidence.Reason, "secret", "the raw quality reason is not projected")
			require.NotContains(t, evidence.RejectionReasons[0], "secret")
		}
	})

	t.Run("a known failure makes partial and the next window stays unknown", func(t *testing.T) {
		a := interval15(t, nil)
		accept(t, a, intSample("a3", 3, 10), intSample("a9", 9, 30), badSample("bad", 12, "read-timeout"))
		first := closeOneWindow(t, a, 15)
		if evidence := evidenceOf(t, first); evidence != nil {
			require.Equal(t, "20", digitsOfExact(t, first.Reduction.Value))
			require.Equal(t, int64(2), evidence.ValidCount)
			require.Equal(t, int64(1), evidence.RejectedCount)
			require.Equal(t, CoveragePartial, evidence.Coverage)
		}

		require.Contains(t, []OfferOutcome{OfferSelected, OfferNotSelected}, a.Replay(intSample("b18", 18, 20)).Outcome)
		require.Contains(t, []OfferOutcome{OfferSelected, OfferNotSelected}, a.Replay(intSample("b24", 24, 40)).Outcome)
		second := closeOneWindow(t, a, 30)
		if evidence := evidenceOf(t, second); evidence != nil {
			require.Equal(t, "30", digitsOfExact(t, second.Reduction.Value), "the next window uses only its own observations")
			require.Equal(t, int64(2), evidence.ValidCount)
			require.Equal(t, int64(0), evidence.RejectedCount)
			require.Equal(t, CoverageUnknown, evidence.Coverage, "no whole-window proof means unknown, never complete")
			require.Equal(t, ReasonCoverageUnverified, evidence.Reason)
			require.NotEqual(t, CoverageComplete, evidence.Coverage)
		}
	})

	t.Run("a silent window is missing without a fabricated cause", func(t *testing.T) {
		a := interval15(t, nil)
		member := closeOneWindow(t, a, 15)
		require.Nil(t, member.Reduction)
		require.Equal(t, MemberMissing, member.Status)
		require.Equal(t, ReasonNoSamples, member.Reason)
		evidence := evidenceOf(t, member)
		if evidence != nil {
			require.Equal(t, CoverageMissing, evidence.Coverage)
			require.Equal(t, ReasonNoSamples, evidence.Reason)
			require.Equal(t, int64(0), evidence.ValidCount)
			require.Equal(t, int64(0), evidence.RejectedCount)
			require.True(t, evidence.FirstObservedAt.IsZero(), "no observation time is invented")
			require.True(t, evidence.LastObservedAt.IsZero())
			require.Empty(t, evidence.Basis, "no basis is claimed without a value")
			require.Empty(t, evidence.ResultValueType)
		}
	})

	t.Run("a slower acquisition has actual counts without whole-window proof", func(t *testing.T) {
		a := interval15(t, nil)
		accept(t, a, intSample("a3", 3, 10))
		member := closeOneWindow(t, a, 15)
		if evidence := evidenceOf(t, member); evidence != nil {
			require.Equal(t, "10", digitsOfExact(t, member.Reduction.Value))
			require.Equal(t, int64(1), evidence.ValidCount)
			require.Equal(t, int64(0), evidence.RejectedCount)
			require.Equal(t, CoverageUnknown, evidence.Coverage)
			require.Equal(t, ReasonCoverageUnverified, evidence.Reason)
		}
	})

	t.Run("a cutover window counts only actual observations", func(t *testing.T) {
		a := interval15(t, func(c *Config) { c.StartsAt = at(7) })
		accept(t, a, intSample("before", 1, 10), badSample("badBefore", 2, "read-timeout"))
		accept(t, a, intSample("a8", 8, 200), intSample("a9", 9, 300))
		outcomes := a.Tick(at(15))
		require.Len(t, outcomes, 1)
		member := outcomes[0].Members[0]
		if evidence := evidenceOf(t, member); evidence != nil {
			require.Equal(t, "250", digitsOfExact(t, member.Reduction.Value))
			require.Equal(t, int64(2), evidence.ValidCount)
			require.Equal(t, int64(0), evidence.RejectedCount, "an observation outside the actual window is not an in-window failure")
			require.Equal(t, at(8), evidence.FirstObservedAt)
			require.Equal(t, at(9), evidence.LastObservedAt)
			require.Equal(t, CoveragePartial, evidence.Coverage)
			require.Equal(t, ReasonBoundaryTruncated, evidence.Reason)
		}
	})

	t.Run("source and result types stay frozen facts", func(t *testing.T) {
		exact := interval15(t, func(c *Config) { c.Members[0].SourceValueType = schema.DataTypeInt16 })
		accept(t, exact, intSample("a3", 3, 10))
		member := closeOneWindow(t, exact, 15)
		if evidence := evidenceOf(t, member); evidence != nil {
			require.Equal(t, schema.DataTypeInt16, evidence.SourceValueType)
			require.Equal(t, measurement.ExactDecimal, evidence.ResultValueType, "an exact mean carries a decimal result")
		}

		float := interval15(t, func(c *Config) { c.Members[0].SourceValueType = schema.DataTypeFloat32 })
		accept(t, float, acceptedSample("f3", 3, mustFloatValue(t, 1.5)))
		member = closeOneWindow(t, float, 15)
		if evidence := evidenceOf(t, member); evidence != nil {
			require.Equal(t, schema.DataTypeFloat32, evidence.SourceValueType)
			require.Equal(t, measurement.ExactFloat64, evidence.ResultValueType, "a float mean carries a float64 result")
			require.Equal(t, measurement.MeanBasisFloat64, evidence.Basis)
		}
	})

	t.Run("replay in another arrival order yields identical evidence", func(t *testing.T) {
		journal := []Sample{intSample("a3", 3, 10), intSample("a6", 6, 20), intSample("a9", 9, 30), badSample("bad", 12, "read-timeout")}
		replay := func(order []Sample) MemberResult {
			a := interval15(t, nil)
			for _, sample := range order {
				require.Contains(t, []OfferOutcome{OfferSelected, OfferNotSelected}, a.Replay(sample).Outcome, sample.SampleID)
			}
			return closeOneWindow(t, a, 15)
		}
		forward := replay(journal)
		backward := replay([]Sample{journal[3], journal[2], journal[1], journal[0]})
		require.Equal(t, forward, backward, "replaying the same accepted journal reproduces the result and evidence")
		if evidence := evidenceOf(t, forward); evidence != nil {
			require.Equal(t, []string{"bad"}, evidence.RejectionReasons)
		}
	})

	t.Run("the replay keeps stable public reasons", func(t *testing.T) {
		journal := []Sample{intSample("a3", 3, 10), badSample("bad", 6, "secret endpoint /dsn")}
		replay := func(order []Sample) MemberResult {
			a := interval15(t, nil)
			for _, sample := range order {
				require.Contains(t, []OfferOutcome{OfferSelected, OfferNotSelected}, a.Replay(sample).Outcome, sample.SampleID)
			}
			return closeOneWindow(t, a, 15)
		}
		forward := replay(journal)
		backward := replay([]Sample{journal[1], journal[0]})
		require.Equal(t, forward, backward)
		if evidence := evidenceOf(t, forward); evidence != nil {
			require.Equal(t, []string{"bad"}, evidence.RejectionReasons)
		}
	})

	t.Run("a text reading is invalid, not a valid numeric sample", func(t *testing.T) {
		a := interval15(t, nil)
		accept(t, a, acceptedSample("state", 3, measurement.NewText("on")))
		member := closeOneWindow(t, a, 15)
		require.Nil(t, member.Reduction)
		require.Equal(t, MemberInvalid, member.Status)
		require.False(t, member.Usable())
		if evidence := evidenceOf(t, member); evidence != nil {
			require.Equal(t, CoverageMissing, evidence.Coverage)
			require.Equal(t, "non-numeric", evidence.Reason)
			require.Equal(t, int64(0), evidence.ValidCount)
			require.Equal(t, int64(1), evidence.RejectedCount)
			require.Equal(t, []string{"non-numeric"}, evidence.RejectionReasons)
			require.Equal(t, at(3), evidence.FirstObservedAt)
		}
	})

	t.Run("a text reading never enters the denominator", func(t *testing.T) {
		a := interval15(t, nil)
		accept(t, a, intSample("a3", 3, 10), acceptedSample("state", 6, measurement.NewText("on")), intSample("a9", 9, 30))
		member := closeOneWindow(t, a, 15)
		reduction := exactMean(t, member)
		if evidence := evidenceOf(t, member); reduction != nil && evidence != nil {
			require.Equal(t, MemberOK, member.Status)
			require.Equal(t, "20", digitsOfExact(t, reduction.Value))
			require.Equal(t, int64(2), evidence.ValidCount)
			require.Equal(t, int64(1), evidence.RejectedCount)
			require.Equal(t, CoveragePartial, evidence.Coverage)
			require.Equal(t, "non-numeric", evidence.Reason)
			require.Equal(t, "a9", member.Sample.SampleID)
		}
	})

	t.Run("invalid and stale qualities keep stable evidence", func(t *testing.T) {
		a := interval15(t, nil)
		invalid := acceptedSample("invalid", 6, measurement.NewInt64(1000))
		invalid.Quality, invalid.QualityReason = schema.QualityInvalid, "secret invalid detail"
		stale := acceptedSample("stale", 9, measurement.NewInt64(2000))
		stale.Quality, stale.QualityReason = schema.QualityStale, "secret stale detail"
		accept(t, a, intSample("good", 3, 10), invalid, stale)
		member := closeOneWindow(t, a, 15)
		reduction := exactMean(t, member)
		if evidence := evidenceOf(t, member); reduction != nil && evidence != nil {
			require.Equal(t, "10", digitsOfExact(t, reduction.Value))
			require.Equal(t, int64(1), evidence.ValidCount)
			require.Equal(t, int64(2), evidence.RejectedCount)
			require.Equal(t, []string{"invalid", "stale"}, evidence.RejectionReasons)
			require.Equal(t, "invalid", evidence.Reason)
			require.Equal(t, CoveragePartial, evidence.Coverage)
		}

		only := interval15(t, nil)
		accept(t, only, stale)
		missing := closeOneWindow(t, only, 15)
		require.Nil(t, missing.Reduction)
		require.Equal(t, MemberStale, missing.Status)
		if evidence := evidenceOf(t, missing); evidence != nil {
			require.Equal(t, CoverageMissing, evidence.Coverage)
			require.Equal(t, "stale", evidence.Reason)
			require.Equal(t, int64(0), evidence.ValidCount)
			require.Equal(t, int64(1), evidence.RejectedCount)
		}
	})

	t.Run("a value that contradicts the frozen source type is excluded", func(t *testing.T) {
		a := interval15(t, func(c *Config) { c.Members[0].SourceValueType = schema.DataTypeInt16 })
		accept(t, a, intSample("a3", 3, 10), acceptedSample("state", 6, measurement.NewText("on")), acceptedSample("f9", 9, mustFloatValue(t, 1.5)))
		member := closeOneWindow(t, a, 15)
		reduction := exactMean(t, member)
		if evidence := evidenceOf(t, member); reduction != nil && evidence != nil {
			require.Equal(t, "10", digitsOfExact(t, reduction.Value))
			require.Equal(t, int64(1), evidence.ValidCount)
			require.Equal(t, int64(2), evidence.RejectedCount)
			require.Equal(t, []string{"type-mismatch"}, evidence.RejectionReasons)
			require.Equal(t, CoveragePartial, evidence.Coverage)
		}
	})

	t.Run("a known frozen type states its typed contract for a silent window", func(t *testing.T) {
		exact := interval15(t, func(c *Config) { c.Members[0].SourceValueType = schema.DataTypeInt16 })
		member := closeOneWindow(t, exact, 15)
		require.Nil(t, member.Reduction)
		if evidence := evidenceOf(t, member); evidence != nil {
			require.Equal(t, schema.DataTypeInt16, evidence.SourceValueType)
			require.Equal(t, measurement.ExactDecimal, evidence.ResultValueType)
			require.Equal(t, measurement.MeanBasisExact, evidence.Basis)
			require.Equal(t, CoverageMissing, evidence.Coverage)
			require.Equal(t, ReasonNoSamples, evidence.Reason)
		}

		float := interval15(t, func(c *Config) { c.Members[0].SourceValueType = schema.DataTypeFloat32 })
		member = closeOneWindow(t, float, 15)
		if evidence := evidenceOf(t, member); evidence != nil {
			require.Equal(t, measurement.ExactFloat64, evidence.ResultValueType)
			require.Equal(t, measurement.MeanBasisFloat64, evidence.Basis)
		}
	})

	t.Run("a successful mean carries no member-level reason", func(t *testing.T) {
		unknown := interval15(t, nil)
		accept(t, unknown, intSample("a3", 3, 10), intSample("a6", 6, 20))
		member := closeOneWindow(t, unknown, 15)
		require.Equal(t, MemberOK, member.Status)
		require.Empty(t, member.Reason, "a usable value has no failure reason")
		if evidence := evidenceOf(t, member); evidence != nil {
			require.Equal(t, CoverageUnknown, evidence.Coverage)
			require.Equal(t, ReasonCoverageUnverified, evidence.Reason, "the coverage note stays in the evidence")
		}

		partial := interval15(t, nil)
		accept(t, partial, intSample("b3", 3, 10), badSample("b6", 6, "read-timeout"))
		member = closeOneWindow(t, partial, 15)
		require.Equal(t, MemberOK, member.Status)
		require.Empty(t, member.Reason)
		if evidence := evidenceOf(t, member); evidence != nil {
			require.Equal(t, CoveragePartial, evidence.Coverage)
			require.Equal(t, "bad", evidence.Reason)
		}
	})

	t.Run("signed and unsigned values never cross a frozen source type", func(t *testing.T) {
		signed := interval15(t, func(c *Config) { c.Members[0].SourceValueType = schema.DataTypeInt64 })
		accept(t, signed, intSample("a3", 3, 10), acceptedSample("u6", 6, measurement.NewUint64(math.MaxUint64)))
		member := closeOneWindow(t, signed, 15)
		reduction := exactMean(t, member)
		if evidence := evidenceOf(t, member); reduction != nil && evidence != nil {
			require.Equal(t, "10", digitsOfExact(t, reduction.Value))
			require.Equal(t, int64(1), evidence.ValidCount)
			require.Equal(t, int64(1), evidence.RejectedCount)
			require.Equal(t, []string{"type-mismatch"}, evidence.RejectionReasons)
		}

		unsigned := interval15(t, func(c *Config) { c.Members[0].SourceValueType = schema.DataTypeUint64 })
		accept(t, unsigned, acceptedSample("u3", 3, measurement.NewUint64(math.MaxUint64)), intSample("n6", 6, -1))
		member = closeOneWindow(t, unsigned, 15)
		reduction = exactMean(t, member)
		if evidence := evidenceOf(t, member); reduction != nil && evidence != nil {
			require.Equal(t, "18446744073709551615", digitsOfExact(t, reduction.Value), "a matching uint64 keeps every digit")
			require.Equal(t, int64(1), evidence.ValidCount)
			require.Equal(t, int64(1), evidence.RejectedCount)
			require.Equal(t, []string{"type-mismatch"}, evidence.RejectionReasons)
		}
	})

	t.Run("a mixed numeric basis keeps the mean-invalid evidence", func(t *testing.T) {
		a := interval15(t, nil)
		accept(t, a, intSample("a3", 3, 10), acceptedSample("f6", 6, mustFloatValue(t, 1.5)))
		member := closeOneWindow(t, a, 15)
		require.Nil(t, member.Reduction)
		require.Equal(t, MemberInvalid, member.Status)
		if evidence := evidenceOf(t, member); evidence != nil {
			require.Equal(t, CoverageMissing, evidence.Coverage)
			require.Equal(t, ReasonMeanInvalid, evidence.Reason)
			require.Equal(t, int64(2), evidence.ValidCount)
			require.Equal(t, int64(0), evidence.RejectedCount)
		}
	})
}
