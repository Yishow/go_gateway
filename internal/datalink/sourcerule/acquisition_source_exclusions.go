package sourcerule

import (
	"cmp"
	"slices"
	"strings"

	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
)

const (
	acquisitionSourceOwnershipIncomplete     = "source-ownership-incomplete"
	acquisitionSourceRuleGeometryUnsupported = "source-rule-geometry-unsupported"
	acquisitionSourceSkipsUnsupported        = "source-skips-unsupported"
	acquisitionSourceLinksIncomplete         = "source-links-incomplete"
	acquisitionSourceLayoutUnsupported       = "source-layout-unsupported"
)

// AcquisitionSourceExclusions derives only source exclusions, not readability.
// A nonempty closed reason means geometry is unknown and cannot authorize gaps.
func AcquisitionSourceExclusions(input AcquisitionSourceRows) (ranges []optimizer.NativeRange, reason string) {
	protocol := acquisitionNativeProtocol(input.Device.Protocol)
	if protocol == "" {
		return nil, "source-protocol-unsupported"
	}
	identity := sourceAcquisitionIdentity(input.Device)
	rules := map[string]*schema.SourceRule{}
	points := map[string]*schema.Point{}
	links := map[string][]*schema.SourceRuleLink{}
	owners := map[string]string{}
	var excluded []optimizer.NativeRange
	for _, rule := range input.Rules {
		if rule == nil || rule.ID == "" || rule.DeviceID != input.Device.ID || rules[rule.ID] != nil {
			return nil, acquisitionSourceOwnershipIncomplete
		}
		rules[rule.ID] = rule
	}
	for _, point := range input.Points {
		if point == nil || point.ID == "" || point.DeviceID != input.Device.ID || points[point.ID] != nil {
			return nil, acquisitionSourceOwnershipIncomplete
		}
		span, err := optimizer.AcquisitionPointSpanForIdentity(identity, acquisitionEffectivePoint(point))
		if err != nil {
			return nil, "source-point-geometry-unsupported"
		}
		points[point.ID] = point
		if !point.Enabled {
			excluded = append(excluded, span)
		}
	}
	for _, link := range input.Links {
		if link == nil || rules[link.RuleID] == nil || points[link.PointID] == nil || owners[link.PointID] != "" || points[link.PointID].Address != link.Address {
			return nil, acquisitionSourceOwnershipIncomplete
		}
		owners[link.PointID] = link.RuleID
		links[link.RuleID] = append(links[link.RuleID], link)
	}
	mixedAllowed := acquisitionMixedLayoutSupported(input.Device)
	for _, rule := range input.Rules {
		spans, reason := acquisitionRuleExclusions(identity, input.Device.Protocol, mixedAllowed, rule, links[rule.ID], points)
		if reason != "" {
			return nil, reason
		}
		excluded = append(excluded, spans...)
	}
	slices.SortFunc(excluded, func(a, b optimizer.NativeRange) int {
		if n := cmp.Compare(a.Area, b.Area); n != 0 {
			return n
		}
		if n := cmp.Compare(a.Start, b.Start); n != 0 {
			return n
		}
		return cmp.Compare(a.End, b.End)
	})
	return slices.Compact(excluded), ""
}

