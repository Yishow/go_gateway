package runtime

import (
	"context"
	"maps"

	"go-gateway/internal/datalink/workspace"
)

type acquisitionActivationSelection struct {
	devices  map[string]bool
	points   map[string]bool
	bindings map[string]map[string][]mappingBinding
}

// activationProjectionAssignments merges only the actual selected installation.
// Original full SQL graphs stay in capture; record membership never becomes a
// source graph or an authority to install absent, unselected points.
func (s *Service) activationProjectionAssignments(projection *workspace.RuntimeProjection, capture *sourceAcquisitionProjection) (bindings map[string][]mappingBinding, projections map[string]runtimeProjectionDeviceState) {
	next := make(map[string][]mappingBinding, len(capture.mappings))
	for id, proof := range capture.mappings {
		if proof.exists {
			next[id] = proof.installed
		}
	}
	states := maps.Clone(capture.states)
	if states == nil {
		states = map[string]runtimeProjectionDeviceState{}
	}
	for id, device := range capture.devices {
		scoped := scopedWorkspaceProjectionForDevice(projection, id, capture.selection.devices, capture.selection.points)
		bindings := capture.selection.bindings[id]
		applied := appliedProjectionPointIDs(scoped, id)
		if len(capture.selection.points) > 0 {
			retained := maps.Clone(capture.states[id].AppliedPointIDs)
			if retained == nil {
				retained = map[string]struct{}{}
			}
			for _, p := range scoped.Points {
				delete(retained, p.ID)
				delete(next, p.ID)
			}
			maps.Copy(retained, applied)
			applied = retained
		} else {
			for pointID, deviceID := range capture.pointDevices {
				if deviceID == id {
					delete(next, pointID)
				}
			}
		}
		maps.Copy(next, bindings)
		state := alignedProjectionDeviceState(projection.Version, applied)
		state.acquisitionGraph = device.graph
		states[id] = state
	}
	return next, states
}

func (s *Service) stageActivationBindings(ctx context.Context, projection *workspace.RuntimeProjection, deviceID string, capture *sourceAcquisitionProjection) error {
	before := s.mappingObservationSnapshot(ctx)
	next, err := workspaceProjectionMappingBindings(projection, nil, ctx)
	if err != nil {
		return err
	}
	s.annotateMappingObservation(ctx, before, next, projection)
	if capture.selection.bindings == nil {
		capture.selection.bindings = map[string]map[string][]mappingBinding{}
	}
	capture.selection.bindings[deviceID] = next
	return nil
}
