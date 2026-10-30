package diagnostics

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/datalink/runtimeview"
)

type mountedDomainEmitter interface {
	EmitDomain(context.Context, []byte, runtimeview.Scope) bool
	EmitDomainUnavailable(context.Context) bool
}

func mountedEmitter(t *testing.T, b *Broker) mountedDomainEmitter {
	t.Helper()
	emitter, ok := any(b).(mountedDomainEmitter)
	if !ok {
		t.Fatal("broker has no mounted, validated domain capture boundary")
	}
	return emitter
}

func sourceDomainJSON(t *testing.T) []byte {
	t.Helper()
	return projectionJSON(t, projectionRef(projectionCases[1]))
}

func TestDomainBrokerPreservesProjectedFields(t *testing.T) {
	b := testBroker(t)
	emitter := mountedEmitter(t, b)
	raw := sourceDomainJSON(t)
	if !emitter.EmitDomain(t.Context(), raw, domainScope()) {
		t.Fatal("first source observation refused")
	}
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	s, err := b.Snapshot(Query{})
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Records) != 1 {
		t.Fatalf("want one event, got %#v", s)
	}
	e := s.Records[0]
	if e.Code != "domain.source.rejected" || e.Fields["workspace_revision"] != "W1" || e.Fields["entity_revision"] != "R1" || e.Fields["rule_id"] != "S" || e.Fields["occurrences"] != "1" || e.Fields["reference_id"] != "E" {
		t.Fatalf("domain proof lost to legacy re-projection: %#v", e)
	}
	if e.InstanceID == "" || e.Sequence != "1" || e.Timestamp.IsZero() {
		t.Fatalf("admission envelope absent: %#v", e)
	}
	// Legacy untrusted Input remains a separate fail-closed boundary.
	b.Emit(Input{Code: "domain.source.rejected", Fields: map[string]any{"rule_id": "forged"}})
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	s, _ = b.Snapshot(Query{})
	if s.Records[1].Code != codeRawSuppressed || len(s.Records[1].Fields) != 0 {
		t.Fatal("legacy domain Input was trusted")
	}
}

func TestDomainBrokerRejectsHostileAndUnknownCandidates(t *testing.T) {
	b := testBroker(t)
	emitter := mountedEmitter(t, b)
	base := sourceDomainJSON(t)
	cases := [][]byte{
		[]byte(strings.TrimSuffix(string(base), "}") + `,"raw_error":"password_fixture","credentials":{"password":"nested_secret"}}`),
		[]byte(strings.Replace(string(base), "domain.source.rejected", "domain.unknown", 1)),
		[]byte(strings.Replace(string(base), `"S"`, `"http://private.example:5432/path"`, 1)),
		[]byte(strings.Replace(string(base), `"S"`, `"`+strings.Repeat("a", 129)+`"`, 1)),
	}
	for _, raw := range cases {
		emitter.EmitDomain(t.Context(), raw, domainScope())
		if err := b.Flush(t.Context()); err != nil {
			t.Fatal(err)
		}
	}
	s, _ := b.Snapshot(Query{})
	if len(s.Records) != len(cases) {
		t.Fatalf("expected suppressed records, got %d", len(s.Records))
	}
	for _, e := range s.Records {
		if e.Code != codeRawSuppressed || len(e.Fields) != 0 {
			t.Fatalf("unsafe event retained: %#v", e)
		}
	}
	raw, _ := json.Marshal(s)
	for _, secret := range []string{"password_fixture", "nested_secret", "private.example"} {
		if strings.Contains(string(raw), secret) {
			t.Fatalf("secret retained: %s", secret)
		}
	}
}

func TestDomainBrokerUnavailableIsFixedAndUnattributed(t *testing.T) {
	b := testBroker(t)
	emitter := mountedEmitter(t, b)
	if !emitter.EmitDomainUnavailable(t.Context()) {
		t.Fatal("unavailable refused")
	}
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	s, _ := b.Snapshot(Query{})
	if len(s.Records) != 1 || s.Records[0].Code != "domain.unavailable" || s.Records[0].Source != "runtime" || len(s.Records[0].Fields) != 0 {
		t.Fatalf("unavailable fabricated identity: %#v", s)
	}
}

