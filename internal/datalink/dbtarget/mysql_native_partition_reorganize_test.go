package dbtarget

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func reorganizeProof() *MySQLManagedInspection {
	return &MySQLManagedInspection{
		Inspection: &TableInspection{
			Status: TableInspectionExists, Schema: "gw", Table: "gw_long_abc",
			Columns: []ColumnInfo{{Name: "bucket_start", DataType: "datetime(6)", PrimaryKey: true}}, IndexesInspected: true,
		},
		Evidence: &MySQLManagedTableEvidence{
			ServerFamily: "mysql", ServerVersion: "8.4.5", CapabilityRevision: "cap", Database: "gw",
			Engine: "InnoDB", Charset: "utf8mb4", Collation: "utf8mb4_0900_bin", TableType: "BASE TABLE",
			PrimaryKey: []string{"bucket_start", "record_id"},
			Partitions: &MySQLPartitionEvidence{
				Method: "RANGE COLUMNS", Expression: "bucket_start",
				Partitions: []MySQLPartitionBoundary{
					{Name: "p_202610", LessThan: "2026-11-01T00:00:00Z", Ordinal: 1},
					{Name: "p_future", LessThan: "MAXVALUE", Ordinal: 2},
				},
			},
		},
		Indexes: []IndexInfo{},
	}
}

func reorganizeAfter(names ...string) *MySQLPartitionEvidence {
	boundaries := map[string]string{
		"p_202611":     "2026-12-01T00:00:00Z",
		"p_202612":     "2027-01-01T00:00:00Z",
		"p_202701":     "2027-02-01T00:00:00Z",
		"p_202702":     "2027-03-01T00:00:00Z",
		"p_202703":     "2027-04-01T00:00:00Z",
		"p_gap_202612": "2026-12-01T00:00:00Z",
		"p_202610":     "2026-12-01T00:00:00Z",
	}
	partitions := make([]MySQLPartitionBoundary, 0, len(names)+2)
	partitions = append(partitions, MySQLPartitionBoundary{Name: "p_202610", LessThan: "2026-11-01T00:00:00Z", Ordinal: 1})
	for index, name := range names {
		partitions = append(partitions, MySQLPartitionBoundary{Name: name, LessThan: boundaries[name], Ordinal: index + 2})
	}
	partitions = append(partitions, MySQLPartitionBoundary{Name: "p_future", LessThan: "MAXVALUE", Ordinal: len(partitions) + 1})
	return &MySQLPartitionEvidence{Method: "RANGE COLUMNS", Expression: "bucket_start", Partitions: partitions}
}

func TestMySQLFuturePartitionReorganizeStatement(t *testing.T) {
	after := reorganizeAfter("p_202611", "p_202612")
	additions, err := validateMySQLFuturePartitionPlan(reorganizeProof(), after, "gw")
	require.NoError(t, err)
	statement, err := mysqlFuturePartitionReorganizeStatement("gw", "gw_long_abc", additions)
	require.NoError(t, err)
	require.Equal(t,
		"ALTER TABLE `gw`.`gw_long_abc` REORGANIZE PARTITION `p_future` INTO ("+
			"PARTITION `p_202611` VALUES LESS THAN ('2026-12-01 00:00:00'), "+
			"PARTITION `p_202612` VALUES LESS THAN ('2027-01-01 00:00:00'), "+
			"PARTITION `p_future` VALUES LESS THAN (MAXVALUE));", statement)
	for _, forbidden := range []string{"DROP", "TRUNCATE", "MERGE", "RENAME"} {
		require.NotContains(t, strings.ToUpper(statement), forbidden)
	}

	gap, err := validateMySQLFuturePartitionPlan(reorganizeProof(), reorganizeAfter("p_gap_202612", "p_202612"), "gw")
	require.NoError(t, err)
	require.Len(t, gap, 2)
	require.Equal(t, "p_gap_202612", gap[0].Name)
}

