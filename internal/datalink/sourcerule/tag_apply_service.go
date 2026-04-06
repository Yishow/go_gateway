package sourcerule

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"
)

var (
	ErrInvalidTagApplyRequest   = errors.New("invalid tag apply request")
	ErrTagApplyRevisionConflict = errors.New("tag apply revision conflict")
)

type ApplyTagCandidatesRequest struct {
	RevisionID   string   `json:"revision_id"`
	CandidateIDs []string `json:"candidate_ids"`
}

type ApplyTagCandidateResult struct {
	CandidateID string `json:"candidate_id"`
	Status      string `json:"status"`
	TagID       string `json:"tag_id,omitempty"`
	MappingID   string `json:"mapping_id,omitempty"`
	Error       string `json:"error,omitempty"`
}

type ApplyTagCandidatesResponse struct {
	SourceRuleID string                    `json:"source_rule_id"`
	RevisionID   string                    `json:"revision_id"`
	Results      []ApplyTagCandidateResult `json:"results"`
}

func (s *Service) ApplyTagCandidates(ctx context.Context, ruleID string, req ApplyTagCandidatesRequest) (*ApplyTagCandidatesResponse, error) {
	if s.tagSvc == nil || s.mappingSvc == nil {
		return nil, fmt.Errorf("tag/mapping 服務未配置，無法套用來源規則標籤")
	}

	candidateIDs, revisionID, err := validateApplyTagCandidatesRequest(req)
	if err != nil {
		return nil, err
	}

	rule, err := s.repo.GetByID(ctx, ruleID)
	if err != nil {
		return nil, fmt.Errorf("取得來源規則失敗: %w", err)
	}
	if rule.RevisionID != revisionID {
		return nil, fmt.Errorf("%w: requested=%s current=%s", ErrTagApplyRevisionConflict, revisionID, rule.RevisionID)
	}

	links, err := s.repo.ListLinks(ctx, rule.ID)
	if err != nil {
		return nil, fmt.Errorf("列出來源規則連結失敗: %w", err)
	}
	previousLinks := cloneSourceRuleLinks(links)
	nextLinks := cloneSourceRuleLinks(links)
	candidateByID, err := s.listCurrentTagCandidatesByID(ctx, rule)
	if err != nil {
		return nil, err
	}
	decisionByCandidateID, err := s.listCurrentTagReviewDecisionsByCandidateID(ctx, rule.ID)
	if err != nil {
		return nil, fmt.Errorf("列出來源規則 tag 審查決策失敗: %w", err)
	}
	linkByAddress := make(map[string]*schema.SourceRuleLink, len(nextLinks))
	for _, link := range nextLinks {
		if link == nil {
			continue
		}
		linkByAddress[normalizeAddressKey(link.Address)] = link
	}

	for _, candidateID := range candidateIDs {
		candidate, exists := candidateByID[candidateID]
		if !exists {
			return nil, fmt.Errorf("%w: tag candidate %s not found in revision %s", ErrInvalidTagApplyRequest, candidateID, rule.RevisionID)
		}
		if linkByAddress[normalizeAddressKey(candidate.Address)] == nil {
			return nil, fmt.Errorf("%w: source rule link not found for candidate %s", ErrInvalidTagApplyRequest, candidateID)
		}
	}

	result := newTagMappingSyncResult()
	response := &ApplyTagCandidatesResponse{
		SourceRuleID: rule.ID,
		RevisionID:   rule.RevisionID,
		Results:      make([]ApplyTagCandidateResult, 0, len(candidateIDs)),
	}
	linksChanged := false

	for _, candidateID := range candidateIDs {
		candidate := candidateByID[candidateID]
		link := linkByAddress[normalizeAddressKey(candidate.Address)]
		pointRecord, err := s.pointSvc.GetByID(ctx, link.PointID)
		if err != nil {
			s.rollbackTagMappingSync(ctx, result)
			return nil, fmt.Errorf("取得來源規則衍生點位失敗: %w", err)
		}

		candidateResult := newTagMappingSyncResult()
		effectiveTagKey, overrideTagID, err := resolveTagApplyDecision(candidate, decisionByCandidateID[candidateID])
		if err != nil {
			response.Results = append(response.Results, ApplyTagCandidateResult{
				CandidateID: candidate.ID,
				Status:      "failed",
				Error:       err.Error(),
			})
			continue
		}

		tagRecord, mappingRecord, err := s.applyRuleTagMapping(ctx, rule, pointRecord, link, candidate, effectiveTagKey, overrideTagID, &candidateResult)
		if err != nil {
			s.rollbackTagMappingSync(ctx, candidateResult)
			response.Results = append(response.Results, ApplyTagCandidateResult{
				CandidateID: candidate.ID,
				Status:      "failed",
				Error:       err.Error(),
			})
			continue
		}

		mergeTagMappingSyncResult(&result, candidateResult)
		link.TagID = stringPtr(tagRecord.ID)
		link.MappingID = stringPtr(mappingRecord.ID)
		link.UpdatedAt = time.Now()
		linksChanged = true
		response.Results = append(response.Results, ApplyTagCandidateResult{
			CandidateID: candidate.ID,
			Status:      "applied",
			TagID:       tagRecord.ID,
			MappingID:   mappingRecord.ID,
		})
	}

	if linksChanged {
		if err := s.replaceRuleLinks(ctx, rule.ID, previousLinks, nextLinks); err != nil {
			s.rollbackTagMappingSync(ctx, result)
			return nil, err
		}
		if err := s.persistCandidateSnapshots(ctx, rule, nextLinks); err != nil {
			restoreLinksErr := s.replaceRuleLinks(ctx, rule.ID, nextLinks, previousLinks)
			s.rollbackTagMappingSync(ctx, result)
			var restoreSnapshotsErr error
			if restoreLinksErr == nil {
				restoreSnapshotsErr = s.persistCandidateSnapshots(ctx, rule, previousLinks)
			}
			switch {
			case restoreLinksErr != nil && restoreSnapshotsErr != nil:
				return nil, fmt.Errorf("refresh apply candidates: %w (restore links: %v, restore snapshots: %v)", err, restoreLinksErr, restoreSnapshotsErr)
			case restoreLinksErr != nil:
				return nil, fmt.Errorf("refresh apply candidates: %w (restore links: %v)", err, restoreLinksErr)
			case restoreSnapshotsErr != nil:
				return nil, fmt.Errorf("refresh apply candidates: %w (restore snapshots: %v)", err, restoreSnapshotsErr)
			default:
				return nil, err
			}
		}
	}

	return response, nil
}

