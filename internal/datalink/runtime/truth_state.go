package runtime

import (
	"strings"
	"time"

	"go-gateway/internal/datalink/device"
)

const (
	// RuntimeTruthStateReady means runtime data is available and not degraded.
	RuntimeTruthStateReady = "ready"
	// RuntimeTruthStateEmpty means runtime data is validly absent.
	RuntimeTruthStateEmpty = "empty"
	// RuntimeTruthStateDegraded means runtime data exists but has degraded quality.
	RuntimeTruthStateDegraded = "degraded"
	// RuntimeTruthStateUnavailable means runtime data cannot be observed right now.
	RuntimeTruthStateUnavailable = "unavailable"
	// RuntimeTruthStateStale means runtime data exists but is older than expected.
	RuntimeTruthStateStale = "stale"
)

// RuntimeTruthState explicitly describes whether runtime data is usable, empty, or degraded.
type RuntimeTruthState struct {
	State       string `json:"state"`
	Empty       bool   `json:"empty"`
	Degraded    bool   `json:"degraded"`
	Unavailable bool   `json:"unavailable"`
	Stale       bool   `json:"stale"`
	Reason      string `json:"reason,omitempty"`
}

// RuntimeStreamStateEvent describes stream availability without requiring clients to infer it.
type RuntimeStreamStateEvent struct {
	DeviceID    string            `json:"device_id"`
	StreamState RuntimeTruthState `json:"stream_state"`
	Timestamp   time.Time         `json:"timestamp"`
	Code        string            `json:"code,omitempty"`
	Message     string            `json:"message,omitempty"`
	Retryable   bool              `json:"retryable"`
	Action      string            `json:"action,omitempty"`
	RequestID   string            `json:"request_id,omitempty"`
}

// RuntimeReadyTruthState returns a ready runtime truth marker.
func RuntimeReadyTruthState() RuntimeTruthState {
	return runtimeTruthState(RuntimeTruthStateReady, "")
}

// RuntimeUnavailableTruthState returns an unavailable runtime truth marker.
func RuntimeUnavailableTruthState(reason string) RuntimeTruthState {
	return runtimeTruthState(RuntimeTruthStateUnavailable, reason)
}

// RuntimeDegradedTruthState returns a degraded runtime truth marker.
func RuntimeDegradedTruthState(reason string) RuntimeTruthState {
	return runtimeTruthState(RuntimeTruthStateDegraded, reason)
}

// RuntimeStaleTruthState returns a stale runtime truth marker.
func RuntimeStaleTruthState(reason string) RuntimeTruthState {
	return runtimeTruthState(RuntimeTruthStateStale, reason)
}

// RuntimeEmptyTruthState returns an empty runtime truth marker.
func RuntimeEmptyTruthState(reason string) RuntimeTruthState {
	return runtimeTruthState(RuntimeTruthStateEmpty, reason)
}

// DeriveSnapshotTruthState derives snapshot truth from collector state.
func DeriveSnapshotTruthState(collectors []DeviceRuntimeStatus, selectedDeviceID string) RuntimeTruthState {
	if len(collectors) == 0 {
		if strings.TrimSpace(selectedDeviceID) != "" {
			return RuntimeEmptyTruthState("selected device has no runtime snapshot")
		}
		return RuntimeEmptyTruthState("runtime has no device snapshot")
	}

	allCollectorsEmpty := true
	hasUnavailable := false
	hasStale := false
	hasDegraded := false
	for _, collector := range collectors {
		if collector.PointsTotal > 0 || collector.LastReadAt != nil {
			allCollectorsEmpty = false
		}
		if collector.AvailabilityStatus == device.AvailabilityStatusUnavailable {
			hasUnavailable = true
			continue
		}
		if collector.PointsStale > 0 || strings.EqualFold(collector.ProjectionAlignment, "stale") {
			hasStale = true
		}
		if collector.Status == "error" ||
			collector.Status == "warning" ||
			collector.PointsError > 0 ||
			collector.BreakerState == "open" ||
			collector.BreakerState == "half-open" {
			hasDegraded = true
		}
	}

	if hasUnavailable {
		return RuntimeUnavailableTruthState("one or more runtime devices are unavailable")
	}
	if allCollectorsEmpty {
		return RuntimeEmptyTruthState("runtime snapshot has no point data")
	}
	if hasStale {
		return RuntimeStaleTruthState("one or more runtime points are stale")
	}
	if hasDegraded {
		return RuntimeDegradedTruthState("one or more runtime collectors are degraded")
	}
	return RuntimeReadyTruthState()
}

func runtimeTruthState(state string, reason string) RuntimeTruthState {
	truth := RuntimeTruthState{
		State:  state,
		Reason: strings.TrimSpace(reason),
	}
	switch state {
	case RuntimeTruthStateEmpty:
		truth.Empty = true
	case RuntimeTruthStateDegraded:
		truth.Degraded = true
	case RuntimeTruthStateUnavailable:
		truth.Unavailable = true
	case RuntimeTruthStateStale:
		truth.Stale = true
	}
	return truth
}
