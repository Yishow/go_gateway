package groupdelivery

import (
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/testutil/lostcommit"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

// TestMySQLSenderFrozenPayloadRestartAfterCommitAckLost verifies that when an insert
// commits on the destination but the acknowledgment is lost, the local store stays in
// retrying without a local receipt. Upon process restart (closing and reopening the
// SQLite store), the exact same frozen payload, digest, and destination are preserved,
// and resending via the production OpenDestination service completes idempotently
// with receipt deduplication, never duplicating rows.
func TestMySQLSenderFrozenPayloadRestartAfterCommitAckLost(t *testing.T) {
	for _, family := range []string{mysqlFixtureFamily, mariaDBFixtureFamily} {
		t.Run(family, func(t *testing.T) {
			fixture := newMySQLSenderFixture(t, family)
			table, proof := fixture.createTables(t)

			store, localDB, sqlitePath := newTestStore(t)
			effectKey := "effect-ack-lost-" + family

			fixture.enqueueCanonicalLong(t, store, table, effectKey, 9007199254740993)

			// Capture original item before any delivery attempt to prove frozen immutability
			beforeItem, err := store.GetOutbox(t.Context(), effectKey)
			require.NoError(t, err)
			require.Equal(t, StatePending, beforeItem.State)
			require.NotEmpty(t, beforeItem.Payload)
			require.NotEmpty(t, beforeItem.PayloadDigest)

			lose := &atomic.Bool{}
			lose.Store(true)
			faults := &lostcommit.Faults{Lose: lose}
			lossyDriver := &mysqldriver.MySQLDriver{}
			lossyDB := lostcommit.WrapFaults(lossyDriver, fixture.dsn, faults)
			t.Cleanup(func() { _ = lossyDB.Close() })

			immediate := func(int) (string, time.Time) {
				return StateRetrying, time.Now().UTC().Add(-time.Second)
			}
			resolver := &staticResolver{
				target: Target{
					DB:    lossyDB,
					Kind:  schema.DatabaseConnectorKindMySQL,
					MySQL: proof,
				},
			}
			sender := NewSender(store, resolver, SenderConfig{
				MaxRetries: 3,
				Backoff:    immediate,
			})

			// 1. First attempt: destination commits, but client ACK is lost
			firstResult, err := sender.Deliver(t.Context(), effectKey)
			require.NoError(t, err)
			require.Equal(t, StateRetrying, firstResult.State)

			// Destination has committed the row and the receipt
			require.Equal(t, 1, fixture.countRows(t, table, ""))
			require.Equal(t, 1, fixture.countRows(t, longTestReceiptTable, "effect_key = ?", effectKey))

			// Local store has no receipt yet because ACK was lost
			require.Zero(t, countRows(t, localDB, "wg_delivery_receipts"), "no local proof before confirmation")

			// 2. Process restart: close store and lossy DB connection
			require.NoError(t, localDB.Close())
			require.NoError(t, lossyDB.Close())

			// Reopen local store from same SQLite file
			reopenedDB := openStoreDB(t, sqlitePath)
			store2 := NewStore(reopenedDB)

			// Reopened outbox item must have exact same payload, digest, and destination as captured before delivery
			item, err := store2.GetOutbox(t.Context(), effectKey)
			require.NoError(t, err)
			require.Equal(t, StateRetrying, item.State)
			require.Equal(t, beforeItem.Payload, item.Payload, "frozen payload must remain byte-for-byte unchanged across restart")
			require.Equal(t, beforeItem.PayloadDigest, item.PayloadDigest, "frozen payload digest must remain unchanged across restart")
			require.Equal(t, beforeItem.Destination, item.Destination, "frozen destination must remain unchanged across restart")

			// 3. Re-deliver over healthy connection via OpenDestination to exercise production saved identity
			resolver2 := mysqlSenderResolver{fixture: fixture, proof: proof}
			sender2 := NewSender(store2, resolver2, SenderConfig{
				MaxRetries: 3,
				Backoff:    immediate,
			})

			// 4. Resend: deduplicated by destination receipt
			secondResult, err := sender2.Deliver(t.Context(), effectKey)
			require.NoError(t, err)
			require.Equal(t, StateCommitted, secondResult.State)

			// Destination still holds exactly one row and one receipt (no duplication)
			require.Equal(t, 1, fixture.countRows(t, table, ""))
			require.Equal(t, 1, fixture.countRows(t, longTestReceiptTable, "effect_key = ?", effectKey))

			// Local store now holds the confirmed receipt
			require.Equal(t, 1, countRows(t, reopenedDB, "wg_delivery_receipts"))

			// Frozen data row was not modified
			recordID := strings.Repeat("e", 64)
			actualValue := fixture.readValueInt(t, table, recordID)
			require.Equal(t, int64(9007199254740993), actualValue)
		})
	}
}

// TestMySQLSenderForeignPayloadSameEffectConflictBlocked verifies that when a foreign
// payload presents an already-committed effect_key with a conflicting payload digest,
// the destination rejects it with a digest mismatch, the sender moves the item to
// blocked with destination-identity-conflict, and the existing row is never overwritten.
func TestMySQLSenderForeignPayloadSameEffectConflictBlocked(t *testing.T) {
	for _, family := range []string{mysqlFixtureFamily, mariaDBFixtureFamily} {
		t.Run(family, func(t *testing.T) {
			fixture := newMySQLSenderFixture(t, family)
			table, proof := fixture.createTables(t)

			store, _, _ := newTestStore(t)
			effectKey := "effect-conflict-" + family

			// 1. Commit and deliver initial row with value 9007199254740993
			fixture.enqueueCanonicalLong(t, store, table, effectKey, 9007199254740993)

			origItem, err := store.GetOutbox(t.Context(), effectKey)
			require.NoError(t, err)

			immediate := func(int) (string, time.Time) {
				return StateRetrying, time.Now().UTC().Add(-time.Second)
			}
			resolver := mysqlSenderResolver{fixture: fixture, proof: proof}
			sender := NewSender(store, resolver, SenderConfig{
				MaxRetries: 3,
				Backoff:    immediate,
			})

			firstResult, err := sender.Deliver(t.Context(), effectKey)
			require.NoError(t, err)
			require.Equal(t, StateCommitted, firstResult.State)
			require.Equal(t, 1, fixture.countRows(t, table, ""))
			require.Equal(t, 1, fixture.countRows(t, longTestReceiptTable, "effect_key = ?", effectKey))

			// 2. Prepare foreign item with same effect_key but conflicting payload
			foreignStore, foreignLocalDB, _ := newTestStore(t)
			fixture.enqueueCanonicalLong(t, foreignStore, table, effectKey, 12345)

			foreignSender := NewSender(foreignStore, resolver, SenderConfig{
				MaxRetries: 3,
				Backoff:    immediate,
			})

			// 3. Delivering conflicting payload must be blocked without overwriting
			foreignResult, err := foreignSender.Deliver(t.Context(), effectKey)
			require.NoError(t, err)
			require.Equal(t, StateBlocked, foreignResult.State)
			require.Equal(t, codeIdentityConflict, foreignResult.Code)

			// 4. Assert foreign store has zero local receipts
			require.Zero(t, countRows(t, foreignLocalDB, "wg_delivery_receipts"), "foreign store must have zero local receipts")

			// 5. Assert receipt payload_digest in destination remains the original digest
			var storedReceiptDigest string
			err = fixture.db.QueryRowContext(t.Context(),
				"SELECT payload_digest FROM `"+fixture.cfg.Database+"`.`"+longTestReceiptTable+"` WHERE effect_key = ?",
				effectKey,
			).Scan(&storedReceiptDigest)
			require.NoError(t, err)
			require.Equal(t, origItem.PayloadDigest, storedReceiptDigest, "receipt payload_digest must remain the original and never be overwritten")

			// 6. Destination data must remain intact and not overwritten
			require.Equal(t, 1, fixture.countRows(t, table, ""))
			recordID := strings.Repeat("e", 64)
			actualValue := fixture.readValueInt(t, table, recordID)
			require.Equal(t, int64(9007199254740993), actualValue)
		})
	}
}
