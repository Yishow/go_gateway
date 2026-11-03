package sourcerule

import (
	"context"
	"fmt"

	"go-gateway/internal/datalink/schema"
)

// TagCandidateSnapshotView contains only the persisted tag evidence for a rule revision.
type TagCandidateSnapshotView struct {
	SourceRuleID      string           `json:"source_rule_id"`
	WorkspaceID       string           `json:"workspace_id,omitempty"`
	WorkspaceRevision string           `json:"workspace_revision,omitempty"`
	RevisionID        string           `json:"revision_id"`
	Tags              CandidateSetView `json:"tags"`
}

// GetTagCandidateViewAtRevision reads tag evidence without refreshing database outputs.
func (s *Service) GetTagCandidateViewAtRevision(ctx context.Context, ruleID, revisionID string, scope *CandidateScopeRequest) (*TagCandidateSnapshotView, error) {
	ctx, release := s.AcquireRuleMutation(ctx, ruleID)
	defer release()
	if scope != nil && s.CandidateScopeConfigured() {
		if err := s.ValidateCandidateScope(ctx, ruleID, *scope); err != nil {
			return nil, err
		}
	}
	rule, err := s.repo.GetByID(ctx, ruleID)
	if err != nil {
		return nil, err
	}
	if rule.RevisionID != revisionID {
		return nil, candidateRevisionConflict()
	}
	snapshots, err := s.repo.ListCandidateSnapshots(ctx, ruleID, revisionID)
	if err != nil {
		return nil, err
	}
	current, err := s.repo.GetByID(ctx, ruleID)
	if err != nil {
		return nil, err
	}
	if current.RevisionID != revisionID {
		return nil, candidateRevisionConflict()
	}
	for _, snapshot := range snapshots {
		if snapshot == nil || snapshot.CandidateType != schema.SourceRuleCandidateTypeTags {
			continue
		}
		candidates, err := decodeCandidatePayload[schema.SourceRuleTagCandidate](snapshot.Payload)
		if err != nil {
			return nil, fmt.Errorf("解析來源規則 tag 候選快照失敗: %w", err)
		}
		views := make([]any, 0, len(candidates))
		for _, candidate := range candidates {
			view, err := s.buildTagCandidateView(ctx, candidate)
			if err != nil {
				return nil, err
			}
			views = append(views, view)
		}
		return &TagCandidateSnapshotView{SourceRuleID: ruleID, RevisionID: revisionID, Tags: CandidateSetView{Status: snapshot.Status, Reason: snapshot.Reason, Candidates: views}}, nil
	}
	return nil, fmt.Errorf("來源規則 %s revision %s 缺少 tags 候選快照", ruleID, revisionID)
}
