package dbtarget

import (
	"context"
	"fmt"
)

func (s *MappingService) Validate(ctx context.Context, connectorID string) (*ValidationResult, error) {
	connector, err := s.connectorRepo.GetByID(ctx, connectorID)
	if err != nil {
		return nil, fmt.Errorf("取得資料庫連接器失敗: %w", err)
	}

	result := &ValidationResult{Ready: true}
	if !connector.Enabled {
		result.Issues = append(result.Issues, ValidationIssue{
			Severity: "warning",
			Code:     "connector_disabled",
			Message:  "資料庫連接器目前為停用狀態",
		})
	}

	// 先取得映射才知道要檢查哪些 schema：映射的 table_schema 可以指向
	// 連接器自身資料庫以外的位置。
	projection, err := listLiveTargetProjection(ctx, s.repo, TargetMappingListFilter{ConnectorID: &connectorID}, s.tagService)
	if err != nil {
		return nil, err
	}

	tables, tableErr := inspectTablesForSchemas(ctx, connector, mappingSchemaNames(projection.Mappings))
	if tableErr != nil {
		result.Ready = false
		result.Issues = append(result.Issues, ValidationIssue{
			Severity: "error",
			Code:     "connector_unreachable",
			Message:  tableErr.Error(),
		})
		return result, nil
	}
	if len(projection.Mappings) == 0 {
		result.Issues = append(result.Issues, ValidationIssue{
			Severity: "warning",
			Code:     "mapping_missing",
			Message:  "此資料庫連接器尚未設定任何輸出映射",
		})
	}

	for _, mapping := range projection.Mappings {
		tagEntity := projection.TagsByMappingID[mapping.ID]
		if tagEntity == nil {
			var err error
			tagEntity, err = s.tagService.GetByID(ctx, mapping.TagID)
			if err != nil {
				result.Ready = false
				result.Issues = append(result.Issues, ValidationIssue{
					Severity:  "error",
					MappingID: mapping.ID,
					TagID:     mapping.TagID,
					Code:      "tag_missing",
					Message:   fmt.Sprintf("標籤不存在: %s", mapping.TagID),
				})
				continue
			}
		}

		issues := validateMappingAgainstTables(connector.Kind, *mapping, *tagEntity, tables)
		for _, issue := range issues {
			if issue.Severity == "error" {
				result.Ready = false
			}
			result.Issues = append(result.Issues, issue)
		}
	}

	return result, nil
}