func newTagMappingSyncResult() tagMappingSyncResult {
	return tagMappingSyncResult{
		updatedMappings: make(map[string]mappingRollbackState),
		updatedTags:     make(map[string]tagRollbackState),
	}
}

func mergeTagMappingSyncResult(target *tagMappingSyncResult, source tagMappingSyncResult) {
	target.createdTagIDs = append(target.createdTagIDs, source.createdTagIDs...)
	target.createdMappingIDs = append(target.createdMappingIDs, source.createdMappingIDs...)
	if len(source.updatedMappings) > 0 {
		if target.updatedMappings == nil {
			target.updatedMappings = make(map[string]mappingRollbackState, len(source.updatedMappings))
		}
		for mappingID, state := range source.updatedMappings {
			if _, exists := target.updatedMappings[mappingID]; exists {
				continue
			}
			target.updatedMappings[mappingID] = state
		}
	}
	if len(source.updatedTags) > 0 {
		if target.updatedTags == nil {
			target.updatedTags = make(map[string]tagRollbackState, len(source.updatedTags))
		}
		for tagID, state := range source.updatedTags {
			if _, exists := target.updatedTags[tagID]; exists {
				continue
			}
			target.updatedTags[tagID] = state
		}
	}
}

func validateApplyTagCandidatesRequest(req ApplyTagCandidatesRequest) ([]string, string, error) {
	revisionID := strings.TrimSpace(req.RevisionID)
	if revisionID == "" {
		return nil, "", fmt.Errorf("%w: revision_id is required", ErrInvalidTagApplyRequest)
	}
	if len(req.CandidateIDs) == 0 {
		return nil, "", fmt.Errorf("%w: candidate_ids is required", ErrInvalidTagApplyRequest)
	}

	seen := make(map[string]struct{}, len(req.CandidateIDs))
	candidateIDs := make([]string, 0, len(req.CandidateIDs))
	for _, candidateID := range req.CandidateIDs {
		normalized := strings.TrimSpace(candidateID)
		if normalized == "" {
			return nil, "", fmt.Errorf("%w: candidate_ids must not contain empty values", ErrInvalidTagApplyRequest)
		}
		if _, exists := seen[normalized]; exists {
			continue
		}
		seen[normalized] = struct{}{}
		candidateIDs = append(candidateIDs, normalized)
	}

	return candidateIDs, revisionID, nil
}

