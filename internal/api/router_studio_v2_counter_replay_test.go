package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

// The real current-definition provider remains available and observes the HTTP
// edit. Recovery must not call it: returning stale test data would hide the bug.
type counterReplayAPIReader struct {
	calls atomic.Int64
	db    *sql.DB
}

func (r *counterReplayAPIReader) CounterSources(ctx context.Context, g *workspace.WriteGroup) (map[string]runtime.CounterSource, error) {
	r.calls.Add(1)
	return grouppipeline.NewMeasurementCounterSources(measurement.NewSQLRepository(r.db)).CounterSources(ctx, g)
}

func counterReplayAPISnapshot(t *testing.T, store *groupdelivery.Store, key groupdelivery.GroupKey) (workspace.WriteGroupAppliedSnapshot, []byte) {
	t.Helper()
	raw, err := store.RuntimeVersion(t.Context(), key)
	require.NoError(t, err)
	require.NotEmpty(t, raw)
	var snap workspace.WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal(raw, &snap))
	require.Len(t, snap.CounterBindings, 1)
	return snap, raw
}

func counterReplayAPIJournal(t *testing.T, db *sql.DB) [][3]string {
	t.Helper()
	rows, err := db.Query(`SELECT sample_id,payload,payload_digest FROM wg_delivery_samples ORDER BY seq`)
	require.NoError(t, err)
	defer rows.Close()
	var result [][3]string
	for rows.Next() {
		var row [3]string
		require.NoError(t, rows.Scan(&row[0], &row[1], &row[2]))
		result = append(result, row)
	}
	require.NoError(t, rows.Err())
	return result
}

