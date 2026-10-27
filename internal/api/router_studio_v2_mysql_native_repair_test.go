package api

import (
	"database/sql"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// TestMySQLNativeMonthlyRepair validates that dropping an active native_monthly
// group table triggers automatic storage repair: generation advances to 2,
// HistoryMayBeLost is marked, the partition sidecar resets to the immutable base 6
// partitions, receipts remain durable (82 total), only pending rows commit (41 actual
// rows), and cumulative exact decimal counter values/deltas read back accurately.
func TestMySQLNativeMonthlyRepair(t *testing.T) {
	for _, family := range []string{mysqlCircuitFamilyMySQL, mysqlCircuitFamilyMariaDB} {
		t.Run(family, func(t *testing.T) {
			runMySQLNativeMonthlyRepair(t, family)
		})
	}
}

func runMySQLNativeMonthlyRepair(t *testing.T, family string) {
	f := newMySQLCircuitPlainFixture(t, mysqlCircuitPlainFixtureOptions{
		Family:          family,
		SeedMeasurement: measurement.MeasurementDefinition{ExactRecordingScale: "0.1"},
		Inputs:          mysqlCircuitDecimalInputs,
		StorageMode:     workspace.WriteGroupStorageModeNativeMonthly,
	})

	scope := workspace.StorageScopeOf(f.instant.group)
	lineageBefore, err := recordingplan.ReadStorageLineage(t.Context(), f.f.db, scope)
	require.NoError(t, err)
	require.NotNil(t, lineageBefore)
	require.Equal(t, 1, lineageBefore.Generation)
	require.False(t, lineageBefore.HistoryMayBeLost)
	base := lineageBefore.Layout.MySQL.NativeMonthly
	require.NotNil(t, base)
	require.Len(t, base.Partitions, 6)

	// Step 1: accept first sample, close 10m window, wait for 41 committed outbox rows.
	f.accept(0)
	f.closeAfter(10 * time.Minute)
	waitForMySQLCircuitOutboxes(t, f, 41)
	assertMySQLPartitionRemoteCounts(t, f, 41)

	oldPayloads := mysqlPartitionPayloadDigests(t, f)
	require.Len(t, oldPayloads, 41)
	immutableSnapshots := make(map[string][]byte, 2)
	for _, profile := range []mysqlCircuitProfile{f.instant, f.cumulative} {
		payload, err := f.store.RuntimeVersion(t.Context(), profile.key)
		require.NoError(t, err)
		immutableSnapshots[profile.id] = payload
	}

	// Step 2: run finite maintenance at next UTC month, expect sidecar partitions == 7.
	anchor, err := time.Parse(time.RFC3339, base.Partitions[0].LessThan)
	require.NoError(t, err)
	f.clock.Store(anchor.AddDate(0, 1, 2).UnixNano())
	require.NoError(t, f.pipe.MaintainMySQLPartitions(t.Context()))

	stateBeforeDrop, err := recordingplan.ReadMySQLPartitionMaintenanceState(t.Context(), f.f.db, scope)
	require.NoError(t, err)
	require.NotNil(t, stateBeforeDrop)
	require.Equal(t, "ready", stateBeforeDrop.Status)
	require.Equal(t, 1, stateBeforeDrop.Generation)
	require.Len(t, stateBeforeDrop.Current.Partitions, 7)

	f.clock.Store(f.base.Add(10 * time.Minute).UnixNano())

	// Step 3: drop ONLY f.table owned group table; leave receipt table intact.
	_, err = f.f.target.ExecContext(t.Context(), "DROP TABLE `"+f.table+"`")
	require.NoError(t, err)

	var receiptCountBefore int
	require.NoError(t, f.f.target.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM `"+dbtarget.EffectReceiptTable+"`").Scan(&receiptCountBefore))
	require.Equal(t, 41, receiptCountBefore)

	// Step 4: accept second sample, close 20m window.
	f.accept(1)
	f.closeAfter(20 * time.Minute)

	// A bounded wait verifies that repair settles rather than merely queuing DDL.
	require.Eventually(t, func() bool {
		lineage, err := recordingplan.ReadStorageLineage(t.Context(), f.f.db, scope)
		if err != nil || lineage == nil {
			return false
		}
		state, err := recordingplan.ReadMySQLPartitionMaintenanceState(t.Context(), f.f.db, scope)
		if err != nil || state == nil {
			return false
		}
		return lineage.Generation == 2 && lineage.HistoryMayBeLost && state.Generation == 2 && state.Status == "ready"
	}, 8*time.Second, 100*time.Millisecond, "native storage repair must advance to generation 2 and ready state")

	waitForMySQLCircuitOutboxes(t, f, 82)

	lineageAfter, err := recordingplan.ReadStorageLineage(t.Context(), f.f.db, scope)
	require.NoError(t, err)
	require.Equal(t, 2, lineageAfter.Generation)
	require.True(t, lineageAfter.HistoryMayBeLost)
	require.Equal(t, lineageBefore.LayoutDigest, lineageAfter.LayoutDigest, "initial layout digest unchanged")
	require.Equal(t, lineageBefore.ConfirmedDigest, lineageAfter.ConfirmedDigest)
	require.Equal(t, lineageBefore.Layout, lineageAfter.Layout)

	stateAfterRepair, err := recordingplan.ReadMySQLPartitionMaintenanceState(t.Context(), f.f.db, scope)
	require.NoError(t, err)
	require.NotNil(t, stateAfterRepair)
	require.Equal(t, 2, stateAfterRepair.Generation)
	require.Len(t, stateAfterRepair.Current.Partitions, 6, "sidecar partitions reset to initial base 6")

	var totalReceipts int
	require.NoError(t, f.f.target.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM `"+dbtarget.EffectReceiptTable+"`").Scan(&totalReceipts))
	require.Equal(t, 82, totalReceipts)

	var actualRows, instantRows, cumulativeRows int
	require.NoError(t, f.f.target.QueryRowContext(t.Context(), `
		SELECT
			COUNT(*),
			COUNT(CASE WHEN profile_kind = 'instantaneous' THEN 1 END),
			COUNT(CASE WHEN profile_kind = 'cumulative' THEN 1 END)
		FROM `+"`"+f.table+"`").Scan(&actualRows, &instantRows, &cumulativeRows))
	require.Equal(t, 41, actualRows, "previous history explicitly not restored into rebuilt table")
	require.Equal(t, 40, instantRows)
	require.Equal(t, 1, cumulativeRows)

	var circuitID, deviceID, pointID, tagID string
	var measurementID, seriesEpoch, unit, rawVal, deltaVal sql.NullString
	var status string
	require.NoError(t, f.f.target.QueryRowContext(t.Context(), `
		SELECT circuit_id, device_id, point_id, tag_id, measurement_id, series_epoch, unit, value_decimal, delta_value, status
		FROM `+"`"+f.table+"`"+` WHERE group_id = ? AND profile_kind = 'cumulative'`, f.cumulative.id).Scan(
		&circuitID, &deviceID, &pointID, &tagID, &measurementID, &seriesEpoch, &unit, &rawVal, &deltaVal, &status))
	require.Equal(t, f.cumulative.group.Members[0].CircuitID, circuitID)
	require.Equal(t, f.cumulative.group.Members[0].DeviceID, deviceID)
	require.Equal(t, f.cumulative.group.Members[0].PointID, pointID)
	require.Equal(t, f.cumulative.group.Members[0].TagID, tagID)
	require.True(t, measurementID.Valid)
	require.Equal(t, f.definition.ID, measurementID.String)
	require.True(t, seriesEpoch.Valid)
	require.Equal(t, f.definition.SeriesEpoch, seriesEpoch.String)
	require.True(t, unit.Valid)
	require.Equal(t, "kWh", unit.String)
	require.Equal(t, "ok", status)
	assertMySQLCircuitSQLDecimal(t, rawVal, mysqlCircuitDecimalEnd)
	assertMySQLCircuitSQLDecimal(t, deltaVal, "0.4")

	currentPayloads := mysqlPartitionPayloadDigests(t, f)
	for key, digest := range oldPayloads {
		require.Equal(t, digest, currentPayloads[key], "old accepted outbox payloads unchanged")
	}
	for _, profile := range []mysqlCircuitProfile{f.instant, f.cumulative} {
		payload, err := f.store.RuntimeVersion(t.Context(), profile.key)
		require.NoError(t, err)
		require.Equal(t, immutableSnapshots[profile.id], payload, "frozen runtime snapshots unchanged")

		ready, err := f.f.groups.Readiness(t.Context(), profile.id)
		require.NoError(t, err)
		require.True(t, ready.SchemaReady, "%+v", ready.Issues)
		require.True(t, ready.Ready, "verified destination is ready: %+v", ready.Issues)
	}
}
