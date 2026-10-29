package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/schema"
)

type runtimeSourcePageResponse struct {
	ContractVersion   int                                       `json:"contract_version"`
	WorkspaceID       string                                    `json:"workspace_id"`
	WorkspaceRevision string                                    `json:"workspace_revision"`
	DeviceID          string                                    `json:"device_id"`
	Kind              string                                    `json:"kind"`
	Query             string                                    `json:"query"`
	Limit             int                                       `json:"limit"`
	ListRevision      string                                    `json:"list_revision"`
	Rules             []runtimeWorkspaceSetupSourceRuleResponse `json:"rules"`
	Mappings          []runtimeWorkspaceSetupMappingResponse    `json:"mappings"`
	NextCursor        *string                                   `json:"next_cursor"`
}
type runtimeSourcePageSnapshot struct {
	Rules    []runtimeWorkspaceSetupSourceRuleResponse
	Mappings []runtimeWorkspaceSetupMappingResponse
	Revision string
}

// The list branch is separate from detail/provider reads: no group resolution,
// per-row detail request, readiness or workspace initialization takes place.
func (h *RuntimeHandler) tryRuntimeSourcePage(c *gin.Context, scope *runtimeview.Scope) bool {
	query, err := url.ParseQuery(c.Request.URL.RawQuery)
	present := false
	for _, key := range []string{"source_list", "q", "limit", "cursor"} {
		_, exists := query[key]
		present = present || exists
	}
	if !present {
		return false
	}
	request, valid := runtimeSourcePageRequestFromQuery(query, scope)
	if err != nil || !valid {
		renderRuntimeHandlerError(c, http.StatusBadRequest, "unavailable", false)
		return true
	}
	h.scopedRuntimeSourcePage(c, *scope, request)
	return true
}
func (h *RuntimeHandler) scopedRuntimeSourcePage(c *gin.Context, scope runtimeview.Scope, request runtimeSourcePageRequest) {
	ctx := c.Request.Context()
	if h.workspaceSvc == nil || h.deviceSvc == nil {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonUnavailable))
		return
	}
	record, err := h.workspaceSvc.ReadSnapshot(ctx)
	if err != nil {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonUnavailable))
		return
	}
	if err = validateRuntimeEvidenceScope(scope, record); err != nil {
		renderRuntimeEvidenceError(c, err)
		return
	}
	device, err := h.deviceSvc.GetByID(ctx, scope.DeviceID)
	if err != nil || device == nil {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonUnavailable))
		return
	}
	if device.ID != scope.DeviceID {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonForeign))
		return
	}
	var cursor runtimeSourcePageCursor
	if request.Cursor != "" {
		var valid bool
		cursor, valid = runtimeSourcePageDecode(request.Cursor)
		if !valid {
			renderRuntimeHandlerError(c, http.StatusBadRequest, "unavailable", false)
			return
		}
	}
	before, err := h.readRuntimeSourcePage(ctx, scope, request)
	if err != nil {
		renderRuntimeEvidenceError(c, err)
		return
	}
	page := runtimeSourcePageResponse{ContractVersion: 1, WorkspaceID: scope.WorkspaceID, WorkspaceRevision: scope.WorkspaceRevision, DeviceID: scope.DeviceID, Kind: request.Kind, Query: request.Query, Limit: request.Limit, ListRevision: before.Revision, Rules: []runtimeWorkspaceSetupSourceRuleResponse{}, Mappings: []runtimeWorkspaceSetupMappingResponse{}}
	total := len(before.Rules) + len(before.Mappings)
	start := 0
	if request.Cursor != "" {
		if cursor.Binding != runtimeSourcePageBinding(scope, request, cursor.Position) || cursor.Revision != before.Revision || cursor.Position >= total || cursor.Position%request.Limit != 0 {
			renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonStale))
			return
		}
		start = cursor.Position
	}
	end := min(start+request.Limit, total)
	if request.Kind == "rules" {
		page.Rules = before.Rules[start:end]
	} else {
		page.Mappings = before.Mappings[start:end]
	}
	if end < total {
		page.NextCursor = runtimeSourcePageNext(scope, request, before.Revision, end)
	}
	after, err := h.readRuntimeSourcePage(ctx, scope, request)
	if err != nil {
		renderRuntimeEvidenceError(c, err)
		return
	}
	if before.Revision != after.Revision {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonStale))
		return
	}
	final, err := h.workspaceSvc.ReadSnapshot(ctx)
	if err != nil {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonUnavailable))
		return
	}
	if !reflect.DeepEqual(record, final) {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonStale))
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: runtimeWorkspaceContextResponse{WorkspaceID: scope.WorkspaceID, WorkspaceRevision: scope.WorkspaceRevision, Devices: []runtimeWorkspaceContextDeviceEntry{}, SourcePage: &page}})
}
func (h *RuntimeHandler) readRuntimeSourcePage(ctx context.Context, scope runtimeview.Scope, request runtimeSourcePageRequest) (*runtimeSourcePageSnapshot, error) {
	if h.sourceRuleReader == nil {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	input, err := h.sourceRuleReader.ListByDeviceIDs(ctx, []string{scope.DeviceID})
	if err != nil {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	rules := make([]*schema.SourceRule, 0, len(input))
	seen := map[string]bool{}
	filtered := request.RuleID == ""
	for _, rule := range input {
		if rule == nil {
			return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
		}
		if rule.DeviceID != scope.DeviceID {
			continue
		}
		if runtimeSafeID(rule.ID) == "" || runtimeSafeID(rule.RevisionID) == "" || rule.Count < 0 {
			return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
		}
		if seen[rule.ID] {
			return nil, runtimeEvidenceError(runtimeview.ReasonStale)
		}
		seen[rule.ID] = true
		if request.RuleID != "" {
			if rule.ID != request.RuleID {
				continue
			}
			if rule.RevisionID != request.RuleRevision {
				return nil, runtimeEvidenceError(runtimeview.ReasonStale)
			}
			filtered = true
		}
		copy := *rule
		rules = append(rules, &copy)
	}
	if !filtered {
		return nil, runtimeEvidenceError(runtimeview.ReasonStale)
	}
	sort.Slice(rules, func(i, j int) bool { return rules[i].ID < rules[j].ID })
	// Freeze raw rule contents before dependent readers may change shared records.
	ruleSnapshot, err := json.Marshal(rules)
	if err != nil {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	result := &runtimeSourcePageSnapshot{Rules: []runtimeWorkspaceSetupSourceRuleResponse{}, Mappings: []runtimeWorkspaceSetupMappingResponse{}}
	var links []*schema.SourceRuleLink
	var points []*schema.Point
	if request.Kind == "rules" {
		for _, rule := range rules {
			result.Rules = append(result.Rules, mapRuntimeSetupSourceRule(rule))
		}
	} else {
		if h.pointSvc == nil {
			return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
		}
		points, err = h.pointSvc.ListByDevice(ctx, scope.DeviceID)
		if err != nil {
			return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
		}
		// Repository slices may be borrowed; never replace their elements.
		points = append([]*schema.Point(nil), points...)
		owned := map[string]*schema.Point{}
		for i, p := range points {
			if p == nil || p.DeviceID != scope.DeviceID || owned[p.ID] != nil {
				return nil, runtimeEvidenceError(runtimeview.ReasonForeign)
			}
			copy := *p
			copy.LastValue, copy.LastReadAt = nil, nil
			copy.LastError = ""
			copy.UpdatedAt = time.Time{}
			points[i] = &copy
			owned[p.ID] = &copy
		}
		sort.Slice(points, func(i, j int) bool { return points[i].ID < points[j].ID })
		for _, rule := range rules {
			rows, readErr := h.sourceRuleReader.ListLinks(ctx, rule.ID)
			if readErr != nil {
				return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
			}
			seenPoints := map[string]bool{}
			for _, link := range rows {
				if link == nil || link.RuleID != rule.ID || runtimeSafeID(link.PointID) == "" || seenPoints[link.PointID] {
					return nil, runtimeEvidenceError(runtimeview.ReasonStale)
				}
				seenPoints[link.PointID] = true
				p := owned[link.PointID]
				if p == nil {
					return nil, runtimeEvidenceError(runtimeview.ReasonForeign)
				}
				if p.Address != link.Address {
					return nil, runtimeEvidenceError(runtimeview.ReasonStale)
				}
				if (link.MappingID != nil && h.mappingReader == nil) || (link.TagID != nil && h.tagReader == nil) {
					return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
				}
				row, readErr := h.mapRuntimeSetupMapping(ctx, link, map[string]*schema.SourceRule{rule.ID: rule})
				if readErr != nil {
					return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
				}
				// Legacy summaries may recover a unique binding, but an explicit
				// link identity must never be replaced in the selectable page.
				if (link.MappingID != nil && (row.ID == "" || row.ID != *link.MappingID)) ||
					(link.TagID != nil && (row.TagID == "" || row.TagID != *link.TagID)) {
					return nil, runtimeEvidenceError(runtimeview.ReasonStale)
				}
				row.RuleRevision = rule.RevisionID
				result.Mappings = append(result.Mappings, row)
				copy := *link
				links = append(links, &copy)
			}
		}
		sort.Slice(result.Mappings, func(i, j int) bool {
			a, b := result.Mappings[i], result.Mappings[j]
			if a.RuleID != b.RuleID {
				return a.RuleID < b.RuleID
			}
			if a.PointID != b.PointID {
				return a.PointID < b.PointID
			}
			return a.ID < b.ID
		})
		sort.Slice(links, func(i, j int) bool {
			if links[i].RuleID != links[j].RuleID {
				return links[i].RuleID < links[j].RuleID
			}
			return links[i].PointID < links[j].PointID
		})
	}
	result.Revision, err = runtimeSourcePageHash([]any{json.RawMessage(ruleSnapshot), links, points, result.Rules, result.Mappings})
	if err != nil {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	query := strings.ToLower(request.Query)
	if query != "" {
		keptRules := result.Rules[:0]
		for _, row := range result.Rules {
			if runtimeSourcePageMatches(query, row.ID, row.StartAddress, row.NamingPrefix, string(row.DataType)) {
				keptRules = append(keptRules, row)
			}
		}
		result.Rules = keptRules
		keptMappings := result.Mappings[:0]
		for _, row := range result.Mappings {
			if runtimeSourcePageMatches(query, row.ID, row.RuleID, row.PointID, row.TagID, row.Address, row.TagKey, row.DisplayName, row.Unit, string(row.TargetType)) {
				keptMappings = append(keptMappings, row)
			}
		}
		result.Mappings = keptMappings
	}
	return result, nil
}
func runtimeSourcePageMatches(query string, fields ...string) bool {
	for _, field := range fields {
		if strings.Contains(strings.ToLower(field), query) {
			return true
		}
	}
	return false
}
