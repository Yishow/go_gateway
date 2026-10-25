package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// TestDeviceRecordingIntegration drives the fixed simulated matrix across 5 SQL modes:
// 1. new wide / plain SQLite
// 2. new wide / plain PostgreSQL
// 3. new mean long / plain SQLite
// 4. new mean long / plain PostgreSQL
// 5. existing managed long TimescaleDB
//
// Verifies arithmetic means (20, 1.5, 9007199254740992.5), NULL/partial/unknown coverage,
// counter baseline/delta/reset/gap, committed delivery replay, hypertable catalog
// proof, and isolated cleanup with stopped workers. The reviewed-cutover cases
// separately prove historical pending destinations across a transition.
func TestDeviceRecordingIntegration(t *testing.T) {
	t.Run("wide_plain_sqlite", func(t *testing.T) {
		runWideIntegrationMode(t, "sqlite")
	})

	t.Run("wide_plain_postgres", func(t *testing.T) {
		if os.Getenv("POSTGRES_DSN") == "" {
			t.Skip("owned PostgreSQL fixture not configured; mode is NOT RUN")
		}
		runWideIntegrationMode(t, "postgres")
	})

	t.Run("mean_long_plain_sqlite", func(t *testing.T) {
		runMeanLongIntegrationMode(t, "sqlite", "plain")
	})

	t.Run("mean_long_plain_postgres", func(t *testing.T) {
		if os.Getenv("POSTGRES_DSN") == "" {
			t.Skip("owned PostgreSQL fixture not configured; mode is NOT RUN")
		}
		runMeanLongIntegrationMode(t, "postgres", "plain")
	})

	t.Run("existing_long_timescale", func(t *testing.T) {
		if os.Getenv("POSTGRES_DSN") == "" {
			t.Skip("owned PostgreSQL fixture not configured; mode is NOT RUN")
		}
		runMeanLongIntegrationMode(t, "postgres", "timescale")
	})

	t.Run("reviewed_new_wide_history_pending", TestReviewedWideMeanTransition)
	t.Run("reviewed_long_history_pending", TestWideToLongCutover)
}

