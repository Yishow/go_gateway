package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

// counterProvenanceView is the counter object both layouts keep in provenance:
// what the delta rests on, the policy evidence and the discontinuities. A wide row
// also keeps the status, the typed exact delta and its actual bounds here, where a
// long record keeps them in columns; every number is a string.
type counterProvenanceView struct {
	Status           string `json:"status"`
	Delta            string `json:"delta"`
	DeltaFrom        string `json:"delta_from"`
	DeltaTo          string `json:"delta_to"`
	Kind             string `json:"kind"`
	Basis            string `json:"basis"`
	RawText          string `json:"raw_text"`
	SpanTotal        bool   `json:"span_total"`
	RolloverModulus  string `json:"rollover_modulus"`
	MaxRatePerSecond string `json:"max_rate_per_second"`
	Unit             string `json:"unit"`
	SeriesEpoch      string `json:"series_epoch"`
	MeasurementID    string `json:"measurement_id"`
	Events           []struct {
		Kind string `json:"kind"`
	} `json:"events"`
}

// counterWindow is the logical outcome of one cumulative window, independent of
// the layout that stored it.
type counterWindow struct {
	status, raw        string
	delta              string
	hasDelta           bool
	deltaFrom, deltaTo time.Time
	spanTotal          bool
	basis, rollover    string
	events             []string
	// provenance is the stored document, kept for exact-text assertions; the
	// layouts differ in it by design, so parity compares everything else.
	provenance string
}

// counterScenario is one meter: its tag, its saved counter definition and the
// readings it sends, each in a 10 minute window counted from the first one.
type counterScenario struct {
	name     string
	tagType  string
	semantic measurement.SemanticKind
	encoding string
	unit     string
	policy   *measurement.CounterPolicy
	readings []counterReading
	want     []wantWindow
	// restart stops the pipeline after the first window closes and starts a new
	// one, the way a restarted gateway would: the baseline must come back from the
	// durable closure, not from memory.
	restart bool
	// wideHas are exact substrings the wide provenance document must contain.
	wideHas []string
}

type counterReading struct {
	window int
	value  any
}

// wantWindow is what the existing long engine must report; the wide layout must
// then report exactly the same logical window. An empty delta is SQL NULL.
type wantWindow struct{ status, raw, delta string }

func parseStoredInstant(t *testing.T, text string) time.Time {
	t.Helper()
	for _, layout := range []string{"2006-01-02T15:04:05.000000Z", "2006-01-02 15:04:05.999999-07", "2006-01-02 15:04:05.999999-07:00"} {
		if parsed, err := time.Parse(layout, text); err == nil {
			return parsed.UTC()
		}
	}
	require.Failf(t, "unparsable stored instant", "%q", text)
	return time.Time{}
}

