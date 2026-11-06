package handlers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"sort"
	"strconv"

	"go-gateway/internal/datalink/anomaly"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"

	"github.com/gin-gonic/gin"
)

// Acquisition anomaly API error codes.
const (
	ErrCodeAcquisitionAnomalyInvalidFilter         = "acquisition_anomaly_filter_invalid"
	ErrCodeAcquisitionAnomalyDeviceNotFound        = "acquisition_anomaly_device_not_found"
	ErrCodeAcquisitionAnomalyUnavailable           = "acquisition_anomaly_unavailable"
	ErrCodeAcquisitionAnomalyCapabilityUnavailable = "acquisition_anomaly_capability_unavailable"
)

// Page bounds of the list route (design condition); out-of-range values are
// rejected, never clamped.
const (
	defaultAcquisitionAnomalyLimit = 50
	maxAcquisitionAnomalyLimit     = 500
)

const (
	acquisitionAnomalyInvalidFilterMessage         = "acquisition anomaly filter is invalid"
	acquisitionAnomalyDeviceNotFoundMessage        = "device was not found"
	acquisitionAnomalyUnavailableMessage           = "acquisition anomalies could not be read"
	acquisitionAnomalyCapabilityUnavailableMessage = "acquisition anomaly recording is not available"
	acquisitionAnomalyRetryAction                  = "Retry the request"
	acquisitionAnomalyReviewAction                 = "Review the request and retry"
	acquisitionAnomalyCapabilityUnavailableAction  = "Restart the gateway with anomaly recording enabled"
	acquisitionAnomalyDeviceNotFoundAction         = "Refresh the device list"
	acquisitionAnomalyFieldState                   = "state"
	acquisitionAnomalyFieldKind                    = "kind"
	acquisitionAnomalyFieldLimit                   = "limit"
	acquisitionAnomalyFieldOffset                  = "offset"
	acquisitionAnomalyStateActive                  = string(anomaly.StateActive)
	acquisitionAnomalyStateEnded                   = string(anomaly.StateEnded)
)

// AcquisitionAnomalyReader reads persisted acquisition anomaly events.
type AcquisitionAnomalyReader interface {
	List(ctx context.Context, filter anomaly.ListFilter) ([]anomaly.Event, int, error)
	ActiveCounts(ctx context.Context) (map[string]int, error)
}

// AcquisitionAnomalyDeviceReader resolves devices for anomaly queries.
type AcquisitionAnomalyDeviceReader interface {
	GetByID(ctx context.Context, id string) (*schema.Device, error)
	List(ctx context.Context, filter device.ListFilter) ([]*schema.Device, error)
}

// AcquisitionAnomalyPointReader resolves point names and addresses for events.
type AcquisitionAnomalyPointReader interface {
	ListByDevice(ctx context.Context, deviceID string) ([]*schema.Point, error)
}

// AcquisitionAnomalyHandler serves the acquisition anomaly history and summary.
type AcquisitionAnomalyHandler struct {
	reader  AcquisitionAnomalyReader
	devices AcquisitionAnomalyDeviceReader
	points  AcquisitionAnomalyPointReader
}

// NewAcquisitionAnomalyHandler creates the handler. A nil reader makes every
// route answer that the capability is unavailable instead of an empty list.
func NewAcquisitionAnomalyHandler(reader AcquisitionAnomalyReader, devices AcquisitionAnomalyDeviceReader, points AcquisitionAnomalyPointReader) *AcquisitionAnomalyHandler {
	return &AcquisitionAnomalyHandler{reader: reader, devices: devices, points: points}
}

