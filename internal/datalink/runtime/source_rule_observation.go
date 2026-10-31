package runtime

import (
	"context"
	"errors"
	"slices"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"
)

// SourceRuleObserver receives a finalized historical operation, not current-state authority.
// It must be nonblocking and must not write configuration. Panics are contained.
type SourceRuleObserver func(context.Context, sourcerule.RuntimeSourceObservation)

type sourceWorkspaceSnapshotReader interface {
	ReadSnapshot(context.Context) (*workspace.Record, error)
}

func (s *Service) sourceObservationSnapshot(ctx context.Context) (record *workspace.Record) {
	defer func() {
		if recover() != nil {
			record = nil
		}
	}()
	if s == nil || s.sourceRuleObserver == nil {
		return nil
	}
	reader, ok := s.workspace.(sourceWorkspaceSnapshotReader)
	if !ok {
		return nil
	}
	record, err := reader.ReadSnapshot(ctx)
	if err != nil {
		return nil
	}
	return record
}

func sourceObservationScope(record *workspace.Record, deviceID string) (runtimeview.Scope, bool) {
	if record == nil || record.ID == "" || record.DatabaseSetupRevision == "" || !slices.Contains(record.OrderedDeviceIDs, deviceID) {
		return runtimeview.Scope{}, false
	}
	return runtimeview.Scope{WorkspaceID: record.ID, WorkspaceRevision: record.DatabaseSetupRevision, DeviceID: deviceID}, true
}

func (s *Service) sourceObservation(ctx context.Context, before *workspace.Record, projection *workspace.RuntimeProjection, req sourcerule.RuntimeReconcileRequest, attempted, applied bool) (out *sourcerule.RuntimeSourceObservation) {
	if s == nil || s.sourceRuleObserver == nil {
		return nil
	}
	out = &sourcerule.RuntimeSourceObservation{}
	defer func() {
		if recover() != nil {
			out = &sourcerule.RuntimeSourceObservation{}
		}
	}()
	scope, ok := sourceObservationScope(before, req.Scope.DeviceID)
	if !ok || !attempted || projection == nil || projection.WorkspaceID != scope.WorkspaceID || projection.Version == "" || req.RuleRevision == "" || req.Scope.RuleID == "" {
		return out
	}
	// The selected rule is held by the existing source mutation guard. Delete's
	// revision came from that guard's predelete snapshot; no new/current row is borrowed.
	found := false
	for _, rule := range projection.Rules {
		if rule != nil && rule.ID == req.Scope.RuleID {
			if found || rule.DeviceID != scope.DeviceID || rule.RevisionID != req.RuleRevision || req.Operation == sourcerule.RuntimeReconcileOperationDelete {
				return out
			}
			found = true
		}
	}
	if !found && req.Operation != sourcerule.RuntimeReconcileOperationDelete {
		return out
	}
	after := s.sourceObservationSnapshot(ctx)
	afterScope, ok := sourceObservationScope(after, scope.DeviceID)
	if !ok || afterScope != scope || !slices.Equal(before.OrderedDeviceIDs, after.OrderedDeviceIDs) {
		return out
	}
	if applied && !s.sourceInputStillInstalled(scope.DeviceID, projection.Version) {
		return out
	}
	id, err := common.NewUUID()
	if err != nil {
		return out
	}
	return &sourcerule.RuntimeSourceObservation{Scope: scope, RuleID: req.Scope.RuleID, RuleRevision: req.RuleRevision, ReferenceID: id, ProjectionVersion: projection.Version, Attempted: attempted, Applied: applied, Verified: true}
}

func (s *Service) sourceInputStillInstalled(deviceID, version string) bool {
	s.projectionMu.RLock()
	state, ok := s.projections[deviceID]
	s.projectionMu.RUnlock()
	return ok && state.RuntimeVersion == version && state.WorkspaceVersion == version && state.Alignment == workspace.RuntimeProjectionAlignmentAligned
}

// FinalizeSourceRuleObservation is called only after all source authorities and
// compensation finish. It never reconstructs a projection (which may bootstrap).
func (s *Service) FinalizeSourceRuleObservation(ctx context.Context, outcome sourcerule.RuntimeReconcileOutcome, finalErr error) {
	if s == nil || s.sourceRuleObserver == nil {
		return
	}
	defer func() { recover() }() //nolint:errcheck // Diagnostics must not change the settled result.
	observation := sourcerule.RuntimeSourceObservation{Result: sourcerule.SourceObservationUnavailable}
	if proof := outcome.Observation; proof != nil && proof.Verified && proof.Attempted && outcome.Scope.RuleID == proof.RuleID && outcome.Scope.DeviceID == proof.Scope.DeviceID {
		current, ok := sourceObservationScope(s.sourceObservationSnapshot(ctx), proof.Scope.DeviceID)
		var typed *modbusshare.Error
		uncertain := errors.As(finalErr, &typed) && (typed.Code == modbusshare.ErrCodeDirtyUnknown || typed.DirtyState == "dirty_unknown")
		if ok && current == proof.Scope && !uncertain {
			if finalErr != nil {
				observation = *proof
				observation.Result = sourcerule.SourceObservationRejected
			} else if outcome.Status == sourcerule.RuntimeReconcileStatusAligned && proof.Applied && s.IsRunning() && s.sourceInputStillInstalled(proof.Scope.DeviceID, proof.ProjectionVersion) {
				observation = *proof
				observation.Result = sourcerule.SourceObservationApplied
			}
		}
	}
	s.sourceRuleObserver(ctx, observation)
}

func sourceObservationAttempt(req sourcerule.RuntimeReconcileRequest, projection *workspace.RuntimeProjection) bool {
	switch req.Operation {
	case sourcerule.RuntimeReconcileOperationCreate:
		return len(sourceRuleProjectionPointIDs(projection, req.Scope.RuleID)) > 0
	case sourcerule.RuntimeReconcileOperationUpdate, sourcerule.RuntimeReconcileOperationDelete, sourcerule.RuntimeReconcileOperationEnable, sourcerule.RuntimeReconcileOperationDisable:
		return len(req.AffectedPointIDs) > 0 || len(sourceRuleProjectionPointIDs(projection, req.Scope.RuleID)) > 0 || (req.ProjectionProof != nil && len(req.ProjectionProof.Points) > 0)
	default:
		return false
	}
}
