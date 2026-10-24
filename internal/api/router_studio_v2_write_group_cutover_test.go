package api

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var cutoverKinds = []string{"sqlite", "postgres", "timescale"}

// TestWideToLongCutover moves the only source of a recording wide group to a
// long profile group at an exact, unaligned instant that lies inside both a wide
// and a long window. Through the production routes and pipeline it requires:
// everything observed before the instant belongs to the wide group, everything
// at or after it to the long group; the windows the instant cuts are marked, not
// passed off as complete; rows already waiting for the old destination still
// arrive there; the wide history, layout and ownership do not change; and a
// restart in the middle of a window changes nothing.
func TestWideToLongCutover(t *testing.T) {
	for _, kind := range cutoverKinds {
		t.Run(kind, func(t *testing.T) {
			r := newCutoverRig(t, kind)
			wideID := r.wide["id"].(string)

			// History written by the wide group before the cutover.
			r.accept(r.wideEffective.Add(2*time.Second), 1.5)
			r.advance(r.wideEffective.Add(16 * time.Second))
			r.waitFor(func() bool { return r.wideRows() == 1 }, "the wide history row is delivered")
			receiptsBefore := r.receipts()
			wideBefore := getGroup(t, r.f, wideID)

			// Deliveries wait, so the old group's rows are still pending at the cutover.
			r.restart(time.Hour)
			boundary := r.at(22*time.Second + 500*time.Millisecond)
			cut := r.cutover(boundary)
			require.Equal(t, boundary.Format(time.RFC3339Nano), cut["boundary"], "the stored instant is the reviewed one, to the microsecond")
			r.reconcile()

			r.accept(r.at(5*time.Second), 5.0)             // wide window [0,15): complete
			r.restart(time.Hour)                           // a restart in the middle of an open window; nothing is waiting yet
			r.accept(boundary.Add(-time.Millisecond), 6.0) // the wide group's last observation
			r.accept(boundary, 7.0)                        // exactly the boundary: the long group's first
			r.accept(r.at(45*time.Second), 8.0)            // long window [40,60): complete
			r.advance(r.at(61 * time.Second))
			r.waitFor(func() bool { return r.outbox()["pending"] == 4 }, "two wide rows and two long records wait")
			frozen := r.outboxDigests()

			r.restart(cutoverWorkerEvery)
			r.waitFor(func() bool { return r.outbox()["sql_committed"] == 5 }, "the waiting rows are delivered")

			wide := r.wideRowsInOrder()
			require.Len(t, wide, 3, "the old group kept its history, its last complete window and the window the boundary cut")
			require.InDelta(t, 5.0, *wide[1].First, 0)
			cutWindow := wide[2]
			require.InDelta(t, 6.0, *cutWindow.First, 0, "the cut window keeps what was accepted before the boundary, and nothing at or after it")
			require.True(t, r.at(15*time.Second).Equal(cutWindow.Start), "it still starts at its own window start: %s", cutWindow.Start)
			evidence := windowOfWide(t, cutWindow.Provenance)
			require.Equal(t, "end", evidence["truncated"], "the wide row says it is cut, and where")
			require.Equal(t, boundary.Format(cutoverTimeText), evidence["end"])
			require.Nil(t, windowOfWide(t, wide[1].Provenance), "a complete window carries no truncation evidence")
			for _, row := range wide {
				require.NotNil(t, row.First)
				require.NotContains(t, []float64{7, 8}, *row.First, "nothing at or after the boundary reaches the old group")
			}
			records := r.longRecords()
			require.Len(t, records, 2)
			first, second := records[0], records[1]
			require.Equal(t, "partial", first.Status, "the window the boundary cut is partial, not ok")
			require.True(t, boundary.Equal(first.Start), "it starts at the real boundary: %s", first.Start)
			require.True(t, r.at(40*time.Second).Equal(first.End), "and ends at the aligned window end: %s", first.End)
			require.True(t, boundary.Equal(first.Observed), "an observation at the boundary belongs to the new owner")
			require.InDelta(t, 7.0, first.Value, 0)
			window := windowOf(t, first.Provenance)
			require.Equal(t, "start", window["truncated"])
			require.Equal(t, r.at(20*time.Second).Format(cutoverTimeText), window["nominal_start"])
			require.True(t, r.at(40*time.Second).Equal(second.Start))
			require.Equal(t, [2]any{"ok", 8.0}, [2]any{second.Status, second.Value})
			require.Nil(t, windowOf(t, second.Provenance), "a complete window carries no truncation evidence")

			require.Equal(t, frozen, r.outboxDigests(), "no waiting row was rewritten")
			for key, value := range receiptsBefore {
				require.Equal(t, value, r.receipts()[key], "the old history's receipt is untouched")
			}
			wideAfter := getGroup(t, r.f, wideID)
			for _, key := range []string{"revision", "applied_revision", "layout_version", "row_policy", "write_policy", "destination", "members", "status"} {
				require.Equal(t, wideBefore[key], wideAfter[key], "the old group's saved %s did not change", key)
			}
			view := r.delivery(wideID)
			require.EqualValues(t, 0, view["skipped_buckets"], "the cut window was written, not skipped")
			var truncated int
			for _, issue := range view["recent_bucket_issues"].([]any) {
				if issue.(map[string]any)["kind"] == "truncated" {
					truncated++
				}
			}
			require.Equal(t, 1, truncated, "the delivery view says which bucket the boundary cut")
		})
	}
}

