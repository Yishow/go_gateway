package sourcerule

import (
	"context"
	"fmt"
	"go-gateway/internal/datalink/schema"
)

func (s *Service) rollbackRuleState(ctx context.Context, rule *schema.SourceRule, links []*schema.SourceRuleLink) error {
	if err := s.storageCompensationError(ctx); err != nil {
		return err
	}

	if rule == nil {
		return nil
	}
	current, err := s.repo.GetByID(ctx, rule.ID)
	if err != nil {
		return err
	}
	if err := s.restoreOwnedRule(ctx, current, rule); err != nil {
		return fmt.Errorf("restore source rule: %w", err)
	}
	if err := s.repo.ReplaceLinks(ctx, rule.ID, links); err != nil {
		return fmt.Errorf("restore source-rule links: %w", err)
	}
	return s.deleteOwnedCandidates(ctx, rule.ID, rule.RevisionID)
}
