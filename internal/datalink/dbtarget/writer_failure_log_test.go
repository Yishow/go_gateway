package dbtarget

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type capturedWriterLog struct {
	level slog.Level
	attrs map[string]slog.Value
}

type writerLogCapture struct {
	mu      sync.Mutex
	records []capturedWriterLog
}

func (*writerLogCapture) Enabled(context.Context, slog.Level) bool { return true }

func (h *writerLogCapture) Handle(_ context.Context, record slog.Record) error {
	entry := capturedWriterLog{level: record.Level, attrs: make(map[string]slog.Value)}
	record.Attrs(func(attr slog.Attr) bool {
		entry.attrs[attr.Key] = attr.Value.Resolve()
		return true
	})
	h.mu.Lock()
	defer h.mu.Unlock()
	h.records = append(h.records, entry)
	return nil
}

func (h *writerLogCapture) WithAttrs([]slog.Attr) slog.Handler { return h }
func (h *writerLogCapture) WithGroup(string) slog.Handler      { return h }

func (h *writerLogCapture) forConnector(id string) []capturedWriterLog {
	h.mu.Lock()
	defer h.mu.Unlock()
	var records []capturedWriterLog
	for _, record := range h.records {
		if record.attrs["connector_id"].String() == id {
			records = append(records, record)
		}
	}
	return records
}

func captureWriterLogs(t *testing.T) *writerLogCapture {
	t.Helper()
	capture := &writerLogCapture{}
	previous := slog.Default()
	slog.SetDefault(slog.New(capture))
	t.Cleanup(func() { slog.SetDefault(previous) })
	return capture
}

func advanceLogTestWorker(t *testing.T, worker *connectorWorker, clock *backoffTestClock, signals <-chan struct{}, at time.Time) {
	t.Helper()
	clock.Set(at)
	worker.ticks <- at
	// A due write and the following idle accounting each invoke this hook;
	// the second signal means the worker has finished emitting flush logs.
	for range 2 {
		select {
		case <-signals:
		case <-time.After(3 * time.Second):
			t.Fatal("worker did not finish the flush")
		}
	}
}

func TestWriter_RepeatedFlushFailuresAreSummarized(t *testing.T) {
	connectors, mappings, connector, _ := backoffTestTarget(t)
	createSnapshotMapping(t, mappings, connector.ID, "failure-log", "d0")
	base := time.Date(2026, 11, 7, 10, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	capture := captureWriterLogs(t)

	var attempts atomic.Int64
	idleSignals := make(chan struct{}, 32)
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{
		FlushTick: make(chan time.Time), Now: clock.Now,
		SuppressMapping: func(_, _ string) bool { idleSignals <- struct{}{}; return false },
		OpenDestination: func(schema.DatabaseConnectorKind, ConnectionConfig) (*datalinkbase.DBManager, error) {
			return nil, fmt.Errorf("attempt %d: %w", attempts.Add(1), context.DeadlineExceeded)
		},
	})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "failure-log", 12, base.Add(time.Second)), ErrSnapshotAccepted)
	worker := writer.workerFor(connector.ID)
	first := base.Add(15 * time.Second)
	for i := range 21 {
		at := first.Add(time.Duration(i) * time.Minute)
		advanceLogTestWorker(t, worker, clock, idleSignals, at)
		logs := capture.forConnector(connector.ID)
		switch {
		case i < 10:
			require.Len(t, logs, 1, "repeated errors before ten minutes must be suppressed")
			require.Equal(t, slog.LevelError, logs[0].level)
			require.Equal(t, "timeout", logs[0].attrs["failure_reason"].String())
			require.EqualValues(t, 1, logs[0].attrs["failure_count"].Any())
			require.Equal(t, first, logs[0].attrs["first_failure_at"].Any())
			require.NotContains(t, logs[0].attrs, "error")
		case i < 20:
			require.Len(t, logs, 2)
			require.Equal(t, slog.LevelWarn, logs[1].level)
			require.EqualValues(t, 11, logs[1].attrs["failure_count"].Int64())
			require.True(t, first.Equal(logs[1].attrs["first_failure_at"].Time()))
		default:
			require.Len(t, logs, 3)
			require.Equal(t, slog.LevelWarn, logs[2].level)
			require.EqualValues(t, 21, logs[2].attrs["failure_count"].Int64())
			require.True(t, first.Equal(logs[2].attrs["first_failure_at"].Time()))
		}
	}
}

