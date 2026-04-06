package sourcerule

import (
	"context"
	"fmt"
	"time"

	"go-gateway/internal/datalink/schema"
)

func (s *Service) markStaleTagReviewDecisions(ctx context.Context, rule *schema.SourceRule, activeCandidates []schema.SourceRuleTagCandidate) error {
	if rule == nil {
		return nil
	}

	decisions, err := s.repo.ListTagReviewDecisions(ctx, rule.ID)
	if err != nil {
		return fmt.Errorf("列出來源規則 tag review decisions 失敗: %w", err)
	}
	if len(decisions) == 0 {
		return nil
	}

	activeIDs := make(map[string]struct{}, len(activeCandidates))
	for _, candidate := range activeCandidates {
		activeIDs[candidate.ID] = struct{}{}
	}

	staleAt := time.Now().UTC()
	for _, decision := range decisions {
		if decision == nil || decision.Stale {
			continue
		}
		if _, ok := activeIDs[decision.CandidateID]; ok {
			continue
		}

		decision.Stale = true
		decision.StaleRevisionID = rule.RevisionID
		decision.StaleAt = &staleAt
		decision.UpdatedAt = staleAt
		if err := s.repo.UpsertTagReviewDecision(ctx, decision); err != nil {
			return fmt.Errorf("標記來源規則 tag review decision stale 失敗: %w", err)
		}
	}

	return nil
}
