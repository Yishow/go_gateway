package recordingplan

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSQLServerPartitionDowntimeClosesAdjacentMonthWithoutFalseGap(t *testing.T) {
	initial, err := NewSQLServerNativeMonthlyPartition(time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC), "pf_fixture", "ps_fixture")
	require.NoError(t, err)
	before, err := json.Marshal(initial)
	require.NoError(t, err)

	// The initial horizon ends May 1. Resuming June 10 closes the single May
	// interval, followed by June and the next three complete months.
	after, err := PlanSQLServerNativeMonthlyExtension(initial, time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	require.Len(t, after.Partitions, 11, "one adjacent boundary plus four month ends is bounded")
	require.Equal(t, SQLServerPartitionBoundary{
		Name: "p_202605", LessThan: "2026-06-01T00:00:00.000000", Ordinal: 6,
	}, after.Partitions[5])
	require.Equal(t, SQLServerPartitionBoundary{
		Name: "p_202609", LessThan: "2026-10-01T00:00:00.000000", Ordinal: 10,
	}, after.Partitions[9])
	require.Equal(t, SQLServerPartitionBoundary{Name: "p_future", LessThan: "MAXVALUE", Ordinal: 11}, after.Partitions[10])
	frozen, err := json.Marshal(initial)
	require.NoError(t, err)
	require.Equal(t, before, frozen, "maintenance cannot change the accepted initial proof")
}