// runWideIntegrationMode tests the managed_wide_profile_v1 layout on SQLite or PostgreSQL.
func runWideIntegrationMode(t *testing.T, kind string) {
	t.Helper()
	ctx := t.Context()
	sources := []meanSource{
		{name: "half", tagType: "int16", role: "gauge"},
		{name: "big", tagType: "uint64", role: "gauge"},
		{name: "tri", tagType: "int16", role: "gauge"},
		{name: "flt", tagType: "float32", role: "gauge"},
		{name: "quiet", tagType: "int16", role: "gauge"},
		{name: "state", tagType: "int16", role: "state"},
		{name: "raw", tagType: "uint64", role: "counter_raw", measurementID: "meas-raw"},
	}
	rig := newWideRig(t, kind, sources)

	// Fixed simulated samples inside the first 15s window [0, 15)
	for _, sample := range []struct {
		source, id string
		offset     time.Duration
		value      any
	}{
		{"half", "half-1", time.Second, 1}, {"half", "half-2", 2 * time.Second, 2},
		{"big", "big-1", time.Second, uint64(9007199254740992)}, {"big", "big-2", 2 * time.Second, uint64(9007199254740993)},
		{"tri", "tri-1", time.Second, 10}, {"tri", "tri-2", 2 * time.Second, 20}, {"tri", "tri-3", 3 * time.Second, 30},
		{"tri", "tri-2", 2 * time.Second, 20}, // resend must not enlarge denominator count
		{"flt", "flt-1", time.Second, 1.0}, {"flt", "flt-2", 2 * time.Second, 2.0},
		{"state", "state-1", time.Second, 10}, {"state", "state-2", 5 * time.Second, 30},
		{"raw", "raw-1", time.Second, uint64(10)}, {"raw", "raw-2", 5 * time.Second, uint64(30)},
	} {
		rig.accept(t, sample.source, sample.id, sample.offset, sample.value)
	}

	badSample := rig.envelope("tri", "tri-bad", 4*time.Second, 1000, schema.QualityBad)
	badSample.QualityReason = "dial tcp 10.0.0.5:502: secret-token"
	require.NoError(t, rig.pipe.AcceptSample(ctx, badSample))

	// Close the window
	rig.closeWindow()
	row := rig.awaitRow(t)

	// Readback & numerical verification
	text := func(source string) string { return row.values[source].String }
	require.Equal(t, "1.5", plainDecimal(text("half")), "1 and 2 mean")
	require.Equal(t, "9007199254740992.5", plainDecimal(text("big")), "exact big mean above 2^53")
	require.Equal(t, "20", plainDecimal(text("tri")), "10, 20, 30 with resend and bad sample")
	require.Equal(t, "1.5", text("flt"))
	require.Equal(t, "30", text("state"), "state keeps newest reading")
	require.Equal(t, "30", text("raw"), "counter_raw keeps newest reading")
	require.False(t, row.values["quiet"].Valid, "quiet optional source is SQL NULL")

	// NULL / partial / unknown coverage checks
	halfEv := row.evidence["half"]
	require.Equal(t, [3]any{"unknown", int64(2), "3"}, [3]any{halfEv.Mean.Coverage, halfEv.Mean.ValidCount, halfEv.Mean.Sum})
	require.Equal(t, [2]any{1, "half-even"}, [2]any{halfEv.Mean.Scale, halfEv.Mean.Rounding})

	triEv := row.evidence["tri"]
	require.Equal(t, [4]any{"partial", int64(3), int64(1), "60"}, [4]any{triEv.Mean.Coverage, triEv.Mean.ValidCount, triEv.Mean.RejectedCount, triEv.Mean.Sum})
	require.Equal(t, []string{"bad"}, triEv.Mean.Reasons)

	quietEv := row.evidence["quiet"]
	require.Equal(t, [3]string{"missing", "no-samples", "missing"}, [3]string{quietEv.Status, quietEv.Reason, quietEv.Mean.Coverage})

	// Verify catalog table existence before cleanup
	target := rig.target(t)
	assertCatalogTableExists(t, target, kind, rig.f.namespace, rig.table)

	assertCommittedDeliveryReplay(t, rig.f.db, target, kind, rig.groupID, rig.snapshot.GroupRevision, rig.f.namespace, rig.table)

	// Verify counter baseline, delta, reset, and gap on wide layout
	assertCounterScenarios(t, kind, wideProfileLayoutVersion)

	// Stop own pipeline worker before cleanup to prevent reconstruction races
	assertTableCleanup(t, rig.pipe, target, kind, rig.f.namespace, rig.table)
}

