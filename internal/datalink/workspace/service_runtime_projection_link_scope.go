package workspace

import "go-gateway/internal/datalink/schema"

func runtimeProjectionLinkPointIDs(links []*schema.SourceRuleLink) map[string]struct{} {
	result := make(map[string]struct{}, len(links))
	for _, link := range links {
		if link == nil || link.PointID == "" {
			continue
		}
		result[link.PointID] = struct{}{}
	}
	return result
}

func runtimeProjectionLinkMappingIDs(links []*schema.SourceRuleLink) map[string]struct{} {
	result := make(map[string]struct{}, len(links))
	for _, link := range links {
		if link == nil || link.MappingID == nil || *link.MappingID == "" {
			continue
		}
		result[*link.MappingID] = struct{}{}
	}
	return result
}
