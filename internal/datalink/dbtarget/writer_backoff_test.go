package dbtarget

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type backoffTestClock struct {
	mu  sync.Mutex
	now time.Time
}

func (c *backoffTestClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *backoffTestClock) Set(now time.Time) {
	c.mu.Lock()
	c.now = now
	c.mu.Unlock()
}

func backoffTestTarget(t *testing.T) (*SQLConnectorRepository, *SQLTargetMappingRepository, *schema.DatabaseConnector, *sql.DB) {
	t.Helper()
	db := openMigratedTestDB(t)
	connectors, mappings := NewSQLConnectorRepository(db), NewSQLTargetMappingRepository(db)
	path := filepath.Join(t.TempDir(), "backoff-target.db")
	target, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = target.Close() })
	_, err = target.ExecContext(t.Context(), `CREATE TABLE snapshot_values (d0 REAL, d1 REAL)`)
	require.NoError(t, err)
	connector := createSnapshotConnector(t, connectors, schema.DatabaseConnectorKindSQLite, fmt.Sprintf(`{"dsn":%q}`, path), 15)
	return connectors, mappings, connector, target
}

func countHistoryStatus(connectorID, status string) int {
	count := 0
	for _, record := range listWriteHistory(connectorID, 100) {
		if record.Status == status {
			count++
		}
	}
	return count
}

func TestWriter_BackoffRetriesLatestSnapshotWithInjectedClock(t *testing.T) {
	connectors, mappings, connector, target := backoffTestTarget(t)
	createSnapshotMapping(t, mappings, connector.ID, "retry", "d0")
	base := time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	tick := make(chan time.Time, 1)
	attempts := make(chan time.Time, 20)
	var failuresLeft atomic.Int64
	failuresLeft.Store(8)
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{
		FlushTick: tick, Now: clock.Now,
		OpenDestination: func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
			attempts <- clock.Now()
			if failuresLeft.Add(-1) >= 0 {
				return nil, errors.New("destination unavailable")
			}
			return openExternalDBManagerFunc(kind, config)
		},
	})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "retry", 10, base.Add(time.Second)), ErrSnapshotAccepted)
	first := base.Add(15 * time.Second)
	clock.Set(first)
	tick <- first
	select {
	case got := <-attempts:
		require.Equal(t, first, got)
	case <-time.After(2 * time.Second):
		t.Fatal("first snapshot was not attempted")
	}
	require.Eventually(t, func() bool { return countHistoryStatus(connector.ID, "failed") == 1 }, 2*time.Second, 10*time.Millisecond)
	failedConnector, err := connectors.GetByID(t.Context(), connector.ID)
	require.NoError(t, err)
	require.Equal(t, deliveryOutcomeFailed, failedConnector.LastWriteStatus)
	require.Contains(t, failedConnector.LastWriteError, "1 秒後重試")
	require.NotContains(t, failedConnector.LastWriteError, "destination unavailable")
	firstHistory, err := json.Marshal(listWriteHistory(connector.ID, 10))
	require.NoError(t, err)
	require.Contains(t, string(firstHistory), `"retry_after_seconds":1`, "the failure history must expose the backoff")

	previous := first
	for i, delay := range []time.Duration{time.Second, 2 * time.Second, 4 * time.Second, 8 * time.Second, 16 * time.Second, 32 * time.Second, 60 * time.Second, 60 * time.Second} {
		if i == 0 {
			require.ErrorIs(t, writer.WriteTagValue(t.Context(), "retry", 11, previous.Add(100*time.Millisecond)), ErrSnapshotAccepted)
		}
		if i == 1 {
			require.ErrorIs(t, writer.WriteTagValue(t.Context(), "retry", 12, previous.Add(100*time.Millisecond)), ErrSnapshotAccepted)
		}
		clock.Set(previous.Add(delay - time.Millisecond))
		tick <- clock.Now()
		select {
		case got := <-attempts:
			t.Fatalf("attempted before backoff ended at %s", got)
		case <-time.After(30 * time.Millisecond):
		}
		next := previous.Add(delay)
		clock.Set(next)
		tick <- next
		select {
		case got := <-attempts:
			require.Equal(t, next, got, "retry delay %s", delay)
		case <-time.After(2 * time.Second):
			t.Fatalf("no retry after %s", delay)
		}
		if i < 7 {
			wantFailures := i + 2
			require.Eventually(t, func() bool { return countHistoryStatus(connector.ID, "failed") == wantFailures }, 2*time.Second, 10*time.Millisecond)
		}
		previous = next
	}
	require.Eventually(t, func() bool {
		var value float64
		return target.QueryRowContext(t.Context(), `SELECT d0 FROM snapshot_values`).Scan(&value) == nil && value == 12
	}, 2*time.Second, 10*time.Millisecond)
	var rows int
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values`).Scan(&rows))
	require.Equal(t, 1, rows, "backoff must not queue intermediate snapshots")
	require.Eventually(t, func() bool { return countHistoryStatus(connector.ID, "success") == 1 }, 2*time.Second, 10*time.Millisecond)
	require.Greater(t, countHistoryStatus(connector.ID, "skipped"), 0, "empty intervals during backoff must still be explained")
	updated, err := connectors.GetByID(t.Context(), connector.ID)
	require.NoError(t, err)
	require.Equal(t, deliveryOutcomeSuccess, updated.LastWriteStatus)

	// A successful delivery resets the next failure to the initial delay.
	failuresLeft.Store(1)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "retry", 13, base.Add(199*time.Second)), ErrSnapshotAccepted)
	clock.Set(base.Add(210 * time.Second))
	tick <- clock.Now()
	select {
	case <-attempts:
	case <-time.After(2 * time.Second):
		t.Fatal("new interval was not attempted")
	}
	require.Eventually(t, func() bool { return countHistoryStatus(connector.ID, "failed") == 9 }, 2*time.Second, 10*time.Millisecond)
	clock.Set(base.Add(211 * time.Second))
	tick <- clock.Now()
	select {
	case got := <-attempts:
		require.Equal(t, base.Add(211*time.Second), got, "backoff must reset to one second after success")
	case <-time.After(2 * time.Second):
		t.Fatal("backoff did not reset after success")
	}
	require.Eventually(t, func() bool {
		return target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values`).Scan(&rows) == nil && rows == 2
	}, 2*time.Second, 10*time.Millisecond)
}

