package diagnostics

import (
	"testing"
	"time"
)

func operationRef(kind, outcome, id string) DomainEventRef {
	v := domainFixture()
	v.Code = DomainEventCode("domain." + kind + "." + outcome)
	v.Stage = "blocked"
	if outcome == "completed" || outcome == "applied" {
		v.Stage, v.Reason, v.Action = "accepted", "none", "none"
	}
	v.Evidence = DomainEvidenceRef{Kind: "operation-result", ReferenceID: id}
	return v
}
func TestDomainOperationReferencesNeverCoalesceAcrossOperations(t *testing.T) {
	for _, pair := range [][2]string{{"repair", "blocked"}, {"repair", "completed"}, {"cutover", "rejected"}, {"cutover", "applied"}} {
		t.Run(pair[0]+pair[1], func(t *testing.T) {
			b := pausedDomainBroker(t)
			clock := &coalescerClock{at: time.Now()}
			b.domain.options.Clock = clock.Now
			for range 3 {
				emitDelivery(t, b, operationRef(pair[0], pair[1], "operation-1"))
			}
			emitDelivery(t, b, operationRef(pair[0], pair[1], "operation-2"))
			if b.queue.count != 2 || b.domain.Stats().Entries != 2 {
				t.Fatalf("real ledger references merged: queue=%d entries=%d", b.queue.count, b.domain.Stats().Entries)
			}
			clock.Advance(RuntimeLogHeartbeat)
			b.domain.Tick(t.Context())
			if b.queue.count != 3 || b.queue.at(2).event.Fields["reference_id"] != "operation-1" || b.queue.at(2).event.Fields["occurrences"] != "3" {
				t.Fatal("summary borrowed another operation proof/count")
			}
		})
	}
}
func TestDomainOperationNewCompletionDoesNotRetireAnotherBlockedOperation(t *testing.T) {
	b := pausedDomainBroker(t)
	for range 3 {
		emitDelivery(t, b, operationRef("repair", "blocked", "old"))
	}
	emitDelivery(t, b, operationRef("repair", "completed", "new"))
	if b.domain.Stats().Entries != 2 || b.queue.count != 2 || b.queue.at(1).event.Fields["reference_id"] != "new" {
		t.Fatal("another operation's completion inferred fault recovery")
	}
}
func TestDomainOperationReferenceKeyChangeDoesNotAffectOtherProducers(t *testing.T) {
	for _, kind := range []string{"source", "collector", "mapping", "quota", "intake"} {
		t.Run(kind, func(t *testing.T) {
			v := projectionRef(projectionCases[1])
			switch kind {
			case "collector":
				v = projectionRef(projectionCases[2])
			case "mapping":
				v = projectionRef(projectionCases[4])
			case "quota":
				v = quotaRef("group", false)
			case "intake":
				v = intakeRef("applied")
			}
			other := v
			other.Evidence.ReferenceID = "another-observation"
			if domainCoalescerKey(v) != domainCoalescerKey(other) {
				t.Fatal("operation isolation changed unrelated producer identity")
			}
		})
	}
}

func TestDomainOperationExplicitRecoveryCannotBorrowAnotherLedgerID(t *testing.T) {
	for _, kind := range []string{"repair", "cutover"} {
		t.Run(kind, func(t *testing.T) {
			failure, success := "blocked", "completed"
			if kind == "cutover" {
				failure, success = "rejected", "applied"
			}
			b := pausedDomainBroker(t)
			old := operationRef(kind, failure, "original-op")
			emitDelivery(t, b, old)
			other := operationRef(kind, success, "other-op")
			if b.domain.Recover(t.Context(), projectionJSON(t, other), projectionJSON(t, old), domainScope()) || b.domain.Stats().Entries != 1 || b.queue.count != 1 {
				t.Fatal("another ledger operation borrowed original failure authority")
			}
			exact := operationRef(kind, success, "original-op")
			if !b.domain.Recover(t.Context(), projectionJSON(t, exact), projectionJSON(t, old), domainScope()) || b.domain.Stats().Entries != 0 || b.queue.count != 2 {
				t.Fatal("explicit exact-operation recovery rejected")
			}
		})
	}
}
