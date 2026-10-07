package diagnostics

import (
	"context"
	"strings"
)

// Subscription retains separate bounded replay and live queues. Ack follows a
// successful network write; its progress never skips an undelivered match.
type Subscription struct {
	broker       *Broker
	query        Query
	replay, live recordQueue
	inflight     uint64
	gaps         []Gap
	wake, done   chan struct{}
	err          error
}

// Subscribe atomically captures replay cutoff and registers future live events.
func (b *Broker) Subscribe(query Query, after string) (*Subscription, error) {
	q, err := validateQuery(query)
	if err != nil {
		return nil, err
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closing || b.streamsClosed {
		return nil, ErrClosed
	}
	if len(b.subs) >= MaxSubscribers {
		return nil, ErrCapacity
	}
	var seq uint64
	var gaps []Gap
	if after != "" {
		instance, n, err := parseCursor(after)
		if err != nil {
			return nil, err
		}
		if instance != b.instance {
			gaps = append(gaps, Gap{Reason: "reset", From: strings.Clone(after), To: b.cursor(0)})
		} else {
			if n > b.latest {
				return nil, ErrInvalidCursor
			}
			seq = n
		}
	}
	if b.ring.count > 0 && seq < b.ring.first().sequence-1 {
		gaps = append(gaps, Gap{Reason: "retention", From: b.cursor(seq), To: b.cursor(b.ring.first().sequence - 1)})
	}
	s := &Subscription{broker: b, query: q, replay: newQueue(ReplayRecords), live: newQueue(LiveRecords), gaps: gaps, wake: make(chan struct{}, 1), done: make(chan struct{})}
	var lost uint64
	for i := range b.ring.count {
		r := b.ring.at(i)
		if r.sequence <= seq || !matches(q, r.event) {
			continue
		}
		for s.replay.count >= ReplayRecords || s.replay.bytes+r.size > ReplayBytes {
			lost = s.replay.pop().sequence
		}
		s.replay.push(r)
	}
	if lost > 0 {
		s.gaps = append(s.gaps, Gap{Reason: "replay_limit", From: b.cursor(seq), To: b.cursor(lost)})
	}
	b.subs[s] = struct{}{}
	return s, nil
}
func (s *Subscription) admitLocked(r record) {
	if s.err != nil || !matches(s.query, r.event) {
		return
	}
	if s.live.count >= LiveRecords || s.live.bytes+r.size > LiveBytes {
		increment(&s.broker.overflow)
		s.stopLocked(ErrOverflow)
		return
	}
	s.live.push(r)
	select {
	case s.wake <- struct{}{}:
	default:
	}
}
func (s *Subscription) stopLocked(err error) {
	if s.err != nil {
		return
	}
	s.err = err
	clear(s.replay.entries)
	clear(s.live.entries)
	s.replay.count = 0
	s.live.count = 0
	s.replay.bytes = 0
	s.live.bytes = 0
	close(s.done)
}

// Next returns a replay event before any live event. Use one reader per subscription.
func (s *Subscription) Next(ctx context.Context) (Event, error) {
	for {
		s.broker.mu.Lock()
		if s.err != nil {
			err := s.err
			s.broker.mu.Unlock()
			return Event{}, err
		}
		if s.inflight != 0 {
			s.broker.mu.Unlock()
			return Event{}, ErrInvalidCursor
		}
		var r record
		if s.replay.count > 0 {
			r = s.replay.pop()
		} else if s.live.count > 0 {
			r = s.live.pop()
		}
		if r.sequence != 0 {
			s.inflight = r.sequence
			s.broker.mu.Unlock()
			return cloneEvent(r.event), nil
		}
		s.broker.mu.Unlock()
		select {
		case <-ctx.Done():
			return Event{}, ctx.Err()
		case <-s.done:
		case <-s.wake:
		}
	}
}

// Ack acknowledges the sole in-flight matching record.
func (s *Subscription) Ack(cursor string) {
	s.broker.mu.Lock()
	defer s.broker.mu.Unlock()
	if cursor == s.broker.cursor(s.inflight) {
		s.inflight = 0
	}
}

// Metadata advertises progress only after earlier matches are delivered or gapped.
// Emit its gap controls before advertising the progress cursor.
func (s *Subscription) Metadata() Metadata {
	s.broker.mu.Lock()
	defer s.broker.mu.Unlock()
	m := s.broker.metadataLocked()
	m.Gaps = append(m.Gaps, s.gaps...)
	progress := s.broker.latest
	if s.inflight > 0 {
		progress = min(progress, s.inflight-1)
	}
	if s.replay.count > 0 {
		progress = min(progress, s.replay.first().sequence-1)
	}
	if s.live.count > 0 {
		progress = min(progress, s.live.first().sequence-1)
	}
	m.Progress = s.broker.cursor(progress)
	return m
}

// Close releases the slot only after the HTTP owner releases this subscription.
func (s *Subscription) Close() {
	s.broker.mu.Lock()
	defer s.broker.mu.Unlock()
	s.stopLocked(ErrClosed)
	delete(s.broker.subs, s)
}
