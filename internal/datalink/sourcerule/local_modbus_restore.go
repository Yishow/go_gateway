package sourcerule

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"go-gateway/internal/datalink/schema"
)

// RestoreLocalModbusMappingState rebuilds persisted Local Modbus mappings from
// source-rule candidate snapshots after a process restart.
func (s *Service) RestoreLocalModbusMappingState(ctx context.Context, writer LocalModbusMappingWriter) error {
	if writer == nil {
		return nil
	}

	rules, err := s.repo.List(ctx, ListFilter{})
	if err != nil {
		return fmt.Errorf("列出來源規則失敗: %w", err)
	}

	mappings, err := s.listRestorableLocalModbusMappings(ctx, rules)
	if err != nil {
		return err
	}
	for _, mappingRecord := range mappings {
		if err := writer.Upsert(ctx, mappingRecord); err != nil {
			return fmt.Errorf("還原 local modbus 映射失敗: %w", err)
		}
	}
	return nil
}

func (s *Service) listRestorableLocalModbusMappings(
	ctx context.Context,
	rules []*schema.SourceRule,
) ([]LocalModbusMappingRecord, error) {
	sort.Slice(rules, func(i, j int) bool {
		if rules[i] == nil {
			return false
		}
		if rules[j] == nil {
			return true
		}
		return rules[i].ID < rules[j].ID
	})

	result := make([]LocalModbusMappingRecord, 0)
	for _, rule := range rules {
		if rule == nil || strings.TrimSpace(rule.RevisionID) == "" {
			continue
		}

		snapshots, err := s.repo.ListCandidateSnapshots(ctx, rule.ID, rule.RevisionID)
		if err != nil {
			return nil, fmt.Errorf("列出 local modbus 候選快照失敗: %w", err)
		}

		localModbusSnapshot := currentLocalModbusSnapshot(snapshots)
		if localModbusSnapshot == nil {
			continue
		}

		candidates, err := decodeCandidatePayload[schema.SourceRuleLocalModbusOutputCandidate](localModbusSnapshot.Payload)
		if err != nil {
			return nil, fmt.Errorf("解析 local modbus 候選快照失敗: %w", err)
		}

		for _, candidate := range candidates {
			if candidate.TagID == nil || candidate.Register == nil || strings.TrimSpace(*candidate.TagID) == "" {
				continue
			}
			if candidate.Status == schema.SourceRuleLocalModbusOutputStatusOutOfSync {
				continue
			}
			if strings.TrimSpace(candidate.BlockingReason) != "" {
				continue
			}

			record := LocalModbusMappingRecord{
				MappingID: derefOptional(candidate.MappingID),
				TagID:     strings.TrimSpace(*candidate.TagID),
				Register:  *candidate.Register,
				DataType:  candidate.DataType,
			}
			if candidate.UpdatedAt != nil {
				record.UpdatedAt = candidate.UpdatedAt.UTC()
			}
			result = append(result, record)
		}
	}

	return result, nil
}

func currentLocalModbusSnapshot(snapshots []*schema.SourceRuleCandidateSnapshot) *schema.SourceRuleCandidateSnapshot {
	for _, snapshot := range snapshots {
		if snapshot == nil || snapshot.CandidateType != schema.SourceRuleCandidateTypeLocalModbusOutputs {
			continue
		}
		return snapshot
	}
	return nil
}
