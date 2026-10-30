package diagnostics

import (
	"testing"
	"time"
)

func TestDomainCoalescerDifferentReasonsKeepIndependentGuidance(t *testing.T) {
	c, _, sink := newCoalescerFixture(t)
	initial := domainFixture()
	ownership := domainFixture()
	ownership.Reason, ownership.Action = "recording-ownership-conflict", "inspect-immutable-evidence"
	c.Observe(t.Context(), projectionJSON(t, initial), domainScope())
	c.Observe(t.Context(), projectionJSON(t, ownership), domainScope())
	if c.Stats().Entries != 2 || len(sink.Events()) != 2 || sink.Events()[1].Fields["reason"] != "recording-ownership-conflict" || sink.Events()[1].Fields["action"] != "inspect-immutable-evidence" {
		t.Fatal("different owner reason was hidden under the first refresh guidance")
	}
	recovery := domainCommitFixture()
	recovery.Code = "domain.delivery.recovered"
	if !c.Recover(t.Context(), projectionJSON(t, recovery), projectionJSON(t, initial), domainScope()) || c.Stats().Entries != 1 {
		t.Fatal("recovery did not clear only its matching reason")
	}
	if c.Recover(t.Context(), projectionJSON(t, recovery), projectionJSON(t, initial), domainScope()) || c.Stats().Entries != 1 {
		t.Fatal("repeated recovery cleared a different reason")
	}
	if !c.Recover(t.Context(), projectionJSON(t, recovery), projectionJSON(t, ownership), domainScope()) || c.Stats().Entries != 0 {
		t.Fatal("remaining reason could not be recovered by explicit evidence")
	}
}

func TestDomainCoalescerWindowAndRevisionIsolation(t *testing.T) {
	c, clock, sink := newCoalescerFixture(t)
	r1 := projectionJSON(t, domainFixture())
	r2 := projectionJSON(t, coalescerRevision(domainFixture(), "R2"))
	if !c.Observe(t.Context(), r1, domainScope()) {
		t.Fatal("first fault refused")
	}
	for range 2 {
		clock.Advance(time.Second)
		if !c.Observe(t.Context(), r1, domainScope()) {
			t.Fatal("repeat refused")
		}
	}
	clock.Advance(time.Second)
	if !c.Observe(t.Context(), r2, domainScope()) {
		t.Fatal("R2 first fault refused")
	}
	if got := sink.Events(); len(got) != 2 || got[0].Fields["occurrences"] != "1" || got[1].Fields["group_revision"] != "R2" {
		t.Fatalf("first events not independent: %#v", got)
	}
	clock.Advance(12 * time.Second)
	c.Tick(t.Context())
	got := sink.Events()
	if len(got) != 3 || got[2].Fields["group_revision"] != "R1" || got[2].Fields["occurrences"] != "3" || got[2].Fields["first_seen_at"] != "2026-10-27T10:00:00Z" || got[2].Fields["last_seen_at"] != "2026-10-27T10:00:02Z" {
		t.Fatalf("R1 three-fault summary wrong: %#v", got)
	}
	c.Tick(t.Context())
	if len(sink.Events()) != 3 {
		t.Fatal("unchanged heartbeat repeated a summary")
	}
	recovery := domainCommitFixture()
	recovery.Code = "domain.delivery.recovered"
	if !c.Recover(t.Context(), projectionJSON(t, recovery), r1, domainScope()) {
		t.Fatal("matching durable recovery refused")
	}
	if c.Stats().Entries != 1 || sink.Events()[3].Code != "domain.delivery.recovered" || sink.Events()[3].Fields["effect_key"] != "E" {
		t.Fatal("recovery must clear R1 only and retain receipt E")
	}
	if c.Recover(t.Context(), projectionJSON(t, recovery), r1, domainScope()) || len(sink.Events()) != 4 {
		t.Fatal("duplicate recovery emitted a new success")
	}
	if !c.Observe(t.Context(), r1, domainScope()) || sink.Events()[4].Fields["occurrences"] != "1" {
		t.Fatal("new fault episode did not start at one")
	}
}

