package diagnostics

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func mappingPair() (failed, recovered DomainEventRef) {
	failed, recovered = projectionRef(projectionCases[4]), projectionRef(projectionCases[5])
	failed.Evidence.ReferenceID = "acquisition-failed"
	recovered.Evidence.ReferenceID = "acquisition-good"
	return failed, recovered
}

func TestDomainMappingHealthySamplesDoNotSpam(t *testing.T) {
	b := pausedDomainBroker(t)
	_, good := mappingPair()
	for range AdmissionRecords + 10 {
		if b.EmitDomain(t.Context(), projectionJSON(t, good), domainScope()) {
			t.Fatal("healthy sample captured recovery without an exact fault")
		}
	}
	if b.queue.count != 0 || b.domain.Stats().Entries != 0 || b.domain.Stats().Dropped != 0 {
		t.Fatalf("healthy samples created state/events/loss: %#v", b.domain.Stats())
	}
}

func TestDomainMappingFinalBlockedSummaryBeforeAcquisitionRecovery(t *testing.T) {
	b := pausedDomainBroker(t)
	bad, good := mappingPair()
	for range 3 {
		b.EmitDomain(t.Context(), projectionJSON(t, bad), domainScope())
	}
	if !b.EmitDomain(t.Context(), projectionJSON(t, good), domainScope()) {
		t.Fatal("actual acquisition result did not recover exact blocked failure")
	}
	if b.queue.count != 3 {
		t.Fatalf("expected first/summary/recovery, got %d", b.queue.count)
	}
	summary, recovery := b.queue.at(1).event, b.queue.at(2).event
	if summary.Code != "domain.mapping.failed" || summary.Fields["stage"] != "blocked" || summary.Fields["occurrences"] != "3" || summary.Fields["reference_id"] != "acquisition-failed" || recovery.Code != "domain.mapping.recovered" || recovery.Fields["stage"] != "acquisition" || recovery.Fields["reference_id"] != "acquisition-good" || recovery.Fields["mapping_id"] != "M" || recovery.Fields["point_id"] != "P" || recovery.Fields["tag_id"] != "T" {
		t.Fatal("summary/recovery lost stage, count or exact executed tuple")
	}
	for range 100 {
		b.EmitDomain(t.Context(), projectionJSON(t, good), domainScope())
	}
	if b.queue.count != 3 || b.domain.Stats().Entries != 0 || b.domain.Stats().Dropped != 0 {
		t.Fatal("repeated healthy samples retained success state or spam/loss")
	}
}

func TestDomainMappingRecoveryRequiresExactRetainedScopeAndRevision(t *testing.T) {
	for _, field := range []string{"workspace", "workspace-revision", "device", "mapping", "point", "tag", "revision"} {
		t.Run(field, func(t *testing.T) {
			b := pausedDomainBroker(t)
			bad, good := mappingPair()
			scope := domainScope()
			b.EmitDomain(t.Context(), projectionJSON(t, bad), scope)
			switch field {
			case "workspace":
				scope.WorkspaceID = "W2"
				good.Binding.WorkspaceID = "W2"
				good.Locator.WorkspaceID = "W2"
			case "workspace-revision":
				scope.WorkspaceRevision = "W2"
				good.Binding.WorkspaceRevision = "W2"
			case "device":
				scope.DeviceID = "D2"
				good.Locator.DeviceID = domainText("D2")
			case "mapping":
				good.Locator.MappingID = domainText("M2")
			case "point":
				good.PointID = domainText("P2")
			case "tag":
				good.TagID = domainText("T2")
			case "revision":
				good.Binding.EntityRevision = "R2"
			}
			if _, err := ParseDomainEventRef(projectionJSON(t, good), scope); err != nil {
				t.Fatal(err)
			}
			b.EmitDomain(t.Context(), projectionJSON(t, good), scope)
			if b.queue.count != 1 || b.domain.Stats().Entries != 1 || b.domain.Stats().Dropped != 0 {
				t.Fatal("unrelated healthy tuple cleared fault or created event/state/loss")
			}
		})
	}
}

