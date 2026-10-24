package api

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const energyStart = uint64(9007199254740992)

// TestDualProfileCutover moves one source to its instantaneous and its
// cumulative profile at once. They are two canonical groups of one circuit, so
// one observation legitimately feeds both; what must never happen is the old
// wide group, or a second group of the same kind, recording it too.
func TestDualProfileCutover(t *testing.T) {
	for _, kind := range cutoverKinds {
		t.Run(kind, func(t *testing.T) {
			r := newCutoverRigWith(t, kind, cutoverOptions{dual: true})
			boundary := r.at(22*time.Second + 500*time.Millisecond)

			dup := r.cutoverBody(boundary, nil)
			dup["targets"] = []map[string]any{dup["targets"].([]map[string]any)[0], dup["targets"].([]map[string]any)[0]}
			code, body := r.post(writeGroupsPath+"/"+r.wide["id"].(string)+"/cutover-review", dup)
			require.Equal(t, http.StatusOK, code)
			dupReview := body["data"].(map[string]any)
			require.Equal(t, false, dupReview["applicable"])
			require.Contains(t, toStrings(dupReview["issues"]), "targets_overlap", "the same group twice is a real duplicate")

			// A second group of the same kind for the same source: either it cannot be
			// saved at all, or the review refuses to give one observation to both.
			again := profileRequest(t, r.f.writeGroupRouterFixture, "instantaneous", 30, "second", "alt", "")
			again["group"].(map[string]any)["layout_version"] = "managed_long_v1"
			again["group"].(map[string]any)["storage_mode"] = "plain"
			again["group"].(map[string]any)["destination"].(map[string]any)["table_schema"] = r.f.namespace
			if code, saved := r.post(ensureCircuitProfilePath, again); code == http.StatusOK {
				twin := saved["data"].(map[string]any)["group"].(map[string]any)
				require.Equal(t, "succeeded", string(confirmLong(t, r.f, twin["id"].(string), previewLong(t, r.f, twin["id"].(string))).Status))
				both := r.cutoverBody(boundary, nil)
				both["targets"] = append(both["targets"].([]map[string]any), map[string]any{"group_id": twin["id"], "expected_revision": getGroup(t, r.f, twin["id"].(string))["revision"]})
				code, body := r.post(writeGroupsPath+"/"+r.wide["id"].(string)+"/cutover-review", both)
				require.Equal(t, http.StatusOK, code)
				require.Contains(t, toStrings(body["data"].(map[string]any)["issues"]), "targets_overlap", "two groups of one kind would record the observation twice")
			} else {
				require.GreaterOrEqual(t, code, http.StatusBadRequest)
			}

			r.restart(time.Hour)
			review := r.review(boundary)
			require.Equal(t, true, review["applicable"], "%v", review["issues"])
			require.Len(t, review["targets"], 2)
			cut := r.cutover(boundary)
			require.Len(t, cut["started"], 2, "both profiles start at the boundary")
			require.Len(t, cut["ended"], 1, "the wide revision ends there")
			r.reconcile()

			r.accept(boundary.Add(-time.Millisecond), energyStart-100)
			r.accept(boundary, energyStart)
			r.accept(r.at(45*time.Second), energyStart+1)
			r.accept(r.at(75*time.Second), energyStart+3)
			r.advance(r.at(121 * time.Second))
			r.waitFor(func() bool { return r.outbox()["pending"] == 8 }, "one wide row, five instantaneous and two cumulative records wait")
			r.restart(cutoverWorkerEvery)
			r.waitFor(func() bool { return r.outbox()["sql_committed"] == 8 }, "they are delivered")

			var instant, cumulative []longRecord
			for _, rec := range r.longRecords() {
				if rec.Profile == "cumulative" {
					cumulative = append(cumulative, rec)
				} else {
					instant = append(instant, rec)
				}
			}
			require.Len(t, instant, 5)
			require.Len(t, cumulative, 2)
			require.Equal(t, [2]string{"partial", "9007199254740992"}, [2]string{instant[0].Status, instant[0].Raw}, "the cut window of the instantaneous profile")
			require.Equal(t, [2]string{"ok", "9007199254740993"}, [2]string{instant[1].Status, instant[1].Raw})
			require.Equal(t, [2]string{"ok", "9007199254740995"}, [2]string{instant[2].Status, instant[2].Raw})
			require.Equal(t, [3]string{"baseline", "9007199254740993", "NULL"}, [3]string{cumulative[0].Status, cumulative[0].Raw, cumulative[0].Delta},
				"the cumulative profile starts a fresh baseline: raw kept, no difference, nothing taken from the wide group's reading")
			require.Equal(t, "start", windowOf(t, cumulative[0].Provenance)["truncated"])
			require.Equal(t, [3]string{"ok", "9007199254740995", "2"}, [3]string{cumulative[1].Status, cumulative[1].Raw, cumulative[1].Delta})
			require.Len(t, r.wideRowsInOrder(), 1, "the wide group wrote only the window the boundary cut, nothing at or after it")

			rollback := r.at(82*time.Second + 500*time.Millisecond)
			review = r.rollbackReview(cut["id"].(string), rollback)
			require.Equal(t, true, review["applicable"], "%v", review["issues"])
			code, body = r.post(writeGroupsPath+"/"+r.wide["id"].(string)+"/cutover",
				r.rollbackBody(cut["id"].(string), rollback, map[string]any{"review_digest": review["digest"]}))
			require.Equal(t, http.StatusOK, code, "%v", body)
			require.Len(t, body["data"].(map[string]any)["ended"], 2, "a rollback ends both profiles")
		})
	}
}

