package runtime

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// DeliveryError identifies the runtime target stage that rejected a value.
type DeliveryError struct {
	Target string
	Stage  string
	Err    error
}

func (e *DeliveryError) Error() string { return fmt.Sprintf("%s %s: %v", e.Target, e.Stage, e.Err) }
func (e *DeliveryError) Unwrap() error { return e.Err }

// NewDeliveryError wraps a target failure so diagnostics do not misclassify it as database delivery.
func NewDeliveryError(target, stage string, err error) error {
	return &DeliveryError{Target: target, Stage: stage, Err: err}
}

// DatabaseDeliveryStage describes one observed stage in the runtime database delivery path.
type DatabaseDeliveryStage string

const (
	// DatabaseDeliveryStageCollected means runtime received a collected point value.
	DatabaseDeliveryStageCollected DatabaseDeliveryStage = "collected"
	// DatabaseDeliveryStageMapped means the collected value reached a mapped tag.
	DatabaseDeliveryStageMapped DatabaseDeliveryStage = "mapped"
	// DatabaseDeliveryStageDBWrite means the value reached the database writer stage.
	DatabaseDeliveryStageDBWrite DatabaseDeliveryStage = "db_write"
	// DatabaseDeliveryStageDBWriteFailed means the database writer rejected the value.
	DatabaseDeliveryStageDBWriteFailed DatabaseDeliveryStage = "db_write_failed"
	// DatabaseDeliveryStageDBWriteSucceeded means the database writer accepted the value.
	DatabaseDeliveryStageDBWriteSucceeded DatabaseDeliveryStage = "db_write_succeeded"
)

// DatabaseDeliveryStatus summarizes the latest runtime database delivery result.
type DatabaseDeliveryStatus string

const (
	// DatabaseDeliveryStatusSucceeded means the latest database delivery reached the writer successfully.
	DatabaseDeliveryStatusSucceeded DatabaseDeliveryStatus = "succeeded"
	// DatabaseDeliveryStatusFailed means the latest database delivery failed after collection and mapping.
	DatabaseDeliveryStatusFailed DatabaseDeliveryStatus = "failed"
)

// DatabaseDeliveryDiagnostic is the latest per point/tag database delivery diagnostic.
type DatabaseDeliveryDiagnostic struct {
	DeviceID          string                  `json:"device_id,omitempty"`
	PointID           string                  `json:"point_id"`
	TagID             string                  `json:"tag_id"`
	Status            DatabaseDeliveryStatus  `json:"status"`
	Stages            []DatabaseDeliveryStage `json:"stages"`
	FailedStage       DatabaseDeliveryStage   `json:"failed_stage,omitempty"`
	Error             string                  `json:"error,omitempty"`
	ObservedAt        time.Time               `json:"observed_at"`
	LastSuccessAt     *time.Time              `json:"last_success_at,omitempty"`
	LastFailureAt     *time.Time              `json:"last_failure_at,omitempty"`
	LastFailureReason string                  `json:"last_failure_reason,omitempty"`
}

// RuntimeFlowStage describes one operator-facing runtime flow stage.
type RuntimeFlowStage string

const (
	// RuntimeFlowStageCollector means a point value reached runtime from collection.
	RuntimeFlowStageCollector RuntimeFlowStage = "collector"
	// RuntimeFlowStageMapping means the point value reached mapping.
	RuntimeFlowStageMapping RuntimeFlowStage = "mapping"
	// RuntimeFlowStageRuntimeProjection means runtime projected the mapped value.
	RuntimeFlowStageRuntimeProjection RuntimeFlowStage = "runtime_projection"
	// RuntimeFlowStageDatabaseDelivery means runtime attempted database delivery.
	RuntimeFlowStageDatabaseDelivery RuntimeFlowStage = "database_delivery"
)

// RuntimeFlowStageStatus describes the status of one runtime flow stage.
type RuntimeFlowStageStatus string

const (
	// RuntimeFlowStageStatusSuccess means the stage completed.
	RuntimeFlowStageStatusSuccess RuntimeFlowStageStatus = "success"
	// RuntimeFlowStageStatusFailed means the stage failed.
	RuntimeFlowStageStatusFailed RuntimeFlowStageStatus = "failed"
	// RuntimeFlowStageStatusUnknown means the stage has no current observation.
	RuntimeFlowStageStatusUnknown RuntimeFlowStageStatus = "unknown"
)

// RuntimeFlowDiagnosticStage is one stage status inside a runtime flow diagnostic.
type RuntimeFlowDiagnosticStage struct {
	Stage      RuntimeFlowStage       `json:"stage"`
	Status     RuntimeFlowStageStatus `json:"status"`
	ObservedAt *time.Time             `json:"observed_at,omitempty"`
	Reason     string                 `json:"reason,omitempty"`
}

