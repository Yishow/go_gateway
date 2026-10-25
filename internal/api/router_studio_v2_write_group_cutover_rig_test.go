package api

import (
	"cmp"
	"database/sql"
	"encoding/json"
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
	"go-gateway/internal/testutil/lostcommit"

	"github.com/stretchr/testify/require"
)

const (
	cutoverWideSeconds = 15
	cutoverLongSeconds = 20
	cutoverWorkerEvery = 20 * time.Millisecond
	cutoverTimeText    = "2006-01-02T15:04:05.000000Z"
)

// cutoverRig is a wide managed group that has been recording, and a long
// profile group for the same source that is saved and confirmed but not applied,
// both on one disposable destination, run by the production pipeline and routes.
// Time is the test's own: windows close when it says they do. Every boundary
// the tests choose lies inside both a wide (15 s) and a long (20 s) window.
type cutoverRig struct {
	t      *testing.T
	f      managedGroupFixture
	kind   string
	dest   *sql.DB
	wide   map[string]any
	long   map[string]any
	pipe   *grouppipeline.Pipeline
	clock  atomic.Int64
	member workspace.WriteGroupMember
	// anchor is a UTC minute (a multiple of both window lengths) far enough ahead
	// that every instant the tests use is still in the future of the real clock.
	anchor time.Time
	// wideEffective is when the wide group started recording.
	wideEffective time.Time
	cumulative    bool
	custom        string
	faults        *lostcommit.Faults
	// targets are every long group the cutover moves the source to; long is the first.
	targets []map[string]any
	second  workspace.WriteGroupMember
}

func (r *cutoverRig) quoted(table string) string {
	if r.kind == "sqlite" {
		return `"` + table + `"`
	}
	return `"` + r.f.namespace + `"."` + table + `"`
}

func (r *cutoverRig) at(offset time.Duration) time.Time { return r.anchor.Add(offset) }

func newCutoverRig(t *testing.T, kind string) *cutoverRig {
	t.Helper()
	return newCutoverRigWith(t, kind, cutoverOptions{})
}

// cutoverOptions shape the rig. cumulative makes the source an exact uint64
// energy counter and the long target a one-minute cumulative profile counting
// it from a saved measurement definition. dual adds the instantaneous profile of
// the same circuit as a second target of the same source. second gives the wide
// group a second source that no target takes.
type cutoverOptions struct {
	cumulative, dual, second bool
	// custom makes the wide group a custom one on a table the test created, with
	// the first source's column declared as given ("not_null", "nullable" or
	// "default"); empty keeps the managed wide group.
	custom string
	// lostCommit opens the local database so a commit can be made to succeed for
	// real and still report failure (rig.faults.Lose).
	lostCommit bool
	// wideSeconds is the saved interval of the wide group; zero keeps the 15 s default.
	wideSeconds int
	// wideTarget makes every target a managed_wide_profile_v1 group instead of a long
	// one: an explicit 15 s mean for the instantaneous profile, and the one-minute
	// counter profile for the cumulative one.
	wideTarget bool
	// unconfirmedTarget leaves the targets saved without the explicit schema confirmation.
	unconfirmedTarget bool
}