// TestCutoverReviewBindsEverythingReviewed: a review is valid only for the
// exact state it read. Changing what it read, even without changing a group's
// revision, makes it stale instead of letting the old digest apply.
func TestCutoverReviewBindsEverythingReviewed(t *testing.T) {
	r := newCutoverRig(t, "sqlite")
	wideID, longID := r.wide["id"].(string), r.long["id"].(string)
	boundary := r.at(22*time.Second + 500*time.Millisecond)

	t.Run("a newer workspace revision does not reuse an older digest", func(t *testing.T) {
		review := r.review(boundary)
		seedProfileSource(t, r.f.writeGroupRouterFixture, "unrelated")
		other := profileRequest(t, r.f.writeGroupRouterFixture, "instantaneous", 30, "second", "other", "unrelated")
		other["group"].(map[string]any)["layout_version"] = "managed_long_v1"
		other["group"].(map[string]any)["storage_mode"] = "plain"
		other["group"].(map[string]any)["destination"].(map[string]any)["table_schema"] = r.f.namespace
		code, _ := r.post(ensureCircuitProfilePath, other)
		require.Equal(t, http.StatusOK, code)

		code, body := r.post(writeGroupsPath+"/"+wideID+"/cutover", r.cutoverBody(boundary, map[string]any{"review_digest": review["digest"]}))
		require.Equal(t, http.StatusConflict, code, "%v", body)
		require.Zero(t, r.localCount(`SELECT COUNT(*) FROM write_group_cutovers`))
	})

	t.Run("a schema confirmation made after the review makes it stale", func(t *testing.T) {
		review := r.review(boundary)
		require.Equal(t, "succeeded", string(confirmLong(t, r.f, longID, previewLong(t, r.f, longID)).Status))
		code, body := r.post(writeGroupsPath+"/"+wideID+"/cutover", r.cutoverBody(boundary, map[string]any{"review_digest": review["digest"]}))
		require.Equal(t, http.StatusConflict, code, "%v", body)
		require.Zero(t, r.localCount(`SELECT COUNT(*) FROM write_group_cutovers`))

		fresh := r.review(boundary)
		require.NotEqual(t, review["digest"], fresh["digest"])
	})
}

