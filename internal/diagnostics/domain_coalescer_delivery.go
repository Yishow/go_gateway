package diagnostics

import (
	"context"
	"encoding/json"
	"math"
	"slices"
	"strings"
	"time"

	"go-gateway/internal/datalink/runtimeview"
)

func deliveryFaultCode(code DomainEventCode) bool {
	switch code {
	case domainDeliveryRetrying, domainDeliveryBlocked, domainDeliveryQuarantined, domainDeliveryUnknown:
		return true
	default:
		return false
	}
}

// The producer supplies only newly committed, exact-receipt-verified transitions.
// This is bounded capture bookkeeping, not a durable deduplication ledger. Already
// and restart observations must be excluded at their native transition authority.
func (c *DomainCoalescer) observeDeliveryCommit(ctx context.Context, data []byte, scope runtimeview.Scope, visible bool) bool {
	if ctx == nil || ctx.Err() != nil {
		if !visible {
			return false
		}
		return c.rejectDeliveryCapture()
	}
	v, err := ParseDomainEventRef(data, scope)
	if err != nil || v.Code != domainDeliveryFirstConfirmed {
		return c.rejectDeliveryCapture()
	}
	successKey := deliverySuccessKey(v)
	if !c.mu.TryLock() {
		if !visible {
			return false
		}
		return c.rejectDeliveryCapture()
	}
	expectedSuccess := c.entries[successKey]
	faults := c.deliveryFaultTokensLocked(v)
	closed := c.closed
	c.mu.Unlock()
	if !visible && (closed || len(faults) == 0) {
		return false
	}
	if closed {
		return c.rejectDeliveryCapture()
	}
	// Only immutable reference bytes are inspected outside the mutex. A final
	// complete-set/pointer check rejects replacement AND newly appearing faults.
	for _, token := range faults {
		f, err := ParseDomainEventRef(token.entry.data, scope)
		if err != nil || !deliveryFaultCode(f.Code) || domainCoalescerKey(f) != token.key || f.Evidence.ReferenceID != v.Evidence.ReferenceID || domainEntityKey(f) != domainEntityKey(v) || f.Binding.WorkspaceRevision != v.Binding.WorkspaceRevision {
			return c.rejectDeliveryCapture()
		}
	}
	slices.SortFunc(faults, func(a, b deliveryFaultToken) int { return strings.Compare(a.key, b.key) })
	now, ok := c.now()
	if !ok {
		return c.rejectDeliveryCapture()
	}
	raw, err := json.Marshal(v)
	if err != nil {
		return c.rejectDeliveryCapture()
	}
	next := c.deliverySuccessEntry(v, raw, scope, now)
	if !c.mu.TryLock() {
		return c.rejectDeliveryCapture()
	}
	if c.closed || c.entries[successKey] != expectedSuccess || !c.deliveryFaultTokensCurrentLocked(v, faults) {
		c.mu.Unlock()
		return c.rejectDeliveryCapture()
	}
	old := c.entries[successKey]
	if old != nil && (old.scope != scope || now.Before(old.last) || (visible && old.occurrences == math.MaxUint64)) {
		c.mu.Unlock()
		return c.rejectDeliveryCapture()
	}
	for _, token := range faults {
		entry := token.entry
		if entry.scope != scope || entry.reference != token.reference || now.Before(entry.last) {
			c.mu.Unlock()
			return c.rejectDeliveryCapture()
		}
	}
	// At most the already-capped fault count plus recovery and capture-first.
	// All reservation/state changes happen together; sinks may still lose outputs
	// independently, through the existing sole OnDrop callback after unlocking.
	batch := make([]record, 0, len(faults)+2)
	bytes, removedBytes := 0, 0
	appendOutput := func(e Event) {
		size := domainOutputBytes(e)
		batch = append(batch, record{event: e, size: size})
		bytes += size
	}
	for _, token := range faults {
		entry := token.entry
		removedBytes += entry.bytes
		if entry.occurrences > entry.reported {
			appendOutput(domainEntryEvent(*entry, true))
		}
	}
	newCount, newBytes := 0, 0
	if visible {
		if len(faults) > 0 {
			recovery := v
			recovery.Code = domainDeliveryRecovered
			recoveryRaw, err := json.Marshal(recovery)
			if err != nil {
				c.mu.Unlock()
				return c.rejectDeliveryCapture()
			}
			appendOutput(ProjectDomainSummary(recoveryRaw, scope, DomainEventSummary{FirstSeenAt: now.Format(time.RFC3339Nano), LastSeenAt: now.Format(time.RFC3339Nano), Occurrences: "1"}))
		}
		if old == nil {
			newCount, newBytes = 1, next.bytes
			appendOutput(domainEntryEvent(next, false))
		}
	}
	if !c.fitsLocked(newCount+len(batch)-len(faults), newBytes+bytes-removedBytes) {
		c.mu.Unlock()
		return c.rejectDeliveryCapture()
	}
	for _, token := range faults {
		delete(c.entries, token.key)
	}
	c.stateBytes -= removedBytes
	if visible {
		if old == nil {
			c.entries[successKey] = &next
			c.stateBytes += next.bytes
		} else {
			old.occurrences++
			old.last = now
		}
	}
	for _, r := range batch {
		c.queue.push(r)
	}
	c.mu.Unlock()
	c.drain(ctx)
	return visible
}

