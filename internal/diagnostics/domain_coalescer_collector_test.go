package diagnostics

import (
	"context"
	"math"
	"strconv"
	"strings"
	"testing"
	"time"
)

func collectorPair() (failed, recovered DomainEventRef) {
	failed, recovered = projectionRef(projectionCases[2]), projectionRef(projectionCases[3])
	failed.Evidence.ReferenceID = "acquisition-failed"
	recovered.Evidence.ReferenceID = "acquisition-good"
	return failed, recovered
}

func TestCollectorDomainBrokerGoodWithoutExactFaultHasNoSuccessSpam(t *testing.T) {
	b := pausedDomainBroker(t)
	_, good := collectorPair()
	for range AdmissionRecords + 10 {
		if b.EmitDomain(t.Context(), projectionJSON(t, good), domainScope()) {
			t.Fatal("healthy poll reported a captured recovery without fault")
		}
	}
	if b.queue.count != 0 || b.domain.Stats().Entries != 0 || b.domain.Stats().Dropped != 0 {
		t.Fatalf("healthy observations retained state/events or invented loss: %#v", b.domain.Stats())
	}
}

func TestCollectorDomainBrokerFinalSummaryThenOneRecovery(t *testing.T) {
	b := pausedDomainBroker(t)
	bad, good := collectorPair()
	for range 3 {
		b.EmitDomain(t.Context(), projectionJSON(t, bad), domainScope())
	}
	if !b.EmitDomain(t.Context(), projectionJSON(t, good), domainScope()) {
		t.Fatal("confirmed good result did not recover exact retained failure")
	}
	if b.queue.count != 3 || b.queue.at(1).event.Code != "domain.collector.failed" || b.queue.at(1).event.Fields["occurrences"] != "3" || b.queue.at(1).event.Fields["reference_id"] != "acquisition-failed" || b.queue.at(2).event.Code != "domain.collector.recovered" || b.queue.at(2).event.Fields["reference_id"] != "acquisition-good" {
		t.Fatal("ordered original failure summary / actual recovery reference lost")
	}
	for range 10 {
		b.EmitDomain(t.Context(), projectionJSON(t, good), domainScope())
	}
	if b.queue.count != 3 || b.domain.Stats().Entries != 0 || b.domain.Stats().Dropped != 0 {
		t.Fatal("healthy repeats retained success state or emitted spam")
	}
}

func TestCollectorDomainBrokerRecoveryRequiresFullScopeAndOriginalRevision(t *testing.T) {
	for _, field := range []string{"workspace", "workspace-revision", "device", "point", "revision"} {
		t.Run(field, func(t *testing.T) {
			b := pausedDomainBroker(t)
			bad, good := collectorPair()
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
			case "point":
				good.PointID = domainText("P2")
			case "revision":
				good.Binding.EntityRevision = "R2"
			}
			b.EmitDomain(t.Context(), projectionJSON(t, good), scope)
			if b.queue.count != 1 || b.domain.Stats().Entries != 1 || b.domain.Stats().Dropped != 0 {
				t.Fatal("unrelated good result recovered another fault, spammed, or counted healthy loss")
			}
		})
	}
}

func TestCollectorDomainBrokerFilteredRecoveryRetiresWarnWithoutInfo(t *testing.T) {
	b := pausedDomainBroker(t)
	b.level = levelRank(levelWarn)
	bad, good := collectorPair()
	for range 3 {
		b.EmitDomain(t.Context(), projectionJSON(t, bad), domainScope())
	}
	if b.EmitDomain(t.Context(), projectionJSON(t, good), domainScope()) {
		t.Fatal("filtered info recovery reported emitted")
	}
	if b.queue.count != 2 || b.queue.at(1).event.Fields["occurrences"] != "3" || b.domain.Stats().Entries != 0 || b.domain.Stats().Dropped != 0 {
		t.Fatal("filtered recovery left fault alive or retained info/loss")
	}
}

func TestCollectorDomainBrokerUnavailableBurstIsOneBoundedIdentityFreeEntry(t *testing.T) {
	b := pausedDomainBroker(t)
	clock := &coalescerClock{at: time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC)}
	b.domain.options.Clock = clock.Now
	for range 5000 {
		b.EmitDomainUnavailable(t.Context())
	}
	if b.queue.count != 1 || b.domain.Stats().Entries != 1 || b.domain.Stats().Dropped != 0 || len(b.queue.first().event.Fields) != 0 {
		t.Fatalf("unknown proof burst spams or retains identity: queue=%d stats=%#v", b.queue.count, b.domain.Stats())
	}
	clock.Advance(RuntimeLogHeartbeat)
	b.domain.Tick(t.Context())
	if b.queue.count != 2 {
		t.Fatal("shared heartbeat did not emit unavailable summary")
	}
	e := b.queue.at(1).event
	if e.Code != "domain.unavailable" || len(e.Fields) != 3 || e.Fields["occurrences"] != "5000" || e.Fields["first_seen_at"] != "2026-10-29T12:00:00Z" || e.Fields["last_seen_at"] != "2026-10-29T12:00:00Z" {
		t.Fatalf("unavailable summary retained anything but safe times/count: %#v", e)
	}
	clock.Advance(RuntimeLogHeartbeat)
	b.domain.Tick(t.Context())
	if b.queue.count != 2 {
		t.Fatal("no new observations became heartbeat spam")
	}
}

