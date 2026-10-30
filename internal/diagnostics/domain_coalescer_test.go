package diagnostics

import (
	"context"
	"math"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

type coalescerClock struct {
	mu sync.Mutex
	at time.Time
}

func (c *coalescerClock) Now() time.Time          { c.mu.Lock(); defer c.mu.Unlock(); return c.at }
func (c *coalescerClock) Advance(d time.Duration) { c.mu.Lock(); c.at = c.at.Add(d); c.mu.Unlock() }

type coalescerSink struct {
	mu     sync.Mutex
	events []Event
	drops  uint64
	reject bool
}

func (s *coalescerSink) Admit(e Event) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.reject {
		return false
	}
	s.events = append(s.events, e)
	return true
}
func (s *coalescerSink) Drop(n uint64) { s.mu.Lock(); s.drops += n; s.mu.Unlock() }
func (s *coalescerSink) Events() []Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]Event(nil), s.events...)
}

func newCoalescerFixture(t *testing.T) (*DomainCoalescer, *coalescerClock, *coalescerSink) {
	t.Helper()
	clock := &coalescerClock{at: time.Date(2026, 10, 27, 10, 0, 0, 0, time.UTC)}
	sink := &coalescerSink{}
	c, err := NewDomainCoalescer(DomainCoalescerOptions{Heartbeat: 15 * time.Second, Clock: clock.Now, TryAdmit: sink.Admit, OnDrop: sink.Drop})
	if err != nil {
		t.Fatal(err)
	}
	return c, clock, sink
}

func coalescerRevision(v DomainEventRef, revision string) DomainEventRef {
	v.Binding.EntityRevision, v.Locator.GroupRevision = revision, domainText(revision)
	if v.Evidence.Receipt != nil {
		v.Evidence.Receipt.GroupRevision = revision
	}
	return v
}

func TestDomainCoalescerCapacityAndDropSeparation(t *testing.T) {
	c, clock, sink := newCoalescerFixture(t)
	// Each first event reserves an entry and an output in the same hard budget.
	// After draining, 1023 states leave one slot; a new pair needs two slots.
	for i := range AdmissionRecords - 1 {
		v := domainFixture()
		v.Locator.GroupID = domainText("G" + strconv.Itoa(i))
		if !c.Observe(t.Context(), projectionJSON(t, v), domainScope()) {
			t.Fatalf("count capacity prematurely refused %d: %#v", i, c.Stats())
		}
	}
	v := domainFixture()
	v.Locator.GroupID = domainText("overflow")
	if c.Observe(t.Context(), projectionJSON(t, v), domainScope()) {
		t.Fatal("map count cap exceeded")
	}
	stats := c.Stats()
	if stats.Entries != 1023 || stats.Queued != 0 || stats.Dropped != 1 || stats.StateBytes+stats.QueueBytes > AdmissionBytes {
		t.Fatalf("high cardinality unbounded: %#v", stats)
	}
	clock.Advance(15 * time.Second)
	c.Tick(t.Context())
	if len(sink.Events()) != 1023 {
		t.Fatal("single occurrences became periodic spam")
	}
	// Long opaque identities force the unchanged byte cap independently of count.
	b, _, _ := newCoalescerFixture(t)
	for i := range AdmissionRecords {
		large := domainFixture()
		large.Locator.GroupID = domainText(strings.Repeat("g", 120) + strconv.Itoa(i))
		large.Evidence.ReferenceID = strings.Repeat("e", 128)
		large.ConnectorID, large.ConnectorRevision = domainText(strings.Repeat("c", 128)), domainText(strings.Repeat("r", 128))
		large.Binding.WorkspaceID, large.Binding.WorkspaceRevision = strings.Repeat("w", 128), strings.Repeat("v", 128)
		large.Locator.WorkspaceID = large.Binding.WorkspaceID
		scope := domainScope()
		scope.WorkspaceID, scope.WorkspaceRevision = large.Binding.WorkspaceID, large.Binding.WorkspaceRevision
		b.Observe(t.Context(), projectionJSON(t, large), scope)
	}
	if s := b.Stats(); s.Entries >= 1024 || s.Dropped == 0 || s.StateBytes+s.QueueBytes > 2<<20 {
		t.Fatalf("byte cap not independent: %#v", s)
	}
	// Admission failure is loss, not a repeat; intentional coalescing is not loss.
	x, timer, out := newCoalescerFixture(t)
	out.reject = true
	x.Observe(t.Context(), projectionJSON(t, domainFixture()), domainScope())
	out.reject = false
	x.Observe(t.Context(), projectionJSON(t, domainFixture()), domainScope())
	if x.Stats().Dropped != 1 {
		t.Fatal("coalesced repeat counted as drop")
	}
	timer.Advance(15 * time.Second)
	x.Tick(t.Context())
	if out.Events()[0].Fields["occurrences"] != "2" || x.Stats().Dropped != 1 {
		t.Fatal("drop was added to occurrence count")
	}
}

