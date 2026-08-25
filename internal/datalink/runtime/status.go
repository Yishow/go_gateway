package runtime

import (
	"context"
	"time"

	"go-gateway/internal/datalink/collector/health"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
)

const statusRefreshInterval = time.Second

// RuntimeStatusSnapshot describes the runtime monitoring snapshot payload.
type RuntimeStatusSnapshot struct {
	Running             bool                            `json:"running"`
	UptimeSeconds       int64                           `json:"uptime_seconds"`
	SnapshotState       RuntimeTruthState               `json:"snapshot_state"`
	Metrics             Stats                           `json:"metrics"`
	Collectors          []DeviceRuntimeStatus           `json:"collectors"`
	DatabaseDelivery    []DatabaseDeliveryDiagnostic    `json:"database_delivery"`
	ModbusShareDelivery []ModbusShareDeliveryDiagnostic `json:"modbus_share_delivery,omitempty"`
	Diagnostics         []RuntimeFlowDiagnostic         `json:"diagnostics"`
}

// DeviceRuntimeStatus describes one device runtime summary shared by snapshot and stream.
type DeviceRuntimeStatus struct {
	DeviceID                   string     `json:"device_id"`
	DeviceName                 string     `json:"device_name,omitempty"`
	Protocol                   string     `json:"protocol,omitempty"`
	Status                     string     `json:"status"`
	AvailabilityStatus         string     `json:"availability_status"`
	AvailabilityReason         *string    `json:"availability_reason,omitempty"`
	Running                    bool       `json:"running"`
	PointsTotal                int        `json:"points_total"`
	PointsHealthy              int        `json:"points_healthy"`
	PointsStale                int        `json:"points_stale"`
	PointsError                int        `json:"points_error"`
	LastReadAt                 *time.Time `json:"last_read_at"`
	LastError                  *string    `json:"last_error"`
	BreakerState               string     `json:"breaker_state"`
	ProjectionAlignment        string     `json:"projection_alignment,omitempty"`
	RuntimeProjectionVersion   string     `json:"runtime_projection_version,omitempty"`
	WorkspaceProjectionVersion string     `json:"workspace_projection_version,omitempty"`
	ProjectionMessage          string     `json:"projection_message,omitempty"`
	ProjectionCode             string     `json:"projection_code,omitempty"`
}

// RuntimeStatusSnapshot returns the derived runtime monitoring snapshot.
func (s *Service) RuntimeStatusSnapshot(ctx context.Context, deviceID string) (RuntimeStatusSnapshot, error) {
	devices, err := s.listRuntimeStatusDevices(ctx, deviceID)
	if err != nil {
		return RuntimeStatusSnapshot{}, err
	}

	points, err := s.listRuntimeStatusPoints(ctx, deviceID)
	if err != nil {
		return RuntimeStatusSnapshot{}, err
	}

	groups, err := s.listRuntimeStatusGroups(ctx)
	if err != nil {
		return RuntimeStatusSnapshot{}, err
	}

	collectors := s.deriveDeviceStatuses(devices, points, groups, deviceID, time.Now())
	s.attachProjectionStates(ctx, collectors, devices)

	return RuntimeStatusSnapshot{
		Running:             s.IsRunning(),
		UptimeSeconds:       s.UptimeSeconds(),
		SnapshotState:       DeriveSnapshotTruthState(collectors, deviceID),
		Metrics:             s.Snapshot(),
		Collectors:          collectors,
		DatabaseDelivery:    s.databaseDeliveryDiagnostics(deviceID),
		ModbusShareDelivery: s.modbusShareDeliveryDiagnostics(deviceID),
		Diagnostics:         s.runtimeFlowDiagnostics(deviceID),
	}, nil
}

// DeviceStatus returns one derived device runtime summary.
func (s *Service) DeviceStatus(ctx context.Context, deviceID string) (DeviceRuntimeStatus, bool, error) {
	snapshot, err := s.RuntimeStatusSnapshot(ctx, deviceID)
	if err != nil {
		return DeviceRuntimeStatus{}, false, err
	}
	if len(snapshot.Collectors) == 0 {
		return DeviceRuntimeStatus{}, false, nil
	}
	return snapshot.Collectors[0], true, nil
}

