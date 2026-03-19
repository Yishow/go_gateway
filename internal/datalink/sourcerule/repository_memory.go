package sourcerule

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"go-gateway/internal/datalink/schema"
)

type MemoryRepository struct {
	mu    sync.RWMutex
	rules map[string]*schema.SourceRule
	links map[string][]*schema.SourceRuleLink
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		rules: make(map[string]*schema.SourceRule),
		links: make(map[string][]*schema.SourceRuleLink),
	}
}

func (r *MemoryRepository) Create(_ context.Context, rule *schema.SourceRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.rules[rule.ID]; exists {
		return fmt.Errorf("source rule %s already exists", rule.ID)
	}
	r.rules[rule.ID] = cloneRule(rule)
	return nil
}

func (r *MemoryRepository) Update(_ context.Context, rule *schema.SourceRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.rules[rule.ID]; !exists {
		return fmt.Errorf("%w: %s", ErrSourceRuleNotFound, rule.ID)
	}
	r.rules[rule.ID] = cloneRule(rule)
	return nil
}

func (r *MemoryRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.rules, id)
	delete(r.links, id)
	return nil
}

func (r *MemoryRepository) GetByID(_ context.Context, id string) (*schema.SourceRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rule, exists := r.rules[id]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrSourceRuleNotFound, id)
	}
	return cloneRule(rule), nil
}

func (r *MemoryRepository) List(_ context.Context, filter ListFilter) ([]*schema.SourceRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]*schema.SourceRule, 0, len(r.rules))
	for _, rule := range r.rules {
		if filter.DeviceID != nil && rule.DeviceID != *filter.DeviceID {
			continue
		}
		if filter.Enabled != nil && rule.Enabled != *filter.Enabled {
			continue
		}
		items = append(items, cloneRule(rule))
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].CreatedAt.Equal(items[j].CreatedAt) {
			return items[i].ID < items[j].ID
		}
		return items[i].CreatedAt.Before(items[j].CreatedAt)
	})

	return items, nil
}

func (r *MemoryRepository) CreateLinks(_ context.Context, links []*schema.SourceRuleLink) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(links) == 0 {
		return nil
	}

	for _, link := range links {
		r.links[link.RuleID] = append(r.links[link.RuleID], cloneLink(link))
	}

	sort.Slice(r.links[links[0].RuleID], func(i, j int) bool {
		return r.links[links[0].RuleID][i].Address < r.links[links[0].RuleID][j].Address
	})
	return nil
}

func (r *MemoryRepository) ListLinks(_ context.Context, ruleID string) ([]*schema.SourceRuleLink, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	links := r.links[ruleID]
	result := make([]*schema.SourceRuleLink, 0, len(links))
	for _, link := range links {
		result = append(result, cloneLink(link))
	}
	return result, nil
}

func (r *MemoryRepository) DeleteLinks(_ context.Context, ruleID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.links, ruleID)
	return nil
}

func cloneRule(rule *schema.SourceRule) *schema.SourceRule {
	if rule == nil {
		return nil
	}
	copy := *rule
	return &copy
}

func cloneLink(link *schema.SourceRuleLink) *schema.SourceRuleLink {
	if link == nil {
		return nil
	}
	copy := *link
	if link.TagID != nil {
		value := *link.TagID
		copy.TagID = &value
	}
	if link.MappingID != nil {
		value := *link.MappingID
		copy.MappingID = &value
	}
	return &copy
}
