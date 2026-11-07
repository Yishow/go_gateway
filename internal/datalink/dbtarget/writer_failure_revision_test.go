package dbtarget

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestWriter_ReportsSameConnectionFailureForNewRevision(t *testing.T) {
	connectors, mappings, connector, _ := backoffTestTarget(t)
	createSnapshotMapping(t, mappings, connector.ID, "revision-failure", "d0")
	base := time.Date(2026, 11, 7, 10, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	capture := captureWriterLogs(t)
	signals := make(chan struct{}, 12)
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{
		FlushTick: make(chan time.Time), Now: clock.Now,
		SuppressMapping: func(_, _ string) bool { signals <- struct{}{}; return false },
		OpenDestination: func(schema.DatabaseConnectorKind, ConnectionConfig) (*datalinkbase.DBManager, error) {
			return nil, context.DeadlineExceeded
		},
	})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "revision-failure", 1, base.Add(time.Second)), ErrSnapshotAccepted)
	worker := writer.workerFor(connector.ID)
	first := base.Add(15 * time.Second)
	advanceLogTestWorker(t, worker, clock, signals, first)
	require.Len(t, capture.forConnector(connector.ID), 1)

	updated, err := connectors.GetByID(t.Context(), connector.ID)
	require.NoError(t, err)
	previous := updated.IdentityRevision
	updated.IdentityRevision = "new-revision"
	require.NoError(t, connectors.UpdateWithExpectedIdentityRevision(t.Context(), updated, previous))
	require.NoError(t, writer.RefreshMappings(t.Context()))
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "revision-failure", 2, first.Add(time.Minute)), ErrSnapshotAccepted)
	second := first.Add(2 * time.Minute)
	advanceLogTestWorker(t, worker, clock, signals, second)
	require.Len(t, capture.forConnector(connector.ID), 2, "the new destination's first failure must be reported")

	advanceLogTestWorker(t, worker, clock, signals, second.Add(10*time.Minute))
	logs := capture.forConnector(connector.ID)
	require.Len(t, logs, 3)
	require.EqualValues(t, 2, logs[2].attrs["failure_count"].Int64())
	require.True(t, second.Equal(logs[2].attrs["first_failure_at"].Time()))
}

func TestWriter_NewRevisionSuccessDoesNotRecoverOldDestination(t *testing.T) {
	connectors, mappings, connector, _ := backoffTestTarget(t)
	createSnapshotMapping(t, mappings, connector.ID, "revision-success", "d0")
	base := time.Date(2026, 11, 7, 10, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	capture := captureWriterLogs(t)
	signals := make(chan struct{}, 10)
	var fail atomic.Bool
	fail.Store(true)
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{
		FlushTick: make(chan time.Time), Now: clock.Now,
		SuppressMapping: func(_, _ string) bool { signals <- struct{}{}; return false },
		OpenDestination: func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
			if fail.Load() {
				return nil, context.DeadlineExceeded
			}
			return openExternalDBManagerFunc(kind, config)
		},
	})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "revision-success", 1, base.Add(time.Second)), ErrSnapshotAccepted)
	worker := writer.workerFor(connector.ID)
	first := base.Add(15 * time.Second)
	advanceLogTestWorker(t, worker, clock, signals, first)
	require.Len(t, capture.forConnector(connector.ID), 1)

	updated, err := connectors.GetByID(t.Context(), connector.ID)
	require.NoError(t, err)
	previous := updated.IdentityRevision
	updated.IdentityRevision = "new-revision"
	require.NoError(t, connectors.UpdateWithExpectedIdentityRevision(t.Context(), updated, previous))
	require.NoError(t, writer.RefreshMappings(t.Context()))
	fail.Store(false)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "revision-success", 2, first.Add(time.Minute)), ErrSnapshotAccepted)
	advanceLogTestWorker(t, worker, clock, signals, first.Add(2*time.Minute))
	require.Len(t, capture.forConnector(connector.ID), 1, "a new destination cannot recover an older destination's failure")
}

func TestWriter_ClosingSuccessfulWriteReportsRecovery(t *testing.T) {
	for _, grouped := range []bool{false, true} {
		t.Run(fmt.Sprintf("grouped=%t", grouped), func(t *testing.T) {
			testWriterClosingSuccessfulWriteReportsRecovery(t, grouped)
		})
	}
}

func testWriterClosingSuccessfulWriteReportsRecovery(t *testing.T, grouped bool) {
	t.Helper()
	connectors, mappings, connector, target := backoffTestTarget(t)
	if grouped {
		createDecoupleGroupMapping(t, mappings, connector.ID, "closing-recovery", "d0")
	} else {
		createSnapshotMapping(t, mappings, connector.ID, "closing-recovery", "d0")
	}
	base := time.Date(2026, 11, 7, 10, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	capture := captureWriterLogs(t)
	signals := make(chan struct{}, 4)
	var fail atomic.Bool
	fail.Store(true)
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{
		FlushTick: make(chan time.Time), Now: clock.Now,
		SuppressMapping: func(_, _ string) bool { signals <- struct{}{}; return false },
		OpenDestination: func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
			if fail.Load() {
				return nil, context.DeadlineExceeded
			}
			return openExternalDBManagerFunc(kind, config)
		},
	})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "closing-recovery", 12, base.Add(time.Second)), ErrSnapshotAccepted)
	worker := writer.workerFor(connector.ID)
	first := base.Add(15 * time.Second)
	advanceLogTestWorker(t, worker, clock, signals, first)
	fail.Store(false)
	clock.Set(first.Add(time.Minute))
	require.NoError(t, writer.Close(t.Context()))
	var value float64
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT d0 FROM snapshot_values`).Scan(&value))
	require.Equal(t, 12.0, value)
	logs := capture.forConnector(connector.ID)
	require.Len(t, logs, 2)
	require.Equal(t, slog.LevelInfo, logs[1].level)
}
