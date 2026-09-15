package runtime

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLegacyMetricsSaturateWithoutOverflow(t *testing.T) {
	svc := &Service{}
	svc.collectedTotal.Store(math.MaxUint64)
	svc.writeSuccess.Store(math.MaxInt64 + 1)
	svc.writeError.Store(math.MaxUint64)
	svc.mappingError.Store(math.MaxUint64)
	svc.pointStateError.Store(1)
	metrics := svc.Metrics()
	require.Equal(t, int64(math.MaxInt64), metrics.TotalReads)
	require.Equal(t, int64(math.MaxInt64), metrics.TotalWrites)
	require.Equal(t, int64(math.MaxInt64), metrics.ErrorCount)
	require.Equal(t, int64(math.MaxInt64), metrics.MappingErrorCount)
	require.Equal(t, int64(math.MaxInt64), metrics.WriteErrorCount)
	require.Equal(t, uint64(math.MaxUint64), svc.Snapshot().CollectedTotal)
	require.Equal(t, int64(6), signedMetric(1, 2, 3))
}