func TestDomainBrokerContentionIsNonblockingAndLossCountedOnce(t *testing.T) {
	b := testBroker(t)
	emitter := mountedEmitter(t, b)
	raw := sourceDomainJSON(t)
	b.mu.Lock()
	done := make(chan struct{})
	go func() { emitter.EmitDomain(context.Background(), raw, domainScope()); close(done) }()
	select {
	case <-done:
	case <-time.After(time.Second):
		b.mu.Unlock()
		t.Fatal("capture blocked on broker lock")
	}
	b.mu.Unlock()
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	s, _ := b.Snapshot(Query{})
	if s.CaptureDropped != 1 || len(s.Records) != 0 || len(s.Gaps) != 1 || s.Gaps[0].Reason != "admission" {
		t.Fatalf("loss double-counted or absent: %#v", s)
	}
}

func TestDomainBrokerFilteredEventsAreNotLoss(t *testing.T) {
	b, err := New(Options{Level: "error"})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = b.Close(context.Background()) })
	emitter := mountedEmitter(t, b)
	for range AdmissionRecords + 10 {
		emitter.EmitDomain(t.Context(), sourceDomainJSON(t), domainScope())
		emitter.EmitDomainUnavailable(t.Context())
	}
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	s, _ := b.Snapshot(Query{})
	if s.CaptureDropped != 0 || len(s.Records) != 0 {
		t.Fatalf("intentional filtering counted as loss: %#v", s)
	}
}

func TestDomainBrokerCloseCountsPendingSummaryNotOccurrences(t *testing.T) {
	b := testBroker(t)
	emitter := mountedEmitter(t, b)
	raw := sourceDomainJSON(t)
	for range 3 {
		emitter.EmitDomain(t.Context(), raw, domainScope())
	}
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	if err := b.Close(t.Context()); err != nil {
		t.Fatal(err)
	}
	s, _ := b.Snapshot(Query{})
	if len(s.Records) != 1 || s.CaptureDropped != 1 {
		t.Fatalf("close loss must be one dirty summary, not occurrences: %#v", s)
	}
	if emitter.EmitDomain(t.Context(), raw, domainScope()) {
		t.Fatal("closed broker accepted domain capture")
	}
	s, _ = b.Snapshot(Query{})
	if s.CaptureDropped != 2 {
		t.Fatalf("closed capture not separately counted: %#v", s)
	}
}

