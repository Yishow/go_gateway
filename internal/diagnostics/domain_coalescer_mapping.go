package diagnostics

import (
	"context"
	"time"

	"go-gateway/internal/datalink/runtimeview"
)

func fixedMappingObservation(v DomainEventRef) bool {
	return v.Code == domainMappingRecovered || (v.Code == domainMappingFailed && v.Reason == runtimeview.ReasonUnavailable && v.Action == runtimeview.ActionRefresh)
}

// observeMapping retires only an exact retained transform fault. Healthy samples
// create no state or output. Reservation is atomic within the existing coalescer;
// each eventual broker/file/subscriber output remains independently loss-accounted.
func (c *DomainCoalescer) observeMapping(ctx context.Context, data []byte, scope runtimeview.Scope, visible bool) bool {
	if ctx == nil || ctx.Err() != nil {
		if !visible {
			return false
		}
		return c.rejectMappingCapture()
	}
	recovered, err := ParseDomainEventRef(data, scope)
	if err != nil || !fixedMappingObservation(recovered) {
		return c.rejectMappingCapture()
	}
	if recovered.Code == domainMappingFailed {
		return c.Observe(ctx, data, scope)
	}
	failed := recovered
	failed.Code = domainMappingFailed
	// Failure and recovery intentionally have different stages. Use the failed
	// stage to locate the original episode, never the incoming acquisition stage.
	failed.Stage = runtimeview.StageBlocked
	failed.Reason = runtimeview.ReasonUnavailable
	failed.Action = runtimeview.ActionRefresh
	key := domainCoalescerKey(failed)
	if !c.mu.TryLock() {
		if !visible {
			return false
		}
		return c.rejectMappingCapture()
	}
	expected := c.entries[key]
	closed := c.closed
	c.mu.Unlock()
	if expected == nil {
		return false
	} // absence of retained fault is not a recovery
	if closed {
		return c.rejectMappingCapture()
	}
	original, err := ParseDomainEventRef(expected.data, scope)
	if err != nil || !domainRecoveryPair(recovered, original) || domainCoalescerKey(original) != key {
		return c.rejectMappingCapture()
	}
	now, ok := c.now()
	if !ok {
		return c.rejectMappingCapture()
	}
	if !c.mu.TryLock() {
		return c.rejectMappingCapture()
	}
	entry := c.entries[key]
	if c.closed || entry != expected || entry.reference != original.Evidence.ReferenceID || entry.scope != scope || now.Before(entry.last) {
		c.mu.Unlock()
		return c.rejectMappingCapture()
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
	// On refusal retain the complete old episode. No success entry is needed: a
	// later healthy sample without this exact fault is intentionally ignored.
	if !c.fitsLocked(count-1, bytes-entry.bytes) {
		c.mu.Unlock()
		return c.rejectMappingCapture()
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

func (c *DomainCoalescer) rejectMappingCapture() bool { c.drop(1); return false }