// TestProfileCutoverFailure: whatever goes wrong, the group that holds the
// sources stays active and unchanged, nothing is written twice, and a later
// correct attempt still works.
func TestProfileCutoverFailure(t *testing.T) {
	r := newCutoverRig(t, "sqlite")
	wideID, longID := r.wide["id"].(string), r.long["id"].(string)
	boundary := r.at(22*time.Second + 500*time.Millisecond)

	untouched := func(t *testing.T) {
		t.Helper()
		require.Zero(t, r.localCount(`SELECT COUNT(*) FROM write_group_cutovers`))
		require.Zero(t, r.localCount(`SELECT COUNT(*) FROM write_group_revision_ends`))
		require.Zero(t, r.localCount(`SELECT COUNT(*) FROM write_group_versions WHERE group_id = ?`, longID), "the long group was never applied")
		require.Equal(t, "", getGroup(t, r.f, longID)["applied_revision"])
		snapshot, err := r.f.groups.ResolveAppliedAt(t.Context(), wideID, time.Now().Add(time.Hour))
		require.NoError(t, err)
		require.True(t, snapshot.EndsAt.IsZero(), "the old group still owns its sources")
		eligibility, err := r.f.groups.IntakeEligibility(t.Context(), wideID)
		require.NoError(t, err)
		require.NotEqual(t, "cutover-ended", eligibility.Reason)
	}

	t.Run("an ordinary apply of the long group is refused while the wide group owns the source", func(t *testing.T) {
		code, body := r.post(writeGroupsPath+"/"+longID+"/apply", scopeBody(t, r.f, longID))
		require.Equal(t, http.StatusConflict, code)
		require.Equal(t, "cutover_required", errorOf(t, body)["code"])
		untouched(t)
	})

	t.Run("a boundary that has passed is refused with its reason", func(t *testing.T) {
		past := time.Now().UTC().Add(-time.Minute)
		review := r.review(past)
		require.Equal(t, false, review["applicable"])
		require.Contains(t, review["issues"], "boundary_not_future")
		code, body := r.post(writeGroupsPath+"/"+wideID+"/cutover", r.cutoverBody(past, map[string]any{"review_digest": review["digest"]}))
		require.Equal(t, http.StatusUnprocessableEntity, code)
		require.Equal(t, "cutover_not_applicable", errorOf(t, body)["code"])
		untouched(t)
	})

	t.Run("a review of another instant does not apply", func(t *testing.T) {
		review := r.review(boundary)
		code, body := r.post(writeGroupsPath+"/"+wideID+"/cutover",
			r.cutoverBody(boundary.Add(time.Second), map[string]any{"review_digest": review["digest"]}))
		require.Equal(t, http.StatusConflict, code)
		require.Equal(t, "revision_mismatch", errorOf(t, body)["code"])
		untouched(t)
	})

	t.Run("a stale source revision does not apply", func(t *testing.T) {
		review := r.review(boundary)
		code, body := r.post(writeGroupsPath+"/"+wideID+"/cutover",
			r.cutoverBody(boundary, map[string]any{"review_digest": review["digest"], "expected_source_revision": "an-older-revision"}))
		require.Equal(t, http.StatusConflict, code)
		require.Equal(t, "revision_mismatch", errorOf(t, body)["code"])
		untouched(t)
	})

	t.Run("a failure while recording the cutover rolls the whole transaction back", func(t *testing.T) {
		_, err := r.f.db.ExecContext(t.Context(), `CREATE TRIGGER fail_cutover BEFORE INSERT ON write_group_cutovers BEGIN SELECT RAISE(ABORT, 'injected'); END`)
		require.NoError(t, err)
		review := r.review(boundary)
		code, _ := r.post(writeGroupsPath+"/"+wideID+"/cutover", r.cutoverBody(boundary, map[string]any{"review_digest": review["digest"]}))
		require.GreaterOrEqual(t, code, http.StatusInternalServerError)
		untouched(t)

		// The old group keeps recording through the failed attempt.
		r.accept(r.wideEffective.Add(2*time.Second), 1.5)
		r.advance(r.wideEffective.Add(16 * time.Second))
		r.waitFor(func() bool { return r.wideRows() == 1 }, "the wide group still delivers")
		require.Empty(t, r.longRecords())

		_, err = r.f.db.ExecContext(t.Context(), `DROP TRIGGER fail_cutover`)
		require.NoError(t, err)
		cut := r.cutover(boundary)
		require.Equal(t, boundary.Format(time.RFC3339Nano), cut["boundary"], "the same reviewed values still apply once the fault is gone")
	})
}

