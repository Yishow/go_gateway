package grouppipeline

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

func TestSQLServerPartitionDefinitionMatchesExactAuthorizedBefore(t *testing.T) {
	expected, err := recordingplan.NewSQLServerNativeMonthlyPartition(
		time.Date(2026, time.October, 15, 0, 0, 0, 0, time.UTC), "pf_circuit_long", "ps_circuit_long",
	)
	require.NoError(t, err)
	actual := dbtargetPartitionDefinitionForTest(expected)
	require.True(t, sqlServerPartitionDefinitionMatches(actual, expected))
	for name, mutate := range map[string]func(*dbtarget.SQLServerManagedPartitionEvidence){
		"foreign function": func(value *dbtarget.SQLServerManagedPartitionEvidence) { value.FunctionName = "pf_foreign" },
		"foreign boundary": func(value *dbtarget.SQLServerManagedPartitionEvidence) {
			value.Partitions[1].LessThan = "2026-11-02T00:00:00.000000"
		},
		"short prefix": func(value *dbtarget.SQLServerManagedPartitionEvidence) {
			value.Partitions = value.Partitions[:len(value.Partitions)-1]
		},
	} {
		t.Run(name, func(t *testing.T) {
			candidate := dbtargetPartitionDefinitionForTest(expected)
			mutate(candidate)
			require.False(t, sqlServerPartitionDefinitionMatches(candidate, expected))
		})
	}
	require.False(t, sqlServerPartitionDefinitionMatches(nil, expected))
}

func dbtargetPartitionDefinitionForTest(expected *recordingplan.SQLServerPartitionEvidence) *dbtarget.SQLServerManagedPartitionEvidence {
	result := &dbtarget.SQLServerManagedPartitionEvidence{
		Method: expected.Method, Column: expected.Column, FunctionName: expected.FunctionName,
		SchemeName: expected.SchemeName, Filegroup: expected.Filegroup,
		FunctionExists: true, SchemeExists: true, Exclusive: true, IndexesAligned: true,
	}
	for _, partition := range expected.Partitions {
		result.Partitions = append(result.Partitions, dbtarget.SQLServerPartitionBoundary{
			Name: partition.Name, LessThan: partition.LessThan, Ordinal: partition.Ordinal,
		})
	}
	return result
}