// ListDeviceAnomalies returns one device's events, newest first.
// @Summary List a device's acquisition anomalies
// @Description Lists persisted disconnected, stale and frozen events of one device, newest first, with the filtered total. Each device targets 5000 events by pruning the oldest ended events first. Active events are always retained, even above the target; ending an event triggers pruning again.
// @Tags Datalink Acquisition Anomalies
// @Produce json
// @Param id path string true "Device ID"
// @Param state query string false "Event state" Enums(active, ended)
// @Param kind query string false "Anomaly kind" Enums(disconnected, stale, frozen)
// @Param limit query int false "Page size (1-500, default 50)" minimum(1) maximum(500)
// @Param offset query int false "Number of events to skip (default 0)" minimum(0)
// @Success 200 {object} AcquisitionAnomalyListResponse
// @Failure 400,404,500,503 {object} AcquisitionAnomalyErrorResponse
// @Router /v1/datalink/devices/{id}/acquisition-anomalies [get]
func (h *AcquisitionAnomalyHandler) ListDeviceAnomalies(c *gin.Context) {
	if h.reader == nil || h.devices == nil {
		renderAcquisitionAnomalyError(c, http.StatusServiceUnavailable, ErrCodeAcquisitionAnomalyCapabilityUnavailable, "")
		return
	}
	filter, field := parseAcquisitionAnomalyFilter(c)
	if field != "" {
		renderAcquisitionAnomalyError(c, http.StatusBadRequest, ErrCodeAcquisitionAnomalyInvalidFilter, field)
		return
	}
	ctx := c.Request.Context()
	if _, err := h.devices.GetByID(ctx, filter.DeviceID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			renderAcquisitionAnomalyError(c, http.StatusNotFound, ErrCodeAcquisitionAnomalyDeviceNotFound, "")
		} else {
			renderAcquisitionAnomalyUnavailable(c, "read device", err)
		}
		return
	}
	events, total, err := h.reader.List(ctx, filter)
	if err != nil {
		if errors.Is(err, anomaly.ErrInvalidFilter) {
			renderAcquisitionAnomalyError(c, http.StatusBadRequest, ErrCodeAcquisitionAnomalyInvalidFilter, "")
		} else {
			renderAcquisitionAnomalyUnavailable(c, "list events", err)
		}
		return
	}
	points, err := h.pointsByID(ctx, filter.DeviceID, events)
	if err != nil {
		renderAcquisitionAnomalyUnavailable(c, "read points", err)
		return
	}
	items := make([]AcquisitionAnomalyItem, 0, len(events))
	for _, ev := range events {
		items = append(items, newAcquisitionAnomalyItem(ev, points[ev.PointID]))
	}
	c.JSON(http.StatusOK, AcquisitionAnomalyListResponse{Success: true, Data: AcquisitionAnomalyListData{Items: items, Total: total, Limit: filter.Limit, Offset: filter.Offset}})
}

// Summary returns the number of active events per device.
// @Summary Summarize active acquisition anomalies per device
// @Description Returns every device's number of active (not yet ended) anomaly events; devices without active events report 0.
// @Tags Datalink Acquisition Anomalies
// @Produce json
// @Success 200 {object} AcquisitionAnomalySummaryResponse
// @Failure 500,503 {object} AcquisitionAnomalyErrorResponse
// @Router /v1/datalink/acquisition-anomalies/summary [get]
func (h *AcquisitionAnomalyHandler) Summary(c *gin.Context) {
	if h.reader == nil || h.devices == nil {
		renderAcquisitionAnomalyError(c, http.StatusServiceUnavailable, ErrCodeAcquisitionAnomalyCapabilityUnavailable, "")
		return
	}
	ctx := c.Request.Context()
	counts, err := h.reader.ActiveCounts(ctx)
	if err != nil {
		renderAcquisitionAnomalyUnavailable(c, "count active events", err)
		return
	}
	devices, err := h.devices.List(ctx, device.ListFilter{})
	if err != nil {
		renderAcquisitionAnomalyUnavailable(c, "list devices", err)
		return
	}
	merged := make(map[string]int, len(devices)+len(counts))
	for _, dev := range devices {
		if dev != nil {
			merged[dev.ID] = 0
		}
	}
	data := AcquisitionAnomalySummaryData{Items: make([]AcquisitionAnomalyDeviceCount, 0, len(merged))}
	for id, n := range counts {
		merged[id] = n
		data.TotalActive += n
	}
	for id, n := range merged {
		data.Items = append(data.Items, AcquisitionAnomalyDeviceCount{DeviceID: id, ActiveCount: n})
	}
	sort.Slice(data.Items, func(i, j int) bool { return data.Items[i].DeviceID < data.Items[j].DeviceID })
	c.JSON(http.StatusOK, AcquisitionAnomalySummaryResponse{Success: true, Data: data})
}

