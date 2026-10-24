package runtime

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

const counterWindow = 10 * time.Minute

// newCounterFixture is an applied managed_long_v1 cumulative group (ten minute
// windows) over a real SQLite long table, with one integer, one unsigned and one
// float counter, each with confirmed counter semantics.
func newCounterFixture(t *testing.T) *longFixture {
	t.Helper()
	f := newLongFixture(t)
	f.members = []workspace.WriteGroupMember{
		{DeviceID: "device-A", PointID: "point-int", TagID: "tag-int", SourceRevision: "src-int", MappingRevision: "map-int", CircuitID: "circuit-a", MeasurementID: ptr("meas-int")},
		{DeviceID: "device-A", PointID: "point-uint", TagID: "tag-uint", SourceRevision: "src-uint", MappingRevision: "map-uint", CircuitID: "circuit-a", MeasurementID: ptr("meas-uint")},
		{DeviceID: "device-A", PointID: "point-float", TagID: "tag-float", SourceRevision: "src-float", MappingRevision: "map-float", CircuitID: "circuit-a", MeasurementID: ptr("meas-float")},
	}
	f.config.Group.Members = f.members
	f.config.Group.RecordingProfile = &workspace.WriteGroupRecordingProfile{
		Kind: workspace.WriteGroupProfileCumulative, Interval: 10, Unit: workspace.WriteGroupIntervalUnitMinute,
	}
	f.config.Group.RowPolicy.IntervalSeconds = int(counterWindow / time.Second)
	f.config.TagTypes = map[string]schema.DataType{
		"tag-int": schema.DataTypeInt64, "tag-uint": schema.DataTypeUint64, "tag-float": schema.DataTypeFloat64,
	}
	f.config.Counters = map[string]CounterSource{}
	for _, member := range f.members {
		f.config.Counters[member.TagID] = CounterSource{
			MeasurementID: *member.MeasurementID, DefinitionRevision: "def-1", SourceBindingRevision: "bind-1",
			SeriesEpoch: "epoch-1", Unit: "kWh", Kind: aggregation.CounterKindMonotonic,
		}
	}
	return f
}

// insertRows commits the captured rows into the real long table the way the
// delivery layer does and returns them keyed by tag and window.
func (f *longFixture) insertRows(t *testing.T) {
	t.Helper()
	for _, row := range f.sink.rows {
		_, err := dbtarget.InsertGroupRow(t.Context(), f.db, dbtarget.GroupInsertRequest{
			Kind: schema.DatabaseConnectorKindSQLite, SchemaName: "main", TableName: longFixtureTable, Row: row.Encoded,
			PayloadDigest: "digest-" + row.Encoded.EffectKey, Strategy: dbtarget.GroupEffectReceipt, CommittedAt: "2026-10-24T00:00:30Z",
		})
		require.NoError(t, err)
	}
}

type storedCell struct {
	class string
	text  sql.NullString
}

// readCell reads one stored value together with its SQLite storage class, so a
// float detour (REAL) is distinguishable from an exact integer or text.
func readCell(t *testing.T, f *longFixture, tag string, bucket time.Duration, column string) storedCell {
	t.Helper()
	var cell storedCell
	require.NoError(t, f.db.QueryRowContext(t.Context(),
		fmt.Sprintf(`SELECT typeof(%[1]s), CAST(%[1]s AS TEXT) FROM %[2]s WHERE tag_id = ? AND bucket_start = ?`, column, longFixtureTable),
		tag, bucketText(bucket)).Scan(&cell.class, &cell.text))
	return cell
}

func (c storedCell) String() string { return c.class + ":" + c.text.String }

