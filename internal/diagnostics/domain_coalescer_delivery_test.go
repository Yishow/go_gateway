package diagnostics

import (
	"context"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/runtimeview"
)

func deliveryCommit(effect string) DomainEventRef {
	v := domainCommitFixture()
	v.Evidence.ReferenceID = effect
	v.Evidence.Receipt.EffectKey, v.Evidence.Receipt.RecordID, v.Evidence.Receipt.PayloadDigest = effect, "record-"+effect, "digest-"+effect
	return v
}
func deliveryFault(effect, code string) DomainEventRef {
	v := domainFixture()
	v.Evidence.ReferenceID, v.Code = effect, DomainEventCode("domain.delivery."+code)
	v.Stage = runtimeview.Stage(code)
	return v
}
func emitDelivery(t *testing.T, b *Broker, v DomainEventRef) bool {
	t.Helper()
	return b.EmitDomain(t.Context(), projectionJSON(t, v), domainScope())
}
func deliveryEvents(b *Broker, code string) []Event {
	var events []Event
	for i := 0; i < b.queue.count; i++ {
		e := b.queue.at(i).event
		if e.Code == code {
			events = append(events, e)
		}
	}
	return events
}
func recordingRefusal() DomainEventRef {
	v := domainFixture()
	v.Code, v.Stage, v.Reason = "domain.recording.refused", runtimeview.StageBlocked, runtimeview.ReasonOwnershipUnavailable
	v.Action = runtimeview.ActionForReason(v.Reason)
	v.Evidence = DomainEvidenceRef{Kind: "recording-refusal", ReferenceID: "observation-uuid", PendingSamples: domainText("0")}
	return v
}

