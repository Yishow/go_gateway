package sourcerule

import (
	"context"
	"sync"
	"sync/atomic"
)

type ruleMutationGuard struct {
	mu    sync.Mutex
	users int
}
type ruleMutationOwnership struct {
	service *Service
	ruleID  string
	active  atomic.Bool
}
type ruleMutationContextKey struct{}

// AcquireRuleMutation coordinates a synchronous rule operation and its nested
// service calls. The returned context must never be shared with another goroutine.
// Never acquire a different rule with this owned context. Bulk operations
// acquire and release one rule at a time using the original context.
// HTTP saves release before workspace readiness/apply. Runtime source wrappers
// retain ownership through reconcile/compensation; their callbacks use unguarded
// source reads. Source read methods remain
// unguarded because workspace mutex holders call them; locking reads would invert
// the workspace/rule lock order. A released context cannot bypass a later guard.
func (s *Service) AcquireRuleMutation(ctx context.Context, ruleID string) (guardCtx context.Context, releaseGuard func()) {
	if owner, ok := ctx.Value(ruleMutationContextKey{}).(*ruleMutationOwnership); ok && owner.service == s && owner.ruleID == ruleID && owner.active.Load() {
		return ctx, func() {}
	}
	s.ruleGuardsMu.Lock()
	if s.ruleGuards == nil {
		s.ruleGuards = make(map[string]*ruleMutationGuard)
	}
	guard := s.ruleGuards[ruleID]
	if guard == nil {
		guard = &ruleMutationGuard{}
		s.ruleGuards[ruleID] = guard
	}
	guard.users++
	s.ruleGuardsMu.Unlock()
	guard.mu.Lock()
	owner := &ruleMutationOwnership{service: s, ruleID: ruleID}
	owner.active.Store(true)
	var once sync.Once
	release := func() {
		once.Do(func() {
			owner.active.Store(false)
			guard.mu.Unlock()
			s.ruleGuardsMu.Lock()
			guard.users--
			if guard.users == 0 {
				delete(s.ruleGuards, ruleID)
			}
			s.ruleGuardsMu.Unlock()
		})
	}
	return context.WithValue(ctx, ruleMutationContextKey{}, owner), release
}
