package snapshot

import (
	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/measurement"
	"testing"

	"github.com/stretchr/testify/require"
)

func mixedPolicyAssembler(t *testing.T, mode Mode, mutate func(*aggregation.CounterBinding)) *Assembler {
	t.Helper()
	c := mixedCounterConfig(mode)
	b := c.Counters["c"]
	if mutate != nil {
		mutate(&b)
	}
	c.Counters["c"] = b
	a, err := NewAssembler(c)
	require.NoError(t, err)
	return a
}
func mixedPolicyResult(t *testing.T, p Plan) MemberResult {
	t.Helper()
	require.NoError(t, p.Err)
	require.NotEmpty(t, p.Outcomes)
	row := p.Outcomes[len(p.Outcomes)-1]
	require.NotEmpty(t, row.Members)
	return row.Members[len(row.Members)-1]
}
func mixedPolicyBaseline(t *testing.T, a *Assembler, value int64) {
	t.Helper()
	mixedOffer(t, a, mixedIntegerSample("c", "baseline", 59, value))
	p := a.Plan(mixedTime(60))
	r := mixedPolicyResult(t, p)
	require.Equal(t, aggregation.CounterBaseline, r.Counter.Status)
	require.Nil(t, r.Counter.Delta)
	require.False(t, r.Usable())
	require.Empty(t, a.CounterCheckpoints())
	a.Commit(p)
}
func mixedPolicyDigits(t *testing.T, r MemberResult) string {
	t.Helper()
	require.NotNil(t, r.Counter)
	require.NotNil(t, r.Counter.Delta)
	digits, ok := r.Counter.Delta.Value.DecimalDigits()
	require.True(t, ok)
	return digits
}
func mixedRollover(rate bool) func(*aggregation.CounterBinding) {
	return func(b *aggregation.CounterBinding) {
		mod := measurement.NewUint64(65536)
		b.RolloverModulus = &mod
		if rate {
			v := measurement.NewUint64(100)
			b.MaxRatePerSecond = &v
		}
	}
}

func TestMixedCounterSignedAndConfirmedRolloverPolicy(t *testing.T) {
	for _, mode := range []Mode{ModeMixedWide, ModeMixedLong} {
		t.Run(string(mode), func(t *testing.T) {
			cases := []struct {
				name           string
				first, last    int64
				mutate         func(*aggregation.CounterBinding)
				digits, reason string
			}{
				{"signed negative", 10, 7, func(b *aggregation.CounterBinding) { b.Kind = aggregation.CounterKindSigned }, "-3", ""},
				{"monotonic decrease", 10, 7, nil, "", aggregation.ReasonDecreaseUnproven},
				{"proven single wrap", 65530, 4, mixedRollover(true), "10", ""},
				{"unproven wrap", 65530, 4, mixedRollover(false), "", aggregation.ReasonRolloverEvidence},
			}
			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					a := mixedPolicyAssembler(t, mode, tc.mutate)
					mixedPolicyBaseline(t, a, tc.first)
					mixedOffer(t, a, mixedIntegerSample("c", "tail", 119, tc.last))
					p := a.Plan(mixedTime(120))
					r := mixedPolicyResult(t, p)
					if tc.digits != "" {
						require.True(t, r.Usable())
						require.Equal(t, aggregation.CounterOK, r.Counter.Status)
						require.Equal(t, tc.digits, mixedPolicyDigits(t, r))
					} else {
						require.False(t, r.Usable())
						require.Nil(t, r.Counter.Delta)
						require.Equal(t, tc.reason, r.Counter.Reason)
					}
					require.Equal(t, "baseline", a.CounterCheckpoints()["c"].Reading.SampleID)
					require.Equal(t, "tail", p.CounterCheckpoints()["c"].Reading.SampleID)
					a.Commit(p)
					require.Equal(t, "tail", a.CounterCheckpoints()["c"].Reading.SampleID)
				})
			}
		})
	}
}

func TestMixedCounterResetUsesFullSequenceAndReliableTail(t *testing.T) {
	for _, mode := range []Mode{ModeMixedWide, ModeMixedLong} {
		t.Run(string(mode), func(t *testing.T) {
			a := mixedPolicyAssembler(t, mode, nil)
			mixedPolicyBaseline(t, a, 90)
			for _, s := range []Sample{mixedIntegerSample("c", "after-reset", 118, 110), mixedIntegerSample("c", "before-reset", 70, 100), mixedIntegerSample("c", "reset", 90, 10)} {
				mixedOffer(t, a, s)
			}
			p := a.Plan(mixedTime(120))
			r := mixedPolicyResult(t, p)
			require.False(t, r.Usable())
			require.Equal(t, aggregation.CounterUncertain, r.Counter.Status)
			require.Nil(t, r.Counter.Delta)
			require.Equal(t, aggregation.ReasonDecreaseUnproven, r.Counter.Reason)
			require.Len(t, r.Counter.Events, 1)
			require.Equal(t, "before-reset", r.Counter.Events[0].FromSampleID)
			require.Equal(t, "reset", r.Counter.Events[0].ToSampleID)
			require.Equal(t, "after-reset", p.CounterCheckpoints()["c"].Reading.SampleID)
			require.Equal(t, "baseline", a.CounterCheckpoints()["c"].Reading.SampleID)
			require.Equal(t, p.Outcomes, a.Plan(mixedTime(120)).Outcomes)
			a.Commit(p)
			mixedOffer(t, a, mixedIntegerSample("c", "next", 179, 125))
			r = mixedPolicyResult(t, a.Plan(mixedTime(180)))
			require.True(t, r.Usable())
			require.Equal(t, "15", mixedPolicyDigits(t, r))
		})
	}
}

func TestMixedCounterContradictoryTailCannotAdvance(t *testing.T) {
	a := mixedPolicyAssembler(t, ModeMixedWide, nil)
	mixedPolicyBaseline(t, a, 90)
	mixedOffer(t, a, mixedIntegerSample("c", "same-time-a", 119, 95))
	mixedOffer(t, a, mixedIntegerSample("c", "same-time-b", 119, 96))
	p := a.Plan(mixedTime(120))
	r := mixedPolicyResult(t, p)
	require.False(t, r.Usable())
	require.Equal(t, aggregation.ReasonContradictory, r.Counter.Reason)
	require.Nil(t, r.Counter.Delta)
	require.Empty(t, p.CounterCheckpoints())
	a.Commit(p)
	require.Equal(t, "baseline", a.CounterCheckpoints()["c"].Reading.SampleID)
}
