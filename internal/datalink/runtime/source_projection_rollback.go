package runtime

import (
	"context"
	"fmt"
	"slices"

	"go-gateway/internal/datalink/sourcerule"
)

type sourceProjectionRollbackSnapshot struct {
	owner            *Service
	deviceID         string
	applied          map[string]bool
	mappings         map[string][]mappingBinding
	expectedMappings map[string]FormatRecoveryBindingProof
}

// SnapshotSourceRuleProjection captures the existing authority for affected
// keys only. The opaque value cannot be supplied through an HTTP request.
func (s *Service) SnapshotSourceRuleProjection(deviceID string, ids []string) any {
	snapshot := &sourceProjectionRollbackSnapshot{owner: s, deviceID: deviceID, applied: map[string]bool{}, mappings: map[string][]mappingBinding{}, expectedMappings: map[string]FormatRecoveryBindingProof{}}
	s.mappingMu.RLock()
	defer s.mappingMu.RUnlock()
	s.projectionMu.RLock()
	defer s.projectionMu.RUnlock()
	state := s.projections[deviceID]
	for _, id := range ids {
		bindings, exists := s.mappingIndex[id]
		snapshot.expectedMappings[id] = FormatRecoveryBindingProof{pointID: id, exists: exists, installed: bindings, values: slices.Clone(bindings)}
		_, snapshot.applied[id] = state.AppliedPointIDs[id]
		if bindings, ok := s.mappingIndex[id]; ok {
			snapshot.mappings[id] = slices.Clone(bindings)
		}
	}
	return snapshot
}

// RestoreSourceRuleProjection publishes old affected keys only after durable
// and scheduler compensation have succeeded. A later foreign installation or
// epoch advance rejects the entire restore before any mapping/scope mutation.
func (s *Service) RestoreSourceRuleProjection(value any, proof *sourcerule.RuntimeProjectionProof) error {
	snapshot, ok := value.(*sourceProjectionRollbackSnapshot)
	if !ok || snapshot.owner != s || s.scheduler == nil {
		return fmt.Errorf("source projection rollback snapshot unavailable")
	}
	if proof == nil {
		if len(snapshot.applied) == 0 {
			return nil
		}
		return fmt.Errorf("source projection rollback ownership unavailable")
	}
	// Follow format-recovery publication's existing lock order.
	s.mappingMu.Lock()
	defer s.mappingMu.Unlock()
	s.pointMetaMu.Lock()
	defer s.pointMetaMu.Unlock()
	s.projectionMu.Lock()
	defer s.projectionMu.Unlock()
	for id, expected := range proof.Points {
		bindings, exists := s.mappingIndex[id]
		expectedBinding := snapshot.expectedMappings[id]
		// New owned identities were absent before the mutation. Missing receipt
		// means no generic publication occurred, never permission to adopt a binding.
		if !sameRecoveryBindingInstall(bindings, expectedBinding, exists) {
			return fmt.Errorf("source projection rollback mapping installation changed")
		}
		current, captured := s.scheduler.CaptureOrdinaryPoint(id)
		if !captured || current.Epoch != proof.Epoch || current.Lease != expected.Lease || (expected.Lease != 0 && current.Revision != expected.Revision) {
			return fmt.Errorf("source projection rollback installation changed")
		}
		if meta, exists := s.pointMetaIndex[id]; exists && meta.DeviceID != snapshot.deviceID {
			return fmt.Errorf("source projection rollback device changed")
		}
	}
	if s.mappingIndex == nil {
		s.mappingIndex = map[string][]mappingBinding{}
	}
	if s.projections == nil {
		s.projections = map[string]runtimeProjectionDeviceState{}
	}
	state := s.projections[snapshot.deviceID]
	for id := range proof.Points {
		delete(s.mappingIndex, id)
		if bindings, exists := snapshot.mappings[id]; exists {
			s.mappingIndex[id] = slices.Clone(bindings)
		}
		delete(state.AppliedPointIDs, id)
		if snapshot.applied[id] {
			if state.AppliedPointIDs == nil {
				state.AppliedPointIDs = map[string]struct{}{}
			}
			state.AppliedPointIDs[id] = struct{}{}
		}
	}
	s.projections[snapshot.deviceID] = state
	return nil
}

// The opaque snapshot travels only with the owning synchronous mutation.
type sourceProjectionRollbackKey struct{}

// recordSourceProjectionMappings records the exact just-published mapping
// installation, not a later lookup. Caller holds mappingMu while publishing.
func recordSourceProjectionMappings(ctx context.Context, next map[string][]mappingBinding, pointIDs map[string]struct{}) {
	snapshot, ok := ctx.Value(sourceProjectionRollbackKey{}).(*sourceProjectionRollbackSnapshot)
	if !ok || snapshot == nil {
		return
	}
	for id := range pointIDs {
		values, exists := next[id]
		snapshot.expectedMappings[id] = FormatRecoveryBindingProof{pointID: id, exists: exists, installed: values, values: slices.Clone(values)}
	}
}

// validateSourceProjectionMappings checks the pre-mutation install (or the last
// owned publication receipt) before any forward write. Caller holds mappingMu.
func (s *Service) validateSourceProjectionMappings(ctx context.Context, pointIDs map[string]struct{}) error {
	value := ctx.Value(sourceProjectionRollbackKey{})
	if value == nil {
		return nil
	}
	snapshot, ok := value.(*sourceProjectionRollbackSnapshot)
	if !ok || snapshot == nil || snapshot.owner != s {
		return fmt.Errorf("source projection mapping ownership unavailable")
	}
	for id := range pointIDs {
		current, exists := s.mappingIndex[id]
		if !sameRecoveryBindingInstall(current, snapshot.expectedMappings[id], exists) {
			return fmt.Errorf("source projection mapping installation changed")
		}
	}
	return nil
}