// parseAcquisitionAnomalyFilter validates the query strictly; unknown values
// are rejected with the field name rather than replaced by defaults.
func parseAcquisitionAnomalyFilter(c *gin.Context) (filter anomaly.ListFilter, invalidField string) {
	filter = anomaly.ListFilter{DeviceID: c.Param("id"), Limit: defaultAcquisitionAnomalyLimit}
	if state, ok := c.GetQuery(acquisitionAnomalyFieldState); ok {
		if state != acquisitionAnomalyStateActive && state != acquisitionAnomalyStateEnded {
			return filter, acquisitionAnomalyFieldState
		}
		filter.State = anomaly.State(state)
	}
	if kind, ok := c.GetQuery(acquisitionAnomalyFieldKind); ok {
		if !anomaly.Kind(kind).Valid() {
			return filter, acquisitionAnomalyFieldKind
		}
		filter.Kind = anomaly.Kind(kind)
	}
	if raw, ok := c.GetQuery(acquisitionAnomalyFieldLimit); ok {
		limit, err := strconv.Atoi(raw)
		if err != nil || limit < 1 || limit > maxAcquisitionAnomalyLimit {
			return filter, acquisitionAnomalyFieldLimit
		}
		filter.Limit = limit
	}
	if raw, ok := c.GetQuery(acquisitionAnomalyFieldOffset); ok {
		offset, err := strconv.Atoi(raw)
		if err != nil || offset < 0 {
			return filter, acquisitionAnomalyFieldOffset
		}
		filter.Offset = offset
	}
	return filter, ""
}

func (h *AcquisitionAnomalyHandler) pointsByID(ctx context.Context, deviceID string, events []anomaly.Event) (map[string]*schema.Point, error) {
	needed := false
	for _, ev := range events {
		needed = needed || ev.PointID != ""
	}
	if !needed || h.points == nil {
		return nil, nil
	}
	points, err := h.points.ListByDevice(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	byID := make(map[string]*schema.Point, len(points))
	for _, p := range points {
		if p != nil {
			byID[p.ID] = p
		}
	}
	return byID, nil
}

func newAcquisitionAnomalyItem(ev anomaly.Event, p *schema.Point) AcquisitionAnomalyItem {
	item := AcquisitionAnomalyItem{
		ID: ev.ID, DeviceID: ev.DeviceID, Kind: string(ev.Kind), State: acquisitionAnomalyStateActive,
		StartedAt: ev.StartedAt, EndedAt: ev.EndedAt, LastValue: ev.LastValue, Detail: ev.Detail,
	}
	if ev.PointID != "" {
		pointID := ev.PointID
		item.PointID = &pointID
	}
	if p != nil {
		name, address := p.Name, p.Address
		item.PointName, item.PointAddress = &name, &address
	}
	if ev.EndedAt != nil {
		item.State = acquisitionAnomalyStateEnded
		duration := ev.EndedAt.Sub(ev.StartedAt).Milliseconds()
		item.DurationMs = &duration
	}
	return item
}

// renderAcquisitionAnomalyUnavailable logs the storage failure once on the
// server and answers with a retryable typed error that never carries raw
// backend text.
func renderAcquisitionAnomalyUnavailable(c *gin.Context, op string, err error) {
	log.Printf("acquisition anomaly API: %s: %v", op, err)
	renderAcquisitionAnomalyError(c, http.StatusInternalServerError, ErrCodeAcquisitionAnomalyUnavailable, "")
}

func renderAcquisitionAnomalyError(c *gin.Context, status int, code, field string) {
	message, action, retryable := acquisitionAnomalyInvalidFilterMessage, acquisitionAnomalyReviewAction, false
	switch code {
	case ErrCodeAcquisitionAnomalyDeviceNotFound:
		message, action = acquisitionAnomalyDeviceNotFoundMessage, acquisitionAnomalyDeviceNotFoundAction
	case ErrCodeAcquisitionAnomalyUnavailable:
		message, action, retryable = acquisitionAnomalyUnavailableMessage, acquisitionAnomalyRetryAction, true
	case ErrCodeAcquisitionAnomalyCapabilityUnavailable:
		message, action = acquisitionAnomalyCapabilityUnavailableMessage, acquisitionAnomalyCapabilityUnavailableAction
	}
	c.JSON(status, AcquisitionAnomalyErrorResponse{Error: AcquisitionAnomalyErrorEnvelope{
		TypedAPIErrorEnvelope: TypedAPIErrorEnvelope{Code: code, Message: message, Retryable: retryable, RequestID: getOrGenerateRequestID(c), Action: action},
		Field:                 field,
	}})
}