// runMeanLongIntegrationMode tests the managed_long_v1 layout on SQLite, PostgreSQL, or TimescaleDB.
func runMeanLongIntegrationMode(t *testing.T, kind, storageMode string) {
	t.Helper()
	ctx := t.Context()
	f := newManagedDestinationFixture(t, kind)
	if storageMode == "timescale" && !timescaleInstalled(t, f) {
		t.Skip("TimescaleDB extension not installed on destination database; mode is NOT RUN")
	}

	seedMeanSources(t, f, meanSources)
	insertMeasurement(t, f, "meas-raw", "raw", "epoch-1", "counter")

	req := meanLongRequest(t, f, meanSources)
	if storageMode != "" && storageMode != "plain" {
		req["group"].(map[string]any)["storage_mode"] = storageMode
	}

	created := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, req)
	require.Equal(t, http.StatusOK, created.Code, created.Body.String())
	_, saved := groupSaveData(t, decodeJSONBody(t, created))
	groupID, table := saved["id"].(string), destinationOf(saved)["table_name"].(string)

	token := previewLong(t, f, groupID)
	if storageMode == "timescale" {
		require.Contains(t, token.Statements[1], "create_hypertable")
		require.Equal(t, "bucket_start", token.Tables[0].Partition.Column)
	}

	operation := confirmLong(t, f, groupID, token)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, operation.Status)

	ready, err := f.groups.Readiness(ctx, groupID)
	require.NoError(t, err)
	require.True(t, ready.Ready, "%+v", ready.Issues)

	applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/apply", scopeBody(t, f, groupID))
	require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())

	snapshot, err := f.groups.ResolveAppliedAt(ctx, groupID, time.Now().Add(time.Minute))
	require.NoError(t, err)

	var clock atomic.Int64
	clock.Store(snapshot.EffectiveAt.Add(time.Second).UnixNano())
	pipe := newMeanPipeline(t, f, &clock, "node/test", true)
	require.Equal(t, grouppipeline.StateActive, pipeStatus(pipe, groupID).State, "%+v", pipe.Status())

	frozen := map[string]workspace.WriteGroupMember{}
	for _, member := range snapshot.Group.Members {
		frozen[member.TagID] = member
	}

	effective := snapshot.EffectiveAt
	envelopeFor := func(source, id string, offset time.Duration, value any, quality schema.QualityFlag) measurement.SampleEnvelope {
		member := frozen["tag-"+source]
		observed := effective.Add(offset)
		return measurement.SampleEnvelope{
			SampleID: id, WorkspaceID: f.record.ID, DeviceID: "device-A", PointID: member.PointID, TagID: member.TagID,
			SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision,
			ObservedAt: observed, ReceivedAt: observed.Add(30 * time.Millisecond), TimeOrigin: "gateway", Value: value, Quality: quality,
		}
	}

	clock.Store(effective.Add(13 * time.Second).UnixNano())
	good := schema.QualityGood
	for _, sample := range []struct {
		source, id string
		offset     time.Duration
		value      any
	}{
		{"half", "half-1", time.Second, 1}, {"half", "half-2", 2 * time.Second, 2},
		{"big", "big-1", time.Second, uint64(9007199254740992)}, {"big", "big-2", 2 * time.Second, uint64(9007199254740993)},
		{"tri", "tri-1", time.Second, 10}, {"tri", "tri-2", 2 * time.Second, 20}, {"tri", "tri-3", 3 * time.Second, 30},
		{"tri", "tri-2", 2 * time.Second, 20}, // resend
		{"flt", "flt-1", time.Second, 1.0}, {"flt", "flt-2", 2 * time.Second, 2.0},
		{"state", "state-1", time.Second, 10}, {"state", "state-2", 5 * time.Second, 30},
		{"raw", "raw-1", time.Second, uint64(10)}, {"raw", "raw-2", 5 * time.Second, uint64(30)},
	} {
		require.NoError(t, pipe.AcceptSample(ctx, envelopeFor(sample.source, sample.id, sample.offset, sample.value, good)), sample.id)
	}

	badRead := envelopeFor("tri", "tri-bad", 4*time.Second, 1000, schema.QualityBad)
	badRead.QualityReason = "dial tcp 10.0.0.5:502: secret-token"
	require.NoError(t, pipe.AcceptSample(ctx, badRead))

	clock.Store(effective.Add(16 * time.Second).UnixNano())

	target := f.target
	if kind == "sqlite" {
		target = openSQLiteTarget(t, f.file)
	}

	fromTable := qualified(kind, f.namespace, table)
	require.Eventually(t, func() bool {
		var n int
		return target.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+fromTable).Scan(&n) == nil && n == len(meanSources)
	}, 15*time.Second, 50*time.Millisecond, "closed window writes one row per source")

	rows := readMeanRows(t, target, fromTable)
	require.Len(t, rows, len(meanSources))

	// Numerical readback verification
	halfRow := rows["tag-half"]
	require.Equal(t, [3]string{"ok", "decimal", "1.5"}, [3]string{halfRow.status, halfRow.valueType, plainDecimal(halfRow.decimal.String)})
	require.Equal(t, [3]any{"unknown", int64(2), "3"}, [3]any{halfRow.prov.Mean.Coverage, halfRow.prov.Mean.ValidCount, halfRow.prov.Mean.Sum})

	bigRow := rows["tag-big"]
	require.Equal(t, [3]string{"ok", "decimal", "9007199254740992.5"}, [3]string{bigRow.status, bigRow.valueType, plainDecimal(bigRow.decimal.String)})
	require.Equal(t, "18014398509481985", bigRow.prov.Mean.Sum)

	triRow := rows["tag-tri"]
	require.Equal(t, [3]string{"ok", "decimal", "20"}, [3]string{triRow.status, triRow.valueType, plainDecimal(triRow.decimal.String)})
	require.Equal(t, [3]any{"partial", int64(3), int64(1)}, [3]any{triRow.prov.Mean.Coverage, triRow.prov.Mean.ValidCount, triRow.prov.Mean.RejectedCount})
	require.Equal(t, []string{"bad"}, triRow.prov.Mean.Reasons)

	fltRow := rows["tag-flt"]
	require.Equal(t, [3]string{"ok", "float64", "1.5"}, [3]string{fltRow.status, fltRow.valueType, plainDecimal(fltRow.float.String)})

	stateRow := rows["tag-state"]
	require.Equal(t, [3]string{"ok", "int64", "30"}, [3]string{stateRow.status, stateRow.valueType, stateRow.integer.String})

	rawRow := rows["tag-raw"]
	require.Equal(t, [3]string{"ok", "uint64", "30"}, [3]string{rawRow.status, rawRow.valueType, rawRow.decimal.String})

	quietRow := rows["tag-quiet"]
	require.Equal(t, "missing", quietRow.status)
	require.False(t, quietRow.decimal.Valid || quietRow.integer.Valid || quietRow.float.Valid)

	// Timescale hypertable catalog proof after actual writes
	if storageMode == "timescale" {
		var dimensions, chunks int
		var timeColumn string
		require.NoError(t, target.QueryRowContext(ctx,
			`SELECT num_dimensions, num_chunks FROM timescaledb_information.hypertables WHERE hypertable_schema = $1 AND hypertable_name = $2`,
			f.namespace, table).Scan(&dimensions, &chunks))
		require.NoError(t, target.QueryRowContext(ctx,
			`SELECT column_name FROM timescaledb_information.dimensions WHERE hypertable_schema = $1 AND hypertable_name = $2`,
			f.namespace, table).Scan(&timeColumn))
		require.Equal(t, 1, dimensions)
		require.Equal(t, "bucket_start", timeColumn)
		require.GreaterOrEqual(t, chunks, 1, "real pipeline write created at least one chunk in hypertable")
		t.Logf("timescale_catalog case=%q schema=%q table=%q dimensions=%d time_column=%q chunks=%d result=PASS", t.Name(), f.namespace, table, dimensions, timeColumn, chunks)
	} else {
		assertCatalogTableExists(t, target, kind, f.namespace, table)
	}

	assertCommittedDeliveryReplay(t, f.db, target, kind, groupID, snapshot.GroupRevision, f.namespace, table)

	// Verify counter baseline, delta, reset, and gap on long layout (timescale tested with hypertable proof)
	if storageMode == "timescale" {
		assertCounterScenarios(t, "postgres", "managed_long_v1", "timescale")
	} else {
		assertCounterScenarios(t, kind, "managed_long_v1")
	}

	// Stop own pipeline worker before cleanup to prevent reconstruction races
	assertTableCleanup(t, pipe, target, kind, f.namespace, table)
}

