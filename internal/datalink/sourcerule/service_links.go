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
	before, err := s.repo.ListLinks(ctx, ruleID)
	if err != nil {
		return fmt.Errorf("取得來源規則連結失敗: %w", err)
	}
	ctx, finishAcquisition, err := s.beginSourceLinkAcquisitionMutation(ctx, ruleID, before, links)
	if err != nil {
		return err
	}
	defer finishAcquisition()
	if err := s.repo.ReplaceLinks(ctx, ruleID, links); err != nil {
		if s.noteTransactionFailure(ctx, err) {
			return s.dirtyUnknown(ctx)
		}
		return fmt.Errorf("替換來源規則連結失敗: %w", err)
	}
	return nil
}

func (s *Service) replaceRuleLinks(ctx context.Context, ruleID string, beforeLinks, nextLinks []*schema.SourceRuleLink) error {
	ctx, finishAcquisition, err := s.beginSourceLinkAcquisitionMutation(ctx, ruleID, beforeLinks, nextLinks)
	if err != nil {
		return err
	}
	defer finishAcquisition()
	if err := s.repo.ReplaceLinks(ctx, ruleID, nextLinks); err != nil {
		s.noteTransactionFailure(ctx, err)
		return fmt.Errorf("重建來源規則連結失敗: %w", err)
	}
	return nil
}