func acquisitionRuleExclusions(identity optimizer.AcquisitionIdentity, protocol schema.ProtocolType, mixedAllowed bool, rule *schema.SourceRule, links []*schema.SourceRuleLink, points map[string]*schema.Point) (ranges []optimizer.NativeRange, reason string) {
	if (rule.LayoutMode != "" && rule.LayoutMode != "homogeneous") || strings.TrimSpace(rule.MixedItems) != "" {
		if !mixedAllowed || rule.LayoutMode != layoutMixed {
			return nil, acquisitionSourceLayoutUnsupported
		}
		return acquisitionMixedRuleExclusions(identity, rule, links, points)
	}
	// Validate the original address BEFORE the legacy rule expander: its decimal
	// helper is permissive and must not overflow into an apparently valid start.
	start, err := optimizer.AcquisitionPointSpanForIdentity(identity, optimizer.EffectivePoint{Address: rule.StartAddress, DataType: rule.DataType})
	if err != nil || rule.Count < 0 {
		return nil, acquisitionSourceRuleGeometryUnsupported
	}
	rawSkip := strings.TrimSpace(rule.SkippedAddresses)
	if rawSkip != "" && !strings.HasPrefix(rawSkip, "[") {
		return nil, acquisitionSourceSkipsUnsupported
	}
	skipped, err := parseSkippedAddresses(rule.SkippedAddresses)
	if err != nil {
		return nil, acquisitionSourceSkipsUnsupported
	}
	excluded := make([]optimizer.NativeRange, 0, len(skipped))
	for _, address := range skipped {
		span, err := optimizer.AcquisitionPointSpanForIdentity(identity, optimizer.EffectivePoint{Address: address, DataType: rule.DataType})
		if err != nil {
			return nil, acquisitionSourceSkipsUnsupported
		}
		excluded = append(excluded, span)
	}
	// Missing cells cannot become an authorization for an unrepresented gap.
	// Bound expansion by the actual enumerated input, then check address math.
	if rule.Count > len(links)+len(skipped) {
		return nil, acquisitionSourceLinksIncomplete
	}
	step := getDataTypeCellSpan(rule.DataType)
	if rule.Count > 0 && rule.Count-1 > (int(^uint(0)>>1)-start.End)/step {
		return nil, acquisitionSourceRuleGeometryUnsupported
	}
	if rule.Count > 0 {
		last, err := offsetAddressForIdentity(rule.StartAddress, (rule.Count-1)*step, protocol, identity)
		if err != nil {
			return nil, acquisitionSourceRuleGeometryUnsupported
		}
		if _, err = optimizer.AcquisitionPointSpanForIdentity(identity, optimizer.EffectivePoint{Address: last, DataType: rule.DataType}); err != nil {
			return nil, acquisitionSourceRuleGeometryUnsupported
		}
	}
	addresses, err := buildPlannedPointAddressesForIdentity(rule.StartAddress, rule.Count, rule.DataType, protocol, identity)
	if err != nil {
		return nil, acquisitionSourceRuleGeometryUnsupported
	}
	if len(addresses) > 0 {
		first, err := optimizer.AcquisitionPointSpanForIdentity(identity, optimizer.EffectivePoint{Address: addresses[0], DataType: rule.DataType})
		if err != nil || first != start {
			return nil, acquisitionSourceRuleGeometryUnsupported
		}
	}
	expected := map[string]bool{}
	for _, address := range addresses {
		if !rule.Enabled {
			span, err := optimizer.AcquisitionPointSpanForIdentity(identity, optimizer.EffectivePoint{Address: address, DataType: rule.DataType})
			if err != nil {
				return nil, acquisitionSourceRuleGeometryUnsupported
			}
			excluded = append(excluded, span)
		}
		if !containsAddress(skipped, address) {
			expected[normalizeAddressKey(address)] = true
		}
	}
	for _, link := range links {
		key := normalizeAddressKey(link.Address)
		point := points[link.PointID]
		if !expected[key] || point == nil || point.DataType != rule.DataType || (!rule.Enabled && point.Enabled) {
			return nil, acquisitionSourceLinksIncomplete
		}
		want, err := optimizer.AcquisitionPointSpanForIdentity(identity, optimizer.EffectivePoint{Address: link.Address, DataType: rule.DataType})
		actual, actualErr := optimizer.AcquisitionPointSpanForIdentity(identity, acquisitionEffectivePoint(point))
		if err != nil || actualErr != nil || want != actual {
			return nil, acquisitionSourceLinksIncomplete
		}
		delete(expected, key)
	}
	if len(expected) != 0 {
		return nil, acquisitionSourceLinksIncomplete
	}
	return excluded, ""
}

func acquisitionEffectivePoint(point *schema.Point) optimizer.EffectivePoint {
	return optimizer.EffectivePoint{ID: point.ID, Address: point.Address, Function: point.Function, DataType: point.DataType}
}

func acquisitionNativeProtocol(protocol schema.ProtocolType) string {
	switch protocol {
	case schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP:
		return "modbus"
	case schema.ProtocolMC3E:
		return "mc"
	case schema.ProtocolFatekFBs:
		return "fatek"
	}
	return ""
}
