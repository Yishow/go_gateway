package snapshot

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func mixedDueMember(t *testing.T, p Plan, seconds int, member string) (Outcome, MemberResult) {
	t.Helper()
	require.NoError(t, p.Err)
	for _, row := range p.Outcomes {
		if !row.DueAt.Equal(mixedTime(seconds)) {
			continue
		}
		for _, result := range row.Members {
			if result.MemberKey == member {
				return row, result
			}
		}
	}
	t.Fatalf("missing member %s at %d", member, seconds)
	return Outcome{}, MemberResult{}
}
func mixedGoodGauges(t *testing.T, a *Assembler, seconds int) {
	t.Helper()
	mixedOffer(t, a, mixedIntegerSample("a", fmt.Sprintf("a-%d", seconds), seconds, 22))
	mixedOffer(t, a, mixedIntegerSample("b", fmt.Sprintf("b-%d", seconds), seconds, 12))
}

func TestMixedMissingResultsFollowRequiredScope(t *testing.T) {
	for _, mode := range []Mode{ModeMixedWide, ModeMixedLong} {
		for _, required := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/required=%t", mode, required), func(t *testing.T) {
				cases := []struct {
					name, reason     string
					first, last      int64
					baseline, sample bool
				}{
					{name: "missing", reason: ReasonNoSamples},
					{name: "first baseline", reason: aggregation.ReasonFirstReading, last: 1000, sample: true},
					{name: "unproven decrease", reason: aggregation.ReasonDecreaseUnproven, first: 1000, last: 900, baseline: true, sample: true},
				}
				for _, tc := range cases {
					t.Run(tc.name, func(t *testing.T) {
						c := mixedCounterConfig(mode)
						c.Members[2].Required = required
						a, err := NewAssembler(c)
						require.NoError(t, err)
						end := 60
						if tc.baseline {
							mixedPolicyBaseline(t, a, tc.first)
							end = 120
						}
						mixedGoodGauges(t, a, end-5)
						if tc.sample {
							mixedOffer(t, a, mixedIntegerSample("c", "counter-tail", end-1, tc.last))
						}
						p := a.Plan(mixedTime(end))
						row, result := mixedDueMember(t, p, end, "c")
						require.False(t, result.Usable(), "unavailable calculation is never raw or zero")
						require.Equal(t, tc.reason, result.Reason)
						require.Nil(t, result.Reduction)
						if result.Counter != nil {
							require.Nil(t, result.Counter.Delta)
						}
						if required {
							require.Equal(t, OutcomeSkipped, row.Kind)
							require.Empty(t, row.EffectKey)
							require.Equal(t, ReasonIncompleteRequired, row.Reason)
						} else {
							require.Equal(t, OutcomeRow, row.Kind)
							require.True(t, row.Partial)
							require.NotEmpty(t, row.EffectKey)
						}
						gaugeRow, gauge := mixedDueMember(t, p, end, "a")
						require.True(t, gauge.Usable())
						require.NotNil(t, gauge.Reduction)
						if mode == ModeMixedLong {
							require.Equal(t, OutcomeRow, gaugeRow.Kind, "counter failure cannot suppress another long row")
						} else {
							require.Equal(t, row.RecordID, gaugeRow.RecordID)
						}
						if tc.sample {
							require.Equal(t, "counter-tail", p.CounterCheckpoints()["c"].Reading.SampleID, "reliable baseline can advance even when output skips")
						} else {
							require.Empty(t, p.CounterCheckpoints())
						}
						before := a.CounterCheckpoints()
						require.Equal(t, p.Outcomes, a.Plan(mixedTime(end)).Outcomes)
						require.Equal(t, before, a.CounterCheckpoints(), "planning does not commit anything")
						a.Commit(p)
						if tc.sample {
							require.Equal(t, "counter-tail", a.CounterCheckpoints()["c"].Reading.SampleID)
						}
					})
				}
			})
		}
	}
}

