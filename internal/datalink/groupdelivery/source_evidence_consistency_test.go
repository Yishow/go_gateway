package groupdelivery

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"

	"github.com/stretchr/testify/require"
)

func TestSourceRecordsRefuseAnythingTheFrozenPayloadDoesNotProve(t *testing.T) {
	cases := []struct {
		name    string
		cause   SourceUnavailableCause
		breakIt func(t *testing.T, db *sql.DB, effect string, row dbtarget.EncodedRow)
	}{
		{"digest differs", SourceCauseDigestMismatch, func(t *testing.T, db *sql.DB, effect string, _ dbtarget.EncodedRow) {
			_, err := db.ExecContext(t.Context(), `UPDATE wg_delivery_outbox SET payload_digest = 'x' WHERE effect_key = ?`, effect)
			require.NoError(t, err)
		}},
		{"payload is not json", SourceCausePayloadInvalid, func(t *testing.T, db *sql.DB, effect string, _ dbtarget.EncodedRow) {
			rewritePayload(t, db, effect, "{not json")
		}},
		{"another device", SourceCauseIdentityMismatch, func(t *testing.T, db *sql.DB, effect string, row dbtarget.EncodedRow) {
			changeCell(t, db, effect, row, "device_id", "dev-other")
		}},
		{"another revision", SourceCauseIdentityMismatch, func(t *testing.T, db *sql.DB, effect string, row dbtarget.EncodedRow) {
			changeCell(t, db, effect, row, "group_revision", "rev-other")
		}},
		{"another workspace", SourceCauseIdentityMismatch, func(t *testing.T, db *sql.DB, effect string, row dbtarget.EncodedRow) {
			changeCell(t, db, effect, row, "workspace_id", "workspace-other")
		}},
		{"another profile", SourceCauseIdentityMismatch, func(t *testing.T, db *sql.DB, effect string, row dbtarget.EncodedRow) {
			changeCell(t, db, effect, row, "profile_kind", "cumulative")
		}},
		{"an unknown status", SourceCauseStatusUnknown, func(t *testing.T, db *sql.DB, effect string, row dbtarget.EncodedRow) {
			changeCell(t, db, effect, row, "status", "great")
		}},
		{"a window that ends before it starts", SourceCauseWindowInvalid, func(t *testing.T, db *sql.DB, effect string, row dbtarget.EncodedRow) {
			changeCell(t, db, effect, row, "bucket_end", timeText(at(-60)))
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store, db, _ := newTestStore(t)
			spec := longSpec{source: sourceA, start: 0, status: "ok", slots: map[string]any{"value_float": 1.0}}
			effect := enqueueLong(t, store, spec)
			row, _ := spec.row(t)
			tc.breakIt(t, db, effect, row)
			got := readSource(t, store, sourceA)
			require.NotNil(t, got.Latest)
			require.False(t, got.Latest.Available)
			require.Equal(t, tc.cause, got.Latest.Cause)
			require.Nil(t, got.Latest.VerifiedSourceRecord, "an unverifiable record exposes no result")
			require.Nil(t, got.LastSQLCommitted)
		})
	}
}

func changeCell(t *testing.T, db *sql.DB, effect string, row dbtarget.EncodedRow, column string, value any) {
	t.Helper()
	for i := range row.Cells {
		if row.Cells[i].Column == column {
			row.Cells[i].Value = value
		}
	}
	payload, _, err := encodeRowPayload(row)
	require.NoError(t, err)
	rewritePayload(t, db, effect, string(payload))
}

// cutRow closes a window a forward cutover cut short, shaped as the assembler, the long row encoder and the
// outbox shape it: the outbox keeps the nominal bucket start, the payload and the table keep the real bounds.
func cutRow(t *testing.T, store *Store, truncated string, shiftStart, endBefore time.Duration, tweak func(window map[string]any)) {
	t.Helper()
	row, outcome := (longSpec{source: sourceA, status: "partial", slots: map[string]any{"value_float": 21.5}}).row(t)
	nominalStart, nominalEnd := outcome.BucketStart, outcome.BucketEnd
	actualStart, actualEnd := nominalStart.Add(shiftStart), nominalEnd.Add(-endBefore)
	outcome.Truncated, outcome.WindowStart, outcome.WindowEnd = truncated, actualStart, actualEnd
	row.BucketStart = actualStart
	window := map[string]any{
		"truncated": truncated, "start": timeText(actualStart), "end": timeText(actualEnd),
		"nominal_start": timeText(nominalStart), "nominal_end": timeText(nominalEnd),
	}
	if tweak != nil {
		tweak(window)
	}
	provenance, err := json.Marshal(map[string]any{
		"type": "float64", "member": row.EntityKey, "profile": map[string]any{"kind": "instantaneous", "interval": 15, "unit": "second"}, "window": window,
	})
	require.NoError(t, err)
	reviewSetCells(&row, map[string]any{"bucket_start": timeText(actualStart), "bucket_end": timeText(actualEnd), "provenance": string(provenance)})
	require.NoError(t, store.CommitClosure(t.Context(), testKey, Closure{
		Destination: testDestination, NextClose: nominalEnd, Buckets: []ClosedBucket{{Outcome: outcome, Row: &row}},
	}))
}

