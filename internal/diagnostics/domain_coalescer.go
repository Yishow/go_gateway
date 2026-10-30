package diagnostics

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go-gateway/internal/datalink/runtimeview"
)

// DomainCoalescerOptions injects the existing heartbeat (runtimeLogHeartbeat at wiring).
// Clock/TryAdmit/OnDrop MUST be nonblocking; callbacks run outside the state mutex.
// OnDrop is the sole adapter to existing capture-loss accounting: TryAdmit must not
// double-account rejected events. OnDrop must not recursively capture diagnostics.
// Limits cover entries plus queued/in-flight outputs together. They may be lowered
// for tests, never raised above the existing admission caps.
type DomainCoalescerOptions struct {
	Heartbeat              time.Duration
	Clock                  func() time.Time
	TryAdmit               func(Event) bool
	OnDrop                 func(uint64)
	RecordLimit, ByteLimit int
}

// DomainCoalescerStats separates encoded retained storage from capture loss.
type DomainCoalescerStats struct {
	Entries, StateBytes, Queued, QueueBytes int
	Dropped                                 uint64
}
type domainCoalescerEntry struct {
	unavailable           bool // fixed identity-free observation, never candidate data
	data                  []byte
	scope                 runtimeview.Scope
	reference             string
	first, last, next     time.Time
	quotaProbeAt          time.Time // scheduling metadata within the existing bounded entry
	occurrences, reported uint64
	bytes                 int
}

// DomainCoalescer is an unmounted, bounded observer, not delivery authority.
type DomainCoalescer struct {
	mu                                   sync.Mutex
	entries                              map[string]*domainCoalescerEntry
	options                              DomainCoalescerOptions
	queue                                recordQueue
	stateBytes, activeBytes, activeCount int
	closed                               bool
	quotaHintNext                        time.Time
	draining                             atomic.Bool
	dropped                              atomic.Uint64
}

// NewDomainCoalescer constructs a coalescer without background workers or replay.
func NewDomainCoalescer(o DomainCoalescerOptions) (*DomainCoalescer, error) {
	if o.RecordLimit == 0 {
		o.RecordLimit = AdmissionRecords
	}
	if o.ByteLimit == 0 {
		o.ByteLimit = AdmissionBytes
	}
	if o.Heartbeat <= 0 || o.Clock == nil || o.TryAdmit == nil || o.OnDrop == nil || o.RecordLimit < 1 || o.RecordLimit > AdmissionRecords || o.ByteLimit < 1 || o.ByteLimit > AdmissionBytes {
		return nil, errors.New("invalid domain coalescer configuration")
	}
	return &DomainCoalescer{entries: map[string]*domainCoalescerEntry{}, options: o, queue: newQueue(o.RecordLimit)}, nil
}

// Observe returns local acceptance (including intentional coalescing), NOT delivery
// success or admission success. A contended/full/closed capture is loss-accounted.
// Recovery codes require Recover's explicit original fault; no inferred clearance.
func (c *DomainCoalescer) Observe(ctx context.Context, data []byte, scope runtimeview.Scope) bool {
	if ctx == nil || ctx.Err() != nil {
		c.drop(1)
		return false
	}
	v, err := ParseDomainEventRef(data, scope)
	if err != nil {
		return c.enqueue(ctx, ProjectDomainEvent(data, scope))
	}
	if v.Code == "domain.collector.recovered" || v.Code == "domain.mapping.recovered" || v.Code == "domain.delivery.recovered" || v.Code == "domain.quota.recovered" {
		c.drop(1)
		return false
	}
	now, ok := c.now()
	if !ok {
		c.drop(1)
		return false
	}
	raw, _ := json.Marshal(v) // canonical validated reference; no raw candidate retention
	key := domainCoalescerKey(v)
	entry := domainCoalescerEntry{data: raw, scope: runtimeview.Scope{WorkspaceID: strings.Clone(scope.WorkspaceID), WorkspaceRevision: strings.Clone(scope.WorkspaceRevision), DeviceID: strings.Clone(scope.DeviceID)}, reference: v.Evidence.ReferenceID, first: now, last: now, next: now.Add(c.options.Heartbeat), occurrences: 1, reported: 1}
	// Encoded storage budget, like the broker: include key/reference/scope and a
	// fixed bookkeeping reserve. This is not a claim about Go allocator heap bytes.
	entry.bytes = len(raw) + len(key) + len(entry.reference) + len(scope.WorkspaceID) + len(scope.WorkspaceRevision) + len(scope.DeviceID) + 256
	e := domainEntryEvent(entry, false)
	size := domainOutputBytes(e)
	if !c.mu.TryLock() {
		c.drop(1)
		return false
	}
	if c.closed {
		c.mu.Unlock()
		c.drop(1)
		return false
	}
	if old := c.entries[key]; old != nil {
		if old.occurrences == math.MaxUint64 || now.Before(old.last) {
			c.mu.Unlock()
			c.drop(1)
			return false
		}
		old.occurrences++
		old.last = now
		c.mu.Unlock()
		c.drain(ctx)
		return true
	}
	if !c.fitsLocked(2, entry.bytes+size) {
		c.mu.Unlock()
		c.drop(1)
		return false
	}
	c.entries[key] = &entry
	c.stateBytes += entry.bytes
	c.queue.push(record{event: e, size: size})
	c.mu.Unlock()
	c.drain(ctx)
	return true
}

