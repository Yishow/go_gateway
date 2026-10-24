package api

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestCustomPartialCutoverNeedsOmittableColumns: the remaining wide revision of a
// custom table stops writing the column of the source that moved. If the table
// declares that column NOT NULL with no default, the first insert after the
// cutover would fail; the review must refuse, with no change anywhere, while the
// old group keeps writing. A column that can be left out (nullable, or with a
// default the table itself declares) is not refused.
func TestCustomPartialCutoverNeedsOmittableColumns(t *testing.T) {
	cases := []struct {
		column  string
		refused bool
	}{{"not_null", true}, {"nullable", false}, {"default", false}}
	for _, kind := range []string{"sqlite", "postgres"} {
		for _, tc := range cases {
			t.Run(kind+"/"+tc.column, func(t *testing.T) {
				r := newCutoverRigWith(t, kind, cutoverOptions{second: true, custom: tc.column})
				wideID := r.wide["id"].(string)
				boundary := r.at(22*time.Second + 500*time.Millisecond)
				full := getGroup(t, r.f, wideID)
				r.saveRemainingDraft(full, "tag-second")
				before := r.cutoverCounts()

				review := r.review(boundary)
				if !tc.refused {
					require.Equal(t, true, review["applicable"], "%v", review["issues"])
					return
				}
				require.Equal(t, false, review["applicable"], "a NOT NULL column without a default cannot be left out")
				require.Contains(t, toStrings(review["issues"]), "remaining_draft_omits_required_column")
				code, body := r.post(writeGroupsPath+"/"+wideID+"/cutover", r.cutoverBody(boundary, map[string]any{"review_digest": review["digest"]}))
				require.Equal(t, http.StatusUnprocessableEntity, code, "%v", body)
				require.Equal(t, "cutover_not_applicable", errorOf(t, body)["code"])
				require.Equal(t, before, r.cutoverCounts(), "nothing was written")
				require.Zero(t, r.localCount(`SELECT COUNT(*) FROM write_group_versions WHERE group_id = ?`, r.long["id"]), "the profile was never applied")
				require.Zero(t, r.localCount(`SELECT COUNT(*) FROM write_group_revision_ends`))

				// The old group is untouched and keeps writing both of its columns.
				r.acceptFor(r.second, r.wideEffective.Add(2*time.Second), 21.5)
				r.accept(r.wideEffective.Add(2*time.Second), 1.5)
				r.advance(r.wideEffective.Add(16 * time.Second))
				r.waitFor(func() bool {
					var n int
					return r.dest.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM `+r.quoted("custom_readings")+` WHERE a_col IS NOT NULL`).Scan(&n) == nil && n == 1
				}, "the custom group still writes the moved source's column")
			})
		}
	}
}

// partialForward moves tag-A to its profile and leaves tag-second recording in
// the wide group, and returns the cutover.
func (r *cutoverRig) partialForward() (cut map[string]any) {
	r.t.Helper()
	full := getGroup(r.t, r.f, r.wide["id"].(string))
	r.saveRemainingDraft(full, "tag-second")
	r.restart(time.Hour)
	cut = r.cutover(r.at(22*time.Second + 500*time.Millisecond))
	r.reconcile()
	return cut
}

// ruinMovedColumn damages, in the destination, the column the moved source used
// to write: the remaining group never touches it, so it stays ready.
func (r *cutoverRig) ruinMovedColumn() {
	r.t.Helper()
	column := `"` + r.member.TargetColumn + `"`
	statement := `ALTER TABLE ` + r.quoted(r.wideTable()) + ` DROP COLUMN ` + column
	if r.kind != "sqlite" {
		statement = `ALTER TABLE ` + r.quoted(r.wideTable()) + ` ALTER COLUMN ` + column + ` TYPE text`
	}
	_, err := r.dest.ExecContext(r.t.Context(), statement)
	require.NoError(r.t, err)
}

// TestRollbackProvesTheWholeRestoredGroup: a rollback restores the full group,
// not the remaining one, so the full group must be proven against the real
// destination before anything is committed. The remaining revision being ready
// proves nothing about a column it no longer writes.
func TestRollbackProvesTheWholeRestoredGroup(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			r := newCutoverRigWith(t, kind, cutoverOptions{second: true})
			wideID := r.wide["id"].(string)
			cut := r.partialForward()
			rollback := r.at(82*time.Second + 500*time.Millisecond)
			r.ruinMovedColumn()
			before := r.cutoverCounts()
			residual := getGroup(t, r.f, wideID)

			review := r.rollbackReview(cut["id"].(string), rollback)
			require.Equal(t, false, review["applicable"], "the group that would come back cannot be written")
			require.Contains(t, toStrings(review["issues"]), "restored_group_not_ready")
			code, body := r.post(writeGroupsPath+"/"+wideID+"/cutover", r.rollbackBody(cut["id"].(string), rollback, map[string]any{"review_digest": review["digest"]}))
			require.Equal(t, http.StatusUnprocessableEntity, code, "%v", body)
			require.Equal(t, before, r.cutoverCounts(), "nothing was written")
			require.Equal(t, residual, getGroup(t, r.f, wideID), "the remaining revision is still the applied one")
			require.Len(t, r.targets, 1)
			require.Zero(t, r.localCount(`SELECT COUNT(*) FROM write_group_revision_ends WHERE cutover_id <> ?`, cut["id"]), "no revision was ended by the refused rollback")

			// The remaining group and the profile go on working.
			r.acceptFor(r.second, r.at(45*time.Second), 108.0)
			r.accept(r.at(45*time.Second), 8.0)
			r.advance(r.at(61 * time.Second))
			r.restart(cutoverWorkerEvery)
			r.waitFor(func() bool {
				return r.scalar(`SELECT COUNT(*) FROM `+r.quoted(r.wideTable())) == 1 && r.scalar(`SELECT COUNT(*) FROM `+r.quoted(r.longTable())) == 2
			}, "both owners still record after the refused rollback")
		})
	}
}

// TestRollbackDoesNotOverwriteAnUnreviewedDraft: once the remaining group has a
// newer saved draft, that draft is somebody's work. A rollback does not replace
// it silently; it says what is in the way and changes nothing.
func TestRollbackDoesNotOverwriteAnUnreviewedDraft(t *testing.T) {
	r := newCutoverRigWith(t, "sqlite", cutoverOptions{second: true})
	wideID := r.wide["id"].(string)
	cut := r.partialForward()
	rollback := r.at(82*time.Second + 500*time.Millisecond)

	draft := getGroup(t, r.f, wideID)
	draft["row_policy"] = map[string]any{"interval_seconds": 30, "allowed_lateness_seconds": 0}
	response := performJSONRequest(t, r.f.router, http.MethodPut, writeGroupsPath+"/"+wideID, groupEditRequest(r.f.writeGroupRouterFixture, getGroupData(t, r.f, wideID), draft))
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	saved := getGroup(t, r.f, wideID)
	require.NotEqual(t, saved["revision"], saved["applied_revision"], "a newer draft exists")
	before := r.cutoverCounts()

	review := r.rollbackReview(cut["id"].(string), rollback)
	require.Equal(t, false, review["applicable"])
	require.Contains(t, toStrings(review["issues"]), "source_has_unapplied_draft")
	code, body := r.post(writeGroupsPath+"/"+wideID+"/cutover", r.rollbackBody(cut["id"].(string), rollback, map[string]any{"review_digest": review["digest"]}))
	require.Equal(t, http.StatusUnprocessableEntity, code, "%v", body)
	require.Equal(t, before, r.cutoverCounts())
	require.Equal(t, saved, getGroup(t, r.f, wideID), "the unreviewed draft is exactly as it was")
}

// A remaining draft that also changes the group's policy is not the difference
// the cutover can carry out: it is refused with its reason, nothing is written.
func TestPartialCutoverRefusesADraftThatChangesPolicy(t *testing.T) {
	r := newCutoverRigWith(t, "sqlite", cutoverOptions{second: true})
	wideID := r.wide["id"].(string)
	full := getGroup(t, r.f, wideID)
	r.saveRemainingDraft(full, "tag-second")
	draft := getGroup(t, r.f, wideID)
	draft["row_policy"] = map[string]any{"interval_seconds": 30, "allowed_lateness_seconds": 0}
	response := performJSONRequest(t, r.f.router, http.MethodPut, writeGroupsPath+"/"+wideID, groupEditRequest(r.f.writeGroupRouterFixture, getGroupData(t, r.f, wideID), draft))
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	require.Equal(t, "succeeded", string(confirmLong(t, r.f, wideID, previewLong(t, r.f, wideID)).Status))
	before := r.cutoverCounts()

	boundary := r.at(22*time.Second + 500*time.Millisecond)
	review := r.review(boundary)
	require.Equal(t, false, review["applicable"])
	require.Contains(t, toStrings(review["issues"]), "remaining_draft_changes_policy")
	code, body := r.post(writeGroupsPath+"/"+wideID+"/cutover", r.cutoverBody(boundary, map[string]any{"review_digest": review["digest"]}))
	require.Equal(t, http.StatusUnprocessableEntity, code, "%v", body)
	require.Equal(t, before, r.cutoverCounts())
}

// TestCutoverCommitAcknowledgementLost makes the local database's transaction
// commit for real and then report failure, which is what a lost response looks
// like to the service, for the forward cutover and for its rollback. The first
// answer is an error, the database holds exactly one cutover with its ends and
// versions, and the same request sent again (after a restart, after the
// boundary passed) answers with that one cutover and writes nothing more.
func TestCutoverCommitAcknowledgementLost(t *testing.T) {
	r := newCutoverRigWith(t, "sqlite", cutoverOptions{lostCommit: true})
	wideID, longID := r.wide["id"].(string), r.long["id"].(string)
	forward := r.at(22*time.Second + 500*time.Millisecond)
	rollback := r.at(82*time.Second + 500*time.Millisecond)
	cutoverPath := writeGroupsPath + "/" + wideID + "/cutover"

	// Only the cutover's own commit may be lost: the pipeline is stopped meanwhile.
	lose := func(send func() (int, map[string]any)) (int, map[string]any) {
		require.NoError(t, r.pipe.Stop(time.Second))
		r.faults.Lose.Store(true)
		code, body := send()
		r.faults.Lose.Store(false)
		return code, body
	}
	recorded := func(kind string) (id string, boundary string, ends []string) {
		require.NoError(t, r.f.db.QueryRowContext(t.Context(), `SELECT id, boundary FROM write_group_cutovers WHERE kind = ?`, kind).Scan(&id, &boundary))
		rows, err := r.f.db.QueryContext(t.Context(), `SELECT group_id || '@' || ends_at FROM write_group_revision_ends WHERE cutover_id = ? ORDER BY 1`, id)
		require.NoError(t, err)
		defer rows.Close()
		for rows.Next() {
			var end string
			require.NoError(t, rows.Scan(&end))
			ends = append(ends, end)
		}
		require.NoError(t, rows.Err())
		return id, boundary, ends
	}

	review := r.review(forward)
	forwardBody := r.cutoverBody(forward, map[string]any{"review_digest": review["digest"]})
	code, body := lose(func() (int, map[string]any) { return r.post(cutoverPath, forwardBody) })
	require.GreaterOrEqual(t, code, http.StatusInternalServerError, "the answer was lost, so the caller sees a failure: %v", body)
	id, boundary, ends := recorded("forward")
	require.Equal(t, forward.UTC().Format(time.RFC3339Nano), boundary, "the committed boundary is exactly the reviewed one")
	require.Equal(t, []string{wideID + "@" + boundary}, ends, "the wide revision ends once, at that instant")
	require.Equal(t, 1, r.localCount(`SELECT COUNT(*) FROM write_group_versions WHERE group_id = ?`, longID), "and the profile has its one version")
	counts := r.cutoverCounts()

	replay := func(t *testing.T, path string, request map[string]any, wantID string) {
		t.Helper()
		code, again := r.post(path, request)
		require.Equal(t, http.StatusOK, code, "%v", again)
		require.Equal(t, wantID, again["data"].(map[string]any)["id"], "the recorded cutover, not a second one")
		require.Equal(t, counts, r.cutoverCounts(), "nothing was written again")
	}
	t.Run("the same request again", func(t *testing.T) { replay(t, cutoverPath, forwardBody, id) })
	t.Run("after a restart and after the boundary passed", func(t *testing.T) {
		r.launch(cutoverWorkerEvery)
		r.f.groups.WithClock(func() time.Time { return forward.Add(time.Hour) })
		t.Cleanup(func() { r.f.groups.WithClock(nil) })
		replay(t, cutoverPath, forwardBody, id)

		// The source still changes hands at exactly that instant, counted once.
		r.reconcile()
		r.accept(forward.Add(-time.Millisecond), 6.0)
		r.accept(forward, 7.0)
		r.advance(r.at(41 * time.Second))
		r.waitFor(func() bool {
			return r.scalar(`SELECT COUNT(*) FROM `+r.quoted(r.wideTable())) == 1 && r.scalar(`SELECT COUNT(*) FROM `+r.quoted(r.longTable())) == 1
		}, "one cut wide window and one cut profile window are written")
		records := r.longRecords()
		require.True(t, forward.Equal(records[0].Start), "the profile starts at the instant: %s", records[0].Start)
		require.True(t, forward.Equal(records[0].Observed), "an observation at the instant is the profile's")
		require.InDelta(t, 6.0, *r.wideRowsInOrder()[0].First, 0, "and the one just before it is the wide group's")
	})
	t.Run("a different intent is refused", func(t *testing.T) {
		other := cloneMap(forwardBody)
		other["expected_workspace_revision"] = "another-revision"
		code, body := r.post(cutoverPath, other)
		require.Equal(t, http.StatusConflict, code, "%v", body)
		require.Equal(t, counts, r.cutoverCounts())
	})

	// The rollback loses its answer the same way.
	rollbackReview := r.rollbackReview(id, rollback)
	rollbackBody := r.rollbackBody(id, rollback, map[string]any{"review_digest": rollbackReview["digest"]})
	code, body = lose(func() (int, map[string]any) { return r.post(cutoverPath, rollbackBody) })
	require.GreaterOrEqual(t, code, http.StatusInternalServerError, "%v", body)
	rolledID, rolledBoundary, rolledEnds := recorded("rollback")
	require.Equal(t, rollback.UTC().Format(time.RFC3339Nano), rolledBoundary)
	require.Equal(t, []string{longID + "@" + rolledBoundary}, rolledEnds, "the profile ends once, at the new instant")
	require.Equal(t, 2, r.localCount(`SELECT COUNT(*) FROM write_group_versions WHERE group_id = ?`, wideID), "the wide group has the version it had and the restored one")
	counts = r.cutoverCounts()
	t.Run("a rollback request again", func(t *testing.T) { replay(t, cutoverPath, rollbackBody, rolledID) })
	t.Run("a rollback after a restart", func(t *testing.T) {
		r.launch(cutoverWorkerEvery)
		replay(t, cutoverPath, rollbackBody, rolledID)
	})
	t.Run("a rollback of another boundary is refused", func(t *testing.T) {
		code, body := r.post(cutoverPath, r.rollbackBody(id, rollback.Add(time.Second), nil))
		require.GreaterOrEqual(t, code, http.StatusBadRequest, "%v", body)
		require.Equal(t, counts, r.cutoverCounts())
	})
}

// TestCustomTableCutWindowIsWrittenAndTold: a custom table has no provenance or
// bucket column and the gateway adds none. The window a cutover cuts is still
// written, with the values that belong to it, and the cut is told by the local
// delivery truth instead.
func TestCustomTableCutWindowIsWrittenAndTold(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			r := newCutoverRigWith(t, kind, cutoverOptions{second: true, custom: "nullable"})
			wideID := r.wide["id"].(string)
			forward := r.at(22*time.Second + 500*time.Millisecond)
			cut := r.partialForward()
			require.Equal(t, forward.UTC().Format(time.RFC3339Nano), cut["boundary"])

			for _, reading := range []struct {
				at     time.Time
				a, two float64
			}{{r.at(5 * time.Second), 5, 105}, {forward.Add(-time.Millisecond), 6, 106}, {forward, 7, 107}} {
				r.acceptFor(r.second, reading.at, reading.two)
				r.accept(reading.at, reading.a)
			}
			r.advance(r.at(61 * time.Second))
			r.restart(cutoverWorkerEvery)
			r.waitFor(func() bool { return r.scalar(`SELECT COUNT(*) FROM `+r.quoted("custom_readings")) == 3 }, "the complete window and both cut windows are written")

			rows, err := r.dest.QueryContext(t.Context(), `SELECT a_col, b_col FROM `+r.quoted("custom_readings")+` ORDER BY b_col`)
			require.NoError(t, err)
			defer rows.Close()
			var got [][2]*float64
			for rows.Next() {
				var a, b *float64
				require.NoError(t, rows.Scan(&a, &b))
				got = append(got, [2]*float64{a, b})
			}
			require.NoError(t, rows.Err())
			require.Len(t, got, 3)
			require.InDelta(t, 106, *got[1][1], 0, "the old group's cut window keeps what it accepted before the boundary")
			require.InDelta(t, 6, *got[1][0], 0)
			require.Nil(t, got[2][0], "the remaining group no longer writes the moved source")
			require.InDelta(t, 107, *got[2][1], 0)

			query := `SELECT COUNT(*) FROM pragma_table_info('custom_readings')`
			if kind != "sqlite" {
				query = `SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = '` + r.f.namespace + `' AND table_name = 'custom_readings'`
			}
			require.Equal(t, 2, r.scalar(query), "no column was added to the table")

			var truncated int
			for _, issue := range r.delivery(wideID)["recent_bucket_issues"].([]any) {
				if issue.(map[string]any)["kind"] == "truncated" {
					truncated++
				}
			}
			require.Equal(t, 2, truncated, "the delivery truth names both cut windows")
		})
	}
}
