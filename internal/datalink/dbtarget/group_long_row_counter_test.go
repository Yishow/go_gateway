package dbtarget

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

var counterWindowStart = time.Date(2026, 10, 24, 10, 0, 0, 0, time.UTC)

const counterWindowLength = 10 * time.Minute

func counterMemberFor(typ measurement.ExactType) LongRowMember {
	binding := aggregation.CounterBinding{
		WorkspaceID: "ws", GroupID: "g1", ProfileRevision: "r1", MemberKey: "m-kwh", MeasurementID: "meas-kwh",
		DefinitionRevision: "def-1", SourceBindingRevision: "bind-1", SeriesEpoch: "epoch-1",
		Kind: aggregation.CounterKindMonotonic, Unit: "kWh",
	}
	return LongRowMember{
		MemberKey: "m-kwh", DeviceID: "dev", PointID: "p-kwh", TagID: "t-kwh", CircuitID: "circuit-a",
		MeasurementID: binding.MeasurementID, Type: typ, Counter: &binding,
	}
}

func counterSpecFor(dialect SQLDialect, columns []ColumnInfo, member LongRowMember) LongRowSpec {
	return LongRowSpec{
		Dialect: dialect, Columns: columns, WorkspaceID: "ws", GroupID: "g1", GroupRevision: "r1",
		ProfileKind: "cumulative", ProfileInterval: 10, ProfileUnit: "minute", Members: []LongRowMember{member},
	}
}

// counterOutcomes runs the real exact engine over one reading per window and
// returns the closed outcomes the boundary would hand to the encoder; nil
// stands for a silent window.
func counterOutcomes(t *testing.T, member LongRowMember, values ...*measurement.ExactValue) []snapshot.Outcome {
	t.Helper()
	var checkpoint *aggregation.CounterCheckpoint
	outcomes := make([]snapshot.Outcome, 0, len(values))
	for i, value := range values {
		window := aggregation.CounterWindow{
			Start: counterWindowStart.Add(time.Duration(i) * counterWindowLength),
		}
		window.End = window.Start.Add(counterWindowLength)
		result := snapshot.MemberResult{MemberKey: member.MemberKey, Status: snapshot.MemberMissing, Reason: snapshot.ReasonNoSamples}
		if value != nil {
			observed := window.Start.Add(9*time.Minute + 123456789*time.Nanosecond)
			sample := snapshot.Sample{
				SampleID: fmt.Sprintf("s-%d", i), MemberKey: member.MemberKey, ObservedAt: observed,
				ReceivedAt: observed.Add(40 * time.Millisecond), Quality: schema.QualityGood, Value: *value,
				SourceRevision: "src-1", MappingRevision: "map-1",
			}
			counter, next, err := aggregation.CloseCounterWindow(checkpoint, *member.Counter, window,
				[]aggregation.CounterReading{{SampleID: sample.SampleID, ObservedAt: observed, ReceivedAt: sample.ReceivedAt, Value: *value}})
			require.NoError(t, err)
			checkpoint = next
			result = snapshot.MemberResult{MemberKey: member.MemberKey, Status: snapshot.MemberOK, Sample: &sample, Counter: &counter}
		}
		outcomes = append(outcomes, snapshot.Outcome{
			Kind: snapshot.OutcomeRow, EntityKey: member.MemberKey, RecordID: fmt.Sprintf("rec-%d", i), EffectKey: fmt.Sprintf("eff-%d", i),
			BucketStart: window.Start, BucketEnd: window.End, Members: []snapshot.MemberResult{result},
		})
	}
	return outcomes
}

func exactDecimal(t *testing.T, digits string) *measurement.ExactValue {
	t.Helper()
	value, err := measurement.NewDecimal(digits)
	require.NoError(t, err)
	return &value
}

func exactUint(v uint64) *measurement.ExactValue {
	value := measurement.NewUint64(v)
	return &value
}

