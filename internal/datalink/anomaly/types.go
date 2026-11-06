// Package anomaly detects and persists acquisition anomalies: device
// disconnection, points that stop updating (stale) and points whose value stays
// unchanged (frozen). Detection only reads runtime observations and never
// issues device reads.
package anomaly

import (
	"context"
	"errors"
	"time"
)

// Kind classifies an acquisition anomaly.
type Kind string

const (
	// KindDisconnected is a device-level event: connection failure or open breaker.
	KindDisconnected Kind = "disconnected"
	// KindStale is a point without a successful read for longer than its threshold.
	KindStale Kind = "stale"
	// KindFrozen is a point whose value has not changed for longer than its threshold.
	KindFrozen Kind = "frozen"
)

// Valid reports whether k is a known anomaly kind.
func (k Kind) Valid() bool {
	return k == KindDisconnected || k == KindStale || k == KindFrozen
}

// State filters events by whether they have ended.
type State string

const (
	// StateActive selects events without an end time.
	StateActive State = "active"
	// StateEnded selects events with an end time.
	StateEnded State = "ended"
)

// Detail codes stored on events. They are stable reason codes, never raw
// backend error text, so operator views cannot leak endpoints or credentials.
const (
	ReasonConnectionFailed  = "connection-failed"
	ReasonCircuitOpen       = "circuit-open"
	ReasonBreakerOpen       = "breaker-open"
	DetailNoSuccessfulRead  = "no-successful-read"
	DetailValueUnchanged    = "value-unchanged"
	DetailMonitoringStopped = "monitoring-stopped"
	DetailFrozenDisabled    = "frozen-detection-disabled"
)

// DefaultRetention is the per-device retention target (design condition).
const DefaultRetention = 5000

// Event is one persisted anomaly with a start and an optional end.
type Event struct {
	ID        int64
	DeviceID  string
	PointID   string // empty for device-level (disconnected) events
	Kind      Kind
	StartedAt time.Time
	EndedAt   *time.Time
	LastValue *string
	Detail    string
}

// ListFilter selects one device's events. Zero Limit means no limit.
type ListFilter struct {
	DeviceID string
	State    State
	Kind     Kind
	Limit    int
	Offset   int
}

var (
	// ErrAlreadyActive means an active event already exists for the device, point and kind.
	ErrAlreadyActive = errors.New("acquisition anomaly already active")
	// ErrNotActive means the event does not exist or has already ended.
	ErrNotActive = errors.New("acquisition anomaly is not active")
	// ErrInvalidEvent rejects malformed events before they reach storage.
	ErrInvalidEvent = errors.New("invalid acquisition anomaly")
	// ErrInvalidFilter rejects unknown list filters instead of ignoring them.
	ErrInvalidFilter = errors.New("invalid acquisition anomaly filter")
)

// Repository persists anomaly events. Open enforces at most one active event
// per device, point and kind and applies per-device retention, deleting the
// oldest ended events first and never deleting active ones.
type Repository interface {
	Open(ctx context.Context, event Event) (Event, error)
	End(ctx context.Context, id int64, endedAt time.Time, lastValue *string, endDetail string) error
	ListActive(ctx context.Context) ([]Event, error)
	List(ctx context.Context, filter ListFilter) ([]Event, int, error)
	ActiveCounts(ctx context.Context) (map[string]int, error)
}

func validateEvent(ev Event) error {
	switch {
	case !ev.Kind.Valid():
		return errors.Join(ErrInvalidEvent, errors.New("unknown kind"))
	case ev.DeviceID == "":
		return errors.Join(ErrInvalidEvent, errors.New("device is required"))
	case ev.Kind == KindDisconnected && ev.PointID != "":
		return errors.Join(ErrInvalidEvent, errors.New("disconnected events are device-level"))
	case ev.Kind != KindDisconnected && ev.PointID == "":
		return errors.Join(ErrInvalidEvent, errors.New("point is required"))
	case ev.StartedAt.IsZero():
		return errors.Join(ErrInvalidEvent, errors.New("start time is required"))
	case ev.EndedAt != nil:
		return errors.Join(ErrInvalidEvent, errors.New("new events must be active"))
	}
	return nil
}

func validateFilter(f ListFilter) error {
	switch {
	case f.DeviceID == "":
		return errors.Join(ErrInvalidFilter, errors.New("device is required"))
	case f.State != "" && f.State != StateActive && f.State != StateEnded:
		return errors.Join(ErrInvalidFilter, errors.New("unknown state"))
	case f.Kind != "" && !f.Kind.Valid():
		return errors.Join(ErrInvalidFilter, errors.New("unknown kind"))
	case f.Limit < 0 || f.Offset < 0:
		return errors.Join(ErrInvalidFilter, errors.New("limit and offset must not be negative"))
	}
	return nil
}