func newCutoverRigWith(t *testing.T, kind string, o cutoverOptions) *cutoverRig {
	t.Helper()
	cumulative := o.cumulative || o.dual
	destination := kind
	if kind == "timescale" {
		destination = "postgres"
	}
	var tagTypes []string
	if cumulative {
		tagTypes = []string{"uint64"}
	}
	var faults *lostcommit.Faults
	var f managedGroupFixture
	if o.lostCommit {
		faults = &lostcommit.Faults{}
		base := newWriteGroupRouterFixtureOn(t, func(path string) *sql.DB { return lostcommit.OpenFaults(path, faults) })
		f = newManagedDestinationFixtureOn(t, destination, base, tagTypes...)
	} else {
		f = newManagedDestinationFixture(t, destination, tagTypes...)
	}
	if kind == "timescale" && !timescaleInstalled(t, f) {
		t.Skip("this destination database has no TimescaleDB extension; the timescale case is NOT RUN here")
	}
	r := &cutoverRig{t: t, f: f, kind: kind, cumulative: cumulative, faults: faults}
	if o.second {
		seedProfileSource(t, f.writeGroupRouterFixture, "second")
	}
	if cumulative {
		_, err := f.db.ExecContext(t.Context(), `
			INSERT INTO measurement_definitions (
				id, workspace_id, device_id, point_id, tag_id, equipment_id, definition_revision, source_binding_revision,
				series_epoch, name, quantity, unit, semantic_kind
			) VALUES ('meas-energy', ?, 'device-A', 'point-A', 'tag-A', 'device-A', 'def-5', 'bind-5', 'epoch-3', 'Energy', 'energy', 'kWh', 'counter')`,
			f.record.ID)
		require.NoError(t, err)
	}
	r.dest = (&recoveryRig{t: t, f: f, kind: kind}).openDestination()
	r.anchor = time.Now().UTC().Add(2 * time.Minute).Truncate(time.Minute)

	request := f.createRequest()
	group := request["group"].(map[string]any)
	group["destination"].(map[string]any)["storage_strategy"] = "managed"
	group["destination"].(map[string]any)["table_schema"] = f.namespace
	group["row_policy"] = map[string]any{"interval_seconds": cmp.Or(o.wideSeconds, cutoverWideSeconds)}
	r.custom = o.custom
	if o.custom != "" {
		r.createCustomTable(o.custom)
		group["destination"] = map[string]any{
			"connector_id": "connector-A", "connector_revision": "connector-1", "table_schema": f.namespace,
			"table_name": "custom_readings", "storage_strategy": "custom",
		}
		group["write_policy"] = map[string]any{"dedupe_capability": "none"}
		group["members"] = []map[string]any{{"device_id": "device-A", "point_id": "point-A", "tag_id": "tag-A", "target_column": "a_col", "required": true}}
	}
	if o.second {
		column := "second"
		if o.custom != "" {
			column = "b_col"
		}
		group["members"] = append(group["members"].([]map[string]any), map[string]any{
			"device_id": "device-A", "point_id": "point-second", "tag_id": "tag-second", "target_column": column, "required": true,
		})
	}
	created := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, request)
	require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
	_, r.wide = groupSaveData(t, decodeJSONBody(t, created))
	wideID := r.wide["id"].(string)
	if o.custom == "" {
		require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, wideID, previewLong(t, f, wideID)).Status)
	}
	applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+wideID+"/apply", scopeBody(t, f, wideID))
	require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())

	profiles := []string{"instantaneous"}
	if cumulative {
		profiles = []string{"cumulative"}
	}
	if o.dual {
		profiles = append(profiles, "instantaneous")
	}
	for _, profile := range profiles {
		long := profileRequest(t, f.writeGroupRouterFixture, "instantaneous", cutoverLongSeconds, "second", "main", "")
		if profile == "cumulative" {
			long = profileRequest(t, f.writeGroupRouterFixture, "cumulative", 1, "minute", "main", "")
			long["group"].(map[string]any)["members"].([]map[string]any)[0]["measurement_id"] = "meas-energy"
		}
		longGroup := long["group"].(map[string]any)
		if o.wideTarget {
			longGroup["layout_version"] = "managed_wide_profile_v1"
			member := longGroup["members"].([]map[string]any)[0]
			member["required"] = false
			if profile == "instantaneous" {
				longGroup["recording_profile"] = map[string]any{"kind": "instantaneous", "interval": cutoverWideSeconds, "unit": "second", "sample_reduction": "mean"}
				member["value_role"] = "gauge"
			}
		} else {
			longGroup["layout_version"] = "managed_long_v1"
			longGroup["storage_mode"] = "plain"
			if kind == "timescale" {
				longGroup["storage_mode"] = "timescale"
			}
		}
		longGroup["destination"].(map[string]any)["table_schema"] = f.namespace
		response := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, long)
		require.Equal(t, http.StatusOK, response.Code, response.Body.String())
		_, saved := groupSaveData(t, decodeJSONBody(t, response))
		longID := saved["id"].(string)
		if !o.unconfirmedTarget {
			require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, longID, previewLong(t, f, longID)).Status)
		}
		r.targets = append(r.targets, saved)
	}
	r.long = r.targets[0]

	snapshot, err := f.groups.ResolveAppliedAt(t.Context(), wideID, time.Now().Add(time.Hour))
	require.NoError(t, err)
	r.wideEffective = snapshot.EffectiveAt
	r.member = snapshot.Group.Members[0]
	for _, m := range snapshot.Group.Members {
		if m.TagID == "tag-second" {
			r.second = m
		}
	}
	r.clock.Store(r.wideEffective.Add(time.Second).UnixNano())
	r.launch(cutoverWorkerEvery)
	return r
}

