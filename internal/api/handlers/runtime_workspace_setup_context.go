package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	mappingpkg "go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"
	sourcerulepkg "go-gateway/internal/datalink/sourcerule"
	tagpkg "go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"
)

// WithSetupContextServices wires the persisted Studio V2 setup readers used by
// runtime context. It intentionally reads saved state only and does not mutate runtime.
func (h *RuntimeHandler) WithSetupContextServices(
	ruleSvc *sourcerulepkg.Service,
	tagSvc *tagpkg.Service,
	mappingSvc *mappingpkg.Service,
	connectorSvc *dbtarget.ConnectorService,
	targetSvc *dbtarget.MappingService,
) *RuntimeHandler {
	if ruleSvc != nil {
		h.sourceRuleReader = ruleSvc
	}
	if tagSvc != nil {
		h.tagReader = tagSvc
	}
	if mappingSvc != nil {
		h.mappingReader = mappingSvc
	}
	if connectorSvc != nil {
		h.dbConnectorReader = connectorSvc
	}
	if targetSvc != nil {
		h.dbTargetReader = targetSvc
	}
	return h
}

func (h *RuntimeHandler) buildWorkspaceSetupContext(ctx context.Context, record *workspace.Record) (runtimeWorkspaceSetupContextResponse, error) {
	response := runtimeWorkspaceSetupContextResponse{
		SourceRules:     []runtimeWorkspaceSetupSourceRuleResponse{},
		Mappings:        []runtimeWorkspaceSetupMappingResponse{},
		DatabaseTargets: []runtimeWorkspaceSetupDatabaseTargetMapping{},
	}

	readiness, err := h.workspaceSvc.Readiness(ctx)
	if err == nil {
		response.ReadinessSummary = readiness
	} else if !errors.Is(err, workspace.ErrReadinessUnavailable) {
		return response, fmt.Errorf("load readiness summary: %w", err)
	}

	rules, err := h.listWorkspaceSetupSourceRules(ctx, record.OrderedDeviceIDs)
	if err != nil {
		return response, err
	}
	response.SourceRules = make([]runtimeWorkspaceSetupSourceRuleResponse, 0, len(rules))
	ruleByID := make(map[string]*schema.SourceRule, len(rules))
	for _, rule := range rules {
		if rule == nil {
			continue
		}
		ruleByID[rule.ID] = rule
		response.SourceRules = append(response.SourceRules, mapRuntimeSetupSourceRule(rule))
	}

	mappings, err := h.listWorkspaceSetupMappings(ctx, rules, ruleByID)
	if err != nil {
		return response, err
	}
	response.Mappings = mappings
	workspaceTagIDs := make(map[string]struct{}, len(mappings))
	for _, item := range mappings {
		tagID := strings.TrimSpace(item.TagID)
		if tagID == "" {
			continue
		}
		workspaceTagIDs[tagID] = struct{}{}
	}

	connectorID := strings.TrimSpace(record.DatabaseConnectorID)
	if connectorID == "" {
		return response, nil
	}

	connectorMissing := false
	if h.dbConnectorReader != nil {
		connector, err := h.dbConnectorReader.GetByID(ctx, connectorID)
		if err != nil {
			if isRuntimeSetupMissingConnectorError(err) {
				connectorMissing = true
			} else {
				return response, fmt.Errorf("load workspace database connector %s: %w", connectorID, err)
			}
		} else {
			config := buildWorkspaceDatabaseConfigResponse(record.ID, connector)
			response.DatabaseConfig = &config
		}
	}
	if connectorMissing {
		return response, nil
	}

	if h.dbTargetReader != nil {
		targets, err := h.dbTargetReader.List(ctx, dbtarget.TargetMappingListFilter{
			ConnectorID: &connectorID,
		})
		if err != nil {
			return response, fmt.Errorf("list workspace database targets: %w", err)
		}
		response.DatabaseTargets = make([]runtimeWorkspaceSetupDatabaseTargetMapping, 0, len(targets))
		for _, target := range targets {
			if target == nil {
				continue
			}
			if _, ok := workspaceTagIDs[strings.TrimSpace(target.TagID)]; !ok {
				continue
			}
			response.DatabaseTargets = append(response.DatabaseTargets, mapRuntimeSetupDatabaseTarget(target))
		}
	}

	return response, nil
}

func isRuntimeSetupMissingConnectorError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "資料庫連接器不存在")
}