func TestWriter_RecordsSkippedIntervalBeforeFirstValue(t *testing.T) {
	connectors, mappings, connector, target := backoffTestTarget(t)
	createSnapshotMapping(t, mappings, connector.ID, "never-observed", "d0")
	base := time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	tick := make(chan time.Time, 1)
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: tick, Now: clock.Now})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	clock.Set(base.Add(15 * time.Second))
	tick <- clock.Now()
	require.Eventually(t, func() bool { return countHistoryStatus(connector.ID, "skipped") == 1 }, 2*time.Second, 10*time.Millisecond)
	var rows int
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values`).Scan(&rows))
	require.Zero(t, rows)
}

func TestWriter_RecordsEmptyIntervalsAfterDelayedTick(t *testing.T) {
	connectors, mappings, connector, target := backoffTestTarget(t)
	createSnapshotMapping(t, mappings, connector.ID, "delayed-tick", "d0")
	base := time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	tick := make(chan time.Time, 1)
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: tick, Now: clock.Now})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	clock.Set(base.Add(31 * time.Second))
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "delayed-tick", 12, clock.Now()), ErrSnapshotAccepted)
	clock.Set(base.Add(45 * time.Second))
	tick <- clock.Now()
	require.Eventually(t, func() bool {
		return countHistoryStatus(connector.ID, "success") == 1 && countHistoryStatus(connector.ID, "skipped") == 2
	}, 2*time.Second, 10*time.Millisecond, "a later value must not hide earlier empty intervals")
	var rows int
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values`).Scan(&rows))
	require.Equal(t, 1, rows)
}

