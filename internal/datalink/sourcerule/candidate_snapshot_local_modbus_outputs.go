package sourcerule

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-gateway/internal/datalink/modbusshare"
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
	for candidateIndex, tagCandidate := range tagCandidates {
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
				if strings.TrimSpace(mappingRecord.MappingID) != "" {
					candidate.MappingID = stringPtr(mappingRecord.MappingID)
				}
				candidate.Register = uint16Ptr(mappingRecord.Register)
				candidate.UpdatedAt = timePtr(mappingRecord.UpdatedAt)
				if mappingRecord.DataType != "" {
					candidate.DataType = mappingRecord.DataType
					candidate.RegisterCount = localModbusRegisterCount(mappingRecord.DataType)
				}
			}
		}
		if rule.ShareEnabled && rule.ShareStartRegister != nil && rule.ShareStride != nil && candidate.TagID != nil && strings.TrimSpace(*candidate.TagID) != "" {
			register, registerErr := configuredLocalModbusRegister(rule, candidate.Address, candidateIndex, candidate.RegisterCount)
			if registerErr != nil {
				candidate.Status = schema.SourceRuleLocalModbusOutputStatusOutOfSync
				candidate.BlockingReason = registerErr.Error()
			} else {
				candidate.Register = register
				candidate.Status = schema.SourceRuleLocalModbusOutputStatusReady
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

func configuredLocalModbusRegister(rule *schema.SourceRule, candidateAddress string, candidateIndex, registerSpan int) (*uint16, error) {
	if rule == nil || rule.ShareStartRegister == nil || rule.ShareStride == nil {
		return nil, fmt.Errorf("local modbus share geometry is not configured")
	}
	startRegister := *rule.ShareStartRegister
	if startRegister < 0 || startRegister > 4294967295 {
		return nil, fmt.Errorf("source-rule share start register is invalid")
	}
	base, err := modbusshare.HumanToZeroBased(uint32(startRegister))
	if err != nil {
		return nil, fmt.Errorf("source-rule share start register is invalid: %w", err)
	}
	stride := *rule.ShareStride
	if stride < registerSpan {
		return nil, &modbusshare.Error{
			Code:      modbusshare.ErrCodeInvalidGeometry,
			Message:   fmt.Sprintf("source-rule share stride %d is smaller than datatype span %d", stride, registerSpan),
			Retryable: false,
			Action:    localModbusStrideAction,
		}
	}
	offset := candidateIndex
	if start, startErr := strconv.Atoi(strings.TrimSpace(rule.StartAddress)); startErr == nil {
		if address, addressErr := strconv.Atoi(strings.TrimSpace(candidateAddress)); addressErr == nil && address >= start {
			offset = address - start
		}
	}
	zeroBased := int(base) + offset*stride
	if zeroBased < 0 || zeroBased > int(modbusshare.MaxRegisterIndex) {
		return nil, fmt.Errorf("source-rule share register exceeds capacity")
	}
	return uint16Ptr(uint16(zeroBased)), nil
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
