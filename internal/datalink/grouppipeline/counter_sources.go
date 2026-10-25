package grouppipeline

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/workspace"
)

// CounterSources supplies the confirmed counter semantics of a cumulative
// group's members, by tag id. A member it does not return is unconfirmed, and
// the group is not built: nothing is guessed from a tag name or a unit.
type CounterSources interface {
	CounterSources(ctx context.Context, group *workspace.WriteGroup) (map[string]runtime.CounterSource, error)
}

// DefinitionReader reads saved measurement definitions.
type DefinitionReader interface {
	GetByID(ctx context.Context, id string) (*measurement.MeasurementDefinition, error)
}

// MeasurementCounterSources confirms counter semantics from saved measurement
// definitions. A definition counts only when it belongs to the member's own
// workspace, device, point and tag, and its saved semantic kind says counter or
// signed_counter (see runtime.CounterSourceFromDefinition).
type MeasurementCounterSources struct{ definitions DefinitionReader }

// NewMeasurementCounterSources reads definitions from the given repository.
func NewMeasurementCounterSources(definitions DefinitionReader) *MeasurementCounterSources {
	return &MeasurementCounterSources{definitions: definitions}
}

// CounterSources returns the confirmed sources of the group's members. A member
// whose definition is missing, belongs to another source or does not confirm
// counter semantics is left out; a failure to read definitions is an error, so
// an outage is not mistaken for an unconfirmed counter.
func (m *MeasurementCounterSources) CounterSources(ctx context.Context, group *workspace.WriteGroup) (map[string]runtime.CounterSource, error) {
	sources := make(map[string]runtime.CounterSource, len(group.Members))
	for _, member := range group.Members {
		if member.MeasurementID == nil || *member.MeasurementID == "" {
			continue
		}
		definition, err := m.definitions.GetByID(ctx, *member.MeasurementID)
		if errors.Is(err, sql.ErrNoRows) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("read measurement definition: %w", err)
		}
		if definition == nil || definition.WorkspaceID != group.WorkspaceID ||
			definition.DeviceID != member.DeviceID || definition.PointID != member.PointID ||
			(definition.TagID != nil && *definition.TagID != member.TagID) {
			continue
		}
		if source, ok := runtime.CounterSourceFromDefinition(*definition); ok {
			sources[member.TagID] = source
		}
	}
	return sources, nil
}

// needsCounterSources reports a group whose members count from confirmed
// counters: a cumulative long or wide profile group, or an instantaneous long or
// wide profile group with a counter_raw member. Only they read the counter provider.
func needsCounterSources(group *workspace.WriteGroup) bool {
	if group == nil || group.RecordingProfile == nil {
		return false
	}
	switch group.EffectiveLayoutVersion() {
	case workspace.WriteGroupLayoutManagedLongV1:
		return group.RecordingProfile.Kind == workspace.WriteGroupProfileCumulative || hasCounterRawMember(group)
	case workspace.WriteGroupLayoutManagedWideProfileV1:
		return group.RecordingProfile.Kind == workspace.WriteGroupProfileCumulative || hasCounterRawMember(group)
	}
	return false
}

// hasCounterRawMember reports a long group with an instantaneous counter_raw
// member. Its confirmed counter meaning is resolved through the same provider
// the cumulative path uses, so the latest raw value never runs on an unconfirmed
// binding.
func hasCounterRawMember(group *workspace.WriteGroup) bool {
	if group == nil || group.RecordingProfile == nil || group.RecordingProfile.Kind != workspace.WriteGroupProfileInstantaneous {
		return false
	}
	for _, member := range group.Members {
		if member.ValueRole == workspace.WriteGroupValueRoleCounterRaw {
			return true
		}
	}
	return false
}
