package adapters

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

var errAcquisitionConnectSuperseded = errors.New("acquisition connection attempt superseded")

// acquisitionGeneration is connection-local cache lifecycle, not source state.
// The managed gate normally serializes operations. This short publication lock
// also rejects late completions from direct callers without locking network I/O.
type acquisitionGeneration struct {
	mu        sync.Mutex
	attempt   uint64
	published atomic.Uint64
}

func (g *acquisitionGeneration) connect(ctx context.Context, connect func(context.Context) error) error {
	attempt, err := g.begin(ctx)
	if err != nil {
		return err
	}
	if err := connect(ctx); err != nil {
		return err
	}
	return g.publish(ctx, attempt)
}

func (g *acquisitionGeneration) begin(ctx context.Context) (uint64, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	if g.attempt == ^uint64(0) {
		return 0, errAcquisitionConnectSuperseded
	}
	g.attempt++
	return g.attempt, nil
}

func (g *acquisitionGeneration) publish(ctx context.Context, attempt uint64) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return err
	}
	// Allocation and this check/store are one ordered lifecycle: a newer begin
	// cannot slip between the current-attempt check and publication.
	if g.attempt != attempt || g.published.Load() >= attempt {
		return errAcquisitionConnectSuperseded
	}
	g.published.Store(attempt)
	return nil
}

// AcquisitionTransportGeneration is the last successfully published connection.
// The getter is atomic and never acquires the adapter/transport operation gate.
func (c *MC3EConnector) AcquisitionTransportGeneration() uint64 {
	return c.acquisitionGeneration.published.Load()
}

// AcquisitionTransportGeneration is the last successfully published connection.
// The getter is atomic and never acquires the adapter/transport operation gate.
func (c *FatekConnector) AcquisitionTransportGeneration() uint64 {
	return c.acquisitionGeneration.published.Load()
}