func TestDomainCoalescerRecoveryRequiresExactFault(t *testing.T) {
	for _, mutate := range []func(*DomainEventRef, *DomainEventRef){
		func(r, _ *DomainEventRef) { *r = coalescerRevision(*r, "R2") },
		func(r, _ *DomainEventRef) { r.Locator.GroupID = domainText("other-G") },
		func(_, f *DomainEventRef) { f.Code, f.Stage = "domain.delivery.blocked", "blocked" },
		func(r, _ *DomainEventRef) {
			r.Code, r.Stage, r.Evidence.Kind, r.Evidence.Receipt = "domain.repair.completed", "accepted", "operation-result", nil
		},
		func(r, _ *DomainEventRef) {
			r.Evidence.Receipt.EffectKey, r.Evidence.ReferenceID = "other-E", "other-E"
		},
	} {
		c, _, sink := newCoalescerFixture(t)
		f := domainFixture()
		r := domainCommitFixture()
		r.Code = "domain.delivery.recovered"
		c.Observe(t.Context(), projectionJSON(t, f), domainScope())
		mutate(&r, &f)
		if c.Recover(t.Context(), projectionJSON(t, r), projectionJSON(t, f), domainScope()) || c.Stats().Entries != 1 || len(sink.Events()) != 1 {
			t.Fatal("unproven recovery erased or relabeled an active fault")
		}
	}
	c, clock, sink := newCoalescerFixture(t)
	f := projectionJSON(t, domainFixture())
	c.Observe(t.Context(), f, domainScope())
	clock.Advance(time.Second)
	c.Observe(t.Context(), f, domainScope())
	r := domainCommitFixture()
	r.Code = "domain.delivery.recovered"
	c.Recover(t.Context(), projectionJSON(t, r), f, domainScope())
	got := sink.Events()
	if len(got) != 3 || got[1].Fields["occurrences"] != "2" || got[2].Code != "domain.delivery.recovered" {
		t.Fatal("recovery lost its pending fault summary or was merged into it")
	}
}

func TestDomainCoalescerNoStaleSummaryAfterRecovery(t *testing.T) {
	c, clock, sink := newCoalescerFixture(t)
	f := domainFixture()
	raw := projectionJSON(t, f)
	c.Observe(t.Context(), raw, domainScope())
	clock.Advance(time.Second)
	c.Observe(t.Context(), raw, domainScope())
	key := domainCoalescerKey(f)
	c.mu.Lock()
	token := c.entries[key]
	snapshot := *token
	token.reported = token.occurrences
	c.mu.Unlock()
	r := domainCommitFixture()
	r.Code = "domain.delivery.recovered"
	if !c.Recover(t.Context(), projectionJSON(t, r), raw, domainScope()) {
		t.Fatal("recovery failed")
	}
	// A Tick snapshot prepared before recovery must not reopen the cleared fault.
	if c.enqueueCurrent(t.Context(), domainEntryEvent(snapshot, true), key, token) || len(sink.Events()) != 2 {
		t.Fatal("stale heartbeat fault emitted after matching recovery")
	}
}

func TestDomainCoalescerRecoveryFamilies(t *testing.T) {
	for _, pair := range [][2]string{{"domain.source.rejected", "domain.source.applied"}, {"domain.collector.failed", "domain.collector.recovered"}, {"domain.mapping.failed", "domain.mapping.recovered"}, {"domain.repair.blocked", "domain.repair.completed"}, {"domain.cutover.rejected", "domain.cutover.applied"}} {
		c, _, sink := newCoalescerFixture(t)
		var failed, recovered DomainEventRef
		for _, tc := range projectionCases {
			if tc.code == pair[0] {
				failed = projectionRef(tc)
			}
			if tc.code == pair[1] {
				recovered = projectionRef(tc)
			}
		}
		f := projectionJSON(t, failed)
		c.Observe(t.Context(), f, domainScope())
		if !c.Recover(t.Context(), projectionJSON(t, recovered), f, domainScope()) || c.Stats().Entries != 0 || len(sink.Events()) != 2 {
			t.Fatalf("typed recovery pair failed: %v", pair)
		}
	}
	c, _, _ := newCoalescerFixture(t)
	var input DomainEventRef
	for _, tc := range projectionCases {
		if tc.code == "domain.counter-replay.quarantined" {
			input = projectionRef(tc)
		}
	}
	f := projectionJSON(t, input)
	c.Observe(t.Context(), f, domainScope())
	r := domainCommitFixture()
	r.Code = "domain.delivery.recovered"
	if c.Recover(t.Context(), projectionJSON(t, r), f, domainScope()) || c.Stats().Entries != 1 {
		t.Fatal("SQL receipt cleared quarantined input samples")
	}
}