func TestDomainMappingFilteredRecoveryRetiresOnlyExactWarnFault(t *testing.T) {
	b := pausedDomainBroker(t)
	b.level = levelRank(levelWarn)
	bad, good := mappingPair()
	if b.EmitDomain(t.Context(), projectionJSON(t, good), domainScope()) || b.domain.Stats().Entries != 0 {
		t.Fatal("filtered no-fault success created state")
	}
	for range 3 {
		b.EmitDomain(t.Context(), projectionJSON(t, bad), domainScope())
	}
	if b.EmitDomain(t.Context(), projectionJSON(t, good), domainScope()) {
		t.Fatal("filtered info reported emitted")
	}
	if b.queue.count != 2 || b.queue.at(1).event.Code != "domain.mapping.failed" || b.queue.at(1).event.Fields["occurrences"] != "3" || b.domain.Stats().Entries != 0 || b.domain.Stats().Dropped != 0 {
		t.Fatal("filtered recovery retained stale fault/info or counted intentional filtering as loss")
	}
	b.domain.Tick(t.Context())
	if b.queue.count != 2 {
		t.Fatal("retired warn fault reappeared")
	}
}

func TestDomainMappingRecoveryRejectsReplacedEpisode(t *testing.T) {
	b := pausedDomainBroker(t)
	bad, good := mappingPair()
	b.EmitDomain(t.Context(), projectionJSON(t, bad), domainScope())
	newer := bad
	newer.Evidence.ReferenceID = strings.Repeat("N", len(bad.Evidence.ReferenceID))
	replaced := false
	b.domain.options.Clock = func() time.Time {
		if !replaced {
			replaced = true
			b.domain.mu.Lock()
			key := domainCoalescerKey(bad)
			entry := *b.domain.entries[key]
			entry.data, entry.reference = projectionJSON(t, newer), newer.Evidence.ReferenceID
			b.domain.entries[key] = &entry
			b.domain.mu.Unlock()
		}
		return time.Now()
	}
	if b.EmitDomain(t.Context(), projectionJSON(t, good), domainScope()) {
		t.Fatal("stale recovery cleared a newly replaced episode")
	}
	if !replaced || b.queue.count != 1 || b.domain.Stats().Dropped != 1 || b.domain.entries[domainCoalescerKey(newer)].reference != newer.Evidence.ReferenceID {
		t.Fatal("new fault token/reference was overwritten or race refusal not counted once")
	}
}

func TestDomainMappingRevisionEpisodesAndLateTickStaySeparate(t *testing.T) {
	b := pausedDomainBroker(t)
	bad, good := mappingPair()
	r2 := bad
	r2.Binding.EntityRevision, r2.Evidence.ReferenceID = "R2", "second-revision"
	for range 3 {
		b.EmitDomain(t.Context(), projectionJSON(t, bad), domainScope())
	}
	old := b.domain.entries[domainCoalescerKey(bad)]
	prepared := domainEntryEvent(*old, true)
	b.EmitDomain(t.Context(), projectionJSON(t, r2), domainScope())
	b.EmitDomain(t.Context(), projectionJSON(t, good), domainScope())
	if b.queue.count != 4 || b.domain.Stats().Entries != 1 || b.domain.entries[domainCoalescerKey(r2)] == nil || b.queue.at(2).event.Fields["occurrences"] != "3" {
		t.Fatal("R1 recovery cleared R2 or lost R1 final summary")
	}
	if b.domain.enqueueCurrent(t.Context(), prepared, domainCoalescerKey(bad), old) {
		t.Fatal("late Tick reopened recovered R1")
	}
	next := bad
	next.Evidence.ReferenceID = "new-episode"
	b.EmitDomain(t.Context(), projectionJSON(t, next), domainScope())
	if b.queue.count != 5 || b.queue.at(4).event.Fields["reference_id"] != "new-episode" || b.queue.at(4).event.Fields["occurrences"] != "1" {
		t.Fatal("new fault inherited recovered episode")
	}
}

