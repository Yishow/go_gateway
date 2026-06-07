package workspace

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"
)

func (s *Service) runtimeProjectionDatabaseTargets(
	ctx context.Context,
	connectorID string,
	tags []*schema.Tag,
) (*schema.DatabaseConnector, []*schema.DatabaseTargetMapping, error) {
	if connectorID == "" {
		return nil, nil, nil
	}

	var connector *schema.DatabaseConnector
	var err error
	if s.projectionConnectors != nil {
		connector, err = s.projectionConnectors.GetByID(ctx, connectorID)
		if err != nil {
			return nil, nil, fmt.Errorf("read workspace projection database connector %s: %w", connectorID, err)
		}
	}
	if s.projectionTargets == nil {
		return connector, nil, nil
	}

	enabled := true
	targets := make([]*schema.DatabaseTargetMapping, 0)
	for _, tagRecord := range tags {
		tagID := tagRecord.ID
		tagTargets, err := s.projectionTargets.List(ctx, dbtarget.TargetMappingListFilter{
			ConnectorID: &connectorID,
			TagID:       &tagID,
			Enabled:     &enabled,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("read workspace projection database targets %s: %w", tagID, err)
		}
		sort.SliceStable(tagTargets, func(i, j int) bool {
			return tagTargets[i].ID < tagTargets[j].ID
		})
		targets = append(targets, tagTargets...)
	}
	return connector, targets, nil
}

func (s *Service) runtimeProjectionPollingGroups(ctx context.Context, points []*schema.Point) ([]*schema.PollingGroup, error) {
	usedGroupIDs := make(map[string]struct{})
	for _, pointRecord := range points {
		if pointRecord.PollingGroupID == nil || *pointRecord.PollingGroupID == "" {
			continue
		}
		usedGroupIDs[*pointRecord.PollingGroupID] = struct{}{}
	}
	if len(usedGroupIDs) == 0 {
		return nil, nil
	}

	groups, err := s.projectionGroups.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("read workspace projection polling groups: %w", err)
	}
	filtered := make([]*schema.PollingGroup, 0, len(usedGroupIDs))
	for _, group := range groups {
		if _, ok := usedGroupIDs[group.ID]; ok {
			filtered = append(filtered, group)
		}
	}
	sort.SliceStable(filtered, func(i, j int) bool {
		return filtered[i].ID < filtered[j].ID
	})
	return filtered, nil
}

type runtimeProjectionVersionInput struct {
	WorkspaceID       string                          `json:"workspace_id"`
	Alignment         RuntimeProjectionAlignment      `json:"alignment"`
	DeviceIDs         []string                        `json:"device_ids"`
	Devices           []*schema.Device                `json:"devices"`
	Rules             []*schema.SourceRule            `json:"rules"`
	RuleLinks         []*schema.SourceRuleLink        `json:"rule_links"`
	Points            []*schema.Point                 `json:"points"`
	PollingGroups     []*schema.PollingGroup          `json:"polling_groups"`
	Mappings          []*schema.Mapping               `json:"mappings"`
	Tags              []*schema.Tag                   `json:"tags"`
	DatabaseConnector *schema.DatabaseConnector       `json:"database_connector,omitempty"`
	DatabaseTargets   []*schema.DatabaseTargetMapping `json:"database_targets"`
}

func runtimeProjectionVersion(projection *RuntimeProjection) (string, error) {
	input := runtimeProjectionVersionInput{
		WorkspaceID:       projection.WorkspaceID,
		Alignment:         projection.Alignment,
		DeviceIDs:         projection.DeviceIDs,
		Devices:           projection.Devices,
		Rules:             projection.Rules,
		RuleLinks:         projection.RuleLinks,
		Points:            projection.Points,
		PollingGroups:     projection.PollingGroups,
		Mappings:          projection.Mappings,
		Tags:              projection.Tags,
		DatabaseConnector: projection.DatabaseConnector,
		DatabaseTargets:   projection.DatabaseTargets,
	}
	encoded, err := json.Marshal(input)
	if err != nil {
		return "", fmt.Errorf("encode workspace runtime projection version: %w", err)
	}
	sum := sha256.Sum256(encoded)
	return fmt.Sprintf("sha256:%x", sum[:8]), nil
}