// Recover clears only an explicitly matched fault.
func (c *DomainCoalescer) Recover(ctx context.Context, recovery, failed []byte, scope runtimeview.Scope) bool {
	if ctx == nil || ctx.Err() != nil {
		c.drop(1)
		return false
	}
	r, re := ParseDomainEventRef(recovery, scope)
	f, fe := ParseDomainEventRef(failed, scope)
	now, validTime := c.now()
	if re != nil || fe != nil || !validTime || !domainRecoveryPair(r, f) || domainEntityKey(r) != domainEntityKey(f) {
		c.drop(1)
		return false
	}
	if !c.mu.TryLock() {
		c.drop(1)
		return false
	}
	key := domainCoalescerKey(f)
	entry := c.entries[key]
	if c.closed || entry == nil || entry.reference != f.Evidence.ReferenceID || now.Before(entry.last) {
		c.mu.Unlock()
		c.drop(1)
		return false
	}
	// Reserve the complete transition before releasing its fault identity. The
	// bounded pair is queued together, so callback reentry cannot place a new
	// episode between the old final summary and its recovery. No callback runs here.
	var batch [2]record
	count, bytes := 0, 0
	if entry.occurrences > entry.reported {
		e := domainEntryEvent(*entry, true)
		batch[count] = record{event: e, size: domainOutputBytes(e)}
		bytes += batch[count].size
		count++
	}
	e := ProjectDomainSummary(recovery, scope, DomainEventSummary{FirstSeenAt: now.Format(time.RFC3339Nano), LastSeenAt: now.Format(time.RFC3339Nano), Occurrences: "1"})
	batch[count] = record{event: e, size: domainOutputBytes(e)}
	bytes += batch[count].size
	count++
	delete(c.entries, key)
	c.stateBytes -= entry.bytes
	fits := c.fitsLocked(count, bytes)
	if fits {
		for _, output := range batch[:count] {
			c.queue.push(output)
		}
	}
	c.mu.Unlock()
	if !fits {
		c.drop(uint64(count)) // never retain a misleading partial transition
	}
	c.drain(ctx)
	return true // logging refusal never rewrites the confirmed lifecycle observation
}

// Tick flushes due summaries through the injected nonblocking admission.
func (c *DomainCoalescer) Tick(ctx context.Context) {
	if ctx == nil || ctx.Err() != nil {
		return
	}
	now, ok := c.now()
	if !ok || !c.mu.TryLock() {
		return
	}
	if c.closed {
		c.mu.Unlock()
		return
	}
	var pending []struct {
		key      string
		token    *domainCoalescerEntry
		snapshot domainCoalescerEntry
	} // bounded, sharing immutable bytes
	for key, e := range c.entries {
		if !now.Before(e.next) {
			e.next = now.Add(c.options.Heartbeat)
			if e.occurrences > e.reported {
				pending = append(pending, struct {
					key      string
					token    *domainCoalescerEntry
					snapshot domainCoalescerEntry
				}{key, e, *e})
				e.reported = e.occurrences
			}
		}
	}
	c.mu.Unlock()
	for _, e := range pending {
		c.enqueueCurrent(ctx, domainEntryEvent(e.snapshot, true), e.key, e.token)
	}
	c.drain(ctx)
}