func (s *Service) listRuntimeStatusDevices(ctx context.Context, deviceID string) ([]*schema.Device, error) {
	if s.deviceSvc != nil {
		if deviceID != "" {
			deviceRecord, err := s.deviceSvc.GetByID(ctx, deviceID)
			if err != nil {
				return nil, err
			}
			return []*schema.Device{deviceRecord}, nil
		}
		return s.deviceSvc.List(ctx, device.ListFilter{Limit: 100000})
	}

	if deviceID == "" {
		return append([]*schema.Device(nil), s.snapshot.Devices...), nil
	}

	for _, deviceRecord := range s.snapshot.Devices {
		if deviceRecord.ID == deviceID {
			return []*schema.Device{deviceRecord}, nil
		}
	}
	return nil, nil
}

func (s *Service) listRuntimeStatusPoints(ctx context.Context, deviceID string) ([]*schema.Point, error) {
	if s.pointSvc != nil {
		filter := point.ListFilter{Limit: 100000}
		if deviceID != "" {
			filter.DeviceID = &deviceID
		}
		return s.pointSvc.List(ctx, filter)
	}

	if deviceID == "" {
		return append([]*schema.Point(nil), s.snapshot.Points...), nil
	}

	points := make([]*schema.Point, 0)
	for _, pointRecord := range s.snapshot.Points {
		if pointRecord.DeviceID == deviceID {
			points = append(points, pointRecord)
		}
	}
	return points, nil
}

func (s *Service) listRuntimeStatusGroups(ctx context.Context) ([]*schema.PollingGroup, error) {
	if s.groupSvc != nil {
		return s.groupSvc.List(ctx)
	}
	return append([]*schema.PollingGroup(nil), s.snapshot.PollingGroups...), nil
}

func (s *Service) deriveDeviceStatuses(
	devices []*schema.Device,
	points []*schema.Point,
	groups []*schema.PollingGroup,
	deviceID string,
	now time.Time,
) []DeviceRuntimeStatus {
	groupIntervalByID := make(map[string]time.Duration, len(groups))
	for _, group := range groups {
		groupIntervalByID[group.ID] = time.Duration(group.IntervalMs) * time.Millisecond
	}

	pointsByDeviceID := make(map[string][]*schema.Point, len(devices))
	for _, pt := range points {
		pointsByDeviceID[pt.DeviceID] = append(pointsByDeviceID[pt.DeviceID], pt)
	}

	collectors := make([]DeviceRuntimeStatus, 0, len(devices))
	for _, dev := range devices {
		if deviceID != "" && dev.ID != deviceID {
			continue
		}

		status := DeviceRuntimeStatus{
			DeviceID:           dev.ID,
			DeviceName:         dev.Name,
			Protocol:           string(dev.Protocol),
			Status:             "idle",
			AvailabilityStatus: device.AvailabilityStatusAvailable,
			BreakerState:       s.runtimeBreakerState(dev.ID),
		}
		if availabilityStatus, availabilityReason := device.AvailabilityOf(dev); availabilityStatus == device.AvailabilityStatusUnavailable {
			status.AvailabilityStatus = availabilityStatus
			if availabilityReason != "" {
				status.AvailabilityReason = &availabilityReason
			}
		}

		var lastReadAt *time.Time
		var lastError *string
		var lastErrorAt *time.Time

		for _, pt := range pointsByDeviceID[dev.ID] {
			status.PointsTotal++

			if pt.LastReadAt == nil {
				status.PointsStale++
				continue
			}

			if lastReadAt == nil || pt.LastReadAt.After(*lastReadAt) {
				lastReadAt = pt.LastReadAt
			}

			if pt.LastError != "" {
				status.PointsError++
				if lastErrorAt == nil || pt.LastReadAt.After(*lastErrorAt) {
					lastErrorAt = pt.LastReadAt
					errText := pt.LastError
					lastError = &errText
				}
				continue
			}

			if now.Sub(*pt.LastReadAt) > effectiveRuntimePointInterval(groupIntervalByID, pt) {
				status.PointsStale++
				continue
			}

			status.PointsHealthy++
		}

		status.LastReadAt = lastReadAt
		status.LastError = lastError
		status.Status = deriveRuntimeDeviceState(status, s.IsRunning())
		status.Running = status.Status == "running" && status.AvailabilityStatus == device.AvailabilityStatusAvailable
		if status.AvailabilityStatus == device.AvailabilityStatusUnavailable {
			status.Running = false
			status.Status = "idle"
		}
		collectors = append(collectors, status)
	}

	return collectors
}