func TestDomainRecordingRefusalHasClosedOriginalInputEvidence(t *testing.T) {
	for _, reason := range []runtimeview.Reason{runtimeview.ReasonOwnershipUnavailable, runtimeview.ReasonOwnershipConflict, runtimeview.ReasonHandoffConflict, runtimeview.ReasonFinalizedOverlap} {
		t.Run(string(reason), func(t *testing.T) {
			v := recordingRefusal()
			v.Reason, v.Action = reason, runtimeview.ActionForReason(reason)
			v.Evidence.PendingSamples = domainText("18446744073709551615")
			raw := projectionJSON(t, v)
			if _, err := ParseDomainEventRef(raw, domainScope()); err != nil {
				t.Fatalf("confirmed typed original-scope refusal unavailable: %v", err)
			}
			e := ProjectDomainEvent(raw, domainScope())
			if e.Code != "domain.recording.refused" || e.Fields["pending_samples"] != "18446744073709551615" || e.Fields["evidence_kind"] != "recording-refusal" || e.Fields["reference_id"] != "observation-uuid" || e.Fields["reason"] != string(reason) || e.Fields["effect_key"] != nil || e.Fields["record_id"] != nil {
				t.Fatal("refusal became durable output/receipt or lost exact input count")
			}
		})
	}
}
func TestDomainRecordingRefusalRejectsUnknownReasonOrMixedEvidence(t *testing.T) {
	for _, kind := range []string{"reason", "missing-count", "negative-count", "leading-zero", "receipt", "wrong-kind", "device-source", "raw-error"} {
		t.Run(kind, func(t *testing.T) {
			v := recordingRefusal()
			switch kind {
			case "reason":
				v.Reason, v.Action = runtimeview.ReasonUnavailable, runtimeview.ActionRefresh
			case "missing-count":
				v.Evidence.PendingSamples = nil
			case "negative-count":
				v.Evidence.PendingSamples = domainText("-1")
			case "leading-zero":
				v.Evidence.PendingSamples = domainText("01")
			case "receipt":
				v.Evidence.Receipt = deliveryCommit("E").Evidence.Receipt
			case "wrong-kind":
				v.Evidence.Kind = "durable-transition"
			case "device-source":
				v.PointID = domainText("P")
			}
			raw := projectionJSON(t, v)
			if kind == "raw-error" {
				raw = []byte(strings.TrimSuffix(string(raw), "}") + `,"raw_error":"password_fixture","params":{"secret":"nested-secret"}}`)
			}
			requireProjectionSuppressed(t, ProjectDomainEvent(raw, domainScope()))
		})
	}
}
func TestDomainDeliveryReceiptProofMismatchCannotRecover(t *testing.T) {
	for _, kind := range []string{"missing", "effect", "group-revision", "connector-revision", "provenance", "timestamp"} {
		t.Run(kind, func(t *testing.T) {
			b := pausedDomainBroker(t)
			emitDelivery(t, b, deliveryFault("E", "retrying"))
			v := deliveryCommit("E")
			switch kind {
			case "missing":
				v.Evidence.Receipt = nil
			case "effect":
				v.Evidence.Receipt.EffectKey = "foreign"
			case "group-revision":
				v.Evidence.Receipt.GroupRevision = "R0"
			case "connector-revision":
				v.Evidence.Receipt.ConnectorRevision = "C0"
			case "provenance":
				v.Evidence.Receipt.Provenance = "test-write"
			case "timestamp":
				v.Evidence.Receipt.CommittedAt = "yesterday password_fixture"
			}
			emitDelivery(t, b, v)
			if len(deliveryEvents(b, "domain.delivery.recovered")) != 0 || len(deliveryEvents(b, "domain.delivery.first-confirmed")) != 0 || b.domain.Stats().Entries != 1 {
				t.Fatal("unmatched/unknown receipt cleared fault or established commitment")
			}
			if e := b.queue.at(1).event; e.Code != codeRawSuppressed || len(e.Fields) != 0 {
				t.Fatal("invalid receipt retained unsafe identity or proof")
			}
		})
	}
}
func TestDomainDeliveryFaultKeysKeepEffectsIndependent(t *testing.T) {
	b := pausedDomainBroker(t)
	clock := &coalescerClock{at: time.Now()}
	b.domain.options.Clock = clock.Now
	for range 3 {
		emitDelivery(t, b, deliveryFault("E1", "retrying"))
	}
	emitDelivery(t, b, deliveryFault("E2", "retrying"))
	if b.queue.count != 2 || b.domain.Stats().Entries != 2 {
		t.Fatal("different durable effects merged under group revision")
	}
	clock.Advance(RuntimeLogHeartbeat)
	b.domain.Tick(t.Context())
	if b.queue.count != 3 || b.queue.at(2).event.Fields["reference_id"] != "E1" || b.queue.at(2).event.Fields["occurrences"] != "3" {
		t.Fatal("summary borrowed another effect occurrence")
	}
}
func TestDomainDeliveryRecoveryRetiresAllExactEffectFaultCodes(t *testing.T) {
	b := pausedDomainBroker(t)
	for _, code := range []string{"retrying", "blocked", "quarantined", "unknown"} {
		for range 3 {
			emitDelivery(t, b, deliveryFault("E1", code))
		}
	}
	emitDelivery(t, b, deliveryFault("E2", "retrying"))
	if !emitDelivery(t, b, deliveryCommit("E1")) {
		t.Fatal("new confirmed transition refused")
	}
	if b.queue.count != 11 || b.domain.Stats().Entries != 2 {
		t.Fatalf("all-code recovery did not retain only E2 + success state: events=%d stats=%#v", b.queue.count, b.domain.Stats())
	}
	for i := 5; i < 9; i++ {
		e := b.queue.at(i).event
		if e.Fields["reference_id"] != "E1" || e.Fields["occurrences"] != "3" {
			t.Fatal("pending failure summary lost before recovery")
		}
	}
	recovered := deliveryEvents(b, "domain.delivery.recovered")
	if len(recovered) != 1 || recovered[0].Fields["effect_key"] != "E1" || len(deliveryEvents(b, "domain.delivery.first-confirmed")) != 1 {
		t.Fatal("recovery guessed another effect or duplicated per fault code")
	}
	if b.domain.entries[domainCoalescerKey(deliveryFault("E2", "retrying"))] == nil {
		t.Fatal("E1 recovery cleared E2")
	}
}
func TestDomainDeliveryRecoveryMatchesExactEffectAndRevision(t *testing.T) {
	for _, field := range []string{"workspace", "workspace-revision", "group", "group-revision", "connector", "connector-revision", "effect"} {
		t.Run(field, func(t *testing.T) {
			b := pausedDomainBroker(t)
			bad, good := deliveryFault("E1", "retrying"), deliveryCommit("E1")
			scope := domainScope()
			emitDelivery(t, b, bad)
			switch field {
			case "workspace":
				scope.WorkspaceID = "W2"
				good.Binding.WorkspaceID = "W2"
				good.Locator.WorkspaceID = "W2"
			case "workspace-revision":
				scope.WorkspaceRevision = "W2"
				good.Binding.WorkspaceRevision = "W2"
			case "group":
				good.Locator.GroupID = domainText("G2")
			case "group-revision":
				good.Binding.EntityRevision = "R2"
				good.Locator.GroupRevision = domainText("R2")
				good.Evidence.Receipt.GroupRevision = "R2"
			case "connector":
				good.ConnectorID = domainText("C2")
			case "connector-revision":
				good.ConnectorRevision = domainText("C2")
				good.Evidence.Receipt.ConnectorRevision = "C2"
			case "effect":
				good = deliveryCommit("E2")
			}
			b.EmitDomain(t.Context(), projectionJSON(t, good), scope)
			if len(deliveryEvents(b, "domain.delivery.recovered")) != 0 || b.domain.entries[domainCoalescerKey(bad)] == nil || b.domain.Stats().Entries != 2 {
				t.Fatal("foreign effect/scope/revision recovered retained fault")
			}
		})
	}
}
func TestDomainDeliveryFirstConfirmationThenBoundedSummary(t *testing.T) {
	b := pausedDomainBroker(t)
	clock := &coalescerClock{at: time.Date(2026, 10, 29, 13, 0, 0, 0, time.UTC)}
	b.domain.options.Clock = clock.Now
	for _, effect := range []string{"E1", "E2", "E3"} {
		emitDelivery(t, b, deliveryCommit(effect))
		clock.Advance(time.Second)
	}
	if b.queue.count != 1 || b.domain.Stats().Entries != 1 {
		t.Fatal("new confirmed effects generated per-row first events")
	}
	first := b.queue.first().event
	if !strings.Contains(first.Message, "newly observed") || !strings.Contains(first.Message, "capture scope") || !strings.Contains(first.Message, "only its own effect") {
		t.Fatal("first-confirmed incorrectly claims first-ever or aggregate receipt proof")
	}
	clock.Advance(RuntimeLogHeartbeat)
	b.domain.Tick(t.Context())
	if b.queue.count != 2 {
		t.Fatal("missing bounded confirmation summary")
	}
	e := b.queue.at(1).event
	if e.Code != "domain.delivery.summary" || e.Fields["occurrences"] != "3" || e.Fields["effect_key"] != "E1" || e.Fields["record_id"] != "record-E1" || e.Fields["last_seen_at"] != "2026-10-29T13:00:02Z" || !strings.Contains(e.Message, "only its own effect") {
		t.Fatal("summary swapped receipt or claimed SQL aggregate count")
	}
}
func TestDomainDeliveryFilteredCommitRetiresWarnFaultsWithoutSuccessState(t *testing.T) {
	b := pausedDomainBroker(t)
	b.level = levelRank(levelWarn)
	for _, code := range []string{"retrying", "blocked"} {
		for range 3 {
			emitDelivery(t, b, deliveryFault("E", code))
		}
	}
	if emitDelivery(t, b, deliveryCommit("E")) {
		t.Fatal("filtered info reported visible")
	}
	if b.queue.count != 4 || b.domain.Stats().Entries != 0 || b.domain.Stats().Dropped != 0 {
		t.Fatal("filtered confirmation retained warning fault, success state or loss")
	}
	emitDelivery(t, b, deliveryCommit("healthy"))
	if b.domain.Stats().Entries != 0 || b.domain.Stats().Dropped != 0 {
		t.Fatal("filtered healthy confirmation created state/loss")
	}
}
func TestDomainDeliveryReservationFailureKeepsCompleteFaultSet(t *testing.T) {
	for _, mode := range []string{"record-cap", "byte-cap", "contention"} {
		t.Run(mode, func(t *testing.T) {
			b := pausedDomainBroker(t)
			for _, code := range []string{"retrying", "blocked"} {
				for range 3 {
					emitDelivery(t, b, deliveryFault("E", code))
				}
			}
			switch mode {
			case "record-cap":
				b.domain.options.RecordLimit = 2
			case "byte-cap":
				b.domain.options.ByteLimit = 1
			case "contention":
				b.domain.mu.Lock()
			}
			accepted := emitDelivery(t, b, deliveryCommit("E"))
			if mode == "contention" {
				b.domain.mu.Unlock()
			}
			if accepted || b.queue.count != 2 || b.domain.Stats().Entries != 2 || b.domain.Stats().Dropped != 1 {
				t.Fatal("failed atomic reservation partially retired faults or created success")
			}
		})
	}
}
func TestDomainDeliveryChangedFaultSetTokenCannotBeCleared(t *testing.T) {
	for _, mode := range []string{"replace", "new-code"} {
		t.Run(mode, func(t *testing.T) {
			b := pausedDomainBroker(t)
			bad := deliveryFault("E", "retrying")
			emitDelivery(t, b, bad)
			changed := false
			b.domain.options.Clock = func() time.Time {
				if !changed {
					changed = true
					b.domain.mu.Lock()
					entry := *b.domain.entries[domainCoalescerKey(bad)]
					if mode == "replace" {
						b.domain.entries[domainCoalescerKey(bad)] = &entry
					} else {
						other := deliveryFault("E", "blocked")
						entry.data = projectionJSON(t, other)
						b.domain.entries[domainCoalescerKey(other)] = &entry
						b.domain.stateBytes += entry.bytes
					}
					b.domain.mu.Unlock()
				}
				return time.Now()
			}
			if emitDelivery(t, b, deliveryCommit("E")) || !changed || b.queue.count != 1 || b.domain.Stats().Dropped != 1 {
				t.Fatal("stale complete-set token cleared replaced/new fault or emitted confirmation")
			}
		})
	}
}
func TestDomainDeliveryFullBrokerLossDoesNotBecomePartialCoalescerState(t *testing.T) {
	b := pausedDomainBroker(t)
	for _, code := range []string{"retrying", "blocked"} {
		for range 3 {
			emitDelivery(t, b, deliveryFault("E", code))
		}
	}
	for b.Emit(Input{Code: "startup.ready"}) {
	}
	lost := b.dropped
	if !emitDelivery(t, b, deliveryCommit("E")) {
		t.Fatal("broker refusal changed verified transition observation")
	}
	s, _ := b.Snapshot(Query{})
	if b.domain.Stats().Entries != 1 || b.domain.Stats().Dropped != 4 || s.CaptureDropped != lost+4 {
		t.Fatal("two summaries/recovery/first output loss not independent or fault set partially retained")
	}
}
func TestDomainDeliveryClosedAndCancelledDoNotInventCommit(t *testing.T) {
	for _, mode := range []string{"cancel", "close"} {
		t.Run(mode, func(t *testing.T) {
			b := pausedDomainBroker(t)
			emitDelivery(t, b, deliveryFault("E", "retrying"))
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			if mode == "cancel" {
				cancel()
			} else {
				b.domain.Close()
			}
			if b.EmitDomain(ctx, projectionJSON(t, deliveryCommit("E")), domainScope()) || b.queue.count != 1 || b.domain.Stats().Dropped != 1 {
				t.Fatal("closed/canceled capture admitted SQL success or lost rejection accounting")
			}
		})
	}
}
func TestDomainDeliveryEmptyRestartHasNoReceiptReplayInput(t *testing.T) {
	b := pausedDomainBroker(t)
	b.domain.Tick(t.Context())
	if b.queue.count != 0 || b.domain.Stats().Entries != 0 {
		t.Fatal("empty capture manufactured startup receipt event")
	}
	// Actual persisted Already/restart exclusions are separately tested at Store's
	// transition authority, not inferred from this coalescer's lack of history.
}