// TestCutoverLostResponse: the apply committed but its answer never arrived.
// The same request sent again, after a restart or after the boundary passed,
// returns the same cutover and writes nothing more; a different intent is never
// answered with that success. Forward and rollback alike.
func TestCutoverLostResponse(t *testing.T) {
	r := newCutoverRig(t, "sqlite")
	wideID := r.wide["id"].(string)
	forward := r.at(22*time.Second + 500*time.Millisecond)
	rollback := r.at(82*time.Second + 500*time.Millisecond)

	forwardReview := r.review(forward)
	forwardBody := r.cutoverBody(forward, map[string]any{"review_digest": forwardReview["digest"]})
	code, first := r.post(writeGroupsPath+"/"+wideID+"/cutover", forwardBody)
	require.Equal(t, http.StatusOK, code, "%v", first)
	applied := first["data"].(map[string]any) // the answer the client never saw
	counts := r.cutoverCounts()

	replay := func(t *testing.T, path string, body map[string]any, original map[string]any) {
		t.Helper()
		code, again := r.post(path, body)
		require.Equal(t, http.StatusOK, code, "%v", again)
		got := again["data"].(map[string]any)
		for _, key := range []string{"id", "boundary", "started", "ended", "kind"} {
			require.Equal(t, original[key], got[key], "the replay answers with the recorded %s", key)
		}
		require.Equal(t, counts, r.cutoverCounts(), "a replay writes nothing")
	}
	cutoverPath := writeGroupsPath + "/" + wideID + "/cutover"

	t.Run("the same request again", func(t *testing.T) { replay(t, cutoverPath, forwardBody, applied) })
	t.Run("after a restart", func(t *testing.T) {
		r.restart(cutoverWorkerEvery)
		replay(t, cutoverPath, forwardBody, applied)
	})
	t.Run("after the boundary passed", func(t *testing.T) {
		r.f.groups.WithClock(func() time.Time { return forward.Add(time.Hour) })
		t.Cleanup(func() { r.f.groups.WithClock(nil) })
		replay(t, cutoverPath, forwardBody, applied)
	})
	t.Run("the result can be read", func(t *testing.T) {
		response := performJSONRequest(t, r.f.router, http.MethodGet, writeGroupsPath+"/"+wideID+"/cutovers/"+applied["id"].(string), nil)
		require.Equal(t, http.StatusOK, response.Code, response.Body.String())
		require.Equal(t, applied["id"], decodeJSONBody(t, response)["data"].(map[string]any)["id"])
		list := performJSONRequest(t, r.f.router, http.MethodGet, writeGroupsPath+"/"+wideID+"/cutovers", nil)
		require.Equal(t, http.StatusOK, list.Code, list.Body.String())
		require.Len(t, decodeJSONBody(t, list)["data"], 1)
	})
	t.Run("a different intent is not answered with that success", func(t *testing.T) {
		elsewhere := cloneMap(forwardBody)
		elsewhere["expected_workspace_revision"] = "another-revision"
		code, body := r.post(cutoverPath, elsewhere)
		require.Equal(t, http.StatusConflict, code, "the digest is the original but the request is not: %v", body)

		later := r.cutoverBody(forward.Add(time.Second), nil)
		code, body = r.post(cutoverPath, later)
		require.GreaterOrEqual(t, code, http.StatusBadRequest, "%v", body)
		require.Equal(t, counts, r.cutoverCounts())
	})

	// The same for a rollback.
	rollbackReview := r.rollbackReview(applied["id"].(string), rollback)
	rollbackBody := r.rollbackBody(applied["id"].(string), rollback, map[string]any{"review_digest": rollbackReview["digest"]})
	code, second := r.post(cutoverPath, rollbackBody)
	require.Equal(t, http.StatusOK, code, "%v", second)
	rolled := second["data"].(map[string]any)
	counts = r.cutoverCounts()
	t.Run("a rollback request again", func(t *testing.T) { replay(t, cutoverPath, rollbackBody, rolled) })
	t.Run("a rollback after a restart", func(t *testing.T) {
		r.restart(cutoverWorkerEvery)
		replay(t, cutoverPath, rollbackBody, rolled)
	})
	t.Run("a rollback of another boundary is refused", func(t *testing.T) {
		code, body := r.post(cutoverPath, r.rollbackBody(applied["id"].(string), rollback.Add(time.Second), nil))
		require.GreaterOrEqual(t, code, http.StatusBadRequest, "%v", body)
		require.Equal(t, counts, r.cutoverCounts())
	})
}

