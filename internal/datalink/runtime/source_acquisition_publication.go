package runtime

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/workspace"
)

func alignedProjectionDeviceState(version string, pointIDs map[string]struct{}) runtimeProjectionDeviceState {
	return runtimeProjectionDeviceState{RuntimeVersion: version, WorkspaceVersion: version, Alignment: workspace.RuntimeProjectionAlignmentAligned, AppliedPointIDs: pointIDs}
}

// publishFullAcquisitionProjection joins the original full mapping/scope write
// to source acceptance. An unknown graph still gets the lawful ordinary install,
// but its publication has no semantic revision and cannot authorize gap reads.
func (s *Service) publishFullAcquisitionProjection(ctx context.Context, projection *workspace.RuntimeProjection) (bool, error) {
	capture := sourceAcquisitionProjectionFrom(ctx)
	if capture == nil || capture.owner != s {
		return false, nil
	}
	before := s.mappingObservationSnapshot(ctx)
	next, err := workspaceProjectionMappingBindings(projection, nil, ctx)
	if err != nil {
		return true, err
	}
	s.annotateMappingObservation(ctx, before, next, projection)
	states := make(map[string]runtimeProjectionDeviceState, len(projection.DeviceIDs))
	for _, id := range projection.DeviceIDs {
		state := alignedProjectionDeviceState(projection.Version, appliedProjectionPointIDs(projection, id))
		if device := capture.devices[id]; device != nil {
			state.acquisitionGraph = device.graph
		}
		states[id] = state
	}
	ids := make([]string, 0, len(capture.devices))
	for id := range capture.devices {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	publications := make([]collector.AcquisitionGeometryPublication, 0, len(ids))
	var graphs []workspace.AcquisitionSourceGraph
	for _, id := range ids {
		device := capture.devices[id]
		if device.removed {
			continue
		}
		if device.invalid || device.update == nil {
			return true, fmt.Errorf("source geometry installation receipt unavailable")
		}
		publications = append(publications, device.update.Seal())
		if device.graph.Complete() {
			graphs = append(graphs, device.graph)
		}
	}
	receipts := make([]collector.AcquisitionGeometryReceipt, len(publications))
	stage := func() error {
		if err := ctx.Err(); err != nil {
			return err
		}
		if !s.mappingMu.TryLock() {
			return fmt.Errorf("source mapping publication busy")
		}
		defer s.mappingMu.Unlock()
		if !s.pointMetaMu.TryLock() {
			return fmt.Errorf("source metadata publication busy")
		}
		defer s.pointMetaMu.Unlock()
		if !s.projectionMu.TryLock() {
			return fmt.Errorf("source scope publication busy")
		}
		defer s.projectionMu.Unlock()
		if !s.fullAcquisitionProjectionCurrent(capture) {
			return fmt.Errorf("source mapping or scope installation changed")
		}
		assign := func() { s.mappingIndex = next; s.projections = states }
		var published bool
		if len(publications) == 0 {
			published = s.scheduler.PublishEmptyAcquisitionProjection(assign)
		} else {
			published = s.scheduler.PublishAcquisitionGeometries(publications, receipts, assign)
		}
		if !published {
			return fmt.Errorf("source geometry installation changed")
		}
		return nil
	}
	if len(graphs) == 0 {
		return true, stage()
	}
	acceptor, ok := s.workspace.(interface {
		AcceptAcquisitionSourceGraphs(context.Context, []workspace.AcquisitionSourceGraph, func() error) error
	})
	if !ok {
		return true, fmt.Errorf("source SQL acceptance unavailable")
	}
	if err := acceptor.AcceptAcquisitionSourceGraphs(ctx, graphs, stage); err != nil {
		for _, receipt := range receipts {
			s.scheduler.RevokeAcquisitionGeometry(receipt)
		}
		return true, err
	}
	return true, nil
}

func (s *Service) fullAcquisitionProjectionCurrent(capture *sourceAcquisitionProjection) bool {
	if !maps.Equal(s.pointMetaIndex, capture.expectedMeta) {
		return false
	}
	if len(s.mappingIndex) != len(capture.mappings) || len(s.projections) != len(capture.states) {
		return false
	}
	for id, expected := range capture.mappings {
		current, exists := s.mappingIndex[id]
		if !sameRecoveryBindingInstall(current, expected, exists) {
			return false
		}
	}
	for id, expected := range capture.states {
		current, exists := s.projections[id]
		if !exists || !sameAcquisitionProjectionState(current, expected) {
			return false
		}
	}
	return true
}

func sameAcquisitionProjectionState(a, b runtimeProjectionDeviceState) bool {
	return a.RuntimeVersion == b.RuntimeVersion && a.WorkspaceVersion == b.WorkspaceVersion && a.Alignment == b.Alignment && a.Message == b.Message && a.Code == b.Code && a.acquisitionGraph == b.acquisitionGraph && maps.Equal(a.AppliedPointIDs, b.AppliedPointIDs)
}
