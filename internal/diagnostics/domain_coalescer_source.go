package diagnostics

import (
	"context"
	"encoding/json"
	"math"
	"strings"
	"time"

	"go-gateway/internal/datalink/runtimeview"
)

// Only the confirmed source producer's fixed pair uses this transition. Other
// reasons are not guessed recovered; generic Observe/Recover remain unchanged.
func fixedSourceObservation(v DomainEventRef) bool {
	return v.Code == "domain.source.applied" || (v.Code == "domain.source.rejected" && v.Reason == runtimeview.ReasonUnavailable && v.Action == runtimeview.ActionRefresh)
}

// observeSource reserves an exact source transition in the existing bounded state
// and output queue. No callback runs under the lock. The initial token snapshot is
// denial-only: a replaced episode between observation and reservation cannot be
// cleared. The final reservation, summaries and new state commit under one lock.
func (c *DomainCoalescer) observeSource(ctx context.Context, data []byte, scope runtimeview.Scope, visible bool) bool {
	if ctx == nil || ctx.Err() != nil {
		if !visible {
			return false
		}
		return c.rejectSourceCapture()
	}
	v, err := ParseDomainEventRef(data, scope)
	if err != nil || !fixedSourceObservation(v) {
		return c.rejectSourceCapture()
	}
	raw, _ := json.Marshal(v)
	key, otherKey := domainCoalescerKey(v), sourceOppositeKey(v)
	if !c.mu.TryLock() {
		if !visible {
			return false
		}
		return c.rejectSourceCapture()
	}
	expected, expectedOther := c.entries[key], c.entries[otherKey]
	closed := c.closed
	c.mu.Unlock()
	if !visible && (closed || expectedOther == nil) {
		return false
	}
	if closed {
		return c.rejectSourceCapture()
	}
	now, ok := c.now()
	if !ok {
		return c.rejectSourceCapture()
	}
	next := c.sourceCoalescerEntry(v, raw, scope, now)
	if !c.mu.TryLock() {
		return c.rejectSourceCapture()
	}
	old, other := c.entries[key], c.entries[otherKey]
	if c.closed || old != expected || other != expectedOther || (old != nil && (old.scope != scope || now.Before(old.last))) || (other != nil && (other.scope != scope || now.Before(other.last))) {
		c.mu.Unlock()
		return c.rejectSourceCapture()
	}
	if old != nil && other == nil {
		if old.occurrences == math.MaxUint64 {
			c.mu.Unlock()
			return c.rejectSourceCapture()
		}
		old.occurrences++
		old.last = now
		c.mu.Unlock()
		c.drain(ctx)
		return true
	}
	// At most two historical pending summaries and the one new transition. A
	// rejected reservation preserves old state; only this new observation is lost.
	var batch [3]record
	count, bytes, removedCount, removedBytes := 0, 0, 0, 0
	for _, entry := range []*domainCoalescerEntry{old, other} {
		if entry == nil {
			continue
		}
		removedCount++
		removedBytes += entry.bytes
		if entry.occurrences > entry.reported && (visible || entry == other) {
			e := domainEntryEvent(*entry, true)
			size := domainOutputBytes(e)
			batch[count] = record{event: e, size: size}
			count++
			bytes += size
		}
	}
	newCount, newBytes := 0, 0
	if visible {
		e := domainEntryEvent(next, false)
		size := domainOutputBytes(e)
		batch[count] = record{event: e, size: size}
		count++
		bytes += size
		newCount, newBytes = 1, next.bytes
	}
	if !c.fitsLocked(newCount+count-removedCount, newBytes+bytes-removedBytes) {
		c.mu.Unlock()
		return c.rejectSourceCapture()
	}
	delete(c.entries, key)
	delete(c.entries, otherKey)
	c.stateBytes -= removedBytes
	if visible {
		c.entries[key] = &next
		c.stateBytes += next.bytes
	}
	for _, r := range batch[:count] {
		c.queue.push(r)
	}
	c.mu.Unlock()
	c.drain(ctx)
	return visible
}

func (c *DomainCoalescer) rejectSourceCapture() bool { c.drop(1); return false }

func sourceOppositeKey(v DomainEventRef) string {
	if v.Code == "domain.source.applied" {
		v.Code = "domain.source.rejected"
		v.Stage = runtimeview.StageBlocked
		v.Reason = runtimeview.ReasonUnavailable
		v.Action = runtimeview.ActionRefresh
	} else {
		v.Code = "domain.source.applied"
		v.Stage = runtimeview.StageAccepted
		v.Reason = runtimeview.ReasonNone
		v.Action = runtimeview.ActionNone
	}
	return domainCoalescerKey(v)
}

func (c *DomainCoalescer) sourceCoalescerEntry(v DomainEventRef, raw []byte, scope runtimeview.Scope, now time.Time) domainCoalescerEntry {
	entry := domainCoalescerEntry{data: raw, scope: runtimeview.Scope{WorkspaceID: strings.Clone(scope.WorkspaceID), WorkspaceRevision: strings.Clone(scope.WorkspaceRevision), DeviceID: strings.Clone(scope.DeviceID)}, reference: v.Evidence.ReferenceID, first: now, last: now, next: now.Add(c.options.Heartbeat), occurrences: 1, reported: 1}
	entry.bytes = len(raw) + len(domainCoalescerKey(v)) + len(entry.reference) + len(scope.WorkspaceID) + len(scope.WorkspaceRevision) + len(scope.DeviceID) + 256
	return entry
}