type domainCaptureLossRecorder interface{ RecordDomainCaptureLoss(uint64) }

func TestDomainDeliveryExplicitCaptureOverflowUsesSingleLossOwner(t *testing.T) {
	b := pausedDomainBroker(t)
	recorder, ok := any(b).(domainCaptureLossRecorder)
	if !ok {
		t.Fatal("broker has no nonblocking capture-loss scalar boundary")
	}
	var nilBroker *Broker
	nilRecorder, ok := any(nilBroker).(domainCaptureLossRecorder)
	if !ok {
		t.Fatal("nil broker does not implement capture-loss boundary")
	}
	nilRecorder.RecordDomainCaptureLoss(1)
	recorder.RecordDomainCaptureLoss(13)
	recorder.RecordDomainCaptureLoss(4)
	recorder.RecordDomainCaptureLoss(0)
	s, _ := b.Snapshot(Query{})
	if s.CaptureDropped != 17 || b.queue.count != 0 || b.domain.Stats().Entries != 0 || b.domain.Stats().Dropped != 0 {
		t.Fatal("bounded observation overflow duplicated losses or fabricated domain/unavailable events")
	}
	recorder.RecordDomainCaptureLoss(^uint64(0))
	s, _ = b.Snapshot(Query{})
	if s.CaptureDropped != ^uint64(0) || b.queue.count != 0 {
		t.Fatal("capture-loss scalar wrapped or created output")
	}
}

