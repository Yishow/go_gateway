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

type switchableOutcomeRepository struct {
	ConnectorRepository
	fail atomic.Bool
}

func (r *switchableOutcomeRepository) UpdateWriteOutcome(ctx context.Context, id string, at time.Time, status, summary string) error {
	if r.fail.Load() {
		return fmt.Errorf("outcome validation failed: %w", ErrValidation)
	}
	return r.ConnectorRepository.UpdateWriteOutcome(ctx, id, at, status, summary)
}

func TestWriter_JoinedOutcomeFailureChangesLogClass(t *testing.T) {
	connectors, mappings, connector, _ := backoffTestTarget(t)
	createSnapshotMapping(t, mappings, connector.ID, "outcome-class-log", "d0")
	base := time.Date(2026, 11, 7, 10, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	capture := captureWriterLogs(t)
	signals := make(chan struct{}, 10)
	outcomes := &switchableOutcomeRepository{ConnectorRepository: connectors}
	writer := NewWriterWithConfig(outcomes, mappings, WriterConfig{
		FlushTick: make(chan time.Time), Now: clock.Now,
		SuppressMapping: func(_, _ string) bool { signals <- struct{}{}; return false },
		OpenDestination: func(schema.DatabaseConnectorKind, ConnectionConfig) (*datalinkbase.DBManager, error) {
			return nil, context.DeadlineExceeded
		},
	})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "outcome-class-log", 12, base.Add(time.Second)), ErrSnapshotAccepted)
	worker := writer.workerFor(connector.ID)
	first := base.Add(15 * time.Second)
	advanceLogTestWorker(t, worker, clock, signals, first)
	require.Len(t, capture.forConnector(connector.ID), 1)

	outcomes.fail.Store(true)
	advanceLogTestWorker(t, worker, clock, signals, first.Add(time.Minute))
	logs := capture.forConnector(connector.ID)
	require.Len(t, logs, 2, "the additional outcome error must not be hidden by the existing timeout")
	require.Equal(t, slog.LevelError, logs[1].level)
	require.NotContains(t, logs[1].attrs, "error", "raw outcome errors must not be logged")
	advanceLogTestWorker(t, worker, clock, signals, first.Add(2*time.Minute))
	require.Len(t, capture.forConnector(connector.ID), 2)

	outcomes.fail.Store(false)
	advanceLogTestWorker(t, worker, clock, signals, first.Add(3*time.Minute))
	require.Len(t, capture.forConnector(connector.ID), 3)
}