// TestProfileRollbackDrain returns the sources to the wide group at a second,
// explicit boundary. Records the long group already accepted still drain to the
// long table, nothing is written twice, and no earlier version, frozen row,
// table or receipt is rewritten or removed.
func TestProfileRollbackDrain(t *testing.T) {
	for _, kind := range cutoverKinds {
		t.Run(kind, func(t *testing.T) {
			r := newCutoverRig(t, kind)
			wideID := r.wide["id"].(string)
			forward := r.at(22*time.Second + 500*time.Millisecond)
			rollback := r.at(82*time.Second + 500*time.Millisecond)
			ended := getGroup(t, r.f, wideID)
			endedRevision := ended["revision"]
			endedVersion := r.versionPayload(wideID, endedRevision.(string))

			r.restart(time.Hour)
			cut := r.cutover(forward)
			r.reconcile()
			r.accept(forward, 7.0)
			r.accept(r.at(45*time.Second), 8.0)
			r.advance(r.at(61 * time.Second))
			r.waitFor(func() bool { return r.outbox()["pending"] == 2 }, "two long records wait")
			frozen := r.outboxDigests()

			// An unreviewed or repeated rollback is not applicable.
			code, body := r.post(writeGroupsPath+"/"+wideID+"/cutover", r.rollbackBody(cut["id"].(string), rollback, map[string]any{"review_digest": "not-the-digest"}))
			require.Equal(t, http.StatusConflict, code, "%v", body)

			review := r.rollbackReview(cut["id"].(string), rollback)
			require.Equal(t, true, review["applicable"], "%v", review["issues"])
			code, body = r.post(writeGroupsPath+"/"+wideID+"/cutover", r.rollbackBody(cut["id"].(string), rollback, map[string]any{"review_digest": review["digest"]}))
			require.Equal(t, http.StatusOK, code, "%v", body)
			require.Equal(t, rollback.Format(time.RFC3339Nano), body["data"].(map[string]any)["boundary"])
			r.reconcile()

			r.accept(r.at(70*time.Second), 9.0)             // long window [60,80): complete
			r.accept(rollback.Add(-time.Millisecond), 10.0) // the long group's last observation
			r.accept(rollback, 11.0)                        // exactly the second boundary: the wide group's again
			r.accept(r.at(100*time.Second), 12.0)           // wide window [90,105): complete
			r.advance(r.at(106 * time.Second))
			r.waitFor(func() bool { return r.outbox()["pending"] == 6 }, "four long records and two wide rows wait")

			r.restart(cutoverWorkerEvery)
			r.waitFor(func() bool { return r.outbox()["sql_committed"] == 6 }, "everything drains to the destination that accepted it")

			records := r.longRecords()
			require.Len(t, records, 4)
			var values []float64
			for _, rec := range records {
				values = append(values, rec.Value)
			}
			require.Equal(t, []float64{7, 8, 9, 10}, values, "the long table holds exactly what it owned, nothing at or after the second boundary")
			require.Equal(t, "partial", records[3].Status)
			require.True(t, rollback.Equal(records[3].End), "the long window the rollback cut ends at the boundary: %s", records[3].End)
			require.Equal(t, "end", windowOf(t, records[3].Provenance)["truncated"])
			wide := r.wideRowsInOrder()
			require.Len(t, wide, 2, "the restored group wrote the window the rollback cut and its next complete window")
			require.InDelta(t, 11.0, *wide[0].First, 0)
			require.True(t, rollback.Equal(wide[0].Start), "a window cut at its start begins at the real boundary: %s", wide[0].Start)
			require.Equal(t, "start", windowOfWide(t, wide[0].Provenance)["truncated"])
			require.InDelta(t, 12.0, *wide[1].First, 0)
			for key, digest := range frozen {
				require.Equal(t, digest, r.outboxDigests()[key], "a payload frozen before the rollback is not rewritten")
			}

			restored := getGroup(t, r.f, wideID)
			require.NotEqual(t, endedRevision, restored["revision"], "the restored group is a new revision")
			require.Equal(t, restored["revision"], restored["applied_revision"])
			for _, key := range []string{"layout_version", "row_policy", "write_policy", "members", "destination"} {
				require.Equal(t, ended[key], restored[key], "the restored group has the saved %s", key)
			}
			require.Equal(t, endedVersion, r.versionPayload(wideID, endedRevision.(string)), "the ended version is not rewritten")
			require.True(t, r.tableExists(r.wideTable()) && r.tableExists(r.longTable()), "no table was removed")

			code, _ = r.post(writeGroupsPath+"/"+wideID+"/cutover-review", r.rollbackBody(cut["id"].(string), rollback.Add(time.Hour), nil))
			require.Equal(t, http.StatusOK, code)
			again := r.rollbackReview(cut["id"].(string), rollback.Add(time.Hour))
			require.Equal(t, false, again["applicable"], "a cutover is rolled back once")
			require.True(t, strings.Contains(strings.Join(toStrings(again["issues"]), ","), "cutover_not_rollbackable"))
		})
	}
}

