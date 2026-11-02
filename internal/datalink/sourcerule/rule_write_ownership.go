package sourcerule

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"go-gateway/internal/datalink/schema"
	"time"
)

func ruleStateKey(rule *schema.SourceRule) (key [32]byte, valid bool) {
	if rule == nil {
		return [32]byte{}, false
	}
	normalized := cloneRule(rule)
	normalized.CreatedAt = normalized.CreatedAt.UTC()
	// SQLite update triggers replace this timestamp with CURRENT_TIMESTAMP.
	// Configuration plus revision is the attempted state; SQL CAS separately
	// compares the actual stored timestamp read with that state.
	normalized.UpdatedAt = time.Time{}
	data, err := json.Marshal(normalized)
	if err != nil {
		return [32]byte{}, false
	}
	return sha256.Sum256(data), true
}
func (s *Service) noteRuleBaseline(ctx context.Context, rule *schema.SourceRule) {
	if j := s.pointJournal(ctx); j != nil && j.ruleTargets == nil {
		key, valid := ruleStateKey(rule)
		if !valid {
			return
		}
		j.ruleTargets = map[[32]byte]bool{key: true}
		j.deviceID = rule.DeviceID
	}
}
func (s *Service) ownsRuleState(ctx context.Context, rule *schema.SourceRule) bool {
	j := s.pointJournal(ctx)
	if j == nil {
		return false
	}
	key, valid := ruleStateKey(rule)
	return valid && j.ruleTargets[key]
}
func (s *Service) updateOwnedRule(ctx context.Context, before, next *schema.SourceRule) error {
	s.noteRuleBaseline(ctx, before)
	if !s.ownsRuleState(ctx, before) {
		return ErrRuleRevisionConflict
	}
	nextKey, valid := ruleStateKey(next)
	if !valid {
		return ErrRuleRevisionConflict
	}
	j := s.pointJournal(ctx)
	if j != nil {
		j.ruleTargets[nextKey] = true
	}
	if repo, ok := s.repo.(ruleRevisionUpdater); ok {
		return repo.UpdateRevision(ctx, next, before)
	}
	// Non-SQL fixture repositories retain the existing shared rule guard contract.
	return s.repo.Update(ctx, next)
}
func (s *Service) restoreOwnedRule(ctx context.Context, current, original *schema.SourceRule) error {
	if !s.ownsRuleState(ctx, current) {
		return ErrRuleRevisionConflict
	}
	if _, valid := ruleStateKey(original); !valid {
		return ErrRuleRevisionConflict
	}
	if repo, ok := s.repo.(ruleRevisionUpdater); ok {
		return repo.UpdateRevision(ctx, original, current)
	}
	return s.repo.Update(ctx, original)
}

func (s *Service) restoreCurrentOwnedRule(ctx context.Context, original *schema.SourceRule) error {
	if err := s.storageCompensationError(ctx); err != nil {
		return err
	}

	current, err := s.repo.GetByID(ctx, original.ID)
	if err != nil {
		return err
	}
	return s.restoreOwnedRule(ctx, current, original)
}
