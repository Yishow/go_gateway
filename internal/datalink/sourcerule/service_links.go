package sourcerule

import (
	"context"
	"fmt"

	"go-gateway/internal/datalink/schema"
)

func (s *Service) ReplaceLinks(ctx context.Context, ruleID string, links []*schema.SourceRuleLink) error {
	ctx, release := s.AcquireRuleMutation(ctx, ruleID)
	defer release()
	if _, err := s.repo.GetByID(ctx, ruleID); err != nil {
		return fmt.Errorf("取得來源規則失敗: %w", err)
	}
	if err := s.repo.ReplaceLinks(ctx, ruleID, links); err != nil {
		return fmt.Errorf("替換來源規則連結失敗: %w", err)
	}
	return nil
}

func (s *Service) replaceRuleLinks(ctx context.Context, ruleID string, _, nextLinks []*schema.SourceRuleLink) error {
	if err := s.repo.ReplaceLinks(ctx, ruleID, nextLinks); err != nil {
		return fmt.Errorf("重建來源規則連結失敗: %w", err)
	}
	return nil
}