func (r *cutoverRig) launch(workerEvery time.Duration) {
	r.t.Helper()
	store := groupdelivery.NewStore(r.f.db)
	targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(r.f.db))
	r.pipe = grouppipeline.New(grouppipeline.Dependencies{
		Groups: r.f.groups, Tags: tag.NewService(tag.NewSQLRepository(r.f.db)), Destinations: targets,
		Inspector: dbtarget.NewReadOnlyTableInspector(targets), Store: store,
		Counters: grouppipeline.NewMeasurementCounterSources(measurement.NewSQLRepository(r.f.db)),
	}, grouppipeline.Config{
		NodeID: "node", Owner: "node/test", Now: func() time.Time { return time.Unix(0, r.clock.Load()).UTC() },
		TickInterval: 20 * time.Millisecond, WorkerInterval: workerEvery, ReconcileInterval: time.Hour, MaxFutureSkew: time.Minute,
	})
	r.f.router = NewRouter(&DatalinkServices{
		Workspace: r.f.workspace, WriteGroups: r.f.groups, DBTarget: targets, WriteGroupDelivery: r.pipe,
		RecordingPlan: recordingplan.NewService(recordingplan.NewSQLRepository(r.f.db)),
	})
	require.NoError(r.t, r.pipe.Start(r.t.Context()))
	pipe := r.pipe
	r.t.Cleanup(func() { _ = pipe.Stop(time.Second) })
}

func (r *cutoverRig) restart(workerEvery time.Duration) {
	r.t.Helper()
	require.NoError(r.t, r.pipe.Stop(time.Second))
	r.launch(workerEvery)
}

// accept offers one reading at an exact instant. The clock follows the reading,
// so a window closes only after everything observed inside it was accepted.
func (r *cutoverRig) accept(observed time.Time, value any) { r.acceptFor(r.member, observed, value) }

// acceptFor offers one reading of the given source at an exact instant.
func (r *cutoverRig) acceptFor(member workspace.WriteGroupMember, observed time.Time, value any) {
	r.t.Helper()
	r.clock.Store(observed.Add(10 * time.Millisecond).UnixNano())
	require.NoError(r.t, r.pipe.AcceptSample(r.t.Context(), measurement.SampleEnvelope{
		SampleID: "cutover-" + member.TagID + observed.Format(time.RFC3339Nano), WorkspaceID: r.f.record.ID, DeviceID: member.DeviceID,
		PointID: member.PointID, TagID: member.TagID, SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision,
		ObservedAt: observed, ReceivedAt: observed.Add(30 * time.Millisecond), TimeOrigin: "gateway", Value: value, Quality: schema.QualityGood,
	}))
}

func (r *cutoverRig) advance(to time.Time) { r.clock.Store(to.UnixNano()) }

func (r *cutoverRig) reconcile() {
	r.t.Helper()
	require.NoError(r.t, r.pipe.Reconcile(r.t.Context()))
}

func (r *cutoverRig) wideTable() string { return destinationOf(r.wide)["table_name"].(string) }
func (r *cutoverRig) longTable() string { return destinationOf(r.long)["table_name"].(string) }

func (r *cutoverRig) scalar(query string) int {
	r.t.Helper()
	var n int
	require.NoError(r.t, r.dest.QueryRowContext(r.t.Context(), query).Scan(&n))
	return n
}

func (r *cutoverRig) wideRows() int {
	return r.scalar(`SELECT COUNT(*) FROM ` + r.quoted(r.wideTable()))
}

// longRecords reads the long table in window order as typed records.
type longRecord struct {
	Profile              string
	Raw, Delta           string
	Status               string
	Start, End, Observed time.Time
	Value                float64
	Provenance           string
}

// windowOf returns the truncation evidence of a record's provenance, or nil.
func windowOf(t *testing.T, provenance string) map[string]any {
	t.Helper()
	var decoded struct {
		Window map[string]any `json:"window"`
	}
	require.NoError(t, json.Unmarshal([]byte(provenance), &decoded))
	return decoded.Window
}

// storedInstant parses a timestamp as either destination renders it as text.
func storedInstant(t *testing.T, text string) time.Time {
	t.Helper()
	if text == "" {
		return time.Time{}
	}
	for _, layout := range []string{cutoverTimeText, "2006-01-02 15:04:05.999999999-07", "2006-01-02 15:04:05.999999999-07:00", time.RFC3339Nano} {
		if parsed, err := time.Parse(layout, text); err == nil {
			return parsed.UTC()
		}
	}
	require.FailNow(t, "unreadable stored instant", text)
	return time.Time{}
}

