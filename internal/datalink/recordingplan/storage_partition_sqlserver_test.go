package recordingplan

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func sqlServerMaintenanceBaseEvidence(t *testing.T) *SQLServerPartitionEvidence {
	t.Helper()
	evidence, err := NewSQLServerNativeMonthlyPartition(time.Date(2026, time.October, 15, 0, 0, 0, 0, time.UTC),
		"pf_gw_long_0123456789abcdef", "ps_gw_long_0123456789abcdef")
	require.NoError(t, err)
	return evidence
}

func TestSQLServerPartitionRollover(t *testing.T) {
	base := sqlServerMaintenanceBaseEvidence(t)
	after, err := PlanSQLServerNativeMonthlyExtension(base, time.Date(2026, time.December, 20, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	require.Len(t, after.Partitions, len(base.Partitions)+2)
	require.Equal(t, base.Partitions[:len(base.Partitions)-1], after.Partitions[:len(base.Partitions)-1])
	require.Equal(t, "p_202702", after.Partitions[len(after.Partitions)-3].Name)
	require.Equal(t, "p_future", after.Partitions[len(after.Partitions)-1].Name)
	again, err := PlanSQLServerNativeMonthlyExtension(after, time.Date(2026, time.December, 20, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	require.Equal(t, after, again)
}

func TestSQLServerPartitionDowntimeIsBoundedAndGapNamed(t *testing.T) {
	base := sqlServerMaintenanceBaseEvidence(t)
	after, err := PlanSQLServerNativeMonthlyExtension(base, time.Date(2027, time.May, 20, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	require.Len(t, after.Partitions, len(base.Partitions)+5)
	require.Equal(t, "p_gap_202705", after.Partitions[len(base.Partitions)-1].Name)
	require.Equal(t, "2027-05-01T00:00:00.000000", after.Partitions[len(base.Partitions)-1].LessThan)
	require.Equal(t, "p_202705", after.Partitions[len(base.Partitions)].Name)
	require.Equal(t, "p_202708", after.Partitions[len(after.Partitions)-2].Name)
	require.Equal(t, "p_future", after.Partitions[len(after.Partitions)-1].Name)
	require.NoError(t, validateCurrentSQLServerPartitionEvidence(after))

	again, err := PlanSQLServerNativeMonthlyExtension(after, time.Date(2027, time.May, 20, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	require.Equal(t, after, again)
}

func TestSQLServerPartitionMaintenanceRejectsMalformedGap(t *testing.T) {
	base := sqlServerMaintenanceBaseEvidence(t)
	malformed := CloneSQLServerPartitionEvidence(base)
	malformed.Partitions[len(malformed.Partitions)-2].Name = "p_202601"
	require.Error(t, validateCurrentSQLServerPartitionEvidence(malformed))

	duplicate := CloneSQLServerPartitionEvidence(base)
	duplicate.Partitions[1].Name = duplicate.Partitions[2].Name
	require.Error(t, validateCurrentSQLServerPartitionEvidence(duplicate))
}
