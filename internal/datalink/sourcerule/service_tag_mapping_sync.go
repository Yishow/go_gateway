package sourcerule

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

func (s *Service) syncRuleTagMappings(ctx context.Context, oldRule, rule *schema.SourceRule, links []*schema.SourceRuleLink, enabled bool, skipMappingIDs ...string) (tagMappingSyncResult, error) {
	result := tagMappingSyncResult{
		updatedMappings: make(map[string]mappingRollbackState),
		updatedTags:     make(map[string]tagRollbackState),
	}
	if s.tagSvc == nil || s.mappingSvc == nil || len(links) == 0 {
		return result, nil
	}

	// Review-first flow defers creation of rule-managed tags/mappings until explicit apply.
	// Existing applied links still need lifecycle updates when the source rule changes.
	for _, link := range links {
		if link == nil || link.TagID == nil || link.MappingID == nil {
			continue
		}
		if strings.TrimSpace(*link.TagID) == "" || strings.TrimSpace(*link.MappingID) == "" {
			continue
		}

		if len(skipMappingIDs) > 0 && *link.MappingID == skipMappingIDs[0] {
			continue
		}

		pointRecord, err := s.pointSvc.GetByID(ctx, link.PointID)
		if err != nil {
			return result, fmt.Errorf("取得衍生點位失敗: %w", err)
		}

		tagRecord, err := s.tagSvc.GetByID(ctx, *link.TagID)
		if err != nil {
			return result, fmt.Errorf("取得來源規則既有標籤失敗: %w", err)
		}

		mappingRecord, err := s.mappingSvc.GetByID(ctx, *link.MappingID)
		if err != nil {
			return result, fmt.Errorf("取得來源規則既有映射失敗: %w", err)
		}
		if mappingRecord.PointID != pointRecord.ID {
			return result, fmt.Errorf("來源規則連結映射 point 不一致: %s", mappingRecord.ID)
		}
		if mappingRecord.TagID != tagRecord.ID {
			return result, fmt.Errorf("來源規則連結映射 tag 不一致: %s", mappingRecord.ID)
		}

		expectedPipeline := s.buildRuleTransformPipeline(rule, pointRecord)
		expectedSignature, err := mappingCandidateSignature(expectedPipeline)
		if err != nil {
			return result, err
		}
		appliedPipeline, err := decodeTransformPipeline(mappingRecord.TransformPipeline)
		if err != nil {
			return result, err
		}
		appliedSignature, err := mappingCandidateSignature(appliedPipeline)
		if err != nil {
			return result, err
		}
		// The saved pipeline remains published until an explicit candidate reapply.
		// A different source proposal must not silently change that pipeline's Tag type.
		targetType := desiredRuleTargetDataType(rule, pointRecord)
		if expectedSignature != appliedSignature {
			targetType = tagRecord.DataType
		}
		tagRecord, err = s.syncRuleManagedTagDataType(ctx, rule, link, tagRecord, targetType, &result)
		if err != nil {
			return result, err
		}

		mappingRecord, err = s.syncRuleManagedMapping(
			ctx,
			oldRule,
			rule,
			pointRecord,
			link,
			mappingRecord,
			enabled,
			s.buildRuleTransformPipeline(rule, pointRecord),
			&result,
		)
		if err != nil {
			return result, err
		}

		link.TagID = common.Ptr(tagRecord.ID)
		link.MappingID = common.Ptr(mappingRecord.ID)
		link.UpdatedAt = time.Now()
	}

	return result, nil
}
