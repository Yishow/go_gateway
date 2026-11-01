package runtime

import (
	"context"
	"fmt"
	"maps"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"
)

// prepareScopedAcquisitionProjection consumes the original operation token and
// successful mutation receipts, never a post-provider lookup of current leases.
func (s *Service) prepareScopedAcquisitionProjection(ctx context.Context, projection *workspace.RuntimeProjection, req sourcerule.RuntimeReconcileRequest) (context.Context, error) {
	operation, ok := ctx.Value(sourceAcquisitionMutationKey{}).(*sourceAcquisitionMutation)
	if !ok || operation == nil || operation.owner != s || !operation.active.Load() || operation.deviceID != req.Scope.DeviceID {
		return ctx, nil
	}
	snapshot, owned := req.ProjectionSnapshot.(*sourceProjectionRollbackSnapshot)
	if !owned || snapshot.owner != s || snapshot.deviceID != req.Scope.DeviceID {
		return ctx, fmt.Errorf("source scope snapshot unavailable")
	}
	graph, _ := projection.AcquisitionSourceGraph(req.Scope.DeviceID)
	var device *sourceAcquisitionDeviceProjection
	for _, d := range projection.Devices {
		if d != nil && d.ID == req.Scope.DeviceID {
			device = &sourceAcquisitionDeviceProjection{device: *d, graph: graph, mutation: operation.token}
			break
		}
	}
	if device == nil {
		return ctx, fmt.Errorf("source geometry device unavailable")
	}
	definition := acquisitionGraphDefinition(graph)
	device.update = s.scheduler.PrepareAcquisitionGeometry(operation.token, device.device, definition, req.Scope.RuleID)
	if device.update == nil {
		return ctx, fmt.Errorf("source geometry mutation unavailable")
	}
	capture := &sourceAcquisitionProjection{owner: s, devices: map[string]*sourceAcquisitionDeviceProjection{device.device.ID: device}, pointDevices: map[string]string{}, scopeRule: req.Scope.RuleID, states: map[string]runtimeProjectionDeviceState{req.Scope.DeviceID: snapshot.expectedState}}
	for _, p := range projection.Points {
		if p != nil {
			capture.pointDevices[p.ID] = p.DeviceID
		}
	}
	if req.ProjectionProof == nil {
		device.invalid = len(req.AffectedPointIDs) > 0 || len(sourceRuleProjectionPointIDs(projection, req.Scope.RuleID)) > 0
	} else {
		points := map[string]bool{}
		for _, p := range projection.Points {
			if p == nil || p.DeviceID != req.Scope.DeviceID {
				continue
			}
			proof, owned := req.ProjectionProof.Points[p.ID]
			if !owned {
				continue
			}
			points[p.ID] = true
			if proof.Deleted || proof.Revision != p.RevisionID {
				return ctx, fmt.Errorf("source geometry mutation point changed")
			}
			if !p.Enabled {
				device.update.Removed(p.ID)
				continue
			}
			profile, err := projectionReadProfile(projection, p)
			if !modernPoint(p) {
				profile, err = s.pointReadProfile(ctx, p)
			}
			if err != nil || !device.update.Installed(p, profile, proof.Lease) {
				return ctx, fmt.Errorf("source geometry mutation install receipt unavailable")
			}
		}
		for id, proof := range req.ProjectionProof.Points {
			if !points[id] {
				if !proof.Deleted {
					return ctx, fmt.Errorf("source geometry mutation point missing")
				}
				capture.pointDevices[id] = req.Scope.DeviceID
				device.update.Removed(id)
			}
		}
	}
	return context.WithValue(ctx, sourceAcquisitionProjectionKey{}, capture), nil
}

func acquisitionGraphDefinition(graph workspace.AcquisitionSourceGraph) collector.AcquisitionGeometryDefinition {
	if !graph.Complete() {
		return collector.AcquisitionGeometryDefinition{}
	}
	definition := collector.AcquisitionGeometryDefinition{WorkspaceID: graph.WorkspaceID(), SourceRevision: graph.IdentityRevision(), Revision: graph.SemanticRevision(), Excluded: graph.Excluded(), Rules: map[string]string{}}
	for _, rule := range graph.Rules() {
		revision, ok := graph.RuleSemanticRevision(rule.ID)
		if !ok {
			return collector.AcquisitionGeometryDefinition{}
		}
		definition.Rules[rule.ID] = revision
	}
	return definition
}

