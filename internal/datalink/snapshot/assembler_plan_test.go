package snapshot

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAssemblerCheckNeverMutatesState(t *testing.T) {
	a := newFixtureAssembler(t, nil)
	good := goodSample("t5", "temperature", 5, 5)

	require.Equal(t, OfferSelected, a.Check(good, at(6)).Outcome)
	require.Zero(t, a.OpenBuckets(), "Check must not create a bucket")
	require.Equal(t, OfferSelected, a.Check(good, at(6)).Outcome, "so the same sample is still new")

	require.Equal(t, OfferSelected, a.Offer(good, at(6)).Outcome)
	require.Equal(t, OfferDuplicate, a.Check(good, at(6)).Outcome)
	changed := good
	changed.QualityReason = "other"
	require.Equal(t, OfferConflict, a.Check(changed, at(6)).Outcome)
	require.Equal(t, OfferNotSelected, a.Check(goodSample("t3", "temperature", 3, 3), at(6)).Outcome)

	require.Equal(t, OfferRejected, a.Check(goodSample("far", "temperature", 3600, 1), at(6)).Outcome)
	require.Equal(t, ReasonUnknownMember, a.Check(goodSample("x", "flow", 3, 1), at(6)).Reason)
	require.Equal(t, 1, a.OpenBuckets())

	a.Offer(goodSample("p5", "pressure", 5, 5), at(6))
	require.Len(t, a.Tick(at(10)), 1)
	require.Equal(t, OfferLate, a.Check(goodSample("late", "temperature", 9, 9), at(11)).Outcome)
	require.Zero(t, a.Late(), "Check does not count late samples; only Offer does")
}

func TestAssemblerPlanDoesNotCloseUntilCommitted(t *testing.T) {
	a := newFixtureAssembler(t, nil)
	offerGood(t, a, "t5", "temperature", 5, 6)
	offerGood(t, a, "p5", "pressure", 5, 6)

	plan := a.Plan(at(10))
	require.False(t, plan.Empty())
	require.Len(t, plan.Outcomes, 1)
	require.Equal(t, OutcomeRow, plan.Outcomes[0].Kind)
	require.Equal(t, at(10), plan.NextClose())
	require.Equal(t, fixtureStart, a.NextClose(), "planning leaves the checkpoint where it was")
	require.Equal(t, 1, a.OpenBuckets())

	again := a.Plan(at(10))
	require.Equal(t, plan.Outcomes, again.Outcomes, "an uncommitted plan is repeatable, so a failed durable write can retry")

	a.Commit(plan)
	require.Equal(t, at(10), a.NextClose())
	require.Zero(t, a.OpenBuckets())
	require.True(t, a.Plan(at(10)).Empty())
	require.Equal(t, OfferLate, a.Offer(goodSample("t9", "temperature", 9, 9), at(11)).Outcome)
}

func TestAssemblerPlanCoversSilentBucketsAndMatchesTick(t *testing.T) {
	planned := newFixtureAssembler(t, nil)
	ticked := newFixtureAssembler(t, nil)
	for _, a := range []*Assembler{planned, ticked} {
		offerGood(t, a, "t15", "temperature", 15, 16)
	}
	plan := planned.Plan(at(30))
	planned.Commit(plan)
	require.Equal(t, ticked.Tick(at(30)), plan.Outcomes)
	require.Equal(t, ticked.NextClose(), planned.NextClose())
	require.Len(t, plan.Outcomes, 3)
	require.Equal(t, OutcomeNoData, plan.Outcomes[0].Kind)
}

func TestAssemblerReplayRebuildsOpenBucketBypassingFutureSkew(t *testing.T) {
	a := newFixtureAssembler(t, nil)
	// Accepted earlier when the gateway clock was already near t=9; replay runs
	// at restart time, so the future-skew rule cannot be applied again.
	require.Equal(t, OfferSelected, a.Replay(goodSample("t9", "temperature", 9, 9)).Outcome)
	require.Equal(t, OfferSelected, a.Replay(goodSample("p9", "pressure", 9, 9)).Outcome)
	require.Equal(t, OfferDuplicate, a.Replay(goodSample("p9", "pressure", 9, 9)).Outcome)

	outcomes := a.Tick(at(10))
	require.Len(t, outcomes, 1)
	require.Equal(t, OutcomeRow, outcomes[0].Kind)

	require.Equal(t, OfferLate, a.Replay(goodSample("old", "temperature", 3, 3)).Outcome, "a replayed sample of a closed bucket is not reintroduced")
	require.Equal(t, ReasonUnknownMember, a.Replay(goodSample("x", "flow", 12, 1)).Reason)
}

func TestAssemblerUntilStopsClosingAtTheSupersedingBoundary(t *testing.T) {
	a := newFixtureAssembler(t, func(c *Config) { c.Until = at(20) })
	require.False(t, a.Done())
	outcomes := a.Tick(at(100))
	require.Len(t, outcomes, 2, "only buckets [0,10) and [10,20) belong to this revision")
	require.Equal(t, at(20), a.NextClose())
	require.True(t, a.Done())
	require.Empty(t, a.Tick(at(200)))

	for name, until := range map[string]time.Time{"unaligned": at(15), "before first": at(0).Add(-10 * time.Second)} {
		config := fixtureConfig()
		config.Members = append([]Member(nil), config.Members...)
		config.Until = until
		_, err := NewAssembler(config)
		require.ErrorIs(t, err, ErrInvalidConfig, name)
	}
}
