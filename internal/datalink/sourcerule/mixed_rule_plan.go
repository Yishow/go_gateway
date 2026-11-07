package sourcerule

import (
	"strings"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
)

// planCreateRulePoints derives the point plan of a new rule. For a mixed rule
// it writes the start_address, count and data_type projections into req.
func planCreateRulePoints(req *CreateRuleRequest, protocol schema.ProtocolType, identities ...optimizer.AcquisitionIdentity) (rulePointPlan, error) {
	identity := sourceIdentityForProtocol(protocol, identities)
	if normalizeLayoutMode(req.LayoutMode) != layoutMixed {
		addresses, err := buildPlannedPointAddressesForIdentity(strings.TrimSpace(strings.ToUpper(req.StartAddress)), req.Count, req.DataType, protocol, identity)
		return rulePointPlan{addresses: addresses}, err
	}
	plan, err := planMixedRuleItems(req.MixedItems, protocol, identity)
	if err != nil {
		return rulePointPlan{}, err
	}
	req.StartAddress, req.Count, req.DataType = plan.points[0].Address, len(plan.points), plan.points[0].DataType
	return plan, nil
}

// planUpdateRulePoints applies the requested layout to next and returns its
// point plan. Omitted mixed_items keep the stored items of a mixed rule.
func planUpdateRulePoints(current, next *schema.SourceRule, req UpdateRuleRequest, protocol schema.ProtocolType, identities ...optimizer.AcquisitionIdentity) (rulePointPlan, error) {
	identity := sourceIdentityForProtocol(protocol, identities)
	layout := normalizeLayoutMode(current.LayoutMode)
	if req.LayoutMode != nil {
		layout = normalizeLayoutMode(*req.LayoutMode)
	}
	hasItems := req.MixedItems != nil && len(*req.MixedItems) > 0
	var skipped []string
	if layout == layoutMixed {
		var err error
		if skipped, err = parseSkippedAddresses(next.SkippedAddresses); err != nil {
			return rulePointPlan{}, err
		}
	}
	if err := validateMixedRuleFields(layout, hasItems, next.ShareEnabled, skipped, next.TargetDataType, next.ScaleMultiplier, next.ScaleOffset); err != nil {
		return rulePointPlan{}, err
	}
	if layout != layoutMixed {
		next.LayoutMode, next.MixedItems = "", ""
		addresses, err := buildPlannedPointAddressesForIdentity(next.StartAddress, next.Count, next.DataType, protocol, identity)
		return rulePointPlan{addresses: addresses}, err
	}
	var items []measurement.MixedItem
	if req.MixedItems != nil {
		items = *req.MixedItems
	} else if isMixedRule(current) {
		stored, err := DecodeMixedItems(current.MixedItems)
		if err != nil {
			return rulePointPlan{}, fieldError("mixed_items", "stored mixed_items are unreadable; send the complete item list")
		}
		items = stored
	}
	plan, err := planMixedRuleItems(items, protocol, identity)
	if err != nil {
		return rulePointPlan{}, err
	}
	next.LayoutMode, next.MixedItems = layoutMixed, plan.encoded
	next.StartAddress, next.Count, next.DataType = plan.points[0].Address, len(plan.points), plan.points[0].DataType
	return plan, nil
}

// persistedRuleAddressPlan is the address set the stored rule owns: its item
// addresses for a mixed rule, the corrected MC range plan otherwise.
func persistedRuleAddressPlan(rule *schema.SourceRule, identities ...optimizer.AcquisitionIdentity) ([]string, error) {
	if !isMixedRule(rule) {
		if len(identities) > 0 && rule != nil && qiqrBitAddress(identities[0], rule.StartAddress) {
			if rule.Count <= 0 {
				return nil, validationError("stored count must be greater than zero")
			}
			addresses, err := buildPlannedPointAddressesForIdentity(rule.StartAddress, rule.Count, rule.DataType, schema.ProtocolMC3E, identities[0])
			if err != nil {
				return nil, err
			}
			skipped, err := parseSkippedAddresses(rule.SkippedAddresses)
			if err != nil {
				return nil, err
			}
			kept := make([]string, 0, len(addresses))
			for _, address := range addresses {
				if !containsAddress(skipped, address) {
					kept = append(kept, address)
				}
			}
			return kept, nil
		}
		return correctedMCAddressPlan(rule)
	}
	items, err := DecodeMixedItems(rule.MixedItems)
	if err != nil {
		return nil, err
	}
	addresses := make([]string, 0, len(items))
	for _, item := range items {
		addresses = append(addresses, normalizeAddressKey(item.Address))
	}
	return addresses, nil
}

// mixedLayout is the in-memory rule layout value; homogeneous stays empty.
func (p rulePointPlan) mixedLayout() string {
	if p.mixed() {
		return layoutMixed
	}
	return ""
}
