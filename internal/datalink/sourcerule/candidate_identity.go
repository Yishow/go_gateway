package sourcerule

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"

	"go-gateway/internal/datalink/schema"
)

const (
	derivedTargetDataTypeScopeKey = "derived_target_data_type"
	databaseConnectorScopeKey     = "database_connector_id"
	databaseSchemaScopeKey        = "database_table_schema"
	databaseTableScopeKey         = "database_table_name"
	databaseColumnScopeKey        = "database_column_name"
)

type tagCandidateSignaturePayload struct {
	TagKey            string                 `json:"tag_key"`
	DisplayName       string                 `json:"display_name"`
	DataType          schema.DataType        `json:"data_type"`
	TransformPipeline []schema.TransformStep `json:"transform_pipeline"`
}

type databaseOutputCandidateSignaturePayload struct {
	TagID           *string                  `json:"tag_id,omitempty"`
	TagKey          string                   `json:"tag_key"`
	DisplayName     string                   `json:"display_name"`
	DataType        schema.DataType          `json:"data_type"`
	ConnectorID     string                   `json:"connector_id,omitempty"`
	TableSchema     string                   `json:"table_schema,omitempty"`
	TableName       string                   `json:"table_name,omitempty"`
	ColumnName      string                   `json:"column_name,omitempty"`
	WriteMode       schema.DatabaseWriteMode `json:"write_mode,omitempty"`
	TimestampColumn *string                  `json:"timestamp_column,omitempty"`
}

func buildTagCandidateIdentity(ruleID, address string, dataType schema.DataType) schema.SourceRuleCandidateIdentity {
	return schema.SourceRuleCandidateIdentity{
		SourceRuleID:           ruleID,
		CandidateType:          schema.SourceRuleCandidateTypeTags,
		CandidateKind:          schema.SourceRuleCandidateKindTag,
		DerivedFromRuleAddress: normalizeAddressKey(address),
		TargetBindingScope: []schema.SourceRuleCandidateScopeField{
			{
				Key:   derivedTargetDataTypeScopeKey,
				Value: string(dataType),
			},
		},
	}
}

func buildDatabaseOutputCandidateIdentity(
	ruleID string,
	address string,
	dataType schema.DataType,
	connectorID string,
	tableSchema string,
	tableName string,
	columnName string,
) schema.SourceRuleCandidateIdentity {
	scope := []schema.SourceRuleCandidateScopeField{
		{
			Key:   derivedTargetDataTypeScopeKey,
			Value: string(dataType),
		},
	}
	appendScope := func(key, value string) {
		if value == "" {
			return
		}
		scope = append(scope, schema.SourceRuleCandidateScopeField{
			Key:   key,
			Value: value,
		})
	}
	appendScope(databaseConnectorScopeKey, connectorID)
	appendScope(databaseSchemaScopeKey, tableSchema)
	appendScope(databaseTableScopeKey, tableName)
	appendScope(databaseColumnScopeKey, columnName)

	return schema.SourceRuleCandidateIdentity{
		SourceRuleID:           ruleID,
		CandidateType:          schema.SourceRuleCandidateTypeDatabaseOutputs,
		CandidateKind:          schema.SourceRuleCandidateKindDatabaseOutput,
		DerivedFromRuleAddress: normalizeAddressKey(address),
		TargetBindingScope:     scope,
	}
}

func candidateID(identity schema.SourceRuleCandidateIdentity) (string, error) {
	normalized := normalizedCandidateIdentity(identity)
	data, err := json.Marshal(normalized)
	if err != nil {
		return "", fmt.Errorf("序列化來源規則候選識別失敗: %w", err)
	}
	return candidateHash(string(normalized.CandidateKind), data), nil
}

func tagCandidateSignature(candidate schema.SourceRuleTagCandidate) (string, error) {
	payload, err := json.Marshal(tagCandidateSignaturePayload{
		TagKey:            candidate.TagKey,
		DisplayName:       candidate.DisplayName,
		DataType:          candidate.DataType,
		TransformPipeline: candidate.TransformPipeline,
	})
	if err != nil {
		return "", fmt.Errorf("序列化來源規則候選簽章失敗: %w", err)
	}
	return candidateHash("signature", payload), nil
}

func databaseOutputCandidateSignature(candidate schema.SourceRuleDatabaseOutputCandidate) (string, error) {
	payload, err := json.Marshal(databaseOutputCandidateSignaturePayload{
		TagID:           candidate.TagID,
		TagKey:          candidate.TagKey,
		DisplayName:     candidate.DisplayName,
		DataType:        candidate.DataType,
		ConnectorID:     candidate.ConnectorID,
		TableSchema:     candidate.TableSchema,
		TableName:       candidate.TableName,
		ColumnName:      candidate.ColumnName,
		WriteMode:       candidate.WriteMode,
		TimestampColumn: candidate.TimestampColumn,
	})
	if err != nil {
		return "", fmt.Errorf("序列化資料庫輸出候選簽章失敗: %w", err)
	}
	return candidateHash("signature", payload), nil
}

func normalizedCandidateIdentity(identity schema.SourceRuleCandidateIdentity) schema.SourceRuleCandidateIdentity {
	normalized := identity
	normalized.DerivedFromRuleAddress = normalizeAddressKey(identity.DerivedFromRuleAddress)
	if len(identity.TargetBindingScope) == 0 {
		normalized.TargetBindingScope = nil
		return normalized
	}

	scope := append([]schema.SourceRuleCandidateScopeField(nil), identity.TargetBindingScope...)
	sort.Slice(scope, func(i, j int) bool {
		if scope[i].Key == scope[j].Key {
			return scope[i].Value < scope[j].Value
		}
		return scope[i].Key < scope[j].Key
	})
	normalized.TargetBindingScope = scope
	return normalized
}

func candidateHash(prefix string, payload []byte) string {
	sum := sha256.Sum256(payload)
	return prefix + ":" + hex.EncodeToString(sum[:])
}
