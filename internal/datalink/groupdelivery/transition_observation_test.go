package groupdelivery

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"go-gateway/internal/datalink/runtimeview"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDeliveryObservationRequiresExactMatchingReceipt(t *testing.T) {
	for _, mismatch := range []bool{false, true} {
		t.Run(map[bool]string{false: "matching", true: "ON_CONFLICT_mismatched_receipt"}[mismatch], func(t *testing.T) {
			base, db, _ := newTestStore(t)
			var seen []map[string]any
			store := observedDeliveryStore(t, base.WithClock(func() time.Time { return at(20) }), &seen)
			item, claim := deliveryObservationRow(t, store)
			if mismatch {
				_, err := db.ExecContext(t.Context(), `INSERT INTO wg_delivery_receipts(effect_key,payload_digest,committed_at) VALUES(?,?,?)`, item.EffectKey, "wrong-digest", timestamp(at(1)))
				require.NoError(t, err)
			}
			require.NoError(t, store.CompleteDelivery(t.Context(), claim, item.PayloadDigest))
			after, err := store.GetOutbox(t.Context(), item.EffectKey)
			require.NoError(t, err)
			require.Equal(t, StateCommitted, after.State, "logging never rewrites actual native result")
			require.Len(t, seen, 1, "new durable transition must produce exactly one safe observation")
			require.Equal(t, !mismatch, seen[0]["Available"])
			if !mismatch {
				require.Equal(t, item.EffectKey, seen[0]["EffectKey"])
				require.Equal(t, "observed-W1", seen[0]["WorkspaceRevision"])
				require.Equal(t, true, seen[0]["ReceiptMatched"])
			} else {
				require.Equal(t, "", seen[0]["EffectKey"], "unproved receipt must not retain attribution")
			}
		})
	}
}

func TestDeliveryObservationQueuedUnknownNeverClaimCommit(t *testing.T) {
	for _, action := range []string{"retry", "unknown", "claimed-block", "unclaimed-block", "quarantine"} {
		t.Run(action, func(t *testing.T) {
			base, _, _ := newTestStore(t)
			var seen []DeliveryObservation
			store := base.WithDeliveryObservation(DeliveryObservationConfig{
				ScopeRevision: func(context.Context, GroupKey) (string, bool) { return "W1", true },
				Observe:       func(o DeliveryObservation) { seen = append(seen, o) },
			})
			item, claim := deliveryObservationRow(t, store)
			require.Empty(t, seen, "queued and sending are not observed delivery transitions")
			var err error
			wantState := StateBlocked
			switch action {
			case "retry":
				wantState = StateRetrying
				err = store.MarkRetry(t.Context(), claim, wantState, "password_fixture", at(30))
			case "unknown":
				wantState = StateUnknown
				err = store.MarkUnknown(t.Context(), claim, "password_fixture")
			case "claimed-block":
				err = store.BlockClaimed(t.Context(), claim, "password_fixture")
			case "quarantine":
				wantState = StateQuarantined
				err = store.MarkQuarantined(t.Context(), claim, "password_fixture")
			case "unclaimed-block":
				require.NoError(t, store.MarkRetry(t.Context(), claim, StateRetrying, "insert-failed", at(30)))
				seen = nil
				err = store.MarkBlocked(t.Context(), item.EffectKey, "password_fixture")
			}
			require.NoError(t, err)
			require.Len(t, seen, 1)
			o := seen[0]
			require.True(t, o.Available)
			require.Equal(t, item.Key, o.Key)
			require.Equal(t, wantState, o.State)
			require.Equal(t, runtimeview.ReasonUnavailable, o.Reason)
			require.False(t, o.ReceiptMatched)
			require.True(t, o.CommittedAt.IsZero())
			raw, err := json.Marshal(o)
			require.NoError(t, err)
			require.NotContains(t, string(raw), "password_fixture")
			after, err := store.GetOutbox(t.Context(), item.EffectKey)
			require.NoError(t, err)
			require.Equal(t, "password_fixture", after.LastErrorCode, "diagnostics does not rewrite native persisted state")
		})
	}
}

func TestDeliveryObservationAlreadyAndRestartDoNotReplayCommit(t *testing.T) {
	base, _, path := newTestStore(t)
	var seen []DeliveryObservation
	config := DeliveryObservationConfig{ScopeRevision: func(context.Context, GroupKey) (string, bool) { return "W1", true }, Observe: func(o DeliveryObservation) { seen = append(seen, o) }}
	store := base.WithDeliveryObservation(config)
	item, claim := deliveryObservationRow(t, store)
	require.NoError(t, store.CompleteDelivery(t.Context(), claim, item.PayloadDigest))
	require.Len(t, seen, 1)
	seen = nil
	require.ErrorIs(t, store.CompleteDelivery(t.Context(), claim, item.PayloadDigest), ErrNotDeliverable)
	_, err := store.Restore(t.Context(), item.Key)
	require.NoError(t, err)
	restarted := NewStore(openStoreDB(t, path)).WithDeliveryObservation(config)
	result, err := NewSender(restarted, nil, SenderConfig{}).Deliver(t.Context(), item.EffectKey)
	require.NoError(t, err)
	require.True(t, result.Already)
	require.Equal(t, StateCommitted, result.State)
	state, err := restarted.ResolveInterrupted(t.Context(), item.EffectKey)
	require.NoError(t, err)
	require.Equal(t, StateCommitted, state)
	require.Empty(t, seen)
}

