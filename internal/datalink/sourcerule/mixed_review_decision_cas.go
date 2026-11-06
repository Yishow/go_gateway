package sourcerule

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-gateway/internal/datalink/schema"
)

type tagReviewDecisionWrite struct {
	before, after *schema.SourceRuleTagReviewDecision
}

type tagReviewDecisionCASRepository interface {
	compareAndSwapTagReviewDecision(context.Context, *schema.SourceRuleTagReviewDecision, *schema.SourceRuleTagReviewDecision) error
}

func (s *Service) compareAndSwapTagReviewDecision(ctx context.Context, expected, next *schema.SourceRuleTagReviewDecision) error {
	repo, ok := s.repo.(tagReviewDecisionCASRepository)
	if !ok || !sameTagReviewDecisionIdentity(expected, next) {
		return ErrRuleRevisionConflict
	}
	err := repo.compareAndSwapTagReviewDecision(ctx, expected, next)
	if err != nil && !errors.Is(err, ErrRuleRevisionConflict) {
		if journal := s.pointJournal(ctx); journal != nil {
			journal.ruleUncertain = true
		}
	}
	return err
}

func (s *Service) restoreMixedTagReviewDecisions(ctx context.Context) error {
	journal := s.pointJournal(ctx)
	if journal == nil {
		return nil
	}
	for i := len(journal.reviewDecisionWrites) - 1; i >= 0; i-- {
		write := journal.reviewDecisionWrites[i]
		if err := s.compareAndSwapTagReviewDecision(ctx, write.after, write.before); err != nil {
			return fmt.Errorf("restore mixed tag review decision: %w", err)
		}
	}
	journal.reviewDecisionWrites = nil
	return nil
}

func sameTagReviewDecisionIdentity(a, b *schema.SourceRuleTagReviewDecision) bool {
	return a != nil && b != nil && a.SourceRuleID == b.SourceRuleID && a.CandidateID == b.CandidateID && a.Action == b.Action && a.TagKey == b.TagKey && optionalDecisionString(a.OverrideTagID) == optionalDecisionString(b.OverrideTagID)
}

func optionalDecisionString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func decisionTime(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.UTC().Format(time.RFC3339Nano)
}

func (r *SQLRepository) compareAndSwapTagReviewDecision(ctx context.Context, expected, next *schema.SourceRuleTagReviewDecision) error {
	var staleAt any
	if next.StaleAt != nil {
		staleAt = decisionTime(next.StaleAt)
	}
	result, err := r.db.ExecContext(ctx, `UPDATE source_rule_tag_review_decisions SET stale=?,stale_revision_id=?,stale_at=?,updated_at=? WHERE source_rule_id=? AND candidate_id=? AND decision_type=? AND tag_key=? AND COALESCE(override_tag_id,'')=? AND stale=? AND COALESCE(stale_revision_id,'')=? AND COALESCE(stale_at,'')=?`, next.Stale, next.StaleRevisionID, staleAt, next.UpdatedAt, expected.SourceRuleID, expected.CandidateID, expected.Action, expected.TagKey, optionalDecisionString(expected.OverrideTagID), expected.Stale, expected.StaleRevisionID, decisionTime(expected.StaleAt))
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return ErrRuleRevisionConflict
	}
	return nil
}

func (r *MemoryRepository) compareAndSwapTagReviewDecision(_ context.Context, expected, next *schema.SourceRuleTagReviewDecision) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	current := r.tagReviewDecisions[expected.SourceRuleID][expected.CandidateID]
	if !sameTagReviewDecisionIdentity(current, expected) || current.Stale != expected.Stale || current.StaleRevisionID != expected.StaleRevisionID || decisionTime(current.StaleAt) != decisionTime(expected.StaleAt) {
		return ErrRuleRevisionConflict
	}
	r.tagReviewDecisions[expected.SourceRuleID][expected.CandidateID] = cloneTagReviewDecision(next)
	return nil
}