// Public configuration mutations, actual Apply, durable journal, fresh runtime
// and destination SQL are all required evidence; HTTP 200 alone proves none of
// the replay assertions. Sensor input uses the existing simulated envelope seam.
func TestStudioV2CounterReplayPublicMeasurementEditPreservesAcceptedEpoch(t *testing.T) {
	f := newManagedDestinationFixture(t, "sqlite", "uint64")
	var clock atomic.Int64
	clock.Store(time.Date(2026, 10, 1, 0, 0, 1, 0, time.UTC).UnixNano())
	now := func() time.Time { return time.Unix(0, clock.Load()).UTC() }
	f.groups.WithClock(now)
	definition := map[string]any{"id": "meas-energy", "workspace_id": f.record.ID, "device_id": "device-A", "point_id": "point-A", "tag_id": "tag-A", "equipment_id": "main", "name": "Energy", "quantity": "energy", "unit": "kWh", "semantic_kind": "counter", "numeric_encoding": "uint64", "definition_revision": "def-1", "source_binding_revision": "bind-1", "series_epoch": "epoch-1"}
	response := performJSONRequest(t, f.router, http.MethodPost, measurementsPath, definition)
	require.Equal(t, http.StatusCreated, response.Code, response.Body.String())
	saved := ensureProfileGroup(t, f, "cumulative", 10, "minute", "meas-energy")
	id, table := saved["id"].(string), destinationOf(saved)["table_name"].(string)
	require.Regexp(t, longTablePattern, table)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, id, previewLong(t, f, id)).Status)
	apply := func() string {
		ready := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+id+"/readiness", nil)
		require.Equal(t, http.StatusOK, ready.Code, ready.Body.String())
		require.Equal(t, true, decodeJSONBody(t, ready)["data"].(map[string]any)["ready"])
		response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/apply", scopeBody(t, f, id))
		require.Equal(t, http.StatusOK, response.Code, response.Body.String())
		return appliedRevisionOf(t, f, id)
	}
	oldRevision := apply()
	key := groupdelivery.GroupKey{WorkspaceID: f.record.ID, GroupID: id, GroupRevision: oldRevision}
	originalStore := groupdelivery.NewStore(f.db)
	old, original := counterReplayAPISnapshot(t, originalStore, key)
	for _, proof := range old.CounterBindings {
		require.Equal(t, "kWh", proof.Unit)
		require.Equal(t, "epoch-1", proof.SeriesEpoch)
		require.Equal(t, "def-1", proof.DefinitionRevision)
	}
	member := old.Group.Members[0]
	reader := &counterReplayAPIReader{db: f.db}
	targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.db))
	initial := grouppipeline.New(grouppipeline.Dependencies{Groups: f.groups, Destinations: targets, Store: originalStore, Counters: reader}, grouppipeline.Config{Now: now})
	clock.Store(old.EffectiveAt.Add(time.Second).UnixNano())
	require.NoError(t, initial.Reconcile(t.Context()))
	envelope := func(sampleID string, at time.Time, value uint64) measurement.SampleEnvelope {
		return measurement.SampleEnvelope{SampleID: sampleID, WorkspaceID: f.record.ID, DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID, SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision, ObservedAt: at, ReceivedAt: at.Add(time.Millisecond), Value: value, Quality: schema.QualityGood}
	}
	first := envelope("c-1", old.EffectiveAt.Add(time.Second), 100)
	second := envelope("c-2", old.EffectiveAt.Add(601*time.Second), 105)
	require.NoError(t, initial.AcceptSample(t.Context(), first))
	clock.Store(second.ObservedAt.Add(time.Second).UnixNano())
	require.NoError(t, initial.AcceptSample(t.Context(), second))
	accepted := counterReplayAPIJournal(t, f.db)
	require.Len(t, accepted, 2)
	stateBefore, err := originalStore.Restore(t.Context(), key)
	require.NoError(t, err)
	require.Len(t, stateBefore.Samples, 2)
	require.True(t, stateBefore.NextClose.IsZero())
	require.Empty(t, stateBefore.Baselines)

	// This is the original user-visible edit, not a direct SQL replacement.
	definition["unit"] = "MWh"
	response = performJSONRequest(t, f.router, http.MethodPut, measurementsPath+"/meas-energy", definition)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	require.Equal(t, true, decodeJSONBody(t, response)["data"].(map[string]any)["epoch_transitioned"])
	current, err := measurement.NewSQLRepository(f.db).GetByID(t.Context(), "meas-energy")
	require.NoError(t, err)
	require.Equal(t, "MWh", current.Unit)
	require.Equal(t, "epoch-2", current.SeriesEpoch)
	require.NotEqual(t, "def-1", current.DefinitionRevision)
	require.Equal(t, oldRevision, appliedRevisionOf(t, f, id))
	_, afterEdit := counterReplayAPISnapshot(t, originalStore, key)
	require.Equal(t, original, afterEdit)
	require.Equal(t, accepted, counterReplayAPIJournal(t, f.db))

	// Reopen the actual durable store and rebuild every boundary; no old runtime
	// object or assembler is reused. Background work only sends finalized rows;
	// all closure times are explicitly driven by the fixed clock below.
	var localPath string
	require.NoError(t, f.db.QueryRow(`SELECT file FROM pragma_database_list WHERE name='main'`).Scan(&localPath))
	local, err := sql.Open("sqlite", localPath+"?_pragma=busy_timeout(5000)")
	require.NoError(t, err)
	local.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, local.Close()) })
	store := groupdelivery.NewStore(local)
	groups := workspace.NewWriteGroupService(workspace.NewService(workspace.NewSQLRepository(local)), workspace.NewSQLWriteGroupRepository(local))
	destinations := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(local))
	start := func(owner string) *grouppipeline.Pipeline {
		p := grouppipeline.New(grouppipeline.Dependencies{Groups: groups, Destinations: destinations, Store: store, Counters: reader}, grouppipeline.Config{Now: now, Owner: owner, TickInterval: time.Hour, ReconcileInterval: time.Hour, WorkerInterval: 10 * time.Millisecond, Delivery: groupdelivery.DispatcherConfig{MaxPartitions: 1}})
		require.NoError(t, p.Start(t.Context()))
		t.Cleanup(func() { require.NoError(t, p.Stop(time.Second)) })
		return p
	}
	pipe := start("counter-replay/first-restart")
	require.NoError(t, pipe.AcceptSample(t.Context(), second), "identical accepted resend stays a duplicate after the definition edit")
	require.Equal(t, accepted, counterReplayAPIJournal(t, local))
	restored, err := store.Restore(t.Context(), key)
	require.NoError(t, err)
	require.Equal(t, stateBefore, restored)
	require.Zero(t, reader.calls.Load())
	clock.Store(old.EffectiveAt.Add(1200 * time.Second).UnixNano())
	pipe.TickAll(t.Context())
	target, err := sql.Open("sqlite", f.file)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, target.Close()) })
	waitRows := func(want int) {
		require.Eventually(t, func() bool {
			var rows, receipts int
			err := target.QueryRow(`SELECT COUNT(*) FROM ` + table).Scan(&rows)
			if err != nil {
				return false
			}
			err = local.QueryRow(`SELECT COUNT(*) FROM wg_delivery_receipts`).Scan(&receipts)
			return err == nil && rows == want && receipts == want
		}, 5*time.Second, 10*time.Millisecond)
	}
	waitRows(2)
	oldEffects := counterReplayAPIEffects(t, local, oldRevision)
	require.Len(t, oldEffects, 2)
	var unit, epoch, raw, status, provenance string
	var delta *string
	require.NoError(t, target.QueryRow(`SELECT unit,series_epoch,CAST(value_decimal AS TEXT),CAST(delta_value AS TEXT),status,provenance FROM `+table+` WHERE group_revision=? AND bucket_start=?`, oldRevision, old.EffectiveAt.Add(600*time.Second).Format(meanTimeLayout)).Scan(&unit, &epoch, &raw, &delta, &status, &provenance))
	require.Equal(t, "kWh", unit)
	require.Equal(t, "epoch-1", epoch)
	require.Equal(t, "105", raw)
	require.NotNil(t, delta)
	require.Equal(t, "5", *delta)
	require.Contains(t, provenance, `"definition_revision":"def-1"`)
	require.NoError(t, pipe.Stop(time.Second))
	oldClosed, err := store.Restore(t.Context(), key)
	require.NoError(t, err)
	require.Empty(t, oldClosed.Samples)
	require.Len(t, oldClosed.Baselines, 1)

	// An explicit fresh saved group revision, reviewed and applied via HTTP, is
	// the only path to the new epoch. It must not reuse the old applied owner.
	response = performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+id, nil)
	require.Equal(t, http.StatusOK, response.Code)
	data, draft := groupSaveData(t, decodeJSONBody(t, response))
	draft["name"] = "Confirm new counter epoch"
	response = performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+id, groupEditRequest(f.writeGroupRouterFixture, data, draft))
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	_, fresh := groupSaveData(t, decodeJSONBody(t, response))
	nextRevision := apply()
	require.Equal(t, fresh["revision"], nextRevision)
	require.NotEqual(t, oldRevision, nextRevision)
	nextKey := key
	nextKey.GroupRevision = nextRevision
	next, nextBytes := counterReplayAPISnapshot(t, store, nextKey)
	for _, proof := range next.CounterBindings {
		require.Equal(t, "MWh", proof.Unit)
		require.Equal(t, "epoch-2", proof.SeriesEpoch)
		require.Equal(t, current.DefinitionRevision, proof.DefinitionRevision)
	}
	clock.Store(next.EffectiveAt.Add(2 * time.Second).UnixNano())
	pipe = start("counter-replay/new-owner")
	nextSample := envelope("c-2", next.EffectiveAt.Add(time.Second), 500)
	nextSample.MeasurementID, nextSample.SeriesEpoch, nextSample.DefinitionRevision, nextSample.SourceBindingRevision = "meas-energy", "epoch-2", current.DefinitionRevision, current.SourceBindingRevision
	require.NoError(t, pipe.AcceptSample(t.Context(), nextSample))
	require.NoError(t, pipe.AcceptSample(t.Context(), nextSample))
	newState, err := store.Restore(t.Context(), nextKey)
	require.NoError(t, err)
	require.Len(t, newState.Samples, 1)
	require.Empty(t, newState.Baselines)
	retained, err := store.Restore(t.Context(), key)
	require.NoError(t, err)
	require.Equal(t, oldClosed, retained)
	clock.Store(next.EffectiveAt.Add(600 * time.Second).UnixNano())
	pipe.TickAll(t.Context())
	waitRows(3)
	require.NoError(t, target.QueryRow(`SELECT unit,series_epoch,CAST(value_decimal AS TEXT),CAST(delta_value AS TEXT),status FROM `+table+` WHERE group_revision=?`, nextRevision).Scan(&unit, &epoch, &raw, &delta, &status))
	require.Equal(t, "MWh", unit)
	require.Equal(t, "epoch-2", epoch)
	require.Equal(t, "500", raw)
	require.Nil(t, delta)
	require.Equal(t, "baseline", status)
	require.NoError(t, pipe.Stop(time.Second))
	// Consumed journal rows may be reclaimed by the normal worker. The accepted
	// bytes were checked before consumption; finalized effects must never change.
	require.Equal(t, oldEffects, counterReplayAPIEffects(t, local, oldRevision))
	_, retainedBytes := counterReplayAPISnapshot(t, store, key)
	require.Equal(t, original, retainedBytes)
	_, retainedNext := counterReplayAPISnapshot(t, store, nextKey)
	require.Equal(t, nextBytes, retainedNext)
	require.Zero(t, reader.calls.Load())
}

func counterReplayAPIEffects(t *testing.T, db *sql.DB, revision string) [][3]string {
	t.Helper()
	rows, err := db.Query(`SELECT effect_key,payload,payload_digest FROM wg_delivery_outbox WHERE group_revision=? ORDER BY bucket_start`, revision)
	require.NoError(t, err)
	defer rows.Close()
	var result [][3]string
	for rows.Next() {
		var row [3]string
		require.NoError(t, rows.Scan(&row[0], &row[1], &row[2]))
		result = append(result, row)
	}
	require.NoError(t, rows.Err())
	return result
}
