package snapshot

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMixedDueUnionLongAndWide(t *testing.T) {
	for _, mode := range []Mode{ModeMixedWide, ModeMixedLong} {
		t.Run(string(mode), func(t *testing.T) {
			c := mixedConfig(mode)
			a, err := NewAssembler(c)
			require.NoError(t, err)
			for _, sec := range []int{10, 25, 40, 55} {
				mixedOffer(t, a, mixedSample("a", mixedTime(sec).String(), sec, 22))
			}
			mixedOffer(t, a, mixedSample("b", "slow", 55, 8))
			p := a.Plan(mixedTime(60))
			require.NoError(t, p.Err)
			if mode == ModeMixedWide {
				require.Len(t, p.Outcomes, 4)
				for i, row := range p.Outcomes {
					require.Equal(t, mixedTime((i+1)*15), row.DueAt)
					require.Equal(t, OutcomeRow, row.Kind)
					require.Len(t, row.Members, 2)
					if i < 3 {
						r := row.Members[1]
						require.Equal(t, MemberNotDue, r.Status)
						require.Equal(t, ReasonNotDue, r.Reason)
						require.Nil(t, r.Sample)
						require.Nil(t, r.Counter)
						require.Nil(t, r.Reduction)
						require.True(t, r.Window.AlignedStart.IsZero())
						require.Equal(t, mixedTime(60), r.Window.NextDue)
					}
				}
			} else {
				require.Len(t, p.Outcomes, 5)
				for i, row := range p.Outcomes[:4] {
					require.Equal(t, "a", row.EntityKey)
					require.Equal(t, mixedTime((i+1)*15), row.DueAt)
				}
				require.Equal(t, "b", p.Outcomes[4].EntityKey)
			}
			require.Equal(t, p.Outcomes, a.Plan(mixedTime(60)).Outcomes, "Plan is pure")
			require.True(t, a.NextClose().IsZero(), "mixed must not reinterpret legacy NextClose")
			checkpoint, ok := p.MixedProgress()
			require.True(t, ok)
			require.Equal(t, mixedTime(60), checkpoint.ClosedThrough)
			require.True(t, p.NextClose().IsZero())
			a.Commit(p)
			require.True(t, a.Plan(mixedTime(60)).Empty())
			require.Empty(t, a.CounterCheckpoints(), "latest and mean must not create baselines")
		})
	}
}

func TestMixedBoundaryFirstPartialAndRestart(t *testing.T) {
	c := mixedConfig(ModeMixedWide)
	c.StartsAt = mixedTime(7)
	c.AllowedLateness = 2 * time.Second
	a, err := NewAssembler(c)
	require.NoError(t, err)
	require.Equal(t, OfferRejected, a.Offer(mixedSample("a", "before", 6, 999), mixedTime(7)).Outcome)
	mixedOffer(t, a, mixedSample("a", "inside", 8, 20))
	mixedOffer(t, a, mixedSample("a", "boundary", 15, 24))
	mixedOffer(t, a, mixedSample("b", "slow", 10, 100))
	require.Empty(t, a.Tick(mixedTime(16)))
	rows := a.Tick(mixedTime(17))
	require.Len(t, rows, 1)
	r := rows[0].Members[0]
	require.Equal(t, mixedTime(0), r.Window.AlignedStart)
	require.Equal(t, mixedTime(15), r.Window.AlignedEnd)
	require.Equal(t, mixedTime(7), r.Window.ActualStart)
	require.Equal(t, "start", r.Window.Truncated)
	require.Equal(t, "inside", r.Sample.SampleID)
	require.Equal(t, CoveragePartial, r.Evidence.Coverage)
	progress, _ := a.MixedProgress()
	b, err := NewAssembler(c)
	require.NoError(t, err)
	require.NoError(t, b.RestoreMixedProgress(progress))
	// The retained slow selector is replayed from its accepted journal after restore.
	mixedOffer(t, b, mixedSample("b", "slow", 10, 100))
	mixedOffer(t, b, mixedSample("a", "boundary", 15, 24))
	require.Equal(t, OfferLate, b.Replay(mixedSample("a", "closed", 8, 20)).Outcome)
	require.Equal(t, a.Plan(mixedTime(62)).Outcomes, b.Plan(mixedTime(62)).Outcomes)
	next, err := b.MemberProgress()
	require.NoError(t, err)
	require.Equal(t, MemberProgress{mixedTime(15), mixedTime(30)}, next["a"])
	require.Equal(t, MemberProgress{mixedTime(0), mixedTime(60)}, next["b"])
}

