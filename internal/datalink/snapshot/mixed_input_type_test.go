package snapshot

import (
	"go-gateway/internal/datalink/measurement"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMixedScaledCounterTailAndRestoreKeepExactInputType(t *testing.T) {
	for _, mode := range []Mode{ModeMixedWide, ModeMixedLong} {
		t.Run(string(mode), func(t *testing.T) {
			c := mixedCounterConfig(mode)
			c.Members = c.Members[2:]
			binding := c.Counters["c"]
			binding.ExactRecordingScale = "0.1"
			c.Counters["c"] = binding
			a, err := NewAssembler(c)
			require.NoError(t, err)
			scaled := func(id string, at int, n int64) Sample {
				s := mixedIntegerSample("c", id, at, n)
				s.Value, err = measurement.ScaleExactInteger(s.Value, "0.1")
				require.NoError(t, err)
				return s
			}
			mixedOffer(t, a, scaled("base", 59, 1000))
			a.Commit(a.Plan(mixedTime(60)))
			saved := a.CounterCheckpoints()
			require.Equal(t, measurement.ExactDecimal, saved["c"].Reading.Value.Type())
			fresh, err := NewAssembler(c)
			require.NoError(t, err)
			require.NoError(t, fresh.RestoreCounterCheckpoints(saved))
			mixedOffer(t, a, scaled("good", 118, 1008))
			mixedOffer(t, a, mixedIntegerSample("c", "unscaled-tail", 119, 1009))
			p := a.Plan(mixedTime(120))
			require.NoError(t, p.Err)
			require.False(t, p.Outcomes[0].Members[0].Usable())
			require.Nil(t, p.Outcomes[0].Members[0].Counter.Delta)
			for _, cp := range p.CounterCheckpoints() {
				require.NotEqual(t, "unscaled-tail", cp.Reading.SampleID)
				require.Equal(t, measurement.ExactDecimal, cp.Reading.Value.Type())
			}
		})
	}
}