func TestSourceRecordsAcceptOnlyProvenCutoverWindows(t *testing.T) {
	for _, tc := range []struct {
		name, truncated       string
		shiftStart, endBefore time.Duration
		wantStart, wantEndGap time.Duration
	}{
		{"start cut", "start", 7500 * time.Millisecond, 0, 7500 * time.Millisecond, 0},
		{"end cut", "end", 0, 4 * time.Second, 0, 4 * time.Second},
		{"both cut", "both", 3 * time.Second, 5 * time.Second, 3 * time.Second, 5 * time.Second},
		{"end cut by one microsecond", "end", 0, time.Microsecond, 0, time.Microsecond},
	} {
		t.Run(tc.name, func(t *testing.T) {
			store, _, _ := newTestStore(t)
			cutRow(t, store, tc.truncated, tc.shiftStart, tc.endBefore, nil)
			got := readSource(t, store, sourceA).Latest
			require.True(t, got.Available, "%+v", got)
			require.Equal(t, at(0).Add(tc.wantStart), got.WindowStart, "the real bound, not the nominal bucket")
			require.Equal(t, at(15).Add(-tc.wantEndGap), got.WindowEnd)
		})
	}
}

func TestSourceRecordsRefuseWindowsTheCutoverContractDoesNotProve(t *testing.T) {
	for _, tc := range []struct {
		name, truncated       string
		shiftStart, endBefore time.Duration
		tweak                 func(map[string]any)
	}{
		{"claims a start cut that did not move the start", "start", 0, 0, nil},
		{"claims an end cut that did not move the end", "end", 0, 0, nil},
		{"wrong cut side", "end", 3 * time.Second, 0, nil},
		{"unknown cut kind", "middle", 3 * time.Second, 0, nil},
		{"nominal start is not the outbox bucket", "start", 3 * time.Second, 0, func(w map[string]any) { w["nominal_start"] = timeText(at(-30)) }},
		{"provenance bounds differ from the payload", "start", 3 * time.Second, 0, func(w map[string]any) { w["start"] = timeText(at(5)) }},
		{"nominal bucket ends before it starts", "start", 3 * time.Second, 0, func(w map[string]any) { w["nominal_end"] = timeText(at(-60)) }},
		{"actual window outside the nominal bucket", "end", 0, -time.Second, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			store, _, _ := newTestStore(t)
			cutRow(t, store, tc.truncated, tc.shiftStart, tc.endBefore, tc.tweak)
			got := readSource(t, store, sourceA).Latest
			require.False(t, got.Available, "%+v", got)
			require.Equal(t, SourceCauseWindowInvalid, got.Cause)
		})
	}
}

func TestSourceRecordsRefuseAShiftedWindowWithoutCutoverEvidence(t *testing.T) {
	store, db, _ := newTestStore(t)
	spec := longSpec{source: sourceA, status: "ok", slots: map[string]any{"value_float": 1.0}}
	effect := enqueueLong(t, store, spec)
	row, _ := spec.row(t)
	row.BucketStart = at(5)
	changeCell(t, db, effect, row, "bucket_start", timeText(at(5)))
	got := readSource(t, store, sourceA).Latest
	require.False(t, got.Available)
	require.Equal(t, SourceCauseWindowInvalid, got.Cause, "a window that differs from the bucket without proof is not accepted")
}

func TestSourceRecordsAreProofOnlyForTheFrozenAppliedDestination(t *testing.T) {
	store, _, _ := newTestStore(t)
	effect := enqueueLong(t, store, longSpec{source: sourceA, start: 0, status: "ok", slots: map[string]any{"value_float": 1.0}})
	confirm(t, store, effect)
	require.True(t, readSource(t, store, sourceA).LastSQLCommitted.Available, "the applied destination is the scope")

	for _, drift := range []struct{ name, column string }{{"connector id", "connector_id"}, {"connector revision", "connector_revision"}} {
		t.Run("the outbox "+drift.name+" drifted after the receipt", func(t *testing.T) {
			store, db, _ := newTestStore(t)
			confirm(t, store, enqueueLong(t, store, longSpec{source: sourceA, start: 0, status: "ok", slots: map[string]any{"value_float": 1.0}}))
			_, err := db.ExecContext(t.Context(), "UPDATE wg_delivery_outbox SET "+drift.column+" = 'drifted'")
			require.NoError(t, err)
			got := readSource(t, store, sourceA)
			require.False(t, got.Latest.Available)
			require.Equal(t, SourceCauseIdentityMismatch, got.Latest.Cause)
			require.Nil(t, got.LastSQLCommitted, "a receipt for another destination is not SQL evidence here")
		})
	}
	t.Run("a scope that names another destination gets no proof", func(t *testing.T) {
		other := scope
		other.ConnectorRevision = "a-draft-destination-revision"
		got, err := store.SourceRecords(t.Context(), other, sourceA)
		require.NoError(t, err)
		require.False(t, got.Latest.Available)
		require.Nil(t, got.LastSQLCommitted)
	})
	t.Run("a scope without the frozen destination is refused", func(t *testing.T) {
		missing := scope
		missing.ConnectorID = ""
		_, err := store.SourceRecords(t.Context(), missing, sourceA)
		require.ErrorIs(t, err, ErrInvalidGroupKey)
	})
}