func (s *Service) listCurrentTagCandidatesByID(ctx context.Context, rule *schema.SourceRule) (map[string]schema.SourceRuleTagCandidate, error) {
	snapshots, err := s.repo.ListCandidateSnapshots(ctx, rule.ID, rule.RevisionID)
	if err != nil {
		return nil, fmt.Errorf("列出來源規則候選快照失敗: %w", err)
	}

	candidates := make(map[string]schema.SourceRuleTagCandidate)
	for _, snapshot := range snapshots {
		if snapshot == nil || snapshot.CandidateType != schema.SourceRuleCandidateTypeTags {
			continue
		}
		items, decodeErr := decodeCandidatePayload[schema.SourceRuleTagCandidate](snapshot.Payload)
		if decodeErr != nil {
			return nil, fmt.Errorf("解析來源規則 tag 候選快照失敗: %w", decodeErr)
		}
		for _, candidate := range items {
			candidates[candidate.ID] = candidate
		}
	}

	return candidates, nil
}

func (s *Service) listCurrentTagReviewDecisionsByCandidateID(ctx context.Context, ruleID string) (map[string]*schema.SourceRuleTagReviewDecision, error) {
	decisions, err := s.repo.ListTagReviewDecisions(ctx, ruleID)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*schema.SourceRuleTagReviewDecision)
	for _, decision := range decisions {
		if decision == nil || decision.Stale {
			continue
		}
		result[decision.CandidateID] = decision
	}
	return result, nil
}

func resolveTagApplyDecision(candidate schema.SourceRuleTagCandidate, decision *schema.SourceRuleTagReviewDecision) (string, string, error) {
	if decision == nil {
		return candidate.TagKey, "", nil
	}

	switch decision.Action {
	case schema.SourceRuleTagReviewDecisionActionRename:
		if strings.TrimSpace(decision.TagKey) == "" {
			return "", "", fmt.Errorf("%w: rename decision for %s is missing tag_key", ErrInvalidTagApplyRequest, candidate.ID)
		}
		return decision.TagKey, "", nil
	case schema.SourceRuleTagReviewDecisionActionOverride:
		if decision.OverrideTagID == nil || strings.TrimSpace(*decision.OverrideTagID) == "" {
			return "", "", fmt.Errorf("%w: override decision for %s is missing override_tag_id", ErrInvalidTagApplyRequest, candidate.ID)
		}
		return decision.TagKey, *decision.OverrideTagID, nil
	case schema.SourceRuleTagReviewDecisionActionSkip:
		return "", "", fmt.Errorf("%w: candidate %s is marked skip", ErrInvalidTagApplyRequest, candidate.ID)
	default:
		return "", "", fmt.Errorf("%w: unsupported decision action %s", ErrInvalidTagApplyRequest, decision.Action)
	}
}

