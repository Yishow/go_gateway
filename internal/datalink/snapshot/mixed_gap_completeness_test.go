package snapshot

import (
	"fmt"
	"go-gateway/internal/datalink/aggregation"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMixedGapSpanIsEvidenceNotCurrentWindowUsage(t *testing.T) {
	for _, mode := range []Mode{ModeMixedWide, ModeMixedLong} {
		for _, required := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/required=%t", mode, required), func(t *testing.T) {
				c := mixedCounterConfig(mode)
				c.Members[2].Required = required
				a, err := NewAssembler(c)
				require.NoError(t, err)
				mixedPolicyBaseline(t, a, 1000)
				mixedGoodGauges(t, a, 115)
				missing := a.Plan(mixedTime(120))
				_, missingResult := mixedDueMember(t, missing, 120, "c")
				require.False(t, missingResult.Usable())
				require.Equal(t, ReasonNoSamples, missingResult.Reason)
				require.Empty(t, missing.CounterCheckpoints())
				a.Commit(missing)
				require.Equal(t, "baseline", a.CounterCheckpoints()["c"].Reading.SampleID)
				mixedGoodGauges(t, a, 175)
				mixedOffer(t, a, mixedIntegerSample("c", "after-gap", 179, 1008))
				p := a.Plan(mixedTime(180))
				row, r := mixedDueMember(t, p, 180, "c")
				require.False(t, r.Usable(), "span total must not enter the current-window result")
				require.Equal(t, string(aggregation.CounterGap), r.Reason)
				require.Equal(t, aggregation.CounterGap, r.Counter.Status)
				require.NotNil(t, r.Counter.Delta)
				require.True(t, r.Counter.Delta.Span)
				require.Equal(t, "8", mixedPolicyDigits(t, r), "actual span total remains evidence, not the selected column value")
				require.Equal(t, mixedTime(59), r.Counter.Delta.From)
				require.Equal(t, mixedTime(179), r.Counter.Delta.To)
				require.Equal(t, mixedTime(120), r.Window.AlignedStart)
				require.Equal(t, mixedTime(180), r.Window.AlignedEnd)
				if required {
					require.Equal(t, OutcomeSkipped, row.Kind)
					require.Empty(t, row.EffectKey)
				} else {
					require.Equal(t, OutcomeRow, row.Kind)
					require.True(t, row.Partial)
				}
				require.Equal(t, "after-gap", p.CounterCheckpoints()["c"].Reading.SampleID)
				require.Equal(t, "baseline", a.CounterCheckpoints()["c"].Reading.SampleID)
				a.Commit(p)
				progress, _ := a.MixedProgress()
				fresh, err := NewAssembler(c)
				require.NoError(t, err)
				require.NoError(t, fresh.RestoreMixedProgress(progress))
				require.NoError(t, fresh.RestoreCounterCheckpoints(a.CounterCheckpoints()))
				mixedGoodGauges(t, fresh, 235)
				mixedOffer(t, fresh, mixedIntegerSample("c", "following", 239, 1013))
				_, next := mixedDueMember(t, fresh.Plan(mixedTime(240)), 240, "c")
				require.True(t, next.Usable())
				require.False(t, next.Counter.Delta.Span)
				require.Equal(t, "5", mixedPolicyDigits(t, next))
			})
		}
	}
}
