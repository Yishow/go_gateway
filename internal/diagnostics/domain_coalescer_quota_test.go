package diagnostics

import (
	"context"
	"testing"
	"time"
)

func TestDomainQuotaGlobalAndGroupFaultsNeverMerge(t *testing.T) {
	b := pausedDomainBroker(t)
	for range 3 {
		emitDelivery(t, b, quotaRef("group", false))
	}
	emitDelivery(t, b, quotaRef("global", false))
	if b.queue.count != 2 || b.domain.Stats().Entries != 2 {
		t.Fatalf("quota scopes merged: queue=%d stats=%#v", b.queue.count, b.domain.Stats())
	}
}
func TestDomainQuotaHintIsBoundedAndThrottledWithoutHealthyState(t *testing.T) {
	b := pausedDomainBroker(t)
	clock := &coalescerClock{at: time.Now()}
	b.domain.options.Clock = clock.Now
	if _, ok := takeQuotaHint(t, b, t.Context()); ok {
		t.Fatal("healthy no-fault hint")
	}
	emitDelivery(t, b, quotaRef("group", false))
	emitDelivery(t, b, quotaRef("global", false))
	h, ok := takeQuotaHint(t, b, t.Context())
	if !ok {
		t.Fatal("retained fault hint missing")
	}
	first := quotaHintRef(t, h)
	if first.Code != "domain.quota.refused" || first.Evidence.Quota == nil {
		t.Fatal("hint not validated original fault")
	}
	if _, ok := takeQuotaHint(t, b, t.Context()); ok {
		t.Fatal("more than one hint per heartbeat")
	}
	clock.Advance(RuntimeLogHeartbeat)
	h2, ok := takeQuotaHint(t, b, t.Context())
	if !ok || quotaHintRef(t, h2).Evidence.Quota.Scope == first.Evidence.Quota.Scope {
		t.Fatal("next retained scope starved")
	}
	if b.domain.Stats().Dropped != 0 || b.queue.count != 2 {
		t.Fatal("hint emitted or lost events")
	}
}
func TestDomainQuotaHintRecoveryFinalizesOnlyOriginalScope(t *testing.T) {
	b := pausedDomainBroker(t)
	clock := &coalescerClock{at: time.Now()}
	b.domain.options.Clock = clock.Now
	for range 3 {
		emitDelivery(t, b, quotaRef("group", false))
	}
	h, ok := takeQuotaHint(t, b, t.Context())
	if !ok {
		t.Fatal("hint missing")
	}
	emitDelivery(t, b, quotaRef("global", false))
	original := quotaHintRef(t, h)
	v := quotaRef(original.Evidence.Quota.Scope, true)
	if !emitQuotaHint(t, b, h, v, domainScope()) {
		t.Fatal("exact measured recovery refused")
	}
	if b.domain.Stats().Entries != 1 || b.queue.count != 4 || b.queue.at(2).event.Code != "domain.quota.refused" || b.queue.at(2).event.Fields["occurrences"] != "3" || b.queue.at(3).event.Code != "domain.quota.recovered" {
		t.Fatal("wrong scope retirement or final summary ordering")
	}
	clock.Advance(RuntimeLogHeartbeat)
	b.domain.Tick(t.Context())
	for _, e := range deliveryEvents(b, "domain.quota.refused") {
		if e.Fields["quota_scope"] == original.Evidence.Quota.Scope && e.Fields["occurrences"] == "3" && e.Fields["last_seen_at"] == nil {
			t.Fatal("bad final summary")
		}
	}
	if emitQuotaHint(t, b, h, v, domainScope()) || len(deliveryEvents(b, "domain.quota.recovered")) != 1 {
		t.Fatal("replayed hint made healthy spam")
	}
}
func TestDomainQuotaHintNoFaultContentionCancellationOrCloseHasNoLoss(t *testing.T) {
	for _, kind := range []string{"empty", "contended", "cancelled", "closed", "error"} {
		t.Run(kind, func(t *testing.T) {
			b := pausedDomainBroker(t)
			ctx := t.Context()
			if kind != "empty" {
				emitDelivery(t, b, quotaRef("group", false))
			}
			switch kind {
			case "contended":
				b.domain.mu.Lock()
				defer b.domain.mu.Unlock()
			case "cancelled":
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			case "closed":
				b.domain.Close()
			case "error":
				b.level = levelRank(levelError)
			}
			before := b.domain.dropped.Load()
			if _, ok := takeQuotaHint(t, b, ctx); ok {
				t.Fatal("noneligible hint returned")
			}
			if b.domain.dropped.Load() != before {
				t.Fatal("no-work hint counted event loss")
			}
		})
	}
}

