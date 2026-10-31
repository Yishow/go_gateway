package grouppipeline

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"

	"github.com/stretchr/testify/require"
	"modernc.org/sqlite"
)

func assertCompleteCutoverDelivery(t *testing.T, f *ownershipPipelineFixture, shape averageShape) {
	t.Helper()
	expected, targetRows := 4, 2
	if shape.wide() {
		expected, targetRows = 3, 1
	}
	require.Equal(t, expected, f.count(f.local, `SELECT COUNT(*) FROM wg_delivery_outbox`))
	require.Equal(t, expected, f.count(f.local, `SELECT COUNT(*) FROM wg_delivery_outbox WHERE state='sql_committed'`))
	require.Equal(t, expected, f.count(f.local, `SELECT COUNT(*) FROM wg_delivery_receipts`))
	rows, err := f.local.QueryContext(t.Context(), `SELECT effect_key,payload_digest,dedupe_capability FROM wg_delivery_outbox`)
	require.NoError(t, err)
	defer rows.Close()
	for rows.Next() {
		var effect, digest, capability string
		require.NoError(t, rows.Scan(&effect, &digest, &capability))
		var localDigest string
		require.NoError(t, f.local.QueryRowContext(t.Context(), `SELECT payload_digest FROM wg_delivery_receipts WHERE effect_key=?`, effect).Scan(&localDigest))
		require.Equal(t, digest, localDigest)
		var count int
		require.NoError(t, f.plainTarget().QueryRowContext(t.Context(), `SELECT COUNT(*) FROM gw_effect_receipts WHERE effect_key=? AND payload_digest=?`, effect, digest).Scan(&count))
		if capability == groupdelivery.DedupeReceipt {
			require.Equal(t, 1, count)
		} else {
			require.Zero(t, count, "no fabricated destination receipt for legacy none strategy")
		}
	}
	require.NoError(t, rows.Err())
	require.Equal(t, targetRows, f.count(f.plainTarget(), `SELECT COUNT(*) FROM `+f.table))
	require.Equal(t, targetRows, f.count(f.plainTarget(), `SELECT COUNT(DISTINCT record_id) FROM `+f.table))
	require.Equal(t, targetRows, f.count(f.plainTarget(), `SELECT COUNT(*) FROM gw_effect_receipts`))
	require.Equal(t, 2, f.count(f.plainTarget(), `SELECT COUNT(*) FROM legacy_readings`))
	if !shape.wide() {
		var missing int
		require.NoError(t, f.plainTarget().QueryRowContext(t.Context(), `SELECT COUNT(*) FROM `+f.table+` WHERE tag_id='tag-kwh-2' AND value_decimal IS NULL`).Scan(&missing))
		require.Equal(t, 1, missing, "the optional missing member is still one accepted long row")
	}
	// A later normal cycle cannot duplicate either accepted partition.
	f.deliver()
	require.Equal(t, targetRows, f.count(f.plainTarget(), `SELECT COUNT(*) FROM `+f.table))
	require.Equal(t, 2, f.count(f.plainTarget(), `SELECT COUNT(*) FROM legacy_readings`))
}