// TestCumulativeProfileExactDeltaThroughBoundary drives cumulative samples through the
// production boundary and reads the committed rows back from a real SQLite table.
func TestCumulativeProfileExactDeltaThroughBoundary(t *testing.T) {
	f := newCounterFixture(t)
	boundary := f.build(t)
	const first, second = counterWindow - time.Minute, 2*counterWindow - time.Minute

	f.accept(t, boundary,
		f.envelope(0, "int-0", at(first), int64(9007199254740992)),
		f.envelope(1, "uint-0", at(first), uint64(9007199254740992)),
		f.envelope(2, "float-0", at(first), 12000.4),
	)
	f.clock = at(counterWindow)
	require.NoError(t, boundary.Tick(t.Context(), f.clock))
	f.clock = at(2 * counterWindow)
	f.accept(t, boundary,
		f.envelope(0, "int-1", at(second), int64(9007199254740993)),
		f.envelope(1, "uint-1", at(second), uint64(9007199254740993)),
		f.envelope(2, "float-1", at(second), 12000.7),
	)
	require.NoError(t, boundary.Tick(t.Context(), at(2*counterWindow)))
	require.Len(t, f.sink.rows, 6, "one record per counter and window")
	f.insertRows(t)

	// The first window only establishes the baseline.
	for _, tag := range []string{"tag-int", "tag-uint", "tag-float"} {
		require.Equal(t, "text:baseline", readCell(t, f, tag, 0, "status").String(), tag)
		for _, column := range []string{"delta_value", "delta_from", "delta_to"} {
			require.Equal(t, "null", readCell(t, f, tag, 0, column).class, "%s %s is NULL, not zero", tag, column)
		}
	}
	// The raw readings keep their exact type and digits in their own slot.
	require.Equal(t, "integer:9007199254740992", readCell(t, f, "tag-int", 0, "value_int").String())
	require.Equal(t, "text:9007199254740992", readCell(t, f, "tag-uint", 0, "value_decimal").String())
	require.Equal(t, "real:12000.4", readCell(t, f, "tag-float", 0, "value_float").String())

	// The second window carries the raw reading and the exact difference.
	require.Equal(t, "integer:9007199254740993", readCell(t, f, "tag-int", counterWindow, "value_int").String())
	require.Equal(t, "text:9007199254740993", readCell(t, f, "tag-uint", counterWindow, "value_decimal").String())
	for _, tag := range []string{"tag-int", "tag-uint"} {
		require.Equal(t, "text:ok", readCell(t, f, tag, counterWindow, "status").String(), tag)
		require.Equal(t, "text:1", readCell(t, f, tag, counterWindow, "delta_value").String(), "%s delta is exactly 1", tag)
		require.Equal(t, "text:"+at(first).Format(longTimeLayout), readCell(t, f, tag, counterWindow, "delta_from").String())
		require.Equal(t, "text:"+at(second).Format(longTimeLayout), readCell(t, f, tag, counterWindow, "delta_to").String())
		require.Equal(t, "text:"+at(second).Format(longTimeLayout), readCell(t, f, tag, counterWindow, "observed_at").String())
	}
	require.Equal(t, "text:ok", readCell(t, f, "tag-float", counterWindow, "status").String())
	require.Equal(t, "text:0.3000000000010914", readCell(t, f, "tag-float", counterWindow, "delta_value").String(),
		"a float source keeps the float's own difference and is labeled as such in provenance")

	// Identity of the comparable series is written on the record.
	for _, column := range []struct{ name, want string }{
		{"series_epoch", "epoch-1"}, {"unit", "kWh"}, {"measurement_id", "meas-uint"}, {"profile_kind", "cumulative"}, {"value_type", "uint64"},
	} {
		require.Equal(t, "text:"+column.want, readCell(t, f, "tag-uint", counterWindow, column.name).String(), column.name)
	}
	provenance := readCell(t, f, "tag-uint", counterWindow, "provenance").text.String
	require.JSONEq(t, `{
		"type":"uint64","member":"[\"device-A\",\"point-uint\",\"tag-uint\"]","sample_id":"uint-1",
		"source_revision":"src-uint","mapping_revision":"map-uint",
		"profile":{"kind":"cumulative","interval":10,"unit":"minute"},
		"counter":{"kind":"counter","basis":"exact","baseline_sample_id":"uint-0","end_sample_id":"uint-1",
			"definition_revision":"def-1","source_binding_revision":"bind-1","profile_revision":"rev-1"}
	}`, provenance)
	floatProvenance := readCell(t, f, "tag-float", counterWindow, "provenance").text.String
	require.Contains(t, floatProvenance, `"basis":"float64"`)
}

// TestCumulativeProfileMissingStaysNullAndGapIsASpan keeps unproven windows empty: a
// silent window is NULL (not zero), and the window after it is a labeled gap
// whose difference is a span total over its real interval.
func TestCumulativeProfileMissingStaysNullAndGapIsASpan(t *testing.T) {
	f := newCounterFixture(t)
	boundary := f.build(t)
	f.accept(t, boundary, f.envelope(1, "uint-0", at(counterWindow-time.Minute), uint64(1000)))
	f.clock = at(3 * counterWindow)
	f.accept(t, boundary, f.envelope(1, "uint-2", at(3*counterWindow-time.Minute), uint64(1002)))
	require.NoError(t, boundary.Tick(t.Context(), at(3*counterWindow)))
	f.insertRows(t)

	require.Equal(t, "text:baseline", readCell(t, f, "tag-uint", 0, "status").String())
	silent := func(column string) string { return readCell(t, f, "tag-uint", counterWindow, column).String() }
	require.Equal(t, "text:missing", silent("status"))
	require.Equal(t, "null:", silent("value_decimal"))
	require.Equal(t, "null:", silent("delta_value"))
	require.Equal(t, "null:", silent("observed_at"), "no time is invented for an absent reading")

	gap := func(column string) string { return readCell(t, f, "tag-uint", 2*counterWindow, column).String() }
	require.Equal(t, "text:gap", gap("status"))
	require.Equal(t, "text:1002", gap("value_decimal"))
	require.Equal(t, "text:2", gap("delta_value"), "a span total over the real interval; the silent window above keeps none of it")
	require.Equal(t, "text:"+at(counterWindow-time.Minute).Format(longTimeLayout), gap("delta_from"))
	require.Equal(t, "text:"+at(3*counterWindow-time.Minute).Format(longTimeLayout), gap("delta_to"))
}

