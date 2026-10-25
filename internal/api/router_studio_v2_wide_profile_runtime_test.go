package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

const wideProfileLayoutVersion = "managed_wide_profile_v1"

// wideProfileRequest is the public circuit-profile request of a managed wide
// profile group. A member is optional unless its source says it is required.
func wideProfileRequest(t *testing.T, f managedGroupFixture, sources []meanSource) map[string]any {
	t.Helper()
	request := meanLongRequest(t, f, sources)
	group := request["group"].(map[string]any)
	group["layout_version"] = wideProfileLayoutVersion
	for i, member := range group["members"].([]map[string]any) {
		member["required"] = sources[i].required
	}
	return request
}

// wideRig is one applied wide profile group over a real destination, built only
// through the public routes and services, with the production pipeline running.
type wideRig struct {
	f        managedGroupFixture
	kind     string
	sources  []meanSource
	groupID  string
	table    string
	columns  map[string]string
	snapshot *workspace.WriteGroupAppliedSnapshot
	frozen   map[string]workspace.WriteGroupMember
	clock    *atomic.Int64
	pipe     *grouppipeline.Pipeline
}

// newWideRig saves, previews, confirms and applies the group, then starts the
// production pipeline on an injected clock one second into the first window.
func newWideRig(t *testing.T, kind string, sources []meanSource) *wideRig {
	t.Helper()
	ctx := t.Context()
	f := newManagedDestinationFixture(t, kind)
	seedMeanSources(t, f, sources)
	for _, source := range sources {
		if source.measurementID != "" {
			insertMeasurement(t, f, source.measurementID, source.name, "epoch-1", "counter")
		}
	}
	created := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, wideProfileRequest(t, f, sources))
	require.Equal(t, http.StatusOK, created.Code, created.Body.String())
	_, saved := groupSaveData(t, decodeJSONBody(t, created))
	rig := &wideRig{f: f, kind: kind, sources: sources, groupID: saved["id"].(string), table: destinationOf(saved)["table_name"].(string), columns: map[string]string{}}
	for _, member := range saved["members"].([]any) {
		entry := member.(map[string]any)
		rig.columns[strings.TrimPrefix(entry["tag_id"].(string), "tag-")] = entry["target_column"].(string)
	}

	token := previewLong(t, f, rig.groupID)
	require.NotEmpty(t, token.Statements, "the wide profile layout is prepared from the frozen result types")
	require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, rig.groupID, token).Status)
	ready, err := f.groups.Readiness(ctx, rig.groupID)
	require.NoError(t, err)
	require.True(t, ready.Ready, "%+v", ready.Issues)
	applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+rig.groupID+"/apply", scopeBody(t, f, rig.groupID))
	require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())
	rig.snapshot, err = f.groups.ResolveAppliedAt(ctx, rig.groupID, time.Now().Add(time.Minute))
	require.NoError(t, err)
	rig.frozen = map[string]workspace.WriteGroupMember{}
	for _, member := range rig.snapshot.Group.Members {
		rig.frozen[member.TagID] = member
	}

	rig.clock = &atomic.Int64{}
	rig.clock.Store(rig.snapshot.EffectiveAt.Add(time.Second).UnixNano())
	rig.pipe = newMeanPipeline(t, f, rig.clock, "node/test", true)
	require.Equal(t, grouppipeline.StateActive, pipeStatus(rig.pipe, rig.groupID).State, "%+v", rig.pipe.Status())
	rig.clock.Store(rig.snapshot.EffectiveAt.Add(13 * time.Second).UnixNano())
	return rig
}

func (r *wideRig) at(offset time.Duration) time.Time { return r.snapshot.EffectiveAt.Add(offset) }

func (r *wideRig) envelope(source, id string, offset time.Duration, value any, quality schema.QualityFlag) measurement.SampleEnvelope {
	member := r.frozen["tag-"+source]
	observed := r.at(offset)
	return measurement.SampleEnvelope{
		SampleID: id, WorkspaceID: r.f.record.ID, DeviceID: "device-A", PointID: member.PointID, TagID: member.TagID,
		SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision,
		ObservedAt: observed, ReceivedAt: observed.Add(30 * time.Millisecond), TimeOrigin: "gateway", Value: value, Quality: quality,
	}
}

func (r *wideRig) accept(t *testing.T, source, id string, offset time.Duration, value any) {
	t.Helper()
	require.NoError(t, r.pipe.AcceptSample(t.Context(), r.envelope(source, id, offset, value, schema.QualityGood)), id)
}

// closeWindow moves the injected clock past the first window's end.
func (r *wideRig) closeWindow() { r.clock.Store(r.at(16 * time.Second).UnixNano()) }

func (r *wideRig) target(t *testing.T) *sql.DB {
	t.Helper()
	if r.kind == "postgres" {
		return r.f.target
	}
	target, err := sql.Open("sqlite", r.f.file)
	require.NoError(t, err)
	t.Cleanup(func() { _ = target.Close() })
	return target
}

func (r *wideRig) from() string { return qualified(r.kind, r.f.namespace, r.table) }

