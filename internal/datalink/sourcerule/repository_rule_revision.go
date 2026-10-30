package sourcerule

import (
	"context"
	"errors"
	"go-gateway/internal/datalink/schema"
)

var ErrRuleRevisionConflict = errors.New("source rule revision changed during mutation")

type ruleRevisionUpdater interface {
	UpdateRevision(context.Context, *schema.SourceRule, *schema.SourceRule) error
}

func (r *SQLRepository) UpdateRevision(ctx context.Context, next, expected *schema.SourceRule) error {
	return r.updateRuleRevision(ctx, next, expected)
}
func (r *MemoryRepository) UpdateRevision(_ context.Context, next, expected *schema.SourceRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	current, ok := r.rules[next.ID]
	if !ok || current.DeviceID != expected.DeviceID || current.RevisionID != expected.RevisionID || !current.UpdatedAt.Equal(expected.UpdatedAt) || ruleStateKey(current) != ruleStateKey(expected) {
		return ErrRuleRevisionConflict
	}
	r.rules[next.ID] = cloneRule(next)
	return nil
}
