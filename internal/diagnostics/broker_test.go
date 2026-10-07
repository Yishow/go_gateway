package diagnostics

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
	"unicode/utf8"
)

func testBroker(t *testing.T) *Broker {
	t.Helper()
	b, err := New(Options{})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := b.Close(ctx); err != nil {
			t.Error(err)
		}
	})
	return b
}
func emitBatch(t *testing.T, b *Broker, n int) {
	t.Helper()
	for range n {
		b.Emit(Input{Code: "startup.ready"})
	}
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
}
func TestEncodedEventUTF8Limit(t *testing.T) {
	e, size := boundEvent(Event{Code: "safe", Message: strings.Repeat("界<", 10000)})
	raw, _ := json.Marshal(e)
	if size > MaxEventBytes || len(raw) > MaxEventBytes || !utf8.Valid(raw) || !e.Truncated {
		t.Fatal("encoded bound or truncation failed")
	}
}
func TestHTTPProjectionAllowsOnlyScalarTemplates(t *testing.T) {
	b := testBroker(t)
	b.RegisterRoutes([]string{"/api/:id"})
	b.Emit(Input{Code: "http.access", Fields: map[string]any{"route": "/secret/password=abc", "method": "secret", "status": map[string]any{"secret": "secret"}, "duration_ms": -1, "request_id": "secret", "body": "secret"}})
	b.Emit(Input{Code: "http.access", Fields: map[string]any{"route": "/api/:id", "method": "GET", "status": 200, "duration_ms": 1.5, "request_id": strings.Repeat("a", 32)}})
	_ = b.Flush(t.Context())
	snap, _ := b.Snapshot(Query{})
	if len(snap.Records) != 2 || len(snap.Records[0].Fields) != 0 || len(snap.Records[1].Fields) != 5 {
		t.Fatal("strict projection failed")
	}
}
func TestAdmissionHasIndependentCountAndByteCaps(t *testing.T) {
	for _, large := range []bool{false, true} {
		t.Run(strconv.FormatBool(large), func(t *testing.T) {
			b := &Broker{instance: strings.Repeat("a", 32), level: 1, routes: map[string]struct{}{}, queue: newQueue(AdmissionRecords), wake: make(chan struct{}, 1)}
			in := Input{Code: "startup.ready"}
			if large {
				route := "/" + strings.Repeat("<", 1023)
				b.routes[route] = struct{}{}
				in = Input{Code: "http.access", Fields: map[string]any{"route": route}}
			}
			for range 2000 {
				b.Emit(in)
			}
			if b.queue.count > AdmissionRecords || b.queue.bytes > AdmissionBytes || b.dropped == 0 {
				t.Fatal("admission overflow not bounded")
			}
			if large && b.queue.count >= AdmissionRecords {
				t.Fatal("byte bound did not apply independently")
			}
		})
	}
}
func TestRingHasIndependentCountAndByteCaps(t *testing.T) {
	for _, large := range []bool{false, true} {
		t.Run(strconv.FormatBool(large), func(t *testing.T) {
			b := testBroker(t)
			in := Input{Code: "startup.ready"}
			if large {
				route := "/" + strings.Repeat("<", 1023)
				b.RegisterRoutes([]string{route})
				in = Input{Code: "http.access", Fields: map[string]any{"route": route}}
			}
			for range 24 {
				for range 100 {
					b.Emit(in)
				}
				_ = b.Flush(t.Context())
			}
			b.mu.Lock()
			defer b.mu.Unlock()
			if b.ring.count > RingRecords || b.ring.bytes > RingBytes || b.evicted == 0 {
				t.Fatal("retention not bounded")
			}
			if large && b.ring.count >= RingRecords {
				t.Fatal("byte eviction missing")
			}
		})
	}
}
func TestConcurrentSequencesAndDistinctInstances(t *testing.T) {
	b := testBroker(t)
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			for range 300 {
				b.Emit(Input{Code: "startup.ready"})
			}
		})
	}
	wg.Wait()
	_ = b.Flush(t.Context())
	s, _ := b.Snapshot(Query{Limit: 500})
	var prev uint64
	for _, e := range s.Records {
		n, _ := strconv.ParseUint(e.Sequence, 10, 64)
		if n <= prev || e.InstanceID != s.InstanceID {
			t.Fatal("sequence identity not monotonic")
		}
		prev = n
	}
	if testBroker(t).instance == b.instance {
		t.Fatal("instance reused")
	}
}
func TestQueryValidationAndUnicodeFolding(t *testing.T) {
	b := testBroker(t)
	b.RegisterRoutes([]string{"/Straße"})
	b.Emit(Input{Code: "http.access", Fields: map[string]any{"route": "/Straße"}})
	_ = b.Flush(t.Context())
	s, err := b.Snapshot(Query{Q: "STRASSE"})
	if err != nil || len(s.Records) != 1 {
		t.Fatal("Unicode literal folding failed")
	}
	for _, q := range []Query{{Limit: 501}, {Limit: -1}, {Q: strings.Repeat("界", 257)}, {Level: "fatal"}, {Source: "secret"}, {Before: b.cursor(999)}} {
		if _, err := b.Snapshot(q); err == nil {
			t.Fatalf("invalid query accepted %#v", q)
		}
	}
	q, err := ValidateQuery(Query{Q: strings.Repeat("ß", 256)})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := b.Snapshot(q); err != nil {
		t.Fatal("double fold expands past input limit")
	}
	empty, _ := b.Snapshot(Query{Q: "absent"})
	raw, _ := json.Marshal(empty)
	if !strings.Contains(string(raw), `"records":[]`) {
		t.Fatal("empty records not array")
	}
}
func TestReplaySeparateAndProgressRequiresAck(t *testing.T) {
	b := testBroker(t)
	emitBatch(t, b, 400)
	s, err := b.Subscribe(Query{}, b.cursor(0))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if s.replay.count != 400 || s.live.count != 0 || s.Metadata().Progress != b.cursor(0) {
		t.Fatal("replay shares live bound or skips progress")
	}
	e, err := s.Next(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if s.Metadata().Progress != b.cursor(0) {
		t.Fatal("unacknowledged write skipped")
	}
	s.Ack(e.Cursor())
	if s.Metadata().Progress != b.cursor(1) {
		t.Fatal("ack did not advance")
	}
	emitBatch(t, b, 1)
	if s.live.count != 1 {
		t.Fatal("post-cutoff event missing")
	}
}
func TestRetentionReplayLimitResetAndFutureCursor(t *testing.T) {
	b := testBroker(t)
	for range 23 {
		emitBatch(t, b, 100)
	}
	s, err := b.Subscribe(Query{}, b.cursor(0))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if s.replay.count != ReplayRecords || s.replay.bytes > ReplayBytes {
		t.Fatal("replay unbounded")
	}
	reasons := map[string]bool{}
	for _, g := range s.Metadata().Gaps {
		reasons[g.Reason] = true
	}
	if !reasons["retention"] || !reasons["replay_limit"] {
		t.Fatal("missing explicit gaps")
	}
	restart, err := b.Subscribe(Query{}, strings.Repeat("f", 32)+":9999")
	if err != nil {
		t.Fatal(err)
	}
	defer restart.Close()
	if restart.Metadata().Gaps[0].Reason != "reset" {
		t.Fatal("restart not explicit")
	}
	if _, err := b.Subscribe(Query{}, b.cursor(9999)); !errors.Is(err, ErrInvalidCursor) {
		t.Fatal("future cursor accepted")
	}
}
func TestFilteredProgressSlowReaderAndCapacity(t *testing.T) {
	b := testBroker(t)
	filtered, _ := b.Subscribe(Query{Level: "error"}, b.cursor(0))
	defer filtered.Close()
	slow, _ := b.Subscribe(Query{}, b.cursor(0))
	defer slow.Close()
	for range 3 {
		emitBatch(t, b, 100)
	}
	if filtered.Metadata().Progress != b.cursor(300) {
		t.Fatal("filtered events froze progress")
	}
	if _, err := slow.Next(t.Context()); !errors.Is(err, ErrOverflow) {
		t.Fatal("slow reader not closed")
	}
	if filtered.Metadata().SubscriberOverflow != 1 {
		t.Fatal("subscriber loss hidden")
	}
	subs := make([]*Subscription, 0, 6)
	for range 6 {
		s, err := b.Subscribe(Query{}, b.cursor(300))
		if err != nil {
			t.Fatal(err)
		}
		subs = append(subs, s)
	}
	if _, err := b.Subscribe(Query{}, ""); !errors.Is(err, ErrCapacity) {
		t.Fatal("ninth client admitted")
	}
	var wg sync.WaitGroup
	for _, s := range subs {
		wg.Go(func() { s.Close(); s.Close() })
	}
	wg.Go(b.CloseSubscriptions)
	wg.Wait()
}