func TestDeliveryObservationInterruptedTransitionRaceDoesNotEmit(t *testing.T) {
	for _, race := range []bool{false, true} {
		t.Run(map[bool]string{false: "actual-transition", true: "native-race-no-update"}[race], func(t *testing.T) {
			base, db, _ := newTestStore(t)
			item, _ := deliveryObservationRow(t, base)
			var seen []DeliveryObservation
			reads := 0
			store := base.WithDeliveryObservation(DeliveryObservationConfig{
				ScopeRevision: func(ctx context.Context, _ GroupKey) (string, bool) {
					reads++
					if race && reads == 1 {
						_, err := db.ExecContext(ctx, `UPDATE wg_delivery_outbox SET state='blocked' WHERE effect_key=?`, item.EffectKey)
						require.NoError(t, err)
					}
					return "W1", true
				},
				Observe: func(o DeliveryObservation) { seen = append(seen, o) },
			})
			_, err := store.ResolveInterrupted(t.Context(), item.EffectKey)
			require.NoError(t, err, "preserve native swallowed ErrNotDeliverable behavior")
			if race {
				require.Empty(t, seen)
			} else {
				require.Len(t, seen, 1)
				require.True(t, seen[0].Available)
				require.Contains(t, []string{StateRetrying, StateUnknown}, seen[0].State)
			}
		})
	}
}

func TestDeliveryObservationScopeChangeBecomesUnavailable(t *testing.T) {
	for _, mode := range []string{"missing", "changed", "panic", "identity-change", "claim-owner-change", "receipt-change", "read-error"} {
		t.Run(mode, func(t *testing.T) {
			base, db, _ := newTestStore(t)
			item, claim := deliveryObservationRow(t, base.WithClock(func() time.Time { return at(20) }))
			var seen []DeliveryObservation
			reads := 0
			store := base.WithClock(func() time.Time { return at(20) }).WithDeliveryObservation(DeliveryObservationConfig{
				ScopeRevision: func(ctx context.Context, _ GroupKey) (string, bool) {
					reads++
					if mode == "missing" {
						return "", false
					}
					if mode == "panic" {
						panic("secret provider error")
					}
					if reads == 2 {
						switch mode {
						case "changed":
							return "W2", true
						case "identity-change":
							_, err := db.ExecContext(ctx, `UPDATE wg_delivery_outbox SET connector_revision='C-foreign' WHERE effect_key=?`, item.EffectKey)
							require.NoError(t, err)
						case "claim-owner-change":
							_, err := db.ExecContext(ctx, `UPDATE wg_delivery_outbox SET claim_owner='foreign-owner' WHERE effect_key=?`, item.EffectKey)
							require.NoError(t, err)
						case "receipt-change":
							_, err := db.ExecContext(ctx, `UPDATE wg_delivery_receipts SET payload_digest='wrong' WHERE effect_key=?`, item.EffectKey)
							require.NoError(t, err)
						case "read-error":
							require.NoError(t, db.Close())
						}
					}
					return "W1", true
				},
				Observe: func(o DeliveryObservation) { seen = append(seen, o) },
			})
			require.NoError(t, store.CompleteDelivery(t.Context(), claim, item.PayloadDigest), "proof failure never changes committed native return")
			require.Equal(t, []DeliveryObservation{{Reason: runtimeview.ReasonUnavailable}}, seen)
		})
	}
}

func TestDeliveryObservationCaptureFailurePreservesNativeOutcomeAndReleasesLocks(t *testing.T) {
	for _, mode := range []string{"nil", "panic", "readback-write", "closed-capture"} {
		t.Run(mode, func(t *testing.T) {
			base, db, _ := newTestStore(t)
			db.SetMaxOpenConns(1)
			item, claim := deliveryObservationRow(t, base)
			reads, calls := 0, 0
			config := DeliveryObservationConfig{ScopeRevision: func(context.Context, GroupKey) (string, bool) { reads++; return "W1", true }}
			if mode != "nil" {
				config.Observe = func(o DeliveryObservation) {
					calls++
					require.True(t, o.Available)
					if mode == "panic" {
						panic("raw private exception")
					}
					if mode == "closed-capture" {
						return
					}
					ctx, cancel := context.WithTimeout(t.Context(), time.Second)
					defer cancel()
					_, err := db.ExecContext(ctx, `UPDATE wg_delivery_outbox SET updated_at=updated_at WHERE effect_key=?`, item.EffectKey)
					require.NoError(t, err, "callback runs after original transaction has released its connection/write lock")
				}
			}
			store := base.WithDeliveryObservation(config)
			require.NoError(t, store.CompleteDelivery(t.Context(), claim, item.PayloadDigest))
			after, err := base.GetOutbox(t.Context(), item.EffectKey)
			require.NoError(t, err)
			require.Equal(t, StateCommitted, after.State)
			if mode == "nil" {
				require.Zero(t, reads)
				require.Zero(t, calls)
			} else {
				require.Equal(t, 2, reads)
				require.Equal(t, 1, calls)
			}
			require.Nil(t, base.deliveryObservation.Observe, "clone does not change original store")
		})
	}
}

