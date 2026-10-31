package runtime

import (
	"context"
	"fmt"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"
)

func (s *Service) applyAffectedSourceRuleProjection(ctx context.Context, projection *workspace.RuntimeProjection, req sourcerule.RuntimeReconcileRequest) error {
	affected := map[string]bool{}
	for _, id := range req.AffectedPointIDs {
		if id != "" {
			affected[id] = true
		}
	}
	for _, id := range sourceRuleProjectionPointIDs(projection, req.Scope.RuleID) {
		affected[id] = true
	}
	if req.ProjectionProof != nil {
		for id := range req.ProjectionProof.Points {
			affected[id] = true
		}
		for id := range affected {
			if _, ok := req.ProjectionProof.Points[id]; !ok {
				return fmt.Errorf("source projection point ownership unavailable")
			}
		}
	}
	if len(affected) == 0 {
		return nil
	}
	present := map[string]bool{}
	var selected []string
	for _, p := range projection.Points {
		if p == nil || !affected[p.ID] {
			continue
		}
		if p.DeviceID != req.Scope.DeviceID {
			return fmt.Errorf("source projection point device changed")
		}
		if req.ProjectionProof != nil {
			proof := req.ProjectionProof.Points[p.ID]
			if proof.Deleted || proof.Revision != p.RevisionID {
				return fmt.Errorf("source projection point revision changed")
			}
		}
		present[p.ID] = true
		selected = append(selected, p.ID)
	}
	// Validate the entire affected batch before changing any point. Publication
	// repeats this CAS under the same lock as its write, covering provider reentry.
	affectedIDs := make(map[string]struct{}, len(affected))
	for id := range affected {
		affectedIDs[id] = struct{}{}
	}
	s.mappingMu.Lock()
	mappingErr := s.validateSourceProjectionMappings(ctx, affectedIDs)
	s.mappingMu.Unlock()
	if mappingErr != nil {
		return mappingErr
	}
	// Preserve public AndPoints validation: deleted identities use this internal
	// mutation-only removal seam, never fabricated entries in the public projection.
	if len(selected) > 0 {
		if err := s.ApplyWorkspaceProjectionForDevicesAndPoints(ctx, projection, []string{req.Scope.DeviceID}, selected); err != nil {
			return err
		}
	}
	for id := range affected {
		if present[id] {
			continue
		}
		if req.ProjectionProof != nil && !req.ProjectionProof.Points[id].Deleted {
			return fmt.Errorf("source projection point unexpectedly missing")
		}
		s.mappingMu.Lock()
		if err := s.validateSourceProjectionMappings(ctx, map[string]struct{}{id: {}}); err != nil {
			s.mappingMu.Unlock()
			return err
		}
		if err := s.removeProjectionPoint(ctx, id); err != nil {
			s.mappingMu.Unlock()
			return err
		}
		delete(s.mappingIndex, id)
		recordSourceProjectionMappings(ctx, nil, map[string]struct{}{id: {}})
		s.mappingMu.Unlock()
		s.projectionMu.Lock()
		if s.projections == nil {
			s.projections = map[string]runtimeProjectionDeviceState{}
		}
		state := s.projections[req.Scope.DeviceID]
		delete(state.AppliedPointIDs, id)
		state.RuntimeVersion = projection.Version
		state.WorkspaceVersion = projection.Version
		state.Alignment = workspace.RuntimeProjectionAlignmentAligned
		state.Message = ""
		state.Code = ""
		s.projections[req.Scope.DeviceID] = state
		s.projectionMu.Unlock()
	}
	return nil
}

func (s *Service) beginSourceRuleProjectionGuard(ctx context.Context, req sourcerule.RuntimeReconcileRequest) context.Context {
	if req.ProjectionSnapshot != nil {
		ctx = context.WithValue(ctx, sourceProjectionRollbackKey{}, req.ProjectionSnapshot)
	}
	if req.ProjectionProof == nil {
		return s.beginSourceProjectionLeaseGuard(ctx, req.Scope.DeviceID)
	}
	proof := *req.ProjectionProof
	proof.Points = make(map[string]sourcerule.RuntimePointProof, len(req.ProjectionProof.Points))
	expected := map[string]uint64{}
	for id, p := range req.ProjectionProof.Points {
		proof.Points[id] = p
		expected[id] = p.Lease
	}
	return context.WithValue(ctx, sourceProjectionLeaseKey{}, &sourceProjectionLeaseGuard{deviceID: req.Scope.DeviceID, expected: expected, proof: &proof, removed: req.PointRemoved})
}

func sourceProjectionPointOwned(g *sourceProjectionLeaseGuard, p *schema.Point) bool {
	if g.proof == nil {
		return true
	}
	proof, ok := g.proof.Points[p.ID]
	return ok && !proof.Deleted && proof.Revision == p.RevisionID
}
