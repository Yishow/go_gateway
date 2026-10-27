package dbtarget

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func transitionPartitions(names ...string) *MySQLPartitionEvidence {
	boundaries := map[string]string{
		"p_202610": "2026-11-01T00:00:00Z",
		"p_202611": "2026-12-01T00:00:00Z",
	}
	parts := make([]MySQLPartitionBoundary, 0, len(names)+1)
	for i, name := range names {
		parts = append(parts, MySQLPartitionBoundary{Name: name, LessThan: boundaries[name], Ordinal: i + 1})
	}
	parts = append(parts, MySQLPartitionBoundary{Name: "p_future", LessThan: "MAXVALUE", Ordinal: len(parts) + 1})
	return &MySQLPartitionEvidence{Method: "RANGE COLUMNS", Expression: "bucket_start", Partitions: parts}
}

func transitionProof(table string, partitions *MySQLPartitionEvidence) *MySQLManagedInspection {
	return &MySQLManagedInspection{
		Inspection: &TableInspection{
			Status: TableInspectionExists, Schema: "gw", Table: table, Indexes: []IndexInfo{}, IndexesInspected: true,
			Columns: []ColumnInfo{{Name: "bucket_start", DataType: "datetime(6)", PrimaryKey: true}},
		},
		Evidence: &MySQLManagedTableEvidence{
			ServerFamily: "mysql", ServerVersion: "8.4.5", CapabilityRevision: "cap", Database: "gw",
			Engine: "InnoDB", Charset: "utf8mb4", Collation: "utf8mb4_0900_bin", TableType: "BASE TABLE",
			PrimaryKey: []string{"bucket_start"}, ColumnCollations: map[string]string{},
			ColumnDefaults: map[string]*string{"bucket_start": nil}, JSONColumns: map[string]bool{},
			Partitions: partitions,
		},
		Indexes: []IndexInfo{},
	}
}

func transitionFacts(partitions *MySQLPartitionEvidence) mysqlManagedTableFacts {
	return mysqlManagedTableFacts{
		status: TableInspectionExists, tableType: "BASE TABLE", engine: "InnoDB",
		charset: "utf8mb4", collation: "utf8mb4_0900_bin",
		columns:    []mysqlManagedColumnFacts{{name: "bucket_start", dataType: "datetime", columnType: "datetime(6)", primaryKey: true}},
		primaryKey: []string{"bucket_start"}, jsonColumns: map[string]bool{}, partitions: partitions,
	}
}

func TestMySQLManagedPartitionOnlyTransitionIsTransient(t *testing.T) {
	expected := transitionProof("gw_long_abc", transitionPartitions("p_202610"))
	err := matchMySQLManagedInspection(expected, transitionFacts(transitionPartitions("p_202610", "p_202611")))
	require.ErrorIs(t, err, ErrMySQLManagedPartitionChanged)
	require.False(t, errors.Is(err, ErrGroupInsertUnsupported))
	require.False(t, errors.Is(err, ErrMySQLManagedPhysicalMismatch))
	require.Equal(t, InsertErrorTransient, ClassifyInsertError(err))
	require.False(t, IsSchemaShapeError(err))
}

func TestMySQLManagedPartitionTransitionKeepsOtherDriftBlocked(t *testing.T) {
	before := transitionPartitions("p_202610")
	after := transitionPartitions("p_202610", "p_202611")

	t.Run("combined column drift", func(t *testing.T) {
		expected := transitionProof("gw_long_abc", before)
		facts := transitionFacts(after)
		facts.columns[0].columnType = "varchar(64)"
		err := matchMySQLManagedInspection(expected, facts)
		require.False(t, errors.Is(err, ErrMySQLManagedPartitionChanged))
		require.ErrorIs(t, err, ErrGroupInsertUnsupported)
		require.ErrorIs(t, err, ErrMySQLManagedPhysicalMismatch)
	})

	t.Run("receipt partitioned", func(t *testing.T) {
		expected := transitionProof(EffectReceiptTable, before)
		err := matchMySQLManagedInspection(expected, transitionFacts(after))
		require.False(t, errors.Is(err, ErrMySQLManagedPartitionChanged))
		require.ErrorIs(t, err, ErrGroupInsertUnsupported)
	})

	t.Run("plain to native", func(t *testing.T) {
		expected := transitionProof("gw_long_abc", nil)
		err := matchMySQLManagedInspection(expected, transitionFacts(after))
		require.False(t, errors.Is(err, ErrMySQLManagedPartitionChanged))
		require.ErrorIs(t, err, ErrGroupInsertUnsupported)
	})

	t.Run("native to plain", func(t *testing.T) {
		expected := transitionProof("gw_long_abc", before)
		err := matchMySQLManagedInspection(expected, transitionFacts(nil))
		require.False(t, errors.Is(err, ErrMySQLManagedPartitionChanged))
		require.ErrorIs(t, err, ErrGroupInsertUnsupported)
	})
}
