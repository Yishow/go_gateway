package runtime

import (
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/testutil/lostcommit"

	"github.com/stretchr/testify/require"
)

func TestMixedWideReceiptRestartKeepsDueAndMean(t *testing.T) {
	f := mixedBoundaryFixture(t)
	localPath := filepath.Join(t.TempDir(), "journal.db")
	local := openDurableDB(t, localPath)
	store := mixedDurableFixture(t, f, local)
	b, err := NewGroupBoundaryContext(t.Context(), f.config)
	require.NoError(t, err)
	f.accept(t, b, f.envelope(0, "mean-a", boundaryAt(1), int64(20)), f.envelope(0, "mean-b", boundaryAt(2), int64(24)))
	require.NoError(t, b.Tick(t.Context(), boundaryAt(15)))
	var effect string
	require.NoError(t, local.QueryRowContext(t.Context(), `SELECT effect_key FROM wg_delivery_outbox`).Scan(&effect))
	original, err := store.GetOutbox(t.Context(), effect)
	require.NoError(t, err)
	require.True(t, original.BucketStart.IsZero())
	require.True(t, original.DueAt.Equal(boundaryAt(15)))
	var seq int
	var name, path string
	require.NoError(t, f.db.QueryRowContext(t.Context(), `PRAGMA database_list`).Scan(&seq, &name, &path))
	lose := &atomic.Bool{}
	target := lostcommit.Open(path, lose)
	t.Cleanup(func() { target.Close() })
	_, err = target.ExecContext(t.Context(), `CREATE TABLE IF NOT EXISTS gw_effect_receipts(effect_key TEXT PRIMARY KEY,payload_digest TEXT NOT NULL,committed_at TEXT NOT NULL)`)
	require.NoError(t, err)
	resolver := &mixedReceiptTarget{db: target}
	config := groupdelivery.SenderConfig{Backoff: func(int) (string, time.Time) { return groupdelivery.StateRetrying, time.Now().Add(-time.Second) }}
	lose.Store(true)
	first, err := groupdelivery.NewSender(store, resolver, config).Deliver(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, groupdelivery.StateRetrying, first.State)
	require.Equal(t, 1, countTable(t, target, wideFixtureTable))
	require.NoError(t, local.Close())
	reopened := openDurableDB(t, localPath)
	restored := groupdelivery.NewStore(reopened)
	lose.Store(false)
	next, err := groupdelivery.NewSender(restored, resolver, config).Deliver(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, groupdelivery.StateCommitted, next.State)
	after, err := restored.GetOutbox(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, original.Payload, after.Payload)
	require.Equal(t, original.PayloadDigest, after.PayloadDigest)
	require.Equal(t, 1, countTable(t, target, wideFixtureTable))
	var value float64
	var due, provenance string
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT v_gauge,due_at,provenance FROM `+wideFixtureTable).Scan(&value, &due, &provenance))
	require.Equal(t, float64(22), value)
	parsedDue, err := time.Parse(time.RFC3339Nano, due)
	require.NoError(t, err)
	require.True(t, boundaryAt(15).Equal(parsedDue))
	require.Contains(t, provenance, `"calculation":"mean"`)
	require.Contains(t, provenance, `"status":"not_due"`)
	require.Zero(t, countTable(t, reopened, "wg_delivery_counter_baselines"))
	var localDigest, targetDigest string
	require.NoError(t, reopened.QueryRowContext(t.Context(), `SELECT payload_digest FROM wg_delivery_receipts WHERE effect_key=?`, effect).Scan(&localDigest))
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT payload_digest FROM gw_effect_receipts WHERE effect_key=?`, effect).Scan(&targetDigest))
	require.Equal(t, original.PayloadDigest, localDigest)
	require.Equal(t, localDigest, targetDigest)
}