func TestLongRowCounterEncodingCarriesRawAndExactDelta(t *testing.T) {
	member := counterMemberFor(measurement.ExactDecimal)
	outcomes := counterOutcomes(t, member, exactDecimal(t, "12000.0"), exactDecimal(t, "12000.4"), exactDecimal(t, "12000.7"))

	for _, dialect := range []SQLDialect{SQLDialectSQLite, SQLDialectPostgres} {
		t.Run(string(dialect), func(t *testing.T) {
			layout, issues := NewLongRowLayout(counterSpecFor(dialect, longColumnsFor(t, string(dialect)), member))
			require.Empty(t, issues)
			rows := make([]map[string]any, len(outcomes))
			for i, outcome := range outcomes {
				encoded, err := layout.EncodeRow(outcome)
				require.NoError(t, err)
				rows[i] = cellMap(encoded)
			}
			timeCell := func(offset time.Duration) any {
				observed := counterWindowStart.Add(offset).Add(9*time.Minute + 123456*time.Microsecond)
				if dialect == SQLDialectPostgres {
					return observed
				}
				return observed.Format(longTimeText)
			}

			// 12000.0 is a baseline: the raw digits survive, and there is no delta, not a zero one.
			require.Equal(t, "baseline", rows[0]["status"])
			require.Equal(t, "12000.0", rows[0]["value_decimal"], "the source's own digits are the stored value")
			require.Nil(t, rows[0]["delta_value"])
			require.Nil(t, rows[0]["delta_from"])
			require.Nil(t, rows[0]["delta_to"])
			require.Equal(t, []string{"value_decimal"}, valueSlotsOf(rows[0]))

			require.Equal(t, "ok", rows[1]["status"])
			require.Equal(t, "12000.4", rows[1]["value_decimal"])
			require.Equal(t, "0.4", rows[1]["delta_value"])
			require.Equal(t, timeCell(0), rows[1]["delta_from"], "the delta starts at the baseline reading's real time")
			require.Equal(t, timeCell(counterWindowLength), rows[1]["delta_to"])
			require.Equal(t, timeCell(counterWindowLength), rows[1]["observed_at"])

			require.Equal(t, "0.3", rows[2]["delta_value"])
			for _, row := range rows {
				require.Equal(t, "epoch-1", row["series_epoch"])
				require.Equal(t, "kWh", row["unit"])
				require.Equal(t, "meas-kwh", row["measurement_id"])
				require.Equal(t, "cumulative", row["profile_kind"])
				require.Equal(t, "decimal", row["value_type"])
			}
			require.JSONEq(t, `{
				"type":"decimal","member":"m-kwh","sample_id":"s-1","source_revision":"src-1","mapping_revision":"map-1",
				"profile":{"kind":"cumulative","interval":10,"unit":"minute"},
				"counter":{"kind":"counter","basis":"exact","raw_text":"12000.4","baseline_sample_id":"s-0","end_sample_id":"s-1",
					"definition_revision":"def-1","source_binding_revision":"bind-1","profile_revision":"r1"}
			}`, fmt.Sprint(rows[1]["provenance"]))
			require.Contains(t, rows[0]["provenance"], `"reason":"first-reading"`)
		})
	}
}

func valueSlotsOf(cells map[string]any) []string {
	var slots []string
	for _, column := range []string{"value_int", "value_decimal", "value_float", "value_bool", "value_text"} {
		if cells[column] != nil {
			slots = append(slots, column)
		}
	}
	return slots
}

func TestLongRowCounterEncodingKeepsLargeIntegersAndFloatBasis(t *testing.T) {
	layoutFor := func(t *testing.T, member LongRowMember) *LongRowLayout {
		layout, issues := NewLongRowLayout(counterSpecFor(SQLDialectSQLite, longColumnsFor(t, "sqlite"), member))
		require.Empty(t, issues)
		return layout
	}
	t.Run("uint64 above 2^53", func(t *testing.T) {
		member := counterMemberFor(measurement.ExactUint64)
		outcomes := counterOutcomes(t, member, exactUint(9007199254740992), exactUint(9007199254740993))
		encoded, err := layoutFor(t, member).EncodeRow(outcomes[1])
		require.NoError(t, err)
		cells := cellMap(encoded)
		require.Equal(t, "9007199254740993", cells["value_decimal"])
		require.Equal(t, "1", cells["delta_value"])
		require.Equal(t, []string{"value_decimal"}, valueSlotsOf(cells), "the exact decimal slot, never value_float")
	})
	t.Run("int64 above 2^53", func(t *testing.T) {
		member := counterMemberFor(measurement.ExactInt64)
		first, second := measurement.NewInt64(9007199254740992), measurement.NewInt64(9007199254740993)
		outcomes := counterOutcomes(t, member, &first, &second)
		encoded, err := layoutFor(t, member).EncodeRow(outcomes[1])
		require.NoError(t, err)
		cells := cellMap(encoded)
		require.Equal(t, int64(9007199254740993), cells["value_int"])
		require.Equal(t, "1", cells["delta_value"])
	})
	t.Run("a float source is labeled as a float difference", func(t *testing.T) {
		member := counterMemberFor(measurement.ExactFloat64)
		first, err := measurement.NewFloat64(12000.4)
		require.NoError(t, err)
		second, err := measurement.NewFloat64(12000.7)
		require.NoError(t, err)
		outcomes := counterOutcomes(t, member, &first, &second)
		cells := cellMap(outcomesEncoded(t, layoutFor(t, member), outcomes[1]))
		require.Equal(t, 12000.7, cells["value_float"])
		require.Equal(t, "0.3000000000010914", cells["delta_value"], "the float's own difference, not a reconstructed 0.3")
		require.Contains(t, cells["provenance"], `"basis":"float64"`)
		require.NotContains(t, cells["provenance"], "raw_text")
	})
}