// TestCumulativeProfileLateSampleKeepsClosedDelta shows a late reading neither
// reopens a closed window nor changes the baseline it already advanced.
func TestCumulativeProfileLateSampleKeepsClosedDelta(t *testing.T) {
	f := newCounterFixture(t)
	boundary := f.build(t)
	f.accept(t, boundary, f.envelope(1, "uint-0", at(counterWindow-time.Minute), uint64(100)))
	f.clock = at(counterWindow)
	require.NoError(t, boundary.Tick(t.Context(), f.clock))

	late := f.envelope(1, "uint-late", at(counterWindow-30*time.Second), uint64(7))
	err := boundary.AcceptSample(t.Context(), late)
	var refused *GroupSampleError
	require.ErrorAs(t, err, &refused)
	require.Equal(t, "late", string(refused.Outcome))

	f.clock = at(2 * counterWindow)
	f.accept(t, boundary, f.envelope(1, "uint-1", at(2*counterWindow-time.Minute), uint64(110)))
	require.NoError(t, boundary.Tick(t.Context(), at(2*counterWindow)))
	f.insertRows(t)
	require.Equal(t, "text:10", readCell(t, f, "tag-uint", counterWindow, "delta_value").String(),
		"the late reading 7 never became part of the baseline")
}

// TestCumulativeProfileGateRefusals pins what a cumulative boundary will not
// guess: unconfirmed semantics, a non numeric source, mixed identities and a
// durable ledger whose baselines are not persisted yet.
func TestCumulativeProfileGateRefusals(t *testing.T) {
	code := func(t *testing.T, f *longFixture) string {
		t.Helper()
		_, err := NewGroupBoundary(f.config)
		var blocked *GroupBoundaryError
		require.ErrorAs(t, err, &blocked)
		return blocked.Code
	}
	t.Run("counter semantics must be confirmed for every member", func(t *testing.T) {
		f := newCounterFixture(t)
		delete(f.config.Counters, "tag-uint")
		require.Equal(t, "counter-semantics-required", code(t, f))
		f = newCounterFixture(t)
		f.config.Counters = nil
		require.Equal(t, "counter-semantics-required", code(t, f))
	})
	t.Run("a kind that is not counter or signed_counter is not semantics", func(t *testing.T) {
		f := newCounterFixture(t)
		source := f.config.Counters["tag-int"]
		source.Kind = "gauge"
		f.config.Counters["tag-int"] = source
		require.Equal(t, "counter-semantics-required", code(t, f))
	})
	t.Run("a measurement that disagrees with the member is refused", func(t *testing.T) {
		f := newCounterFixture(t)
		source := f.config.Counters["tag-int"]
		source.MeasurementID = "another-measurement"
		f.config.Counters["tag-int"] = source
		require.Equal(t, "counter-semantics-conflict", code(t, f))
	})
	t.Run("only numeric sources can be counters", func(t *testing.T) {
		f := newCounterFixture(t)
		f.config.TagTypes["tag-int"] = schema.DataTypeString
		require.Equal(t, "counter-type-unsupported", code(t, f))
	})
	t.Run("a sample claiming another epoch or definition is not compared", func(t *testing.T) {
		f := newCounterFixture(t)
		boundary := f.build(t)
		for field, mutate := range map[string]func(*measurement.SampleEnvelope){
			"epoch":      func(e *measurement.SampleEnvelope) { e.SeriesEpoch = "epoch-2" },
			"definition": func(e *measurement.SampleEnvelope) { e.DefinitionRevision = "def-9" },
			"binding":    func(e *measurement.SampleEnvelope) { e.SourceBindingRevision = "bind-9" },
			"measure":    func(e *measurement.SampleEnvelope) { e.MeasurementID = "other" },
		} {
			envelope := f.envelope(1, "uint-"+field, at(counterWindow-time.Minute), uint64(5))
			mutate(&envelope)
			var refused *GroupSampleError
			require.ErrorAs(t, boundary.AcceptSample(t.Context(), envelope), &refused, field)
			require.Equal(t, "revision-mismatch", refused.Reason, field)
		}
		matching := f.envelope(1, "uint-ok", at(counterWindow-time.Minute), uint64(5))
		matching.SeriesEpoch, matching.DefinitionRevision, matching.SourceBindingRevision = "epoch-1", "def-1", "bind-1"
		require.NoError(t, boundary.AcceptSample(t.Context(), matching), "an identity that matches is accepted")
	})
}
