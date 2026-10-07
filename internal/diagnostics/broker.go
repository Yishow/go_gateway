package diagnostics

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"math"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

// Options configures safe capture and its optional independently queued file sink.
type Options struct {
	Level string
	File  *FileSink
}

// Broker owns only bounded safe event representations.
type Broker struct {
	mu                                                                       sync.Mutex
	routesMu                                                                 sync.RWMutex
	routes                                                                   map[string]struct{}
	instance                                                                 string
	level                                                                    int
	file                                                                     *FileSink
	queue, ring                                                              recordQueue
	next, latest, completed, dropped, evicted, overflow, dropFirst, dropLast uint64
	subs                                                                     map[*Subscription]struct{}
	closing, streamsClosed                                                   bool
	wake, done, changed                                                      chan struct{}
}

// New creates a process-unique broker and starts its bounded dispatcher.
func New(o Options) (*Broker, error) {
	level := o.Level
	if level == "" {
		level = "info"
	}
	if levelRank(level) < 0 {
		return nil, ErrInvalidQuery
	}
	id := make([]byte, 16)
	if _, err := rand.Read(id); err != nil {
		return nil, err
	}
	b := &Broker{instance: hex.EncodeToString(id), level: levelRank(level), file: o.File, routes: map[string]struct{}{"unmatched": {}}, queue: newQueue(AdmissionRecords), ring: newQueue(RingRecords), subs: make(map[*Subscription]struct{}), wake: make(chan struct{}, 1), done: make(chan struct{}), changed: make(chan struct{})}
	go b.run()
	return b, nil
}

// RegisterRoutes admits only application route templates, never request paths.
func (b *Broker) RegisterRoutes(routes []string) {
	b.routesMu.Lock()
	defer b.routesMu.Unlock()
	for _, route := range routes {
		if len(b.routes) >= 4096 {
			break
		}
		if route != "" && len(route) <= 1024 && utf8.ValidString(route) && !strings.ContainsAny(route, "?\r\n\x00") {
			b.routes[strings.Clone(route)] = struct{}{}
		}
	}
}

// Emit projects before admission and never waits for disk or network IO.
func (b *Broker) Emit(in Input) bool {
	b.routesMu.RLock()
	e := project(in, b.routes)
	b.routesMu.RUnlock()
	if levelRank(e.Level) < b.level {
		return false
	}
	e.InstanceID = b.instance
	e.Timestamp = time.Now().UTC()
	e, size := boundEvent(e)
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closing || b.next == math.MaxUint64 {
		return false
	}
	b.next++
	e.Sequence = strconv.FormatUint(b.next, 10)
	if b.queue.count >= AdmissionRecords || b.queue.bytes+size > AdmissionBytes {
		increment(&b.dropped)
		if b.dropFirst == 0 {
			b.dropFirst = b.next
		}
		b.dropLast = b.next
		b.signal()
		return false
	}
	b.queue.push(record{event: e, sequence: b.next, size: size})
	b.signal()
	return true
}
func (b *Broker) signal() {
	select {
	case b.wake <- struct{}{}:
	default:
	}
}
func (b *Broker) changedLocked()         { close(b.changed); b.changed = make(chan struct{}) }
func (b *Broker) cursor(n uint64) string { return b.instance + ":" + strconv.FormatUint(n, 10) }
func (b *Broker) run() {
	defer close(b.done)
	for range b.wake {
		for {
			b.mu.Lock()
			if b.queue.count == 0 {
				b.latest = b.next
				b.completed = b.next
				b.changedLocked()
				closing := b.closing
				b.mu.Unlock()
				if closing {
					return
				}
				break
			}
			r := b.queue.pop()
			for b.ring.count >= RingRecords || b.ring.bytes+r.size > RingBytes {
				b.ring.pop()
				increment(&b.evicted)
			}
			b.ring.push(r)
			b.latest = r.sequence
			for s := range b.subs {
				s.admitLocked(r)
			}
			b.mu.Unlock()
			if b.file != nil {
				b.file.enqueue(r)
			}
			b.mu.Lock()
			b.completed = r.sequence
			b.changedLocked()
			b.mu.Unlock()
		}
	}
}

// Flush waits for already admitted records, then the file writer.
func (b *Broker) Flush(ctx context.Context) error {
	b.mu.Lock()
	target := b.next
	for b.completed < target {
		changed := b.changed
		b.mu.Unlock()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-changed:
		}
		b.mu.Lock()
	}
	b.mu.Unlock()
	if b.file != nil {
		return b.file.Flush(ctx)
	}
	return nil
}

// CloseSubscriptions stops web readers immediately while diagnostics can continue.
func (b *Broker) CloseSubscriptions() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.streamsClosed = true
	for s := range b.subs {
		s.stopLocked(ErrClosed)
	}
}

// Close drains accepted records; a deadline-limited wait is safely retryable.
func (b *Broker) Close(ctx context.Context) error {
	b.CloseSubscriptions()
	b.mu.Lock()
	b.closing = true
	b.signal()
	b.mu.Unlock()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-b.done:
	}
	if b.file != nil {
		return b.file.Close(ctx)
	}
	return nil
}

// Snapshot reads a bounded chronological page at one atomic append boundary.
func (b *Broker) Snapshot(query Query) (Snapshot, error) {
	q, err := validateQuery(query)
	if err != nil {
		return Snapshot{}, err
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	before := uint64(math.MaxUint64)
	if q.Before != "" {
		instance, n, err := parseCursor(q.Before)
		if err != nil || instance != b.instance || n > b.latest {
			return Snapshot{}, ErrInvalidCursor
		}
		before = n
	}
	out := Snapshot{Metadata: b.metadataLocked(), Records: make([]Event, 0, q.Limit)}
	for i := b.ring.count - 1; i >= 0 && len(out.Records) < q.Limit; i-- {
		r := b.ring.at(i)
		if r.sequence < before && matches(q, r.event) {
			out.Records = append(out.Records, cloneEvent(r.event))
		}
	}
	slices.Reverse(out.Records)
	return out, nil
}
func (b *Broker) metadataLocked() Metadata {
	m := Metadata{InstanceID: b.instance, Latest: b.cursor(b.latest), Progress: b.cursor(b.latest), CaptureDropped: b.dropped, Evicted: b.evicted, SubscriberOverflow: b.overflow, SinkHealth: "disabled"}
	if b.ring.count > 0 {
		m.Oldest = b.cursor(b.ring.first().sequence)
	}
	if b.dropped > 0 {
		m.Gaps = append(m.Gaps, Gap{Reason: "admission", From: b.cursor(b.dropFirst), To: b.cursor(b.dropLast)})
	}
	if b.overflow > 0 {
		m.Gaps = append(m.Gaps, Gap{Reason: "subscriber_overflow"})
	}
	if b.file != nil {
		s := b.file.health()
		m.FileDropped = s.Dropped
		m.FileRecoveredGaps = s.RecoveredGaps
		m.DiskErrors = s.DiskErrors
		m.SinkHealth = "healthy"
		if s.Degraded {
			m.SinkHealth = "degraded"
		}
		if s.Dropped > 0 {
			m.Gaps = append(m.Gaps, Gap{Reason: "file_loss", From: s.LostFrom, To: s.LostTo})
		}
	}
	return m
}
