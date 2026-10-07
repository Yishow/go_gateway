package diagnostics

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/cases"
)

// Stable errors never include input, file paths or raw application details.
var (
	ErrInvalidQuery  = errors.New("invalid diagnostic query")
	ErrInvalidCursor = errors.New("invalid diagnostic cursor")
	ErrCapacity      = errors.New("diagnostic stream capacity reached")
	ErrClosed        = errors.New("diagnostic stream closed")
	ErrOverflow      = errors.New("diagnostic stream overflow")
)

// Query bounds a current-process history query; Level is minimum severity.
type Query struct {
	Level, Source, Q, Before string
	Limit                    int
}

// Gap identifies an explicitly unavailable interval or process reset.
type Gap struct {
	Reason string `json:"reason"`
	From   string `json:"from,omitempty"`
	To     string `json:"to,omitempty"`
}

// Metadata reports retained boundaries, delivered progress and separate sink loss.
type Metadata struct {
	InstanceID         string `json:"instance_id"`
	Oldest             string `json:"oldest_cursor"`
	Latest             string `json:"latest_cursor"`
	Progress           string `json:"progress_cursor"`
	CaptureDropped     uint64 `json:"capture_dropped"`
	Evicted            uint64 `json:"evicted"`
	SubscriberOverflow uint64 `json:"subscriber_overflow"`
	FileDropped        uint64 `json:"file_dropped"`
	FileRecoveredGaps  uint64 `json:"file_recovered_gaps"`
	DiskErrors         uint64 `json:"disk_errors"`
	SinkHealth         string `json:"sink_health"`
	Gaps               []Gap  `json:"gaps,omitempty"`
}

// Snapshot contains chronological safe records at an atomic append boundary.
type Snapshot struct {
	Metadata
	Records []Event `json:"records"`
}

func validateQuery(q Query) (Query, error) {
	if q.Level != "" && levelRank(q.Level) < 0 {
		return q, ErrInvalidQuery
	}
	if q.Limit == 0 {
		q.Limit = 200
	}
	if q.Limit < 1 || q.Limit > 500 || !utf8.ValidString(q.Q) || utf8.RuneCountInString(q.Q) > 256 {
		return q, ErrInvalidQuery
	}
	switch q.Source {
	case "", "startup", "runtime", "shutdown", "http", "standard", "slog", "application":
	default:
		return q, ErrInvalidQuery
	}
	q.Q = strings.Clone(cases.Fold().String(q.Q))
	q.Level = strings.Clone(q.Level)
	q.Source = strings.Clone(q.Source)
	q.Before = strings.Clone(q.Before)
	return q, nil
}

// ValidateQuery validates external limits without exposing normalized text whose
// Unicode case folding could expand past the original character limit.
func ValidateQuery(q Query) (Query, error) { _, err := validateQuery(q); return q, err }
func matches(q Query, e Event) bool {
	if q.Level != "" && levelRank(e.Level) < levelRank(q.Level) || q.Source != "" && q.Source != e.Source {
		return false
	}
	if q.Q == "" {
		return true
	}
	fold := cases.Fold()
	if strings.Contains(fold.String(e.Message), q.Q) || strings.Contains(fold.String(e.Code), q.Q) {
		return true
	}
	for _, v := range e.Fields {
		if s, ok := v.(string); ok && strings.Contains(fold.String(s), q.Q) {
			return true
		}
	}
	return false
}
func parseCursor(cursor string) (instanceID string, number uint64, err error) {
	instance, sequence, ok := strings.Cut(cursor, ":")
	if !ok || !opaqueID(instance) || sequence == "" || len(sequence) > 20 || len(sequence) > 1 && sequence[0] == '0' {
		return "", 0, ErrInvalidCursor
	}
	for _, c := range sequence {
		if c < '0' || c > '9' {
			return "", 0, ErrInvalidCursor
		}
	}
	n, err := strconv.ParseUint(sequence, 10, 64)
	if err != nil {
		return "", 0, ErrInvalidCursor
	}
	return instance, n, nil
}

// ValidateCursor checks opaque-instance and decimal-sequence wire syntax.
func ValidateCursor(cursor string) error { _, _, err := parseCursor(cursor); return err }
