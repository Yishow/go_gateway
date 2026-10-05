package sourcerule

import (
	"context"
	"fmt"

	"go-gateway/internal/datalink/schema"
)

func (s *Service) SyncDerivedPointState(ctx context.Context) error {
	rules, err := s.repo.List(ctx, ListFilter{})
	if err != nil {
		return fmt.Errorf("列出來源規則失敗: %w", err)
	}

	for _, rule := range rules {
		if err := s.SyncRuleDerivedState(ctx, rule.ID); err != nil {
			return err
		}
	}

	return nil
}

// SyncRuleDerivedState brings one rule's links, tags and mappings to the state
// a restart would give them. Saving a mapping calls it so the revision a write
// group binds to is already the settled one; otherwise the first restart would
// adopt the mapping and invalidate every group built on it.
func (s *Service) SyncRuleDerivedState(ctx context.Context, ruleID string) error {
	ctx, release := s.AcquireRuleMutation(ctx, ruleID)
	defer release()
	rule, err := s.repo.GetByID(ctx, ruleID)
	if err != nil {
		return fmt.Errorf("取得來源規則失敗: %w", err)
	}
	return s.syncDerivedRuleState(ctx, rule)
}

func (s *Service) syncDerivedRuleState(ctx context.Context, rule *schema.SourceRule, skipMappingIDs ...string) error {
	enabled, err := s.runtimeRuleEnabled(ctx, rule)
	if err != nil {
		return err
	}
	if _, err := s.syncRuleLinksEnabled(ctx, rule.ID, enabled); err != nil {
		return err
	}
	currentLinks, linkErr := s.repo.ListLinks(ctx, rule.ID)
	if linkErr != nil {
		return fmt.Errorf("取得來源規則連結失敗: %w", linkErr)
	}
	nextLinks := cloneSourceRuleLinks(currentLinks)
	syncResult, syncErr := s.syncRuleTagMappings(ctx, rule, rule, nextLinks, enabled, skipMappingIDs...)
	if syncErr != nil {
		return s.finishRollbackErrors(ctx, fmt.Errorf("同步來源規則標籤映射失敗: %w", syncErr), s.rollbackTagMappingSync(ctx, syncResult))
	}
	if err := s.replaceRuleLinks(ctx, rule.ID, currentLinks, nextLinks); err != nil {
		return s.finishRollbackErrors(ctx, err, s.rollbackTagMappingSync(ctx, syncResult))
	}
	return nil
}