// runCumulativeWindows saves, previews, confirms and applies one cumulative
// profile group on the given layout through the public routes and services, sends
// the scenario's readings through the production pipeline, and reads every closed
// window back from the real destination.
func runCumulativeWindows(t *testing.T, kind, layout string, sc counterScenario, storageModes ...string) []counterWindow {
	t.Helper()
	ctx := t.Context()
	f := newManagedDestinationFixture(t, kind, sc.tagType)
	tagID := "tag-A"
	require.NoError(t, measurement.NewService(measurement.NewSQLRepository(f.db)).CreateMeasurement(ctx, &measurement.MeasurementDefinition{
		ID: "meas-energy", WorkspaceID: f.record.ID, DeviceID: "device-A", PointID: "point-A", TagID: &tagID, Name: "Meter", Quantity: "meter",
		Unit: sc.unit, SemanticKind: sc.semantic, NumericEncoding: sc.encoding, CounterPolicy: sc.policy,
	}), "the exact counter evidence is saved through the measurement service")

	request := profileRequest(t, f.writeGroupRouterFixture, "cumulative", 10, "minute", "main", "")
	group := request["group"].(map[string]any)
	group["layout_version"] = layout
	if len(storageModes) > 0 && storageModes[0] == "timescale" {
		group["storage_mode"] = "timescale"
	}
	group["destination"].(map[string]any)["table_schema"] = f.namespace
	group["members"].([]map[string]any)[0]["measurement_id"] = "meas-energy"
	// A wide member is optional unless it is marked required: the long layout has
	// no such gate, so parity is over members that may be missing in a window.
	group["members"].([]map[string]any)[0]["required"] = false
	created := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, request)
	require.Equal(t, http.StatusOK, created.Code, created.Body.String())
	_, saved := groupSaveData(t, decodeJSONBody(t, created))
	groupID, table := saved["id"].(string), destinationOf(saved)["table_name"].(string)
	column, _ := saved["members"].([]any)[0].(map[string]any)["target_column"].(string)

	require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, groupID, previewLong(t, f, groupID)).Status,
		"the %s layout prepares and confirms its cumulative table", layout)
	ready, err := f.groups.Readiness(ctx, groupID)
	require.NoError(t, err)
	require.True(t, ready.Ready, "%+v", ready.Issues)
	applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/apply", scopeBody(t, f, groupID))
	require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())
	snapshot, err := f.groups.ResolveAppliedAt(ctx, groupID, time.Now().Add(time.Hour))
	require.NoError(t, err)
	member := snapshot.Group.Members[0]

	var clock atomic.Int64
	clock.Store(snapshot.EffectiveAt.Add(time.Second).UnixNano())
	pipe := newMeanPipeline(t, f, &clock, "node/test", true)
	require.Equal(t, grouppipeline.StateActive, pipeStatus(pipe, groupID).State, "%+v", pipe.Status())

	target := f.target
	if kind == "sqlite" {
		target = openSQLiteTarget(t, f.file)
	}
	from := qualified(kind, f.namespace, table)
	awaitRows := func(n int) {
		require.Eventually(t, func() bool {
			var count int
			return target.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+from).Scan(&count) == nil && count == n
		}, 15*time.Second, 50*time.Millisecond, "%d closed windows reach the %s table", n, layout)
	}
	windowEnd := func(window int) int64 {
		return snapshot.EffectiveAt.Add(time.Duration(window+1)*10*time.Minute + time.Second).UnixNano()
	}
	last, restarted := 0, false
	for i, reading := range sc.readings {
		if sc.restart && !restarted && reading.window > last {
			clock.Store(windowEnd(last))
			awaitRows(last + 1)
			require.NoError(t, pipe.Stop(time.Second))
			pipe = newMeanPipeline(t, f, &clock, "node/test", true)
			require.Equal(t, grouppipeline.StateActive, pipeStatus(pipe, groupID).State, "%+v", pipe.Status())
			restarted = true
		}
		observed := snapshot.EffectiveAt.Add(time.Duration(reading.window)*10*time.Minute + 9*time.Minute)
		clock.Store(observed.Add(time.Second).UnixNano())
		require.NoError(t, pipe.AcceptSample(ctx, measurement.SampleEnvelope{
			SampleID: "reading-" + string(rune('a'+i)), WorkspaceID: f.record.ID, DeviceID: "device-A", PointID: member.PointID, TagID: member.TagID,
			SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision,
			ObservedAt: observed, ReceivedAt: observed.Add(30 * time.Millisecond), TimeOrigin: "gateway",
			Value: reading.value, Quality: schema.QualityGood,
		}))
		last = max(last, reading.window)
	}
	clock.Store(windowEnd(last))
	awaitRows(last + 1)

	if len(storageModes) > 0 && storageModes[0] == "timescale" {
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
		require.GreaterOrEqual(t, chunks, 1, "timescale cumulative hypertable wrote real chunks")
		t.Logf("counter_timescale_catalog case=%q schema=%q table=%q dimensions=%d time_column=%q chunks=%d result=PASS", t.Name(), f.namespace, table, dimensions, timeColumn, chunks)
	}
	t.Logf("counter_destination case=%q scenario=%q schema=%q table=%q layout=%q storage=%q sql_rows=%d", t.Name(), sc.name, f.namespace, table, layout, snapshot.Group.StorageMode, last+1)

	if layout == "managed_long_v1" {
		return readLongCounterWindows(t, target, from)
	}
	return readWideCounterWindows(t, target, from, column)
}

func counterWindowOf(t *testing.T, counter *counterProvenanceView) counterWindow {
	t.Helper()
	window := counterWindow{}
	if counter == nil {
		return window
	}
	window.spanTotal, window.basis, window.rollover = counter.SpanTotal, counter.Basis, counter.RolloverModulus
	for _, event := range counter.Events {
		window.events = append(window.events, event.Kind)
	}
	return window
}

