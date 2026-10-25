package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// transitionEvidence is one member's evidence in the provenance of a wide row.
type transitionEvidence struct {
	Member  string                 `json:"member"`
	Status  string                 `json:"status"`
	Window  map[string]any         `json:"window"`
	Mean    *storedMean            `json:"mean"`
	Counter *counterProvenanceView `json:"counter"`
}

// transitionRow is one row of the new wide profile table as the destination stores it.
type transitionRow struct {
	Start    time.Time
	Value    sql.NullString
	Evidence []transitionEvidence
}

// targetRows reads the new wide profile table in window order.
func (r *cutoverRig) targetRows() []transitionRow {
	r.t.Helper()
	column := r.long["members"].([]any)[0].(map[string]any)["target_column"].(string)
	rows, err := r.dest.QueryContext(r.t.Context(), `SELECT CAST(bucket_start AS TEXT), CAST("`+column+`" AS TEXT), CAST(provenance AS TEXT)
		FROM `+r.quoted(r.longTable())+` ORDER BY bucket_start`)
	require.NoError(r.t, err)
	defer rows.Close()
	var out []transitionRow
	for rows.Next() {
		var row transitionRow
		var start, provenance string
		require.NoError(r.t, rows.Scan(&start, &row.Value, &provenance))
		row.Start = storedInstant(r.t, start)
		require.NoError(r.t, json.Unmarshal([]byte(provenance), &row.Evidence), provenance)
		out = append(out, row)
	}
	require.NoError(r.t, rows.Err())
	return out
}

// frozenOutbox is what a closed window froze: the group revision, its destination and the payload with its digest.
type frozenOutbox struct {
	GroupID, GroupRevision, ConnectorID, ConnectorRevision, TableSchema, TableName, PayloadDigest, Payload string
}

// frozenOutbox reads those facts of every outbox row, keyed by effect.
func (r *cutoverRig) frozenOutbox() map[string]frozenOutbox {
	r.t.Helper()
	rows, err := r.f.db.QueryContext(r.t.Context(), `SELECT effect_key, group_id, group_revision, connector_id, connector_revision, table_schema, table_name, payload_digest, payload FROM wg_delivery_outbox`)
	require.NoError(r.t, err)
	defer rows.Close()
	out := map[string]frozenOutbox{}
	for rows.Next() {
		var key string
		var row frozenOutbox
		require.NoError(r.t, rows.Scan(&key, &row.GroupID, &row.GroupRevision, &row.ConnectorID, &row.ConnectorRevision, &row.TableSchema, &row.TableName, &row.PayloadDigest, &row.Payload))
		out[key] = row
	}
	require.NoError(r.t, rows.Err())
	return out
}

func floatOf(t *testing.T, value sql.NullString) float64 {
	t.Helper()
	require.True(t, value.Valid, "the member column holds a value")
	parsed, err := strconv.ParseFloat(value.String, 64)
	require.NoError(t, err, value.String)
	return parsed
}

