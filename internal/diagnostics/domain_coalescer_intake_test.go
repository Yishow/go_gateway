package diagnostics

import (
	"testing"
	"time"
)

func TestDomainIntakeSameStateHasNoHeartbeatSpam(t *testing.T) {
	b := pausedDomainBroker(t)
	clock := &coalescerClock{at: time.Now()}
	b.domain.options.Clock = clock.Now
	for range 3 {
		if !emitDelivery(t, b, intakeRef("applied")) {
			t.Fatal("verified intake refused")
		}
	}
	clock.Advance(RuntimeLogHeartbeat)
	b.domain.Tick(t.Context())
	if b.queue.count != 1 || b.domain.Stats().Entries != 1 || b.domain.Stats().Dropped != 0 {
		t.Fatalf("healthy repeated intake produced summary or loss: queue=%d stats=%#v", b.queue.count, b.domain.Stats())
	}
}
func TestDomainIntakeMeaningfulTransitionsKeepOneBoundedCurrentState(t *testing.T) {
	b := pausedDomainBroker(t)
	for _, state := range []string{"applied", "paused", "paused", "applied", "retiring"} {
		emitDelivery(t, b, intakeRef(state))
	}
	if b.queue.count != 4 || b.domain.Stats().Entries != 1 {
		t.Fatalf("state changes lost or obsolete states retained: %d %#v", b.queue.count, b.domain.Stats())
	}
	want := []string{"domain.group.applied", "domain.group.paused", "domain.group.applied", "domain.group.retiring"}
	for i, code := range want {
		if b.queue.at(i).event.Code != code {
			t.Fatalf("transition %d=%s", i, b.queue.at(i).event.Code)
		}
	}
}
func TestDomainIntakeRevisionStatesStayIndependent(t *testing.T) {
	b := pausedDomainBroker(t)
	a := intakeRef("applied")
	emitDelivery(t, b, a)
	other := intakeRef("retiring")
	other.Binding.EntityRevision = "R2"
	other.Locator.GroupRevision = domainText("R2")
	other.Evidence.Intake.GroupRevision = "R2"
	emitDelivery(t, b, other)
	if b.queue.count != 2 || b.domain.Stats().Entries != 2 {
		t.Fatal("different original revisions shared intake state")
	}
}
func TestDomainIntakeFilteredStateDoesNotAllocateOrLose(t *testing.T) {
	b := pausedDomainBroker(t)
	b.level = levelRank(levelWarn)
	for range 5 {
		if emitDelivery(t, b, intakeRef("applied")) {
			t.Fatal("filtered intake accepted")
		}
	}
	if b.queue.count != 0 || b.domain.Stats().Entries != 0 || b.domain.Stats().Dropped != 0 {
		t.Fatal("intentional filter allocated or lost")
	}
}

func TestDomainIntakeAtomicCapacityAndContentionRetainOldState(t *testing.T) {
	for _, kind := range []string{"count", "bytes", "contended", "closed"} {
		t.Run(kind, func(t *testing.T) {
			b := pausedDomainBroker(t)
			emitDelivery(t, b, intakeRef("applied"))
			switch kind {
			case "count":
				b.domain.options.RecordLimit = 1
			case "bytes":
				b.domain.options.ByteLimit = 1
			case "contended":
				b.domain.mu.Lock()
			case "closed":
				b.domain.Close()
			}
			got := emitDelivery(t, b, intakeRef("paused"))
			if kind == "contended" {
				b.domain.mu.Unlock()
			}
			if got || b.queue.count != 1 || b.domain.Stats().Dropped != 1 {
				t.Fatalf("partial intake transition/capture accounting: got=%v queue=%d stats=%#v", got, b.queue.count, b.domain.Stats())
			}
			if kind != "closed" && b.domain.Stats().Entries != 1 {
				t.Fatal("old intake state lost")
			}
		})
	}
}
func TestDomainIntakeDifferentWorkspaceAndConnectorStayIndependent(t *testing.T) {
	for _, kind := range []string{"workspace-revision", "connector-revision"} {
		t.Run(kind, func(t *testing.T) {
			b := pausedDomainBroker(t)
			emitDelivery(t, b, intakeRef("applied"))
			v := intakeRef("paused")
			scope := domainScope()
			if kind == "workspace-revision" {
				v.Binding.WorkspaceRevision = "W2"
				scope.WorkspaceRevision = "W2"
			} else {
				v.ConnectorRevision = domainText("C2")
			}
			b.EmitDomain(t.Context(), projectionJSON(t, v), scope)
			if b.domain.Stats().Entries != 2 || b.queue.count != 2 {
				t.Fatal("scope proof merged")
			}
		})
	}
}
