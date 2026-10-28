package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/workspace"
)

func TestCutoverPartialOwnersSurviveDelayedForwardAndRollback(t *testing.T) {
	for _, wide := range []bool{false, true} {
		name := "long"
		if wide {
			name = "wide"
		}
		t.Run(name, func(t *testing.T) {
			r := newScopeAPIRig(t, cutoverOptions{second: true, wideTarget: wide})
			// The retained-owner rollback graph exceeds five seconds under race instrumentation.
			r.reconcileBudget = 15 * time.Second
			forward, rollback := r.at(22500*time.Millisecond), r.at(82500*time.Millisecond)
			targetEnd, targetRollbackStart := r.at(40*time.Second), r.at(80*time.Second)
			if wide {
				targetEnd, targetRollbackStart = r.at(30*time.Second), r.at(75*time.Second)
			}
			wideID, targetID := r.wide["id"].(string), r.long["id"].(string)
			full := getGroup(t, r.f, wideID)
			originalKey := partialAPIKey(r, full)
			originalVersion := r.versionPayload(wideID, originalKey.GroupRevision)
			r.saveRemainingDraft(full, "tag-second")
			r.advance(r.at(21 * time.Second))
			partialAPIOfferTriples(t, r, "forward", forward, 6, 7)
			originals := partialAPIOriginals(t, r)
			require.Len(t, originals, 6)
			cut := r.cutover(forward)
			require.Len(t, cut["started"], 2, "moved target and retained source each get their reviewed revision")
			retained := partialAPIKey(r, getGroup(t, r.f, wideID))
			target := partialAPIKey(r, getGroup(t, r.f, targetID))
			require.NotEqual(t, originalKey, retained)
			partialAPIOwner(t, r, r.member, forward.Add(-time.Millisecond), originalKey)
			partialAPIOwner(t, r, r.second, forward.Add(-time.Millisecond), originalKey)
			partialAPIOwner(t, r, r.member, forward, target)
			partialAPIOwner(t, r, r.second, forward, retained)
			require.Zero(t, r.localCount(`SELECT COUNT(*) FROM wg_delivery_sample_handoffs`))
			partialAPIRefuseUnready(t, r, forward)
			partialAPIStaleTick(t, r, r.at(41*time.Second), 1)
			r.reconcile()
			require.Equal(t, originals, partialAPIOriginals(t, r))
			require.GreaterOrEqual(t, r.flush(r.at(41*time.Second)).Delivered, 3)
			partialAPITransfers(t, r, map[string]groupdelivery.GroupKey{"forward-a": target, "forward-b": retained})
			partialAPIWideRow(t, r, r.at(15*time.Second), forward, 6, 106, "end", "forward-a-before", "forward-b-before")
			partialAPIWideRow(t, r, forward, r.at(30*time.Second), 0, 107, "start", "", "forward-b-after")
			partialAPITargetValue(t, r, wide, forward, targetEnd, r.at(23*time.Second), 7, "start", 2)
			wideHistory, frozen, receipts := r.wideRowsInOrder(), r.frozenOutbox(), r.receipts()
			longHistory, wideTargetHistory := []longRecord(nil), []transitionRow(nil)
			if wide {
				wideTargetHistory = r.targetRows()
			} else {
				longHistory = r.longRecords()
			}
			retainedVersion, targetVersion := r.versionPayload(wideID, retained.GroupRevision), r.versionPayload(targetID, target.GroupRevision)

			// Both current owners accept the next future boundary before learning
			// about rollback; the unselected source is not omitted from the handoff.
			r.advance(r.at(81 * time.Second))
			partialAPIOfferTriples(t, r, "rollback", rollback, 10, 12)
			allOriginals := partialAPIOriginals(t, r)
			for seq, bytes := range originals {
				require.Equal(t, bytes, allOriginals[seq])
			}
			originals = allOriginals
			require.Len(t, originals, 12)
			review := r.rollbackReview(cut["id"].(string), rollback)
			require.Equal(t, true, review["applicable"], "%v", review["issues"])
			code, body := r.post(writeGroupsPath+"/"+wideID+"/cutover", r.rollbackBody(cut["id"].(string), rollback, map[string]any{"review_digest": review["digest"]}))
			require.Equal(t, http.StatusOK, code, "%v", body)
			require.Len(t, body["data"].(map[string]any)["ended"], 2)
			restored := partialAPIKey(r, getGroup(t, r.f, wideID))
			require.NotEqual(t, originalKey, restored)
			require.NotEqual(t, retained, restored)
			partialAPIOwner(t, r, r.member, rollback.Add(-time.Millisecond), target)
			partialAPIOwner(t, r, r.second, rollback.Add(-time.Millisecond), retained)
			partialAPIOwner(t, r, r.member, rollback, restored)
			partialAPIOwner(t, r, r.second, rollback, restored)
			require.Equal(t, 4, r.localCount(`SELECT COUNT(*) FROM wg_delivery_sample_handoffs`))
			partialAPIRefuseUnready(t, r, rollback)
			partialAPIStaleTick(t, r, r.at(106*time.Second), 2)
			r.reconcile()
			require.Equal(t, originals, partialAPIOriginals(t, r))
			require.Positive(t, r.flush(r.at(106*time.Second)).Delivered)
			partialAPITransfers(t, r, map[string]groupdelivery.GroupKey{
				"forward-a": target, "forward-b": retained, "rollback-a": restored, "rollback-b": restored,
			})
			partialAPIWideRow(t, r, r.at(75*time.Second), rollback, 0, 110, "end", "", "rollback-b-before")
			partialAPIWideRow(t, r, rollback, r.at(90*time.Second), 12, 112, "start", "rollback-a-after", "rollback-b-after")
			partialAPITargetValue(t, r, wide, targetRollbackStart, rollback, rollback.Add(-time.Millisecond), 10, "end", 1)
			require.Subset(t, r.wideRowsInOrder(), wideHistory)
			if wide {
				require.Subset(t, r.targetRows(), wideTargetHistory)
			} else {
				require.Subset(t, r.longRecords(), longHistory)
			}
			for key, row := range frozen {
				require.Equal(t, row, r.frozenOutbox()[key])
			}
			for key, receipt := range receipts {
				require.Equal(t, receipt, r.receipts()[key])
			}
			require.Equal(t, originalVersion, r.versionPayload(wideID, originalKey.GroupRevision))
			require.Equal(t, retainedVersion, r.versionPayload(wideID, retained.GroupRevision))
			require.Equal(t, targetVersion, r.versionPayload(targetID, target.GroupRevision))
			stable, stableReceipts, stableDelivery := r.frozenOutbox(), r.receipts(), barrierAPIDeliveryState(t, r.cutoverRig)
			r.reconcile()
			require.Zero(t, r.flush(r.at(106*time.Second)).Delivered)
			require.Equal(t, stable, r.frozenOutbox())
			require.Equal(t, stableReceipts, r.receipts())
			require.Equal(t, stableDelivery, barrierAPIDeliveryState(t, r.cutoverRig))
		})
	}
}

