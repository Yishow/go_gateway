package main

import (
	"net/http"
	"testing"
	"time"

	"go-gateway/internal/api"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestProductionDeliveryOperatorResolution(t *testing.T) {
	for _, decision := range []string{"retry", "skip"} {
		t.Run(decision, func(t *testing.T) {
			env, clock := newOutageEnv(t), &testClock{}
			base := time.Now().UTC().Truncate(10 * time.Second).Add(-time.Minute)
			clock.set(base.Add(-time.Second))
			env.services.writeGroups.WithClock(clock.now)
			group := env.createAndApply(t, "A")
			pipe := env.pipeline(clock, "operator")
			clock.set(base.Add(time.Second))
			require.NoError(t, pipe.Reconcile(t.Context()))
			env.feed(t, pipe, clock, []*workspace.WriteGroup{group}, base, 21.5, 100)
			env.feed(t, pipe, clock, []*workspace.WriteGroup{group}, base.Add(10*time.Second), 22.5, 101)
			store := groupdelivery.NewStore(env.db)
			heads, err := store.ReadyHeads(t.Context(), time.Now().Add(time.Minute))
			require.NoError(t, err)
			require.Len(t, heads, 1)
			key := heads[0].EffectKey
			before, err := store.GetOutbox(t.Context(), key)
			require.NoError(t, err)
			// Simulate the sender's durable outcome, retaining the real production row.
			_, err = env.db.ExecContext(t.Context(), `UPDATE wg_delivery_outbox SET state='blocked', last_error_code='missing-table' WHERE effect_key=?`, key)
			require.NoError(t, err)
			heads, err = store.ReadyHeads(t.Context(), time.Now().Add(time.Minute))
			require.NoError(t, err)
			require.Empty(t, heads, "blocked head gates later row")
			router := api.NewRouter(&api.DatalinkServices{Workspace: env.services.workspace, WriteGroups: env.services.writeGroups, WriteGroupDelivery: pipe})
			status, body := httpJSON(t, router, http.MethodGet, testWriteBase+"/write-groups/"+group.ID+"/delivery", nil)
			require.Equal(t, http.StatusOK, status, "%+v", body)
			data := body["data"].(map[string]any)
			attention, ok := data["attention"].([]any)
			require.True(t, ok, "attention must expose durable rows")
			require.Len(t, attention, 1)
			request := map[string]any{"decision_id": "operator-" + decision, "effect_key": key, "expected_state": "blocked", "expected_state_revision": 0, "payload_digest": before.PayloadDigest, "resolution": decision, "reason": "Repaired table or consciously abandoned this row", "confirm_skip": decision == "skip"}
			path := testWriteBase + "/write-groups/" + group.ID + "/delivery/resolve"
			foreign := env.createAndApply(t, "B")
			status, _ = httpJSON(t, router, http.MethodPost, testWriteBase+"/write-groups/"+foreign.ID+"/delivery/resolve", request)
			require.Equal(t, http.StatusNotFound, status)
			stale := copyOperatorRequest(request)
			stale["expected_state"] = "quarantined"
			status, _ = httpJSON(t, router, http.MethodPost, path, stale)
			require.Equal(t, http.StatusConflict, status)
			unsafe := copyOperatorRequest(request)
			unsafe["resolution"] = "skip"
			unsafe["confirm_skip"] = false
			status, _ = httpJSON(t, router, http.MethodPost, path, unsafe)
			require.Equal(t, http.StatusBadRequest, status)
			status, body = httpJSON(t, router, http.MethodPost, path, request)
			require.Equal(t, http.StatusOK, status, "%+v", body)
			status, body = httpJSON(t, router, http.MethodPost, path, request)
			require.Equal(t, http.StatusOK, status, "same decision is idempotent")
			require.True(t, body["data"].(map[string]any)["duplicate"].(bool))
			conflicted := copyOperatorRequest(request)
			conflicted["reason"] = "different decision"
			status, _ = httpJSON(t, router, http.MethodPost, path, conflicted)
			require.Equal(t, http.StatusConflict, status)
			after, err := store.GetOutbox(t.Context(), key)
			require.NoError(t, err)
			require.Equal(t, before.Destination, after.Destination)
			require.Equal(t, before.PayloadDigest, after.PayloadDigest)
			require.Equal(t, before.Payload, after.Payload)
			require.Equal(t, before.RecordID, after.RecordID)
			if decision == "retry" {
				require.Equal(t, groupdelivery.StatePending, after.State)
			} else {
				require.Equal(t, groupdelivery.StateSkipped, after.State)
			}
			var count int
			require.NoError(t, env.db.QueryRowContext(t.Context(), `SELECT count(*) FROM workspace_audit_history WHERE id=? AND event_type='write_group_delivery_resolution' AND reference_id=?`, request["decision_id"], key).Scan(&count))
			require.Equal(t, 1, count)
			heads, err = store.ReadyHeads(t.Context(), time.Now().Add(time.Minute))
			require.NoError(t, err)
			require.Len(t, heads, 1)
			if decision == "retry" {
				_, err = env.db.ExecContext(t.Context(), `UPDATE wg_delivery_outbox SET state='blocked', claim_epoch=claim_epoch+1 WHERE effect_key=?`, key)
				require.NoError(t, err)
				staleAgain := copyOperatorRequest(request)
				staleAgain["decision_id"] = "stale-blocked-again"
				status, _ = httpJSON(t, router, http.MethodPost, path, staleAgain)
				require.Equal(t, http.StatusConflict, status, "same state after a new attempt is a different decision snapshot")
			}
			next := heads[0].EffectKey
			_, err = env.db.ExecContext(t.Context(), `UPDATE wg_delivery_outbox SET state='unknown' WHERE effect_key=?`, next)
			require.NoError(t, err)
			unknown := copyOperatorRequest(request)
			unknown["decision_id"] = "unknown-" + decision
			unknown["effect_key"] = next
			unknown["expected_state"] = "unknown"
			status, _ = httpJSON(t, router, http.MethodPost, path, unknown)
			require.Equal(t, http.StatusConflict, status, "unknown cannot be blindly retried or skipped")
		})
	}
}

func copyOperatorRequest(input map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range input {
		out[k] = v
	}
	return out
}
