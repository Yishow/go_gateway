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

type counterReadingAt struct {
	offset time.Duration
	value  measurement.ExactValue
}

// counterOutcomesFromWindows runs the real engine over the readings of each
// window (an empty window is silent) and returns the closed outcomes.
func counterOutcomesFromWindows(t *testing.T, member LongRowMember, windows ...[]counterReadingAt) []snapshot.Outcome {
	t.Helper()
	var checkpoint *aggregation.CounterCheckpoint
	outcomes := make([]snapshot.Outcome, 0, len(windows))
	for i, readings := range windows {
		window := aggregation.CounterWindow{Start: counterWindowStart.Add(time.Duration(i) * counterWindowLength)}
		window.End = window.Start.Add(counterWindowLength)
		result := snapshot.MemberResult{MemberKey: member.MemberKey, Status: snapshot.MemberMissing, Reason: snapshot.ReasonNoSamples}
		if len(readings) > 0 {
			var input []aggregation.CounterReading
			var newest snapshot.Sample
			for j, reading := range readings {
				observed := window.Start.Add(reading.offset)
				id := fmt.Sprintf("s-%d-%d", i, j)
				input = append(input, aggregation.CounterReading{SampleID: id, ObservedAt: observed, ReceivedAt: observed, Value: reading.value})
				if newest.ObservedAt.IsZero() || observed.After(newest.ObservedAt) {
					newest = snapshot.Sample{SampleID: id, MemberKey: member.MemberKey, ObservedAt: observed, ReceivedAt: observed, Quality: schema.QualityGood, Value: reading.value}
				}
			}
			counter, next, err := aggregation.CloseCounterWindow(checkpoint, *member.Counter, window, input)
			require.NoError(t, err)
			checkpoint = next
			result = snapshot.MemberResult{MemberKey: member.MemberKey, Status: snapshot.MemberOK, Sample: &newest, Counter: &counter}
		}
		outcomes = append(outcomes, snapshot.Outcome{
			Kind: snapshot.OutcomeRow, EntityKey: member.MemberKey, RecordID: fmt.Sprintf("rec-%d", i), EffectKey: fmt.Sprintf("eff-%d", i),
			BucketStart: window.Start, BucketEnd: window.End, Members: []snapshot.MemberResult{result},
		})
	}
	return outcomes
}

func uintAt(offset time.Duration, v uint64) counterReadingAt {
	return counterReadingAt{offset: offset, value: measurement.NewUint64(v)}
}

// discontinuityScenarios are the 2.3 fixtures as closed outcomes: a reset, a
// signed decrease, a proven wrap, a gap span and a flapping counter.
func discontinuityScenarios(t *testing.T) map[string]struct {
	member   LongRowMember
	outcomes []snapshot.Outcome
} {
	t.Helper()
	signed := counterMemberFor(measurement.ExactInt64)
	signedBinding := *signed.Counter
	signedBinding.Kind = aggregation.CounterKindSigned
	signed.Counter = &signedBinding

	wrap := counterMemberFor(measurement.ExactUint64)
	wrapBinding := *wrap.Counter
	modulus, rate := measurement.NewUint64(65536), measurement.NewUint64(100)
	wrapBinding.RolloverModulus, wrapBinding.MaxRatePerSecond = &modulus, &rate
	wrap.Counter = &wrapBinding

	plain := counterMemberFor(measurement.ExactUint64)
	flapping := make([]counterReadingAt, 0, 40)
	for i := range 40 {
		value := uint64(100)
		if i%2 == 1 {
			value = 10
		}
		flapping = append(flapping, uintAt(time.Duration(i)*10*time.Second, value))
	}
	ten, seven := measurement.NewInt64(10), measurement.NewInt64(7)
	type scenario = struct {
		member   LongRowMember
		outcomes []snapshot.Outcome
	}
	return map[string]scenario{
		"reset":    {plain, counterOutcomesFromWindows(t, plain, []counterReadingAt{uintAt(9*time.Minute, 90)}, []counterReadingAt{uintAt(2*time.Minute, 100), uintAt(5*time.Minute, 10), uintAt(8*time.Minute, 110)}, []counterReadingAt{uintAt(5*time.Minute, 125)})},
		"signed":   {signed, counterOutcomesFromWindows(t, signed, []counterReadingAt{{9 * time.Minute, ten}}, []counterReadingAt{{9 * time.Minute, seven}})},
		"rollover": {wrap, counterOutcomesFromWindows(t, wrap, []counterReadingAt{uintAt(9*time.Minute, 65530)}, []counterReadingAt{uintAt(9*time.Minute, 4)})},
		"gap":      {plain, counterOutcomesFromWindows(t, plain, []counterReadingAt{uintAt(0, 1000)}, nil, []counterReadingAt{uintAt(0, 1002)})},
		"flapping": {plain, counterOutcomesFromWindows(t, plain, []counterReadingAt{uintAt(9*time.Minute, 5)}, flapping)},
	}
}

