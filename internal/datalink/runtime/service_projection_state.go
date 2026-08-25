package runtime

import (
	"context"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

type runtimeProjectionDeviceState struct {
	RuntimeVersion   string
	WorkspaceVersion string
	Alignment        workspace.RuntimeProjectionAlignment
	Message          string
	Code             string
}

func (s *Service) replaceProjectionAligned(projection *workspace.RuntimeProjection) {
	if s == nil || projection == nil {
		return
	}

	next := make(map[string]runtimeProjectionDeviceState, len(projection.DeviceIDs))
	for _, deviceID := range projection.DeviceIDs {
		if deviceID == "" {
			continue
		}
		next[deviceID] = runtimeProjectionDeviceState{
			RuntimeVersion:   projection.Version,
			WorkspaceVersion: projection.Version,
			Alignment:        workspace.RuntimeProjectionAlignmentAligned,
		}
	}

	s.projectionMu.Lock()
	s.projections = next
	s.projectionMu.Unlock()
}

func (s *Service) markDeviceProjectionAligned(deviceID string, projection *workspace.RuntimeProjection) {
	if s == nil || projection == nil || deviceID == "" {
		return
	}

	s.projectionMu.Lock()
	if s.projections == nil {
		s.projections = make(map[string]runtimeProjectionDeviceState)
	}
	s.projections[deviceID] = runtimeProjectionDeviceState{
		RuntimeVersion:   projection.Version,
		WorkspaceVersion: projection.Version,
		Alignment:        workspace.RuntimeProjectionAlignmentAligned,
	}
	s.projectionMu.Unlock()
}

func (s *Service) clearDeviceProjection(deviceID string) {
	if s == nil || deviceID == "" {
		return
	}

	s.projectionMu.Lock()
	delete(s.projections, deviceID)
	s.projectionMu.Unlock()
}

// MarkDeviceProjectionRestartRequired records that a persisted device change cannot be applied in place.
func (s *Service) MarkDeviceProjectionRestartRequired(ctx context.Context, deviceID, message string) {
	if s == nil || deviceID == "" {
		return
	}

	workspaceVersion := ""
	if s.workspace != nil {
		if projection, err := s.workspace.RuntimeProjection(ctx); err == nil && projection != nil {
			workspaceVersion = projection.Version
		}
	}

	s.projectionMu.Lock()
	if s.projections == nil {
		s.projections = make(map[string]runtimeProjectionDeviceState)
	}
	state := s.projections[deviceID]
	state.WorkspaceVersion = workspaceVersion
	state.Alignment = workspace.RuntimeProjectionAlignmentRestartRequired
	state.Message = message
	s.projections[deviceID] = state
	s.projectionMu.Unlock()
}

func (s *Service) markDeviceProjectionAlignedFromWorkspace(ctx context.Context, deviceRecord *schema.Device) {
	if s == nil || s.workspace == nil || deviceRecord == nil {
		return
	}
	projection, err := s.workspace.RuntimeProjection(ctx)
	if err != nil {
		return
	}
	s.markDeviceProjectionAligned(deviceRecord.ID, projection)
}

func (s *Service) projectionStatesForDevices(ctx context.Context, devices []*schema.Device) map[string]runtimeProjectionDeviceState {
	if s == nil || s.workspace == nil || !s.IsRunning() || len(devices) == 0 {
		return nil
	}

	latestProjection, err := s.workspace.RuntimeProjection(ctx)
	if err != nil {
		states := make(map[string]runtimeProjectionDeviceState, len(devices))
		for _, deviceRecord := range devices {
			if deviceRecord == nil {
				continue
			}
			states[deviceRecord.ID] = runtimeProjectionDeviceState{
				Alignment: workspace.RuntimeProjectionAlignmentDegraded,
				Code:      "runtime_projection_unavailable",
				Message:   "workspace runtime projection is unavailable",
			}
		}
		return states
	}
	if latestProjection == nil {
		return nil
	}

	latestDevices := make(map[string]struct{}, len(latestProjection.DeviceIDs))
	for _, deviceID := range latestProjection.DeviceIDs {
		latestDevices[deviceID] = struct{}{}
	}

	states := make(map[string]runtimeProjectionDeviceState, len(devices))
	for _, deviceRecord := range devices {
		if deviceRecord == nil {
			continue
		}
		state := s.deviceProjectionState(deviceRecord.ID)
		if state.WorkspaceVersion == "" {
			state.WorkspaceVersion = latestProjection.Version
		}

		if state.Alignment == workspace.RuntimeProjectionAlignmentRestartRequired {
			state.WorkspaceVersion = latestProjection.Version
			states[deviceRecord.ID] = state
			continue
		}
		if _, exists := latestDevices[deviceRecord.ID]; !exists {
			state.Alignment = workspace.RuntimeProjectionAlignmentStale
			state.WorkspaceVersion = latestProjection.Version
			state.Message = "device is not in the latest workspace projection"
			states[deviceRecord.ID] = state
			continue
		}
		if state.RuntimeVersion == latestProjection.Version {
			state.Alignment = workspace.RuntimeProjectionAlignmentAligned
			state.WorkspaceVersion = latestProjection.Version
			states[deviceRecord.ID] = state
			continue
		}

		state.Alignment = workspace.RuntimeProjectionAlignmentStale
		state.WorkspaceVersion = latestProjection.Version
		if state.Message == "" {
			state.Message = "runtime projection is older than persisted workspace projection"
		}
		states[deviceRecord.ID] = state
	}
	return states
}

func (s *Service) deviceProjectionState(deviceID string) runtimeProjectionDeviceState {
	s.projectionMu.RLock()
	defer s.projectionMu.RUnlock()
	if s.projections == nil {
		return runtimeProjectionDeviceState{}
	}
	return s.projections[deviceID]
}