func TestDomainCoalescerCancelCloseAndContention(t *testing.T) {
	c, _, sink := newCoalescerFixture(t)
	raw := projectionJSON(t, domainFixture())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if c.Observe(ctx, raw, domainScope()) || c.Stats().Entries != 0 {
		t.Fatal("canceled observation was retained")
	}
	c.mu.Lock()
	done := make(chan bool, 1)
	go func() { done <- c.Observe(context.Background(), raw, domainScope()) }()
	select {
	case accepted := <-done:
		if accepted {
			t.Error("contended observation accepted")
		}
	case <-time.After(time.Second):
		t.Error("writer waited on coalescer lock")
	}
	c.mu.Unlock()
	c.Observe(t.Context(), raw, domainScope())
	c.Observe(t.Context(), raw, domainScope())
	c.Close()
	c.Close()
	if s := c.Stats(); s.Entries != 0 || s.StateBytes != 0 || s.Queued != 0 || s.Dropped != 3 {
		t.Fatalf("close must separately count one lost summary: %#v", s)
	}
	if c.Observe(t.Context(), raw, domainScope()) || len(sink.Events()) != 1 {
		t.Fatal("closed coalescer emitted")
	}
	c.Tick(t.Context())
	if c.Stats().Dropped != 4 {
		t.Fatal("closed capture loss missing or counted twice")
	}
}

func TestDomainCoalescerClockCountAndConstructor(t *testing.T) {
	clock := &coalescerClock{at: time.Now().UTC()}
	sink := &coalescerSink{}
	valid := DomainCoalescerOptions{Heartbeat: 15 * time.Second, Clock: clock.Now, TryAdmit: sink.Admit, OnDrop: sink.Drop}
	for _, mutate := range []func(*DomainCoalescerOptions){
		func(o *DomainCoalescerOptions) { o.Heartbeat = 0 }, func(o *DomainCoalescerOptions) { o.Heartbeat = -1 },
		func(o *DomainCoalescerOptions) { o.Clock = nil }, func(o *DomainCoalescerOptions) { o.TryAdmit = nil }, func(o *DomainCoalescerOptions) { o.OnDrop = nil },
		func(o *DomainCoalescerOptions) { o.RecordLimit = AdmissionRecords + 1 }, func(o *DomainCoalescerOptions) { o.ByteLimit = AdmissionBytes + 1 },
	} {
		o := valid
		mutate(&o)
		if c, err := NewDomainCoalescer(o); err == nil || c != nil {
			t.Fatal("invalid config must fail closed")
		}
	}
	c, timer, out := newCoalescerFixture(t)
	raw := projectionJSON(t, domainFixture())
	c.Observe(t.Context(), raw, domainScope())
	timer.Advance(-time.Second)
	if c.Observe(t.Context(), raw, domainScope()) {
		t.Fatal("regressing clock corrupted chronology")
	}
	timer.Advance(16 * time.Second)
	c.mu.Lock()
	for _, entry := range c.entries {
		entry.occurrences, entry.reported = math.MaxUint64, math.MaxUint64-1
	}
	c.mu.Unlock()
	if c.Observe(t.Context(), raw, domainScope()) {
		t.Fatal("occurrence count wrapped")
	}
	c.Tick(t.Context())
	if got := out.Events(); got[len(got)-1].Fields["occurrences"] != "18446744073709551615" {
		t.Fatal("counter lost decimal precision")
	}
}