func TestWriter_RecoveryLogRequiresSuccessfulWrite(t *testing.T) {
	for _, grouped := range []bool{false, true} {
		t.Run(fmt.Sprintf("grouped=%t", grouped), func(t *testing.T) {
			connectors, mappings, connector, target := backoffTestTarget(t)
			if grouped {
				createDecoupleGroupMapping(t, mappings, connector.ID, "recovery-log", "d0")
			} else {
				createSnapshotMapping(t, mappings, connector.ID, "recovery-log", "d0")
			}
			base := time.Date(2026, 11, 7, 10, 0, 0, 0, time.UTC)
			clock := &backoffTestClock{now: base}
			capture := captureWriterLogs(t)
			signals := make(chan struct{}, 8)
			var fail, skip atomic.Bool
			fail.Store(true)
			writer := NewWriterWithConfig(connectors, mappings, WriterConfig{
				FlushTick: make(chan time.Time), Now: clock.Now,
				SuppressMapping: func(_, _ string) bool { signals <- struct{}{}; return skip.Load() },
				OpenDestination: func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
					if fail.Load() {
						return nil, context.DeadlineExceeded
					}
					return openExternalDBManagerFunc(kind, config)
				},
			})
			t.Cleanup(func() { _ = writer.Close(context.Background()) })
			require.ErrorIs(t, writer.WriteTagValue(t.Context(), "recovery-log", 12, base.Add(time.Second)), ErrSnapshotAccepted)
			worker := writer.workerFor(connector.ID)
			first := base.Add(15 * time.Second)
			advanceLogTestWorker(t, worker, clock, signals, first)
			require.Len(t, capture.forConnector(connector.ID), 1)

			fail.Store(false)
			skip.Store(true)
			advanceLogTestWorker(t, worker, clock, signals, first.Add(time.Minute))
			require.Len(t, capture.forConnector(connector.ID), 1, "a skipped snapshot is not recovery")

			skip.Store(false)
			clock.Set(first.Add(2 * time.Minute))
			worker.ticks <- clock.Now()
			select {
			case <-signals:
			case <-time.After(3 * time.Second):
				t.Fatal("worker did not finish the idle tick")
			}
			require.Len(t, capture.forConnector(connector.ID), 1, "an empty flush is not recovery")
			require.ErrorIs(t, writer.WriteTagValue(t.Context(), "recovery-log", 13, first.Add(3*time.Minute)), ErrSnapshotAccepted)
			advanceLogTestWorker(t, worker, clock, signals, first.Add(4*time.Minute))
			logs := capture.forConnector(connector.ID)
			require.Len(t, logs, 2)
			require.Equal(t, slog.LevelInfo, logs[1].level)
			require.EqualValues(t, 1, logs[1].attrs["failure_count"].Int64())
			var value float64
			require.NoError(t, target.QueryRowContext(t.Context(), `SELECT d0 FROM snapshot_values`).Scan(&value))
			require.Equal(t, 13.0, value)
		})
	}
}

func TestWriter_FailureClassChangesProduceNewErrors(t *testing.T) {
	connectors, mappings, connector, target := backoffTestTarget(t)
	createSnapshotMapping(t, mappings, connector.ID, "class-log", "d0")
	base := time.Date(2026, 11, 7, 10, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	capture := captureWriterLogs(t)
	signals := make(chan struct{}, 16)
	var attempts atomic.Int64
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{
		FlushTick: make(chan time.Time), Now: clock.Now,
		SuppressMapping: func(_, _ string) bool { signals <- struct{}{}; return false },
		OpenDestination: func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
			switch attempt := attempts.Add(1); attempt {
			case 1, 2:
				return nil, fmt.Errorf("timeout attempt %d: %w", attempt, context.DeadlineExceeded)
			case 3, 4:
				return nil, &net.OpError{Op: fmt.Sprintf("dial-%d", attempt), Net: "tcp", Err: syscall.ECONNREFUSED}
			case 5, 6:
				return nil, fmt.Errorf("validation attempt %d: %w", attempt, ErrValidation)
			default:
				return openExternalDBManagerFunc(kind, config)
			}
		},
	})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "class-log", 12, base.Add(time.Second)), ErrSnapshotAccepted)
	worker := writer.workerFor(connector.ID)
	first := base.Add(15 * time.Second)
	for i, want := range []int{1, 1, 2, 2, 3, 3, 4} {
		advanceLogTestWorker(t, worker, clock, signals, first.Add(time.Duration(i)*time.Minute))
		logs := capture.forConnector(connector.ID)
		require.Len(t, logs, want)
		if i == 2 || i == 4 {
			require.Equal(t, slog.LevelError, logs[len(logs)-1].level)
		}
		if i == 6 {
			require.Equal(t, slog.LevelInfo, logs[3].level)
			require.EqualValues(t, 6, logs[3].attrs["failure_count"].Int64())
		}
	}
	var value float64
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT d0 FROM snapshot_values`).Scan(&value))
	require.Equal(t, 12.0, value)
}

func TestWriter_FailureLogsAreIsolatedPerConnector(t *testing.T) {
	connectors, mappings, firstConnector, _ := backoffTestTarget(t)
	secondConnector := *firstConnector
	secondConnector.ID, secondConnector.Name = "another-connector", "another-destination"
	require.NoError(t, connectors.Create(t.Context(), &secondConnector))
	createSnapshotMapping(t, mappings, firstConnector.ID, "first-log", "d0")
	createSnapshotMapping(t, mappings, secondConnector.ID, "second-log", "d1")
	base := time.Date(2026, 11, 7, 10, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	capture := captureWriterLogs(t)
	signals := make(chan struct{}, 16)
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{
		FlushTick: make(chan time.Time), Now: clock.Now,
		SuppressMapping: func(_, _ string) bool { signals <- struct{}{}; return false },
		OpenDestination: func(schema.DatabaseConnectorKind, ConnectionConfig) (*datalinkbase.DBManager, error) {
			return nil, context.DeadlineExceeded
		},
	})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "first-log", 1, base.Add(time.Second)), ErrSnapshotAccepted)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "second-log", 2, base.Add(time.Second)), ErrSnapshotAccepted)
	firstWorker, secondWorker := writer.workerFor(firstConnector.ID), writer.workerFor(secondConnector.ID)
	first := base.Add(15 * time.Second)
	advanceLogTestWorker(t, firstWorker, clock, signals, first)
	advanceLogTestWorker(t, firstWorker, clock, signals, first.Add(time.Minute))
	advanceLogTestWorker(t, secondWorker, clock, signals, first.Add(2*time.Minute))
	require.Len(t, capture.forConnector(firstConnector.ID), 1)
	secondLogs := capture.forConnector(secondConnector.ID)
	require.Len(t, secondLogs, 1)
	require.Equal(t, slog.LevelError, secondLogs[0].level)
}