// rowCount counts the rows of the group's table.
func (r *wideRig) rowCount(t *testing.T, target *sql.DB) int {
	t.Helper()
	var n int
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM `+r.from()).Scan(&n))
	return n
}

type wideProvenance struct {
	Member     string      `json:"member"`
	Status     string      `json:"status"`
	Reason     string      `json:"reason"`
	SampleID   string      `json:"sample_id"`
	Role       string      `json:"role"`
	SourceType string      `json:"source_type"`
	ResultType string      `json:"result_type"`
	Mean       *storedMean `json:"mean"`
}

// wideProfileRow is the one row of a closed window: the member columns as text (NULL
// stays NULL), the bucket start and the member evidence by source name.
type wideProfileRow struct {
	values      map[string]sql.NullString
	bucketStart string
	evidence    map[string]wideProvenance
	rawEvidence string
}

// awaitRow waits for exactly one row and reads it back.
func (r *wideRig) awaitRow(t *testing.T) wideProfileRow {
	t.Helper()
	target := r.target(t)
	require.Eventually(t, func() bool {
		var n int
		return target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM `+r.from()).Scan(&n) == nil && n == 1
	}, 15*time.Second, 50*time.Millisecond, "the closed window writes exactly one row")
	row := wideProfileRow{values: map[string]sql.NullString{}, evidence: map[string]wideProvenance{}}
	names := make([]string, 0, len(r.sources))
	dest := make([]any, 0, 2+len(r.sources))
	dest = append(dest, &row.bucketStart, &row.rawEvidence)
	selects := make([]string, 0, 2+len(r.sources))
	selects = append(selects, `bucket_start`, `provenance`)
	values := make([]sql.NullString, len(r.sources))
	for i, source := range r.sources {
		names = append(names, source.name)
		selects = append(selects, `CAST("`+r.columns[source.name]+`" AS TEXT)`)
		dest = append(dest, &values[i])
	}
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT `+strings.Join(selects, ", ")+` FROM `+r.from()).Scan(dest...))
	for i, name := range names {
		row.values[name] = values[i]
	}
	var entries []wideProvenance
	require.NoError(t, json.Unmarshal([]byte(row.rawEvidence), &entries), row.rawEvidence)
	for _, entry := range entries {
		var key []string
		require.NoError(t, json.Unmarshal([]byte(entry.Member), &key), entry.Member)
		row.evidence[strings.TrimPrefix(key[2], "tag-")] = entry
	}
	require.Len(t, row.evidence, len(r.sources), "every member is described in the row's evidence")
	return row
}

// bucketLog reads the local closure record of the first window.
func (r *wideRig) bucketLog(t *testing.T) (kind, reason string) {
	t.Helper()
	require.Eventually(t, func() bool {
		return r.f.db.QueryRowContext(t.Context(), `SELECT kind, reason FROM wg_delivery_buckets WHERE group_id = ?`, r.groupID).Scan(&kind, &reason) == nil
	}, 15*time.Second, 50*time.Millisecond, "the window is recorded as closed")
	return kind, reason
}

// plainDecimal drops the trailing zeros PostgreSQL may keep from a NUMERIC scale,
// so an exact digit string compares by value on both destinations.
func plainDecimal(text string) string {
	if !strings.Contains(text, ".") {
		return text
	}
	return strings.TrimSuffix(strings.TrimRight(text, "0"), ".")
}

func wideKinds() []string { return []string{"sqlite", "postgres"} }

// TestWideProfileRecording drives a managed_wide_profile_v1 group through the
// public routes and services (save, schema preview, confirmation, Apply), then
// through the production pipeline into a real destination table, and reads the
// one row of the closed 15s window back. Every case runs on SQLite and, where a
// disposable PostgreSQL is configured, on PostgreSQL; an unconfigured PostgreSQL
// is a skip, never a pass.
func TestWideProfileRecording(t *testing.T) {
	for _, kind := range wideKinds() {
		t.Run(kind, func(t *testing.T) {
			t.Run("numeric gauges average, state and counter_raw keep the newest reading", func(t *testing.T) {
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
				for _, reading := range []struct {
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
					rig.accept(t, reading.source, reading.id, reading.offset, reading.value)
				}
				bad := rig.envelope("tri", "tri-bad", 4*time.Second, 1000, schema.QualityBad)
				bad.QualityReason = "dial tcp 10.0.0.5:502: secret-token"
				require.NoError(t, rig.pipe.AcceptSample(t.Context(), bad))
				rig.closeWindow()

				row := rig.awaitRow(t)
				require.Equal(t, rig.at(0).UTC().Format(time.RFC3339Nano), row.bucketStart, "one fixed UTC window")
				text := func(source string) string { return row.values[source].String }
				require.Equal(t, "1.5", plainDecimal(text("half")), "the exact mean of 1 and 2")
				require.Equal(t, "9007199254740992.5", plainDecimal(text("big")), "the exact mean above 2^53, never a float")
				require.Equal(t, "20", plainDecimal(text("tri")), "10, 20 and 30 with a resend and a bad read")
				require.Equal(t, "1.5", text("flt"))
				require.Equal(t, "30", text("state"), "a state keeps its newest reading, not the mean 20")
				require.Equal(t, "30", text("raw"), "a confirmed counter_raw keeps its newest raw value")
				require.False(t, row.values["quiet"].Valid, "an optional member without a sample is SQL NULL, never zero")

				half, big, tri, quiet := row.evidence["half"], row.evidence["big"], row.evidence["tri"], row.evidence["quiet"]
				require.Equal(t, [4]any{"mean", "int16", "decimal", "exact"}, [4]any{half.Mean.Reduction, half.Mean.SourceType, half.Mean.ResultType, half.Mean.Basis})
				require.Equal(t, [3]any{"unknown", int64(2), "3"}, [3]any{half.Mean.Coverage, half.Mean.ValidCount, half.Mean.Sum}, "no whole-window proof is never complete")
				require.Equal(t, [2]any{1, "half-even"}, [2]any{half.Mean.Scale, half.Mean.Rounding}, "the exact scale and rounding rule are carried as evidence")
				require.Equal(t, "18014398509481985", big.Mean.Sum, "the exact sum is carried as text")
				require.Equal(t, [4]any{"partial", int64(3), int64(1), "60"}, [4]any{tri.Mean.Coverage, tri.Mean.ValidCount, tri.Mean.RejectedCount, tri.Mean.Sum})
				require.Equal(t, []string{"bad"}, tri.Mean.Reasons)
				require.Equal(t, [3]any{"float64", "float32", "float64"}, [3]any{row.evidence["flt"].Mean.Basis, row.evidence["flt"].Mean.SourceType, row.evidence["flt"].Mean.ResultType})
				require.Equal(t, [3]string{"missing", "no-samples", "missing"}, [3]string{quiet.Status, quiet.Reason, quiet.Mean.Coverage})
				require.Equal(t, [3]string{"state", "int16", "int64"}, [3]string{row.evidence["state"].Role, row.evidence["state"].SourceType, row.evidence["state"].ResultType})
				require.Equal(t, [3]string{"counter_raw", "uint64", "uint64"}, [3]string{row.evidence["raw"].Role, row.evidence["raw"].SourceType, row.evidence["raw"].ResultType})
				for _, latest := range []string{"state", "raw"} {
					require.Nil(t, row.evidence[latest].Mean, "a latest member carries no mean block")
				}
				for _, leaked := range []string{"secret-token", "10.0.0.5", "dial tcp"} {
					require.NotContains(t, row.rawEvidence, leaked, "a private driver error never reaches operator evidence")
				}
				for source, want := range map[string]string{"half": "int16", "big": "uint64", "flt": "float32", "state": "int16", "raw": "uint64"} {
					var tagType string
					require.NoError(t, rig.f.db.QueryRowContext(t.Context(), `SELECT data_type FROM tags WHERE id = ?`, "tag-"+source).Scan(&tagType))
					require.Equal(t, want, tagType, "the Tag's own source type is never rewritten by the result")
				}
			})

			t.Run("a window where nothing was observed is still one NULL row", func(t *testing.T) {
				rig := newWideRig(t, kind, []meanSource{
					{name: "half", tagType: "int16", role: "gauge"}, {name: "state", tagType: "int16", role: "state"},
				})
				rig.closeWindow()
				row := rig.awaitRow(t)
				for _, source := range []string{"half", "state"} {
					require.False(t, row.values[source].Valid, "%s: no sample is NULL, never 0 or an older value", source)
					require.Equal(t, [2]string{"missing", "no-samples"}, [2]string{row.evidence[source].Status, row.evidence[source].Reason})
				}
				require.Equal(t, "missing", row.evidence["half"].Mean.Coverage)
				logKind, _ := rig.bucketLog(t)
				require.Equal(t, "row", logKind, "unlike a legacy wide group, a silent window is a row, not no_data")
			})

			t.Run("an explicit required member without a value skips the whole row and says why", func(t *testing.T) {
				rig := newWideRig(t, kind, []meanSource{
					{name: "half", tagType: "int16", role: "gauge", required: true}, {name: "state", tagType: "int16", role: "state"},
				})
				rig.accept(t, "state", "state-1", time.Second, 7)
				rig.closeWindow()
				logKind, reason := rig.bucketLog(t)
				require.Equal(t, [2]string{"skipped", "incomplete-required-member"}, [2]string{logKind, reason})
				require.Zero(t, rig.rowCount(t, rig.target(t)), "a skipped window writes no row")
			})

			t.Run("an optional member missing next to a required one that has data still writes the row", func(t *testing.T) {
				rig := newWideRig(t, kind, []meanSource{
					{name: "half", tagType: "int16", role: "gauge", required: true}, {name: "state", tagType: "int16", role: "state"},
				})
				rig.accept(t, "half", "half-1", time.Second, 4)
				rig.accept(t, "half", "half-2", 2*time.Second, 5)
				rig.closeWindow()
				row := rig.awaitRow(t)
				require.Equal(t, "4.5", plainDecimal(row.values["half"].String))
				require.False(t, row.values["state"].Valid)
				require.Equal(t, "no-samples", row.evidence["state"].Reason)
			})
		})
	}
}
