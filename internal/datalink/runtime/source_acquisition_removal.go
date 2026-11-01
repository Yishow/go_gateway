package runtime

import (
	"context"

	"go-gateway/internal/datalink/schema"
)

func (s *Service) removeAcquisitionProjectionDevice(ctx context.Context, id string) bool {
	capture := sourceAcquisitionProjectionFrom(ctx)
	if capture == nil || capture.owner != s || capture.devices[id] == nil {
		return false
	}
	device := capture.devices[id]
	if device.removed {
		return true
	}
	s.pointMetaMu.Lock()
	defer s.pointMetaMu.Unlock()
	s.projectionMu.Lock()
	defer s.projectionMu.Unlock()
	expected, present := capture.states[id]
	current, exists := s.projections[id]
	if present != exists || !sameAcquisitionProjectionState(expected, current) {
		return false
	}
	for pointID, meta := range s.pointMetaIndex {
		if meta.DeviceID == id && capture.expectedMeta[pointID] != meta {
			return false
		}
	}
	for pointID, meta := range capture.expectedMeta {
		if meta.DeviceID == id && s.pointMetaIndex[pointID] != meta {
			return false
		}
	}
	if !s.scheduler.RemoveAcquisitionDevice(device.mutation, func() { s.removeDevicePointMetaLocked(id); delete(s.projections, id) }) {
		return false
	}
	device.removed = true
	for pointID, meta := range capture.expectedMeta {
		if meta.DeviceID == id {
			delete(capture.expectedMeta, pointID)
		}
	}
	delete(capture.states, id)
	return true
}

func recordAcquisitionPointMetadata(ctx context.Context, id string, p *schema.Point) {
	capture := sourceAcquisitionProjectionFrom(ctx)
	if capture == nil || capture.expectedMeta == nil {
		return
	}
	if p == nil {
		delete(capture.expectedMeta, id)
	} else {
		capture.expectedMeta[id] = pointMeta{DeviceID: p.DeviceID, Address: p.Address}
	}
}
