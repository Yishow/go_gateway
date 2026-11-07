package dbtarget

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"net"
	"slices"
	"strings"
	"time"
)

const (
	failureSummaryInterval = 10 * time.Minute
	failureConnection      = "connection"
	failureTimeout         = "timeout"
	failureValidation      = "validation"
)

// connectorFailureLog is owned by one connector worker, independently of its
// retry counter, which is capped to limit the backoff.
type connectorFailureLog struct {
	class        string
	revision     string
	firstFailure time.Time
	lastReport   time.Time
	count        int
}

func (state *connectorFailureLog) selectRevision(revision string) {
	if state.revision != revision {
		*state = connectorFailureLog{revision: revision}
	}
}

func (state *connectorFailureLog) failed(connectorID string, err error, now time.Time) {
	class := snapshotFailureClass(err)
	first := state.count == 0
	state.count++
	if first {
		state.firstFailure = now
	}
	if first || class != state.class {
		state.class = class
		state.lastReport = now
		slog.Error("flush database target snapshot", "connector_id", connectorID, "failure_reason", snapshotFailureReason(err),
			"failure_count", state.count, "first_failure_at", state.firstFailure)
		return
	}
	if now.Sub(state.lastReport) >= failureSummaryInterval {
		state.lastReport = now
		slog.Warn("database target snapshot is still failing", "connector_id", connectorID,
			"failure_count", state.count, "first_failure_at", state.firstFailure)
	}
}

func snapshotFailureReason(err error) string {
	if errors.Is(err, ErrDestinationUnreachable) {
		return failureConnection
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return failureTimeout
	}
	var networkErr net.Error
	if errors.As(err, &networkErr) {
		if networkErr.Timeout() {
			return failureTimeout
		}
		return failureConnection
	}
	return "write"
}

func (state *connectorFailureLog) recovered(connectorID string) {
	if state.count == 0 {
		return
	}
	slog.Info("database target snapshot recovered", "connector_id", connectorID, "failure_count", state.count)
	*state = connectorFailureLog{}
}

// Classify the underlying cause using sentinels, driver types and codes, never
// error text, which can change between attempts for the same failure.
func snapshotFailureClass(err error) string {
	classes := make(map[string]struct{})
	collectSnapshotFailureClasses(err, classes)
	return strings.Join(slices.Sorted(maps.Keys(classes)), "|")
}

func collectSnapshotFailureClasses(err error, classes map[string]struct{}) {
	if err == nil {
		return
	}
	if joined, ok := err.(interface{ Unwrap() []error }); ok {
		for _, cause := range joined.Unwrap() {
			collectSnapshotFailureClasses(cause, classes)
		}
		return
	}
	var networkErr net.Error
	if errors.As(err, &networkErr) {
		if networkErr.Timeout() {
			classes[failureTimeout] = struct{}{}
		} else {
			classes[failureConnection] = struct{}{}
		}
		return
	}
	if cause := errors.Unwrap(err); cause != nil {
		collectSnapshotFailureClasses(cause, classes)
		return
	}
	var class string
	switch {
	case errors.Is(err, ErrValidation):
		class = failureValidation
	case errors.Is(err, context.DeadlineExceeded):
		class = failureTimeout
	case errors.Is(err, context.Canceled):
		class = "canceled"
	case errors.Is(err, ErrDestinationUnreachable):
		class = failureConnection
	case errors.Is(err, ErrDestinationBlocked):
		class = "blocked"
	default:
		if insertClass := ClassifyInsertError(err); insertClass != InsertErrorTransient {
			class = string(insertClass)
		} else {
			class = fmt.Sprintf("%T", err)
		}
	}
	classes[class] = struct{}{}
}
