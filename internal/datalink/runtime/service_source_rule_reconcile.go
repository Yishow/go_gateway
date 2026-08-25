package runtime

import (
	"context"

	"go-gateway/internal/datalink/sourcerule"
)

// ReconcileSourceRule applies the latest workspace projection after a source-rule lifecycle change.
func (s *Service) ReconcileSourceRule(ctx context.Context, req sourcerule.RuntimeReconcileRequest) sourcerule.RuntimeReconcileOutcome {
	if s == nil || !s.IsRunning() {
		return sourcerule.RuntimeReconcileOutcome{
			Status:  sourcerule.RuntimeReconcileStatusNotRunning,
			Scope:   req.Scope,
			Message: "runtime is not running",
		}
	}
	if s.workspace == nil {
		return sourcerule.RuntimeReconcileOutcome{
			Status:  sourcerule.RuntimeReconcileStatusDeferred,
			Scope:   req.Scope,
			Message: "workspace runtime projection unavailable",
		}
	}

	projection, err := s.workspace.RuntimeProjection(ctx)
	if err != nil {
		return sourcerule.RuntimeReconcileOutcome{
			Status:  sourcerule.RuntimeReconcileStatusStale,
			Scope:   req.Scope,
			Code:    "runtime_projection_unavailable",
			Message: "workspace runtime projection is unavailable",
		}
	}
	if err := s.applyWorkspaceProjectionForDevice(ctx, projection, req.Scope.DeviceID); err != nil {
		return sourcerule.RuntimeReconcileOutcome{
			Status:  sourcerule.RuntimeReconcileStatusStale,
			Scope:   req.Scope,
			Code:    "runtime_projection_apply_failed",
			Message: "workspace runtime projection could not be applied",
		}
	}

	return sourcerule.RuntimeReconcileOutcome{
		Status: sourcerule.RuntimeReconcileStatusAligned,
		Scope:  req.Scope,
	}
}
