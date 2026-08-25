package sourcerule

import (
	"context"
	"errors"
	"fmt"

	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/schema"
)

// RuntimeReconcileStatus describes the outcome of a source-rule projection.
type RuntimeReconcileStatus string

// RuntimeReconcileStatus values describe whether runtime projection is aligned.
const (
	RuntimeReconcileStatusAligned         RuntimeReconcileStatus = "aligned"
	RuntimeReconcileStatusDeferred        RuntimeReconcileStatus = "deferred"
	RuntimeReconcileStatusRestartRequired RuntimeReconcileStatus = "restart-required"
	RuntimeReconcileStatusStale           RuntimeReconcileStatus = "stale"
	RuntimeReconcileStatusFailed          RuntimeReconcileStatus = "failed"
	RuntimeReconcileStatusNotRunning      RuntimeReconcileStatus = "not-running"
)

const retrySourceRuleSaveAction = "retry the source-rule save"

// RuntimeReconcileOperation identifies the source-rule mutation being projected.
type RuntimeReconcileOperation string

// RuntimeReconcileOperation values identify supported source-rule mutations.
const (
	RuntimeReconcileOperationEnable  RuntimeReconcileOperation = "enable"
	RuntimeReconcileOperationDisable RuntimeReconcileOperation = "disable"
	RuntimeReconcileOperationUpdate  RuntimeReconcileOperation = "update"
	RuntimeReconcileOperationDelete  RuntimeReconcileOperation = "delete"
	RuntimeReconcileOperationCreate  RuntimeReconcileOperation = "create"
)

// RuntimeReconcileScope identifies the source rule and device being projected.
type RuntimeReconcileScope struct {
	RuleID   string `json:"rule_id"`
	DeviceID string `json:"device_id"`
}

// RuntimeReconcileRequest describes a source-rule projection request.
type RuntimeReconcileRequest struct {
	Operation RuntimeReconcileOperation `json:"operation"`
	Scope     RuntimeReconcileScope     `json:"scope"`
}

// RuntimeReconcileOutcome reports the source-rule projection result.
type RuntimeReconcileOutcome struct {
	Status  RuntimeReconcileStatus `json:"status"`
	Scope   RuntimeReconcileScope  `json:"scope"`
	Code    string                 `json:"code,omitempty"`
	Message string                 `json:"message,omitempty"`
}

// RuntimeRuleReconciler reconciles generic source-rule runtime state.
type RuntimeRuleReconciler interface {
	ReconcileSourceRule(ctx context.Context, req RuntimeReconcileRequest) RuntimeReconcileOutcome
}

// ShareRuntimeReconciler hands generic source-rule mutations to the
// authoritative Local Modbus projection seam when one is configured.
type ShareRuntimeReconciler interface {
	ReconcileSourceRuleShare(context.Context, RuntimeReconcileRequest) (RuntimeReconcileOutcome, error)
}

// ShareRuntimeSnapshot is the reversible Share projection captured around a
// destructive source-rule mutation. Implementations may leave it empty when
// their projection is already transactionally staged.
type ShareRuntimeSnapshot struct {
	Mappings []modbusshare.TagMirrorMapping
	Memory   []byte
}

// ShareRuntimeSnapshotter provides an exact Share projection compensation seam.
type ShareRuntimeSnapshotter interface {
	SnapshotSourceRuleShare(context.Context, RuntimeReconcileRequest) (ShareRuntimeSnapshot, error)
	RestoreSourceRuleShare(context.Context, ShareRuntimeSnapshot) (RuntimeReconcileOutcome, error)
}

// ShareRuntimeCompensator applies the inverse Share operation after a source
// rule mutation fails. It is optional for staged/transactional authorities.
type ShareRuntimeCompensator interface {
	CompensateSourceRuleShare(context.Context, RuntimeReconcileRequest) (RuntimeReconcileOutcome, error)
}

// ShareRuntimeFailCloser is invoked before reporting dirty_unknown when a
// runtime/durable compensation cannot prove the final projection state.
type ShareRuntimeFailCloser interface {
	FailClosed(context.Context)
}