func pausedDomainBroker(t *testing.T) *Broker {
	t.Helper()
	b := &Broker{instance: strings.Repeat("a", 32), level: 1, routes: map[string]struct{}{}, queue: newQueue(AdmissionRecords), ring: newQueue(RingRecords), wake: make(chan struct{}, 1)}
	var err error
	b.domain, err = NewDomainCoalescer(DomainCoalescerOptions{Heartbeat: RuntimeLogHeartbeat, Clock: time.Now, TryAdmit: b.tryAdmitDomain, OnDrop: b.recordDomainLoss})
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestDomainBrokerAdmissionCapsCountLossExactlyOnce(t *testing.T) {
	for _, large := range []bool{false, true} {
		t.Run(strconv.FormatBool(large), func(t *testing.T) {
			b := pausedDomainBroker(t)
			input := Input{Code: "startup.ready"}
			if large {
				route := "/" + strings.Repeat("<", 1023)
				b.routes[route] = struct{}{}
				input = Input{Code: "http.access", Fields: map[string]any{"route": route}}
			}
			for b.Emit(input) {
			}
			// The rejection that finds the legacy cap is already one accounted loss.
			before := b.dropped
			if large {
				b.queue.bytes = AdmissionBytes
			} // exact byte-cap boundary, not a count-only refusal
			mountedEmitter(t, b).EmitDomain(t.Context(), sourceDomainJSON(t), domainScope())
			s, _ := b.Snapshot(Query{})
			if s.CaptureDropped != before+1 || b.domain.Stats().Dropped != 1 {
				t.Fatalf("one domain rejection counted more/less than once: metadata=%#v coalescer=%#v", s, b.domain.Stats())
			}
			if b.queue.count > AdmissionRecords || b.queue.bytes > AdmissionBytes {
				t.Fatal("existing admission cap expanded")
			}
			if large && b.queue.count >= AdmissionRecords {
				t.Fatal("byte case did not reach independent cap")
			}
		})
	}
}

func TestDomainBrokerCoalescerCapsAndRevisionIsolation(t *testing.T) {
	b := pausedDomainBroker(t)
	emitter := mountedEmitter(t, b)
	for i := range AdmissionRecords + 1 {
		ref := projectionRef(projectionCases[1])
		ref.Locator.RuleID = domainText("rule-" + strconv.Itoa(i))
		emitter.EmitDomain(t.Context(), projectionJSON(t, ref), domainScope())
	}
	stats := b.domain.Stats()
	s, _ := b.Snapshot(Query{})
	if stats.Entries != AdmissionRecords-1 || stats.StateBytes+stats.QueueBytes > AdmissionBytes || stats.Queued != 0 || stats.Dropped != 2 || s.CaptureDropped != 2 {
		t.Fatalf("bounded cardinality/loss mismatch: %#v %#v", stats, s)
	}
	c := pausedDomainBroker(t)
	out := mountedEmitter(t, c)
	first := projectionRef(projectionCases[1])
	second := first
	second.Binding.EntityRevision = "R2"
	for range 3 {
		out.EmitDomain(t.Context(), projectionJSON(t, first), domainScope())
	}
	out.EmitDomain(t.Context(), projectionJSON(t, second), domainScope())
	if c.queue.count != 2 || c.domain.Stats().Entries != 2 || c.domain.Stats().Dropped != 0 {
		t.Fatal("repeat/revision identity was merged or counted as loss")
	}
	c.domain.mu.Lock()
	for _, entry := range c.domain.entries {
		entry.next = time.Now().Add(-time.Second)
	}
	c.domain.mu.Unlock()
	c.domain.Tick(t.Context())
	if c.queue.count != 3 || c.queue.at(2).event.Fields["occurrences"] != "3" || c.queue.at(2).event.Fields["entity_revision"] != "R1" {
		t.Fatal("summary/revision occurrence evidence changed")
	}
}

func TestDomainBrokerCanceledAndClosedLossVisibleAfterDispatcher(t *testing.T) {
	b := testBroker(t)
	emitter := mountedEmitter(t, b)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if emitter.EmitDomain(ctx, sourceDomainJSON(t), domainScope()) {
		t.Fatal("canceled observation accepted")
	}
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	if err := b.Close(t.Context()); err != nil {
		t.Fatal(err)
	}
	for range 2 {
		if emitter.EmitDomainUnavailable(t.Context()) {
			t.Fatal("closed fixed capture accepted")
		}
	}
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	s, _ := b.Snapshot(Query{})
	if s.CaptureDropped != 3 || len(s.Records) != 0 || s.Latest != b.cursor(1) {
		t.Fatalf("closed/canceled loss hidden or cursor advanced after drain: %#v", s)
	}
}

func TestDomainBrokerConcurrentLossConservation(t *testing.T) {
	b := testBroker(t)
	emitter := mountedEmitter(t, b)
	const observations = 100
	var wg sync.WaitGroup
	for i := range observations {
		ref := projectionRef(projectionCases[1])
		ref.Locator.RuleID = domainText("rule-" + strconv.Itoa(i))
		raw := projectionJSON(t, ref)
		wg.Go(func() { emitter.EmitDomain(context.Background(), raw, domainScope()) })
	}
	wg.Wait()
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	s, _ := b.Snapshot(Query{})
	if uint64(len(s.Records))+s.CaptureDropped != observations {
		t.Fatalf("capture loss not conserved: got %d records and %d lost", len(s.Records), s.CaptureDropped)
	}
	var prior uint64
	for _, e := range s.Records {
		n, _ := strconv.ParseUint(e.Sequence, 10, 64)
		if n <= prior {
			t.Fatal("domain gap fold broke sequence ordering")
		}
		prior = n
	}
}

func TestDomainBrokerSourceRecoveryClosesOnlyExactObservedFault(t *testing.T) {
	b := pausedDomainBroker(t)
	emitter := mountedEmitter(t, b)
	fault := projectionRef(projectionCases[1])
	applied := projectionRef(projectionCases[0])
	applied.Evidence.ReferenceID = "applied-E"
	for range 3 {
		emitter.EmitDomain(t.Context(), projectionJSON(t, fault), domainScope())
	}
	emitter.EmitDomain(t.Context(), projectionJSON(t, applied), domainScope())
	if b.queue.count != 3 || b.queue.at(1).event.Code != "domain.source.rejected" || b.queue.at(1).event.Fields["occurrences"] != "3" || b.queue.at(2).event.Code != "domain.source.applied" {
		t.Fatalf("recovery must place final old-fault summary before applied, got count %d", b.queue.count)
	}
	// Repeated confirmed applied observations are not new recoveries or first spam.
	emitter.EmitDomain(t.Context(), projectionJSON(t, applied), domainScope())
	if b.queue.count != 3 {
		t.Fatal("repeated applied observation resent first/recovery event")
	}
	b.domain.mu.Lock()
	for _, entry := range b.domain.entries {
		entry.next = time.Now().Add(-time.Second)
	}
	b.domain.mu.Unlock()
	b.domain.Tick(t.Context())
	for i := 3; i < b.queue.count; i++ {
		if b.queue.at(i).event.Code == "domain.source.rejected" {
			t.Fatal("old rejected episode emitted after applied")
		}
	}
	if b.queue.count != 4 || b.queue.at(3).event.Fields["occurrences"] != "2" {
		t.Fatal("post-recovery success repeat did not retain bounded summary")
	}
}

func TestDomainBrokerSourceRecoveryDoesNotClearOtherScopeOrRevision(t *testing.T) {
	for _, change := range []string{"workspace", "workspace-revision", "device", "rule", "entity-revision"} {
		t.Run(change, func(t *testing.T) {
			b := pausedDomainBroker(t)
			emitter := mountedEmitter(t, b)
			fault := projectionRef(projectionCases[1])
			applied := projectionRef(projectionCases[0])
			scope := domainScope()
			for range 3 {
				emitter.EmitDomain(t.Context(), projectionJSON(t, fault), scope)
			}
			switch change {
			case "workspace":
				scope.WorkspaceID = "W2"
				applied.Binding.WorkspaceID = "W2"
				applied.Locator.WorkspaceID = "W2"
			case "workspace-revision":
				scope.WorkspaceRevision = "W2"
				applied.Binding.WorkspaceRevision = "W2"
			case "device":
				scope.DeviceID = "D2"
				applied.Locator.DeviceID = domainText("D2")
			case "rule":
				applied.Locator.RuleID = domainText("S2")
			case "entity-revision":
				applied.Binding.EntityRevision = "R2"
			}
			emitter.EmitDomain(t.Context(), projectionJSON(t, applied), scope)
			if b.queue.count != 2 {
				t.Fatal("foreign/revision applied observation cleared another fault")
			}
			b.domain.mu.Lock()
			for _, entry := range b.domain.entries {
				entry.next = time.Now().Add(-time.Second)
			}
			b.domain.mu.Unlock()
			b.domain.Tick(t.Context())
			if b.queue.count != 3 || b.queue.at(2).event.Code != "domain.source.rejected" || b.queue.at(2).event.Fields["occurrences"] != "3" {
				t.Fatal("other-scope/revision fault was falsely recovered")
			}
		})
	}
}

func TestDomainBrokerSourceRecoveryRejectsReplacedEpisode(t *testing.T) {
	b := pausedDomainBroker(t)
	emitter := mountedEmitter(t, b)
	fault := projectionRef(projectionCases[1])
	applied := projectionRef(projectionCases[0])
	emitter.EmitDomain(t.Context(), projectionJSON(t, fault), domainScope())
	newer := fault
	newer.Evidence.ReferenceID = "N"
	var replaced bool
	b.domain.options.Clock = func() time.Time {
		// Inject replacement after the adapter snapshots the old immutable reference,
		// exactly when existing Recover reads its clock before checking current state.
		if !replaced {
			replaced = true
			b.domain.mu.Lock()
			entry := b.domain.entries[domainCoalescerKey(fault)]
			updated := *entry
			updated.data = projectionJSON(t, newer)
			updated.reference = "N"
			b.domain.entries[domainCoalescerKey(fault)] = &updated
			b.domain.mu.Unlock()
		}
		return time.Now()
	}
	if emitter.EmitDomain(t.Context(), projectionJSON(t, applied), domainScope()) {
		t.Fatal("stale recovery snapshot was accepted")
	}
	if b.queue.count != 1 || b.domain.Stats().Dropped != 1 || b.domain.entries[domainCoalescerKey(newer)].reference != "N" {
		t.Fatal("replaced rejection episode was cleared or applied emitted")
	}
}
