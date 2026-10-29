package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/schema"
)

// Selection is only a read identity, never a locator or mutation capability.
type runtimeSourceSelection struct{ RuleID, RuleRevision, PointID, MappingID string }

func runtimeSourceRequestSelection(c *gin.Context, scope *runtimeview.Scope) (*runtimeSourceSelection, bool) {
	query, err := url.ParseQuery(c.Request.URL.RawQuery)
	keys := []string{"rule_id", "rule_revision", "point_id", "mapping_id"}
	present := false
	for _, key := range keys {
		_, exists := query[key]
		present = present || exists
	}
	if !present && err == nil {
		return nil, true
	}
	valid := err == nil && scope != nil
	for _, key := range keys {
		values, exists := query[key]
		if key == "mapping_id" && !exists {
			continue
		}
		valid = valid && len(values) == 1 && values[0] != "" && runtimeSafeID(values[0]) == values[0]
	}
	if !valid {
		renderRuntimeHandlerError(c, http.StatusBadRequest, "unavailable", false)
		return nil, false
	}
	return &runtimeSourceSelection{RuleID: query.Get("rule_id"), RuleRevision: query.Get("rule_revision"), PointID: query.Get("point_id"), MappingID: query.Get("mapping_id")}, true
}

// These records are an ephemeral comparison snapshot, not persisted or execution
// proof. Telemetry is excluded because collection updates it independently.
type runtimeSourceRecords struct {
	Rule    schema.SourceRule
	Link    schema.SourceRuleLink
	Point   schema.Point
	Polling *schema.PollingGroup
	Mapping *schema.Mapping
	Tag     *schema.Tag
}

func (h *RuntimeHandler) readRuntimeSourceRecords(ctx context.Context, scope runtimeview.Scope, selected runtimeSourceSelection) (*runtimeSourceRecords, error) {
	if h.sourceRuleReader == nil || h.pointSvc == nil {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	rules, err := h.sourceRuleReader.ListByDeviceIDs(ctx, []string{scope.DeviceID})
	if err != nil {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	var rule *schema.SourceRule
	for _, candidate := range rules {
		if candidate == nil || candidate.ID != selected.RuleID {
			continue
		}
		if rule != nil {
			return nil, runtimeEvidenceError(runtimeview.ReasonStale)
		}
		copy := *candidate
		rule = &copy
	}
	if rule == nil || rule.RevisionID != selected.RuleRevision {
		return nil, runtimeEvidenceError(runtimeview.ReasonStale)
	}
	if rule.DeviceID != scope.DeviceID {
		return nil, runtimeEvidenceError(runtimeview.ReasonForeign)
	}
	links, err := h.sourceRuleReader.ListLinks(ctx, selected.RuleID)
	if err != nil {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	var link *schema.SourceRuleLink
	for _, candidate := range links {
		if candidate == nil || candidate.PointID != selected.PointID {
			continue
		}
		if candidate.RuleID != selected.RuleID || link != nil {
			return nil, runtimeEvidenceError(runtimeview.ReasonStale)
		}
		copy := *candidate
		link = &copy
	}
	if link == nil {
		return nil, runtimeEvidenceError(runtimeview.ReasonStale)
	}
	p, err := h.pointSvc.GetByID(ctx, selected.PointID)
	if err != nil || p == nil {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	if p.ID != selected.PointID || p.DeviceID != scope.DeviceID {
		return nil, runtimeEvidenceError(runtimeview.ReasonForeign)
	}
	if p.Address != link.Address {
		return nil, runtimeEvidenceError(runtimeview.ReasonStale)
	}
	records := &runtimeSourceRecords{Rule: *rule, Link: *link, Point: *p}
	records.Point.LastReadAt, records.Point.LastValue = nil, nil
	records.Point.LastError = ""
	records.Point.UpdatedAt = time.Time{}
	if p.PollingGroupID != nil && *p.PollingGroupID != "" {
		if h.groupSvc == nil {
			return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
		}
		group, readErr := h.groupSvc.GetByID(ctx, *p.PollingGroupID)
		if readErr != nil || group == nil {
			return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
		}
		if group.ID != *p.PollingGroupID {
			return nil, runtimeEvidenceError(runtimeview.ReasonStale)
		}
		copy := *group
		records.Polling = &copy
	}
	if selected.MappingID != "" {
		if err := h.readRuntimeSourceMapping(ctx, records, selected.MappingID); err != nil {
			return nil, err
		}
	}
	// Detach optional pointer values before a later repository observation can
	// mutate shared fixture/storage objects. JSON never leaves this read adapter.
	raw, err := json.Marshal(records)
	if err != nil {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	var detached runtimeSourceRecords
	if json.Unmarshal(raw, &detached) != nil {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	return &detached, nil
}

func (h *RuntimeHandler) readRuntimeSourceMapping(ctx context.Context, records *runtimeSourceRecords, mappingID string) error {
	if h.mappingReader == nil || h.tagReader == nil {
		return runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	if records.Link.MappingID != nil && *records.Link.MappingID != mappingID {
		return runtimeEvidenceError(runtimeview.ReasonStale)
	}
	m, err := h.mappingReader.GetByID(ctx, mappingID)
	if err != nil || m == nil {
		return runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	// A dependent Tag reader must not mutate an already observed mapping.
	mappingCopy := *m
	m = &mappingCopy
	if m.ID != mappingID || m.PointID != records.Point.ID || runtimeSafeID(m.TagID) == "" {
		return runtimeEvidenceError(runtimeview.ReasonStale)
	}
	if records.Link.TagID != nil && *records.Link.TagID != m.TagID {
		return runtimeEvidenceError(runtimeview.ReasonStale)
	}
	tag, err := h.tagReader.GetByID(ctx, m.TagID)
	if err != nil || tag == nil {
		return runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	if tag.ID != m.TagID || (!isRuleManagedWorkspaceTagOwnedBy(tag, records.Rule.ID, records.Link.Address) && !isLegacyDirectWorkspaceBinding(&records.Link, m, tag)) {
		return runtimeEvidenceError(runtimeview.ReasonStale)
	}
	tagCopy := *tag
	records.Mapping, records.Tag = m, &tagCopy
	return nil
}
