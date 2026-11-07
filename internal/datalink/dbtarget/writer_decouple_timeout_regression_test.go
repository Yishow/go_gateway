package dbtarget

import (
	"context"
	"database/sql"
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

func TestSnapshotWorker_HungOpenFailsAndStaysIsolated(t *testing.T) {
	for _, grouped := range []bool{false, true} {
		t.Run(fmt.Sprintf("grouped=%t", grouped), func(t *testing.T) {
			mainDB := openMigratedTestDB(t)
			connectors := NewSQLConnectorRepository(mainDB)
			mappings := NewSQLTargetMappingRepository(mainDB)
			path := filepath.Join(t.TempDir(), "healthy.db")
			target, err := sql.Open("sqlite", path)
			require.NoError(t, err)
			t.Cleanup(func() { _ = target.Close() })
			_, err = target.ExecContext(t.Context(), `CREATE TABLE snapshot_values (d0 REAL)`)
			require.NoError(t, err)
			bad := createSnapshotConnector(t, connectors, schema.DatabaseConnectorKindPostgres, `{"host":"127.0.0.1","database":"fixture"}`, 15)
			good := createSnapshotConnector(t, connectors, schema.DatabaseConnectorKindSQLite, fmt.Sprintf(`{"dsn":%q}`, path), 15)
			if grouped {
				createDecoupleGroupMapping(t, mappings, bad.ID, "hung", "d0")
			} else {
				createSnapshotMapping(t, mappings, bad.ID, "hung", "d0")
			}
			createSnapshotMapping(t, mappings, good.ID, "healthy", "d0")

			base := time.Date(2026, 11, 5, 10, 0, 0, 0, time.UTC)
			clock := &backoffTestClock{now: base}
			tick := make(chan time.Time, 1)
			entered, release := make(chan struct{}), make(chan struct{})
			releaseOpen := sync.OnceFunc(func() { close(release) })
			var opens atomic.Int64
			writer := NewWriterWithConfig(connectors, mappings, WriterConfig{
				FlushTick: tick, Now: clock.Now,
				OpenDestination: func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
					if kind == schema.DatabaseConnectorKindPostgres {
						if opens.Add(1) == 1 {
							close(entered)
						}
						<-release
						return nil, errors.New("fixture released")
					}
					return openExternalDBManager(kind, config)
				},
			})
			t.Cleanup(func() { releaseOpen(); _ = writer.Close(context.Background()) })
			require.ErrorIs(t, writer.WriteTagValue(t.Context(), "hung", 10, base.Add(time.Second)), ErrSnapshotAccepted)
			require.ErrorIs(t, writer.WriteTagValue(t.Context(), "healthy", 12, base.Add(time.Second)), ErrSnapshotAccepted)
			clock.Set(base.Add(15 * time.Second))
			tick <- clock.Now()
			select {
			case <-entered:
			case <-time.After(3 * time.Second):
				t.Fatal("the unreachable connector was not attempted")
			}

			done := make(chan error, 1)
			go func() { done <- writer.WriteTagValue(t.Context(), "hung", 11, base.Add(16*time.Second)) }()
			select {
			case err := <-done:
				require.ErrorIs(t, err, ErrSnapshotAccepted)
			case <-time.After(50 * time.Millisecond):
				t.Fatal("ingest waited for the unreachable connector")
			}
			require.Eventually(t, func() bool {
				var value float64
				return target.QueryRowContext(t.Context(), `SELECT d0 FROM snapshot_values`).Scan(&value) == nil && value == 12
			}, 3*time.Second, 10*time.Millisecond, "the healthy connector must keep writing")

			historyService := NewConnectorService(connectors, mappings)
			failedHistory := func() []WriteHistoryRecord {
				records, listErr := historyService.ListWriteHistory(t.Context(), bad.ID, 200)
				if listErr != nil {
					return nil
				}
				var failures []WriteHistoryRecord
				for _, record := range records {
					if record.Status == "failed" {
						failures = append(failures, record)
					}
				}
				return failures
			}
			require.Eventually(t, func() bool {
				updated, getErr := connectors.GetByID(t.Context(), bad.ID)
				return getErr == nil && updated.LastWriteStatus == "failed" && len(failedHistory()) == 1
			}, 7*time.Second, 10*time.Millisecond, "a connection attempt that never returns must be reported as failed")
			failures := failedHistory()
			require.Equal(t, 1, failures[0].RetryAfterSeconds)
			require.NotNil(t, failures[0].NextRetryAt)
			if grouped {
				updated, getErr := connectors.GetByID(t.Context(), bad.ID)
				require.NoError(t, getErr)
				require.Equal(t, "failed", updated.LastFlushStatus)
			}

			for i := range 3 {
				clock.Set(base.Add(time.Duration(100+100*i) * time.Second))
				tick <- clock.Now()
				require.Eventually(t, func() bool { return len(failedHistory()) >= i+2 }, 3*time.Second, 10*time.Millisecond)
			}
			require.EqualValues(t, 1, opens.Load(), "retries must not start more connection attempts while the first is stuck")
			closeCtx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
			defer cancel()
			closeErr := writer.Close(closeCtx)
			require.Error(t, closeErr, "the pending snapshot must not be reported as delivered")
			require.NotErrorIs(t, closeErr, context.DeadlineExceeded, "a stuck connector must not prevent bounded shutdown")
		})
	}
}