func TestWriter_CloseDoesNotRetryDuringBackoff(t *testing.T) {
	connectors, mappings, connector, _ := backoffTestTarget(t)
	createSnapshotMapping(t, mappings, connector.ID, "close-backoff", "d0")
	base := time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	tick, attempts := make(chan time.Time, 1), make(chan struct{}, 2)
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: tick, Now: clock.Now,
		OpenDestination: func(schema.DatabaseConnectorKind, ConnectionConfig) (*datalinkbase.DBManager, error) {
			attempts <- struct{}{}
			return nil, errors.New("destination unavailable")
		},
	})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "close-backoff", 10, base.Add(time.Second)), ErrSnapshotAccepted)
	clock.Set(base.Add(15 * time.Second))
	tick <- clock.Now()
	select {
	case <-attempts:
	case <-time.After(2 * time.Second):
		t.Fatal("first write was not attempted")
	}
	require.Eventually(t, func() bool { return countHistoryStatus(connector.ID, "failed") == 1 }, 2*time.Second, 10*time.Millisecond)
	clock.Set(base.Add(15500 * time.Millisecond))
	require.ErrorContains(t, writer.Close(context.Background()), "backoff")
	require.Empty(t, attempts, "shutdown must not bypass the active backoff")
}

func TestWriter_OutcomeRecordingFailureDoesNotRetryCommittedRow(t *testing.T) {
	connectors, mappings, connector, target := backoffTestTarget(t)
	createSnapshotMapping(t, mappings, connector.ID, "committed", "d0")
	base := time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	tick := make(chan time.Time, 1)
	var attempts atomic.Int64
	writer := NewWriterWithConfig(failingWriteOutcomeRepository{connectors}, mappings, WriterConfig{FlushTick: tick, Now: clock.Now,
		OpenDestination: func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
			attempts.Add(1)
			return openExternalDBManagerFunc(kind, config)
		},
	})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "committed", 12, base.Add(time.Second)), ErrSnapshotAccepted)
	clock.Set(base.Add(15 * time.Second))
	tick <- clock.Now()
	require.Eventually(t, func() bool { return countHistoryStatus(connector.ID, "success") == 1 }, 2*time.Second, 10*time.Millisecond)
	clock.Set(base.Add(30 * time.Second))
	tick <- clock.Now()
	require.Eventually(t, func() bool { return countHistoryStatus(connector.ID, "skipped") == 1 }, 2*time.Second, 10*time.Millisecond)
	require.EqualValues(t, 1, attempts.Load())
	var rows int
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values`).Scan(&rows))
	require.Equal(t, 1, rows)
}

func TestWriter_BackoffIsIndependentPerConnector(t *testing.T) {
	connectors, mappings, fast, target := backoffTestTarget(t)
	slow := createSnapshotConnector(t, connectors, schema.DatabaseConnectorKindPostgres, `{"host":"192.0.2.1","database":"lost"}`, 15)
	createSnapshotMapping(t, mappings, slow.ID, "slow-retry", "d0")
	interval := 1
	require.NoError(t, mappings.Create(t.Context(), &schema.DatabaseTargetMapping{
		ID: "fast-mapping", TagID: "fast-tag", ConnectorID: fast.ID, TableName: "snapshot_values",
		ColumnName: "d0", WriteMode: schema.DatabaseWriteModeInsert, WriteIntervalSeconds: &interval,
		Enabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}))
	base := time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	tick, slowAttempts := make(chan time.Time, 1), make(chan struct{}, 4)
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: tick, Now: clock.Now,
		OpenDestination: func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
			if kind == schema.DatabaseConnectorKindPostgres {
				slowAttempts <- struct{}{}
				return nil, errors.New("destination unavailable")
			}
			return openExternalDBManagerFunc(kind, config)
		},
	})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "slow-retry", 10, base.Add(time.Second)), ErrSnapshotAccepted)
	for i, offset := range []time.Duration{15 * time.Second, 16 * time.Second} {
		clock.Set(base.Add(offset))
		tick <- clock.Now()
		select {
		case <-slowAttempts:
		case <-time.After(2 * time.Second):
			t.Fatalf("slow connector attempt %d did not run", i+1)
		}
		want := i + 1
		require.Eventually(t, func() bool { return countHistoryStatus(slow.ID, "failed") == want }, 2*time.Second, 10*time.Millisecond)
	}
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "fast-tag", 55, base.Add(16100*time.Millisecond)), ErrSnapshotAccepted)
	clock.Set(base.Add(17 * time.Second))
	tick <- clock.Now()
	require.Eventually(t, func() bool {
		var value float64
		return target.QueryRowContext(t.Context(), `SELECT d0 FROM snapshot_values`).Scan(&value) == nil && value == 55
	}, 2*time.Second, 10*time.Millisecond)
	select {
	case <-slowAttempts:
		t.Fatal("slow connector retried before its two-second backoff elapsed")
	case <-time.After(30 * time.Millisecond):
	}
}

func TestWriter_RecordsSkipWhenIntervalHasNoNewValue(t *testing.T) {
	connectors, mappings, connector, target := backoffTestTarget(t)
	createSnapshotMapping(t, mappings, connector.ID, "idle", "d0")
	base := time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	tick := make(chan time.Time, 1)
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: tick, Now: clock.Now})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	for i, value := range []float64{10, 11, 12} {
		require.ErrorIs(t, writer.WriteTagValue(t.Context(), "idle", value, base.Add(time.Duration(1+4*i)*time.Second)), ErrSnapshotAccepted)
	}
	clock.Set(base.Add(15 * time.Second))
	tick <- clock.Now()
	require.Eventually(t, func() bool { return countHistoryStatus(connector.ID, "success") == 1 }, 2*time.Second, 10*time.Millisecond)
	clock.Set(base.Add(30 * time.Second))
	tick <- clock.Now()
	require.Eventually(t, func() bool { return countHistoryStatus(connector.ID, "skipped") == 1 }, 2*time.Second, 10*time.Millisecond)
	var count int
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values`).Scan(&count))
	require.Equal(t, 1, count)
	history, err := NewConnectorService(connectors, mappings).ListWriteHistory(t.Context(), connector.ID, 10)
	require.NoError(t, err)
	encoded, err := json.Marshal(history)
	require.NoError(t, err)
	require.Contains(t, string(encoded), "no_new_value", "the skipped interval must explain why no row was written")
}

