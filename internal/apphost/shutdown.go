package apphost

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ErrShutdownTimeout is an unmet notification deadline, never proof of completion.
var ErrShutdownTimeout = errors.New("shutdown notification deadline exceeded")

// Phase stops a resource once and separately confirms that original work ended.
type Phase struct {
	Name       string
	Stop, Wait func(context.Context) error
}

// ShutdownState contains safe phase metadata rather than arbitrary error strings.
type ShutdownState struct {
	Process, Phase string
	Failed         bool
}

// Shutdown owns one cancellation budget and one ordered completion operation.
type Shutdown struct {
	once         sync.Once
	mu, notifyMu sync.Mutex
	budget       time.Duration
	phases       []Phase
	onState      func(ShutdownState)
	state        ShutdownState
	done         chan struct{}
	err          error
	timedOut     bool
	stopping     atomic.Bool
}

// NewShutdown creates a coordinator without starting cleanup.
func NewShutdown(budget time.Duration, phases []Phase, onState func(ShutdownState)) *Shutdown {
	return &Shutdown{budget: budget, phases: phases, onState: onState, done: make(chan struct{})}
}

// Begin is safe for concurrent tray/signal/session requests.
func (s *Shutdown) Begin() { s.once.Do(func() { s.stopping.Store(true); go s.run() }) }

// IsStopping closes new admission before asynchronous phases start.
func (s *Shutdown) IsStopping() bool { return s.stopping.Load() }

// Done closes only when all resource completion has been observed.
func (s *Shutdown) Done() <-chan struct{} { return s.done }

// State returns the current safe process projection.
func (s *Shutdown) State() ShutdownState { s.mu.Lock(); defer s.mu.Unlock(); return s.state }

// Wait never starts another Stop or grants another shutdown budget.
func (s *Shutdown) Wait(ctx context.Context) error {
	select {
	case <-s.done:
		s.mu.Lock()
		defer s.mu.Unlock()
		return s.err
	default:
	}
	select {
	case <-s.done:
		s.mu.Lock()
		defer s.mu.Unlock()
		return s.err
	case <-ctx.Done():
		return ctx.Err()
	}
}
func (s *Shutdown) publish(phase string, failed bool) {
	s.notifyMu.Lock()
	defer s.notifyMu.Unlock()
	s.mu.Lock()
	process := "stopping"
	if s.timedOut {
		process = "stop-timeout"
	}
	s.state = ShutdownState{Process: process, Phase: phase, Failed: failed}
	state := s.state
	s.mu.Unlock()
	if s.onState != nil {
		s.onState(state)
	}
}
func (s *Shutdown) notifyDeadline(ctx context.Context, done chan<- struct{}) {
	defer close(done)
	<-ctx.Done()
	if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return
	}
	s.notifyMu.Lock()
	defer s.notifyMu.Unlock()
	s.mu.Lock()
	s.timedOut = true
	s.state.Process = "stop-timeout"
	state := s.state
	s.mu.Unlock()
	if s.onState != nil {
		s.onState(state)
	}
}
func (s *Shutdown) run() {
	ctx, cancel := context.WithTimeout(context.Background(), s.budget)
	defer cancel()
	noticeDone := make(chan struct{})
	go s.notifyDeadline(ctx, noticeDone)
	var failures []error
	for _, phase := range s.phases {
		s.publish(phase.Name, false)
		if phase.Stop != nil {
			if err := phase.Stop(ctx); err != nil {
				failures = append(failures, fmt.Errorf("%s: %w", phase.Name, err))
				s.publish(phase.Name, true)
			}
		}
		// A canceled observer is not a completed worker. Never release its storage.
		if phase.Wait != nil {
			if err := phase.Wait(context.WithoutCancel(ctx)); err != nil {
				failures = append(failures, fmt.Errorf("%s completion: %w", phase.Name, err))
				s.publish(phase.Name, true)
			}
		}
	}
	cancel()
	<-noticeDone
	s.notifyMu.Lock()
	defer s.notifyMu.Unlock()
	s.mu.Lock()
	if s.timedOut {
		failures = append(failures, ErrShutdownTimeout)
	}
	s.err = errors.Join(failures...)
	s.state = ShutdownState{Process: "stopped"}
	if s.err != nil {
		s.state.Process = "stop-failed"
		s.state.Failed = true
	}
	state := s.state
	s.mu.Unlock()
	if s.onState != nil {
		s.onState(state)
	}
	close(s.done)
}
