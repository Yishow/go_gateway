package sourcerule

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema"
)

var (
	ErrInvalidOutputApplyRequest   = errors.New("invalid output apply request")
	ErrOutputApplyRevisionConflict = errors.New("output apply revision conflict")
)

const (
	outputApplyCodeConflict             = "conflict"
	outputApplyCodeConnectorUnavailable = "connector_unavailable"
	outputApplyCodeOutOfSync            = "out_of_sync"
	outputApplyCodeRegisterMissing      = "register_missing"
	outputApplyCodeSchemaMissing        = "schema_missing"
	outputApplyCodeTagMissing           = "tag_missing"
	outputApplyCodeTypeConflict         = "type_conflict"
)

type ApplyOutputCandidatesRequest struct {
	RevisionID   string   `json:"revision_id"`
	CandidateIDs []string `json:"candidate_ids"`
}

type ApplyOutputCandidateResult struct {
	CandidateID string `json:"candidate_id"`
	Status      string `json:"status"`
	Code        string `json:"code,omitempty"`
	Reason      string `json:"reason,omitempty"`
	MappingID   string `json:"mapping_id,omitempty"`
	ConnectorID string `json:"connector_id,omitempty"`
}

type ApplyOutputCandidatesResponse struct {
	SourceRuleID string                       `json:"source_rule_id"`
	RevisionID   string                       `json:"revision_id"`
	Results      []ApplyOutputCandidateResult `json:"results"`
}

func (s *Service) ApplyDatabaseOutputCandidates(
	ctx context.Context,
	ruleID string,
	req ApplyOutputCandidatesRequest,
) (*ApplyOutputCandidatesResponse, error) {
	candidateIDs, revisionID, err := validateApplyOutputCandidatesRequest(req)
	if err != nil {
		return nil, err
	}

	rule, err := s.repo.GetByID(ctx, ruleID)
	if err != nil {
		return nil, fmt.Errorf("取得來源規則失敗: %w", err)
	}
	if rule.RevisionID != revisionID {
		return nil, fmt.Errorf("%w: requested=%s current=%s", ErrOutputApplyRevisionConflict, revisionID, rule.RevisionID)
	}

	snapshots, err := s.repo.ListCandidateSnapshots(ctx, rule.ID, rule.RevisionID)
	if err != nil {
		return nil, fmt.Errorf("列出來源規則候選快照失敗: %w", err)
	}
	databaseSnapshot, err := outputSnapshotByType(snapshots, schema.SourceRuleCandidateTypeDatabaseOutputs)
	if err != nil {
		return nil, err
	}
	candidates, err := decodeCandidatePayload[schema.SourceRuleDatabaseOutputCandidate](databaseSnapshot.Payload)
	if err != nil {
		return nil, fmt.Errorf("解析來源規則 database output 候選快照失敗: %w", err)
	}

	candidateByID := make(map[string]schema.SourceRuleDatabaseOutputCandidate, len(candidates))
	for _, candidate := range candidates {
		candidateByID[candidate.ID] = candidate
	}

	validator := s.databaseTargetConnectorValidator()
	validationCache := map[string]*DatabaseTargetConnectorValidation{}
	validationErrors := map[string]error{}

	response := &ApplyOutputCandidatesResponse{
		SourceRuleID: rule.ID,
		RevisionID:   rule.RevisionID,
		Results:      make([]ApplyOutputCandidateResult, 0, len(candidateIDs)),
	}
	for _, candidateID := range candidateIDs {
		candidate, exists := candidateByID[candidateID]
		if !exists {
			response.Results = append(response.Results, ApplyOutputCandidateResult{
				CandidateID: candidateID,
				Status:      "failed",
				Code:        "validation",
				Reason:      fmt.Sprintf("candidate %s not found in revision %s", candidateID, rule.RevisionID),
			})
			continue
		}

		if candidate.Status == schema.SourceRuleOutputStatusBlocked || candidate.Status == schema.SourceRuleOutputStatusOutOfSync {
			response.Results = append(response.Results, ApplyOutputCandidateResult{
				CandidateID: candidate.ID,
				Status:      "failed",
				Code:        classifyDatabaseBlockingCode(candidate.BlockingReason),
				Reason:      defaultReason(candidate.BlockingReason, "candidate is blocked"),
				MappingID:   derefOptional(candidate.MappingID),
				ConnectorID: candidate.ConnectorID,
			})
			continue
		}
		if candidate.MappingID == nil || strings.TrimSpace(*candidate.MappingID) == "" {
			response.Results = append(response.Results, ApplyOutputCandidateResult{
				CandidateID: candidate.ID,
				Status:      "failed",
				Code:        "schema_missing",
				Reason:      "database mapping scope is not configured",
			})
			continue
		}
		if strings.TrimSpace(candidate.ConnectorID) == "" {
			response.Results = append(response.Results, ApplyOutputCandidateResult{
				CandidateID: candidate.ID,
				Status:      "failed",
				Code:        "schema_missing",
				Reason:      "connector_id is required",
				MappingID:   *candidate.MappingID,
			})
			continue
		}

		if validator != nil {
			if _, ok := validationCache[candidate.ConnectorID]; !ok {
				validation, validateErr := validator.Validate(ctx, candidate.ConnectorID)
				if validateErr != nil {
					validationErrors[candidate.ConnectorID] = validateErr
				} else {
					validationCache[candidate.ConnectorID] = validation
				}
			}
			if validateErr, hasErr := validationErrors[candidate.ConnectorID]; hasErr {
				response.Results = append(response.Results, ApplyOutputCandidateResult{
					CandidateID: candidate.ID,
					Status:      "failed",
					Code:        "connector_unavailable",
					Reason:      validateErr.Error(),
					MappingID:   *candidate.MappingID,
					ConnectorID: candidate.ConnectorID,
				})
				continue
			}
			if validation := validationCache[candidate.ConnectorID]; validation != nil {
				if issueCode, issueMessage, blocked := matchValidationIssue(validation, *candidate.MappingID); blocked {
					response.Results = append(response.Results, ApplyOutputCandidateResult{
						CandidateID: candidate.ID,
						Status:      "failed",
						Code:        mapValidationIssueCode(issueCode),
						Reason:      issueMessage,
						MappingID:   *candidate.MappingID,
						ConnectorID: candidate.ConnectorID,
					})
					continue
				}
			}
		}

		response.Results = append(response.Results, ApplyOutputCandidateResult{
			CandidateID: candidate.ID,
			Status:      "success",
			MappingID:   *candidate.MappingID,
			ConnectorID: candidate.ConnectorID,
		})
	}

	return response, nil
}

