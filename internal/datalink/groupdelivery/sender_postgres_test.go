package groupdelivery

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/testutil/lostcommit"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

// Live PostgreSQL delivery fixtures; skipped without POSTGRES_DSN. Each test
// uses its own schema, dropped afterwards.

func newPostgresDeliveryFixture(t *testing.T, receiptTable bool) (*deliveryFixture, string, *sql.DB) {
	t.Helper()
	dsn := strings.TrimSpace(os.Getenv("POSTGRES_DSN"))
	if dsn == "" {
		t.Skip("Skipping PostgreSQL delivery test: POSTGRES_DSN not set")
	}
	admin, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = admin.Close() })
	require.NoError(t, admin.PingContext(t.Context()))

	schemaName := fmt.Sprintf("gw_wgd_%d", time.Now().UnixNano())
	_, err = admin.ExecContext(t.Context(), `CREATE SCHEMA "`+schemaName+`"`)
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = admin.ExecContext(context.Background(), `DROP SCHEMA "`+schemaName+`" CASCADE`) })
	_, err = admin.ExecContext(t.Context(), `CREATE TABLE "`+schemaName+`".readings (
		record_id TEXT, counter BIGINT NOT NULL, label TEXT, ratio DOUBLE PRECISION, running BOOLEAN,
		optional TEXT, bucket_start TIMESTAMPTZ)`)
	require.NoError(t, err)
	if receiptTable {
		require.NoError(t, dbtarget.CreateEffectReceiptTable(t.Context(), admin, schema.DatabaseConnectorKindPostgres, schemaName))
	}

	store, local, _ := newTestStore(t)
	lose := &atomic.Bool{}
	faults := &lostcommit.Faults{Lose: lose}
	target := lostcommit.WrapFaults(stdlib.GetDefaultDriver(), dsn, faults)
	t.Cleanup(func() { _ = target.Close() })
	resolver := &staticResolver{target: Target{DB: target, Kind: schema.DatabaseConnectorKindPostgres}}
	immediate := func(int) (string, time.Time) { return StateRetrying, time.Now().UTC().Add(-time.Second) }
	return &deliveryFixture{
		store: store, local: local, target: admin, lose: lose, faults: faults, resolver: resolver,
		sender: NewSender(store, resolver, SenderConfig{MaxRetries: 3, Backoff: immediate}),
	}, schemaName, admin
}

func (f *deliveryFixture) enqueueSchema(t *testing.T, schemaName, capability string, counter uint64) string {
	t.Helper()
	destination := testDestination
	destination.TableSchema = schemaName
	destination.DedupeCapability = capability
	bucket := rowBucket("", t0, counter)
	require.NoError(t, f.store.CommitClosure(t.Context(), testKey,
		Closure{Destination: destination, NextClose: at(10), Buckets: []ClosedBucket{bucket}}))
	return bucket.Outcome.EffectKey
}

func pgRows(t *testing.T, admin *sql.DB, schemaName string) int {
	t.Helper()
	var n int
	require.NoError(t, admin.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM "`+schemaName+`".readings`).Scan(&n))
	return n
}

func TestPostgresTargetReceiptIdentitySenderDeliversAndConvergesAfterLostResponse(t *testing.T) {
	f, schemaName, admin := newPostgresDeliveryFixture(t, true)
	effect := f.enqueueSchema(t, schemaName, DedupeReceipt, 9007199254740993)

	f.lose.Store(true)
	first, err := f.sender.Deliver(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, StateRetrying, first.State)
	require.Equal(t, 1, pgRows(t, admin, schemaName))
	require.Zero(t, countRows(t, f.local, "wg_delivery_receipts"), "no local proof before destination confirmation")

	f.lose.Store(false)
	second, err := f.sender.Deliver(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, StateCommitted, second.State)
	require.Equal(t, 1, pgRows(t, admin, schemaName), "exactly one effect in real PostgreSQL")
	require.Equal(t, 1, countRows(t, f.local, "wg_delivery_receipts"))

	var counter int64
	require.NoError(t, admin.QueryRowContext(t.Context(), `SELECT counter FROM "`+schemaName+`".readings`).Scan(&counter))
	require.Equal(t, int64(9007199254740993), counter)
}

func TestPostgresTargetReceiptIdentityCustomTableWithoutDedupeStopsAtUnknown(t *testing.T) {
	f, schemaName, admin := newPostgresDeliveryFixture(t, false)
	effect := f.enqueueSchema(t, schemaName, DedupeNone, 5)

	f.lose.Store(true)
	first, err := f.sender.Deliver(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, StateUnknown, first.State)
	f.lose.Store(false)
	_, err = f.sender.Deliver(t.Context(), effect)
	require.ErrorIs(t, err, ErrNotDeliverable)
	require.Equal(t, 1, pgRows(t, admin, schemaName), "no blind second insert")
	require.Zero(t, countRows(t, f.local, "wg_delivery_receipts"))
}
