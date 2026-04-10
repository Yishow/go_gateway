package sourcerule

import (
	"context"
	"fmt"
	"sort"
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
) ([]schema.SourceRuleDatabaseOutputCandidate, schema.SourceRuleCandidateStatus, string, error) {
	states, err := s.listEffectiveTagReviewStates(ctx, rule.ID, tagCandidates)
	if err != nil {
		return nil, "", "", err
	}
	mappingsByTagID, err := s.listDatabaseTargetMappingsByTagID(ctx)
	if err != nil {
		return nil, "", "", err
	}
	connectorsByID, err := s.listDatabaseTargetConnectorsByID(ctx, mappingsByTagID)
	if err != nil {
		return nil, "", "", err
	}
	validationIndex, err := s.listDatabaseTargetValidationIndex(ctx, mappingsByTagID)
	if err != nil {
		return nil, "", "", err
	}

	candidates := make([]schema.SourceRuleDatabaseOutputCandidate, 0, len(tagCandidates))
	for _, tagCandidate := range tagCandidates {
		state, ok := states[tagCandidate.ID]
		if !ok || !state.Include {
			continue
		}

		mappings := effectiveStateMappings(state, mappingsByTagID)
		if len(mappings) == 0 {
			candidate, buildErr := buildDatabaseOutputCandidate(rule.ID, state, nil, nil, validationIndex)
			if buildErr != nil {
				return nil, "", "", buildErr
			}
			candidates = append(candidates, candidate)
			continue
		}

		for _, mappingRecord := range mappings {
			candidate, buildErr := buildDatabaseOutputCandidate(
				rule.ID,
				state,
				mappingRecord,
				connectorsByID[mappingRecord.ConnectorID],
				validationIndex,
			)
			if buildErr != nil {
				return nil, "", "", buildErr
			}
			candidates = append(candidates, candidate)
		}
	}

	return databaseOutputSnapshotState(candidates)
}

func buildDatabaseOutputCandidate(
	ruleID string,
	state effectiveTagReviewState,
	mappingRecord *schema.DatabaseTargetMapping,
	connectorRecord *schema.DatabaseConnector,
	validationIndex databaseTargetValidationIndex,
) (schema.SourceRuleDatabaseOutputCandidate, error) {
	candidate := schema.SourceRuleDatabaseOutputCandidate{
		Address:              state.Address,
		PointID:              state.PointID,
		TagID:                cloneOptionalString(state.TagID),
		TagKey:               state.TagKey,
		DisplayName:          state.DisplayName,
		DataType:             state.DataType,
		WriteIntervalSeconds: defaultDatabaseWriteIntervalPointer(),
		Status:               schema.SourceRuleOutputStatusReady,
	}
	candidate.GroupKey, candidate.ColumnName = inferDatabaseGroupAndColumn(state.TagKey)
	if mappingRecord != nil {
		candidate.MappingID = stringPtr(mappingRecord.ID)
		candidate.ConnectorID = mappingRecord.ConnectorID
		candidate.TableSchema = mappingRecord.TableSchema
		candidate.TableName = mappingRecord.TableName
		candidate.ColumnName = mappingRecord.ColumnName
		candidate.GroupKey = cloneOptionalString(mappingRecord.GroupKey)
		candidate.WriteMode = mappingRecord.WriteMode
		candidate.TimestampColumn = cloneOptionalString(mappingRecord.TimestampColumn)
		candidate.WriteIntervalSeconds = effectiveDatabaseWriteIntervalPointer(mappingRecord, connectorRecord)
		applyDatabaseTargetValidation(&candidate, mappingRecord, validationIndex)
	}
	candidate.Identity = buildDatabaseOutputCandidateIdentity(
		ruleID,
		state.Address,
		state.DataType,
		candidate.ConnectorID,
		candidate.GroupKey,
		candidate.TableSchema,
		candidate.TableName,
		candidate.ColumnName,
		candidate.WriteIntervalSeconds,
	)

	var err error
	candidate.ID, err = candidateID(candidate.Identity)
	if err != nil {
		return schema.SourceRuleDatabaseOutputCandidate{}, err
	}
	candidate.ProposedSignature, err = databaseOutputCandidateSignature(candidate)
	if err != nil {
		return schema.SourceRuleDatabaseOutputCandidate{}, err
	}
	return candidate, nil
}