// Close rejects new capture, discards queued events/one pending summary per dirty
// episode, and accounts that loss separately from occurrences. It does not wait on
// a callback already in flight, replay receipts, flush fake recovery, or spawn work.
func (c *DomainCoalescer) Close() {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return
	}
	c.closed = true
	lost := uint64(c.queue.count)
	for _, e := range c.entries {
		if e.occurrences > e.reported {
			lost++
		}
	}
	c.entries = nil
	c.stateBytes = 0
	for c.queue.count > 0 {
		c.queue.pop()
	}
	c.mu.Unlock()
	c.drop(lost)
}

// Stats returns a bounded accounting snapshot.
func (c *DomainCoalescer) Stats() DomainCoalescerStats {
	c.mu.Lock()
	defer c.mu.Unlock()
	return DomainCoalescerStats{Entries: len(c.entries), StateBytes: c.stateBytes, Queued: c.queue.count + c.activeCount, QueueBytes: c.queue.bytes + c.activeBytes, Dropped: c.dropped.Load()}
}

func (c *DomainCoalescer) now() (now time.Time, ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()
	now = c.options.Clock().UTC()
	return now, !now.IsZero() && now.Year() >= 1 && now.Year() <= 9999
}

func (c *DomainCoalescer) enqueue(ctx context.Context, e Event) bool {
	return c.enqueueCurrent(ctx, e, "", nil)
}

// A prepared Tick snapshot cannot reopen an episode that recovered or was replaced.
func (c *DomainCoalescer) enqueueCurrent(ctx context.Context, e Event, key string, expected *domainCoalescerEntry) bool {
	size := domainOutputBytes(e)
	if ctx == nil || ctx.Err() != nil || !c.mu.TryLock() {
		c.drop(1)
		return false
	}
	if c.closed || (expected != nil && c.entries[key] != expected) || !c.fitsLocked(1, size) {
		c.mu.Unlock()
		c.drop(1)
		return false
	}
	c.queue.push(record{event: e, size: size})
	c.mu.Unlock()
	c.drain(ctx)
	return true
}

// All retained records share one cap, including an output currently in TryAdmit.
func (c *DomainCoalescer) fitsLocked(records, bytes int) bool {
	return len(c.entries)+c.queue.count+c.activeCount+records <= c.options.RecordLimit && c.stateBytes+c.queue.bytes+c.activeBytes+bytes <= c.options.ByteLimit
}

// One caller drains a bounded batch; reentrant capture only enqueues. No callback
// runs under a mutex and no goroutine is created per event (or by this module).
func (c *DomainCoalescer) drain(ctx context.Context) {
	if !c.mu.TryLock() {
		return
	}
	if !c.draining.CompareAndSwap(false, true) {
		c.mu.Unlock()
		return
	}
	for range c.options.RecordLimit {
		if c.queue.count == 0 {
			break
		}
		r := c.queue.pop()
		c.activeCount++
		c.activeBytes += r.size
		c.mu.Unlock()
		accepted := false
		if ctx != nil && ctx.Err() == nil {
			accepted = c.tryAdmit(r.event)
		}
		c.mu.Lock()
		c.activeCount--
		c.activeBytes -= r.size
		if !accepted {
			c.mu.Unlock()
			c.drop(1)
			c.mu.Lock()
		}
	}
	// An enqueue cannot miss the handoff between an empty queue and releasing
	// drain ownership. The batch limit still bounds work under continuous input.
	c.draining.Store(false)
	c.mu.Unlock()
}