// saveRemainingDraft saves, through the ordinary update route, the wide group
// with only the sources that are not moved, and confirms its schema again.
func (r *cutoverRig) saveRemainingDraft(full map[string]any, keepTag string) {
	r.t.Helper()
	wideID := r.wide["id"].(string)
	group := getGroup(r.t, r.f, wideID)
	var kept []any
	for _, m := range full["members"].([]any) {
		if m.(map[string]any)["tag_id"] == keepTag {
			kept = append(kept, m)
		}
	}
	group["members"] = kept
	data := getGroupData(r.t, r.f, wideID)
	response := performJSONRequest(r.t, r.f.router, http.MethodPut, writeGroupsPath+"/"+wideID, groupEditRequest(r.f.writeGroupRouterFixture, data, group))
	require.Equal(r.t, http.StatusOK, response.Code, response.Body.String())
	if r.custom == "" {
		require.Equal(r.t, "succeeded", string(confirmLong(r.t, r.f, wideID, previewLong(r.t, r.f, wideID)).Status))
	}
}

// TestPartialSourceCutover moves one of a wide group's two sources and leaves
// the other recording in the same table. The remaining wide revision is saved
// first as an ordinary draft and takes effect at the same instant as the long
// profile; nothing that was not chosen is moved.
func TestPartialSourceCutover(t *testing.T) {
	for _, kind := range cutoverKinds {
		t.Run(kind, func(t *testing.T) {
			r := newCutoverRigWith(t, kind, cutoverOptions{second: true})
			wideID := r.wide["id"].(string)
			forward := r.at(22*time.Second + 500*time.Millisecond)
			rollback := r.at(82*time.Second + 500*time.Millisecond)
			full := getGroup(t, r.f, wideID)

			// Without a saved draft the second source would be moved silently: refused.
			unreviewed := r.review(forward)
			require.Equal(t, false, unreviewed["applicable"])
			require.Contains(t, toStrings(unreviewed["issues"]), "unselected_sources")

			// A draft that keeps the source the profile takes, and drops the other, is not the difference chosen.
			r.saveRemainingDraft(full, "tag-A")
			wrong := r.review(forward)
			require.Equal(t, false, wrong["applicable"])
			require.Contains(t, toStrings(wrong["issues"]), "remaining_draft_mismatch", "nothing is moved or left behind by omission")

			r.saveRemainingDraft(full, "tag-second")
			review := r.review(forward)
			require.Equal(t, true, review["applicable"], "%v", review["issues"])
			require.Len(t, review["moved"], 1)
			require.Equal(t, "tag-A", review["moved"].([]any)[0].(map[string]any)["tag_id"])
			require.Len(t, review["remaining"], 1)
			require.Equal(t, "tag-second", review["remaining"].([]any)[0].(map[string]any)["tag_id"])

			r.restart(time.Hour)
			cut := r.cutover(forward)
			require.Len(t, cut["started"], 2, "the long profile and the remaining wide revision")
			r.reconcile()
			for _, reading := range []struct {
				at     time.Time
				a, two float64
			}{
				{r.at(5 * time.Second), 5, 105}, {forward.Add(-time.Millisecond), 6, 106}, {forward, 7, 107}, {r.at(45 * time.Second), 8, 108},
			} {
				r.acceptFor(r.second, reading.at, reading.two)
				r.accept(reading.at, reading.a)
			}
			r.advance(r.at(61 * time.Second))
			r.waitFor(func() bool { return r.outbox()["pending"] == 6 }, "four wide rows and two long records wait")
			r.restart(cutoverWorkerEvery)
			r.waitFor(func() bool { return r.outbox()["sql_committed"] == 6 }, "they are delivered")

			wide := r.wideRowsInOrder()
			require.Len(t, wide, 4)
			read := func(row wideRow) (p [2]*float64) { return [2]*float64{row.First, row.Second} }
			require.InDelta(t, 5, *read(wide[0])[0], 0)
			require.InDelta(t, 105, *read(wide[0])[1], 0)
			require.InDelta(t, 6, *read(wide[1])[0], 0, "the old revision's cut window keeps both sources it owned before the boundary")
			require.InDelta(t, 106, *read(wide[1])[1], 0)
			require.Nil(t, wide[2].First, "the remaining revision does not record the moved source")
			require.InDelta(t, 107, *wide[2].Second, 0, "the unselected source goes on without a gap")
			require.True(t, forward.Equal(wide[2].Start), "the remaining revision's first window begins at the boundary: %s", wide[2].Start)
			require.Equal(t, "start", windowOfWide(t, wide[2].Provenance)["truncated"])
			require.Nil(t, wide[3].First)
			require.InDelta(t, 108, *wide[3].Second, 0)
			records := r.longRecords()
			require.Len(t, records, 2)
			require.Equal(t, [2]any{"partial", 7.0}, [2]any{records[0].Status, records[0].Value}, "only the chosen source went to the profile")
			now := getGroup(t, r.f, wideID)
			require.NotEqual(t, full["revision"], now["revision"])
			require.Len(t, now["members"], 1)
			require.Equal(t, full["row_policy"], now["row_policy"], "policy and layout are the ones the group had")
			require.Equal(t, full["destination"].(map[string]any)["table_name"], now["destination"].(map[string]any)["table_name"])

			// A rollback restores the full group and ends both the profile and the remainder.
			rollbackReview := r.rollbackReview(cut["id"].(string), rollback)
			require.Equal(t, true, rollbackReview["applicable"], "%v", rollbackReview["issues"])
			code, body := r.post(writeGroupsPath+"/"+wideID+"/cutover", r.rollbackBody(cut["id"].(string), rollback, map[string]any{"review_digest": rollbackReview["digest"]}))
			require.Equal(t, http.StatusOK, code, "%v", body)
			require.Len(t, body["data"].(map[string]any)["ended"], 2, "the profile and the remaining revision end at the new boundary")
			r.reconcile()
			for _, reading := range []struct {
				at     time.Time
				a, two float64
			}{
				{rollback.Add(-time.Millisecond), 10, 110}, {rollback, 11, 111}, {r.at(100 * time.Second), 12, 112},
			} {
				r.acceptFor(r.second, reading.at, reading.two)
				r.accept(reading.at, reading.a)
			}
			r.advance(r.at(106 * time.Second))
			r.waitFor(func() bool { return r.outbox()["sql_committed"] == 11 }, "three wide rows and two long records are delivered")

			wide = r.wideRowsInOrder()
			require.Len(t, wide, 7)
			require.Nil(t, wide[4].First, "the remaining revision ends at the new boundary")
			require.InDelta(t, 110, *wide[4].Second, 0)
			require.InDelta(t, 11, *wide[5].First, 0, "the restored group records both sources again from the new boundary")
			require.InDelta(t, 111, *wide[5].Second, 0)
			require.True(t, rollback.Equal(wide[5].Start))
			require.InDelta(t, 12, *wide[6].First, 0)
			require.InDelta(t, 112, *wide[6].Second, 0)
			require.Len(t, r.longRecords(), 4, "the profile recorded nothing at or after the new boundary")
		})
	}
}