func (h *RuntimeHandler) listWorkspaceSetupSourceRules(ctx context.Context, deviceIDs []string) ([]*schema.SourceRule, error) {
	if h.sourceRuleReader == nil || len(deviceIDs) == 0 {
		return nil, nil
	}
	rules, err := h.sourceRuleReader.ListByDeviceIDs(ctx, deviceIDs)
	if err != nil {
		return nil, fmt.Errorf("list workspace source rules: %w", err)
	}
	return rules, nil
}

func (h *RuntimeHandler) listWorkspaceSetupMappings(
	ctx context.Context,
	rules []*schema.SourceRule,
	ruleByID map[string]*schema.SourceRule,
) ([]runtimeWorkspaceSetupMappingResponse, error) {
	if h.sourceRuleReader == nil || len(rules) == 0 {
		return []runtimeWorkspaceSetupMappingResponse{}, nil
	}

	result := make([]runtimeWorkspaceSetupMappingResponse, 0)
	for _, rule := range rules {
		if rule == nil {
			continue
		}
		links, err := h.sourceRuleReader.ListLinks(ctx, rule.ID)
		if err != nil {
			return nil, fmt.Errorf("list workspace source rule links %s: %w", rule.ID, err)
		}
		for _, link := range links {
			if link == nil {
				continue
			}
			mappingResponse, err := h.mapRuntimeSetupMapping(ctx, link, ruleByID)
			if err != nil {
				return nil, err
			}
			result = append(result, mappingResponse)
		}
	}

	return result, nil
}

func (h *RuntimeHandler) mapRuntimeSetupMapping(
	ctx context.Context,
	link *schema.SourceRuleLink,
	ruleByID map[string]*schema.SourceRule,
) (runtimeWorkspaceSetupMappingResponse, error) {
	rule := ruleByID[link.RuleID]
	response := runtimeWorkspaceSetupMappingResponse{
		RuleID:  link.RuleID,
		PointID: link.PointID,
		Address: link.Address,
		Enabled: true,
	}
	if rule != nil {
		response.DeviceID = rule.DeviceID
		response.TargetType = rule.DataType
		response.Enabled = rule.Enabled
		if rule.TargetDataType != nil {
			response.TargetType = *rule.TargetDataType
		}
	}

	mappingRecord, tagRecord, err := h.resolveRuntimeSetupBinding(ctx, link)
	if err != nil {
		return response, err
	}
	if mappingRecord != nil {
		response.ID = mappingRecord.ID
		response.TagID = mappingRecord.TagID
		response.Enabled = mappingRecord.Enabled
		response.Status = mappingRecord.Status
	}
	if tagRecord == nil {
		tagRecord, err = h.resolveRuntimeSetupDirectTag(ctx, link)
		if err != nil {
			return response, err
		}
	}
	if tagRecord != nil {
		if response.TagID == "" {
			response.TagID = tagRecord.ID
		}
		response.TagKey = tagRecord.Key
		response.DisplayName = tagRecord.DisplayName
		response.Unit = tagRecord.Unit
		if response.TargetType == "" {
			response.TargetType = tagRecord.DataType
		}
	}

	return response, nil
}

func (h *RuntimeHandler) resolveRuntimeSetupBinding(ctx context.Context, link *schema.SourceRuleLink) (*schema.Mapping, *schema.Tag, error) {
	if link == nil {
		return nil, nil, nil
	}
	if mappingRecord, tagRecord, ok, err := h.resolveRuntimeSetupDirectMapping(ctx, link); err != nil {
		return nil, nil, err
	} else if ok {
		return mappingRecord, tagRecord, nil
	}
	pointID := strings.TrimSpace(link.PointID)
	if pointID == "" || h.mappingReader == nil {
		return nil, nil, nil
	}
	records, err := h.mappingReader.List(ctx, mappingpkg.ListFilter{PointID: &pointID})
	if err != nil {
		return nil, nil, fmt.Errorf("list workspace mappings for point %s: %w", pointID, err)
	}
	matches := make([]struct {
		mapping *schema.Mapping
		tag     *schema.Tag
	}, 0, len(records))
	for _, record := range records {
		mappingRecord, tagRecord, ok, matchErr := h.matchRuntimeSetupMapping(ctx, link, record)
		if matchErr != nil {
			return nil, nil, matchErr
		}
		if ok {
			matches = append(matches, struct {
				mapping *schema.Mapping
				tag     *schema.Tag
			}{mapping: mappingRecord, tag: tagRecord})
		}
	}
	if len(matches) != 1 {
		return nil, nil, nil
	}
	return matches[0].mapping, matches[0].tag, nil
}

