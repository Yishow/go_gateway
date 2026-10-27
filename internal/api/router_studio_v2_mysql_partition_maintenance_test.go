package api

import (
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestMySQLPartitionMaintenanceIntegration(t *testing.T) {
	runMySQLPartitionMaintenanceIntegration(t, mysqlCircuitFamilyMySQL)
}

func runMySQLPartitionMaintenanceIntegration(t *testing.T, family string) {
	f := newMySQLCircuitPlainFixture(t, mysqlCircuitPlainFixtureOptions{
		Family:          family,
		SeedMeasurement: measurement.MeasurementDefinition{ExactRecordingScale: "0.1"},
		Inputs:          mysqlCircuitDecimalInputs, StorageMode: workspace.WriteGroupStorageModeNativeMonthly,
	})
	f.faultTargets.failOpen.Store(1_000_000)
	f.accept(0)
	f.closeAfter(10 * time.Minute)
	requireEventuallyCumulativeRetry(t, f)
	assertMySQLPartitionJournalSample(t, f, mysqlCircuitDecimalBaselineID, mysqlCircuitDecimalBaseline)
	require.NoError(t, f.pipe.Stop(5*time.Second))
	oldPayloads := mysqlPartitionPayloadDigests(t, f)
	require.Len(t, oldPayloads, 41)
	assertMySQLPartitionRemoteCounts(t, f, 0)
	immutable := map[string][]byte{}
	for _, profile := range []mysqlCircuitProfile{f.instant, f.cumulative} {
		payload, err := f.store.RuntimeVersion(t.Context(), profile.key)
		require.NoError(t, err)
		immutable[profile.id] = payload
	}
	lineageBefore, err := recordingplan.ReadStorageLineage(t.Context(), f.f.db, workspace.StorageScopeOf(f.instant.group))
	require.NoError(t, err)
	base := lineageBefore.Layout.MySQL.NativeMonthly
	anchor, err := time.Parse(time.RFC3339, base.Partitions[0].LessThan)
	require.NoError(t, err)
	f.faultTargets.failOpen.Store(0)
	f.clock.Store(anchor.AddDate(0, 1, 2).UnixNano())
	require.NoError(t, f.pipe.MaintainMySQLPartitions(t.Context()))
	state, err := recordingplan.ReadMySQLPartitionMaintenanceState(t.Context(), f.f.db, workspace.StorageScopeOf(f.instant.group))
	require.NoError(t, err)
	require.NotNil(t, state)
	require.Len(t, state.Current.Partitions, len(base.Partitions)+1)
	require.Equal(t, "ready", state.Status)
	revision := state.Revision
	require.NoError(t, f.pipe.MaintainMySQLPartitions(t.Context()))
	state, err = recordingplan.ReadMySQLPartitionMaintenanceState(t.Context(), f.f.db, state.Scope)
	require.NoError(t, err)
	require.Equal(t, revision, state.Revision, "the same horizon requires no repeated DDL or ledger transition")
	assertMySQLPartitionRemoteCounts(t, f, 0)
	f.clock.Store(f.base.Add(10 * time.Minute).UnixNano())
	require.NoError(t, f.pipe.Start(t.Context()))
	waitForMySQLCircuitOutboxes(t, f, 41)
	f.accept(1)
	f.closeAfter(20 * time.Minute)
	waitForMySQLCircuitOutboxes(t, f, 82)
	assertMySQLPartitionJournalSample(t, f, mysqlCircuitDecimalEndID, mysqlCircuitDecimalEnd)
	currentPayloads := mysqlPartitionPayloadDigests(t, f)
	for key, digest := range oldPayloads {
		require.Equal(t, digest, currentPayloads[key], "accepted payloads cannot be recomputed for a new partition layout")
	}
	require.NoError(t, f.pipe.Stop(5*time.Second))
	f.clock.Store(anchor.AddDate(10, 0, 2).UnixNano())
	require.NoError(t, f.pipe.MaintainMySQLPartitions(t.Context()))
	f.clock.Store(f.base.Add(20 * time.Minute).UnixNano())
	state, err = recordingplan.ReadMySQLPartitionMaintenanceState(t.Context(), f.f.db, state.Scope)
	require.NoError(t, err)
	require.Len(t, state.Current.Partitions, len(base.Partitions)+6, "ten years of downtime adds only one gap and four future months")
	require.Equal(t, "p_gap_"+anchor.AddDate(10, 0, 0).Format("200601"), state.Current.Partitions[len(base.Partitions)].Name)
	assertMySQLPartitionRemoteCounts(t, f, 82)
	for _, profile := range []mysqlCircuitProfile{f.instant, f.cumulative} {
		payload, err := f.store.RuntimeVersion(t.Context(), profile.key)
		require.NoError(t, err)
		require.Equal(t, immutable[profile.id], payload, "maintenance cannot edit frozen runtime snapshots")
		ready, err := f.f.groups.Readiness(t.Context(), profile.id)
		require.NoError(t, err)
		require.True(t, ready.SchemaReady, "%+v", ready.Issues)
		require.True(t, ready.Ready, "verified destination is ready: %+v", ready.Issues)
	}
	lineageAfter, err := recordingplan.ReadStorageLineage(t.Context(), f.f.db, state.Scope)
	require.NoError(t, err)
	require.Equal(t, lineageBefore, lineageAfter, "partition maintenance cannot change the original confirmation or generation")
	rows := assertMySQLCircuitDecimalOutboxAndReceipts(t, f)
	assertMySQLCircuitDecimalSQLRows(t, f, rows)
}

func assertMySQLPartitionJournalSample(t *testing.T, f *mysqlCircuitPlainFixture, sampleID, expected string) {
	t.Helper()
	// Use the real journal bytes captured synchronously at intake; consumed
	// samples may already have been reclaimed by the existing worker.
	for _, profile := range []mysqlCircuitProfile{f.instant, f.cumulative} {
		var payload string
		for _, entry := range f.acceptedJournal {
			if entry.groupID == profile.id && entry.sampleID == sampleID {
				payload = entry.payload
			}
		}
		require.NotEmpty(t, payload)
		var sample mysqlCircuitPlainJournalPayload
		require.NoError(t, json.Unmarshal([]byte(payload), &sample))
		require.Equal(t, measurement.ExactDecimal, sample.Value.Type())
		digits, ok := sample.Value.DecimalDigits()
		require.True(t, ok)
		require.Equal(t, expected, digits)
	}
}

func mysqlPartitionPayloadDigests(t *testing.T, f *mysqlCircuitPlainFixture) map[string]string {
	t.Helper()
	rows, err := f.f.db.QueryContext(t.Context(), `SELECT effect_key, payload_digest FROM wg_delivery_outbox
		WHERE group_id IN (?, ?)`, f.instant.id, f.cumulative.id)
	require.NoError(t, err)
	defer rows.Close()
	result := map[string]string{}
	for rows.Next() {
		var key, digest string
		require.NoError(t, rows.Scan(&key, &digest))
		result[key] = digest
	}
	require.NoError(t, rows.Err())
	return result
}

func assertMySQLPartitionRemoteCounts(t *testing.T, f *mysqlCircuitPlainFixture, expected int) {
	t.Helper()
	for _, table := range []string{f.table, dbtarget.EffectReceiptTable} {
		var count int
		require.NoError(t, f.f.target.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM `"+table+"`").Scan(&count))
		require.Equal(t, expected, count)
	}
}
