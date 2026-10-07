// Package diagnostics projects application events into bounded, safe local sinks.
package diagnostics

import "time"

// Capacity constants are hard count and encoded-byte bounds.
const (
	MaxEventBytes    = 8 << 10
	AdmissionRecords = 1024
	AdmissionBytes   = 2 << 20
	RingRecords      = 2000
	RingBytes        = 4 << 20
	ReplayRecords    = 500
	ReplayBytes      = 1 << 20
	LiveRecords      = 256
	LiveBytes        = 512 << 10
	MaxSubscribers   = 8
	MaxFileBytes     = 5 << 20
)

// Input contains a template code and candidate untrusted scalar fields.
// Raw messages, errors and nested values have no admission path.
type Input struct {
	Code   string
	Fields map[string]any
}

// Event is a safe application projection. Stored maps are never exposed directly.
type Event struct {
	InstanceID string         `json:"instance_id"`
	Sequence   string         `json:"sequence"`
	Timestamp  time.Time      `json:"timestamp"`
	Level      string         `json:"level"`
	Source     string         `json:"source"`
	Code       string         `json:"code"`
	Message    string         `json:"message"`
	Fields     map[string]any `json:"fields,omitempty"`
	Truncated  bool           `json:"truncated"`
}

// Cursor is independent of JavaScript numeric precision.
func (e Event) Cursor() string { return e.InstanceID + ":" + e.Sequence }

// Record aliases the safe event representation.
type Record = Event
