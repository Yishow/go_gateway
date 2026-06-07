package sourcerule

import (
	"context"

	"go-gateway/internal/datalink/schema"
)

type RuntimeReconcileStatus string

const (
	RuntimeReconcileStatusAligned         RuntimeReconcileStatus = "aligned"
	RuntimeReconcileStatusDeferred        RuntimeReconcileStatus = "deferred"
	RuntimeReconcileStatusRestartRequired RuntimeReconcileStatus = "restart-required"
	RuntimeReconcileStatusStale           RuntimeReconcileStatus = "stale"
	RuntimeReconcileStatusNotRunning      RuntimeReconcileStatus = "not-running"
)

type RuntimeReconcileOperation string

const (
	RuntimeReconcileOperationEnable  RuntimeReconcileOperation = "enable"
	RuntimeReconcileOperationDisable RuntimeReconcileOperation = "disable"
	RuntimeReconcileOperationUpdate  RuntimeReconcileOperation = "update"
	RuntimeReconcileOperationDelete  RuntimeReconcileOperation = "delete"
)

type RuntimeReconcileScope struct {
	RuleID   string `json:"rule_id"`
	DeviceID string `json:"device_id"`
}

type RuntimeReconcileRequest struct {
	Operation RuntimeReconcileOperation `json:"operation"`
	Scope     RuntimeReconcileScope     `json:"scope"`
}

type RuntimeReconcileOutcome struct {
	Status  RuntimeReconcileStatus `json:"status"`
	Scope   RuntimeReconcileScope  `json:"scope"`
	Message string                 `json:"message,omitempty"`
}

type RuntimeRuleReconciler interface {
	ReconcileSourceRule(ctx context.Context, req RuntimeReconcileRequest) RuntimeReconcileOutcome
}

func (s *Service) EnableWithRuntimeReconcile(ctx context.Context, id string) (RuntimeReconcileOutcome, error) {
	return s.setEnabledWithRuntimeReconcile(ctx, id, true, RuntimeReconcileOperationEnable)
}

func (s *Service) DisableWithRuntimeReconcile(ctx context.Context, id string) (RuntimeReconcileOutcome, error) {
	return s.setEnabledWithRuntimeReconcile(ctx, id, false, RuntimeReconcileOperationDisable)
}

func (s *Service) UpdateWithRuntimeReconcile(
	ctx context.Context,
	id string,
	req UpdateRuleRequest,
) (*schema.SourceRule, RuntimeReconcileOutcome, error) {
	rule, err := s.Update(ctx, id, req)
	if err != nil {
		return nil, RuntimeReconcileOutcome{}, err
	}
	outcome := s.reconcileSourceRuleRuntime(ctx, RuntimeReconcileOperationUpdate, rule.ID, rule.DeviceID)
	return rule, outcome, nil
}

func (s *Service) DeleteWithRuntimeReconcile(ctx context.Context, id string) (RuntimeReconcileOutcome, error) {
	rule, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return RuntimeReconcileOutcome{}, err
	}
	if err := s.Delete(ctx, id); err != nil {
		return RuntimeReconcileOutcome{}, err
	}
	return s.reconcileSourceRuleRuntime(ctx, RuntimeReconcileOperationDelete, rule.ID, rule.DeviceID), nil
}

func (s *Service) setEnabledWithRuntimeReconcile(
	ctx context.Context,
	id string,
	enabled bool,
	operation RuntimeReconcileOperation,
) (RuntimeReconcileOutcome, error) {
	rule, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return RuntimeReconcileOutcome{}, err
	}
	if err := s.setEnabled(ctx, id, enabled); err != nil {
		return RuntimeReconcileOutcome{}, err
	}
	return s.reconcileSourceRuleRuntime(ctx, operation, rule.ID, rule.DeviceID), nil
}

func (s *Service) reconcileSourceRuleRuntime(
	ctx context.Context,
	operation RuntimeReconcileOperation,
	ruleID string,
	deviceID string,
) RuntimeReconcileOutcome {
	scope := RuntimeReconcileScope{RuleID: ruleID, DeviceID: deviceID}
	reconciler, ok := s.runtimeSync.(RuntimeRuleReconciler)
	if s.runtimeSync == nil || !ok {
		return RuntimeReconcileOutcome{
			Status:  RuntimeReconcileStatusDeferred,
			Scope:   scope,
			Message: "runtime reconcile unavailable",
		}
	}

	outcome := reconciler.ReconcileSourceRule(ctx, RuntimeReconcileRequest{
		Operation: operation,
		Scope:     scope,
	})
	if outcome.Scope.RuleID == "" {
		outcome.Scope = scope
	}
	if outcome.Status == "" {
		outcome.Status = RuntimeReconcileStatusStale
		outcome.Message = "runtime reconcile returned empty status"
	}
	return outcome
}
