package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/groupdelivery"
)

// The spec's rollback40 example keeps two legitimate native I/C acceptances
// before40, but sends exactly40 only to the newly restored wide revision.
// This does not reinterpret raw legacy samples as counter-bound handoffs.
func TestCutoverPublicIndependentProfilesRollback40SQLite(t *testing.T) {
	for _, wide := range []bool{false, true} {
		name := "long"
		if wide {
			name = "wide"
		}
		t.Run(name, func(t *testing.T) {
			r := newScopeAPIRig(t, cutoverOptions{dual: true, wideTarget: wide})
			wideID := r.wide["id"].(string)
			oldRevision := getGroup(t, r.f, wideID)["applied_revision"].(string)
			oldVersion := r.versionPayload(wideID, oldRevision)
			r.advance(r.at(5 * time.Second))
			require.NoError(t, r.offer("historical-wide", r.member, r.at(5*time.Second), uint64(1)))
			require.Equal(t, 1, r.flush(r.at(16*time.Second)).Delivered)
			history, frozen, receipts := r.wideRowsInOrder(), r.frozenOutbox(), r.receipts()
			require.Len(t, history, 1)
			require.Len(t, frozen, 1)
			r.advance(r.at(21 * time.Second))
			forward := r.cutover(r.at(22500 * time.Millisecond))
			require.Len(t, forward["started"], 2)
			r.reconcile()
			profiles := map[string]groupdelivery.GroupKey{}
			for _, target := range r.targets {
				group := getGroup(t, r.f, target["id"].(string))
				kind := group["recording_profile"].(map[string]any)["kind"].(string)
				profiles[kind] = groupdelivery.GroupKey{WorkspaceID: r.f.record.ID, GroupID: group["id"].(string), GroupRevision: group["applied_revision"].(string)}
			}
			require.Len(t, profiles, 2)
			immutable := map[groupdelivery.GroupKey][2]string{}
			for _, key := range []groupdelivery.GroupKey{
				{WorkspaceID: r.f.record.ID, GroupID: wideID, GroupRevision: oldRevision},
				profiles["instantaneous"], profiles["cumulative"],
			} {
				runtime, err := r.store.RuntimeVersion(t.Context(), key)
				require.NoError(t, err)
				require.NotEmpty(t, runtime, "each owner has a frozen runtime descriptor before intake")
				immutable[key] = [2]string{r.versionPayload(key.GroupID, key.GroupRevision), string(runtime)}
			}
			r.advance(r.at(38500 * time.Millisecond))
			observed, rollback := r.at(39999*time.Millisecond), r.at(40*time.Second)
			require.NoError(t, r.offer("same-observation", r.member, observed, uint64(100)), "one observation legitimately enters both profile scopes")
			originals := independentProfileJournal(t, r)
			require.Len(t, originals, 2, "source identity alone must not deduplicate I/C")
			source := groupdelivery.OwnershipSource{WorkspaceID: r.f.record.ID, DeviceID: r.member.DeviceID, PointID: r.member.PointID, TagID: r.member.TagID}
			owner, err := r.store.ResolveCutoverOwnership(t.Context(), source, observed)
			require.NoError(t, err)
			require.Equal(t, map[groupdelivery.GroupKey]bool{profiles["instantaneous"]: true, profiles["cumulative"]: true}, owner.Owners)
			review := r.rollbackReview(forward["id"].(string), rollback)
			require.Equal(t, true, review["applicable"], "%v", review)
			code, body := r.post(writeGroupsPath+"/"+wideID+"/cutover", r.rollbackBody(forward["id"].(string), rollback, map[string]any{"review_digest": review["digest"]}))
			require.Equal(t, http.StatusOK, code, "%v", body)
			require.Len(t, body["data"].(map[string]any)["ended"], 2)
			restored := getGroup(t, r.f, wideID)
			newKey := groupdelivery.GroupKey{WorkspaceID: r.f.record.ID, GroupID: wideID, GroupRevision: restored["applied_revision"].(string)}
			require.NotEqual(t, oldRevision, newKey.GroupRevision)
			owner, err = r.store.ResolveCutoverOwnership(t.Context(), source, rollback)
			require.NoError(t, err)
			require.Equal(t, map[groupdelivery.GroupKey]bool{newKey: true}, owner.Owners, "the exact rollback boundary has only the restored owner")
			beforeRefusal := barrierAPIDeliveryState(t, r.cutoverRig)
			require.ErrorIs(t, r.offer("rollback-exact", r.member, rollback, uint64(200)), groupdelivery.ErrOwnershipUnavailable)
			require.Equal(t, beforeRefusal, barrierAPIDeliveryState(t, r.cutoverRig), "an absent restored worker cannot acknowledge via either former profile")
			require.Equal(t, originals, independentProfileJournal(t, r))
			require.Equal(t, history, r.wideRowsInOrder())
			require.Equal(t, frozen, r.frozenOutbox())
			require.Equal(t, receipts, r.receipts())
			r.advance(r.at(41 * time.Second))
			r.reconcile()
			require.NoError(t, r.offer("rollback-exact", r.member, rollback, uint64(200)))
			require.Equal(t, 1, r.localCount(`SELECT COUNT(*) FROM wg_delivery_samples WHERE sample_id='rollback-exact' AND group_id=? AND group_revision=?`, wideID, newKey.GroupRevision))
			require.Equal(t, 1, r.localCount(`SELECT COUNT(*) FROM wg_delivery_samples WHERE sample_id='rollback-exact'`))
			require.Zero(t, r.localCount(`SELECT COUNT(*) FROM wg_delivery_sample_handoffs`), "pre-boundary native I/C observations retain their original scopes")
			require.GreaterOrEqual(t, r.flush(r.at(61*time.Second)).Delivered, 3)
			require.Equal(t, originals, independentProfileJournal(t, r), "all original payload, proof and acceptance bytes survive rollback")
			var effects []string
			for _, kind := range []string{"instantaneous", "cumulative"} {
				key := profiles[kind]
				require.Equal(t, 1, r.localCount(`SELECT COUNT(*) FROM wg_delivery_samples WHERE sample_id='same-observation' AND group_id=? AND group_revision=? AND consumed=1`, key.GroupID, key.GroupRevision))
				effects = append(effects, independentProfileSQL(t, r, wide, kind, key, observed, rollback))
			}
			require.NotEqual(t, effects[0], effects[1], "independent I/C effects are not deduplicated")
			wideRows := r.wideRowsInOrder()
			require.Len(t, wideRows, 2)
			require.Equal(t, history[0], wideRows[0])
			require.Equal(t, rollback, wideRows[1].Start)
			require.Equal(t, float64(200), *wideRows[1].First)
			var provenance []map[string]any
			require.NoError(t, json.Unmarshal([]byte(wideRows[1].Provenance), &provenance))
			require.Equal(t, "rollback-exact", provenance[0]["sample_id"])
			require.Equal(t, rollback, storedInstant(t, provenance[0]["observed_at"].(string)))
			require.Equal(t, oldVersion, r.versionPayload(wideID, oldRevision))
			for key, bytes := range immutable {
				runtime, err := r.store.RuntimeVersion(t.Context(), key)
				require.NoError(t, err)
				require.Equal(t, bytes, [2]string{r.versionPayload(key.GroupID, key.GroupRevision), string(runtime)}, "rollback preserves source and I/C versions, including original counter proof bytes")
			}
			for key, row := range frozen {
				require.Equal(t, row, r.frozenOutbox()[key], "rollback never changes frozen history or destination")
			}
			for key, receipt := range receipts {
				require.Equal(t, receipt, r.receipts()[key])
			}
			complete, delivered, delivery := r.frozenOutbox(), r.receipts(), barrierAPIDeliveryState(t, r.cutoverRig)
			r.reconcile()
			require.Zero(t, r.flush(r.at(61*time.Second)).Delivered)
			require.Equal(t, complete, r.frozenOutbox())
			require.Equal(t, delivered, r.receipts())
			require.Equal(t, delivery, barrierAPIDeliveryState(t, r.cutoverRig), "repeated drain cannot change journals, checkpoints, proofs or effects")
			require.Equal(t, wideRows, r.wideRowsInOrder())
		})
	}
}

