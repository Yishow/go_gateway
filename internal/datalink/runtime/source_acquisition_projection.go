package runtime

import (
	"context"
	"maps"
	"slices"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

type sourceAcquisitionProjectionKey struct{}
type sourceAcquisitionDeviceProjection struct {
	device   schema.Device
	mutation collector.AcquisitionGeometryMutation
	graph    workspace.AcquisitionSourceGraph
	update   *collector.AcquisitionGeometryUpdate
	invalid  bool
	removed  bool
}
type sourceAcquisitionProjection struct {
	scopeRule    string
	observations []*acquisitionObservationInstall
	owner        *Service
	devices      map[string]*sourceAcquisitionDeviceProjection
	pointDevices map[string]string
	expectedMeta map[string]pointMeta
	mappings     map[string]FormatRecoveryBindingProof
	states       map[string]runtimeProjectionDeviceState
}

// beginFullAcquisitionProjection retains original graph evidence before any
// profile/provider work. Capture uses the existing mapping -> metadata -> scope
// -> scheduler lock order; all hashing/parsing stays outside those locks.
func (s *Service) beginFullAcquisitionProjection(ctx context.Context, projection *workspace.RuntimeProjection) (projectionContext context.Context, finishProjection func()) {
	capture := &sourceAcquisitionProjection{owner: s, devices: map[string]*sourceAcquisitionDeviceProjection{}, pointDevices: map[string]string{}, expectedMeta: map[string]pointMeta{}, mappings: map[string]FormatRecoveryBindingProof{}, states: map[string]runtimeProjectionDeviceState{}}
	for _, device := range projection.Devices {
		if device == nil {
			continue
		}
		graph, _ := projection.AcquisitionSourceGraph(device.ID)
		capture.devices[device.ID] = &sourceAcquisitionDeviceProjection{graph: graph, device: *device}
	}
	s.mappingMu.RLock()
	s.pointMetaMu.RLock()
	s.projectionMu.RLock()
	for id, bindings := range s.mappingIndex {
		capture.mappings[id] = FormatRecoveryBindingProof{pointID: id, exists: true, installed: bindings, values: slices.Clone(bindings)}
	}
	for id, meta := range s.pointMetaIndex {
		capture.pointDevices[id] = meta.DeviceID
		capture.expectedMeta[id] = meta
	}
	for id, state := range s.projections {
		state.AppliedPointIDs = maps.Clone(state.AppliedPointIDs)
		capture.states[id] = state
	}
	requested := make([]string, 0, len(capture.devices))
	for id := range capture.devices {
		requested = append(requested, id)
	}
	for id, mutation := range s.scheduler.BeginFullAcquisitionGeometryMutation(requested) {
		if capture.devices[id] == nil {
			capture.devices[id] = &sourceAcquisitionDeviceProjection{device: schema.Device{ID: id}}
		}
		capture.devices[id].mutation = mutation
	}
	s.projectionMu.RUnlock()
	s.pointMetaMu.RUnlock()
	s.mappingMu.RUnlock()
	for _, point := range projection.Points {
		if point != nil {
			capture.pointDevices[point.ID] = point.DeviceID
		}
	}
	for _, device := range capture.devices {
		if device.device.Status != schema.DeviceStatusActive {
			continue
		}
		if !device.graph.Complete() {
			device.update = s.scheduler.PrepareAcquisitionGeometry(device.mutation, device.device, collector.AcquisitionGeometryDefinition{}, "")
			continue
		}
		definition := collector.AcquisitionGeometryDefinition{Revision: device.graph.SemanticRevision(), Excluded: device.graph.Excluded(), Rules: map[string]string{}}
		for _, rule := range device.graph.Rules() {
			revision, ok := device.graph.RuleSemanticRevision(rule.ID)
			if !ok {
				device.invalid = true
				break
			}
			definition.Rules[rule.ID] = revision
		}
		device.update = s.scheduler.PrepareAcquisitionGeometry(device.mutation, device.graph.Device(), definition, "")
	}
	finish := func() {
		for _, device := range capture.devices {
			s.scheduler.FinishAcquisitionGeometryMutation(device.mutation)
		}
	}
	return context.WithValue(ctx, sourceAcquisitionProjectionKey{}, capture), finish
}

func sourceAcquisitionProjectionFrom(ctx context.Context) *sourceAcquisitionProjection {
	capture, ok := ctx.Value(sourceAcquisitionProjectionKey{}).(*sourceAcquisitionProjection)
	if !ok {
		return nil
	}
	return capture
}

func recordAcquisitionGeometryInstall(ctx context.Context, p *schema.Point, profile connector.ReadProfile, lease uint64) {
	capture := sourceAcquisitionProjectionFrom(ctx)
	if capture == nil {
		return
	}
	if device := capture.devices[p.DeviceID]; device != nil && device.update != nil && !device.update.Installed(p, profile, lease) {
		device.invalid = true
	}
}

func recordAcquisitionGeometryRemoval(ctx context.Context, id string) {
	capture := sourceAcquisitionProjectionFrom(ctx)
	if capture == nil {
		return
	}
	if device := capture.devices[capture.pointDevices[id]]; device != nil && device.update != nil && !device.update.Removed(id) {
		device.invalid = true
	}
}

func sourceAcquisitionExpectedLease(ctx context.Context, deviceID, pointID string) (uint64, bool) {
	capture := sourceAcquisitionProjectionFrom(ctx)
	if capture == nil {
		return 0, false
	}
	if deviceID == "" {
		deviceID = capture.pointDevices[pointID]
	}
	if device := capture.devices[deviceID]; device != nil {
		if device.removed {
			return 0, true
		}
		if device.update != nil {
			return device.update.ExpectedLease(pointID)
		}
	}
	return 0, false
}

func (s *Service) applyAcquisitionProjectionDevice(ctx context.Context, device *schema.Device) bool {
	capture := sourceAcquisitionProjectionFrom(ctx)
	if capture != nil && capture.owner == s {
		if expected := capture.devices[device.ID]; expected != nil {
			return s.scheduler.ApplyAcquisitionDevice(expected.update)
		}
	}
	s.scheduler.AddDevice(device)
	return true
}
