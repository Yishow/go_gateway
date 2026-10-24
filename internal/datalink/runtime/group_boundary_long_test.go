package runtime

import (
	"encoding/json"
	"math/rand/v2"
	"slices"
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

const maxUint64Text = "18446744073709551615"

func TestInstantaneousProfileWindowEdges(t *testing.T) {
	f := newLongFixture(t)
	boundary := f.build(t)
	f.accept(t, boundary,
		f.envelope(0, "a-start", at(0), 1.0),
		f.envelope(0, "a-last", at(15*time.Second-time.Microsecond), 2.0),
		f.envelope(0, "a-next-window", at(15*time.Second), 3.0),
		f.envelope(1, "b-int", at(5*time.Second), int64(9007199254740993)),
		f.envelope(2, "run", at(7*time.Second), true),
		f.envelope(3, "batch", at(9*time.Second), "run-42"),
		f.envelope(4, "count", at(11*time.Second), uint64(18446744073709551615)),
	)

	require.NoError(t, boundary.Tick(t.Context(), at(14*time.Second)))
	require.Empty(t, f.sink.rows, "the window is still open")
	require.NoError(t, boundary.Tick(t.Context(), at(15*time.Second)))
	first := longRows(t, f.sink)
	require.Len(t, first, 5, "every enabled member has one row for the window")

	a := first[longRowKey{"tag-A", bucketText(0)}]
	require.Equal(t, bucketText(15*time.Second), a["bucket_end"])
	require.Equal(t, 2.0, a["value_float"], "the newest observation inside [start,end) wins")
	require.Equal(t, at(15*time.Second-time.Microsecond).Format(longTimeLayout), a["observed_at"], "the real observation time is kept")
	require.Equal(t, at(15*time.Second-time.Microsecond+40*time.Millisecond).Format(longTimeLayout), a["received_at"])
	require.Equal(t, "ok", a["status"])
	require.Equal(t, "instantaneous", a["profile_kind"])
	require.Equal(t, "float64", a["value_type"])
	require.Equal(t, []string{"value_float"}, valueSlots(a))
	require.Equal(t, "circuit-a", a["circuit_id"])
	require.Nil(t, a["delta_value"])
	require.Nil(t, a["delta_from"])
	require.Nil(t, a["delta_to"])

	require.Equal(t, int64(9007199254740993), first[longRowKey{"tag-B", bucketText(0)}]["value_int"])
	require.Equal(t, int64(1), first[longRowKey{"tag-run", bucketText(0)}]["value_bool"])
	require.Equal(t, "run-42", first[longRowKey{"tag-batch", bucketText(0)}]["value_text"])
	count := first[longRowKey{"tag-count", bucketText(0)}]
	require.Equal(t, maxUint64Text, count["value_decimal"], "uint64 stays exact and never passes through float")
	require.Equal(t, []string{"value_decimal"}, valueSlots(count))
	require.Equal(t, "uint64", count["value_type"])
	require.Equal(t, "measurement-count", count["measurement_id"])

	// The sample observed exactly at the end belongs to the next window.
	f.sink.rows = nil
	require.NoError(t, boundary.Tick(t.Context(), at(30*time.Second)))
	second := longRows(t, f.sink)
	require.Len(t, second, 5)
	next := second[longRowKey{"tag-A", bucketText(15 * time.Second)}]
	require.Equal(t, 3.0, next["value_float"])
	require.Equal(t, at(15*time.Second).Format(longTimeLayout), next["observed_at"])
	require.Equal(t, "missing", second[longRowKey{"tag-B", bucketText(15 * time.Second)}]["status"], "no value is carried across windows")
}

func TestInstantaneousProfileEligibility(t *testing.T) {
	f := newLongFixture(t)
	f.members[0].MaxAgeSeconds = ptr(5)
	f.config.Group.Members = f.members
	boundary := f.build(t)

	good := f.envelope(0, "a-good", at(3*time.Second), 1.5) // 12s before the end: stale for a 5s limit
	bad := f.envelope(1, "b-bad", at(10*time.Second), int64(7))
	bad.Quality, bad.QualityReason, bad.Value = schema.QualityBad, "read-timeout", nil
	goodOlder := f.envelope(2, "run-good", at(2*time.Second), true)
	badNewer := f.envelope(2, "run-bad", at(12*time.Second), nil)
	badNewer.Quality, badNewer.QualityReason = schema.QualityBad, "device-fault"
	mismatch := f.envelope(3, "batch-wrong-type", at(8*time.Second), int64(5)) // text member, integer payload
	f.accept(t, boundary, good, bad, goodOlder, badNewer, mismatch)
	require.NoError(t, boundary.Tick(t.Context(), at(15*time.Second)))
	rows := longRows(t, f.sink)

	stale := rows[longRowKey{"tag-A", bucketText(0)}]
	require.Equal(t, "stale", stale["status"])
	require.Empty(t, valueSlots(stale), "a stale reading is NULL, never reused")
	require.Contains(t, provenanceText(t, stale), "max-age-exceeded")
	require.Equal(t, at(3*time.Second).Format(longTimeLayout), stale["observed_at"], "the real observation time of the stale reading is kept")

	failed := rows[longRowKey{"tag-B", bucketText(0)}]
	require.Equal(t, "bad", failed["status"])
	require.Equal(t, "bad", failed["quality"])
	require.Empty(t, valueSlots(failed), "a failed read is never written as zero")
	require.Contains(t, provenanceText(t, failed), "read-timeout")

	older := rows[longRowKey{"tag-run", bucketText(0)}]
	require.Equal(t, "ok", older["status"], "the newest eligible sample is used even when a newer one is bad")
	require.Equal(t, int64(1), older["value_bool"])

	invalid := rows[longRowKey{"tag-batch", bucketText(0)}]
	require.Equal(t, "invalid", invalid["status"])
	require.Empty(t, valueSlots(invalid))
	require.Contains(t, provenanceText(t, invalid), "type-mismatch")

	untouched := rows[longRowKey{"tag-count", bucketText(0)}]
	require.Equal(t, "missing", untouched["status"])
	require.Nil(t, untouched["observed_at"], "no sample means no invented time")
	require.Nil(t, untouched["received_at"])
	require.Equal(t, "uint64", untouched["value_type"], "the declared source type survives a missing window")
}

func TestLongProfileIndependentMissing(t *testing.T) {
	f := newLongFixture(t)
	boundary := f.build(t)
	// Member B and the required counter have nothing in the first window.
	f.accept(t, boundary,
		f.envelope(0, "a1", at(2*time.Second), 10.5),
		f.envelope(2, "run1", at(4*time.Second), false),
	)
	require.NoError(t, boundary.Tick(t.Context(), at(15*time.Second)))
	require.NoError(t, boundary.Tick(t.Context(), at(30*time.Second))) // a window without any sample
	require.Len(t, f.sink.rows, 10, "required or missing members never suppress each other, and silent windows still emit rows")
	require.Empty(t, f.sink.outcomes, "no member is skipped or reported as no-data")

	rows := longRows(t, f.sink)
	require.Equal(t, 10.5, rows[longRowKey{"tag-A", bucketText(0)}]["value_float"])
	require.Equal(t, int64(0), rows[longRowKey{"tag-run", bucketText(0)}]["value_bool"], "false is a value, not a gap")
	for _, tag := range []string{"tag-B", "tag-batch", "tag-count"} {
		missing := rows[longRowKey{tag, bucketText(0)}]
		require.Equal(t, "missing", missing["status"], tag)
		require.Empty(t, valueSlots(missing), tag)
		require.Contains(t, provenanceText(t, missing), "no-samples", tag)
	}
	for _, tag := range []string{"tag-A", "tag-B", "tag-run", "tag-batch", "tag-count"} {
		require.Equal(t, "missing", rows[longRowKey{tag, bucketText(15 * time.Second)}]["status"], tag)
	}

	// Members of one window never share an identity, so no row overwrites another.
	recordIDs, effectKeys := map[string]struct{}{}, map[string]struct{}{}
	for _, row := range f.sink.rows {
		require.NotEmpty(t, row.Encoded.RecordID)
		require.NotEmpty(t, row.Encoded.EffectKey)
		recordIDs[row.Encoded.RecordID] = struct{}{}
		effectKeys[row.Encoded.EffectKey] = struct{}{}
		require.Equal(t, row.Encoded.RecordID, cellsOf(row)["record_id"])
	}
	require.Len(t, recordIDs, 10)
	require.Len(t, effectKeys, 10)
}

func TestLongProfileRowsAreIndependentOfArrivalOrder(t *testing.T) {
	build := func(order []int) []GroupRow {
		f := newLongFixture(t)
		boundary := f.build(t)
		envelopes := []measurement.SampleEnvelope{
			f.envelope(0, "a-1", at(2*time.Second), 1.0),
			f.envelope(0, "a-2", at(9*time.Second), 2.0),
			f.envelope(0, "a-tie-1", at(12*time.Second), 3.0),
			f.envelope(0, "a-tie-2", at(12*time.Second), 4.0), // same instant: the greater sample ID wins
			f.envelope(1, "b-1", at(3*time.Second), int64(5)),
			f.envelope(4, "c-1", at(14*time.Second), uint64(7)),
			f.envelope(3, "d-1", at(1*time.Second), "text"),
		}
		for _, i := range order {
			require.NoError(t, boundary.AcceptSample(t.Context(), envelopes[i]))
		}
		// A resend of an accepted sample is a no-op.
		require.NoError(t, boundary.AcceptSample(t.Context(), envelopes[order[0]]))
		require.NoError(t, boundary.Tick(t.Context(), at(15*time.Second)))
		rows := slices.Clone(f.sink.rows)
		slices.SortFunc(rows, func(a, b GroupRow) int { return compareStrings(a.Encoded.RecordID, b.Encoded.RecordID) })
		return rows
	}
	base := build([]int{0, 1, 2, 3, 4, 5, 6})
	require.Len(t, base, 5)
	forwardJSON, err := json.Marshal(base)
	require.NoError(t, err)
	for _, cells := range base {
		if cellsOf(cells)["tag_id"] == "tag-A" {
			require.Equal(t, 4.0, cellsOf(cells)["value_float"], "the greater sample ID breaks a same-instant tie")
			require.Contains(t, cellsOf(cells)["provenance"], "a-tie-2")
		}
	}
	rng := rand.New(rand.NewPCG(7, 11))
	orders := [][]int{{6, 5, 4, 3, 2, 1, 0}, rng.Perm(7), rng.Perm(7), rng.Perm(7)}
	for _, order := range orders {
		rows := build(order)
		encoded, err := json.Marshal(rows)
		require.NoError(t, err)
		require.JSONEq(t, string(forwardJSON), string(encoded), "arrival order %v must not change any row or identity", order)
	}
}

func compareStrings(a, b string) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

func TestProfileLateSample(t *testing.T) {
	f := newLongFixture(t)
	f.config.Group.RowPolicy.AllowedLatenessSeconds = 5
	boundary := f.build(t)
	f.accept(t, boundary, f.envelope(0, "a-on-time", at(10*time.Second), 1.0))

	// Inside the allowed lateness the window is still open.
	require.NoError(t, boundary.Tick(t.Context(), at(19*time.Second)))
	require.Empty(t, f.sink.rows)
	f.accept(t, boundary, f.envelope(0, "a-in-grace", at(14*time.Second), 2.0))

	require.NoError(t, boundary.Tick(t.Context(), at(20*time.Second)))
	require.Len(t, f.sink.rows, 5)
	closed := longRows(t, f.sink)
	require.Equal(t, 2.0, closed[longRowKey{"tag-A", bucketText(0)}]["value_float"])
	beforeRows, err := json.Marshal(f.sink.rows)
	require.NoError(t, err)
	require.Zero(t, boundary.assembler.Late())

	// After the window closed, a late sample is counted, never merged.
	lateSample := f.envelope(0, "a-late", at(13*time.Second), 9.0)
	err = boundary.AcceptSample(t.Context(), lateSample)
	var refused *GroupSampleError
	require.ErrorAs(t, err, &refused)
	require.Equal(t, snapshot.OfferLate, refused.Outcome)
	require.Equal(t, uint64(1), boundary.assembler.Late())

	require.NoError(t, boundary.Tick(t.Context(), at(60*time.Second)))
	for key, cells := range longRows(t, &capturedGroupSink{rows: f.sink.rows[:5]}) {
		require.Equal(t, closed[key], cells, "a finalized row is never rewritten")
	}
	afterFirst, err := json.Marshal(f.sink.rows[:5])
	require.NoError(t, err)
	require.JSONEq(t, string(beforeRows), string(afterFirst))
	require.Len(t, f.sink.rows, 5+2*5, "windows [15,30) and [30,45) close by 60s with a 5s grace; the closed window is not emitted again")
	rewindowed := longRows(t, &capturedGroupSink{rows: f.sink.rows[5:]})
	for key := range rewindowed {
		require.NotEqual(t, bucketText(0), key.bucket)
	}
}

func TestProfileSampleGates(t *testing.T) {
	f := newLongFixture(t)
	boundary := f.build(t)

	stale := f.envelope(0, "a-old-revision", at(3*time.Second), 1.0)
	stale.MappingRevision = "map-A-old"
	var refused *GroupSampleError
	require.ErrorAs(t, boundary.AcceptSample(t.Context(), stale), &refused)
	require.Equal(t, snapshot.ReasonRevisionMismatch, refused.Reason)

	wrongSource := f.envelope(0, "a-old-source", at(3*time.Second), 1.0)
	wrongSource.SourceRevision = "src-A-old"
	require.ErrorAs(t, boundary.AcceptSample(t.Context(), wrongSource), &refused)
	require.Equal(t, snapshot.ReasonRevisionMismatch, refused.Reason)

	future := f.envelope(0, "a-future", f.clock.Add(time.Hour), 1.0)
	require.ErrorAs(t, boundary.AcceptSample(t.Context(), future), &refused)
	require.Equal(t, snapshot.ReasonFutureSkew, refused.Reason)

	foreign := f.envelope(0, "a-foreign", at(3*time.Second), 1.0)
	foreign.WorkspaceID = "workspace-B"
	require.ErrorAs(t, boundary.AcceptSample(t.Context(), foreign), &refused)

	f.accept(t, boundary, f.envelope(0, "a-ok", at(3*time.Second), 1.0))
	conflict := f.envelope(0, "a-ok", at(3*time.Second), 2.0)
	require.ErrorAs(t, boundary.AcceptSample(t.Context(), conflict), &refused)
	require.Equal(t, snapshot.ReasonIdentityConflict, refused.Reason)

	require.NoError(t, boundary.Tick(t.Context(), at(15*time.Second)))
	rows := longRows(t, f.sink)
	require.Equal(t, 1.0, rows[longRowKey{"tag-A", bucketText(0)}]["value_float"], "refused samples never reach a row")
	require.Equal(t, "ok", rows[longRowKey{"tag-A", bucketText(0)}]["status"])
}