func assertCatalogTableExists(t *testing.T, target *sql.DB, kind, namespace, table string) {
	t.Helper()
	var count int
	if kind == "sqlite" {
		require.NoError(t, target.QueryRowContext(t.Context(),
			`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name = ?`, table).Scan(&count))
	} else {
		require.NoError(t, target.QueryRowContext(t.Context(),
			`SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = $1 AND table_name = $2`,
			namespace, table).Scan(&count))
	}
	require.Equal(t, 1, count, "catalog proves table exists before cleanup")
	t.Logf("catalog case=%q schema=%q table=%q count=%d result=PASS", t.Name(), namespace, table, count)
}

// assertCommittedDeliveryReplay retries a committed effect through the real
// Sender. SQL acknowledgement loss is exercised separately by the existing
// TestAverageCommitAckLost and TestProfileCommitAckLost cases.
func assertCommittedDeliveryReplay(t *testing.T, db, target *sql.DB, kind, groupID, expectedRevision, expectedNamespace, expectedTable string) {
	t.Helper()
	ctx := t.Context()
	var effectKey string
	require.NoError(t, db.QueryRowContext(ctx,
		`SELECT effect_key FROM wg_delivery_outbox WHERE group_id = ? ORDER BY effect_key LIMIT 1`, groupID).Scan(&effectKey))
	store := groupdelivery.NewStore(db)
	var before groupdelivery.OutboxItem
	require.Eventually(t, func() bool {
		var err error
		before, err = store.GetOutbox(ctx, effectKey)
		return err == nil && before.State == groupdelivery.StateCommitted
	}, 15*time.Second, 20*time.Millisecond, "the real pipeline settles its SQL delivery")
	require.Equal(t, groupID, before.Key.GroupID)
	require.Equal(t, expectedRevision, before.Key.GroupRevision)
	require.Equal(t, expectedNamespace, before.Destination.TableSchema)
	require.Equal(t, expectedTable, before.Destination.TableName)
	require.NotEmpty(t, before.PayloadDigest)

	// Committed replay must return before target resolution. A nil resolver
	// would fail this test if the sender attempted another destination write.
	sender := groupdelivery.NewSender(store, nil, groupdelivery.SenderConfig{Owner: "node/matrix-replay"})
	result, err := sender.Deliver(ctx, effectKey)
	require.NoError(t, err)
	require.True(t, result.Already)
	require.Equal(t, groupdelivery.StateCommitted, result.State)
	after, err := store.GetOutbox(ctx, effectKey)
	require.NoError(t, err)
	require.Equal(t, before, after, "committed replay preserves payload, destination, revisions and effect")

	receiptTable := qualified(kind, expectedNamespace, "gw_effect_receipts")
	var receiptDigest string
	require.NoError(t, target.QueryRowContext(ctx,
		`SELECT payload_digest FROM `+receiptTable+` WHERE effect_key = $1`, effectKey).Scan(&receiptDigest))
	require.Equal(t, before.PayloadDigest, receiptDigest)
	var receiptCount, rowCount int
	require.NoError(t, target.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+receiptTable+` WHERE effect_key = $1`, effectKey).Scan(&receiptCount))
	require.NoError(t, target.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+qualified(kind, expectedNamespace, expectedTable)+` WHERE record_id = $1`, before.RecordID).Scan(&rowCount))
	require.Equal(t, 1, receiptCount)
	require.Equal(t, 1, rowCount, "the replay leaves one SQL row for this effect")
	evidence, err := json.Marshal(map[string]any{
		"case": t.Name(), "group_revision": before.Key.GroupRevision, "schema": expectedNamespace, "table": expectedTable,
		"connector_id": before.Destination.ConnectorID, "connector_revision": before.Destination.ConnectorRevision,
		"effect_key": effectKey, "payload_digest": before.PayloadDigest, "readback": "PASS", "committed_replay": "PASS",
		"sql_rows_for_record": rowCount, "receipts_for_effect": receiptCount,
	})
	require.NoError(t, err)
	t.Logf("delivery_evidence %s", evidence)
}

