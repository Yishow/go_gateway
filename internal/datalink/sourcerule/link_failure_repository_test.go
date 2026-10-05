package sourcerule

import (
	"context"
	"errors"
	"go-gateway/internal/datalink/schema"
)

func (r *failingLinkRepository) CreateLinks(ctx context.Context, links []*schema.SourceRuleLink) error {
	if r.failCreateLinks {
		r.failCreateLinks = false
		return errors.New("create links failed")
	}
	return r.MemoryRepository.CreateLinks(ctx, links)
}

func (r *failingLinkRepository) ReplaceLinks(ctx context.Context, ruleID string, links []*schema.SourceRuleLink) error {
	if r.failCreateLinks {
		r.failCreateLinks = false
		return errors.New("create links failed")
	}
	return r.MemoryRepository.ReplaceLinks(ctx, ruleID, links)
}

func (r *rollbackSourceRuleRepository) CreateLinks(ctx context.Context, links []*schema.SourceRuleLink) error {
	if r.createLinksFailures > 0 {
		r.createLinksFailures--
		return errors.New("injected source-rule link write failure")
	}
	return r.MemoryRepository.CreateLinks(ctx, links)
}

func (r *rollbackSourceRuleRepository) ReplaceLinks(ctx context.Context, ruleID string, links []*schema.SourceRuleLink) error {
	if r.createLinksFailures > 0 {
		r.createLinksFailures--
		return errors.New("injected source-rule link write failure")
	}
	return r.MemoryRepository.ReplaceLinks(ctx, ruleID, links)
}
