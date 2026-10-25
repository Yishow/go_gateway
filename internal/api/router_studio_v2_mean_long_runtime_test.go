package api

import (
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

	"github.com/stretchr/testify/require"
)

// meanTimeLayout is the fixed text a SQLite long table stores an instant in.
const meanTimeLayout = "2006-01-02T15:04:05.000000Z"

// meanSource is one fixed simulated source of the mean flow: a mapped point and
// tag of device-A with a chosen tag type and recording role. The sources are
// source data only; every applied fact below comes from the real Apply.
type meanSource struct {
	name, tagType, role, measurementID string
	// required marks a wide profile member whose missing value skips the row.
	required bool
}

var meanSources = []meanSource{
	{name: "half", tagType: "int16", role: "gauge"},
	{name: "big", tagType: "uint64", role: "gauge"},
	{name: "tri", tagType: "int16", role: "gauge"},
	{name: "flt", tagType: "float32", role: "gauge"},
	{name: "quiet", tagType: "int16", role: "gauge"},
	{name: "state", tagType: "int16", role: "state"},
	{name: "raw", tagType: "uint64", role: "counter_raw", measurementID: "meas-raw"},
}

func seedMeanSources(t *testing.T, f managedGroupFixture, sources []meanSource) {
	t.Helper()
	for _, source := range sources {
		seedProfileSource(t, f.writeGroupRouterFixture, source.name)
		_, err := f.db.ExecContext(t.Context(), `UPDATE tags SET data_type=? WHERE id=?`, source.tagType, "tag-"+source.name)
		require.NoError(t, err)
	}
}

// insertMeasurement saves a counter definition for one seeded source. A
// confirmed one needs a unit, an epoch and both revisions, never a name.
func insertMeasurement(t *testing.T, f managedGroupFixture, id, source, epoch, kind string) {
	t.Helper()
	_, err := f.db.ExecContext(t.Context(), `
		INSERT INTO measurement_definitions (
			id, workspace_id, device_id, point_id, tag_id, equipment_id, definition_revision, source_binding_revision,
			series_epoch, name, quantity, unit, semantic_kind
		) VALUES (?, ?, 'device-A', ?, ?, 'device-A', 'def-1', 'bind-1', ?, 'Energy', 'energy', 'kWh', ?)`,
		id, f.record.ID, "point-"+source, "tag-"+source, epoch, kind)
	require.NoError(t, err)
}

// meanLongRequest builds the public circuit-profile request of a managed long
// mean group: every member states its role and the server derives the types.
func meanLongRequest(t *testing.T, f managedGroupFixture, sources []meanSource) map[string]any {
	t.Helper()
	request := profileRequest(t, f.writeGroupRouterFixture, "instantaneous", 15, "second", "main", "")
	group := request["group"].(map[string]any)
	group["layout_version"] = "managed_long_v1"
	group["destination"].(map[string]any)["table_schema"] = f.namespace
	group["recording_profile"].(map[string]any)["sample_reduction"] = "mean"
	members := make([]map[string]any, 0, len(sources))
	for _, source := range sources {
		member := map[string]any{
			"device_id": "device-A", "point_id": "point-" + source.name, "tag_id": "tag-" + source.name, "required": true,
			"circuit_id": "main", "circuit_name": "Circuit main", "value_role": source.role,
		}
		if source.measurementID != "" {
			member["measurement_id"] = source.measurementID
		}
		members = append(members, member)
	}
	group["members"] = members
	return request
}

