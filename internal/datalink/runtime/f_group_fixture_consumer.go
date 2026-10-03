//go:build f_write_group_fixture

package runtime

import (
	"context"
	"sync"
	"time"

	"go-gateway/internal/datalink/collector"
)

type fixtureConsumerGate struct {
	mu     sync.Mutex
	paused bool
	active int
	idle   chan struct{}
}

func newFixtureConsumerGate() *fixtureConsumerGate {
	idle := make(chan struct{})
	close(idle)
	return &fixtureConsumerGate{idle: idle}
}

func (g *fixtureConsumerGate) begin() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.paused {
		return false
	}
	if g.active == 0 {
		g.idle = make(chan struct{})
	}
	g.active++
	return true
}

func (g *fixtureConsumerGate) end() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.active--
	if g.active == 0 {
		close(g.idle)
	}
}

func (g *fixtureConsumerGate) pause(ctx context.Context) error {
	g.mu.Lock()
	g.paused = true
	idle := g.idle
	g.mu.Unlock()

	select {
	case <-idle:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

var fixtureConsumerGates sync.Map // map[*Service]*fixtureConsumerGate

func fixtureConsumerGateFor(s *Service) *fixtureConsumerGate {
	if existing, ok := fixtureConsumerGates.Load(s); ok {
		gate, ok := existing.(*fixtureConsumerGate)
		if !ok {
			panic("fixture consumer gate has an unexpected type")
		}
		return gate
	}
	candidate := newFixtureConsumerGate()
	actual, _ := fixtureConsumerGates.LoadOrStore(s, candidate)
	gate, ok := actual.(*fixtureConsumerGate)
	if !ok {
		panic("fixture consumer gate has an unexpected type")
	}
	return gate
}

func fixtureConsumerBegin(s *Service) bool {
	return fixtureConsumerGateFor(s).begin()
}

func fixtureConsumerEnd(s *Service) {
	fixtureConsumerGateFor(s).end()
}

func (s *Service) consumeCollectedValue(runtimeCtx context.Context, cv collector.CollectedValue) {
	if !fixtureConsumerBegin(s) {
		return
	}
	s.collectedTotal.Add(1)
	func() {
		defer fixtureConsumerEnd(s)
		ctx, cancel := context.WithTimeout(runtimeCtx, 5*time.Second)
		defer cancel()
		s.handleCollectedValue(ctx, cv)
	}()
}

func pauseFixtureConsumer(ctx context.Context, s *Service) error {
	return fixtureConsumerGateFor(s).pause(ctx)
}

func drainFixtureValues(values <-chan collector.CollectedValue) {
	for {
		select {
		case _, ok := <-values:
			if !ok {
				return
			}
		default:
			return
		}
	}
}
