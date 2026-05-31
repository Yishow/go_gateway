package sourcerule

import (
	"context"
	"fmt"

	"go-gateway/internal/datalink/schema"
)

func (s *Service) ReplaceLinks(ctx context.Context, ruleID string, links []*schema.SourceRuleLink) error {
	if _, err := s.repo.GetByID(ctx, ruleID); err != nil {
		return fmt.Errorf("取得來源規則失敗: %w", err)
	}
	if err := s.repo.DeleteLinks(ctx, ruleID); err != nil {
		return fmt.Errorf("刪除來源規則連結失敗: %w", err)
	}
	if err := s.repo.CreateLinks(ctx, links); err != nil {
		return fmt.Errorf("建立來源規則連結失敗: %w", err)
	}
	return nil
}
