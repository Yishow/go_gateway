package sourcerule

import (
	"encoding/json"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/fatek"
)

const (
	layoutHomogeneous = string(measurement.LayoutModeHomogeneous)
	layoutMixed       = string(measurement.LayoutModeMixed)
)

// plannedRulePoint is one derived point of a rule. Mixed items keep their
// item_id when their address changes.
type plannedRulePoint struct {
	ItemID     string
	Address    string
	DataType   schema.DataType
	DataFormat string
}

// rulePointPlan is the address plan of a rule. A homogeneous plan only lists
// addresses; a mixed plan also carries per-item types and the stored items.
type rulePointPlan struct {
	layout    string
	encoded   string
	addresses []string
	points    []plannedRulePoint
	byAddress map[string]plannedRulePoint
}

func (p rulePointPlan) mixed() bool { return p.layout == layoutMixed }

func (p rulePointPlan) dataType(address string, fallback schema.DataType) schema.DataType {
	if item, ok := p.byAddress[normalizeAddressKey(address)]; ok {
		return item.DataType
	}
	return fallback
}

func (p rulePointPlan) dataFormat(address, fallback string) string {
	if item, ok := p.byAddress[normalizeAddressKey(address)]; ok && item.DataFormat != "" {
		return item.DataFormat
	}
	return fallback
}

// normalizeLayoutMode returns the canonical layout; empty means homogeneous.
func normalizeLayoutMode(raw string) string {
	value := strings.ToLower(strings.TrimSpace(raw))
	if value == "" {
		return layoutHomogeneous
	}
	return value
}

// storedLayoutMode maps a rule's layout to its NOT NULL column value.
func storedLayoutMode(rule *schema.SourceRule) string {
	return normalizeLayoutMode(rule.LayoutMode)
}

// ruleLayoutMode maps a stored layout to the in-memory rule value. Homogeneous
// rules keep the empty value they had before layouts were persisted.
func ruleLayoutMode(stored string) string {
	if layout := normalizeLayoutMode(stored); layout != layoutHomogeneous {
		return layout
	}
	return ""
}

func nullableMixedItems(rule *schema.SourceRule) any {
	if strings.TrimSpace(rule.MixedItems) == "" {
		return nil
	}
	return rule.MixedItems
}

func isMixedRule(rule *schema.SourceRule) bool {
	return rule != nil && normalizeLayoutMode(rule.LayoutMode) == layoutMixed
}

// DecodeMixedItems decodes the stored mixed_items of a source rule.
func DecodeMixedItems(raw string) ([]measurement.MixedItem, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	var items []measurement.MixedItem
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return nil, fmt.Errorf("decode mixed_items: %w", err)
	}
	return items, nil
}

// validateMixedRuleFields rejects rule-level settings that a mixed layout
// would silently ignore or that its consumers do not support yet.
func validateMixedRuleFields(layout string, hasItems, shareEnabled bool, skipped []string, target *schema.DataType, multiplier, offset *float64) error {
	switch {
	case layout != layoutHomogeneous && layout != layoutMixed:
		return fieldError("layout_mode", fmt.Sprintf("unsupported layout_mode %q; use homogeneous or mixed", layout))
	case layout == layoutHomogeneous && hasItems:
		return fieldError("mixed_items", "mixed_items requires layout_mode mixed")
	case layout == layoutHomogeneous:
		return nil
	case shareEnabled:
		return fieldError("share_enabled", "Modbus Share is not supported for mixed layout rules")
	case len(skipped) > 0:
		return fieldError("skipped_addresses", "mixed layout rules list their items instead of skipped addresses")
	case target != nil:
		return fieldError("target_data_type", "rule-level target_data_type is not supported for mixed layout rules")
	case multiplier != nil:
		return fieldError("scale_multiplier", "rule-level scale is not supported for mixed layout rules")
	case offset != nil:
		return fieldError("scale_offset", "rule-level scale is not supported for mixed layout rules")
	}
	return nil
}

// planMixedRuleItems validates mixed items and expands each into one point.
// Items may use different areas; spans come from the acquisition planner so
// the overlap check matches what acquisition would read.
func planMixedRuleItems(items []measurement.MixedItem, protocol schema.ProtocolType, identities ...optimizer.AcquisitionIdentity) (rulePointPlan, error) {
	if len(items) == 0 {
		return rulePointPlan{}, fieldError("mixed_items", "mixed layout must contain at least one item")
	}
	native := acquisitionNativeProtocol(protocol)
	if native == "" {
		return rulePointPlan{}, fieldError("layout_mode", fmt.Sprintf("mixed layout is not supported for protocol %s", protocol))
	}
	identity := sourceIdentityForProtocol(protocol, identities)
	normalized := make([]measurement.MixedItem, len(items))
	spans := make([]optimizer.NativeRange, len(items))
	plan := rulePointPlan{layout: layoutMixed, byAddress: make(map[string]plannedRulePoint, len(items))}
	seen := make(map[string]struct{}, len(items))
	for index, raw := range items {
		item, span, err := normalizeMixedItem(index, raw, protocol, identity)
		if err != nil {
			return rulePointPlan{}, err
		}
		if _, dup := seen[item.ItemID]; dup {
			return rulePointPlan{}, fieldError(fmt.Sprintf("mixed_items[%d].item_id", index), fmt.Sprintf("duplicate item_id %q", item.ItemID), item.ItemID)
		}
		seen[item.ItemID] = struct{}{}
		span = mixedPhysicalSpan(protocol, span)
		for other := range index {
			if spans[other].Area == span.Area && spans[other].Start < span.End && span.Start < spans[other].End {
				a := normalized[other]
				return rulePointPlan{}, fieldError("mixed_items", fmt.Sprintf("item %q (%s %s) overlaps item %q (%s %s)",
					a.ItemID, a.Address, a.DataType, item.ItemID, item.Address, item.DataType), a.ItemID, item.ItemID)
			}
		}
		normalized[index], spans[index] = item, span
		point := plannedRulePoint{ItemID: item.ItemID, Address: item.Address, DataType: item.DataType, DataFormat: item.DataFormat}
		plan.points = append(plan.points, point)
		plan.addresses = append(plan.addresses, item.Address)
		plan.byAddress[normalizeAddressKey(item.Address)] = point
	}
	encoded, err := json.Marshal(normalized)
	if err != nil {
		return rulePointPlan{}, fmt.Errorf("encode mixed_items: %w", err)
	}
	plan.encoded = string(encoded)
	return plan, nil
}