// newMeanPipeline builds the production pipeline over the real services. The
// counter provider is the production one, or absent for the missing-provider case.
func newMeanPipeline(t *testing.T, f managedGroupFixture, clock *atomic.Int64, owner string, counters bool) *grouppipeline.Pipeline {
	t.Helper()
	targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.db))
	deps := grouppipeline.Dependencies{
		Groups: f.groups, Tags: tag.NewService(tag.NewSQLRepository(f.db)), Destinations: targets,
		Inspector: dbtarget.NewReadOnlyTableInspector(targets), Store: groupdelivery.NewStore(f.db),
	}
	if counters {
		deps.Counters = grouppipeline.NewMeasurementCounterSources(measurement.NewSQLRepository(f.db))
	}
	pipe := grouppipeline.New(deps, grouppipeline.Config{
		NodeID: "node", Owner: owner, Now: func() time.Time { return time.Unix(0, clock.Load()).UTC() },
		TickInterval: 20 * time.Millisecond, WorkerInterval: 20 * time.Millisecond, ReconcileInterval: time.Hour, MaxFutureSkew: time.Minute,
	})
	require.NoError(t, pipe.Start(t.Context()))
	t.Cleanup(func() { _ = pipe.Stop(time.Second) })
	return pipe
}

type storedMean struct {
	Reduction       string   `json:"reduction"`
	Role            string   `json:"role"`
	SourceType      string   `json:"source_type"`
	ResultType      string   `json:"result_type"`
	Version         string   `json:"version"`
	Method          string   `json:"method"`
	Basis           string   `json:"basis"`
	ValidCount      int64    `json:"valid_count"`
	RejectedCount   int64    `json:"rejected_count"`
	FirstObservedAt string   `json:"first_observed_at"`
	LastObservedAt  string   `json:"last_observed_at"`
	Coverage        string   `json:"coverage"`
	Reason          string   `json:"reason"`
	Reasons         []string `json:"reasons"`
	Sum             string   `json:"sum"`
	Scale           int      `json:"scale"`
	Rounding        string   `json:"rounding"`
}

type storedProvenance struct {
	Reason     string          `json:"reason"`
	Role       string          `json:"role"`
	SourceType string          `json:"source_type"`
	ResultType string          `json:"result_type"`
	Mean       *storedMean     `json:"mean"`
	Counter    json.RawMessage `json:"counter"`
}

type storedMeanRow struct {
	status, valueType, bucketStart, rawProvenance  string
	decimal, integer, float, measurement, observed sql.NullString
	deltaValue, deltaFrom, deltaTo                 sql.NullString
	prov                                           storedProvenance
}

func readMeanRows(t *testing.T, target *sql.DB, table string) map[string]storedMeanRow {
	t.Helper()
	rows, err := target.QueryContext(t.Context(), `
		SELECT tag_id, status, value_type, bucket_start, CAST(value_decimal AS TEXT), CAST(value_int AS TEXT), CAST(value_float AS TEXT),
			measurement_id, observed_at, CAST(delta_value AS TEXT), CAST(delta_from AS TEXT), CAST(delta_to AS TEXT), provenance
		FROM `+table)
	require.NoError(t, err)
	defer rows.Close()
	got := map[string]storedMeanRow{}
	for rows.Next() {
		var tagID, provenance string
		var row storedMeanRow
		require.NoError(t, rows.Scan(&tagID, &row.status, &row.valueType, &row.bucketStart, &row.decimal, &row.integer, &row.float,
			&row.measurement, &row.observed, &row.deltaValue, &row.deltaFrom, &row.deltaTo, &provenance))
		require.NoError(t, json.Unmarshal([]byte(provenance), &row.prov), provenance)
		row.rawProvenance = provenance
		got[tagID] = row
	}
	require.NoError(t, rows.Err())
	return got
}

