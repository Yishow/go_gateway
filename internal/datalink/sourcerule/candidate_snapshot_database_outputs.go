package sourcerule

import (
	"context"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema"
)

type effectiveTagReviewState struct {
	Include     bool
	Address     string
	PointID     string
	TagID       *string
	TagKey      string
	DisplayName string
	DataType    schema.DataType
}

func (s *Service) buildDatabaseOutputCandidates(
	ctx context.Context,
	rule *schema.SourceRule,
	tagCandidates []schema.SourceRuleTagCandidate,
) ([]schema.SourceRuleDatabaseOutputCandidate, error) {
	states, err := s.listEffectiveTagReviewStates(ctx, rule.ID, tagCandidates)
	if err != nil {
		return nil, err
	}

	candidates := make([]schema.SourceRuleDatabaseOutputCandidate, 0, len(tagCandidates))
	for _, tagCandidate := range tagCandidates {
		state, ok := states[tagCandidate.ID]
		if !ok || !state.Include {
			continue
		}

		candidate := schema.SourceRuleDatabaseOutputCandidate{
			Identity:    buildDatabaseOutputCandidateIdentity(rule.ID, state.Address, state.DataType),
			Address:     state.Address,
			PointID:     state.PointID,
			TagID:       cloneOptionalString(state.TagID),
			TagKey:      state.TagKey,
			DisplayName: state.DisplayName,
			DataType:    state.DataType,
		}
		candidate.ID, err = candidateID(candidate.Identity)
		if err != nil {
			return nil, err
		}
		candidate.ProposedSignature, err = databaseOutputCandidateSignature(candidate)
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, candidate)
	}

	return candidates, nil
}

func (s *Service) listEffectiveTagReviewStates(
	ctx context.Context,
	ruleID string,
	tagCandidates []schema.SourceRuleTagCandidate,
) (map[string]effectiveTagReviewState, error) {
	decisionsByCandidateID, err := s.listCurrentTagReviewDecisionsByCandidateID(ctx, ruleID)
	if err != nil {
		return nil, fmt.Errorf("列出來源規則 tag 審查決策失敗: %w", err)
	}

	states := make(map[string]effectiveTagReviewState, len(tagCandidates))
	for _, candidate := range tagCandidates {
		state := effectiveTagReviewState{
			Include:     true,
			Address:     candidate.Address,
			PointID:     candidate.PointID,
			TagID:       cloneOptionalString(candidate.TagID),
			TagKey:      candidate.TagKey,
			DisplayName: candidate.DisplayName,
			DataType:    candidate.DataType,
		}

		decision := decisionsByCandidateID[candidate.ID]
		if decision != nil {
			switch decision.Action {
			case schema.SourceRuleTagReviewDecisionActionRename:
				if strings.TrimSpace(decision.TagKey) == "" {
					return nil, fmt.Errorf("rename decision for %s is missing tag_key", candidate.ID)
				}
				state.TagID = nil
				state.TagKey = decision.TagKey
			case schema.SourceRuleTagReviewDecisionActionSkip:
				state.Include = false
			case schema.SourceRuleTagReviewDecisionActionOverride:
				if decision.OverrideTagID == nil || strings.TrimSpace(*decision.OverrideTagID) == "" {
					return nil, fmt.Errorf("override decision for %s is missing override_tag_id", candidate.ID)
				}
				if strings.TrimSpace(decision.TagKey) == "" {
					return nil, fmt.Errorf("override decision for %s is missing tag_key", candidate.ID)
				}
				state.TagID = stringPtr(*decision.OverrideTagID)
				state.TagKey = decision.TagKey
			default:
				return nil, fmt.Errorf("unsupported tag review decision action: %s", decision.Action)
			}
		}

		states[candidate.ID] = state
	}

	return states, nil
}
