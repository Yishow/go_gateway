package snapshot

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/measurement"
	"testing"
)

func TestMixedCounterPolicyDoesNotAliasCallerOrCheckpointViews(t *testing.T) {
	for _, boundary := range []string{"constructor", "plan view", "committed view", "restore input"} {
		t.Run(boundary, func(t *testing.T) {
			c := mixedCounterConfig(ModeMixedWide)
			binding := c.Counters["c"]
			mixedRollover(true)(&binding)
			c.Counters["c"] = binding
			a, err := NewAssembler(c)
			require.NoError(t, err)
			if boundary == "constructor" {
				*binding.RolloverModulus = measurement.NewUint64(16)
			}
			mixedOffer(t, a, mixedIntegerSample("c", "baseline", 59, 65530))
			p := a.Plan(mixedTime(60))
			if boundary == "plan view" {
				view := p.CounterCheckpoints()
				*view["c"].Binding.RolloverModulus = measurement.NewUint64(16)
			}
			a.Commit(p)
			if boundary == "committed view" {
				view := a.CounterCheckpoints()
				*view["c"].Binding.RolloverModulus = measurement.NewUint64(16)
			}
			if boundary == "restore input" {
				checkpoint := a.CounterCheckpoints()
				fresh := mixedPolicyAssembler(t, ModeMixedWide, mixedRollover(true))
				progress, _ := a.MixedProgress()
				require.NoError(t, fresh.RestoreMixedProgress(progress))
				require.NoError(t, fresh.RestoreCounterCheckpoints(checkpoint))
				*checkpoint["c"].Binding.MaxRatePerSecond = measurement.NewUint64(1)
				a = fresh
			}
			mixedOffer(t, a, mixedIntegerSample("c", "tail", 119, 4))
			r := mixedPolicyResult(t, a.Plan(mixedTime(120)))
			require.True(t, r.Usable(), "external mutation must not rewrite frozen counter policy: %s", r.Reason)
			require.Equal(t, "10", mixedPolicyDigits(t, r))
		})
	}
}
