package workspace

import (
	"context"
	"fmt"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"
	"strings"
)

func (s *Service) downstreamReadinessIssues(ctx context.Context, deviceIDs []string, connectorID string, coverage writeGroupReadinessCoverage) ([]ReadinessIssue, error) {
	if s.readinessRules == nil || len(deviceIDs) == 0 {
		return nil, nil
	}

	rules, err := s.readinessRules.ListByDeviceIDs(ctx, deviceIDs)
	if err != nil {
		return nil, fmt.Errorf("list workspace source rules for readiness: %w", err)
	}
	issues := make([]ReadinessIssue, 0, len(rules))
	for _, rule := range rules {
		if rule == nil {
			continue
		}
		links, err := s.readinessRules.ListLinks(ctx, rule.ID)
		if err != nil {
			return nil, fmt.Errorf("list workspace source rule links for readiness %s: %w", rule.ID, err)
		}
		for _, link := range links {
			issue, ok, err := s.downstreamLinkIssue(ctx, connectorID, link, coverage)
			if err != nil {
				return nil, err
			}
			if ok {
				issues = append(issues, issue)
			}
		}
	}

	return issues, nil
}

func (s *Service) downstreamLinkIssue(ctx context.Context, connectorID string, link *schema.SourceRuleLink, coverage writeGroupReadinessCoverage) (ReadinessIssue, bool, error) {
	if link == nil {
		return ReadinessIssue{}, false, nil
	}

	pointID := strings.TrimSpace(link.PointID)
	if pointID == "" {
		return ReadinessIssue{
			Code:     "point-missing",
			Severity: ReadinessSeverityBlocking,
			Step:     ReadinessStep2,
			Scope:    strings.TrimSpace(link.RuleID),
			Message:  "source rule link is missing its derived point",
		}, true, nil
	}
	binding, err := s.resolveReadinessBinding(ctx, link)
	if err != nil {
		return ReadinessIssue{}, false, err
	}

	if !binding.tagValid {
		message := "derived point is missing its persisted tag"
		if binding.tagID != "" {
			message = "derived point is not bound to a persisted tag owned by this source rule"
		}
		return ReadinessIssue{
			Code:     "tag-missing",
			Severity: ReadinessSeverityBlocking,
			Step:     ReadinessStep3,
			Scope:    pointID,
			Message:  message,
		}, true, nil
	}

	if !binding.mappingValid {
		message := "derived point is missing its persisted mapping"
		if binding.mappingID != "" {
			message = "derived point is not bound to a persisted mapping owned by this source rule"
		}
		return ReadinessIssue{
			Code:     "mapping-missing",
			Severity: ReadinessSeverityBlocking,
			Step:     ReadinessStep3,
			Scope:    pointID,
			Message:  message,
		}, true, nil
	}

	if connectorID == "" || coverage.members[writeGroupReadinessMemberKey{connectorID: connectorID, pointID: pointID, tagID: binding.tagID}] {
		return ReadinessIssue{}, false, nil
	}
	if s.readinessMappings == nil {
		return ReadinessIssue{}, false, ErrReadinessUnavailable
	}

	rows, err := s.readinessMappings.List(ctx, dbtarget.TargetMappingListFilter{
		ConnectorID: &connectorID,
		TagID:       &binding.tagID,
	})
	if err != nil {
		return ReadinessIssue{}, false, fmt.Errorf("list database targets for readiness tag %s: %w", binding.tagID, err)
	}
	if len(rows) == 0 {
		return ReadinessIssue{
			Code:     "database-target-missing",
			Severity: ReadinessSeverityBlocking,
			Step:     ReadinessStep4,
			Scope:    pointID,
			Message:  "derived point is missing its persisted database target",
		}, true, nil
	}

	return ReadinessIssue{}, false, nil
}
