package handlers

import "time"

// AcquisitionAnomalyItem is one anomaly event. point_* fields are null for
// device-level events or points that no longer exist; duration_ms is set only
// once the event has ended.
type AcquisitionAnomalyItem struct {
	ID           int64      `json:"id"`
	DeviceID     string     `json:"device_id"`
	PointID      *string    `json:"point_id"`
	PointName    *string    `json:"point_name"`
	PointAddress *string    `json:"point_address"`
	Kind         string     `json:"kind" enums:"disconnected,stale,frozen"`
	State        string     `json:"state" enums:"active,ended"`
	StartedAt    time.Time  `json:"started_at"`
	EndedAt      *time.Time `json:"ended_at"`
	DurationMs   *int64     `json:"duration_ms"`
	LastValue    *string    `json:"last_value"`
	Detail       string     `json:"detail"`
}

// AcquisitionAnomalyListData is one page of a device's events with the
// filtered total.
type AcquisitionAnomalyListData struct {
	Items  []AcquisitionAnomalyItem `json:"items"`
	Total  int                      `json:"total"`
	Limit  int                      `json:"limit"`
	Offset int                      `json:"offset"`
}

// AcquisitionAnomalyListResponse wraps a device's anomaly page.
type AcquisitionAnomalyListResponse struct {
	Success bool                       `json:"success" example:"true"`
	Data    AcquisitionAnomalyListData `json:"data"`
}

// AcquisitionAnomalyDeviceCount is one device's number of active events.
type AcquisitionAnomalyDeviceCount struct {
	DeviceID    string `json:"device_id"`
	ActiveCount int    `json:"active_count"`
}

// AcquisitionAnomalySummaryData lists every device's active event count,
// including explicit zeros.
type AcquisitionAnomalySummaryData struct {
	Items       []AcquisitionAnomalyDeviceCount `json:"items"`
	TotalActive int                             `json:"total_active"`
}

// AcquisitionAnomalySummaryResponse wraps the per-device active counts.
type AcquisitionAnomalySummaryResponse struct {
	Success bool                          `json:"success" example:"true"`
	Data    AcquisitionAnomalySummaryData `json:"data"`
}

// AcquisitionAnomalyErrorEnvelope is a typed error naming the rejected field.
type AcquisitionAnomalyErrorEnvelope struct {
	TypedAPIErrorEnvelope
	Field string `json:"field,omitempty"`
}

// AcquisitionAnomalyErrorResponse is the failure shape of the anomaly routes.
type AcquisitionAnomalyErrorResponse struct {
	Success bool                            `json:"success" example:"false"`
	Error   AcquisitionAnomalyErrorEnvelope `json:"error"`
}