func TestDomainDeliveryRevisionCountsAndLateSummaryRemainSeparate(t *testing.T) {
	b := pausedDomainBroker(t)
	bad := deliveryFault("E", "retrying")
	for range 3 {
		emitDelivery(t, b, bad)
	}
	old := b.domain.entries[domainCoalescerKey(bad)]
	pending := domainEntryEvent(*old, true)
	r2 := bad
	r2.Binding.EntityRevision = "R2"
	r2.Locator.GroupRevision = domainText("R2")
	emitDelivery(t, b, r2)
	emitDelivery(t, b, deliveryCommit("E"))
	if b.domain.entries[domainCoalescerKey(r2)] == nil || b.domain.entries[domainCoalescerKey(bad)] != nil {
		t.Fatal("R1 receipt changed R2 fault state")
	}
	summaries := deliveryEvents(b, "domain.delivery.retrying")
	if len(summaries) != 3 || summaries[2].Fields["entity_revision"] != "R1" || summaries[2].Fields["occurrences"] != "3" {
		t.Fatal("R1/R2 observations or final summary merged")
	}
	if b.domain.enqueueCurrent(t.Context(), pending, domainCoalescerKey(bad), old) {
		t.Fatal("late prepared Tick reopened recovered effect")
	}
}

func TestDomainDeliveryReentryCannotSplitSummaryRecoveryAndFirst(t *testing.T) {
	b := pausedDomainBroker(t)
	bad := deliveryFault("E", "retrying")
	for range 3 {
		emitDelivery(t, b, bad)
	}
	original := b.domain.options.TryAdmit
	reentered := false
	b.domain.options.TryAdmit = func(e Event) bool {
		if e.Code == "domain.delivery.retrying" && e.Fields["occurrences"] == "3" {
			if !b.domain.mu.TryLock() {
				t.Error("callback invoked under coalescer mutex")
				return false
			}
			b.domain.mu.Unlock()
			reentered = true
			emitDelivery(t, b, deliveryFault("other-effect", "blocked"))
		}
		return original(e)
	}
	emitDelivery(t, b, deliveryCommit("E"))
	if !reentered || b.queue.count != 5 || b.queue.at(1).event.Code != "domain.delivery.retrying" || b.queue.at(2).event.Code != "domain.delivery.recovered" || b.queue.at(3).event.Code != "domain.delivery.first-confirmed" || b.queue.at(4).event.Fields["reference_id"] != "other-effect" {
		t.Fatal("callback reentry interleaved another effect into reserved transition")
	}
}