func TestDomainCoalescerNoSuccessSpamOrRestartReplay(t *testing.T) {
	c, clock, sink := newCoalescerFixture(t)
	commit := projectionJSON(t, domainCommitFixture())
	for range 20 {
		c.Observe(t.Context(), commit, domainScope())
	}
	if len(sink.Events()) != 1 {
		t.Fatal("per-sample success spam")
	}
	clock.Advance(15 * time.Second)
	c.Tick(t.Context())
	if got := sink.Events(); len(got) != 2 || got[1].Code != "domain.delivery.summary" || got[1].Fields["occurrences"] != "20" {
		t.Fatal("repeat success was mislabeled as a new first commitment")
	}
	c.Close()
	fresh, _, empty := newCoalescerFixture(t)
	fresh.Tick(t.Context())
	if len(empty.Events()) != 0 {
		t.Fatal("new coalescer replayed historical receipts")
	}
	r := domainCommitFixture()
	r.Code = "domain.delivery.recovered"
	if fresh.Observe(t.Context(), projectionJSON(t, r), domainScope()) {
		t.Fatal("recovery bypassed explicit matching fault")
	}
}

func TestDomainCoalescerOutputQueueBounds(t *testing.T) {
	for _, limit := range []struct{ records, bytes int }{{2, AdmissionBytes}, {64, 1024}} {
		clock := &coalescerClock{at: time.Now().UTC()}
		sink := &coalescerSink{}
		c, err := NewDomainCoalescer(DomainCoalescerOptions{Heartbeat: 15 * time.Second, Clock: clock.Now, TryAdmit: sink.Admit, OnDrop: sink.Drop, RecordLimit: limit.records, ByteLimit: limit.bytes})
		if err != nil {
			t.Fatal(err)
		}
		c.draining.Store(true) // Hold the bounded drain, without a blocking admission callback.
		for range 100 {
			c.Observe(t.Context(), []byte(`{"raw_error":"password_fixture"}`), domainScope())
		}
		before := c.Stats()
		if before.Entries != 0 || before.Queued > limit.records || before.QueueBytes > limit.bytes || before.Dropped == 0 {
			t.Fatalf("output queue not independently bounded: %#v", before)
		}
		if limit.records == 2 && before.Queued != 2 {
			t.Fatal("exact output record cap not reached")
		}
		if limit.bytes == 1024 && before.Queued >= 64 {
			t.Fatal("output byte cap did not apply")
		}
		c.Close()
		c.draining.Store(false)
		if after := c.Stats(); after.Queued != 0 || after.QueueBytes != 0 || after.Dropped != 100 {
			t.Fatalf("queued close loss not accounted: %#v", after)
		}
	}
}

func TestDomainCoalescerCallbacksReenterWithoutLocks(t *testing.T) {
	var c *DomainCoalescer
	clock := &coalescerClock{at: time.Now().UTC()}
	sink := &coalescerSink{}
	other := domainFixture()
	other.Locator.GroupID = domainText("other-G")
	otherRaw := projectionJSON(t, other)
	reentered := false
	var err error
	c, err = NewDomainCoalescer(DomainCoalescerOptions{Heartbeat: 15 * time.Second, Clock: clock.Now,
		TryAdmit: func(e Event) bool {
			_ = c.Stats()
			if !reentered {
				reentered = true
				c.Observe(context.Background(), otherRaw, domainScope())
			}
			return sink.Admit(e)
		},
		OnDrop: func(n uint64) { _ = c.Stats(); sink.Drop(n) },
	})
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})
	raw := projectionJSON(t, domainFixture())
	go func() { c.Observe(context.Background(), raw, domainScope()); close(done) }()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("callback deadlocked on coalescer mutex")
	}
	if len(sink.Events()) != 2 || c.Stats().Entries != 2 {
		t.Fatal("reentrant bounded enqueue lost its first event")
	}
	c.options.TryAdmit = func(Event) bool { panic("private_fixture") }
	c.options.OnDrop = func(uint64) { panic("private_fixture") }
	c.Tick(t.Context())
	third := domainFixture()
	third.Locator.GroupID = domainText("third-G")
	c.Observe(t.Context(), projectionJSON(t, third), domainScope())
	if c.Stats().Dropped != 1 {
		t.Fatal("callback failure escaped loss accounting")
	}
}