func effectiveRuntimePointInterval(groupIntervalByID map[string]time.Duration, pt *schema.Point) time.Duration {
	interval := time.Second
	if pt == nil || pt.PollingGroupID == nil {
		return interval
	}

	groupInterval, exists := groupIntervalByID[*pt.PollingGroupID]
	if !exists || groupInterval <= 0 {
		return interval
	}
	return groupInterval
}

func deriveRuntimeDeviceState(status DeviceRuntimeStatus, running bool) string {
	switch {
	case status.PointsError > 0 || status.BreakerState == "open":
		return "error"
	case status.PointsStale > 0 || status.BreakerState == "half-open":
		return "warning"
	case running:
		return "running"
	default:
		return "idle"
	}
}

func (s *Service) attachProjectionStates(ctx context.Context, collectors []DeviceRuntimeStatus, devices []*schema.Device) {
	states := s.projectionStatesForDevices(ctx, devices)
	if len(states) == 0 {
		return
	}

	for idx := range collectors {
		state, ok := states[collectors[idx].DeviceID]
		if !ok || state.Alignment == "" {
			continue
		}
		collectors[idx].ProjectionAlignment = string(state.Alignment)
		collectors[idx].RuntimeProjectionVersion = state.RuntimeVersion
		collectors[idx].WorkspaceProjectionVersion = state.WorkspaceVersion
		collectors[idx].ProjectionMessage = state.Message
		collectors[idx].ProjectionCode = state.Code
	}
}

func (s *Service) runtimeBreakerState(deviceID string) string {
	if s.scheduler == nil {
		return "closed"
	}

	state, exists := s.scheduler.GetDeviceBreakerState(deviceID)
	if !exists {
		return "closed"
	}

	switch state {
	case health.StateDead:
		return "open"
	case health.StateUnstable:
		return "half-open"
	default:
		return "closed"
	}
}

func (s *Service) statusLoop() {
	defer s.wg.Done()

	ticker := time.NewTicker(statusRefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.refreshDerivedStatuses()
		}
	}
}

func (s *Service) refreshDerivedStatuses() {
	s.statusSubscriberMu.RLock()
	hasSubscribers := len(s.statusSubscribers) > 0
	s.statusSubscriberMu.RUnlock()
	if !hasSubscribers {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	snapshot, err := s.RuntimeStatusSnapshot(ctx, "")
	if err != nil {
		return
	}

	for _, collector := range snapshot.Collectors {
		s.emitStatusIfChanged(collector)
	}
}

func (s *Service) emitStatusIfChanged(status DeviceRuntimeStatus) {
	if s.recordLatestStatus(status) {
		s.broadcastStatusEvent(status)
	}
}

func (s *Service) recordLatestStatus(status DeviceRuntimeStatus) bool {
	s.lastStatusMu.Lock()
	defer s.lastStatusMu.Unlock()
	if s.lastStatuses == nil {
		s.lastStatuses = make(map[string]DeviceRuntimeStatus)
	}

	previous, exists := s.lastStatuses[status.DeviceID]
	if exists && runtimeStatusEqual(previous, status) {
		return false
	}

	s.lastStatuses[status.DeviceID] = status
	return true
}

func runtimeStatusEqual(left, right DeviceRuntimeStatus) bool {
	if left.DeviceID != right.DeviceID ||
		left.DeviceName != right.DeviceName ||
		left.Protocol != right.Protocol ||
		left.Status != right.Status ||
		left.PointsTotal != right.PointsTotal ||
		left.PointsHealthy != right.PointsHealthy ||
		left.PointsStale != right.PointsStale ||
		left.PointsError != right.PointsError ||
		left.BreakerState != right.BreakerState ||
		left.ProjectionAlignment != right.ProjectionAlignment ||
		left.RuntimeProjectionVersion != right.RuntimeProjectionVersion ||
		left.WorkspaceProjectionVersion != right.WorkspaceProjectionVersion ||
		left.ProjectionMessage != right.ProjectionMessage ||
		left.ProjectionCode != right.ProjectionCode {
		return false
	}

	if !runtimeTimeEqual(left.LastReadAt, right.LastReadAt) {
		return false
	}

	return runtimeStringPtrEqual(left.LastError, right.LastError)
}

func runtimeTimeEqual(left, right *time.Time) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return left.Equal(*right)
}

func runtimeStringPtrEqual(left, right *string) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}