func partialAPIKey(r *scopeAPIRig, group map[string]any) groupdelivery.GroupKey {
	return groupdelivery.GroupKey{WorkspaceID: r.f.record.ID, GroupID: group["id"].(string), GroupRevision: group["applied_revision"].(string)}
}

func partialAPIOwner(t *testing.T, r *scopeAPIRig, member workspace.WriteGroupMember, observed time.Time, expected groupdelivery.GroupKey) {
	t.Helper()
	owners, err := r.store.ResolveCutoverOwnership(t.Context(), groupdelivery.OwnershipSource{
		WorkspaceID: r.f.record.ID, DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID,
	}, observed)
	require.NoError(t, err)
	require.Equal(t, map[groupdelivery.GroupKey]bool{expected: true}, owners.Owners)
}

func partialAPIOfferTriples(t *testing.T, r *scopeAPIRig, phase string, boundary time.Time, before, after float64) {
	t.Helper()
	for _, input := range []struct {
		name string
		at   time.Time
		val  float64
	}{{"after", boundary.Add(500 * time.Millisecond), after}, {"equal", boundary, after}, {"before", boundary.Add(-time.Millisecond), before}} {
		require.NoError(t, r.offer(phase+"-a-"+input.name, r.member, input.at, input.val))
		require.NoError(t, r.offer(phase+"-b-"+input.name, r.second, input.at, input.val+100))
	}
}