func TestWriter_GroupedRetryUsesLatestValues(t *testing.T) {
	connectors, mappings, connector, target := backoffTestTarget(t)
	group := "meter"
	for i, tagID := range []string{"retry-a", "retry-b"} {
		require.NoError(t, mappings.Create(t.Context(), &schema.DatabaseTargetMapping{
			ID: "mapping-" + tagID, TagID: tagID, ConnectorID: connector.ID,
			TableName: "snapshot_values", ColumnName: fmt.Sprintf("d%d", i), GroupKey: &group,
			WriteMode: schema.DatabaseWriteModeInsert, Enabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}))
	}
	base := time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	tick, attempts := make(chan time.Time, 1), make(chan time.Time, 4)
	var failuresLeft atomic.Int64
	failuresLeft.Store(1)
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: tick, Now: clock.Now,
		OpenDestination: func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
			attempts <- clock.Now()
			if failuresLeft.Add(-1) >= 0 {
				return nil, errors.New("destination unavailable")
			}
			return openExternalDBManagerFunc(kind, config)
		},
	})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "retry-a", 10, base.Add(time.Second)), ErrSnapshotAccepted)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "retry-b", 20, base.Add(time.Second)), ErrSnapshotAccepted)
	clock.Set(base.Add(15 * time.Second))
	tick <- clock.Now()
	select {
	case <-attempts:
	case <-time.After(2 * time.Second):
		t.Fatal("first grouped write was not attempted")
	}
	require.Eventually(t, func() bool { return countHistoryStatus(connector.ID, "failed") == 1 }, 2*time.Second, 10*time.Millisecond)
	failed, err := connectors.GetByID(t.Context(), connector.ID)
	require.NoError(t, err)
	require.Equal(t, deliveryOutcomeFailed, failed.LastWriteStatus, "grouped failures must update the connector write outcome")
	require.Contains(t, failed.LastWriteError, "1 秒後重試")
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "retry-a", 11, base.Add(15100*time.Millisecond)), ErrSnapshotAccepted)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "retry-a", 12, base.Add(15200*time.Millisecond)), ErrSnapshotAccepted)
	clock.Set(base.Add(16 * time.Second))
	tick <- clock.Now()
	select {
	case <-attempts:
	case <-time.After(2 * time.Second):
		t.Fatal("grouped snapshot was not retried")
	}
	require.Eventually(t, func() bool {
		var a, b sql.NullFloat64
		return target.QueryRowContext(t.Context(), `SELECT d0, d1 FROM snapshot_values`).Scan(&a, &b) == nil && a.Valid && b.Valid && a.Float64 == 12 && b.Float64 == 20
	}, 2*time.Second, 10*time.Millisecond)
	var rows int
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values`).Scan(&rows))
	require.Equal(t, 1, rows)
	require.Eventually(t, func() bool {
		recovered, err := connectors.GetByID(t.Context(), connector.ID)
		return err == nil && recovered.LastWriteStatus == deliveryOutcomeSuccess
	}, 2*time.Second, 10*time.Millisecond, "connector write outcome must recover after a successful retry")
}

func TestWriter_GroupedRetryCombinesLatestValuesAcrossIntervals(t *testing.T) {
	connectors, mappings, connector, target := backoffTestTarget(t)
	group := "meter"
	for i, tagID := range []string{"blocked-a", "blocked-b"} {
		require.NoError(t, mappings.Create(t.Context(), &schema.DatabaseTargetMapping{
			ID: "mapping-" + tagID, TagID: tagID, ConnectorID: connector.ID,
			TableName: "snapshot_values", ColumnName: fmt.Sprintf("d%d", i), GroupKey: &group,
			WriteMode: schema.DatabaseWriteModeInsert, Enabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}))
	}
	base := time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	tick := make(chan time.Time, 1)
	entered, release := make(chan struct{}), make(chan struct{})
	var releaseOnce sync.Once
	releaseOpen := func() { releaseOnce.Do(func() { close(release) }) }
	var calls atomic.Int64
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: tick, Now: clock.Now,
		OpenDestination: func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
			if calls.Add(1) == 1 {
				close(entered)
				<-release
				return nil, errors.New("destination unavailable")
			}
			return openExternalDBManagerFunc(kind, config)
		},
	})
	t.Cleanup(func() { releaseOpen(); _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "blocked-a", 10, base.Add(time.Second)), ErrSnapshotAccepted)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "blocked-b", 20, base.Add(time.Second)), ErrSnapshotAccepted)
	clock.Set(base.Add(15 * time.Second))
	tick <- clock.Now()
	select {
	case <-entered:
	case <-time.After(2 * time.Second):
		t.Fatal("first grouped write did not start")
	}
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "blocked-a", 11, base.Add(31*time.Second)), ErrSnapshotAccepted)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "blocked-b", 21, base.Add(46*time.Second)), ErrSnapshotAccepted)
	clock.Set(base.Add(47 * time.Second))
	releaseOpen()
	require.Eventually(t, func() bool { return countHistoryStatus(connector.ID, "failed") == 1 }, 2*time.Second, 10*time.Millisecond)
	clock.Set(base.Add(48 * time.Second))
	tick <- clock.Now()
	require.Eventually(t, func() bool { return countHistoryStatus(connector.ID, "success") >= 1 }, 2*time.Second, 10*time.Millisecond)
	var a, b sql.NullFloat64
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT d0, d1 FROM snapshot_values ORDER BY rowid LIMIT 1`).Scan(&a, &b))
	require.True(t, a.Valid && b.Valid, "one retry must contain the latest value for each grouped mapping")
	require.Equal(t, 11.0, a.Float64)
	require.Equal(t, 21.0, b.Float64)
	var rows int
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values`).Scan(&rows))
	require.Equal(t, 1, rows)
}