func (r *cutoverRig) longRecords() []longRecord {
	r.t.Helper()
	cast := func(column string) string { return `CAST(` + column + ` AS TEXT)` }
	rows, err := r.dest.QueryContext(r.t.Context(), `SELECT profile_kind, COALESCE(`+cast("value_decimal")+`, `+cast("value_float")+`, ''), COALESCE(`+cast("delta_value")+`, 'NULL'), status, `+cast("bucket_start")+`, `+cast("bucket_end")+`,
		COALESCE(`+cast("observed_at")+`, ''), COALESCE(value_float, 0), `+cast("provenance")+` FROM `+r.quoted(r.longTable())+` ORDER BY bucket_start, profile_kind`)
	require.NoError(r.t, err)
	defer rows.Close()
	var out []longRecord
	for rows.Next() {
		var rec longRecord
		var start, end, observed string
		require.NoError(r.t, rows.Scan(&rec.Profile, &rec.Raw, &rec.Delta, &rec.Status, &start, &end, &observed, &rec.Value, &rec.Provenance))
		rec.Start, rec.End, rec.Observed = storedInstant(r.t, start), storedInstant(r.t, end), storedInstant(r.t, observed)
		out = append(out, rec)
	}
	require.NoError(r.t, rows.Err())
	return out
}

func (r *cutoverRig) outbox() map[string]int {
	r.t.Helper()
	rows, err := r.f.db.QueryContext(r.t.Context(), `SELECT state, COUNT(*) FROM wg_delivery_outbox GROUP BY state`)
	require.NoError(r.t, err)
	defer rows.Close()
	out := map[string]int{}
	for rows.Next() {
		var state string
		var n int
		require.NoError(r.t, rows.Scan(&state, &n))
		out[state] = n
	}
	require.NoError(r.t, rows.Err())
	return out
}

func (r *cutoverRig) waitFor(cond func() bool, message string) {
	r.t.Helper()
	(&recoveryRig{t: r.t, f: r.f}).waitFor(cond, message)
}

// cutoverBody is the request of a review or an apply.
func (r *cutoverRig) cutoverBody(boundary time.Time, extra map[string]any) map[string]any {
	r.t.Helper()
	wideID := r.wide["id"].(string)
	current := getGroup(r.t, r.f, wideID)
	data, _ := groupSaveData(r.t, decodeJSONBody(r.t, performJSONRequest(r.t, r.f.router, http.MethodGet, writeGroupsPath+"/"+wideID, nil)))
	targets := make([]map[string]any, 0, len(r.targets))
	for _, saved := range r.targets {
		longGroup := getGroup(r.t, r.f, saved["id"].(string))
		targets = append(targets, map[string]any{"group_id": longGroup["id"], "expected_revision": longGroup["revision"]})
	}
	body := map[string]any{
		"workspace_id": r.f.record.ID, "expected_workspace_revision": data["workspace_revision"],
		"expected_source_revision": current["revision"], "expected_connector_revision": "connector-1",
		"targets":  targets,
		"boundary": boundary.UTC().Format(time.RFC3339Nano),
	}
	for k, v := range extra {
		body[k] = v
	}
	return body
}

func (r *cutoverRig) post(path string, body map[string]any) (status int, data map[string]any) {
	r.t.Helper()
	response := performJSONRequest(r.t, r.f.router, http.MethodPost, path, body)
	return response.Code, decodeJSONBody(r.t, response)
}

// review asks for the difference at the boundary and returns the data.
func (r *cutoverRig) review(boundary time.Time) map[string]any {
	r.t.Helper()
	code, body := r.post(writeGroupsPath+"/"+r.wide["id"].(string)+"/cutover-review", r.cutoverBody(boundary, nil))
	require.Equal(r.t, http.StatusOK, code, "%v", body)
	return body["data"].(map[string]any)
}

// cutover reviews and applies the forward cutover at boundary.
func (r *cutoverRig) cutover(boundary time.Time) map[string]any {
	r.t.Helper()
	review := r.review(boundary)
	require.Equal(r.t, true, review["applicable"], "%v", review["issues"])
	code, body := r.post(writeGroupsPath+"/"+r.wide["id"].(string)+"/cutover",
		r.cutoverBody(boundary, map[string]any{"review_digest": review["digest"]}))
	require.Equal(r.t, http.StatusOK, code, "%v", body)
	return body["data"].(map[string]any)
}

// createCustomTable creates the table a custom wide group records into. The
// first source's column a_col is declared as the test asks; the second source's
// b_col is always writable. The gateway never creates or alters it.
func (r *cutoverRig) createCustomTable(mode string) {
	r.t.Helper()
	floating := "REAL"
	if r.kind != "sqlite" {
		floating = "DOUBLE PRECISION"
	}
	column := map[string]string{
		"not_null": floating + " NOT NULL", "nullable": floating, "default": floating + " NOT NULL DEFAULT 0",
	}[mode]
	require.NotEmpty(r.t, column, mode)
	_, err := r.dest.ExecContext(r.t.Context(), `CREATE TABLE `+r.quoted("custom_readings")+` (a_col `+column+`, b_col `+floating+` NOT NULL)`)
	require.NoError(r.t, err)
}
