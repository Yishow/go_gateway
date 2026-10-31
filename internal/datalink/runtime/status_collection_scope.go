package runtime

import (
	"maps"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

// Applied scope is recorded with the existing successful projection state.
// Editable points, historical reads and mapping caches are not that authority.
func appliedProjectionPointIDs(projection *workspace.RuntimeProjection, deviceID string) map[string]struct{} {
	ids := make(map[string]struct{})
	active := false
	for _, device := range projection.Devices {
		if device != nil && device.ID == deviceID && device.Status == schema.DeviceStatusActive {
			active = true
			break
		}
	}
	if active {
		for _, point := range projection.Points {
			if point != nil && point.DeviceID == deviceID && point.Enabled {
				ids[point.ID] = struct{}{}
			}
		}
	}
	return ids
}

func (s *Service) appliedCollectionScopes() map[string]map[string]struct{} {
	s.projectionMu.RLock()
	defer s.projectionMu.RUnlock()
	scopes := make(map[string]map[string]struct{}, len(s.projections))
	for id, state := range s.projections {
		if state.AppliedPointIDs != nil {
			scopes[id] = maps.Clone(state.AppliedPointIDs)
		}
	}
	return scopes
}
