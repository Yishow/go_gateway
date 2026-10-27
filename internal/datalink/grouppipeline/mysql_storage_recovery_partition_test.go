package grouppipeline

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
)

func TestMySQLStorageRecovery_MirrorMySQLPartitionEvidence_PreservesNil(t *testing.T) {
	require.Nil(t, mirrorMySQLPartitionEvidence(nil), "plain table evidence must preserve nil partition pointer")
}

func TestMySQLStorageRecovery_MirrorMySQLPartitionEvidence_MirrorsOrderedBoundariesAndIsolates(t *testing.T) {
	original := &dbtarget.MySQLPartitionEvidence{
		Method:     "RANGE COLUMNS",
		Expression: "bucket_start",
		Partitions: []dbtarget.MySQLPartitionBoundary{
			{Name: "p_past", LessThan: "2026-10-01T00:00:00Z", Ordinal: 1},
			{Name: "p_2026_10", LessThan: "2026-11-01T00:00:00Z", Ordinal: 2},
			{Name: "p_future", LessThan: "MAXVALUE", Ordinal: 3},
		},
	}

	mirrored := mirrorMySQLPartitionEvidence(original)
	require.NotNil(t, mirrored)
	require.Equal(t, "RANGE COLUMNS", mirrored.Method)
	require.Equal(t, "bucket_start", mirrored.Expression)
	require.Len(t, mirrored.Partitions, 3)

	expected := []recordingplan.MySQLPartitionBoundary{
		{Name: "p_past", LessThan: "2026-10-01T00:00:00Z", Ordinal: 1},
		{Name: "p_2026_10", LessThan: "2026-11-01T00:00:00Z", Ordinal: 2},
		{Name: "p_future", LessThan: "MAXVALUE", Ordinal: 3},
	}
	require.Equal(t, expected, mirrored.Partitions)

	// Verify slice and pointer isolation
	original.Partitions[0].Name = "mutated"
	original.Partitions = append(original.Partitions, dbtarget.MySQLPartitionBoundary{Name: "extra"})
	require.Equal(t, "p_past", mirrored.Partitions[0].Name, "mirrored partitions slice must not be modified by original slice mutations")
	require.Len(t, mirrored.Partitions, 3)
}

func TestMySQLRestoreTargetInspection_MirrorsPartitions(t *testing.T) {
	inspection := &dbtarget.MySQLManagedInspection{
		Inspection: &dbtarget.TableInspection{
			Status: dbtarget.TableInspectionExists,
			Schema: "gwtest",
			Table:  "gw_long_test",
		},
		Evidence: &dbtarget.MySQLManagedTableEvidence{
			ServerFamily:       "mysql",
			ServerVersion:      "8.4.13",
			CapabilityRevision: "cap-1",
			Database:           "gwtest",
			Engine:             "InnoDB",
			Charset:            "utf8mb4",
			Collation:          "utf8mb4_0900_bin",
			TableType:          "BASE TABLE",
			PrimaryKey:         []string{"bucket_start"},
			ColumnCollations:   map[string]string{},
			ColumnDefaults:     map[string]*string{},
			JSONColumns:        map[string]bool{},
			Partitions: &dbtarget.MySQLPartitionEvidence{
				Method:     "RANGE COLUMNS",
				Expression: "bucket_start",
				Partitions: []dbtarget.MySQLPartitionBoundary{
					{Name: "p_past", LessThan: "2026-10-01T00:00:00Z", Ordinal: 1},
					{Name: "p_current", LessThan: "2026-11-01T00:00:00Z", Ordinal: 2},
				},
			},
		},
	}

	state := mysqlRestoreTargetInspection(inspection)
	require.NotNil(t, state.MySQL)
	require.NotNil(t, state.MySQL.Partitions)
	require.Equal(t, "RANGE COLUMNS", state.MySQL.Partitions.Method)
	require.Equal(t, "bucket_start", state.MySQL.Partitions.Expression)
	require.Len(t, state.MySQL.Partitions.Partitions, 2)
	require.Equal(t, "p_past", state.MySQL.Partitions.Partitions[0].Name)
	require.Equal(t, "p_current", state.MySQL.Partitions.Partitions[1].Name)
}
