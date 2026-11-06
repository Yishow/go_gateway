package measurement

import (
	"encoding/json"
	"errors"
	"strings"

	"go-gateway/internal/datalink/schema"
)

// SourceRulePointFormat returns the rule default for a particular point. Mixed
// rules own the item format; point-level overrides still take precedence when
// the caller resolves the point's format intent.
func SourceRulePointFormat(rule *schema.SourceRule, point *schema.Point) (string, error) {
	if rule == nil || point == nil {
		return "", errors.New("invalid source rule point")
	}
	if !strings.EqualFold(strings.TrimSpace(rule.LayoutMode), string(LayoutModeMixed)) {
		return rule.DataFormat, nil
	}
	var items []MixedItem
	if json.Unmarshal([]byte(rule.MixedItems), &items) != nil {
		return "", errors.New("invalid mixed source items")
	}
	address := strings.ToUpper(strings.TrimSpace(point.Address))
	format := rule.DataFormat
	found := false
	for _, item := range items {
		if strings.ToUpper(strings.TrimSpace(item.Address)) != address {
			continue
		}
		if found || strings.TrimSpace(item.ItemID) == "" || item.DataType != point.DataType {
			return "", errors.New("ambiguous mixed source point")
		}
		found = true
		if item.DataFormat != "" {
			format = item.DataFormat
		}
	}
	if !found {
		return "", errors.New("mixed source point not found")
	}
	return format, nil
}