func readLongCounterWindows(t *testing.T, target *sql.DB, from string) []counterWindow {
	t.Helper()
	rows, err := target.QueryContext(t.Context(), `SELECT status, CAST(value_decimal AS TEXT), CAST(value_int AS TEXT), CAST(value_float AS TEXT), CAST(delta_value AS TEXT), CAST(delta_from AS TEXT), CAST(delta_to AS TEXT), CAST(provenance AS TEXT) FROM `+from+` ORDER BY bucket_start`)
	require.NoError(t, err)
	defer rows.Close()
	var windows []counterWindow
	for rows.Next() {
		var status string
		var decimal, integer, float, delta, deltaFrom, deltaTo sql.NullString
		var provenance string
		require.NoError(t, rows.Scan(&status, &decimal, &integer, &float, &delta, &deltaFrom, &deltaTo, &provenance))
		var document struct {
			Counter *counterProvenanceView `json:"counter"`
		}
		require.NoError(t, json.Unmarshal([]byte(provenance), &document), provenance)
		window := counterWindowOf(t, document.Counter)
		window.status, window.provenance = status, provenance
		for _, raw := range []sql.NullString{decimal, integer, float} {
			if raw.Valid {
				window.raw = raw.String
			}
		}
		if delta.Valid {
			window.delta, window.hasDelta = delta.String, true
		}
		if deltaFrom.Valid && deltaTo.Valid {
			window.deltaFrom, window.deltaTo = parseStoredInstant(t, deltaFrom.String), parseStoredInstant(t, deltaTo.String)
		}
		windows = append(windows, window)
	}
	require.NoError(t, rows.Err())
	return windows
}

func readWideCounterWindows(t *testing.T, target *sql.DB, from, column string) []counterWindow {
	t.Helper()
	rows, err := target.QueryContext(t.Context(), `SELECT CAST("`+column+`" AS TEXT), provenance FROM `+from+` ORDER BY bucket_start`)
	require.NoError(t, err)
	defer rows.Close()
	var windows []counterWindow
	for rows.Next() {
		var raw sql.NullString
		var provenance string
		require.NoError(t, rows.Scan(&raw, &provenance))
		var entries []struct {
			Status  string                 `json:"status"`
			Counter *counterProvenanceView `json:"counter"`
		}
		require.NoError(t, json.Unmarshal([]byte(provenance), &entries), provenance)
		require.Len(t, entries, 1)
		window := counterWindowOf(t, entries[0].Counter)
		window.status, window.raw, window.provenance = entries[0].Status, raw.String, provenance
		if counter := entries[0].Counter; counter != nil {
			window.status = counter.Status
			if counter.Delta != "" {
				window.delta, window.hasDelta = counter.Delta, true
			}
			if counter.DeltaFrom != "" && counter.DeltaTo != "" {
				window.deltaFrom, window.deltaTo = parseStoredInstant(t, counter.DeltaFrom), parseStoredInstant(t, counter.DeltaTo)
			}
		}
		windows = append(windows, window)
	}
	require.NoError(t, rows.Err())
	return windows
}