func TestDomainQuotaHintRejectsNewerFaultAndABA(t *testing.T) {
	for _, kind := range []string{"repeat", "replaced", "ABA"} {
		t.Run(kind, func(t *testing.T) {
			b := pausedDomainBroker(t)
			emitDelivery(t, b, quotaRef("group", false))
			h, ok := takeQuotaHint(t, b, t.Context())
			if !ok {
				t.Fatal("hint missing")
			}
			if kind == "repeat" {
				emitDelivery(t, b, quotaRef("group", false))
			} else {
				original := quotaHintRef(t, h)
				key := domainCoalescerKey(original)
				b.domain.mu.Lock()
				old := b.domain.entries[key]
				delete(b.domain.entries, key)
				b.domain.stateBytes -= old.bytes
				b.domain.mu.Unlock()
				next := quotaRef("group", false)
				if kind == "replaced" {
					next.Evidence.ReferenceID = "next-observation"
				}
				emitDelivery(t, b, next)
			}
			if emitQuotaHint(t, b, h, quotaRef("group", true), domainScope()) || b.domain.Stats().Entries != 1 || len(deliveryEvents(b, "domain.quota.recovered")) != 0 {
				t.Fatal("stale probe retired newer episode")
			}
		})
	}
}
func TestDomainQuotaHintRejectsEveryForeignProofLeaf(t *testing.T) {
	for _, kind := range []string{"workspace", "workspace-revision", "group", "revision", "connector", "connector-revision", "scope", "reference", "mutated-hint"} {
		t.Run(kind, func(t *testing.T) {
			b := pausedDomainBroker(t)
			emitDelivery(t, b, quotaRef("group", false))
			h, ok := takeQuotaHint(t, b, t.Context())
			if !ok {
				t.Fatal("hint missing")
			}
			v := quotaRef("group", true)
			scope := domainScope()
			switch kind {
			case "workspace":
				v.Binding.WorkspaceID = "W2"
				v.Locator.WorkspaceID = "W2"
				scope.WorkspaceID = "W2"
			case "workspace-revision":
				v.Binding.WorkspaceRevision = "W2"
				scope.WorkspaceRevision = "W2"
			case "group":
				v.Locator.GroupID = domainText("G2")
			case "revision":
				v.Binding.EntityRevision = "R2"
				v.Locator.GroupRevision = domainText("R2")
			case "connector":
				v.ConnectorID = domainText("C2")
			case "connector-revision":
				v.ConnectorRevision = domainText("C2")
			case "scope":
				v.Evidence.Quota.Scope = "global"
			case "reference":
				v.Evidence.ReferenceID = "other"
			case "mutated-hint":
				r := quotaHintRef(t, h)
				r.Evidence.Quota.Scope = "global"
			}
			if emitQuotaHint(t, b, h, v, scope) || b.domain.Stats().Entries != 1 || len(deliveryEvents(b, "domain.quota.recovered")) != 0 {
				t.Fatal("foreign or altered hint retired original fault")
			}
		})
	}
}
func TestDomainQuotaFilteredRecoveryRetiresWarnWithoutSuccessState(t *testing.T) {
	b := pausedDomainBroker(t)
	b.level = levelRank(levelWarn)
	for range 3 {
		emitDelivery(t, b, quotaRef("group", false))
	}
	h, ok := takeQuotaHint(t, b, t.Context())
	if !ok {
		t.Fatal("warn level lost recovery processing")
	}
	if emitQuotaHint(t, b, h, quotaRef("group", true), domainScope()) {
		t.Fatal("filtered info reported visible")
	}
	if b.domain.Stats().Entries != 0 || b.domain.Stats().Dropped != 0 || b.queue.count != 2 || b.queue.at(1).event.Fields["occurrences"] != "3" || len(deliveryEvents(b, "domain.quota.recovered")) != 0 {
		t.Fatalf("filtered recovery did not finalize exact warning: %#v queue=%d", b.domain.Stats(), b.queue.count)
	}
}
func TestDomainQuotaAtomicReservationFailureRetainsFault(t *testing.T) {
	for _, limit := range []string{"count", "bytes"} {
		t.Run(limit, func(t *testing.T) {
			b := pausedDomainBroker(t)
			for range 3 {
				emitDelivery(t, b, quotaRef("group", false))
			}
			h, ok := takeQuotaHint(t, b, t.Context())
			if !ok {
				t.Fatal("hint missing")
			}
			if limit == "count" {
				b.domain.options.RecordLimit = 1
			} else {
				b.domain.options.ByteLimit = 1
			}
			if emitQuotaHint(t, b, h, quotaRef("group", true), domainScope()) || b.domain.Stats().Entries != 1 || b.domain.Stats().Dropped != 1 || b.queue.count != 1 {
				t.Fatalf("half transition retained or loss duplicated: %#v queue=%d", b.domain.Stats(), b.queue.count)
			}
		})
	}
}

