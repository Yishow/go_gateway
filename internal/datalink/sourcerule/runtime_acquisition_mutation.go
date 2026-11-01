package sourcerule

import (
	"context"
	"sync"
)

// sourceAcquisitionMutationLifecycle is an optional native invalidation seam.
// Begin must precede durable source writes; finalize never substitutes for a
// complete SQL-backed reconciliation. It is separate from pure snapshots.
type sourceAcquisitionMutationLifecycle interface {
	BeginSourceRuleAcquisitionMutation(context.Context, RuntimeReconcileScope) (context.Context, func())
}

type sourceAcquisitionMutationContextKey struct{}
type sourceAcquisitionMutationContext struct {
	owner    *Service
	deviceID string
	active   bool
}

func (s *Service) beginSourceAcquisitionMutation(ctx context.Context, ruleID, deviceID string) (nextContext context.Context, finalize func()) {
	if current, ok := ctx.Value(sourceAcquisitionMutationContextKey{}).(*sourceAcquisitionMutationContext); ok && current != nil && current.owner == s && current.deviceID == deviceID && current.active {
		return ctx, func() {}
	}
	if lifecycle, ok := s.runtimeSync.(sourceAcquisitionMutationLifecycle); ok {
		next, finish := lifecycle.BeginSourceRuleAcquisitionMutation(ctx, RuntimeReconcileScope{RuleID: ruleID, DeviceID: deviceID})
		owner := &sourceAcquisitionMutationContext{owner: s, deviceID: deviceID, active: true}
		var once sync.Once
		return context.WithValue(next, sourceAcquisitionMutationContextKey{}, owner), func() { once.Do(func() { owner.active = false; finish() }) }
	}
	return ctx, func() {}
}