// ShareRuntimeReconcilerFunc adapts a function to ShareRuntimeReconciler.
type ShareRuntimeReconcilerFunc func(context.Context, RuntimeReconcileRequest) (RuntimeReconcileOutcome, error)

// ReconcileSourceRuleShare invokes the adapted Share reconciliation function.
func (f ShareRuntimeReconcilerFunc) ReconcileSourceRuleShare(ctx context.Context, req RuntimeReconcileRequest) (RuntimeReconcileOutcome, error) {
	return f(ctx, req)
}

// SetShareRuntimeReconciler configures the Local Modbus authority handoff.
func (s *Service) SetShareRuntimeReconciler(reconciler ShareRuntimeReconciler) {
	s.shareRuntimeSync = reconciler
}

func (s *Service) reconcileSourceRuleAuthorities(ctx context.Context, operation RuntimeReconcileOperation, ruleID, deviceID string) (RuntimeReconcileOutcome, error) {
	scope := RuntimeReconcileScope{RuleID: ruleID, DeviceID: deviceID}
	outcome := s.reconcileSourceRuleRuntime(ctx, operation, ruleID, deviceID)
	if outcome.Status == RuntimeReconcileStatusStale || outcome.Status == RuntimeReconcileStatusFailed {
		return outcome, &modbusshare.Error{Code: modbusshare.ErrCodeReconcileFailed, Message: "generic runtime projection is not aligned", Retryable: true, Action: retrySourceRuleSaveAction}
	}
	if s.shareRuntimeSync == nil {
		return outcome, nil
	}
	shareOutcome, err := s.shareRuntimeSync.ReconcileSourceRuleShare(ctx, RuntimeReconcileRequest{Operation: operation, Scope: scope})
	if err != nil {
		var shareErr *modbusshare.Error
		if !errors.As(err, &shareErr) {
			return shareOutcome, &modbusshare.Error{Code: modbusshare.ErrCodeReconcileFailed, Message: "share runtime projection reconcile failed", Retryable: true, Action: "retry the source-rule save"}
		}
		return shareOutcome, err
	}
	if shareOutcome.Status != RuntimeReconcileStatusAligned {
		return shareOutcome, &modbusshare.Error{Code: modbusshare.ErrCodeReconcileFailed, Message: "share runtime projection is not aligned", Retryable: true, Action: "retry the source-rule save"}
	}
	return shareOutcome, nil
}

// CreateWithRuntimeReconcile persists a source rule and hands its derived
// runtime projections to their configured authorities.
func (s *Service) CreateWithRuntimeReconcile(ctx context.Context, req CreateRuleRequest) (*schema.SourceRule, RuntimeReconcileOutcome, error) {
	rule, err := s.Create(ctx, req)
	if err != nil {
		return nil, RuntimeReconcileOutcome{}, err
	}
	outcome, err := s.reconcileSourceRuleAuthorities(ctx, RuntimeReconcileOperationCreate, rule.ID, rule.DeviceID)
	if err != nil {
		var rollbackErrs []error
		if rollbackErr := s.Delete(ctx, rule.ID); rollbackErr != nil {
			rollbackErrs = append(rollbackErrs, rollbackErr)
		}
		rollback := s.reconcileSourceRuleRuntime(ctx, RuntimeReconcileOperationDelete, rule.ID, rule.DeviceID)
		if rollback.Status != RuntimeReconcileStatusAligned && rollback.Status != RuntimeReconcileStatusDeferred {
			rollbackErrs = append(rollbackErrs, fmt.Errorf("source-rule runtime compensation did not align: %s", rollback.Status))
		}
		if compensator, ok := s.shareRuntimeSync.(ShareRuntimeCompensator); ok {
			shareRollback, shareErr := compensator.CompensateSourceRuleShare(ctx, RuntimeReconcileRequest{
				Operation: RuntimeReconcileOperationDelete,
				Scope:     RuntimeReconcileScope{RuleID: rule.ID, DeviceID: rule.DeviceID},
			})
			if shareErr != nil {
				rollbackErrs = append(rollbackErrs, shareErr)
			}
			if shareRollback.Status != RuntimeReconcileStatusAligned && shareRollback.Status != RuntimeReconcileStatusDeferred {
				rollbackErrs = append(rollbackErrs, fmt.Errorf("source-rule Share compensation did not align: %s", shareRollback.Status))
			}
		}
		if len(rollbackErrs) > 0 {
			return rule, outcome, s.dirtyUnknown(ctx)
		}
	}
	return rule, outcome, err
}