// TestReviewedWideMeanTransition moves the source of a recording legacy
// 60 second wide group to a new managed_wide_profile_v1 mean group through the
// existing reviewed cutover: the production routes, pipeline and a disposable
// SQLite destination, with the test's own UTC clock. The legacy wide table is
// never rewritten; its history, its pending rows and its saved revision stay.
func TestReviewedWideMeanTransition(t *testing.T) {
	t.Run("the cut windows, the first mean window and the frozen pending rows", func(t *testing.T) {
		r := newCutoverRigWith(t, "sqlite", cutoverOptions{wideSeconds: 60, wideTarget: true})
		wideID := r.wide["id"].(string)
		wideBefore := getGroup(t, r.f, wideID)

		// Deliveries wait. The legacy group closes its first 60 second window before the transition, so a
		// closed row is already pending, frozen to the old connector, table and payload, when the transition is reviewed.
		r.restart(time.Hour)
		r.accept(r.wideEffective.Add(2*time.Second), 1.5)
		r.advance(r.wideEffective.Add(61 * time.Second))
		r.waitFor(func() bool { return r.outbox()["pending"] == 1 }, "the legacy history row is closed and waiting")
		require.Zero(t, r.wideRows(), "it has not reached the old table yet")
		waiting := r.frozenOutbox()
		require.Len(t, waiting, 1)
		for _, row := range waiting {
			require.Equal(t, [4]string{wideID, wideBefore["revision"].(string), "connector-A", "connector-1"}, [4]string{row.GroupID, row.GroupRevision, row.ConnectorID, row.ConnectorRevision})
			require.Equal(t, r.wideTable(), row.TableName, "frozen to the legacy table")
			require.NotEmpty(t, row.PayloadDigest)
		}
		boundary := r.at(7 * time.Second) // inside the nominal [0,15) window of the new group and [0,60) of the old one
		review := r.review(boundary)
		require.Equal(t, true, review["applicable"], "the reviewed cutover accepts the known wide mean target: %v", review["issues"])
		cut := r.cutover(boundary)
		require.Equal(t, boundary.Format(time.RFC3339Nano), cut["boundary"])
		r.reconcile()

		r.accept(r.at(3*time.Second), 3.0) // the legacy group's last observation: its window [0,60) is cut at 7
		r.accept(boundary, 4.0)            // exactly the boundary: the new group's first
		r.accept(r.at(9*time.Second), 8.0) // new window [0,15), actual [7,15)
		r.accept(r.at(13*time.Second), 12.0)
		r.accept(r.at(17*time.Second), 20.0) // new window [15,30): complete
		r.accept(r.at(25*time.Second), 30.0)
		r.advance(r.at(61 * time.Second)) // past the legacy window's nominal end, which is when its cut row closes
		r.waitFor(func() bool { return r.outbox()["pending"] == 6 }, "the history row, the cut legacy row and the four mean rows wait")
		frozen := r.outboxDigests()

		r.restart(cutoverWorkerEvery)
		r.waitFor(func() bool { return r.outbox()["sql_committed"] == 6 }, "the waiting rows are delivered to the destination that accepted them")

		for key, before := range waiting {
			require.Equal(t, before, r.frozenOutbox()[key], "the row that was pending before the transition keeps its revision, destination and payload")
		}
		legacy := r.wideRowsInOrder()
		require.Len(t, legacy, 2, "the legacy table holds its history and the window the boundary cut, nothing from the new owner")
		require.InDelta(t, 1.5, *legacy[0].First, 0, "the row pending before the transition arrived in the old table")
		require.InDelta(t, 3.0, *legacy[1].First, 0, "only what was accepted before the boundary")
		require.True(t, r.at(0).Equal(legacy[1].Start), "it still starts at its own window start: %s", legacy[1].Start)
		cutEvidence := windowOfWide(t, legacy[1].Provenance)
		require.Equal(t, "end", cutEvidence["truncated"])
		require.Equal(t, boundary.Format(cutoverTimeText), cutEvidence["end"])

		rows := r.targetRows()
		require.Len(t, rows, 4, "the new group wrote [7,15), [15,30) and the two silent windows after them")
		first, second := rows[0], rows[1]
		require.True(t, boundary.Equal(first.Start), "the cut window starts at the real boundary: %s", first.Start)
		require.InDelta(t, 8.0, floatOf(t, first.Value), 0, "the mean of 4, 8 and 12 over [7,15) only; the legacy 3 does not enter")
		firstWindow := first.Evidence[0].Window
		require.Equal(t, "start", firstWindow["truncated"])
		require.Equal(t, r.at(0).Format(cutoverTimeText), firstWindow["nominal_start"], "nominal and actual bounds are both kept")
		require.Equal(t, [2]any{"partial", int64(3)}, [2]any{first.Evidence[0].Mean.Coverage, first.Evidence[0].Mean.ValidCount}, "a truncated window is partial")
		require.True(t, r.at(15*time.Second).Equal(second.Start))
		require.InDelta(t, 25.0, floatOf(t, second.Value), 0)
		require.Nil(t, second.Evidence[0].Window, "a complete window carries no truncation evidence")
		require.Equal(t, [2]any{"unknown", int64(2)}, [2]any{second.Evidence[0].Mean.Coverage, second.Evidence[0].Mean.ValidCount}, "no whole-window proof is never complete")
		for _, silent := range rows[2:] {
			require.False(t, silent.Value.Valid, "a window without a reading is a NULL row, never a zero or an older value")
			require.Equal(t, "missing", silent.Evidence[0].Mean.Coverage)
		}

		require.Equal(t, frozen, r.outboxDigests(), "no pending payload or destination was rewritten")
		wideAfter := getGroup(t, r.f, wideID)
		for _, key := range []string{"revision", "applied_revision", "layout_version", "row_policy", "write_policy", "destination", "members", "status"} {
			require.Equal(t, wideBefore[key], wideAfter[key], "the legacy group's saved %s did not change", key)
		}
		require.Equal(t, float64(60), wideAfter["row_policy"].(map[string]any)["interval_seconds"])
		target := getGroup(t, r.f, r.long["id"].(string))
		require.Equal(t, target["revision"], target["applied_revision"], "the new group is applied by the reviewed cutover")
	})

	t.Run("a counter starts a fresh baseline with a NULL delta in the new profile", func(t *testing.T) {
		const first = uint64(9007199254740992) // 2^53: a float64 would lose the next step
		r := newCutoverRigWith(t, "sqlite", cutoverOptions{cumulative: true, wideTarget: true, wideSeconds: 60})
		boundary := r.at(7 * time.Second)
		r.restart(time.Hour)
		review := r.review(boundary)
		require.Equal(t, true, review["applicable"], "the known wide counter target is accepted: %v", review["issues"])
		r.cutover(boundary)
		r.reconcile()

		r.accept(boundary.Add(-time.Second), first-100) // the legacy group's reading: not the baseline of the new profile
		r.accept(boundary, first)                       // the new profile's first reading, in its minute window [0,60)
		r.accept(r.at(75*time.Second), first+1)         // window [60,120)
		r.advance(r.at(121 * time.Second))
		r.restart(cutoverWorkerEvery)
		r.waitFor(func() bool { return len(r.targetRows()) == 2 }, "both cumulative windows are delivered")

		rows := r.targetRows()
		baseline, next := rows[0].Evidence[0].Counter, rows[1].Evidence[0].Counter
		require.NotNil(t, baseline)
		require.NotNil(t, next)
		require.Equal(t, [2]string{"baseline", ""}, [2]string{baseline.Status, baseline.Delta}, "the first window has no difference, even though it was cut")
		require.Equal(t, "start", rows[0].Evidence[0].Window["truncated"])
		require.Equal(t, [2]string{"ok", "1"}, [2]string{next.Status, next.Delta}, "the next window is the exact difference from the profile's own baseline")
	})

	t.Run("a wide mean target that is not ready is refused for its own reason, not as an unknown layout", func(t *testing.T) {
		r := newCutoverRigWith(t, "sqlite", cutoverOptions{wideSeconds: 60, wideTarget: true, unconfirmedTarget: true})
		targetID := r.long["id"].(string)
		readiness := performJSONRequest(t, r.f.router, http.MethodGet, writeGroupsPath+"/"+targetID+"/readiness", nil)
		require.Equal(t, http.StatusOK, readiness.Code, readiness.Body.String())
		require.Equal(t, false, decodeJSONBody(t, readiness)["data"].(map[string]any)["ready"], "the precondition: its schema was never confirmed")

		review := r.review(r.at(7 * time.Second))
		require.Equal(t, false, review["applicable"])
		issues := toStrings(review["issues"])
		require.Contains(t, issues, "target_not_ready", "full readiness is still required of a known wide target")
		require.NotContains(t, issues, "target_not_long")
		require.Zero(t, r.localCount(`SELECT COUNT(*) FROM write_group_cutovers`))
		require.Equal(t, "", getGroup(t, r.f, targetID)["applied_revision"])
	})

	t.Run("a refused transition leaves the legacy owner recording, and a correct one still works", func(t *testing.T) {
		r := newCutoverRigWith(t, "sqlite", cutoverOptions{wideSeconds: 60, wideTarget: true})
		wideID, targetID := r.wide["id"].(string), r.long["id"].(string)
		boundary := r.at(7 * time.Second)

		untouched := func(t *testing.T) {
			t.Helper()
			require.Zero(t, r.localCount(`SELECT COUNT(*) FROM write_group_cutovers`))
			require.Zero(t, r.localCount(`SELECT COUNT(*) FROM write_group_revision_ends`))
			require.Equal(t, "", getGroup(t, r.f, targetID)["applied_revision"], "the new group was never applied")
			snapshot, err := r.f.groups.ResolveAppliedAt(t.Context(), wideID, time.Now().Add(time.Hour))
			require.NoError(t, err)
			require.True(t, snapshot.EndsAt.IsZero(), "the legacy group still owns its source")
		}

		code, body := r.post(writeGroupsPath+"/"+wideID+"/cutover", r.cutoverBody(boundary, map[string]any{"review_digest": "not-the-digest"}))
		require.Equal(t, http.StatusConflict, code, "%v", body)
		untouched(t)

		past := time.Now().UTC().Add(-time.Minute)
		pastReview := r.review(past)
		require.Equal(t, false, pastReview["applicable"], "a boundary that has passed is not applicable")
		require.NotContains(t, toStrings(pastReview["issues"]), "target_not_long", "the refusal is the boundary, not the layout of a known target")
		untouched(t)

		// The refusals changed nothing about recording: the legacy owner still takes a reading and its window closes.
		r.accept(r.at(3*time.Second), 3.0)
		r.advance(r.at(61 * time.Second))
		r.waitFor(func() bool { return r.wideRows() == 1 }, "the legacy owner keeps recording and delivering after the refusals")
		require.InDelta(t, 3.0, *r.wideRowsInOrder()[0].First, 0)

		r.cutover(r.at(67 * time.Second))
		require.Equal(t, 1, r.localCount(`SELECT COUNT(*) FROM write_group_cutovers`))
	})

	t.Run("a rollback is a reviewed return that deletes neither history", func(t *testing.T) {
		r := newCutoverRigWith(t, "sqlite", cutoverOptions{wideSeconds: 60, wideTarget: true})
		wideID := r.wide["id"].(string)
		forward, rollback := r.at(7*time.Second), r.at(37*time.Second)
		endedRevision := getGroup(t, r.f, wideID)["revision"].(string)
		endedVersion := r.versionPayload(wideID, endedRevision)
		r.accept(r.wideEffective.Add(2*time.Second), 1.5) // history the legacy group wrote before the transition
		r.advance(r.wideEffective.Add(61 * time.Second))
		r.waitFor(func() bool { return r.wideRows() == 1 }, "the legacy history row is delivered")

		r.restart(time.Hour)
		cut := r.cutover(forward)
		r.reconcile()
		r.accept(r.at(3*time.Second), 3.0)
		r.accept(forward, 4.0)
		r.accept(r.at(9*time.Second), 8.0) // the new group's window [7,15)
		r.accept(r.at(33*time.Second), 50.0)
		review := r.rollbackReview(cut["id"].(string), rollback)
		require.Equal(t, true, review["applicable"], "%v", review["issues"])
		code, body := r.post(writeGroupsPath+"/"+wideID+"/cutover", r.rollbackBody(cut["id"].(string), rollback, map[string]any{"review_digest": review["digest"]}))
		require.Equal(t, http.StatusOK, code, "%v", body)
		r.reconcile()

		r.accept(rollback, 11.0)             // exactly the second boundary: the legacy group's again
		r.accept(r.at(65*time.Second), 12.0) // legacy window [60,120)
		r.advance(r.at(121 * time.Second))
		r.restart(cutoverWorkerEvery)
		r.waitFor(func() bool { return r.outbox()["sql_committed"] == 7 }, "both tables drain what they own")

		rows := r.targetRows()
		require.Len(t, rows, 3, "[7,15), the silent [15,30) and the window the rollback cut")
		require.InDelta(t, 6.0, floatOf(t, rows[0].Value), 0, "the mean of 4 and 8 over [7,15)")
		require.False(t, rows[1].Value.Valid)
		require.InDelta(t, 50.0, floatOf(t, rows[2].Value), 0, "the window the rollback cut keeps its reading")
		require.Equal(t, "end", rows[2].Evidence[0].Window["truncated"])
		require.Equal(t, rollback.Format(cutoverTimeText), rows[2].Evidence[0].Window["end"])
		legacy := r.wideRowsInOrder()
		require.Len(t, legacy, 4, "the legacy history, the window the transition cut, and what the restored group wrote")
		require.InDelta(t, 1.5, *legacy[0].First, 0, "the history written before the transition")
		require.InDelta(t, 3.0, *legacy[1].First, 0)
		require.InDelta(t, 11.0, *legacy[2].First, 0)
		require.True(t, rollback.Equal(legacy[2].Start), "a window cut at its start begins at the real boundary: %s", legacy[2].Start)
		require.InDelta(t, 12.0, *legacy[3].First, 0)
		require.Equal(t, endedVersion, r.versionPayload(wideID, endedRevision), "the ended legacy version is not rewritten")
		require.True(t, r.tableExists(r.wideTable()) && r.tableExists(r.longTable()), "no table was removed")
		require.Equal(t, 1, r.localCount(`SELECT COUNT(*) FROM write_group_versions WHERE group_id = ?`, r.long["id"].(string)), "the new group's version is kept")
	})
}