type deliveryFaultToken struct {
	key, reference string
	entry          *domainCoalescerEntry
}

func deliverySuccessKey(v DomainEventRef) string {
	return "\x00domain.delivery.success\x00" + domainEntityKey(v) + "\x00" + v.Binding.WorkspaceRevision
}

func deliveryFaultKeyMatches(key string, v DomainEventRef) bool {
	prefix := domainEntityKey(v) + "\x00" + v.Binding.WorkspaceRevision + "\x00"
	if !strings.HasPrefix(key, prefix) || !strings.HasSuffix(key, "\x00"+v.Evidence.ReferenceID) {
		return false
	}
	rest := strings.TrimPrefix(key, prefix)
	for _, stage := range []string{string(runtimeview.StageRetrying), string(runtimeview.StageBlocked), string(runtimeview.StageQuarantined), string(runtimeview.StageUnknown)} {
		if strings.HasPrefix(rest, stage+"\x00domain.delivery."+stage+"\x00") {
			return true
		}
	}
	return false
}

func (c *DomainCoalescer) deliveryFaultTokensLocked(v DomainEventRef) []deliveryFaultToken {
	var tokens []deliveryFaultToken // temporary, bounded by the existing capped entries
	for key, entry := range c.entries {
		if deliveryFaultKeyMatches(key, v) {
			tokens = append(tokens, deliveryFaultToken{key: key, reference: entry.reference, entry: entry})
		}
	}
	return tokens
}

func (c *DomainCoalescer) deliveryFaultTokensCurrentLocked(v DomainEventRef, tokens []deliveryFaultToken) bool {
	count := 0
	for key := range c.entries {
		if deliveryFaultKeyMatches(key, v) {
			count++
		}
	}
	if count != len(tokens) {
		return false
	}
	for _, token := range tokens {
		if c.entries[token.key] != token.entry {
			return false
		}
	}
	return true
}

func (c *DomainCoalescer) deliverySuccessEntry(v DomainEventRef, raw []byte, scope runtimeview.Scope, now time.Time) domainCoalescerEntry {
	entry := domainCoalescerEntry{data: raw, scope: runtimeview.Scope{WorkspaceID: strings.Clone(scope.WorkspaceID), WorkspaceRevision: strings.Clone(scope.WorkspaceRevision), DeviceID: strings.Clone(scope.DeviceID)}, reference: v.Evidence.ReferenceID, first: now, last: now, next: now.Add(c.options.Heartbeat), occurrences: 1, reported: 1}
	entry.bytes = len(raw) + len(deliverySuccessKey(v)) + len(entry.reference) + len(scope.WorkspaceID) + len(scope.WorkspaceRevision) + len(scope.DeviceID) + 256
	return entry // retain this first receipt unchanged; it never proves other effects
}

func (c *DomainCoalescer) rejectDeliveryCapture() bool { c.drop(1); return false }