func TestMySQLFuturePartitionReorganizeRefusesPlanMutations(t *testing.T) {
	cases := map[string]func(*MySQLManagedInspection, *MySQLPartitionEvidence){
		"prefix name": func(_ *MySQLManagedInspection, after *MySQLPartitionEvidence) { after.Partitions[0].Name = "p_other" },
		"prefix boundary": func(_ *MySQLManagedInspection, after *MySQLPartitionEvidence) {
			after.Partitions[0].LessThan = "2026-10-01T00:00:00Z"
		},
		"prefix ordinal": func(_ *MySQLManagedInspection, after *MySQLPartitionEvidence) { after.Partitions[0].Ordinal = 9 },
		"missing prefix": func(_ *MySQLManagedInspection, after *MySQLPartitionEvidence) {
			after.Partitions = after.Partitions[1:]
		},
		"missing final future": func(_ *MySQLManagedInspection, after *MySQLPartitionEvidence) {
			after.Partitions = after.Partitions[:len(after.Partitions)-1]
		},
		"final not maxvalue": func(_ *MySQLManagedInspection, after *MySQLPartitionEvidence) {
			after.Partitions[len(after.Partitions)-1].LessThan = "2027-02-01T00:00:00Z"
		},
		"foreign method":     func(_ *MySQLManagedInspection, after *MySQLPartitionEvidence) { after.Method = "HASH" },
		"foreign expression": func(_ *MySQLManagedInspection, after *MySQLPartitionEvidence) { after.Expression = "other" },
		"non first of month": func(_ *MySQLManagedInspection, after *MySQLPartitionEvidence) {
			after.Partitions[1].LessThan = "2026-12-02T00:00:00Z"
		},
		"offset boundary": func(_ *MySQLManagedInspection, after *MySQLPartitionEvidence) {
			after.Partitions[1].LessThan = "2026-12-01T00:00:00+00:00"
		},
		"invalid addition name": func(_ *MySQLManagedInspection, after *MySQLPartitionEvidence) { after.Partitions[1].Name = "p_bad" },
		"duplicate addition":    func(_ *MySQLManagedInspection, after *MySQLPartitionEvidence) { after.Partitions[2].Name = "p_202611" },
		"addition ordinal":      func(_ *MySQLManagedInspection, after *MySQLPartitionEvidence) { after.Partitions[1].Ordinal = 9 },
		"not increasing": func(_ *MySQLManagedInspection, after *MySQLPartitionEvidence) {
			after.Partitions[2].LessThan = "2026-12-01T00:00:00Z"
		},
		"too many months": func(_ *MySQLManagedInspection, after *MySQLPartitionEvidence) {
			*after = *reorganizeAfter("p_202611", "p_202612", "p_202701", "p_202702", "p_202703")
		},
		"two gaps": func(_ *MySQLManagedInspection, after *MySQLPartitionEvidence) {
			*after = *reorganizeAfter("p_gap_202612", "p_gap_202612", "p_202701")
		},
		"gap not first": func(_ *MySQLManagedInspection, after *MySQLPartitionEvidence) {
			*after = *reorganizeAfter("p_202611", "p_gap_202612")
		},
		"gap only": func(_ *MySQLManagedInspection, after *MySQLPartitionEvidence) {
			*after = *reorganizeAfter("p_gap_202612")
		},
		"foreign table":            func(proof *MySQLManagedInspection, _ *MySQLPartitionEvidence) { proof.Inspection.Table = "other_table" },
		"foreign database":         func(proof *MySQLManagedInspection, _ *MySQLPartitionEvidence) { proof.Evidence.Database = "other" },
		"missing proof partitions": func(proof *MySQLManagedInspection, _ *MySQLPartitionEvidence) { proof.Evidence.Partitions = nil },
		"before final renamed": func(proof *MySQLManagedInspection, _ *MySQLPartitionEvidence) {
			proof.Evidence.Partitions.Partitions[1].Name = "p_max"
		},
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			proof := reorganizeProof()
			after := reorganizeAfter("p_202611", "p_202612")
			mutate(proof, after)
			_, err := validateMySQLFuturePartitionPlan(proof, after, "gw")
			require.Error(t, err)
		})
	}
}