// TestMeanLongRecording saves, previews, confirms and applies a managed long
// mean group through the public routes and services, then runs the production
// pipeline over the proof Apply froze. Nothing about the applied state is built
// by the test: the closed 15s window is read back from the real SQLite table.
func TestMeanLongRecording(t *testing.T) {
	ctx := t.Context()
	f := newManagedDestinationFixture(t, "sqlite")
	seedMeanSources(t, f, meanSources)
	insertMeasurement(t, f, "meas-raw", "raw", "epoch-1", "counter")

	created := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, meanLongRequest(t, f, meanSources))
	require.Equal(t, http.StatusOK, created.Code, created.Body.String())
	_, saved := groupSaveData(t, decodeJSONBody(t, created))
	groupID, table := saved["id"].(string), destinationOf(saved)["table_name"].(string)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, groupID, previewLong(t, f, groupID)).Status)

	// The confirmed schema and connector stay as they are while a counter_raw
	// member's saved definition changes under the confirmed group: readiness
	// refuses a definition that no longer confirms counter semantics, or that now
	// names another source, and only that evidence is named. The group is ready
	// again once the definition is restored.
	readiness := func() (bool, []string) {
		result, err := f.groups.Readiness(ctx, groupID)
		require.NoError(t, err)
		codes := make([]string, 0, len(result.Issues))
		for _, issue := range result.Issues {
			codes = append(codes, issue.Code)
		}
		return result.Ready, codes
	}
	setRawMeasurement := func(epoch, point, tag string) {
		_, err := f.db.ExecContext(ctx, `UPDATE measurement_definitions SET series_epoch=?, point_id=?, tag_id=? WHERE id='meas-raw'`, epoch, point, tag)
		require.NoError(t, err)
	}
	for _, tc := range []struct{ name, epoch, point, tag, code string }{
		{"unconfirmed", "", "point-raw", "tag-raw", "counter-measurement-unconfirmed"},
		{"foreign", "epoch-1", "point-state", "tag-state", "counter-measurement-foreign"},
	} {
		setRawMeasurement(tc.epoch, tc.point, tc.tag)
		isReady, codes := readiness()
		t.Logf("%s counter evidence: readiness issues %v", tc.name, codes)
		require.False(t, isReady, tc.name)
		require.Contains(t, codes, tc.code, "%s: readiness names the counter evidence it refuses", tc.name)
		setRawMeasurement("epoch-1", "point-raw", "tag-raw")
	}
	isReady, codes := readiness()
	require.True(t, isReady, "%v", codes)
	applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/apply", scopeBody(t, f, groupID))
	require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())
	snapshot, err := f.groups.ResolveAppliedAt(ctx, groupID, time.Now().Add(time.Minute))
	require.NoError(t, err)

	frozen := map[string]workspace.WriteGroupMember{}
	for _, member := range snapshot.Group.Members {
		frozen[member.TagID] = member
	}
	for _, want := range []struct {
		source string
		role   workspace.WriteGroupValueRole
		from   schema.DataType
		to     measurement.ExactType
	}{
		{"half", "gauge", schema.DataTypeInt16, measurement.ExactDecimal},
		{"big", "gauge", schema.DataTypeUint64, measurement.ExactDecimal},
		{"tri", "gauge", schema.DataTypeInt16, measurement.ExactDecimal},
		{"flt", "gauge", schema.DataTypeFloat32, measurement.ExactFloat64},
		{"quiet", "gauge", schema.DataTypeInt16, measurement.ExactDecimal},
		{"state", "state", schema.DataTypeInt16, measurement.ExactInt64},
		{"raw", "counter_raw", schema.DataTypeUint64, measurement.ExactUint64},
	} {
		member := frozen["tag-"+want.source]
		require.Equal(t, [3]any{want.role, want.from, want.to}, [3]any{member.ValueRole, member.SourceValueType, member.ResultValueType}, want.source)
	}

	var clock atomic.Int64
	clock.Store(snapshot.EffectiveAt.Add(time.Second).UnixNano())

	// A confirmed counter_raw needs the production counter provider: without it,
	// or once its measurement no longer confirms counter semantics, the group is
	// never built, and a non-empty measurement id alone proves nothing.
	missing := newMeanPipeline(t, f, &clock, "node/no-provider", false)
	require.Equal(t, [2]string{grouppipeline.StateBlocked, "counter-semantics-required"},
		[2]string{pipeStatus(missing, groupID).State, pipeStatus(missing, groupID).Reason}, "%+v", missing.Status())
	require.NoError(t, missing.Stop(time.Second))
	_, err = f.db.ExecContext(ctx, `UPDATE measurement_definitions SET series_epoch='' WHERE id='meas-raw'`)
	require.NoError(t, err)
	downgraded := newMeanPipeline(t, f, &clock, "node/downgraded", true)
	require.Equal(t, [2]string{grouppipeline.StateBlocked, "counter-semantics-required"},
		[2]string{pipeStatus(downgraded, groupID).State, pipeStatus(downgraded, groupID).Reason}, "%+v", downgraded.Status())
	require.NoError(t, downgraded.Stop(time.Second))
	_, err = f.db.ExecContext(ctx, `UPDATE measurement_definitions SET series_epoch='epoch-1' WHERE id='meas-raw'`)
	require.NoError(t, err)

	pipe := newMeanPipeline(t, f, &clock, "node/test", true)
	require.Equal(t, grouppipeline.StateActive, pipeStatus(pipe, groupID).State, "%+v", pipe.Status())
	frozenLayout, err := groupdelivery.NewStore(f.db).RuntimeVersion(ctx, groupdelivery.GroupKey{
		WorkspaceID: snapshot.WorkspaceID, GroupID: groupID, GroupRevision: snapshot.AppliedRevision,
	})
	require.NoError(t, err)
	require.NotEmpty(t, frozenLayout, "Apply froze the verified layout the pipeline built from")

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
	accept := func(source, id string, offset time.Duration, value any, quality schema.QualityFlag) error {
		return pipe.AcceptSample(ctx, envelopeFor(source, id, offset, value, quality))
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
		{"tri", "tri-2", 2 * time.Second, 20}, // a resend of the same sample id never adds to the denominator
		{"flt", "flt-1", time.Second, 1.0}, {"flt", "flt-2", 2 * time.Second, 2.0},
		{"state", "state-1", time.Second, 10}, {"state", "state-2", 5 * time.Second, 30},
		{"raw", "raw-1", time.Second, uint64(10)}, {"raw", "raw-2", 5 * time.Second, uint64(30)},
	} {
		require.NoError(t, accept(sample.source, sample.id, sample.offset, sample.value, good), sample.id)
	}
	badRead := envelopeFor("tri", "tri-bad", 4*time.Second, 1000, schema.QualityBad)
	badRead.QualityReason = "dial tcp 10.0.0.5:502: secret-token"
	require.NoError(t, pipe.AcceptSample(ctx, badRead))

	// The process restarts with the window still open: the new pipeline rebuilds
	// the group from the frozen proof and replays the journaled samples, so the
	// closed window below is the one a restart would have produced.
	require.NoError(t, pipe.Stop(time.Second))
	pipe = newMeanPipeline(t, f, &clock, "node/test", true)
	require.Equal(t, grouppipeline.StateActive, pipeStatus(pipe, groupID).State, "%+v", pipe.Status())
	foreign := measurement.SampleEnvelope{
		SampleID: "raw-other-epoch", WorkspaceID: f.record.ID, DeviceID: "device-A", PointID: "point-raw", TagID: "tag-raw",
		SourceRevision: frozen["tag-raw"].SourceRevision, MappingRevision: frozen["tag-raw"].MappingRevision,
		ObservedAt: effective.Add(6 * time.Second), TimeOrigin: "gateway", Value: uint64(99), Quality: good, SeriesEpoch: "epoch-other",
	}
	require.Error(t, pipe.AcceptSample(ctx, foreign), "a reading of another counter series never enters the confirmed member")
	clock.Store(effective.Add(16 * time.Second).UnixNano())

	target, err := sql.Open("sqlite", f.file)
	require.NoError(t, err)
	t.Cleanup(func() { _ = target.Close() })
	require.Eventually(t, func() bool {
		var n int
		return target.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+table).Scan(&n) == nil && n == len(meanSources)
	}, 15*time.Second, 50*time.Millisecond, "the closed window writes one record per source")
	rows := readMeanRows(t, target, table)
	require.Len(t, rows, len(meanSources))

	at := func(offset time.Duration) string { return effective.Add(offset).UTC().Format(meanTimeLayout) }
	for _, source := range []struct {
		name, valueType, digits, sum, coverage, reason string
		valid, rejected                                int64
		scale                                          int
		last                                           time.Duration
	}{
		{"half", "decimal", "1.5", "3", "unknown", "coverage-unverified", 2, 0, 1, 2 * time.Second},
		{"big", "decimal", "9007199254740992.5", "18014398509481985", "unknown", "coverage-unverified", 2, 0, 1, 2 * time.Second},
		{"tri", "decimal", "20", "60", "partial", "bad", 3, 1, 0, 4 * time.Second},
	} {
		row := rows["tag-"+source.name]
		require.Equal(t, [4]string{"ok", source.valueType, source.digits, at(0)}, [4]string{row.status, row.valueType, row.decimal.String, row.bucketStart}, source.name)
		require.False(t, row.integer.Valid || row.float.Valid, "%s: only the decimal slot is written", source.name)
		mean := row.prov.Mean
		require.NotNil(t, mean, source.name)
		require.Equal(t, storedMean{
			Reduction: "mean", Role: "gauge", SourceType: string(frozen["tag-"+source.name].SourceValueType), ResultType: "decimal",
			Version: "arithmetic-mean-v1", Method: "arithmetic-mean", Basis: "exact", ValidCount: source.valid, RejectedCount: source.rejected,
			FirstObservedAt: at(time.Second), LastObservedAt: at(source.last), Coverage: source.coverage, Reason: source.reason,
			Reasons: mean.Reasons, Sum: source.sum, Scale: source.scale, Rounding: "half-even",
		}, *mean, source.name)
	}
	require.Equal(t, []string{"bad"}, rows["tag-tri"].prov.Mean.Reasons, "the bad reading is rejected with a stable reason, not zero")
	for _, leaked := range []string{"secret-token", "10.0.0.5", "dial tcp"} {
		require.NotContains(t, rows["tag-tri"].rawProvenance, leaked, "a private driver error never reaches operator evidence")
	}
	require.Empty(t, rows["tag-half"].prov.Mean.Reasons)

	float := rows["tag-flt"]
	require.Equal(t, [4]string{"ok", "float64", "1.5", "float64"}, [4]string{float.status, float.valueType, float.float.String, float.prov.Mean.Basis})
	require.False(t, float.decimal.Valid || float.integer.Valid)
	require.Equal(t, [2]string{"float32", "float64"}, [2]string{float.prov.Mean.SourceType, float.prov.Mean.ResultType})

	quiet := rows["tag-quiet"]
	require.Equal(t, [3]string{"missing", "decimal", "no-samples"}, [3]string{quiet.status, quiet.valueType, quiet.prov.Reason})
	require.False(t, quiet.decimal.Valid || quiet.integer.Valid || quiet.float.Valid || quiet.observed.Valid, "a silent window is SQL NULL, never zero")
	require.Equal(t, [3]any{"missing", "no-samples", int64(0)}, [3]any{quiet.prov.Mean.Coverage, quiet.prov.Mean.Reason, quiet.prov.Mean.ValidCount})
	require.Empty(t, quiet.prov.Mean.FirstObservedAt+quiet.prov.Mean.LastObservedAt, "no observation, no invented time")

	state, raw := rows["tag-state"], rows["tag-raw"]
	require.Equal(t, [3]string{"ok", "int64", "30"}, [3]string{state.status, state.valueType, state.integer.String}, "a state keeps its newest reading, not the mean 20")
	require.False(t, state.decimal.Valid)
	require.Equal(t, [3]string{"state", "int16", "int64"}, [3]string{state.prov.Role, state.prov.SourceType, state.prov.ResultType})
	require.Equal(t, [4]string{"ok", "uint64", "30", "meas-raw"}, [4]string{raw.status, raw.valueType, raw.decimal.String, raw.measurement.String},
		"a confirmed counter_raw keeps its newest raw value, not the mean 20")
	require.Equal(t, [3]string{"counter_raw", "uint64", "uint64"}, [3]string{raw.prov.Role, raw.prov.SourceType, raw.prov.ResultType})
	for _, latest := range []storedMeanRow{state, raw} {
		require.Nil(t, latest.prov.Mean, "a latest member carries no mean block")
		require.Empty(t, latest.prov.Counter, "counter_raw carries no counter window")
		require.False(t, latest.deltaValue.Valid || latest.deltaFrom.Valid || latest.deltaTo.Valid, "no delta or baseline is produced")
	}

	var tagTypes []string
	typeRows, err := f.db.QueryContext(ctx, `SELECT id || '=' || data_type FROM tags WHERE id LIKE 'tag-%' ORDER BY id`)
	require.NoError(t, err)
	defer typeRows.Close()
	for typeRows.Next() {
		var text string
		require.NoError(t, typeRows.Scan(&text))
		tagTypes = append(tagTypes, text)
	}
	require.NoError(t, typeRows.Err())
	require.Equal(t, []string{
		"tag-A=float32", "tag-big=uint64", "tag-flt=float32", "tag-half=int16", "tag-quiet=int16", "tag-raw=uint64", "tag-state=int16", "tag-tri=int16",
	}, tagTypes, "the Tag's own source type is never rewritten by the mean")
}

