package storage

import (
	"context"
	"errors"
	"io"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type blockedBatchWriter struct {
	Writer
	entered chan context.Context
	release chan struct{}
	calls   atomic.Int32
	closed  atomic.Int32
	failure error
}

func (w *blockedBatchWriter) WriteBatch(ctx context.Context, _ []TimeSeriesRecord) error {
	w.calls.Add(1)
	if w.entered != nil {
		w.entered <- ctx
		<-w.release
	}
	return w.failure
}
func (w *blockedBatchWriter) Close() error { w.closed.Add(1); return nil }

func TestBatchWriterShutdownTimeoutKeepsTimerStorageOwned(t *testing.T) {
	underlying := &blockedBatchWriter{entered: make(chan context.Context, 1), release: make(chan struct{})}
	release := sync.OnceFunc(func() { close(underlying.release) })
	defer release()
	bw := NewBatchWriter(underlying, BatchWriterConfig{BatchSize: 10, FlushInterval: time.Millisecond, WriteTimeout: time.Minute})
	require.NoError(t, bw.Write(t.Context(), TimeSeriesRecord{TagID: "accepted"}))
	var timerCtx context.Context
	select {
	case timerCtx = <-underlying.entered:
	case <-time.After(time.Second):
		t.Fatal("timer did not enter storage")
	}
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Millisecond)
	defer cancel()
	require.ErrorIs(t, bw.CloseContext(ctx), context.DeadlineExceeded)
	require.ErrorIs(t, timerCtx.Err(), context.Canceled, "Close cancels the owned timer operation")
	require.ErrorIs(t, bw.CloseContext(ctx), context.DeadlineExceeded)
	require.ErrorIs(t, bw.WaitClosed(ctx), context.DeadlineExceeded)
	require.Zero(t, underlying.closed.Load(), "must not close storage under a live timer")
	require.ErrorIs(t, bw.Write(t.Context(), TimeSeriesRecord{}), io.ErrClosedPipe)
	release()
	require.NoError(t, bw.WaitClosed(t.Context()))
	require.NoError(t, bw.Close())
	require.EqualValues(t, 1, underlying.closed.Load())
	require.EqualValues(t, 1, underlying.calls.Load())
}

func TestBatchWriterShutdownConcurrentCloseFlushesOnceAndPreservesError(t *testing.T) {
	failure := errors.New("test flush failed")
	underlying := &blockedBatchWriter{failure: failure}
	bw := NewBatchWriter(underlying, BatchWriterConfig{BatchSize: 10, FlushInterval: time.Hour, WriteTimeout: time.Second})
	require.NoError(t, bw.Write(t.Context(), TimeSeriesRecord{TagID: "buffered"}))
	results := make(chan error, 20)
	var wg sync.WaitGroup
	for range cap(results) {
		wg.Go(func() { results <- bw.CloseContext(t.Context()) })
	}
	wg.Wait()
	close(results)
	for err := range results {
		require.ErrorIs(t, err, failure)
	}
	require.ErrorIs(t, bw.WaitClosed(t.Context()), failure)
	require.EqualValues(t, 1, underlying.calls.Load())
	require.EqualValues(t, 1, underlying.closed.Load())
}

func TestBatchWriterShutdownDoesNotCloseUnderSynchronousWrite(t *testing.T) {
	underlying := &blockedBatchWriter{entered: make(chan context.Context, 1), release: make(chan struct{})}
	release := sync.OnceFunc(func() { close(underlying.release) })
	defer release()
	bw := NewBatchWriter(underlying, BatchWriterConfig{BatchSize: 1, FlushInterval: time.Hour, WriteTimeout: time.Second})
	written := make(chan error, 1)
	go func() { written <- bw.Write(t.Context(), TimeSeriesRecord{TagID: "in-flight"}) }()
	select {
	case <-underlying.entered:
	case <-time.After(time.Second):
		t.Fatal("write did not enter storage")
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	require.ErrorIs(t, bw.CloseContext(ctx), context.Canceled)
	require.Zero(t, underlying.closed.Load())
	release()
	require.NoError(t, <-written)
	require.NoError(t, bw.WaitClosed(t.Context()))
	require.EqualValues(t, 1, underlying.closed.Load())
}

// ctxRecordingWriter fails like SQLite when its context ends and records what it stored.
type ctxRecordingWriter struct {
	Writer
	block   chan struct{}
	entered chan struct{}
	mu      sync.Mutex
	stored  []string
}

func (w *ctxRecordingWriter) WriteBatch(ctx context.Context, records []TimeSeriesRecord) error {
	if w.block != nil {
		w.entered <- struct{}{}
		w.block = nil
		<-ctx.Done()
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, r := range records {
		w.stored = append(w.stored, r.TagID)
	}
	return nil
}
func (w *ctxRecordingWriter) Close() error { return nil }

func TestBatchWriterShutdownRetainsRecordsFromCanceledTimerFlush(t *testing.T) {
	underlying := &ctxRecordingWriter{block: make(chan struct{}), entered: make(chan struct{}, 1)}
	bw := NewBatchWriter(underlying, BatchWriterConfig{BatchSize: 10, FlushInterval: time.Millisecond, WriteTimeout: time.Minute})
	require.NoError(t, bw.Write(t.Context(), TimeSeriesRecord{TagID: "accepted"}))
	select {
	case <-underlying.entered:
	case <-time.After(time.Second):
		t.Fatal("timer did not enter storage")
	}
	require.NoError(t, bw.CloseContext(t.Context()))
	require.Equal(t, []string{"accepted"}, underlying.stored, "canceled timer batch must reach the final flush")
}

func TestBatchWriterShutdownFinalFlushOutlivesExpiredCallerContext(t *testing.T) {
	underlying := &ctxRecordingWriter{}
	bw := NewBatchWriter(underlying, BatchWriterConfig{BatchSize: 10, FlushInterval: time.Hour, WriteTimeout: time.Second})
	require.NoError(t, bw.Write(t.Context(), TimeSeriesRecord{TagID: "buffered"}))
	expired, cancel := context.WithCancel(t.Context())
	cancel()
	_ = bw.CloseContext(expired)
	require.NoError(t, bw.WaitClosed(t.Context()))
	require.Equal(t, []string{"buffered"}, underlying.stored, "caller wait deadline must not abort the final flush")
}