func TestDeliveryObservationPreservesKnownOwnerReasonsAndRejectsUnsafeFields(t *testing.T) {
	for _, reason := range []string{"recording-ownership-unavailable", "recording-ownership-conflict", "recording-handoff-proof-conflict", "recording-cutover-finalized-overlap"} {
		t.Run(reason, func(t *testing.T) {
			base, _, _ := newTestStore(t)
			var seen []DeliveryObservation
			store := base.WithDeliveryObservation(DeliveryObservationConfig{ScopeRevision: func(context.Context, GroupKey) (string, bool) { return "W1", true }, Observe: func(o DeliveryObservation) { seen = append(seen, o) }})
			_, claim := deliveryObservationRow(t, store)
			require.NoError(t, store.BlockClaimed(t.Context(), claim, reason))
			require.Len(t, seen, 1)
			require.Equal(t, runtimeview.ParseReason(reason), seen[0].Reason)
			raw, err := json.Marshal(seen[0])
			require.NoError(t, err)
			for _, disallowed := range []string{"Destination", "TableName", "Database", "Payload\"", "LastErrorCode", "Error"} {
				require.False(t, strings.Contains(string(raw), disallowed), string(raw))
			}
		})
	}
}

func TestDeliveryObservationCorruptFrozenProofIsUnavailableWithoutChangingCommit(t *testing.T) {
	for _, mode := range []string{"record-id", "effect-key", "entity-key", "payload-digest", "undecodable-payload"} {
		t.Run(mode, func(t *testing.T) {
			base, db, _ := newTestStore(t)
			var seen []DeliveryObservation
			store := base.WithDeliveryObservation(DeliveryObservationConfig{
				ScopeRevision: func(context.Context, GroupKey) (string, bool) { return "W1", true },
				Observe:       func(o DeliveryObservation) { seen = append(seen, o) },
			})
			item, claim := deliveryObservationRow(t, store)
			var err error
			switch mode {
			case "record-id":
				_, err = db.ExecContext(t.Context(), `UPDATE wg_delivery_outbox SET record_id=? WHERE effect_key=?`, "corrupt-record", item.EffectKey)
			case "effect-key":
				claim.EffectKey = "corrupt-effect"
				_, err = db.ExecContext(t.Context(), `UPDATE wg_delivery_outbox SET effect_key=? WHERE effect_key=?`, claim.EffectKey, item.EffectKey)
			case "entity-key":
				_, err = db.ExecContext(t.Context(), `UPDATE wg_delivery_outbox SET entity_key=? WHERE effect_key=?`, "corrupt-entity", item.EffectKey)
			case "payload-digest":
				_, err = db.ExecContext(t.Context(), `UPDATE wg_delivery_outbox SET payload_digest=? WHERE effect_key=?`, strings.Repeat("0", 64), item.EffectKey)
			case "undecodable-payload":
				payload := "not an encoded row"
				digest := sha256.Sum256([]byte(payload))
				_, err = db.ExecContext(t.Context(), `UPDATE wg_delivery_outbox SET payload=?,payload_digest=? WHERE effect_key=?`, payload, hex.EncodeToString(digest[:]), item.EffectKey)
			}
			require.NoError(t, err)
			corrupt, err := store.GetOutbox(t.Context(), claim.EffectKey)
			require.NoError(t, err)
			require.NoError(t, store.CompleteDelivery(t.Context(), claim, corrupt.PayloadDigest), "diagnostic proof must not alter native completion")
			after, err := store.GetOutbox(t.Context(), claim.EffectKey)
			require.NoError(t, err)
			require.Equal(t, StateCommitted, after.State)
			matched, err := store.receiptMatches(t.Context(), sourceRow{state: after.State, effectKey: after.EffectKey, digest: after.PayloadDigest, committedAt: timestamp(after.CommittedAt)})
			require.NoError(t, err)
			require.True(t, matched, "matching receipt alone cannot validate corrupted frozen identity/payload")
			require.Equal(t, []DeliveryObservation{{Reason: runtimeview.ReasonUnavailable}}, seen)
		})
	}
}
