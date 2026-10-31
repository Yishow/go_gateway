package runtime

import (
	"context"
	"database/sql"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/testutil/lostcommit"

	"github.com/stretchr/testify/require"
)

type mixedReceiptTarget struct {
	db    *sql.DB
	calls int
}

func (r *mixedReceiptTarget) Resolve(context.Context, groupdelivery.OutboxItem) (groupdelivery.Target, error) {
	r.calls++
	return groupdelivery.Target{DB: r.db, Kind: schema.DatabaseConnectorKindSQLite}, nil
}

func TestMixedDurableFrozenDeltaReceiptRecovery(t *testing.T) {
	for _, scaled := range []bool{false, true} {
		for _, mode := range []string{"lost-response", "local-receipt-failure", "unknown-no-dedupe", "receipt-conflict"} {
			t.Run(map[bool]string{false: "unscaled", true: "scaled"}[scaled]+"/"+mode, func(t *testing.T) {
				f := mixedCounterBoundaryFixture(t, scaled)
				if mode == "unknown-no-dedupe" {
					f.config.Group.WritePolicy.DedupeCapability = groupdelivery.DedupeNone
				}
				localPath := filepath.Join(t.TempDir(), "journal.db")
				local := openDurableDB(t, localPath)
				store := mixedDurableFixture(t, f, local)
				b, err := NewGroupBoundaryContext(t.Context(), f.config)
				require.NoError(t, err)
				f.accept(t, b, f.envelope(0, "baseline", boundaryAt(14), uint64(9007199254740993)))
				require.NoError(t, b.Tick(t.Context(), boundaryAt(15)))
				f.accept(t, b, f.envelope(0, "delta-tail", boundaryAt(29), uint64(9007199254741001)))
				require.NoError(t, b.Tick(t.Context(), boundaryAt(30)))
				var effect string
				require.NoError(t, local.QueryRowContext(t.Context(), `SELECT effect_key FROM wg_delivery_outbox`).Scan(&effect))
				before, err := store.GetOutbox(t.Context(), effect)
				require.NoError(t, err)
				require.True(t, before.DueAt.Equal(boundaryAt(30)))
				original, err := groupdelivery.DecodeRowPayload(before.Payload)
				require.NoError(t, err)
				require.Equal(t, before.RecordID, original.RecordID)
				var baselineBefore string
				require.NoError(t, local.QueryRowContext(t.Context(), `SELECT payload FROM wg_delivery_counter_baselines`).Scan(&baselineBefore))
				var seq int
				var name, targetPath string
				require.NoError(t, f.db.QueryRowContext(t.Context(), `PRAGMA database_list`).Scan(&seq, &name, &targetPath))
				lose := &atomic.Bool{}
				target := lostcommit.Open(targetPath, lose)
				t.Cleanup(func() { target.Close() })
				_, err = target.ExecContext(t.Context(), `CREATE TABLE IF NOT EXISTS gw_effect_receipts(effect_key TEXT PRIMARY KEY,payload_digest TEXT NOT NULL,committed_at TEXT NOT NULL)`)
				require.NoError(t, err)
				resolver := &mixedReceiptTarget{db: target}
				config := groupdelivery.SenderConfig{Backoff: func(int) (string, time.Time) { return groupdelivery.StateRetrying, time.Now().Add(-time.Second) }}
				sender := groupdelivery.NewSender(store, resolver, config)
				if mode == "local-receipt-failure" {
					_, err = local.ExecContext(t.Context(), `CREATE TRIGGER mixed_receipt_failure BEFORE INSERT ON wg_delivery_receipts BEGIN SELECT RAISE(ABORT,'receipt failure'); END`)
					require.NoError(t, err)
				} else {
					lose.Store(true)
				}
				first, err := sender.Deliver(t.Context(), effect)
				if mode == "local-receipt-failure" {
					require.ErrorIs(t, err, groupdelivery.ErrLocalReceiptFailed)
				} else {
					require.NoError(t, err)
				}
				expectedState := groupdelivery.StateRetrying
				if mode == "unknown-no-dedupe" {
					expectedState = groupdelivery.StateUnknown
				}
				if mode != "local-receipt-failure" {
					require.Equal(t, expectedState, first.State)
				}
				require.Equal(t, 1, countTable(t, target, longFixtureTable))
				require.Zero(t, countTable(t, local, "wg_delivery_receipts"))
				lose.Store(false)
				if mode == "receipt-conflict" {
					_, err = target.ExecContext(t.Context(), `UPDATE gw_effect_receipts SET payload_digest='different-accepted-payload' WHERE effect_key=?`, effect)
					require.NoError(t, err)
				}
				if mode == "local-receipt-failure" {
					_, err = local.ExecContext(t.Context(), `DROP TRIGGER mixed_receipt_failure`)
					require.NoError(t, err)
				}
				// Reopen the accepted local evidence; no assembler or current draft is
				// consulted by the sender to reconstruct this already finalized result.
				require.NoError(t, local.Close())
				reopened := openDurableDB(t, localPath)
				recoveredStore := groupdelivery.NewStore(reopened)
				if mode == "local-receipt-failure" {
					state, err := recoveredStore.ResolveInterrupted(t.Context(), effect)
					require.NoError(t, err)
					require.Equal(t, groupdelivery.StateRetrying, state)
				}
				recovered := groupdelivery.NewSender(recoveredStore, resolver, config)
				// Deliberately mutate the caller's now-irrelevant draft calculation.
				f.config.Group.Members[0].Cadence.Interval = 60
				f.config.Group.Members[0].Calculation = "latest"
				for i := 0; i < 2; i++ {
					result, err := recovered.Deliver(t.Context(), effect)
					switch mode {
					case "unknown-no-dedupe":
						require.ErrorIs(t, err, groupdelivery.ErrNotDeliverable)
						require.Equal(t, groupdelivery.StateUnknown, result.State)
						attention, err := recoveredStore.Attention(t.Context(), before.Key.GroupID)
						require.NoError(t, err)
						require.Len(t, attention, 1)
						require.Equal(t, effect, attention[0].EffectKey)
					case "receipt-conflict":
						if i == 0 {
							require.NoError(t, err)
						} else {
							require.ErrorIs(t, err, groupdelivery.ErrNotDeliverable)
						}
						require.Equal(t, groupdelivery.StateBlocked, result.State)
					default:
						require.NoError(t, err)
						require.Equal(t, groupdelivery.StateCommitted, result.State)
					}
					digest, found, err := dbtarget.ReadEffectReceiptDigest(t.Context(), target, schema.DatabaseConnectorKindSQLite, "main", effect)
					require.NoError(t, err)
					require.Equal(t, mode != "unknown-no-dedupe", found)
					if found {
						if mode == "receipt-conflict" {
							require.Equal(t, "different-accepted-payload", digest)
						} else {
							require.Equal(t, before.PayloadDigest, digest)
						}
					}
					require.Equal(t, 1, countTable(t, target, longFixtureTable), "inspection/retry cannot create a second result")
				}
				after, err := recoveredStore.GetOutbox(t.Context(), effect)
				require.NoError(t, err)
				require.Equal(t, before.Payload, after.Payload)
				require.Equal(t, before.PayloadDigest, after.PayloadDigest)
				var baselineAfter string
				require.NoError(t, reopened.QueryRowContext(t.Context(), `SELECT payload FROM wg_delivery_counter_baselines`).Scan(&baselineAfter))
				require.Equal(t, baselineBefore, baselineAfter, "delivery cannot recalculate or advance baseline")
				var value, provenance string
				var delta sql.NullString
				require.NoError(t, target.QueryRowContext(t.Context(), `SELECT value_decimal,delta_value,provenance FROM `+longFixtureTable).Scan(&value, &delta, &provenance))
				expected := "8"
				if scaled {
					expected = "0.8"
				}
				require.Equal(t, expected, value)
				require.False(t, delta.Valid, "mixed result is the selected value, not legacy delta side-column")
				rawText := "9007199254741001"
				if scaled {
					rawText = "900719925474100.1"
				}
				require.Contains(t, provenance, rawText)
				require.Contains(t, provenance, `"calculation":"delta"`)
				require.Contains(t, provenance, `"interval_seconds":15`)
				switch mode {
				case "unknown-no-dedupe":
					require.Equal(t, 1, resolver.calls)
					require.Zero(t, countTable(t, reopened, "wg_delivery_receipts"))
				case "receipt-conflict":
					require.Zero(t, countTable(t, reopened, "wg_delivery_receipts"))
					require.Equal(t, "destination-identity-conflict", after.LastErrorCode)
				default:
					require.Equal(t, 1, countTable(t, reopened, "wg_delivery_receipts"))
					var digest string
					require.NoError(t, reopened.QueryRowContext(t.Context(), `SELECT payload_digest FROM wg_delivery_receipts WHERE effect_key=?`, effect).Scan(&digest))
					require.Equal(t, before.PayloadDigest, digest)
				}
			})
		}
	}
}