func defaultDatabaseWriteIntervalPointer() *int {
	seconds := 15
	return &seconds
}

func effectiveDatabaseWriteIntervalPointer(
	mappingRecord *schema.DatabaseTargetMapping,
	connectorRecord *schema.DatabaseConnector,
) *int {
	if mappingRecord != nil && mappingRecord.WriteIntervalSeconds != nil {
		return cloneOptionalInt(mappingRecord.WriteIntervalSeconds)
	}
	if connectorRecord != nil && connectorRecord.DefaultWriteIntervalSeconds > 0 {
		seconds := connectorRecord.DefaultWriteIntervalSeconds
		return &seconds
	}
	return defaultDatabaseWriteIntervalPointer()
}

func inferDatabaseGroupAndColumn(tagKey string) (*string, string) {
	trimmed := strings.TrimSpace(tagKey)
	prefix, suffix, found := strings.Cut(trimmed, "/")
	if !found {
		return nil, ""
	}
	groupKey := strings.TrimSpace(prefix)
	columnName := normalizeDatabaseCandidateColumnName(suffix)
	if groupKey == "" || columnName == "" {
		return nil, ""
	}
	return stringPtr(groupKey), columnName
}

func normalizeDatabaseCandidateColumnName(value string) string {
	replaced := strings.ReplaceAll(strings.TrimSpace(value), "/", "_")
	if replaced == "" {
		return ""
	}
	return strings.ToLower(replaced)
}

func (s *Service) listDatabaseTargetMappingsByTagID(
	ctx context.Context,
) (map[string][]*schema.DatabaseTargetMapping, error) {
	reader := s.databaseTargetMappingReader()
	if reader == nil {
		return map[string][]*schema.DatabaseTargetMapping{}, nil
	}

	mappings, err := reader.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("列出資料庫目標映射失敗: %w", err)
	}

	result := make(map[string][]*schema.DatabaseTargetMapping)
	for _, mappingRecord := range mappings {
		if mappingRecord == nil || strings.TrimSpace(mappingRecord.TagID) == "" {
			continue
		}
		tagID := strings.TrimSpace(mappingRecord.TagID)
		result[tagID] = append(result[tagID], mappingRecord)
	}

	for tagID := range result {
		sort.Slice(result[tagID], func(i, j int) bool {
			left := result[tagID][i]
			right := result[tagID][j]
			if left.ConnectorID != right.ConnectorID {
				return left.ConnectorID < right.ConnectorID
			}
			if left.TableSchema != right.TableSchema {
				return left.TableSchema < right.TableSchema
			}
			if left.TableName != right.TableName {
				return left.TableName < right.TableName
			}
			if left.ColumnName != right.ColumnName {
				return left.ColumnName < right.ColumnName
			}
			return left.ID < right.ID
		})
	}

	return result, nil
}

func (s *Service) listDatabaseTargetConnectorsByID(
	ctx context.Context,
	mappingsByTagID map[string][]*schema.DatabaseTargetMapping,
) (map[string]*schema.DatabaseConnector, error) {
	reader := s.databaseTargetConnectorReader()
	if reader == nil {
		return map[string]*schema.DatabaseConnector{}, nil
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

	connectorsByID := make(map[string]*schema.DatabaseConnector, len(connectorIDs))
	for connectorID := range connectorIDs {
		connector, err := reader.GetByID(ctx, connectorID)
		if err != nil {
			return nil, fmt.Errorf("取得資料庫連接器 %s 失敗: %w", connectorID, err)
		}
		if connector == nil {
			continue
		}
		connectorsByID[connectorID] = connector
	}
	return connectorsByID, nil
}

func effectiveStateMappings(
	state effectiveTagReviewState,
	mappingsByTagID map[string][]*schema.DatabaseTargetMapping,
) []*schema.DatabaseTargetMapping {
	if state.TagID == nil {
		return nil
	}
	return mappingsByTagID[*state.TagID]
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