func partialAPIOriginals(t *testing.T, r *scopeAPIRig) map[int64]string {
	t.Helper()
	rows, err := r.f.db.QueryContext(t.Context(), `SELECT seq,group_id,group_revision,sample_id,payload,payload_digest,created_at FROM wg_delivery_samples WHERE acceptance_kind='native' ORDER BY seq`)
	require.NoError(t, err)
	defer rows.Close()
	result := map[int64]string{}
	for rows.Next() {
		var seq int64
		var group, revision, id, payload, digest, created string
		require.NoError(t, rows.Scan(&seq, &group, &revision, &id, &payload, &digest, &created))
		result[seq] = strings.Join([]string{group, revision, id, payload, digest, created}, "|")
	}
	require.NoError(t, rows.Err())
	return result
}

func partialAPIRefuseUnready(t *testing.T, r *scopeAPIRig, boundary time.Time) {
	t.Helper()
	before := barrierAPIDeliveryState(t, r.cutoverRig)
	for _, member := range []workspace.WriteGroupMember{r.member, r.second} {
		require.ErrorIs(t, r.offer("unready-"+member.TagID, member, boundary, 999.0), groupdelivery.ErrOwnershipUnavailable)
	}
	require.Equal(t, before, barrierAPIDeliveryState(t, r.cutoverRig))
}

func partialAPIStaleTick(t *testing.T, r *scopeAPIRig, at time.Time, workers int) {
	t.Helper()
	before := barrierAPIDeliveryState(t, r.cutoverRig)
	r.advance(at)
	r.pipe.TickAll(t.Context())
	require.Len(t, r.reports, workers)
	for range workers {
		require.ErrorIs(t, <-r.reports, groupdelivery.ErrOwnershipConflict)
	}
	require.Equal(t, before, barrierAPIDeliveryState(t, r.cutoverRig))
}

func partialAPITransfers(t *testing.T, r *scopeAPIRig, targets map[string]groupdelivery.GroupKey) {
	t.Helper()
	rows, err := r.f.db.QueryContext(t.Context(), `SELECT s.sample_id,s.payload,s.payload_digest,s.created_at,s.consumed,t.group_id,t.group_revision,t.origin_sample_id,t.payload,t.payload_digest,t.created_at,t.consumed,h.transfer_id,t.handoff_id
		FROM wg_delivery_sample_handoffs h JOIN wg_delivery_samples s ON s.seq=h.source_seq JOIN wg_delivery_samples t ON t.seq=h.target_seq WHERE h.state='accepted' ORDER BY s.sample_id`)
	require.NoError(t, err)
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id, payload, digest, created, group, revision, origin, copied, copiedDigest, copiedCreated, transfer, handoff string
		var sourceConsumed, targetConsumed bool
		require.NoError(t, rows.Scan(&id, &payload, &digest, &created, &sourceConsumed, &group, &revision, &origin, &copied, &copiedDigest, &copiedCreated, &targetConsumed, &transfer, &handoff))
		prefix := strings.TrimSuffix(strings.TrimSuffix(id, "-equal"), "-after")
		target, exists := targets[prefix]
		require.True(t, exists, "only boundary-or-later inputs transfer: %s", id)
		require.Equal(t, target.GroupID, group)
		require.Equal(t, target.GroupRevision, revision)
		require.Equal(t, id, origin)
		require.Equal(t, []string{payload, digest, created}, []string{copied, copiedDigest, copiedCreated})
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer, handoff)
		require.False(t, sourceConsumed)
		require.True(t, targetConsumed)
		ids = append(ids, id)
	}
	require.NoError(t, rows.Err())
	var expectedIDs []string
	for prefix := range targets {
		expectedIDs = append(expectedIDs, prefix+"-equal", prefix+"-after")
	}
	require.ElementsMatch(t, expectedIDs, ids)
	require.NoError(t, rows.Close())
	require.Equal(t, len(targets)*5, r.localCount(`SELECT COUNT(*) FROM wg_delivery_samples`), "each three originals has exactly two successor copies")
	require.Equal(t, len(targets)*3, r.localCount(`SELECT COUNT(*) FROM wg_delivery_samples WHERE consumed=1`))
	require.Zero(t, r.localCount(`SELECT COUNT(*) FROM (SELECT COALESCE(NULLIF(origin_sample_id,''),sample_id) FROM wg_delivery_samples WHERE consumed=1 GROUP BY COALESCE(NULLIF(origin_sample_id,''),sample_id) HAVING COUNT(*)<>1)`))
}

