package sourcerule

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
)

type UpsertTagReviewDecisionRequest struct {
	CandidateID   string                                   `json:"candidate_id"`
	Action        schema.SourceRuleTagReviewDecisionAction `json:"action"`
	TagKey        string                                   `json:"tag_key,omitempty"`
	OverrideTagID string                                   `json:"override_tag_id,omitempty"`
}

func (s *Service) UpsertTagReviewDecision(ctx context.Context, ruleID string, req UpsertTagReviewDecisionRequest) (*schema.SourceRuleTagReviewDecision, error) {
	rule, err := s.repo.GetByID(ctx, ruleID)
	if err != nil {
		return nil, fmt.Errorf("取得來源規則失敗: %w", err)
	}

	candidate, err := s.findCurrentTagCandidate(ctx, rule, req.CandidateID)
	if err != nil {
		return nil, err
	}

	decision := &schema.SourceRuleTagReviewDecision{
		SourceRuleID:    rule.ID,
		CandidateID:     candidate.ID,
		Action:          req.Action,
		Stale:           false,
		StaleRevisionID: "",
		StaleAt:         nil,
		UpdatedAt:       time.Now().UTC(),
	}

	switch req.Action {
	case schema.SourceRuleTagReviewDecisionActionRename:
		tagKey := strings.TrimSpace(req.TagKey)
		if err := tag.ValidateTagKey(tagKey); err != nil {
			return nil, fmt.Errorf("rename tag key 無效: %w", err)
		}
		decision.TagKey = tagKey
	case schema.SourceRuleTagReviewDecisionActionSkip:
		// No extra payload needed.
	case schema.SourceRuleTagReviewDecisionActionOverride:
		if s.tagSvc == nil {
			return nil, fmt.Errorf("tag 服務未配置，無法建立 override decision")
		}
		if strings.TrimSpace(req.OverrideTagID) == "" {
			return nil, fmt.Errorf("override_tag_id is required")
		}
		tagRecord, err := s.tagSvc.GetByID(ctx, req.OverrideTagID)
		if err != nil {
			return nil, fmt.Errorf("取得 override tag 失敗: %w", err)
		}
		decision.OverrideTagID = stringPtr(tagRecord.ID)
		decision.TagKey = tagRecord.Key
	default:
		return nil, fmt.Errorf("unsupported tag review decision action: %s", req.Action)
	}

	if err := s.repo.UpsertTagReviewDecision(ctx, decision); err != nil {
		return nil, err
	}
	return s.repo.GetTagReviewDecision(ctx, rule.ID, candidate.ID)
}

func (s *Service) GetTagReviewDecision(ctx context.Context, ruleID, candidateID string) (*schema.SourceRuleTagReviewDecision, error) {
	decision, err := s.repo.GetTagReviewDecision(ctx, ruleID, candidateID)
	if err != nil {
		return nil, err
	}
	return decision, nil
}

func (s *Service) ListTagReviewDecisions(ctx context.Context, ruleID string) ([]*schema.SourceRuleTagReviewDecision, error) {
	return s.repo.ListTagReviewDecisions(ctx, ruleID)
}

func (s *Service) findCurrentTagCandidate(ctx context.Context, rule *schema.SourceRule, candidateID string) (*schema.SourceRuleTagCandidate, error) {
	snapshots, err := s.repo.ListCandidateSnapshots(ctx, rule.ID, rule.RevisionID)
	if err != nil {
		return nil, fmt.Errorf("列出來源規則候選快照失敗: %w", err)
	}

	for _, snapshot := range snapshots {
		if snapshot == nil || snapshot.CandidateType != schema.SourceRuleCandidateTypeTags {
			continue
		}
		candidates, err := decodeCandidatePayload[schema.SourceRuleTagCandidate](snapshot.Payload)
		if err != nil {
			return nil, fmt.Errorf("解析來源規則 tag 候選快照失敗: %w", err)
		}
		for _, candidate := range candidates {
			if candidate.ID == candidateID {
				candidateCopy := candidate
				return &candidateCopy, nil
			}
		}
	}

	return nil, fmt.Errorf("來源規則 revision %s 不存在 tag candidate %s", rule.RevisionID, candidateID)
}