var counterScenarios = []counterScenario{
	{
		name: "electricity above 2^53 keeps an exact difference of 1", tagType: "uint64", semantic: measurement.SemanticKindCounter, encoding: "uint64", unit: "kWh",
		readings: []counterReading{{0, uint64(9007199254740992)}, {1, uint64(9007199254740993)}},
		want:     []wantWindow{{"baseline", "9007199254740992", ""}, {"ok", "9007199254740993", "1"}},
		wideHas:  []string{`"status":"ok"`, `"delta":"1"`, `"kind":"counter"`, `"basis":"exact"`, `"unit":"kWh"`, `"series_epoch":"epoch-1"`},
	},
	{
		name: "the baseline survives a restart", tagType: "uint64", semantic: measurement.SemanticKindCounter, encoding: "uint64", unit: "kWh",
		readings: []counterReading{{0, uint64(9007199254740992)}, {1, uint64(9007199254740993)}},
		want:     []wantWindow{{"baseline", "9007199254740992", ""}, {"ok", "9007199254740993", "1"}},
		restart:  true,
	},
	{
		name: "water keeps the float basis", tagType: "float64", semantic: measurement.SemanticKindCounter, encoding: "float64", unit: "m3",
		readings: []counterReading{{0, 100.25}, {1, 100.75}},
		want:     []wantWindow{{"baseline", "100.25", ""}, {"ok", "100.75", "0.5"}},
	},
	{
		name: "gas is an exact integer difference", tagType: "int64", semantic: measurement.SemanticKindCounter, encoding: "int64", unit: "m3",
		readings: []counterReading{{0, int64(1000)}, {1, int64(1042)}},
		want:     []wantWindow{{"baseline", "1000", ""}, {"ok", "1042", "42"}},
	},
	{
		name: "an unexplained decrease is NULL, never absolute or zero", tagType: "uint64", semantic: measurement.SemanticKindCounter, encoding: "uint64", unit: "kWh",
		readings: []counterReading{{0, uint64(100)}, {1, uint64(40)}},
		want:     []wantWindow{{"baseline", "100", ""}, {"reset_or_uncertain", "40", ""}},
	},
	{
		name: "a missing window is a gap-labeled span total over the real interval", tagType: "uint64", semantic: measurement.SemanticKindCounter, encoding: "uint64", unit: "kWh",
		readings: []counterReading{{0, uint64(100)}, {2, uint64(130)}},
		want:     []wantWindow{{"baseline", "100", ""}, {"missing", "", ""}, {"gap", "130", "30"}},
	},
	{
		name: "a signed counter keeps a negative difference", tagType: "int64", semantic: measurement.SemanticKindSignedCounter, encoding: "int64", unit: "kWh",
		policy:   &measurement.CounterPolicy{AllowNegative: true},
		readings: []counterReading{{0, int64(100)}, {1, int64(70)}},
		want:     []wantWindow{{"baseline", "100", ""}, {"ok", "70", "-30"}},
	},
	{
		name: "a proven rollover is compensated from the saved policy", tagType: "uint64", semantic: measurement.SemanticKindCounter, encoding: "uint64", unit: "kWh",
		policy:   &measurement.CounterPolicy{RolloverModulus: "65536", MaxRatePerSecond: "100"},
		readings: []counterReading{{0, uint64(65530)}, {1, uint64(4)}},
		want:     []wantWindow{{"baseline", "65530", ""}, {"ok", "4", "10"}},
	},
}

// TestWideCounterParity runs the same exact meters through a long and a wide
// profile group, each by the public preview, confirmation and Apply and the
// production pipeline, and requires the same logical windows: baseline NULL, the
// exact difference with its actual bounds, reset NULL, the gap span total, signed
// and rollover results. The long layout is the reference; the wide layout keeps
// the raw reading in its value column and the typed delta, bounds, status and
// policy in provenance.counter.
func TestWideCounterParity(t *testing.T) {
	for _, kind := range wideKinds() {
		t.Run(kind, func(t *testing.T) {
			for _, sc := range counterScenarios {
				t.Run(sc.name, func(t *testing.T) {
					long := runCumulativeWindows(t, kind, "managed_long_v1", sc)
					require.Len(t, long, len(sc.want))
					for i, want := range sc.want {
						require.Equal(t, [3]any{want.status, want.raw, want.delta}, [3]any{long[i].status, long[i].raw, long[i].delta},
							"the long reference, window %d", i)
						require.Equal(t, want.delta != "", long[i].hasDelta, "an absent difference is NULL, never zero: window %d", i)
					}

					wide := runCumulativeWindows(t, kind, wideProfileLayoutVersion, sc)
					for _, text := range sc.wideHas {
						require.Contains(t, wide[len(wide)-1].provenance, text, "the delta, status and policy are exact strings in provenance.counter")
					}
					if len(sc.wideHas) > 0 {
						require.NotRegexp(t, `"delta":[-0-9]`, wide[len(wide)-1].provenance, "no exact value passes through a JSON number")
					}
					for i := range long {
						long[i].provenance, wide[i].provenance = "", ""
					}
					require.Equal(t, long, wide, "both layouts store the same logical counter windows")
				})
			}
		})
	}
}