// TestWideToCumulativeCutover moves an exact energy counter from a wide group to
// a cumulative profile. The counter's history belongs to the wide group: the new
// profile starts a fresh baseline at the cut window, whose first record has raw
// value and no difference, and never subtracts a reading the old group took. The
// next window then reports the exact difference.
func TestWideToCumulativeCutover(t *testing.T) {
	const first = uint64(9007199254740992) // 2^53: a float64 would lose the next step
	for _, kind := range cutoverKinds {
		t.Run(kind, func(t *testing.T) {
			r := newCutoverRigWith(t, kind, cutoverOptions{cumulative: true})
			boundary := r.at(22*time.Second + 500*time.Millisecond)
			r.restart(time.Hour)
			r.cutover(boundary)
			r.reconcile()

			r.accept(boundary.Add(-time.Second), first-100) // the wide group's reading: not the baseline of the new profile
			r.accept(boundary, first)                       // the new profile's first reading, inside its one-minute window [0,60)
			r.accept(r.at(75*time.Second), first+1)         // window [60,120)
			r.advance(r.at(121 * time.Second))
			r.restart(cutoverWorkerEvery)
			r.waitFor(func() bool {
				var n int
				return r.dest.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM `+r.quoted(r.longTable())).Scan(&n) == nil && n == 2
			}, "both cumulative records are delivered")

			rows, err := r.dest.QueryContext(t.Context(), `SELECT status, CAST(value_decimal AS TEXT), COALESCE(CAST(delta_value AS TEXT), 'NULL'), CAST(provenance AS TEXT)
				FROM `+r.quoted(r.longTable())+` ORDER BY bucket_start`)
			require.NoError(t, err)
			defer rows.Close()
			type record struct{ status, raw, delta, provenance string }
			var got []record
			for rows.Next() {
				var rec record
				require.NoError(t, rows.Scan(&rec.status, &rec.raw, &rec.delta, &rec.provenance))
				got = append(got, rec)
			}
			require.NoError(t, rows.Err())
			require.Len(t, got, 2)
			require.Equal(t, [3]string{"baseline", "9007199254740992", "NULL"}, [3]string{got[0].status, got[0].raw, got[0].delta},
				"the first record of the profile is a baseline with no difference, even though its window was cut")
			require.Equal(t, "start", windowOf(t, got[0].provenance)["truncated"], "the cut shows in the evidence")
			require.Equal(t, [3]string{"ok", "9007199254740993", "1"}, [3]string{got[1].status, got[1].raw, got[1].delta},
				"the next window is the exact difference from the profile's own baseline")
		})
	}
}