func TestMySQLFuturePartitionReorganizeInputValidation(t *testing.T) {
	err := ReorganizeOpenMySQLFuturePartition(t.Context(), nil, nil, "gw", nil, nil)
	require.Error(t, err)
	require.False(t, strings.Contains(err.Error(), "ALTER"))
}

func TestMySQLFuturePartitionFirstOfMonthCanonical(t *testing.T) {
	expected, err := parseMySQLFirstOfMonth("2026-12-01T00:00:00Z")
	require.NoError(t, err)
	require.Equal(t, "2026-12-01T00:00:00Z", expected.Format(time.RFC3339Nano))
	for _, invalid := range []string{"2026-12-01T00:00:01Z", "2026-12-02T00:00:00Z", "2026-12-01T00:00:00+00:00", "MAXVALUE"} {
		_, err := parseMySQLFirstOfMonth(invalid)
		require.Error(t, err, invalid)
	}
}

func TestMySQLFuturePartitionPlanRejectsPatternOnlyFalsePositives(t *testing.T) {
	cases := map[string]func(*MySQLManagedInspection, *MySQLPartitionEvidence){
		"month name from unrelated month": func(_ *MySQLManagedInspection, after *MySQLPartitionEvidence) {
			after.Partitions[1].Name = "p_209901"
		},
		"skipped month": func(_ *MySQLManagedInspection, after *MySQLPartitionEvidence) {
			after.Partitions[2].Name = "p_202701"
			after.Partitions[2].LessThan = "2027-02-01T00:00:00Z"
		},
		"gap name from unrelated month": func(_ *MySQLManagedInspection, after *MySQLPartitionEvidence) {
			*after = *reorganizeAfter("p_gap_202612", "p_202612")
			after.Partitions[1].Name = "p_gap_209901"
		},
		"addition collides with prefix": func(_ *MySQLManagedInspection, after *MySQLPartitionEvidence) {
			*after = *reorganizeAfter("p_202610", "p_202611")
		},
		"before method": func(proof *MySQLManagedInspection, _ *MySQLPartitionEvidence) {
			proof.Evidence.Partitions.Method = "HASH"
		},
		"before expression": func(proof *MySQLManagedInspection, _ *MySQLPartitionEvidence) {
			proof.Evidence.Partitions.Expression = "other"
		},
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			proof := reorganizeProof()
			after := reorganizeAfter("p_202611", "p_202612")
			mutate(proof, after)
			_, err := validateMySQLFuturePartitionPlan(proof, after, "gw")
			require.Error(t, err)
		})
	}

	t.Run("year below mysql datetime", func(t *testing.T) {
		proof := reorganizeProof()
		proof.Evidence.Partitions.Partitions = []MySQLPartitionBoundary{
			{Name: "p_099910", LessThan: "0999-11-01T00:00:00Z", Ordinal: 1},
			{Name: "p_future", LessThan: "MAXVALUE", Ordinal: 2},
		}
		after := &MySQLPartitionEvidence{Method: "RANGE COLUMNS", Expression: "bucket_start", Partitions: []MySQLPartitionBoundary{
			{Name: "p_099910", LessThan: "0999-11-01T00:00:00Z", Ordinal: 1},
			{Name: "p_099911", LessThan: "0999-12-01T00:00:00Z", Ordinal: 2},
			{Name: "p_future", LessThan: "MAXVALUE", Ordinal: 3},
		}}
		_, err := validateMySQLFuturePartitionPlan(proof, after, "gw")
		require.Error(t, err)
	})
}
