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
	tagCandidates []schema.SourceRuleTagCandidate,
) ([]schema.SourceRuleLocalModbusOutputCandidate, error) {
	states, err := s.listEffectiveTagReviewStates(ctx, rule.ID, tagCandidates)
	if err != nil {
		return nil, err
	}
	mappingsByTagID, err := s.listLocalModbusMappingsByTagID(ctx)
	if err != nil {
		return nil, err
	}

	candidates := make([]schema.SourceRuleLocalModbusOutputCandidate, 0, len(tagCandidates))
	for _, tagCandidate := range tagCandidates {
		state, ok := states[tagCandidate.ID]
		if !ok || !state.Include {
			continue
		}

		missingTag := false
		candidate := schema.SourceRuleLocalModbusOutputCandidate{
			Address:        state.Address,
			PointID:        state.PointID,
			TagID:          cloneOptionalString(state.TagID),
			TagKey:         state.TagKey,
			DisplayName:    state.DisplayName,
			DataType:       state.DataType,
			RegisterCount:  localModbusRegisterCount(state.DataType),
			Status:         schema.SourceRuleLocalModbusOutputStatusDeferred,
			BlockingReason: "",
		}
		if s.tagSvc != nil && state.TagID != nil && strings.TrimSpace(*state.TagID) != "" {
			tagRecord, err := s.tagSvc.GetByID(ctx, *state.TagID)
			if err != nil {
				if !errors.Is(err, tag.ErrTagNotFound) {
					return nil, fmt.Errorf("取得 local modbus 候選標籤失敗: %w", err)
				}
				missingTag = true
			}
			if !missingTag {
				candidate.TagKey = tagRecord.Key
				candidate.DisplayName = tagRecord.DisplayName
				candidate.DataType = tagRecord.DataType
				candidate.RegisterCount = localModbusRegisterCount(tagRecord.DataType)
			} else {
				candidate.Status = schema.SourceRuleLocalModbusOutputStatusOutOfSync
				candidate.BlockingReason = fmt.Sprintf("%s: %s", localModbusLinkedTagNotFoundReason, *state.TagID)
			}
		}
		if candidate.TagID != nil && strings.TrimSpace(*candidate.TagID) != "" {
			if mappingRecord, ok := mappingsByTagID[strings.TrimSpace(*candidate.TagID)]; ok {
				candidate.Register = uint16Ptr(mappingRecord.Register)
				candidate.UpdatedAt = timePtr(mappingRecord.UpdatedAt)
				if mappingRecord.DataType != "" {
					candidate.DataType = mappingRecord.DataType
					candidate.RegisterCount = localModbusRegisterCount(mappingRecord.DataType)
				}
			}
		}

		candidate.Identity = buildLocalModbusOutputCandidateIdentity(rule.ID, state.Address, candidate.DataType, candidate.Register)
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
