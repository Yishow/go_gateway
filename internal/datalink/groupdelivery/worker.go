package groupdelivery

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ErrShutdownForced means the worker did not drain by the deadline and its
// in-flight delivery was canceled; the row stays recoverable.
var ErrShutdownForced = errors.New("delivery worker shutdown hit its deadline")

// ErrWorkerRunning means Start was called on a running worker.
var ErrWorkerRunning = errors.New("delivery worker already running")

// WorkerConfig drives the delivery loop.
type WorkerConfig struct {
	// Interval is the pause between delivery cycles.
	Interval time.Duration
	// NodeID and Owner identify this incarnation ("<node>/<incarnation>") for
	// startup recovery of earlier incarnations' claims.
	NodeID string
	Owner  string
	// OnError receives cycle errors.
	OnError func(error)
	// ReclaimEvery and ReclaimRetention drive periodic reclamation of finished
	// data; zero ReclaimEvery disables it. Retention must be positive then.
	ReclaimEvery     time.Duration
	ReclaimRetention time.Duration
}

// Worker runs the dispatcher in the background with a bounded shutdown.
type Worker struct {
	store      *Store
	dispatcher *Dispatcher
	config     WorkerConfig

	mu          sync.Mutex
	running     bool
	stopping    bool
	graceCancel context.CancelFunc
	workCancel  context.CancelFunc
	done        chan struct{}
}

// NewWorker builds a worker around a dispatcher.
func NewWorker(store *Store, dispatcher *Dispatcher, config WorkerConfig) *Worker {
	if config.Interval <= 0 {
		config.Interval = time.Second
	}
	return &Worker{store: store, dispatcher: dispatcher, config: config}
}

// Start settles claims left by a previous incarnation or an expired lease, then
// begins delivering. Recovery runs before the first cycle so a restarted
// gateway never waits on, or races with, a worker that no longer exists.
func (w *Worker) Start(ctx context.Context) error {
	return w.StartWithLifetime(ctx, ctx)
}

// StartWithLifetime uses startup for cancellable recovery and lifetime for the
// delivery loop. Startup cancellation never starts new delivery after recovery;
// canceling a completed start request does not end the separately owned loop.
func (w *Worker) StartWithLifetime(startup, lifetime context.Context) error {
	w.mu.Lock()
	if w.running {
		w.mu.Unlock()
		return ErrWorkerRunning
	}
	graceCtx, graceCancel := context.WithCancel(lifetime)
	workCtx, workCancel := context.WithCancel(context.WithoutCancel(lifetime))
	done := make(chan struct{})
	w.running, w.stopping, w.graceCancel, w.workCancel, w.done = true, false, graceCancel, workCancel, done
	w.mu.Unlock()
	recoveryCtx, recoveryCancel := context.WithCancel(startup)
	stopRecovery := context.AfterFunc(graceCtx, recoveryCancel)
	defer stopRecovery()
	defer recoveryCancel()
	if w.config.NodeID != "" {
		if _, err := w.store.RecoverStaleClaims(recoveryCtx, w.config.NodeID, w.config.Owner, w.dispatcher.config.Now()); err != nil {
			graceCancel()
			w.finish(workCancel, done)
			return err
		}
	}
	if err := recoveryCtx.Err(); err != nil {
		w.finish(workCancel, done)
		return err
	}
	go w.loop(workCtx, graceCtx, workCancel, done)
	return nil
}

func (w *Worker) finish(workCancel context.CancelFunc, done chan struct{}) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.done == done {
		w.graceCancel()
		w.running = false
	}
	workCancel()
	close(done)
}

func (w *Worker) loop(workCtx, graceCtx context.Context, workCancel context.CancelFunc, done chan struct{}) {
	defer w.finish(workCancel, done)
	ticker := time.NewTicker(w.config.Interval)
	defer ticker.Stop()
	var lastReclaim time.Time
	for {
		if graceCtx.Err() != nil {
			return
		}
		w.sweep(workCtx)
		if w.config.ReclaimEvery > 0 && w.config.ReclaimRetention > 0 && time.Since(lastReclaim) >= w.config.ReclaimEvery {
			lastReclaim = time.Now()
			if _, err := w.store.Reclaim(workCtx, w.config.ReclaimRetention); err != nil && w.config.OnError != nil {
				w.config.OnError(err)
			}
		}
		if _, err := w.dispatcher.runCycle(workCtx, graceCtx); err != nil && w.config.OnError != nil {
			w.config.OnError(err)
		}
		select {
		case <-graceCtx.Done():
			return
		case <-ticker.C:
		}
	}
}

// sweep settles claims whose lease has expired while the process keeps running
// (a failed local state update, a stalled sibling worker). Without it such a row
// would hold its whole partition until the next restart.
func (w *Worker) sweep(ctx context.Context) {
	node := w.config.NodeID
	if node == "" {
		node = "\x00no-node" // matches no real owner prefix; only lease expiry applies
	}
	if _, err := w.store.RecoverStaleClaims(ctx, node, w.config.Owner, w.dispatcher.config.Now()); err != nil && w.config.OnError != nil {
		w.config.OnError(err)
	}
}

// Stop preserves the duration-based API using one total wait budget.
func (w *Worker) Stop(deadline time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()
	return w.StopContext(ctx)
}

// StopContext stops starting rows once, allows in-flight delivery until the
// original context expires, then cancels it. Returning a timeout does not mean
// workers finished; callers must observe WaitStopped before closing the store.
func (w *Worker) StopContext(ctx context.Context) error {
	w.mu.Lock()
	done := w.done
	if w.running && !w.stopping {
		w.stopping = true
		w.graceCancel()
		workCancel := w.workCancel
		go func() {
			select {
			case <-done:
			case <-ctx.Done():
				workCancel()
			}
		}()
	}
	w.mu.Unlock()
	if err := waitWorker(ctx, done); err != nil {
		return errors.Join(ErrShutdownForced, err)
	}
	return nil
}

// WaitStopped observes the current worker's completion without initiating or
// retrying shutdown. It is safe after a StopContext timeout.
func (w *Worker) WaitStopped(ctx context.Context) error {
	w.mu.Lock()
	done := w.done
	w.mu.Unlock()
	return waitWorker(ctx, done)
}

func waitWorker(ctx context.Context, done <-chan struct{}) error {
	if done == nil {
		return nil
	}
	select {
	case <-done:
		return nil
	default:
	}
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