func independentProfileJournal(t *testing.T, r *scopeAPIRig) map[string]string {
	t.Helper()
	rows, err := r.f.db.QueryContext(t.Context(), `SELECT group_id,group_revision,payload,payload_digest,created_at FROM wg_delivery_samples WHERE sample_id='same-observation' ORDER BY group_id`)
	require.NoError(t, err)
	defer rows.Close()
	out := map[string]string{}
	for rows.Next() {
		var group, revision, payload, digest, accepted string
		require.NoError(t, rows.Scan(&group, &revision, &payload, &digest, &accepted))
		out[group+"@"+revision] = payload + "|" + digest + "|" + accepted
	}
	require.NoError(t, rows.Err())
	return out
}

// Read actual destination rows. Long baseline delta is SQL NULL; wide retains
// its raw value and records the absent difference in counter provenance.
func independentProfileSQL(t *testing.T, r *scopeAPIRig, wide bool, kind string, key groupdelivery.GroupKey, observed, rollback time.Time) string {
	t.Helper()
	var recordID, start, provenance string
	if !wide {
		var end, when, raw, status, delta string
		var count int
		require.NoError(t, r.dest.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM `+r.quoted(r.longTable())+` WHERE group_id=? AND observed_at IS NOT NULL`, key.GroupID).Scan(&count))
		require.Equal(t, 1, count, "each independent profile produces one SQL row for its native observation")
		require.NoError(t, r.dest.QueryRowContext(t.Context(), `SELECT record_id,CAST(bucket_start AS TEXT),CAST(bucket_end AS TEXT),CAST(observed_at AS TEXT),CAST(value_decimal AS TEXT),status,COALESCE(CAST(delta_value AS TEXT),'NULL'),provenance FROM `+r.quoted(r.longTable())+` WHERE group_id=? AND observed_at IS NOT NULL`, key.GroupID).Scan(&recordID, &start, &end, &when, &raw, &status, &delta, &provenance))
		require.Equal(t, observed, storedInstant(t, when))
		require.Equal(t, rollback, storedInstant(t, end))
		require.Equal(t, r.at(22500*time.Millisecond), storedInstant(t, start))
		require.Equal(t, "100", raw)
		if kind == "cumulative" {
			require.Equal(t, "baseline", status)
			require.Equal(t, "NULL", delta)
		}
	} else {
		var target map[string]any
		for _, candidate := range r.targets {
			if candidate["id"] == key.GroupID {
				target = candidate
			}
		}
		require.NotNil(t, target)
		column := target["members"].([]any)[0].(map[string]any)["target_column"].(string)
		table := destinationOf(target)["table_name"].(string)
		var count int
		require.NoError(t, r.dest.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM `+r.quoted(table)+` WHERE group_id=? AND json_extract(provenance,'$[0].sample_id')='same-observation'`, key.GroupID).Scan(&count))
		require.Equal(t, 1, count, "each independent profile produces one SQL row for its native observation")
		var value sql.NullString
		require.NoError(t, r.dest.QueryRowContext(t.Context(), `SELECT record_id,CAST(bucket_start AS TEXT),CAST("`+column+`" AS TEXT),provenance FROM `+r.quoted(table)+` WHERE group_id=? AND json_extract(provenance,'$[0].sample_id')='same-observation'`, key.GroupID).Scan(&recordID, &start, &value, &provenance))
		var entries []struct {
			Observed string                 `json:"observed_at"`
			Window   map[string]any         `json:"window"`
			Counter  *counterProvenanceView `json:"counter"`
			Mean     *storedMean            `json:"mean"`
		}
		require.NoError(t, json.Unmarshal([]byte(provenance), &entries))
		require.Len(t, entries, 1)
		require.Equal(t, observed, storedInstant(t, entries[0].Observed))
		require.Equal(t, rollback, storedInstant(t, entries[0].Window["end"].(string)))
		if kind == "cumulative" {
			require.True(t, value.Valid)
			require.Equal(t, "100", value.String, "wide counter value column retains raw digits")
			require.NotNil(t, entries[0].Counter)
			require.Equal(t, "baseline", entries[0].Counter.Status)
			require.Empty(t, entries[0].Counter.Delta, "the new cumulative profile starts without a delta")
			require.Empty(t, entries[0].Counter.DeltaFrom)
			require.Equal(t, r.at(22500*time.Millisecond), storedInstant(t, start))
		} else {
			require.True(t, value.Valid)
			require.Equal(t, "100", value.String)
			require.NotNil(t, entries[0].Mean)
			require.EqualValues(t, 1, entries[0].Mean.ValidCount)
			require.Equal(t, r.at(30*time.Second), storedInstant(t, start))
		}
	}
	var effect string
	require.NoError(t, r.f.db.QueryRowContext(t.Context(), `SELECT effect_key FROM wg_delivery_outbox WHERE record_id=? AND group_id=? AND group_revision=? AND state='sql_committed'`, recordID, key.GroupID, key.GroupRevision).Scan(&effect))
	return effect
}