func TestDomainMappingRecoveryReservationAndAdmissionLoss(t *testing.T) {
	for _, failure := range []string{"record-capacity", "byte-capacity", "contention", "admission"} {
		t.Run(failure, func(t *testing.T) {
			b := pausedDomainBroker(t)
			bad, good := mappingPair()
			for range 3 {
				b.EmitDomain(t.Context(), projectionJSON(t, bad), domainScope())
			}
			switch failure {
			case "record-capacity":
				b.domain.options.RecordLimit = 1
			case "byte-capacity":
				b.domain.options.ByteLimit = 1
			case "contention":
				b.domain.mu.Lock()
			case "admission":
				for b.Emit(Input{Code: "startup.ready"}) {
				}
			}
			before, lost := b.queue.count, b.dropped
			accepted := b.EmitDomain(t.Context(), projectionJSON(t, good), domainScope())
			if failure == "contention" {
				b.domain.mu.Unlock()
			}
			snapshot, _ := b.Snapshot(Query{})
			if failure == "admission" {
				if !accepted || b.domain.Stats().Entries != 0 || snapshot.CaptureDropped != lost+2 {
					t.Fatal("individual broker refusals lost/doubled or coalescer retained half-transition")
				}
				b.EmitDomain(t.Context(), projectionJSON(t, good), domainScope())
				if b.domain.Stats().Dropped != 2 {
					t.Fatal("healthy repetition counted as loss")
				}
			} else if accepted || b.domain.Stats().Entries != 1 || snapshot.CaptureDropped != lost+1 {
				t.Fatal("failed reservation partially transitioned or loss not single")
			}
			if b.queue.count != before {
				t.Fatal("refused output bypassed original admission cap")
			}
		})
	}
}

func TestDomainMappingCancelAndCloseNeverClearOrFakeRecovery(t *testing.T) {
	for _, mode := range []string{"cancel", "close-between-locks"} {
		t.Run(mode, func(t *testing.T) {
			b := pausedDomainBroker(t)
			bad, good := mappingPair()
			b.EmitDomain(t.Context(), projectionJSON(t, bad), domainScope())
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			if mode == "cancel" {
				cancel()
			} else {
				b.domain.options.Clock = func() time.Time { b.domain.Close(); return time.Now() }
			}
			if b.EmitDomain(ctx, projectionJSON(t, good), domainScope()) {
				t.Fatal("canceled/closed capture claimed recovery")
			}
			if b.queue.count != 1 || b.domain.Stats().Dropped != 1 {
				t.Fatal("closed/canceled result emitted recovery or double counted capture")
			}
			if mode == "cancel" && b.domain.Stats().Entries != 1 {
				t.Fatal("cancellation cleared fault")
			}
			if mode == "close-between-locks" && b.domain.Stats().Entries != 0 {
				t.Fatal("closed coalescer retained state")
			}
		})
	}
}

func TestDomainMappingCallbacksRunOutsideLockAndPreserveTransitionOrder(t *testing.T) {
	b := pausedDomainBroker(t)
	bad, good := mappingPair()
	for range 3 {
		b.EmitDomain(t.Context(), projectionJSON(t, bad), domainScope())
	}
	original := b.domain.options.TryAdmit
	reentered := false
	b.domain.options.TryAdmit = func(e Event) bool {
		if e.Code == "domain.mapping.failed" && e.Fields["occurrences"] == "3" {
			if !b.domain.mu.TryLock() {
				t.Error("callback invoked under state mutex")
				return false
			}
			b.domain.mu.Unlock()
			reentered = true
			next := bad
			next.Evidence.ReferenceID = "new-episode"
			b.EmitDomain(t.Context(), projectionJSON(t, next), domainScope())
		}
		return original(e)
	}
	b.EmitDomain(t.Context(), projectionJSON(t, good), domainScope())
	if !reentered || b.queue.count != 4 || b.queue.at(1).event.Fields["occurrences"] != "3" || b.queue.at(2).event.Code != "domain.mapping.recovered" || b.queue.at(3).event.Fields["reference_id"] != "new-episode" || b.domain.Stats().Entries != 1 {
		t.Fatal("reentry interleaved replacement fault into final summary/recovery pair")
	}
}

func TestDomainMappingOtherReasonsAndHostileCandidatesFailClosed(t *testing.T) {
	b := pausedDomainBroker(t)
	bad, good := mappingPair()
	bad.Reason, bad.Action = "unsupported", "inspect-settings"
	b.EmitDomain(t.Context(), projectionJSON(t, bad), domainScope())
	b.EmitDomain(t.Context(), projectionJSON(t, good), domainScope())
	if b.queue.count != 1 || b.domain.Stats().Entries != 1 || b.domain.Stats().Dropped != 0 {
		t.Fatal("GOOD guessed recovery for another reason")
	}
	hostile := []byte(strings.TrimSuffix(string(projectionJSON(t, good)), "}") + `,"raw_error":"password_fixture","params":{"value":"secret-transform"}}`)
	b.EmitDomain(t.Context(), hostile, domainScope())
	e := b.queue.at(1).event
	if e.Code != codeRawSuppressed || len(e.Fields) != 0 {
		t.Fatal("mapping path bypassed closed scalar projection")
	}
}

