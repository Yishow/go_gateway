package storage

import (
	"context"
	"errors"
	"io"
)

// Close closes the writer, observing actual completion without a wait deadline.
func (bw *BatchWriter) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), bw.config.WriteTimeout)
	defer cancel()
	err := bw.CloseContext(ctx)
	if err != nil && errors.Is(err, ctx.Err()) {
		// The legacy Close contract joins ownership before returning. Context-aware
		// owners can use CloseContext to receive an earlier deadline notification.
		return errors.Join(err, bw.WaitClosed(context.WithoutCancel(ctx)))
	}
	return err
}

// CloseContext starts one final flush/close sequence and bounds only this
// caller's wait. A timer or writer ignoring cancellation retains ownership of
// the underlying storage until it really returns.
func (bw *BatchWriter) CloseContext(ctx context.Context) error {
	bw.lifecycle.Lock()
	if bw.closeDone == nil {
		bw.closing.Store(true)
		bw.closeDone = make(chan struct{})
		close(bw.stopCh)
		bw.timerCancel()
		go bw.finishClose(ctx, bw.closeDone)
	}
	done := bw.closeDone
	bw.lifecycle.Unlock()
	return bw.waitClosed(ctx, done)
}

// WaitClosed observes the original CloseContext operation without repeating it.
func (bw *BatchWriter) WaitClosed(ctx context.Context) error {
	bw.lifecycle.Lock()
	done := bw.closeDone
	bw.lifecycle.Unlock()
	if done == nil {
		return errors.New("batch writer has not been closed")
	}
	return bw.waitClosed(ctx, done)
}

func (bw *BatchWriter) waitClosed(ctx context.Context, done <-chan struct{}) error {
	select {
	case <-done:
		return bw.closeErr
	default:
	}
	select {
	case <-done:
		return bw.closeErr
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (bw *BatchWriter) finishClose(ctx context.Context, done chan struct{}) {
	bw.wg.Wait()
	bw.mu.Lock()
	flushErr := bw.flushLocked(ctx)
	bw.mu.Unlock()
	bw.closeErr = errors.Join(flushErr, bw.underlying.Close())
	close(done)
}

func (bw *BatchWriter) checkOpen() error {
	if bw.closing.Load() {
		return io.ErrClosedPipe
	}
	return nil
}
