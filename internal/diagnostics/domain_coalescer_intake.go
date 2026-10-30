package diagnostics

import (
	"context"
	"encoding/json"
	"strings"

	"go-gateway/internal/datalink/runtimeview"
)

func intakeObservation(v DomainEventRef) bool {
	return v.Code == "domain.group.applied" || v.Code == "domain.group.paused" || v.Code == "domain.group.retiring"
}

// observeIntake retains one last verified state per exact original scope. Identical
// observations do not count as transitions or generate healthy heartbeat spam.
// Historical emitted records remain immutable when the retained state changes.
func (c *DomainCoalescer) observeIntake(ctx context.Context, data []byte, scope runtimeview.Scope) bool {
	if ctx == nil || ctx.Err() != nil {
		return c.rejectIntakeCapture()
	}
	v, err := ParseDomainEventRef(data, scope)
	if err != nil || !intakeObservation(v) {
		return c.rejectIntakeCapture()
	}
	key := "\x00domain.intake\x00" + domainEntityKey(v) + "\x00" + v.Binding.WorkspaceRevision
	if !c.mu.TryLock() {
		return c.rejectIntakeCapture()
	}
	expected, closed := c.entries[key], c.closed
	c.mu.Unlock()
	if closed {
		return c.rejectIntakeCapture()
	}
	var prior DomainEventRef
	if expected != nil {
		prior, err = ParseDomainEventRef(expected.data, scope)
		if err != nil || !intakeObservation(prior) {
			return c.rejectIntakeCapture()
		}
	}
	now, ok := c.now()
	if !ok {
		return c.rejectIntakeCapture()
	}
	raw, _ := json.Marshal(v)
	next := domainCoalescerEntry{data: raw, scope: runtimeview.Scope{WorkspaceID: strings.Clone(scope.WorkspaceID), WorkspaceRevision: strings.Clone(scope.WorkspaceRevision), DeviceID: strings.Clone(scope.DeviceID)}, reference: v.Evidence.ReferenceID, first: now, last: now, next: now.Add(c.options.Heartbeat), occurrences: 1, reported: 1}
	next.bytes = len(raw) + len(key) + len(next.reference) + len(scope.WorkspaceID) + len(scope.WorkspaceRevision) + len(scope.DeviceID) + 256
	event := domainEntryEvent(next, false)
	size := domainOutputBytes(event)
	if !c.mu.TryLock() {
		return c.rejectIntakeCapture()
	}
	old := c.entries[key]
	if c.closed || old != expected || (old != nil && (old.scope != scope || now.Before(old.last))) {
		c.mu.Unlock()
		return c.rejectIntakeCapture()
	}
	if old != nil && prior.Code == v.Code {
		c.mu.Unlock()
		return true
	}
	removedCount, removedBytes := 0, 0
	if old != nil {
		removedCount, removedBytes = 1, old.bytes
	}
	if !c.fitsLocked(2-removedCount, next.bytes+size-removedBytes) {
		c.mu.Unlock()
		return c.rejectIntakeCapture()
	}
	c.entries[key] = &next
	c.stateBytes += next.bytes - removedBytes
	c.queue.push(record{event: event, size: size})
	c.mu.Unlock()
	c.drain(ctx)
	return true
}
func (c *DomainCoalescer) rejectIntakeCapture() bool { c.drop(1); return false }
