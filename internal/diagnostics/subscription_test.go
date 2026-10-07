package diagnostics

import (
	"errors"
	"strings"
	"testing"
)

func TestReplayAndLiveIndependentByteLimits(t *testing.T) {
	b := testBroker(t)
	route := "/" + strings.Repeat("<", 1023)
	b.RegisterRoutes([]string{route})
	emit := func(n int) {
		for range n {
			b.Emit(Input{Code: "http.access", Fields: map[string]any{"route": route}})
		}
		_ = b.Flush(t.Context())
	}
	for range 5 {
		emit(100)
	}
	s, _ := b.Subscribe(Query{}, b.cursor(0))
	defer s.Close()
	if s.replay.bytes > ReplayBytes || s.replay.count >= ReplayRecords {
		t.Fatal("replay byte bound missing")
	}
	found := false
	for _, g := range s.Metadata().Gaps {
		found = found || g.Reason == "replay_limit"
	}
	if !found {
		t.Fatal("byte replay loss hidden")
	}
	live, _ := b.Subscribe(Query{}, b.cursor(500))
	defer live.Close()
	emit(100)
	if _, err := live.Next(t.Context()); !errors.Is(err, ErrOverflow) {
		t.Fatal("live byte cap missing")
	}
}
func TestReturnedFieldsCannotMutateRing(t *testing.T) {
	b := testBroker(t)
	b.RegisterRoutes([]string{"/safe"})
	b.Emit(Input{Code: "http.access", Fields: map[string]any{"route": "/safe"}})
	_ = b.Flush(t.Context())
	first, _ := b.Snapshot(Query{})
	first.Records[0].Fields["route"] = "injected-secret"
	s, _ := b.Subscribe(Query{}, b.cursor(0))
	defer s.Close()
	e, _ := s.Next(t.Context())
	e.Fields["route"] = "injected-secret"
	again, _ := b.Snapshot(Query{})
	if again.Records[0].Fields["route"] != "/safe" {
		t.Fatal("returned map aliases safe storage")
	}
}
func TestNoArbitraryFormattingAndShutdownPhases(t *testing.T) {
	b := testBroker(t)
	b.Emit(Input{Code: "http.access", Fields: map[string]any{"method": panicStringer{}, "route": panicStringer{}, "status": panicStringer{}, "request_id": panicStringer{}}})
	for _, phase := range []string{"startup", "http", "runtime", "pipeline", "share", "connections", "database", "diagnostics", "tray", "owner"} {
		b.Emit(Input{Code: "shutdown.finalizing", Fields: map[string]any{"phase": phase}})
	}
	_ = b.Flush(t.Context())
	snap, _ := b.Snapshot(Query{})
	if len(snap.Records[0].Fields) != 0 || len(snap.Records) != 11 {
		t.Fatal("scalar projection failed")
	}
	for _, e := range snap.Records[1:] {
		if e.Fields["phase"] == nil {
			t.Fatal("shutdown phase omitted")
		}
	}
}

type panicStringer struct{}

func (panicStringer) String() string { panic("must not format raw fields") }

func TestSnapshotToSubscribeConcurrentAppendHasNoGap(t *testing.T) {
	b := testBroker(t)
	emitBatch(t, b, 50)
	done := make(chan struct{})
	go func() {
		defer close(done)
		for range 200 {
			b.Emit(Input{Code: "startup.ready"})
		}
	}()
	snapshot, err := b.Snapshot(Query{Limit: 500})
	if err != nil {
		t.Fatal(err)
	}
	sub, err := b.Subscribe(Query{}, snapshot.Latest)
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Close()
	<-done
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	_, start, err := parseCursor(snapshot.Latest)
	if err != nil {
		t.Fatal(err)
	}
	for want := start + 1; want <= 250; want++ {
		e, err := sub.Next(t.Context())
		if err != nil {
			t.Fatal(err)
		}
		if e.Cursor() != b.cursor(want) {
			t.Fatalf("unexplained boundary gap: want%s got%s", b.cursor(want), e.Cursor())
		}
		sub.Ack(e.Cursor())
	}
	if sub.Metadata().Progress != b.cursor(250) {
		t.Fatal("progress did not reach delivered boundary")
	}
}