// RuntimeFlowDiagnostic summarizes the latest operator-facing flow state for one scope.
type RuntimeFlowDiagnostic struct {
	Scope                 string                       `json:"scope"`
	DeviceID              string                       `json:"device_id,omitempty"`
	PointID               string                       `json:"point_id,omitempty"`
	TagID                 string                       `json:"tag_id,omitempty"`
	LastSuccessAt         *time.Time                   `json:"last_success_at,omitempty"`
	LastFailureAt         *time.Time                   `json:"last_failure_at,omitempty"`
	LatestSuccessfulStage RuntimeFlowStage             `json:"latest_successful_stage,omitempty"`
	FailureStage          RuntimeFlowStage             `json:"failure_stage,omitempty"`
	FailureReason         string                       `json:"failure_reason,omitempty"`
	Stages                []RuntimeFlowDiagnosticStage `json:"stages"`
}

func (s *Service) databaseDeliveryDiagnostics(deviceID string) []DatabaseDeliveryDiagnostic {
	if s == nil {
		return nil
	}

	s.databaseDeliveryMu.RLock()
	defer s.databaseDeliveryMu.RUnlock()

	filterDeviceID := strings.TrimSpace(deviceID)
	diagnostics := make([]DatabaseDeliveryDiagnostic, 0, len(s.databaseDelivery))
	for _, diagnostic := range s.databaseDelivery {
		if filterDeviceID != "" && diagnostic.DeviceID != filterDeviceID {
			continue
		}
		diagnostic.Stages = append([]DatabaseDeliveryStage(nil), diagnostic.Stages...)
		diagnostic.LastSuccessAt = cloneTimePtr(diagnostic.LastSuccessAt)
		diagnostic.LastFailureAt = cloneTimePtr(diagnostic.LastFailureAt)
		diagnostics = append(diagnostics, diagnostic)
	}

	sort.Slice(diagnostics, func(i, j int) bool {
		left := diagnostics[i]
		right := diagnostics[j]
		if !left.ObservedAt.Equal(right.ObservedAt) {
			return left.ObservedAt.Before(right.ObservedAt)
		}
		if left.PointID != right.PointID {
			return left.PointID < right.PointID
		}
		return left.TagID < right.TagID
	})
	return diagnostics
}

func (s *Service) recordDatabaseDeliveryDiagnostic(diagnostic DatabaseDeliveryDiagnostic) {
	if s == nil {
		return
	}

	diagnostic.DeviceID = strings.TrimSpace(diagnostic.DeviceID)
	diagnostic.PointID = strings.TrimSpace(diagnostic.PointID)
	diagnostic.TagID = strings.TrimSpace(diagnostic.TagID)
	diagnostic.Error = strings.TrimSpace(diagnostic.Error)
	if diagnostic.PointID == "" || diagnostic.TagID == "" {
		return
	}
	if diagnostic.ObservedAt.IsZero() {
		diagnostic.ObservedAt = time.Now()
	}
	diagnostic.ObservedAt = diagnostic.ObservedAt.UTC()
	diagnostic.Stages = append([]DatabaseDeliveryStage(nil), diagnostic.Stages...)

	s.databaseDeliveryMu.Lock()
	defer s.databaseDeliveryMu.Unlock()
	if s.databaseDelivery == nil {
		s.databaseDelivery = make(map[string]DatabaseDeliveryDiagnostic)
	}
	key := databaseDeliveryKey(diagnostic.PointID, diagnostic.TagID)
	s.databaseDelivery[key] = mergeDatabaseDeliveryContext(diagnostic, s.databaseDelivery[key])
}

func databaseDeliveryKey(pointID, tagID string) string {
	return strings.TrimSpace(pointID) + "\x00" + strings.TrimSpace(tagID)
}

func mergeDatabaseDeliveryContext(next DatabaseDeliveryDiagnostic, previous DatabaseDeliveryDiagnostic) DatabaseDeliveryDiagnostic {
	next.LastSuccessAt = cloneTimePtr(previous.LastSuccessAt)
	next.LastFailureAt = cloneTimePtr(previous.LastFailureAt)
	next.LastFailureReason = previous.LastFailureReason

	if next.LastSuccessAt == nil && previous.Status == DatabaseDeliveryStatusSucceeded {
		next.LastSuccessAt = cloneObservedAt(previous.ObservedAt)
	}
	if next.LastFailureAt == nil && previous.Status == DatabaseDeliveryStatusFailed {
		next.LastFailureAt = cloneObservedAt(previous.ObservedAt)
		next.LastFailureReason = firstNonEmpty(previous.LastFailureReason, previous.Error)
	}

	switch next.Status {
	case DatabaseDeliveryStatusSucceeded:
		next.LastSuccessAt = cloneObservedAt(next.ObservedAt)
	case DatabaseDeliveryStatusFailed:
		next.LastFailureAt = cloneObservedAt(next.ObservedAt)
		next.LastFailureReason = firstNonEmpty(next.Error, next.LastFailureReason)
	}
	return next
}