func assertCounterScenarios(t *testing.T, kind, layout string, storageModes ...string) {
	t.Helper()
	// Scenarios: 0=baseline+delta, 4=reset(decrease unproven), 5=gap(missing window span total)
	for _, scIdx := range []int{0, 4, 5} {
		sc := counterScenarios[scIdx]
		windows := runCumulativeWindows(t, kind, layout, sc, storageModes...)
		require.Len(t, windows, len(sc.want), "%s: window count mismatch", sc.name)
		for i, want := range sc.want {
			require.Equal(t, want.status, windows[i].status, "%s: window %d status", sc.name, i)
			require.Equal(t, want.raw, windows[i].raw, "%s: window %d exact raw", sc.name, i)
			require.Equal(t, want.delta, windows[i].delta, "%s: window %d delta", sc.name, i)
			require.Equal(t, want.delta != "", windows[i].hasDelta, "%s: window %d hasDelta", sc.name, i)
		}
		t.Logf("counter case=%q scenario=%q storage=%q windows=%d baseline_delta_reset_gap=PASS", t.Name(), sc.name, storageModes, len(windows))
	}
}

func assertTableCleanup(t *testing.T, pipe *grouppipeline.Pipeline, target *sql.DB, kind, namespace, table string) {
	t.Helper()
	// Stop own pipeline worker before dropping table to eliminate reconstruction race
	require.NoError(t, pipe.Stop(time.Second), "pipeline worker must be stopped before table cleanup")

	fromTable := qualified(kind, namespace, table)
	_, err := target.ExecContext(t.Context(), `DROP TABLE `+fromTable)
	require.NoError(t, err)

	var count int
	if kind == "sqlite" {
		require.NoError(t, target.QueryRowContext(t.Context(),
			`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name = ?`, table).Scan(&count))
	} else {
		require.NoError(t, target.QueryRowContext(t.Context(),
			`SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = $1 AND table_name = $2`,
			namespace, table).Scan(&count))
	}
	require.Zero(t, count, "managed table is dropped during cleanup")
	t.Logf("cleanup case=%q schema=%q table=%q catalog_count=%d stopped_worker=true result=PASS", t.Name(), namespace, table, count)
}