func TestLongRowCounterEncodingNeverWritesAZeroForWhatIsUnknown(t *testing.T) {
	member := counterMemberFor(measurement.ExactUint64)
	layout, issues := NewLongRowLayout(counterSpecFor(SQLDialectSQLite, longColumnsFor(t, "sqlite"), member))
	require.Empty(t, issues)
	// 100 then 10 (a decrease), then a silent window, then 20 after the gap.
	outcomes := counterOutcomes(t, member, exactUint(100), exactUint(10), nil, exactUint(20))
	// Only the span after the silent window has a proven difference (20 against the
	// reset window's last reliable 10); everything unknown stays NULL.
	want := []struct{ status, delta string }{{"baseline", ""}, {"reset_or_uncertain", ""}, {"missing", ""}, {"gap", "10"}}
	for i, outcome := range outcomes {
		encoded, err := layout.EncodeRow(outcome)
		require.NoError(t, err)
		cells := cellMap(encoded)
		require.Equal(t, want[i].status, cells["status"], "window %d", i)
		if want[i].delta != "" {
			require.Equal(t, want[i].delta, cells["delta_value"], "window %d", i)
			continue
		}
		for _, column := range []string{"delta_value", "delta_from", "delta_to"} {
			require.Nil(t, cells[column], "window %d %s", i, column)
		}
		require.Equal(t, "epoch-1", cells["series_epoch"], "the series identity is on every record, even an empty one")
	}
	silent, err := layout.EncodeRow(outcomes[2])
	require.NoError(t, err)
	require.Nil(t, cellMap(silent)["value_decimal"])
	require.Nil(t, cellMap(silent)["observed_at"])
	require.Contains(t, cellMap(silent)["provenance"], `"reason":"no-samples"`)
	require.Contains(t, cellMap(outcomesEncoded(t, layout, outcomes[1]))["provenance"], `"reason":"decrease-unproven"`)

	t.Run("a counter result needs a bound member", func(t *testing.T) {
		unbound := member
		unbound.Counter = nil
		bare, issues := NewLongRowLayout(counterSpecFor(SQLDialectSQLite, longColumnsFor(t, "sqlite"), unbound))
		require.Empty(t, issues)
		_, err := bare.EncodeRow(outcomes[1])
		require.Equal(t, "counter-unbound", rowCode(t, err), "the delta is never dropped silently")
	})
}

func outcomesEncoded(t *testing.T, layout *LongRowLayout, outcome snapshot.Outcome) EncodedRow {
	t.Helper()
	encoded, err := layout.EncodeRow(outcome)
	require.NoError(t, err)
	return encoded
}

// TestLongRowCounterSQLiteRealCatalogReadback commits decimal progression rows
// into a real SQLite long table and reads each stored value with its storage
// class, proving the raw digits and exact deltas never become REAL.
func TestLongRowCounterSQLiteRealCatalogReadback(t *testing.T) {
	target := newInspectionSQLiteTarget(t, longDDL(t, "sqlite")...)
	fixture := newManagedDestinationFixture(t, schema.DatabaseConnectorKindSQLite, ConnectionConfig{"dsn": target})
	inspection, err := fixture.inspector.InspectTable(t.Context(), fixture.connector.ID, "main", indexTestTable)
	require.NoError(t, err)
	member := counterMemberFor(measurement.ExactDecimal)
	layout, issues := NewLongRowLayout(counterSpecFor(SQLDialectSQLite, inspection.Columns, member))
	require.Empty(t, issues)

	db, err := sql.Open("sqlite", target)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	for _, outcome := range counterOutcomes(t, member, exactDecimal(t, "12000.0"), exactDecimal(t, "12000.4"), exactDecimal(t, "12000.7")) {
		_, err := InsertGroupRow(t.Context(), db, GroupInsertRequest{
			Kind: schema.DatabaseConnectorKindSQLite, SchemaName: "main", TableName: indexTestTable,
			Row: outcomesEncoded(t, layout, outcome), PayloadDigest: "digest-" + outcome.RecordID,
			Strategy: GroupEffectReceipt, CommittedAt: "2026-10-24T00:00:30Z",
		})
		require.NoError(t, err)
	}
	read := func(record, column string) string {
		var class string
		var text sql.NullString
		require.NoError(t, db.QueryRowContext(t.Context(),
			fmt.Sprintf(`SELECT typeof(%[1]s), CAST(%[1]s AS TEXT) FROM %[2]s WHERE record_id = ?`, column, indexTestTable), record).Scan(&class, &text))
		return class + ":" + text.String
	}
	require.Equal(t, "text:12000.0", read("rec-0", "value_decimal"))
	require.Equal(t, "null:", read("rec-0", "delta_value"))
	require.Equal(t, "text:baseline", read("rec-0", "status"))
	require.Equal(t, "text:12000.4", read("rec-1", "value_decimal"))
	require.Equal(t, "text:0.4", read("rec-1", "delta_value"))
	require.Equal(t, "text:12000.7", read("rec-2", "value_decimal"))
	require.Equal(t, "text:0.3", read("rec-2", "delta_value"))
	require.Equal(t, "text:2026-10-24T10:19:00.123456Z", read("rec-1", "delta_to"))
	require.Equal(t, "text:2026-10-24T10:09:00.123456Z", read("rec-1", "delta_from"))
}