func (s *Service) ApplyLocalModbusOutputCandidates(
	ctx context.Context,
	ruleID string,
	req ApplyOutputCandidatesRequest,
) (*ApplyOutputCandidatesResponse, error) {
	candidateIDs, revisionID, err := validateApplyOutputCandidatesRequest(req)
	if err != nil {
		return nil, err
	}

	rule, err := s.repo.GetByID(ctx, ruleID)
	if err != nil {
		return nil, fmt.Errorf("取得來源規則失敗: %w", err)
	}
	if rule.RevisionID != revisionID {
		return nil, fmt.Errorf("%w: requested=%s current=%s", ErrOutputApplyRevisionConflict, revisionID, rule.RevisionID)
	}

	snapshots, err := s.repo.ListCandidateSnapshots(ctx, rule.ID, rule.RevisionID)
	if err != nil {
		return nil, fmt.Errorf("列出來源規則候選快照失敗: %w", err)
	}
	localModbusSnapshot, err := outputSnapshotByType(snapshots, schema.SourceRuleCandidateTypeLocalModbusOutputs)
	if err != nil {
		return nil, err
	}

	response := &ApplyOutputCandidatesResponse{
		SourceRuleID: rule.ID,
		RevisionID:   rule.RevisionID,
		Results:      make([]ApplyOutputCandidateResult, 0, len(candidateIDs)),
	}
	if localModbusSnapshot.Status == schema.SourceRuleCandidateStatusDeferred {
		for _, candidateID := range candidateIDs {
			response.Results = append(response.Results, ApplyOutputCandidateResult{
				CandidateID: candidateID,
				Status:      "skipped",
				Code:        "deferred",
				Reason:      defaultReason(localModbusSnapshot.Reason, localModbusOutputsDeferredReason),
			})
		}
		return response, nil
	}

	candidates, err := decodeCandidatePayload[schema.SourceRuleLocalModbusOutputCandidate](localModbusSnapshot.Payload)
	if err != nil {
		return nil, fmt.Errorf("解析來源規則 local modbus 候選快照失敗: %w", err)
	}
	candidateByID := make(map[string]schema.SourceRuleLocalModbusOutputCandidate, len(candidates))
	for _, candidate := range candidates {
		candidateID := strings.TrimSpace(candidate.ID)
		if candidateID == "" {
			continue
		}
		candidateByID[candidateID] = candidate
	}

	for _, candidateID := range candidateIDs {
		candidate, exists := candidateByID[candidateID]
		if !exists {
			response.Results = append(response.Results, ApplyOutputCandidateResult{
				CandidateID: candidateID,
				Status:      "failed",
				Code:        "validation",
				Reason:      fmt.Sprintf("candidate %s not found in revision %s", candidateID, rule.RevisionID),
			})
			continue
		}
		if code, reason, blocked := verifyLocalModbusApplyCandidate(candidate); blocked {
			response.Results = append(response.Results, ApplyOutputCandidateResult{
				CandidateID: candidate.ID,
				Status:      "failed",
				Code:        code,
				Reason:      reason,
			})
			continue
		}
		response.Results = append(response.Results, ApplyOutputCandidateResult{
			CandidateID: candidateID,
			Status:      "success",
		})
	}

	return response, nil
}

