package diagnostics

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestDomainBrokerMemoryFileAndSubscriberUseSameProjection(t *testing.T) {
	opts := fileOptions(t)
	sink, err := OpenFileSink(opts)
	if err != nil {
		t.Fatal(err)
	}
	b, err := New(Options{File: sink})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = b.Close(context.Background()) })
	emitter := mountedEmitter(t, b)
	sub, err := b.Subscribe(Query{}, "")
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Close()
	emitter.EmitDomain(t.Context(), sourceDomainJSON(t), domainScope())
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	hostile := []byte(strings.TrimSuffix(string(sourceDomainJSON(t)), "}") + `,"raw_error":"password_fixture","nested":{"credential":"nested_secret"}}`)
	emitter.EmitDomain(t.Context(), hostile, domainScope())
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	snapshot, _ := b.Snapshot(Query{})
	if len(snapshot.Records) != 2 {
		t.Fatalf("want 2 projected records: %#v", snapshot)
	}
	for _, want := range snapshot.Records {
		ctx, cancel := context.WithTimeout(t.Context(), time.Second)
		got, err := sub.Next(ctx)
		cancel()
		if err != nil {
			t.Fatal(err)
		}
		a, _ := json.Marshal(want)
		c, _ := json.Marshal(got)
		if string(a) != string(c) {
			t.Fatalf("subscriber changed projection: %s != %s", a, c)
		}
		sub.Ack(got.Cursor())
	}
	raw, err := os.ReadFile(opts.Path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), `"code":"domain.source.rejected"`) || !strings.Contains(string(raw), `"workspace_revision":"W1"`) || !strings.Contains(string(raw), `"code":"raw.suppressed"`) {
		t.Fatalf("file lost safe domain projection: %s", raw)
	}
	for _, secret := range []string{"password_fixture", "nested_secret"} {
		if strings.Contains(string(raw), secret) {
			t.Fatalf("file retained %s", secret)
		}
	}
}