func TestCollectorDomainBrokerUnavailableFilteringAndLossAreDistinct(t *testing.T) {
	b := pausedDomainBroker(t)
	b.level = levelRank(levelError)
	for range 100 {
		b.EmitDomainUnavailable(t.Context())
	}
	if b.domain.Stats().Entries != 0 || b.domain.Stats().Dropped != 0 {
		t.Fatal("intentional filtering retained unavailable state or loss")
	}
	b.level = levelRank(levelInfo)
	b.domain.mu.Lock()
	done := make(chan struct{})
	go func() { b.EmitDomainUnavailable(context.Background()); close(done) }()
	select {
	case <-done:
	case <-time.After(time.Second):
		b.domain.mu.Unlock()
		t.Fatal("unavailable capture blocked on lock")
	}
	b.domain.mu.Unlock()
	s, _ := b.Snapshot(Query{})
	if s.CaptureDropped != 1 || b.domain.Stats().Dropped != 1 || b.queue.count != 0 {
		t.Fatal("unavailable contention loss absent or double counted")
	}
}

func TestCollectorDomainBrokerUnavailableSharesOriginalCapacity(t *testing.T) {
	b := pausedDomainBroker(t)
	for i := range AdmissionRecords - 1 {
		v := domainFixture()
		v.Locator.GroupID = domainText("G" + strconv.Itoa(i))
		b.EmitDomain(t.Context(), projectionJSON(t, v), domainScope())
	}
	if b.EmitDomainUnavailable(t.Context()) {
		t.Fatal("unavailable bypassed the existing coalescer capacity")
	}
	s, _ := b.Snapshot(Query{})
	if b.domain.Stats().Entries != AdmissionRecords-1 || b.domain.Stats().Dropped != 1 || s.CaptureDropped != 1 {
		t.Fatal("shared capacity loss changed or new map was allocated")
	}
}

func TestCollectorDomainBrokerUnavailableCountOverflowIsLossNotWrap(t *testing.T) {
	b := pausedDomainBroker(t)
	b.EmitDomainUnavailable(t.Context())
	if b.domain.Stats().Entries != 1 {
		t.Fatal("unavailable does not use shared bounded state")
	}
	for _, entry := range b.domain.entries {
		entry.occurrences = math.MaxUint64
		entry.reported = math.MaxUint64
	}
	if b.EmitDomainUnavailable(t.Context()) {
		t.Fatal("occurrence overflow accepted")
	}
	if b.queue.count != 1 || b.domain.Stats().Dropped != 1 {
		t.Fatal("count overflow wrapped or double counted")
	}
	raw := projectionJSON(t, b.queue.first().event)
	if strings.Contains(string(raw), "workspace") || strings.Contains(string(raw), "point_id") {
		t.Fatal("unavailable borrowed identity")
	}
}

func TestCollectorDomainBrokerRecoveryRejectsReplacedEpisode(t *testing.T) {
	b := pausedDomainBroker(t)
	bad, good := collectorPair()
	b.EmitDomain(t.Context(), projectionJSON(t, bad), domainScope())
	newer := bad
	newer.Evidence.ReferenceID = strings.Repeat("N", len(bad.Evidence.ReferenceID))
	replaced := false
	b.domain.options.Clock = func() time.Time {
		if !replaced {
			replaced = true
			b.domain.mu.Lock()
			key := domainCoalescerKey(bad)
			entry := b.domain.entries[key]
			updated := *entry
			updated.data = projectionJSON(t, newer)
			updated.reference = newer.Evidence.ReferenceID
			b.domain.entries[key] = &updated
			b.domain.mu.Unlock()
		}
		return time.Now()
	}
	if b.EmitDomain(t.Context(), projectionJSON(t, good), domainScope()) {
		t.Fatal("old observation cleared a replaced failure episode")
	}
	if b.queue.count != 1 || b.domain.Stats().Dropped != 1 || b.domain.entries[domainCoalescerKey(newer)].reference != newer.Evidence.ReferenceID {
		t.Fatal("replacement race emitted false recovery or changed new fault")
	}
}

