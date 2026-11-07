package dbtarget

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

type connectorWorker struct {
	ticks      chan time.Time
	failures   int
	retryAt    time.Time
	idle       map[string]idleIntervalState
	failureLog connectorFailureLog
}

// superviseLoop refreshes configuration outside ingest and signals independent
// workers without waiting for any destination to complete.
func (w *Writer) superviseLoop() {
	var latest time.Time
	for {
		select {
		case now, ok := <-w.flushTick:
			if !ok {
				return
			}
			if latest.IsZero() || now.After(latest) {
				latest = now
			}
			if w.mappingRepo != nil && w.connectorRepo != nil && w.refreshing.CompareAndSwap(false, true) {
				w.refreshVersion.Add(1)
				w.refreshReq <- struct{}{}
			}
			w.signalConnectorWorkers(latest)
		case <-w.refreshDone:
			if !latest.IsZero() {
				w.signalConnectorWorkers(latest)
			}
		case <-w.stopCh:
			return
		}
	}
}

func (w *Writer) refreshLoop() {
	for {
		select {
		case <-w.refreshReq:
			ctx, cancel := context.WithTimeout(w.baseCtx, 5*time.Second)
			err := w.RefreshMappings(ctx)
			cancel()
			w.refreshVersion.Add(1)
			w.refreshing.Store(false)
			if err != nil && w.baseCtx.Err() == nil {
				slog.Error("refresh database target mappings", "error", err)
			}
			select {
			case w.refreshDone <- struct{}{}:
			default:
			}
		case <-w.stopCh:
			return
		}
	}
}

func (w *Writer) signalConnectorWorkers(now time.Time) {
	ids := w.indexedConnectorIDs()
	for id := range w.connectorsWithPendingData() {
		ids[id] = struct{}{}
	}
	for id := range ids {
		worker := w.workerFor(id)
		select {
		case worker.ticks <- now:
		default:
			// A blocked connector needs only the most recent wakeup.
			select {
			case <-worker.ticks:
			default:
			}
			select {
			case worker.ticks <- now:
			default:
			}
		}
	}
}

func (w *Writer) workerFor(connectorID string) *connectorWorker {
	w.workersMu.Lock()
	defer w.workersMu.Unlock()
	if worker := w.workers[connectorID]; worker != nil {
		return worker
	}
	worker := &connectorWorker{ticks: make(chan time.Time, 1), idle: w.initialIdleIntervals(connectorID, w.now())}
	w.workers[connectorID] = worker
	w.wg.Go(func() { w.runConnectorWorker(connectorID, worker) })
	return worker
}

func (w *Writer) runConnectorWorker(connectorID string, worker *connectorWorker) {
	for {
		select {
		case now := <-worker.ticks:
			if !worker.retryAt.IsZero() && w.now().Before(worker.retryAt) {
				w.recordIdleIntervals(connectorID, worker, now)
				continue
			}
			wrote, err := w.flushConnector(w.workerCtx, connectorID, now, false, retryDelay(worker.failures))
			var retryErr *snapshotRetryError
			if errors.As(err, &retryErr) {
				worker.failureLog.selectRevision(retryErr.connectorRevision)
				worker.failures = min(worker.failures+1, 6)
				worker.retryAt = retryErr.retryAt
			} else {
				worker.failures = 0
				worker.retryAt = time.Time{}
			}
			if err != nil {
				worker.failureLog.failed(connectorID, err, w.now().UTC())
			} else if _, sameRevision := wrote[worker.failureLog.revision]; sameRevision && worker.failureLog.revision != "" {
				worker.failureLog.recovered(connectorID)
			}
			w.recordIdleIntervals(connectorID, worker, now)
		case <-w.stopCh:
			return
		}
	}
}

// Close stops and joins connector workers before flushing the remaining data.
// The caller's context bounds the wait; an uncooperative opener may outlive it.
func (w *Writer) Close(ctx context.Context) error {
	w.closeOnce.Do(func() {
		w.mu.Lock()
		w.closed = true
		w.mu.Unlock()
		close(w.stopCh)
		w.stopTick()
		w.cancelBase()
		go func() {
			w.wg.Wait()
			w.cancelWorkers()
			w.closeErr = w.flushClosing(ctx, w.now())
			close(w.closeDone)
		}()
	})
	select {
	case <-w.closeDone:
		return w.closeErr
	case <-ctx.Done():
		w.cancelWorkers()
		return ctx.Err()
	}
}

// WaitClosed observes the original close operation without starting new work.
func (w *Writer) WaitClosed(ctx context.Context) error {
	w.mu.Lock()
	closed := w.closed
	w.mu.Unlock()
	if !closed {
		return errors.New("database target writer has not been closed")
	}
	select {
	case <-w.closeDone:
		return w.closeErr
	default:
	}
	select {
	case <-w.closeDone:
		return w.closeErr
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (w *Writer) flushClosing(ctx context.Context, now time.Time) error {
	var failures []error
	for connectorID := range w.connectorsWithPendingData() {
		w.workersMu.Lock()
		worker := w.workers[connectorID]
		w.workersMu.Unlock()
		if worker != nil && now.Before(worker.retryAt) {
			failures = append(failures, fmt.Errorf("connector %s has a snapshot pending during backoff", connectorID))
			continue
		}
		wrote, err := w.flushConnector(ctx, connectorID, now, true, 0)
		if err != nil {
			failures = append(failures, err)
		} else if worker != nil {
			if _, sameRevision := wrote[worker.failureLog.revision]; sameRevision && worker.failureLog.revision != "" {
				worker.failureLog.recovered(connectorID)
			}
		}
	}
	return errors.Join(failures...)
}