// EnableWithRuntimeReconcile enables a rule and reconciles its runtime state.
func (s *Service) EnableWithRuntimeReconcile(ctx context.Context, id string) (RuntimeReconcileOutcome, error) {
	return s.setEnabledWithRuntimeReconcile(ctx, id, true, RuntimeReconcileOperationEnable)
}

// DisableWithRuntimeReconcile disables a rule and reconciles its runtime state.
func (s *Service) DisableWithRuntimeReconcile(ctx context.Context, id string) (RuntimeReconcileOutcome, error) {
	return s.setEnabledWithRuntimeReconcile(ctx, id, false, RuntimeReconcileOperationDisable)
}

// UpdateWithRuntimeReconcile updates a rule and reconciles its runtime state.
func (s *Service) UpdateWithRuntimeReconcile(
	ctx context.Context,
	id string,
	req UpdateRuleRequest,
) (*schema.SourceRule, RuntimeReconcileOutcome, error) {
	snapshot, err := s.snapshotRule(ctx, id)
	if err != nil {
		return nil, RuntimeReconcileOutcome{}, err
	}
	rule, err := s.Update(ctx, id, req)
	if err != nil {
		if rollbackErr := s.restoreRuleSnapshot(ctx, snapshot); rollbackErr != nil {
			return nil, RuntimeReconcileOutcome{}, s.dirtyUnknown(ctx)
		}
		return nil, RuntimeReconcileOutcome{}, err
	}
	outcome, err := s.reconcileSourceRuleAuthorities(ctx, RuntimeReconcileOperationUpdate, rule.ID, rule.DeviceID)
	if err != nil {
		return nil, outcome, s.compensateFailedShareMutation(ctx, snapshot, RuntimeReconcileOperationUpdate, err)
	}
	return rule, outcome, err
}

// DeleteWithRuntimeReconcile deletes a rule and reconciles its runtime state.
func (s *Service) DeleteWithRuntimeReconcile(ctx context.Context, id string) (RuntimeReconcileOutcome, error) {
	snapshot, err := s.snapshotRule(ctx, id)
	if err != nil {
		return RuntimeReconcileOutcome{}, err
	}
	rule := snapshot.rule
	if err := s.Delete(ctx, id); err != nil {
		if rollbackErr := s.restoreRuleSnapshot(ctx, snapshot); rollbackErr != nil {
			return RuntimeReconcileOutcome{}, s.dirtyUnknown(ctx)
		}
		return RuntimeReconcileOutcome{}, err
	}
	outcome, err := s.reconcileSourceRuleAuthorities(ctx, RuntimeReconcileOperationDelete, rule.ID, rule.DeviceID)
	if err != nil {
		return outcome, s.compensateFailedShareMutation(ctx, snapshot, RuntimeReconcileOperationDelete, err)
	}
	return outcome, nil
}

func (s *Service) setEnabledWithRuntimeReconcile(
	ctx context.Context,
	id string,
	enabled bool,
	operation RuntimeReconcileOperation,
) (RuntimeReconcileOutcome, error) {
	snapshot, err := s.snapshotRule(ctx, id)
	if err != nil {
		return RuntimeReconcileOutcome{}, err
	}
	rule := snapshot.rule
	if err := s.setEnabled(ctx, id, enabled); err != nil {
		if rollbackErr := s.restoreRuleSnapshot(ctx, snapshot); rollbackErr != nil {
			return RuntimeReconcileOutcome{}, s.dirtyUnknown(ctx)
		}
		return RuntimeReconcileOutcome{}, err
	}
	outcome, err := s.reconcileSourceRuleAuthorities(ctx, operation, rule.ID, rule.DeviceID)
	if err != nil {
		return outcome, s.compensateFailedShareMutation(ctx, snapshot, operation, err)
	}
	return outcome, nil
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