func TestDomainMappingHeartbeatCountsAndFirstLastRemainObserved(t *testing.T) {
	b := pausedDomainBroker(t)
	clock := &coalescerClock{at: time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC)}
	b.domain.options.Clock = clock.Now
	bad, good := mappingPair()
	b.EmitDomain(t.Context(), projectionJSON(t, bad), domainScope())
	clock.Advance(time.Second)
	b.EmitDomain(t.Context(), projectionJSON(t, bad), domainScope())
	clock.Advance(time.Second)
	b.EmitDomain(t.Context(), projectionJSON(t, bad), domainScope())
	clock.Advance(RuntimeLogHeartbeat - 2*time.Second)
	b.domain.Tick(t.Context())
	if b.queue.count != 2 {
		t.Fatal("existing heartbeat did not emit exactly one repeated-fault summary")
	}
	summary := b.queue.at(1).event
	if summary.Fields["occurrences"] != "3" || summary.Fields["first_seen_at"] != "2026-10-29T12:00:00Z" || summary.Fields["last_seen_at"] != "2026-10-29T12:00:02Z" {
		t.Fatal("summary timestamps/count no longer describe actual observations")
	}
	b.EmitDomain(t.Context(), projectionJSON(t, good), domainScope())
	if b.queue.count != 3 || b.queue.at(2).event.Code != "domain.mapping.recovered" || b.domain.Stats().Entries != 0 {
		t.Fatal("recovery duplicated an already reported summary or left fault pending")
	}
	clock.Advance(RuntimeLogHeartbeat)
	b.domain.Tick(t.Context())
	if b.queue.count != 3 || b.domain.Stats().Dropped != 0 {
		t.Fatal("recovered mapping kept emitting heartbeat/loss")
	}
}

func TestDomainMappingConcurrentRecoveryAndTickNeverDuplicateOrCrossTuple(t *testing.T) {
	b := pausedDomainBroker(t)
	sink := &coalescerSink{}
	b.domain.options.TryAdmit = sink.Admit
	const scopes = 32
	goods := make([][]byte, scopes)
	for i := range scopes {
		bad, good := mappingPair()
		id := "mapping-" + strconv.Itoa(i)
		bad.Locator.MappingID, good.Locator.MappingID = domainText(id), domainText(id)
		bad.Evidence.ReferenceID, good.Evidence.ReferenceID = "bad-"+id, "good-"+id
		raw := projectionJSON(t, bad)
		b.EmitDomain(t.Context(), raw, domainScope())
		b.EmitDomain(t.Context(), raw, domainScope())
		goods[i] = projectionJSON(t, good)
	}
	var wg sync.WaitGroup
	for _, raw := range goods {
		for range 2 {
			wg.Go(func() { b.EmitDomain(context.Background(), raw, domainScope()) })
		}
	}
	wg.Go(func() { b.domain.Tick(context.Background()) })
	wg.Wait()
	// Retry only observations refused by lock contention; do not invent healthy
	// events for the episodes already retired by either concurrent observation.
	for _, raw := range goods {
		b.EmitDomain(t.Context(), raw, domainScope())
	}
	b.domain.Tick(t.Context()) // flush bounded outputs deferred by contention
	recovered := map[string]int{}
	for _, e := range sink.Events() {
		if e.Code != "domain.mapping.recovered" {
			continue
		}
		id, ok := e.Fields["mapping_id"].(string)
		if !ok || e.Fields["reference_id"] != "good-"+id || e.Fields["point_id"] != "P" || e.Fields["tag_id"] != "T" || e.Fields["entity_revision"] != "R1" {
			t.Fatal("concurrent recovery borrowed another executed tuple/reference")
		}
		recovered[id]++
	}
	if len(recovered) != scopes || b.domain.Stats().Entries != 0 {
		t.Fatalf("matched final recovery/state incomplete: recovered=%d entries=%d", len(recovered), b.domain.Stats().Entries)
	}
	for _, count := range recovered {
		if count != 1 {
			t.Fatal("concurrent observation emitted duplicate recovery")
		}
	}
}