func TestCutoverLongSQLiteBusyRollbackWaitsForPersistedRetryThenDrains(t *testing.T) {
	f := newHandoffPipelineFixture(t, shapeMeanLong)
	f.readyTarget()
	target := f.pipe.boundaries[boundaryKey(f.target.GroupID, f.target.AppliedRevision)]
	f.clock = at(23)
	require.NoError(t, target.boundary.AcceptSample(t.Context(), envelope(f.target.Group, "target", at(23), uint64(7))))
	require.NoError(t, target.boundary.Tick(t.Context(), at(30)))
	var effect string
	require.NoError(t, f.local.QueryRowContext(t.Context(), `SELECT effect_key FROM wg_delivery_outbox ORDER BY effect_key LIMIT 1`).Scan(&effect))
	item, err := f.store.GetOutbox(t.Context(), effect)
	require.NoError(t, err)
	row, err := groupdelivery.DecodeRowPayload(item.Payload)
	require.NoError(t, err)
	// This disposable connection deliberately refuses waiting, while a real
	// independent writer holds SQLite's lock. Production busy_timeout is unchanged.
	f.dest.target.SetMaxOpenConns(1)
	_, err = f.dest.target.ExecContext(t.Context(), `PRAGMA busy_timeout=0`)
	require.NoError(t, err)
	blocker, err := f.plainTarget().BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer blocker.Rollback()
	_, err = blocker.ExecContext(t.Context(), `UPDATE `+f.table+` SET record_id=record_id WHERE 0`)
	require.NoError(t, err)
	_, err = dbtarget.InsertGroupRow(t.Context(), f.dest.target, dbtarget.GroupInsertRequest{
		Kind: f.dest.kind, SchemaName: item.Destination.TableSchema, TableName: item.Destination.TableName,
		Row: row, PayloadDigest: item.PayloadDigest, Strategy: item.Destination.DedupeCapability,
		RecordKeyColumn: item.Destination.RecordKeyColumn, CommittedAt: time.Now().UTC().Format(time.RFC3339Nano),
	})
	var phase *dbtarget.GroupInsertError
	require.ErrorAs(t, err, &phase)
	require.Equal(t, dbtarget.InsertPhasePreCommit, phase.Phase)
	var busy *sqlite.Error
	require.ErrorAs(t, err, &busy)
	require.Equal(t, 5, busy.Code()&255)
	now := time.Now().UTC()
	due := now.Add(time.Minute)
	config := f.pipe.senderConfig()
	config.Now = func() time.Time { return now }
	config.Backoff = func(int) (string, time.Time) { return groupdelivery.StateRetrying, due }
	sender := groupdelivery.NewSender(f.store, destinationResolver{destinations: f.dest}, config)
	result, err := sender.Deliver(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, groupdelivery.StateRetrying, result.State)
	item, err = f.store.GetOutbox(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, due, item.NextRetryAt)
	require.Equal(t, 1, item.RetryCount)
	require.Equal(t, 0, f.count(f.plainTarget(), `SELECT COUNT(*) FROM `+f.table))
	require.Equal(t, 0, f.count(f.plainTarget(), `SELECT COUNT(*) FROM gw_effect_receipts`))
	require.Equal(t, 0, f.count(f.local, `SELECT COUNT(*) FROM wg_delivery_receipts`))
	require.NoError(t, blocker.Rollback())
	probe, err := f.dest.target.BeginTx(t.Context(), nil)
	require.NoError(t, err, "both precommit attempts cleaned the native transaction")
	require.NoError(t, probe.Rollback())
	dispatcher := groupdelivery.NewDispatcher(f.store, sender, groupdelivery.DispatcherConfig{Now: func() time.Time { return now }, MaxPartitions: 1})
	_, err = dispatcher.RunOnce(t.Context())
	require.NoError(t, err)
	item, err = f.store.GetOutbox(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, groupdelivery.StateRetrying, item.State, "not yet due must not be retried")
	now = due
	for range 4 {
		report, err := dispatcher.RunOnce(t.Context())
		require.NoError(t, err)
		require.Zero(t, report.Unknown)
		require.Zero(t, report.Blocked)
		require.Zero(t, report.Quarantined)
	}
	require.Equal(t, 2, f.count(f.local, `SELECT COUNT(*) FROM wg_delivery_outbox WHERE state='sql_committed'`))
	require.Equal(t, 2, f.count(f.local, `SELECT COUNT(*) FROM wg_delivery_receipts`))
	require.Equal(t, 2, f.count(f.plainTarget(), `SELECT COUNT(*) FROM `+f.table))
	require.Equal(t, 2, f.count(f.plainTarget(), `SELECT COUNT(*) FROM gw_effect_receipts`))
	item, err = f.store.GetOutbox(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, 1, item.RetryCount, "one proven rollback retry, no unknown retry")
	require.Equal(t, 1, f.count(f.plainTarget(), `SELECT COUNT(*) FROM `+f.table+` WHERE value_decimal=7`))
	report, err := dispatcher.RunOnce(t.Context())
	require.NoError(t, err)
	require.Zero(t, report.Partitions)
}