func (c *DomainCoalescer) tryAdmit(e Event) (accepted bool) {
	defer func() {
		if recover() != nil {
			accepted = false
		}
	}()
	return c.options.TryAdmit(e)
}
func (c *DomainCoalescer) drop(n uint64) {
	if n == 0 {
		return
	}
	for {
		old := c.dropped.Load()
		next := old + n
		if math.MaxUint64-old < n {
			next = math.MaxUint64
		}
		if c.dropped.CompareAndSwap(old, next) {
			break
		}
	}
	defer func() { _ = recover() }() // logging callbacks cannot fail a writer
	c.options.OnDrop(n)
}

func domainEntityKey(v DomainEventRef) string {
	// Workspace observation revision is intentionally NOT latest-draft identity.
	parts := []string{v.Binding.WorkspaceID, v.Binding.EntityRevision, strconv.Itoa(v.Locator.Step)}
	for _, p := range []*string{v.Locator.DeviceID, v.Locator.RuleID, v.Locator.MappingID, v.Locator.GroupID, v.PointID, v.TagID, v.ConnectorID, v.ConnectorRevision} {
		if p == nil {
			parts = append(parts, "")
		} else {
			parts = append(parts, *p)
		}
	}
	return strings.Join(parts, "\x00")
}
func domainCoalescerKey(v DomainEventRef) string {
	// Owner reason is part of the fault code identity; distinct repair guidance
	// cannot disappear under the first occurrence's canonical reason/action.
	key := domainEntityKey(v) + "\x00" + v.Binding.WorkspaceRevision + "\x00" + string(v.Stage) + "\x00" + string(v.Code) + "\x00" + string(v.Reason)
	if deliveryFaultCode(v.Code) {
		key += "\x00" + v.Evidence.ReferenceID // distinct durable effects never share fault episodes
	}
	if v.Evidence.Kind == "operation-result" {
		key += "\x00" + v.Evidence.ReferenceID // distinct repair/cutover ledger operations never merge
	}
	if v.Code == "domain.quota.refused" && v.Evidence.Quota != nil {
		key += "\x00" + v.Evidence.Quota.Scope // global and group capacity are distinct
	}
	return key
}
func domainEntryEvent(e domainCoalescerEntry, summary bool) Event {
	if e.unavailable {
		return unavailableCoalescerEvent(e, summary)
	}
	data := e.data
	if summary {
		var v DomainEventRef
		_ = json.Unmarshal(data, &v)
		if v.Code == "domain.delivery.first-confirmed" {
			v.Code = "domain.delivery.summary"
			data, _ = json.Marshal(v)
		}
	}
	return ProjectDomainSummary(data, e.scope, DomainEventSummary{FirstSeenAt: e.first.Format(time.RFC3339Nano), LastSeenAt: e.last.Format(time.RFC3339Nano), Occurrences: strconv.FormatUint(e.occurrences, 10)})
}
func domainOutputBytes(e Event) int {
	e.InstanceID = "ffffffffffffffffffffffffffffffff"
	e.Sequence = "18446744073709551615"
	e.Timestamp = time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC)
	raw, err := json.Marshal(e)
	if err != nil || len(raw) > MaxEventBytes {
		return AdmissionBytes + 1
	}
	return len(raw)
}
func domainRecoveryPair(r, f DomainEventRef) bool {
	if f.Reason == runtimeview.ReasonNone {
		return false
	}
	switch r.Code {
	case "domain.source.applied":
		return f.Code == "domain.source.rejected"
	case "domain.collector.recovered":
		return f.Code == "domain.collector.failed"
	case "domain.mapping.recovered":
		return f.Code == "domain.mapping.failed"
	case "domain.repair.completed":
		return f.Code == "domain.repair.blocked" && r.Evidence.ReferenceID == f.Evidence.ReferenceID
	case "domain.cutover.applied":
		return f.Code == "domain.cutover.rejected" && r.Evidence.ReferenceID == f.Evidence.ReferenceID
	case "domain.delivery.recovered":
		return (f.Code == "domain.delivery.retrying" || f.Code == "domain.delivery.blocked" || f.Code == "domain.delivery.quarantined" || f.Code == "domain.delivery.unknown") && r.Evidence.ReferenceID == f.Evidence.ReferenceID
	default:
		return false
	}
}
