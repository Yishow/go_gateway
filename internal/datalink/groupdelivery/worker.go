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

// forcedStopGrace bounds how long Stop waits after canceling in-flight work.
const forcedStopGrace = 2 * time.Second

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
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.running {
		return ErrWorkerRunning
	}
	if w.config.NodeID != "" {
		if _, err := w.store.RecoverStaleClaims(ctx, w.config.NodeID, w.config.Owner, w.dispatcher.config.Now()); err != nil {
			return err
		}
	}
	graceCtx, graceCancel := context.WithCancel(ctx)
	workCtx, workCancel := context.WithCancel(context.WithoutCancel(ctx))
	done := make(chan struct{})
	w.running, w.graceCancel, w.workCancel, w.done = true, graceCancel, workCancel, done
	go w.loop(workCtx, graceCtx, workCancel, done)
	return nil
}

func (w *Worker) loop(workCtx, graceCtx context.Context, workCancel context.CancelFunc, done chan struct{}) {
	defer func() {
		// The loop can end without Stop (its start context was canceled); the
		// worker must then be restartable and its work context released.
		w.mu.Lock()
		if w.done == done {
			w.running = false
		}
		w.mu.Unlock()
		workCancel()
		close(done)
	}()
	ticker := time.NewTicker(w.config.Interval)
	defer ticker.Stop()
	var lastReclaim time.Time
	for {
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

// Stop stops starting new rows, waits up to deadline for in-flight deliveries
// to finish, then cancels them. Rows that were interrupted stay recoverable:
// they are settled by the next incarnation's startup recovery. Stopping a
// worker that is not running does nothing.
func (w *Worker) Stop(deadline time.Duration) error {
	w.mu.Lock()
	if !w.running {
		w.mu.Unlock()
		return nil
	}
	graceCancel, workCancel, done := w.graceCancel, w.workCancel, w.done
	w.running = false
	w.mu.Unlock()

	graceCancel()
	timer := time.NewTimer(deadline)
	defer timer.Stop()
	select {
	case <-done:
		workCancel()
		return nil
	case <-timer.C:
	}
	workCancel()
	select {
	case <-done:
	case <-time.After(forcedStopGrace):
	}
	return ErrShutdownForced
}
