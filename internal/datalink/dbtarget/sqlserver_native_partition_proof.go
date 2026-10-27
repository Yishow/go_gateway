package dbtarget

import (
	"slices"
	"strings"
	"time"
)

const (
	sqlServerManagedPartitionBoundaryLayout = "2006-01-02T15:04:05.000000"
	sqlServerManagedPartitionMinimumRanges  = 6
	sqlServerManagedPartitionFutureBoundary = "MAXVALUE"
	sqlServerNativeRangeRight               = "RANGE RIGHT"
	sqlServerNativePastPartitionName        = "p_past"
	sqlServerNativeFuturePartitionName      = "p_future"
)

func sqlServerManagedNativeMonthlyEvidenceMatchesForTable(
	tableName string,
	evidence *SQLServerManagedPartitionEvidence,
) bool {
	return sqlServerManagedNativeMonthlyDefinitionMatchesForTable(tableName, evidence) &&
		evidence.FunctionExists && evidence.SchemeExists && evidence.Exclusive && evidence.IndexesAligned
}

func sqlServerManagedNativeMonthlyDefinitionMatchesForTable(
	tableName string,
	evidence *SQLServerManagedPartitionEvidence,
) bool {
	if strings.TrimSpace(tableName) == "" || evidence == nil || evidence.Method != sqlServerNativeRangeRight || evidence.Column != longColBucketStart ||
		evidence.FunctionName != "pf_"+tableName ||
		evidence.SchemeName != "ps_"+tableName || evidence.Filegroup != sqlServerNativePrimaryFilegroup ||
		len(evidence.Partitions) < sqlServerManagedPartitionMinimumRanges {
		return false
	}

	first := evidence.Partitions[0]
	previous, ok := sqlServerManagedPartitionMonthStart(first.LessThan)
	if !ok || first.Name != sqlServerNativePastPartitionName || first.Ordinal != 1 {
		return false
	}

	expected := make([]SQLServerPartitionBoundary, 0, len(evidence.Partitions))
	expected = append(expected, first)
	for index := 1; index < len(evidence.Partitions)-1; index++ {
		boundary := evidence.Partitions[index]
		at, valid := sqlServerManagedPartitionMonthStart(boundary.LessThan)
		if !valid || !at.After(previous) {
			return false
		}
		name := "p_" + previous.Format("200601")
		if !at.Equal(previous.AddDate(0, 1, 0)) {
			if !at.After(previous.AddDate(0, 1, 0)) {
				return false
			}
			name = "p_gap_" + at.Format("200601")
		}
		if boundary.Name != name || boundary.Ordinal != index+1 {
			return false
		}
		expected = append(expected, SQLServerPartitionBoundary{
			Name: name, LessThan: at.Format(sqlServerManagedPartitionBoundaryLayout), Ordinal: index + 1,
		})
		previous = at
	}
	expected = append(expected, SQLServerPartitionBoundary{
		Name: sqlServerNativeFuturePartitionName, LessThan: sqlServerManagedPartitionFutureBoundary, Ordinal: len(evidence.Partitions),
	})
	return slices.Equal(evidence.Partitions, expected)
}

func sqlServerManagedPartitionMonthStart(value string) (time.Time, bool) {
	monthStart, err := time.Parse(sqlServerManagedPartitionBoundaryLayout, value)
	if err != nil || monthStart.Format(sqlServerManagedPartitionBoundaryLayout) != value ||
		monthStart.Year() < 1 || monthStart.Year() > 9999 || monthStart.Day() != 1 ||
		monthStart.Hour() != 0 || monthStart.Minute() != 0 || monthStart.Second() != 0 || monthStart.Nanosecond() != 0 {
		return time.Time{}, false
	}
	return monthStart, true
}
