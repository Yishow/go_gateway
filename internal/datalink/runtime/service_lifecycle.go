package runtime

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Start starts runtime workers with a service-owned context. Request values
// survive request cancellation; Stop, rather than the request, owns their end.
func (s *Service) Start(ctx context.Context) error {
	s.lifecycle.Lock()
	if s.running.Load() || s.startDone != nil || s.stopDone != nil {
		s.lifecycle.Unlock()
		return fmt.Errorf("runtime is already started or stopping")
	}
	s.running.Store(true)
	s.startedAt.Store(time.Now().UnixNano())
	s.startDone = make(chan struct{})
	runtimeCtx, cancel := context.WithCancel(context.WithoutCancel(ctx))
	s.consumeCancel = cancel
	startCtx, startCancel := context.WithCancel(ctx)
	s.startCancel = startCancel
	s.lifecycle.Unlock()

	err := s.bootstrap(startCtx)
	if err == nil {
		err = startCtx.Err()
	}
	startCancel()
	s.lifecycle.Lock()
	defer s.lifecycle.Unlock()
	switch {
	case err != nil:
		s.running.Store(false)
		s.startedAt.Store(0)
		cancel()
	case s.stopDone == nil:
		s.wg.Go(func() { s.consumeLoop(runtimeCtx) })
		s.wg.Go(s.statusLoop) //nolint:contextcheck // The status worker is owned by stopCh.
	default:
		err = fmt.Errorf("runtime stopped during startup")
	}
	close(s.startDone)
	if err != nil && s.stopDone == nil {
		s.startDone = nil
	}
	return err
}

// Stop starts one ordered shutdown. A deadline only bounds this caller's wait;
// resources remain owned until all original workers actually finish. It also
// releases resources acquired by construction or a failed Start.
func (s *Service) Stop(ctx context.Context) error {
	s.lifecycle.Lock()
	if s.stopDone == nil {
		s.stopDone = make(chan struct{})
		s.running.Store(false)
		s.startedAt.Store(0)
		close(s.stopCh)
		if s.startCancel != nil {
			s.startCancel()
		}
		if s.consumeCancel != nil {
			s.consumeCancel()
		}
		go s.finishStop(ctx, s.startDone, s.stopDone)
	}
	done := s.stopDone
	s.lifecycle.Unlock()
	return s.waitStopped(ctx, done)
}

// WaitStopped observes completion of the original shutdown without re-running
// Stop or allocating another grace budget.
func (s *Service) WaitStopped(ctx context.Context) error {
	s.lifecycle.Lock()
	done := s.stopDone
	s.lifecycle.Unlock()
	if done == nil {
		return fmt.Errorf("runtime has not been stopped")
	}
	return s.waitStopped(ctx, done)
}

func (s *Service) waitStopped(ctx context.Context, done <-chan struct{}) error {
	select {
	case <-done:
		return s.stopErr
	default:
	}
	select {
	case <-done:
		return s.stopErr
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *Service) finishStop(ctx context.Context, started, done chan struct{}) {
	if started != nil {
		<-started
	}
	// A timed-out protocol read may still own its connection. Observe its actual
	// completion before allowing the caller to close shared connections or DBs.
	observation := context.WithoutCancel(ctx)
	stopErr := s.scheduler.StopContext(ctx)
	schedulerErr := s.scheduler.WaitStopped(observation)
	if stopErr != nil && !errors.Is(stopErr, ctx.Err()) {
		schedulerErr = errors.Join(schedulerErr, stopErr)
	}
	s.wg.Wait()
	var targetErr error
	if closer, ok := s.target.(interface{ Close(context.Context) error }); ok && closer != nil {
		targetErr = closer.Close(ctx)
	}
	var writerErr error
	if closer, ok := s.writer.(interface {
		CloseContext(context.Context) error
		WaitClosed(context.Context) error
	}); ok {
		closeErr := closer.CloseContext(ctx)
		writerErr = closer.WaitClosed(observation)
		if writerErr == nil && closeErr != nil && !errors.Is(closeErr, ctx.Err()) {
			writerErr = closeErr
		}
	} else {
		// A failed flush must not bypass Close and leak its resource. This is legacy
		// in-memory storage; this path does not claim durable recovery or delivery.
		writerErr = errors.Join(s.writer.Flush(ctx), s.writer.Close())
	}
	s.stopErr = errors.Join(schedulerErr, targetErr, writerErr)
	close(done)
}
