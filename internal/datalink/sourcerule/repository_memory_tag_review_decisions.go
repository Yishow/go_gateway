package sourcerule

import (
	"context"
	"sort"
	"time"

	"go-gateway/internal/datalink/schema"
)

func (r *MemoryRepository) UpsertTagReviewDecision(_ context.Context, decision *schema.SourceRuleTagReviewDecision) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cloned := cloneTagReviewDecision(decision)
	if cloned == nil {
		return nil
	}
	if cloned.CreatedAt.IsZero() {
		cloned.CreatedAt = time.Now()
	}
	cloned.UpdatedAt = time.Now()

	if r.tagReviewDecisions[cloned.SourceRuleID] == nil {
		r.tagReviewDecisions[cloned.SourceRuleID] = make(map[string]*schema.SourceRuleTagReviewDecision)
	}
	if existing := r.tagReviewDecisions[cloned.SourceRuleID][cloned.CandidateID]; existing != nil {
		cloned.CreatedAt = existing.CreatedAt
	}
	r.tagReviewDecisions[cloned.SourceRuleID][cloned.CandidateID] = cloned
	return nil
}

func (r *MemoryRepository) GetTagReviewDecision(_ context.Context, ruleID, candidateID string) (*schema.SourceRuleTagReviewDecision, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	decisions := r.tagReviewDecisions[ruleID]
	if decisions == nil {
		return nil, nil
	}
	return cloneTagReviewDecision(decisions[candidateID]), nil
}

func (r *MemoryRepository) ListTagReviewDecisions(_ context.Context, ruleID string) ([]*schema.SourceRuleTagReviewDecision, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	decisions := r.tagReviewDecisions[ruleID]
	items := make([]*schema.SourceRuleTagReviewDecision, 0, len(decisions))
	for _, decision := range decisions {
		items = append(items, cloneTagReviewDecision(decision))
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CandidateID < items[j].CandidateID
	})
	return items, nil
}

func cloneTagReviewDecision(decision *schema.SourceRuleTagReviewDecision) *schema.SourceRuleTagReviewDecision {
	if decision == nil {
		return nil
	}
	copy := *decision
	if decision.OverrideTagID != nil {
		value := *decision.OverrideTagID
		copy.OverrideTagID = &value
	}
	if decision.StaleAt != nil {
		value := *decision.StaleAt
		copy.StaleAt = &value
	}
	return &copy
}
