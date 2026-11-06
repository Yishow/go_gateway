package dbtarget

import (
	"context"
	"database/sql"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
	"modernc.org/sqlite"
)

func reviewSQLiteBusy(t *testing.T) error {
	t.Helper()
	path := filepath.Join(t.TempDir(), "lock.db")
	first, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = first.Close() })
	_, err = first.ExecContext(t.Context(), "CREATE TABLE items (value INTEGER)")
	require.NoError(t, err)
	tx, err := first.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	_, err = tx.ExecContext(t.Context(), "INSERT INTO items VALUES (1)")
	require.NoError(t, err)
	second, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = second.Close() })
	_, busy := second.ExecContext(t.Context(), "INSERT INTO items VALUES (2)")
	var sqliteErr *sqlite.Error
	require.ErrorAs(t, busy, &sqliteErr)
	require.Equal(t, 5, sqliteErr.Code()&0xff)
	require.NoError(t, tx.Rollback())
	return busy
}

type reviewBusyMappings struct {
	TargetMappingRepository
	busy       error
	fail       atomic.Bool
	persistent atomic.Bool
}

func (r *reviewBusyMappings) List(ctx context.Context, filter TargetMappingListFilter) ([]*schema.DatabaseTargetMapping, error) {
	if r.persistent.Load() || r.fail.CompareAndSwap(true, false) {
		return nil, r.busy
	}
	return r.TargetMappingRepository.List(ctx, filter)
}

func TestWriter_PersistentSQLiteBusyRemainsVisible(t *testing.T) {
	connectors, mappings, connector, _ := backoffTestTarget(t)
	createSnapshotMapping(t, mappings, connector.ID, "persistent-busy", "d0")
	repo := &reviewBusyMappings{TargetMappingRepository: mappings, busy: reviewSQLiteBusy(t)}
	writer := NewWriterWithConfig(connectors, repo, WriterConfig{FlushTick: make(chan time.Time)})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	repo.persistent.Store(true)
	require.ErrorIs(t, writer.RefreshMappings(t.Context()), repo.busy)
	accepted, err := writer.AcceptTagValue(t.Context(), "persistent-busy", 12, time.Now())
	require.ErrorIs(t, err, repo.busy)
	require.False(t, accepted)
	repo.persistent.Store(false)
	require.NoError(t, writer.RefreshMappings(t.Context()))
}

type reviewBusyOutcomes struct {
	ConnectorRepository
	busy      error
	failWrite atomic.Bool
	failFlush atomic.Bool
}

func (r *reviewBusyOutcomes) UpdateWriteOutcome(ctx context.Context, id string, at time.Time, status, summary string) error {
	if r.failWrite.CompareAndSwap(true, false) {
		return r.busy
	}
	return r.ConnectorRepository.UpdateWriteOutcome(ctx, id, at, status, summary)
}

func (r *reviewBusyOutcomes) UpdateFlushOutcome(ctx context.Context, id string, at time.Time, status, summary string) error {
	if r.failFlush.CompareAndSwap(true, false) {
		return r.busy
	}
	return r.ConnectorRepository.UpdateFlushOutcome(ctx, id, at, status, summary)
}

func TestWriter_RefreshRetriesTransientSQLiteBusy(t *testing.T) {
	connectors, mappings, connector, _ := backoffTestTarget(t)
	createSnapshotMapping(t, mappings, connector.ID, "busy-refresh", "d0")
	repo := &reviewBusyMappings{TargetMappingRepository: mappings, busy: reviewSQLiteBusy(t)}
	writer := NewWriterWithConfig(connectors, repo, WriterConfig{FlushTick: make(chan time.Time)})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	repo.fail.Store(true)
	require.NoError(t, writer.RefreshMappings(t.Context()))
	accepted, err := writer.AcceptTagValue(t.Context(), "busy-refresh", 12, time.Now())
	require.NoError(t, err)
	require.True(t, accepted)
}

func TestWriter_OutcomesRetryTransientSQLiteBusy(t *testing.T) {
	connectors, mappings, connector, _ := backoffTestTarget(t)
	repo := &reviewBusyOutcomes{ConnectorRepository: connectors, busy: reviewSQLiteBusy(t)}
	writer := NewWriterWithConfig(repo, mappings, WriterConfig{FlushTick: make(chan time.Time)})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	repo.failWrite.Store(true)
	repo.failFlush.Store(true)
	now := time.Now()
	require.NoError(t, writer.recordWriteOutcome(t.Context(), connector, now, deliveryOutcomeSuccess, ""))
	require.NoError(t, writer.recordFlushOutcome(t.Context(), connector, now, deliveryOutcomeSuccess, ""))
	updated, err := connectors.GetByID(t.Context(), connector.ID)
	require.NoError(t, err)
	require.Equal(t, deliveryOutcomeSuccess, updated.LastWriteStatus)
	require.Equal(t, deliveryOutcomeSuccess, updated.LastFlushStatus)
}