func TestLongRowCounterEncodingCarriesEventsAndSpan(t *testing.T) {
	for _, dialect := range []SQLDialect{SQLDialectSQLite, SQLDialectPostgres} {
		t.Run(string(dialect), func(t *testing.T) {
			rows := map[string][]map[string]any{}
			for name, scenario := range discontinuityScenarios(t) {
				layout, issues := NewLongRowLayout(counterSpecFor(dialect, longColumnsFor(t, string(dialect)), scenario.member))
				require.Empty(t, issues)
				for _, outcome := range scenario.outcomes {
					rows[name] = append(rows[name], cellMap(outcomesEncoded(t, layout, outcome)))
				}
			}
			provenance := func(name string, window int) provenanceView {
				return decodeProvenance(t, rows[name][window]["provenance"])
			}

			reset := rows["reset"][1]
			require.Equal(t, "reset_or_uncertain", reset["status"])
			require.Equal(t, "110", reset["value_decimal"])
			require.Nil(t, reset["delta_value"])
			event := provenance("reset", 1).Counter.Events[0]
			require.Equal(t, [4]string{"decrease", "100", "10", "decrease-unproven"}, [4]string{event.Kind, event.FromValue, event.ToValue, event.Reason})
			require.Equal(t, "15", rows["reset"][2]["delta_value"], "the next window starts from the last reliable reading")

			require.Equal(t, "-3", rows["signed"][1]["delta_value"], "a negative exact delta is stored with its sign")
			require.Equal(t, int64(7), rows["signed"][1]["value_int"])
			require.Equal(t, "ok", rows["signed"][1]["status"])

			require.Equal(t, "10", rows["rollover"][1]["delta_value"])
			wrapped := provenance("rollover", 1).Counter
			require.Equal(t, "rollover", wrapped.Events[0].Kind)
			require.Equal(t, [2]string{"65536", "100"}, [2]string{wrapped.RolloverModulus, wrapped.MaxRatePerSecond})

			gap := rows["gap"]
			require.Equal(t, "missing", gap[1]["status"])
			require.Nil(t, gap[1]["delta_value"])
			require.Equal(t, "gap", gap[2]["status"])
			require.Equal(t, "2", gap[2]["delta_value"])
			require.True(t, provenance("gap", 2).Counter.SpanTotal)
			require.False(t, provenance("gap", 0).Counter.SpanTotal)
			if dialect == SQLDialectPostgres {
				require.Equal(t, counterWindowStart, gap[2]["delta_from"])
				require.Equal(t, counterWindowStart.Add(2*counterWindowLength), gap[2]["delta_to"])
			} else {
				require.Equal(t, counterWindowStart.Format(longTimeText), gap[2]["delta_from"])
				require.Equal(t, counterWindowStart.Add(2*counterWindowLength).Format(longTimeText), gap[2]["delta_to"])
			}

			flapping := provenance("flapping", 1).Counter
			require.Len(t, flapping.Events, 16)
			require.Equal(t, 4, flapping.EventsOmitted)
		})
	}
}

// TestLongRowCounterDiscontinuitiesSQLiteReadback commits every scenario into a
// real SQLite long table and reads the stored values back with their classes.
func TestLongRowCounterDiscontinuitiesSQLiteReadback(t *testing.T) {
	target := newInspectionSQLiteTarget(t, longDDL(t, "sqlite")...)
	fixture := newManagedDestinationFixture(t, schema.DatabaseConnectorKindSQLite, ConnectionConfig{"dsn": target})
	inspection, err := fixture.inspector.InspectTable(t.Context(), fixture.connector.ID, "main", indexTestTable)
	require.NoError(t, err)
	db, err := sql.Open("sqlite", target)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	for name, scenario := range discontinuityScenarios(t) {
		layout, issues := NewLongRowLayout(counterSpecFor(SQLDialectSQLite, inspection.Columns, scenario.member))
		require.Empty(t, issues)
		for _, outcome := range scenario.outcomes {
			outcome.RecordID, outcome.EffectKey = name+"-"+outcome.RecordID, name+"-"+outcome.EffectKey
			row := outcomesEncoded(t, layout, outcome)
			_, err := InsertGroupRow(t.Context(), db, GroupInsertRequest{
				Kind: schema.DatabaseConnectorKindSQLite, SchemaName: "main", TableName: indexTestTable, Row: row,
				PayloadDigest: "digest-" + row.RecordID, Strategy: GroupEffectReceipt, CommittedAt: "2026-10-24T00:00:30Z",
			})
			require.NoError(t, err, row.RecordID)
		}
	}
	read := func(record, column string) string {
		var class string
		var text sql.NullString
		require.NoError(t, db.QueryRowContext(t.Context(),
			fmt.Sprintf(`SELECT typeof(%[1]s), CAST(%[1]s AS TEXT) FROM %[2]s WHERE record_id = ?`, column, indexTestTable), record).Scan(&class, &text))
		return class + ":" + text.String
	}
	require.Equal(t, "text:reset_or_uncertain", read("reset-rec-1", "status"))
	require.Equal(t, "null:", read("reset-rec-1", "delta_value"))
	require.Equal(t, "text:15", read("reset-rec-2", "delta_value"))
	require.Equal(t, "text:-3", read("signed-rec-1", "delta_value"))
	require.Equal(t, "integer:7", read("signed-rec-1", "value_int"))
	require.Equal(t, "text:10", read("rollover-rec-1", "delta_value"))
	require.Equal(t, "null:", read("gap-rec-1", "delta_value"))
	require.Equal(t, "text:2", read("gap-rec-2", "delta_value"))
	require.Equal(t, "text:gap", read("gap-rec-2", "status"))
	require.Equal(t, "text:2026-10-24T10:20:00.000000Z", read("gap-rec-2", "delta_to"))
	require.Equal(t, "text:2026-10-24T10:00:00.000000Z", read("gap-rec-2", "delta_from"))
	require.Equal(t, "integer:1", read("gap-rec-2", "json_extract(provenance, '$.counter.span_total')"))
}