func TestCollectorDomainBrokerRevisionEpisodesAndLateTickStaySeparate(t *testing.T) {
	b := pausedDomainBroker(t)
	bad, good := collectorPair()
	r2 := bad
	r2.Binding.EntityRevision = "R2"
	r2.Evidence.ReferenceID = "second-revision"
	for range 3 {
		b.EmitDomain(t.Context(), projectionJSON(t, bad), domainScope())
	}
	old := b.domain.entries[domainCoalescerKey(bad)]
	prepared := domainEntryEvent(*old, true)
	b.EmitDomain(t.Context(), projectionJSON(t, r2), domainScope())
	b.EmitDomain(t.Context(), projectionJSON(t, good), domainScope())
	if b.queue.count != 4 || b.domain.Stats().Entries != 1 || b.domain.entries[domainCoalescerKey(r2)] == nil {
		t.Fatal("R1 recovery cleared R2 or lost final R1 summary")
	}
	if b.domain.enqueueCurrent(t.Context(), prepared, domainCoalescerKey(bad), old) {
		t.Fatal("late Tick reopened recovered R1 failure")
	}
	next := bad
	next.Evidence.ReferenceID = "new-failure-episode"
	b.EmitDomain(t.Context(), projectionJSON(t, next), domainScope())
	if b.queue.count != 5 || b.queue.at(4).event.Fields["reference_id"] != "new-failure-episode" || b.queue.at(4).event.Fields["occurrences"] != "1" {
		t.Fatal("new failure reused recovered episode")
	}
}

func TestCollectorDomainBrokerRecoveryReservationAndAdmissionLoss(t *testing.T) {
	for _, failure := range []string{"capacity", "contention", "admission"} {
		t.Run(failure, func(t *testing.T) {
			b := pausedDomainBroker(t)
			bad, good := collectorPair()
			for range 3 {
				b.EmitDomain(t.Context(), projectionJSON(t, bad), domainScope())
			}
			switch failure {
			case "capacity":
				b.domain.options.RecordLimit = 1
			case "contention":
				b.domain.mu.Lock()
			case "admission":
				for b.Emit(Input{Code: "startup.ready"}) {
				}
			}
			before := b.queue.count
			lost := b.dropped
			accepted := b.EmitDomain(t.Context(), projectionJSON(t, good), domainScope())
			if failure == "contention" {
				b.domain.mu.Unlock()
			}
			snapshot, _ := b.Snapshot(Query{})
			if failure == "admission" {
				if !accepted || b.domain.Stats().Entries != 0 || snapshot.CaptureDropped != lost+2 {
					t.Fatal("per-output admission loss double counted or recovered state retained")
				}
				b.EmitDomain(t.Context(), projectionJSON(t, good), domainScope())
				if b.domain.Stats().Dropped != 2 {
					t.Fatal("healthy repetition counted as additional loss")
				}
			} else if accepted || b.domain.Stats().Entries != 1 || snapshot.CaptureDropped != lost+1 {
				t.Fatal("failed reservation partially recovered or loss not single")
			}
			if b.queue.count != before {
				t.Fatal("rejected output bypassed bounded admission")
			}
		})
	}
}

func TestCollectorDomainBrokerUnknownProofAdmissionAndCloseLoss(t *testing.T) {
	b := pausedDomainBroker(t)
	b.mu.Lock()
	for range 50 {
		b.EmitDomainUnavailable(t.Context())
	}
	b.mu.Unlock()
	snapshot, _ := b.Snapshot(Query{})
	if snapshot.CaptureDropped != 1 || b.domain.Stats().Entries != 1 || b.domain.Stats().Dropped != 1 || b.queue.count != 0 {
		t.Fatal("first unknown-output rejection confused with 49 coalesced observations")
	}
	b.domain.Close()
	snapshot, _ = b.Snapshot(Query{})
	if snapshot.CaptureDropped != 2 {
		t.Fatal("close must account one unreported unknown summary, not observations")
	}
	if b.EmitDomainUnavailable(t.Context()) {
		t.Fatal("closed unknown observer accepted")
	}
	snapshot, _ = b.Snapshot(Query{})
	if snapshot.CaptureDropped != 3 {
		t.Fatal("closed capture was not counted once")
	}
}

func TestCollectorDomainBrokerOtherReasonsAndHostileCandidatesRemainSeparate(t *testing.T) {
	b := pausedDomainBroker(t)
	bad, good := collectorPair()
	bad.Reason = "unsupported"
	bad.Action = "inspect-settings"
	b.EmitDomain(t.Context(), projectionJSON(t, bad), domainScope())
	b.EmitDomain(t.Context(), projectionJSON(t, good), domainScope())
	if b.queue.count != 1 || b.domain.Stats().Entries != 1 || b.domain.Stats().Dropped != 0 {
		t.Fatal("GOOD guessed recovery for another reason")
	}
	hostile := []byte(strings.TrimSuffix(string(projectionJSON(t, good)), "}") + `,"value":"private-sample","raw_error":"password-fixture"}`)
	b.EmitDomain(t.Context(), hostile, domainScope())
	e := b.queue.at(1).event
	if e.Code != codeRawSuppressed || len(e.Fields) != 0 {
		t.Fatal("collector path bypassed closed scalar projection")
	}
}
