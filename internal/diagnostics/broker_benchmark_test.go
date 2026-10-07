package diagnostics

import (
	"context"
	"testing"
)

func BenchmarkSanitizedCapture(b *testing.B) {
	broker, err := New(Options{})
	if err != nil {
		b.Fatal(err)
	}
	defer func() { _ = broker.Close(context.Background()) }()
	b.ReportAllocs()
	for b.Loop() {
		broker.Emit(Input{Code: "startup.ready"})
	}
	if err := broker.Flush(b.Context()); err != nil {
		b.Fatal(err)
	}
	broker.mu.Lock()
	b.ReportMetric(float64(broker.queue.bytes), "queue-B")
	b.ReportMetric(float64(broker.ring.bytes), "ring-B")
	b.ReportMetric(float64(broker.ring.count), "retained-records")
	broker.mu.Unlock()
}
