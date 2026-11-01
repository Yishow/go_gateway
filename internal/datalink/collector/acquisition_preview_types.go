package collector

import "time"

// AcquisitionPreview contains only detached, safe scoped metadata.
type AcquisitionPreview struct {
	WorkspaceID       string                    `json:"workspace_id"`
	DeviceID          string                    `json:"device_id"`
	WorkspaceRevision string                    `json:"workspace_revision"`
	Status            string                    `json:"status"`
	Reason            string                    `json:"reason,omitempty"`
	PreviewRevision   string                    `json:"preview_revision,omitempty"`
	Estimate          bool                      `json:"estimate,omitempty"`
	Groups            []AcquisitionPreviewGroup `json:"groups"`
	OutputLoss        *AcquisitionPreviewLoss   `json:"output_loss,omitempty"`
}

// AcquisitionPreviewGroup preserves a real polling group's independent cadence.
type AcquisitionPreviewGroup struct {
	PollingGroupID        string                           `json:"polling_group_id"`
	PollIntervalMS        int64                            `json:"poll_interval_ms"`
	PlanRevision          string                           `json:"plan_revision"`
	DeclaredPlanRevision  string                           `json:"declared_plan_revision"`
	CapabilityRevision    string                           `json:"capability_revision"`
	Basis                 string                           `json:"basis"`
	AcquisitionPointCount int                              `json:"acquisition_point_count"`
	UniqueWireUnits       []AcquisitionPreviewUnits        `json:"unique_wire_units"`
	PlannedRequests       int                              `json:"planned_requests"`
	RequestLimits         []AcquisitionPreviewRequestLimit `json:"request_limits"`
	PointErrors           []AcquisitionPreviewPointError   `json:"point_errors"`
	FallbackReasons       []string                         `json:"fallback_reasons"`
	Actual                AcquisitionActual                `json:"actual"`
}

// AcquisitionPreviewUnits counts physical address cells, never read components.
type AcquisitionPreviewUnits struct {
	Area               string `json:"area"`
	Function           string `json:"function"`
	AddressUnitBits    int    `json:"address_unit_bits"`
	UniqueAddressUnits int    `json:"unique_address_units"`
}

// AcquisitionPreviewPointError retains rejected source membership safely.
type AcquisitionPreviewPointError struct {
	PointID string `json:"point_id"`
	Code    string `json:"code"`
}

// AcquisitionActual separates wire attempts from logical final outcomes.
type AcquisitionActual struct {
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
	*AcquisitionActualMetrics
}

// AcquisitionActualMetrics exists only for a matching completed scheduled scan.
type AcquisitionActualMetrics struct {
	AcquisitionID        string                     `json:"acquisition_id"`
	ObservedAt           time.Time                  `json:"observed_at"`
	DeclaredPlanRevision string                     `json:"declared_plan_revision"`
	InitialPlanRevision  string                     `json:"initial_plan_revision"`
	Phases               []AcquisitionActualPhase   `json:"phases"`
	WireAttempts         uint64                     `json:"wire_attempts"`
	RetryWireAttempts    uint64                     `json:"retry_wire_attempts"`
	FallbackWireAttempts uint64                     `json:"fallback_wire_attempts"`
	LogicalRequests      AcquisitionLogicalOutcomes `json:"logical_requests"`
	CycleLatencyMS       float64                    `json:"cycle_latency_ms"`
	Overrun              bool                       `json:"overrun"`
	LastGoodAgeMS        *float64                   `json:"last_good_age_ms"`
}

// AcquisitionActualPhase keeps fallback fingerprints distinct from the initial plan.
type AcquisitionActualPhase struct {
	Phase        string `json:"phase"`
	PlanRevision string `json:"plan_revision"`
}

// AcquisitionLogicalOutcomes omits counts when no request was entered.
type AcquisitionLogicalOutcomes struct {
	Status string `json:"status"`
	*AcquisitionLogicalCounts
}

// AcquisitionLogicalCounts counts one final outcome per entered request.
type AcquisitionLogicalCounts struct {
	Succeeded uint64 `json:"succeeded"`
	Failed    uint64 `json:"failed"`
	Canceled  uint64 `json:"canceled"`
}

// AcquisitionPreviewLoss exposes only retained matching detail, not global totals.
type AcquisitionPreviewLoss struct {
	Status          string                         `json:"status"`
	Reason          string                         `json:"reason,omitempty"`
	HistoryComplete bool                           `json:"history_complete"`
	RetainedDetail  []AcquisitionPreviewLossDetail `json:"retained_detail,omitempty"`
}

// AcquisitionPreviewLossDetail is explicitly incomplete output-channel loss.
type AcquisitionPreviewLossDetail struct {
	PointID string `json:"point_id"`
	Losses  uint64 `json:"losses"`
}

// AcquisitionPreviewRequestLimit is copied from the canonical request, without recomputation.
type AcquisitionPreviewRequestLimit struct {
	Kind               string `json:"kind"`
	Area               string `json:"area"`
	Function           string `json:"function"`
	NativeUnitLimit    int    `json:"native_unit_limit"`
	RequestFrameBytes  int    `json:"request_frame_bytes"`
	ResponseFrameBytes int    `json:"response_frame_bytes"`
	FrameCostKnown     bool   `json:"frame_cost_known"`
}