func (s *Service) applyRuleTagMapping(
	ctx context.Context,
	rule *schema.SourceRule,
	pointRecord *schema.Point,
	link *schema.SourceRuleLink,
	candidate schema.SourceRuleTagCandidate,
	effectiveTagKey string,
	overrideTagID string,
	result *tagMappingSyncResult,
) (*schema.Tag, *schema.Mapping, error) {
	transformPipeline := s.buildRuleTransformPipeline(rule, pointRecord)

	tagRecord, createdTag, err := s.resolveRuleTag(ctx, rule, pointRecord, link, effectiveTagKey, overrideTagID, result)
	if err != nil {
		return nil, nil, err
	}
	if createdTag {
		result.createdTagIDs = append(result.createdTagIDs, tagRecord.ID)
	}
	if err := s.validateTagAvailability(ctx, tagRecord, pointRecord.ID); err != nil {
		return nil, nil, err
	}

	pointID := pointRecord.ID
	pointMappings, err := s.mappingSvc.List(ctx, mapping.ListFilter{PointID: &pointID})
	if err != nil {
		return nil, nil, fmt.Errorf("查詢 point 既有映射失敗: %w", err)
	}
	if len(pointMappings) > 1 {
		return nil, nil, fmt.Errorf("point %s 存在多條映射，無法套用來源規則標籤", pointRecord.ID)
	}
	if len(pointMappings) == 1 {
		if pointMappings[0].TagID != tagRecord.ID {
			return nil, nil, fmt.Errorf("point %s 已綁定其他 tag，無法套用來源規則標籤", pointRecord.ID)
		}
		mappingRecord, applyErr := s.applyRuleManagedMapping(ctx, rule, pointRecord, link, pointMappings[0], rule.Enabled, transformPipeline, result)
		if applyErr != nil {
			return nil, nil, applyErr
		}
		return tagRecord, mappingRecord, nil
	}

	tagID := tagRecord.ID
	tagMappings, err := s.mappingSvc.List(ctx, mapping.ListFilter{TagID: &tagID})
	if err != nil {
		return nil, nil, fmt.Errorf("查詢 tag 既有映射失敗: %w", err)
	}
	if len(tagMappings) > 1 {
		return nil, nil, fmt.Errorf("tag %s 存在多條映射，無法套用來源規則標籤", tagRecord.Key)
	}
	if len(tagMappings) == 1 {
		if tagMappings[0].PointID != pointRecord.ID {
			return nil, nil, fmt.Errorf("tag %s 已綁定其他 point", tagRecord.Key)
		}
		mappingRecord, applyErr := s.applyRuleManagedMapping(ctx, rule, pointRecord, link, tagMappings[0], rule.Enabled, transformPipeline, result)
		if applyErr != nil {
			return nil, nil, applyErr
		}
		return tagRecord, mappingRecord, nil
	}

	ruleCandidateID, proposedSignature, err := ruleManagedMappingMetadata(rule, pointRecord, link, transformPipeline)
	if err != nil {
		return nil, nil, err
	}
	mappingEnabled := rule.Enabled
	mappingStatus := schema.MappingStatusActive
	if !mappingEnabled {
		mappingStatus = schema.MappingStatusDraft
	}
	mappingRecord, err := s.mappingSvc.Create(ctx, mapping.CreateMappingRequest{
		PointID:              pointRecord.ID,
		TagID:                tagRecord.ID,
		Enabled:              &mappingEnabled,
		TransformPipeline:    transformPipeline,
		Status:               &mappingStatus,
		RuleCandidateID:      ruleCandidateID,
		ProposedSignature:    proposedSignature,
		LastAppliedSignature: proposedSignature,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("建立來源規則映射失敗: %w", err)
	}
	result.createdMappingIDs = append(result.createdMappingIDs, mappingRecord.ID)
	_ = candidate
	return tagRecord, mappingRecord, nil
}

func (s *Service) applyRuleManagedMapping(
	ctx context.Context,
	rule *schema.SourceRule,
	pointRecord *schema.Point,
	link *schema.SourceRuleLink,
	mappingRecord *schema.Mapping,
	enabled bool,
	transformPipeline []schema.TransformStep,
	result *tagMappingSyncResult,
) (*schema.Mapping, error) {
	ruleCandidateID, proposedSignature, err := ruleManagedMappingMetadata(rule, pointRecord, link, transformPipeline)
	if err != nil {
		return nil, err
	}
	nextPipelineJSON, err := encodeTransformPipeline(transformPipeline)
	if err != nil {
		return nil, fmt.Errorf("序列化來源規則映射轉換管線失敗: %w", err)
	}
	if mappingRecord.RuleCandidateID != "" && mappingRecord.RuleCandidateID != ruleCandidateID {
		return nil, fmt.Errorf("映射已綁定其他來源規則候選，無法重新套用")
	}
	if err := rememberMappingRollbackState(result, mappingRecord); err != nil {
		return nil, err
	}

	status := schema.MappingStatusActive
	if !enabled {
		status = schema.MappingStatusDraft
	}
	blockingReason := ""
	if mappingRecord.TransformPipeline == nextPipelineJSON &&
		mappingRecord.Enabled == enabled &&
		mappingRecord.Status == status &&
		mappingRecord.RuleCandidateID == ruleCandidateID &&
		mappingRecord.ProposedSignature == proposedSignature &&
		mappingRecord.LastAppliedSignature == proposedSignature &&
		mappingRecord.BlockingReason == blockingReason {
		return mappingRecord, nil
	}

	updateReq := mapping.UpdateMappingRequest{}
	if mappingRecord.TransformPipeline != nextPipelineJSON {
		updateReq.TransformPipeline = transformPipeline
	}
	if mappingRecord.Enabled != enabled {
		updateReq.Enabled = &enabled
	}
	if mappingRecord.Status != status {
		updateReq.Status = &status
	}
	if mappingRecord.RuleCandidateID != ruleCandidateID {
		updateReq.RuleCandidateID = &ruleCandidateID
	}
	if mappingRecord.ProposedSignature != proposedSignature {
		updateReq.ProposedSignature = &proposedSignature
	}
	if mappingRecord.LastAppliedSignature != proposedSignature {
		updateReq.LastAppliedSignature = &proposedSignature
	}
	if mappingRecord.BlockingReason != blockingReason {
		updateReq.BlockingReason = &blockingReason
	}

	updated, err := s.mappingSvc.Update(ctx, mappingRecord.ID, updateReq)
	if err != nil {
		return nil, fmt.Errorf("套用來源規則映射失敗: %w", err)
	}
	return updated, nil
}