func TestDomainCoalescerConcurrentLifecycle(t *testing.T) {
	c, _, sink := newCoalescerFixture(t)
	var inputs [][]byte
	for i := range 16 {
		v := domainFixture()
		v.Locator.GroupID = domainText("G" + strconv.Itoa(i))
		inputs = append(inputs, projectionJSON(t, v))
	}
	var wg sync.WaitGroup
	for worker := range 8 {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for i := range 64 {
				c.Observe(context.Background(), inputs[(i+worker)%16], domainScope())
				if i%7 == 0 {
					c.Tick(context.Background())
				}
				if worker == 0 && i == 32 {
					c.Close()
				}
			}
		}(worker)
	}
	wg.Wait()
	c.Close()
	stats := c.Stats()
	if stats.Entries != 0 || stats.StateBytes != 0 || stats.Queued != 0 || stats.QueueBytes != 0 || stats.Dropped == 0 || stats.Dropped != sink.drops {
		t.Fatalf("concurrent lifecycle accounting inconsistent: %#v / drops %d", stats, sink.drops)
	}
}

func TestDomainCoalescerClockFailureAndInputOwnership(t *testing.T) {
	c, clock, sink := newCoalescerFixture(t)
	raw := projectionJSON(t, domainFixture())
	c.Observe(t.Context(), raw, domainScope())
	for i := range raw {
		raw[i] = ' '
	}
	clock.Advance(time.Second)
	c.Observe(t.Context(), projectionJSON(t, domainFixture()), domainScope())
	clock.Advance(14 * time.Second)
	c.Tick(t.Context())
	if sink.Events()[1].Fields["reference_id"] != "E" {
		t.Fatal("retained reference aliased producer input")
	}
	for _, now := range []func() time.Time{func() time.Time { return time.Time{} }, func() time.Time { panic("private_clock_fixture") }} {
		c.options.Clock = now
		if c.Observe(t.Context(), projectionJSON(t, domainFixture()), domainScope()) {
			t.Fatal("missing clock fabricated observation time")
		}
	}
	if c.Stats().Dropped != 2 {
		t.Fatal("clock failures were not separately accounted")
	}
}

func TestDomainCoalescerAggregateRecordBudget(t *testing.T) {
	for _, active := range []bool{false, true} {
		t.Run(strconv.FormatBool(active), func(t *testing.T) {
			c, _, sink := newCoalescerFixture(t)
			c.options.RecordLimit = 4
			capture := func(group string) bool {
				v := domainFixture()
				v.Locator.GroupID = domainText(group)
				return c.Observe(t.Context(), projectionJSON(t, v), domainScope())
			}
			assertFull := func() {
				if !capture("G2") || capture("G3") {
					t.Error("aggregate capacity must accept second pair and reject third")
				}
				stats := c.Stats()
				if stats.Entries != 2 || stats.Queued != 2 || stats.Entries+stats.Queued != 4 || stats.Dropped != 1 {
					t.Errorf("entries plus queued/in-flight outputs exceeded one budget: %#v", stats)
				}
			}
			if active {
				reentered := false
				c.options.TryAdmit = func(e Event) bool {
					if !reentered {
						reentered = true
						assertFull()
					}
					return sink.Admit(e)
				}
			} else {
				c.draining.Store(true)
			}
			if !capture("G1") {
				t.Fatal("first pair refused")
			}
			if !active {
				assertFull()
				c.draining.Store(false)
				c.Tick(t.Context())
			}
			if len(sink.Events()) != 2 || c.Stats().Entries != 2 || c.Stats().Dropped != 1 {
				t.Fatal("aggregate refusal altered the two accepted episodes")
			}
		})
	}
}

func TestDomainCoalescerInFlightOutputConsumesSharedSlot(t *testing.T) {
	c, _, sink := newCoalescerFixture(t)
	c.options.RecordLimit = 2
	c.options.TryAdmit = func(e Event) bool {
		if e.Code == "domain.delivery.retrying" {
			if c.Observe(t.Context(), []byte(`{"raw_error":"password_fixture"}`), domainScope()) {
				t.Error("in-flight output slot was reused for queued suppression")
			}
			stats := c.Stats()
			if stats.Entries != 1 || stats.Queued != 1 || stats.Dropped != 1 {
				t.Errorf("in-flight output escaped aggregate cap: %#v", stats)
			}
		}
		return sink.Admit(e)
	}
	if !c.Observe(t.Context(), projectionJSON(t, domainFixture()), domainScope()) || len(sink.Events()) != 1 || sink.drops != 1 {
		t.Fatal("active slot refusal lost original first event or drop accounting")
	}
}
