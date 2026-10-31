package snapshot

import (
	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/measurement"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCounterStateCopyPreservesIdentityAndOwnsPolicy(t *testing.T) {
	binding := counterBinding("c")
	mixedRollover(true)(&binding)
	before, err := binding.Key()
	require.NoError(t, err)
	copied := cloneCounterBinding(binding)
	after, err := copied.Key()
	require.NoError(t, err)
	require.Equal(t, before, after)
	require.NotSame(t, binding.RolloverModulus, copied.RolloverModulus)
	require.NotSame(t, binding.MaxRatePerSecond, copied.MaxRatePerSecond)
	original := aggregation.CounterCheckpoint{Binding: binding}
	checkpoint := cloneCounterCheckpoint(original)
	*checkpoint.Binding.RolloverModulus = measurement.NewUint64(16)
	*checkpoint.Binding.MaxRatePerSecond = measurement.NewUint64(1)
	again, err := binding.Key()
	require.NoError(t, err)
	require.Equal(t, before, again)
	zero := cloneCounterBinding(counterBinding("plain"))
	require.Nil(t, zero.RolloverModulus)
	require.Nil(t, zero.MaxRatePerSecond)
}
