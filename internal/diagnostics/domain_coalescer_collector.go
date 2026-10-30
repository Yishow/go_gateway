package diagnostics

import (
	"context"
	"math"
	"strconv"
	"time"

	"go-gateway/internal/datalink/runtimeview"
)

func fixedCollectorObservation(v DomainEventRef) bool {
	return v.Code == "domain.collector.recovered" || (v.Code == "domain.collector.failed" && v.Reason == runtimeview.ReasonUnavailable && v.Action == runtimeview.ActionRefresh)
}

// observeCollector does not retain healthy samples. Only an exact, already
// observed failure can produce recovery. Reservation/state is atomic here; each
// eventual broker/file/subscriber output remains independently loss-accounted.
func (c *DomainCoalescer) observeCollector(ctx context.Context, data []byte, scope runtimeview.Scope, visible bool) bool {
	if ctx == nil || ctx.Err() != nil {
		if !visible {
			return false
		}
		return c.rejectCollectorCapture()
	}
	recovered, err := ParseDomainEventRef(data, scope)
	if err != nil || !fixedCollectorObservation(recovered) {
		return c.rejectCollectorCapture()
	}
	if recovered.Code == "domain.collector.failed" {
		return c.Observe(ctx, data, scope)
	}
	failed := recovered
	failed.Code = "domain.collector.failed"
	failed.Reason = runtimeview.ReasonUnavailable
	failed.Action = runtimeview.ActionRefresh
	key := domainCoalescerKey(failed)
	if !c.mu.TryLock() {
		if !visible {
			return false
		}
		return c.rejectCollectorCapture()
	}
	expected := c.entries[key]
	closed := c.closed
	c.mu.Unlock()
	if expected == nil {
		return false
	} // successful sampling is not per-sample log traffic
	if closed {
		return c.rejectCollectorCapture()
	}
	original, err := ParseDomainEventRef(expected.data, scope)
	if err != nil || !domainRecoveryPair(recovered, original) {
		return c.rejectCollectorCapture()
	}
	now, ok := c.now()
	if !ok {
		return c.rejectCollectorCapture()
	}
	if !c.mu.TryLock() {
		return c.rejectCollectorCapture()
	}
	entry := c.entries[key]
	if c.closed || entry != expected || entry.reference != original.Evidence.ReferenceID || entry.scope != scope || now.Before(entry.last) {
		c.mu.Unlock()
		return c.rejectCollectorCapture()
	}
	var batch [2]record
	count, bytes := 0, 0
	if entry.occurrences > entry.reported {
		e := domainEntryEvent(*entry, true)
		size := domainOutputBytes(e)
		batch[count] = record{event: e, size: size}
		count++
		bytes += size
	}
	if visible {
		e := ProjectDomainSummary(data, scope, DomainEventSummary{FirstSeenAt: now.Format(time.RFC3339Nano), LastSeenAt: now.Format(time.RFC3339Nano), Occurrences: "1"})
		size := domainOutputBytes(e)
		batch[count] = record{event: e, size: size}
		count++
		bytes += size
	}
	if !c.fitsLocked(count-1, bytes-entry.bytes) {
		c.mu.Unlock()
		return c.rejectCollectorCapture()
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

func (c *DomainCoalescer) rejectCollectorCapture() bool { c.drop(1); return false }

// No valid attributed domain key can start with NUL. This one shared entry stores
// only fixed unavailable bookkeeping, not guessed device/revision associations.
const unavailableCoalescerKey = "\x00domain.unavailable"

func (c *DomainCoalescer) observeUnavailable(ctx context.Context) bool {
	if ctx == nil || ctx.Err() != nil {
		return c.rejectCollectorCapture()
	}
	now, ok := c.now()
	if !ok {
		return c.rejectCollectorCapture()
	}
	if !c.mu.TryLock() {
		return c.rejectCollectorCapture()
	}
	if c.closed {
		c.mu.Unlock()
		return c.rejectCollectorCapture()
	}
	if entry := c.entries[unavailableCoalescerKey]; entry != nil {
		if entry.occurrences == math.MaxUint64 || now.Before(entry.last) {
			c.mu.Unlock()
			return c.rejectCollectorCapture()
		}
		entry.occurrences++
		entry.last = now
		c.mu.Unlock()
		c.drain(ctx)
		return true
	}
	entry := domainCoalescerEntry{unavailable: true, first: now, last: now, next: now.Add(c.options.Heartbeat), occurrences: 1, reported: 1, bytes: len(unavailableCoalescerKey) + 256}
	e := domainEntryEvent(entry, false)
	size := domainOutputBytes(e)
	if !c.fitsLocked(2, entry.bytes+size) {
		c.mu.Unlock()
		return c.rejectCollectorCapture()
	}
	c.entries[unavailableCoalescerKey] = &entry
	c.stateBytes += entry.bytes
	c.queue.push(record{event: e, size: size})
	c.mu.Unlock()
	c.drain(ctx)
	return true
}

func unavailableCoalescerEvent(entry domainCoalescerEntry, summary bool) Event {
	e := project(Input{Code: "domain.unavailable"}, nil)
	if summary {
		e.Fields = map[string]any{"first_seen_at": entry.first.Format(time.RFC3339Nano), "last_seen_at": entry.last.Format(time.RFC3339Nano), "occurrences": strconv.FormatUint(entry.occurrences, 10)}
	}
	return e
}