func TestMixedRequiredNotDueAndOptionalMissingHaveDifferentReasons(t *testing.T) {
	for _, mode := range []Mode{ModeMixedWide, ModeMixedLong} {
		t.Run(string(mode), func(t *testing.T) {
			c := mixedCounterConfig(mode)
			c.Members[0].Required = false
			a, err := NewAssembler(c)
			require.NoError(t, err)
			mixedOffer(t, a, mixedIntegerSample("b", "present", 10, 12))
			p := a.Plan(mixedTime(15))
			row, missing := mixedDueMember(t, p, 15, "a")
			require.Equal(t, OutcomeRow, row.Kind)
			require.Equal(t, ReasonNoSamples, missing.Reason)
			require.False(t, missing.Usable())
			require.Equal(t, CoverageMissing, missing.Evidence.Coverage)
			if mode == ModeMixedWide {
				_, notDue := mixedDueMember(t, p, 15, "c")
				require.Equal(t, MemberNotDue, notDue.Status)
				require.Equal(t, ReasonNotDue, notDue.Reason)
				require.Nil(t, notDue.Sample)
				require.Nil(t, notDue.Counter)
				require.True(t, notDue.Window.AlignedStart.IsZero())
			} else {
				require.Len(t, p.Outcomes, 2)
				for _, out := range p.Outcomes {
					require.NotEqual(t, "c", out.EntityKey)
				}
			}
		})
	}
}

func TestMixedMissingLatestNeverCarriesPreviousValue(t *testing.T) {
	c := mixedConfig(ModeMixedWide)
	c.Members[0].Calculation = CalculationLatest
	c.Members[0].Required = false
	a, err := NewAssembler(c)
	require.NoError(t, err)
	mixedOffer(t, a, mixedSample("a", "old", 10, 7))
	a.Commit(a.Plan(mixedTime(15)))
	row, r := mixedDueMember(t, a.Plan(mixedTime(30)), 30, "a")
	require.Equal(t, OutcomeRow, row.Kind)
	require.False(t, r.Usable())
	require.Equal(t, ReasonNoSamples, r.Reason)
	require.Nil(t, r.Sample)
	bad := mixedSample("a", "bad", 40, 999)
	bad.Quality = schema.QualityBad
	bad.QualityReason = "transport-failed"
	mixedOffer(t, a, bad)
	_, r = mixedDueMember(t, a.Plan(mixedTime(45)), 45, "a")
	require.False(t, r.Usable())
	require.Equal(t, "transport-failed", r.Reason)
	require.Equal(t, "bad", r.Sample.SampleID, "raw failure evidence is not a selected result")
}

func TestMixedSkippedGaugeDoesNotLoseReliableCounterProgress(t *testing.T) {
	c := mixedCounterConfig(ModeMixedWide)
	a, err := NewAssembler(c)
	require.NoError(t, err)
	mixedPolicyBaseline(t, a, 1000)
	// A required mean is absent, but C's reliable 8 remains part of the skip closure.
	mixedOffer(t, a, mixedIntegerSample("b", "only-gauge", 115, 12))
	mixedOffer(t, a, mixedIntegerSample("c", "reliable-tail", 119, 1008))
	p := a.Plan(mixedTime(120))
	row, r := mixedDueMember(t, p, 120, "c")
	require.Equal(t, OutcomeSkipped, row.Kind)
	require.True(t, r.Usable())
	require.Equal(t, "8", mixedPolicyDigits(t, r))
	require.Empty(t, row.EffectKey)
	require.Equal(t, "reliable-tail", p.CounterCheckpoints()["c"].Reading.SampleID)
	require.Equal(t, "baseline", a.CounterCheckpoints()["c"].Reading.SampleID)
	a.Commit(p)
	mixedGoodGauges(t, a, 175)
	mixedOffer(t, a, mixedIntegerSample("c", "next-window", 179, 1013))
	row, r = mixedDueMember(t, a.Plan(mixedTime(180)), 180, "c")
	require.Equal(t, OutcomeRow, row.Kind)
	require.Equal(t, "5", mixedPolicyDigits(t, r))
}
