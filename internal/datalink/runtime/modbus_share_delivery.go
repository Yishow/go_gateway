package runtime

import (
	"sort"
	"time"
)

// ModbusShareDeliveryDiagnostic is the latest per point/tag Share delivery result.
type ModbusShareDeliveryDiagnostic struct {
	DeviceID   string                    `json:"device_id,omitempty"`
	PointID    string                    `json:"point_id"`
	TagID      string                    `json:"tag_id"`
	Status     ModbusShareDeliveryStatus `json:"status"`
	Stage      string                    `json:"stage"`
	Error      string                    `json:"error,omitempty"`
	ObservedAt time.Time                 `json:"observed_at"`
}

// ModbusShareDeliveryStatus describes the latest local Share target result.
type ModbusShareDeliveryStatus string

const (
	// ModbusShareDeliveryStatusSucceeded means the Share target accepted the value.
	ModbusShareDeliveryStatusSucceeded ModbusShareDeliveryStatus = "succeeded"
	// ModbusShareDeliveryStatusFailed means the Share target rejected the value.
	ModbusShareDeliveryStatusFailed ModbusShareDeliveryStatus = "failed"
)

func (s *Service) recordModbusShareDeliveryDiagnostic(d ModbusShareDeliveryDiagnostic) {
	if s == nil || d.PointID == "" || d.TagID == "" {
		return
	}
	if d.ObservedAt.IsZero() {
		d.ObservedAt = time.Now().UTC()
	}
	d.ObservedAt = d.ObservedAt.UTC()
	s.databaseDeliveryMu.Lock()
	defer s.databaseDeliveryMu.Unlock()
	if s.modbusShareDelivery == nil {
		s.modbusShareDelivery = make(map[string]ModbusShareDeliveryDiagnostic)
	}
	s.modbusShareDelivery[databaseDeliveryKey(d.PointID, d.TagID)] = d
}

func (s *Service) modbusShareDeliveryDiagnostics(deviceID string) []ModbusShareDeliveryDiagnostic {
	if s == nil {
		return nil
	}
	s.databaseDeliveryMu.RLock()
	defer s.databaseDeliveryMu.RUnlock()
	out := make([]ModbusShareDeliveryDiagnostic, 0, len(s.modbusShareDelivery))
	for _, d := range s.modbusShareDelivery {
		if deviceID == "" || d.DeviceID == deviceID {
			out = append(out, d)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ObservedAt.Before(out[j].ObservedAt) })
	return out
}
