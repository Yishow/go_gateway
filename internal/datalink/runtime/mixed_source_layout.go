package runtime

import (
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

// Legacy mixed points also need their rule layout attached to the scheduled
// profile; their legacy decoding policy and identity remain unchanged.
func projectionMixedPoint(projection *workspace.RuntimeProjection, point *schema.Point) bool {
	if projection == nil || point == nil {
		return false
	}
	for _, link := range projection.RuleLinks {
		if link == nil || link.PointID != point.ID {
			continue
		}
		for _, rule := range projection.Rules {
			if rule != nil && rule.ID == link.RuleID && rule.DeviceID == point.DeviceID && rule.LayoutMode == "mixed" {
				return true
			}
		}
	}
	return false
}