func normalizeMixedItem(index int, item measurement.MixedItem, protocol schema.ProtocolType, identity optimizer.AcquisitionIdentity) (measurement.MixedItem, optimizer.NativeRange, error) {
	field := fmt.Sprintf("mixed_items[%d]", index)
	item.ItemID = strings.TrimSpace(item.ItemID)
	if item.ItemID == "" {
		return item, optimizer.NativeRange{}, fieldError(field+".item_id", "item_id is required")
	}
	item.Name = strings.TrimSpace(item.Name)
	if reason := unsupportedMixedItemSetting(item); reason != "" {
		return item, optimizer.NativeRange{}, fieldError(field, reason, item.ItemID)
	}
	if err := measurement.ValidateMixedLayout([]measurement.MixedItem{item}); err != nil {
		return item, optimizer.NativeRange{}, fieldError(field, err.Error(), item.ItemID)
	}
	if item.DataType == schema.DataTypeString {
		return item, optimizer.NativeRange{}, fieldError(field, "string items are not supported in mixed layout", item.ItemID)
	}
	if err := validateRuleDataType(item.DataType); err != nil {
		return item, optimizer.NativeRange{}, fieldError(field, err.Error(), item.ItemID)
	}
	item.DataFormat = normalizeDataFormat(item.DataFormat)
	if err := validateRuleDataFormat(item.DataFormat); err != nil {
		return item, optimizer.NativeRange{}, fieldError(field+".data_format", err.Error(), item.ItemID)
	}
	address := strings.ToUpper(strings.TrimSpace(item.Address))
	if address == "" {
		return item, optimizer.NativeRange{}, fieldError(field+".address", "address is required", item.ItemID)
	}
	if protocol == schema.ProtocolFatekFBs {
		area, offset, err := fatek.ParseComponentAddress(address)
		if err != nil {
			return item, optimizer.NativeRange{}, fieldError(field+".address", fmt.Sprintf("invalid address %s for %s: %v", address, item.DataType, err), item.ItemID)
		}
		item.Address = fmt.Sprintf("%s%d", area, offset)
	} else {
		planned, err := buildPlannedPointAddressesForIdentity(address, 1, item.DataType, protocol, identity)
		if err != nil || len(planned) != 1 {
			return item, optimizer.NativeRange{}, fieldError(field+".address", fmt.Sprintf("invalid address %s for %s: %v", address, item.DataType, err), item.ItemID)
		}
		item.Address = planned[0]
	}
	span, err := optimizer.AcquisitionPointSpanForIdentity(identity, optimizer.EffectivePoint{Address: item.Address, DataType: item.DataType})
	if err != nil {
		return item, optimizer.NativeRange{}, fieldError(field+".address", fmt.Sprintf("invalid address %s for %s: %v", item.Address, item.DataType, err), item.ItemID)
	}
	return item, span, nil
}

func mixedPhysicalSpan(protocol schema.ProtocolType, span optimizer.NativeRange) optimizer.NativeRange {
	if protocol != schema.ProtocolFatekFBs {
		return span
	}
	switch span.Area {
	case "DR":
		span.Area = "R"
	case "DD":
		span.Area = "D"
	case "DF":
		span.Area = "F"
	case "WM", "DWM":
		span.Area = "M"
	case "WS", "DWS":
		span.Area = "S"
	case "WX", "DWX":
		span.Area = "X"
	case "WY", "DWY":
		span.Area = "Y"
	}
	return span
}

// unsupportedMixedItemSetting names item settings the source rule pipeline
// does not apply yet, so saving them would be a silent no-op.
func unsupportedMixedItemSetting(item measurement.MixedItem) string {
	switch {
	case item.RegisterOffset != 0:
		return "register_offset is not used by source rule items; give each item its own address"
	case item.BitOffset != nil || item.BitLength != nil:
		return "bit slices are not supported for source rule items"
	case item.TargetDataType != nil:
		return "per-item target_data_type is not supported for source rule items"
	case item.ScaleMultiplier != nil || item.ScaleOffset != nil:
		return "per-item scale is not supported for source rule items"
	case item.CounterPolicy != nil:
		return "counter_policy is not supported for source rule items"
	case len(item.StateMap) > 0 || len(item.BitmaskLabels) > 0:
		return "state_map and bitmask_labels are not supported for source rule items"
	}
	return ""
}
