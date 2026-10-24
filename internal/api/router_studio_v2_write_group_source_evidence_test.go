package api

import (
	"database/sql"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// evidenceRun is an applied long group running in the production pipeline with a controllable clock.
type evidenceRun struct {
	f        managedGroupFixture
	kind     string
	groupID  string
	table    string
	snapshot *workspace.WriteGroupAppliedSnapshot
	pipe     *grouppipeline.Pipeline
	clock    *atomic.Int64
	target   *sql.DB
	delivery http.Handler
}

func (r evidenceRun) from() string { return qualified(r.kind, r.f.namespace, r.table) }

func (r evidenceRun) setClock(at time.Time) { r.clock.Store(at.UnixNano()) }

func (r evidenceRun) accept(t *testing.T, id string, observed time.Time, value any) {
	t.Helper()
	member := r.snapshot.Group.Members[0]
	r.setClock(observed.Add(time.Second))
	require.NoError(t, r.pipe.AcceptSample(t.Context(), measurement.SampleEnvelope{
		SampleID: id, WorkspaceID: r.f.record.ID, DeviceID: "device-A", PointID: "point-A", TagID: "tag-A",
		SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision,
		ObservedAt: observed, ReceivedAt: observed.Add(30 * time.Millisecond), TimeOrigin: "gateway", Value: value, Quality: schema.QualityGood,
	}))
}

func (r evidenceRun) waitRows(t *testing.T, want int) {
	t.Helper()
	require.Eventually(t, func() bool {
		var n int
		return r.target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM `+r.from()).Scan(&n) == nil && n == want
	}, 15*time.Second, 50*time.Millisecond, "the closed windows reach the shared long table")
}

// committedEvidence reads the API answer once the newest record is SQL committed with a matching receipt.
func (r evidenceRun) committedEvidence(t *testing.T) map[string]any {
	t.Helper()
	var item map[string]any
	require.Eventually(t, func() bool {
		response := performJSONRequest(t, r.delivery, http.MethodGet, writeGroupsPath+"/"+r.groupID+"/delivery", nil)
		if response.Code != http.StatusOK {
			return false
		}
		data := decodeJSONBody(t, response)["data"].(map[string]any)
		evidence, ok := data["source_evidence"].(map[string]any)
		if !ok || len(evidence["sources"].([]any)) != 1 {
			return false
		}
		item = evidence["sources"].([]any)[0].(map[string]any)
		latest, _ := item["latest"].(map[string]any)
		return latest != nil && latest["delivery"] == "sql_committed"
	}, 15*time.Second, 50*time.Millisecond, "the API reports the newest record as SQL committed with a matching receipt")
	return item
}

func startEvidenceRun(t *testing.T, kind string, f managedGroupFixture, groupID string, counters bool) evidenceRun {
	t.Helper()
	table := destinationOf(groupFromRoute(t, f, groupID))["table_name"].(string)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, groupID, previewLong(t, f, groupID)).Status)
	ready, err := f.groups.Readiness(t.Context(), groupID)
	require.NoError(t, err)
	require.True(t, ready.Ready, "%+v", ready.Issues)
	applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/apply", scopeBody(t, f, groupID))
	require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())
	snapshot, err := f.groups.ResolveAppliedAt(t.Context(), groupID, time.Now().Add(time.Hour))
	require.NoError(t, err)

	var clock atomic.Int64
	clock.Store(snapshot.EffectiveAt.Add(time.Second).UnixNano())
	targetSvc := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.db))
	deps := grouppipeline.Dependencies{
		Groups: f.groups, Tags: tag.NewService(tag.NewSQLRepository(f.db)), Destinations: targetSvc,
		Inspector: dbtarget.NewReadOnlyTableInspector(targetSvc), Store: groupdelivery.NewStore(f.db),
	}
	if counters {
		deps.Counters = grouppipeline.NewMeasurementCounterSources(measurement.NewSQLRepository(f.db))
	}
	pipe := grouppipeline.New(deps, grouppipeline.Config{
		NodeID: "node", Owner: "node/test", Now: func() time.Time { return time.Unix(0, clock.Load()).UTC() },
		TickInterval: 20 * time.Millisecond, WorkerInterval: 20 * time.Millisecond, ReconcileInterval: time.Hour, MaxFutureSkew: time.Minute,
	})
	require.NoError(t, pipe.Start(t.Context()))
	t.Cleanup(func() { _ = pipe.Stop(time.Second) })
	require.Equal(t, grouppipeline.StateActive, pipeStatus(pipe, groupID).State, "%+v", pipe.Status())

	target := f.target
	if kind == "sqlite" {
		target, err = sql.Open("sqlite", f.file)
		require.NoError(t, err)
		t.Cleanup(func() { _ = target.Close() })
	}
	return evidenceRun{
		f: f, kind: kind, groupID: groupID, table: table, snapshot: snapshot, pipe: pipe, clock: &clock, target: target,
		delivery: NewRouter(&DatalinkServices{Workspace: f.workspace, WriteGroups: f.groups, WriteGroupDelivery: pipe}),
	}
}

func groupFromRoute(t *testing.T, f managedGroupFixture, groupID string) map[string]any {
	t.Helper()
	response := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+groupID, nil)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	_, group := groupSaveData(t, decodeJSONBody(t, response))
	return group
}

func parseEvidenceTime(t *testing.T, value any) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339Nano, value.(string))
	require.NoError(t, err)
	return parsed.UTC()
}

// TestCircuitProfileDeliveryEvidence applies long groups through the real route, runs the production pipeline and
// sender over a real destination and compares the per source evidence the delivery API reports with what the
// destination actually holds: the stored status, raw presence and NULL difference, the window, and the receipt.
func TestCircuitProfileDeliveryEvidence(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind+"/instantaneous with a missing window", func(t *testing.T) {
			f := newManagedDestinationFixture(t, kind)
			group := ensureLong(t, f, "instantaneous", "second", 15, "")
			run := startEvidenceRun(t, kind, f, group["id"].(string), false)

			effective := run.snapshot.EffectiveAt
			run.accept(t, "sample-1", effective.Add(2*time.Second), 21.5)
			run.setClock(effective.Add(16 * time.Second))
			run.waitRows(t, 1)
			first := run.committedEvidence(t)
			latest := first["latest"].(map[string]any)
			require.Equal(t, true, latest["available"])
			require.Equal(t, "ok", latest["status"])
			require.Equal(t, true, latest["raw_present"])
			require.Nil(t, latest["delta_known"], "an instantaneous record has no difference")
			require.Equal(t, run.snapshot.AppliedRevision, latest["group_revision"])

			run.setClock(effective.Add(31 * time.Second))
			run.waitRows(t, 2)
			second := run.committedEvidence(t)
			latest = second["latest"].(map[string]any)
			require.Equal(t, "missing", latest["status"])
			require.Equal(t, "no-samples", latest["reason"])
			require.Equal(t, false, latest["raw_present"], "a window without a sample is stored as NULL")
			require.NotNil(t, latest["committed_at"], "the receipt matched this record")
			require.Equal(t, latest["effect_key"], second["last_sql_committed"].(map[string]any)["effect_key"])

			var storedStatus string
			var storedValue sql.NullFloat64
			var storedStart any
			require.NoError(t, run.target.QueryRowContext(t.Context(),
				`SELECT status, value_float, bucket_start FROM `+run.from()+` ORDER BY bucket_start DESC LIMIT 1`).Scan(&storedStatus, &storedValue, &storedStart))
			require.Equal(t, "missing", storedStatus)
			require.False(t, storedValue.Valid, "the destination holds NULL, matching raw_present=false")
			if at, ok := storedStart.(time.Time); ok {
				require.Equal(t, at.UTC(), parseEvidenceTime(t, latest["window_start"]))
			}
			body := performJSONRequest(t, run.delivery, http.MethodGet, writeGroupsPath+"/"+run.groupID+"/delivery", nil).Body.String()
			require.NotContains(t, body, "21.5", "the evidence carries no reading")
			require.NotContains(t, body, "password")
		})

		t.Run(kind+"/cumulative baseline then difference", func(t *testing.T) {
			f := newManagedDestinationFixture(t, kind, "uint64")
			_, err := f.db.ExecContext(t.Context(), `
				INSERT INTO measurement_definitions (
					id, workspace_id, device_id, point_id, tag_id, equipment_id, definition_revision, source_binding_revision,
					series_epoch, name, quantity, unit, semantic_kind
				) VALUES ('meas-energy', ?, 'device-A', 'point-A', 'tag-A', 'device-A', 'def-5', 'bind-5', 'epoch-3', 'Energy', 'energy', 'kWh', 'counter')`,
				f.record.ID)
			require.NoError(t, err)
			request := profileRequest(t, f.writeGroupRouterFixture, "cumulative", 10, "minute", "main", "")
			saved := request["group"].(map[string]any)
			saved["layout_version"] = "managed_long_v1"
			saved["destination"].(map[string]any)["table_schema"] = f.namespace
			saved["members"].([]map[string]any)[0]["measurement_id"] = "meas-energy"
			created := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, request)
			require.Equal(t, http.StatusOK, created.Code, created.Body.String())
			_, group := groupSaveData(t, decodeJSONBody(t, created))
			run := startEvidenceRun(t, kind, f, group["id"].(string), true)

			effective := run.snapshot.EffectiveAt
			run.accept(t, "kwh-1", effective.Add(9*time.Minute), uint64(9007199254740992))
			run.setClock(effective.Add(10*time.Minute + time.Second))
			run.waitRows(t, 1)
			baseline := run.committedEvidence(t)["latest"].(map[string]any)
			require.Equal(t, "baseline", baseline["status"])
			require.Equal(t, "first-reading", baseline["reason"])
			require.Equal(t, true, baseline["raw_present"], "the first raw reading reached SQL")
			require.Equal(t, false, baseline["delta_known"], "no comparable baseline: the difference is unknown, not zero")
			var delta sql.NullString
			require.NoError(t, run.target.QueryRowContext(t.Context(), `SELECT CAST(delta_value AS TEXT) FROM `+run.from()).Scan(&delta))
			require.False(t, delta.Valid, "the destination holds NULL, matching delta_known=false")

			run.accept(t, "kwh-2", effective.Add(19*time.Minute), uint64(9007199254740993))
			run.setClock(effective.Add(20*time.Minute + time.Second))
			run.waitRows(t, 2)
			next := run.committedEvidence(t)
			ok := next["latest"].(map[string]any)
			require.Equal(t, "ok", ok["status"])
			require.Equal(t, true, ok["delta_known"])
			require.NoError(t, run.target.QueryRowContext(t.Context(), `SELECT CAST(delta_value AS TEXT) FROM `+run.from()+` ORDER BY bucket_start DESC LIMIT 1`).Scan(&delta))
			require.Equal(t, "1", delta.String, "the stored exact difference matches delta_known=true")
			require.NotEqual(t, baseline["effect_key"], ok["effect_key"], "the newest record replaced the baseline record")
			require.NotContains(t, performJSONRequest(t, run.delivery, http.MethodGet, writeGroupsPath+"/"+run.groupID+"/delivery", nil).Body.String(), "9007199254740993")
		})
	}
}
