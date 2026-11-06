package sourcerule

import (
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
)

// acquisitionMixedLayoutSupported reports whether the device's saved
// declaration names a random-capable profile. Undeclared, safe-v1, Modbus and
// unreadable configurations keep mixed rules excluded (fail closed).
func acquisitionMixedLayoutSupported(device schema.Device) bool {
	if acquisitionNativeProtocol(device.Protocol) == "modbus" {
		return false
	}
	config, err := connector.ResolveAcquisitionConfig(device.Protocol, device.ConnectionConfig)
	return err == nil && config.Declared() && config.RandomPointLimit() > 0
}

// acquisitionMixedRuleExclusions checks that a mixed rule owns exactly one
// enabled-state-consistent point per stored item, at the item's address and
// type. Its items are ordinary points for the optimizer, so the rule adds
// exclusions only when disabled.
func acquisitionMixedRuleExclusions(identity optimizer.AcquisitionIdentity, rule *schema.SourceRule, links []*schema.SourceRuleLink, points map[string]*schema.Point) ([]optimizer.NativeRange, string) {
	items, err := DecodeMixedItems(rule.MixedItems)
	if err != nil || len(items) == 0 {
		return nil, acquisitionSourceRuleGeometryUnsupported
	}
	expected := make(map[string]optimizer.EffectivePoint, len(items))
	var excluded []optimizer.NativeRange
	for _, item := range items {
		want := optimizer.EffectivePoint{Address: item.Address, DataType: item.DataType}
		span, err := optimizer.AcquisitionPointSpanForIdentity(identity, want)
		key := normalizeAddressKey(item.Address)
		if _, duplicate := expected[key]; err != nil || duplicate {
			return nil, acquisitionSourceRuleGeometryUnsupported
		}
		expected[key] = want
		if !rule.Enabled {
			excluded = append(excluded, span)
		}
	}
	for _, link := range links {
		key := normalizeAddressKey(link.Address)
		want, ok := expected[key]
		point := points[link.PointID]
		if !ok || point == nil || point.DataType != want.DataType || (!rule.Enabled && point.Enabled) {
			return nil, acquisitionSourceLinksIncomplete
		}
		wantSpan, err := optimizer.AcquisitionPointSpanForIdentity(identity, want)
		actual, actualErr := optimizer.AcquisitionPointSpanForIdentity(identity, acquisitionEffectivePoint(point))
		if err != nil || actualErr != nil || wantSpan != actual {
			return nil, acquisitionSourceLinksIncomplete
		}
		delete(expected, key)
	}
	if len(expected) != 0 {
		return nil, acquisitionSourceLinksIncomplete
	}
	return excluded, ""
}