func TestSourceRecordsRefuseContradictoryStoredResults(t *testing.T) {
	cumulative := scope
	cumulative.ProfileKind = "cumulative"
	decimal := map[string]any{"value_decimal": "100"}
	for _, tc := range []struct {
		name  string
		spec  longSpec
		scope SourceScope
		edit  map[string]any
		cause SourceUnavailableCause
	}{
		{"two value slots hold a reading", longSpec{source: sourceA, status: "ok", slots: map[string]any{"value_float": 1.5}}, scope, map[string]any{"value_int": int64(2)}, SourceCausePayloadInvalid},
		{"the reading sits in the slot of another type", longSpec{source: sourceA, status: "ok", slots: map[string]any{"value_float": 1.5}}, scope, map[string]any{"value_type": "int64"}, SourceCausePayloadInvalid},
		{"an unknown value type", longSpec{source: sourceA, status: "ok", slots: map[string]any{"value_float": 1.5}}, scope, map[string]any{"value_type": "complex"}, SourceCausePayloadInvalid},
		{"provenance names another value type", longSpec{source: sourceA, status: "ok", slots: map[string]any{"value_float": 1.5}}, scope, map[string]any{"value_type": "decimal"}, SourceCausePayloadInvalid},
		{"a stored reading with status ok but no value", longSpec{source: sourceA, status: "ok"}, scope, nil, SourceCauseStatusUnknown},
		{"a bad reading that stores a value", longSpec{source: sourceA, status: "bad", reason: "quality-bad", slots: map[string]any{"value_float": 1.5}}, scope, nil, SourceCauseStatusUnknown},
		{"an instantaneous record that stores a difference", longSpec{source: sourceA, status: "ok", slots: map[string]any{"value_float": 1.5}, delta: true}, scope, nil, SourceCauseStatusUnknown},
		{"a baseline that stores a difference", longSpec{source: sourceA, status: "baseline", profile: "cumulative", slots: decimal, delta: true}, cumulative, nil, SourceCauseStatusUnknown},
		{"an uncertain result that stores a difference", longSpec{source: sourceA, status: "reset_or_uncertain", profile: "cumulative", slots: decimal, delta: true}, cumulative, nil, SourceCauseStatusUnknown},
		{"a difference without its bounds", longSpec{source: sourceA, status: "ok", profile: "cumulative", slots: decimal, delta: true}, cumulative, map[string]any{"delta_from": nil, "delta_to": nil}, SourceCausePayloadInvalid},
		{"a difference that is not a number", longSpec{source: sourceA, status: "ok", profile: "cumulative", slots: decimal, delta: true}, cumulative, map[string]any{"delta_value": "lots"}, SourceCausePayloadInvalid},
		{"a cumulative record of a text value", longSpec{source: sourceA, status: "ok", profile: "cumulative", slots: map[string]any{"value_text": "x"}}, cumulative, nil, SourceCausePayloadInvalid},
	} {
		t.Run(tc.name, func(t *testing.T) {
			store, db, _ := newTestStore(t)
			effect := enqueueLong(t, store, tc.spec)
			if tc.edit != nil {
				row, _ := tc.spec.row(t)
				editCells(t, db, effect, row, tc.edit)
			}
			got, err := store.SourceRecords(t.Context(), tc.scope, sourceA)
			require.NoError(t, err)
			require.NotNil(t, got.Latest)
			require.False(t, got.Latest.Available, "%+v", got.Latest)
			require.Equal(t, tc.cause, got.Latest.Cause)
		})
	}

	t.Run("zero, false and a stored difference are valid", func(t *testing.T) {
		store, _, _ := newTestStore(t)
		enqueueLong(t, store, longSpec{source: sourceA, start: 0, status: "ok", profile: "cumulative", slots: map[string]any{"value_decimal": "0"}, delta: true})
		got, err := store.SourceRecords(t.Context(), cumulative, sourceA)
		require.NoError(t, err)
		require.True(t, got.Latest.Available, "%+v", got.Latest)
		require.True(t, got.Latest.RawPresent)
		require.True(t, *got.Latest.DeltaKnown)
	})
}

func editCells(t *testing.T, db *sql.DB, effect string, row dbtarget.EncodedRow, edits map[string]any) {
	t.Helper()
	for column, value := range edits {
		for i := range row.Cells {
			if row.Cells[i].Column == column {
				row.Cells[i].Value = value
			}
		}
	}
	payload, _, err := encodeRowPayload(row)
	require.NoError(t, err)
	rewritePayload(t, db, effect, string(payload))
}
