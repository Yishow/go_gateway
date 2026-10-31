package modbusshare

import (
	"context"
	"github.com/google/uuid"
)

type formatRecoveryCompletion struct {
	key               string
	lease             uint64
	generation, epoch uint64
}

// PublishFormatRecovery runs while a zero-write SQL reservation pins the already
// durable source. TryLock avoids inversion with membership writers. The callback
// must only conditionally activate the already-staged owned lease: no SQL,
// provider reads, listener start, or recording start may run under this lock.
func (s *Service) PublishFormatRecovery(ctx context.Context, fence FormatRecoveryFence, key string, lease uint64, activate func() bool) (HydrationState, error) {
	if err := ctx.Err(); err != nil {
		return HydrationState{}, err
	}
	if !s.mu.TryLock() {
		return HydrationState{}, NewError(ErrCodeRevisionConflict, "format recovery lifecycle is busy", true)
	}
	defer s.mu.Unlock()
	if key == "" || activate == nil || !s.matchesFormatRecoveryFenceLocked(fence) || fence.hydration.State != HydrationStatePending || s.hydrationGeneration == ^uint64(0) {
		return HydrationState{}, NewError(ErrCodeRevisionConflict, "format recovery publication changed", true)
	}
	if !activate() {
		return HydrationState{}, NewError(ErrCodeReconcileFailed, "format recovery installation changed", true)
	}
	state := s.hydration
	state.State = HydrationStateReady
	state.Readiness = true
	state.ReadinessToken = uuid.NewString()
	state.SettingsRevision = s.settings.SettingsRevision
	s.advanceHydrationGenerationLocked()
	s.hydration = state
	s.formatRecoveryCompletion = &formatRecoveryCompletion{key: key, lease: lease, generation: s.hydrationGeneration, epoch: s.lifecycleEpoch}
	return state, nil
}

// ReplayFormatRecovery is read-only and never reinstalls or issues a new token.
// The caller has revalidated the original SQL evidence; validLease performs only
// an atomic check against the same installed profile/incarnation, without I/O.
func (s *Service) ReplayFormatRecovery(key string, validLease func(uint64) bool) (HydrationState, bool) {
	if !s.mu.TryRLock() {
		return HydrationState{}, false
	}
	defer s.mu.RUnlock()
	done := s.formatRecoveryCompletion
	if key == "" || done == nil || done.key != key || done.generation != s.hydrationGeneration || done.epoch != s.lifecycleEpoch || s.hydrationErr != nil || !s.hydration.Readiness || s.hydration.State != HydrationStateReady || validLease == nil || !validLease(done.lease) {
		return HydrationState{}, false
	}
	return s.hydration, true
}

func (f FormatRecoveryFence) Pending() bool {
	return f.hydration.State == HydrationStatePending && !f.hydration.Readiness
}

// RevokeFormatRecoveryPublication handles zero-write reservation cleanup failure.
// It never clears a newer generation/token, including the same operation replay.
// Callers separately suspend only their exact owned runtime lease.
func (s *Service) RevokeFormatRecoveryPublication(key, token string, lease uint64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	done := s.formatRecoveryCompletion
	if done == nil || done.key != key || done.lease != lease || done.generation != s.hydrationGeneration || done.epoch != s.lifecycleEpoch || token == "" || s.hydration.ReadinessToken != token {
		return false
	}
	s.advanceHydrationGenerationLocked()
	s.hydration.State = HydrationStateFailed
	s.hydration.Readiness = false
	s.hydration.ReadinessToken = ""
	return true
}
