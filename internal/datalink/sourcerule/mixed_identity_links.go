package sourcerule

import (
	"strings"

	"go-gateway/internal/datalink/schema"
)

// matchRulePlanLinks reuses a mixed point only for the same persisted item ID.
// Layout changes without item identities keep the existing address-based path.
func matchRulePlanLinks(rule *schema.SourceRule, plan rulePointPlan, links []*schema.SourceRuleLink) (map[string]*schema.SourceRuleLink, error) {
	byAddress := make(map[string]*schema.SourceRuleLink, len(links))
	for _, link := range links {
		byAddress[normalizeAddressKey(link.Address)] = link
	}
	if !isMixedRule(rule) || !plan.mixed() {
		return byAddress, nil
	}
	items, err := DecodeMixedItems(rule.MixedItems)
	if err != nil {
		return nil, &pointAddressPreflightError{cause: fieldError("mixed_items", "stored item identities are unreadable")}
	}
	oldByID := make(map[string]*schema.SourceRuleLink, len(items))
	seenIDs, seenAddresses := map[string]bool{}, map[string]bool{}
	for _, item := range items {
		id, address := strings.TrimSpace(item.ItemID), normalizeAddressKey(item.Address)
		if id == "" || address == "" || seenIDs[id] || seenAddresses[address] {
			return nil, &pointAddressPreflightError{cause: fieldError("mixed_items", "stored item identities are ambiguous")}
		}
		seenIDs[id], seenAddresses[address] = true, true
		oldByID[id] = byAddress[address]
	}
	matched := make(map[string]*schema.SourceRuleLink, len(plan.points))
	for _, item := range plan.points {
		if link := oldByID[item.ItemID]; link != nil {
			matched[normalizeAddressKey(item.Address)] = link
		}
	}
	return matched, nil
}