func scopedAcquisitionProjection(ctx context.Context) *sourceAcquisitionProjection {
	capture := sourceAcquisitionProjectionFrom(ctx)
	if capture == nil || capture.scopeRule == "" {
		return nil
	}
	return capture
}

func (s *Service) publishScopedAcquisitionProjection(ctx context.Context, projection *workspace.RuntimeProjection, affected map[string]struct{}) error {
	capture := scopedAcquisitionProjection(ctx)
	if capture == nil {
		return nil
	}
	var device *sourceAcquisitionDeviceProjection
	for _, candidate := range capture.devices {
		device = candidate
	}
	if device == nil || device.invalid || device.update == nil {
		return fmt.Errorf("source geometry scoped receipt unavailable")
	}
	next, err := workspaceProjectionMappingBindings(projection, affected, ctx)
	if err != nil {
		return err
	}
	s.annotateMappingObservation(ctx, s.mappingObservationSnapshot(ctx), next, projection)
	expected := capture.states[device.device.ID]
	state := alignedProjectionDeviceState(projection.Version, maps.Clone(expected.AppliedPointIDs))
	if state.AppliedPointIDs == nil {
		state.AppliedPointIDs = map[string]struct{}{}
	}
	for id := range affected {
		delete(state.AppliedPointIDs, id)
	}
	for _, p := range projection.Points {
		if p != nil && p.Enabled {
			if _, ok := affected[p.ID]; ok {
				state.AppliedPointIDs[p.ID] = struct{}{}
			}
		}
	}
	state.acquisitionGraph = device.graph
	publishedState := state
	publishedState.AppliedPointIDs = maps.Clone(state.AppliedPointIDs)
	snapshot, owned := ctx.Value(sourceProjectionRollbackKey{}).(*sourceProjectionRollbackSnapshot)
	if !owned || snapshot.owner != s {
		return fmt.Errorf("source scope snapshot unavailable")
	}
	publication := device.update.Seal()
	var receipt collector.AcquisitionGeometryReceipt
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
		if err := s.validateSourceProjectionMappings(ctx, affected); err != nil {
			return err
		}
		currentState, statePresent := s.projections[device.device.ID]
		if statePresent != snapshot.expectedStatePresent || !sameAcquisitionProjectionState(currentState, expected) {
			return fmt.Errorf("source scope installation changed")
		}
		var ok bool
		receipt, ok = s.scheduler.PublishAcquisitionGeometry(publication, func() {
			for id := range affected {
				delete(s.mappingIndex, id)
			}
			maps.Copy(s.mappingIndex, next)
			s.projections[device.device.ID] = state
		})
		if !ok {
			return fmt.Errorf("source geometry scoped installation changed")
		}
		snapshot.expectedState = publishedState
		snapshot.expectedStatePresent = true
		snapshot.acquisitionReceipt = receipt
		if operation, ok := ctx.Value(sourceAcquisitionMutationKey{}).(*sourceAcquisitionMutation); ok && operation != nil && operation.owner == s && operation.active.Load() {
			operation.forward = receipt
		}
		recordSourceProjectionMappings(ctx, next, affected)
		return nil
	}
	if !device.graph.Complete() {
		err = stage()
	} else {
		acceptor, ok := s.workspace.(interface {
			AcceptAcquisitionSourceGraph(context.Context, workspace.AcquisitionSourceGraph, func() error) error
		})
		if !ok {
			return fmt.Errorf("source SQL acceptance unavailable")
		}
		err = acceptor.AcceptAcquisitionSourceGraph(ctx, device.graph, stage)
	}
	if err != nil {
		s.scheduler.RevokeAcquisitionGeometry(receipt)
		return err
	}
	for _, observation := range capture.observations {
		s.finishAcquisitionObservation(ctx, observation)
	}
	return nil
}
