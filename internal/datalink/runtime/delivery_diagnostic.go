package runtime

import (
	"sort"
	"strings"
	"time"
)

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
	DeviceID    string                  `json:"device_id,omitempty"`
	PointID     string                  `json:"point_id"`
	TagID       string                  `json:"tag_id"`
	Status      DatabaseDeliveryStatus  `json:"status"`
	Stages      []DatabaseDeliveryStage `json:"stages"`
	FailedStage DatabaseDeliveryStage   `json:"failed_stage,omitempty"`
	Error       string                  `json:"error,omitempty"`
	ObservedAt  time.Time               `json:"observed_at"`
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
	s.databaseDelivery[databaseDeliveryKey(diagnostic.PointID, diagnostic.TagID)] = diagnostic
}

func databaseDeliveryKey(pointID, tagID string) string {
	return strings.TrimSpace(pointID) + "\x00" + strings.TrimSpace(tagID)
}
