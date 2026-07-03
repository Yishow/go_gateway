package workspace

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"
)

func (s *Service) databaseRowGroupReadinessIssues(ctx context.Context, record *Record, connector *schema.DatabaseConnector) ([]ReadinessIssue, error) {
	if record == nil || connector == nil || databaseConnectorWriteMode(connector) != schema.DatabaseWriteModeUpsert {
		return nil, nil
	}
	if len(record.DatabaseRowGroups) == 0 || s.readinessRules == nil || s.readinessMappings == nil {
		return nil, nil
	}

	pointTags, err := s.readinessPointTags(ctx, record.OrderedDeviceIDs)
	if err != nil {
		return nil, err
	}
	refs := databaseTargetRefsByPoint(record.DatabaseTargetRefs, record.DatabaseRowGroups)
	issues := make([]ReadinessIssue, 0)
	for _, group := range record.DatabaseRowGroups {
		if len(group.UniqueKeyColumns) > 0 {
			continue
		}
		reused, err := s.rowGroupReusesColumn(ctx, connector.ID, group, pointTags, refs)
		if err != nil {
			return nil, err
		}
		if reused {
			issues = append(issues, ReadinessIssue{
				Code:     "database-row-group-upsert-unsafe",
				Severity: ReadinessSeverityBlocking,
				Step:     ReadinessStep4,
				Scope:    group.ID,
				Message:  "row group reuses database columns under upsert without a stable uniqueness key",
			})
		}
	}
	return issues, nil
}

func (s *Service) readinessPointTags(ctx context.Context, deviceIDs []string) (map[string]string, error) {
	rules, err := s.readinessRules.ListByDeviceIDs(ctx, deviceIDs)
	if err != nil {
		return nil, fmt.Errorf("list workspace source rules for row group readiness: %w", err)
	}
	out := map[string]string{}
	for _, rule := range rules {
		if rule == nil {
			continue
		}
		links, err := s.readinessRules.ListLinks(ctx, rule.ID)
		if err != nil {
			return nil, fmt.Errorf("list workspace source rule links for row group readiness %s: %w", rule.ID, err)
		}
		for _, link := range links {
			if link != nil && link.TagID != nil {
				out[strings.TrimSpace(link.PointID)] = strings.TrimSpace(*link.TagID)
			}
		}
	}
	return out, nil
}

func (s *Service) rowGroupReusesColumn(ctx context.Context, connectorID string, group DatabaseRowGroup, pointTags map[string]string, refs map[string]string) (bool, error) {
	counts := map[string]int{}
	for _, pointID := range group.MemberPointIDs {
		if refs[pointID] != group.ID {
			continue
		}
		tagID := pointTags[pointID]
		if tagID == "" {
			continue
		}
		rows, err := s.readinessMappings.List(ctx, dbtarget.TargetMappingListFilter{
			ConnectorID: &connectorID,
			TagID:       &tagID,
		})
		if err != nil {
			return false, fmt.Errorf("list row group database targets for readiness tag %s: %w", tagID, err)
		}
		for _, row := range rows {
			if row != nil && row.Enabled {
				counts[row.ColumnName]++
			}
		}
	}
	for _, count := range counts {
		if count > 1 {
			return true, nil
		}
	}
	return false, nil
}

func databaseTargetRefsByPoint(refs []DatabaseTargetRef, groups []DatabaseRowGroup) map[string]string {
	out := make(map[string]string, len(refs))
	for _, ref := range refs {
		pointID := strings.TrimSpace(ref.PointID)
		rowGroupID := strings.TrimSpace(ref.RowGroupID)
		if databaseRowGroupContainsPoint(groups, rowGroupID, pointID) {
			out[pointID] = rowGroupID
		}
	}
	return out
}

func databaseConnectorWriteMode(connector *schema.DatabaseConnector) schema.DatabaseWriteMode {
	var config map[string]any
	if err := json.Unmarshal([]byte(connector.ConnectionConfig), &config); err != nil {
		return schema.DatabaseWriteModeInsert
	}
	if value, ok := config["write_mode"].(string); ok && strings.TrimSpace(value) != "" {
		return schema.DatabaseWriteMode(value)
	}
	return schema.DatabaseWriteModeInsert
}