func TestSnapshotWorker_LateConnectionIsClosedWithoutWriting(t *testing.T) {
	connectors, mappings, connector, target := backoffTestTarget(t)
	createSnapshotMapping(t, mappings, connector.ID, "late-open", "d0")
	base := time.Date(2026, 11, 5, 10, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	tick := make(chan time.Time, 1)
	entered, release := make(chan struct{}), make(chan struct{})
	releaseOpen := sync.OnceFunc(func() { close(release) })
	opened := make(chan *datalinkbase.DBManager, 1)
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{
		FlushTick: tick, Now: clock.Now,
		OpenDestination: func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
			close(entered)
			<-release
			manager, err := openExternalDBManager(kind, config)
			opened <- manager
			return manager, err
		},
	})
	t.Cleanup(func() { releaseOpen(); _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "late-open", 10, base.Add(time.Second)), ErrSnapshotAccepted)
	clock.Set(base.Add(15 * time.Second))
	tick <- clock.Now()
	select {
	case <-entered:
	case <-time.After(3 * time.Second):
		t.Fatal("the connector was not attempted")
	}
	require.Eventually(t, func() bool {
		updated, err := connectors.GetByID(t.Context(), connector.ID)
		return err == nil && updated.LastWriteStatus == "failed"
	}, 7*time.Second, 10*time.Millisecond, "the stuck opener must time out before it returns")
	releaseOpen()
	var manager *datalinkbase.DBManager
	select {
	case manager = <-opened:
		require.NotNil(t, manager)
	case <-time.After(3 * time.Second):
		t.Fatal("the delayed connector did not return")
	}
	require.Eventually(t, func() bool { return manager.DB() == nil }, 3*time.Second, 10*time.Millisecond, "a late connection must be closed")
	var count int
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values`).Scan(&count))
	require.Zero(t, count, "a timed-out connection must not write after it returns")
}

func TestSnapshotWorker_CanceledFlushDoesNotAttemptOrRecordWrite(t *testing.T) {
	connectors, mappings, connector, _ := backoffTestTarget(t)
	createSnapshotMapping(t, mappings, connector.ID, "canceled", "d0")
	resetWriteHistory(connector.ID)
	var opens atomic.Int64
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{
		FlushTick: make(chan time.Time),
		OpenDestination: func(schema.DatabaseConnectorKind, ConnectionConfig) (*datalinkbase.DBManager, error) {
			opens.Add(1)
			return nil, errors.New("unexpected connection")
		},
	})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	at := time.Date(2026, 11, 5, 10, 0, 1, 0, time.UTC)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "canceled", 1, at), ErrSnapshotAccepted)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	require.ErrorIs(t, writer.flushDue(ctx, at.Add(time.Minute)), context.Canceled)
	require.Zero(t, opens.Load())
	require.Empty(t, listWriteHistory(connector.ID, 10))
	_, pending := writer.connectorsWithPendingData()[connector.ID]
	require.True(t, pending)
}
