package collector

import (
	"context"
	"errors"
	"slices"
	"sync"
)

var errActiveCapacityWait = errors.New("execution capacity wait")

// activeCapacity serializes fair admission only, never device ownership or I/O.
// Scheduled waiters are bounded by one dispatched scan per device. Explicit
// manual callers share this pool without creating an internal pending queue.
type activeCapacity struct {
	mu            sync.Mutex
	limit, active int
	waiters       []*capacityWaiter
}
type capacityWaiter struct {
	ctx     context.Context
	ready   chan struct{}
	granted bool
}
type activeLease struct {
	capacity *activeCapacity
	held     bool
}

func newActiveCapacity(limit int) *activeCapacity { return &activeCapacity{limit: limit} }

func (c *activeCapacity) acquire(ctx context.Context) (*activeLease, error) {
	if err := c.take(ctx); err != nil {
		return nil, err
	}
	return &activeLease{capacity: c, held: true}, nil
}

func (c *activeCapacity) take(ctx context.Context) error {
	c.mu.Lock()
	if err := ctx.Err(); err != nil {
		c.mu.Unlock()
		return err
	}
	if c.active < c.limit && len(c.waiters) == 0 {
		c.active++
		c.mu.Unlock()
		return nil
	}
	waiter := &capacityWaiter{ctx: ctx, ready: make(chan struct{})}
	c.waiters = append(c.waiters, waiter)
	c.mu.Unlock()
	select {
	case <-waiter.ready:
	case <-ctx.Done():
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := ctx.Err(); err != nil {
		if waiter.granted {
			c.active--
		} else if index := slices.Index(c.waiters, waiter); index >= 0 {
			c.waiters = slices.Delete(c.waiters, index, index+1)
		}
		c.grantLocked()
		return err
	}
	return nil
}

func (c *activeCapacity) grantLocked() {
	for c.active < c.limit && len(c.waiters) > 0 {
		waiter := c.waiters[0]
		c.waiters[0] = nil
		c.waiters = c.waiters[1:]
		if waiter.ctx.Err() != nil {
			continue
		}
		c.active++
		waiter.granted = true
		close(waiter.ready)
	}
}

func (l *activeLease) release() {
	if !l.held {
		return
	}
	l.held = false
	l.capacity.mu.Lock()
	l.capacity.active--
	l.capacity.grantLocked()
	l.capacity.mu.Unlock()
}

func (l *activeLease) Suspend() { l.release() }
func (l *activeLease) Resume(ctx context.Context) error {
	if l.held {
		return ctx.Err()
	}
	if err := l.capacity.take(ctx); err != nil {
		return &activeCapacityWaitError{cause: err}
	}
	l.held = true
	return nil
}

// A private marker separates admission waiting from an invoked connector's
// existing dial/transport/validation failure. A shared budget may contain other
// requests' attempts, so its absolute Used count is not health authority.
type activeCapacityWaitError struct{ cause error }

func (e *activeCapacityWaitError) Error() string {
	return errActiveCapacityWait.Error() + ": " + e.cause.Error()
}
func (e *activeCapacityWaitError) Unwrap() error        { return e.cause }
func (e *activeCapacityWaitError) Is(target error) bool { return target == errActiveCapacityWait }

func capacityWaitResult(err error, invoked bool) error {
	var wait *activeCapacityWaitError
	if invoked && errors.As(err, &wait) {
		return wait.cause
	}
	return err
}