func TestDomainCoalescerRecoverReentryPreservesEpisodeOrder(t *testing.T) {
	c, clock, sink := newCoalescerFixture(t)
	raw := projectionJSON(t, domainFixture())
	c.Observe(t.Context(), raw, domainScope())
	clock.Advance(time.Second)
	c.Observe(t.Context(), raw, domainScope())
	reentered := false
	c.options.TryAdmit = func(e Event) bool {
		accepted := sink.Admit(e)
		if !reentered && e.Fields["occurrences"] == "2" {
			reentered = true
			clock.Advance(time.Second)
			if !c.Observe(t.Context(), raw, domainScope()) {
				t.Error("new episode refused")
			}
		}
		return accepted
	}
	recovery := domainCommitFixture()
	recovery.Code = "domain.delivery.recovered"
	if !c.Recover(t.Context(), projectionJSON(t, recovery), raw, domainScope()) {
		t.Fatal("recovery refused")
	}
	got := sink.Events()
	if len(got) != 4 || got[1].Fields["occurrences"] != "2" || got[2].Code != "domain.delivery.recovered" || got[2].Fields["first_seen_at"] != "2026-10-27T10:00:01Z" || got[3].Code != "domain.delivery.retrying" || got[3].Fields["first_seen_at"] != "2026-10-27T10:00:02Z" || got[3].Fields["occurrences"] != "1" {
		t.Fatalf("old recovery must precede new episode fault: %#v", got)
	}
	if c.Stats().Entries != 1 || c.Stats().Dropped != 0 {
		t.Fatal("old recovery cleared or lost the new episode")
	}
}

func TestDomainCoalescerRecoveryBatchCapacityFailure(t *testing.T) {
	for _, capKind := range []string{"records", "bytes"} {
		t.Run(capKind, func(t *testing.T) {
			c, clock, sink := newCoalescerFixture(t)
			raw := projectionJSON(t, domainFixture())
			c.Observe(t.Context(), raw, domainScope())
			c.Observe(t.Context(), raw, domainScope())
			c.draining.Store(true)
			if !c.Observe(t.Context(), []byte(`{"raw_error":"password_fixture"}`), domainScope()) {
				t.Fatal("suppressed queue setup failed")
			}
			recovery := domainCommitFixture()
			recovery.Code = "domain.delivery.recovered"
			r := projectionJSON(t, recovery)
			if capKind == "records" {
				c.options.RecordLimit = 2
			} else {
				now := clock.Now().Format(time.RFC3339Nano)
				output := ProjectDomainSummary(r, domainScope(), DomainEventSummary{FirstSeenAt: now, LastSeenAt: now, Occurrences: "1"})
				c.options.ByteLimit = c.queue.bytes + domainOutputBytes(output)
			}
			c.options.OnDrop = func(n uint64) { _ = c.Stats(); sink.Drop(n) }
			if !c.Recover(t.Context(), r, raw, domainScope()) {
				t.Fatal("capture capacity rewrote a confirmed lifecycle observation")
			}
			stats := c.Stats()
			if stats.Entries != 0 || stats.Queued != 1 || stats.Dropped != 2 {
				t.Fatalf("partial recovery batch retained or loss miscounted: %#v", stats)
			}
			c.draining.Store(false)
			c.Tick(t.Context())
			got := sink.Events()
			if len(got) != 2 || got[1].Code != "raw.suppressed" || sink.drops != 2 {
				t.Fatal("capacity failure must lose both transition outputs, not orphan a final summary")
			}
		})
	}
}
