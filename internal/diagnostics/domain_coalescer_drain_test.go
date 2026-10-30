package diagnostics

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestDomainCoalescerConcurrentDrainDoesNotStrandAcceptedOutput(t *testing.T) {
	for round := range 20000 {
		var admitted, dropped atomic.Uint64
		c, err := NewDomainCoalescer(DomainCoalescerOptions{
			Heartbeat: time.Second, Clock: time.Now, RecordLimit: 2,
			TryAdmit: func(Event) bool { admitted.Add(1); return true },
			OnDrop:   func(n uint64) { dropped.Add(n) },
		})
		if err != nil {
			t.Fatal(err)
		}
		// Two outputs cannot exhaust the drain's two-record batch and leave a
		// third output pending. No Tick or later producer repairs the handoff.
		start := make(chan struct{})
		var workers sync.WaitGroup
		for range 2 {
			workers.Go(func() {
				<-start
				c.enqueue(t.Context(), Event{Code: "domain.unavailable"})
			})
		}
		close(start)
		workers.Wait()
		stats := c.Stats()
		if stats.Queued != 0 || admitted.Load()+dropped.Load() != 2 {
			t.Fatalf("round %d: drain left accepted output without a caller: queued=%d admitted=%d dropped=%d", round, stats.Queued, admitted.Load(), dropped.Load())
		}
	}
}
