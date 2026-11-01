package runtime

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"
)

// RestoreSourceRuleProjectionContext retains the original native compensation
// semantics, adding complete geometry only for the active owning operation and
// its original SQL witness. Legacy callers cannot promote unknown metadata.
func (s *Service) RestoreSourceRuleProjectionContext(ctx context.Context, value any, proof *sourcerule.RuntimeProjectionProof) (resultErr error) {
	snapshot, ok := value.(*sourceProjectionRollbackSnapshot)
	if !ok || snapshot == nil || snapshot.owner != s {
		return fmt.Errorf("source projection rollback snapshot unavailable")
	}
	// Acceptance can fail before publishing a rollback receipt (notably for
	// skipped-only rules with no changed leases). Never leave the operation's
	// earlier, looser forward geometry live when its pending token finishes.
	defer func() {
		if resultErr != nil && s.scheduler != nil {
			s.scheduler.RevokeAcquisitionGeometry(snapshot.acquisitionReceipt)
		}
	}()
	operation, ok := ctx.Value(sourceAcquisitionMutationKey{}).(*sourceAcquisitionMutation)
	if !ok || operation == nil || operation.owner != s || operation.deviceID != snapshot.deviceID || !operation.active.Load() || !snapshot.originalState.acquisitionGraph.Complete() {
		return s.RestoreSourceRuleProjection(value, proof)
	}
	if !s.scheduler.CanRestoreAcquisitionGeometry(operation.token) {
		return s.RestoreSourceRuleProjection(value, proof)
	}
	var epoch *uint64
	if proof == nil {
		if len(snapshot.applied) != 0 {
			return fmt.Errorf("source projection rollback ownership unavailable")
		}
		if !snapshot.acquisitionReceipt.Published() {
			return s.RestoreSourceRuleProjection(value, proof)
		}
		proof = &sourcerule.RuntimeProjectionProof{Points: map[string]sourcerule.RuntimePointProof{}}
	} else {
		epoch = &proof.Epoch
	}
	graph := snapshot.originalState.acquisitionGraph
	update := s.scheduler.PrepareAcquisitionGeometryRollback(operation.token, snapshot.acquisitionReceipt, graph.Device(), acquisitionGraphDefinition(graph), epoch)
	if update == nil {
		return fmt.Errorf("source geometry rollback receipt unavailable")
	}
	for id, point := range proof.Points {
		if !update.Restored(id, point.Revision, point.Lease) {
			return fmt.Errorf("source geometry compensated point receipt unavailable")
		}
	}
	publication := update.Seal()
	next := make(map[string][]mappingBinding, len(snapshot.mappings))
	for id, bindings := range snapshot.mappings {
		next[id] = slices.Clone(bindings)
	}
	state := snapshot.expectedState
	state.AppliedPointIDs = maps.Clone(state.AppliedPointIDs)
	if state.AppliedPointIDs == nil {
		state.AppliedPointIDs = map[string]struct{}{}
	}
	for id := range proof.Points {
		delete(state.AppliedPointIDs, id)
		if snapshot.applied[id] {
			state.AppliedPointIDs[id] = struct{}{}
		}
	}
	state.acquisitionGraph = graph
	var receipt collector.AcquisitionGeometryReceipt
	stage := func() error {
		if err := ctx.Err(); err != nil {
			return err
		}
		if !s.mappingMu.TryLock() {
			return fmt.Errorf("source rollback mapping busy")
		}
		defer s.mappingMu.Unlock()
		if !s.pointMetaMu.TryLock() {
			return fmt.Errorf("source rollback metadata busy")
		}
		defer s.pointMetaMu.Unlock()
		if !s.projectionMu.TryLock() {
			return fmt.Errorf("source rollback scope busy")
		}
		defer s.projectionMu.Unlock()
		for id := range proof.Points {
			current, exists := s.mappingIndex[id]
			if !sameRecoveryBindingInstall(current, snapshot.expectedMappings[id], exists) {
				return fmt.Errorf("source rollback mapping installation changed")
			}
			if meta, exists := s.pointMetaIndex[id]; exists && meta.DeviceID != snapshot.deviceID {
				return fmt.Errorf("source rollback point device changed")
			}
		}
		currentState, present := s.projections[snapshot.deviceID]
		if present != snapshot.expectedStatePresent || !sameAcquisitionProjectionState(currentState, snapshot.expectedState) {
			return fmt.Errorf("source rollback scope installation changed")
		}
		var accepted bool
		receipt, accepted = s.scheduler.PublishAcquisitionGeometry(publication, func() {
			for id := range proof.Points {
				delete(s.mappingIndex, id)
				if bindings, exists := next[id]; exists {
					s.mappingIndex[id] = bindings
				}
			}
			s.projections[snapshot.deviceID] = state
		})
		if !accepted {
			return fmt.Errorf("source rollback geometry installation changed")
		}
		return nil
	}
	acceptor, ok := s.workspace.(interface {
		AcceptAcquisitionSourceGraph(context.Context, workspace.AcquisitionSourceGraph, func() error) error
	})
	if !ok {
		return fmt.Errorf("source rollback SQL acceptance unavailable")
	}
	if err := acceptor.AcceptAcquisitionSourceGraph(ctx, graph, stage); err != nil {
		s.scheduler.RevokeAcquisitionGeometry(receipt)
		return err
	}
	return nil
}
