package collector

import (
	"math"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

// PointObservationBinding is immutable execution metadata for one point install.
// RuleRevision is the original source-rule revision, never the fingerprint.
// This ephemeral value is not permission, configuration or delivery authority.
type PointObservationBinding struct {
	WorkspaceID, WorkspaceRevision, DeviceID, PointID string
	RuleID, RuleRevision, ConfigFingerprint           string
}

// AcquisitionObservation contains only final normalized acquisition metadata.
// Raw error text, sample values, bytes and connection settings have no field.
type AcquisitionObservation struct {
	Binding           PointObservationBinding
	DeviceID          string
	PointID           string
	AcquisitionID     string
	ConfigFingerprint string
	Quality           schema.QualityFlag
	QualityReason     string
	ReceivedAt        time.Time
}

// HasAcquisitionObserver reports immutable scheduler construction configuration.
func (s *Scheduler) HasAcquisitionObserver() bool {
	return s != nil && s.config.AcquisitionObserver != nil
}

// AddPointWithObservationLease installs a native point with empty observation
// proof. Its lease belongs to this scheduler's lifetime, not the point name.
// Exhaustion returns zero while preserving ordinary point installation.
func (s *Scheduler) AddPointWithObservationLease(point *schema.Point) uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.installPointLocked(point)
}

func (s *Scheduler) installPointLocked(point *schema.Point) uint64 {
	groupID := ""
	if point.PollingGroupID != nil {
		groupID = *point.PollingGroupID
	}
	lease := uint64(0)
	if s.pointObservationLease < math.MaxUint64 {
		s.pointObservationLease++
		lease = s.pointObservationLease
	}
	s.pointInfos[point.ID] = pointInfo{
		ID: point.ID, DeviceID: point.DeviceID, Address: point.Address,
		Function: point.Function, DataType: point.DataType, DataFormat: point.DataFormat,
		Mode: point.Mode, PollingGroupID: groupID, observationLease: lease,
	}
	return lease
}

// BindPointObservation attaches proof only to the exact current installation.
// It cannot add a point, change acquisition configuration or revive an old lease.
func (s *Scheduler) BindPointObservation(pointID string, lease uint64, binding PointObservationBinding) bool {
	if s == nil || lease == 0 || binding.WorkspaceID == "" || binding.WorkspaceRevision == "" || binding.RuleID == "" || binding.RuleRevision == "" || binding.ConfigFingerprint == "" {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	pt, ok := s.pointInfos[pointID]
	if !ok || pt.observationLease != lease || pt.ID != binding.PointID || pt.DeviceID != binding.DeviceID {
		return false
	}
	cfg, ok := s.deviceConfigs[pt.DeviceID]
	if !ok || cfg.ID != binding.DeviceID {
		return false
	}
	device := schema.Device{ID: cfg.ID, Protocol: cfg.Protocol, ConnectionConfig: cfg.Config}
	point := schema.Point{ID: pt.ID, DeviceID: pt.DeviceID, Address: pt.Address, Function: pt.Function, DataType: pt.DataType, DataFormat: pt.DataFormat, Mode: pt.Mode, PollingGroupID: pollingGroupPointer(pt.PollingGroupID)}
	fingerprint, err := measurement.AcquisitionConfigFingerprint(device, point)
	if err != nil || fingerprint != binding.ConfigFingerprint {
		return false
	}
	pt.observationBinding = binding
	s.pointInfos[pointID] = pt
	return true
}

func (s *Scheduler) observeAcquisition(pt pointInfo, value CollectedValue) {
	if !s.HasAcquisitionObserver() {
		return
	}
	// The callback contract is nonblocking. Panic containment prevents a faulty
	// optional observer from changing the native acquisition/retry result.
	defer func() { recover() }() //nolint:errcheck // An observer panic must not change acquisition.
	s.config.AcquisitionObserver(AcquisitionObservation{
		Binding: pt.observationBinding, DeviceID: value.DeviceID, PointID: value.PointID,
		AcquisitionID: value.AcquisitionID, ConfigFingerprint: value.ConfigFingerprint,
		Quality: value.Quality, QualityReason: value.QualityReason, ReceivedAt: value.ReceivedAt,
	})
}
