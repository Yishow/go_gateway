package sourcerule

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"go-gateway/internal/datalink/schema"
)

const (
	localModbusOutputsDeferredReason = "local modbus output review flow is deferred until the output phase is configured"
)

type candidateSnapshotPayload struct {
	Candidates any `json:"candidates"`
}

func (s *Service) ListCandidateSnapshots(ctx context.Context, ruleID string) ([]*schema.SourceRuleCandidateSnapshot, error) {
	rule, err := s.repo.GetByID(ctx, ruleID)
	if err != nil {
		return nil, fmt.Errorf("取得來源規則失敗: %w", err)
	}
	return s.ListCandidateSnapshotsByRevision(ctx, rule.ID, rule.RevisionID)
}

func (s *Service) ListCandidateSnapshotsByRevision(ctx context.Context, ruleID, revisionID string) ([]*schema.SourceRuleCandidateSnapshot, error) {
	snapshots, err := s.repo.ListCandidateSnapshots(ctx, ruleID, revisionID)
	if err != nil {
		return nil, fmt.Errorf("列出來源規則候選快照失敗: %w", err)
	}
	return snapshots, nil
}

func (s *Service) persistCandidateSnapshots(ctx context.Context, rule *schema.SourceRule, links []*schema.SourceRuleLink) error {
	snapshots, err := s.buildCandidateSnapshots(ctx, rule, links)
	if err != nil {
		return err
	}
	if err := s.repo.ReplaceCandidateSnapshots(ctx, snapshots); err != nil {
		return fmt.Errorf("儲存來源規則候選快照失敗: %w", err)
	}
	tagCandidates, err := decodeCurrentTagCandidates(snapshots)
	if err != nil {
		return err
	}
	if err := s.markStaleTagReviewDecisions(ctx, rule, tagCandidates); err != nil {
		return err
	}
	if err := s.refreshLocalModbusConflictSnapshots(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Service) buildCandidateSnapshots(ctx context.Context, rule *schema.SourceRule, links []*schema.SourceRuleLink) ([]*schema.SourceRuleCandidateSnapshot, error) {
	if rule == nil {
		return nil, fmt.Errorf("source rule is nil")
	}
	if strings.TrimSpace(rule.RevisionID) == "" {
		return nil, fmt.Errorf("source rule revision id is empty")
	}

	tagCandidates, err := s.buildTagCandidates(ctx, rule, links)
	if err != nil {
		return nil, err
	}
	tagPayload, err := marshalCandidateSnapshotPayload(tagCandidates)
	if err != nil {
		return nil, err
	}
	databaseOutputCandidates, databaseOutputStatus, databaseOutputReason, err := s.buildDatabaseOutputCandidates(ctx, rule, tagCandidates)
	if err != nil {
		return nil, err
	}
	databaseOutputPayload, err := marshalCandidateSnapshotPayload(databaseOutputCandidates)
	if err != nil {
		return nil, err
	}
	localModbusOutputCandidates, err := s.buildLocalModbusOutputCandidates(ctx, rule, tagCandidates)
	if err != nil {
		return nil, err
	}
	localModbusOutputPayload, err := marshalCandidateSnapshotPayload(localModbusOutputCandidates)
	if err != nil {
		return nil, err
	}

	generatedAt := time.Now().UTC()
	return []*schema.SourceRuleCandidateSnapshot{
		{
			SourceRuleID:  rule.ID,
			RevisionID:    rule.RevisionID,
			CandidateType: schema.SourceRuleCandidateTypeTags,
			Payload:       tagPayload,
			Status:        schema.SourceRuleCandidateStatusReady,
			GeneratedAt:   generatedAt,
		},
		{
			SourceRuleID:  rule.ID,
			RevisionID:    rule.RevisionID,
			CandidateType: schema.SourceRuleCandidateTypeDatabaseOutputs,
			Payload:       databaseOutputPayload,
			Status:        databaseOutputStatus,
			Reason:        databaseOutputReason,
			GeneratedAt:   generatedAt,
		},
		{
			SourceRuleID:  rule.ID,
			RevisionID:    rule.RevisionID,
			CandidateType: schema.SourceRuleCandidateTypeLocalModbusOutputs,
			Payload:       localModbusOutputPayload,
			Status:        schema.SourceRuleCandidateStatusDeferred,
			Reason:        localModbusOutputsDeferredReason,
			GeneratedAt:   generatedAt,
		},
	}, nil
}

func (s *Service) buildTagCandidates(ctx context.Context, rule *schema.SourceRule, links []*schema.SourceRuleLink) ([]schema.SourceRuleTagCandidate, error) {
	candidates := make([]schema.SourceRuleTagCandidate, 0, len(links))
	for _, link := range links {
		if link == nil {
			continue
		}

		pointRecord, err := s.pointSvc.GetByID(ctx, link.PointID)
		if err != nil {
			return nil, fmt.Errorf("取得來源規則候選點位失敗: %w", err)
		}

		tagKey := buildPointName(rule.NamingPrefix, link.Address)
		displayName := pointRecord.Name
		dataType := desiredRuleTargetDataType(rule, pointRecord)
		transformPipeline := s.buildRuleTransformPipeline(rule, pointRecord)
		tagID := cloneOptionalString(link.TagID)
		mappingID := cloneOptionalString(link.MappingID)
		candidate := schema.SourceRuleTagCandidate{
			Identity:          buildTagCandidateIdentity(rule.ID, link.Address, dataType),
			Address:           link.Address,
			PointID:           link.PointID,
			TagID:             tagID,
			MappingID:         mappingID,
			TagKey:            tagKey,
			DisplayName:       displayName,
			DataType:          dataType,
			TransformPipeline: transformPipeline,
		}
		candidate.ID, err = candidateID(candidate.Identity)
		if err != nil {
			return nil, err
		}
		candidate.ProposedSignature, err = tagCandidateSignature(candidate)
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

func marshalCandidateSnapshotPayload(candidates any) (string, error) {
	payload, err := json.Marshal(candidateSnapshotPayload{Candidates: candidates})
	if err != nil {
		return "", fmt.Errorf("序列化來源規則候選快照失敗: %w", err)
	}
	return string(payload), nil
}

func cloneOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	return stringPtr(*value)
}

func cloneOptionalInt(value *int) *int {
	if value == nil {
		return nil
	}
	normalized := *value
	return &normalized
}

func decodeCurrentTagCandidates(snapshots []*schema.SourceRuleCandidateSnapshot) ([]schema.SourceRuleTagCandidate, error) {
	for _, snapshot := range snapshots {
		if snapshot == nil || snapshot.CandidateType != schema.SourceRuleCandidateTypeTags {
			continue
		}
		candidates, err := decodeCandidatePayload[schema.SourceRuleTagCandidate](snapshot.Payload)
		if err != nil {
			return nil, fmt.Errorf("解析來源規則 tag 候選快照失敗: %w", err)
		}
		return candidates, nil
	}
	return []schema.SourceRuleTagCandidate{}, nil
}
