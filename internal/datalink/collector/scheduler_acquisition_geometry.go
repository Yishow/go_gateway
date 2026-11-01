package collector

import (
	"go-gateway/internal/datalink/optimizer"
)

// AcquisitionGeometryReceipt is an opaque capture of the current publication
// and entire installed device slot set. Capture alone never grants readability.
// Copies are immutable; receipts are not source revisions or sample identities.
type AcquisitionGeometryReceipt struct {
	authorized    bool
	ownedMutation *acquisitionGeometryMutation
	owner         *Scheduler
	deviceID      string
	device        deviceConfig
	present       bool
	publication   *acquisitionGeometry
	slots         map[string]pointInfo
}

type acquisitionGeometry struct {
	revision string
	excluded []optimizer.NativeRange
	rules    map[string]string
	device   deviceConfig
	slots    map[string]pointInfo
}

type acquisitionGeometryState struct {
	publication *acquisitionGeometry
	pending     map[*acquisitionGeometryMutation]struct{}
}

type acquisitionGeometryMutation struct {
	owner    *Scheduler
	deviceID string
	marker   *acquisitionGeometry
	before   AcquisitionGeometryReceipt
}

// AcquisitionGeometryMutation is an opaque, operation-local quarantine token.
// Ending one token cannot end another operation or restore prior authorization.
type AcquisitionGeometryMutation struct{ token *acquisitionGeometryMutation }

// CaptureAcquisitionGeometry is pure: it captures no provider state and does
// not quarantine a device. The caller must use the explicit mutation lifecycle.
func (s *Scheduler) CaptureAcquisitionGeometry(deviceID string) AcquisitionGeometryReceipt {
	if s == nil || deviceID == "" {
		return AcquisitionGeometryReceipt{}
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.captureAcquisitionGeometryLocked(deviceID)
}

func (s *Scheduler) captureAcquisitionGeometryLocked(deviceID string) AcquisitionGeometryReceipt {
	device, present := s.deviceConfigs[deviceID]
	receipt := AcquisitionGeometryReceipt{owner: s, deviceID: deviceID, device: device, present: present, slots: map[string]pointInfo{}}
	if state := s.acquisitionGeometry[deviceID]; state != nil {
		receipt.publication = state.publication
	}
	for id, point := range s.pointInfos {
		if point.DeviceID == deviceID {
			receipt.slots[id] = acquisitionSlotProof(point)
		}
	}
	receipt.authorized = receipt.publication != nil && s.currentAcquisitionGeometryLocked(deviceID) == receipt.publication
	return receipt
}

func acquisitionSlotProof(point pointInfo) pointInfo {
	// Optional diagnostics and mapping bindings cannot become source authority.
	point.observationBinding = PointObservationBinding{}
	return point
}

// BeginAcquisitionGeometryMutation invalidates gap authorization before durable
// source mutation. It never changes a Point, lease, profile, mapping or identity.
func (s *Scheduler) BeginAcquisitionGeometryMutation(deviceID string) AcquisitionGeometryMutation {
	if s == nil || deviceID == "" {
		return AcquisitionGeometryMutation{}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.beginAcquisitionGeometryMutationLocked(deviceID)
}

func (s *Scheduler) beginAcquisitionGeometryMutationLocked(deviceID string) AcquisitionGeometryMutation {
	before := s.captureAcquisitionGeometryLocked(deviceID)
	if s.acquisitionGeometry == nil {
		s.acquisitionGeometry = map[string]*acquisitionGeometryState{}
	}
	state := s.acquisitionGeometry[deviceID]
	if state == nil {
		state = &acquisitionGeometryState{pending: map[*acquisitionGeometryMutation]struct{}{}}
		s.acquisitionGeometry[deviceID] = state
	}
	marker := &acquisitionGeometry{}
	token := &acquisitionGeometryMutation{owner: s, deviceID: deviceID, marker: marker, before: before}
	state.publication = marker
	state.pending[token] = struct{}{}
	return AcquisitionGeometryMutation{token: token}
}

// FinishAcquisitionGeometryMutation releases only this operation's quarantine.
// With no accepted new publication, unknown remains unknown even after success,
// preflight refusal, cancellation or compensation; an old graph is never revived.
func (s *Scheduler) FinishAcquisitionGeometryMutation(receipt AcquisitionGeometryMutation) {
	token := receipt.token
	if s == nil || token == nil || token.owner != s {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if state := s.acquisitionGeometry[token.deviceID]; state != nil {
		if _, owned := state.pending[token]; !owned {
			return
		}
		delete(state.pending, token)
		// Failed/seed writes to a device not installed here cannot accumulate
		// retained unknown cache entries. Pending overlaps remain quarantined.
		if _, installed := s.deviceConfigs[token.deviceID]; !installed && len(state.pending) == 0 && (state.publication == nil || state.publication.revision == "") {
			delete(s.acquisitionGeometry, token.deviceID)
		}
	}
}

func (s *Scheduler) currentAcquisitionGeometryLocked(deviceID string) *acquisitionGeometry {
	state := s.acquisitionGeometry[deviceID]
	if state == nil || len(state.pending) != 0 || state.publication == nil {
		return nil
	}
	geometry := state.publication
	if geometry.revision == "" || s.deviceConfigs[deviceID] != geometry.device || !s.acquisitionSlotsCurrentLocked(deviceID, geometry.slots) {
		return nil
	}
	return geometry
}

// The full-set comparison catches removal paths that intentionally do not
// advance pointObservationLease, and catches extra foreign slots as well.
// Admission holds only scheduler.mu, performs no I/O, hashing or allocations.
func (s *Scheduler) acquisitionSlotsCurrentLocked(deviceID string, expected map[string]pointInfo) bool {
	return s.acquisitionSlotsMatchLocked(deviceID, expected, true)
}

func (s *Scheduler) acquisitionSlotsMatchLocked(deviceID string, expected map[string]pointInfo, usable bool) bool {
	count := 0
	for id, point := range s.pointInfos {
		if point.DeviceID != deviceID {
			continue
		}
		captured, exists := expected[id]
		if !exists || (usable && (point.recoveryPending || !s.installationConnectionCurrentLocked(point))) || acquisitionSlotProof(point) != captured {
			return false
		}
		count++
	}
	return count == len(expected)
}