func partialAPIWideRow(t *testing.T, r *scopeAPIRig, start, end time.Time, first, second float64, truncation, firstID, secondID string) {
	t.Helper()
	var found []wideRow
	for _, row := range r.wideRowsInOrder() {
		if row.Start.Equal(start) {
			found = append(found, row)
		}
	}
	require.Len(t, found, 1)
	row := found[0]
	if firstID == "" {
		require.Nil(t, row.First, "retained revision cannot write the moved source")
	} else {
		require.NotNil(t, row.First)
		require.Equal(t, first, *row.First)
	}
	require.NotNil(t, row.Second)
	require.Equal(t, second, *row.Second)
	window := windowOfWide(t, row.Provenance)
	require.Equal(t, truncation, window["truncated"])
	require.Equal(t, end, storedInstant(t, window["end"].(string)))
	var evidence []map[string]any
	require.NoError(t, json.Unmarshal([]byte(row.Provenance), &evidence))
	var ids []string
	for _, entry := range evidence {
		id := entry["sample_id"].(string)
		ids = append(ids, id)
		observed := start.Add(500 * time.Millisecond)
		if strings.HasSuffix(id, "-before") {
			observed = end.Add(-time.Millisecond)
		}
		require.Equal(t, observed, storedInstant(t, entry["observed_at"].(string)))
	}
	want := []string{secondID}
	if firstID != "" {
		want = append(want, firstID)
	}
	require.ElementsMatch(t, want, ids)
}

func partialAPITargetValue(t *testing.T, r *scopeAPIRig, wide bool, start, end, observed time.Time, expected float64, truncation string, samples int) {
	t.Helper()
	if wide {
		column := r.long["members"].([]any)[0].(map[string]any)["target_column"].(string)
		var count int
		require.NoError(t, r.dest.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM `+r.quoted(r.longTable())+` WHERE bucket_start=?`, start.Format(time.RFC3339Nano)).Scan(&count))
		require.Equal(t, 1, count, "the target writes exactly one row for this bucket")
		var value float64
		var provenance string
		require.NoError(t, r.dest.QueryRowContext(t.Context(), `SELECT "`+column+`",provenance FROM `+r.quoted(r.longTable())+` WHERE bucket_start=?`, start.Format(time.RFC3339Nano)).Scan(&value, &provenance))
		require.Equal(t, expected, value)
		var entries []map[string]any
		require.NoError(t, json.Unmarshal([]byte(provenance), &entries))
		require.Len(t, entries, 1)
		require.Equal(t, observed, storedInstant(t, entries[0]["observed_at"].(string)))
		require.Equal(t, float64(samples), entries[0]["mean"].(map[string]any)["valid_count"])
		window := windowOfWide(t, provenance)
		require.Equal(t, start, storedInstant(t, window["start"].(string)))
		require.Equal(t, end, storedInstant(t, window["end"].(string)))
		require.Equal(t, truncation, window["truncated"])
		return
	}
	var found []longRecord
	for _, row := range r.longRecords() {
		if row.Start.Equal(start) {
			found = append(found, row)
		}
	}
	require.Len(t, found, 1)
	require.Equal(t, expected, found[0].Value)
	require.Equal(t, observed, found[0].Observed)
	require.Equal(t, end, found[0].End)
	require.Equal(t, truncation, windowOf(t, found[0].Provenance)["truncated"])
}
