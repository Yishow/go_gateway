package aggregation

import (
	"time"

	"go-gateway/internal/datalink/measurement"
)

// maxCounterEvents bounds the events kept per window so a flapping counter
// cannot grow a record without limit; the rest are only counted.
const maxCounterEvents = 16

// CounterEventKind names what was found between two adjacent readings.
type CounterEventKind string

const (
	// EventDecrease is a drop that was not proven to be a rollover.
	EventDecrease CounterEventKind = "decrease"
	// EventRollover is a drop compensated as one proven wrap of the modulus.
	EventRollover CounterEventKind = "rollover"
	// EventContradictory is two different values observed at one instant.
	EventContradictory CounterEventKind = "contradictory-same-time"
	// EventNonMonotonicTime is a reading observed before the one it follows.
	EventNonMonotonicTime CounterEventKind = "non-monotonic-time"
	// EventMixedBasis is a pair that has no exact difference.
	EventMixedBasis CounterEventKind = "mixed-numeric-basis"
)

// CounterEvent is one discontinuity or compensation, with the exact readings on
// both sides. It is part of the closed record and is serialized with it, so a
// retry reuses the same evidence and never recomputes it.
type CounterEvent struct {
	Kind CounterEventKind `json:"kind"`
	// Reason says why a drop was not compensated; empty for a compensation.
	Reason       string                 `json:"reason,omitempty"`
	FromSampleID string                 `json:"from_sample_id"`
	ToSampleID   string                 `json:"to_sample_id"`
	From         time.Time              `json:"from"`
	To           time.Time              `json:"to"`
	FromValue    measurement.ExactValue `json:"from_value"`
	ToValue      measurement.ExactValue `json:"to_value"`
}

// eventLog keeps the first maxCounterEvents events and counts the rest.
type eventLog struct {
	events  []CounterEvent
	omitted int
}

func (l *eventLog) add(kind CounterEventKind, reason string, from, to CounterReading) {
	if len(l.events) >= maxCounterEvents {
		l.omitted++
		return
	}
	l.events = append(l.events, CounterEvent{
		Kind: kind, Reason: reason, FromSampleID: from.SampleID, ToSampleID: to.SampleID,
		From: from.ObservedAt.UTC(), To: to.ObservedAt.UTC(), FromValue: from.Value, ToValue: to.Value,
	})
}
