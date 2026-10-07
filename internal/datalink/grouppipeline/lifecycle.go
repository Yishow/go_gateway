package grouppipeline

import (
	"context"
	"errors"
	"time"
)

// Stop preserves the duration API with one deadline shared by all phases.
func (p *Pipeline) Stop(deadline time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()
	return p.StopContext(ctx)
}

// StopContext ends intake once, then the delivery worker. A timeout never
// releases ownership of a still-live worker or the durable store it uses.
func (p *Pipeline) StopContext(ctx context.Context) error {
	p.lifecycle.Lock()
	p.beginStopLocked(ctx)
	done := p.stopDone
	p.lifecycle.Unlock()
	return waitPipeline(ctx, done)
}

// beginStopLocked also rolls back a canceled startup after a worker launched.
// The caller holds lifecycle. An existing stop always keeps its original budget.
func (p *Pipeline) beginStopLocked(ctx context.Context) {
	if p.stopDone == nil {
		p.stopping = true
		p.stopDone = make(chan struct{})
		if p.startCancel != nil {
			p.startCancel()
		}
		if p.cancel != nil {
			p.cancel()
		}
		go p.finishStop(ctx, p.startDone, p.stopDone)
	}
}

// WaitStopped observes the original shutdown, including late completion after
// the notification deadline, without invoking Stop again.
func (p *Pipeline) WaitStopped(ctx context.Context) error {
	p.lifecycle.Lock()
	done := p.stopDone
	started := p.started
	p.lifecycle.Unlock()
	if done == nil && started {
		return errors.New("group pipeline has not been stopped")
	}
	return waitPipeline(ctx, done)
}

func waitPipeline(ctx context.Context, done <-chan struct{}) error {
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

func (p *Pipeline) finishStop(ctx context.Context, started, done chan struct{}) {
	if started != nil {
		<-started
	}
	p.lifecycle.Lock()
	loopDone, worker := p.loopDone, p.worker
	p.lifecycle.Unlock()
	if loopDone != nil {
		<-loopDone
	}
	p.intakeWG.Wait()
	if worker != nil {
		if err := worker.StopContext(ctx); err != nil && !errors.Is(err, ctx.Err()) {
			p.report(err)
		}
		p.report(worker.WaitStopped(context.WithoutCancel(ctx)))
	}
	p.lifecycle.Lock()
	p.started = false
	close(done)
	p.lifecycle.Unlock()
}
