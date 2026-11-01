package runtime

import (
	"context"
	"sync"
	"sync/atomic"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/sourcerule"
)

type sourceAcquisitionMutationKey struct{}
type sourceAcquisitionMutation struct {
	owner    *Service
	deviceID string
	token    collector.AcquisitionGeometryMutation
	forward  collector.AcquisitionGeometryReceipt
	active   atomic.Bool
}

// BeginSourceRuleAcquisitionMutation explicitly quarantines gap authorization
// before a source mutation. It neither changes installed points nor grants a new
// graph; every exit releases only this operation's token. Nested point writes
// may reuse the same context while its original synchronous owner remains active.
func (s *Service) BeginSourceRuleAcquisitionMutation(ctx context.Context, scope sourcerule.RuntimeReconcileScope) (mutationContext context.Context, finishMutation func()) {
	if s == nil || s.scheduler == nil || scope.DeviceID == "" {
		return ctx, func() {}
	}
	if current, ok := ctx.Value(sourceAcquisitionMutationKey{}).(*sourceAcquisitionMutation); ok && current != nil && current.owner == s && current.deviceID == scope.DeviceID && current.active.Load() {
		return ctx, func() {}
	}
	operation := &sourceAcquisitionMutation{owner: s, deviceID: scope.DeviceID, token: s.scheduler.BeginAcquisitionGeometryMutation(scope.DeviceID)}
	operation.active.Store(true)
	var once sync.Once
	finish := func() {
		once.Do(func() { operation.active.Store(false); s.scheduler.FinishAcquisitionGeometryMutation(operation.token) })
	}
	return context.WithValue(ctx, sourceAcquisitionMutationKey{}, operation), finish
}

// BeginAcquisitionConfigurationMutation is the invalidate-only lifecycle for
// direct Point service writes. The outer source owner, when present, retains
// its pending token through reconciliation and compensation. Telemetry does not
// call this method; direct writes cannot grant new complete source geometry.
func (s *Service) BeginAcquisitionConfigurationMutation(ctx context.Context, deviceID string) func() {
	_, finish := s.BeginSourceRuleAcquisitionMutation(ctx, sourcerule.RuntimeReconcileScope{DeviceID: deviceID})
	return finish
}

// rejectSourceAcquisitionPublication revokes only this synchronous operation's
// forward receipt. A successfully compensated old geometry has another pointer
// and survives. This native safety action is independent of observer enablement.
func (s *Service) rejectSourceAcquisitionPublication(ctx context.Context, finalErr error) {
	if s == nil || s.scheduler == nil || finalErr == nil {
		return
	}
	operation, ok := ctx.Value(sourceAcquisitionMutationKey{}).(*sourceAcquisitionMutation)
	if ok && operation != nil && operation.owner == s && operation.active.Load() {
		s.scheduler.RevokeAcquisitionGeometry(operation.forward)
	}
}