func TestDomainBrokerOwnedHeartbeatAndShutdown(t *testing.T) {
	if RuntimeLogHeartbeat != 15*time.Second {
		t.Fatal("existing heartbeat was changed")
	}
	b := testBroker(t)
	emitter := mountedEmitter(t, b)
	sub, err := b.Subscribe(Query{}, "")
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Close()
	raw := sourceDomainJSON(t)
	for range 3 {
		emitter.EmitDomain(t.Context(), raw, domainScope())
	}
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	first, err := sub.Next(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	sub.Ack(first.Cursor())
	// The first observation may fall just after the process heartbeat phase; at
	// most the next two unchanged 15s ticks are needed, with no elapsed-time UI rule.
	ctx, cancel := context.WithTimeout(t.Context(), 2*RuntimeLogHeartbeat+2*time.Second)
	defer cancel()
	summary, err := sub.Next(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if summary.Fields["occurrences"] != "3" || summary.Fields["first_seen_at"] != first.Fields["first_seen_at"] || summary.Fields["last_seen_at"] == nil {
		t.Fatalf("owned ticker did not emit exact summary: %#v", summary)
	}
	sub.Ack(summary.Cursor())
	if err := b.Close(t.Context()); err != nil {
		t.Fatal(err)
	}
	select {
	case <-b.done:
	default:
		t.Fatal("owned dispatcher/timer did not terminate")
	}
	if stats := b.domain.Stats(); stats.Entries != 0 || stats.Queued != 0 || stats.Dropped != 0 {
		t.Fatalf("heartbeat/close invented loss or retained state: %#v", stats)
	}
}

func TestDomainBrokerSourceWorkspaceRevisionIsolatesOccurrencesAndRecovery(t *testing.T) {
	b := pausedDomainBroker(t)
	emitter := mountedEmitter(t, b)
	w1 := domainScope()
	w2 := w1
	w2.WorkspaceRevision = "W2"
	first := projectionRef(projectionCases[1])
	second := first
	second.Binding.WorkspaceRevision = "W2"
	second.Evidence.ReferenceID = "E2"
	for range 2 {
		emitter.EmitDomain(t.Context(), projectionJSON(t, first), w1)
	}
	for range 3 {
		emitter.EmitDomain(t.Context(), projectionJSON(t, second), w2)
	}
	if b.queue.count != 2 || b.domain.Stats().Entries != 2 {
		t.Fatal("W2 observation/reference was coalesced into W1")
	}
	applied := projectionRef(projectionCases[0])
	applied.Binding.WorkspaceRevision = "W2"
	applied.Evidence.ReferenceID = "A2"
	emitter.EmitDomain(t.Context(), projectionJSON(t, applied), w2)
	if b.queue.count != 4 || b.queue.at(2).event.Fields["workspace_revision"] != "W2" || b.queue.at(2).event.Fields["reference_id"] != "E2" || b.queue.at(2).event.Fields["occurrences"] != "3" || b.queue.at(3).event.Code != "domain.source.applied" {
		t.Fatal("W2 recovery lost own original reference/count")
	}
	b.domain.mu.Lock()
	for _, entry := range b.domain.entries {
		entry.next = time.Now().Add(-time.Second)
	}
	b.domain.mu.Unlock()
	b.domain.Tick(t.Context())
	if b.queue.count != 5 || b.queue.at(4).event.Code != "domain.source.rejected" || b.queue.at(4).event.Fields["workspace_revision"] != "W1" || b.queue.at(4).event.Fields["occurrences"] != "2" {
		t.Fatal("W2 applied cleared W1 fault or emitted later W2 fault")
	}
}

func TestDomainBrokerNewRejectionFinalizesPriorAppliedSummary(t *testing.T) {
	b := pausedDomainBroker(t)
	emitter := mountedEmitter(t, b)
	applied := projectionRef(projectionCases[0])
	fault := projectionRef(projectionCases[1])
	for range 3 {
		emitter.EmitDomain(t.Context(), projectionJSON(t, applied), domainScope())
	}
	emitter.EmitDomain(t.Context(), projectionJSON(t, fault), domainScope())
	if b.queue.count != 3 || b.queue.at(1).event.Code != "domain.source.applied" || b.queue.at(1).event.Fields["occurrences"] != "3" || b.queue.at(2).event.Code != "domain.source.rejected" {
		t.Fatal("old pending success was not finalized before new rejection")
	}
	b.domain.mu.Lock()
	for _, entry := range b.domain.entries {
		entry.next = time.Now().Add(-time.Second)
	}
	b.domain.mu.Unlock()
	b.domain.Tick(t.Context())
	if b.queue.count != 3 {
		t.Fatal("old success summary was emitted after new rejection")
	}
}

func TestDomainBrokerSourceTransitionCapacityAndContentionStayWhole(t *testing.T) {
	for _, failure := range []string{"state-capacity", "state-contention", "admission-full"} {
		t.Run(failure, func(t *testing.T) {
			b := pausedDomainBroker(t)
			emitter := mountedEmitter(t, b)
			fault := projectionRef(projectionCases[1])
			applied := projectionRef(projectionCases[0])
			raw := projectionJSON(t, fault)
			for range 3 {
				emitter.EmitDomain(t.Context(), raw, domainScope())
			}
			switch failure {
			case "state-capacity":
				b.domain.options.RecordLimit = 2 // first transition needs retained success plus both outputs
			case "state-contention":
				b.domain.mu.Lock()
			case "admission-full":
				for b.Emit(Input{Code: "startup.ready"}) {
				}
			}
			priorCount := b.queue.count
			priorLoss := b.dropped
			accepted := emitter.EmitDomain(t.Context(), projectionJSON(t, applied), domainScope())
			if failure == "state-contention" {
				b.domain.mu.Unlock()
			}
			snapshot, _ := b.Snapshot(Query{})
			if failure == "admission-full" {
				if !accepted || b.domain.Stats().Entries != 1 || b.domain.entries[domainCoalescerKey(fault)] != nil || snapshot.CaptureDropped != priorLoss+2 {
					t.Fatal("full admission lost/duplicated transition state or loss")
				}
				emitter.EmitDomain(t.Context(), projectionJSON(t, applied), domainScope())
				if b.queue.count != priorCount || b.domain.Stats().Dropped != 2 {
					t.Fatal("repeated applied resent lost recovery instead of coalescing")
				}
			} else {
				if accepted || b.queue.count != priorCount || b.domain.entries[domainCoalescerKey(fault)] == nil || b.domain.entries[domainCoalescerKey(applied)] != nil || snapshot.CaptureDropped != priorLoss+1 {
					t.Fatal("failed reservation partially cleared fault or emitted applied")
				}
			}
		})
	}
}

func TestDomainBrokerSourceAppliedDoesNotClaimOtherReasonRecovery(t *testing.T) {
	b := pausedDomainBroker(t)
	emitter := mountedEmitter(t, b)
	fault := projectionRef(projectionCases[1])
	fault.Reason = "unsupported"
	fault.Action = "inspect-settings"
	applied := projectionRef(projectionCases[0])
	for range 2 {
		emitter.EmitDomain(t.Context(), projectionJSON(t, fault), domainScope())
	}
	emitter.EmitDomain(t.Context(), projectionJSON(t, applied), domainScope())
	if b.domain.entries[domainCoalescerKey(fault)] == nil || b.queue.count != 2 {
		t.Fatal("fixed unavailable recovery guessed another reason episode")
	}
}

func TestDomainBrokerSourceTransitionInvalidatesLatePreparedSummaries(t *testing.T) {
	for _, firstApplied := range []bool{false, true} {
		t.Run(strconv.FormatBool(firstApplied), func(t *testing.T) {
			b := pausedDomainBroker(t)
			emitter := mountedEmitter(t, b)
			first, second := projectionRef(projectionCases[1]), projectionRef(projectionCases[0])
			if firstApplied {
				first, second = second, first
			}
			for range 2 {
				emitter.EmitDomain(t.Context(), projectionJSON(t, first), domainScope())
			}
			key := domainCoalescerKey(first)
			entry := b.domain.entries[key]
			prepared := domainEntryEvent(*entry, true)
			emitter.EmitDomain(t.Context(), projectionJSON(t, second), domainScope())
			before := b.queue.count
			if b.domain.enqueueCurrent(t.Context(), prepared, key, entry) {
				t.Fatal("late prepared summary reopened completed episode")
			}
			if b.queue.count != before || b.domain.Stats().Dropped != 1 {
				t.Fatal("late prepared capture was admitted or loss missing")
			}
		})
	}
}

func TestDomainBrokerFilteredSourceApplyFinalizesFaultWithoutRetainingInfo(t *testing.T) {
	b := pausedDomainBroker(t)
	b.level = levelRank(levelWarn)
	emitter := mountedEmitter(t, b)
	fault, applied := projectionRef(projectionCases[1]), projectionRef(projectionCases[0])
	for range 3 {
		emitter.EmitDomain(t.Context(), projectionJSON(t, fault), domainScope())
	}
	if emitter.EmitDomain(t.Context(), projectionJSON(t, applied), domainScope()) {
		t.Fatal("filtered source output reported admitted")
	}
	if b.queue.count != 2 || b.queue.at(1).event.Code != "domain.source.rejected" || b.queue.at(1).event.Fields["occurrences"] != "3" {
		t.Fatal("filtered applied observation did not finalize its pending warn summary")
	}
	if stats := b.domain.Stats(); stats.Entries != 0 || stats.StateBytes != 0 || stats.Dropped != 0 {
		t.Fatalf("filtered transition retained info state or counted intentional loss: %#v", stats)
	}
	b.domain.Tick(t.Context())
	for range AdmissionRecords + 1 {
		if emitter.EmitDomain(t.Context(), projectionJSON(t, applied), domainScope()) {
			t.Fatal("filtered no-match observation accepted")
		}
	}
	if b.queue.count != 2 || b.domain.Stats().Entries != 0 || b.domain.Stats().Dropped != 0 {
		t.Fatal("filtered no-match events grew state, emitted, or counted as loss")
	}
}
