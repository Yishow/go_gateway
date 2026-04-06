package sourcerule

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
)

const localModbusLinkedTagNotFoundReason = "linked tag not found"

func (s *Service) buildLocalModbusOutputCandidates(
	ctx context.Context,
	rule *schema.SourceRule,
	links []*schema.SourceRuleLink,
) ([]schema.SourceRuleLocalModbusOutputCandidate, error) {
	if s.tagSvc == nil {
		return []schema.SourceRuleLocalModbusOutputCandidate{}, nil
	}

	mappingsByTagID, err := s.listLocalModbusMappingsByTagID(ctx)
	if err != nil {
		return nil, err
	}

	candidates := make([]schema.SourceRuleLocalModbusOutputCandidate, 0, len(links))
	for _, link := range links {
		if link == nil || link.TagID == nil || strings.TrimSpace(*link.TagID) == "" {
			continue
		}

		pointRecord, err := s.pointSvc.GetByID(ctx, link.PointID)
		if err != nil {
			return nil, fmt.Errorf("取得 local modbus 候選點位失敗: %w", err)
		}

		missingTag := false
		tagRecord, err := s.tagSvc.GetByID(ctx, *link.TagID)
		if err != nil {
			if !errors.Is(err, tag.ErrTagNotFound) {
				return nil, fmt.Errorf("取得 local modbus 候選標籤失敗: %w", err)
			}
			missingTag = true
		}

		candidate := schema.SourceRuleLocalModbusOutputCandidate{
			Address:        link.Address,
			PointID:        link.PointID,
			TagID:          cloneOptionalString(link.TagID),
			TagKey:         buildPointName(rule.NamingPrefix, link.Address),
			DisplayName:    pointRecord.Name,
			DataType:       desiredRuleTargetDataType(rule, pointRecord),
			RegisterCount:  localModbusRegisterCount(desiredRuleTargetDataType(rule, pointRecord)),
			Status:         schema.SourceRuleLocalModbusOutputStatusDeferred,
			BlockingReason: "",
		}
		if !missingTag {
			candidate.TagKey = tagRecord.Key
			candidate.DisplayName = tagRecord.DisplayName
			candidate.DataType = tagRecord.DataType
			candidate.RegisterCount = localModbusRegisterCount(tagRecord.DataType)
		} else {
			candidate.Status = schema.SourceRuleLocalModbusOutputStatusOutOfSync
			candidate.BlockingReason = fmt.Sprintf("%s: %s", localModbusLinkedTagNotFoundReason, *link.TagID)
		}
		if mappingRecord, ok := mappingsByTagID[*link.TagID]; ok {
			candidate.Register = uint16Ptr(mappingRecord.Register)
			candidate.UpdatedAt = timePtr(mappingRecord.UpdatedAt)
			if mappingRecord.DataType != "" {
				candidate.DataType = mappingRecord.DataType
				candidate.RegisterCount = localModbusRegisterCount(mappingRecord.DataType)
			}
		}

		candidate.Identity = buildLocalModbusOutputCandidateIdentity(rule.ID, link.Address, candidate.DataType, candidate.Register)
		candidate.ID, err = candidateID(candidate.Identity)
		if err != nil {
			return nil, err
		}
		candidate.ProposedSignature, err = localModbusOutputCandidateSignature(candidate)
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, candidate)
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Address < candidates[j].Address
	})
	return candidates, nil
}

func (s *Service) listLocalModbusMappingsByTagID(ctx context.Context) (map[string]LocalModbusMappingRecord, error) {
	reader := s.localModbusMappingReader()
	if reader == nil {
		return map[string]LocalModbusMappingRecord{}, nil
	}

	mappings, err := reader.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("列出 local modbus 映射失敗: %w", err)
	}

	result := make(map[string]LocalModbusMappingRecord, len(mappings))
	for _, mappingRecord := range mappings {
		if strings.TrimSpace(mappingRecord.TagID) == "" {
			continue
		}
		result[strings.TrimSpace(mappingRecord.TagID)] = mappingRecord
	}
	return result, nil
}

func localModbusRegisterCount(dataType schema.DataType) int {
	switch dataType {
	case schema.DataTypeInt32, schema.DataTypeUint32, schema.DataTypeFloat32:
		return 2
	case schema.DataTypeInt64, schema.DataTypeUint64, schema.DataTypeFloat64:
		return 4
	default:
		return 1
	}
}

func uint16Ptr(value uint16) *uint16 {
	return &value
}

func timePtr(value time.Time) *time.Time {
	return &value
}
