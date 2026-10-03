package sourcerule

import (
	"context"
	"encoding/json"
	"fmt"

	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"
)

type mappingCandidateSignaturePayload struct {
	TransformPipeline []schema.TransformStep `json:"transform_pipeline"`
}

func buildMappingCandidateIdentity(ruleID, address string, dataType schema.DataType) schema.SourceRuleCandidateIdentity {
	_ = dataType
	return schema.SourceRuleCandidateIdentity{
		SourceRuleID:           ruleID,
		CandidateType:          schema.SourceRuleCandidateTypeMappings,
		CandidateKind:          schema.SourceRuleCandidateKindMapping,
		DerivedFromRuleAddress: normalizeAddressKey(address),
	}
}

// canonicalTransformPipeline drops steps that cannot change a value (a scale of
// x1 + 0), so a rule that spells the neutral scale out and a mapping that has no
// step at all describe the same candidate. An empty result is never nil.
func canonicalTransformPipeline(transformPipeline []schema.TransformStep) []schema.TransformStep {
	steps := make([]schema.TransformStep, 0, len(transformPipeline))
	for _, step := range transformPipeline {
		if isNeutralScaleStep(step) {
			continue
		}
		steps = append(steps, step)
	}
	return steps
}

func isNeutralScaleStep(step schema.TransformStep) bool {
	if step.Type != schema.TransformScale {
		return false
	}
	scale, scaleOK := step.Params["scale"].(float64)
	offset, offsetOK := step.Params["offset"].(float64)
	return scaleOK && offsetOK && scale == 1 && offset == 0
}

func mappingCandidateSignature(transformPipeline []schema.TransformStep) (string, error) {
	payload, err := json.Marshal(mappingCandidateSignaturePayload{TransformPipeline: canonicalTransformPipeline(transformPipeline)})
	if err != nil {
		return "", fmt.Errorf("序列化來源規則映射候選簽章失敗: %w", err)
	}
	return candidateHash("signature", payload), nil
}

func ruleManagedMappingMetadata(rule *schema.SourceRule, pointRecord *schema.Point, link *schema.SourceRuleLink, transformPipeline []schema.TransformStep) (resolvedCandidateID, proposedSignature string, metadataErr error) {
	identity := buildMappingCandidateIdentity(rule.ID, link.Address, desiredRuleTargetDataType(rule, pointRecord))
	id, err := candidateID(identity)
	if err != nil {
		return "", "", err
	}
	signature, err := mappingCandidateSignature(transformPipeline)
	if err != nil {
		return "", "", err
	}
	return id, signature, nil
}

func (s *Service) syncRuleManagedMapping(
	ctx context.Context,
	oldRule *schema.SourceRule,
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
	if mappingRecord.RuleCandidateID != "" && mappingRecord.RuleCandidateID != ruleCandidateID {
		return nil, fmt.Errorf("映射已綁定其他來源規則候選，無法自動同步")
	}

	appliedPipeline, err := decodeTransformPipeline(mappingRecord.TransformPipeline)
	if err != nil {
		return nil, fmt.Errorf("解析既有來源規則映射轉換管線失敗: %w", err)
	}
	lastAppliedSignature, err := mappingCandidateSignature(appliedPipeline)
	if err != nil {
		return nil, err
	}
	if oldRule == nil && mappingRecord.RuleCandidateID == "" && lastAppliedSignature != proposedSignature {
		return nil, fmt.Errorf("映射存在手動編輯的轉換管線，無法自動覆蓋")
	}

	status := schema.MappingStatusActive
	if !enabled {
		status = schema.MappingStatusDraft
	}
	blockingReason := ""
	if proposedSignature != lastAppliedSignature {
		status = schema.MappingStatusOutOfSync
		blockingReason = "來源規則修訂已變更映射候選，需明確重新套用後才能覆蓋既有 pipeline"
	} else {
		lastAppliedSignature = proposedSignature
	}

	// A pipeline whose signature matches is left as confirmed: rewriting an
	// equivalent one would change the mapping revision write groups are bound to.
	// A different one is never overwritten here; that needs an explicit re-apply.
	if mappingRecord.Enabled == enabled &&
		mappingRecord.Status == status &&
		mappingRecord.RuleCandidateID == ruleCandidateID &&
		mappingRecord.ProposedSignature == proposedSignature &&
		mappingRecord.LastAppliedSignature == lastAppliedSignature &&
		mappingRecord.BlockingReason == blockingReason {
		return mappingRecord, nil
	}
	if err := rememberMappingRollbackState(result, mappingRecord); err != nil {
		return nil, err
	}

	updateReq := mapping.UpdateMappingRequest{}
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
	if mappingRecord.LastAppliedSignature != lastAppliedSignature {
		updateReq.LastAppliedSignature = &lastAppliedSignature
	}
	if mappingRecord.BlockingReason != blockingReason {
		updateReq.BlockingReason = &blockingReason
	}

	updated, err := s.mappingSvc.Update(ctx, mappingRecord.ID, updateReq)
	if err != nil {
		return nil, fmt.Errorf("同步來源規則映射設定失敗: %w", err)
	}
	return updated, nil
}

func rememberMappingRollbackState(result *tagMappingSyncResult, mappingRecord *schema.Mapping) error {
	if result.updatedMappings == nil {
		result.updatedMappings = make(map[string]mappingRollbackState)
	}
	if _, exists := result.updatedMappings[mappingRecord.ID]; exists {
		return nil
	}
	previousPipeline, err := decodeTransformPipeline(mappingRecord.TransformPipeline)
	if err != nil {
		return fmt.Errorf("解析既有來源規則映射轉換管線失敗: %w", err)
	}
	result.updatedMappings[mappingRecord.ID] = mappingRollbackState{
		enabled:              mappingRecord.Enabled,
		transformPipeline:    previousPipeline,
		status:               mappingRecord.Status,
		ruleCandidateID:      mappingRecord.RuleCandidateID,
		proposedSignature:    mappingRecord.ProposedSignature,
		lastAppliedSignature: mappingRecord.LastAppliedSignature,
		blockingReason:       mappingRecord.BlockingReason,
	}
	return nil
}