func TestDomainDeliveryReceiptDoesNotClearReplayInputsOrOwnerRefusal(t *testing.T) {
	b := pausedDomainBroker(t)
	replay := projectionRef(projectionCases[16])
	replay.Evidence.ReferenceID = "E"
	refusal := recordingRefusal()
	refusal.Evidence.ReferenceID = "E"
	emitDelivery(t, b, replay)
	emitDelivery(t, b, refusal)
	emitDelivery(t, b, deliveryCommit("E"))
	if b.domain.Stats().Entries != 3 || len(deliveryEvents(b, "domain.delivery.recovered")) != 0 || b.domain.entries[domainCoalescerKey(replay)] == nil || b.domain.entries[domainCoalescerKey(refusal)] == nil {
		t.Fatal("SQL receipt cleared independent accepted-input quarantine or ownership refusal")
	}
	if b.queue.at(0).event.Fields["evidence_kind"] != "counter-replay-inputs" || b.queue.at(1).event.Fields["evidence_kind"] != "recording-refusal" || b.queue.at(1).event.Fields["pending_samples"] != "0" {
		t.Fatal("input observations became finalized output evidence or zero became unavailable")
	}
}

type domainWarningEligibility interface{ DomainWarningsEnabled() bool }

func TestDomainDeliveryWarningEligibilityIsImmutableAndNilSafe(t *testing.T) {
	var nilBroker *Broker
	for _, tc := range []struct {
		name   string
		broker *Broker
		want   bool
	}{
		{"nil", nilBroker, false}, {"uninitialized", &Broker{}, false},
		{"error", pausedDomainBroker(t), false}, {"warn", pausedDomainBroker(t), true}, {"info", pausedDomainBroker(t), true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name == "error" {
				tc.broker.level = levelRank(levelError)
			}
			if tc.name == "warn" {
				tc.broker.level = levelRank(levelWarn)
			}
			query, ok := any(tc.broker).(domainWarningEligibility)
			if !ok {
				t.Fatal("broker lacks immutable warning eligibility query for preprojection producer gating")
			}
			if tc.broker != nil {
				tc.broker.mu.Lock()
				defer tc.broker.mu.Unlock()
			}
			done := make(chan bool, 1)
			go func() { done <- query.DomainWarningsEnabled() }()
			select {
			case got := <-done:
				if got != tc.want {
					t.Fatalf("warning eligibility=%v want %v", got, tc.want)
				}
			case <-time.After(time.Second):
				t.Fatal("immutable eligibility query blocked on broker lock")
			}
			if tc.broker != nil && (tc.broker.domainLoss.Load() != 0 || tc.broker.queue.count != 0) {
				t.Fatal("eligibility query mutated capture state or loss")
			}
		})
	}
}
