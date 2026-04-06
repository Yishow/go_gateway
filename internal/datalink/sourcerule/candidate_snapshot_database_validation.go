package sourcerule

import (
	"context"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema"
)

type databaseTargetValidationIndex struct {
	connectorIssues map[string][]DatabaseTargetValidationIssue
	mappingIssues   map[string][]DatabaseTargetValidationIssue
}

func applyDatabaseTargetValidation(
	candidate *schema.SourceRuleDatabaseOutputCandidate,
	mappingRecord *schema.DatabaseTargetMapping,
	validationIndex databaseTargetValidationIndex,
) {
	if candidate == nil || mappingRecord == nil {
		return
	}

	issues := append([]DatabaseTargetValidationIssue(nil), validationIndex.connectorIssues[mappingRecord.ConnectorID]...)
	issues = append(issues, validationIndex.mappingIssues[mappingRecord.ID]...)
	if len(issues) == 0 {
		return
	}

	candidate.Status = schema.SourceRuleOutputStatusOutOfSync
	candidate.BlockingReason = joinValidationIssueMessages(issues)
}

func (s *Service) listDatabaseTargetValidationIndex(
	ctx context.Context,
	mappingsByTagID map[string][]*schema.DatabaseTargetMapping,
) (databaseTargetValidationIndex, error) {
	validator := s.databaseTargetConnectorValidator()
	if validator == nil {
		return databaseTargetValidationIndex{
			connectorIssues: map[string][]DatabaseTargetValidationIssue{},
			mappingIssues:   map[string][]DatabaseTargetValidationIssue{},
		}, nil
	}

	connectorIDs := make(map[string]struct{})
	for _, mappings := range mappingsByTagID {
		for _, mappingRecord := range mappings {
			if mappingRecord == nil || strings.TrimSpace(mappingRecord.ConnectorID) == "" {
				continue
			}
			connectorIDs[mappingRecord.ConnectorID] = struct{}{}
		}
	}

	index := databaseTargetValidationIndex{
		connectorIssues: make(map[string][]DatabaseTargetValidationIssue),
		mappingIssues:   make(map[string][]DatabaseTargetValidationIssue),
	}
	for connectorID := range connectorIDs {
		validation, err := validator.Validate(ctx, connectorID)
		if err != nil {
			return databaseTargetValidationIndex{}, fmt.Errorf("驗證資料庫連接器 %s 失敗: %w", connectorID, err)
		}
		if validation == nil {
			continue
		}
		for _, issue := range validation.Issues {
			if strings.TrimSpace(issue.Severity) != "error" {
				continue
			}
			if strings.TrimSpace(issue.MappingID) == "" {
				index.connectorIssues[connectorID] = append(index.connectorIssues[connectorID], issue)
				continue
			}
			index.mappingIssues[issue.MappingID] = append(index.mappingIssues[issue.MappingID], issue)
		}
	}

	return index, nil
}

func databaseOutputSnapshotState(
	candidates []schema.SourceRuleDatabaseOutputCandidate,
) ([]schema.SourceRuleDatabaseOutputCandidate, schema.SourceRuleCandidateStatus, string, error) {
	for _, candidate := range candidates {
		if candidate.Status == schema.SourceRuleOutputStatusBlocked || candidate.Status == schema.SourceRuleOutputStatusOutOfSync {
			reason := candidate.BlockingReason
			if strings.TrimSpace(reason) == "" {
				reason = "one or more database output candidates require context repair"
			}
			return candidates, schema.SourceRuleCandidateStatusBlocked, reason, nil
		}
	}
	return candidates, schema.SourceRuleCandidateStatusReady, "", nil
}

func joinValidationIssueMessages(issues []DatabaseTargetValidationIssue) string {
	messages := make([]string, 0, len(issues))
	for _, issue := range issues {
		message := strings.TrimSpace(issue.Message)
		if message == "" {
			continue
		}
		messages = append(messages, message)
	}
	if len(messages) == 0 {
		return ""
	}
	return strings.Join(messages, "; ")
}
