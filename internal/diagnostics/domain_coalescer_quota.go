package diagnostics

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"time"

	"go-gateway/internal/datalink/runtimeview"
)

// QuotaRecoveryHint carries a detached validated original reference, not permission
// or current health. Its private entry/generation token cannot be forged by native
// packages. Mutation of Ref invalidates the completion rather than changing history.
type QuotaRecoveryHint struct {
	Ref        DomainEventRef
	entry      *domainCoalescerEntry
	generation uint64
}

// TakeQuotaRecoveryHint selects at most one retained quota episode per existing
// heartbeat. No fault, filtering, contention, cancellation or shutdown does no
// work and records no capture loss. The caller alone owns the subsequent pure read.
func (b *Broker) TakeQuotaRecoveryHint(ctx context.Context) (QuotaRecoveryHint, bool) {
	if !b.DomainWarningsEnabled() || ctx == nil || ctx.Err() != nil {
		return QuotaRecoveryHint{}, false
	}
	c := b.domain
	now, ok := c.now()
	if !ok || !c.mu.TryLock() {
		return QuotaRecoveryHint{}, false
	}
	if c.closed || now.Before(c.quotaHintNext) {
		c.mu.Unlock()
		return QuotaRecoveryHint{}, false
	}
	var chosen *domainCoalescerEntry
	var chosenKey string
	for key, e := range c.entries {
		if !strings.Contains(key, "\x00blocked\x00domain.quota.refused\x00") {
			continue
		}
		if chosen == nil || e.quotaProbeAt.Before(chosen.quotaProbeAt) || (e.quotaProbeAt.Equal(chosen.quotaProbeAt) && key < chosenKey) {
			chosen, chosenKey = e, key
		}
	}
	if chosen == nil {
		c.mu.Unlock()
		return QuotaRecoveryHint{}, false
	}
	raw, scope, generation := chosen.data, chosen.scope, chosen.occurrences
	chosen.quotaProbeAt = now
	c.quotaHintNext = now.Add(c.options.Heartbeat)
	c.mu.Unlock()
	ref, err := ParseDomainEventRef(raw, scope)
	if err != nil || ref.Code != "domain.quota.refused" {
		return QuotaRecoveryHint{}, false
	}
	return QuotaRecoveryHint{Ref: ref, entry: chosen, generation: generation}, true
}

// EmitQuotaRecovery retires only the original retained episode after the adapter
// supplies exact configured native capacity evidence. This does not resume intake
// or imply SQL commitment. A filtered info recovery still finalizes warn faults.
func (b *Broker) EmitQuotaRecovery(ctx context.Context, hint QuotaRecoveryHint, data []byte, scope runtimeview.Scope) bool {
	if !b.DomainWarningsEnabled() {
		return false
	}
	return b.domain.recoverQuota(ctx, hint, data, scope, b.DomainInfoEnabled())
}

func (c *DomainCoalescer) recoverQuota(ctx context.Context, hint QuotaRecoveryHint, data []byte, scope runtimeview.Scope, visible bool) bool {
	if ctx == nil || ctx.Err() != nil || hint.entry == nil {
		return c.rejectQuotaCapture()
	}
	recovered, err := ParseDomainEventRef(data, scope)
	if err != nil || recovered.Code != "domain.quota.recovered" {
		return c.rejectQuotaCapture()
	}
	// Ref is public and mutable. Canonical equality with immutable retained bytes
	// prevents callers from borrowing a private token for a different identity.
	raw, err := json.Marshal(hint.Ref)
	if err != nil || !bytes.Equal(raw, hint.entry.data) {
		return c.rejectQuotaCapture()
	}
	original, err := ParseDomainEventRef(hint.entry.data, scope)
	if err != nil || original.Code != "domain.quota.refused" || domainEntityKey(original) != domainEntityKey(recovered) || original.Binding.WorkspaceRevision != recovered.Binding.WorkspaceRevision || original.Evidence.ReferenceID != recovered.Evidence.ReferenceID || original.Evidence.Quota.Scope != recovered.Evidence.Quota.Scope {
		return c.rejectQuotaCapture()
	}
	now, ok := c.now()
	if !ok {
		return c.rejectQuotaCapture()
	}
	key := domainCoalescerKey(original)
	if !c.mu.TryLock() {
		return c.rejectQuotaCapture()
	}
	entry := c.entries[key]
	if c.closed || entry != hint.entry || entry.occurrences != hint.generation || entry.reference != original.Evidence.ReferenceID || entry.scope != scope || now.Before(entry.last) {
		c.mu.Unlock()
		return c.rejectQuotaCapture()
	}
	var batch [2]record
	count, size := 0, 0
	if entry.occurrences > entry.reported {
		e := domainEntryEvent(*entry, true)
		n := domainOutputBytes(e)
		batch[count] = record{event: e, size: n}
		count++
		size += n
	}
	if visible {
		e := ProjectDomainSummary(data, scope, DomainEventSummary{FirstSeenAt: now.Format(time.RFC3339Nano), LastSeenAt: now.Format(time.RFC3339Nano), Occurrences: "1"})
		n := domainOutputBytes(e)
		batch[count] = record{event: e, size: n}
		count++
		size += n
	}
	// Reservation/state retirement is atomic here, not across broker/file/SSE.
	// A refusal retains all old state; sinks account any later output loss separately.
	if !c.fitsLocked(count-1, size-entry.bytes) {
		c.mu.Unlock()
		return c.rejectQuotaCapture()
	}
	delete(c.entries, key)
	c.stateBytes -= entry.bytes
	for _, r := range batch[:count] {
		c.queue.push(r)
	}
	c.mu.Unlock()
	c.drain(ctx)
	return visible
}
func (c *DomainCoalescer) rejectQuotaCapture() bool { c.drop(1); return false }