func TestDomainQuotaNewRefusalDuringProbeCannotBeCleared(t *testing.T) {
	b := pausedDomainBroker(t)
	emitDelivery(t, b, quotaRef("group", false))
	h, ok := takeQuotaHint(t, b, t.Context())
	if !ok {
		t.Fatal("hint missing")
	}
	nested := false
	b.domain.options.Clock = func() time.Time {
		if !nested {
			nested = true
			emitDelivery(t, b, quotaRef("group", false))
		}
		return time.Now()
	}
	if emitQuotaHint(t, b, h, quotaRef("group", true), domainScope()) || b.domain.Stats().Entries != 1 || len(deliveryEvents(b, "domain.quota.recovered")) != 0 {
		t.Fatal("concurrent newer refusal lost to old probe")
	}
}
func TestDomainQuotaFullBrokerLossIsSeparateFromAtomicRetirement(t *testing.T) {
	b := pausedDomainBroker(t)
	for range 3 {
		emitDelivery(t, b, quotaRef("group", false))
	}
	h, ok := takeQuotaHint(t, b, t.Context())
	if !ok {
		t.Fatal("hint missing")
	}
	b.mu.Lock()
	got := emitQuotaHint(t, b, h, quotaRef("group", true), domainScope())
	b.mu.Unlock()
	if !got || b.domain.Stats().Entries != 0 || b.domain.Stats().Dropped != 2 || b.queue.count != 1 {
		t.Fatalf("sink loss rewrote coalescer transition or was double counted: got=%v stats=%#v queue=%d", got, b.domain.Stats(), b.queue.count)
	}
}
func TestDomainQuotaHintConcurrentRequestsShareOneHeartbeat(t *testing.T) {
	b := pausedDomainBroker(t)
	emitDelivery(t, b, quotaRef("group", false))
	result := make(chan bool, 32)
	for range 32 {
		go func() { _, ok := b.TakeQuotaRecoveryHint(t.Context()); result <- ok }()
	}
	accepted := 0
	for range 32 {
		if <-result {
			accepted++
		}
	}
	if accepted != 1 || b.domain.Stats().Dropped != 0 {
		t.Fatalf("concurrent hint fanout: accepted=%d stats=%#v", accepted, b.domain.Stats())
	}
}
func TestDomainQuotaFreshHintCanRecoverAfterRejectedOldGeneration(t *testing.T) {
	b := pausedDomainBroker(t)
	clock := &coalescerClock{at: time.Now()}
	b.domain.options.Clock = clock.Now
	emitDelivery(t, b, quotaRef("group", false))
	old, ok := takeQuotaHint(t, b, t.Context())
	if !ok {
		t.Fatal("initial hint missing")
	}
	emitDelivery(t, b, quotaRef("group", false))
	if emitQuotaHint(t, b, old, quotaRef("group", true), domainScope()) {
		t.Fatal("old generation accepted")
	}
	clock.Advance(RuntimeLogHeartbeat)
	fresh, ok := takeQuotaHint(t, b, t.Context())
	if !ok || !emitQuotaHint(t, b, fresh, quotaRef("group", true), domainScope()) {
		t.Fatal("fresh exact proof could not retire retained fault")
	}
	if b.domain.Stats().Entries != 0 || len(deliveryEvents(b, "domain.quota.recovered")) != 1 {
		t.Fatal("fresh recovery inconsistent")
	}
}
