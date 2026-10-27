package api

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// TestMySQLCircuitNativeMonthlyIntegration uses simulated UTC event time across
// a month boundary and actual MySQL partition placement and SQL readback.
func TestMySQLCircuitNativeMonthlyIntegration(t *testing.T) {
	runMySQLCircuitNativeMonthlyIntegration(t, mysqlCircuitFamilyMySQL)
}

func runMySQLCircuitNativeMonthlyIntegration(t *testing.T, family string) {
	now := time.Now().UTC()
	next := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	f := newMySQLCircuitPlainFixture(t, mysqlCircuitPlainFixtureOptions{
		Family:          family,
		SeedMeasurement: measurement.MeasurementDefinition{ExactRecordingScale: "0.1"},
		Inputs:          mysqlCircuitDecimalInputs, StorageMode: workspace.WriteGroupStorageModeNativeMonthly,
		BaseTime: next.Add(-10 * time.Minute),
	})
	for _, profile := range []mysqlCircuitProfile{f.instant, f.cumulative} {
		require.Equal(t, workspace.WriteGroupStorageModeNativeMonthly, profile.group.StorageMode)
		proof := profile.snapshot.RuntimeLayout.MySQL
		require.NotNil(t, proof.Group.Evidence.Partitions)
		require.Len(t, proof.Group.Evidence.Partitions.Partitions, 6)
		require.Nil(t, proof.Receipt.Evidence.Partitions)
	}
	f.accept(0)
	f.closeAfter(10 * time.Minute)
	waitForMySQLCircuitOutboxes(t, f, 41)
	f.accept(1)
	f.closeAfter(20 * time.Minute)
	waitForMySQLCircuitOutboxes(t, f, 82)
	assertMySQLCircuitDecimalJournal(t, f)
	rows := assertMySQLCircuitDecimalOutboxAndReceipts(t, f)
	assertMySQLCircuitDecimalSQLRows(t, f, rows)

	for _, month := range []time.Time{next.AddDate(0, -1, 0), next} {
		partition := fmt.Sprintf("p_%04d%02d", month.Year(), month.Month())
		var count int
		var earliest, latest time.Time
		err := f.f.target.QueryRowContext(t.Context(), "SELECT COUNT(*), MIN(bucket_start), MAX(bucket_start) FROM `"+
			f.table+"` PARTITION (`"+partition+"`)").Scan(&count, &earliest, &latest)
		require.NoError(t, err)
		require.Equal(t, 41, count, "each UTC month holds forty instantaneous rows and one cumulative row")
		require.Equal(t, month.Month(), earliest.UTC().Month())
		require.Equal(t, month.Month(), latest.UTC().Month())
	}
	var receiptPartition sql.NullString
	require.NoError(t, f.f.target.QueryRowContext(t.Context(), `SELECT PARTITION_NAME FROM information_schema.PARTITIONS
		WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME=?`, dbtarget.EffectReceiptTable).Scan(&receiptPartition))
	require.False(t, receiptPartition.Valid, "receipts stay in an ordinary InnoDB table")
}