func verifyLocalModbusApplyCandidate(candidate schema.SourceRuleLocalModbusOutputCandidate) (code string, reason string, blocked bool) {
	switch {
	case candidate.Status == schema.SourceRuleLocalModbusOutputStatusBlockedConflict:
		return outputApplyCodeConflict, defaultReason(candidate.BlockingReason, "candidate is blocked by conflict"), true
	case candidate.Status == schema.SourceRuleLocalModbusOutputStatusOutOfSync:
		return outputApplyCodeOutOfSync, defaultReason(candidate.BlockingReason, "candidate is out of sync"), true
	case candidate.TagID == nil || strings.TrimSpace(*candidate.TagID) == "":
		return outputApplyCodeTagMissing, "local modbus tag is not configured", true
	case candidate.Register == nil:
		return outputApplyCodeRegisterMissing, "local modbus register is not configured", true
	default:
		return "", "", false
	}
}

func validateApplyOutputCandidatesRequest(req ApplyOutputCandidatesRequest) (candidateIDs []string, revisionID string, err error) {
	revisionID = strings.TrimSpace(req.RevisionID)
	if revisionID == "" {
		return nil, "", fmt.Errorf("%w: revision_id is required", ErrInvalidOutputApplyRequest)
	}
	if len(req.CandidateIDs) == 0 {
		return nil, "", fmt.Errorf("%w: candidate_ids is required", ErrInvalidOutputApplyRequest)
	}

	seen := make(map[string]struct{}, len(req.CandidateIDs))
	candidateIDs = make([]string, 0, len(req.CandidateIDs))
	for _, candidateID := range req.CandidateIDs {
		normalized := strings.TrimSpace(candidateID)
		if normalized == "" {
			return nil, "", fmt.Errorf("%w: candidate_ids must not contain empty values", ErrInvalidOutputApplyRequest)
		}
		if _, exists := seen[normalized]; exists {
			continue
		}
		seen[normalized] = struct{}{}
		candidateIDs = append(candidateIDs, normalized)
	}

	return candidateIDs, revisionID, nil
}

func outputSnapshotByType(
	snapshots []*schema.SourceRuleCandidateSnapshot,
	candidateType schema.SourceRuleCandidateType,
) (*schema.SourceRuleCandidateSnapshot, error) {
	for _, snapshot := range snapshots {
		if snapshot != nil && snapshot.CandidateType == candidateType {
			return snapshot, nil
		}
	}
	return nil, fmt.Errorf("%w: candidate snapshot %s is missing", ErrInvalidOutputApplyRequest, candidateType)
}

func matchValidationIssue(validation *DatabaseTargetConnectorValidation, mappingID string) (
	issueCode string,
	issueMessage string,
	blocked bool,
) {
	if validation == nil {
		return "", "", false
	}
	for _, issue := range validation.Issues {
		if strings.TrimSpace(issue.Severity) != "error" {
			continue
		}
		if issue.MappingID == "" || issue.MappingID == mappingID {
			return issue.Code, defaultReason(issue.Message, "database validation failed"), true
		}
	}
	if !validation.Ready {
		return outputApplyCodeConnectorUnavailable, "connector is not ready", true
	}
	return "", "", false
}

func mapValidationIssueCode(code string) string {
	normalized := strings.TrimSpace(strings.ToLower(code))
	switch normalized {
	case "connector_unreachable", "connector_unavailable", "connector_auth_failed":
		return outputApplyCodeConnectorUnavailable
	case "table_missing", "column_missing", "timestamp_missing", "timestamp_column_missing", "timestamp_column_not_unique":
		return outputApplyCodeSchemaMissing
	case "column_type_mismatch", "type_conflict":
		return outputApplyCodeTypeConflict
	case "mapping_missing":
		return outputApplyCodeSchemaMissing
	default:
		return outputApplyCodeConflict
	}
}

func classifyDatabaseBlockingCode(reason string) string {
	normalized := strings.ToLower(strings.TrimSpace(reason))
	switch {
	case normalized == "":
		return outputApplyCodeConflict
	case strings.Contains(normalized, "table"), strings.Contains(normalized, "欄位"), strings.Contains(normalized, "schema"):
		return outputApplyCodeSchemaMissing
	case strings.Contains(normalized, "unreachable"), strings.Contains(normalized, "連線"), strings.Contains(normalized, "connector"):
		return outputApplyCodeConnectorUnavailable
	case strings.Contains(normalized, "type"):
		return outputApplyCodeTypeConflict
	default:
		return outputApplyCodeConflict
	}
}

func defaultReason(reason, fallback string) string {
	trimmed := strings.TrimSpace(reason)
	if trimmed != "" {
		return trimmed
	}
	return fallback
}

func derefOptional(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}