func (s *Service) runtimeFlowDiagnostics(deviceID string) []RuntimeFlowDiagnostic {
	deliveryDiagnostics := s.databaseDeliveryDiagnostics(deviceID)
	if len(deliveryDiagnostics) == 0 {
		return nil
	}

	diagnostics := make([]RuntimeFlowDiagnostic, 0, len(deliveryDiagnostics))
	for _, diagnostic := range deliveryDiagnostics {
		diagnostics = append(diagnostics, runtimeFlowDiagnosticFromDatabaseDelivery(diagnostic))
	}
	return diagnostics
}

func runtimeFlowDiagnosticFromDatabaseDelivery(diagnostic DatabaseDeliveryDiagnostic) RuntimeFlowDiagnostic {
	failureReason := firstNonEmpty(diagnostic.LastFailureReason, diagnostic.Error)
	flow := RuntimeFlowDiagnostic{
		Scope:                 runtimeFlowScope(diagnostic),
		DeviceID:              diagnostic.DeviceID,
		PointID:               diagnostic.PointID,
		TagID:                 diagnostic.TagID,
		LastSuccessAt:         cloneTimePtr(diagnostic.LastSuccessAt),
		LastFailureAt:         cloneTimePtr(diagnostic.LastFailureAt),
		LatestSuccessfulStage: RuntimeFlowStageRuntimeProjection,
		FailureReason:         failureReason,
		Stages: []RuntimeFlowDiagnosticStage{
			successRuntimeFlowStage(RuntimeFlowStageCollector, diagnostic.ObservedAt),
			successRuntimeFlowStage(RuntimeFlowStageMapping, diagnostic.ObservedAt),
			successRuntimeFlowStage(RuntimeFlowStageRuntimeProjection, diagnostic.ObservedAt),
		},
	}

	if flow.LastSuccessAt == nil && diagnostic.Status == DatabaseDeliveryStatusSucceeded {
		flow.LastSuccessAt = cloneObservedAt(diagnostic.ObservedAt)
	}
	if flow.LastFailureAt == nil && diagnostic.Status == DatabaseDeliveryStatusFailed {
		flow.LastFailureAt = cloneObservedAt(diagnostic.ObservedAt)
	}

	switch diagnostic.Status {
	case DatabaseDeliveryStatusSucceeded:
		flow.LatestSuccessfulStage = RuntimeFlowStageDatabaseDelivery
		flow.FailureReason = failureReason
		flow.Stages = append(flow.Stages, successRuntimeFlowStage(RuntimeFlowStageDatabaseDelivery, diagnostic.ObservedAt))
	case DatabaseDeliveryStatusFailed:
		flow.FailureStage = RuntimeFlowStageDatabaseDelivery
		flow.Stages = append(flow.Stages, RuntimeFlowDiagnosticStage{
			Stage:      RuntimeFlowStageDatabaseDelivery,
			Status:     RuntimeFlowStageStatusFailed,
			ObservedAt: cloneObservedAt(diagnostic.ObservedAt),
			Reason:     failureReason,
		})
	default:
		flow.LatestSuccessfulStage = ""
		flow.FailureReason = ""
		flow.Stages = append(flow.Stages, RuntimeFlowDiagnosticStage{
			Stage:  RuntimeFlowStageDatabaseDelivery,
			Status: RuntimeFlowStageStatusUnknown,
		})
	}
	return flow
}

func successRuntimeFlowStage(stage RuntimeFlowStage, observedAt time.Time) RuntimeFlowDiagnosticStage {
	return RuntimeFlowDiagnosticStage{
		Stage:      stage,
		Status:     RuntimeFlowStageStatusSuccess,
		ObservedAt: cloneObservedAt(observedAt),
	}
}

func runtimeFlowScope(diagnostic DatabaseDeliveryDiagnostic) string {
	if diagnostic.DeviceID != "" {
		return "device:" + diagnostic.DeviceID
	}
	if diagnostic.PointID != "" {
		return "point:" + diagnostic.PointID
	}
	if diagnostic.TagID != "" {
		return "tag:" + diagnostic.TagID
	}
	return "runtime"
}

func cloneObservedAt(observedAt time.Time) *time.Time {
	if observedAt.IsZero() {
		return nil
	}
	t := observedAt.UTC()
	return &t
}

func cloneTimePtr(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	t := value.UTC()
	return &t
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