func TestMixedCadencePhaseAndBoundedCatchup(t *testing.T) {
	c := mixedConfig(ModeMixedWide)
	c.StartsAt = mixedTime(7)
	c.MaxBucketsPerTick = 2
	c.Members = append(c.Members, Member{Key: "five", Calculation: CalculationLatest, Interval: 5 * time.Second, SourceValueType: c.Members[0].SourceValueType}, Member{Key: "tenmin", Calculation: CalculationLatest, Interval: 600 * time.Second, SourceValueType: c.Members[0].SourceValueType})
	a, err := NewAssembler(c)
	require.NoError(t, err)
	next, err := a.MemberProgress()
	require.NoError(t, err)
	require.Equal(t, mixedTime(10), next["five"].End)
	require.Equal(t, mixedTime(15), next["a"].End)
	require.Equal(t, mixedTime(60), next["b"].End)
	require.Equal(t, mixedTime(600), next["tenmin"].End)
	p := a.Plan(mixedTime(600))
	require.NoError(t, p.Err)
	require.Len(t, p.Outcomes, 2)
	require.Equal(t, mixedTime(10), p.Outcomes[0].DueAt)
	require.Equal(t, mixedTime(15), p.Outcomes[1].DueAt)
	a.Commit(p)
	p = a.Plan(mixedTime(600))
	require.Len(t, p.Outcomes, 2)
	require.Equal(t, mixedTime(20), p.Outcomes[0].DueAt)
}

func TestMixedLongDueFailureDoesNotSuppressOtherMember(t *testing.T) {
	c := mixedConfig(ModeMixedLong)
	a, err := NewAssembler(c)
	require.NoError(t, err)
	mixedOffer(t, a, mixedSample("a", "present", 55, 22))
	rows := a.Tick(mixedTime(60))
	require.Len(t, rows, 5)
	require.Equal(t, OutcomeRow, rows[3].Kind)
	require.Equal(t, OutcomeSkipped, rows[4].Kind)
	require.Equal(t, ReasonIncompleteRequired, rows[4].Reason)
	require.Equal(t, ReasonNoSamples, rows[4].Members[0].Reason)
	require.NotEqual(t, MemberNotDue, rows[4].Members[0].Status)
}

func TestMixedProgressAndSamplesStayMemberAndGroupScoped(t *testing.T) {
	c := mixedConfig(ModeMixedWide)
	a, err := NewAssembler(c)
	require.NoError(t, err)
	c.GroupID = "other"
	b, err := NewAssembler(c)
	require.NoError(t, err)
	// Even the identical observation ID on two members has independent selection.
	mixedOffer(t, a, mixedSample("a", "same", 10, 20))
	mixedOffer(t, a, mixedSample("b", "same", 10, 100))
	a.Tick(mixedTime(15))
	mixedOffer(t, b, mixedSample("a", "same", 10, 20))
	mixedOffer(t, b, mixedSample("b", "same", 10, 100))
	progress, _ := a.MixedProgress()
	fresh, err := NewAssembler(c)
	require.NoError(t, err)
	require.Error(t, fresh.RestoreMixedProgress(progress))
	require.Equal(t, OfferLate, a.Check(mixedSample("a", "new", 10, 20), mixedTime(20)).Outcome)
	require.Equal(t, OfferSelected, b.Check(mixedSample("a", "new", 11, 20), mixedTime(20)).Outcome)
	require.Equal(t, OfferSelected, a.Check(mixedSample("b", "new", 20, 101), mixedTime(20)).Outcome)
}

func TestMixedLongCatchupBudgetKeepsDueEventAtomic(t *testing.T) {
	c := mixedCounterConfig(ModeMixedLong)
	c.MaxBucketsPerTick = 1
	a, err := NewAssembler(c)
	require.NoError(t, err)
	p := a.Plan(mixedTime(600))
	require.NoError(t, p.Err)
	require.Len(t, p.Outcomes, 2, "one complete event, never an unbounded catch-up")
	for _, row := range p.Outcomes {
		require.Equal(t, mixedTime(15), row.DueAt)
	}
	a.Commit(p)
	next := a.Plan(mixedTime(600))
	require.Len(t, next.Outcomes, 2)
	require.Equal(t, mixedTime(30), next.Outcomes[0].DueAt)
}

func TestMixedOfferRetainsSelectorIdentityAndRevisionChecks(t *testing.T) {
	c := mixedConfig(ModeMixedWide)
	c.Members[0].SourceRevision = "source-1"
	c.Members[0].MappingRevision = "mapping-1"
	a, err := NewAssembler(c)
	require.NoError(t, err)
	s := mixedSample("a", "id", 10, 20)
	require.Equal(t, ReasonRevisionMismatch, a.Check(s, s.ObservedAt).Reason)
	s.SourceRevision = "source-1"
	s.MappingRevision = "mapping-1"
	require.Equal(t, OfferSelected, a.Check(s, s.ObservedAt).Outcome)
	require.Zero(t, a.OpenBuckets())
	mixedOffer(t, a, s)
	require.Equal(t, OfferDuplicate, a.Offer(s, s.ObservedAt).Outcome)
	changed := s
	changed.QualityReason = "changed"
	require.Equal(t, OfferConflict, a.Offer(changed, s.ObservedAt).Outcome)
	future := s
	future.SampleID = "future"
	future.ObservedAt = mixedTime(20)
	require.Equal(t, ReasonFutureSkew, a.Check(future, mixedTime(19)).Reason)
	unknown := s
	unknown.MemberKey = "unknown"
	require.Equal(t, ReasonUnknownMember, a.Check(unknown, s.ObservedAt).Reason)
	a.Tick(mixedTime(15))
	require.Equal(t, OfferLate, a.Check(s, mixedTime(20)).Outcome)
	require.Zero(t, a.Late())
	require.Equal(t, OfferLate, a.Offer(s, mixedTime(20)).Outcome)
	require.Equal(t, uint64(1), a.Late())
}