func (h *RuntimeHandler) resolveRuntimeSetupDirectMapping(ctx context.Context, link *schema.SourceRuleLink) (*schema.Mapping, *schema.Tag, bool, error) {
	if link.MappingID == nil || strings.TrimSpace(*link.MappingID) == "" {
		return nil, nil, false, nil
	}
	mappingRecord, err := h.loadRuntimeSetupMapping(ctx, strings.TrimSpace(*link.MappingID))
	if err != nil {
		return nil, nil, false, err
	}
	if mappingRecord == nil {
		return nil, nil, false, nil
	}
	return h.matchRuntimeSetupMapping(ctx, link, mappingRecord)
}

func (h *RuntimeHandler) matchRuntimeSetupMapping(ctx context.Context, link *schema.SourceRuleLink, mappingRecord *schema.Mapping) (*schema.Mapping, *schema.Tag, bool, error) {
	if mappingRecord == nil {
		return nil, nil, false, nil
	}
	if strings.TrimSpace(mappingRecord.PointID) != strings.TrimSpace(link.PointID) {
		return nil, nil, false, nil
	}
	if strings.TrimSpace(mappingRecord.TagID) == "" {
		return nil, nil, false, nil
	}
	tagRecord, err := h.loadRuntimeSetupTag(ctx, mappingRecord.TagID)
	if err != nil {
		return nil, nil, false, err
	}
	if tagRecord == nil || (!isRuleManagedWorkspaceTagOwnedBy(tagRecord, link.RuleID, link.Address) && !isLegacyDirectWorkspaceBinding(link, mappingRecord, tagRecord)) {
		return nil, nil, false, nil
	}
	return mappingRecord, tagRecord, true, nil
}

func (h *RuntimeHandler) resolveRuntimeSetupDirectTag(ctx context.Context, link *schema.SourceRuleLink) (*schema.Tag, error) {
	if link == nil || link.TagID == nil || strings.TrimSpace(*link.TagID) == "" {
		return nil, nil
	}
	tagRecord, err := h.loadRuntimeSetupTag(ctx, strings.TrimSpace(*link.TagID))
	if err != nil {
		return nil, err
	}
	if tagRecord == nil || !isRuleManagedWorkspaceTagOwnedBy(tagRecord, link.RuleID, link.Address) {
		return nil, nil
	}
	return tagRecord, nil
}

func (h *RuntimeHandler) loadRuntimeSetupMapping(ctx context.Context, id string) (*schema.Mapping, error) {
	if h.mappingReader == nil {
		return nil, nil
	}
	mappingRecord, err := h.mappingReader.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, mappingpkg.ErrMappingNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("load workspace mapping %s: %w", id, err)
	}
	return mappingRecord, nil
}

func (h *RuntimeHandler) loadRuntimeSetupTag(ctx context.Context, id string) (*schema.Tag, error) {
	if h.tagReader == nil {
		return nil, nil
	}
	tagRecord, err := h.tagReader.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, tagpkg.ErrTagNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("load workspace tag %s: %w", id, err)
	}
	return tagRecord, nil
}

func mapRuntimeSetupSourceRule(rule *schema.SourceRule) runtimeWorkspaceSetupSourceRuleResponse {
	return runtimeWorkspaceSetupSourceRuleResponse{
		ID:             rule.ID,
		DeviceID:       rule.DeviceID,
		StartAddress:   rule.StartAddress,
		Count:          rule.Count,
		DataType:       rule.DataType,
		NamingPrefix:   rule.NamingPrefix,
		Enabled:        rule.Enabled,
		RevisionID:     rule.RevisionID,
		TargetDataType: rule.TargetDataType,
	}
}

func mapRuntimeSetupDatabaseTarget(target *schema.DatabaseTargetMapping) runtimeWorkspaceSetupDatabaseTargetMapping {
	return runtimeWorkspaceSetupDatabaseTargetMapping{
		ID:                   target.ID,
		TagID:                target.TagID,
		ConnectorID:          target.ConnectorID,
		TableSchema:          target.TableSchema,
		TableName:            target.TableName,
		ColumnName:           target.ColumnName,
		WriteMode:            target.WriteMode,
		TimestampColumn:      target.TimestampColumn,
		GroupKey:             target.GroupKey,
		WriteIntervalSeconds: target.WriteIntervalSeconds,
		Enabled:              target.Enabled,
	}
}