// TestMeanLongRecordingRefusesUnprovenCounterRaw pins the save-time proof of a
// counter_raw member: only a saved definition of the same workspace, device,
// point and tag that confirms counter semantics is accepted; a missing, blank,
// foreign, non-counter or evidence-less one is refused and nothing is saved.
func TestMeanLongRecordingRefusesUnprovenCounterRaw(t *testing.T) {
	f := newManagedDestinationFixture(t, "sqlite")
	seedMeanSources(t, f, []meanSource{{name: "raw", tagType: "uint64"}, {name: "other", tagType: "uint64"}})
	insertMeasurement(t, f, "meas-gauge", "raw", "epoch-1", "gauge")
	insertMeasurement(t, f, "meas-no-epoch", "raw", "", "counter")
	insertMeasurement(t, f, "meas-foreign", "other", "epoch-1", "counter")
	insertMeasurement(t, f, "meas-ok", "raw", "epoch-1", "counter")

	for _, tc := range []struct {
		name          string
		measurementID *string
		status        int
		code          string
	}{
		{"no measurement id", nil, http.StatusUnprocessableEntity, "WRITE_GROUP_INVALID"},
		{"an empty measurement id", ptr(""), http.StatusUnprocessableEntity, "WRITE_GROUP_INVALID"},
		{"a blank measurement id", ptr("   "), http.StatusUnprocessableEntity, "WRITE_GROUP_INVALID"},
		{"an unknown measurement", ptr("meas-ghost"), http.StatusNotFound, "WRITE_GROUP_NOT_FOUND"},
		{"another source's measurement", ptr("meas-foreign"), http.StatusNotFound, "WRITE_GROUP_NOT_FOUND"},
		{"a gauge measurement", ptr("meas-gauge"), http.StatusUnprocessableEntity, "WRITE_GROUP_INVALID"},
		{"a counter without series evidence", ptr("meas-no-epoch"), http.StatusUnprocessableEntity, "WRITE_GROUP_INVALID"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			source := meanSource{name: "raw", role: "counter_raw"}
			request := meanLongRequest(t, f, []meanSource{source})
			if tc.measurementID != nil {
				request["group"].(map[string]any)["members"].([]map[string]any)[0]["measurement_id"] = *tc.measurementID
			}
			response := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, request)
			require.Equal(t, tc.status, response.Code, response.Body.String())
			assertWriteGroupError(t, decodeJSONBody(t, response), tc.code)
			require.Zero(t, countProfileGroups(t, f.writeGroupRouterFixture), "nothing is saved")
		})
	}

	request := meanLongRequest(t, f, []meanSource{{name: "raw", role: "counter_raw", measurementID: "meas-ok"}})
	response := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	_, group := groupSaveData(t, decodeJSONBody(t, response))
	member := group["members"].([]any)[0].(map[string]any)
	require.Equal(t, [4]any{"counter_raw", "uint64", "uint64", "meas-ok"},
		[4]any{member["value_role"], member["source_value_type"], member["result_value_type"], member["measurement_id"]})
}

func ptr[T any](value T) *T { return &value }
