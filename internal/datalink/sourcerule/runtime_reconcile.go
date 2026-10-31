package sourcerule

import (
	"context"
	"errors"
	"fmt"
	"go-gateway/internal/datalink/common"
	"strings"

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
	ProjectionSnapshot any                       `json:"-"`
	AffectedPointIDs   []string                  `json:"-"`
	ProjectionProof    *RuntimeProjectionProof   `json:"-"`
	PointRemoved       func(string, uint64)      `json:"-"`
	Operation          RuntimeReconcileOperation `json:"operation"`
	Scope              RuntimeReconcileScope     `json:"scope"`
	// RuleRevision is the actual mutation input, not an HTTP field.
	RuleRevision   string                      `json:"-"`
	PointInstalled func(*schema.Point, uint64) `json:"-"`
}

// RuntimeReconcileOutcome reports the source-rule projection result.
type RuntimeReconcileOutcome struct {
	Status      RuntimeReconcileStatus    `json:"status"`
	Scope       RuntimeReconcileScope     `json:"scope"`
	Code        string                    `json:"code,omitempty"`
	Message     string                    `json:"message,omitempty"`
	Observation *RuntimeSourceObservation `json:"-"`
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

func (s *Service) reconcileSourceRuleAuthorities(ctx context.Context, operation RuntimeReconcileOperation, ruleID, deviceID, ruleRevision string) (RuntimeReconcileOutcome, error) {
	scope := RuntimeReconcileScope{RuleID: ruleID, DeviceID: deviceID}
	outcome := s.reconcileSourceRuleRuntime(ctx, operation, ruleID, deviceID, ruleRevision)
	if outcome.Status == RuntimeReconcileStatusStale || outcome.Status == RuntimeReconcileStatusFailed {
		return outcome, &modbusshare.Error{Code: modbusshare.ErrCodeReconcileFailed, Message: "generic runtime projection is not aligned", Retryable: true, Action: retrySourceRuleSaveAction}
	}
	if s.shareRuntimeSync == nil {
		return outcome, nil
	}
	shareOutcome, err := s.shareRuntimeSync.ReconcileSourceRuleShare(ctx, RuntimeReconcileRequest{Operation: operation, Scope: scope})
	// Share alignment does not substitute for generic acquisition evidence.
	shareOutcome.Observation = outcome.Observation
	if err != nil {
		if s.noteTransactionFailure(ctx, err) {
			return shareOutcome, s.dirtyUnknown(ctx)
		}
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
func (s *Service) CreateWithRuntimeReconcile(ctx context.Context, req CreateRuleRequest) (rule *schema.SourceRule, outcome RuntimeReconcileOutcome, err error) {
	if strings.TrimSpace(req.ID) == "" {
		id, err := common.NewUUID()
		if err != nil {
			return nil, RuntimeReconcileOutcome{}, err
		}
		req.ID = id
	}
	ctx, release := s.AcquireRuleMutation(ctx, strings.TrimSpace(req.ID))
	defer release()
	projectionSnapshot := s.snapshotEmptyRuntimeProjection(req.DeviceID)
	if projectionSnapshot != nil {
		ctx = context.WithValue(ctx, sourceProjectionSnapshotKey{}, projectionSnapshot)
	}
	rule, err = s.Create(ctx, req)
	if err != nil {
		return nil, RuntimeReconcileOutcome{}, err
	}
	defer func() { s.finalizeSourceRuleObservation(ctx, outcome, err) }()
	outcome, err = s.reconcileSourceRuleAuthorities(ctx, RuntimeReconcileOperationCreate, rule.ID, rule.DeviceID, rule.RevisionID)
	if err != nil {
		if s.storageCompensationError(ctx) != nil {
			return rule, outcome, err
		}
		var rollbackErrs []error
		if rollbackErr := s.Delete(ctx, rule.ID); rollbackErr != nil {
			return rule, outcome, s.dirtyUnknown(ctx)
		}
		if rollbackErr := s.restoreCapturedRuntimeProjection(ctx, projectionSnapshot); rollbackErr != nil {
			return rule, outcome, s.dirtyUnknown(ctx)
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
) (rule *schema.SourceRule, outcome RuntimeReconcileOutcome, err error) {
	ctx, release := s.AcquireRuleMutation(ctx, id)
	defer release()
	snapshot, err := s.snapshotRule(ctx, id)
	if err != nil {
		return nil, RuntimeReconcileOutcome{}, err
	}
	ctx = withSourceAffectedPoints(ctx, snapshot)
	rule, err = s.Update(ctx, id, req)
	if err != nil {
		// Inner compensation already failed closed. Replaying an unconditional
		// outer snapshot here would overwrite the newer point revision it protected.
		var uncertain *modbusshare.Error
		if errors.As(err, &uncertain) && uncertain.Code == modbusshare.ErrCodeDirtyUnknown {
			return nil, RuntimeReconcileOutcome{}, err
		}
		var formatValidation *formatPreflightError
		var addressValidation *pointAddressPreflightError
		if errors.As(err, &formatValidation) || errors.As(err, &addressValidation) {
			return nil, RuntimeReconcileOutcome{}, err
		}
		var validationErr *mcPreMutationValidationError
		if !errors.Is(err, errMCAddressPlanConflict) && !errors.As(err, &validationErr) {
			if rollbackErr := s.restoreRuleSnapshot(ctx, snapshot); rollbackErr != nil {
				return nil, RuntimeReconcileOutcome{}, s.dirtyUnknown(ctx)
			}
		}
		return nil, RuntimeReconcileOutcome{}, err
	}
	defer func() { s.finalizeSourceRuleObservation(ctx, outcome, err) }()
	outcome, err = s.reconcileSourceRuleAuthorities(ctx, RuntimeReconcileOperationUpdate, rule.ID, rule.DeviceID, rule.RevisionID)
	if err != nil {
		return nil, outcome, s.compensateFailedShareMutation(ctx, snapshot, RuntimeReconcileOperationUpdate, err)
	}
	return rule, outcome, err
}

// DeleteWithRuntimeReconcile deletes a rule and reconciles its runtime state.
func (s *Service) DeleteWithRuntimeReconcile(ctx context.Context, id string) (outcome RuntimeReconcileOutcome, err error) {
	ctx, release := s.AcquireRuleMutation(ctx, id)
	defer release()
	snapshot, err := s.snapshotRule(ctx, id)
	if err != nil {
		return RuntimeReconcileOutcome{}, err
	}
	ctx = withSourceAffectedPoints(ctx, snapshot)
	rule := snapshot.rule
	if err := s.Delete(ctx, id); err != nil {
		var uncertain *modbusshare.Error
		if errors.As(err, &uncertain) && uncertain.Code == modbusshare.ErrCodeDirtyUnknown {
			return RuntimeReconcileOutcome{}, err
		}
		if rollbackErr := s.restoreRuleSnapshot(ctx, snapshot); rollbackErr != nil {
			return RuntimeReconcileOutcome{}, s.dirtyUnknown(ctx)
		}
		return RuntimeReconcileOutcome{}, err
	}
	defer func() { s.finalizeSourceRuleObservation(ctx, outcome, err) }()
	outcome, err = s.reconcileSourceRuleAuthorities(ctx, RuntimeReconcileOperationDelete, rule.ID, rule.DeviceID, rule.RevisionID)
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
) (outcome RuntimeReconcileOutcome, err error) {
	ctx, release := s.AcquireRuleMutation(ctx, id)
	defer release()
	snapshot, err := s.snapshotRule(ctx, id)
	if err != nil {
		return RuntimeReconcileOutcome{}, err
	}
	ctx = withSourceAffectedPoints(ctx, snapshot)
	rule := snapshot.rule
	if err := s.setEnabled(ctx, id, enabled); err != nil {
		var uncertain *modbusshare.Error
		if errors.As(err, &uncertain) && uncertain.Code == modbusshare.ErrCodeDirtyUnknown {
			return RuntimeReconcileOutcome{}, err
		}
		if rollbackErr := s.restoreRuleSnapshot(ctx, snapshot); rollbackErr != nil {
			return RuntimeReconcileOutcome{}, s.dirtyUnknown(ctx)
		}
		return RuntimeReconcileOutcome{}, err
	}
	defer func() { s.finalizeSourceRuleObservation(ctx, outcome, err) }()
	outcome, err = s.reconcileSourceRuleAuthorities(ctx, operation, rule.ID, rule.DeviceID, rule.RevisionID)
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
	ruleRevision ...string,
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

	req := RuntimeReconcileRequest{Operation: operation, Scope: scope, ProjectionProof: s.sourceProjectionProof(ctx), PointRemoved: func(id string, epoch uint64) { s.observeProjectionRemoval(ctx, id, epoch) }, PointInstalled: func(p *schema.Point, lease uint64) { s.observeRuntimeInstallation(ctx, p, lease) }}
	req.ProjectionSnapshot = ctx.Value(sourceProjectionSnapshotKey{})
	req.AffectedPointIDs, _ = ctx.Value(sourceAffectedPointsKey{}).([]string)
	if len(ruleRevision) == 1 {
		req.RuleRevision = ruleRevision[0]
	}
	outcome := reconciler.ReconcileSourceRule(ctx, req)
	if outcome.Scope.RuleID == "" {
		outcome.Scope = scope
	}
	if outcome.Status == "" {
		outcome.Status = RuntimeReconcileStatusStale
		outcome.Message = "runtime reconcile returned empty status"
	}
	return outcome
}
