package sourcerule

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
)

func (s *Service) ensureRuleTagMapping(
	ctx context.Context,
	oldRule *schema.SourceRule,
	rule *schema.SourceRule,
	pointRecord *schema.Point,
	link *schema.SourceRuleLink,
	enabled bool,
	result *tagMappingSyncResult,
) (*schema.Tag, *schema.Mapping, error) {
	return s.ensureRuleTagMappingWithTarget(
		ctx,
		oldRule,
		rule,
		pointRecord,
		link,
		enabled,
		buildPointName(rule.NamingPrefix, link.Address),
		"",
		result,
	)
}

func (s *Service) ensureRuleTagMappingWithTarget(
	ctx context.Context,
	oldRule *schema.SourceRule,
	rule *schema.SourceRule,
	pointRecord *schema.Point,
	link *schema.SourceRuleLink,
	enabled bool,
	effectiveTagKey string,
	overrideTagID string,
	result *tagMappingSyncResult,
) (*schema.Tag, *schema.Mapping, error) {
	transformPipeline := s.buildRuleTransformPipeline(rule, pointRecord)
	targetDataType := desiredRuleTargetDataType(rule, pointRecord)

	pointID := pointRecord.ID
	pointMappings, err := s.mappingSvc.List(ctx, mapping.ListFilter{PointID: &pointID})
	if err != nil {
		return nil, nil, fmt.Errorf("查詢 point 既有映射失敗: %w", err)
	}
	if len(pointMappings) > 1 {
		return nil, nil, fmt.Errorf("point %s 存在多條映射，無法同步來源規則", pointRecord.ID)
	}
	if len(pointMappings) == 1 {
		tagRecord, getErr := s.tagSvc.GetByID(ctx, pointMappings[0].TagID)
		if getErr != nil {
			return nil, nil, fmt.Errorf("取得既有映射標籤失敗: %w", getErr)
		}
		tagRecord, syncErr := s.syncRuleManagedTagDataType(ctx, rule, link, tagRecord, targetDataType, result)
		if syncErr != nil {
			return nil, nil, syncErr
		}
		mappingRecord, syncErr := s.syncRuleManagedMapping(ctx, oldRule, rule, pointRecord, link, pointMappings[0], enabled, transformPipeline, result)
		if syncErr != nil {
			return nil, nil, syncErr
		}
		return tagRecord, mappingRecord, nil
	}

	tagRecord, createdTag, err := s.resolveRuleTag(ctx, rule, pointRecord, link, effectiveTagKey, overrideTagID, result)
	if err != nil {
		return nil, nil, err
	}
	if createdTag {
		result.createdTagIDs = append(result.createdTagIDs, tagRecord.ID)
	}
	if err := s.validateTagAvailability(ctx, tagRecord, pointRecord.ID); err != nil {
		return nil, nil, err
	}

	tagID := tagRecord.ID
	tagMappings, err := s.mappingSvc.List(ctx, mapping.ListFilter{TagID: &tagID})
	if err != nil {
		return nil, nil, fmt.Errorf("查詢 tag 既有映射失敗: %w", err)
	}
	if len(tagMappings) > 1 {
		return nil, nil, fmt.Errorf("tag %s 存在多條映射，無法同步來源規則", tagRecord.Key)
	}
	if len(tagMappings) == 1 {
		if tagMappings[0].PointID != pointRecord.ID {
			return nil, nil, fmt.Errorf("tag %s 已綁定其他 point", tagRecord.Key)
		}
		mappingRecord, syncErr := s.syncRuleManagedMapping(ctx, oldRule, rule, pointRecord, link, tagMappings[0], enabled, transformPipeline, result)
		if syncErr != nil {
			return nil, nil, syncErr
		}
		return tagRecord, mappingRecord, nil
	}

	ruleCandidateID, proposedSignature, err := ruleManagedMappingMetadata(rule, pointRecord, link, transformPipeline)
	if err != nil {
		return nil, nil, err
	}
	mappingEnabled := enabled
	mappingStatus := schema.MappingStatusActive
	if !mappingEnabled {
		mappingStatus = schema.MappingStatusDraft
	}
	mappingRecord, err := s.mappingSvc.Create(ctx, mapping.CreateMappingRequest{
		PointID:              pointRecord.ID,
		TagID:                tagRecord.ID,
		Enabled:              &mappingEnabled,
		TransformPipeline:    transformPipeline,
		Status:               &mappingStatus,
		RuleCandidateID:      ruleCandidateID,
		ProposedSignature:    proposedSignature,
		LastAppliedSignature: proposedSignature,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("建立來源規則映射失敗: %w", err)
	}
	result.createdMappingIDs = append(result.createdMappingIDs, mappingRecord.ID)
	return tagRecord, mappingRecord, nil
}

func (s *Service) resolveRuleTag(
	ctx context.Context,
	rule *schema.SourceRule,
	pointRecord *schema.Point,
	link *schema.SourceRuleLink,
	effectiveTagKey string,
	overrideTagID string,
	result *tagMappingSyncResult,
) (*schema.Tag, bool, error) {
	targetDataType := desiredRuleTargetDataType(rule, pointRecord)

	if strings.TrimSpace(overrideTagID) != "" {
		tagRecord, err := s.tagSvc.GetByID(ctx, overrideTagID)
		if err != nil {
			return nil, false, fmt.Errorf("取得 override tag 失敗: %w", err)
		}
		tagRecord, err = s.syncRuleManagedTagDataType(ctx, rule, link, tagRecord, targetDataType, result)
		if err != nil {
			return nil, false, err
		}
		return tagRecord, false, nil
	}

	if link.TagID != nil && strings.TrimSpace(*link.TagID) != "" {
		tagRecord, err := s.tagSvc.GetByID(ctx, *link.TagID)
		if err == nil {
			tagRecord, err = s.syncRuleManagedTagDataType(ctx, rule, link, tagRecord, targetDataType, result)
			if err != nil {
				return nil, false, err
			}
			return tagRecord, false, nil
		}
		if !errors.Is(err, tag.ErrTagNotFound) {
			return nil, false, fmt.Errorf("取得來源規則既有標籤失敗: %w", err)
		}
	}

	expectedKey := strings.TrimSpace(effectiveTagKey)
	if expectedKey == "" {
		expectedKey = buildPointName(rule.NamingPrefix, link.Address)
	}
	tagRecord, err := s.tagSvc.GetByKey(ctx, expectedKey)
	if err == nil {
		if tagRecord.Status == schema.TagStatusRetired {
			return nil, false, fmt.Errorf("自動產生的標籤鍵已被 retired tag 佔用: %s", expectedKey)
		}
		tagRecord, err = s.syncRuleManagedTagDataType(ctx, rule, link, tagRecord, targetDataType, result)
		if err != nil {
			return nil, false, err
		}
		return tagRecord, false, nil
	}
	if !errors.Is(err, tag.ErrTagNotFound) {
		return nil, false, fmt.Errorf("查詢來源規則標籤失敗: %w", err)
	}

	created, err := s.tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         expectedKey,
		DisplayName: pointRecord.Name,
		DataType:    targetDataType,
		Labels: map[string]string{
			ruleManagedTagLabelSource:  ruleManagedTagLabelValue,
			ruleManagedTagLabelRuleID:  rule.ID,
			ruleManagedTagLabelAddress: normalizeAddressKey(link.Address),
		},
	})
	if err != nil {
		return nil, false, fmt.Errorf("建立來源規則標籤失敗: %w", err)
	}
	return created, true, nil
}
